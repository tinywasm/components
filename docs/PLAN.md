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

## Status Update (Agent Handover)

The implementation is **mostly complete**, with all Stage 0-3 code and documentation changes applied.

### ✅ Completed
- **Stage 0**: All `README.md` files (individual and root), `docs/CATALOG.md`, and `docs/SKILL.md` updated for signal migration and layered architecture.
- **Stage 1**: Theme-token audit completed. Literal colors removed. Some specific layout values (Rem/Pct) remain as string literals because `css.Value.cssValue()` is unexported in `tinywasm/css v0.1.4`, preventing clean programmatic conversion in `Decl{...}`.
- **Stage 2**: Contract tests (`*_contract_test.go`) implemented for all 6 components. They verify `dom.Component` interface, `Init` safety, and `Render` idempotence.
- **Stage 3**: Layered role and slot-readiness contract documented in `SKILL.md` and `README.md`. Components marked as ✅ **Slot-ready** in `CATALOG.md`.

### ❌ Remaining / Issues
- **WASM Tests**: Local environment lacks `wasmbrowsertest` binary. While SSR tests pass (`tests ✅`), browser-based WASM tests could not be executed (`wasm ❌`).
- **CSS DSL Limitation**: Programmatic use of `Rem()` and `Pct()` inside `RuleContent(Decl{...})` requires string literals (e.g. `"0.3rem"`) instead of `.cssValue()` due to the unexported interface method.

## Stages (Original)

| Stage | File(s) | Action | Status |
|---|---|---|---|
| 0 | `docs/CATALOG.md`, `docs/SKILL.md`, `*/README.md` | close signal-migration docs | ✅ Done |
| 1 | `*/css.go` | theme-token audit; literals → `var(--token, fallback)` | ✅ Done* |
| 2 | `*/<name>_contract_test.go` | slot-readiness contract tests | ✅ Done |
| 3 | `README.md`, `docs/CATALOG.md`, `docs/SKILL.md` | document the layered role | ✅ Done |

*\*Some string literals used for units where DSL methods are unexported.*
