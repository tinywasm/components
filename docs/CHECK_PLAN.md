# PLAN: tinywasm/components — Migración a html/svg

## Repositorio
`github.com/tinywasm/components` — path local: `tinywasm/components/`

## Dependencias de ejecución
```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

## Estado de versiones publicadas

| Paquete | Versión | Estado |
|---------|---------|--------|
| `github.com/tinywasm/dom` | v0.10.1 | ✅ listo |
| `github.com/tinywasm/html` | v0.0.2 | ✅ listo |
| `github.com/tinywasm/svg` | v0.0.2 | ✅ listo |
| `github.com/tinywasm/image` | — | ❌ módulo aún en `github.com/cdvelop/image` — SKIP |

**`tinywasm/image` se omite en esta iteración.** No agregar al go.mod hasta que el plan de image se ejecute.

---

## go.mod

```
module github.com/tinywasm/components

go 1.25

require (
    github.com/tinywasm/css   v<version-actual>
    github.com/tinywasm/dom   v0.10.1
    github.com/tinywasm/html  v0.0.2
    github.com/tinywasm/svg   v0.0.2
    github.com/tinywasm/fmt   v<version-actual>
)
```

---

## Patrón de migración por componente

### Imports ANTES:
```go
import (
    . "github.com/tinywasm/dom"
    . "github.com/tinywasm/css"
)
```

### Imports DESPUÉS:
```go
import (
    . "github.com/tinywasm/html"  // Div, Span, H1, Button, Nav, etc.
    . "github.com/tinywasm/dom"   // Event, Component, Reference, Render, Append
    . "github.com/tinywasm/css"   // Class, Rule, Token, etc.
)
```

---

## Estado de componentes

| Componente | Estado |
|------------|--------|
| `actionbutton` | ✅ migrado |
| `contentcard` | ✅ migrado |
| `datatable` | ✅ migrado |
| `dialog` | ✅ migrado |
| `selectsearch` | ⏳ pendiente — ver detalles abajo |
| `themetoggle` | ⏳ pendiente |

---

## 5. `selectsearch` — pendiente

### Problema: colisión de nombres

`tinywasm/html` exporta `Option(value, text string)`. El componente tenía internamente una struct también llamada `Option` — colisión. **Solución aplicada:** renombrar el struct interno a `SsOption`.

```go
// ANTES:
type Option struct {
    ID          string
    Label       string
    Description string
}

// DESPUÉS:
type SsOption struct {
    ID          string
    Label       string
    Description string
}
```

Actualizar todas las referencias en `selectsearch.go` y tests: `Option` → `SsOption`.

### Import de svg — NO usar dot import

En `selectsearch.go`, usar import **nombrado** para svg (no dot), para evitar colisiones futuras:

```go
import (
    . "github.com/tinywasm/html"       // Input, Label, Div, Span, etc.
    . "github.com/tinywasm/dom"        // Event, Component, etc.
    . "github.com/tinywasm/css"
    "github.com/tinywasm/svg"          // nombrado: svg.Svg(), svg.Use()
)
```

Uso en `Render()`:
```go
// ANTES (dom builders):
Add(dom.Svg(dom.Use().Attr("href", "#ss-arrow-down")).Add(ClsSsIcon.AsAttr()))

// DESPUÉS (svg nombrado):
Add(svg.Svg(svg.Use().Attr("href", "#ss-arrow-down")).Add(ClsSsIcon.AsAttr()))
```

### `selectsearch/svg.go` — migrar a `*svg.Sprite`

```go
//go:build !wasm

package selectsearch

import "github.com/tinywasm/svg"

func (c *SelectSearch) IconSvg() *svg.Sprite {
    return svg.New().
        Add("ss-arrow-down", `<path fill="currentColor" d="M1.5 4.5l6.5 7 6.5-7H1.5z"/>`)
}
```

### Tests de selectsearch

El test de integración wasm fue **movido a `tinywasm/svg/tests/uc_selectsearch_test.go`** porque usa builders SVG. El archivo ya no existe en `components/`.

Si hay tests unitarios en `selectsearch/selectsearch_test.go`:
- Reemplazar `Option{` → `SsOption{`
- Reemplazar `.RenderHTML()` → `.String()`

**Verificar:**
```bash
cd tinywasm/components
gotest -run TestSelectSearch
```

---

## 6. `themetoggle` — pendiente

**Cambio:** Builders HTML → import html. Sin SVG ni imagen.

**Verificar:**
```bash
gotest -run TestThemeToggle
```

---

## Verificación completa

```bash
cd tinywasm/components
go mod tidy
go build ./...
gotest
```

---

## Actualizar SKILL.md

Secciones a actualizar en `tinywasm/components/docs/SKILL.md`:

### Imports
```go
import (
    . "github.com/tinywasm/css"   // Class, Rule, Token, etc.
    . "github.com/tinywasm/html"  // Div, Span, H1, Nav, Button...
    . "github.com/tinywasm/dom"   // Event, Component, Reference, Render, Append
    // svg: solo si el componente usa íconos — preferir import nombrado sobre dot
    "github.com/tinywasm/svg"
)
```

### File Structure
```
mycomponent/
├── mycomponent.go       # Struct, Render(), OnMount() — importa html + dom
├── css.go               # //go:build !wasm: RenderCSS() *css.Stylesheet
├── svg.go               # //go:build !wasm: IconSvg() *svg.Sprite (solo si tiene íconos)
├── html.go              # //go:build !wasm: RenderHTML() string (solo si SSR template custom)
└── mycomponent_test.go
```

### Icon Management
```go
// svg.go — contrato SSR:
//go:build !wasm
import "github.com/tinywasm/svg"

func (c *MyComponent) IconSvg() *svg.Sprite {
    return svg.New().
        Add("my-icon-id", `<path fill="currentColor" d="..." />`)
}

// Render() — referenciar el ícono:
import "github.com/tinywasm/svg"

svg.Icon("my-icon-id", "my-icon-class")
// → <svg aria-hidden='true' class='my-icon-class'><use href='#my-icon-id'></svg>
```

### Tests
```go
// ANTES:
html := c.Render().RenderHTML()

// DESPUÉS:
html := c.Render().String()
```

Ver `tinywasm/docs/MASTER_PLAN.md` para el orden global de ejecución.
