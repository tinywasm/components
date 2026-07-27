package fieldset

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/widget"
)

// NameFieldset is the widget name for fieldset.
const NameFieldset = widget.Name("fieldset")

// Parts of the fieldset widget.
const (
	PartLabel = widget.Part("label")
	PartError = widget.Part("error")
	PartInput = widget.Part("input")
)

// Fieldset is the skin's receiver type. The SSR collector instantiates it to call
// RenderCSS(); it carries no state and paints no markup of its own.
type Fieldset struct {
	dom.Element
}

// WidgetName returns the widget identity name.
func (f *Fieldset) WidgetName() widget.Name {
	return NameFieldset
}

// WidgetKind returns the widget ARIA kind.
func (f *Fieldset) WidgetKind() widget.Kind {
	return widget.Form
}
