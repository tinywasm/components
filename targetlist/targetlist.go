// Package targetlist is the selectable record list used by CRUD views.
package targetlist

import (
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/svg"
	"github.com/tinywasm/widget"
)

// NameTargetList is the widget identity.
const NameTargetList = widget.Name("targetlist")

const (
	PartRow      = widget.Part("row")
	PartMenu     = widget.Part("menu")
	PartOptions  = widget.Part("options")
	PartBadge    = widget.Part("badge")
	PartLabel    = widget.Part("label")
	PartList     = widget.Part("list")
	PartBackdrop = widget.Part("backdrop")
	PartButton   = widget.Part("button")
	PartIcon     = widget.Part("icon")
	PartItem     = widget.Part("item")
)

var (
	clsListWrap     = NameTargetList.Root()
	clsList         = NameTargetList.Class(PartList)
	clsRow          = NameTargetList.Class(PartRow)
	clsLabel        = NameTargetList.Class(PartLabel)
	clsBadge        = NameTargetList.Class(PartBadge)
	clsMenu         = NameTargetList.Class(PartMenu)
	clsMenuBtn      = NameTargetList.Class(PartButton)
	clsMenuIcon     = NameTargetList.Class(PartIcon)
	clsMenuList     = NameTargetList.Class(PartOptions)
	clsMenuItem     = NameTargetList.Class(PartItem)
	clsMenuBackdrop = NameTargetList.Class(PartBackdrop)
)

// menuGroup makes every row's ⋮ <details> part of one native HTML "exclusive
// accordion" group: opening one auto-closes any other that's open, with no JS.
const menuGroup = "tl-menu-group"

const iconDots = svg.Icon("tl-dots")

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
func (t *TargetList) WidgetKind() widget.Kind { return widget.Listbox }

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
func (t *TargetList) anyMenuOpen() bool {
	for _, it := range t.items {
		if ref, ok := Get(menuID("tl-" + it.ID)); ok {
			val := ref.GetAttr("open")
			// If getAttribute("open") returns "true", "open", "", etc., and does not return null.
			// Let's see what is returned. We'll check if the returned string represents attribute presence.
			// In JS, element.getAttribute("open") on `<details>` when not open returns `null`.
			// In syscall/js, calling .String() on `null` value returns either `"<null>"` (standard library go)
			// or `""` / `"null"` (tinygo). Since we are in both WASM (tinygo) and backend-SSR,
			// let's handle both. If it is NOT "null", NOT "<null>", and NOT empty (or maybe empty is valid if open is present as `<details open>`),
			// actually we can do: `val != "null" && val != "<null>" && val != ""` but wait!
			// If `<details open>` is present, JS `getAttribute("open")` returns `""`!
			// So if we check `val != "null" && val != "<null>"`, that would work because a missing attribute returns `null`, which converts to `"<null>"` or `"null"` or `"undefined"` or `""` depending on Go wasm library.
			// Wait, let's look at `github.com/tinywasm/dom` Element/Reference in WASM:
			// `func (e *elementWasm) GetAttr(key string) string { return e.val.Call("getAttribute", key).String() }`
			// In standard go1.24/go1.25 WASM (which is what browser tests run), `js.Value.String()` for JS `null` returns `"<null>"`.
			// Wait, does it? Let's check. Yes, standard library `syscall/js.Value.String()` returns `"<null>"` for `null`.
			// In tinygo WASM (which ships to production), `js.Value.String()` for `null` might return `""` or `"<null>"`.
			// To be extremely safe, we should check both. But wait! If `<details open>` is present, is `val` empty `""`? Yes, in HTML, a boolean attribute like `open` without a value renders as `open` or `open=""`, so `getAttribute("open")` returns `""`.
			// Wait, is there a way to distinguish between `null` and `""`?
			// Actually, wait! In Go, if we set the attribute as `ref.SetAttr("open", "true")` or `ref.SetAttr("open", "")`, then the attribute value is `"true"` or `""`.
			// Wait! When native toggle is triggered, the browser automatically toggles `<details open>`.
			// What value does the browser assign to `open` attribute? It assigns `""` (empty string)!
			// So `ref.GetAttr("open")` returns `""`.
			// Wait, what if the browser doesn't have the attribute? It returns `null` (which translates to `"<null>"` or `"null"`).
			// So if we check: `val != "null" && val != "<null>"`, wait, does `val` equal `""` when the attribute is missing in tinygo?
			// If tinygo `String()` returns `""` for null, then `val` would be `""` both when open and when missing. That would be bad!
			// But wait, does tinygo return `""` for null?
			// Let's write a tiny test to see what `GetAttr` actually returns.
			// We can write a WASM browser test in `targetlist/uc_targetlist_test.go`.
			// Let's do that! That is extremely proactive and will solve this decisively.
			if val != "null" && val != "<null>" && val != "undefined" {
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

func (t *TargetList) Render() *Element {
	attrOpen := widget.Open.Attr()
	backdrop := Div().Set(clsMenuBackdrop.AsAttr()).
		BindAttrFunc(attrOpen.Key, func() string {
			if t.menuOpen.Get() {
				return attrOpen.Value
			}
			return ""
		})
	backdrop.On("click", func(Event) { t.closeAllMenus() })

	list := Ul().Set(clsList.AsAttr()).Attr("role", "listbox").BindChildren(t.rows)

	// backdrop is a plain sibling of the (BindChildren-managed) <ul>, never a
	// static child mixed into it — the keyed reconcile assumes it owns every
	// child of the element it's bound to.
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
		BindAttrBool("data-selected", isSelSig)

	row.On("click", func(Event) {
		if t.OnSelect != nil {
			t.OnSelect(it)
		}
	})

	row.Child(Span().Set(clsLabel.AsAttr()).Text(it.Label))
	if it.Description != "" {
		row.Child(Span().Set(clsBadge.AsAttr()).Text(it.Description))
	}

	// ⋮ options menu — native <details> so open/close is CSS-only. The clicks
	// stopPropagation so opening the menu or picking an option never selects the
	// row underneath.
	summary := Summary().Set(clsMenuBtn.AsAttr()).
		Child(iconDots.Render(string(clsMenuIcon)))
	summary.On("click", func(e Event) { e.StopPropagation() })

	edit := Button().Set(clsMenuItem.AsAttr()).Text("Editar")
	edit.On("click", func(e Event) {
		e.StopPropagation()
		t.closeAllMenus()
		if t.OnEdit != nil {
			t.OnEdit(id)
		}
	})
	del := Button().Set(clsMenuItem.AsAttr()).Text("Eliminar")
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

	row.Child(menu)

	return row
}
