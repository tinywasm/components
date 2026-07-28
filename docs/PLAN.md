---
PLAN: "components: eliminar el último escape hatch (css.Raw) usando el vocabulario de overlay"
EXECUTOR: jules
---

> Este plan se despacha con el flujo CodeJob. Ver skill: **agents-workflow**.

# Plan — `tinywasm/components`: cerrar la puerta de `css.Raw`

## 🚦 0. Bloqueo previo — no empieces sin esto

Este plan requiere **`github.com/tinywasm/widget` v0.3.0 publicado**, que aporta el vocabulario
de overlay: `Backdrop(Scope)`, `Above()`, `Scrim()`, `Hidden()`, `Shown()`.

Plan de esa versión: <https://github.com/tinywasm/widget/blob/main/docs/PLAN.md>

Comprobación obligatoria:

```bash
go get github.com/tinywasm/widget@v0.3.0
go doc github.com/tinywasm/widget/style.Backdrop
```

Si falla, **detente y repórtalo**. No aproximes el overlay con `Fixed()` (es la excepción
*no-reflota*, no posicionamiento) ni conserves `css.Raw` "mientras tanto".

**Nota de nombres:** el constructor se llama `Backdrop`, no `Overlay` — `style.Overlay` ya
existía como el nivel máximo del enum `Elevation` (`Raise(Overlay)`, usado hoy por
`modaldialog`). No lo confundas ni intentes usar `style.Overlay(...)` como función.

---

## ⚠️ 1. Alcance

Quedan **dos** usos de CSS crudo en el repo, ambos introducidos porque el vocabulario no
existía. Este plan los elimina y **cierra la puerta con un test** para que no vuelvan.

**PROHIBIDO:**

| Prohibición | Motivo |
|---|---|
| Dejar cualquier `css.Raw(`, `css.RawItem(` o `css.RawRule(` en el repo | Son las tres puertas del escape hatch. `RawRule` ya estaba vigilada; las otras dos no, y por ahí entró esto. |
| Importar `github.com/tinywasm/css` en un componente | Tras este cambio ningún componente lo necesita. |
| Escribir un selector CSS como literal (`".targetlist__backdrop"`) | Toda clase se deriva de `widget.Name`. |
| Usar `:has(` | La visibilidad la decide un estado escrito por Go (§3.1). |
| Escribir un `z-index`, `position`, `top/left/right/bottom` o `inset` a mano | Lo aportan `Backdrop()` y `Above()`. |
| Tocar la anatomía (`Name`, `Part`) de ningún componente | Ya está migrada y publicada. |
| Usar `go test` | En este repo se usa `gotest`. |

---

## 2. Estado actual: los dos bloques a eliminar

```go
// targetlist/css.go — cazaclics del menú ⋮
func (t *TargetList) RenderCSS() *css.Stylesheet {
	var (
		root     = "." + clsListWrap.String()
		menu     = "." + clsMenu.String()
		backdrop = "." + clsMenuBackdrop.String()
	)
	return css.NewStylesheet(css.Raw(
		backdrop + "{display:none;position:fixed;top:0;left:0;right:0;bottom:0;z-index:4;}" +
			root + ":has(" + menu + "[open]) " + backdrop + "{display:block;}",
	))
}

// modaldialog/css.go — velo del diálogo
func (m *ModalDialog) RenderCSS() *css.Stylesheet {
	backdrop := "." + clsModalBackdrop.String()
	return css.NewStylesheet(css.Raw(
		backdrop + "{position:absolute;top:0;left:0;width:100%;height:100%;" +
			"background-color:color-mix(in srgb, var(--color-surface) 60%, transparent);z-index:1;}",
	))
}
```

---

## Etapa 1 — `modaldialog` (el fácil: no necesita estado)

`modaldialog.Render()` ya envuelve todo en `Show(m.visible, …)`, así que el backdrop **solo
existe en el DOM cuando el diálogo está abierto**. No hace falta `Hidden()`/`Shown()`.

1. **Borra el método `RenderCSS()` entero** de `modaldialog/css.go`, y el import de
   `github.com/tinywasm/css`.
2. En `Style()`, añade la parte:

```go
.Part(PartBackdrop,
	style.Backdrop(style.Parent),
	style.Scrim(),
).
```

3. El panel del diálogo debe quedar por encima del velo — añade `style.Above()` a la parte
   `PartPanel`:

```go
.Part(PartPanel,
	style.Above(),
	// …las opciones que ya tenga
).
```

---

## Etapa 2 — `targetlist` (necesita que Go rastree el estado)

### 2.1 El problema, explicado

Los menús ⋮ son `<details>` **nativos**: los abre y cierra el navegador. Hoy Go **no sabe** si
hay alguno abierto — por eso el CSS usaba `:has(.targetlist__menu[open])`.

Con el vocabulario nuevo la visibilidad se decide por un estado que **escribe Go**
(`widget.Open` → `data-open="true"`). Así que `targetlist` tiene que empezar a rastrearlo. Es la
parte con más trabajo real de este plan; no la saltes ni la aproximes.

### 2.2 Rastrear "hay algún menú abierto"

En `targetlist/targetlist.go`:

1. Añade un campo al struct `TargetList`:
   ```go
   menuOpen *SignalBool
   ```
2. Inicialízalo donde se inicializan los demás signals (junto a `rows`/`Selected`, en `Init`):
   ```go
   t.menuOpen = NewBool(false)
   ```
3. En el `<details>` de cada fila, escucha el evento nativo `toggle` y refleja el estado:
   ```go
   menu.On("toggle", func(e Event) { t.menuOpen.Set(t.anyMenuOpen()) })
   ```
4. Añade el helper, junto a `closeAllMenus()` y con la misma forma (recorre `t.items` y consulta
   el DOM por `menuID`):
   ```go
   // anyMenuOpen informa si alguna fila tiene su ⋮ <details> abierto. El navegador
   // posee ese estado; esto lo refleja en Go para que la hoja pueda seleccionarlo.
   func (t *TargetList) anyMenuOpen() bool {
       for _, it := range t.items {
           if ref, ok := Get(menuID("tl-" + it.ID)); ok {
               if _, open := ref.Attr("open"); open {
                   return true
               }
           }
       }
       return false
   }
   ```
   Si `dom.Element`/`ref` no expone una lectura de atributo con esa firma, **usa la que exista y
   repórtalo en el PR**; no inventes una API en `dom`.
5. `closeAllMenus()` debe dejar el signal en `false` al terminar:
   ```go
   t.menuOpen.Set(false)
   ```

### 2.3 Publicar el estado en el backdrop

En `Render()`, sobre el elemento `backdrop`:

```go
attrOpen := widget.Open.Attr()
backdrop := Div().Set(clsMenuBackdrop.AsAttr()).
	BindAttrFunc(attrOpen.Key, func() string {
		if t.menuOpen.Get() {
			return attrOpen.Value
		}
		return ""
	})
backdrop.On("click", func(Event) { t.closeAllMenus() })
```

**⚠️ Usa `BindAttrFunc`, NUNCA `BindAttrBoolFunc`.** `BindAttrBoolFunc` emite `data-open=""` (sin
valor) y el selector generado es `[data-open="true"]`: no casan, el estilo nunca se aplica y nada
falla. Es el mismo footgun que ya se documentó en `tinywasm/form`.

### 2.4 El estilo

1. **Borra el método `RenderCSS()` entero** de `targetlist/css.go` y el import de
   `github.com/tinywasm/css`.
2. En `Style()`:

```go
.Part(PartBackdrop,
	style.Backdrop(style.Viewport),
	style.Hidden(),
).
.Part(PartOptions,
	style.Above(),
	// …las opciones que ya tenga (Stack, On, Raise, Clip)
).
.When(widget.Open, PartBackdrop,
	style.Shown(),
).
```

El backdrop **no** lleva `Scrim()`: es un cazaclics invisible, no un velo.

---

## Etapa 3 — Cerrar la puerta en el conformance

Archivo: `conformance_test.go`.

Hoy la guardia vigila una sola de las tres puertas:

```go
if ident.Name == "RawRule" { t.Errorf("%s: uses forbidden RawRule", path) }
```

`css.Raw` y `css.RawItem` pasaron sin que saltara. Añade las dos que faltan, junto a las
existentes:

```go
if ident.Name == "Raw" {
	t.Errorf("%s: uses forbidden css.Raw (escape hatch — report the gap upstream in widget/style)", path)
}
if ident.Name == "RawItem" {
	t.Errorf("%s: uses forbidden css.RawItem (escape hatch — report the gap upstream in widget/style)", path)
}
```

Añade también, en el mismo recorrido, un fallo si algún `.go` del repo importa
`github.com/tinywasm/css`: tras este plan ningún componente lo necesita, y ése es el chequeo que
cierra la clase entera de problema en vez de perseguir nombres uno a uno.

---

## Etapa 4 — `go.mod`

`css` deja de ser una dependencia directa. Ejecuta `go mod tidy` y confirma que queda en el
bloque `// indirect` (llega por `widget/style`).

---

## Etapa 5 — Tests

1. Los tests de emparejamiento por paquete que ya existen (`*/‌*_test.go`,
   `TestPairMarkupAndStylesheet`) **deben seguir pasando sin modificarse**. Son bidireccionales:
   cada clase del markup tiene regla y cada regla tiene markup.
2. Añade a `targetlist/targetlist_test.go`:
   - El HTML renderizado **no** contiene `data-open="true"` con todos los menús cerrados.
   - La hoja emitida contiene `display: none` para la parte `backdrop` y `display: block` bajo
     `[data-open="true"]`.
   - La hoja **no** contiene `:has(`.
3. Añade a `modaldialog/modaldialog_test.go`: la hoja contiene `position: absolute` y un
   `background-color` con `color-mix(` para la parte `backdrop`.
4. **Todo archivo de test que llame a `.Style()` debe llevar `//go:build !wasm`** en su primera
   línea. `Style()` vive en `css.go`, que es `!wasm`; sin el tag, la pata WASM de `gotest` no
   compila. Ya pasó una vez.

---

## 6. Criterios de aceptación — verificables con grep

1. `gotest` en verde, incluida la pata WASM (`race ✅, tests ✅`, sin `wasm ❌`).
2. `grep -rn "css.Raw(\|css.RawItem(\|css.RawRule(\|RawRule(" --include='*.go' .` → **vacío**
   fuera de `conformance_test.go` (donde aparecen como cadenas vigiladas).
3. `grep -rn '"github.com/tinywasm/css"' --include='*.go' .` → **vacío**.
4. `grep -n "tinywasm/css" go.mod` → solo en el bloque `// indirect`.
5. `grep -rn ":has(" --include='*.go' .` → **vacío**.
6. `grep -rn "func (.*) RenderCSS()" --include='*.go' .` → **vacío**.
7. `grep -rnE '"\.[a-z]+__' --include='*.go' .` → **vacío** (ningún selector como literal).
8. `grep -rn "z-index\|position:\s*fixed\|position:\s*absolute" --include='*.go' .` → **vacío**.
9. `GOOS=js GOARCH=wasm go build ./...` compila.

---

## 7. Checklist de calidad Go (obligatorio)

- **Sin strings repetidos**: toda clase de `widget.Name`; toda clave de estado de
  `widget.Open.Attr()`.
- **Sin stdlib** en código compartido WASM; DOM solo por `github.com/tinywasm/dom`.
- **Errores** con `github.com/tinywasm/fmt`.
- **Cero `any`, cero `map`** en API nueva.
- `//go:build !wasm` se conserva en todo `css.go` y `svg.go`.

---

## 8. Tabla de etapas

| # | Etapa | Archivos | Gate |
|---|---|---|---|
| 0 | *(bloqueo)* `widget` v0.3.0 publicado | — | `go doc …/style.Backdrop` |
| 1 | `modaldialog` | `modaldialog/css.go` | compila |
| 2 | `targetlist` | `targetlist/css.go`, `targetlist/targetlist.go` | compila |
| 3 | Cerrar la puerta | `conformance_test.go` | compila |
| 4 | Dependencias | `go.mod`, `go.sum` | `go mod tidy` limpio |
| 5 | Tests | `targetlist/targetlist_test.go`, `modaldialog/modaldialog_test.go` | `gotest` verde |

Secuenciales. La 5 es el gate real.
