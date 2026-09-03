// Package targetdate is targetlist's sibling for rows that need a prominent
// leading badge — an hour, a day, anything view.Item.LeadTop/Main/Bottom
// carries — instead of a plain label. Same multi-selection mechanics as
// targetlist.
package targetdate

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

// badgeChars mirrors targetlist's own budget — see that package's comment.
const badgeChars = 16

// NameTargetDate is the widget identity.
const NameTargetDate = widget.Name("targetdate")

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
	PartBadge          = widget.Part("badge")
	PartLabel          = widget.Part("label")
	PartList           = widget.Part("list")
	PartLead           = widget.Part("lead")
	PartLeadStack      = widget.Part("lead-stack")
	PartLeadTop        = widget.Part("lead-top")
	PartLeadMain       = widget.Part("lead-main")
	PartLeadBottom     = widget.Part("lead-bottom")
)

var (
	clsListWrap       = NameTargetDate.Root()
	clsList           = NameTargetDate.Class(PartList)
	clsRow            = NameTargetDate.Class(PartRow)
	clsContent        = NameTargetDate.Class(PartContent)
	clsCheck          = NameTargetDate.Class(PartCheck)
	clsCheckTrash     = NameTargetDate.Class(PartCheckTrash)
	clsCheckPencil    = NameTargetDate.Class(PartCheckPencil)
	clsCheckAll       = NameTargetDate.Class(PartCheckAll)
	clsCheckAllTrash  = NameTargetDate.Class(PartCheckAllTrash)
	clsCheckAllPencil = NameTargetDate.Class(PartCheckAllPencil)
	clsCheckAllCount  = NameTargetDate.Class(PartCheckAllCount)
	clsLabel          = NameTargetDate.Class(PartLabel)
	clsBadge          = NameTargetDate.Class(PartBadge)
	clsLead           = NameTargetDate.Class(PartLead)
	clsLeadStack      = NameTargetDate.Class(PartLeadStack)
	clsLeadTop        = NameTargetDate.Class(PartLeadTop)
	clsLeadMain       = NameTargetDate.Class(PartLeadMain)
	clsLeadBottom     = NameTargetDate.Class(PartLeadBottom)
)

// Item is view.Item — see targetlist.Item's comment for why this is an
// alias, not a copy. TargetDate is the one row that actually reads
// LeadTop/Main/Bottom.
type Item = view.Item

// TargetDate is a selectable list of records with a prominent leading badge
// (LeadTop/Main/Bottom) instead of a plain label lead-in.
type TargetDate struct {
	Element

	Selected *SignalString

	OnSelect func(it Item)

	items []Item
	rows  *SignalNodes
	sel   listselect.Mode
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
}

func (t *TargetDate) Init(_ Ctx) { t.ensure() }

func (t *TargetDate) SetSelectMode(on bool)        { t.sel.SetOn(on) }
func (t *TargetDate) SetDanger(on bool)            { t.sel.SetDanger(on) }
func (t *TargetDate) OnCheckedChange(fn func(int)) { t.sel.OnChange = fn }

func (t *TargetDate) itemIDs() []string {
	ids := make([]string, len(t.items))
	for i, it := range t.items {
		ids[i] = it.ID
	}
	return ids
}

func (t *TargetDate) CheckedIDs() []string {
	return t.sel.CheckedIDs(t.itemIDs())
}

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

func (t *TargetDate) Render() *Element {
	list := Ul().Set(clsList.AsAttr()).Attr("role", "listbox").BindChildren(t.rows)
	return Div().Set(clsListWrap.AsAttr()).
		BindStateFunc(widget.Open, func() bool { return t.sel.On().Get() }).
		Child(t.buildMasterCheck()).
		Child(list)
}

func (t *TargetDate) buildMasterCheck() *Element {
	allChecked := DeriveBool(func() bool {
		_ = t.sel.Changed().Get() // re-read after every toggle (see Mode.Changed)
		n := t.sel.Count()
		return n > 0 && n == len(t.items)
	})
	m := Span().Set(clsCheckAll.AsAttr()).
		Attr("role", "checkbox").
		BindAttrBool("aria-checked", allChecked).
		// glyph selector: trash while the danger tone is armed, pencil
		// otherwise — mirrors the row check's Invalid/Selected split, but
		// keyed on the MODE (sel.Danger), not on "this row is checked".
		BindState(widget.Invalid, DeriveBool(func() bool { return t.sel.On().Get() && t.sel.Danger().Get() })).
		BindState(widget.Selected, DeriveBool(func() bool { return t.sel.On().Get() && !t.sel.Danger().Get() })).
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

func (t *TargetDate) buildRow(it Item) *Element {
	id := it.ID
	key := "td-" + id

	// See targetlist.buildRow for the split: isEditCheck is the narrow
	// "marked for edit in selection mode"; isSel widens it with the
	// normal-mode "loaded record" highlight for the ROW only. The CHECK box
	// binds the narrow one, so a row merely loaded in normal mode reveals no
	// glyph.
	isEditCheckSig := DeriveBool(func() bool {
		_ = t.sel.Changed().Get() // re-read IsChecked after every tap (see Mode.Changed)
		return t.sel.On().Get() && !t.sel.Danger().Get() && t.sel.IsChecked(id)
	})
	isDangerSig := DeriveBool(func() bool {
		_ = t.sel.Changed().Get() // re-read IsChecked after every tap (see Mode.Changed)
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

	row.On("click", func(Event) {
		if t.sel.On().Get() {
			t.sel.Toggle(id)
			return
		}
		if t.OnSelect != nil {
			t.OnSelect(it)
		}
	})

	lead := Div().Set(clsLead.AsAttr()).Child(
		Div().Set(clsLeadStack.AsAttr()).Child(
			Span().Set(clsLeadTop.AsAttr()).Text(it.LeadTop),
			Span().Set(clsLeadMain.AsAttr()).Text(it.LeadMain),
			Span().Set(clsLeadBottom.AsAttr()).Text(it.LeadBottom),
		),
	)

	// The box owns which glyph shows and its colour via its own
	// Selected (edit) / Invalid (delete) state, written only in selection
	// mode. Nothing here hangs off the ROW's state. See targetlist.buildRow.
	check := Span().Set(clsCheck.AsAttr()).
		BindState(widget.Selected, isEditCheckSig).
		BindState(widget.Invalid, isDangerSig).
		Child(trash.Ref.Render(string(clsCheckTrash))).
		Child(pencil.Ref.Render(string(clsCheckPencil)))

	content := Div().Set(clsContent.AsAttr()).
		Child(check).
		Child(Span().Set(clsLabel.AsAttr()).Text(it.Label))
	if it.Description != "" {
		content.Child(Span().Set(clsBadge.AsAttr()).
			Attr("title", it.Description).
			Text(fmt.Convert(it.Description).Truncate(badgeChars).String()))
	}

	row.Child(lead)
	row.Child(content)

	return row
}
