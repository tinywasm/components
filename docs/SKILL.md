---
name: component-creation
description: Standard for creating reusable, efficient, WebAssembly-ready UI components with strict naming, styling, and event handling conventions.
---

# SKILL: tinywasm/components

A catalog of reusable, efficient, and WebAssembly-ready UI components. Each component encapsulates styling, rendering, and event handling with a consistent API.

## Naming Convention

Component names **must be composed of at least two words** — single-word names are reserved for primitive builder functions in `tinywasm/dom` (e.g. `Div`, `Button`, `Input`).

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
    ├── css.go           # !wasm only: RenderCSS()
    ├── svg.go           # !wasm only: IconSvg()
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
- Do NOT style form elements — use `tinywasm/form`.

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

### Imports — use dot notation for clean code

**Always import `tinywasm/css` and `tinywasm/dom` with dot notation (`.`)** to make CSS DSL and DOM builders directly readable.

```go
import (
    . "github.com/tinywasm/css"  // Dot import: Class, ColorPrimary, Space4, etc.
    . "github.com/tinywasm/dom"  // Dot import: Div, Button, Element, Event, etc.
)
```

This makes the code read as intent (e.g., `Div(clsRoot, ...)`) rather than noise (e.g., `dom.Div(css.Class(...), ...)`).

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

### Class constants — PRIVATE

Declare **PRIVATE** (lowercase) `css.Class` constants at package scope. These are internal implementation details; the component's public API is its struct fields and methods, not CSS internals.

Both `mycomponent.go` (HTML side) and `css.go` (CSS side) reference the same constant — a rename is a compile error, not a silent drift.

```go
package mycomponent

import (
    . "github.com/tinywasm/css"
    . "github.com/tinywasm/dom"
)

var (
    clsRoot  Class = "mycomponent"
    clsTitle Class = "mycomponent-title"
    clsBody  Class = "mycomponent-body"
)

type MyComponent struct {
    Element
    Title string
}

func (c *MyComponent) Render() *Element {
    return Div(clsRoot.AsAttr(),
        H2(clsTitle.AsAttr(), c.Title),
        Div(clsBody.AsAttr()),
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
- `css.go`: Contains `RenderCSS() *Stylesheet` (required for styling).
- `svg.go`: Contains `IconSvg() map[string]string` (optional).
- `js.go`: Contains `RenderJS() []*js.Script` (optional).
- `html.go`: Contains `RenderHTML() *Element` (optional, for custom SSR templates).

No `.css` embed. No `//go:embed`. No `var css string`. No `//go:embed *.css`.

**`css.go` Example:**
```go
//go:build !wasm

package mycomponent

import . "github.com/tinywasm/css"

func (c *MyComponent) RenderCSS() *Stylesheet {
    return New(
        Rule(clsRoot,
            Display(Flex_),
            FlexDirection(Column),
            Gap(Space4),
            Padding(Space6),
            Background(ColorSurface),
            BorderRadius(RadiusMd),
        ),
        Rule(clsTitle,
            FontSize(TextLg),
            FontWeight(FontWeightBold),
            Color(ColorOnSurface),
        ),
        Rule(clsRoot.Hover(),
            BoxShadow(ShadowMd),
        ),
    )
}
```

**Key patterns:**
- Use dot import (`. "github.com/tinywasm/css"`) for clean DSL syntax
- Reference the **same private class constants** declared in `mycomponent.go` (e.g., `clsRoot`, not `ClsRoot`)
- Use `.AsAttr()` method on classes to convert them to DOM attributes
- `RenderCSS()` ships **component-scoped** CSS. `assetmin` discovers the component at build time, calls `.RenderCSS().String()`, and routes the result to the `middle` slot of `<head>`

Do **not** declare a `RootCSS()` function in a component package. `RootCSS()` is reserved for the app or `tinywasm/dom` (single-override rule — third-party `RootCSS()` is silently ignored by `assetmin` with a warning).

---

## 3. Icon Management (`IconSvgProvider`)

The framework injects the SVG sprite **directly into `<body>`** at server time. There is no public `/assets/icons.svg` URL — the sprite lives only in memory and inline in the HTML.

Pipeline: **`IconSvg()` in `svg.go`** → sprite built in memory → **injected inline in HTML** → `<svg><use href="#id">` in `Render()` resolves with zero network requests.

> **MANDATORY:** `IconSvg()` MUST be in `svg.go` (`//go:build !wasm`).
> SVG strings are dead code on WASM — never define icons in the main file.

> **MANDATORY:** All paths and shapes MUST include `fill="currentColor"` (or `stroke="currentColor"`) so CSS can control icon color via `fill` or `color` on any ancestor.

**`svg.go` — register the icon:**
```go
func (c *MyComponent) IconSvg() map[string]string {
    return map[string]string{
        // Do NOT include the wrapping <svg> tag — the system adds it.
        // Only internal content: paths, circles, etc.
        // Default viewBox is "0 0 16 16". Include viewBox="..." to override.
        "my-icon-id": `<path fill="currentColor" d="..." />`,
    }
}
```

**`mycomponent.go` `Render()` — reference the icon:**
```go
dom.Svg(dom.Use().Attr("href", "#my-icon-id")).Class("my-icon")
```

**Icon styles in `css.go` `RenderCSS()`:**
```go
var ClsIcon css.Class = "my-icon"

// inside RenderCSS():
Rule(ClsIcon,
    Width(Rem(1)),
    Height(Rem(1)),
    Fill(CurrentColor),
    Transition("transform", DurationFast),
),
```

---

## 4. Tests

All tests run through `gotest` — the tinywasm test runner. Do **not** use raw `go test`.

Install once:
```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

Run tests:
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
    html := c.Render().RenderHTML()

    if !Contains(html, string(ClsRoot)) {
        t.Error("expected root class")
    }
    if !Contains(html, "Hello") {
        t.Error("expected title text")
    }
}

func TestMyComponent_CSS(t *testing.T) {
    inst := &MyComponent{}
    sheet := inst.RenderCSS()
    if sheet == nil {
        t.Fatal("RenderCSS returned nil")
    }
    css := sheet.String()
    if !Contains(css, string(ClsRoot)) {
        t.Error("expected root class in CSS output")
    }
}
```

---

## Integration

Components are consumed by `tinywasm/assetmin` via the **compile-and-invoke** pipeline:

1. `assetmin` discovers all component modules in the project.
2. It generates a single combined `main.go` that imports every discovered component.
3. It compiles and runs that `main.go` **once** (one `go run` for all components combined), captures the aggregated JSON output, and routes each component's assets into the bundle.

This means:
- `assetmin` automatically infers the receiver type from the `Render*` methods; the boilerplate `SSRInstance()` is **not** required.
- A compile error in any component fails the whole extraction run. The compiler error is surfaced verbatim so the developer sees it as a normal `go build` failure.
- Extraction is cached by content hash of all component Go files. In steady state (no file changes) it costs ~0 ms.

**Slot routing:**
- Component `RenderCSS()` → `middle` slot (loaded after the document `:root`).
- App's `RootCSS()` (or `dom`'s default) → `open` slot.
- Root project's `RenderCSS()` → `close` slot (last word, can override anything).
