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

const iconArrowDown = svg.Icon("ss-arrow-down")

// SsOption represents a selectable item.
type SsOption struct {
	ID          string // unique identifier, returned in OnSelect
	Label       string // visible text
	Description string // optional badge shown on the right
}

type SelectSearch struct {
	Element                                  // value embed — NEVER pointer (TinyGo heap constraint)
	Placeholder string                       // text shown when nothing is selected
	Options     []SsOption                   // initial static options
	OnSelect    func(id, description string) // called when user picks an option
	OnSearch    func(term string) []SsOption // called when ALL local options are filtered out

	// Internal state signals
	selectedLabel *SignalString
	query         *SignalString
	isOpen        *SignalBool
	rows          *SignalNodes
}

func (c *SelectSearch) Init(_ Ctx) {
	c.selectedLabel = NewString("")
	c.query = NewString("")
	c.isOpen = NewBool(false)
	c.rows = NewNodes(c.buildRows("")...)
}

// SetOptions replaces the option list — safe to call after Init/Render,
// e.g. once options from an async source (fetch, MCP call) arrive.
// Preserves the current search query filter, if any.
func (c *SelectSearch) SetOptions(options []SsOption) {
	c.Options = options
	c.rows.Set(c.buildRows(c.query.Get()))
}

func (c *SelectSearch) Render() *Element {
	headerTextSig := DeriveString(func() string {
		sel := c.selectedLabel.Get()
		if sel != "" {
			return sel
		}
		if c.Placeholder != "" {
			return c.Placeholder
		}
		return "Select..."
	})

	toggle := Input("checkbox").Set(ClsSsToggle.AsAttr()).
		ID("ss-toggle").
		BindAttrBool("checked", c.isOpen).
		On("change", func(e Event) {
			c.isOpen.Set(e.TargetChecked())
		})

	// BindText sets textContent which would erase child elements.
	// Wrap header text in a Span so the SVG icon survives as a sibling.
	header := Label().Set(ClsSsHeader.AsAttr()).
		Attr("for", "ss-toggle").
		Child(
			Span().BindText(headerTextSig),
			iconArrowDown.Render(string(ClsSsIcon)),
		)

	searchInput := Input("search").
		Set(ClsSsSearch.AsAttr()).
		ID("ss-search").
		Attr("placeholder", "Search...").
		Bind(c.query).
		Autofocus().
		On("input", func(e Event) {
			term := e.TargetValue()
			// query is already updated by Bind(c.query) in WASM,
			// but we need to trigger the rows update.

			if term != "" {
				allHidden := true
				for _, opt := range c.Options {
					if fmt.Matches(opt.Label, term) || fmt.Matches(opt.Description, term) {
						allHidden = false
						break
					}
				}
				if allHidden && c.OnSearch != nil {
					if newOpts := c.OnSearch(term); len(newOpts) > 0 {
						c.Options = append(c.Options, newOpts...)
					}
				}
			}

			c.rows.Set(c.buildRows(term))
		})

	optList := Ul().Set(ClsSsOptions.AsAttr()).ID("ss-options").
		BindChildren(c.rows)

	dropdown := Show(c.isOpen, func() *Element {
		return Div().Set(ClsSsDropdown.AsAttr()).
			Child(searchInput).
			Child(optList)
	})

	return Div().Set(ClsSsBox.AsAttr()).
		Child(toggle).
		Child(header).
		Child(dropdown)
}

func (c *SelectSearch) buildRows(term string) []*Element {
	var rows []*Element
	for _, opt := range c.Options {
		if term != "" && !fmt.Matches(opt.Label, term) && !fmt.Matches(opt.Description, term) {
			continue
		}

		o := opt // capture loop variable
		item := Li().Set(ClsSsOption.AsAttr()).
			Key(opt.ID).
			ID("ss-opt-"+opt.ID). // required for wirePendingEvents to attach the click handler
			Child(Span().Set(ClsSsLabel.AsAttr()).Text(opt.Label)).
			On("click", func(e Event) {
				c.selectedLabel.Set(o.Label)
				c.isOpen.Set(false)
				c.query.Set("")
				c.rows.Set(c.buildRows(""))
				if c.OnSelect != nil {
					c.OnSelect(o.ID, o.Description)
				}
			})

		if opt.Description != "" {
			item.Child(Span().Set(ClsSsDesc.AsAttr()).Text(opt.Description))
		}
		rows = append(rows, item)
	}
	return rows
}
