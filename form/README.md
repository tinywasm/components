# Form Component

Wrapper for form fields.

## Import
`"github.com/tinywasm/components/form"`

## Usage

```go
f := &form.Form{
    Fields: []dom.Component{
        &input.Input{Label: "Email"},
        &button.Button{Text: "Submit"},
    },
    OnSubmit: func(e dom.Event) {
        e.PreventDefault()
        // Handle submit
    },
}
```

## Properties
- `Fields` ([]dom.Component): List of form components.
- `OnSubmit` (func(dom.Event)): Submit handler.

---
[Back to Catalog](../docs/CATALOG.md)
