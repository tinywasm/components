package card

import (
	"github.com/tinywasm/dom"
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

	card := dom.Div().
		Class("card")

	if c.Header != nil {
		card.Add(dom.Div().Class("card-header").Add(c.Header))
	}

	if c.Body != nil {
		card.Add(dom.Div().Class("card-body").Add(c.Body))
	}

	if c.Footer != nil {
		card.Add(dom.Div().Class("card-footer").Add(c.Footer))
	}

	return card
}
