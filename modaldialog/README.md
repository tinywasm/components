# ModalDialog

Modal dialog overlay with backdrop and close button.

## Import
`"github.com/tinywasm/components/modaldialog"`

## Usage

```go
import (
    "github.com/tinywasm/components/modaldialog"
    . "github.com/tinywasm/html"
)

m := &modaldialog.ModalDialog{
    Title: "Confirm Action",
    Content: P("Are you sure?"),
}

// To open
m.Open()
```

## Properties
- `Title` (string): Modal title.
- `Content` (dom.Component): Modal body content.

## Methods
- `Open()`: Shows the modal.

---
[Back to Catalog](../docs/CATALOG.md)
