# PLAN: tinywasm/components — Migración a html/svg/image

## Repositorio
`github.com/tinywasm/components` — path local: `tinywasm/components/`

## Dependencias de ejecución
```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

## Prerequisitos (ejecutar ANTES de este plan)
1. `tinywasm/dom` publicado con: `String()` en lugar de `RenderHTML()`, builders eliminados
2. `tinywasm/html` publicado con todos los builders
3. `tinywasm/svg` publicado con builders + `*Sprite`
4. `tinywasm/image` publicado con `Img()` fluido + `*ImageAsset`

---

## Objetivo

Migrar todos los componentes de `tinywasm/components` para:
1. Usar `. "github.com/tinywasm/html"` en lugar de `. "github.com/tinywasm/dom"` para builders
2. Usar `. "github.com/tinywasm/svg"` para íconos SVG
3. Usar `tinywasm/image` para imágenes (si las hay)
4. Renombrar `.RenderHTML()` → `.String()` en todos los tests

Es un **break change limpio** — sin aliases, sin compatibilidad hacia atrás. El ecosistema no está publicado.

---

## Actualizar go.mod

```
module github.com/tinywasm/components

go 1.25

require (
    github.com/tinywasm/css    v<nueva-version>
    github.com/tinywasm/dom    v<nueva-version>  // aún necesario: Event, Component, Reference
    github.com/tinywasm/html   v<nueva-version>  // NUEVO
    github.com/tinywasm/svg    v<nueva-version>  // NUEVO
    github.com/tinywasm/image  v<nueva-version>  // NUEVO (solo si hay imágenes)
    github.com/tinywasm/fmt    v<version-actual>
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

### Tests ANTES:
```go
html := c.Render().RenderHTML()
```

### Tests DESPUÉS:
```go
html := c.Render().String()
```

---

## Componentes a migrar

### 1. `actionbutton/button.go`

**Imports actuales:** `. "github.com/tinywasm/dom"`, `. "github.com/tinywasm/css"`  
**Cambio:** Agregar `. "github.com/tinywasm/html"`, mantener dom para Event/Component.  
**Builders usados (verificar en el archivo):** `Div`, `Button`, `Span` — mover a html import.

**Tests `actionbutton/button_test.go`:**  
Buscar y reemplazar: `.RenderHTML()` → `.String()`

**Verificar:**
```bash
cd tinywasm/components
gotest -run TestActionButton
```

---

### 2. `contentcard/card.go`

**Cambio:** Mismo patrón. Builders HTML → import html.  
**Tests `contentcard/card_test.go`:**  
- Reemplazar `.RenderHTML()` → `.String()`  
- Tipo `simpleComponent` implementa `dom.Component` — verificar que `RenderHTML() string` → `String() string`

```go
// ANTES en test:
func (s *simpleComponent) RenderHTML() string { return s.html }

// DESPUÉS:
func (s *simpleComponent) String() string { return s.html }
```

**Verificar:**
```bash
gotest -run TestContentCard
```

---

### 3. `datatable/table.go`

**Cambio:** Builders HTML → import html.  
**Verificar:**
```bash
gotest -run TestDataTable
```

---

### 4. `dialog/modal.go`

**Cambio:** Builders HTML → import html.  
**Tests `dialog/modal_test.go`:**  
- Reemplazar `.RenderHTML()` → `.String()`
- Tipo `simpleComponent`: `RenderHTML() string` → `String() string`

**Verificar:**
```bash
gotest -run TestDialog
```

---

### 5. `selectsearch/selectsearch.go` + `selectsearch/svg.go`

**Cambio principal:** Este componente YA tiene `svg.go` con `IconSvg() map[string]string`.  
Migrar a nueva API de `tinywasm/svg`:

**Archivo `selectsearch/svg.go` ANTES:**
```go
//go:build !wasm

package selectsearch

func (c *SelectSearch) IconSvg() map[string]string {
    return map[string]string{
        "ss-arrow-down": `<path fill="currentColor" d="M1.5 4.5l6.5 7 6.5-7H1.5z"/>`,
    }
}
```

**Archivo `selectsearch/svg.go` DESPUÉS:**
```go
//go:build !wasm

package selectsearch

import "github.com/tinywasm/svg"

func (c *SelectSearch) IconSvg() *svg.Sprite {
    return svg.New().
        Add("ss-arrow-down", `<path fill="currentColor" d="M1.5 4.5l6.5 7 6.5-7H1.5z"/>`)
}
```

**Tests `selectsearch/selectsearch_test.go`:**  
Buscar y reemplazar todas las ocurrencias de `.RenderHTML()` → `.String()`

**Verificar:**
```bash
gotest -run TestSelectSearch
```

---

### 6. `themetoggle/themetoggle.go`

**Cambio:** Builders HTML → import html.  
**Verificar:**
```bash
gotest -run TestThemeToggle
```

---

## Actualizar SKILL.md

El archivo `tinywasm/components/docs/SKILL.md` debe actualizarse para reflejar los nuevos imports:

**Sección "Imports — use dot notation":**
```go
// ANTES:
import (
    . "github.com/tinywasm/css"
    . "github.com/tinywasm/dom"
)

// DESPUÉS:
import (
    . "github.com/tinywasm/css"   // Class, Rule, Token, etc.
    . "github.com/tinywasm/html"  // Div, Span, H1, Button, Nav...
    . "github.com/tinywasm/svg"   // Icon(), Svg(), Path()... (si usa íconos)
    . "github.com/tinywasm/dom"   // Event, Component, Reference, Render, Append
)
```

**Sección "File Structure":**
```
mycomponent/
├── mycomponent.go       # Struct, Render(), OnMount() — importa html + dom
├── css.go               # //go:build !wasm: RenderCSS() *css.Stylesheet
├── svg.go               # //go:build !wasm: IconSvg() *svg.Sprite
├── html.go              # //go:build !wasm: RenderHTML() *html.HTML (solo si SSR template custom)
├── image.go             # //go:build !wasm: RenderImages() *image.ImageAsset (solo si hay imágenes)
└── mycomponent_test.go
```

**Sección "Icon Management":**  
Actualizar ejemplo de `IconSvg()` de `map[string]string` a `*svg.Sprite`:
```go
// DESPUÉS:
import "github.com/tinywasm/svg"

func (c *MyComponent) IconSvg() *svg.Sprite {
    return svg.New().
        Add("my-icon-id", `<path fill="currentColor" d="..." />`)
}
```

---

## Verificación completa

```bash
cd tinywasm/components
go build ./...
gotest
```

Todos los tests deben pasar. Si algún test falla por `RenderHTML` o por builders faltantes, es un síntoma de migración incompleta.

## Documentación a Actualizar

### `components/docs/SKILL.md`

Este es el documento más importante — guía la creación de componentes futuros. Actualizar:

1. **Sección "Imports"** — reemplazar el bloque de imports:
   ```go
   // ANTES:
   import (
       . "github.com/tinywasm/css"
       . "github.com/tinywasm/dom"
   )

   // DESPUÉS:
   import (
       . "github.com/tinywasm/css"   // Class, Rule, Token, etc.
       . "github.com/tinywasm/html"  // Div, Span, H1, Nav, Button...
       . "github.com/tinywasm/svg"   // Icon(), Svg() — solo si usa íconos
       . "github.com/tinywasm/dom"   // Event, Component, Reference, Render, Append
   )
   ```

2. **Sección "File Structure"** — actualizar a:
   ```
   mycomponent/
   ├── mycomponent.go       # Struct, Render(), OnMount() — importa html + dom
   ├── css.go               # //go:build !wasm: RenderCSS() *css.Stylesheet
   ├── svg.go               # //go:build !wasm: IconSvg() *svg.Sprite (solo si tiene íconos)
   ├── html.go              # //go:build !wasm: RenderHTML() *html.HTML (solo si SSR template custom)
   ├── image.go             # //go:build !wasm: RenderImages() *image.ImageAsset (solo si hay imágenes)
   └── mycomponent_test.go
   ```

3. **Sección "Icon Management"** — actualizar todo el ejemplo:
   ```go
   // svg.go NUEVO:
   import "github.com/tinywasm/svg"

   func (c *MyComponent) IconSvg() *svg.Sprite {
       return svg.New().
           Add("my-icon-id", `<path fill="currentColor" d="..." />`)
   }

   // Render() — usar svg.Icon():
   import . "github.com/tinywasm/svg"

   svg.Icon("my-icon-id", "my-icon-class")
   // → <svg aria-hidden="true" class="my-icon-class"><use href="#my-icon-id"></use></svg>
   ```

4. **Sección "Tests"** — actualizar todos los ejemplos:
   ```go
   // ANTES:
   html := c.Render().RenderHTML()

   // DESPUÉS:
   html := c.Render().String()
   ```

### `components/docs/CATALOG.md`

Agregar nota al inicio bajo "Theme":
```markdown
> Components import `tinywasm/html` for HTML builders and `tinywasm/svg` for icons.
> See [SKILL.md](./SKILL.md) for the updated import pattern.
```

Ver `tinywasm/docs/MASTER_PLAN.md` para el orden global de ejecución.
