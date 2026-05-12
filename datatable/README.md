# Table Component

Data table.

## Import
`"github.com/tinywasm/components/table"`

## Usage

```go
t := &table.Table{
    Headers: []string{"ID", "Name", "Role"},
    Rows: [][]string{
        {"1", "Alice", "Admin"},
        {"2", "Bob", "User"},
    },
}
```

## Properties
- `Headers` ([]string): Column headers.
- `Rows` ([][]string): Data rows.

---
[Back to Catalog](../docs/CATALOG.md)
