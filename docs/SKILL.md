---
name: component-creation
description: Standard for creating reusable, efficient, WebAssembly-ready UI components. Pattern based on RenderCSS() and the widget package visual-contract.
---

# SKILL: tinywasm/components

A catalog of reusable, efficient, and WebAssembly-ready UI components in the TinyWasm ecosystem.

---

## Guide to Styling with Visual Contract (The One Page Guide)

With the construction harness closed, styling components in the TinyWasm ecosystem is done by declaring `RenderCSS() *css.Stylesheet`, built with the typed `style.Sheet` DSL. It prevents local color patching, manual class strings, and media query overheads.

### Quick Reference

| To achieve… | Use |
|---|---|
| A card over the page background | `On(Panel)` |
| A sunken well inside a card | `On(Sunken)` |
| Vertical layout rhythm | `Stack(Space2)` |
| Two panels that stack on mobile | `Split(style.RatioTwoThirds, Space2)` |
| A grid that auto-fits its columns | `Grid(TrackMd, Space2)` |
| Make element take up all available height | `Fill()` |
| Enable internal scroll instead of expanding | `Scrolls()` |
| Keep layout fixed (do not reflow on mobile) | `Fixed()` |
| Highlight a selected row / part | `When(widget.Selected, partRow, On(Selected))` |

---

## File Structure

Each component resides in its own folder within `tinywasm/components`. There are **no `.css` files** — styles and assets are expressed as Go code in backend-only files (`css.go`, `svg.go`, etc.).

```
tinywasm/components/
└── mycomponent/
    ├── mycomponent.go   # Struct, Render() + optional Init(ctx) — shared WASM + SSR
    ├── css.go           # !wasm only: RenderCSS() *css.Stylesheet visual sheet
    ├── svg.go           # !wasm only: IconSvg() *sprite.Sprite (optional)
    └── mycomponent_test.go
```

---

## Icon Management (`svg.go`)

The framework injects the SVG sprite **directly into `<body>`** at server time.

Pipeline: **`IconSvg()` in `svg.go`** → sprite built in memory → **injected inline in HTML** → `<svg><use href="#id">` in `Render()` resolves with zero network requests.

> **MANDATORY:** `svg.go` MUST have `//go:build !wasm` build tag.
> SVG definition code is dead code on WASM — never define icons outside `svg.go`.

> **MANDATORY:** All paths and shapes MUST include `fill="currentColor"` (or `stroke="currentColor"`) so CSS can control icon color via `fill` or `color` on any ancestor.

**`svg.go` — register the icon (backend only):**
```go
//go:build !wasm

package mycomponent

import (
	"github.com/tinywasm/svg"
	"github.com/tinywasm/svg/sprite"
)

func (m *MyComponent) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		sprite.Define(iconX, "0 0 16 16",
			sprite.Path(`M2 2h5v5H2zm7 0h5v5H9zm-7 7h5v5H2zm7 0h5v5H9z`),
		),
	)
}
```

---

## Example: Anatomy and Styling of `TargetList`

A standard component implements `widget.Widget` to declare its identity, and optionally declares `RenderCSS()` in a tagged `!wasm` file to define its visual rules.

```go
const (
	nameTargetList = widget.Name("targetlist")
	partRow        = widget.Part("row")
	partMenu       = widget.Part("menu")
	partOptions    = widget.Part("options")
)

func (l *TargetList) WidgetName() widget.Name { return nameTargetList }
func (l *TargetList) WidgetKind() widget.Kind { return widget.Listbox }

func (l *TargetList) RenderCSS() *css.Stylesheet {
	return style.Of(nameTargetList).
		Root(Stack(Space1), On(Sunken), Scrolls(), Round(RadiusMd)).
		Part(partRow, Row(Space2), On(Panel), Pad(Space2), Round(RadiusSm)).
		Part(partMenu, Stack(Space0), On(Panel), Raise(Floating), Clip()).
		When(widget.Selected, partRow, On(Selected)).
		Cue(widget.Hover, partRow, On(PanelHover)).
		Stylesheet()
}
```

---

## Component Contract

A component implements **only**:

- `Render() *dom.Element` — describes the structure ONCE; dynamic parts are signal bindings.
- `Init(ctx dom.Ctx)` — optional; runs ONCE before first render (load storage, fetch, subscribe).

There is **NO** `OnMount`/`OnUpdate`/`OnUnmount` and **NO** manual `Update()`.

State the UI shows lives in **typed signals**; changing one patches only the bound DOM node — never re-render a whole element, never use a Virtual DOM.

### Naming Constraint

A component's Go struct name and folder/package name must be at least two words, identifying its style/class:

```
✅ ModalDialog (modaldialog/)   ThemeToggle (themetoggle/)   ActionButton (actionbutton/)
❌ Dialog (dialog/)              Toggle (toggle/)              Button (button/)
```

### WASM / TinyGo Constraints

- **No Go stdlib** (`fmt`/`strings`/`errors`): use `github.com/tinywasm/fmt`.
- **WASM files must never import** `github.com/tinywasm/css` or `github.com/tinywasm/widget/style`. Move all styling code into `css.go` under `//go:build !wasm`.
- Use `switch` instead of `map` for performance and heap optimization in TinyGo.
