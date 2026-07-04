# ActionButton

Versatile button component with variant support.

## Import
`"github.com/tinywasm/components/actionbutton"`

## Usage

```go
import "github.com/tinywasm/components/actionbutton"

btn := &actionbutton.ActionButton{
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

---
[Back to Catalog](../docs/CATALOG.md)
