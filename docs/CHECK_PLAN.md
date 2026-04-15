# PLAN: tinywasm/components — Theme migration + Form deprecation

## Context

`tinywasm/components` provides reusable UI components for the tinywasm ecosystem.

This plan covers three changes:

1. **Theme migration**: All component CSS files must consume `--color-*` tokens from
   `tinywasm/dom`'s `theme.css`. Components that use hardcoded colors or inconsistent
   variable names must be updated.

2. **Form deprecation**: The `form/` sub-package is a thin wrapper with no added value.
   Forms in the tinywasm ecosystem are built exclusively with `tinywasm/form`.
   The `components/form` package must be removed and its documentation corrected.

3. **Documentation update**: All README files, CATALOG.md, and CREATION.md must reflect
   the current state: theme token usage, form removal, and correct component list.

**Prerequisite:** Execute `tinywasm/dom` PLAN.md first — `theme.css` and `CssVars`
must be published before this plan is executed.

---

## Part 1 — CSS theme migration

### Standard token reference (from `tinywasm/dom` theme.css)

All component CSS must use only these variables (never hardcode colors):

| Variable              | Semantic role                        |
|-----------------------|--------------------------------------|
| `--color-primary`     | Main text / foreground               |
| `--color-secondary`   | Brand accent (Go cyan `#00ADD8`)     |
| `--color-tertiary`    | Muted text, borders                  |
| `--color-quaternary`  | Deep background / shadows            |
| `--color-gray`        | Neutral surface                      |
| `--color-selection`   | Selected / active state              |
| `--color-hover`       | Hover accent                         |
| `--color-success`     | Success feedback                     |
| `--color-error`       | Error / danger feedback              |
| `--mag-pri`           | Primary spacing                      |
| `--mag-sec`           | Secondary spacing                    |
| `--mag-cua`           | Quaternary spacing                   |

### Components to migrate

#### `card/card.css` — hardcoded colors, must be replaced

Current (hardcoded):
```css
.card { border: 1px solid #e5e7eb; background-color: white; }
.card-header { border-bottom: 1px solid #e5e7eb; }
.card-footer { background-color: #f9fafb; border-top: 1px solid #e5e7eb; }
```

Replace with:
```css
.card {
    border: 1px solid var(--color-tertiary, #8B949E);
    background-color: var(--color-gray, #0D1117);
    border-radius: 0.5rem;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
}

.card-header {
    padding: var(--mag-pri, 0.5rem);
    border-bottom: 1px solid var(--color-tertiary, #8B949E);
    font-weight: bold;
    color: var(--color-primary, #E6EDF3);
}

.card-body {
    padding: var(--mag-pri, 0.5rem);
    flex-grow: 1;
    color: var(--color-primary, #E6EDF3);
}

.card-footer {
    padding: var(--mag-pri, 0.5rem);
    border-top: 1px solid var(--color-tertiary, #8B949E);
    background-color: var(--color-quaternary, #161B22);
    border-bottom-left-radius: 0.5rem;
    border-bottom-right-radius: 0.5rem;
    color: var(--color-primary, #E6EDF3);
}
```

#### `modal/modal.css` — hardcoded colors, must be replaced

Locate all instances of hardcoded hex values (e.g. `#fff`, `#000`, `rgba(0,0,0,...)`,
`#f3f4f6`, etc.) and replace with the corresponding token + fallback. Example:

```css
/* before */
.modal-content { background: #fff; }
.modal-backdrop { background: rgba(0,0,0,0.5); }

/* after */
.modal-content { background: var(--color-gray, #0D1117); }
.modal-backdrop { background: color-mix(in srgb, var(--color-quaternary, #161B22), transparent 40%); }
```

Apply the same pattern to all remaining hardcoded values in `modal.css`.

#### `nav/nav.css` — hardcoded colors, must be replaced

```css
/* before */
.nav { background-color: #1f2937; color: white; }

/* after */
.nav {
    background-color: var(--color-quaternary, #161B22);
    color: var(--color-primary, #E6EDF3);
}
```

#### `button/button.css` — already uses `--color-*`, verify and keep

`button/button.css` already uses `var(--color-secondary)`, `var(--color-tertiary)`, etc.
Verify all values are from the standard token list. No hardcoded hex colors should remain.

#### `input/input.stlib.css` and `input/input.custom.css` — already use `--color-*`, verify and keep

Both input CSS files already use `var(--color-*)` tokens. Verify no hardcoded colors remain.

#### `table/table.css` — audit and migrate

Audit `table/table.css` for hardcoded colors. Replace with `--color-*` tokens + fallbacks.

#### `form/form.css` — DELETE (see Part 2)

---

## Part 2 — Remove `components/form` package

### Reason

`components/form` is a thin wrapper around `dom.Form()` with no added value:
- It only wraps a `[]dom.Component` slice and an `OnSubmit` handler
- All actual form functionality (field layout, validation, submission, labels) lives in
  `tinywasm/form`, which is the **only** form library in the tinywasm ecosystem
- Having `components/form` creates confusion and a false impression that it is the
  correct way to build forms

### Action: delete entire `form/` sub-package

Files to delete:
- `form/form.go`
- `form/form.css`
- `form/form_test.go`
- `form/README.md`
- `form/ssr.go`

### Action: remove `palette.go` from root package

`palette.go` (the `CssVars` type) is being migrated to `tinywasm/dom`. Once the dom
PLAN.md is executed and `dom.CssVars` is published, delete:
- `components/palette.go`
- `components/registry.go` — if unused after form removal, evaluate deletion

Verify no other file in the module imports `palette.go` types before deleting.

---

## Part 3 — Documentation update

### `docs/CATALOG.md` — update

- Remove the Form entry
- Add clarification block:

```markdown
## Forms

Forms are NOT part of `tinywasm/components`.
Use `github.com/tinywasm/form` directly — it is the only form library in the
tinywasm ecosystem and provides field layout, validation, labels and submission.
```

- Add note about theme tokens at the top:

```markdown
## Theme

All components consume CSS custom properties from `tinywasm/dom`'s `theme.css`.
Inject `dom.ThemeCSS` into your page `<head>` once via the site builder.
Components do not define colors — they inherit from the theme.
```

### `docs/CREATION.md` — update CSS section

Replace the CSS guidelines section with:

```markdown
## CSS guidelines

- All colors must use `--color-*` CSS custom properties from `tinywasm/dom` theme.
- Never hardcode hex values for colors. Always provide a fallback:
  `var(--color-secondary, #00ADD8)`.
- Spacing must use `--mag-pri`, `--mag-sec`, `--mag-cua` variables.
- CSS class names must be prefixed with the component name to avoid collisions.
- CSS lives in `<component>.css`, embedded in `ssr.go` via `//go:embed`.
- Do NOT create or embed form-related CSS — use `tinywasm/form`.
```

### `README.md` — update

- Remove Form from the component list
- Add a "Forms" section pointing to `tinywasm/form`
- Update the theme/installation section to mention `dom.ThemeCSS`

### `form/README.md` — replace content before deleting file

Before deleting, update the file to redirect (in case anyone has the old URL):

```markdown
# Form — DEPRECATED in tinywasm/components

The `components/form` package has been removed.

Use `github.com/tinywasm/form` directly — it is the standard form library
for the tinywasm ecosystem.
```

Then delete the file along with the rest of the `form/` package.

---

## Execution order for Jules

1. Verify `github.com/tinywasm/dom` has `CssVars`, `ThemeCSS` and `theme.css`
   published (run `go get github.com/tinywasm/dom@latest`)

2. **Migrate CSS** (Part 1):
   - Update `card/card.css`
   - Update `modal/modal.css`
   - Update `nav/nav.css`
   - Audit and update `table/table.css`
   - Verify `button/button.css` — no hardcoded colors
   - Verify `input/input.stlib.css`, `input/input.custom.css` — no hardcoded colors

3. **Remove form package** (Part 2):
   - Delete `form/` directory entirely
   - Delete `components/palette.go`
   - Evaluate `components/registry.go` — delete if unused

4. **Update documentation** (Part 3):
   - Update `docs/CATALOG.md`
   - Update `docs/CREATION.md`
   - Update `README.md`

5. Run `go test ./...` — all tests must pass (form tests are gone, others must pass)

6. Commit: `feat: migrate CSS to dom theme tokens, remove deprecated form package`

---

## Notes for Jules

- CSS fallback values must match the `DefaultCssVars()` defaults in `tinywasm/dom`:
  primary `#E6EDF3`, secondary `#00ADD8`, tertiary `#8B949E`, quaternary `#161B22`,
  gray `#0D1117`, selection `#654FF0`, error `#E34F26`, success `#3FB950`.
- Do NOT import `tinywasm/form` from this package — components have no dependency on it.
  The documentation only references it by import path for the user.
- The `registry.go` file is currently unused. Delete it if no component references it
  after the form removal.
- `palette.go` can only be deleted after confirming `tinywasm/dom` has published
  `CssVars`. If dom is not yet updated, keep `palette.go` and add a TODO comment.
