package themetoggle

import (
	. "github.com/tinywasm/css"
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

var (
	clsTsBtn Class = "ts-btn"
)

// storageKey identifica la entrada de localStorage del componente.
const storageKey = "tinywasm-themeswitch"

// Theme representa el estado de tema del componente.
// Solo hay dos estados: "dark" y "light", escritos literalmente en data-theme.
type Theme string

const (
	ThemeDark  Theme = "dark"
	ThemeLight Theme = "light"
)

// ThemeToggle es un botón flotante que cicla entre los 3 modos de tema.
// Restaura automáticamente el tema guardado en localStorage al montarse.
//
//	ts := &themetoggle.ThemeToggle{}
//	Append("body", ts)
type ThemeToggle struct {
	Element
	theme *SignalString // "", "dark", "light"
}

func (t *ThemeToggle) Render() *Element {
	// labelSig is used twice (title + aria-label) → a named shared computed. Auto-tracked: no deps list.
	labelSig := DeriveString(func() string { return label(Theme(t.theme.Get())) })

	return Button().
		BindTextFunc(func() string { return icon(Theme(t.theme.Get())) }). // computed icon, auto-tracked
		BindAttr("title", labelSig).
		BindAttr("aria-label", labelSig).
		Set(clsTsBtn.AsAttr()).
		On("click", func(Event) {
			t.onClick()
		})
}

// cycle alterna entre los 2 estados. Switch (no map) — TinyGo.
func cycle(current Theme) Theme {
	if current == ThemeLight {
		return ThemeDark
	}
	return ThemeLight
}

// icon retorna el símbolo visible del botón. Switch (no map) — TinyGo.
func icon(theme Theme) string {
	if theme == ThemeLight {
		return "☀"
	}
	return "🌙"
}

// label retorna el nombre del modo actual (usado como tooltip y aria-label).
func label(theme Theme) string {
	if theme == ThemeLight {
		return "light"
	}
	return "dark"
}

func valid(t Theme) bool {
	return t == ThemeDark || t == ThemeLight
}
