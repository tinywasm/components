//go:build !wasm

package herobanner

import (
	"fmt"
	"strings"
	"testing"

	"webtyp.com/image"
)

func TestHeroBanner_AutoRotateAndReducedMotion(t *testing.T) {
	hb := &HeroBanner{
		Title:    "Clínica de Excelencia",
		Subtitle: "Cuidando de tu salud",
		Images:   []string{"/img/hero1.jpg", "/img/hero2.jpg"},
	}

	sheet := hb.RenderCSS()
	if sheet == nil {
		t.Fatal("RenderCSS returned nil")
	}

	cssStr := sheet.String()

	if !strings.Contains(cssStr, "tw-auto-rotate") {
		t.Errorf("RenderCSS output does not contain keyframe animation 'tw-auto-rotate': %s", cssStr)
	}

	if !strings.Contains(cssStr, "prefers-reduced-motion") {
		t.Errorf("RenderCSS output does not contain 'prefers-reduced-motion': %s", cssStr)
	}
}

// TestHeroBanner_SlidesAreVoidElements pins the markup shape of the media
// layer. NewElement("img") has no notion of HTML's content model, so a missing
// NoCloseTag() silently produces `<img ...></img>` — invalid markup that
// browsers paper over and nobody notices until a strict parser chokes on it.
func TestHeroBanner_SlidesAreVoidElements(t *testing.T) {
	hb := &HeroBanner{Images: []string{"/img/a.jpg", "/img/b.jpg"}}

	got := hb.Render().String()

	if strings.Contains(got, "</img>") {
		t.Errorf("slides emit a closing tag for the void element <img>:\n%s", got)
	}
}

// TestHeroBanner_OnlyFirstSlideIsEager keeps the hero from requesting every
// photograph in the rotation before first paint: only the layer that is
// actually visible on load is the LCP candidate.
func TestHeroBanner_OnlyFirstSlideIsEager(t *testing.T) {
	hb := &HeroBanner{Images: []string{"/img/a.jpg"}}

	got := hb.Render().String()

	if n := strings.Count(got, `loading='eager'`); n != 1 {
		t.Errorf("expected exactly 1 eager slide, got %d:\n%s", n, got)
	}
	if n := strings.Count(got, `loading='lazy'`); n != autoRotateLayers-1 {
		t.Errorf("expected %d lazy slides, got %d:\n%s", autoRotateLayers-1, n, got)
	}
}

func TestHeroBanner_ResponsiveImages(t *testing.T) {
	hb := &HeroBanner{
		Images: []string{"/img/hero.jpg"},
	}

	got := hb.Render().String()

	expectedSrcset := fmt.Sprintf(
		"srcset='/img/hero.S.jpg %dw, /img/hero.M.jpg %dw, /img/hero.L.jpg %dw'",
		image.VariantS.Width(), image.VariantM.Width(), image.VariantL.Width(),
	)
	if !strings.Contains(got, expectedSrcset) {
		t.Errorf("expected output to contain %q, got:\n%s", expectedSrcset, got)
	}

	expectedSizes := `sizes='100vw'`
	if !strings.Contains(got, expectedSizes) {
		t.Errorf("expected output to contain %q, got:\n%s", expectedSizes, got)
	}

	// alt should be empty in all 6 layers
	if n := strings.Count(got, `alt=''`); n != autoRotateLayers {
		t.Errorf("expected %d empty alt attributes, got %d:\n%s", autoRotateLayers, n, got)
	}
}

func TestHeroBanner_NoImages(t *testing.T) {
	hb := &HeroBanner{
		Title: "Banner title",
	}

	got := hb.Render().String()

	if strings.Contains(got, clsHeroMedia.String()) {
		t.Errorf("expected no media layer when Images is empty, got:\n%s", got)
	}
}
