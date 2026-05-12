# TinyWasm Components
<img src="docs/img/badges.svg">

A catalog of reusable, efficient, and WebAssembly-ready UI components for the TinyWasm ecosystem.

## Features

-   **Zero Boilerplate**: Components handle their own state, styling, and event binding.
-   **SSR/CSR Split**: CSS and heavy assets are server-side only; WASM binaries remain tiny.
-   **Fluent API**: Uses `tinywasm/dom` for type-safe, declarative UI construction.
-   **Theme Integration**: Consumes canonical tokens from `tinywasm/dom`.
-   **Explicit Registration**: Only pay for what you use (tree-shakeable).
-   **No Magic**: Standard Go structs and interfaces.

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

## Available Components

See [Component Catalog](docs/CATALOG.md) for full documentation.

-   **ActionButton**: Primary/secondary actions with variants.
-   **ContentCard**: Content container with header/body/footer.
-   **DataTable**: Data tables with headers and rows.
-   **NavBar**: Navigation menu with icon support.
-   **Dialog**: Modal dialog overlays.
-   **ThemeToggle**: Theme switcher (light/dark/auto).

## Forms

Forms are NOT part of `tinywasm/components`. Use `github.com/tinywasm/form` directly — it is the standard form library for the tinywasm ecosystem.

## Development

See [Component Skill Guide](docs/SKILL.md) to learn how to build your own components.

### Prerequisites

-   Go 1.25+
-   TinyGo (for WASM compilation)

### Running Tests

```bash
go test ./...
```
