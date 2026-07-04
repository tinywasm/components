package contentcard

import (
	. "github.com/tinywasm/css"
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

var (
	clsCard       Class = "card"
	clsCardHeader Class = "card-header"
	clsCardBody   Class = "card-body"
	clsCardFooter Class = "card-footer"
)

type ContentCard struct {
	Element
	Header Component
	Body   Component
	Footer Component
}

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
