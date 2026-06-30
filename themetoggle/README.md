# ThemeToggle

Signal-driven component to toggle the application theme (light / dark / auto).

## Features

- 3 states: `auto` (OS preference), `dark`, and `light`.
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
