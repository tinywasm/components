package herobanner

import (
	"testing"

	"github.com/tinywasm/dom"
)

var _ dom.Component = (*HeroBanner)(nil)

func TestHeroBanner_RenderIdempotent(t *testing.T) {
	hb := &HeroBanner{
		Title:    "Clínica de Excelencia",
		Subtitle: "Cuidando de tu salud y la de tu familia siempre",
		Images:   []string{"/img/hero1.jpg", "/img/hero2.jpg"},
	}

	first := hb.Render().String()
	second := hb.Render().String()

	if first != second {
		t.Errorf("Render() not idempotent: first=%q second=%q", first, second)
	}
}
