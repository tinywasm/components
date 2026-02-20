package input

import (
	"github.com/tinywasm/dom"
)

type Input struct {
	*dom.Element
	Label       string
	Placeholder string
	Value       string
	Type        string // text, email, password
	Required    bool
	OnInput     func(dom.Event)
}

func (i *Input) Render() *dom.Element {
	if i.Element == nil {
		i.Element = &dom.Element{}
	}
	container := dom.Div().Class("input-group")

	inputId := i.GetID() + "-field"
	if i.Label != "" {
		container.Add(
			dom.Label().
				Attr("for", inputId).
				Text(i.Label),
		)
	}

	inputType := i.Type
	if inputType == "" {
		inputType = "text"
	}

	input := dom.Input(inputType).
		ID(inputId).
		Attr("type", inputType).
		Attr("placeholder", i.Placeholder).
		Attr("value", i.Value)

	if i.Required {
		input.Required()
	}

	if i.OnInput != nil {
		input.On("input", i.OnInput)
	}

	container.Add(input.AsElement())
	return container
}
