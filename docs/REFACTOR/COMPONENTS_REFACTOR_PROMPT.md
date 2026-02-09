# Components Refactor - Execution Prompt

**Library**: `github.com/tinywasm/components`
**Location**: `components/`
**Status**: Ready to execute (after DOM refactor)
**Estimated Time**: 2-3 days
**Priority**: 🟡 High (depends on DOM completion)

---

## Context

You are populating the `tinywasm/components` library with reusable UI components that:
1. Follow **explicit registration** pattern (no auto-magic)
2. Use **SSR/CSR split** with build tags
3. Implement **7 essential components** (Phase 1)
4. Work with the **refactored DOM API** (fluent builder, Elm pattern)

**CRITICAL: Respect Existing CSS**: The component folders already contain `.css` files (e.g., `button/button.css`). These files **MUST NOT** be deleted or overwritten. They must be used as the styles for the components being created.

**Goal**: Create a catalog of production-ready components that reduce boilerplate by 50%+.

---

## Prerequisites

### ✅ Verify DOM Refactor Complete
Before starting, verify DOM has been refactored:
- [ ] `dom.Div()` returns chainable `*Element`
- [ ] `ViewRenderer` interface exists in `dom/interface.dom.go`
- [ ] Example `dom/web/client.go` uses new API
- [ ] `gotest` passes in `dom/`

**If not complete**: Wait for DOM refactor or work with DOM agent to coordinate.

---

## Current State

### Files to Review
- `docs/COMPONENT_CREATION.md` - Current guide
- `components.go` - Empty placeholder
- `palette.go` - CssVars utility (keep this)

### Current Issues
- ❌ No component logic exists (only existing CSS assets in folders)
- ❌ No registration mechanism
- ❌ No examples of SSR/CSR split

---

## Approved Decisions

### Decision 1: Explicit Registration
**Pattern**: Components registered explicitly via `site.RegisterHandlers()`

```go
// User's code
import (
    "github.com/tinywasm/components/button"
    "github.com/tinywasm/components/card"
)

site.RegisterHandlers(
    &button.Button{},
    &card.Card{},
)
```

**With convenience import** (optional):
```go
import _ "github.com/tinywasm/components/all"
// Auto-registers all components
```

**Requirements**:
- No `init()` functions for registration
- Tree-shakeable (only import what you use)
- `all.go` provides convenience for prototyping

---

### Decision 2: Component Structure
**Pattern**: Follow COMPONENT_CREATION.md strictly

```
components/
├── button/
│   ├── button.go       # Shared: struct, Render() or RenderHTML()
│   ├── button.css      # Component styles
│   ├── ssr.go          # //go:build !wasm - RenderCSS(), IconSvg()
│   └── front.go        # //go:build wasm - OnMount() if needed
```

**Requirements**:
- One folder per component
- Mandatory: `{name}.go`, `{name}.css`, `ssr.go`
- Optional: `front.go` (only if client-side interactivity)
- All components embed `dom.BaseComponent`

---

### Decision 3: SSR/CSR Split
**Pattern**: Build tags prevent server code in WASM

```go
// button.go (shared)
package button

type Button struct {
    dom.BaseComponent
    Text    string
    Variant string // primary, secondary, danger
}

func (b *Button) Render() dom.Node {
    return dom.Button().
        ID(b.ID()).
        Class("btn btn-"+b.Variant).
        Text(b.Text).
        ToNode()
}

// ssr.go (backend only)
//go:build !wasm

//go:embed button.css
var css string

func (b *Button) RenderCSS() string { return css }

// front.go (optional, client only)
//go:build wasm

func (b *Button) OnMount() {
    // Add client-side enhancements if needed
}
```

---

## Phase 1 Components (Essential 7)

Implement these components in priority order:

### 1. Button (`components/button/`)
**Purpose**: Primary action button with variants

```go
type Button struct {
    dom.BaseComponent
    Text    string
    Variant string // "primary", "secondary", "danger"
    OnClick func(dom.Event)
}

func (b *Button) Render() dom.Node {
    btn := dom.Button().
        ID(b.ID()).
        Class("btn btn-"+b.Variant).
        Text(b.Text)

    if b.OnClick != nil {
        btn.OnClick(b.OnClick)
    }

    return btn.ToNode()
}
```

**CSS** (`button.css`):
```css
.btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
    font-size: 1rem;
}

.btn-primary {
    background: var(--color-primary, #3b82f6);
    color: white;
}

.btn-secondary {
    background: var(--color-secondary, #6b7280);
    color: white;
}

.btn-danger {
    background: var(--color-error, #ef4444);
    color: white;
}

.btn:hover {
    opacity: 0.9;
}
```

---

### 2. Card (`components/card/`)
**Purpose**: Container with header/body/footer sections

```go
type Card struct {
    dom.BaseComponent
    Header dom.Component
    Body   dom.Component
    Footer dom.Component
}

func (c *Card) Render() dom.Node {
    card := dom.Div().
        ID(c.ID()).
        Class("card")

    if c.Header != nil {
        card.Append(dom.Div().Class("card-header").Append(c.Header))
    }

    if c.Body != nil {
        card.Append(dom.Div().Class("card-body").Append(c.Body))
    }

    if c.Footer != nil {
        card.Append(dom.Div().Class("card-footer").Append(c.Footer))
    }

    return card.ToNode()
}
```

---

### 3. Input (`components/input/`)
**Purpose**: Text input with validation

```go
type Input struct {
    dom.BaseComponent
    Label       string
    Placeholder string
    Value       string
    Type        string // text, email, password
    Required    bool
    OnInput     func(dom.Event)
}

func (i *Input) Render() dom.Node {
    container := dom.Div().Class("input-group")

    if i.Label != "" {
        container.Append(
            dom.Tag("label").
                Attr("for", i.ID()).
                Text(i.Label),
        )
    }

    input := dom.Tag("input").
        ID(i.ID()).
        Attr("type", i.Type).
        Attr("placeholder", i.Placeholder).
        Attr("value", i.Value)

    if i.Required {
        input.Attr("required", "required")
    }

    if i.OnInput != nil {
        input.OnInput(i.OnInput)
    }

    container.Append(input)
    return container.ToNode()
}
```

---

### 4. Nav (`components/nav/`)
**Purpose**: Navigation menu (moved from site)

```go
type Nav struct {
    dom.BaseComponent
    Items []NavItem
}

type NavItem struct {
    Label string
    Route string // e.g., "users", "products"
    Icon  string // optional icon ID
}

func (n *Nav) Render() dom.Node {
    nav := dom.Tag("nav").
        ID(n.ID()).
        Class("nav")

    ul := dom.Tag("ul").Class("nav-list")

    for _, item := range n.Items {
        li := dom.Tag("li").Class("nav-item")

        link := dom.Tag("a").
            Attr("href", "#"+item.Route).
            Text(item.Label)

        if item.Icon != "" {
            // Add icon using SVG sprite
            link.Append(
                dom.Tag("svg").Class("icon").Append(
                    dom.Tag("use").Attr("href", "#"+item.Icon),
                ),
            )
        }

        li.Append(link)
        ul.Append(li)
    }

    nav.Append(ul)
    return nav.ToNode()
}
```

---

### 5. Modal (`components/modal/`)
**Purpose**: Dialog/modal overlay

```go
type Modal struct {
    dom.BaseComponent
    Title   string
    Content dom.Component
    Visible bool
}

func (m *Modal) Render() dom.Node {
    class := "modal"
    if !m.Visible {
        class += " hidden"
    }

    return dom.Div().
        ID(m.ID()).
        Class(class).
        Append(
            dom.Div().Class("modal-backdrop").OnClick(m.Close),
        ).
        Append(
            dom.Div().Class("modal-content").
                Append(
                    dom.Div().Class("modal-header").
                        Append(dom.Tag("h2").Text(m.Title)).
                        Append(
                            dom.Button().Text("×").
                                Class("modal-close").
                                OnClick(m.Close),
                        ),
                ).
                Append(
                    dom.Div().Class("modal-body").Append(m.Content),
                ),
        ).
        ToNode()
}

func (m *Modal) Close(e dom.Event) {
    m.Visible = false
    m.Update()
}

func (m *Modal) Open() {
    m.Visible = true
    m.Update()
}
```

---

### 6. Table (`components/table/`)
**Purpose**: Data table with sorting

```go
type Table struct {
    dom.BaseComponent
    Headers []string
    Rows    [][]string
}

func (t *Table) Render() dom.Node {
    table := dom.Tag("table").
        ID(t.ID()).
        Class("table")

    // Header
    thead := dom.Tag("thead")
    tr := dom.Tag("tr")
    for _, header := range t.Headers {
        tr.Append(dom.Tag("th").Text(header))
    }
    thead.Append(tr)
    table.Append(thead)

    // Body
    tbody := dom.Tag("tbody")
    for _, row := range t.Rows {
        tr := dom.Tag("tr")
        for _, cell := range row {
            tr.Append(dom.Tag("td").Text(cell))
        }
        tbody.Append(tr)
    }
    table.Append(tbody)

    return table.ToNode()
}
```

---

### 7. Form (`components/form/`)
**Purpose**: Form wrapper with validation

```go
type Form struct {
    dom.BaseComponent
    Fields   []dom.Component
    OnSubmit func(dom.Event)
}

func (f *Form) Render() dom.Node {
    form := dom.Tag("form").
        ID(f.ID()).
        Class("form")

    if f.OnSubmit != nil {
        form.OnSubmit(f.OnSubmit)
    }

    for _, field := range f.Fields {
        form.Append(field)
    }

    return form.ToNode()
}
```

---

## Implementation Tasks

### Task 1: Create Registration System

**File**: `components/registry.go`
```go
package components

import "github.com/tinywasm/dom"

// Registry stores component factories (optional, for future use)
var registry = make(map[string]func() dom.Component)

func Register(name string, factory func() dom.Component) {
    registry[name] = factory
}

func Get(name string) dom.Component {
    if factory, ok := registry[name]; ok {
        return factory()
    }
    return nil
}

func List() []string {
    names := make([]string, 0, len(registry))
    for name := range registry {
        names = append(names, name)
    }
    return names
}
```

---

### Task 2: Create Convenience Import

**File**: `components/all/all.go`
```go
package all

import (
    _ "github.com/tinywasm/components/button"
    _ "github.com/tinywasm/components/card"
    _ "github.com/tinywasm/components/input"
    _ "github.com/tinywasm/components/nav"
    _ "github.com/tinywasm/components/modal"
    _ "github.com/tinywasm/components/table"
    _ "github.com/tinywasm/components/form"
)

// Importing this package auto-registers all components
```

**Usage**:
```go
import _ "github.com/tinywasm/components/all"
```

---

### Task 3: Implement Each Component

**IMPORTANT**: Many component folders already contain a `{name}.css` file. **DO NOT DELETE OR OVERWRITE** these files. Use them as the basis for the component's styling and embed them in `ssr.go`.

For each of the 7 components, create:

1. **{name}/{name}.go** - Struct and Render() method
2. **{name}/{name}.css** - Component styles
3. **{name}/ssr.go** - `//go:build !wasm` with RenderCSS()
4. **{name}/front.go** - `//go:build wasm` with OnMount() (if needed)

**Template** (`button/button.go`):
```go
package button

import (
    "github.com/tinywasm/dom"
)

type Button struct {
    dom.BaseComponent
    Text    string
    Variant string
    OnClick func(dom.Event)
}

func (b *Button) Render() dom.Node {
    btn := dom.Button().
        ID(b.ID()).
        Class("btn btn-"+b.Variant).
        Text(b.Text)

    if b.OnClick != nil {
        btn.OnClick(b.OnClick)
    }

    return btn.ToNode()
}

func (b *Button) Children() []dom.Component {
    return nil
}
```

**Template** (`button/ssr.go`):
```go
//go:build !wasm

package button

import _ "embed"

//go:embed button.css
var css string

func (b *Button) RenderCSS() string {
    return css
}

// Optional: Add icons if needed
func (b *Button) IconSvg() map[string]string {
    return nil
}
```

---

### Task 4: Update Base Utilities

**File**: `components/base.go`
```go
package components

import "github.com/tinywasm/dom"

// BaseComponent provides common utilities for component authors
// (Currently just an alias, but can be extended later)
type BaseComponent = dom.BaseComponent

// Props pattern (optional utility)
type Props map[string]any

func (p Props) Get(key string) any {
    return p[key]
}

func (p Props) String(key string) string {
    if v, ok := p[key].(string); ok {
        return v
    }
    return ""
}

func (p Props) Int(key string) int {
    if v, ok := p[key].(int); ok {
        return v
    }
    return 0
}
```

---

### Task 5: Update Documentation

**Update**: `docs/COMPONENT_CREATION.md`
- Add section on using the new fluent builder API
- Add examples using DSL vs string HTML
- Update lifecycle hooks section (OnUpdate now available)

**Create**: `docs/COMPONENT_CATALOG.md`
- Document all 7 Phase 1 components
- Include usage examples for each
- Show customization patterns

**Example section**:
```markdown
## Button

Versatile button component with variant support.

### Usage

```go
import "github.com/tinywasm/components/button"

btn := &button.Button{
    Text: "Click me",
    Variant: "primary",
    OnClick: func(e dom.Event) {
        fmt.Println("Clicked!")
    },
}

btn.Render("body")
```

### Variants
- `primary` - Primary action (blue)
- `secondary` - Secondary action (gray)
- `danger` - Destructive action (red)

### Styling
Customize via CSS variables:
- `--color-primary`
- `--color-secondary`
- `--color-error`

---

### Task 6: Create Example App

**File**: `examples/component-showcase/web/client.go`
```go
//go:build wasm

package main

import (
    "github.com/tinywasm/components/button"
    "github.com/tinywasm/components/card"
    "github.com/tinywasm/components/input"
    "github.com/tinywasm/components/modal"
    "github.com/tinywasm/dom"
    "github.com/tinywasm/fmt"
)

func main() {
    // Showcase all components
    renderShowcase()
    select {}
}

func renderShowcase() {
    // Button showcase
    btnCard := &card.Card{
        Header: &CardHeader{Title: "Buttons"},
        Body: &ButtonShowcase{},
    }
    dom.Render("app", btnCard)

    // Modal example
    modal := &modal.Modal{
        Title: "Example Modal",
        Content: &ModalContent{},
    }
    dom.Append("app", modal)

    // Open modal button
    openBtn := &button.Button{
        Text: "Open Modal",
        Variant: "primary",
        OnClick: func(e dom.Event) {
            modal.Open()
        },
    }
    dom.Append("app", openBtn)
}

// ... helper components ...
```

---

### Task 7: Write Tests

**For each component**, create tests:

`button/button_test.go`:
```go
package button

import (
    "testing"
)

func TestButton_Render(t *testing.T) {
    btn := &Button{
        Text: "Click",
        Variant: "primary",
    }

    node := btn.Render()

    if node.Tag != "button" {
        t.Error("expected button tag")
    }

    // Verify classes
    hasClass := false
    for _, attr := range node.Attrs {
        if attr.Key == "class" && strings.Contains(attr.Value, "btn-primary") {
            hasClass = true
        }
    }

    if !hasClass {
        t.Error("expected btn-primary class")
    }
}
```

---

## Success Criteria

Before marking as complete, verify:

### ✅ Functionality
- [ ] All 7 Phase 1 components implemented
- [ ] Each component uses new DOM fluent API
- [ ] SSR/CSR split with build tags works
- [ ] Components render correctly
- [ ] Explicit registration works
- [ ] Convenience import `all.go` works

### ✅ Testing
- [ ] Each component has unit tests
- [ ] Example app runs and displays all components
- [ ] Run `gotest` in components directory
- [ ] TinyGo build succeeds

### ✅ Binary Size
- [ ] Example WASM <500KB
- [ ] No standard library imports in WASM files
- [ ] CSS embedded in ssr.go only (not in WASM)

### ✅ Documentation
- [ ] COMPONENT_CREATION.md updated
- [ ] COMPONENT_CATALOG.md created
- [ ] Each component has usage examples
- [ ] README.md updated

---

## Files to Create

| File | Description |
|------|-------------|
| `registry.go` | Component registration system |
| `base.go` | Shared utilities |
| `all/all.go` | Convenience import |
| `button/*` | Button component (4 files) |
| `card/*` | Card component (3 files) |
| `input/*` | Input component (4 files) |
| `nav/*` | Navigation component (4 files) |
| `modal/*` | Modal component (4 files) |
| `table/*` | Table component (3 files) |
| `form/*` | Form component (3 files) |
| `docs/COMPONENT_CATALOG.md` | Component documentation |
| `examples/component-showcase/` | Example app |

---

## Implementation Order

1. **Registry system** (foundation)
2. **Button** (simplest, validates pattern)
3. **Input** (validates form patterns)
4. **Card** (validates composition)
5. **Nav** (most complex, moved from site)
6. **Modal** (validates state management)
7. **Table** (validates data rendering)
8. **Form** (validates form integration)
9. **Convenience import** (all.go)
10. **Example app** (documentation)
11. **Tests** (validation)
12. **Documentation** (guides)

---

## Questions/Ambiguities?

If you encounter decisions not covered here:
1. **Read** [COMPONENTS_STRUCTURE.md](./COMPONENTS_STRUCTURE.md) for full context
2. **Check** [COMPONENT_CREATION.md](../components/docs/COMPONENT_CREATION.md) for patterns
3. **Follow principles**: Reusable, composable, minimal

---

## Completion

When done:
1. Commit: `feat(components): implement Phase 1 component catalog (7 components)`
2. Run `gotest` and paste output
3. Measure WASM: `ls -lh examples/component-showcase/client.wasm`
4. Report results and move to Phase 3 (Site)

---

**Status**: Ready to execute after DOM refactor. Begin when DOM Phase 1 complete.
