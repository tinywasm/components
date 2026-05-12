package themeswitch

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/dom"
)

var (
	ClsTsBtn css.Class = "ts-btn"
)

// storageKey identifica la entrada de localStorage del componente.
const storageKey = "tinywasm-themeswitch"

// Theme representa el estado de tema del componente.
// ThemeAuto ("") = sin atributo data-theme → OS preference via @media.
// Los valores "dark" y "light" se escriben literalmente en data-theme.
type Theme string

const (
	ThemeAuto  Theme = ""      // no data-theme attribute → OS preference
	ThemeDark  Theme = "dark"
	ThemeLight Theme = "light"
)

// ThemeSwitch es un botón flotante que cicla entre los 3 modos de tema.
// Restaura automáticamente el tema guardado en localStorage al montarse.
//
//   ts := &themeswitch.ThemeSwitch{}
//   dom.Append("body", ts)
type ThemeSwitch struct {
	dom.Element
}

func (t *ThemeSwitch) Render() *dom.Element {
	current := Theme(dom.GetDocumentAttr("data-theme"))
	return dom.Button(icon(current)).
		Add(dom.Class(ClsTsBtn)).
		Attr("title", label(current)).
		Attr("aria-label", label(current)).
		On("click", t.onClick) // implementado por build tag
}

// cycle define el orden de los 3 estados. Switch (no map) — TinyGo.
func cycle(current Theme) Theme {
	switch current {
	case ThemeDark:
		return ThemeLight
	case ThemeLight:
		return ThemeAuto
	default: // ThemeAuto ("") o cualquier valor inesperado
		return ThemeDark
	}
}

// icon retorna el símbolo visible del botón. Switch (no map) — TinyGo.
func icon(theme Theme) string {
	switch theme {
	case ThemeDark:
		return "🌙"
	case ThemeLight:
		return "☀"
	default: // ThemeAuto ("")
		return "◑"
	}
}

// label retorna el nombre del modo actual (usado como tooltip y aria-label).
func label(theme Theme) string {
	switch theme {
	case ThemeDark:
		return "dark"
	case ThemeLight:
		return "light"
	default: // ThemeAuto ("")
		return "auto"
	}
}

func valid(t Theme) bool {
	return t == ThemeAuto || t == ThemeDark || t == ThemeLight
}
