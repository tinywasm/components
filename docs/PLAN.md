# PLAN: `ThemeSwitch` — `tinywasm/components/themeswitch`

## Contexto

Componente visual para alternar el tema de la aplicación (light / dark / auto).
Usa `dom.SetDocumentAttr`/`dom.GetDocumentAttr` para escribir `data-theme` en
`<html>`, y `dom.LocalStorage*` para persistir la preferencia del usuario.

Se inyecta como botón flotante fijo en la esquina superior derecha.

Propósito principal: **herramienta de desarrollo** para probar componentes en
diferentes modos sin cambiar la preferencia del OS.

**Este paquete posee:**
- El tipo `Theme` y sus constantes (`ThemeAuto/ThemeDark/ThemeLight`)
- Los bloques CSS `[data-theme="light"]` y `[data-theme="dark"]`
- La lógica de ciclo (`auto → dark → light → auto`) y persistencia

`dom` solo expone las primitivas del DOM (`SetDocumentAttr`, `GetDocumentAttr`,
`LocalStorage*`). La semántica de tema vive aquí.

---

## Decisiones confirmadas

| Decisión | Valor |
|----------|-------|
| Nombre | `ThemeSwitch` — paquete `themeswitch` |
| Implementación | `Render() *dom.Element` con `.On("click", ...)` — canónico |
| CSS | Tokens `--color-*` del tema activo + bloques `[data-theme]` propios |
| Estados | 3: `auto → dark → light → auto` (ciclo local, sin depender de `dom`) |
| `ThemeAuto` | `Theme("")` — mapeado directo a "sin atributo"; coherente con `GetDocumentAttr` retornando `""` cuando no hay atributo |
| Posición | `fixed`, esquina superior derecha |

---

## Prerequisito de instalación

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

`gotest` corre los tests en browser real — necesario porque `Init()` y el handler
de click usan `dom.LocalStorage*` que solo funciona en entorno WASM.

---

## Estructura de archivos

```
components/themeswitch/
├── themeswitch.go         — Theme type + struct + Render() + helpers puros (sin build tag)
├── themeswitch_wasm.go    — Init() + onClick() — usan dom.LocalStorage* y dom.SetDocumentAttr
├── themeswitch_backend.go — Init() + onClick() no-op para SSR
├── themeswitch.css        — .ts-btn + bloques [data-theme="light"] y [data-theme="dark"]
├── ssr.go                 — //go:build !wasm — RenderCSS(), IconSvg()
├── themeswitch_test.go    — tests SSR (label, cycle, valid, RenderCSS)
├── uc_themeswitch_test.go — tests WASM en browser real (API pública)
├── web/
│   └── client.go
└── README.md
```

**Por qué split `_wasm.go` / `_backend.go` para `Init/onClick`:** ambas funciones
llaman a `dom.LocalStorage*` que es wasm-only (sin stub backend, decisión del PLAN
de `dom`). Para que `themeswitch.go` (sin tag) compile en SSR, las operaciones de
storage se aíslan en archivos taggeados.

**Sin `syscall/js` en ningún archivo.** Todo el acceso al browser pasa por `dom`.

---

## Implementación

### `themeswitch.go` (sin build tag) — Theme type + struct + Render + helpers puros

```go
package themeswitch

import "github.com/tinywasm/dom"

// storageKey identifica la entrada de localStorage del componente.
const storageKey = "tinywasm-themeswitch"

// Theme representa el estado de tema del componente.
// ThemeAuto ("") = sin atributo data-theme → OS preference via @media.
// Los valores "dark" y "light" se escriben literalmente en data-theme.
type Theme string

const (
    ThemeAuto  Theme = ""      // no data-theme attribute → OS preference
    ThemeDark  Theme = "dark"
    ThemeLight Theme = "light"
)

// ThemeSwitch es un botón flotante que cicla entre los 3 modos de tema.
//
//   ts := &themeswitch.ThemeSwitch{}
//   ts.Init()                  // restaura tema guardado (no-op en SSR)
//   dom.Append("body", ts)
type ThemeSwitch struct {
    dom.Element
}

func (t *ThemeSwitch) Render() *dom.Element {
    current := Theme(dom.GetDocumentAttr("data-theme"))
    return dom.Button(label(current)).
        Class("ts-btn").
        On("click", t.onClick) // implementado por build tag
}

// cycle define el orden de los 3 estados. Switch (no map) — TinyGo.
func cycle(current Theme) Theme {
    switch current {
    case ThemeDark:
        return ThemeLight
    case ThemeLight:
        return ThemeAuto
    default: // ThemeAuto ("") o cualquier valor inesperado
        return ThemeDark
    }
}

// label retorna el texto visible del botón. Switch (no map) — TinyGo.
func label(theme Theme) string {
    switch theme {
    case ThemeDark:
        return "🌙 dark"
    case ThemeLight:
        return "☀ light"
    default: // ThemeAuto ("")
        return "☀/🌙 auto"
    }
}

func valid(t Theme) bool {
    return t == ThemeAuto || t == ThemeDark || t == ThemeLight
}
```

**Por qué `ThemeAuto = ""`:**
`dom.GetDocumentAttr("data-theme")` retorna `""` cuando el atributo no existe.
`Theme("") == ThemeAuto` — no se necesita conversión ni comparación especial.
`dom.SetDocumentAttr("data-theme", "")` elimina el atributo — coherente con
la semántica de "modo auto = sin override". El ciclo y la persistencia funcionan
sin casos especiales.

### `themeswitch_wasm.go` (`//go:build wasm`) — Init + click handler

```go
//go:build wasm

package themeswitch

import "github.com/tinywasm/dom"

// Init restaura el tema guardado en localStorage. Si el storage no está disponible
// o la entrada está corrupta, sale limpiamente sin modificar el tema.
func (t *ThemeSwitch) Init() {
    if !dom.LocalStorageAvailable() {
        return // storage bloqueado — modo auto por defecto
    }
    saved, err := dom.LocalStorageGet(storageKey)
    if err != nil || saved == "" {
        return
    }
    theme := Theme(saved)
    if !valid(theme) {
        dom.LocalStorageDel(storageKey) // best-effort cleanup, error ignorado
        return
    }
    dom.SetDocumentAttr("data-theme", string(theme))
}

func (t *ThemeSwitch) onClick(dom.Event) {
    current := Theme(dom.GetDocumentAttr("data-theme"))
    next := cycle(current)
    dom.SetDocumentAttr("data-theme", string(next)) // "" elimina el atributo para ThemeAuto
    // Persistencia best-effort: el tema se aplica aunque el storage falle.
    if next == ThemeAuto {
        dom.LocalStorageDel(storageKey) // error ignorado — tema ya aplicado
    } else {
        dom.LocalStorageSet(storageKey, string(next)) // error ignorado — tema ya aplicado
    }
    t.Update()
}
```

### `themeswitch_backend.go` (`//go:build !wasm`) — stubs SSR

```go
//go:build !wasm

package themeswitch

import "github.com/tinywasm/dom"

// En SSR no hay localStorage ni clicks. Init es no-op para que el código de
// la app que llama ts.Init() compile y se comporte coherentemente en build !wasm.
func (t *ThemeSwitch) Init()             {}
func (t *ThemeSwitch) onClick(dom.Event) {}
```

### `themeswitch.css`

El CSS tiene dos responsabilidades:
1. Estilos del botón `.ts-btn`
2. Definir los overrides de tema `[data-theme]` — movidos desde `dom/theme.css`

```css
/* ── Botón flotante ────────────────────────────────────────────── */
.ts-btn {
    position: fixed;
    top: 1rem;
    right: 1rem;
    z-index: 9999;
    padding: var(--mag-pri) calc(var(--mag-pri) * 2);
    border-radius: 0.4em;
    border: none;
    cursor: pointer;
    font-size: 0.85rem;
    background: var(--color-secondary);
    color: var(--color-primary);
    opacity: 0.9;
    transition: opacity 0.2s;
}

.ts-btn:hover {
    opacity: 1;
    background: var(--color-selection);
    color: var(--color-gray);
}

/* ── Manual light override ([data-theme="light"] on <html>) ─────── */
[data-theme="light"] {
    --color-primary:    #1C1C1E;
    --color-secondary:  #00ADD8;
    --color-tertiary:   #6E6E73;
    --color-quaternary: #F2F2F7;
    --color-gray:       #FFFFFF;
    --color-selection:  #654FF0;
    --color-hover:      #B8860B;
}

/* ── Manual dark override ([data-theme="dark"] on <html>) ──────── */
[data-theme="dark"] {
    --color-primary:    #E6EDF3;
    --color-secondary:  #00ADD8;
    --color-tertiary:   #8B949E;
    --color-quaternary: #161B22;
    --color-gray:       #0D1117;
    --color-selection:  #654FF0;
    --color-hover:      #F7DF1E;
}
```

**Por qué los bloques `[data-theme]` viven aquí y no en `dom/theme.css`:**
El modo auto funciona solo con `@media prefers-color-scheme` (que permanece en
`dom/theme.css`). Los overrides manuales solo son necesarios cuando `ThemeSwitch`
está en la app. Moverlos aquí evita que `dom` cargue CSS para una feature que
puede no usarse, y mantiene toda la lógica de tema en este paquete.

**Orden de carga:** el servidor SSR inyecta `themeswitch.RenderCSS()` en `<head>`
antes del WASM — los tokens están disponibles desde el primer frame.

**Por qué `--color-secondary` como fondo del botón:** es el cyan de Go (`#00ADD8`),
siempre presente y reconocible en ambos modos. `--color-primary` como texto adapta
el contraste automáticamente.

### `ssr.go` (`//go:build !wasm`)

```go
//go:build !wasm

package themeswitch

import _ "embed"

//go:embed themeswitch.css
var css string

func (t *ThemeSwitch) RenderCSS() string          { return css }
func (t *ThemeSwitch) IconSvg() map[string]string { return nil }
```

---

## Uso desde `web/client.go` de cualquier componente

```go
//go:build wasm

package main

import (
    "github.com/tinywasm/components/themeswitch"
    . "github.com/tinywasm/dom"
)

func main() {
    ts := &themeswitch.ThemeSwitch{}
    ts.Init()          // restaura tema guardado antes del primer render
    Render("app", &App{})
    Append("body", ts)
    select {}
}
```

---

## Restricciones TinyGo a respetar

- **Sin `map[string]string`** en código que compile a WASM — usar `switch`/`if-else`.
- **Sin stdlib Go** (`strings`, `fmt`, etc.) — usar `github.com/tinywasm/fmt`.
- **Embed solo en `!wasm`** — los `//go:embed` van en `ssr.go`.
- El struct **embebe `dom.Element` como valor**, nunca como puntero.

---

## Tests requeridos

### `themeswitch_test.go` (`//go:build !wasm`) — lógica pura

| Test | Verifica |
|------|----------|
| `TestLabel_AllThemes_NonEmpty` | `label()` retorna string no vacío para los 3 themes |
| `TestCycle_AutoToDark` | `cycle(ThemeAuto) == ThemeDark` |
| `TestCycle_DarkToLight` | `cycle(ThemeDark) == ThemeLight` |
| `TestCycle_LightToAuto` | `cycle(ThemeLight) == ThemeAuto` |
| `TestValid_AcceptsAllThemes` | `valid()` retorna `true` para los 3 valores legales |
| `TestValid_RejectsInvalid` | `valid("xyz")` retorna `false` |
| `TestRenderCSS_NotEmpty` | `(&ThemeSwitch{}).RenderCSS()` retorna el CSS embebido |

### `uc_themeswitch_test.go` (`//go:build wasm`) — integración browser

| Test | Verifica |
|------|----------|
| `TestThemeSwitch_Init_NoSavedValue_StaysAuto` | Sin entrada en localStorage → `dom.GetDocumentAttr("data-theme") == ""` |
| `TestThemeSwitch_Init_RestoresDark` | Pre-cargar `"dark"` en localStorage + `Init()` → `dom.GetDocumentAttr("data-theme") == "dark"` |
| `TestThemeSwitch_Init_InvalidValue_Cleans` | Pre-cargar `"xyz"` + `Init()` → entrada borrada de localStorage |
| `TestThemeSwitch_Click_CyclesAndPersists` | Simular click → `dom.GetDocumentAttr("data-theme")` avanza ciclo + `LocalStorageGet` refleja el nuevo valor |
| `TestThemeSwitch_Click_AutoDeletesEntry` | Llegar a `ThemeAuto` por click → `v, err := dom.LocalStorageGet(storageKey); v == "" && err == nil` |

**Convención de ubicación:**
- `uc_themeswitch_test.go` — tests de API pública (`package themeswitch_test`). Todos los tests WASM van aquí.
- `themeswitch_test.go` — tests de lógica pura que usan internals (`package themeswitch`).

**Cleanup en cada test WASM:** `dom.LocalStorageDel(storageKey)` y
`dom.SetDocumentAttr("data-theme", "")` en setUp/tearDown para aislamiento.

---

## Checklist de implementación

Prerequisito: `go install github.com/tinywasm/devflow/cmd/gotest@latest`

**Bloqueador:** este componente depende de `LocalStorage*`, `SetDocumentAttr` y
`GetDocumentAttr` de `tinywasm/dom`. Implementar `dom` primero
(ver `dom/docs/PLAN.md`). Verificar también que `dom/theme.css` ya no tiene los
bloques `[data-theme]` antes de añadirlos a `themeswitch.css`.

- [ ] Crear `themeswitch.go` (sin tag) — `type Theme string` + constantes, `ThemeSwitch`, `Render()`, `cycle()`, `label()`, `valid()`, `storageKey`
- [ ] Crear `themeswitch_wasm.go` (`wasm`) — `Init()` y `onClick()` con `dom.LocalStorage*` y `dom.SetDocumentAttr/GetDocumentAttr`
- [ ] Crear `themeswitch_backend.go` (`!wasm`) — `Init()` y `onClick()` no-op
- [ ] Crear `themeswitch.css` con `.ts-btn` + bloques `[data-theme="light"]` y `[data-theme="dark"]`
- [ ] Crear `ssr.go` (`!wasm`) con `RenderCSS()` e `IconSvg()`
- [ ] Crear `themeswitch_test.go` (`!wasm`) con los 7 tests de lógica pura (raíz — accede a helpers internos)
- [ ] Crear `uc_themeswitch_test.go` (`wasm`) con los 5 tests de integración browser
- [ ] Crear `web/client.go` con ejemplo de uso
- [ ] Crear `README.md`
- [ ] Verificar con `gotest` (corre browser + ambos build tags)
- [ ] Actualizar `selectsearch/web/client.go` para usar `ThemeSwitch` como demo
- [ ] `gopush 'feat(components): add ThemeSwitch component'`
