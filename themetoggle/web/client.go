//go:build wasm

package main

import (
	. "webtyp.com/dom"
	. "webtyp.com/html"
)

type App struct {
	Element
}

func (a *App) Render() *Element {
	return Div().Child(
		H1().Text("ThemeSwitch Demo"),
		P().Text("Use the button in the top right corner to cycle through themes."),
	)
}

func main() {
	// ThemeToggle is initialized internally
	Render("app", &App{})
	select {}
}
