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
type TsTheme string

const (
	TsThemeAuto  TsTheme = ""
	TsThemeDark  TsTheme = "dark"
	TsThemeLight TsTheme = "light"
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

// cycle alterna entre los 3 estados. Switch (no map) — TinyGo.
func cycle(current TsTheme) TsTheme {
	switch current {
	case TsThemeAuto:
		return TsThemeDark
	case TsThemeDark:
		return TsThemeLight
	default:
		return TsThemeAuto
	}
}

// icon retorna el símbolo visible del botón. Switch (no map) — TinyGo.
func icon(theme TsTheme) string {
	switch theme {
	case TsThemeDark:
		return "🌙"
	case TsThemeLight:
		return "☀"
	default:
		return "🌓"
	}
}

// label retorna el nombre del modo actual (usado como tooltip y aria-label).
func label(theme TsTheme) string {
	switch theme {
	case TsThemeDark:
		return "dark"
	case TsThemeLight:
		return "light"
	default:
		return "auto"
	}
}

func valid(t TsTheme) bool {
	return t == TsThemeAuto || t == TsThemeDark || t == TsThemeLight
}
