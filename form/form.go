package form

import (
	"github.com/tinywasm/dom"
)

type Form struct {
	*dom.Element
	Fields   []dom.Component
	OnSubmit func(dom.Event)
}

func (f *Form) Render() *dom.Element {
	if f.Element == nil {
		f.Element = &dom.Element{}
	}

	form := dom.Form().
		Class("form")

	if f.OnSubmit != nil {
		form.OnSubmit(f.OnSubmit)
	}

	for _, field := range f.Fields {
		form.Add(field)
	}

	return form.AsElement()
}
