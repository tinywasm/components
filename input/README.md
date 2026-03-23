# Input Component

Text input with label and validation support.

## Import
`"github.com/tinywasm/components/input"`

## Usage

```go
inp := &input.Input{
    Label: "Username",
    Placeholder: "Enter username",
    Type: "text",
    Required: true,
    OnInput: func(e dom.Event) {
        // Handle input
    },
}
```

## Properties
- `Label` (string): Label text.
- `Placeholder` (string): Placeholder text.
- `Type` (string): HTML input type (`text`, `email`, `password`, etc).
- `Value` (string): Initial value.
- `Required` (bool): Whether input is required.
- `OnInput` (func(dom.Event)): Input handler.

---
[Back to Catalog](../docs/CATALOG.md)
