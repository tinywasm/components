// Package targetlist is the selectable record list used by CRUD views.
package targetlist

import (
	. "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/svg"
	"github.com/tinywasm/view"
	"github.com/tinywasm/widget"
)

// badgeChars is the badge's budget, calibrated against the --chip-width the
// skin gives it. Truncate counts the three-byte ellipsis inside this number.
//
// The count is BYTES, not runes, so an accented character costs two and a cut
// can land mid-character. Badges carry identifiers and addresses in practice,
// which is why this is acceptable here.
const badgeChars = 16

// NameTargetList is the widget identity.
const NameTargetList = widget.Name("targetlist")

const (
	PartRow     = widget.Part("row")
	PartOptions = widget.Part("options")
	PartBadge   = widget.Part("badge")
	PartLabel   = widget.Part("label")
	PartList    = widget.Part("list")
	// PartMenu (the <details> wrapper) and PartBackdrop are gone with the
	// overlay: the trigger is a plain button and the options are an in-flow
	// accordion, so there is no wrapper to style and nothing floating for a
	// backdrop to dismiss.
	PartButton     = widget.Part("button")
	PartIcon       = widget.Part("icon")
	PartItemDanger = widget.Part("item-danger")
	PartItemIcon   = widget.Part("item-icon")
	PartItemLabel  = widget.Part("item-label")
)

var (
	clsListWrap       = NameTargetList.Root()
	clsList           = NameTargetList.Class(PartList)
	clsRow            = NameTargetList.Class(PartRow)
	clsLabel          = NameTargetList.Class(PartLabel)
	clsBadge          = NameTargetList.Class(PartBadge)
	clsMenuBtn        = NameTargetList.Class(PartButton)
	clsMenuIcon       = NameTargetList.Class(PartIcon)
	clsMenuList       = NameTargetList.Class(PartOptions)
	clsMenuItemDanger = NameTargetList.Class(PartItemDanger)
	clsMenuItemIcon   = NameTargetList.Class(PartItemIcon)
	clsMenuItemLabel  = NameTargetList.Class(PartItemLabel)
)

// openMenu holds the id of the row whose ⋮ options are expanded, or "" when
// none is. One signal instead of one per row is what makes the accordion
// exclusive by construction: only one id can be in it, so opening a row closes
// whichever was open with no bookkeeping and no DOM walk.
//
// This replaced a native <details name="…"> group. The native element gave the
// exclusivity for free but owned the open state in the DOM, which meant Go had
// to read it back out (GetAttr("open") != "<null>") and force it closed by
// removing the attribute — and the options had to live INSIDE the <details> to
// be hidden by it, which is exactly the nesting that made them an overlay
// trapped in the list's Scroll() region. See css.go's PartOptions comment.

const iconDots = svg.Icon("tl-dots")

// iconDelete backs the ⋮ menu's only option. Rendered in the markup on every
// device — desktop keeps it alongside the text label, hidden by CSS; mobile
// is where it becomes the visible affordance and the label hides instead, so
// the item reads as an icon button matching crudview's own floating action
// button rather than a dropdown row. (Editar's pencil left the menu with the
// lock it existed to undo — see crudview's Stage-3 comment in targetlist.go.)
const iconDelete = svg.Icon("tl-delete")

// Item is view.Item, not a copy: a shared shape means crudview.filter's
// []view.Item flows straight into SetItems, and a host swapping this widget
// for targetdate (also view.Item-based) needs no re-mapping either. TargetList
// itself only ever reads ID/Label/Description — LeadTop/Main/Bottom are
// targetdate's slot, ignored here.
type Item = view.Item

// TargetList is a selectable list of records with a per-row options menu.
type TargetList struct {
	Element

	// Selected holds the id of the highlighted row. Optional — created if nil so a
	// host can share it (e.g. a CRUD view binding the form to the same signal).
	Selected *SignalString

	// Row callbacks. All optional.
	OnSelect func(it Item)   // row body clicked
	OnDelete func(id string) // ⋮ → Eliminar

	items    []Item
	rows     *SignalNodes
	openMenu *SignalString
}

func (t *TargetList) WidgetName() widget.Name { return NameTargetList }
func (t *TargetList) WidgetKind() widget.Kind { return widget.Combobox }

// ensure lazily creates the reactive state so a host may call SetItems before the
// framework mounts the component (both Init and SetItems are safe in any order).
func (t *TargetList) ensure() {
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

func (t *TargetList) Init(_ Ctx) { t.ensure() }

// SetItems replaces the visible rows. Safe to call from a host on every filter or
// reload; rows go through the keyed reconcile so their bindings stay wired.
func (t *TargetList) SetItems(items []Item) {
	t.ensure()
	t.items = items
	nodes := make([]*Element, 0, len(items))
	for _, it := range items {
		nodes = append(nodes, t.buildRow(it))
	}
	t.rows.Set(nodes)
}

// Items returns the current items (the data behind the rendered rows).
func (t *TargetList) Items() []Item { return t.items }

// Count reports how many rows are currently rendered (used by hosts/tests).
func (t *TargetList) Count() int { return len(t.items) }

// closeAllMenus collapses whichever row is expanded. One assignment: the open
// state is a single id in Go now, not an `open` attribute spread across the
// DOM, so there is nothing to walk and nothing to remove.
func (t *TargetList) closeAllMenus() {
	t.openMenu.Set("")
}

// CloseMenus is closeAllMenus, exported for a host that shares Selected (see
// crudview) to call when ITS OWN cancel path clears the selection. The ⋮
// trigger sets Selected directly (see buildRow) so the expanded options read
// as belonging to a highlighted row, but the two are separate signals —
// clearing Selected elsewhere does not, by itself, collapse a row left open. A
// host that lets a row's ⋮ drive its own "active" state is responsible for
// closing it back up when that state resets.
func (t *TargetList) CloseMenus() {
	t.closeAllMenus()
}

func (t *TargetList) Render() *Element {
	// No backdrop. It existed to catch the outside-tap that dismissed a
	// floating options panel; the options are in the row now, so there is
	// nothing floating to dismiss — and a viewport-sized element over an
	// in-flow accordion would eat every tap meant for the controls inside it,
	// which is the same trap usermenu already documents for its own mobile
	// accordion.
	list := Ul().Set(clsList.AsAttr()).Attr("role", "listbox").BindChildren(t.rows)

	return Div().Set(clsListWrap.AsAttr()).Child(list)
}

func (t *TargetList) buildRow(it Item) *Element {
	id := it.ID
	key := "tl-" + id

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

	// ⋮ trigger. A plain button, not a <summary>: the options it expands are a
	// sibling further down this row, not its children, so there is no <details>
	// left for a summary to belong to. Toggling one id in openMenu keeps the
	// accordion exclusive without touching any other row.
	//
	// StopPropagation so the row's own OnSelect (full select-and-navigate, see
	// crudview's selectAction) never fires from a ⋮ tap — expanding the options
	// must not also jump the mobile strip to the form panel.
	//
	// It still sets Selected directly: the amber highlight is what ties the
	// expanded Eliminar to the record it acts on. This also flips
	// crudview's own active()/Open state for free — Selected is the same
	// *SignalString crudview binds its action button's icon to — so the button
	// reads "cancel" the instant the menu opens, with no separate wiring here.
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

	// The options are a sibling of the trigger, not its child, and the LAST
	// child of the row: Width(Full)+KeepSize makes them take a wrapped line of
	// their own beneath the label, which only works if nothing follows them on
	// that line. RevealedBy(Open) in css.go hides them until this row's id is
	// the one in openMenu.
	options := Div().Set(clsMenuList.AsAttr()).
		BindStateFunc(widget.Open, func() bool { return t.openMenu.Get() == id }).
		Child(del)

	// The trigger is the row's LAST in-flow child (before options), not the
	// first: it used to lead so the mobile master-detail sliver (which shows
	// only the row's leading edge) could still reach it — see the PartButton
	// comment in css.go for the measured history. That placement is gone
	// because it also meant the ⋮ could never sit at the trailing edge on a
	// wide screen: PartLabel's Grow() already claims 100% of the row's free
	// space during flex resolution, so a margin-auto push on a DOM-leading
	// trigger has nothing left to distribute — a DOM-trailing trigger is the
	// only placement that actually lands at the trailing edge. The sliver
	// trade-off is conscious, not overlooked: a phone user loses one-tap ⋮
	// access from the sliver and returns to the full list to delete, same
	// as reaching any other off-sliver control.
	row.Child(Span().Set(clsLabel.AsAttr()).Text(it.Label))
	if it.Description != "" {
		row.Child(Span().Set(clsBadge.AsAttr()).
			Attr("title", it.Description). // the untruncated text stays reachable
			Text(fmt.Convert(it.Description).Truncate(badgeChars).String()))
	}
	row.Child(trigger)
	row.Child(options)

	return row
}
