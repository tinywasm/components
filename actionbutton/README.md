# Button Component

Versatile button component with variant support.

## Import
`"github.com/tinywasm/components/button"`

## Usage

```go
btn := &button.Button{
    Text: "Save Changes",
    Variant: "primary",
    OnClick: func(e dom.Event) {
        fmt.Println("Saved!")
    },
}
```

## Properties
- `Text` (string): The button label.
- `Variant` (string): Visual style (`primary`, `secondary`, `danger`).
- `OnClick` (func(dom.Event)): Click handler.

## CSS Variables
- `--color-primary`
- `--color-secondary`
- `--color-error`

---
[Back to Catalog](../docs/CATALOG.md)
