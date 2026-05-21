# PLAN — Eliminar `SSRInstance()` de los componentes

## Objetivo

Quitar el boilerplate `func SSRInstance() *T { return &T{} }` de cada
sub-componente. El contrato real de un componente SSR es la firma de sus
métodos `Render*`; la función accesoria sólo existe para que el extractor
externo construya el receiver, lo cual puede deducirse de la firma del método.

## Justificación

- `SSRInstance` no aporta lógica — cuerpo siempre `return &T{}`.
- Reduce un símbolo público por componente y baja la barrera de entrada para
  añadir uno nuevo.
- El contrato implícito (`RenderCSS() *Stylesheet` sobre un tipo receiver
  zero-value-safe) ya es la convención de facto en este paquete.

## Precondición técnica

Borrar `SSRInstance` rompe la compilación del extractor SSR upstream mientras
éste invoque `pkg.SSRInstance()` en el `main.go` que genera. Aplicar este plan
sólo cuando el extractor ya descubra el tipo receiver automáticamente desde la
firma del método (o acepte su ausencia como fallback). Verificación previa:

```bash
# El extractor debe procesar un módulo sin SSRInstance sin error.
go test ./... -run TestExtract_NoSSRInstanceFunction
```

## Componentes afectados

| Componente | Archivo | Tipo |
|---|---|---|
| dialog | `dialog/ssr.go` | `*DialogWidget` |
| themetoggle | `themetoggle/ssr.go` | `*ThemeToggle` |
| selectsearch | `selectsearch/ssr.go` | `*SelectSearch` |
| contentcard | `contentcard/ssr.go` | `*ContentCard` |
| navbar | `navbar/ssr.go` | `*NavBar` |
| datatable | `datatable/ssr.go` | `*DataTable` |
| actionbutton | `actionbutton/ssr.go` | `*ActionButton` |

En cada uno: borrar la función `SSRInstance` y los imports que queden huérfanos.

## SKILL.md

[docs/SKILL.md](SKILL.md) menciona `SSRInstance` como parte de "embedding rules"
/ "PLAN.md workflow". Acciones:

- Eliminar la regla que exige declarar `SSRInstance`.
- Documentar que el tipo receiver se infiere de la firma del método `Render*`.
- Actualizar el snippet "componente mínimo" — sólo `RenderCSS` (y opcionalmente
  `RenderHTML`, `RenderJS`, `IconSvg`).
- Sincronizar con `llmskill` para propagar el skill actualizado.

## Tests

- `grep -rn SSRInstance components/` debe quedar vacío tras los cambios.
- `go test ./...` en `tinywasm/components` verde.
- Verificación visual vía MCP browser de los componentes (themetoggle, navbar,
  datatable) cargados en una página SSR — sin regresiones.

## Migración adicional: `RenderJS()` → `[]*js.Script` (breaking change)

Como parte del mismo barrido de simplificación del contrato SSR, los
componentes que implementen `RenderJS()` deben migrar de `string` a
`[]*js.Script` (ver tipo `Script` en `github.com/tinywasm/js`). Reglas:

- Bundle global: `&js.Script{Content: content}` (Name vacío).
- Standalone crudo (escape hatch): `&js.Script{Name: "raw.js", Content: content}`.
- **Recomendado para SW/Worker:** constructores tipados de `tinywasm/js`
  (`js.ServiceWorker(name, handler)`, `js.WebWorker(name, handler)`) —
  el componente implementa la lógica como interfaz Go y el framework genera
  el JS-shim. Cero JS escrita en el componente.

Estado actual: ningún componente del paquete implementa `RenderJS()` todavía,
así que la migración aquí se limita a:

- Documentar la firma nueva en `docs/SKILL.md` cuando se hable de capacidades
  opcionales.
- Si se añade un componente con JS (p.ej. un widget con worker) usar ya la
  firma nueva.

Precondición técnica: `tinywasm/js`, `tinywasm/dom` y `tinywasm/assetmin`
deben estar publicados con el contrato `[]*js.Script`. Verificación:

```bash
go list -m github.com/tinywasm/js github.com/tinywasm/dom github.com/tinywasm/assetmin
```

## Migración adicional: `ssr.go` → split por extensión (breaking change)

El motor de `assetmin` deja de reconocer el nombre `ssr.go` y pasa a descubrir
assets por archivos con nombre de extensión, todos `//go:build !wasm`:
`css.go` (RootCSS/RenderCSS), `js.go` (RenderJS), `html.go` (RenderHTML),
`svg.go` (IconSvg). Ver el stage homónimo en `assetmin/docs/PLAN.md`.

Cada componente de este paquete tiene `RenderCSS()` + `IconSvg()` en su
`ssr.go`, así que **se parte en dos archivos**: `css.go` + `svg.go`. El
contenido se mueve literal (mismo receiver, mismo `//go:embed`); no cambia
lógica, solo ubicación.

Componentes afectados (cada `ssr.go` → `css.go` + `svg.go`):
`selectsearch`, `actionbutton`, `dialog`, `themetoggle`, `contentcard`,
`navbar`, `datatable`.

`docs/SKILL.md` (skill `components` en devflow) debe actualizar la sección
"File Structure" y "SSR File" para reflejar los nuevos nombres y resincronizar
con `llmskill`.

Precondición: `assetmin` publicado con la whitelist `ssrSourceFiles`. Aplicar
en el mismo PR coordinado que el cambio de motor para no dejar el extractor
roto.

## Stages

| # | Tarea | Done |
|---|---|---|
| 1 | Confirmar precondición técnica (test de extractor sin SSRInstance) | [ ] |
| 2 | Borrar `SSRInstance` en los 7 componentes listados | [ ] |
| 3 | Actualizar `docs/SKILL.md` y sincronizar con `llmskill` | [ ] |
| 4 | Confirmar precondición `[]*js.Script` publicada en js/dom/assetmin | [ ] |
| 5 | Documentar firma nueva de `RenderJS()` en `docs/SKILL.md` | [ ] |
| 6 | `go test ./...` verde en `tinywasm/components` | [ ] |
| 7 | Verificación visual vía MCP browser | [ ] |
| 8 | Confirmar precondición: `assetmin` con whitelist `ssrSourceFiles` publicado | [ ] |
| 9 | Partir cada `ssr.go` en `css.go` + `svg.go` (7 componentes) | [ ] |
| 10 | Actualizar "File Structure" + "SSR File" en `docs/SKILL.md` y resync `llmskill` | [ ] |
| 11 | `go test ./...` verde tras el split | [ ] |
