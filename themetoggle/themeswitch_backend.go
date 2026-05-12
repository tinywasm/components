//go:build !wasm

package themetoggle

import "github.com/tinywasm/dom"

// En SSR no hay localStorage ni clicks. Stubs no-op para compilación correcta.
func (t *ThemeToggle) OnMount()          {}
func (t *ThemeToggle) onClick(dom.Event) {}
