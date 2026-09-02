// Package targetdate is targetlist's sibling for rows that need a prominent
// leading badge — an hour, a day, anything view.Item.LeadTop/Main/Bottom
// carries — instead of a plain label. Same multi-selection mechanics as
// targetlist.
package targetdate

import (
	. "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/svg"
	"github.com/tinywasm/view"
	"github.com/tinywasm/widget"

	"github.com/tinywasm/components/listselect"
)

// badgeChars mirrors targetlist's own budget — see that package's comment.
const badgeChars = 16

// NameTargetDate is the widget identity.
const NameTargetDate = widget.Name("targetdate")

const (
	PartRow        = widget.Part("row")
	PartContent    = widget.Part("content")
	PartCheck      = widget.Part("check")
	PartBadge      = widget.Part("badge")
	PartLabel      = widget.Part("label")
	PartList       = widget.Part("list")
	PartLead       = widget.Part("lead")
	PartLeadStack  = widget.Part("lead-stack")
	PartLeadTop    = widget.Part("lead-top")
	PartLeadMain   = widget.Part("lead-main")
	PartLeadBottom = widget.Part("lead-bottom")
)

var (
	clsListWrap   = NameTargetDate.Root()
	clsList       = NameTargetDate.Class(PartList)
	clsRow        = NameTargetDate.Class(PartRow)
	clsContent    = NameTargetDate.Class(PartContent)
	clsCheck      = NameTargetDate.Class(PartCheck)
	clsLabel      = NameTargetDate.Class(PartLabel)
	clsBadge      = NameTargetDate.Class(PartBadge)
	clsLead       = NameTargetDate.Class(PartLead)
	clsLeadStack  = NameTargetDate.Class(PartLeadStack)
	clsLeadTop    = NameTargetDate.Class(PartLeadTop)
	clsLeadMain   = NameTargetDate.Class(PartLeadMain)
	clsLeadBottom = NameTargetDate.Class(PartLeadBottom)
)

const iconCheck = svg.Icon("td-check")

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

func (t *TargetDate) SetSelectMode(on bool)       { t.sel.SetOn(on) }
func (t *TargetDate) OnCheckedChange(fn func(int)) { t.sel.OnChange = fn }

func (t *TargetDate) CheckedIDs() []string {
	ids := make([]string, len(t.items))
	for i, it := range t.items {
		ids[i] = it.ID
	}
	return t.sel.CheckedIDs(ids)
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
		Child(list)
}

func (t *TargetDate) buildRow(it Item) *Element {
	id := it.ID
	key := "td-" + id

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

	lead := Div().Set(clsLead.AsAttr()).Child(
		Div().Set(clsLeadStack.AsAttr()).Child(
			Span().Set(clsLeadTop.AsAttr()).Text(it.LeadTop),
			Span().Set(clsLeadMain.AsAttr()).Text(it.LeadMain),
			Span().Set(clsLeadBottom.AsAttr()).Text(it.LeadBottom),
		),
	)

	check := Span().Set(clsCheck.AsAttr()).Child(iconCheck.Render(string(clsCheck)))

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
