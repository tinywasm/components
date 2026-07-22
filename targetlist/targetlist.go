// Package targetlist is the selectable record list used by CRUD views.
//
// It renders rows of Item{Label, Description}, tracks a single selected row, and
// gives each row a ⋮ options menu (Editar / Eliminar) in its top-right corner.
// It owns ONLY its own look (rows, badge, selected state, menu); the layout that
// hosts it decides where it sits. Reactive: call SetItems to (re)populate; the
// selected highlight follows the shared Selected signal.
package targetlist

import (
	. "github.com/tinywasm/css"
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/svg"
)

var (
	clsList     Class = "tl-list"
	clsRow      Class = "tl-row"
	clsRowOn    Class = "tl-row-on"
	clsLabel    Class = "tl-label"
	clsBadge    Class = "tl-badge"
	clsMenu     Class = "tl-menu"
	clsMenuBtn  Class = "tl-menu-btn"
	clsMenuIcon Class = "tl-menu-icon"
	clsMenuList Class = "tl-menu-list"
	clsMenuItem Class = "tl-menu-item"
)

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

	items []Item
	rows  *SignalNodes
}

// ensure lazily creates the reactive state so a host may call SetItems before the
// framework mounts the component (both Init and SetItems are safe in any order).
func (t *TargetList) ensure() {
	if t.rows == nil {
		t.rows = NewNodes()
	}
	if t.Selected == nil {
		t.Selected = NewString("")
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

func (t *TargetList) Render() *Element {
	return Ul().Set(clsList.AsAttr()).BindChildren(t.rows)
}

func (t *TargetList) buildRow(it Item) *Element {
	id := it.ID
	key := "tl-" + id

	row := Li().Set(clsRow.AsAttr()).
		ID(key).
		Key(key).
		BindClass(string(clsRowOn), DeriveBool(func() bool {
			return t.Selected.Get() == id
		}))
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
		if t.OnEdit != nil {
			t.OnEdit(id)
		}
	})
	del := Button().Set(clsMenuItem.AsAttr()).Text("Eliminar")
	del.On("click", func(e Event) {
		e.StopPropagation()
		if t.OnDelete != nil {
			t.OnDelete(id)
		}
	})

	row.Child(Details().Set(clsMenu.AsAttr()).
		Child(summary).
		Child(Div().Set(clsMenuList.AsAttr()).Child(edit, del)))

	return row
}
