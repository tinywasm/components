package card

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/dom"
)

var (
	ClsCard       css.Class = "card"
	ClsCardHeader css.Class = "card-header"
	ClsCardBody   css.Class = "card-body"
	ClsCardFooter css.Class = "card-footer"
)

type Card struct {
	*dom.Element
	Header dom.Component
	Body   dom.Component
	Footer dom.Component
}

func (c *Card) Render() *dom.Element {
	if c.Element == nil {
		c.Element = &dom.Element{}
	}

	card := dom.Div(dom.Class(ClsCard))

	if c.Header != nil {
		card.Add(dom.Div(dom.Class(ClsCardHeader)).Add(c.Header))
	}

	if c.Body != nil {
		card.Add(dom.Div(dom.Class(ClsCardBody)).Add(c.Body))
	}

	if c.Footer != nil {
		card.Add(dom.Div(dom.Class(ClsCardFooter)).Add(c.Footer))
	}

	return card
}
