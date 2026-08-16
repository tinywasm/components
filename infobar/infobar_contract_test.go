package infobar

import (
	"testing"

	"github.com/tinywasm/dom"
)

var _ dom.Component = (*InfoBar)(nil)

func TestInfoBar_RenderIdempotent(t *testing.T) {
	ib := &InfoBar{
		Items: []InfoItem{
			{Icon: IconPhone, Text: "+56 9 1234 5678", Href: "tel:+56912345678"},
			{Icon: IconMail, Text: "contacto@clinica.cl", Href: "mailto:contacto@clinica.cl"},
			{Icon: IconPin, Text: "Av. Providencia 1234"},
		},
	}

	first := ib.Render().String()
	second := ib.Render().String()

	if first != second {
		t.Errorf("Render() not idempotent: first=%q second=%q", first, second)
	}
}
