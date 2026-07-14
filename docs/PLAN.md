# PLAN — `tinywasm/components`: SVG icon harness migration (selectsearch)

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Master plan: https://github.com/tinywasm/tinywasm/blob/main/docs/SVG_ICON_HARNESS_MASTER_PLAN.md
> Repo rules: `AGENTS.md` at this repo's root — read it first (especially the
> "Build tags belong to the consumer" and "SVG icons" sections).
>
> **Queue note:** `docs/CHECK_PLAN.md` (slot-ready catalog) is also pending in
> this repo. THIS plan runs first — it is small, mechanical, and touches
> `selectsearch`, which CHECK_PLAN also modifies. Do not execute CHECK_PLAN.
>
> **GATE:** requires `tinywasm/svg` already split into `svg` (Icon reference)
> + `svg/sprite` (definition) — published. Update `go.mod` to that version
> first; if `github.com/tinywasm/svg/sprite` does not exist as an import path,
> STOP and report.

## Context (zero-context summary)

This ecosystem is isomorphic Go: the same packages compile to backend (Go) and
browser (TinyGo → WASM). **Every byte of the WASM binary counts.**

`tinywasm/svg` is split in two packages (NOT by build tag inside the library —
the library itself never uses `//go:build`; the consumer decides what to
import and tags its own files):

- `github.com/tinywasm/svg` — shared reference: `const iconX = svg.Icon("comp-x")`,
  rendered with `iconX.Render(class)` → `<svg aria-hidden="true" focusable="false" class="..."><use href="#comp-x"/></svg>`.
  This is the ONLY thing a WASM-reachable file may import from this library.
- `github.com/tinywasm/svg/sprite` — backend-only definition:
  `sprite.Define(iconX, viewBox, sprite.Path(...))`, returned by
  `IconSvg() *sprite.Sprite` from a `//go:build !wasm` file. The SSR pipeline
  extracts and injects the sprite inline into `<body>`.

`components/selectsearch` violates this today:

1. `selectsearch/svg.go` has NO build tag → the path data and the (old,
   unsplit) sprite machinery compile into the WASM binary.
2. `selectsearch/selectsearch.go:85` builds the reference by hand with a raw
   string: `svg.Svg().Child(svg.Use().Attr("href", "#ss-arrow-down")).Set(ClsSsIcon.AsAttr())`.
   `svg.Svg()`/`svg.Use()` no longer exist.

## Stages

### Stage 1 — typed reference in `selectsearch/selectsearch.go`

Next to the existing typed class constants (around line 16, where
`ClsSsIcon Class = "ss-icon"` is declared), add:

```go
const iconArrowDown = svg.Icon("ss-arrow-down")
```

Replace line 85:

```go
// BEFORE
svg.Svg().Child(svg.Use().Attr("href", "#ss-arrow-down")).Set(ClsSsIcon.AsAttr()),
// AFTER
iconArrowDown.Render(string(ClsSsIcon)),
```

Keep the surrounding comment about wrapping header text in a `Span` — it is
still valid. The `svg` import stays (now for `svg.Icon` only, root package).

### Stage 2 — tag and migrate `selectsearch/svg.go` to use `svg/sprite`

Replace the whole file with:

```go
//go:build !wasm

package selectsearch

import "github.com/tinywasm/svg/sprite"

func (c *SelectSearch) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		sprite.Define(iconArrowDown, "0 0 16 16",
			sprite.Path("M1.5 4.5l6.5 7 6.5-7H1.5z"),
		),
	)
}
```

The package-level `arrowDown` var is deleted (the geometry now binds directly
to the shared `iconArrowDown` constant, which lives in the untagged
`selectsearch.go` from Stage 1 and is visible here too since both files are
package `selectsearch`).

**This `//go:build !wasm` tag is YOUR responsibility as the consumer** — the
library does not and cannot enforce it. Forgetting it does not fail to
compile (both `svg` and `svg/sprite` build fine for any target); it silently
grows the WASM binary. That is exactly what Stage 4's leak check catches.

### Stage 3 — `AGENTS.md` consistency check

`AGENTS.md` already contains the "Build tags belong to the consumer" and "SVG
icons" sections. Verify the code now matches them exactly; if you find a
discrepancy between AGENTS.md and this plan, this plan wins — report the
discrepancy.

### Stage 4 — tests and mandatory leak-check verification

- Update any test asserting the old markup. Search:
  `grep -rn "ss-arrow-down" --include='*_test.go' .` — assertions must expect
  `href="#ss-arrow-down"` produced by `iconArrowDown.Render`, and
  `IconSvg()` tests keep asserting the `<symbol id="ss-arrow-down">` output.
- Run `gotest` (never `go test`).
- Leak/one-path checks (all must be empty) — run in THIS repo's root:

```bash
GOOS=js GOARCH=wasm go build ./...
GOOS=js GOARCH=wasm go list -deps ./selectsearch | grep tinywasm/svg/sprite   # MUST be empty — the key leak check
grep -rn 'svg.Svg()\|svg.Use()' --include='*.go' .
grep -rn '"#ss-arrow-down"' --include='*.go' . | grep -v _test
```

The second command is the substitute for a compile-time guarantee: since
`svg/sprite` has no build tag of its own, only this dependency-graph check
proves the WASM build never reaches it.

## Anti-footguns (do NOT do)

- **Do NOT convert `themetoggle` to SVG icons.** Its `icon()` function returns
  text glyphs by design (switch over theme state). Out of scope.
- **Do NOT add icons to `actionbutton`, `contentcard`, `datatable`,
  `modaldialog`.** They have none; this plan only fixes `selectsearch`.
- **Do NOT touch `css.go` files** — their `//go:build !wasm` tags are already
  correct.
- **Do NOT add stdlib imports** to any untagged file (`tinywasm/fmt` etc. only).
- Never run `gopush` or `codejob`.

## Stages table

| # | Stage | Files | Done |
|---|---|---|---|
| 1 | Typed reference + Render | `selectsearch/selectsearch.go` | ☐ |
| 2 | Tag + migrate to `svg/sprite` | `selectsearch/svg.go` | ☐ |
| 3 | AGENTS.md consistency | `AGENTS.md` (verify only) | ☐ |
| 4 | Tests + mandatory leak check | `selectsearch/*_test.go` | ☐ |
