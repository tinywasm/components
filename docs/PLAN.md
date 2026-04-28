# PLAN: Fix SelectSearch UX Bugs

## Context

Three UX bugs exist in `selectsearch`:

1. **No auto-focus** on the search input when the dropdown opens.
2. **Dropdown closes on each keystroke**, forcing the user to click the header again after every character typed.
3. **Cursor jumps to position 0** on every keystroke while typing in the search input.

---

## Root Cause Analysis

### Bug 1 — No auto-focus on open

**Mechanism:** The dropdown is toggled via a hidden CSS checkbox (`.ss-toggle`). When the user clicks the header label, the browser natively checks/unchecks the checkbox, and CSS reveals `.ss-dropdown` via `.ss-toggle:checked ~ .ss-dropdown`. There is **no JavaScript `change` listener** on the toggle to call `Focus()` on the search input after the dropdown becomes visible.

**Effect:** The user must manually click the search field before they can type.

**Status: Fixed** — `change` listener added in `front.go`.

---

### Bug 2 — Dropdown closes on each keystroke

**Mechanism:** The `input` event handler in `front.go` calls `c.Update()` after updating `c.filterTerm`. `c.Update()` replaces the component's `outerHTML` with the result of a fresh `Render()` call. The `Render()` function created the toggle checkbox without a `checked` attribute, so the CSS rule `.ss-toggle:checked ~ .ss-dropdown { display: block }` no longer applied, and the dropdown hid.

**Status: Fixed** — `isOpen bool` field added to the struct; `Render()` emits `checked` on the toggle when `isOpen == true`.

---

### Bug 3 — Cursor jumps to position 0 on every keystroke

**Mechanism:** `dom.Update()` (in `tinywasm/dom`) replaces the entire component via `outerHTML`. This destroys the active DOM element (the search input the user was typing into), moving browser focus to `document.body`. After the replacement, `front.go` calls `Focus()` on the new element to restore interaction. However, calling `focus()` on a fresh input element (even with a `value` attribute set) positions the cursor at position 0 in some environments, not at the end of the typed text.

**Effect:** The cursor resets on every character, making the input field unusable for multi-character queries.

**Root fix location:** `tinywasm/dom` — `Update()` must snapshot `document.activeElement` before `outerHTML` replacement and restore focus + cursor via `setSelectionRange` afterwards. See [tinywasm/dom/docs/PLAN.md](../../../dom/docs/PLAN.md).

**Status: Partially mitigated** — `Focus()` in the `input` handler is a workaround that keeps interaction alive but does not fix the cursor reset. The real fix is in `tinywasm/dom`.

---

## Affected Files

| File | Change needed |
|------|--------------|
| `selectsearch/selectsearch.go` | Add `isOpen bool` field; conditionally emit `checked` on toggle in `Render()` ✅ |
| `selectsearch/front.go` | Add `change` listener on toggle; set `isOpen` in `input`/`click` handlers ✅ |
| `tinywasm/dom` `dom_frontend.go` | Preserve active element focus + cursor across `outerHTML` replacement ⏳ |

---

## Tests

Tests are in `selectsearch/selectsearch_test.go`. Run with:

```bash
gotest -run TestSelectSearch
```

| Test | Bug covered |
|------|-------------|
| `TestSelectSearch_OpenState_RendersChecked` | Bug 2 — `isOpen=true` must emit `checked` on toggle |
| `TestSelectSearch_ClosedState_NoChecked` | Bug 2 inverse — `isOpen=false` must not emit `checked` |
| `TestSelectSearch_OpenWithFilter_RendersCheckedAndFiltered` | Bug 2 + filtering — open state preserved with active filter |
| `TestSelectSearch_AfterSelection_RendersUnchecked` | Bug 2 — dropdown collapses after option selection |
| `TestSelectSearch_SearchInput_ValueMatchesFilterTerm` | **Bug 3** — search input must emit `value='<filterTerm>'`; precondition for cursor restoration once `tinywasm/dom` is fixed |
| `TestSelectSearch_Render` *(pre-existing)* | Basic HTML structure |
| `TestSelectSearch_SelectedValue` *(pre-existing)* | Selected label shown in header |
| `TestSelectSearch_Filtering` *(pre-existing)* | Option filtering logic |

> **Bug 1 (auto-focus):** Verified at render-contract level — if `checked` is present in the HTML, `OnMount` calls `Focus()` safely. Full event-loop coverage requires a WASM integration test.
>
> **Bug 3 (cursor reset):** The dom-level WASM integration test is specified in `tinywasm/dom/docs/PLAN.md` → `TestUpdate_PreservesActiveElementFocus`.

---

## No API Changes

All changes are internal to the `selectsearch` package. `SelectSearch`, `Option`, `OnSelect`, and `OnSearch` signatures remain unchanged. No consumer code needs to be updated.
