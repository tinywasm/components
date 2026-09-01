// Package targetdate is targetlist's sibling for rows that need a prominent
// leading badge — an hour, a day, anything view.Item.LeadTop/Main/Bottom
// carries — instead of a plain label. Same selection/⋮-menu mechanics as
// targetlist, copied rather than shared: the two rows diverge enough in
// markup (the three-line badge has no targetlist equivalent) that factoring
// a common base would cost more indirection than the ~80 duplicated lines
// save. Both satisfy crudview.ListView (see layout/crudview), so a host
// picks whichever row shape fits its data without crudview knowing which one
// it got.
package targetdate

import (
	. "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/svg"
	"github.com/tinywasm/view"
	"github.com/tinywasm/widget"
)

// badgeChars mirrors targetlist's own budget — see that package's comment.
const badgeChars = 16

// NameTargetDate is the widget identity.
const NameTargetDate = widget.Name("targetdate")

const (
	PartRow        = widget.Part("row")
	PartContent    = widget.Part("content")
	PartOptions    = widget.Part("options")
	PartBadge      = widget.Part("badge")
	PartLabel      = widget.Part("label")
	PartList       = widget.Part("list")
	PartLead       = widget.Part("lead")
	PartLeadStack  = widget.Part("lead-stack")
	PartLeadTop    = widget.Part("lead-top")
	PartLeadMain   = widget.Part("lead-main")
	PartLeadBottom = widget.Part("lead-bottom")

	PartButton     = widget.Part("button")
	PartIcon       = widget.Part("icon")
	PartItemDanger = widget.Part("item-danger")
	PartItemIcon   = widget.Part("item-icon")
	PartItemLabel  = widget.Part("item-label")
)

var (
	clsListWrap       = NameTargetDate.Root()
	clsList           = NameTargetDate.Class(PartList)
	clsRow            = NameTargetDate.Class(PartRow)
	clsContent        = NameTargetDate.Class(PartContent)
	clsLabel          = NameTargetDate.Class(PartLabel)
	clsBadge          = NameTargetDate.Class(PartBadge)
	clsLead           = NameTargetDate.Class(PartLead)
	clsLeadStack      = NameTargetDate.Class(PartLeadStack)
	clsLeadTop        = NameTargetDate.Class(PartLeadTop)
	clsLeadMain       = NameTargetDate.Class(PartLeadMain)
	clsLeadBottom     = NameTargetDate.Class(PartLeadBottom)
	clsMenuBtn        = NameTargetDate.Class(PartButton)
	clsMenuIcon       = NameTargetDate.Class(PartIcon)
	clsMenuList       = NameTargetDate.Class(PartOptions)
	clsMenuItemDanger = NameTargetDate.Class(PartItemDanger)
	clsMenuItemIcon   = NameTargetDate.Class(PartItemIcon)
	clsMenuItemLabel  = NameTargetDate.Class(PartItemLabel)
)

const iconDots = svg.Icon("td-dots")
const iconDelete = svg.Icon("td-delete")

// Item is view.Item — see targetlist.Item's comment for why this is an
// alias, not a copy. TargetDate is the one row that actually reads
// LeadTop/Main/Bottom.
type Item = view.Item

// TargetDate is a selectable list of records with a per-row options menu and
// a prominent leading badge (LeadTop/Main/Bottom) instead of a plain label
// lead-in. Same field/method surface as targetlist.TargetList on purpose —
// see crudview.ListView.
type TargetDate struct {
	Element

	Selected *SignalString

	OnSelect func(it Item)
	OnDelete func(id string)

	items    []Item
	rows     *SignalNodes
	openMenu *SignalString
}

func (t *TargetDate) WidgetName() widget.Name { return NameTargetDate }
func (t *TargetDate) WidgetKind() widget.Kind { return widget.Combobox }

func (t *TargetDate) ensure() {
	if t.rows == nil {
		t.rows = NewNodes()
	}
	if t.Selected == nil {
		t.Selected = NewString("")
	}
	if t.openMenu == nil {
		t.openMenu = NewString("")
	}
}

func (t *TargetDate) Init(_ Ctx) { t.ensure() }

// SetItems replaces the visible rows — see targetlist.TargetList.SetItems.
func (t *TargetDate) SetItems(items []Item) {
	t.ensure()
	t.items = items
	nodes := make([]*Element, 0, len(items))
	for _, it := range items {
		nodes = append(nodes, t.buildRow(it))
	}
	t.rows.Set(nodes)
}

func (t *TargetDate) Items() []Item { return t.items }
func (t *TargetDate) Count() int    { return len(t.items) }

func (t *TargetDate) closeAllMenus() {
	t.openMenu.Set("")
}

// CloseMenus — see targetlist.TargetList.CloseMenus for why a host needs this.
func (t *TargetDate) CloseMenus() {
	t.closeAllMenus()
}

func (t *TargetDate) Render() *Element {
	list := Ul().Set(clsList.AsAttr()).Attr("role", "listbox").BindChildren(t.rows)
	return Div().Set(clsListWrap.AsAttr()).Child(list)
}

func (t *TargetDate) buildRow(it Item) *Element {
	id := it.ID
	key := "td-" + id

	isSelSig := DeriveBool(func() bool {
		return t.Selected.Get() == id
	})

	row := Li().Set(clsRow.AsAttr()).
		ID(key).
		Key(key).
		Attr("role", "option").
		BindAttrBool("aria-selected", isSelSig).
		BindState(widget.Selected, isSelSig)

	row.On("click", func(Event) {
		if t.OnSelect != nil {
			t.OnSelect(it)
		}
	})

	trigger := Button().Set(clsMenuBtn.AsAttr()).
		Attr("aria-label", "Opciones").
		Child(iconDots.Render(string(clsMenuIcon)))
	trigger.On("click", func(e Event) {
		e.StopPropagation()
		if t.openMenu.Get() == id {
			t.openMenu.Set("")
		} else {
			t.openMenu.Set(id)
		}
		if t.Selected != nil {
			t.Selected.Set(id)
		}
	})

	del := Button().Set(clsMenuItemDanger.AsAttr()).
		Child(iconDelete.Render(string(clsMenuItemIcon))).
		Child(Span().Set(clsMenuItemLabel.AsAttr()).Text("Eliminar"))
	del.On("click", func(e Event) {
		e.StopPropagation()
		t.closeAllMenus()
		if t.OnDelete != nil {
			t.OnDelete(id)
		}
	})

	options := Div().Set(clsMenuList.AsAttr()).
		BindStateFunc(widget.Open, func() bool { return t.openMenu.Get() == id }).
		Child(del)

	// The leading badge — LeadTop/Main/Bottom, e.g. "Vie" / "10" / "Jul 26".
	// Two nested elements, not one: PartLead is the square (MediaBox
	// AspectSquare + ControlBox — see css.go) and a square's flow centers
	// ONE child, it cannot also stack three lines in a column — Stack() and
	// MediaBox() both claim the same flow slot, so combining them on one
	// Part silently drops whichever was declared first. PartLeadStack is
	// that one child: the actual column, sized to fit inside the square.
	lead := Div().Set(clsLead.AsAttr()).Child(
		Div().Set(clsLeadStack.AsAttr()).Child(
			Span().Set(clsLeadTop.AsAttr()).Text(it.LeadTop),
			Span().Set(clsLeadMain.AsAttr()).Text(it.LeadMain),
			Span().Set(clsLeadBottom.AsAttr()).Text(it.LeadBottom),
		),
	)

	// content holds everything EXCEPT lead: label, the doctor badge, the
	// trigger, the options panel. It carries the row's only padding (see
	// css.go's PartContent) so lead — PartRow's other, padding-less direct
	// child — sits flush against the row's leading/top/bottom edges instead
	// of floating inside a uniform inset on every side.
	//
	// trigger is content's LAST in-flow child, after the label: a growing
	// flex item (clsLabel has Grow() — see css.go) already claims 100% of
	// content's free space during flex resolution, before any margin-auto
	// trick gets a share, so a DOM-trailing trigger is what actually lands
	// at the trailing edge — not a PushEnd on a leading one.
	content := Div().Set(clsContent.AsAttr()).
		Child(Span().Set(clsLabel.AsAttr()).Text(it.Label))
	if it.Description != "" {
		content.Child(Span().Set(clsBadge.AsAttr()).
			Attr("title", it.Description).
			Text(fmt.Convert(it.Description).Truncate(badgeChars).String()))
	}
	content.Child(trigger)
	content.Child(options)

	row.Child(lead)
	row.Child(content)

	return row
}
