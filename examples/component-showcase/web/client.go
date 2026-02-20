//go:build wasm

package main

import (
	"github.com/tinywasm/components/button"
	"github.com/tinywasm/components/card"
	"github.com/tinywasm/components/form"
	"github.com/tinywasm/components/input"
	"github.com/tinywasm/components/modal"
	"github.com/tinywasm/components/nav"
	"github.com/tinywasm/components/table"
	"github.com/tinywasm/dom"
)

func main() {
	// Initialize the showcase app
	app := &ShowcaseApp{
		Element: &dom.Element{},
	}
	app.Init()

	// 1. Inject CSS using dom.Style
	css := `
		body { font-family: sans-serif; margin: 0; background: #fdfdfd; color: #333; }
		.container { padding: 2rem; max-width: 900px; margin: 0 auto; }
		section { margin-bottom: 2rem; border-bottom: 1px solid #eee; padding-bottom: 2rem; }
		h1, h2 { color: #2c3e50; }
		.btn { cursor: pointer; padding: 0.5rem 1rem; border-radius: 4px; border: none; margin-right: 0.5rem; }
		.btn-primary { background: #3498db; color: white; }
		.btn-secondary { background: #95a5a6; color: white; }
		.card { border: 1px solid #ddd; border-radius: 4px; padding: 1rem; background: white; }
		.card-header { font-weight: bold; border-bottom: 1px solid #eee; margin-bottom: 0.5rem; padding-bottom: 0.5rem; }
		.card-footer { border-top: 1px solid #eee; margin-top: 0.5rem; padding-top: 0.5rem; color: #777; font-size: 0.9em; }
		.nav { background: #2c3e50; padding: 1rem; color: white; border-radius: 4px; }
		.nav-list { list-style: none; padding: 0; margin: 0; display: flex; gap: 1rem; }
		.nav-item a { color: white; text-decoration: none; }
		.nav-item.active a { font-weight: bold; text-decoration: underline; }
		.modal { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.5); display: flex; justify-content: center; align-items: center; }
		.modal.hidden { display: none; }
		.modal-content { background: white; padding: 2rem; border-radius: 8px; min-width: 300px; }
		.input-group { margin-bottom: 1rem; }
		.input-group label { display: block; margin-bottom: 0.25rem; }
		.input-group input { width: 100%; padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
		table { width: 100%; border-collapse: collapse; }
		th, td { text-align: left; padding: 0.5rem; border-bottom: 1px solid #eee; }
	`
	dom.Append("head", dom.Style(css))

	// 2. Render Main Container
	dom.Render("body", app)

	select {}
}

type ShowcaseApp struct {
	*dom.Element
	modal *modal.Modal
}

func (a *ShowcaseApp) Init() {
	// Initialize modal state
	a.modal = &modal.Modal{
		Title:   "Example Modal",
		Content: dom.P("This is a modal content"),
		Visible: false,
	}
}

func (a *ShowcaseApp) Render() *dom.Element {
	return dom.Div(
		dom.H1("TinyWASM Component Showcase"),

		// Nav Showcase
		dom.Section(
			dom.H2("Navigation"),
			&nav.Nav{
				Items: []nav.NavItem{
					{Label: "Home", Route: "home"},
					{Label: "Components", Route: "components"},
					{Label: "About", Route: "about"},
				},
			},
		),

		// Button Showcase
		dom.Section(
			dom.H2("Buttons"),
			&button.Button{
				Text:    "Primary Button",
				Variant: "primary",
				OnClick: func(e dom.Event) { dom.Log("Primary clicked!") },
			},
			&button.Button{
				Text:    "Secondary Button",
				Variant: "secondary",
				OnClick: func(e dom.Event) { dom.Log("Secondary clicked!") },
			},
		),

		// Card Showcase
		dom.Section(
			dom.H2("Cards"),
			&card.Card{
				Header: dom.Div().Text("Card Header"),
				Body:   dom.P("This is a simple card component."),
				Footer: dom.Div().Text("Footer content"),
			},
		),

		// Input & Form Showcase
		dom.Section(
			dom.H2("Forms & Inputs"),
			&form.Form{
				Fields: []dom.Component{
					&input.Input{Label: "Username", Placeholder: "jdoe", Type: "text"},
					&input.Input{Label: "Password", Placeholder: "secret", Type: "password"},
					&button.Button{
						Text:    "Login",
						Variant: "primary",
						OnClick: func(e dom.Event) {
							e.PreventDefault()
							dom.Log("Login attempt")
						},
					},
				},
			},
		),

		// Table Showcase
		dom.Section(
			dom.H2("Tables"),
			&table.Table{
				Headers: []string{"ID", "Product", "Price"},
				Rows: [][]string{
					{"1", "Laptop", "$999"},
					{"2", "Mouse", "$25"},
					{"3", "Monitor", "$199"},
				},
			},
		),

		// Modal Showcase
		dom.Section(
			dom.H2("Modals"),
			&button.Button{
				Text:    "Open Modal",
				Variant: "primary",
				OnClick: func(e dom.Event) {
					a.modal.Open()
				},
			},
			a.modal,
		),
	).Class("container")
}
