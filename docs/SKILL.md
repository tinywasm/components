---
name: component-creation
description: Standard for creating reusable, efficient, WebAssembly-ready UI components with strict naming, styling, and event handling conventions.
---

# SKILL: tinywasm/components

A catalog of reusable, efficient, and WebAssembly-ready UI components. Each component encapsulates styling, rendering, and event handling with a consistent API.

## Naming Convention

Component names **must be composed of at least two words** — single-word names are reserved for primitive builder functions in `tinywasm/html` (e.g. `Div`, `Button`, `Input`).

```
✅ ThemeSwitch   NavBar   DataTable   UserCard   SearchBox
❌ Switch        Nav      Table       Card       Search
```

This applies to both the Go struct name (`ThemeSwitch`) and the folder/package name (`themeswitch`). Two-word names make it immediately clear that something is a reusable component rather than a dom primitive or a stdlib type.

---

## File Structure

Each component resides in its own folder within `tinywasm/components`. There are **no `.css` files** — styles are expressed as Go code in backend-only files (`css.go`, `svg.go`, etc.).

```
tinywasm/components/
└── mycomponent/
    ├── mycomponent.go   # Struct, Render(), OnMount() — shared WASM + SSR
    ├── css.go           # !wasm only: RenderCSS() *css.Stylesheet
    ├── svg.go           # !wasm only: IconSvg() *svg.Sprite (optional)
    ├── html.go          # !wasm only: RenderHTML() string (optional)
    └── mycomponent_test.go
```

> There is NO `front.go` and NO `.css` file. WASM interactivity lives in `mycomponent.go`
> via `OnMount()`. CSS lives in `css.go` as a typed Go DSL — never as an embedded file.

---

## CSS Guidelines

CSS is written using the `tinywasm/css` typed DSL. All design decisions reference named tokens — never raw hex values, never raw `rem`/`px` magic numbers.

- **No `.css` files.** No `//go:embed`. CSS is Go code.
- **Use token constants** — `ColorPrimary`, `Space4`, `RadiusMd`, etc. (see `tinywasm/css/tokens.go`).
- **Never hardcode values** — `Rem(0.5)` is acceptable only when no scale token matches; flag it in the PR as a candidate for tokenization.
- **Exported `css.Class` constants** declare every class name. The same constant is used by both the HTML emission side and the CSS emission side of the component.
- **Class name prefix** — follow `<component>-*` convention (e.g. `"mycomponent-header"`).
- A component **must not** declare a `:root {}` block via `css.Root(...)`. That is reserved for the app or `tinywasm/dom`. Use only `Rule()` and `Selector()`.
- Do NOT style form elements — use `github.com/tinywasm/form`.

Available token groups (from `tinywasm/css/tokens.go`):

| Group | Examples |
|---|---|
| Color — brand | `ColorPrimary`, `ColorSecondary`, `ColorError` |
| Color — theme | `ColorBackground`, `ColorSurface`, `ColorOnSurface`, `ColorMuted` |
| Typography | `TextSm`, `TextBase`, `TextLg`, `FontWeightBold`, `LeadingNormal` |
| Spacing | `Space1` … `Space12` |
| Border radius | `RadiusSm`, `RadiusMd`, `RadiusLg`, `RadiusFull` |
| Elevation | `ShadowSm` … `ShadowXl` |
| Motion | `DurationFast`, `DurationBase`, `EaseInOut` |
| Z-index | `ZBase`, `ZDropdown`, `ZModal` |

---

## 1. Main File (`mycomponent.go`)

Contains the struct definition, private `css.Class` constants, `Render()` (shared SSR + WASM), and `OnMount()` (WASM only — no build tag; TinyGo eliminates dead code).

### Imports

**Import `tinywasm/css`, `tinywasm/html` and `tinywasm/dom` with dot notation (`.`)** to make CSS DSL, HTML and DOM builders directly readable. For `tinywasm/svg`, prefer a named import to avoid collisions.

```go
import (
    . "github.com/tinywasm/css"   // Class, Rule, Token, etc.
    . "github.com/tinywasm/html"  // Div, Span, H1, Nav, Button...
    . "github.com/tinywasm/dom"   // Event, Component, Reference, Render, Append
    "github.com/tinywasm/svg"     // named import for icons
)
```

This makes the code read as intent (e.g., `Div(clsRoot, ...)`) rather than noise.

### Embedding Rule — CRITICAL for TinyGo/WASM

**Always embed `Element` as a VALUE, never as a pointer.**

```go
// ✅ CORRECT — value embed, zero GC overhead, no nil risk
type MyComponent struct {
    Element
    Title string
}

// ❌ WRONG — pointer embed causes 2 heap allocations, requires nil-guard,
//            risks nil panic in WASM, wastes GC cycles in TinyGo
type MyComponent struct {
    *Element
}
```

### Class constants

Declare `css.Class` constants at package scope. Exported constants (starting with uppercase) are preferred if they need to be accessible from other packages (e.g., tests or parent components).

```go
package mycomponent

import (
    . "github.com/tinywasm/css"
    . "github.com/tinywasm/dom"
    . "github.com/tinywasm/html"
)

var (
    ClsRoot  Class = "mycomponent"
    ClsTitle Class = "mycomponent-title"
    ClsBody  Class = "mycomponent-body"
)

type MyComponent struct {
    Element
    Title string
}

func (c *MyComponent) Render() *Element {
    return Div(ClsRoot.AsAttr(),
        H2(ClsTitle.AsAttr(), c.Title),
        Div(ClsBody.AsAttr()),
    )
}

// OnMount wires events after the component is injected into the DOM.
// TinyGo eliminates this as dead code on SSR builds — no build tag needed.
func (c *MyComponent) OnMount() {
    id := c.GetID()
    if el, ok := Get(id); ok {
        el.On("click", func(e Event) {
            // handle click
        })
    }
}
```

---

## 2. Backend Files (`css.go`, `svg.go`, `js.go`, `html.go`)

**CRITICAL:** These files MUST have the `//go:build !wasm` build tag.

Assets are split into separate files by type for better organization and discovery by `assetmin`:
- `css.go`: Contains `RenderCSS() *css.Stylesheet` (required for styling).
- `svg.go`: Contains `IconSvg() *svg.Sprite` (optional).
- `js.go`: Contains `RenderJS() []*js.Script` (optional).
- `html.go`: Contains `RenderHTML() string` (optional, for custom SSR templates).

No `.css` embed. No `//go:embed`. No `var css string`.

---

## 3. Icon Management (`svg.go`)

The framework injects the SVG sprite **directly into `<body>`** at server time.

Pipeline: **`IconSvg()` in `svg.go`** → sprite built in memory → **injected inline in HTML** → `<svg><use href="#id">` in `Render()` resolves with zero network requests.

> **MANDATORY:** `IconSvg()` MUST be in `svg.go` (`//go:build !wasm`) and return `*svg.Sprite`.
> SVG strings are dead code on WASM — never define icons in the main file.

> **MANDATORY:** All paths and shapes MUST include `fill="currentColor"` (or `stroke="currentColor"`) so CSS can control icon color via `fill` or `color` on any ancestor.

**`svg.go` — register the icon:**
```go
//go:build !wasm
import "github.com/tinywasm/svg"

func (c *MyComponent) IconSvg() *svg.Sprite {
    return svg.New().
        Add("my-icon-id", `<path fill="currentColor" d="..." />`)
}
```

**`mycomponent.go` `Render()` — reference the icon:**
```go
import "github.com/tinywasm/svg"

svg.Icon("my-icon-id", "my-icon-class")
// → <svg aria-hidden='true' class='my-icon-class'><use href='#my-icon-id'></svg>
```

---

## 4. Tests

All tests run through `gotest` — the tinywasm test runner.

```bash
gotest
```

Test files: `mycomponent_test.go` in the package root (`package mycomponent`).

```go
package mycomponent

import (
    "testing"
    . "github.com/tinywasm/fmt"
)

func TestMyComponent_Render(t *testing.T) {
    c := &MyComponent{Title: "Hello"}
    html := c.Render().String() // Use .String() to get the HTML representation

    if !Contains(html, string(ClsRoot)) {
        t.Error("expected root class")
    }
}
```

---

## Integration

Components are consumed by `tinywasm/assetmin` via the **compile-and-invoke** pipeline.
- `assetmin` automatically infers the receiver type from the `Render*` methods.
- Extraction is cached by content hash of all component Go files.
