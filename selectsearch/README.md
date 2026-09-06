# SelectSearch

Signal-driven searchable dropdown with static options, live filtering, and optional DB search callback.

## Features

- Signal-driven open/close state using `Show`.
- Live search filtering with `BindChildren` for surgical list updates.
- Badge descriptions for options.
- Callback for custom search (e.g., database).
- Two-way bound search input with `.Autofocus()`.
- Accessible via labels and semantic HTML.

## Usage

```go
import "webtyp.com/components/selectsearch"

ss := &selectsearch.SelectSearch{
    Placeholder: "Choose an option...",
    Options: []selectsearch.SsOption{
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
- `Options []SsOption`: Initial list of options.
- `OnSelect func(id, description string)`: Callback triggered when an option is selected.
- `OnSearch func(term string) []SsOption`: Callback triggered when all local options are filtered out.

### SsOption Struct

- `ID string`: Unique identifier for the option.
- `Label string`: Visible text for the option.
- `Description string`: Optional secondary text (badge).

## Filterable

`SelectSearch` implements `widget.Filterable`: picking an option calls the
registered sink with the option's `ID`. This is in addition to `OnSelect`
(which also gets `Description`) — a host that only needs the generic
narrowing contract (e.g. `webtyp/layout/crudview`'s `Filter` slot) can drop
a `*SelectSearch` in without any bespoke wiring, the same way it accepts a
`*searchbar.SearchBar` today.
