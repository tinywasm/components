---
PLAN: "feat!: listselect lego piece replaces the per-row options menu"
TAG: v0.6.0
EXECUTOR: jules
REVIEWER: none
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.
>
> Forma parte de una ola: `docs/BULK_ACTIONS_MASTER_PLAN.md` en la raíz del
> monorepo. Este plan es **independiente**. Es **puerta** para `layout`.
>
> **Es un cambio breaking.** Ver §8.

# Plan — Modo selección como pieza lego compartida

## 0. Prerrequisito

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

Tests con `gotest`. **Nunca `go test`.**

## 1. Por qué

Hoy cada fila dibuja su propio `⋮` que despliega, en la propia fila, un menú
con una única opción: Eliminar. Tres defectos: N botones idénticos para una
acción que casi nunca se usa; sólo actúa sobre un registro; y no hay dónde
crecer.

El menú por fila se va. En su hueco entra un **check**, visible sólo en modo
selección. Las acciones se mudan al pie, que las pone `crudview` (otro plan).

```
MODO NORMAL      [ Pc Taller                    192.168.122.30 ]
                 tap → OnSelect (carga la ficha)

MODO SELECCIÓN   [ ☑ Pc Taller                  192.168.122.30 ]
                   ↑ el check ocupa el hueco donde vivía el ⋮
                 tap en CUALQUIER punto de la fila → marca/desmarca
```

## 2. La decisión de arquitectura: una pieza lego, no código duplicado

`targetlist` y `targetdate` necesitan **exactamente** el mismo modo selección.
La tentación es escribirlo dos veces. **No lo hagas.**

Este repositorio ya resolvió este problema una vez: `components/listgap` es una
pieza lego que posee el contenedor de lista compartido, y los dos componentes
la ensamblan en vez de re-declarar su forma. Lee
`components/listgap/listgap.go` antes de empezar — su comentario de paquete
explica el patrón y por qué los valores están sin exportar.

Haz lo mismo aquí: **crea `components/listselect`**, y que los dos componentes
la ensamblen.

Esto no es preferencia de estilo. Es la regla del harness:

> **The glue is written once, in the library that owns it.** If every
> application would write the same wiring, that wiring belongs to a piece —
> not to the applications.

Y su consecuencia directa: si el estado de selección vive copiado en dos
sitios, el día que uno arregle un fallo el otro se queda con él.

## 3. Etapa 1 — La pieza `components/listselect`

Paquete nuevo, **plano** (sin subcarpetas), dividido por extensión:

- **`components/listselect/listselect.go`** — neutral (sin build tag): el
  estado y el comportamiento. Compila a WASM.
- **`components/listselect/css.go`** — `//go:build !wasm`: el skin del check.

### `listselect.go`

```go
// Package listselect owns one concern: the multi-selection mode a record list
// enters when its host is about to act on several rows at once. It is a lego
// piece — targetlist and targetdate assemble it, they do not re-declare it.
//
// Sibling of listgap, and for the same reason: two lists that must stay
// visually and behaviourally interchangeable (crudview swaps one for the
// other) cannot each own a private copy of the rule.
package listselect

// Mode is the selection state of one list. The zero value is a usable list in
// normal mode — a list is a list until its host says otherwise.
type Mode struct {
	on *dom.SignalBool

	// checked is a SLICE, never a map. A map pulls TinyGo's hashing and
	// runtime machinery into the WASM binary, and this type compiles to the
	// browser. A record list holds tens of rows, so a linear scan costs
	// microseconds — the map was never buying anything here.
	//
	// []string and not []fmt.KeyValue: this is a SET, so there is no value to
	// carry. fmt.KeyValue is the map-free shape for an actual key→value pair;
	// inventing a Value of "true" would be a stringly-typed bool.
	checked []string

	// OnChange fires after every toggle with the current count, so a host can
	// label its commit button ("🗑 3") and disable it at zero.
	OnChange func(n int)
}

// On reports the signal a component binds its root state to, so the stylesheet
// can reveal the checks. Never a bool: the skin has to react.
func (m *Mode) On() *dom.SignalBool

// SetOn enters or leaves selection mode. Leaving ALWAYS clears the marks:
// a mode the user cancelled must not leave a hidden selection behind for the
// next entry to inherit silently.
func (m *Mode) SetOn(on bool)

// Toggle marks or unmarks one id and fires OnChange.
func (m *Mode) Toggle(id string)

// IsChecked answers for one id — what a row binds its check state to.
func (m *Mode) IsChecked(id string) bool

// CheckedIDs returns the marked ids in the order given by ids, which the
// caller passes as its current render order.
//
// Ordering is NOT optional and NOT the caller's problem to remember: checked
// accumulates in TAP order, and a host building a confirmation message from
// tap order would list rows in an order that matches nothing on screen.
// Taking the render order as a parameter is what makes the wrong version
// unwritable — there is no accessor that returns the raw slice.
func (m *Mode) CheckedIDs(ids []string) []string
```

**Deliberadamente NO hay `ClearChecked()`.** `SetOn(false)` ya limpia, y el
único consumidor (`crudview`) siempre sale del modo al terminar. Dos caminos
para lo mismo violan el principio 4 del harness ("one way to do each thing") y
exportarían maquinaria que nadie usa (principio 5, "minimal surface"). Si algún
día hiciera falta limpiar sin salir del modo, se añade **entonces**, con su
consumidor delante.

**Anti-footgun — dos, y las dos son del binario WASM:**

1. **Nada de mapas.** `listselect.go` es neutral y compila a WASM; un `map[...]`
   arrastra la maquinaria de hashing de TinyGo al binario. Usa slices. Cuando
   de verdad necesites pares clave→valor, el tipo del ecosistema es
   `fmt.KeyValue` (`{Key, Value string}`) en un slice — nunca un mapa. Para un
   conjunto, como aquí, la forma es un `[]string` a secas.
2. **Nada de stdlib.** Usa `github.com/tinywasm/fmt` si necesitas formatear.
   `tinywasm/fmt` **no** tiene `Itoa` ni `FormatInt` — para intercalar un
   entero, `fmt.Sprintf`.

### `css.go`

```go
// Apply adds the check part's skin to s and returns s for chaining. Both
// target* lists call this instead of hand-writing the block, so the check
// cannot drift between them.
func Apply(s *style.Sheet, check widget.Part) *style.Sheet
```

Requisitos del check — el usuario pidió literalmente **"muy distintivo"**:

- `IconBox(IconMd)` + `KeepSize()`: cuadrado que no encoge cuando la fila se
  estrecha.
- Reposo: `As(Inset)` + `Round(RadiusSm)` — un hueco que se lee como
  "aquí se puede marcar".
- Marcado: `As(Accent)`. El ámbar es el único "esto está elegido" del chasis;
  la fila seleccionada y el ítem de navegación actual ya lo llevan.
- `Animate(MotionFast)`: marcar y desmarcar es una transición, no un
  parpadeo. `MotionFast` es el escalón nombrado para *"immediate highlight"*.
- **Oculto salvo en modo selección.** El componente enlaza su raíz con
  `BindStateFunc(widget.Open, …)` leyendo `Mode.On()`, y la hoja revela el
  check desde ese estado de la raíz. Mira cómo lo hace
  `components/selectsearch/css.go` con `WhenWithin(widget.Open, "", …)` antes
  de inventar nada.

## 4. Etapa 2 — Borrar el menú por fila

En **`components/targetlist/targetlist.go`** y
**`components/targetdate/targetdate.go`**, los dos:

| Se borra | Dónde |
|---|---|
| `PartOptions`, `PartButton`, `PartIcon`, `PartItemDanger`, `PartItemIcon`, `PartItemLabel` | bloque `const` de partes |
| `clsMenuBtn`, `clsMenuIcon`, `clsMenuList`, `clsMenuItemDanger`, `clsMenuItemIcon`, `clsMenuItemLabel` | bloque `var` de clases |
| `openMenu *SignalString` | campo del struct |
| `closeAllMenus()` y `CloseMenus()` | métodos |
| El `trigger` (`⋮`), el `del` (Eliminar) y el div `options` | dentro de `buildRow` |
| `OnDelete func(id string)` | campo del struct — el borrado ya no nace en la fila |
| `iconDots`, `iconDelete` | constantes de icono |
| Sus entradas de sprite | `svg.go` |
| Todas las reglas de esas partes | `css.go` |

**Criterio verificable:**

```
grep -rn "openMenu\|CloseMenus\|PartItemDanger\|iconDots\|iconDelete\|OnDelete" components/targetlist components/targetdate
```

no devuelve nada.

Un agente flojo tiende a dejar el código viejo "por si acaso". Aquí no: si
queda un `⋮` en la fila, el rediseño no sirve, porque el usuario ve dos sitios
para la misma acción.

## 5. Etapa 3 — Ensamblar la pieza en los dos componentes

Idéntico en ambos. En el struct:

```go
sel listselect.Mode
```

Superficie pública de cada componente — **exactamente estos tres métodos**, con
la misma firma en los dos, porque `crudview` los intercambia:

```go
func (t *TargetList) SetSelectMode(on bool)        { t.sel.SetOn(on) }
func (t *TargetList) CheckedIDs() []string          { ... }
func (t *TargetList) OnCheckedChange(fn func(int))  { t.sel.OnChange = fn }
```

`CheckedIDs()` construye el orden de render a partir de `t.items` y se lo pasa
a `t.sel.CheckedIDs(order)` — así el orden sale del componente, que es quien lo
sabe, y la pieza no tiene que adivinarlo.

En `Render()`: `BindStateFunc(widget.Open, …)` sobre la raíz, leyendo
`t.sel.On()`.

En `buildRow`:

- Parte nueva `PartCheck = widget.Part("check")`, **primer** hijo en flujo, en
  el hueco que dejó el `⋮`.
- Se pinta **siempre** en el marcado; lo muestra/oculta el CSS. No lo
  construyas condicionalmente, o el reconciliado por clave recrea la fila
  entera al cambiar de modo.
- `BindStateFunc(widget.Selected, func() bool { return t.sel.IsChecked(id) })`.
  **No** hace falta un `widget.State` nuevo.
- Un solo handler de clic, con una rama:

```go
row.On("click", func(Event) {
    if t.sel.On().Get() {
        t.sel.Toggle(id)
        return
    }
    if t.OnSelect != nil {
        t.OnSelect(it)
    }
})
```

**No** dos handlers, **no** `StopPropagation` desde el check: el check no es un
objetivo aparte, la fila entera lo es.

En `css.go` de cada componente: `listselect.Apply(s, PartCheck)`.

## 6. Etapa 4 — Glifo

**`components/targetlist/svg.go`** y **`components/targetdate/svg.go`**:
registra una marca de verificación con el prefijo que ya usa cada paquete
(`tl-` / `td-`), siguiendo el patrón de los iconos que quites. Se pinta dentro
de `PartCheck` y hereda su `currentColor`: **no** le pongas `As()` propio.

## 7. Etapa 5 — Tests

### `components/listselect/listselect_test.go`

| Test | Comprueba |
|---|---|
| `TestZeroValueIsNormalMode` | `Mode{}` recién creado no está en modo selección |
| `TestToggleMarksAndUnmarks` | Dos toggles → desmarcado |
| `TestLeavingModeClearsTheMarks` | Marca dos, `SetOn(false)`, `SetOn(true)` → vacío |
| `TestCheckedIDsFollowTheGivenOrder` | Marca la 3ª y **luego** la 1ª, pide con el orden de render → `[1ª, 3ª]`, no `[3ª, 1ª]` (orden de pulsación) |
| `TestOnChangeReportsTheCount` | Marca, marca, desmarca → 1, 2, 1 |

`TestCheckedIDsFollowTheGivenOrder` es el que evita el fallo silencioso más
probable: devolver el slice interno tal cual, que está en orden de pulsación y
no coincide con lo que el usuario ve en pantalla. Marca en orden inverso a
propósito, o el test pasa por casualidad.

### Tests de consumidor, en cada componente

`components/targetlist/selectmode_test.go` y el equivalente en `targetdate`.

> **Regla del harness que aplica aquí:** *"An API is not published until a
> consumer-shaped test, inside the library itself, proves it."* Estos tests
> usan el componente **real** y la pieza **real** — nada de dobles de
> `listselect`. Si escribir el test resulta incómodo, la API es incómoda y has
> encontrado el defecto antes de publicarlo.

| Test | Comprueba |
|---|---|
| `TestRowsCarryNoOptionsMenu` | El HTML de una fila no contiene `…__button`, `…__options` ni `…__item-danger` |
| `TestSelectModeOffFiresOnSelect` | Modo apagado: un clic llama a `OnSelect` una vez y **no** marca |
| `TestSelectModeOnTogglesInsteadOfSelecting` | Modo encendido: un clic marca y **no** llama a `OnSelect` |
| `TestCheckIsInTheMarkupWhenModeIsOff` | El check está en el HTML aunque el modo esté apagado (lo esconde el CSS) |
| `TestCheckedIDsFollowRenderOrder` | A través del componente real, con `SetItems` |
| `TestSheetValidates` | `sheet().Validate()` sin errores — ya existe; que siga pasando |

**Anti-footgun:** sólo la librería estándar de testing. Nada de `testify` ni
`gomega`.

## 8. Lo que rompe

`OnDelete` y `CloseMenus()` desaparecen de ambos componentes. Quien los usa:

```
layout/crudview/crudview.go                        → lo arregla el plan de layout
app-demo/modules/medicalhistory/medicalhistory.go  → fase C de la ola
```

## 9. Criterios de aceptación

- [ ] `gotest` en verde (vet, race, cover, wasm).
- [ ] `grep -rn "openMenu\|CloseMenus\|PartItemDanger\|PartOptions\|iconDots\|iconDelete\|OnDelete" components/targetlist components/targetdate`
      → **sin resultados**.
- [ ] **La lógica de selección existe una sola vez.** La copia se delataría
      como un campo de estado propio en el struct del componente:
      `grep -rn "checked \[\]string\|checked map\[" components/targetlist components/targetdate`
      → **sin resultados**. Los tres métodos públicos deben delegar en
      `t.sel`, no llevar estado paralelo.
- [ ] `grep -rn "ClearChecked" components/` → **sin resultados**.
- [ ] Los dos componentes exponen la **misma** superficie:
      `SetSelectMode`, `CheckedIDs`, `OnCheckedChange`. Si divergen,
      `crudview` no puede intercambiarlos, que es justo lo que hace hoy.
- [ ] `components/listselect` es plano: `find components/listselect -mindepth 1 -type d`
      → vacío.
- [ ] `components/listselect/css.go` lleva `//go:build !wasm`;
      `listselect.go` **no** lleva build tag y no importa stdlib.
- [ ] **Ni un solo mapa en código que compile a WASM:**
      `grep -rn "map\[" components/listselect/listselect.go components/targetlist/targetlist.go components/targetdate/targetdate.go`
      → **sin resultados**.

## 10. Etapas

| # | Etapa | Ficheros | Depende de |
|---|---|---|---|
| 1 | Pieza `listselect` | `components/listselect/listselect.go`, `css.go` | — |
| 2 | Borrar el menú por fila | `targetlist.go`, `targetdate.go`, sus `css.go` y `svg.go` | — |
| 3 | Ensamblar la pieza | `targetlist.go`, `targetdate.go`, sus `css.go` | 1, 2 |
| 4 | Glifo | `targetlist/svg.go`, `targetdate/svg.go` | 2 |
| 5 | Tests | `listselect/listselect_test.go`, `*/selectmode_test.go` | 3, 4 |
