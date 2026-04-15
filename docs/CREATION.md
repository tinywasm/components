# Component Creation Guide (TinyWasm Components)

This guide establishes the standard for creating reusable components in `tinywasm/components`.

## File Structure

Each component must reside in its own folder within `tinywasm/components` and consist of at least 2 files (and optionally 3 for WASM logic):

```
tinywasm/components/
└── mycomponent/
    ├── mycomponent.go   # Shared structure, RenderHTML
    ├── ssr.go           # Backend: CSS, initial JS (build tag !wasm)
    └── front.go         # Frontend: WASM logic (build tag wasm)
```

## CSS guidelines

- All colors must use `--color-*` CSS custom properties from `tinywasm/dom` theme.
- Never hardcode hex values for colors. Always provide a fallback:
  `var(--color-secondary, #00ADD8)`.
- Spacing must use `--mag-pri`, `--mag-sec`, `--mag-cua` variables.
- CSS class names must be prefixed with the component name to avoid collisions.
- CSS lives in `<component>.css`, embedded in `ssr.go` via `//go:embed`.
- Do NOT create or embed form-related CSS — use `tinywasm/form`.

## 1. Main File (`mycomponent.go`)

Contains the struct definition and common HTML rendering logic.

```go
package mycomponent

import (
    "github.com/tinywasm/dom"
)

type MyComponent struct {
    *dom.Element
    Title string
}

func (c *MyComponent) Render() *dom.Element {
    if c.Element == nil {
        c.Element = &dom.Element{}
    }
    return dom.Div().
        Class("my-class").
        Text(c.Title)
}
```

## 2. Backend File (`ssr.go`)

**CRITICAL:** This file must have the `//go:build !wasm` tag. Here, CSS and any logic that should not reach the WASM client are defined to keep the binary lightweight.

By using separate build tags, we can leverage `go:embed` to include CSS from a separate file (e.g., `mycomponent.css`). This is highly recommended as it makes it much easier for developers to review and edit CSS or JS (if any server-side JS is needed) compared to inline strings.

```go
//go:build !wasm

package mycomponent

import _ "embed"

//go:embed mycomponent.css
var css string

// RenderCSS returns the specific CSS for this component.
func (c *MyComponent) RenderCSS() string {
    return css
}
```

## 3. Icon Management (IconSvgProvider)

To register SVG icons in the site's global sprite (accessible via `<use href="#id">`), the component must implement the `IconSvgProvider` interface.

> [!CRITICAL]
> **MANDATORY:** The `IconSvg()` method MUST be in `ssr.go` (with build tag `//go:build !wasm`).
> SVG strings are dead code on the WASM client and unnecessarily increase the binary size.
> Never define icons in the component's main file.

```go
func (c *MyComponent) IconSvg() map[string]string {
    return map[string]string{
        // DO NOT include the wrapping <svg> tag. The system handles it automatically.
        // Only include the internal content (paths, circles, etc).
        // A default viewBox of "0 0 16 16" is assumed.
        // If you need another size, the system will automatically detect a viewBox="..." attribute in your string.
        "my-icon-id": `<path d="..." />`, 
    }
}
```

The map key (`"my-icon-id"`) will be the ID in the sprite.

## 4. Frontend File (`front.go`) (Optional)

If the component requires client-side interactivity (event listeners, subsequent DOM manipulation), use this file with the `//go:build wasm` tag.

```go
//go:build wasm

package mycomponent

import (
    "github.com/tinywasm/dom"
    "github.com/tinywasm/fmt"
)

// OnMount runs automatically when the component's HTML has been injected into the DOM.
// It is the correct place to add event listeners using tinywasm/dom.
func (c *MyComponent) OnMount() {
    // 1. Get reference to the element (using the ID provided by Element)
    id := c.GetID()
    dom.Log("Component mounted with ID:", id)
}
```

### Rules for Interactivity (`tinywasm/dom`)

1.  **Always use `tinywasm/dom`**: Avoid direct `syscall/js` to maintain size optimization.
2.  **Use `OnMount`**: Do not attempt to find elements in `RenderHTML` (if using string based rendering), as they do not yet exist in the DOM. With the fluent API, `OnMount` is handled internally by `dom` via `Render`.
3.  **Unique IDs**: Ensure IDs generated in `Render` are unique (using `c.ID()` as a prefix) to be able to find them in `OnMount`.

## Integration
When using the component in a module, `tinywasm/site` will handle collecting the CSS (since `site` runs on the server) and serving it, while the WASM client will only receive the necessary structure and logic.
