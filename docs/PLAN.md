---
PLAN: "feat(searchbar): extract crudview's fixed search bar into a slot-ready component"
TAG: v0.4.0
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
>
> **Stage M2 of a 3-repo change.** Master plan:
> `app-releases/docs/LAYOUT_UNIFICATION_MASTER_PLAN.md`.

# Plan — `components`: a `SearchBar` component

## Why

`github.com/tinywasm/layout/crudview` paints a search bar **inline and
unconditionally**: markup in `crudview.go`, styles in `css.go`, magnifier sprite
in `svg.go`. That hardcodes one filter shape into a reusable layout skeleton. A
consumer who wants to filter by date (a calendar), by category (a select), or by
nothing at all has no way in.

The fix spans three repositories:

| | Repo | Plan | What |
|---|---|---|---|
| M1 | `widget` | `widget/docs/PLAN.md` | names the seam: `widget.Filterable` |
| **M2** | **`components`** | **this plan** | `SearchBar` implements it |
| M3 | `layout` | `layout/docs/PLAN.md` | one skeleton, filter as a slot |

## Dependency — read this first

**M1 must already be published.** Verify before starting:

```sh
go doc github.com/tinywasm/widget Filterable
```

Expected:

```go
type Filterable interface{ OnFilterChange(func(term string)) }
```

If that fails, stop — the gate is not released yet.

⚠️ **Do NOT declare a local `FilterSource`/`Filterable` interface in this
repository if the lookup fails.** Recreating an upstream symbol downstream is
the exact defect this change exists to remove. Report the gap and stop.

⚠️ **Do NOT touch `crudview`.** `layout` is a different repository and is not
part of this dispatch. This plan only ADDS a package here; `layout` keeps
compiling against the current release until M3 runs.

---

## Rules this repository enforces (read before writing code)

These are checked by `conformance_test.go` at the module root and will fail the
build if violated.

- **No Go stdlib in shared (untagged) files.** Use `github.com/tinywasm/fmt`
  instead of `fmt`/`strings`/`strconv`/`errors`. `_test.go` files and
  `//go:build !wasm` files are exempt.
- **`github.com/tinywasm/css` may only be imported from a file named
  `css.go`.** Importing it anywhere else is a hard test failure.
- **`css.go` declares exactly one method named `RenderCSS`** and no method named
  `Style`.
- **No `css.Raw` / `css.RawItem`**, no hex/`rgb()`/`hsl()` colour literals, no
  viewport units (`vw`/`vh`), no bare colour names — everything goes through the
  `widget/style` DSL.
- **Components MUST NOT declare `RootCSS()`.** `:root` tokens belong to the app.
- **Value embed `dom.Element`**, never `*dom.Element`.
- **Two-word component name**: `SearchBar` / package `searchbar`. ✅
- **SSR split**: `RenderCSS` in `css.go`, `IconSvg` in `svg.go`, both with
  `//go:build !wasm`. Neither may reach the WASM binary.
- **No `map`** in shared code — use a `switch` or a `[]fmt.KeyValue` scan.
  (`_test.go` and `!wasm` files may use maps.)
- **No `OnMount`.** A component implements `Render() *dom.Element` and
  optionally `Init(dom.Ctx)`. State lives in signals.

---

## Stage 1 — `searchbar/searchbar.go`

Create the directory `searchbar/` and this file. **No build tag** — it compiles
for both WASM and SSR.

```go
package searchbar

import (
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/svg"
	"github.com/tinywasm/widget"
)

// NameSearchBar is the widget identity; it produces the class prefix
// "searchbar" and the part classes "searchbar__icon", "searchbar__glyph",
// "searchbar__input".
const NameSearchBar = widget.Name("searchbar")

const (
	PartIcon  = widget.Part("icon")  // the square coloured cap holding the magnifier
	PartGlyph = widget.Part("glyph") // the magnifier <svg> itself
	PartInput = widget.Part("input") // the text field, the body of the bar
)

var (
	clsSearchBar = NameSearchBar.Root()
	clsIcon      = NameSearchBar.Class(PartIcon)
	clsGlyph     = NameSearchBar.Class(PartGlyph)
	clsInput     = NameSearchBar.Class(PartInput)
)

// iconMagnifier is registered by IconSvg in svg.go.
const iconMagnifier = svg.Icon("searchbar-magnifier")

// defaultPlaceholder is what the field says when the host sets none.
const defaultPlaceholder = "Search…"

// SearchBar is a single-control filter bar: a magnifier cap followed by a text
// field. It holds no list and knows nothing about what it filters — it reports
// the term and the host decides what that means.
//
// It satisfies widget.Filterable, which is the whole reason it is swappable: a
// host holds the seam, not this type, so a calendar or a select that also
// implements Filterable takes its place with no change to the host.
type SearchBar struct {
	Element // value embed — NEVER *dom.Element (TinyGo heap constraint)

	// Placeholder is the field's placeholder text. Empty uses defaultPlaceholder.
	Placeholder string

	onFilter func(term string)
}

// Compile-time proof that the seam is satisfied. Keep this line: if the
// upstream interface ever changes shape, this is what fails, here, instead of
// failing at a host's type assertion where it would silently evaluate false and
// leave the bar wired to nothing.
var _ widget.Filterable = (*SearchBar)(nil)

func (s *SearchBar) WidgetName() widget.Name { return NameSearchBar }

// Form, not Combobox: the bar has no popup and no options list. Combobox would
// license Open/Selected states this control can never be in.
func (s *SearchBar) WidgetKind() widget.Kind { return widget.Form }

// OnFilterChange implements widget.Filterable: it registers the sink for every
// keystroke. The host calls it once while wiring its filter slot; passing nil
// clears the sink.
//
// The signature is fixed by widget.Filterable — do not add a parameter, do not
// return anything, do not add a companion getter.
func (s *SearchBar) OnFilterChange(fn func(term string)) { s.onFilter = fn }

func (s *SearchBar) Render() *Element {
	root := Div().Set(clsSearchBar.AsAttr())

	// A <label> so a click on the cap focuses nothing accidentally and the
	// magnifier is announced as decoration, not as a button.
	root.Child(Label().Set(clsIcon.AsAttr()).
		Child(iconMagnifier.Render(string(clsGlyph))))

	placeholder := s.Placeholder
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}

	// size="1" collapses the input's intrinsic width (the UA default is 20
	// characters): the Row flexbox wraps on the flex base size, and ~240px of
	// intrinsic width made the bar wrap onto two lines inside narrow hosts. The
	// flex Grow still sizes the rendered field; size only feeds the intrinsic
	// measurement.
	input := Input("search").Set(clsInput.AsAttr()).
		Attr("placeholder", placeholder).
		Attr("size", "1")
	input.On("input", func(e Event) {
		if s.onFilter != nil {
			s.onFilter(e.TargetValue())
		}
	})
	root.Child(input)

	return root
}
```

### Acceptance

- `grep -rn "OnMount" searchbar/` → empty.
- `grep -rn "\*dom.Element\|\*Element$" searchbar/searchbar.go` → no pointer embed.
- `searchbar/searchbar.go` imports neither `github.com/tinywasm/css` nor
  `github.com/tinywasm/widget/style`.

---

## Stage 2 — `searchbar/css.go`

The rules below are **moved verbatim** from `crudview/css.go` (its `search`,
`search-icon` and `search-input` parts) — including the reasoning, because that
reasoning is the contract. Do not "improve" the values.

```go
//go:build !wasm

package searchbar

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the search bar's visual contract using the style DSL.
func (s *SearchBar) RenderCSS() *css.Stylesheet {
	return style.For(s).
		// The bar is ONE control, not a card holding two loose pieces: the
		// magnifier is the bar's leading cap, the input its body, and a gap or a
		// card of its own between them saws the bar back into separate boxes.
		// The root carries the radius and clips, so cap and body stay square and
		// still read as one rounded bar.
		// ControlBox pins the bar to --control-height, the token every control
		// in the ecosystem answers to — that is what lets a host stack this bar
		// against its own buttons and have the heights agree by construction.
		Root(
			style.Row(style.SpaceNone),
			style.Round(style.RadiusMd),
			style.HideOverflow(),
			style.ControlBox(),
			style.KeepSize(),
		).
		// The magnifier is the bar's square cap: aspect-ratio, not padding, sets
		// the width — a padded box drifts off the control token (the old
		// Pad(Space2)+icon-box measured 40px against the host's 66), while the
		// square derives from the same --control-height as everything else.
		Part(PartIcon,
			style.As(style.Primary),
			style.MediaBox(style.AspectSquare),
			style.ControlBox(),
			style.KeepSize(),
		).
		// A bare <svg> with no box falls back to 300x150; IconBox pins it.
		Part(PartGlyph,
			style.IconBox(style.IconMd),
		).
		// The input is the body of the bar: it grows into whatever the cap
		// leaves and answers to the same control height, so cap and body can
		// never drift apart vertically — the mismatch that left a 25px field
		// floating in the middle of a 72px strip.
		Part(PartInput,
			style.As(style.Inset),
			style.Pad(style.Space2),
			style.Grow(),
			style.ControlBox(),
		).
		Stylesheet()
}
```

### Acceptance

- No `When(...)` call appears in this file — `widget.Form` allows only
  `Invalid` plus the universal `Disabled`/`Locked`/`Busy`, and this bar declares
  no state. `TestKindAllowsEveryState` parses this file and will reject any
  `When` whose state the Kind disallows.
- No string literal in this file matches `#[0-9a-f]{3,6}`, `rgb(`, `hsl(`,
  `\d+(vw|vh)`, or a `var(--…)` outside the token catalogue.

---

## Stage 3 — `searchbar/svg.go`

The magnifier is **moved** from `crudview/svg.go`, path data unchanged. The
sprite id changes from `icon-crud-search` to `searchbar-magnifier` because
ownership moves with it.

```go
//go:build !wasm

package searchbar

import (
	"github.com/tinywasm/svg/sprite"
)

// IconSvg registers the bar's magnifier. Method receiver (not a free function)
// so tinywasm/ssr detects a single receiver type for the package and emits
// RenderCSS + IconSvg together.
func (s *SearchBar) IconSvg() *sprite.Sprite {
	// A FILL shape (closed outline), never stroke lines: the sprite renders
	// every path with fill="currentColor" and no stroke, so a line-only path
	// (e.g. "M8 1v14") has zero area and is invisible. The lens hole renders via
	// the inner circle subpath winding opposite the outer.
	return sprite.NewSprite(
		sprite.Define(iconMagnifier, "0 0 512 512", sprite.Path("M416 208c0 45.9-14.9 88.3-40 122.7L502.6 457.4c12.5 12.5 12.5 32.8 0 45.3s-32.8 12.5-45.3 0L330.7 376c-34.4 25.2-76.8 40-122.7 40C93.1 416 0 322.9 0 208S93.1 0 208 0S416 93.1 416 208zM208 352a144 144 0 1 0 0-288 144 144 0 1 0 0 288z")),
	)
}
```

---

## Stage 4 — `searchbar/searchbar_test.go`

No build tag (backend by default). Use `strings` and `testing` from the stdlib —
**test files are exempt from the no-stdlib rule; do not "fix" these imports.**

Write exactly these cases:

1. `TestSearchBar_RendersBarAndInput` — `(&SearchBar{}).Render().String()`
   contains `searchbar`, `searchbar__icon`, `searchbar__glyph`,
   `searchbar__input`, and `type='search'`.
2. `TestSearchBar_DefaultPlaceholder` — with `Placeholder` unset the markup
   contains `placeholder='Search…'`; with `Placeholder: "Buscar..."` it
   contains `placeholder='Buscar...'` and NOT `Search…`.
3. `TestSearchBar_IntrinsicWidthCollapsed` — the markup contains `size='1'`.
   Guards the wrap regression called out in the Render comment.
4. `TestSearchBar_OnFilterChangeIsOptional` — `Render()` on a bar with no sink
   registered must not panic.
4b. `TestSearchBar_SatisfiesFilterable` — assign to a
   `var f widget.Filterable = &SearchBar{}`, register a sink through it, and
   assert the sink is the one the bar stores. This is the consumer-shaped proof
   the harness requires: the bar must be usable **through the interface**, not
   only through its concrete type.
5. `TestSearchBar_RenderCSSEmitsEveryPart` — `(&SearchBar{}).RenderCSS().String()`
   contains `.searchbar `, `.searchbar__icon`, `.searchbar__glyph` and
   `.searchbar__input`, and the `.searchbar ` and `.searchbar__input` blocks
   each contain `min-height: var(--control-height`.

> Attribute quoting: `tinywasm/dom` renders attributes with **single quotes**
> (`class='…'`). Assert with single quotes or the tests silently never match.

---

## Stage 5 — register the component in the module's conformance lists

Edit `conformance_test.go` at the module root — three edits, all mechanical:

1. Add to the import block, in alphabetical position (after `selectsearch`):
   ```go
   "github.com/tinywasm/components/searchbar"
   ```
2. In `TestEveryPackageEmits`, add `&searchbar.SearchBar{},` to the slice —
   alphabetically, after `&selectsearch.SelectSearch{},`.
3. In `TestKindAllowsEveryState`, add `"searchbar": &searchbar.SearchBar{},` to
   the map — after the `"selectsearch"` entry.

⚠️ Both lists must be updated. Adding to only one leaves the new package's
`css.go` unparsed by the state check and the omission is silent.

---

## Stage 6 — documentation

### `searchbar/README.md`

Follow the shape of `selectsearch/README.md`. It MUST cover:

- What it is: a one-control filter bar (magnifier cap + text field), sized by
  `--control-height`.
- Fields: `Placeholder`.
- That it implements `widget.Filterable`, and the sentence that makes it
  reusable: *the bar reports the term and nothing else, so any control that can
  produce a term — a calendar, a select — can replace it in the same host slot
  without the host changing.*
- A usage snippet:
  ```go
  bar := &searchbar.SearchBar{Placeholder: "Buscar..."}
  bar.OnFilterChange(func(term string) { /* filter your list */ })
  ```
- Parts table: `searchbar` / `searchbar__icon` / `searchbar__glyph` /
  `searchbar__input`.

### `docs/CATALOG.md`

Insert a new entry between `SelectSearch` and `ThemeToggle`, matching the
existing format exactly (`---` separator, `— ✅ Slot-ready`, detail link):

```markdown
## [SearchBar](../searchbar/README.md) — ✅ Slot-ready
One-control filter bar: a magnifier cap and a text field sized by
`--control-height`. Reports each keystroke through `OnFilterChange(term)` and
knows nothing about what it filters, so a host can swap it for a calendar or a
select in the same slot.
[Detailed Documentation →](../searchbar/README.md)

---
```

---

## Stages table

| # | File(s) | What lands |
|---|---|---|
| 1 | `searchbar/searchbar.go` (new) | struct, parts, `Render`, `OnFilterChange`, widget identity |
| 2 | `searchbar/css.go` (new, `!wasm`) | `RenderCSS` — bar/cap/glyph/input rules moved from crudview |
| 3 | `searchbar/svg.go` (new, `!wasm`) | `IconSvg` — magnifier moved from crudview |
| 4 | `searchbar/searchbar_test.go` (new) | 5 cases listed in Stage 4 |
| 5 | `conformance_test.go` | import + both component lists |
| 6 | `searchbar/README.md` (new), `docs/CATALOG.md` | docs |

Stages 1–3 must land together (2 and 3 reference identifiers from 1). Stages 4–6
are independent of each other.

---

## Definition of done

1. `gotest` green at the module root — this runs `vet`, `race`, the WASM build
   and the conformance suite.
2. `grep -rn "icon-crud-search" .` → empty (the id did not travel).
2b. `grep -rn "interface{ OnFilterChange\|FilterSource" searchbar/` → empty:
   the interface is used from `widget`, never redeclared here.
3. The three new files exist and carry the right build tags: `searchbar.go`
   none, `css.go` and `svg.go` both `//go:build !wasm`.
4. `docs/CATALOG.md` lists `SearchBar`.

## Out of scope

- **Any change to `crudview` or anything under `github.com/tinywasm/layout`.**
  Different repository, `layout/docs/PLAN.md`, runs after this one.
- **Declaring the `Filterable` interface here.** It lives in `widget`
  (`widget/docs/PLAN.md`, the gate). This package implements it and nothing more.
- **Removing `SearchPlaceholder()` from `github.com/tinywasm/view`.** It is
  still part of that interface; the `layout` plan stops calling it, and
  retiring it is a `view` decision tracked there.
- **A debounce on the input.** The current bar filters an in-memory slice on
  every keystroke and that is fast enough; adding a timer here would change
  behaviour under cover of a move.
- **Clear ("×") button, keyboard shortcuts, or search history.** New features,
  not an extraction.
