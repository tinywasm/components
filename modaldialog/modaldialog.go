package modaldialog

import (
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/widget"
)

// NameModalDialog is the widget name for modal dialog.
const NameModalDialog = widget.Name("modaldialog")

const (
	PartBackdrop = widget.Part("backdrop")
	PartPanel    = widget.Part("panel")
	PartHeader   = widget.Part("header")
	PartBody     = widget.Part("body")
	PartActions  = widget.Part("actions")
	PartClose    = widget.Part("close")
)

var (
	clsModal         = NameModalDialog.Root()
	clsModalBackdrop = NameModalDialog.Class(PartBackdrop)
	clsModalContent  = NameModalDialog.Class(PartPanel)
	clsModalHeader   = NameModalDialog.Class(PartHeader)
	clsModalBody     = NameModalDialog.Class(PartBody)
	clsModalClose    = NameModalDialog.Class(PartClose)
)

// ModalDialog represents a modal dialog component.
type ModalDialog struct {
	Element
	Title   string
	Content Component

	// HideClose drops the "×" from the header. Set it when the dialog's own
	// content already offers an explicit way out — a confirmation with
	// Cancel/Confirm buttons, where a third exit is one more thing to read and
	// says nothing Cancel does not. Clicking the backdrop still closes.
	HideClose bool

	visible *SignalBool
}

func (m *ModalDialog) WidgetName() widget.Name { return NameModalDialog }
func (m *ModalDialog) WidgetKind() widget.Kind { return widget.Dialog }

func (m *ModalDialog) Init(_ Ctx) {
	m.visible = NewBool(false)
}

func (m *ModalDialog) Render() *Element {
	header := Div().Set(clsModalHeader.AsAttr()).
		Child(H2().Text(m.Title))
	if !m.HideClose {
		header.Child(Button().
			Text("×").
			Set(clsModalClose.AsAttr()).
			On("click", func(e Event) { m.visible.Set(false) }))
	}

	modalContent := Div().Set(clsModalContent.AsAttr()).
		Child(header).
		Child(
			Div().Set(clsModalBody.AsAttr()).Child(m.Content),
		)

	modal := Div().Set(clsModal.AsAttr()).
		Attr("role", "dialog").
		Attr("aria-modal", "true").
		Child(Div().Set(clsModalBackdrop.AsAttr()).
			On("click", func(e Event) { m.visible.Set(false) })).
		Child(modalContent)
	return Show(m.visible, modal)
}

func (m *ModalDialog) Open() {
	m.visible.Set(true)
}

// Close hides the dialog programmatically.
func (m *ModalDialog) Close() {
	m.visible.Set(false)
}
