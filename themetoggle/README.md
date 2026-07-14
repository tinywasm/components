# ThemeToggle

Signal-driven component to toggle the application theme (dark / light).

## Features

- 2 states: `dark` and `light`. Defaults to `light` when nothing is saved.
  There is no `auto`/OS-preference state — if an app wants to seed the theme
  from `prefers-color-scheme`, resolve that outside this component (e.g. by
  writing to `localStorage`/`data-theme` before mounting `ThemeToggle`).
- Automatic persistence in `localStorage` via `Init`.
- Fixed floating button.
- Signal-driven: updates icon and labels surgically without re-rendering the whole button.

## Usage

```go
import (
	"github.com/tinywasm/components/themetoggle"
	"github.com/tinywasm/dom"
)

func main() {
	ts := &themetoggle.ThemeToggle{}
	dom.Append("body", ts)
	select {}
}
```
