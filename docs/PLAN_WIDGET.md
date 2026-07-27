---
PLAN: "`tinywasm/widget` (nueva): el contrato visual — anatomía, estados y disposición"
---

> Depende de: [`docs/PLAN.md`](PLAN.md) (§3 el contrato, §5 justificación del repo).
> Bloquea a: [`PLAN_CSS`](PLAN_CSS.md), [`PLAN_SSR`](PLAN_SSR.md),
> [`PLAN_COMPONENTS`](PLAN_COMPONENTS.md) y las etapas de `layout`.
> Dependencias de la librería: **solo `github.com/tinywasm/fmt`** en el paquete raíz.

---

## Por qué existe

Hoy nadie posee el contrato visual. La consecuencia medible está en
[`PLAN.md` §1](PLAN.md): tres catálogos de tokens paralelos, una paleta local con nombres de
variables que no existen, 33 `RawRule` y estados disfrazados de clases CSS.

`widget` nombra **una sola cosa**: de qué piezas está hecho un componente visual, en qué
estados puede estar, y cómo se disponen esas piezas. Ni datos, ni transporte, ni DOM, ni
texto CSS.

---

## Estructura de paquetes

Se parte por **paquete**, con el mismo criterio que `tinywasm/svg` / `tinywasm/svg/sprite`:

```
github.com/tinywasm/widget          // sin build tag — VIAJA al binario WASM
    widget.go       Name, Part, Class
    state.go        State (del widget), Cue (del navegador)
    kind.go         Kind — el enum ARIA-APG
    capability.go   Widget, Selectable, Dismissible, Expandable

github.com/tinywasm/widget/style    // //go:build !wasm — NUNCA viaja al WASM
    sheet.go        Sheet, Of, Root, Part, When, Cue, Styler
    flow.go         Stack, Row, Split, Grid, Center, Cover, Reel, Frame
    scale.go        Space, Radius, TextSize, Weight, Elevation, Size, Ratio, Track
    surface.go      Surface + resolución a la tripleta de tokens
    except.go       Fill, Scrolls, Fixed, Flush, Clip
    emit.go         emisión a css.Stylesheet (interno)
```

Regla de tamaño, verificable: el paquete raíz **no debe superar ~150 líneas**. Si crece, es
que se le está colando emisión o lógica de aspecto.

---

## Etapa 1 — El paquete raíz (compartido)

### 1.1 Identidad y anatomía

```go
package widget

// Name identifica un widget. Es el prefijo de TODA clase que emite, así que dos
// widgets no pueden colisionar aunque elijan el mismo nombre de parte.
type Name string

// Part es una ranura nombrada de la anatomía de un widget (vocabulario Open UI).
// Es local al widget: "row", "menu", "header". Nunca lleva prefijo.
type Part string

// Class es un nombre de clase CSS. NO tiene constructor público: la única forma
// de obtener una es derivarla de un Name. Escribir Class("lo-que-sea") no compila
// fuera de este paquete.
type Class string

func (c Class) String() string      { return string(c) }
func (c Class) AsAttr() fmt.KeyValue { return fmt.KeyValue{Key: "class", Value: string(c)} }

// Root es la clase exterior del widget.
func (n Name) Root() Class { return Class(n) }

// Class deriva la clase de una parte: "targetlist__row".
func (n Name) Class(p Part) Class { return Class(string(n) + "__" + string(p)) }
```

> **La decisión que cierra el agujero:** `Class` es un tipo **no exportable por
> construcción** — su representación es `string` pero no hay conversión pública. Hoy
> `crudview` escribe `clsBtnCrud Class = "cv-btn-crud"` a mano y `css.go` repite el mismo
> literal por convención. Con esto, markup y hoja de estilos **derivan la clase del mismo
> símbolo**: desincronizarlos deja de ser posible.

### 1.2 Estados — dos tipos, no uno

```go
// State es un estado que POSEE el widget: lo escribe Go, lo lee la hoja de estilos.
type State uint8

const (
	Selected State = iota
	Disabled
	Locked   // solo lectura, pero plenamente legible
	Invalid
	Busy
	Open     // desplegado / expandido
	Current  // ítem de navegación activo
)

// Attr devuelve el atributo que el DOM escribe y sobre el que la hoja selecciona.
// Markup y CSS coinciden por construcción, no por convención.
func (s State) Attr() fmt.KeyValue // {Key: "data-selected", Value: "true"}

// Cue es un estado que posee el NAVEGADOR. Solo se estiliza; no se puede escribir
// desde Go — por eso es un tipo distinto y no tiene Attr().
type Cue uint8

const (
	Hover Cue = iota
	Focus
	Press
	Target
)
```

Por qué dos tipos: `Hover` no es escribible desde Go, `Selected` sí. Un solo `State` con
ambos permitiría escribir `el.Set(Hover.Attr())` — un estado ilegal, representable. Principio
3.

Esto además elimina la clase-estado que hoy existe: `clsBtnCrudIconHidden` en
`crudview/crudview.go:26` no es una clase, es `State.Open` del botón toggle.

### 1.3 `Kind` — el enum ARIA-APG

```go
// Kind es el tipo de widget según WAI-ARIA Authoring Practices. Determina el rol,
// los estados válidos y el teclado esperado. Cerrado a propósito: si un widget no
// encaja en ninguno, casi siempre es que son dos widgets.
type Kind uint8

const (
	Region Kind = iota // contenedor sin semántica de interacción
	Listbox            // targetlist
	Menu               // el menú ⋮
	Dialog             // modaldialog
	Disclosure         // <details> desplegable
	Tabs
	Toolbar
	Grid
	Combobox
	Form
	Alert              // toasts de platformd
)
```

Beneficio directo: el rol ARIA y los atributos obligatorios (`aria-selected`,
`aria-expanded`, `aria-modal`) los emite el renderer a partir de `Kind` — la accesibilidad
deja de ser un retrofit y pasa a salir de la firma.

### 1.4 Contratos de capacidad

```go
// Widget es la identidad. Es lo único obligatorio.
type Widget interface {
	WidgetName() Name
	WidgetKind() Kind
}

// Capacidades — cada costura asevera solo la que necesita (patrón de la casa,
// el mismo de view.Saver / view.Deleter).
type Selectable interface{ Select(id string) }
type Dismissible interface{ Dismiss() }
type Expandable interface{ Expand(open bool) }
```

`Render()` **no** está aquí a propósito: obligaría a `widget` a importar `dom`, y entonces
`form`/`view` —que solo necesitan nombrar `Invalid` y `Selected`— arrastrarían el DOM. El
markup lo cubre `dom.Component`, que ya existe.

---

## Etapa 2 — `widget/style`: las escalas cerradas

Todas son enums. **Ninguna acepta un `string`, un `int` ni una unidad.**

```go
package style

type Space uint8   // Space0 … Space12 → tokens css.Space*
type Radius uint8  // RadiusNone, RadiusSm, RadiusMd, RadiusLg, RadiusFull
type TextSize uint8// TextXs … Text2xl
type Weight uint8  // WeightRegular, WeightMedium, WeightBold
type Elevation uint8 // Flat, Raised, Floating, Overlay
type Ratio uint8   // Half, TwoThirds, ThreeQuarters  (proporción de Split)
type Track uint8   // TrackSm, TrackMd, TrackLg       (pista mínima de Grid)

// Size es la ÚNICA medida de ancho. Relativa al contenedor, nunca al viewport —
// las unidades vw/vh no existen en esta API (ver ROADMAP: el bug de la franja gris).
type Size uint8
const (
	Content Size = iota // se ajusta a su contenido
	Prose               // ancho legible (~65ch)
	Third
	Half
	TwoThirds
	Full                // 100% del contenedor
)
```

**No existe `Height`.** Si solo se declara el ancho, el alto es automático — que es
exactamente el default correcto en CSS y el que el código actual pelea a mano con
`min-height:0`, `height:100%` y `grid-template` en tres bugs distintos del ROADMAP.

---

## Etapa 3 — `Surface`: el color desaparece del vocabulario

```go
// Surface es una decisión visual completa: fondo, texto y borde resueltos como
// UNA tripleta. No hay forma de elegir un fondo sin su texto — que es la causa
// del bug de rightpanel/css.go:17, donde ColorOnSurface (texto) se usó de fondo.
type Surface uint8

const (
	Page     Surface = iota // el lienzo de la aplicación
	Panel                   // tarjeta que flota sobre la página
	Sunken                  // pozo hundido dentro de un panel  (el "cInset" de crudview)
	Accent                  // relleno de marca
	Secondary
	Selected                // fila seleccionada
	Success
	Danger
	Muted                   // texto secundario
	Disabled
)

func On(s Surface) Opt
```

Resolución (la tabla es la única fuente de verdad; vive en `style`, no en cada componente):

| Surface | fondo | texto | borde |
|---|---|---|---|
| `Page` | `ColorBackground` | `ColorOnSurface` | — |
| `Panel` | `ColorSurface` | `ColorOnSurface` | `ColorOutline` † |
| `Sunken` | `ColorSurfaceSunken` † | `ColorOnSurface` | — |
| `Accent` | `ColorPrimary` | `ColorOnPrimary` ‡ | — |
| `Secondary` | `ColorSecondary` | `ColorOnSecondary` | — |
| `Selected` | `ColorSelection` | `ColorOnSelection` | — |
| `Success` | `ColorSuccess` | `ColorOnSuccess` † | — |
| `Danger` | `ColorError` | `ColorOnError` † | — |
| `Muted` | *transparente* | `ColorMuted` | — |
| `Disabled` | `ColorDisabled` † | `ColorOnDisabled` † | — |

† Tokens que **hoy no existen** y que `crudview` inventó con nombres de Material. Se añaden en
[PLAN_CSS](PLAN_CSS.md) Etapa 1.
‡ El par `ColorPrimary`/`ColorOnPrimary` hoy **no contrasta** — por eso `crudview` lo descarta
y hardcodea `#ffffff`. Se corrige con un test de contraste en `css` ([PLAN_CSS](PLAN_CSS.md)
Etapa 2).

Cada `Surface` define **una vez** su variante `Hover`, `Focus` y `Press`. Eso resuelve de raíz
lo que el ROADMAP describe como *"Hover color was inconsistent (`ColorPrimary`, hardcoded
`filter: brightness`, ad-hoc backgrounds) across components"*.

---

## Etapa 4 — `Flow`: las primitivas de disposición

Vocabulario tomado de *Every Layout*. **Todas son intrínsecamente responsivas: ninguna emite
un `@media`.**

```go
func Stack(gap Space) Opt              // ritmo vertical; hijos a ancho completo
func Row(gap Space) Opt                // horizontal; envuelve cuando no cabe
func Split(r Ratio, gap Space) Opt     // dos paneles; se apila bajo su propio ancho
func Grid(min Track, gap Space) Opt    // auto-fit + minmax; NO se elige nº de columnas
func Center() Opt                      // columna centrada, tope Prose/Content
func Cover() Opt                       // llena el contenedor; un hijo centrado
func Reel(gap Space) Opt               // tira horizontal con scroll-snap
func Frame(r Ratio) Opt                // caja de proporción fija (media)
```

Emisión, con los mecanismos estándar de [`PLAN.md` §2.5](PLAN.md):

```css
@layer primitives {
  .fl-stack  { display:flex; flex-direction:column; min-height:0 }
  .fl-stack > * + * { margin-block-start: var(--gap) }
  .fl-row    { display:flex; flex-wrap:wrap; gap: var(--gap); align-items:center }
  .fl-grid   { display:grid; gap:var(--gap);
               grid-template-columns: repeat(auto-fit, minmax(min(var(--track),100%),1fr)) }
  .fl-split  { container-type:inline-size; display:grid; gap:var(--gap);
               grid-template-columns: var(--ratio) 1fr }
  @container (max-width: 40rem) { .fl-split { grid-template-columns: 1fr } }
}
```

Tres cosas que valen el cambio entero:

1. **Se emite una vez.** Hoy `Display(Flex_) + FlexDirection(Column) + MinHeight(Str("0"))`
   está copiado literal en 6 reglas solo de `crudview/css.go`.
2. **`@container`, no `@media`.** `Split` reacciona a *su* ancho. El bug del ROADMAP
   (*"sizing a scroll-snap child in `vw` against a narrower container overflows it"*) deja de
   ser escribible.
3. **`Split` sustituye 45 líneas de `crudview`.** El bloque `Media("(max-width: 640px)")` con
   `direction:rtl`, `order`, `scroll-snap-align` y 20 líneas de comentario justificándolo se
   borra entero: el orden físico form-izquierda / lista-derecha es una propiedad de la
   primitiva, no un truco de flujo bidireccional.

---

## Etapa 5 — Las excepciones (lo único que se declara)

> *"todos los sitios se buscan que sean responsivos… ¿entonces para qué declararlo? debería
> ser al contrario: declarar cuando no se requiera."*

Exactamente cinco. Son pocas, tipadas y **greppables** — se puede auditar de un vistazo qué
partes de la app se salen del comportamiento estándar:

```go
func Fill() Opt     // además del ancho, toma el alto disponible.
                    // Emite el conjunto correcto COMPLETO (height:100%, min-height:0,
                    // el track 1fr del grid padre) — los tres bugs de alto del ROADMAP
                    // en un solo símbolo.
func Scrolls() Opt  // desborda internamente en vez de crecer. Implica Fill().
func Fixed() Opt    // NO reflota: mantiene su disposición bajo cualquier ancho.
func Flush() Opt    // sin radio ni margen: pega a ras del contenedor padre.
                    // (el caso "Square, not rounded" de crudview/css.go:41)
func Clip() Opt     // recorta a los hijos (overflow:hidden).
```

Nada más. Si aparece la necesidad de una sexta, es señal de que falta una primitiva `Flow`,
no una excepción.

---

## Etapa 6 — `Sheet`: la única forma de emitir estilo

```go
type Sheet struct{ /* privado */ }

// Of abre el bloque de estilo de un widget. Todas sus reglas quedan scoped a n.
func Of(n widget.Name) *Sheet

func (s *Sheet) Root(opts ...Opt) *Sheet
func (s *Sheet) Part(p widget.Part, opts ...Opt) *Sheet
func (s *Sheet) When(st widget.State, p widget.Part, opts ...Opt) *Sheet
func (s *Sheet) Cue(c widget.Cue, p widget.Part, opts ...Opt) *Sheet

// Styler es la capacidad "este widget tiene aspecto". La asevera el recolector SSR
// (ver PLAN_SSR): una interfaz, no una expresión regular sobre el nombre.
type Styler interface {
	widget.Widget
	Style() *Sheet
}
```

No existe `Sheet.Media(...)`, `Sheet.Raw(...)` ni `Sheet.Selector(...)`. **No hay escape.**
Un selector arbitrario es la puerta por la que vuelve todo lo demás.

### Ejemplo completo — `targetlist` entero

```go
const (
	nameTargetList = widget.Name("targetlist")
	partRow        = widget.Part("row")
	partMenu       = widget.Part("menu")
)

func (l *TargetList) WidgetName() widget.Name { return nameTargetList }
func (l *TargetList) WidgetKind() widget.Kind { return widget.Listbox }

func (l *TargetList) Style() *style.Sheet {
	return style.Of(nameTargetList).
		Root(Stack(Space1), On(Sunken), Scrolls(), Round(RadiusMd)).
		Part(partRow, Row(Space2), On(Panel), Pad(Space2), Round(RadiusSm)).
		Part(partMenu, Stack(Space0), On(Panel), Raise(Floating), Clip()).
		When(widget.Selected, partRow, On(Selected)).
		When(widget.Disabled, partRow, On(Disabled)).
		Cue(widget.Hover, partRow, On(PanelHover))
}
```

Nueve líneas. No hay un color, ni una longitud, ni un breakpoint, ni un selector.

---

## Etapa 7 — Emisión estándar y determinista

`Sheet` produce CSS con orden de capas fijo, decidido por la librería:

```css
@layer tokens, primitives, widgets, states;
```

Garantías, todas verificables en test:

- **Cero `!important`.** El orden lo decide la capa, no la especificidad.
- **Especificidad plana.** Toda regla de widget es `.clase` o `.clase[data-x]`. Nunca
  descendientes anidados como `."+string(clsTitle)+" h1"` — el patrón de hoy, que hace que el
  estilo dependa de la forma del markup.
- **Salida estable.** El mismo `Sheet` produce byte a byte el mismo texto: los diffs de CSS
  se vuelven revisables.
- **Sin `@media` emitido por un widget.** Solo `@container`, y solo desde una primitiva.

---

## Etapa 8 — La prueba que autoriza a publicar

> *"An API is not published until a consumer-shaped test, inside the library itself, proves
> it."*

`widget/style/consumer_test.go` construye, **dentro de esta librería**, un widget maestro-detalle
completo con la pila real (`css` real, `html` real, `dom` fake) y asevera:

1. Toda clase presente en el markup existe en la hoja, **y viceversa** — el par que hoy nadie
   verifica y que dejó pasar `cv-btn-crud-icon-hidden` como estado disfrazado de clase.
2. La hoja no contiene: `!important`, `@media`, ningún literal de color, ninguna unidad
   `vw`/`vh`.
3. Las capas aparecen en el orden declarado.
4. Cada `Surface` usada resuelve a un token que **existe** en el catálogo de `css` — el test
   que habría atrapado `--color-surface-variant` el día que se escribió.
5. Emisión determinista: dos ejecuciones, salida idéntica.
6. `GOOS=js GOARCH=wasm go list -deps` sobre un consumidor de ejemplo **no** contiene
   `widget/style`.

Si escribir ese test resulta incómodo, la API es incómoda y el defecto se encontró antes de
publicar.

---

## Riesgos

| Riesgo | Mitigación |
|---|---|
| El vocabulario resulta insuficiente y alguien pide reabrir un escape | `rightpanel` (166 líneas) es el canario deliberado — ver [PLAN.md §8 Etapa 2](PLAN.md). Si falta vocabulario **se amplía el enum aguas arriba y se publica**; jamás se añade `Raw`. |
| `Surface` no cubre un caso legítimo (p. ej. un gráfico con series de color) | Es un caso real y distinto: se resuelve con un catálogo `DataColor` aparte, cerrado y accesible por índice — nunca reabriendo el color libre. |
| Soporte de `@container` | Disponible en todos los navegadores objetivo desde 2023. `Split` degrada a columna única sin container queries: peor, pero nunca roto. |
| `widget` crece hasta volverse un framework | El límite de ~150 líneas del paquete raíz es la alarma. Datos, transporte o DOM dentro de `widget` = la pieza dejó de ser lego. |
