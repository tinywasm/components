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
	// This is just a showcase client
	renderShowcase()
	select {}
}

func renderShowcase() {
	// Button showcase
	btn := &button.Button{
		Text:    "Click Me",
		Variant: "primary",
		OnClick: func(e dom.Event) {
			dom.Log("Button clicked!")
		},
	}
	dom.Append("app", btn)

	// Card showcase
	c := &card.Card{
		Header: dom.Tag("h3").Text("Card Header").ToComponent(),
		Body:   dom.Tag("p").Text("Card Body Content").ToComponent(),
		Footer: dom.Tag("div").Text("Card Footer").ToComponent(),
	}
	dom.Append("app", c)

	// Input showcase
	inp := &input.Input{
		Label:       "Username",
		Placeholder: "Enter username",
		OnInput: func(e dom.Event) {
			dom.Log("Input changed")
		},
	}
	dom.Append("app", inp)

	// Nav showcase
	n := &nav.Nav{
		Items: []nav.NavItem{
			{Label: "Home", Route: "home"},
			{Label: "Profile", Route: "profile"},
		},
	}
	dom.Append("app", n)

	// Modal showcase
	modalContent := dom.Tag("p").Text("This is a modal content").ToComponent()
	m := &modal.Modal{
		Title:   "Example Modal",
		Content: modalContent,
		Visible: false,
	}
	dom.Append("app", m)

	openBtn := &button.Button{
		Text:    "Open Modal",
		Variant: "secondary",
		OnClick: func(e dom.Event) {
			m.Open()
		},
	}
	dom.Append("app", openBtn)

	// Table showcase
	t := &table.Table{
		Headers: []string{"ID", "Name", "Role"},
		Rows: [][]string{
			{"1", "Alice", "Admin"},
			{"2", "Bob", "User"},
		},
	}
	dom.Append("app", t)

	// Form showcase
	f := &form.Form{
		Fields: []dom.Component{
			&input.Input{Label: "Email", Type: "email"},
			&input.Input{Label: "Password", Type: "password"},
			&button.Button{Text: "Login", Variant: "primary", OnClick: func(e dom.Event) {
				dom.Log("Form submitted")
			}},
		},
	}
	dom.Append("app", f)
}
