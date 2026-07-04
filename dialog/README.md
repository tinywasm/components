# DialogWidget

Modal dialog overlay with backdrop and close button.

## Import
`"github.com/tinywasm/components/dialog"`

## Usage

```go
import (
    "github.com/tinywasm/components/dialog"
    . "github.com/tinywasm/html"
)

m := &dialog.DialogWidget{
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
