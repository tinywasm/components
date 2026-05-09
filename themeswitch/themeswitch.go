package themeswitch

import "github.com/tinywasm/dom"

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
	return dom.Button(label(current)).
		Class("ts-btn").
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

// label retorna el texto visible del botón. Switch (no map) — TinyGo.
func label(theme Theme) string {
	switch theme {
	case ThemeDark:
		return "🌙 dark"
	case ThemeLight:
		return "☀ light"
	default: // ThemeAuto ("")
		return "☀/🌙 auto"
	}
}

func valid(t Theme) bool {
	return t == ThemeAuto || t == ThemeDark || t == ThemeLight
}
