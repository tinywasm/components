package form

import "github.com/tinywasm/fmt"

func New(id string) *Form {
	return &Form{
		ID:           id,
		Name:         id,
		Autocomplete: false,
		Fieldsets:    []*Fieldset{},
	}
}

func (b *Form) SetAutocomplete(enable bool) *Form {
	b.Autocomplete = enable
	return b
}

func (b *Form) AddFieldset(legend string, tabIndex int, fields ...Field) *Form {
	b.Fieldsets = append(b.Fieldsets, &Fieldset{
		Legend:   legend,
		TabIndex: tabIndex,
		Fields:   fields,
		Class:    "max-height-100", // default class as seen in examples
	})
	return b
}

// Helper to quickly add a single input in a fieldset (common pattern)
func (b *Form) AddInput(name, label string, typ InputType, required bool) *Form {
	// Auto-calculate tabIndex based on count
	idx := len(b.Fieldsets)

	f := &InputField{
		ID:       "form." + b.ID + "." + name, // convention from example
		Name:     name,
		Label:    label,
		Type:     typ,
		Required: required,
	}

	return b.AddFieldset(label, idx, f)
}

func (b *Form) AddRadioGroup(legend string, name string, options ...Option) *Form {
	idx := len(b.Fieldsets)
	group := &RadioGroupField{
		Name:    name,
		Options: options,
	}
	fieldset := &Fieldset{
		Legend:   legend,
		TabIndex: idx,
		Fields:   []Field{group},
	}
	// override class for radio groups if needed, usually they don't block height?
	// keeping simple for now.
	b.Fieldsets = append(b.Fieldsets, fieldset)
	return b
}

func (b *Form) Render() string {
	inner := ""
	for _, fs := range b.Fieldsets {
		fieldsHtml := ""
		for _, f := range fs.Fields {
			fieldsHtml += f.Render()
		}

		legendHtml := ""
		if fs.Legend != "" {
			// For input fields, the label is usually inside the legend or next to it.
			// In the example: <legend><label for="...">SKU</label></legend> ... <input ...>
			// But the InputField.Render() generates <label><input> or similar?

			// Let's look at the example target:
			// <legend><label for="form.catalog.sku">SKU</label></legend>
			// <input type="text" id="form.catalog.sku" name="catalog_sku" required>

			// My Field implementations generate the label+input pair.
			// But the wrapper fieldset also has a legend.

			// If the field is an InputField, we might want to put its Label in the Legend?
			// The example pattern is specific:
			// Legend contains the label for the input.

			// Simplified approach for generic form:
			// Just render legend as text.
			legendHtml = fmt.Sprintf("<legend>%s</legend>", fs.Legend)
		}

		cls := ""
		if fs.Class != "" {
			cls = fmt.Sprintf(` class="%s"`, fs.Class)
		}

		inner += fmt.Sprintf(`<fieldset tabindex="%d"%s>%s%s</fieldset>`, fs.TabIndex, cls, legendHtml, fieldsHtml)
	}

	auto := "off"
	if b.Autocomplete {
		auto = "on"
	}

	return fmt.Sprintf(
		`<form class="form-distributed-fields" id="%s" name="%s" autocomplete="%s">%s</form>`,
		b.ID, b.Name, auto, inner,
	)
}

// RenderHTML implements site.Component interface
func (b *Form) RenderHTML() string {
	return b.Render()
}

// Helpers for options
func Radio(value, label string) Option {
	return Option{Value: value, Label: label}
}
