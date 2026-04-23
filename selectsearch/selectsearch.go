package selectsearch

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
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
	selectedLabel string
	filterTerm    string
}

func (c *SelectSearch) matches(opt Option, term string) bool {
	if term == "" {
		return true
	}
	return fmt.Contains(fmt.Convert(opt.Label).ToLower().String(), term) ||
		fmt.Contains(fmt.Convert(opt.Description).ToLower().String(), term)
}

func (c *SelectSearch) Render() *dom.Element {
	id := c.GetID()

	headerText := c.Placeholder
	if c.selectedLabel != "" {
		headerText = c.selectedLabel
	}
	if headerText == "" {
		headerText = "Select..."
	}

	// Hidden checkbox — drives the CSS toggle
	toggle := dom.Input("checkbox").
		ID(id+"-toggle").
		Class("ss-toggle")

	// Header label — clicking it toggles the checkbox
	header := dom.Label().
		Attr("for", id+"-toggle").
		Class("ss-header").
		Text(headerText).
		Add(dom.Svg(dom.Use().Attr("href", "#ss-arrow-down")).Class("ss-icon"))

	// Search input
	searchInput := dom.Input("search").
		ID(id+"-search").
		Class("ss-search").
		Attr("placeholder", "Search...").
		Attr("value", c.filterTerm)

	// Options list
	optList := dom.Div().Class("ss-options").ID(id + "-options")

	filterTerm := fmt.Convert(c.filterTerm).ToLower().String()
	for _, opt := range c.Options {
		if !c.matches(opt, filterTerm) {
			continue
		}

		item := dom.Div().
			ID(id + "-opt-" + opt.ID).
			Class("ss-option").
			Attr("data-id", opt.ID).
			Attr("data-description", opt.Description).
			Add(dom.Span().Class("ss-label").Text(opt.Label))

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
