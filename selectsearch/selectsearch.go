package selectsearch

import (
	. "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/svg"
	"github.com/tinywasm/widget"
)

// NameSelectSearch is the widget name.
const NameSelectSearch = widget.Name("selectsearch")

const (
	PartToggle     = widget.Part("toggle")
	PartDropdown   = widget.Part("dropdown")
	PartHeader     = widget.Part("header")
	PartHeaderText = widget.Part("header-text")
	PartIcon       = widget.Part("icon")
	PartGlyph      = widget.Part("glyph")
	PartSearch     = widget.Part("search")
	PartOptions    = widget.Part("options")
	PartOption     = widget.Part("option")
	PartText       = widget.Part("text")
	PartLabel      = widget.Part("label")
	PartSublabel   = widget.Part("sublabel")
	PartDesc       = widget.Part("desc")
)

var (
	ClsSsBox        = NameSelectSearch.Root()
	ClsSsToggle     = NameSelectSearch.Class(PartToggle)
	ClsSsDropdown   = NameSelectSearch.Class(PartDropdown)
	ClsSsHeader     = NameSelectSearch.Class(PartHeader)
	ClsSsHeaderText = NameSelectSearch.Class(PartHeaderText)
	ClsSsIcon       = NameSelectSearch.Class(PartIcon)
	ClsSsGlyph      = NameSelectSearch.Class(PartGlyph)
	ClsSsSearch     = NameSelectSearch.Class(PartSearch)
	ClsSsOptions    = NameSelectSearch.Class(PartOptions)
	ClsSsOption     = NameSelectSearch.Class(PartOption)
	ClsSsText       = NameSelectSearch.Class(PartText)
	ClsSsLabel      = NameSelectSearch.Class(PartLabel)
	ClsSsSublabel   = NameSelectSearch.Class(PartSublabel)
	ClsSsDesc       = NameSelectSearch.Class(PartDesc)
)

const iconArrowDown = svg.Icon("ss-arrow-down")

// SsOption represents a selectable item.
type SsOption struct {
	ID          string // unique identifier, returned in OnSelect
	Label       string // visible text
	Sublabel    string // optional second line under Label — position only, no assumed content
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

	onFilter func(term string) // set via OnFilterChange — satisfies widget.Filterable
}

func (c *SelectSearch) WidgetName() widget.Name { return NameSelectSearch }
func (c *SelectSearch) WidgetKind() widget.Kind { return widget.Combobox }

var _ widget.Filterable = (*SelectSearch)(nil)

// OnFilterChange implements widget.Filterable: it registers the sink called
// with the picked option's ID whenever a selection is made. This is a
// SEPARATE, additive wiring path from OnSelect — OnSelect still gets
// (id, description) for a consumer that needs both; OnFilterChange exists so
// a host that only knows the generic Filterable contract (e.g.
// tinywasm/layout/crudview's Filter slot) can drop a *SelectSearch into the
// same seam a *searchbar.SearchBar fills today, with no bespoke glue.
//
// The signature is fixed by widget.Filterable — do not add a parameter, do
// not return anything, do not add a companion getter (see searchbar.go's
// OnFilterChange for the same rule stated for SearchBar).
func (c *SelectSearch) OnFilterChange(fn func(term string)) { c.onFilter = fn }

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
			checked := e.TargetChecked()
			c.isOpen.Set(checked)
			if checked {
				if ref, ok := Get("ss-search"); ok {
					ref.Focus()
				}
			}
		})

	// The icon is a filled square cap around a white glyph, FIRST child —
	// searchbar's own PartIcon/PartGlyph layout (that package's Render puts
	// its cap before its input too; see its css.go for the flush-square
	// recipe PartHeader below mirrors), not a bare svg trailing the text:
	// a bare <svg> painted straight onto the header read as an unstyled
	// stray mark, disconnected from the rest of the chassis, and trailing
	// it put the cap on the wrong edge for this chassis' own convention.
	//
	// BindText sets textContent which would erase child elements, so the
	// text gets its own Span — and that Span carries the header's padding
	// (see PartHeaderText in css.go): the icon must NOT be padded away
	// from the header's edges, or it stops reading as a flush cap.
	icon := Div().Set(ClsSsIcon.AsAttr()).Child(iconArrowDown.Render(string(ClsSsGlyph)))
	header := Label().Set(ClsSsHeader.AsAttr()).
		Attr("for", "ss-toggle").
		Child(
			icon,
			Span().Set(ClsSsHeaderText.AsAttr()).BindText(headerTextSig),
		)

	searchInput := Input("search").
		Set(ClsSsSearch.AsAttr()).
		ID("ss-search").
		Attr("placeholder", "Search...").
		Attr("role", "combobox").
		BindAttrBool("aria-expanded", c.isOpen).
		Attr("aria-controls", "ss-options").
		Bind(c.query).
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
		Attr("role", "listbox").
		BindChildren(c.rows)

	dropdown := Div().Set(ClsSsDropdown.AsAttr()).
		Child(searchInput).
		Child(optList)

	return Div().Set(ClsSsBox.AsAttr()).
		Child(toggle).
		Child(header).
		Child(Show(c.isOpen, dropdown))
}

// selectOption is the single place an option becomes "chosen" — today only
// a mouse click reaches it, but every future input path (keyboard, a future
// OnSearch auto-pick) commits through here too, so OnSelect and the
// Filterable sink can never fire out of step with each other.
func (c *SelectSearch) selectOption(o SsOption) {
	c.selectedLabel.Set(o.Label)
	c.isOpen.Set(false)
	c.query.Set("")
	c.rows.Set(c.buildRows(""))
	if c.OnSelect != nil {
		c.OnSelect(o.ID, o.Description)
	}
	if c.onFilter != nil {
		c.onFilter(o.ID)
	}
}

func (c *SelectSearch) buildRows(term string) []*Element {
	var rows []*Element
	for _, opt := range c.Options {
		if term != "" && !fmt.Matches(opt.Label, term) && !fmt.Matches(opt.Description, term) {
			continue
		}

		o := opt // capture loop variable

		// text stacks Label over the optional Sublabel — a second line under
		// the name, not a sibling beside it, which is why it is its own Grow()
		// column instead of two more spans loose in the row's own Row().
		text := Div().Set(ClsSsText.AsAttr()).
			Child(Span().Set(ClsSsLabel.AsAttr()).Text(opt.Label))
		if opt.Sublabel != "" {
			text.Child(Span().Set(ClsSsSublabel.AsAttr()).Text(opt.Sublabel))
		}

		item := Li().Set(ClsSsOption.AsAttr()).
			Key(opt.ID).
			ID("ss-opt-"+opt.ID). // required for wirePendingEvents to attach the click handler
			Attr("role", "option").
			Child(text).
			On("click", func(e Event) { c.selectOption(o) })

		if opt.Description != "" {
			item.Child(Span().Set(ClsSsDesc.AsAttr()).Text(opt.Description))
		}
		rows = append(rows, item)
	}
	return rows
}
