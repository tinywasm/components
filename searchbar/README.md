# SearchBar

A one-control filter bar: a magnifier cap followed by a text field, both sized
by `--control-height` so the bar lines up with any other control in the
ecosystem. It holds no list and filters nothing itself — it reports each
keystroke and the host decides what that means.

## Usage

```go
import "github.com/tinywasm/components/searchbar"

bar := &searchbar.SearchBar{Placeholder: "Buscar..."}
bar.OnFilterChange(func(term string) { /* filter your list */ })
```

## API

### SearchBar Struct

- `Placeholder string`: Text shown in the field when it is empty. If left empty,
  the bar falls back to `"Search…"`.

## Filterable

`SearchBar` implements `widget.Filterable`. The bar reports the term and nothing
else, so any control that can produce a term — a calendar, a select — can
replace it in the same host slot without the host changing.

## Parts

| Class | Element | Role |
|---|---|---|
| `searchbar` | root | the bar as one rounded, clipped control |
| `searchbar__icon` | `<label>` | the square primary cap holding the magnifier |
| `searchbar__glyph` | `<svg>` | the magnifier itself, `currentColor` |
| `searchbar__input` | `<input type="search">` | the text field, the bar's body |

The root and the input both answer to `--control-height`, so cap and body can
never drift apart vertically.
