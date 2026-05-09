//go:build wasm

package themeswitch

import (
	"testing"

	"github.com/tinywasm/dom"
)

func setUp() {
	dom.LocalStorageDel(storageKey)
	dom.SetDocumentAttr("data-theme", "")
}

func TestThemeSwitch_OnMount_NoSavedValue_StaysAuto(t *testing.T) {
	setUp()
	ts := &ThemeSwitch{}
	ts.OnMount()

	got := dom.GetDocumentAttr("data-theme")
	if got != "" {
		t.Errorf("expected empty data-theme, got %q", got)
	}
}

func TestThemeSwitch_OnMount_RestoresDark(t *testing.T) {
	setUp()
	dom.LocalStorageSet(storageKey, "dark")

	ts := &ThemeSwitch{}
	ts.OnMount()

	got := dom.GetDocumentAttr("data-theme")
	if got != "dark" {
		t.Errorf("expected data-theme=dark, got %q", got)
	}
}

func TestThemeSwitch_OnMount_InvalidValue_Cleans(t *testing.T) {
	setUp()
	dom.LocalStorageSet(storageKey, "xyz")

	ts := &ThemeSwitch{}
	ts.OnMount()

	got := dom.GetDocumentAttr("data-theme")
	if got != "" {
		t.Errorf("expected empty data-theme after invalid value, got %q", got)
	}

	val, err := dom.LocalStorageGet(storageKey)
	if err != nil {
		t.Fatalf("LocalStorageGet failed: %v", err)
	}
	if val != "" {
		t.Errorf("expected localStorage to be cleaned, but got %q", val)
	}
}

func TestThemeSwitch_Click_CyclesAndPersists(t *testing.T) {
	setUp()
	ts := &ThemeSwitch{}
	ts.OnMount()

	// Auto ("") -> click -> Dark
	ts.onClick(nil)
	if got := dom.GetDocumentAttr("data-theme"); got != "dark" {
		t.Errorf("expected dark after 1st click, got %q", got)
	}
	if val, _ := dom.LocalStorageGet(storageKey); val != "dark" {
		t.Errorf("expected dark in localStorage, got %q", val)
	}

	// Dark -> click -> Light
	ts.onClick(nil)
	if got := dom.GetDocumentAttr("data-theme"); got != "light" {
		t.Errorf("expected light after 2nd click, got %q", got)
	}
	if val, _ := dom.LocalStorageGet(storageKey); val != "light" {
		t.Errorf("expected light in localStorage, got %q", val)
	}

	// Light -> click -> Auto
	ts.onClick(nil)
	if got := dom.GetDocumentAttr("data-theme"); got != "" {
		t.Errorf("expected auto after 3rd click, got %q", got)
	}
	val, _ := dom.LocalStorageGet(storageKey)
	if val != "" {
		t.Errorf("expected empty localStorage after cycling to auto, got %q", val)
	}
}
