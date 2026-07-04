# DataTable

Data table for structured information.

## Import
`"github.com/tinywasm/components/datatable"`

## Usage

```go
import "github.com/tinywasm/components/datatable"

t := &datatable.DataTable{
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
