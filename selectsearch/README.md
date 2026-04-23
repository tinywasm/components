# SelectSearch

Searchable dropdown with static options, live filtering, and optional DB search callback.

## Features

- CSS-first toggle (no JS needed for open/close)
- Live search filtering
- Badge descriptions for options
- Callback for custom search (e.g., database)
- Accessible via labels and semantic HTML

## Usage

```go
import "github.com/tinywasm/components/selectsearch"

ss := &selectsearch.SelectSearch{
    Placeholder: "Choose an option...",
    Options: []selectsearch.Option{
        {ID: "1", Label: "Option 1", Description: "First option"},
        {ID: "2", Label: "Option 2", Description: "Second option"},
    },
    OnSelect: func(id, description string) {
        fmt.Printf("Selected: %s\n", id)
    },
}
```

## API

### SelectSearch Struct

- `Placeholder string`: Text shown when no option is selected.
- `Options []Option`: Initial list of options.
- `OnSelect func(id, description string)`: Callback triggered when an option is selected.
- `OnSearch func(term string) []Option`: Callback triggered when all local options are filtered out.

### Option Struct

- `ID string`: Unique identifier for the option.
- `Label string`: Visible text for the option.
- `Description string`: Optional secondary text (badge).
