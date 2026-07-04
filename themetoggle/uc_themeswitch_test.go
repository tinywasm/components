//go:build wasm

package themetoggle

import (
	"testing"

	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
)

func setUp() {
	dom.LocalStorageDel(storageKey)
	dom.SetDocumentAttr("data-theme", "")
}

func TestThemeToggle_Init_NoSavedValue_DefaultsDark(t *testing.T) {
	setUp()
	ts := &ThemeToggle{}
	ts.Init(nil)

	got := dom.GetDocumentAttr("data-theme")
	if got != "dark" {
		t.Errorf("expected data-theme=dark by default, got %q", got)
	}
	if ts.theme.Get() != "dark" {
		t.Errorf("expected signal=dark by default, got %q", ts.theme.Get())
	}
}

func TestThemeToggle_Init_RestoresLight(t *testing.T) {
	setUp()
	dom.LocalStorageSet(storageKey, "light")

	ts := &ThemeToggle{}
	ts.Init(nil)

	got := dom.GetDocumentAttr("data-theme")
	if got != "light" {
		t.Errorf("expected data-theme=light, got %q", got)
	}
	if ts.theme.Get() != "light" {
		t.Errorf("expected signal=light, got %q", ts.theme.Get())
	}
}

func TestThemeToggle_Init_InvalidValue_DefaultsDark(t *testing.T) {
	setUp()
	dom.LocalStorageSet(storageKey, "xyz")

	ts := &ThemeToggle{}
	ts.Init(nil)

	got := dom.GetDocumentAttr("data-theme")
	if got != "dark" {
		t.Errorf("expected data-theme=dark after invalid value, got %q", got)
	}
	if ts.theme.Get() != "dark" {
		t.Errorf("expected signal=dark, got %q", ts.theme.Get())
	}
}

func TestThemeToggle_Render_Initial(t *testing.T) {
	setUp()
	ts := &ThemeToggle{}
	ts.Init(nil)
	el := ts.Render()

	// default theme is dark → icon is 🌙
	got := el.String()
	if !fmt.Contains(got, icon(ThemeDark)) {
		t.Errorf("expected icon %s in rendered element, got %s", icon(ThemeDark), got)
	}
}
