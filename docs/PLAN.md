# PLAN — Typed CSS migration for tinywasm/components

## Goal

Migrate every component in `tinywasm/components/` to the typed CSS DSL. After this plan:
- No component contains a `.css` file.
- No component uses `//go:embed *.css`.
- Every `RenderCSS()` returns `*css.Stylesheet`.
- Class names are declared as exported `css.Class` constants and shared between the HTML and CSS sides of the same component.

## Scope (big-bang)

Seven components migrate in the same cycle:

| Component | Current `.css` LOC (approx) |
|---|---|
| `button`        | ~180 |
| `card`          | ~30 |
| `modal`         | tbd |
| `nav`           | tbd |
| `selectsearch`  | tbd |
| `table`         | tbd |
| `themeswitch`   | tbd |


## Why big-bang

Decided in the design discussion: gradual coexistence would require assetmin to keep both the AST extractor and the invoke path alive, doubling the surface that this initiative is meant to shrink. Migrating all components at once lets the AST extractor be deleted in the same cycle.

## Component template (canonical form)

Each component's `ssr.go` becomes the **single** SSR file (no separate `.css`, no second Go file unless a component is exceptionally large — in which case a `ssr_styles.go` in the same package is acceptable as an exception, not a pattern).

```go
//go:build !wasm
package button

import . "github.com/tinywasm/css"

// Exported classes — consumed by both the HTML side (button.go) and the CSS below.
var (
    ClsBtn       Class = "btn"
    ClsPrimary   Class = "btn-primary"
    ClsSecondary Class = "btn-secondary"
    ClsDanger    Class = "btn-danger"
)

// SSRInstance satisfies the assetmin invoke convention.
func SSRInstance() *Button { return &Button{} }

func (b *Button) RenderCSS() *Stylesheet {
    return New(
        Rule(ClsBtn,
            Padding(Rem(0.5), Rem(1)),
            Border(None),
            BorderRadius(RadiusSm),
            Cursor(Pointer),
            FontSize(TextBase),
        ),
        Rule(ClsPrimary,
            Background(ColorPrimary),
            Color(ColorOnPrimary),
        ),
        Rule(ClsSecondary,
            Background(ColorMuted),
            Color(ColorOnSurface),
        ),
        Rule(ClsDanger,
            Background(ColorError),
            Color(ColorOnSurface),
        ),
        Rule(ClsBtn.Hover(),
            Opacity(0.9),
        ),
    )
}

func (b *Button) RenderHTML() string { /* unchanged */ }
func (b *Button) RenderJS() string   { /* unchanged */ }
func (b *Button) IconSvg() map[string]string { return nil }
```

And `button.go` (compiles for both wasm and !wasm) consumes the **same** class constant:

```go
package button

import (
    "github.com/tinywasm/css"
    "github.com/tinywasm/dom"
)

func (b *Button) Render() dom.Node {
    return dom.Button(dom.Classes(css.ClsBtn, css.ClsPrimary), b.Label)
}
```

## Per-component migration steps

Repeat for each of the 7 components:

1. Audit the existing `.css` file. List every selector, every declaration, every token reference, every magic value.
2. Identify selectors → choose Go names following the pattern `Cls<Variant>` (`ClsPrimary`, `ClsHeader`, `ClsBody`, `ClsFooter`, etc.). Declare them as exported `css.Class` constants.
3. Replace each token reference `var(--space-4, 1rem)` with the matching token (`Space4`).
4. Replace each magic value with the closest scale token. If a value has no token (e.g. `.4rem` for ad-hoc padding), either:
   - round to the nearest `Space*` token, or
   - emit `Rem(0.4)` — flagged in the PR description as a candidate for token addition if it recurs.
5. Translate selectors to DSL form. Pseudo-classes via `Class.Hover()` / `Class.Focus()` / `Class.Disabled()`. Anything else (attribute selectors like `button[name*="btn"]`) goes through `Selector("button[name*=\"btn\"]")`.
6. `@keyframes` → `Keyframes("name", At("0%", ...decls), At("100%", ...decls))` using the typed builders from `tinywasm/css` (see `tinywasm/css/docs/PLAN.md` for the API addition that this plan depends on). Token references inside frame declarations behave like any other DSL rule — renaming the token breaks the build. Data-URI background images → `BackgroundImage(Str("url(...)"))` (raw string is fine for SVG payloads).
7. Add `SSRInstance()`.
8. Update any `.go` file in the component that emitted class names as string literals to instead reference `css.Cls<...>` via `dom.Class()` / `dom.Classes()`.
9. Delete the `.css` file.
10. Remove the `//go:embed` directive and the `var css string` declaration.

## Special cases

- **`themeswitch`**: emits CSS that targets `[data-theme="dark"]` / `[data-theme="light"]` outside of any class scope. Use `Selector("[data-theme=\"dark\"]")` and `Selector("[data-theme=\"light\"]")` with `Bind()`-style declarations (the same primitive used by `tinywasm/css/ssr.go` for OS dark mode).
- **`button`**: contains (a) `@keyframes pulse-url` referencing `ColorSecondary` — migrated via `Keyframes("pulse-url", At(...)...)` from `tinywasm/css`; depends on the keyframes API addition tracked in `tinywasm/css/docs/PLAN.md`. (b) Base64-encoded SVG data URIs for icon backgrounds — legitimate raw payloads, wrap them with `BackgroundImage(Str(...))`. Consider migrating them to the `IconSvg()` channel instead, since assetmin already has icon plumbing; out of scope for this PR.
- **`modal` / `nav` / `table` / `selectsearch`**: audit needed first. If any contain `@container` or CSS layers, those features must first be added to `tinywasm/css` (as typed constructors or via `Raw()` as a temporary escape hatch) — they are not introduced inside the component packages.

## Files removed

- `components/button/button.css`
- `components/card/card.css`
- `components/modal/modal.css`
- `components/nav/nav.css`
- `components/selectsearch/selectsearch.css`
- `components/table/table.css`
- `components/themeswitch/themeswitch.css`
- All `//go:embed *.css` directives in each `ssr.go`.

## Files modified

- 7 `ssr.go` files — rewritten as DSL.
- Each component's main `.go` file — class-name string literals replaced by `css.Cls<...>` references via `dom.Class*()` adapters.
- `components/docs/CREATION.md` — replace the "embed your CSS" instructions with the DSL pattern (link to `tinywasm/css/README.md` and `tinywasm/css/docs/PLAN.md` for API reference).
- `components/docs/CATALOG.md` — update any code samples.

## Steps (order)

1. Wait for `tinywasm/css` (keyframes addition — `Keyframes()`/`At()`/`KeyframeStep`), `tinywasm/dom`, and `tinywasm/assetmin` plans to land. They provide, respectively: the keyframes builders required by `button`; the `RenderCSS() *css.Stylesheet` contract; and the `SSRInstance()` invoke convention. The core DSL and token catalog from the original `tinywasm/css` plan are already published.
2. Migrate `card` first — smallest CSS, sanity-check the DSL coverage.
3. Migrate `button` — largest, most token-heavy, exercises pseudo-classes + `@keyframes` + attribute selectors + data-URI backgrounds.
4. Migrate the remaining five in any order.
5. Run `go test ./...` per component; ensure rendered `RenderCSS().String()` is visually equivalent to the legacy `.css` (golden diff acceptable for whitespace).
6. Update `components/docs/CREATION.md` and `CATALOG.md`.

## Acceptance

- No `.css` file under `components/`.
- No `//go:embed` directive in any component's `ssr.go`.
- Every component implements `RenderCSS() *css.Stylesheet`.
- Every component exposes `SSRInstance()` per the assetmin convention.
- For every class used by a component, exactly one `css.Class` constant declares its name. No string literal class names in HTML emission paths.
- Visual regression of a demo page (`components/docs` examples) against current output is zero or whitespace-only.

## Out of scope

- Restructuring component packages.
- Adding new components or variants.
- Replacing data-URI SVG backgrounds with the icon system (separate cleanup PR).
