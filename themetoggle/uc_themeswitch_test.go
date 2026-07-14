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

func TestThemeToggle_Init_NoSavedValue_DefaultsLight(t *testing.T) {
	setUp()
	ts := &ThemeToggle{}
	ts.Init(nil)

	got := dom.GetDocumentAttr("data-theme")
	if got != "light" {
		t.Errorf("expected data-theme=light by default, got %q", got)
	}
	if ts.theme.Get() != "light" {
		t.Errorf("expected signal=light by default, got %q", ts.theme.Get())
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

func TestThemeToggle_Init_InvalidValue_DefaultsLight(t *testing.T) {
	setUp()
	dom.LocalStorageSet(storageKey, "xyz")

	ts := &ThemeToggle{}
	ts.Init(nil)

	got := dom.GetDocumentAttr("data-theme")
	if got != "light" {
		t.Errorf("expected data-theme=light after invalid value, got %q", got)
	}
	if ts.theme.Get() != "light" {
		t.Errorf("expected signal=light, got %q", ts.theme.Get())
	}
}

func TestThemeToggle_Render_Initial(t *testing.T) {
	setUp()
	ts := &ThemeToggle{}
	ts.Init(nil)
	el := ts.Render()

	// default theme is light → icon is ☀️
	got := el.String()
	if !fmt.Contains(got, icon(TsThemeLight)) {
		t.Errorf("expected icon %s in rendered element, got %s", icon(TsThemeLight), got)
	}
}
