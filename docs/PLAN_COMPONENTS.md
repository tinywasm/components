---
PLAN: "`tinywasm/components` (+ `tinywasm/form`): migrar los widgets al contrato visual"
---

> Depende de: `github.com/tinywasm/widget v0.1.0` y `github.com/tinywasm/css v0.2.0`
> (**ambos ya publicados**), y de que `ssr` exponga `style.Styler` ([`PLAN_SSR`](PLAN_SSR.md)).
> Bloquea a: [`PLAN.md §8`](PLAN.md) (la migración de `layout`) — un `crudview` migrado no puede
> componer widgets sin migrar.
>
> **Ejecución: un solo cambio, no por etapas.** Con `widget`/`css` ya publicados no queda
> ninguna dependencia externa pendiente — todos los componentes (`fieldset`, `targetlist`,
> `modaldialog` y el resto) y `tinywasm/form` migran en la misma ventana de trabajo, cubiertos
> por un único `components/conformance_test.go` que se escribe primero y pasa a verde al final
> del mismo cambio, no etapa por etapa.

---

## Por qué este es el plan que resuelve el síntoma reportado

> *"llms/agentes que crean componentes en `tinywasm/components` o `layout` los crean
> hardcodeando colores, estilos, variables"*

`components` es donde ocurre. Y ocurre por una razón estructural, no por descuido: **es el
punto más bajo de la cadena donde un hueco de API se manifiesta, y el punto donde nadie tiene
autoridad para publicar aguas arriba.** El propio arnés lo anticipa:

> *"An API gap always surfaces at the leaf (the application), where the agent has no authority
> to publish upstream — so it patches locally. Technical debt is then not an accident: the
> workflow guarantees it."*

Por eso este plan iba **después** de `widget` y `css`: migrar `components` antes de que
existiera el vocabulario correcto solo habría reubicado los parches. Ahora que ambos están
publicados (`widget v0.1.0`, `css v0.2.0`), ese orden ya se cumplió — lo que queda es un solo
cambio que cubre todo `components` + `form` a la vez, no una migración componente por
componente con esperas entre medio.

---

## Auditoría ejecutable (se escribe primero, cubre todo el árbol de una vez)

Antes de migrar el primer componente, añadir `components/conformance_test.go` que falla hoy y
fija el objetivo — se corre contra **todos** los paquetes de componente a la vez, no se activa
progresivamente a medida que cada uno se migra:

1. Cero literales de color (`#rrggbb`, `rgb(`, `hsl(`, nombres CSS).
2. Cero `RawRule(`.
3. Cero `Media(`.
4. Cero unidades de viewport (`vw`, `vh`, `vmin`, `vmax`).
5. Toda `var(--…)` referenciada existe en el catálogo de `css` v0.2.0.
6. Cero constantes de clase escritas a mano (`Class = "..."`); toda clase se deriva de un
   `widget.Name`.

La regla 5 es la que hay que ejecutar **primero y sobre todo el ecosistema**: produce la lista
definitiva de tokens fantasma. En `layout` ya se conocían cuatro y ya se resolvieron en
`css v0.2.0` ([`PLAN.md §1.2-1.3`](PLAN.md), [`PLAN_CSS`](PLAN_CSS.md)); si `components` revela
más huecos, se publican como parche de `css` antes de terminar este mismo cambio — no se
posponen ni se parchean localmente.

La salida de esta primera pasada no es código: es la lista de huecos que faltan aguas arriba
(si los hay) antes de tocar los componentes.

---

## `fieldset`

Es pura superficie y estado, sin disposición compleja — el primero en escribirse dentro de
este mismo cambio, no una etapa aparte con su propio ciclo de revisión.

| Hoy | Después |
|---|---|
| Clases propias para el estado bloqueado | `widget.Locked` → `data-locked` |
| Tinte "frosted glass" con `ColorSurface` a mano | `On(Sunken)` |
| Hover con `ColorHover` (arreglado en el ROADMAP tras haber divergido) | `Cue(widget.Hover, …)` — resuelto por la `Surface`, una sola vez |
| Mensaje de error posicionado en absoluto con offsets a mano | parte `error` + `On(Danger)` |
| El chip de etiqueta que "cruza el borde" y obligó a `PaddingTop(Space3)` en `crudview` | parte `label` con `Raise(Raised)`; el `crudview` deja de compensarlo desde fuera |

Ese último punto es representativo del daño que causa la falta de anatomía: hoy `crudview`
**compensa desde fuera** un detalle interno de `fieldset`, con un comentario de 8 líneas
explicando por qué el padding superior es `Space3` y los otros tres lados `Space2`. Con una
anatomía nombrada, ese ajuste vive en `fieldset`, donde se decide.

**`tinywasm/form`** cambia aquí también, y es un cambio pequeño: en vez de emitir clases
propias para validez y bloqueo, emite `widget.Invalid` y `widget.Locked`. `form` pasa a
importar `widget` (que solo depende de `fmt`) y **no** importa `css`. Ese es exactamente el
motivo de que `widget` sea un repo aparte ([`PLAN.md §5.3`](PLAN.md)).

---

## `targetlist`

El más ilustrativo: es un `Listbox` de ARIA-APG y hoy no lo declara.

```go
const (
	nameTargetList = widget.Name("targetlist")
	partRow        = widget.Part("row")
	partMenu       = widget.Part("menu")
	partOptions    = widget.Part("options")
)

func (l *TargetList) WidgetName() widget.Name { return nameTargetList }
func (l *TargetList) WidgetKind() widget.Kind { return widget.Listbox }

func (l *TargetList) Style() *style.Sheet {
	return style.Of(nameTargetList).
		Root(Stack(Space1), On(Sunken), Scrolls(), Round(RadiusMd)).
		Part(partRow, Row(Space2), On(Panel), Pad(Space2), Round(RadiusSm)).
		Part(partMenu, Stack(Space0), On(Panel), Raise(Floating), Clip()).
		When(widget.Selected, partRow, On(Selected)).
		Cue(widget.Hover, partRow, On(PanelHover))
}
```

Dos deudas que se saldan al declarar `Kind` y anatomía:

- **Accesibilidad.** `role="listbox"`, `role="option"` y `aria-selected` los emite el renderer
  a partir de `Kind`. Hoy la lista es un conjunto de `div` sin semántica: no es navegable por
  teclado ni por lector de pantalla.
- **El menú ⋮ que "vuelca hacia arriba en la última fila".** El ROADMAP lo resuelve hoy con
  CSS a medida y `:has()`. Es un `Menu` de APG y su posicionamiento es responsabilidad de la
  parte `menu` con `Raise(Floating)` — no del contenedor que la aloja.

El grupo acordeón `<details name="tl-menu-group">` se conserva tal cual: es una primitiva
nativa del navegador, cuesta cero JS y hace exactamente lo correcto. No todo tiene que
convertirse en API.

---

## `modaldialog`

`Kind = widget.Dialog`. Aporta `aria-modal`, gestión de foco y cierre con `Esc` desde la firma
en vez de a mano. La API pública (`Open()`/`Close()`) no cambia — `crudview` no se entera.

Su anatomía: `backdrop`, `panel`, `header`, `body`, `actions`. Con ella, el hack que hoy
documenta `crudview/css.go:305-310` desaparece:

> *"Fixed, not static: takes this mount point (and modaldialog's own hidden-state anchor div
> inside it) out of `clsModuleContent`'s grid item participation entirely, so it can never be
> auto-placed into an implicit row (which was doubling the bottom gutter)"*

Un widget de overlay que necesita que su **consumidor** lo saque del flujo de grid tiene un
defecto de anatomía. `Dialog` se posiciona solo; el punto de montaje deja de existir como
concepto en `crudview`.

---

## El resto de los componentes, mismo cambio

Para cada componente restante, mismo procedimiento, todos en el mismo cambio que `fieldset`/
`targetlist`/`modaldialog` — no uno por PR:

1. `Name` + `Kind` + partes nombradas.
2. Reemplazar toda constante de color local por `On(Surface)`.
3. Reemplazar toda clase-estado por `When(State, …)`.
4. Reemplazar todo `Media`/`vw`/`vh` por la primitiva `Flow` correspondiente.
5. Borrar el bloque `Root(Declare(...))` local: la escala la posee `css`, no el componente.

Al terminar, `components/conformance_test.go` (escrito al principio) pasa a verde para todo el
repo de una sola vez, y queda como guardia permanente en CI: **cualquier componente nuevo —lo
escriba un humano o un agente— falla el build si hardcodea un color**. Ese test es la respuesta
operativa al problema que originó todos estos planes.

---

## La guía final, en una página

Con el arnés cerrado, la documentación de `components` se reduce a una tabla y un ejemplo:

| Quiero… | Uso |
|---|---|
| una tarjeta sobre la página | `On(Panel)` |
| un pozo hundido dentro de la tarjeta | `On(Sunken)` |
| ritmo vertical | `Stack(Space2)` |
| dos paneles que se apilan en móvil | `Split(style.RatioTwoThirds, Space2)` |
| una rejilla que decide sus columnas sola | `Grid(TrackMd, Space2)` |
| que además tome el alto disponible | `Fill()` |
| que desborde por dentro en vez de crecer | `Scrolls()` |
| que **no** se adapte | `Fixed()` |
| marcar una fila seleccionada | `When(widget.Selected, partRow, On(Selected))` |

Más el ejemplo de `targetlist` de más arriba, entero. Eso es toda la guía.

Lo que **se borra** de la documentación actual: las reglas de nombres de proveedores SSR (las
sustituye el tipo, ver [`PLAN_SSR`](PLAN_SSR.md)), las advertencias sobre `RawRule`s adyacentes
y su `;`, y cualquier "acuérdate de…". Si hay que recordarlo, es un hueco del arnés, no una
línea de manual.

---

## Criterios de aceptación

1. `components/conformance_test.go` en verde y ejecutándose en CI.
2. Ningún componente declara un `Token`, una `Class` literal ni un color.
3. Cada componente declara su `Kind` y emite los atributos ARIA correspondientes.
4. `GOOS=js GOARCH=wasm go list -deps ./...` no contiene `widget/style` ni `css`.
5. `gotest` en verde en `components`, `form` y `layout`.
6. La guía de `components` cabe en una página.
