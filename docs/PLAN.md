---
PLAN: "components: rename the SSR style entry point from Style() to RenderCSS()"
EXECUTOR: jules
STATUS: running
SESSION: 14427496519939977001
---

> This plan is dispatched via the CodeJob workflow. See skill: **agents-workflow**.

# Plan — `tinywasm/components`: one CSS entry point, `RenderCSS()`

Eight packages currently expose their visual sheet as `Style() *style.Sheet`. The rest of the
ecosystem — `tinywasm/css`, `tinywasm/layout`, every app module — exposes it as
`RenderCSS() *css.Stylesheet`. This plan collapses that into one.

**Only the signature changes. The `style` DSL stays exactly as it is** — every rule body is
copied verbatim; a `.Stylesheet()` call is appended to the chain.

---

## 🚦 0. Blocking gate — do not start without this

This plan requires `github.com/tinywasm/ssr` to be **published** without its `Style()`
detection branch (see
[`ssr/docs/PLAN.md`](https://github.com/tinywasm/ssr/blob/main/docs/PLAN.md)).

Until that ships, the released extractor **hard-errors** on any package that imports
`github.com/tinywasm/widget/style` and does not declare a `Style()` method — which is exactly
the state this plan produces. Renaming first would break asset extraction for every consumer.

**Mandatory check before stage 1:**

```bash
go list -m -versions github.com/tinywasm/ssr
```

Take the newest version and confirm the branch is gone:

```bash
go doc github.com/tinywasm/ssr 2>/dev/null
grep -rn "widgetStylePkg\|HasStyle" "$(go env GOMODCACHE)"/github.com/tinywasm/ssr@*/invoke.go
```

That `grep` must print **nothing** for the newest `ssr@*` directory. If it prints matches, or if
no `ssr` version newer than `v0.0.24` exists, **stop and report it**. Do not add a local
`replace` to work around the gate.

---

## ⚠️ 1. Scope — read this before touching anything

Eight `css.go` files, their tests, `go.mod`, the root conformance test, and `docs/SKILL.md`. One
change, no staged waits; the gate is a single `gotest` green at the end.

**FORBIDDEN — do not do any of this:**

| Prohibition | Reason |
|---|---|
| Changing any styling rule | This is a rename. Every `style.On(...)`, `style.Pad(...)`, `Part(...)`, `When(...)`, `Cue(...)` argument stays byte-identical. A visual diff is a bug. |
| Dropping the `style` DSL for hand-written CSS | The DSL is the typed harness and it stays. Only the method name and return type change. |
| Dot-importing `github.com/tinywasm/css` | Use a named import. A dot-import would collide with the free function `css.RenderCSS()` and force contortions — that collision is the reason the old `AGENTS.md` rule existed. |
| Removing the `github.com/tinywasm/widget/style` import | It is still needed: `style.Of(...)` builds the sheet. |
| Touching `//go:build !wasm` | Every `css.go` and `svg.go` keeps it. `css` and `widget/style` must never reach the WASM binary. |
| Adding a `Style()` shim that calls `RenderCSS()` (or the reverse) | That recreates the two paths this plan deletes. |
| Renaming component types, `Name*` constants, or `widget.Part*` values | Out of scope. |
| Using `go test` | This repo uses `gotest`. |

---

## 2. The transformation — exactly one shape

Current (`fieldset/css.go`):

```go
//go:build !wasm

package fieldset

import (
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// Style defines the fieldset widget visual contract using the style DSL.
func (f *Fieldset) Style() *style.Sheet {
	return style.Of(widget.NameField).
		Root(
			style.On(style.Panel),
			style.Round(style.RadiusMd),
		).
		Part(widget.PartLabel,
			style.On(style.Accent),
		)
}
```

After:

```go
//go:build !wasm

package fieldset

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS defines the fieldset widget visual contract using the style DSL.
func (f *Fieldset) RenderCSS() *css.Stylesheet {
	return style.Of(widget.NameField).
		Root(
			style.On(style.Panel),
			style.Round(style.RadiusMd),
		).
		Part(widget.PartLabel,
			style.On(style.Accent),
		).
		Stylesheet()
}
```

Three edits per file, nothing else:

1. Add `"github.com/tinywasm/css"` to the import block.
2. `func (r *T) Style() *style.Sheet` → `func (r *T) RenderCSS() *css.Stylesheet`, keeping the
   **existing receiver name** unchanged.
3. Append `.Stylesheet()` to the end of the returned chain.

Update the doc comment's first word from `Style` to `RenderCSS`; keep the rest of the sentence.

**Watch the chain terminator.** `Stylesheet()` goes on the **last** call of the builder chain,
after the final `)` of the last `Part`/`When`/`Cue`. Some of these chains are long — read to the
end of the `return` statement before editing.

---

## 3. Stage 1 — the eight `css.go` files

Each keeps its own receiver name. Apply §2 to all of them:

| File | Receiver | Signature after |
|---|---|---|
| `actionbutton/css.go` | `b *ActionButton` | `func (b *ActionButton) RenderCSS() *css.Stylesheet` |
| `contentcard/css.go` | `c *ContentCard` | `func (c *ContentCard) RenderCSS() *css.Stylesheet` |
| `datatable/css.go` | `t *DataTable` | `func (t *DataTable) RenderCSS() *css.Stylesheet` |
| `fieldset/css.go` | `f *Fieldset` | `func (f *Fieldset) RenderCSS() *css.Stylesheet` |
| `modaldialog/css.go` | `m *ModalDialog` | `func (m *ModalDialog) RenderCSS() *css.Stylesheet` |
| `selectsearch/css.go` | `c *SelectSearch` | `func (c *SelectSearch) RenderCSS() *css.Stylesheet` |
| `targetlist/css.go` | `t *TargetList` | `func (t *TargetList) RenderCSS() *css.Stylesheet` |
| `themetoggle/css.go` | `t *ThemeToggle` | `func (t *ThemeToggle) RenderCSS() *css.Stylesheet` |

`contentcard/css.go` and `modaldialog/css.go` import only `widget/style` today (no `widget`) —
they gain `css` and keep it that way. Do not add a `widget` import they do not use.

---

## 4. Stage 2 — `go.mod`

`github.com/tinywasm/css v0.3.0` is currently an **indirect** requirement. After stage 1 it is
imported directly by eight packages, so `go mod tidy` moves it into the direct `require` block.
Run `go mod tidy` and commit the result. Do not change the version.

Expected direct block afterwards:

```
require (
	github.com/tinywasm/css v0.3.0
	github.com/tinywasm/dom v0.11.4
	github.com/tinywasm/fmt v0.25.5
	github.com/tinywasm/html v0.0.6
	github.com/tinywasm/svg v0.1.8
	github.com/tinywasm/widget v0.3.0
)
```

Do **not** add a `replace` directive.

---

## 5. Stage 3 — tests

Every call site of the form `X.Style().Stylesheet()` becomes `X.RenderCSS()`. The
`.Stylesheet()` call moved into the method, so it must **not** remain at the call site.

| File | Lines (approximate) | Current | After |
|---|---|---|---|
| `actionbutton/button_test.go` | 83 | `(&ActionButton{}).Style().Stylesheet().String()` | `(&ActionButton{}).RenderCSS().String()` |
| `contentcard/card_test.go` | 108 | `c.Style().Stylesheet().String()` | `c.RenderCSS().String()` |
| `datatable/table_test.go` | 101 | `dt.Style().Stylesheet().String()` | `dt.RenderCSS().String()` |
| `fieldset/css_test.go` | 17, 51 | `(&Fieldset{}).Style().Stylesheet().String()` / `f.Style().Stylesheet().String()` | `(&Fieldset{}).RenderCSS().String()` / `f.RenderCSS().String()` |
| `modaldialog/modaldialog_test.go` | 67, 124 | `md.Style().Stylesheet().String()` | `md.RenderCSS().String()` |
| `selectsearch/selectsearch_test.go` | 144 | `ss.Style().Stylesheet().String()` | `ss.RenderCSS().String()` |
| `targetlist/targetlist_test.go` | 70, 132 | `tl.Style().Stylesheet().String()` | `tl.RenderCSS().String()` |
| `themetoggle/themeswitch_test.go` | 69, 121 | `ts.Style().Stylesheet()` / `tt.Style().Stylesheet().String()` | `ts.RenderCSS()` / `tt.RenderCSS().String()` |

Line numbers are a starting point, not an authority — resolve every occurrence with
`grep -rn "\.Style()" .` and leave none behind.

**Assertions do not change.** The emitted CSS is identical, so every `strings.Contains` check
keeps its expected substring. If an assertion starts failing, the rule body was altered in
stage 1 — fix the rule, not the assertion.

**Anti-footgun:** several tests use a local variable named `css` (`css := tl.Style()...`).
Go imports are file-scoped, so a `css` variable in `targetlist_test.go` does **not** collide with
the `css` package imported by `targetlist/css.go`. Leave those variable names alone — renaming
them is churn.

`themetoggle/themeswitch_test.go:67` is named `TestRenderCSS_NotEmpty` and
`fieldset/css_test.go:16` is named `TestRenderCSS_StylesFieldset`. Both names become accurate
again; keep them.

---

## 6. Stage 4 — root conformance test

`conformance_test.go` is an AST walker over the repo. Add one check to it so the old shape can
never come back:

> For every package directory containing a `css.go`, the file declares **exactly one** method
> named `RenderCSS` and **no** method named `Style`.

Fail with a message that states the file and what was found, e.g.
`components/foo/css.go: declares Style(); the SSR CSS entry point is RenderCSS() *css.Stylesheet`.

Resolve `css.go` by walking the repo the same way the existing checks in that file do — reuse
its helpers rather than adding a second directory walk. Detect the method through the AST
(`*ast.FuncDecl` with a non-nil `Recv`), never by matching selector text.

---

## 7. Stage 5 — `docs/SKILL.md`

This file teaches the pattern, so it currently teaches the wrong one. Update:

| Line | Current | After |
|---|---|---|
| 3 (`description:`) | `Pattern based on style.Styler and widget package visual-contract.` | `Pattern based on RenderCSS() and the widget package visual-contract.` |
| 14 | `…is done by implementing style.Styler using the typed style.Sheet DSL.` | `…is done by declaring RenderCSS() *css.Stylesheet, built with the typed style.Sheet DSL.` |
| 40 | `├── css.go           # !wasm only: Style() *style.Sheet visual sheet` | `├── css.go           # !wasm only: RenderCSS() *css.Stylesheet visual sheet` |
| 82 | `…and optional style.Styler in a tagged !wasm file…` | `…and optionally declares RenderCSS() in a tagged !wasm file…` |
| 95 | `func (l *TargetList) Style() *style.Sheet {` | `func (l *TargetList) RenderCSS() *css.Stylesheet {` — and the example's chain gains `.Stylesheet()` plus the `css` import |

Any other `Styler` / `Style()` mention in `docs/SKILL.md`, `docs/CATALOG.md`,
`docs/DOCUMENTATION.md` or a component `README.md` gets the same treatment. Find them with
`grep -rn "Styler\|Style()" docs/ */README.md README.md`.

**Do not** change `fieldset/fieldset.go:9` — its comment already says `RenderCSS()` and becomes
correct on its own.

---

## 8. Acceptance criteria — grep-verifiable

1. `gotest` green.
2. `grep -rn "func .*) Style()" --include='*.go' .` → **empty**.
3. `grep -rn "\.Style()\|Styler" --include='*.go' .` → **empty**.
4. `grep -rn "Styler\|Style()" docs/ *.md */README.md` → **empty**. Both spellings, not just
   `Styler`: `docs/SKILL.md` carries each of them in different places (§7).
5. `grep -rln "func .*) RenderCSS() \*css.Stylesheet" --include='css.go' .` → **exactly 8 files**.
6. `grep -rn "Stylesheet()" --include='*_test.go' .` → **empty** (the call moved into the method).
7. `grep -c "go:build !wasm" */css.go` → **1 for each of the 8 files**.
8. `GOOS=js GOARCH=wasm go list -deps ./...` contains **neither** `github.com/tinywasm/css` **nor**
   `github.com/tinywasm/widget/style`.
9. `grep -nE '^[[:space:]]*replace' go.mod` → **empty**.
10. The emitted CSS is unchanged: for each component, the string produced by `RenderCSS()` equals
    what `Style().Stylesheet()` produced before. Verified by every existing assertion in stage 3
    passing **without editing its expected substrings**.

---

## 9. Go quality checklist (mandatory)

- No repeated string literals: classes still derive from `widget.Name` via `.Root()` /
  `.Class(Part)`. No new `"…"` class literal appears anywhere.
- Errors via `github.com/tinywasm/fmt`, never stdlib `errors`/`fmt`.
  **Anti-footgun:** `conformance_test.go` is a `!wasm` test that already imports stdlib
  `fmt`, `go/ast`, `go/parser`, `go/token`. That is legitimate for repo tooling — do **not**
  "fix" those imports.
- `//go:build !wasm` preserved on every `css.go` and `svg.go`.
- Zero `any`, zero `map` in new API.
- No color literal (`#rrggbb`, `rgb(`, `hsl(`), no viewport unit (`vw`/`vh`), no `Media(`, no
  `RawRule(`, no `Str(` is introduced. The existing conformance test already enforces this —
  it must keep passing untouched.

---

## 10. Stages table

| # | Stage | Files | Gate |
|---|---|---|---|
| 0 | *(gate)* `ssr` published without the `Style()` branch | — | `grep` in §0 prints nothing |
| 1 | Rename in the eight `css.go` | `*/css.go` | `go build ./...` |
| 2 | `go.mod` | `go.mod`, `go.sum` | `go mod tidy` clean |
| 3 | Test call sites | 8 `*_test.go` | compiles |
| 4 | Conformance guard | `conformance_test.go` | `gotest` green |
| 5 | Docs | `docs/SKILL.md`, other `.md` | §8.4 empty |

Sequential. Stage 4 is the real gate.

---

## 11. Downstream — informational, not this agent's work

Once this is published, [`tinywasm/layout`](https://github.com/tinywasm/layout) can run its own
visual-contract migration; its `docs/PLAN.md` gate points at the version published from this
plan. `tinywasm/widget` deletes the now-unused `style.Styler` interface under its own plan. Do
not attempt either from this repo.
