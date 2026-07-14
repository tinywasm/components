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

// TsTheme representa el estado de tema del componente.
// Solo hay 2 estados: dark o light. No existe un tercer estado "auto" —
// si una app quiere seguir la preferencia del SO, eso se resuelve fuera de
// este componente (p. ej. eligiendo el valor inicial antes de montar
// ThemeToggle), no como un estado más del ciclo.
type TsTheme string

const (
	TsThemeDark  TsTheme = "dark"
	TsThemeLight TsTheme = "light"
)

// defaultTheme es el tema por defecto cuando no hay nada guardado en
// localStorage (o en SSR, donde no hay localStorage).
const defaultTheme = TsThemeLight

// ThemeToggle es un botón flotante que alterna entre dark y light.
// Restaura automáticamente el tema guardado en localStorage al montarse.
//
//	ts := &themetoggle.ThemeToggle{}
//	Append("body", ts)
type ThemeToggle struct {
	Element
	theme *SignalString // "dark" | "light"
}

func (t *ThemeToggle) Render() *Element {
	// labelSig is used twice (title + aria-label) → a named shared computed. Auto-tracked: no deps list.
	labelSig := DeriveString(func() string { return label(TsTheme(t.theme.Get())) })

	return Button().
		BindTextFunc(func() string { return icon(TsTheme(t.theme.Get())) }). // computed icon, auto-tracked
		BindAttr("title", labelSig).
		BindAttr("aria-label", labelSig).
		Set(clsTsBtn.AsAttr()).
		On("click", func(Event) {
			t.onClick()
		})
}

// toggle alterna entre los 2 estados. Switch (no map) — TinyGo.
func toggle(current TsTheme) TsTheme {
	switch current {
	case TsThemeDark:
		return TsThemeLight
	default:
		return TsThemeDark
	}
}

// icon retorna el símbolo visible del botón. Switch (no map) — TinyGo.
func icon(theme TsTheme) string {
	switch theme {
	case TsThemeDark:
		return "🌙"
	default:
		return "☀️" // U+2600 + U+FE0F (variation selector) forces emoji presentation, same size as 🌙
	}
}

// label retorna el nombre del modo actual (usado como tooltip y aria-label).
func label(theme TsTheme) string {
	switch theme {
	case TsThemeDark:
		return "dark"
	default:
		return "light"
	}
}

func valid(t TsTheme) bool {
	return t == TsThemeDark || t == TsThemeLight
}
