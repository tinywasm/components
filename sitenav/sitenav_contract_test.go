package sitenav

import (
	"testing"

	"github.com/tinywasm/dom"
)

var _ dom.Component = (*SiteNav)(nil)

func TestSiteNav_RenderIdempotent(t *testing.T) {
	sn := &SiteNav{
		WideLogoSrc:    "/images/logo-wide.png",
		CompactLogoSrc: "/images/logo-compact.png",
		LogoAlt:        "Clínica Médica",
		Links: []NavItem{
			{Label: "Inicio", Href: "/", Active: true},
			{Label: "Servicios", Href: "/servicios"},
			{Label: "Contacto", Href: "/contacto"},
		},
	}

	first := sn.Render().String()
	second := sn.Render().String()

	if first != second {
		t.Errorf("Render() not idempotent: first=%q second=%q", first, second)
	}
}
