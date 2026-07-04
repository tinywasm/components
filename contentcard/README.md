# ContentCard

Container with header, body, and footer sections.

## Import
`"github.com/tinywasm/components/contentcard"`

## Usage

```go
import (
    "github.com/tinywasm/components/contentcard"
    . "github.com/tinywasm/html"
)

c := &contentcard.ContentCard{
    Header: H3("Card Title"),
    Body:   P("Content goes here..."),
    Footer: Small("Last updated: today"),
}
```

## Properties
- `Header` (dom.Component): Content for header section.
- `Body` (dom.Component): Main content.
- `Footer` (dom.Component): Footer content.

---
[Back to Catalog](../docs/CATALOG.md)
