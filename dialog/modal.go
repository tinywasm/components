package dialog

import (
	. "github.com/tinywasm/css"
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

var (
	clsModal         Class = "modal"
	clsModalHidden   Class = "hidden"
	clsModalBackdrop Class = "modal-backdrop"
	clsModalContent  Class = "modal-content"
	clsModalHeader   Class = "modal-header"
	clsModalClose    Class = "modal-close"
	clsModalBody     Class = "modal-body"
)

// DialogWidget represents a modal dialog component.
type DialogWidget struct {
	Element
	Title   string
	Content Component
	Visible bool
}

func (m *DialogWidget) Render() *Element {
	modal := Div().Add(clsModal.AsAttr())
	if !m.Visible {
		modal.Add(clsModalHidden.AsAttr())
	}

	// Create backdrop
	backdrop := Div().
		Add(clsModalBackdrop.AsAttr())

	backdrop.On("click", func(e Event) {
		m.Close(e)
	})

	// Create close button
	closeBtn := Button().
		Text("×").
		Add(clsModalClose.AsAttr())

	closeBtn.On("click", func(e Event) {
		m.Close(e)
	})

	modalContent := Div().Add(clsModalContent.AsAttr()).
		Add(
			Div().Add(clsModalHeader.AsAttr()).
				Add(H2().Text(m.Title)).
				Add(closeBtn),
		).
		Add(
			Div().Add(clsModalBody.AsAttr()).Add(m.Content),
		)

	return modal.
		Add(backdrop).
		Add(modalContent)
}

func (m *DialogWidget) Close(e Event) {
	m.Visible = false
	Update(m)
}

func (m *DialogWidget) Open() {
	m.Visible = true
	Update(m)
}
