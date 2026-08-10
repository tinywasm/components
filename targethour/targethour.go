// Package targethour is targetlist's sibling for rows that need a prominent
// leading badge — an hour, a day, anything view.Item.LeadTop/Main/Bottom
// carries — instead of a plain label. Same selection/⋮-menu mechanics as
// targetlist, copied rather than shared: the two rows diverge enough in
// markup (the three-line badge has no targetlist equivalent) that factoring
// a common base would cost more indirection than the ~80 duplicated lines
// save. Both satisfy crudview.ListView (see layout/crudview), so a host
// picks whichever row shape fits its data without crudview knowing which one
// it got.
package targethour

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

// NameTargetHour is the widget identity.
const NameTargetHour = widget.Name("targethour")

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
	clsListWrap       = NameTargetHour.Root()
	clsList           = NameTargetHour.Class(PartList)
	clsRow            = NameTargetHour.Class(PartRow)
	clsContent        = NameTargetHour.Class(PartContent)
	clsLabel          = NameTargetHour.Class(PartLabel)
	clsBadge          = NameTargetHour.Class(PartBadge)
	clsLead           = NameTargetHour.Class(PartLead)
	clsLeadStack      = NameTargetHour.Class(PartLeadStack)
	clsLeadTop        = NameTargetHour.Class(PartLeadTop)
	clsLeadMain       = NameTargetHour.Class(PartLeadMain)
	clsLeadBottom     = NameTargetHour.Class(PartLeadBottom)
	clsMenuBtn        = NameTargetHour.Class(PartButton)
	clsMenuIcon       = NameTargetHour.Class(PartIcon)
	clsMenuList       = NameTargetHour.Class(PartOptions)
	clsMenuItemDanger = NameTargetHour.Class(PartItemDanger)
	clsMenuItemIcon   = NameTargetHour.Class(PartItemIcon)
	clsMenuItemLabel  = NameTargetHour.Class(PartItemLabel)
)

const iconDots = svg.Icon("th-dots")
const iconDelete = svg.Icon("th-delete")

// Item is view.Item — see targetlist.Item's comment for why this is an
// alias, not a copy. TargetHour is the one row that actually reads
// LeadTop/Main/Bottom.
type Item = view.Item

// TargetHour is a selectable list of records with a per-row options menu and
// a prominent leading badge (LeadTop/Main/Bottom) instead of a plain label
// lead-in. Same field/method surface as targetlist.TargetList on purpose —
// see crudview.ListView.
type TargetHour struct {
	Element

	Selected *SignalString

	OnSelect func(it Item)
	OnDelete func(id string)

	items    []Item
	rows     *SignalNodes
	openMenu *SignalString
}

func (t *TargetHour) WidgetName() widget.Name { return NameTargetHour }
func (t *TargetHour) WidgetKind() widget.Kind { return widget.Combobox }

func (t *TargetHour) ensure() {
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

func (t *TargetHour) Init(_ Ctx) { t.ensure() }

// SetItems replaces the visible rows — see targetlist.TargetList.SetItems.
func (t *TargetHour) SetItems(items []Item) {
	t.ensure()
	t.items = items
	nodes := make([]*Element, 0, len(items))
	for _, it := range items {
		nodes = append(nodes, t.buildRow(it))
	}
	t.rows.Set(nodes)
}

func (t *TargetHour) Items() []Item { return t.items }
func (t *TargetHour) Count() int    { return len(t.items) }

func (t *TargetHour) closeAllMenus() {
	t.openMenu.Set("")
}

// CloseMenus — see targetlist.TargetList.CloseMenus for why a host needs this.
func (t *TargetHour) CloseMenus() {
	t.closeAllMenus()
}

func (t *TargetHour) Render() *Element {
	list := Ul().Set(clsList.AsAttr()).Attr("role", "listbox").BindChildren(t.rows)
	return Div().Set(clsListWrap.AsAttr()).Child(list)
}

func (t *TargetHour) buildRow(it Item) *Element {
	id := it.ID
	key := "th-" + id

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
