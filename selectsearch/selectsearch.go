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
		BindAttrBool("checked", c.isOpen).
		On("change", func(e Event) {
			c.isOpen.Set(e.TargetChecked())
		})

	// For(toggle) auto-generates toggle's ID and sets for= — no manual string needed
	header := Label().Set(ClsSsHeader.AsAttr()).
		For(toggle).
		BindText(headerTextSig).
		Child(svg.Svg().Child(svg.Use().Attr("href", "#ss-arrow-down")).Set(ClsSsIcon.AsAttr()))

	searchInput := Input("search").
		Set(ClsSsSearch.AsAttr()).
		Attr("placeholder", "Search...").
		Bind(c.query).
		Autofocus()

	// Rebuild rows when query changes
	DeriveString(func() string {
		term := c.query.Get()
		c.isOpen.Set(true)

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

		c.rows.Set(c.buildRows(term))
		return term
	})

	optList := Div().Set(ClsSsOptions.AsAttr()).
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
		item := Div().Set(ClsSsOption.AsAttr()).
			Key(opt.ID).
			Child(Span().Set(ClsSsLabel.AsAttr()).Text(opt.Label)).
			On("click", func(e Event) {
				c.selectedLabel.Set(o.Label)
				c.isOpen.Set(false)
				c.query.Set("")
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
