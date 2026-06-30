//go:build wasm

package themetoggle

import (
	"testing"

	"github.com/tinywasm/dom"
)

func setUp() {
	dom.LocalStorageDel(storageKey)
	dom.SetDocumentAttr("data-theme", "")
}

func TestThemeToggle_Init_NoSavedValue_StaysAuto(t *testing.T) {
	setUp()
	ts := &ThemeToggle{}
	ts.Init(nil)

	got := dom.GetDocumentAttr("data-theme")
	if got != "" {
		t.Errorf("expected empty data-theme, got %q", got)
	}
	if ts.theme.Get() != "" {
		t.Errorf("expected empty signal, got %q", ts.theme.Get())
	}
}

func TestThemeToggle_Init_RestoresDark(t *testing.T) {
	setUp()
	dom.LocalStorageSet(storageKey, "dark")

	ts := &ThemeToggle{}
	ts.Init(nil)

	got := dom.GetDocumentAttr("data-theme")
	if got != "dark" {
		t.Errorf("expected data-theme=dark, got %q", got)
	}
	if ts.theme.Get() != "dark" {
		t.Errorf("expected signal=dark, got %q", ts.theme.Get())
	}
}

func TestThemeToggle_Init_InvalidValue_StaysAuto(t *testing.T) {
	setUp()
	dom.LocalStorageSet(storageKey, "xyz")

	ts := &ThemeToggle{}
	ts.Init(nil)

	got := dom.GetDocumentAttr("data-theme")
	if got != "" {
		t.Errorf("expected empty data-theme after invalid value, got %q", got)
	}
	if ts.theme.Get() != "" {
		t.Errorf("expected empty signal, got %q", ts.theme.Get())
	}
}

func TestThemeToggle_Render_Initial(t *testing.T) {
	setUp()
	ts := &ThemeToggle{}
	ts.Init(nil)
	el := ts.Render()

	// icon for auto is ◑
	if el.GetText() != icon(ThemeAuto) {
		t.Errorf("expected icon %s, got %s", icon(ThemeAuto), el.GetText())
	}
}
