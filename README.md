# TinyWasm Components
<img src="docs/img/badges.svg">

A catalog of reusable, efficient, and WebAssembly-ready UI components for the TinyWasm ecosystem.

## Features

-   **One way to build**: a component is just `Render()` (+ optional `Init(ctx)`). No
    `OnMount`/`OnUpdate`/`OnUnmount`, no manual `Update()` — impossible to wire wrong.
-   **Fine-grained reactivity**: dynamic state lives in **typed signals**; changing one patches only
    the bound DOM node — no whole-element re-render, no Virtual DOM.
-   **No generics**: concrete typed signals (`SignalString`/`SignalBool`/`SignalNodes`), matching the
    ecosystem's "cero any, cero map" rule — readable, TinyGo-friendly.
-   **SSR/CSR Split**: CSS and heavy assets are server-side only; WASM binaries remain tiny.
-   **Fluent API**: Uses `tinywasm/dom` for type-safe, declarative UI construction.
-   **Theme Integration**: Consumes canonical tokens from `tinywasm/dom`.
-   **Explicit Registration**: Only pay for what you use (tree-shakeable).

## Installation

```bash
go get github.com/tinywasm/components
```

To enable the default theme, inject `dom.ThemeCSS` into your page `<head>` once via your site builder.

## Quick Start

### 1. Import Components

```go
import (
    "github.com/tinywasm/components/actionbutton"
    "github.com/tinywasm/components/contentcard"
    . "github.com/tinywasm/dom"
)
```

### 2. Instantiate and Render

```go
func main() {
    // Create a button
    btn := &actionbutton.ActionButton{
        Text: "Click Me",
        Variant: "primary",
        OnClick: func(e Event) {
            println("Button clicked!")
        },
    }

    // Render to the DOM
    Append("app", btn)
}
```

## The Component API (build your own)

A component implements **only** `Render() *Element`, plus an optional `Init(ctx Ctx)` that runs **once**
before the first render. Dynamic state lives in a **signal**; you change it with `Set`, and the bound
DOM updates itself — you never call `Update()`.

**Start here — the simplest reactive component.** A signal, bound to text, changed on click:

```go
import (
    . "github.com/tinywasm/dom"
    . "github.com/tinywasm/html"
)

type LikeButton struct {
    Element                  // value-embed, never *Element
    label *SignalString      // dynamic state is a signal, not a plain field
}

func (c *LikeButton) Init(_ Ctx) { c.label = NewString("Like") } // starting state

func (c *LikeButton) Render() *Element {
    return Button().
        BindText(c.label).                                   // show the signal
        On("click", func(Event) { c.label.Set("Liked ♥") })  // change it → the button updates itself
}
```

Three ideas: **state is a signal**, **`BindText` shows it**, **`Set` changes it**. That's the whole loop.

**Growing up — computed text + show/hide.** When a value is *computed* from state, pass a function to
`BindTextFunc` (no dependency list — it auto-detects the signals you read). `Toggle()` flips a bool:

```go
type Disclosure struct {
    Element
    open *SignalBool
}
func (c *Disclosure) Init(_ Ctx) { c.open = NewBool(false) }

func (c *Disclosure) Render() *Element {
    return Div(
        Button().
            BindTextFunc(func() string { if c.open.Get() { return "Hide" }; return "Show" }).
            On("click", func(Event) { c.open.Toggle() }),
        Show(c.open, func() *Element { return Div("details…") }), // mounts/unmounts reactively
    )
}
```

**Signals** (no generics — concrete typed cells):

| Type | Create | For |
|---|---|---|
| `*SignalString` | `NewString("")` | text, attribute values, two-way input |
| `*SignalBool` | `NewBool(false)` | class/attr toggles, `Show` conditions (`Toggle()` to flip) |
| `*SignalNodes` | `NewNodes(...)` | keyed lists of rendered rows |

**Bindings** (each patches only its own DOM spot):

- Show a signal directly: `BindText`, `BindAttr`, `BindClass`, `BindAttrBool`, `Bind` (two-way input,
  cursor/IME safe), `BindChildren` (keyed list), `Autofocus`, `Show(cond, render)`.
- Show a *computed* value (function, auto-tracked — no deps list): `BindTextFunc`, `BindAttrFunc`,
  `BindClassFunc`, `BindAttrBoolFunc`. (`DeriveString`/`DeriveBool` for a named shared computed.)

For async/teardown, set signals from a goroutine in `Init` and use `ctx.OnCleanup(fn)`.

See [Component Skill Guide](docs/SKILL.md) for the full standard.

## Available Components

See [Component Catalog](docs/CATALOG.md) for full documentation.

-   **ActionButton**: Primary/secondary actions with variants.
-   **ContentCard**: Content container with header/body/footer.
-   **DataTable**: Data tables with headers and rows.
-   **ModalDialog**: Centered modal overlay with backdrop.
-   **SelectSearch**: Searchable dropdown with live filtering.
-   **ThemeToggle**: Theme switcher (dark/light).

## Where components fit

The ecosystem is layered:

```
tinywasm/components  →  raw reusable pieces (this repo). No layout knowledge.
tinywasm/layout       →  published layout skeletons (platformd shell, rightpanel,
                          crudview) with named slots (Form, Detail, HeadControls, …).
consumer's composition root (e.g. config/layouts)
                       →  preconfigures a layout once, injecting components into
                          its slots. Modules pick a preconfigured layout; they
                          never assemble components by hand.
```

Two rules keep this layering intact:

- **`components` never imports `tinywasm/layout`.** The reverse (layout importing
  a component) is allowed but assembly belongs to the consumer, not this repo.
- **Assembly never lives in modules.** A module consumes a preconfigured layout;
  it does not construct `ActionButton`/`DataTable`/etc. by hand inside a slot.

Every cataloged component is **slot-ready**: it can be dropped into any layout
slot as a `dom.Component` with zero per-component configuration, and it is
**theme-driven**: it inherits branding from `RootCSS` with no hardcoded colors.
See [docs/SKILL.md](docs/SKILL.md) for the slot-readiness contract and
[docs/CATALOG.md](docs/CATALOG.md) for per-component status.

## Forms

Forms are NOT part of `tinywasm/components`. Use `github.com/tinywasm/form` directly — it is the standard form library for the tinywasm ecosystem.

## Development

See [Component Skill Guide](docs/SKILL.md) to learn how to build your own components.

### Prerequisites

-   Go 1.25+
-   TinyGo (for WASM compilation)

### Running Tests

Use `gotest` (the tinywasm runner — WASM tests run against a real DOM), never `go test`:

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
gotest
```
