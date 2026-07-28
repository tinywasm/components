---
PLAN: "components: migrate the eight component sheets to widget v0.4.0"
EXECUTOR: unassigned
STATUS: draft
---

> This plan is dispatched via the CodeJob workflow. See skill: **agents-workflow**.

# Plan — `tinywasm/components`: migrate to `widget` v0.4.0

`tinywasm/widget` v0.4.0 is a closed-API breaking release. **All eight component
packages fail to compile against it.** This plan is the mechanical part of that
migration plus the three places where it is *not* mechanical.

The upstream rename table is authoritative and is not restated here:
[`widget/docs/MIGRATION.md`](https://github.com/tinywasm/widget/blob/main/docs/MIGRATION.md).
Exact emitted output is in
[`widget/docs/SPECS.md`](https://github.com/tinywasm/widget/blob/main/docs/SPECS.md).

Every count below was measured by running
`go get github.com/tinywasm/widget@v0.4.0 && go build ./...` against this
repository, not estimated.

---

## ⚠️ 1. Scope

Eight `css.go` files, their tests, `go.mod`, `conformance_test.go`, and the docs
that teach the pattern.

**FORBIDDEN — do not do any of this:**

| Prohibition | Reason |
|---|---|
| Changing a rule's *intent* | This is a migration. A part that was a `Panel` stays a `Panel`. Where the new API changes emitted output (§4), that is upstream's decision, not a licence to restyle. |
| Adding a compatibility shim or alias | v0.4.0 is a single breaking release with no deprecation period. Recreating the old names defeats it. |
| Dropping the `style` DSL for hand-written CSS | The DSL is the typed harness and it stays. |
| Touching `//go:build !wasm` | Every `css.go` and `svg.go` keeps it. `css` and `widget/style` must never reach the WASM binary. |
| Renaming component types or `Name*` constants | Out of scope. They stay — see §2.2. |
| Reproducing an old shape with new names | Where v0.4.0 collapses several calls into one (`Interactive`, `RevealedBy`), use the collapsed form. |
| Using `go test` | This repo uses `gotest`. |

---

## 2. Mechanical rewrites

Find-and-replace, no judgement. Counts are actual occurrences across the eight
`css.go` files.

### 2.1 Renamed symbols

| Old | New | Uses |
|---|---|---|
| `style.On(…)` | `style.As(…)` | 41 |
| `style.Of(Name)` | `style.For(receiver)` — **takes the widget, not the name** | 8 |
| `style.Sunken` | `style.Inset` | 6 |
| `style.Accent` | `style.Primary` | 3 |
| `style.Space0` | `style.SpaceNone` | 4 |
| `style.Text(…)` | `style.FontSize(…)` | 4 |
| `style.Scrolls()` | `style.Scroll()` | 2 |
| `style.Fixed()` | `style.KeepSize()` | 2 |
| `style.Cover()` | `style.FillCentered()` | 1 |
| `style.Scrim()` | `style.Veil()` | 1 |
| `style.Overlay` | `style.Popover` | 1 |
| `style.Selected` (surface) | `style.Highlight` | 1 |

Unchanged, and used heavily — do not touch: `Pad` (15), `Round` (9), `Stack`,
`Row`, `Raise`, `Fill`, `Clip`, `Backdrop`, `Width`, `FontWeight`, `Space1/2/3`,
`Radius*`, `Text*`, `Weight*`, `Raised`, `Floating`, `Page`, `Panel`, `Danger`,
`Secondary`, `Parent`, `Viewport`.

### 2.2 `Of(Name)` → `For(Widget)`

The sheet now needs the widget's `Kind`, not just its `Name`: stacking and
validation both derive from it.

```go
// before
func (t *TargetList) RenderCSS() *css.Stylesheet {
    return style.Of(NameTargetList).

// after
func (t *TargetList) RenderCSS() *css.Stylesheet {
    return style.For(t).
```

The `Name*` constants stay — `Render()` still uses them to build class
attributes. They simply stop being the sheet's entry point.

**`fieldset` is safe.** It passes `widget.NameField` — a name it does not own —
but its `WidgetName()` already returns `widget.NameField`, so `style.For(f)`
produces the identical name. Verified; no special handling needed.

Some packages will no longer need their `widget` import once the name argument
goes. Let `goimports` decide; do not add an import that is unused.

### 2.3 Collapse the hover pairs into `Interactive`

Seven `Cue(widget.Hover, …)` calls exist, each pairing a base surface with the
`*Hover` twin of the same family. Those twins are unexported in v0.4.0:

```go
// before
Part(PartRow, style.Row(style.Space2), style.On(style.Panel)).
Cue(widget.Hover, PartRow, style.On(style.PanelHover)).

// after
Part(PartRow, style.Row(style.Space2), style.Interactive(style.Panel)).
```

| Family | Packages | Uses |
|---|---|---|
| `Panel` | `datatable`, `targetlist`, `selectsearch` | 3 |
| `Primary` (was `Accent`) | `actionbutton`, `themetoggle` | 2 |
| `Secondary` | `actionbutton` | 1 |
| `Danger` | `actionbutton` | 1 |

**This is a deliberate gain, not a like-for-like swap.** `Interactive` emits
hover **and** focus-visible **and** press; the old code had hover only. The focus
states are new — that is the accessibility improvement the upstream release
exists to deliver — but look at them once rather than merging blind.

`Interactive` is rejected on `Page` and `Inactive`. No component uses either
interactively, so nothing here is affected.

### 2.4 `Above()` is deleted

`modaldialog` calls it once. Stacking now derives from `Kind`: the package
declares `widget.Dialog`, which resolves to `--z-modal`. Delete the call and
change nothing else.

This is also a fix. The old hardcoded `z-index: 101` sat *below* a sticky element
at `--z-sticky: 200`, so a sticky header could cover an open modal.

---

## 3. ⚠️ Not mechanical — `targetlist` will panic

**A compile-clean migration will not catch this.** After the renames,
`targetlist` builds and then aborts at emission. Reproduced against v0.4.0:

```
VALIDATE: sheet targetlist: part "panel": state Open is not meaningful for kind Listbox
PANIC:    widget/style: sheet targetlist: part "panel": state Open is not meaningful for kind Listbox
```

`Kind.Allows()` permits `Selected` and `Current` for a `Listbox` (plus the
universal `Disabled`/`Locked`/`Busy`). `Open` belongs to `Menu`, `Dialog`,
`Disclosure` and `Combobox`.

`targetlist` declares `widget.Listbox` and uses `Open` to show and hide a
backdrop. Because `ssr` calls `RenderCSS()` at build time, this aborts asset
extraction for the **whole application**, not just this package.

**Three ways out. Decide before starting — it changes the diff.**

| Option | What it means | Cost |
|---|---|---|
| **A. Change the `Kind` to `Combobox`** | `Combobox` allows `Open`, `Selected` and `Invalid`. A list with a search field that expands over a backdrop *is* a combobox in WAI-ARIA terms. | `Role()` goes from `listbox` to `combobox`; confirm the markup's ARIA attributes agree |
| **B. Split the widget** | The overlay panel becomes its own `Kind` — `Disclosure` or `Dialog` — and the list stays a `Listbox`. | Two sheets, two `Name`s. Truest to the anatomy, largest change |
| **C. Widen `Allows` upstream** | Let `Listbox` accept `Open`. | Another `widget` release, and it is the wrong fix: ARIA gives the open/closed state to the combobox that *owns* a listbox, not to the listbox |

**Recommendation: A.** `selectsearch` already declares `widget.Combobox` for the
same shape, so the suite already treats "list plus expanding panel" as a
combobox. Avoid C — it would widen the very table that caught this bug.

Whichever is chosen, the `Hidden()`/`Shown()` pair also collapses:

```go
// before — a pair split across two rules, with an ordering rule to remember
Part(PartBackdrop, style.Backdrop(style.Viewport), style.Hidden()).
When(widget.Open, PartBackdrop, style.Shown())

// after — one call
Part(PartBackdrop, style.Backdrop(style.Viewport), style.RevealedBy(widget.Open))
```

**The other seven are clean.** Measured `Kind`/state pairs:

| Package | `Kind` | states used | verdict |
|---|---|---|---|
| `fieldset` | `Form` | `Invalid`, `Locked` | OK — `Invalid` allowed, `Locked` universal |
| `targetlist` | `Listbox` | `Open`, `Selected` | **`Open` rejected** |
| the other six | `Region` / `Grid` / `Dialog` / `Combobox` | none | OK |

---

## 4. ⚠️ Not mechanical — two silent appearance changes

Neither produces a compile error.

### 4.1 Surfaces now carry a radius

In v0.4.0 a surface resolves background, text, border **and radius**. Every
`As(Panel)` gains `border-radius: var(--radius-md)`; `As(Primary)`, `As(Inset)`
and the rest gain `--radius-sm`, unless the rule already overrides it.

Measured exposure — surface applications versus explicit `Round()` calls:

| Package | surfaces | explicit `Round()` |
|---|---|---|
| `targetlist` | 10 | 2 |
| `selectsearch` | 8 | 1 |
| `actionbutton` | 6 | 1 |
| `fieldset` | 6 | 2 |
| `datatable` | 3 | **0** |
| `contentcard` | 2 | 1 |
| `modaldialog` | 2 | 1 |
| `themetoggle` | 2 | 1 |

About thirty rules gain a corner radius they did not have. Most will look right —
that is the point of the change — but **`datatable` has no explicit radius
anywhere and its cells will now round**. Add `style.Round(style.RadiusNone)`
wherever square is intended, and look at a rendered table before merging.

Padding is unaffected: it was never folded into surfaces, and the fifteen `Pad()`
calls stay exactly as they are.

### 4.2 Focus rings arrive globally

`tinywasm/css` now owns a single global `:focus-visible` rule using
`--color-focus-ring`, and `Interactive()` adds a per-family focus background.
Components that previously showed no focus affordance will start showing one.
That is the intended fix — verify it looks deliberate rather than reverting it.

### 4.3 The palette changed

Step 1 pulls a newer `css` transitively. Its contrast-corrected palette changes
several colours — `--color-primary`, `--color-success` and `--color-error` among
them — because the previous values failed WCAG AA at the colours that actually
rendered. Expect a visual diff on brand colours and do not treat it as a
regression.

---

## 5. `actionbutton/button.go`

The only non-`css.go` file that mentions the migrated API:

```go
var variantCls widget.Class
```

`widget.Class` still exists in v0.4.0 with the same shape, so this compiles
unchanged. **No action required** — listed only so it is not mistaken for an
oversight during review.

---

## 6. Implementation order

One package per commit, so a visual regression is bisectable.

| # | Stage | Files | Gate |
|---|---|---|---|
| 0 | **Decide §3** — the `targetlist` `Kind` | — | recorded in the PR description |
| 1 | Bump | `go.mod`, `go.sum` | `go mod tidy` clean |
| 2 | The five straightforward packages: `contentcard`, `themetoggle`, `actionbutton`, `datatable`, `fieldset` | 5 `css.go` | each compiles |
| 3 | `modaldialog` — adds `Above()` deletion and the `Cover`/`Scrim`/`Fixed`/`Overlay` renames | `modaldialog/css.go` | compiles |
| 4 | `selectsearch` — largest surface count | `selectsearch/css.go` | compiles |
| 5 | `targetlist` — the §3 decision plus `RevealedBy` | `targetlist/css.go` | **emits without panicking** |
| 6 | Test call sites | 8 `*_test.go` | compiles |
| 7 | `conformance_test.go` — add the §7 guards | root | `gotest` green |
| 8 | Docs | `docs/SKILL.md`, `docs/CATALOG.md`, `*/README.md` | §8.5 empty |

Stage 1 is `go get github.com/tinywasm/widget@v0.4.0 && go mod tidy`. Do not add
a `replace` directive.

Stage 5 is the real gate: it is the only stage whose failure mode is a panic
rather than a compile error.

---

## 7. Test strategy

| Test | Asserts |
|---|---|
| `TestEveryPackageEmits` | every package's `RenderCSS()` runs on a **zero value** without panicking. This is what catches §3; a compile-only check does not. |
| `TestKindAllowsEveryState` | table-driven: for each package, every state passed to `When()` satisfies its own `Kind.Allows()`, so §3 cannot recur when a component gains a state |
| `TestNoRemovedSymbols` | no `css.go` contains `style.On(`, `style.Of(`, `Hidden()`, `Shown()`, `Above()`, `Scrim()`, `Cover()`, `Fixed()`, `Scrolls()`, `style.Accent`, `style.Sunken` |
| `TestNoHoverCuePairs` | no `Cue(widget.Hover, …)` remains — all seven collapsed into `Interactive` |
| existing `conformance_test.go` | extend, do not replace. Its no-colour-literal / no-viewport-unit checks must keep passing untouched. |

`TestEveryPackageEmits` matters most. The migration's failure mode is *compiling
successfully and then aborting extraction*, and only calling `RenderCSS()` on
each package finds it.

Existing assertions in the per-package tests **will** need their expected
substrings updated where §4 changes the output — that is the one place this
migration legitimately edits an assertion, unlike the previous plan. Change the
expectation only after confirming the new output is what §4 predicts.

---

## 8. Acceptance criteria — grep-verifiable

1. `gotest` green.
2. `go build ./...` → clean. The failing-package count goes **8 → 0**.
3. `grep -rn "style\.On(\|style\.Of(\|style\.Accent\|style\.Sunken\|Hidden()\|Shown()\|Above()\|Scrim()\|Cover()\|style\.Fixed()\|Scrolls()" --include='*.go' .` → **empty**.
4. `grep -rn "Cue(widget.Hover" --include='*.go' .` → **empty**.
5. `grep -rn "style\.Of\|style\.On\|Styler" docs/ *.md */README.md` → **empty**.
6. `grep -rl "style.For(" --include='css.go' .` → **exactly 8 files**.
7. `grep -c "go:build !wasm" */css.go` → **1 for each of the 8**.
8. `GOOS=js GOARCH=wasm go list -deps ./...` contains **neither** `github.com/tinywasm/css` **nor** `github.com/tinywasm/widget/style`.
9. `grep -nE '^[[:space:]]*replace' go.mod` → **empty**.
10. `go.mod` requires `github.com/tinywasm/widget v0.4.0` and no `v0.3.x`.

---

## 9. Go quality checklist (mandatory)

- No repeated string literals: classes still derive from `widget.Name` via
  `.Root()` / `.Class(Part)`. No new `"…"` class literal appears anywhere.
- Errors via `github.com/tinywasm/fmt`, never stdlib `errors`/`fmt`.
  **Anti-footgun:** `conformance_test.go` is a `!wasm` test that already imports
  stdlib `fmt`, `go/ast`, `go/parser`, `go/token`. That is legitimate repo
  tooling — do **not** "fix" those imports.
- `//go:build !wasm` preserved on every `css.go` and `svg.go`.
- Zero `any`, zero `map` in new API.
- No colour literal (`#rrggbb`, `rgb(`, `hsl(`), no viewport unit (`vw`/`vh`) is
  introduced. The existing conformance test enforces this and must keep passing
  untouched.

---

## 10. Coordination

- **`tinywasm/widget` v0.4.0** — published. Nothing pending.
- **`tinywasm/css`** — pulled transitively by stage 1. The palette change is
  intended (§4.3).
- **`tinywasm/ssr`** — has an unshipped plan (E-7) to recover per producer and
  name the offending package on panic. Until it lands, a §3-style panic surfaces
  as a stack trace inside generated code. **Worth landing first if this migration
  is executed by an agent**, since it turns a confusing failure into a named one.
- **`tinywasm/form`** — the other `widget` consumer. `NameField` and its parts are
  unchanged, but any sheet it owns needs this same migration. Not in scope here;
  check it before tagging the suite.

---

## 11. Note on the previous plan

The plan this file replaces — `Style()` → `RenderCSS()` — is **complete**. Its
acceptance criteria were re-verified before overwriting: no `Style()` method
remains, no `.Style()` or `Styler` reference remains in any `.go` file, and all
eight `css.go` declare `RenderCSS() *css.Stylesheet`. Its only surviving mentions
were inside the plan document itself.
