---
name: component-creation
description: Standard for creating reusable, efficient, WebAssembly-ready UI components. ONE construction pattern — Render() + optional Init(ctx); reactive state via typed signals (no generics); no OnMount/Update.
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
    ├── mycomponent.go   # Struct, Render() + optional Init(ctx) — shared WASM + SSR
    ├── css.go           # !wasm only: RenderCSS() *css.Stylesheet
    ├── svg.go           # !wasm only: IconSvg() *svg.Sprite (optional)
    ├── html.go          # !wasm only: RenderHTML() string (optional)
    └── mycomponent_test.go
```

> There is NO `front.go` and NO `.css` file. Interactivity is declared **inside `Render()`** with
> `.On(...)` handlers that mutate **typed signals** — there is NO `OnMount`/`OnUpdate`/`OnUnmount`
> and NO manual `Update()`. CSS lives in `css.go` as a typed Go DSL — never as an embedded file.

---

## CSS Guidelines

CSS is written using the `tinywasm/css` typed DSL. All design decisions reference named tokens — never raw hex values, never raw `rem`/`px` magic numbers.

- **No `.css` files.** No `//go:embed`. CSS is Go code.
- **Use token constants** — `ColorPrimary`, `Space4`, `RadiusMd`, etc. (see `tinywasm/css/tokens.go`).
- **Never hardcode values** — `Rem(0.5)` is acceptable only when no scale token matches; flag it in the PR as a candidate for tokenization.
- **`css.Class` constants** declare every class name, shared by the HTML and CSS emission sides. Keep
  them **unexported** unless another package genuinely needs them (minimal public surface).
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

Contains the struct definition, private `css.Class` constants, `Render()` (shared SSR + WASM), and an
optional `Init(ctx Ctx)` (one-time setup). **No `OnMount`/`OnUpdate`/`OnUnmount`, no `Update()`.**

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

Declare `css.Class` constants at package scope. **Keep them unexported** unless another package
genuinely needs them — export only what a component *user* types (minimal public surface).

```go
package mycomponent

import (
    . "github.com/tinywasm/css"
    . "github.com/tinywasm/dom"
    . "github.com/tinywasm/html"
)

var (
    clsRoot  Class = "mycomponent"
    clsTitle Class = "mycomponent-title"
    clsBody  Class = "mycomponent-body"
)

type MyComponent struct {
    Element                 // value-embed, never *Element
    Title string            // static config: a plain field is fine
    open  *SignalBool       // DYNAMIC UI state lives in a typed signal, never a plain field
}

// Init runs ONCE before the first render (optional). Create/seed signals here — and load storage,
// fetch, or subscribe. Setting a signal (even from a goroutine) patches the bound DOM directly.
func (c *MyComponent) Init(_ Ctx) {
    c.open = NewBool(false)
}

func (c *MyComponent) Render() *Element {
    return Div(clsRoot.AsAttr(),
        H2(clsTitle.AsAttr(), c.Title),
        // Events are declared HERE. The handler only mutates the signal — the DOM patches itself.
        Button("Toggle").On("click", func(Event) { c.open.Toggle() }),
        // Reactive structure: Show mounts/unmounts the body when `open` flips.
        Show(c.open, func() *Element { return Div(clsBody.AsAttr(), "content") }),
    )
}
```

There is no `OnMount`, no `dom.Get(id)`, no `Update()`. You **cannot** forget to refresh — changing a
signal IS the refresh, and it patches only the bound node (no whole-element re-render, no Virtual DOM).

---

## 1b. State & Reactivity — the ONE way

UI state lives in **typed signals** (the DOM boundary is `string`/`bool`, so signals are concrete —
**never generics / `Signal[T]`**, matching the ecosystem's `tinywasm/fmt` codec rule "cero any, cero map"):

| Signal | Create | Use for |
|---|---|---|
| `*SignalString` | `NewString("")` | text, attribute values, two-way input |
| `*SignalBool` | `NewBool(false)` | class/attr toggles, `Show` conditions |
| `*SignalNodes` | `NewNodes(...)` | lists of rendered rows (`*Element`) |

Read/write with `Get()` / `Set(v)`; `*SignalBool` also has `Toggle()`.

**Bind a signal to one DOM location** (the binding patches only that spot):

```go
el.BindText(sig)                       // textContent tracks the signal
el.BindAttr("title", sig)              // attribute value
el.BindClass("active", boolSig)        // toggle a class
el.BindAttrBool("disabled", boolSig)   // boolean attribute (disabled/checked/hidden)
in.Bind(strSig)                        // two-way <input>/<textarea> (cursor + IME safe)
in.Autofocus()                         // focus this node when it first appears
Show(boolSig, render)                  // mount/unmount a subtree
ul.BindChildren(nodesSig)              // keyed list; build []*Element in a loop, each .Key(id)
```

**Computed values — pass a function, no dependency list.** The framework detects which signals the
function reads (auto-tracking) and re-runs it when they change:

```go
btn.BindTextFunc(func() string { if c.open.Get() { return "Hide" }; return "Show" })
row.BindClassFunc("error", func() bool { return c.msg.Get() != "" })
// also: BindAttrFunc, BindAttrBoolFunc. For a NAMED computed shared across binds: DeriveString(fn).
```

You never assemble a deps list — reading a signal inside the closure is enough. That makes computed
UI **impossible to get stale**.

**Async / one-time setup** goes in `Init(ctx Ctx)`; setting a signal from a goroutine patches the DOM.
Register teardown with `ctx.OnCleanup(fn)`. `Init` is the only optional hook — and you never name its
interface, you just write the method.

---

## 2. Backend Files (`css.go`, `svg.go`, `js.go`, `html.go`)

**CRITICAL:** These files MUST have the `//go:build !wasm` build tag.

Assets are split into separate files by type for better organization and discovery by `assetmin`:
- `css.go`: Contains `RenderCSS() *css.Stylesheet` (required for styling).
- `svg.go`: Contains `IconID() string` (if used in UIModule) and `IconSvg() *svg.Sprite` (required if icons present).
- `js.go`: Contains `RenderJS() []*js.Script` (optional).
- `html.go`: Contains `RenderHTML() string` (optional, for custom SSR templates).

No `.css` embed. No `//go:embed`. No `var css string`.

---

## 3. Icon Management (`svg.go`)

The framework injects the SVG sprite **directly into `<body>`** at server time.

Pipeline: **`IconID()` + `IconSvg()` in `svg.go`** → sprite built in memory → **injected inline in HTML** → `<svg><use href="#id">` in `Render()` resolves with zero network requests.

> **MANDATORY:** `svg.go` MUST have `//go:build !wasm` build tag.
> **MANDATORY:** `IconID()` returns the icon ID as string; `IconSvg()` builds the sprite.
> SVG definition code is dead code on WASM — never define icons outside `svg.go`.

> **MANDATORY:** All paths and shapes MUST include `fill="currentColor"` (or `stroke="currentColor"`) so CSS can control icon color via `fill` or `color` on any ancestor.

**`svg.go` — register the icon (backend only):**
```go
//go:build !wasm

package mycomponent

import "github.com/tinywasm/svg"

func (m *MyComponent) IconID() string {
	return "my-icon-id"
}

func (m *MyComponent) IconSvg() *svg.Sprite {
	return svg.NewSprite(
		svg.Define("my-icon-id", "0 0 16 16",
			svg.Path(`M2 2h5v5H2zm7 0h5v5H9zm-7 7h5v5H2zm7 0h5v5H9z`),
		),
	)
}
```

**`mycomponent.go` `Render()` — reference the icon:**
```go
func (c *MyComponent) Render() dom.Element {
	// Use Icon() method from interfaced UIModule
	// or manually reference the icon ID
	return Svg().Attr("aria-hidden", "true").
		Child(Use().Attr("href", "#"+c.IconID()))
}
```

---

## 4. Tests

All tests run through `gotest` — the tinywasm test runner (WASM tests run against a real DOM).
External agents install it first:

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
gotest
```

Use `gotest`, never `go test`. Stdlib assertions only (no testify). Dual WASM/stdlib via build tags
sharing one runner. Cover the **frequent use cases**: load-on-init (no flash), two-way input
(IME/cursor safe), derived value, conditional `Show`, keyed list, and any no-recursion regression.

Test files: `mycomponent_test.go` in the package root (`package mycomponent`).

```go
package mycomponent

import (
    "testing"
    . "github.com/tinywasm/fmt"
)

func TestMyComponent_Render(t *testing.T) {
    c := &MyComponent{Title: "Hello"}
    c.Init(nil)                 // seed signals (the engine calls Init once before the first render)
    html := c.Render().String() // .String() yields the static HTML (SSR parity)

    if !Contains(html, string(clsRoot)) {
        t.Error("expected root class")
    }
}
```

---

## Integration

Components are consumed by `tinywasm/assetmin` via the **compile-and-invoke** pipeline.
- `assetmin` automatically infers the receiver type from the `Render*` methods.
- Extraction is cached by content hash of all component Go files.
