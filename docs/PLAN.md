---
PLAN: "API simplificada css/widget — un contrato único para construir componentes visuales reutilizables"
---

> Plan maestro. Define el problema, la respuesta de diseño y el reparto de trabajo
> entre librerías. Cada librería involucrada tiene su propio plan:
>
> - [`docs/PLAN_WIDGET.md`](PLAN_WIDGET.md) — `tinywasm/widget` (**nueva**): el contrato visual.
> - [`docs/PLAN_CSS.md`](PLAN_CSS.md) — `tinywasm/css`: catálogo de tokens y motor de emisión.
> - [`docs/PLAN_SSR.md`](PLAN_SSR.md) — `tinywasm/ssr`: detección tipada de proveedores de estilo.
> - [`docs/PLAN_COMPONENTS.md`](PLAN_COMPONENTS.md) — `tinywasm/components`: migración de widgets.
> - Este archivo cubre además las etapas dentro de `tinywasm/layout`.
>
> Reglas del repo: `AGENTS.md` en la raíz. Principio rector:
> [CONSTRUCTION_HARNESS](https://github.com/tinywasm/.github) — *el código tipado y explícito
> **es** el arnés*.
>
> **Publicado:** `github.com/tinywasm/widget v0.1.0` y `github.com/tinywasm/css v0.2.0` ya están
> etiquetados. `tinywasm/ssr`, `tinywasm/components` (+ `tinywasm/form`) y las etapas de
> `tinywasm/layout` (§8) se ejecutan **de una sola vez, no por etapas escalonadas**: no hay
> canario ni periodo de coexistencia entre versiones — se migra todo el árbol en el mismo
> cambio y se fija `go.mod` directamente a `css v0.2.0` / `widget v0.1.0`. El único gate real es
> el de siempre: `gotest` en verde y el chequeo de dependencias WASM.

---

## Resumen ejecutivo

El problema no es que los agentes "no lean la documentación". El problema es que **la API
de CSS actual permite escribir lo incorrecto**, y en varios casos *obliga* a hacerlo. Un
arnés que se puede evadir no es un arnés.

La propuesta tiene tres movimientos:

1. **Eliminar el color, la longitud y el breakpoint del vocabulario.** Si no existe un tipo
   para escribir `#3f88bf`, `.4em` o `(max-width: 640px)`, nadie los escribe — ni un humano
   ni un LLM. Se declara *intención* (`On(Panel)`, `Pad(Space2)`), no *valor*.
2. **Invertir el default: responsivo, fluido y alto-automático por defecto.** No se declara
   "quiero que sea responsivo"; se declara la **excepción** (`Fixed()`, `Fill()`, `Scrolls()`).
   Las excepciones son pocas, tipadas y *greppables*.
3. **Un contrato único de widget = identidad + capacidades aseveradas en la costura**, no una
   interfaz gorda. Es exactamente el patrón que la casa ya usa (`view.Saver`/`view.Deleter`),
   extendido a lo visual.

Superficie actual usada solo en los 3 `css.go` de este repo: **81 símbolos exportados
distintos** + **33 `RawRule`** (el agujero sin tipar). Superficie propuesta: **~22
constructores**, **0 escapes**.

---

## 1. Diagnóstico — con evidencia de este repositorio

No es una impresión. Es lo que hay hoy en `crudview/css.go`, `platformd/css.go` y
`rightpanel/css.go`.

### 1.1 El catálogo tipado existe, pero la API no obliga a usarlo

`tinywasm/css` publica un catálogo `Token` tipado (`ColorPrimary`, `ColorSurface`, `Space2`…).
Sin embargo `crudview/css.go:13-23` declara su propia paleta como **strings crudos**:

```go
const (
	cPanel  = "var(--color-background, #ffffff)"
	cInset  = "var(--color-surface-variant, #d7d7dd)"
	cBorder = "var(--color-outline-variant, #cfcfd6)"
	cAccent = "var(--color-primary, #3f88bf)"
	cOnAcc  = "#ffffff"
	cDisBg  = "var(--color-outline-variant, #c2c1c1)"
	cDisFg  = "var(--color-on-surface-variant, #6e6e73)"
)
```

Esto es posible **solo porque `Background()` y `Color()` aceptan `Str(string)`**. Mientras
exista `Str`, el token tipado es opcional — y lo opcional, en un flujo asistido por LLM, no
ocurre.

### 1.2 Tres de esas variables **no existen** en ningún lado

`--color-surface-variant`, `--color-outline-variant` y `--color-on-surface-variant` son
nombres de Material Design, no del catálogo de `tinywasm/css`. Verificado: no se declaran en
`css/tokens.go` ni en ningún `Root(Declare(...))` del ecosistema.

Consecuencia real, no teórica: **el marco gris, las hairlines y el estado deshabilitado de
`crudview` nunca siguen el tema.** Siempre resuelven al hex de fallback. En modo oscuro
siguen siendo grises claros. El bug es invisible porque `var()` con fallback **nunca falla
en voz alta** — es un fallo silencioso, que el principio 6 prohíbe explícitamente.

### 1.3 Pseudo-tokens que nadie puede tematizar

`--cv-title-height`, `--cv-controls-height`, `--cv-detail-width` se **usan** dentro de
`RawRule("... var(--cv-title-height, 8vh) ...")` pero **nunca se `Declare()`n**. Son
variables fantasma: no tienen símbolo Go, no aparecen en autocompletado, no se pueden
sobrescribir desde una app, y su único valor real es el fallback embebido en un string.

### 1.4 Cada paquete reinventa su propia escala

`platformd/tokens.go` declara 7 tokens `--pd-*`; `rightpanel/css.go` declara 10 tokens
`--rp-*`; `crudview` declara 3 fantasma `--cv-*`. Todos describen lo mismo: alto de título,
alto de controles, ancho de panel, gap, color de borde. **Tres catálogos paralelos para un
concepto.** Es exactamente la duplicación que el arnés llama *"pegamento que toda aplicación
escribiría igual"*: pertenece a la pieza, no a los consumidores.

Peor: `rightpanel/css.go:17` hace `Declare(tokenAsideBg, ColorOnSurface.Var())` — usa el
color **de texto** como color **de fondo**. Nada lo impide, porque `Background()` acepta
cualquier `Token`. Un par fg/bg emparejado por tipo lo haría imposible de escribir.

### 1.5 El escape hatch no es una salida de emergencia; es la vía principal

**33 `RawRule`** en tres archivos. Y no para casos exóticos: para `grid-template`, `gap`,
`direction`, `scroll-snap-type`, `flex`, `order` — es decir, **para el layout**, que es
precisamente lo que la librería debería resolver. La API tipada cubre lo fácil (`Display`,
`Color`) y delega lo difícil al string.

Además el escape hatch tiene sus propias trampas documentadas *en el código*:

```go
// NOTE: adjacent RawRules are concatenated without a separator, so
// grid-template and gap MUST share one RawRule with an explicit ';'.
```

Eso es un "acuérdate de…" en un comentario: por definición del arnés, **un agujero**.

### 1.6 Bugs de la propia librería, parcheados río abajo

`crudview/css.go:174-181` documenta que `Padding(a,b,c,d)` **está roto**:

> *"nor a single multi-token `Padding(a,b,c,d)` call (`joinValues()`'s output loses its
> spaces somewhere in the CSS pipeline: `padding:var(...)var(...)var(...)`, no separators —
> both verified via the live stylesheet)"*

Y `crudview` lo evita usando 4 longhands. Las reglas lego son explícitas: *"Never wrap a
library to fix its behaviour"* y *"A missing contract at a boundary is a defect in the
library, not in the consumer"*. El defecto es de `css` y se arregla en `css`
([PLAN_CSS](PLAN_CSS.md), Etapa 3).

Lo mismo con `ColorOnPrimary`: `crudview` lo descarta a mano —

> *"Not --color-on-primary: some themes set on-primary to a near black that is unreadable on
> the primary fill; white is the safe universal."*

— y hardcodea `#ffffff`. Si un par `primary`/`on-primary` no contrasta, **el par está mal
definido en el catálogo**. Se arregla arriba, con un test de contraste en `css`.

### 1.7 El breakpoint es un string, teniendo el token al lado

`Media("(max-width: 640px)")` en `crudview/css.go:343` y `rightpanel/css.go:143`, mientras
`css.BpSm` vale exactamente `640px`. Dos fuentes de verdad para el mismo umbral, que
divergirán.

### 1.8 Fallo silencioso estructural en la detección SSR

De `AGENTS.md`: `tinywasm/ssr` detecta el CSS de un paquete **por regex sobre el nombre de la
función**. Si se llama `GenerateCSS` en vez de `RenderCSS`, *"is silently never emitted — the
component renders with zero styling and nothing fails at build time"*. Y además exige "un solo
receptor por paquete", regla que no está en ningún tipo.

Un contrato es un tipo, no una expresión regular. Se resuelve en [PLAN_SSR](PLAN_SSR.md).

### 1.9 Resumen del diagnóstico

| Síntoma | Causa raíz en la API | Principio violado |
|---|---|---|
| Colores hardcodeados | `Str(string)` acepta cualquier cosa | 1 — tipado sobre `any` |
| Variables inexistentes | `var()` con fallback nunca falla | 6 — fallar en compilación |
| Tokens fantasma `--cv-*` | Ningún tipo obliga a `Declare` | 3 — estados ilegales inescribibles |
| 3 catálogos paralelos | La escala no la posee nadie | 9 — piezas lego |
| 33 `RawRule` | El vocabulario tipado es incompleto | 1, 4 |
| `Padding` roto, parcheado abajo | Defecto no reportado aguas arriba | 9 |
| Breakpoint como string | Token existente no obligatorio | 4 — una sola forma |
| CSS no emitido por nombre de función | Contrato por regex, no por tipo | 6 |

---

## 2. ¿Existe un framework agnóstico a la tecnología para construir UI?

Respuesta corta: **no existe uno adoptable tal cual, pero sí existen cuatro vocabularios
estandarizados que resuelven exactamente lo que estás preguntando.** Conviene robarlos por
nombre: cuestan cero, y un agente que ya conoce esos estándares conoce tu API sin leer nada.

Lo que **no** sirve como base:

- **Flutter / Compose Multiplatform / React Native** — no son agnósticos: son dueños del
  renderer. Cambian el problema por otro más grande.
- **Adaptive Cards (Microsoft) y similares** — sí son agnósticos, pero degeneran en un
  esquema JSON sobre un catálogo cerrado de tarjetas. Sirve para notificaciones, no para
  construir un sistema de componentes.

Lo que **sí** sirve, y qué toma de cada uno esta propuesta:

### 2.1 W3C Design Tokens Community Group (DTCG) — *los valores*

Formato estándar para nombrar y serializar decisiones de diseño (`color.surface`,
`space.4`, `radius.md`, con tipo y valor). `css.Token` **ya es esto** en la práctica.

Lo que aporta alinearse: (a) el catálogo se vuelve importable/exportable desde Figma o Style
Dictionary; (b) da un **nombre canónico a los tokens que hoy faltan** — `surface-variant`,
`outline`, `on-error` no son inventos de Material, son huecos reales del catálogo, y DTCG
dice cómo llamarlos.

### 2.2 Open UI (W3C) — *la anatomía*

Describe cada componente como **`name` + `parts` nombradas + `states`**. Es literalmente el
contrato que estás buscando: no dice cómo se ve, dice **de qué piezas está hecho y en qué
estados puede estar**.

De aquí salen `widget.Name`, `widget.Part` y `widget.State`. Ventaja concreta: hoy
`crudview` declara 20 `Class` a mano (`clsBtnCrudIconHidden` es un *estado* disfrazado de
clase). Con anatomía tipada, la clase se **deriva** y el estado es un `data-*` — imposible
de colisionar entre paquetes, imposible de desincronizar entre markup y hoja de estilos.

### 2.3 WAI-ARIA Authoring Practices (APG) — *el comportamiento*

Es la cosa más parecida a "un contrato agnóstico y normativo para widgets reutilizables" que
existe de verdad. Define, por tipo de widget (`listbox`, `combobox`, `dialog`, `disclosure`,
`grid`, `tabs`, `toolbar`, `menu`…), sus roles, estados y teclado esperados.

Úsalo como **la lista cerrada de qué puede ser un widget**. Beneficio secundario grande: la
accesibilidad deja de ser un retrofit — sale de la firma. `targetlist` es un `listbox`; el
menú ⋮ es un `menu`; `modaldialog` es un `dialog`. Nombrarlo así obliga al markup correcto.

### 2.4 Every Layout + Intrinsic Web Design — *la disposición*

*Every Layout* (Bell & Pickering) y *Intrinsic Web Design* (Jen Simmons) son **la respuesta
publicada a tu pregunta "¿para qué declarar responsivo?"**. Proponen un conjunto pequeño de
**primitivas de layout intrínsecamente responsivas, sin un solo media query**: Stack,
Cluster, Sidebar, Switcher, Cover, Grid (`auto-fit` + `minmax`), Reel, Frame, Center.

De aquí sale `style.Flow`. Y de aquí sale el default invertido: una primitiva **ya** se
adapta; lo que se declara es cuándo **no** debe hacerlo.

### 2.5 Dos mecanismos CSS que hacen todo esto viable

- **`@layer` (cascade layers)** — el estándar que elimina la ambigüedad de la cascada. Con
  `@layer tokens, primitives, widgets, states;` el orden lo decide la capa, no la
  especificidad. Se acabaron los `!important` y los "por qué gana esta regla".
- **`@container` (container queries)** — un widget responde al **contenedor**, no al
  viewport. Esto no es estética: es el bug que ya pagaste. `ROADMAP.md` documenta que los
  paneles en `100vw`/`90vw` desbordaban un contenedor de `96vw`, y la corrección fue pasar a
  `%`. Con container queries ese tipo de bug **no se puede escribir**.

### 2.6 Precedente para separar comportamiento de aspecto

Radix Primitives, React Aria, Headless UI y Melt UI: todos publican **contrato de
comportamiento sin aspecto**. Es la validación externa de la decisión clave de esta
propuesta — `view.Presenter` (datos/conducta) y `widget` (anatomía/aspecto) son **contratos
distintos**, unidos en la costura, no fusionados.

---

## 3. El contrato único

> *"quiero limitarla a un contrato único. ya tengo una firma de presentación, pero no incluye
> el cómo se verá. no sé si es adecuado que los datos siempre estén en este contrato."*

**No, los datos no deben estar en el contrato visual.** Y "contrato único" no significa una
interfaz gorda: significa **una identidad única + capacidades aseveradas en la costura**.

Es el patrón que la casa ya tiene escrito: *"Capability bag + type assertion at the seam is
the assembly pattern"*, y que `view` ya aplica con `Presenter` + `Saver` + `Deleter`.

```go
// IDENTIDAD — obligatoria, mínima, sin dependencias.
type Widget interface {
	WidgetName() Name        // deriva TODAS las clases de sus partes
}

// CAPACIDADES — cada costura asevera solo lo que necesita.
dom.Component      → Render() *dom.Element      // markup      (wasm + ssr)
style.Styler       → Style() *style.Sheet       // aspecto     (!wasm)
view.Presenter     → datos y conducta           // datos       (agnóstico de UI)
widget.Selectable  → Select(id string)          // interacción
widget.Dismissible → Dismiss()
```

```flowchart TD
A[Módulo de dominio] --> B[view.Presenter<br/>datos + conducta<br/>NO sabe cómo se ve]
A --> C[widget.Widget<br/>identidad: WidgetName]
C --> D[dom.Component<br/>Render → markup + data-state]
C --> E[style.Styler<br/>Style → Sheet · solo !wasm]
D --> F[Binario WASM<br/>solo Name/Part/State: strings]
E --> G[Hoja SSR<br/>@layer determinista]
B --> H[Costura: el renderer asevera<br/>Saver / Deleter / Selectable]
D --> H
```

Por qué los datos fuera:

- **Un widget que nombra sus datos solo es reutilizable para esos datos.** Es exactamente por
  qué `crudview` hoy no es reutilizable: está soldado a `view.Presenter` *y* a la forma
  `targetlist.Item`.
- La proyección ya está resuelta arriba: `view.Itemizer` (`Item() Item`) es el único código
  de vista que un modelo escribe. Ese es el punto de contacto correcto — un widget consume
  `[]view.Item`, no un `Presenter`.
- **Regla operativa:** si un widget necesita saber *de qué* es la lista, el diseño está mal.
  Solo debe saber *cuántas filas*, *qué muestra cada una* y *en qué estado está*.

---

## 4. La API simplificada

Detalle completo y firmas en [`docs/PLAN_WIDGET.md`](PLAN_WIDGET.md). Aquí, el criterio.

### 4.1 Lo que se elimina del vocabulario

No hay forma de escribirlo porque **no existe el tipo**:

| Se elimina | Por qué | Reemplazo |
|---|---|---|
| Color literal (`Str("#fff")`) | Fuente de todo hardcodeo | `On(Surface)` — tripleta fg/bg/borde emparejada |
| Longitud literal (`.4em`, `8vh`, `16px`) | Escalas divergentes | `Space`, `Radius`, `Text` — enums cerrados |
| `Height(...)` | *"si solo declaras el ancho, el alto se sobreentiende"* | no existe; `Fill()` es la excepción |
| `Media(...)` | Responsivo es el default | primitivas fluidas + `@container` |
| `Display/Position/Float/...` | Es *cómo*, no *qué* | lo decide la primitiva `Flow` |
| `RawRule(...)` | Agujero sin tipar | si falta algo, es defecto de la pieza — se reporta |

### 4.2 Lo que queda: intención, ~22 constructores

```go
// Una primitiva de disposición. Todas son fluidas y reflotan solas.
Stack(gap)              // ritmo vertical
Row(gap)                // horizontal; envuelve solo
Split(ratio, gap)       // detalle|lista; se apila bajo su propio ancho
Grid(minTrack, gap)     // auto-fit; no se elige número de columnas
Center() · Cover() · Reel(gap) · Frame(ratio)

// Superficie: fondo + texto + borde SIEMPRE como una tripleta.
On(Page|Panel|Sunken|Accent|Selected|Danger|Success|Muted|Disabled)

// Medida: UN eje. El alto siempre es automático.
Width(Content|Prose|Half|Third|TwoThirds|Full|Screen)

// Texto, espacio, decoración: solo desde la escala.
Text(size) · Weight(w) · Pad(space) · Inset(space) · Round(radius) · Raise(elevation)

// LAS EXCEPCIONES — lo único que se declara explícitamente.
Fill()      // además, toma el alto disponible
Scrolls()   // desborda internamente en vez de crecer
Fixed()     // NO reflota
Flush()     // sin radio: pega a ras del contenedor padre
Clip()      // recorta a los hijos
```

### 4.3 Cómo se lee — el mismo panel de `crudview`, antes y después

```go
// ── HOY (crudview/css.go, ~40 líneas para dos partes) ────────────────────
Rule(clsAsideWrap,
	Flex(None), Display(Flex_), FlexDirection(Column),
	MinWidth(Str("0")), MinHeight(Str("0")),
	Background(Str(cPanel)),                     // ← hex escondido tras var()
	BorderRadius(Str(".4em")),                   // ← magic number
	Padding(Space1),
	RawRule("gap: "+Space1.Var()),               // ← escape hatch para un gap
)
Media("(max-width: 640px)",                      // ← breakpoint como string
	Rule(clsAsideWrap,
		RawRule("direction:ltr; flex:0 0 100%; scroll-snap-align:start; order:1"),
	),
)

// ── PROPUESTA ────────────────────────────────────────────────────────────
s.Part(partAside, Stack(Space1), On(Panel), Round(RadiusMd), Pad(Space1), Fill())
```

El reflow móvil no aparece: `Split` ya se apila bajo su propio ancho vía `@container`. La
tira `Reel` con scroll-snap tampoco: es la primitiva que `Split` usa al colapsar. **No hay
`direction:rtl` que razonar** — el orden físico lo garantiza la primitiva, no un truco de
flujo bidireccional que costó un plan entero.

### 4.4 Por qué esto reduce el binario WASM (y por qué el CSS también encoge)

Son dos mecanismos distintos; conviene no confundirlos:

- **Binario WASM.** Todo `widget/style` lleva `//go:build !wasm`. El lado navegador solo
  carga `widget`: `Name`, `Part`, `State` — strings y un `uint8`. Se verifica con el mismo
  chequeo de grafo que `AGENTS.md` ya prescribe para `svg/sprite`:
  ```bash
  GOOS=js GOARCH=wasm go list -deps ./... | grep tinywasm/widget/style   # DEBE estar vacío
  ```
- **Hoja CSS.** Las primitivas se emiten **una vez** en `@layer primitives`, y cada parte de
  cada widget las referencia. Cuarenta widgets con `Stack` producen **un** bloque, no
  cuarenta. Hoy `Display(Flex_) + FlexDirection(Column) + MinHeight(0)` está copiado
  literalmente en 6 reglas solo en `crudview/css.go`.

---

## 5. ¿Hace falta un repo nuevo? — sí, uno: `tinywasm/widget`

> *"¿sería necesario crear otro repo con la firma? widget por ejemplo? justifica"*

**Sí, uno.** Y la justificación tiene que pasar el filtro de las reglas lego, no el gusto.

### 5.1 ¿Hay una responsabilidad que hoy no tiene dueño?

Sí, y hay prueba forense de ello. La regla dice:

> *"A missing contract at a boundary is a defect in the library, not in the consumer. If two
> libraries meet and there is no type to name the thing that crosses between them, the type is
> missing upstream. Do not declare a local intersection to paper over it."*

`crudview/css.go:13-23` **es exactamente esa intersección local**: una paleta declarada abajo,
con nombres de variables de otro sistema de diseño, porque arriba no existía el tipo que
nombrara "la superficie hundida de un panel". Y no es un caso aislado: `platformd`,
`rightpanel` y cada componente de `tinywasm/components` repiten el patrón con su propio
prefijo. El síntoma que reportaste — *"los agentes hardcodean colores"* — es el efecto, no
la causa. La causa es que **nadie posee el contrato visual**.

### 5.2 ¿Puede vivir dentro de una pieza existente?

Se evaluaron las tres candidatas:

| Candidata | Por qué no |
|---|---|
| `tinywasm/css` | `css` emite **texto CSS**: es `!wasm` por naturaleza. El contrato (`Part`, `State`) debe cruzar al WASM porque el markup lo escribe. Meterlos juntos reproduce exactamente la trampa que `AGENTS.md` documenta para `svg`: *"compiles for WASM too, so forgetting the `!wasm` tag does NOT fail the build — it silently ships every path string into the browser bundle"*. |
| `tinywasm/dom` | Arrastraría el vocabulario visual al binario, y obligaría a `form`/`view` a depender de `dom` para nombrar un estado. Dirección de dependencia invertida. |
| `tinywasm/components` | Es un **consumidor**. Un contrato que vive en un consumidor no es un contrato: es una convención con suerte. |

### 5.3 El argumento decisivo: la dirección de dependencias

`widget` es la **única** pieza que puede estar por debajo de todas las demás, porque no
depende de nada salvo `tinywasm/fmt`:

```flowchart TD
W[tinywasm/widget<br/>Name · Part · State · Class<br/>solo depende de fmt] --> C[tinywasm/css<br/>tokens + emisión]
W --> D[tinywasm/dom<br/>markup + data-state]
W --> F[tinywasm/form<br/>estados Invalid/Locked]
W --> V[tinywasm/view<br/>estado Selected]
C --> S[tinywasm/widget/style<br/>build !wasm]
S --> CO[tinywasm/components]
D --> CO
CO --> L[tinywasm/layout]
```

`form` necesita nombrar `Invalid` y `Locked`. `view` necesita nombrar `Selected`. Si esos
símbolos viven en `css`, entonces **una librería de datos pasa a depender de una librería de
estilo** — y eso no se arregla después. Una pieza neutra en la base es la única forma sin
ciclo ni dependencia mal orientada.

### 5.4 La alternativa más barata, y por qué no se recomienda como destino final

Existe precedente en la casa para partir por **paquete** en vez de por repo:
`tinywasm/svg` (compartido) y `tinywasm/svg/sprite` (backend). Se podría empezar como
`github.com/tinywasm/css/widget`, con un repo menos.

Sirve como **paso 0 si se quiere validar la API antes de publicar un repo**, pero no como
destino: en Go el módulo es la unidad de dependencia, así que `form` y `view` acabarían con
una arista al módulo `css` de todos modos — el problema de 5.3, solo que menos visible.
Recomendación: repo propio desde el inicio; la ruta de subpaquete solo si se quiere un
prototipo desechable de una tarde.

### 5.5 Qué posee y qué NO posee `tinywasm/widget`

| Posee | No posee |
|---|---|
| `Name`, `Part`, `Class` (derivada, nunca escrita a mano) | Datos de dominio |
| `State` (widget) y `Cue` (navegador) | Transporte / `router.Caller` |
| Interfaces de capacidad: `Widget`, `Selectable`, `Dismissible`, `Expandable` | Emisión de CSS (vive en `widget/style`) |
| El enum de tipos ARIA-APG (`Listbox`, `Dialog`, …) | Construcción de DOM (vive en `dom`/`html`) |

Si alguna vez se le añade un campo de datos o una llamada de red, la pieza dejó de ser lego.

---

## 6. Qué bugs históricos deja de ser posible escribir

Todos salen de `docs/ROADMAP.md` — ya se pagaron una vez.

| Bug ya sufrido | Por qué no vuelve |
|---|---|
| *"Desktop panels didn't fill the stage's height — grid row track used `none`/auto"* | `Fill()` emite el conjunto correcto completo; no hay track que elegir |
| *"a gray strip showed… sizing a scroll-snap child in `vw` against a narrower container"* | No existen unidades de viewport en la API; toda medida es relativa al contenedor |
| *"adjacent `RawRule`s are concatenated without a separator"* | No hay `RawRule` |
| *"`Padding(a,b,c,d)`'s output loses its spaces"* | Se arregla en `css` ([PLAN_CSS](PLAN_CSS.md), Etapa 3) y `Pad` toma un solo `Space` |
| *"Hover color was inconsistent (`ColorPrimary`, hardcoded `filter: brightness`, ad-hoc backgrounds) across components"* | El hover lo resuelve `On(Surface)`, definido una vez por superficie |
| *"the form was entering from the right… fixed by adding `direction:rtl`"* | El orden físico lo garantiza la primitiva `Split`, no un truco de flujo |
| CSS nunca emitido por nombrar la función `GenerateCSS` | Interfaz `Styler` tipada ([PLAN_SSR](PLAN_SSR.md)) — error de compilación |
| Marco gris que ignora el tema oscuro (§1.2) | No hay `var()` con fallback escribible; el token existe o no compila |

---

## 7. Reparto por librería

| Librería | Cambio | Plan | Versión | Estado |
|---|---|---|---|---|
| `tinywasm/widget` | **Nueva.** Contrato visual + `widget/style` | [PLAN_WIDGET](PLAN_WIDGET.md) | **v0.1.0** | ✅ Publicado |
| `tinywasm/css` | Tokens faltantes, pares con contraste, `Class` como alias, cierre del DSL viejo | [PLAN_CSS](PLAN_CSS.md) | **v0.2.0** | ✅ Publicado |
| `tinywasm/ssr` | `Styler` tipado en vez de regex sobre el nombre | [PLAN_SSR](PLAN_SSR.md) | — | Pendiente, un solo cambio |
| `tinywasm/components` | Migración de `targetlist`, `fieldset`, `modaldialog` | [PLAN_COMPONENTS](PLAN_COMPONENTS.md) | — | Pendiente, un solo cambio |
| `tinywasm/form` | Emitir `State.Invalid`/`State.Locked` en vez de clases propias | dentro de PLAN_COMPONENTS | — | Pendiente, junto con `components` |
| **`tinywasm/layout`** | Migrar `crudview`, `platformd`, `rightpanel` sobre `css v0.2.0` + `widget v0.1.0` | **este archivo, §8** | — | Pendiente, un solo cambio |

Los tres pendientes ya no están escalonados por dependencia — `widget` y `css`, que eran el
bloqueo real, están publicados. `ssr`, `components`/`form` y `layout` se ejecutan cada uno de
una sola vez (sin canario ni periodo de coexistencia); solo `components` debe cerrarse antes de
`layout` porque `crudview` compone widgets de `components` directamente.

---

## 8. Migración de `tinywasm/layout` — un solo cambio, no por etapas

`tinywasm/widget` (v0.1.0) y `tinywasm/css` (v0.2.0) ya están publicados y probados con su test
de forma-consumidor. No hay razón para escalonar la migración de `rightpanel`, `crudview` y
`platformd` en pasos sucesivos que esperan uno al otro — los tres se migran **en el mismo
cambio**, sobre el mismo `go.mod`, y se validan juntos con un único `gotest` al final. La única
secuencia real es la lógica de un solo commit: `go.mod` primero, después los tres paquetes,
después el test que los cubre a todos.

**`go.mod`:**

```go
require (
	github.com/tinywasm/css    v0.2.0
	github.com/tinywasm/widget v0.1.0
)
```

Retirar los `replace` de desarrollo local hacia `../css` que ya no hagan falta.

**Auditoría ejecutable** — añadir `layout_conformance_test.go` en la raíz, cubriendo los tres
paquetes a la vez (no falla-luego-arregla-uno-por-uno; se corre contra el árbol completo desde
el primer commit de la migración):

1. Ningún `.go` de este repo contiene un literal de color (`#rrggbb`, `rgb(`, `hsl(`).
2. Ningún `.go` contiene `RawRule(`.
3. Toda `var(--…)` referenciada existe en el catálogo de `css` v0.2.0.
4. Ningún `Media(` con un umbral literal que ya tenga token `Bp*`.

**Los tres paquetes, migrados juntos:**

`rightpanel` (166 líneas) es el más simple y sirve para confirmar que el vocabulario de
`widget/style` alcanza sin escape hatch — si no alcanza, **se corrige la API (aguas arriba, en
`widget`/`css`) y se publica una versión nueva, nunca se reabre `RawRule` localmente**. Sus 10
tokens `--rp-*` desaparecen: `tokenAsideBg`/`tokenBorderColor`/`tokenBg` se vuelven `On(Panel)`;
`tokenMainWidth`/`tokenAsideWidth` se vuelven `Split(style.RatioTwoThirds, Space2)`.

`crudview` — anatomía derivada de la estructura actual, nombres según Open UI:

| Clase actual | Parte propuesta | Disposición |
|---|---|---|
| `cv-module-content` | *root* | `Split(style.RatioTwoThirds, Space2)`, `On(Accent)`, `Flush()` |
| `cv-article-contend` | `detail` | `Stack(Space2)`, `Fill()` |
| `cv-box-content` | `fields` | `On(Sunken)`, `Pad(Space2)`, `Scrolls()`, `Round(RadiusMd)` |
| `cv-aside-wrap` | `aside` | `Stack(Space1)`, `On(Panel)`, `Pad(Space1)`, `Fill()` |
| `cv-lista-box` | `list` | `On(Sunken)`, `Scrolls()`, `Round(RadiusMd)` |
| `cv-aside-search` | `search` | `Row(Space0)` |
| `cv-btn-crud` | `action` | `On(Accent)`, `Round(RadiusMd)` |
| `cv-btn-crud-icon-hidden` | — | **desaparece**: es `State.Open`, no una clase |
| `cv-back` | `back` | desaparece con `Split`: el reflow ya no necesita botón de vuelta |

Bloque `Media("(max-width: 640px)")` completo (25 líneas + 20 de comentario explicando
`direction:rtl`): **se borra**. Lo cubre `Split`.

`platformd` (521 líneas, el más grande y el que más `vw`/`vh` usa). Sus 7 tokens `--pd-*` se
resuelven: `tokenMenuSize`/`tokenHeaderHeight` pasan a la escala `Space`; `tokenContentHeight`
(`97vh`, `calc(100vh - 2.8rem)`) desaparece con `Cover()` + `Fill()`. Borrar
`platformd/tokens.go` entero y los bloques `Root(Declare(...))` de `rightpanel`/`platformd`.

**Test de forma-consumidor, en el mismo cambio, no como paso final aparte:**

`crudview/consumer_test.go` (753 líneas) ya recorre la pila real. Extenderlo para aseverar la
hoja emitida: que no contenga `!important`, que sus `@layer` estén en orden, y que cada clase
presente en el markup exista en la hoja **y viceversa** — el par que hoy nadie verifica y que
permitió que `cv-btn-crud-icon-hidden` fuera un estado disfrazado. Este test se escribe junto
con la migración de `crudview`, no después.

---

## 9. Criterios de aceptación

1. `grep -rn "RawRule\|#[0-9a-fA-F]\{3,6\}\|Str(" --include=*.go .` → **vacío** en este repo.
2. `GOOS=js GOARCH=wasm go list -deps ./...` no contiene `tinywasm/widget/style` ni
   `tinywasm/css`.
3. `platformd/tokens.go` no existe; no queda ningún `--pd-*`, `--rp-*`, `--cv-*`.
4. La hoja emitida no contiene `!important` ni ningún `@media` escrito a mano por un widget.
5. `gotest` en verde en `layout`, `components`, `css`, `widget`, `ssr`.
6. Verificación visual en vivo (MCP screenshot), escritorio y móvil emulado, claro y oscuro
   — incluida la comprobación de que el marco gris **ahora sí** cambia en modo oscuro (§1.2).

### La prueba ácida

Dar a un agente sin contexto la firma de `style.Sheet` y este único ejemplo:

```go
func (l *TargetList) Style() *style.Sheet {
	return style.Of(l.WidgetName()).
		Root(Stack(Space1), On(Sunken), Scrolls()).
		Part(partRow, Row(Space2), On(Panel), Pad(Space2), Round(RadiusSm)).
		When(State.Selected, partRow, On(Selected)).
		Cue(Cue.Hover, partRow, On(Hover))
}
```

Si produce un widget correcto sin preguntar nada y sin poder hardcodear un color, el arnés
está cerrado. Si necesita preguntar *"¿puedo declarar esta variable localmente?"* —la pregunta
que originó este plan— el arnés sigue abierto.

---

## 10. Costo, riesgo y qué **no** se hace

**Costo honesto:** son 5 librerías y ~1.200 líneas de CSS reescritas. No es una tarde. `widget`
y `css` (las dos que de verdad tenían que ir primero, porque todo lo demás las consume) ya están
publicadas; lo que queda (`ssr`, `components`/`form`, `layout`) se ejecuta cada uno de una sola
vez — sin fases aditivas intermedias ni periodo de coexistencia con el código viejo — porque los
únicos consumidores de esas tres piezas son otras piezas de este mismo árbol, no terceros: no
hay nadie a quien darle una rampa de migración.

**Riesgo principal:** que el vocabulario de ~22 constructores resulte insuficiente y aparezca
la tentación de reabrir `RawRule`. Mitigación: `rightpanel` (el paquete más pequeño de `layout`)
es la primera prueba real del vocabulario dentro del mismo cambio — si ahí falta algo, **se
amplía el enum aguas arriba en `widget`/`css`, se publica una versión nueva y se sigue** (nunca
se reabre el escape hatch localmente); no hace falta parar la migración completa para eso, la
publicación de una versión menor no rompe lo que ya se migró. Un `RawRule` reintroducido en
cualquier repo invalida todo el plan, porque un arnés evadible no es un arnés.

**Fuera de alcance, a propósito:**

- No se toca `tinywasm/dom` ni el contrato de señales. La reactividad ya funciona y su
  contrato (`Init(ctx)` + `Render()`) ya es un arnés cerrado.
- No se toca `view.Presenter`. Su separación datos/UI es correcta; el problema nunca estuvo
  ahí — es justamente lo que valida que lo visual sea un contrato aparte.
- No se construye un runtime multiplataforma. El objetivo es un contrato **agnóstico de
  renderer**, no un segundo renderer. Si algún día llega uno nativo, `widget` es lo que le
  permitiría reutilizar los mismos widgets; pero eso es consecuencia, no objetivo.
