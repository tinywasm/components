package form

import (
	"github.com/tinywasm/components"
	dom "github.com/tinywasm/components/internal/dom"
)

type Form struct {
	components.BaseComponent
	Fields   []dom.Component
	OnSubmit func(dom.Event)
}

func (f *Form) Render() dom.Node {
	form := dom.Tag("form").
		ID(f.ID()).
		Class("form")

	if f.OnSubmit != nil {
		form.OnSubmit(f.OnSubmit)
	}

	for _, field := range f.Fields {
		form.Append(field)
	}

	return form.ToNode()
}
