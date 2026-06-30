//go:build wasm

package themetoggle

import . "github.com/tinywasm/dom"

func (t *ThemeToggle) Init(_ Ctx) {
	t.theme = NewString("") // ThemeAuto
	if LocalStorageAvailable() {
		if s, err := LocalStorageGet(storageKey); err == nil && valid(Theme(s)) {
			t.theme.Set(s)                             // value ready before first paint → correct icon, no flash
			SetDocumentAttr("data-theme", string(s)) // apply stored theme
		}
	}
}

func (t *ThemeToggle) onClick() {
	next := cycle(Theme(t.theme.Get()))
	SetDocumentAttr("data-theme", string(next)) // applies the theme
	if LocalStorageAvailable() {
		if next == ThemeAuto {
			LocalStorageDel(storageKey)
		} else {
			LocalStorageSet(storageKey, string(next))
		}
	}
	t.theme.Set(string(next)) // patches icon + labels surgically
}
