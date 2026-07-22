package modaldialog

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

// ModalDialog represents a modal dialog component.
type ModalDialog struct {
	Element
	Title   string
	Content Component
	visible *SignalBool
}

func (m *ModalDialog) Init(_ Ctx) {
	m.visible = NewBool(false)
}

func (m *ModalDialog) Render() *Element {
	// Create close button
	closeBtn := Button().
		Text("×").
		Set(clsModalClose.AsAttr()).
		On("click", func(e Event) { m.visible.Set(false) })

	modalContent := Div().Set(clsModalContent.AsAttr()).
		Child(
			Div().Set(clsModalHeader.AsAttr()).
				Child(H2().Text(m.Title)).
				Child(closeBtn),
		).
		Child(
			Div().Set(clsModalBody.AsAttr()).Child(m.Content),
		)

	return Show(m.visible, func() *Element {
		return Div().Set(clsModal.AsAttr()).
			Child(Div().Set(clsModalBackdrop.AsAttr()).
				On("click", func(e Event) { m.visible.Set(false) })).
			Child(modalContent)
	})
}

func (m *ModalDialog) Open() {
	m.visible.Set(true)
}

// Close hides the dialog programmatically — for a host that needs to dismiss
// it after an action completes (e.g. a confirm button), not just via the
// built-in backdrop/× close.
func (m *ModalDialog) Close() {
	m.visible.Set(false)
}
