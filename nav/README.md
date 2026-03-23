# Nav Component

Navigation menu with support for icons.

## Import
`"github.com/tinywasm/components/nav"`

## Usage

```go
n := &nav.Nav{
    Items: []nav.NavItem{
        {Label: "Dashboard", Route: "dashboard", Icon: "icon-home"},
        {Label: "Settings", Route: "settings", Icon: "icon-cog"},
    },
}
```

## Properties
- `Items` ([]NavItem): List of navigation items.
    - `Label` (string): Display text.
    - `Route` (string): Target route/hash.
    - `Icon` (string): SVG sprite ID for icon (optional).

---
[Back to Catalog](../docs/CATALOG.md)
