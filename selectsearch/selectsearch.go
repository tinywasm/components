package selectsearch

import (
	"github.com/tinywasm/dom"
	"strings"
)

// Option represents a selectable item.
type Option struct {
	ID          string // unique identifier, returned in OnSelect
	Label       string // visible text
	Description string // optional badge shown on the right
}

type SelectSearch struct {
	dom.Element          // value embed — NEVER pointer (TinyGo heap constraint)
	Placeholder string   // text shown when nothing is selected
	Options     []Option // initial static options
	OnSelect    func(id, description string) // called when user picks an option
	OnSearch    func(term string) []Option   // called when ALL local options are filtered out

	// Internal state
	SelectedLabel string
	FilterTerm    string
}

func (c *SelectSearch) Render() *dom.Element {
	id := c.GetID()

	headerText := c.Placeholder
	if c.SelectedLabel != "" {
		headerText = c.SelectedLabel
	}
	if headerText == "" {
		headerText = "Select..."
	}

	// Hidden checkbox — drives the CSS toggle
	toggle := &dom.Element{}
	toggle.Attr("type", "checkbox").
		ID(id+"-toggle").
		Class("ss-toggle")

	// Header label — clicking it toggles the checkbox
	header := dom.Label(id+"-toggle").
		Class("ss-header").
		Text(headerText).
		Add(dom.Svg(dom.Use().Attr("href", "#ss-arrow-down")).Class("ss-icon"))

	// Search input (raw dom.Element — no form/input dependency)
	searchInput := &dom.Element{}
	searchInput.Attr("type", "search").
		ID(id+"-search").
		Class("ss-search").
		Attr("placeholder", "Search...").
		Attr("value", c.FilterTerm)

	// Options list
	optList := dom.Div().Class("ss-options").ID(id + "-options")

	filterTerm := strings.ToLower(c.FilterTerm)
	for _, opt := range c.Options {
		match := filterTerm == "" ||
			strings.Contains(strings.ToLower(opt.Label), filterTerm) ||
			strings.Contains(strings.ToLower(opt.Description), filterTerm)

		if !match {
			continue
		}

		item := dom.Div().
			Class("ss-option").
			Attr("data-id", opt.ID).
			Attr("data-description", opt.Description).
			Add(dom.Span().Class("ss-label").Text(opt.Label)).
			On("click", func(e dom.Event) {
				c.SelectedLabel = opt.Label
				if c.OnSelect != nil {
					c.OnSelect(opt.ID, opt.Description)
				}
				c.FilterTerm = ""
				dom.Update(c)
			})

		if opt.Description != "" {
			item.Add(dom.Span().Class("ss-desc").Text(opt.Description))
		}
		optList.Add(item)
	}

	dropdown := dom.Div().Class("ss-dropdown").
		Add(searchInput).
		Add(optList)

	return dom.Div().
		Class("ss-box").
		Add(toggle).
		Add(header).
		Add(dropdown)
}
