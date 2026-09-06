package fieldset

import (
	"webtyp.com/dom"
	"webtyp.com/widget"
)

// Fieldset is the skin's receiver type. The SSR collector instantiates it to call
// RenderCSS(); it carries no state and paints no markup of its own.
type Fieldset struct {
	dom.Element
}

// WidgetName returns the widget identity name.
func (f *Fieldset) WidgetName() widget.Name {
	return widget.NameField
}

// WidgetKind returns the widget ARIA kind.
func (f *Fieldset) WidgetKind() widget.Kind {
	return widget.Form
}
