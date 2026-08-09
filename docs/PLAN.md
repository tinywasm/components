---
PLAN: "feat(selectsearch): satisfy widget.Filterable + fix mobile/desktop dropdown anchoring"
TAG: v0.5.0
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 7226293707427292253
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
>
> **Stage A of a 2-repo change.** `tinywasm/layout` has a follow-up plan
> (`layout/docs/PLAN.md`) that adds a demo module built on top of what this
> plan produces, and depends on this plan being merged and published
> (`github.com/tinywasm/components` bumped to `v0.5.0`) first.

# Plan — `components`: `selectsearch` becomes a real `Filterable` control

## 0. Why

`SelectSearch` (`tinywasm/components/selectsearch`) is going to become the
standard "search and pick one" control used across several product modules
(searching doctors/patients, filtering an agenda, narrowing a list by the
picked item). Today it falls short of that role in two concrete, verified
ways:

1. **It does not satisfy `widget.Filterable`.** `github.com/tinywasm/widget`
   defines (`widget/capability.go:26`):
   ```go
   type Filterable interface{ OnFilterChange(func(term string)) }
   ```
   with the doc comment explicitly naming this control's use case: *"a select
   emits the chosen id"*. `searchbar.SearchBar` implements this
   (`searchbar/searchbar.go:54,68`); `SelectSearch` has no `OnFilterChange`
   method and no `var _ widget.Filterable = (*SelectSearch)(nil)` assertion.
   Concretely, this means `SelectSearch` **cannot** be dropped into a host
   slot typed to accept any `Filterable` control (e.g.
   `tinywasm/layout/crudview.Config.Filter`, which type-asserts for
   `widget.Filterable` at `crudview.go:110`) — it renders but never wires up.

2. **Its dropdown has no responsive treatment.** `selectsearch/css.go`
   applies `style.Raise(style.Floating)` to `PartDropdown`
   **unconditionally** (same rule at every viewport) and its `Root` has no
   `style.Anchor()` — so the floating panel has no positioned ancestor to
   hang from, and there is no mobile-specific override at all. Compare to
   `github.com/tinywasm/components/usermenu` (`usermenu/css.go`), the
   **proven working reference** for this exact pattern: `Root(Anchor())`,
   the flyout part is a **plain in-flow block by default** (mobile: no
   `Raise`/`Flyout` at all — it behaves like an accordion, no overlay), and
   `Raise(Floating)` + `Flyout(side)` are added **only** `On(css.Tablet,
   ...)` / `On(css.Desktop, ...)`. `selectsearch` never adopted this
   escalation — it floats (badly, unanchored) at every breakpoint including
   the phone, where there is no room to spare for a mispositioned overlay.

This plan fixes both, using `usermenu` as the literal template for the CSS
fix (do not invent a different mechanism).

---

## 1. `selectsearch/selectsearch.go` — implement `widget.Filterable`

Read the current file in full before editing
(`github.com/tinywasm/components/selectsearch/selectsearch.go`).

### 1a. Add the sink field and setter

Add to the `SelectSearch` struct (after the existing `OnSearch` field):

```go
	// Internal state signals
	selectedLabel *SignalString
	query         *SignalString
	isOpen        *SignalBool
	rows          *SignalNodes

	onFilter func(term string) // set via OnFilterChange — satisfies widget.Filterable
```

Add the import (the file already imports `"github.com/tinywasm/widget"` for
`widget.Name`/`widget.Kind` — reuse that import, do not add a second one):

Add the method, placed right after `WidgetKind()`:

```go
// OnFilterChange implements widget.Filterable: it registers the sink called
// with the picked option's ID whenever a selection is made. This is a
// SEPARATE, additive wiring path from OnSelect — OnSelect still gets
// (id, description) for a consumer that needs both; OnFilterChange exists so
// a host that only knows the generic Filterable contract (e.g.
// tinywasm/layout/crudview's Filter slot) can drop a *SelectSearch into the
// same seam a *searchbar.SearchBar fills today, with no bespoke glue.
//
// The signature is fixed by widget.Filterable — do not add a parameter, do
// not return anything, do not add a companion getter (see searchbar.go's
// OnFilterChange for the same rule stated for SearchBar).
func (c *SelectSearch) OnFilterChange(fn func(term string)) { c.onFilter = fn }
```

Add the compile-time assertion next to the existing one:

```go
func (c *SelectSearch) WidgetName() widget.Name { return NameSelectSearch }
func (c *SelectSearch) WidgetKind() widget.Kind { return widget.Combobox }

var _ widget.Filterable = (*SelectSearch)(nil)
```

### 1b. Extract `selectOption` and call both sinks from it

Today the option's click handler (inside `buildRows`) inlines the "commit a
selection" logic. Extract it to a method so `OnSelect` and the new
`onFilter` sink can never drift out of sync, and so it is unit-testable
without simulating a real DOM click event:

Replace the body of the `On("click", ...)` closure inside `buildRows`
(currently):

```go
			On("click", func(e Event) {
				c.selectedLabel.Set(o.Label)
				c.isOpen.Set(false)
				c.query.Set("")
				c.rows.Set(c.buildRows(""))
				if c.OnSelect != nil {
					c.OnSelect(o.ID, o.Description)
				}
			})
```

with:

```go
			On("click", func(e Event) { c.selectOption(o) })
```

and add the new method (place it right before `buildRows`):

```go
// selectOption is the single place an option becomes "chosen" — today only
// a mouse click reaches it, but every future input path (keyboard, a future
// OnSearch auto-pick) commits through here too, so OnSelect and the
// Filterable sink can never fire out of step with each other.
func (c *SelectSearch) selectOption(o SsOption) {
	c.selectedLabel.Set(o.Label)
	c.isOpen.Set(false)
	c.query.Set("")
	c.rows.Set(c.buildRows(""))
	if c.OnSelect != nil {
		c.OnSelect(o.ID, o.Description)
	}
	if c.onFilter != nil {
		c.onFilter(o.ID)
	}
}
```

**Acceptance for Stage 1**: `go build ./...` succeeds; `grep -n "onFilter\|OnFilterChange\|selectOption" selectsearch/selectsearch.go` shows the field, the method, the assertion, and the extracted helper.

---

## 2. `selectsearch/css.go` — anchor + responsive escalation

Read `github.com/tinywasm/components/usermenu/css.go` in full first — it is
the template this stage copies the pattern from, not just an inspiration.
Also read `github.com/tinywasm/components/targetlist/css.go` for the
`RenderCSS()`/`sheet()` split used in Stage 3 below.

Replace `selectsearch/css.go` entirely with:

```go
//go:build !wasm

package selectsearch

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the selectsearch visual contract using the style DSL.
func (c *SelectSearch) RenderCSS() *css.Stylesheet {
	return c.sheet().Stylesheet()
}

// sheet is split out from RenderCSS so tests can call Validate() on the
// *style.Sheet directly, without the *css.Stylesheet conversion — the same
// split targetlist/css.go uses (see targetlist_test.go's TestSheetValidates).
func (c *SelectSearch) sheet() *style.Sheet {
	return style.For(c).
		// Anchor: what lets PartDropdown's Flyout (added below, tablet/desktop
		// only) hang from THIS box instead of some positioned ancestor further
		// up the tree. Mirrors usermenu's Root(Anchor()) exactly — same
		// primitive, same reason ("a corner pin needs a positioned ancestor to
		// hang from").
		Root(
			style.Anchor(),
			style.Stack(style.Space1),
			style.As(style.Panel),
			style.Round(style.RadiusMd),
		).
		Part(PartToggle,
			style.As(style.Panel),
		).
		// Mobile-first base: the dropdown is a plain IN-FLOW block here — no
		// Raise, no Flyout. Same choice usermenu's PartPanel makes for a phone:
		// there is no room to spare for an anchored floating overlay on a
		// narrow viewport, so opening the control expands it in place instead
		// (an accordion, not a popover). Raise(Floating)+Flyout are added ONLY
		// from Tablet up, in the On(...) blocks below — do not add them here.
		Part(PartDropdown,
			style.Stack(style.Space1),
			style.As(style.Panel),
			style.HideOverflow(),
		).
		Part(PartHeader,
			style.Row(style.Space2),
			style.Pad(style.Space2),
		).
		Part(PartIcon,
			style.As(style.Panel),
		).
		Part(PartSearch,
			style.Pad(style.Space2),
			style.As(style.Inset),
		).
		Part(PartOptions,
			style.Stack(style.SpaceNone),
			style.Scroll(),
		).
		Part(PartOption,
			style.Row(style.Space2),
			style.Pad(style.Space2),
			style.Interactive(style.Panel),
		).
		Part(PartLabel,
			style.As(style.Panel),
		).
		Part(PartDesc,
			style.As(style.Inset),
			style.FontSize(style.TextXs),
		).
		// From a tablet up, the dropdown becomes a floating panel hanging off
		// the Root anchor — the identical escalation usermenu's PartPanel uses,
		// gated on the identical two breakpoints, for the identical reason.
		On(css.Tablet, PartDropdown,
			style.Raise(style.Floating),
			style.Flyout(style.SideStart),
		).
		On(css.Desktop, PartDropdown,
			style.Raise(style.Floating),
			style.Flyout(style.SideStart),
		)
}
```

**Do not** add a `Docked(Viewport, ...)` mobile treatment — that mechanism
was deliberately removed from this exact class of problem elsewhere in this
repo (see `docs/LAST_PLAN_EXECUTED.md`'s Stage 2 note: *"The mobile
`Docked(Viewport, …)` that used to be here was the escape hatch"* for
`targetlist`). The correct mobile behavior is "stay in flow", not "become a
viewport-fixed sheet" — that is what the base (untagged) `Part(PartDropdown,
...)` rule above already gives you by simply not declaring `Raise`/`Flyout`
there.

**Acceptance for Stage 2**: `go build ./...` succeeds.

---

## 3. Tests

### 3a. `selectsearch/selectsearch_test.go` — add

Add these two tests to the existing file (do not remove
`TestSelectSearch_InitTwiceSafe` / `TestSelectSearch_RenderIdempotent`):

```go
func TestSelectSearch_SatisfiesFilterable(t *testing.T) {
	var _ widget.Filterable = (*SelectSearch)(nil)
}

func TestSelectSearch_OnFilterChange_FiresOnSelection(t *testing.T) {
	c := &SelectSearch{Options: []SsOption{
		{ID: "p1", Label: "Juan Pérez", Description: "09:00"},
	}}
	c.Init(nil)

	var got string
	fired := 0
	c.OnFilterChange(func(term string) {
		got = term
		fired++
	})

	c.selectOption(c.Options[0])

	if fired != 1 {
		t.Fatalf("expected OnFilterChange to fire exactly once, got %d", fired)
	}
	if got != "p1" {
		t.Fatalf("expected OnFilterChange term %q, got %q", "p1", got)
	}
}

func TestSelectSearch_OnSelect_StillFiresAlongsideFilterable(t *testing.T) {
	c := &SelectSearch{Options: []SsOption{
		{ID: "p1", Label: "Juan Pérez", Description: "09:00"},
	}}
	c.Init(nil)

	var gotID, gotDesc string
	c.OnSelect = func(id, description string) {
		gotID, gotDesc = id, description
	}
	c.OnFilterChange(func(string) {}) // both sinks registered — neither must suppress the other

	c.selectOption(c.Options[0])

	if gotID != "p1" || gotDesc != "09:00" {
		t.Fatalf("expected OnSelect(%q, %q), got (%q, %q)", "p1", "09:00", gotID, gotDesc)
	}
}
```

Add `"github.com/tinywasm/widget"` to this file's imports if not already
present.

### 3b. `selectsearch/css_test.go` (new file, or add to existing `css_test.go` if one already has content — check first)

```go
package selectsearch

import "testing"

func TestSheetValidates(t *testing.T) {
	c := &SelectSearch{}
	c.Init(nil)
	if errs := c.sheet().Validate(); len(errs) > 0 {
		t.Errorf("selectsearch sheet must validate, got:\n%v", errs)
	}
}
```

**Acceptance for Stage 3**: `gotest` (never `go test`) — all green, including the 3 new tests.

---

## 4. Documentation

### 4a. `selectsearch/README.md` — add a `## Filterable` section

Insert after the existing `## API` section, mirroring `searchbar/README.md`'s
own `## Filterable` section format exactly:

```markdown
## Filterable

`SelectSearch` implements `widget.Filterable`: picking an option calls the
registered sink with the option's `ID`. This is in addition to `OnSelect`
(which also gets `Description`) — a host that only needs the generic
narrowing contract (e.g. `tinywasm/layout/crudview`'s `Filter` slot) can drop
a `*SelectSearch` in without any bespoke wiring, the same way it accepts a
`*searchbar.SearchBar` today.
```

### 4b. `docs/CATALOG.md` — update the SelectSearch entry

Change the existing entry (`docs/CATALOG.md`, "## [SelectSearch]" section) to
mention the new capability. Replace:

```markdown
## [SelectSearch](../selectsearch/README.md) — ✅ Slot-ready
Signal-driven searchable dropdown with static options, live filtering, and optional DB search callback. Uses `BindChildren` for efficient list updates and `Show` for the dropdown.
[Detailed Documentation →](../selectsearch/README.md)
```

with:

```markdown
## [SelectSearch](../selectsearch/README.md) — ✅ Slot-ready
Signal-driven searchable dropdown with static options, live filtering, and optional DB search callback. Satisfies `widget.Filterable` — drop-in for any host slot that accepts a filter control (same seam `searchbar.SearchBar` fills). Uses `BindChildren` for efficient list updates and `Show` for the dropdown.
[Detailed Documentation →](../selectsearch/README.md)
```

**Acceptance for Stage 4**: both files updated; no other prose in either file touched.

---

## 5. Final checklist

- [ ] `selectsearch.go`: `onFilter` field, `OnFilterChange` method,
      `var _ widget.Filterable = (*SelectSearch)(nil)`, `selectOption`
      extracted and called from `buildRows`'s click handler.
- [ ] `css.go`: `Root` gains `style.Anchor()`; `PartDropdown`'s base rule has
      **no** `Raise`/`Flyout`; both added only via `On(css.Tablet, ...)` and
      `On(css.Desktop, ...)`; file split into `RenderCSS()` + `sheet()`.
- [ ] New/updated tests in `selectsearch_test.go` and `css_test.go`, all
      passing under `gotest`.
- [ ] `README.md` and `docs/CATALOG.md` updated per Stage 4.
- [ ] `go build ./...` succeeds; `gotest` all green.
- [ ] Manual browser check (via the demo at `selectsearch/web/client.go` —
      `tinywasm` dev server, no changes needed to that file for this plan):
      at a phone width, the dropdown opens in-flow (pushes the page, no
      floating panel); at tablet/desktop width, it floats anchored directly
      under the toggle, not offset or clipped.

## Stages

| Stage | File(s) | Depends on |
|---|---|---|
| 1 | `selectsearch/selectsearch.go` | — |
| 2 | `selectsearch/css.go` | — (independent of 1) |
| 3 | `selectsearch/selectsearch_test.go`, `selectsearch/css_test.go` | 1, 2 |
| 4 | `selectsearch/README.md`, `docs/CATALOG.md` | 1, 2 |
| 5 | Final checklist + `gotest` + manual browser check | 3, 4 |
