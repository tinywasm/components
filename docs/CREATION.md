# Component Creation Guide (TinyWasm Components)

This guide establishes the standard for creating reusable components in `tinywasm/components`.

## File Structure

Each component must reside in its own folder within `tinywasm/components` and consist of at least 2 files:

```
tinywasm/components/
└── mycomponent/
    ├── mycomponent.go        # Shared struct, Render(), OnMount()
    ├── mycomponent.css       # Component-scoped styles
    ├── mycomponent_test.go   # Tests
    └── ssr.go                # Backend only: CSS embed, IconSvg() — build tag !wasm
```

> There is NO `front.go`. WASM interactivity lives in `mycomponent.go` via `OnMount()`.
> The build system separates concerns via build tags on `ssr.go`, not by splitting files.

## CSS Guidelines

- All colors MUST use `--color-*` CSS custom properties from `tinywasm/dom` theme.
- Never hardcode hex values. Always provide a fallback:
  `var(--color-secondary, #00ADD8)`.
- Spacing MUST use `--mag-pri`, `--mag-sec`, `--mag-cua` variables.
- CSS class names MUST be prefixed with the component name to avoid collisions: `mycomponent-*`.
- CSS lives in `<component>.css`, embedded in `ssr.go` via `//go:embed`.
- Do NOT create or embed form-related CSS — use `tinywasm/form`.

## 1. Main File (`mycomponent.go`)

Contains the struct definition, `Render()` (shared SSR + WASM), and `OnMount()` (WASM only — no build tag needed; dead code eliminated by TinyGo).

### Embedding Rule — CRITICAL for TinyGo/WASM

**Always embed `dom.Element` as a VALUE, never as a pointer.**

```go
// ✅ CORRECT — value embed, zero GC overhead, no nil risk
type MyComponent struct {
    dom.Element
    Title string
}

// ❌ WRONG — pointer embed causes 2 heap allocations, requires nil-guard,
//            risks nil panic in WASM, wastes GC cycles in TinyGo
type MyComponent struct {
    *dom.Element
}
```

**Why value embedding:**
- TinyGo's GC is conservative and simple — fewer heap objects = fewer GC pauses
- One allocation instead of two (struct + Element separately)
- Better cache locality — Element fields are contiguous with the struct
- No nil-guard boilerplate, no nil panic risk in production WASM

### Example

```go
package mycomponent

import "github.com/tinywasm/dom"

type MyComponent struct {
    dom.Element  // value embed — never pointer
    Title string
}

func (c *MyComponent) Render() *dom.Element {
    return dom.Div().
        Class("mycomponent").
        Text(c.Title)
}

// OnMount wires events after the component is injected into the DOM.
// Called automatically by tinywasm/dom — no build tag needed.
// TinyGo eliminates this as dead code on SSR builds.
func (c *MyComponent) OnMount() {
    id := c.GetID()
    if el, ok := dom.Get(id); ok {
        el.On("click", func(e dom.Event) {
            // handle click
        })
    }
}
```

## 2. Backend File (`ssr.go`)

**CRITICAL:** This file MUST have the `//go:build !wasm` build tag.
SVG strings and embedded CSS are dead weight in the WASM binary — keep them here.

```go
//go:build !wasm

package mycomponent

import _ "embed"

//go:embed mycomponent.css
var css string

func (c *MyComponent) RenderCSS() string {
    return css
}
```

## 3. Icon Management (`IconSvgProvider`)

To register SVG icons in the site's global sprite (accessible via `<use href="#id">`), implement `IconSvgProvider` in `ssr.go`.

> **MANDATORY:** `IconSvg()` MUST be in `ssr.go` (`//go:build !wasm`).
> SVG strings are dead code on WASM — never define icons in the main file.

```go
// In ssr.go (already has !wasm build tag)
func (c *MyComponent) IconSvg() map[string]string {
    return map[string]string{
        // Do NOT include the wrapping <svg> tag — the system adds it.
        // Only internal content: paths, circles, etc.
        // Default viewBox is "0 0 16 16".
        // If you need a different size, include viewBox="..." in your string.
        "my-icon-id": `<path d="..." />`,
    }
}
```

The map key (`"my-icon-id"`) becomes the sprite ID, referenced as `<use href="#my-icon-id">`.

## 4. Tests (`mycomponent_test.go`)

Tests run on the backend (no build tag needed — default is `!wasm` when not targeting WASM).
Call `Render().RenderHTML()` and assert on the HTML string.

```go
package mycomponent

import (
    "strings"
    "testing"
)

func TestMyComponent_Render(t *testing.T) {
    c := &MyComponent{Title: "Hello"}
    html := c.Render().RenderHTML()

    if !strings.Contains(html, "mycomponent") {
        t.Error("expected mycomponent class")
    }
    if !strings.Contains(html, "Hello") {
        t.Error("expected title text")
    }
}
```

## Integration

`tinywasm/site` collects `RenderCSS()` and `IconSvg()` from all registered components (SSR build, `!wasm`). The WASM client receives only the struct logic and `Render()`/`OnMount()` — CSS and SVG strings never reach the binary.
