//go:build !wasm

package sitenav

import (
	"strings"
	"testing"

	"github.com/tinywasm/js"
)

func TestSiteNav_AccessibilityAttributes(t *testing.T) {
	sn := &SiteNav{
		WideLogoSrc:    "/logo.png",
		CompactLogoSrc: "/logo-sm.png",
		LogoAlt:        "Logo",
		Links: []NavItem{
			{Label: "Home", Href: "/"},
		},
	}

	html := sn.Render().String()

	if !strings.Contains(html, `aria-controls='sitenav-menu'`) {
		t.Errorf("expected aria-controls='sitenav-menu' in markup, got: %s", html)
	}

	if !strings.Contains(html, `aria-expanded='false'`) {
		t.Errorf("expected aria-expanded='false' in markup, got: %s", html)
	}

	if !strings.Contains(html, `id='sitenav-menu'`) {
		t.Errorf("expected id='sitenav-menu' in markup, got: %s", html)
	}
}

func TestSiteNav_RenderJSIdempotentAndSafe(t *testing.T) {
	sn := &SiteNav{}
	scripts := sn.RenderJS()

	if len(scripts) != 1 || scripts[0].Content == "" {
		t.Fatal("RenderJS returned no script content")
	}
	jsStr := scripts[0].Content

	if !strings.Contains(jsStr, "__sitenavInit") {
		t.Errorf("RenderJS does not include idempotency guard __sitenavInit")
	}

	if !strings.Contains(jsStr, "document.addEventListener") {
		t.Errorf("RenderJS does not use document event delegation")
	}
}

func TestSiteNav_RenderJSSatisfiesSSRJSProvider(t *testing.T) {
	// sitec.RegisterComponents type-asserts for RenderJS() []*js.Script at
	// runtime; a value-slice return silently fails that assertion instead of
	// erroring, so this pins the pointer-slice shape as a compile-time check.
	var _ interface {
		RenderJS() []*js.Script
	} = (*SiteNav)(nil)
}
