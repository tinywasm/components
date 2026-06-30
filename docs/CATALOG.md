# Components Catalog

This catalog documents the available reusable components in `tinywasm/components`.

[← Back to Main README](../README.md)

## Theme

All components consume CSS custom properties from `tinywasm/dom`'s `theme.css`.
Inject `dom.ThemeCSS` into your page `<head>` once via the site builder.
Components do not define colors — they inherit from the theme.

## Overview

All components follow the [Component Creation Guide](./SKILL.md).

---

## [ActionButton](../actionbutton/README.md)
Versatile button component with variant support (primary, secondary, danger).
[Detailed Documentation →](../actionbutton/README.md)

---

## [ContentCard](../contentcard/README.md)
Container with header, body, and footer sections.
[Detailed Documentation →](../contentcard/README.md)

---

## [NavBar](../navbar/README.md)
Navigation menu with support for icons.
[Detailed Documentation →](../navbar/README.md)

---

## [Dialog](../dialog/README.md)
Modal dialog overlay with backdrop and close button.
[Detailed Documentation →](../dialog/README.md)

---

## [DataTable](../datatable/README.md)
Data table for structured information with headers and rows.
[Detailed Documentation →](../datatable/README.md)

---

## [SelectSearch](../selectsearch/README.md)
Signal-driven searchable dropdown with static options, live filtering, and optional DB search callback. Uses `BindChildren` for efficient list updates and `Show` for the dropdown.
[Detailed Documentation →](../selectsearch/README.md)

---

## [ThemeToggle](../themetoggle/README.md)
Signal-driven floating button that cycles between `auto → dark → light` theme modes. Persists preference in `localStorage` via `Init`. Uses derived signals for labels and icon updates.
[Detailed Documentation →](../themetoggle/README.md)

---

## Forms

Forms are NOT part of `tinywasm/components`.
Use `github.com/tinywasm/form` directly — it is the only form library in the
tinywasm ecosystem and provides field layout, validation, labels and submission.
