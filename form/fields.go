package form

import "github.com/tinywasm/fmt"

// InputField represents a standard input field
type InputField struct {
	ID       string
	Name     string
	Label    string
	Type     InputType
	Required bool
	Value    string
}

func (f *InputField) Render() string {
	req := ""
	if f.Required {
		req = " required"
	}
	val := ""
	if f.Value != "" {
		val = fmt.Sprintf(" value=\"%s\"", f.Value)
	}

	idAttr := ""
	if f.ID != "" {
		idAttr = fmt.Sprintf(" id=\"%s\"", f.ID)
	}

	return fmt.Sprintf(
		`<label for="%s">%s</label><input type="%s"%s name="%s"%s%s>`,
		f.Name, f.Label, f.Type, idAttr, f.Name, val, req,
	)
}

// RadioGroupField represents a group of radio buttons
type RadioGroupField struct {
	Name    string
	Options []Option
}

type Option struct {
	Value string
	Label string
}

func (f *RadioGroupField) Render() string {
	out := ""
	for _, opt := range f.Options {
		out += fmt.Sprintf(
			`<label class="block-label"><input type="radio" name="%s" value="%s"><span>%s</span></label>`,
			f.Name, opt.Value, opt.Label,
		)
	}
	return out
}
