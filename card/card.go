package card

import (
	"github.com/tinywasm/components"
	dom "github.com/tinywasm/components/internal/dom"
)

type Card struct {
	components.BaseComponent
	Header dom.Component
	Body   dom.Component
	Footer dom.Component
}

func (c *Card) Render() dom.Node {
	card := dom.Div().
		ID(c.ID()).
		Class("card")

	if c.Header != nil {
		card.Append(dom.Div().Class("card-header").Append(c.Header))
	}

	if c.Body != nil {
		card.Append(dom.Div().Class("card-body").Append(c.Body))
	}

	if c.Footer != nil {
		card.Append(dom.Div().Class("card-footer").Append(c.Footer))
	}

	return card.ToNode()
}
