package contentcard

import (
	. "webtyp.com/dom"
	. "webtyp.com/html"
	"webtyp.com/widget"
)

// NameContentCard is the widget name.
const NameContentCard = widget.Name("contentcard")

const (
	PartHeader = widget.Part("header")
	PartBody   = widget.Part("body")
	PartFooter = widget.Part("footer")
)

var (
	clsCard       = NameContentCard.Root()
	clsCardHeader = NameContentCard.Class(PartHeader)
	clsCardBody   = NameContentCard.Class(PartBody)
	clsCardFooter = NameContentCard.Class(PartFooter)
)

type ContentCard struct {
	Element
	Header Component
	Body   Component
	Footer Component
}

func (c *ContentCard) WidgetName() widget.Name { return NameContentCard }
func (c *ContentCard) WidgetKind() widget.Kind { return widget.Region }

func (c *ContentCard) Render() *Element {
	card := Div().Set(clsCard.AsAttr())

	if c.Header != nil {
		card.Child(Div().Set(clsCardHeader.AsAttr()).Child(c.Header))
	}

	if c.Body != nil {
		card.Child(Div().Set(clsCardBody.AsAttr()).Child(c.Body))
	}

	if c.Footer != nil {
		card.Child(Div().Set(clsCardFooter.AsAttr()).Child(c.Footer))
	}

	return card
}
