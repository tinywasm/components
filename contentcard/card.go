package contentcard

import (
	. "github.com/tinywasm/css"
	. "github.com/tinywasm/dom"
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
	card := Div(clsCard.AsAttr())

	if c.Header != nil {
		card.Add(Div(clsCardHeader.AsAttr()).Add(c.Header))
	}

	if c.Body != nil {
		card.Add(Div(clsCardBody.AsAttr()).Add(c.Body))
	}

	if c.Footer != nil {
		card.Add(Div(clsCardFooter.AsAttr()).Add(c.Footer))
	}

	return card
}
