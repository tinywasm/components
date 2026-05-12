package modal

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/dom"
)

var (
	ClsModal         css.Class = "modal"
	ClsModalHidden   css.Class = "hidden"
	ClsModalBackdrop css.Class = "modal-backdrop"
	ClsModalContent  css.Class = "modal-content"
	ClsModalHeader   css.Class = "modal-header"
	ClsModalClose    css.Class = "modal-close"
	ClsModalBody     css.Class = "modal-body"
)

type Modal struct {
	*dom.Element
	Title   string
	Content dom.Component
	Visible bool
}

func (m *Modal) Render() *dom.Element {
	if m.Element == nil {
		m.Element = &dom.Element{}
	}

	modal := dom.Div().Add(dom.Class(ClsModal))
	if !m.Visible {
		modal.Add(dom.Class(ClsModalHidden))
	}

	// Create backdrop
	backdrop := dom.Div().
		Add(dom.Class(ClsModalBackdrop))

	backdrop.On("click", func(e dom.Event) {
		m.Close(e)
	})

	// Create close button
	closeBtn := dom.Button().
		Text("×").
		Add(dom.Class(ClsModalClose))

	closeBtn.On("click", func(e dom.Event) {
		m.Close(e)
	})

	modalContent := dom.Div().Add(dom.Class(ClsModalContent)).
		Add(
			dom.Div().Add(dom.Class(ClsModalHeader)).
				Add(dom.H2().Text(m.Title)).
				Add(closeBtn),
		).
		Add(
			dom.Div().Add(dom.Class(ClsModalBody)).Add(m.Content),
		)

	return modal.
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
