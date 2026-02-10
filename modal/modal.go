package modal

import (
	"github.com/tinywasm/components"
	"github.com/tinywasm/dom"
)

type Modal struct {
	components.BaseComponent
	Title   string
	Content dom.Component
	Visible bool
}

func (m *Modal) Render() dom.Node {
	class := "modal"
	if !m.Visible {
		class += " hidden"
	}

	// Create backdrop
	backdrop := dom.Div().
		Class("modal-backdrop")

	backdrop.OnClick(func(e dom.Event) {
		m.Close(e)
	})

	// Create close button
	closeBtn := dom.Button().
		Text("×").
		Class("modal-close")

	closeBtn.OnClick(func(e dom.Event) {
		m.Close(e)
	})

	modalContent := dom.Div().Class("modal-content").
		Append(
			dom.Div().Class("modal-header").
				Append(dom.Tag("h2").Text(m.Title)).
				Append(closeBtn),
		).
		Append(
			dom.Div().Class("modal-body").Append(m.Content),
		)

	return dom.Div().
		ID(m.ID()).
		Class(class).
		Append(backdrop).
		Append(modalContent).
		ToNode()
}

func (m *Modal) Close(e dom.Event) {
	m.Visible = false
	dom.Update(m)
}

func (m *Modal) Open() {
	m.Visible = true
	dom.Update(m)
}
