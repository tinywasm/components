# PLAN: Create `selectsearch` Component

## Context

Module: `github.com/tinywasm/components`
New package: `github.com/tinywasm/components/selectsearch`
Location: `tinywasm/components/selectsearch/`
Go version: 1.25.2
Final WASM compiler: TinyGo (binary size and heap allocation are critical constraints)

**Dependencies available in go.mod:**
- `github.com/tinywasm/dom v0.7.2`
- `github.com/tinywasm/fmt v0.23.5` (indirect)

Do NOT add new dependencies. Do NOT import `tinywasm/form`.

## What This Component Does

`selectsearch` is a dropdown with integrated search. It replaces a native `<select>` with:
- A visible "selected value" header that toggles the dropdown open/close
- A live search box that filters options as the user types
- A list of selectable options (each with an ID, a label, and an optional description badge)
- A callback for DB search when all local options are filtered out
- A callback triggered on option selection

**UX behavior (CSS-first):**
- The dropdown toggle (open/close) is driven by CSS using a hidden `<input type="checkbox">` — no JS/WASM needed for toggle.
- Filtering is driven by `OnMount()` event wiring (WASM only) — cannot be done in CSS.

## File Structure to Create

```
tinywasm/components/selectsearch/
├── selectsearch.go        # Struct, Option type, Render(), OnMount()
├── selectsearch.css       # Component styles
├── selectsearch_test.go   # Tests
└── ssr.go                 # //go:build !wasm — RenderCSS(), IconSvg()
```

## Data Types

```go
// Option represents a selectable item.
type Option struct {
    ID          string // unique identifier, returned in OnSelect
    Label       string // visible text
    Description string // optional badge shown on the right
}
```

## Struct Definition

```go
type SelectSearch struct {
    dom.Element          // value embed — NEVER pointer (TinyGo heap constraint)
    Placeholder string   // text shown when nothing is selected
    Options     []Option // initial static options
    OnSelect    func(id, description string) // called when user picks an option
    OnSearch    func(term string) []Option   // called when ALL local options are filtered out
}
```

**Critical:** `dom.Element` is embedded as a VALUE (not `*dom.Element`). This is mandatory:
- TinyGo GC is conservative — fewer heap objects = fewer pauses
- Value embed = 1 allocation (struct inline), pointer embed = 2 allocations (struct + Element)
- No nil-guard needed, no nil panic risk in production WASM

## CSS-First Toggle Mechanism

The dropdown open/close state is controlled by a hidden `<input type="checkbox">` with id `<componentID>-toggle`.

```html
<!-- Simplified structure -->
<div id="<id>" class="ss-box">
  <input type="checkbox" id="<id>-toggle" class="ss-toggle">
  <label for="<id>-toggle" class="ss-header">Placeholder text</label>
  <div class="ss-dropdown">
    <input type="search" class="ss-search" placeholder="Search...">
    <div class="ss-options">
      <!-- one div.ss-option per Option -->
      <div class="ss-option" data-id="..." data-description="...">
        <span class="ss-label">Label</span>
        <span class="ss-desc">Description</span>
      </div>
    </div>
  </div>
</div>
```

CSS controls visibility:
```css
.ss-dropdown { display: none; }
.ss-toggle:checked ~ .ss-dropdown { display: block; }
```

The `<label for="<id>-toggle">` acts as the clickable header — toggling the checkbox state without any JS.

## `Render()` Implementation

```go
func (c *SelectSearch) Render() *dom.Element {
    id := c.GetID()
    placeholder := c.Placeholder
    if placeholder == "" {
        placeholder = "Select..."
    }

    // Hidden checkbox — drives the CSS toggle
    toggle := &dom.Element{}
    toggle.Attr("type", "checkbox").
        ID(id+"-toggle").
        Class("ss-toggle")

    // Header label — clicking it toggles the checkbox
    header := dom.Label(id+"-toggle").
        Class("ss-header").
        Text(placeholder)

    // Search input (raw dom.Element — no form/input dependency)
    searchInput := &dom.Element{}
    searchInput.Attr("type", "search").
        ID(id+"-search").
        Class("ss-search").
        Attr("placeholder", "Search...")

    // Options list
    optList := dom.Div().Class("ss-options").ID(id + "-options")
    for _, opt := range c.Options {
        item := dom.Div().
            Class("ss-option").
            Attr("data-id", opt.ID).
            Attr("data-description", opt.Description).
            Add(dom.Span().Class("ss-label").Text(opt.Label))
        if opt.Description != "" {
            item.Add(dom.Span().Class("ss-desc").Text(opt.Description))
        }
        optList.Add(item)
    }

    dropdown := dom.Div().Class("ss-dropdown").
        Add(searchInput).
        Add(optList)

    return dom.Div().
        Class("ss-box").
        Add(toggle).
        Add(header).
        Add(dropdown)
}
```

## `OnMount()` — WASM Event Wiring

`OnMount()` is called by `tinywasm/dom` after the HTML is injected into the DOM.
No build tag needed — TinyGo eliminates it as dead code in SSR builds.

```go
func (c *SelectSearch) OnMount() {
    id := c.GetID()

    // Wire search filtering
    if searchEl, ok := dom.Get(id + "-search"); ok {
        searchEl.On("input", func(e dom.Event) {
            term := e.Target().Value()
            c.filterOptions(term)
        })
    }

    // Wire option click — close dropdown and call OnSelect
    if optionsEl, ok := dom.Get(id + "-options"); ok {
        optionsEl.On("click", func(e dom.Event) {
            target := e.Target()
            optID := target.Attr("data-id")
            optDesc := target.Attr("data-description")
            if optID != "" {
                // Update header text
                if label := target.Text(); label != "" {
                    if headerEl, ok := dom.Get(id + "-toggle"); ok {
                        _ = headerEl // uncheck the toggle
                        // tinywasm/dom Reference.SetAttr or similar to uncheck
                    }
                }
                if c.OnSelect != nil {
                    c.OnSelect(optID, optDesc)
                }
            }
        })
    }
}

// filterOptions hides options that don't match term.
// Calls OnSearch when ALL options are hidden (lost == len(options)).
func (c *SelectSearch) filterOptions(term string) {
    // Implementation: iterate c.Options, show/hide via dom.Get + style
    // Count hidden; if hidden == len(c.Options) && c.OnSearch != nil → call c.OnSearch(term)
}
```

**Note to agent:** Review the `dom.Reference` interface in `tinywasm/dom/reference.go` to understand the exact API for getting/setting element properties (value, style, attributes) from WASM. Use only what that interface exposes — do not use `syscall/js` directly.

## `ssr.go`

```go
//go:build !wasm

package selectsearch

import _ "embed"

//go:embed selectsearch.css
var css string

func (c *SelectSearch) RenderCSS() string { return css }

func (c *SelectSearch) IconSvg() map[string]string {
    return map[string]string{
        // Arrow down icon for the dropdown header
        // viewBox 0 0 16 16 (default)
        "ss-arrow-down": `<path fill-rule="evenodd" d="M1.5 4.5l6.5 7 6.5-7H1.5z"/>`,
    }
}
```

## `selectsearch.css`

Use CSS custom properties from `tinywasm/dom` theme. All class names prefixed `ss-`.

Key rules to implement:
```css
.ss-toggle { display: none; }

.ss-dropdown { display: none; }
.ss-toggle:checked ~ .ss-dropdown { display: block; }

.ss-header {
    background: var(--color-secondary, #00ADD8);
    color: var(--color-gray, #fff);
    padding: var(--mag-pri, 0.5rem) calc(var(--mag-pri, 0.5rem) * 2);
    cursor: pointer;
    border-radius: 0.4em;
    display: block;
    position: relative;
}

/* Arrow icon via SVG sprite */
.ss-header::after {
    content: "";
    /* use background SVG sprite reference if available, else CSS arrow */
    position: absolute;
    right: var(--mag-pri, 0.5rem);
    top: 50%;
    transform: translateY(-50%);
    width: 1em;
    height: 1em;
}

.ss-toggle:checked ~ .ss-header::after {
    transform: translateY(-50%) rotate(180deg);
}

.ss-search {
    width: 100%;
    border: 0.2em solid var(--color-tertiary, #6E6E73);
    border-radius: 0.4em 0.4em 0 0;
    padding: var(--mag-pri, 0.5rem);
    font-size: 1rem;
}

.ss-search:focus { outline: none; }

.ss-options {
    max-height: 240px;
    overflow-y: auto;
    background: var(--color-quaternary, #F2F2F7);
}

.ss-option {
    padding: var(--mag-pri, 0.5rem) calc(var(--mag-pri, 0.5rem) * 2);
    border-bottom: 1px solid var(--color-tertiary, #6E6E73);
    cursor: pointer;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.ss-option:hover { background: var(--color-hover, #B8860B); color: var(--color-gray, #fff); }

.ss-desc {
    font-size: 0.8em;
    background: var(--color-quaternary, #F2F2F7);
    border-radius: 0.3em;
    padding: 0.2em 0.4em;
    color: var(--color-primary, #1C1C1E);
}

/* Hidden options (set via style.display by OnMount filter) */
.ss-option[hidden] { display: none; }

/* Scrollbar */
.ss-options::-webkit-scrollbar { width: 0.4em; background: none; }
.ss-options::-webkit-scrollbar-thumb { background: var(--color-secondary, #00ADD8); }
```

## Tests (`selectsearch_test.go`)

```go
package selectsearch

import (
    "strings"
    "testing"
)

func TestSelectSearch_Render(t *testing.T) {
    c := &SelectSearch{
        Placeholder: "Choose category",
        Options: []Option{
            {ID: "1", Label: "Automobiles", Description: "auto"},
            {ID: "2", Label: "Film & Animation", Description: "anime"},
        },
    }

    html := c.Render().RenderHTML()

    if !strings.Contains(html, "ss-box") {
        t.Error("expected ss-box class")
    }
    if !strings.Contains(html, "ss-toggle") {
        t.Error("expected ss-toggle checkbox")
    }
    if !strings.Contains(html, "Choose category") {
        t.Error("expected placeholder text")
    }
    if !strings.Contains(html, "Automobiles") {
        t.Error("expected option label")
    }
    if !strings.Contains(html, `data-id='1'`) {
        t.Error("expected option data-id")
    }
    if !strings.Contains(html, "auto") {
        t.Error("expected option description")
    }
}

func TestSelectSearch_DefaultPlaceholder(t *testing.T) {
    c := &SelectSearch{}
    html := c.Render().RenderHTML()
    if !strings.Contains(html, "Select...") {
        t.Error("expected default placeholder")
    }
}
```

## CATALOG.md Update

Add an entry to `tinywasm/components/docs/CATALOG.md`:

```markdown
## 7. [SelectSearch](../selectsearch/README.md)
Searchable dropdown with static options, live filtering, and optional DB search callback.
[Detailed Documentation →](../selectsearch/README.md)
```

## Reference: `dom.Reference` Interface

Before implementing `OnMount()`, read `tinywasm/dom/reference.go` to understand the exact API available for manipulating DOM elements from WASM (getting values, setting styles, attributes). Use only that interface — do not use `syscall/js` directly.

## Reference: `dom.Event` Interface

Read `tinywasm/dom/event.go` to understand `e.Target()` and what methods are available on the event target (e.g., `Value()`, `Attr()`, `Text()`). Use only the documented interface.
