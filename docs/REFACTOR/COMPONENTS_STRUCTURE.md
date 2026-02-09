# Components - Structure & Registration Plan

**Parent Plan**: [API_STANDARDIZATION.md](./API_STANDARDIZATION.md)
**Library**: `github.com/tinywasm/components`
**Status**: Draft - Depends on DOM Plan Approval
**Priority**: 🟡 High (Blocked by DOM decisions)

---

## Current State Analysis

### What Exists ✅
- `COMPONENT_CREATION.md` guide with clear structure
- Build-tag separation (`ssr.go` with `!wasm`, `front.go` with `wasm`)
- CSS embedding via `//go:embed`
- Icon sprite integration via `IconSvgProvider`
- `BaseComponent` pattern (embed in user components)

### What's Missing ❌
- **No registration mechanism**: Components must be manually passed to `site.RegisterHandlers()`
- **No component catalog**: No way to discover available components
- **Unclear lifecycle ownership**: Is `OnMount` part of DOM or Components?
- **No reusable utilities**: Each component reinvents state management, props, etc.

### Directory Structure (Current)
```
tinywasm/components/
├── go.mod
├── components.go          # Empty
├── palette.go             # CssVars utility
├── docs/
│   └── COMPONENT_CREATION.md
└── (no actual components yet)
```

**Problem**: Library exists but has no components. It's a spec without implementation.

---

## Single Responsibility

**Components' ONLY job**: Provide a catalog of reusable UI components with SSR/CSR support.

**NOT Components' job**:
- DOM manipulation primitives (that's `dom`)
- Application routing (that's `site`)
- Business logic (that's `modules`)
- Data fetching (that's `crudp`)

**Boundary**: Components are "dumb UI" - they receive props, render markup, handle local UI interactions. They don't "know" about the business domain.

---

## Proposed Structure

### Directory Layout

```
tinywasm/components/
├── go.mod
├── registry.go            # Registration system
├── palette.go             # CSS variables (keep)
├── base.go                # Shared utilities for component authors
│
├── button/
│   ├── button.go          # Struct + RenderHTML() or Render()
│   ├── button.css         # Component-specific styles
│   ├── ssr.go             # //go:build !wasm - RenderCSS(), IconSvg()
│   └── front.go           # //go:build wasm - OnMount() if needed
│
├── card/
│   ├── card.go
│   ├── card.css
│   └── ssr.go
│
├── input/
│   ├── input.go
│   ├── input.css
│   ├── ssr.go
│   └── front.go           # For validation, masking, etc.
│
└── nav/                   # Navigation component (from site?)
    ├── nav.go
    ├── nav.css
    ├── ssr.go
    └── front.go
```

**Key Principles**:
1. **One folder per component** (SRP at file level)
2. **Mandatory files**: `{component}.go`, `{component}.css`, `ssr.go`
3. **Optional file**: `front.go` (only if client-side interactivity needed)
4. **Build tags strictly enforced**: SSR code never reaches WASM

---

## Registration Strategy

### Problem Statement

**Goal**: Components must be registered so that:
- `site` can collect their CSS during SSR
- `site` can include their icons in the sprite
- Developers can instantiate them by name (optional)

**Constraint**: Minimize boilerplate. Developers shouldn't write registration code if avoidable.

### Alternatives

#### ⭐ **Alternative A: Explicit Registration (RECOMMENDED)**

**Pattern**: Components are explicitly passed to `site.RegisterHandlers()`.

```go
// In user's server.go / client.go
import (
    "github.com/tinywasm/components/button"
    "github.com/tinywasm/components/card"
    "github.com/tinywasm/site"
)

func main() {
    site.RegisterHandlers(
        &button.Button{},
        &card.Card{},
        // ... user's modules ...
    )
    site.Mount()
}
```

**Pros**:
- ✅ Explicit: Clear what's included (no hidden magic)
- ✅ Tree-shakeable: Unused components don't bloat binary
- ✅ Testable: Can register subsets in tests
- ✅ No init() functions (Go best practice: avoid init)

**Cons**:
- ❌ Boilerplate: Must import and list each component
- ❌ Easy to forget: Developer might forget to register a component

**Mitigation**: Provide a "kitchen sink" import for convenience:
```go
import _ "github.com/tinywasm/components/all" // Registers everything
```

---

#### **Alternative B: Auto-Registration via init()**

**Pattern**: Each component registers itself via `init()`.

```go
// In components/button/button.go
func init() {
    components.Register("button", func() dom.Component {
        return &Button{}
    })
}
```

**Usage** (user's code):
```go
import (
    _ "github.com/tinywasm/components/button" // Auto-registers
    _ "github.com/tinywasm/components/card"
)

func main() {
    site.Mount() // Components already registered
}
```

**Pros**:
- ✅ Zero boilerplate: Just import, auto-registers
- ✅ Can't forget: If you import it, it's registered

**Cons**:
- ❌ Magic: Registration happens invisibly
- ❌ Init order issues: Go doesn't guarantee `init()` order
- ❌ Hard to test: Can't un-register components in tests
- ❌ Not tree-shakeable: Import = includes in binary even if unused

**Justification for rejection**: Goes against Go best practices. Init functions should be rare and only for truly global setup.

---

#### **Alternative C: Code Generation**

**Pattern**: A tool scans `components/` and generates `registry.gen.go`.

```bash
$ tinygo generate ./components
# Generates components/registry.gen.go
```

```go
// registry.gen.go (auto-generated)
package components

func init() {
    autoRegister(
        &button.Button{},
        &card.Card{},
        // ... all discovered components ...
    )
}
```

**Pros**:
- ✅ Zero boilerplate: Automatic discovery
- ✅ Tree-shakeable: Generator can inspect imports
- ✅ No runtime cost: Generated code is static

**Cons**:
- ❌ Requires tooling: Must run generator before build
- ❌ Complexity: Need to write/maintain generator
- ❌ Surprises: "I added a component, why isn't it working?" (forgot to generate)

**Justification for rejection**: Over-engineering for current scale. Revisit if component catalog grows to 50+ components.

---

### 🎯 Recommended Decision

**Use Alternative A (Explicit Registration)** because:
1. **Go-idiomatic**: Explicit is better than implicit
2. **Tree-shakeable**: Only pay for what you use (critical for WASM)
3. **Testable**: Full control over what's registered
4. **Simple**: No tooling, no magic, just function calls

**Compromise**: Provide convenience import for rapid prototyping:
```go
import _ "github.com/tinywasm/components/all"
```

This gives developers choice: explicit control OR convenience.

---

## Component Lifecycle Ownership

### Question: Who Owns OnMount/OnUnmount?

**Context**: Components can have `OnMount()` methods in `front.go`. But `dom` is what calls them.

**Clarification**:
- **DOM owns the mechanism**: `dom` detects `Mountable` interface and calls the hooks
- **Components own the implementation**: Each component decides what to do in `OnMount()`

**Analogy**: DOM is the operating system, components are applications. The OS provides lifecycle hooks (startup, shutdown), apps implement them.

**No change needed**: Current architecture is correct. Components implement `Mountable`, DOM calls it.

---

## Component Base Utilities

### Problem
Component authors might need common patterns:
- Props (passing data to components)
- Children (composing components)
- Default CSS (boilerplate styles)

### Proposed: `components/base.go`

```go
package components

import "github.com/tinywasm/dom"

// BaseComponent provides common utilities for component authors
type BaseComponent struct {
    dom.BaseComponent // Inherit ID, SetID, etc.
}

// WithPrefix sets a semantic ID prefix for debugging
func (c *BaseComponent) WithPrefix(prefix string) *BaseComponent {
    c.SetPrefix(prefix)
    return c
}

// Props pattern (optional)
type Props map[string]any

func (c *BaseComponent) WithProps(p Props) *BaseComponent {
    // Components can type-assert props as needed
    return c
}

// Example component using base utilities
type Button struct {
    BaseComponent
    Text string
    OnClickHandler func(dom.Event)
}

func (b *Button) Render() dom.Node {
    return dom.Button(
        dom.ID(b.ID()),
        dom.Text(b.Text),
        dom.OnClick(b.OnClickHandler),
    )
}
```

**Benefit**: Component authors get sensible defaults without reinventing patterns.

---

## SSR/CSR Split (Critical for Binary Size)

### Current Pattern (from COMPONENT_CREATION.md)

✅ **This is already correct**. Just emphasizing:

```go
// button.go (shared)
package button

type Button struct {
    dom.BaseComponent
    Text string
}

func (b *Button) Render() dom.Node {
    return dom.Button(dom.Text(b.Text))
}
```

```go
// ssr.go (backend only)
//go:build !wasm

package button

import _ "embed"

//go:embed button.css
var css string

func (b *Button) RenderCSS() string {
    return css
}

func (b *Button) IconSvg() map[string]string {
    return map[string]string{
        "btn-icon": `<path d="..." />`,
    }
}
```

```go
// front.go (client only) - ONLY if interactive
//go:build wasm

package button

func (b *Button) OnMount() {
    // Add client-side interactivity if needed
}
```

**Key Rules**:
1. **Shared code** (struct, `Render()`) → `{component}.go`
2. **CSS, Icons** → `ssr.go` with `//go:build !wasm`
3. **Client interactivity** → `front.go` with `//go:build wasm` (optional)

**Enforcement**: CI/CD should verify:
- No `go:embed` in WASM-tagged files
- No CSS strings in shared files
- All components have `ssr.go`

---

## Component Catalog (Initial Set)

### Phase 1: Essential Components (MVP)

**Goal**: Provide 80% of common UI needs.

1. **button/** - Basic button with variants (primary, secondary, danger)
2. **card/** - Container with header/body/footer
3. **input/** - Text input with validation
4. **nav/** - Navigation menu (responsive)
5. **modal/** - Dialog/modal overlay
6. **table/** - Data table with sorting
7. **form/** - Form wrapper with validation

**Why these?**: Most web apps need these primitives. Covers CRUD UIs.

### Phase 2: Advanced Components (Post-MVP)

8. **dropdown/** - Select/dropdown menu
9. **tabs/** - Tabbed interface
10. **toast/** - Notification/toast messages
11. **spinner/** - Loading indicators
12. **badge/** - Small label/badge
13. **avatar/** - User avatar with fallback
14. **breadcrumb/** - Navigation breadcrumb

**Timeline**:
- Phase 1: 1-2 weeks (7 components × 2 days each)
- Phase 2: As needed based on user feedback

---

## Integration with Site

### How Site Uses Components

**Server** (`site.Render()`):
1. Call `handler.registeredModules` to get all registered handlers
2. For each handler implementing `CSSProvider`, collect CSS via `RenderCSS()`
3. For each handler implementing `IconSvgProvider`, collect icons via `IconSvg()`
4. Pass to `assetmin` for minification and serving

**Client** (`site.Mount()`):
1. Hydrate initial module (HTML already in page from SSR)
2. Call `dom.Hydrate()` which calls `OnMount()` hooks

**No changes needed**: Current architecture already supports this. Just need to populate `components/` with actual components.

---

## Props and Children Pattern

### Props (Optional Pattern)

Components can accept props via struct fields:

```go
type Button struct {
    components.BaseComponent
    Text    string
    Variant string // "primary", "secondary", "danger"
}

func (b *Button) Render() dom.Node {
    class := "btn btn-" + b.Variant
    return dom.Button(
        dom.Class(class),
        dom.Text(b.Text),
    )
}

// Usage
btn := &Button{Text: "Submit", Variant: "primary"}
btn.Mount("body")
```

**No framework magic**: Just struct fields. Simple, type-safe, Go-like.

### Children (Composition)

Components can compose via `dom.Component` fields:

```go
type Card struct {
    components.BaseComponent
    Header dom.Component
    Body   dom.Component
}

func (c *Card) Render() dom.Node {
    return dom.Div(
        dom.Class("card"),
        c.Header, // Component as child
        c.Body,
    )
}

// Usage
card := &Card{
    Header: &CardHeader{Title: "Hello"},
    Body: &CardBody{Content: "World"},
}
```

**Note**: This relies on DOM's DSL supporting `Component` as child (already implemented in uncommitted changes).

---

## Testing Strategy

### Unit Tests (No Browser Required)

Test components in isolation:

```go
func TestButton_Render(t *testing.T) {
    btn := &Button{Text: "Click me", Variant: "primary"}
    node := btn.Render()

    // Assert node structure
    if node.Tag != "button" { t.Error("expected button tag") }
    if node.Attrs[0].Value != "btn btn-primary" { t.Error("wrong class") }
}
```

**No DOM needed**: Pure unit test of `Render()` logic.

### Integration Tests (WASM)

Test components in browser using `dom_backend.go` (mock DOM):

```go
//go:build !wasm

func TestButton_OnMount(t *testing.T) {
    // Mock DOM
    dom := newMockDOM()
    btn := &Button{Text: "Click"}

    dom.Render("body", btn)

    // Verify OnMount was called (if implemented)
}
```

**Alternative**: Use `playwright-go` or similar for E2E tests.

---

## Migration from Modules to Components

### Problem
Current `site/example/modules/` has components mixed with modules:
- `contact` is really a **module** (page-level, has business logic)
- But hypothetically, if we had a `button` there, it should move to `components/`

### Rule of Thumb

**Is it a Component or Module?**

| Question | Component | Module |
|----------|-----------|--------|
| Does it talk to backend (CRUD)? | ❌ No | ✅ Yes |
| Is it reusable across pages? | ✅ Yes | ❌ No (page-specific) |
| Does it have local state only? | ✅ Yes | ❌ No (domain state) |
| Lives in `components/`? | ✅ Yes | ❌ No (`modules/`) |

**Example**:
- `Button`, `Card`, `Input` → **Components** (reusable UI)
- `UserManagement`, `ProductCatalog` → **Modules** (pages with backend logic)

**Migration**: No migration needed yet. Just clarify the distinction going forward.

---

## API Summary (Final)

### Package Structure

```go
package components

import "github.com/tinywasm/dom"

// BaseComponent provides utilities for component authors
type BaseComponent struct {
    dom.BaseComponent
}

// Registry (internal, used by site)
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

// All registers all components (convenience)
func All() []any {
    return []any{
        &button.Button{},
        &card.Card{},
        // ... all components ...
    }
}
```

### Component Example (Button)

```go
// button/button.go
package button

import (
    "github.com/tinywasm/components"
    "github.com/tinywasm/dom"
)

type Button struct {
    components.BaseComponent
    Text    string
    Variant string // primary, secondary, danger
}

func (b *Button) Render() dom.Node {
    return dom.Button(
        dom.ID(b.ID()),
        dom.Class("btn btn-"+b.Variant),
        dom.Text(b.Text),
    )
}

// ssr.go
//go:build !wasm

//go:embed button.css
var css string

func (b *Button) RenderCSS() string { return css }

// front.go (optional)
//go:build wasm

func (b *Button) OnMount() {
    // Client-side enhancements if needed
}
```

### Usage

```go
// In user's code
import "github.com/tinywasm/components/button"

btn := &button.Button{
    Text: "Submit",
    Variant: "primary",
}
btn.Mount("body")
```

**Lines of code**: ~5 to use a component (vs ~20 to create one from scratch).

---

## Open Questions for Approval

### ❓ Q1: Registration Strategy
**Approve Alternative A (Explicit Registration)?**
- [ ] Yes, use explicit registration
- [ ] No, use auto-registration (Alternative B)
- [ ] No, use code generation (Alternative C)

### ❓ Q2: Convenience Import
**Provide `components/all` for quick imports?**
- [ ] Yes, include it (easy to use, opt-in)
- [ ] No, force explicit imports (tree-shaking only)

### ❓ Q3: Component Priority
**Which components should be Phase 1 (MVP)?**
- [ ] Approve the list (button, card, input, nav, modal, table, form)
- [ ] Add: _______________
- [ ] Remove: _______________

### ❓ Q4: Props Pattern
**Use struct fields for props (simple, Go-like)?**
- [ ] Yes, use struct fields
- [ ] No, propose alternative: _______________

---

## Next Steps After Approval

1. Create `components/registry.go` with registration system
2. Implement Phase 1 components (button, card, input, nav, modal, table, form)
3. Update `COMPONENT_CREATION.md` with new patterns
4. Create examples for each component
5. Write unit tests
6. Integrate with `site` (already compatible)

**Estimated Effort**:
- Registry: 1 day
- Each component: ~2 days (design + implement + test + docs)
- Phase 1 total: ~2 weeks

---

**Ready for approval?** Once you answer the questions above, I'll finalize the [SITE_ORCHESTRATION.md](./SITE_ORCHESTRATION.md) plan.
