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
- `Href` (string): When set, renders as `<a href>` instead of `<button>` —
  no click handler, works before WASM loads and with JavaScript disabled.
  Use this for navigation (e.g. an OAuth login link). `OnClick` is ignored
  when `Href` is non-empty.
- `OnClick` (func(dom.Event)): Click handler. Ignored when `Href` is set.

## Navigation link example

```go
loginLink := &actionbutton.ActionButton{
    Text:    "Sign in with Google",
    Variant: "primary",
    Href:    "/oauth/google",
}
```

---
[Back to Catalog](../docs/CATALOG.md)
