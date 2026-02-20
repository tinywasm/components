package modal

import (
	"github.com/tinywasm/dom"
)

type Modal struct {
	*dom.Element
	Title   string
	Content dom.Component
	Visible bool
}

func (m *Modal) Render() *dom.Element {
	class := "modal"
	if !m.Visible {
		class += " hidden"
	}

	// Create backdrop
	backdrop := dom.Div().
		Class("modal-backdrop")

	backdrop.On("click", func(e dom.Event) {
		m.Close(e)
	})

	// Create close button
	closeBtn := dom.Button().
		Text("×").
		Class("modal-close")

	closeBtn.On("click", func(e dom.Event) {
		m.Close(e)
	})

	modalContent := dom.Div().Class("modal-content").
		Add(
			dom.Div().Class("modal-header").
				Add(dom.H2().Text(m.Title)).
				Add(closeBtn),
		).
		Add(
			dom.Div().Class("modal-body").Add(m.Content),
		)

	if m.Element == nil {
		m.Element = &dom.Element{}
	}

	return dom.Div().
		Class(class).
		Add(backdrop).
		Add(modalContent)
}

func (m *Modal) Close(e dom.Event) {
	m.Visible = false
	dom.Update(m)
}

func (m *Modal) Open() {
	m.Visible = true
	dom.Update(m)
}
