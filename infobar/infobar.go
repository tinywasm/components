package infobar

import (
	. "webtyp.com/dom"
	. "webtyp.com/html"
	"webtyp.com/svg"
	"webtyp.com/widget"
)

// NameInfoBar is the widget name for infobar.
const NameInfoBar = widget.Name("infobar")

const (
	PartItem = widget.Part("item")
	PartIcon = widget.Part("icon")
	PartText = widget.Part("text")
)

var (
	clsInfoBar  = NameInfoBar.Root()
	clsInfoItem = NameInfoBar.Class(PartItem)
	clsInfoIcon = NameInfoBar.Class(PartIcon)
	clsInfoText = NameInfoBar.Class(PartText)
)

// Default icons provided by infobar (geometry provided in svg.go).
const (
	IconPhone = svg.Icon("infobar-phone")
	IconMail  = svg.Icon("infobar-mail")
	IconPin   = svg.Icon("infobar-pin")
	IconClock = svg.Icon("infobar-clock")
)

// InfoItem represents an entry in the contact/information bar.
type InfoItem struct {
	Icon svg.Icon
	Text string
	Href string
}

// InfoBar renders a top strip of contact and location information.
type InfoBar struct {
	Element
	Items []InfoItem
}

func (ib *InfoBar) WidgetName() widget.Name { return NameInfoBar }
func (ib *InfoBar) WidgetKind() widget.Kind { return widget.Region }

func (ib *InfoBar) Render() *Element {
	bar := Div().Set(clsInfoBar.AsAttr())

	for _, item := range ib.Items {
		var content *Element
		if item.Href != "" {
			content = A(item.Href)
		} else {
			content = Span()
		}

		content.Set(clsInfoItem.AsAttr())

		if item.Icon != "" {
			content.Child(item.Icon.Render(string(clsInfoIcon)))
		}

		if item.Text != "" {
			content.Child(Span().Set(clsInfoText.AsAttr()).Text(item.Text))
		}

		bar.Child(content)
	}

	return bar
}
