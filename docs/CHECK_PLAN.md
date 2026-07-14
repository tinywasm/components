# tinywasm/components — Plan: slot-ready catalog for preconfigured layouts

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Supersedes the previous signal-migration plan: verified executed in code
> (no `OnMount`/`Update()` remains; `themetoggle` and `selectsearch` are
> signal-driven). Stage 0 closes its pending documentation, then this plan
> aligns the catalog with the new layered layout architecture.

## Context (zero-context summary)

The ecosystem adopted a layered UI architecture (see
`tinywasm/layout/docs/PLAN.md` and its consumer plan in
`veltylabs/mjosefa-cms/docs/PLAN.md`):

- `tinywasm/components` — **raw reusable pieces** (this repo): buttons, cards,
  dialogs, tables, selects. No layout knowledge.
- `tinywasm/layout` — published layout skeletons (`platformd` shell,
  `rightpanel`, new `crudview`). `crudview` is self-contained and does **not**
  import this repo.
- **The consumer's composition root** (e.g. `mjosefa-cms/config/layouts`)
  preconfigures layouts once and injects components into their slots
  (`Form`, `Detail`, `HeadControls`, `Aside`, `UserBlock`, …). Modules pick a
  preconfigured layout; they never assemble components by hand.

For that to work, every cataloged component must be **slot-ready** and
**theme-driven**: drop it into any layout slot and it inherits the app's
branding (set once via `RootCSS`) with zero per-component configuration. That
is the contract this plan makes explicit, verified, and documented.

**Dependency rule (final): `components` never imports `tinywasm/layout`.**
The reverse (a high-level layout importing a component) is allowed but
currently unused; assembly belongs to the consumer.

**Ecosystem pillars:** minimal WASM binary, avoid allocations, zero consumer
boilerplate, reusable architecture. `gotest` only (install:
`go install github.com/tinywasm/devflow/cmd/gotest@latest`).

## Stage 0 — close the executed signal migration

The code work of the previous plan is done. Verify and finish its doc closure:

- `docs/CATALOG.md` and `docs/SKILL.md` describe the signal-based API for
  `themetoggle` and `selectsearch` (no `OnMount`/`Update` anywhere in docs).
- Each component `README.md` shows `Signal`-based usage; root `README.md`
  index is current.
- If any of the above is already done, this stage is a no-op — do not churn.

## Stage 1 — theme-token audit (branding must reach every component)

Apps brand via `RootCSS()` overriding the canonical CSS custom properties; a
component with hardcoded colors silently escapes the brand (the Pa100T
consumer sets accent `#3f88bf` / grays `#e9e9e9`/`#c2c1c1` as tokens).

- Audit every `css.go` in the catalog (`actionbutton`, `contentcard`,
  `datatable`, `dialog`, `selectsearch`, `themetoggle`, `navbar` if present):
  colors, radii, spacing come from the canonical theme tokens; literal color
  values are only allowed as `var(--token, <fallback>)` fallbacks.
- Acceptance is mechanical: `grep -rn "#[0-9a-fA-F]\{3,6\}" */css.go` returns
  only token fallbacks inside `var(...)`.

## Stage 2 — slot-readiness contract, enforced by tests

A component is slot-ready when a layout can hold it as `dom.Component` with no
extra ceremony:

- Embeds `dom.Element` **by value** (never `*dom.Element`).
- Implements `Render() *dom.Element`; optional `Init(ctx dom.Ctx)` runs once
  (guarded); no other lifecycle.
- All dynamic state in typed signals; configuration is exported struct fields
  (data), not constructors with behavior — so a consumer's composition root
  can preconfigure it declaratively.

Add one compile-time + behavior contract test per package
(`<name>_contract_test.go`): asserts interface satisfaction, double-`Init`
safety, and `Render` idempotence (two calls, same shape). Shared helper in a
single `internal`-free support file at repo root if needed (no new deps).

## Stage 3 — documentation of the layered role

- `README.md` + `docs/CATALOG.md`: new section **"Where components fit"** with
  the layer diagram (components → injected into layout slots → assembled once
  in the consumer's composition root, e.g. `config/layouts`) and the two rules:
  components never import `layout`; assembly never lives in modules.
- `docs/SKILL.md`: add the slot-readiness contract (Stage 2 bullets) to the
  component-creation standards, and the theme-token rule (Stage 1).
- `docs/CATALOG.md`: mark each entry slot-ready once its contract test lands.

## Harness checklist (mandatory)

- No stdlib in wasm-compiled code (`tinywasm/fmt`/`json`); no `any`/`map`/
  generics in public APIs; `switch` over `map`.
- Value embedding; typed signals; CSS/SVG only in `css.go`/`svg.go` (`!wasm`).
- No repeated string literals in logic; class names are typed vars.
- No public API breaks: existing consumers (`mjosefa-cms`, layout demo)
  compile unmodified.

## Acceptance criteria

1. `gotest ./...` green (stdlib + browser suites).
2. Token grep (Stage 1) clean; a consumer overriding accent/gray tokens
   restyles every cataloged component (verified in the existing demo pages).
3. Every cataloged component has a passing contract test (Stage 2).
4. CATALOG/README/SKILL document the layered role; no doc mentions
   `OnMount`/`Update()`.

## Status Update (Agent Handover — corrected 2026-07-11)

A prior handover note claimed Stage 2 and Stage 3 were done; verification found
Stage 2 had **zero** contract test files anywhere in the repo, and Stage 3 had
**no** "Where components fit" section and **no** slot-ready markers in
`docs/CATALOG.md`. Both were completed for real in this pass, and a stale
`NavBar` catalog entry (no `navbar/` package exists in this repo) was removed.

### ✅ Completed (verified)
- **Stage 0**: `README.md`, `docs/CATALOG.md`, `docs/SKILL.md`, and per-component
  `README.md` files describe the signal-based API; no `OnMount`/`Update()` usage
  remains (mentions in docs are explanatory — "there is no `Update()`").
- **Stage 1**: `grep -rn "#[0-9a-fA-F]\{3,6\}" */css.go` returns nothing —
  clean. Some Rem/Pct values remain as string literals inside `Decl{...}`
  because `css.Value.cssValue()` is unexported in `tinywasm/css v0.1.4`; this
  is a units literal, not a color/token escape, so it doesn't violate Stage 1.
- **Stage 2**: Added `<name>_contract_test.go` for all 6 components
  (`actionbutton`, `contentcard`, `datatable`, `dialog`, `selectsearch`,
  `themetoggle`), each asserting `dom.Component` satisfaction, double-`Init`
  safety, and `Render` shape-idempotence (ids normalized before comparison,
  since `Show`/`BindChildren` assign fresh auto-ids per render).
- **Stage 3**: Added a "Where components fit" section with the layer diagram
  to `README.md` and `docs/CATALOG.md`; added a "Slot-Readiness Contract"
  section (§5) and theme-token rule to `docs/SKILL.md`; marked all 6 catalog
  entries ✅ **Slot-ready**.
- `gotest -no-cache` (full suite, `wasmbrowsertest` present locally) →
  `vet ✅, race ✅, tests ✅, wasm ✅, coverage: 71.1%` after all changes.

### ❌ Remaining / Issues
- None. (A previous version of this note claimed WASM tests couldn't run
  locally — that was wrong: `gotest ./...` skips the vet/wasm/badges stages
  because it received args, per `gotest --help`; the plain `gotest -no-cache`
  invocation runs the full suite including `wasm ✅`.)

## Stages (Original)

| Stage | File(s) | Action | Status |
|---|---|---|---|
| 0 | `docs/CATALOG.md`, `docs/SKILL.md`, `*/README.md` | close signal-migration docs | ✅ Done |
| 1 | `*/css.go` | theme-token audit; literals → `var(--token, fallback)` | ✅ Done* |
| 2 | `*/<name>_contract_test.go` | slot-readiness contract tests | ✅ Done |
| 3 | `README.md`, `docs/CATALOG.md`, `docs/SKILL.md` | document the layered role | ✅ Done |

*\*Some unit string literals remain where `tinywasm/css` DSL methods are unexported.*
