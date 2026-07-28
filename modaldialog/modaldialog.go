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
	visible *SignalBool
}

func (m *ModalDialog) WidgetName() widget.Name { return NameModalDialog }
func (m *ModalDialog) WidgetKind() widget.Kind { return widget.Dialog }

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
			Attr("role", "dialog").
			Attr("aria-modal", "true").
			Child(Div().Set(clsModalBackdrop.AsAttr()).
				On("click", func(e Event) { m.visible.Set(false) })).
			Child(modalContent)
	})
}

func (m *ModalDialog) Open() {
	m.visible.Set(true)
}

// Close hides the dialog programmatically.
func (m *ModalDialog) Close() {
	m.visible.Set(false)
}
