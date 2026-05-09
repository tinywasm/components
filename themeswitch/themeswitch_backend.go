//go:build !wasm

package themeswitch

import "github.com/tinywasm/dom"

// En SSR no hay localStorage ni clicks. Stubs no-op para compilación correcta.
func (t *ThemeSwitch) OnMount()          {}
func (t *ThemeSwitch) onClick(dom.Event) {}
