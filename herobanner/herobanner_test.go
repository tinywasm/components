//go:build !wasm

package herobanner

import (
	"strings"
	"testing"
)

func TestHeroBanner_PrefersReducedMotion(t *testing.T) {
	hb := &HeroBanner{
		Title:    "Clínica de Excelencia",
		Subtitle: "Cuidando de tu salud",
		Image:    "/img/hero1.jpg",
	}

	sheet := hb.RenderCSS()
	if sheet == nil {
		t.Fatal("RenderCSS returned nil")
	}

	cssStr := sheet.String()
	if !strings.Contains(cssStr, "prefers-reduced-motion") {
		t.Errorf("RenderCSS output does not contain 'prefers-reduced-motion': %s", cssStr)
	}
}
