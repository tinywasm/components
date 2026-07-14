//go:build wasm

package themetoggle

import . "github.com/tinywasm/dom"

func (t *ThemeToggle) Init(_ Ctx) {
	t.theme = NewString(string(defaultTheme))
	if LocalStorageAvailable() {
		if s, err := LocalStorageGet(storageKey); err == nil && valid(TsTheme(s)) {
			t.theme.Set(s) // value ready before first paint → correct icon, no flash
		}
	}
	SetDocumentAttr("data-theme", t.theme.Get()) // apply theme so colors match the icon
}

func (t *ThemeToggle) onClick() {
	next := toggle(TsTheme(t.theme.Get()))
	SetDocumentAttr("data-theme", string(next)) // applies the theme
	if LocalStorageAvailable() {
		LocalStorageSet(storageKey, string(next))
	}
	t.theme.Set(string(next)) // patches icon + labels surgically
}
