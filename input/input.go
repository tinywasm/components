package input

import (
	"github.com/tinywasm/components"
	dom "github.com/tinywasm/components/internal/dom"
)

type Input struct {
	components.BaseComponent
	Label       string
	Placeholder string
	Value       string
	Type        string // text, email, password
	Required    bool
	OnInput     func(dom.Event)
}

func (i *Input) Render() dom.Node {
	container := dom.Div().Class("input-group")

	if i.Label != "" {
		container.Append(
			dom.Tag("label").
				Attr("for", i.ID()).
				Text(i.Label),
		)
	}

	inputType := i.Type
	if inputType == "" {
		inputType = "text"
	}

	input := dom.Tag("input").
		ID(i.ID()).
		Attr("type", inputType).
		Attr("placeholder", i.Placeholder).
		Attr("value", i.Value)

	if i.Required {
		input.Attr("required", "required")
	}

	if i.OnInput != nil {
		input.OnInput(i.OnInput)
	}

	container.Append(input)
	return container.ToNode()
}
