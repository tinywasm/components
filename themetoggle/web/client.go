//go:build wasm

package main

import (
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

type App struct {
	Element
}

func (a *App) Render() *Element {
	return Div(
		H1("ThemeSwitch Demo"),
		P("Use the button in the top right corner to cycle through themes."),
	)
}

func main() {
	// ThemeToggle is initialized internally
	Render("app", &App{})
	select {}
}
