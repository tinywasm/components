// Package targetlist is the selectable record list used by CRUD views.
package targetlist

import (
	. "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/svg"
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
	PartRow        = widget.Part("row")
	PartMenu       = widget.Part("menu")
	PartOptions    = widget.Part("options")
	PartBadge      = widget.Part("badge")
	PartLabel      = widget.Part("label")
	PartList       = widget.Part("list")
	PartBackdrop   = widget.Part("backdrop")
	PartButton     = widget.Part("button")
	PartIcon       = widget.Part("icon")
	PartItem       = widget.Part("item")
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
	clsMenu           = NameTargetList.Class(PartMenu)
	clsMenuBtn        = NameTargetList.Class(PartButton)
	clsMenuIcon       = NameTargetList.Class(PartIcon)
	clsMenuList       = NameTargetList.Class(PartOptions)
	clsMenuItem       = NameTargetList.Class(PartItem)
	clsMenuItemDanger = NameTargetList.Class(PartItemDanger)
	clsMenuItemIcon   = NameTargetList.Class(PartItemIcon)
	clsMenuItemLabel  = NameTargetList.Class(PartItemLabel)
	clsMenuBackdrop   = NameTargetList.Class(PartBackdrop)
)

// menuGroup makes every row's ⋮ <details> part of one native HTML "exclusive
// accordion" group: opening one auto-closes any other that's open, with no JS.
const menuGroup = "tl-menu-group"

const iconDots = svg.Icon("tl-dots")

// iconEdit and iconDelete back the ⋮ menu's two options. Rendered in the
// markup on every device — desktop keeps them alongside the text label,
// hidden by CSS; mobile is where they become the visible affordance and the
// label hides instead, so the item reads as an icon button matching
// crudview's own floating action button rather than a dropdown row.
const (
	iconEdit   = svg.Icon("tl-edit")
	iconDelete = svg.Icon("tl-delete")
)

// Item is one selectable record: an id, a visible label, and an optional badge.
type Item struct {
	ID          string
	Label       string
	Description string
}

// TargetList is a selectable list of records with a per-row options menu.
type TargetList struct {
	Element

	// Selected holds the id of the highlighted row. Optional — created if nil so a
	// host can share it (e.g. a CRUD view binding the form to the same signal).
	Selected *SignalString

	// Row callbacks. All optional.
	OnSelect func(it Item)   // row body clicked
	OnEdit   func(id string) // ⋮ → Editar
	OnDelete func(id string) // ⋮ → Eliminar

	items    []Item
	rows     *SignalNodes
	menuOpen *SignalBool
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
	if t.menuOpen == nil {
		t.menuOpen = NewBool(false)
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

// menuID is the DOM id of a row's ⋮ <details> — used to close it
// programmatically (closeAllMenus) since native <details> only closes on a
// summary click or an explicit attribute removal, never on an outside click.
func menuID(key string) string { return key + ".menu" }

// anyMenuOpen informa si alguna fila tiene su ⋮ <details> abierto. El navegador
// posee ese estado; esto lo refleja en Go para que la hoja pueda seleccionarlo.
//
// GetAttr wraps js.Value.String(): a present attribute (even "<details open>",
// whose value is "") returns its string value; a missing one returns the
// stringified JS null, "<null>" — that's the only case treated as closed.
func (t *TargetList) anyMenuOpen() bool {
	for _, it := range t.items {
		if ref, ok := Get(menuID("tl-" + it.ID)); ok {
			if ref.GetAttr("open") != "<null>" {
				return true
			}
		}
	}
	return false
}

// closeAllMenus force-closes every row's ⋮ menu. Wired to: picking Editar/
// Eliminar (native <details> does not close itself on that), and a full-page
// backdrop that is shown while any menu is open, so a click anywhere outside
// a menu closes it too.
func (t *TargetList) closeAllMenus() {
	for _, it := range t.items {
		if ref, ok := Get(menuID("tl-" + it.ID)); ok {
			ref.RemoveAttr("open")
		}
	}
	t.menuOpen.Set(false)
}

// CloseMenus is closeAllMenus, exported for a host that shares Selected (see
// crudview) to call when ITS OWN cancel path clears the selection. The ⋮
// summary's click handler sets Selected directly (see buildRow) so the
// floating options menu reads as belonging to a highlighted row, but native
// <details> open state is not tracked by that signal — clearing Selected
// elsewhere does not, by itself, close a menu left open on some row. A host
// that lets a row's ⋮ drive its own "active" state is responsible for
// closing it back up when that state resets.
func (t *TargetList) CloseMenus() {
	t.closeAllMenus()
}

func (t *TargetList) Render() *Element {
	backdrop := Div().Set(clsMenuBackdrop.AsAttr()).
		BindStateFunc(widget.Open, func() bool { return t.menuOpen.Get() })
	backdrop.On("click", func(Event) { t.closeAllMenus() })

	list := Ul().Set(clsList.AsAttr()).Attr("role", "listbox").BindChildren(t.rows)

	// backdrop is a plain sibling of the (BindChildren-managed) <ul>, never a
	// static child mixed into it — the keyed reconcile assumes it owns every
	// child of the element it's bound to.
	//
	// Order is part of the stacking contract: the backdrop and a row's open
	// options both sit on the widget overlay level (var(--z-dropdown)), and
	// among equals the LATER sibling paints on top — the backdrop is first so
	// the options never end up underneath it. Do not reorder.
	return Div().Set(clsListWrap.AsAttr()).Child(backdrop, list)
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

	row.Child(Span().Set(clsLabel.AsAttr()).Text(it.Label))

	// ⋮ options menu — native <details> so open/close is CSS-only. The clicks
	// stopPropagation so the row's own OnSelect (full select-and-navigate,
	// see crudview's selectAction) never fires from a ⋮ tap — opening the
	// menu must not also jump the mobile strip to the form panel.
	//
	// It still sets Selected directly, though: with the mobile row menu now
	// docked to the viewport corner (crudview/rightpanel PLAN.md) rather than
	// anchored to the row, there is nothing else on screen linking the
	// floating Editar/Eliminar icons back to the row they act on. Highlighting
	// the row is that link. This also flips crudview's own active()/Open
	// state for free — Selected is the same *SignalString crudview binds its
	// action button's icon to — so the button reads "cancel" the instant the
	// menu opens, with no separate wiring needed here.
	summary := Summary().Set(clsMenuBtn.AsAttr()).
		Child(iconDots.Render(string(clsMenuIcon)))
	summary.On("click", func(e Event) {
		e.StopPropagation()
		if t.Selected != nil {
			t.Selected.Set(id)
		}
	})

	edit := Button().Set(clsMenuItem.AsAttr()).
		Child(iconEdit.Render(string(clsMenuItemIcon))).
		Child(Span().Set(clsMenuItemLabel.AsAttr()).Text("Editar"))
	edit.On("click", func(e Event) {
		e.StopPropagation()
		t.closeAllMenus()
		if t.OnEdit != nil {
			t.OnEdit(id)
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

	menu := Details().Set(clsMenu.AsAttr()).
		ID(menuID(key)).
		Attr("name", menuGroup).
		Child(summary).
		Child(Div().Set(clsMenuList.AsAttr()).Child(edit, del))

	menu.On("toggle", func(e Event) {
		t.menuOpen.Set(t.anyMenuOpen())
	})

	// The menu sits on the label's line and the badge after it, so a row that
	// runs out of width wraps the badge rather than the affordance.
	row.Child(menu)
	if it.Description != "" {
		row.Child(Span().Set(clsBadge.AsAttr()).
			Attr("title", it.Description). // the untruncated text stays reachable
			Text(fmt.Convert(it.Description).Truncate(badgeChars).String()))
	}

	return row
}
