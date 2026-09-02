// Package targetlist is the selectable record list used by CRUD views.
package targetlist

import (
	. "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/svg"
	"github.com/tinywasm/view"
	"github.com/tinywasm/widget"

	"github.com/tinywasm/components/listselect"
)

// badgeChars is the badge's budget, calibrated against the --chip-width the
// skin gives it. Truncate counts the three-byte ellipsis inside this number.
const badgeChars = 16

// NameTargetList is the widget identity.
const NameTargetList = widget.Name("targetlist")

const (
	PartRow   = widget.Part("row")
	PartCheck = widget.Part("check")
	PartBadge = widget.Part("badge")
	PartLabel = widget.Part("label")
	PartList  = widget.Part("list")
)

var (
	clsListWrap = NameTargetList.Root()
	clsList     = NameTargetList.Class(PartList)
	clsRow      = NameTargetList.Class(PartRow)
	clsCheck    = NameTargetList.Class(PartCheck)
	clsLabel    = NameTargetList.Class(PartLabel)
	clsBadge    = NameTargetList.Class(PartBadge)
)

const iconCheck = svg.Icon("tl-check")

// Item is view.Item, not a copy: a shared shape means crudview.filter's
// []view.Item flows straight into SetItems, and a host swapping this widget
// for targetdate (also view.Item-based) needs no re-mapping either. TargetList
// itself only ever reads ID/Label/Description — LeadTop/Main/Bottom are
// targetdate's slot, ignored here.
type Item = view.Item

// TargetList is a selectable list of records with a multi-selection mode.
type TargetList struct {
	Element

	// Selected holds the id of the highlighted row. Optional — created if nil so a
	// host can share it (e.g. a CRUD view binding the form to the same signal).
	Selected *SignalString

	// Row callbacks. Optional.
	OnSelect func(it Item) // row body clicked

	items []Item
	rows  *SignalNodes
	sel   listselect.Mode
}

func (t *TargetList) WidgetName() widget.Name { return NameTargetList }
func (t *TargetList) WidgetKind() widget.Kind { return widget.Combobox }

func (t *TargetList) ensure() {
	if t.rows == nil {
		t.rows = NewNodes()
	}
	if t.Selected == nil {
		t.Selected = NewString("")
	}
}

func (t *TargetList) Init(_ Ctx) { t.ensure() }

func (t *TargetList) SetSelectMode(on bool)       { t.sel.SetOn(on) }
func (t *TargetList) OnCheckedChange(fn func(int)) { t.sel.OnChange = fn }

func (t *TargetList) CheckedIDs() []string {
	ids := make([]string, len(t.items))
	for i, it := range t.items {
		ids[i] = it.ID
	}
	return t.sel.CheckedIDs(ids)
}

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
	list := Ul().Set(clsList.AsAttr()).Attr("role", "listbox").BindChildren(t.rows)

	return Div().Set(clsListWrap.AsAttr()).
		BindStateFunc(widget.Open, func() bool { return t.sel.On().Get() }).
		Child(list)
}

func (t *TargetList) buildRow(it Item) *Element {
	id := it.ID
	key := "tl-" + id

	isSelSig := DeriveBool(func() bool {
		if t.sel.On().Get() {
			return t.sel.IsChecked(id)
		}
		return t.Selected.Get() == id
	})

	row := Li().Set(clsRow.AsAttr()).
		ID(key).
		Key(key).
		Attr("role", "option").
		BindAttrBool("aria-selected", isSelSig).
		BindState(widget.Selected, isSelSig)

	row.On("click", func(Event) {
		if t.sel.On().Get() {
			t.sel.Toggle(id)
			return
		}
		if t.OnSelect != nil {
			t.OnSelect(it)
		}
	})

	check := Span().Set(clsCheck.AsAttr()).Child(iconCheck.Render(string(clsCheck)))

	row.Child(check)
	row.Child(Span().Set(clsLabel.AsAttr()).Text(it.Label))
	if it.Description != "" {
		row.Child(Span().Set(clsBadge.AsAttr()).
			Attr("title", it.Description).
			Text(fmt.Convert(it.Description).Truncate(badgeChars).String()))
	}

	return row
}
