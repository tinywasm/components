---
name: component-creation
description: Standard for creating reusable, efficient, WebAssembly-ready UI components. Pattern based on style.Styler and widget package visual-contract.
---

# SKILL: tinywasm/components

A catalog of reusable, efficient, and WebAssembly-ready UI components in the TinyWasm ecosystem.

---

## Guide to Styling with Visual Contract (The One Page Guide)

With the construction harness closed, styling components in the TinyWasm ecosystem is done by implementing `style.Styler` using the typed `style.Sheet` DSL. It prevents local color patching, manual class strings, and media query overheads.

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

## Example: Anatomy and Styling of `TargetList`

A standard component implements `widget.Widget` to declare its identity, and optional `style.Styler` in a tagged `!wasm` file to define its visual rules.

```go
const (
	nameTargetList = widget.Name("targetlist")
	partRow        = widget.Part("row")
	partMenu       = widget.Part("menu")
	partOptions    = widget.Part("options")
)

func (l *TargetList) WidgetName() widget.Name { return nameTargetList }
func (l *TargetList) WidgetKind() widget.Kind { return widget.Listbox }

func (l *TargetList) Style() *style.Sheet {
	return style.Of(nameTargetList).
		Root(Stack(Space1), On(Sunken), Scrolls(), Round(RadiusMd)).
		Part(partRow, Row(Space2), On(Panel), Pad(Space2), Round(RadiusSm)).
		Part(partMenu, Stack(Space0), On(Panel), Raise(Floating), Clip()).
		When(widget.Selected, partRow, On(Selected)).
		Cue(widget.Hover, partRow, On(PanelHover))
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
