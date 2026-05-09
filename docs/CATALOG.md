# Components Catalog

This catalog documents the available reusable components in `tinywasm/components`.

[← Back to Main README](../README.md)

## Theme

All components consume CSS custom properties from `tinywasm/dom`'s `theme.css`.
Inject `dom.ThemeCSS` into your page `<head>` once via the site builder.
Components do not define colors — they inherit from the theme.

## Overview

All components follow the [Component Creation Guide](./CREATION.md).

---

## 1. [Button](../button/README.md)
Versatile button component with variant support.
[Detailed Documentation →](../button/README.md)

---

## 2. [Card](../card/README.md)
Container with header, body, and footer sections.
[Detailed Documentation →](../card/README.md)

---

## 3. [Input](../input/README.md)
Text input with label and validation support.
[Detailed Documentation →](../input/README.md)

---

## 4. [Nav](../nav/README.md)
Navigation menu with support for icons.
[Detailed Documentation →](../nav/README.md)

---

## 5. [Modal](../modal/README.md)
Dialog overlay with backdrop and close button.
[Detailed Documentation →](../modal/README.md)

---

## 6. [Table](../table/README.md)
Data table for structured information.
[Detailed Documentation →](../table/README.md)

---

## 7. [SelectSearch](../selectsearch/README.md)
Searchable dropdown with static options, live filtering, and optional DB search callback.
[Detailed Documentation →](../selectsearch/README.md)

---

## 8. [ThemeSwitch](../themeswitch/README.md)
Floating button that cycles between `auto → dark → light` theme modes. Persists preference in `localStorage`. Owns the `[data-theme]` CSS overrides — include its CSS in `<head>` via `RenderCSS()`.
[Detailed Documentation →](../themeswitch/README.md)

---

## Forms

Forms are NOT part of `tinywasm/components`.
Use `github.com/tinywasm/form` directly — it is the only form library in the
tinywasm ecosystem and provides field layout, validation, labels and submission.
