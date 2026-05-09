//go:build wasm

package main

import (
	"github.com/tinywasm/components/themeswitch"
	. "github.com/tinywasm/dom"
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
	ts := &themeswitch.ThemeSwitch{}
	Render("app", &App{})
	Append("body", ts)
	select {}
}
