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
	dom.Element                              // value embed — NEVER pointer (TinyGo heap constraint)
	Placeholder string                       // text shown when nothing is selected
	Options     []Option                     // initial static options
	OnSelect    func(id, description string) // called when user picks an option
	OnSearch    func(term string) []Option   // called when ALL local options are filtered out

	// Internal state
	selectedLabel string
	filterTerm    string
	isOpen        bool
	searchID      string // stable across re-renders for auto-focus in OnMount
}


func (c *SelectSearch) Render() *dom.Element {
	headerText := c.Placeholder
	if c.selectedLabel != "" {
		headerText = c.selectedLabel
	}
	if headerText == "" {
		headerText = "Select..."
	}

	toggle := dom.Input("checkbox").Class("ss-toggle")
	if c.isOpen {
		toggle.Attr("checked", "")
	}
	toggle.On("change", func(e dom.Event) {
		c.isOpen = e.TargetChecked()
		c.Update()
	})

	// For(toggle) auto-generates toggle's ID and sets for= — no manual string needed
	header := dom.Label().
		For(toggle).
		Class("ss-header").
		Text(headerText).
		Add(dom.Svg(dom.Use().Attr("href", "#ss-arrow-down")).Class("ss-icon"))

	searchInput := dom.Input("search").
		Class("ss-search").
		Attr("placeholder", "Search...").
		Attr("value", c.filterTerm)

	// Stable ID: generated once on first render, reused so OnMount can focus it
	if c.searchID == "" {
		c.searchID = searchInput.GetID()
	} else {
		searchInput.ID(c.searchID)
	}

	searchInput.On("input", func(e dom.Event) {
		term := e.TargetValue()
		c.filterTerm = term
		c.isOpen = true

		allHidden := true
		for _, opt := range c.Options {
			if term == "" || fmt.Matches(opt.Label, term) || fmt.Matches(opt.Description, term) {
				allHidden = false
				break
			}
		}

		if allHidden && c.OnSearch != nil {
			newOptions := c.OnSearch(term)
			if len(newOptions) > 0 {
				c.Options = append(c.Options, newOptions...)
			}
		}

		c.Update()
	})

	optList := dom.Div().Class("ss-options")

	for _, opt := range c.Options {
		if c.filterTerm != "" && !fmt.Matches(opt.Label, c.filterTerm) && !fmt.Matches(opt.Description, c.filterTerm) {
			continue
		}

		o := opt // capture loop variable for the closure
		item := dom.Div().
			Class("ss-option").
			Add(dom.Span().Class("ss-label").Text(opt.Label))

		item.On("click", func(e dom.Event) {
			c.selectedLabel = o.Label
			c.isOpen = false
			c.filterTerm = ""
			if c.OnSelect != nil {
				c.OnSelect(o.ID, o.Description)
			}
			c.Update()
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

// No build tag needed — TinyGo eliminates this as dead code in SSR builds.
func (c *SelectSearch) OnMount() {
	if c.isOpen && c.searchID != "" {
		if el, ok := dom.Get(c.searchID); ok {
			el.Focus()
		}
	}
}
