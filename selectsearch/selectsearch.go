package selectsearch

import (
	. "github.com/tinywasm/css"
	. "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/svg"
)

var (
	ClsSsBox      Class = "ss-box"
	ClsSsToggle   Class = "ss-toggle"
	ClsSsDropdown Class = "ss-dropdown"
	ClsSsHeader   Class = "ss-header"
	ClsSsIcon     Class = "ss-icon"
	ClsSsSearch   Class = "ss-search"
	ClsSsOptions  Class = "ss-options"
	ClsSsOption   Class = "ss-option"
	ClsSsLabel    Class = "ss-label"
	ClsSsDesc     Class = "ss-desc"
)

// SsOption represents a selectable item.
type SsOption struct {
	ID          string // unique identifier, returned in OnSelect
	Label       string // visible text
	Description string // optional badge shown on the right
}

type SelectSearch struct {
	Element                                    // value embed — NEVER pointer (TinyGo heap constraint)
	Placeholder string                         // text shown when nothing is selected
	Options     []SsOption                     // initial static options
	OnSelect    func(id, description string)   // called when user picks an option
	OnSearch    func(term string) []SsOption   // called when ALL local options are filtered out

	// Internal state
	selectedLabel string
	filterTerm    string
	isOpen        bool
	searchID      string // stable across re-renders for auto-focus in OnMount
}

func (c *SelectSearch) Render() *Element {
	headerText := c.Placeholder
	if c.selectedLabel != "" {
		headerText = c.selectedLabel
	}
	if headerText == "" {
		headerText = "Select..."
	}

	toggle := Input("checkbox").Add(ClsSsToggle.AsAttr())
	if c.isOpen {
		toggle.Attr("checked", "")
	}
	toggle.On("change", func(e Event) {
		c.isOpen = e.TargetChecked()
		c.Update()
	})

	// For(toggle) auto-generates toggle's ID and sets for= — no manual string needed
	header := Label().Add(ClsSsHeader.AsAttr()).
		For(toggle).
		Text(headerText).
		Add(svg.Svg(svg.Use().Attr("href", "#ss-arrow-down")).Add(ClsSsIcon.AsAttr()))

	searchInput := Input("search").
		Add(ClsSsSearch.AsAttr()).
		Attr("placeholder", "Search...").
		Attr("value", c.filterTerm)

	// Stable ID: generated once on first render, reused so OnMount can focus it
	if c.searchID == "" {
		c.searchID = searchInput.GetID()
	} else {
		searchInput.ID(c.searchID)
	}

	searchInput.On("input", func(e Event) {
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

	optList := Div().Add(ClsSsOptions.AsAttr())

	for _, opt := range c.Options {
		if c.filterTerm != "" && !fmt.Matches(opt.Label, c.filterTerm) && !fmt.Matches(opt.Description, c.filterTerm) {
			continue
		}

		o := opt // capture loop variable for the closure
		item := Div().Add(ClsSsOption.AsAttr()).
			Add(Span().Add(ClsSsLabel.AsAttr()).Text(opt.Label))

		item.On("click", func(e Event) {
			c.selectedLabel = o.Label
			c.isOpen = false
			c.filterTerm = ""
			if c.OnSelect != nil {
				c.OnSelect(o.ID, o.Description)
			}
			c.Update()
		})

		if opt.Description != "" {
			item.Add(Span().Add(ClsSsDesc.AsAttr()).Text(opt.Description))
		}
		optList.Add(item)
	}

	dropdown := Div().Add(ClsSsDropdown.AsAttr()).
		Add(searchInput).
		Add(optList)

	return Div().Add(ClsSsBox.AsAttr()).
		Add(toggle).
		Add(header).
		Add(dropdown)
}

// No build tag needed — TinyGo eliminates this as dead code in SSR builds.
func (c *SelectSearch) OnMount() {
	if c.isOpen && c.searchID != "" {
		if el, ok := Get(c.searchID); ok {
			el.Focus()
		}
	}
}
