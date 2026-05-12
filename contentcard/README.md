# Card Component

Container with header, body, and footer sections.

## Import
`"github.com/tinywasm/components/card"`

## Usage

```go
c := &card.Card{
    Header: dom.Tag("h3").Text("Card Title").ToComponent(),
    Body:   dom.Tag("p").Text("Content goes here...").ToComponent(),
    Footer: dom.Tag("small").Text("Last updated: today").ToComponent(),
}
```

## Properties
- `Header` (dom.Component): Content for header section.
- `Body` (dom.Component): Main content.
- `Footer` (dom.Component): Footer content.

---
[Back to Catalog](../docs/CATALOG.md)
