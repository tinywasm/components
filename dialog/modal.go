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
	visible *SignalBool
}

func (m *DialogWidget) Init(_ Ctx) {
	m.visible = NewBool(false)
}

func (m *DialogWidget) Render() *Element {
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

func (m *DialogWidget) Open() {
	m.visible.Set(true)
}
