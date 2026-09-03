// Package targethour is targetlist's sibling for a day's booked slots: each
// row leads with a prominent hour (HH:MM) and may carry a status tint
// (pending / confirmed / attended). Same multi-selection mechanics as
// targetlist/targetdate — it assembles components/listselect, it does not
// re-declare it.
package targethour

import (
	. "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/view"
	"github.com/tinywasm/widget"

	"github.com/tinywasm/components/listselect"
	"github.com/tinywasm/icons/pencil"
	"github.com/tinywasm/icons/trash"
)

const badgeChars = 16

const NameTargetHour = widget.Name("targethour")

const (
	PartRow            = widget.Part("row")
	PartContent        = widget.Part("content")
	PartCheck          = widget.Part("check")
	PartCheckTrash     = widget.Part("check-trash")
	PartCheckPencil    = widget.Part("check-pencil")
	PartCheckAll       = widget.Part("check-all")
	PartCheckAllTrash  = widget.Part("check-all-trash")
	PartCheckAllPencil = widget.Part("check-all-pencil")
	PartCheckAllCount  = widget.Part("check-all-count")
	PartHour           = widget.Part("hour")
	PartLabel          = widget.Part("label")
	PartBadge          = widget.Part("badge")
	PartList           = widget.Part("list")
)

var (
	clsListWrap       = NameTargetHour.Root()
	clsList           = NameTargetHour.Class(PartList)
	clsRow            = NameTargetHour.Class(PartRow)
	clsContent        = NameTargetHour.Class(PartContent)
	clsCheck          = NameTargetHour.Class(PartCheck)
	clsCheckTrash     = NameTargetHour.Class(PartCheckTrash)
	clsCheckPencil    = NameTargetHour.Class(PartCheckPencil)
	clsCheckAll       = NameTargetHour.Class(PartCheckAll)
	clsCheckAllTrash  = NameTargetHour.Class(PartCheckAllTrash)
	clsCheckAllPencil = NameTargetHour.Class(PartCheckAllPencil)
	clsCheckAllCount  = NameTargetHour.Class(PartCheckAllCount)
	clsHour           = NameTargetHour.Class(PartHour)
	clsLabel          = NameTargetHour.Class(PartLabel)
	clsBadge          = NameTargetHour.Class(PartBadge)
)

type Item = view.Item

// Status is a row's booking state — drives the tint only. The zero value
// (StatusPending) paints no tint, exactly like a plain targetlist row.
type Status uint8

const (
	StatusPending Status = iota // no tint
	StatusConfirmed             // confirmed by reception
	StatusAttended              // patient already attended
)

type TargetHour struct {
	Element

	Selected *SignalString
	OnSelect func(it Item)

	// StatusOf maps a row to its booking state for the tint. Optional — nil
	// means every row is StatusPending (no tint). The host owns the mapping
	// from its own model / view.Item to this typed enum, so the library holds
	// no localized status strings.
	StatusOf func(it Item) Status

	items []Item
	rows  *SignalNodes
	sel   listselect.Mode
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
}

func (t *TargetHour) Init(_ Ctx) { t.ensure() }

func (t *TargetHour) SetSelectMode(on bool)        { t.sel.SetOn(on) }
func (t *TargetHour) SetDanger(on bool)            { t.sel.SetDanger(on) }
func (t *TargetHour) OnCheckedChange(fn func(int)) { t.sel.OnChange = fn }

func (t *TargetHour) itemIDs() []string {
	ids := make([]string, len(t.items))
	for i, it := range t.items {
		ids[i] = it.ID
	}
	return ids
}

func (t *TargetHour) CheckedIDs() []string {
	return t.sel.CheckedIDs(t.itemIDs())
}

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

func (t *TargetHour) Render() *Element {
	list := Ul().Set(clsList.AsAttr()).Attr("role", "listbox").BindChildren(t.rows)
	return Div().Set(clsListWrap.AsAttr()).
		BindStateFunc(widget.Open, func() bool { return t.sel.On().Get() }).
		Child(t.buildMasterCheck()).
		Child(list)
}

func (t *TargetHour) buildMasterCheck() *Element {
	allChecked := DeriveBool(func() bool {
		_ = t.sel.Changed().Get() // re-read after every toggle (see Mode.Changed)
		n := t.sel.Count()
		return n > 0 && n == len(t.items)
	})
	someChecked := DeriveBool(func() bool {
		_ = t.sel.Changed().Get()
		n := t.sel.Count()
		return n > 0 && n < len(t.items)
	})

	m := Span().Set(clsCheckAll.AsAttr()).
		Attr("role", "checkbox").
		BindAttrBool("aria-checked", allChecked).
		// glyph selector: trash while the danger tone is armed, pencil
		// otherwise — mirrors the row check's Invalid/Selected split, but
		// keyed on the MODE (sel.Danger), not on "this row is checked".
		BindState(widget.Invalid, DeriveBool(func() bool { return t.sel.On().Get() && t.sel.Danger().Get() })).
		BindState(widget.Selected, DeriveBool(func() bool { return t.sel.On().Get() && !t.sel.Danger().Get() })).
		// fill: solid when ALL marked, lighter wash when SOME marked.
		BindState(widget.Locked, allChecked).
		BindState(widget.Busy, someChecked).
		Child(trash.Ref.Render(string(clsCheckAllTrash))).
		Child(pencil.Ref.Render(string(clsCheckAllPencil))).
		Child(Span().Set(clsCheckAllCount.AsAttr()).
			BindTextFunc(func() string {
				_ = t.sel.Changed().Get()
				return fmt.Sprintf("%d / %d", t.sel.Count(), len(t.items))
			}))

	m.On("click", func(Event) {
		if n := t.sel.Count(); n > 0 && n == len(t.items) {
			t.sel.Clear()
			return
		}
		t.sel.CheckAll(t.itemIDs())
	})
	return m
}

func (t *TargetHour) buildRow(it Item) *Element {
	id := it.ID
	key := "th-" + id

	isEditCheckSig := DeriveBool(func() bool {
		_ = t.sel.Changed().Get()
		return t.sel.On().Get() && !t.sel.Danger().Get() && t.sel.IsChecked(id)
	})
	isDangerSig := DeriveBool(func() bool {
		_ = t.sel.Changed().Get()
		return t.sel.On().Get() && t.sel.Danger().Get() && t.sel.IsChecked(id)
	})
	isSelSig := DeriveBool(func() bool {
		_ = t.sel.Changed().Get()
		if t.sel.On().Get() {
			return isEditCheckSig.Get()
		}
		return t.Selected.Get() == id
	})
	isMarkedSig := DeriveBool(func() bool { return isSelSig.Get() || isDangerSig.Get() })

	row := Li().Set(clsRow.AsAttr()).
		ID(key).
		Key(key).
		Attr("role", "option").
		BindAttrBool("aria-selected", isMarkedSig).
		BindState(widget.Selected, isSelSig).
		BindState(widget.Invalid, isDangerSig)

	st := StatusPending
	if t.StatusOf != nil {
		st = t.StatusOf(it)
	}
	row.BindStateFunc(widget.Locked, func() bool { return st == StatusConfirmed })
	row.BindStateFunc(widget.Busy, func() bool { return st == StatusAttended })

	row.On("click", func(Event) {
		if t.sel.On().Get() {
			t.sel.Toggle(id)
			return
		}
		if t.OnSelect != nil {
			t.OnSelect(it)
		}
	})

	hour := Span().Set(clsHour.AsAttr()).Text(it.LeadMain)

	check := Span().Set(clsCheck.AsAttr()).
		BindState(widget.Selected, isEditCheckSig).
		BindState(widget.Invalid, isDangerSig).
		Child(trash.Ref.Render(string(clsCheckTrash))).
		Child(pencil.Ref.Render(string(clsCheckPencil)))

	content := Div().Set(clsContent.AsAttr()).
		Child(hour).
		Child(check).
		Child(Span().Set(clsLabel.AsAttr()).Text(it.Label))
	if it.Description != "" {
		content.Child(Span().Set(clsBadge.AsAttr()).
			Attr("title", it.Description).
			Text(fmt.Convert(it.Description).Truncate(badgeChars).String()))
	}

	row.Child(content)
	return row
}
