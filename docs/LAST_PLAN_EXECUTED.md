---
PLAN: "fix(targetlist): el desplegable ⋮ deja de abrirse encima de su propia fila"
TAG: v0.5.0
---

> **NO DESPACHAR TODAVÍA.** Este plan espera revisión, y además **depende de que
> la etapa A esté publicada**.
>
> **Etapa B de un cambio en 2 repos.**
>
> | | Repo | Plan | Qué |
> |---|---|---|---|
> | A | `widget` | `widget/docs/PLAN.md` | el sheet aprende el árbol de partes; `Validate()` deja de callarse |
> | **B** | **`components`** | **este plan** | `targetlist` adopta la construcción legal; `usermenu` gana el test que nunca tuvo |

# Plan — `components`: `targetlist`, `usermenu` y el contenedor de bloque

## El bug, medido

En desktop (1440x900), abrir el menú `⋮` de una fila pinta Editar/Eliminar
**encima de la fila que lo abrió**, tapándole la etiqueta.

```
row      top 113.2   bottom 163.2   (alto 50)
summary  top 118.0   bottom 142.0   (alto 24)
options  top 142.0                  → 21.2px POR ENCIMA del fondo de su fila
label    top 126.2   bottom 150.2   → el desplegable le come 8.2px de texto

options.offsetParent === .targetlist__menu     ← no es la fila
```

En móvil el mismo desplegable **no** sufre esto: allí `PartOptions` es
`Docked(Viewport, …)` → `position: fixed`, que se resuelve contra la pantalla y
no consulta esta cadena. (El solape que se veía en móvil era otro: el
`Backdrop` sin `Veil()`, ya arreglado.)

## Por qué el arreglo no vive aquí

La causa es que `Docked(Parent, …)` en `PartMenu` convierte al `<details>` en el
contenedor de bloque del `Flyout`, desplazando al `Anchor()` de la fila. El
`100%` del `Flyout` mide 24px (el disparador) en vez de 50px (la fila).

Este repositorio **no puede arreglarlo bien**, sólo compensarlo. Cualquier
número que pusiéramos aquí sería un parche sobre un contrato que `widget` no
expresa — que es justo el bucle que
`app-releases/docs/CONSTRUCTION_HARNESS.md` describe:

> *Un hueco de API siempre aparece en la hoja (la aplicación), donde el agente no
> tiene autoridad para publicar aguas arriba — así que parchea en local. La deuda
> técnica no es entonces un accidente: el flujo la garantiza.*

Y la regla directa:

> *Un contrato que falta en una frontera es un defecto de la librería, no del
> consumidor.*

La prueba de que la regla ya se filtró a este repo está escrita a mano en
`targetlist/css.go`:

> *No `Anchor()` aquí — `Docked` ya lo hace contenedor de bloque, y los dos se
> pelean por `position`.*

Ese comentario es una regla memorizada. Cuando la etapa A cierre el hueco, **se
borra**: pasa a ser algo que el compilador/`Validate()` dice, no algo que el
próximo lector tiene que recordar.

## El test ya está escrito y en rojo

`targetlist/anchor_contract_test.go` (nuevo, ya en el repo).
**`TestFlyoutHangsFromTheRowNotFromTheTrigger`** — hoy 🔴:

```
.targetlist__menu sits between the Anchor (.targetlist__row) and the Flyout
(.targetlist__options) and is position: absolute pinned on one block edge only
(inset-block-start: "var(--space-1,0.25rem)", inset-block-end: "auto").
```

Cómo está construido, y por qué así:

- **Deriva la anidación del markup real**, recorriendo la salida de `buildRow()`
  con una pila de etiquetas. No lleva la cadena `row > menu > options` escrita a
  mano: si alguien reestructura la fila, el test **re-deriva** en vez de quedarse
  validando una suposición que ya no es cierta.
- **Lee el CSS real emitido**, no una expectativa.
- **La aserción es agnóstica al arreglo.** No dice "el ancestro posicionado más
  cercano debe ser el Anchor" — esa frase prohibiría la opción 2b de la etapa A.
  Dice lo que **toda** solución correcta cumple: entre el `Flyout` y su `Anchor`,
  ningún elemento posicionado puede estar fijado a un solo borde de bloque (con
  el otro en `auto`), porque entonces lo dimensiona su propio contenido y el
  `100%` nunca llega al fondo de la fila.
- **Móvil queda exento a propósito**, y el test sólo mira las reglas base
  (fuera de `@media`), con el motivo escrito en el propio comentario.

---

## Etapa 1 — dependencia: verificar antes de escribir nada

```sh
go doc github.com/tinywasm/widget/style Within
```

Si eso falla, **parar**. La etapa A no está publicada.

> ⚠️ **No declarar aquí una versión local de `Within`, ni un helper que
> reproduzca lo que hace.** Recrear aguas abajo un símbolo que falta aguas
> arriba es el defecto exacto que este cambio existe para eliminar.
>
> ⚠️ **No "arreglar" el solape con un `PadEdge`, un margen negativo, ni un
> `Space` inventado en `PartOptions`.** Cualquier número que compense 21.2px es
> un parche que se romperá en cuanto cambie el alto de la fila o del disparador.

## Etapa 2 — adoptar la construcción elegida

La etapa A deja tres candidatas. **La decisión se cierra aquí, con medidas**,
porque el dato que falta es de este repo: cuánto *sliver* móvil queda.

### 2a — el disparador vuelve al flujo *(la que hay que medir primero)*

Quitar `Docked` de `PartMenu` y dejar el `<details>` como primer hijo en flujo de
la fila.

- El ancestro posicionado más cercano del `Flyout` pasa a ser **la fila**, y
  `Flyout` cumple su documentación por construcción.
- Se puede quitar el `PadInline(Space8)` de `PartLabel`: existía **sólo** para
  esquivar el icono `Docked`.
- El comentario de `PartMenu` que justifica el `Docked` dice que menú y badge
  salen del flujo *"así que la etiqueta es lo único que dimensiona la fila"*.
  Revisar ese razonamiento: la fila ya tiene `min-height: var(--control-height)`
  = 50px y el menú mide 24px, así que **en flujo tampoco la dimensiona**. El que
  sí necesita salir del flujo es el badge, que envuelve a su propia línea.

**Medir, no suponer:** la fila tiene `Pad(Space3)` = 12px, así que en flujo el
icono arranca a 12px en vez de los 4px de `Docked(Space1)`. En el *sliver* de
~37.5px que deja `MasterDetail(Most)` a 375px de viewport eso son 36px contra
28px. Entra, pero sin margen.

**Comprobar en el navegador, en móvil, con la tira desplazada al panel de
formulario, que el `⋮` sigue siendo alcanzable.** Ese es el requisito que el
comentario de `PartMenu` protege — *"la única palanca que desbloquea el
formulario ahora de sólo lectura, en el panel que el usuario acaba de dejar"* —
y no se puede sacrificar por arreglar el desktop. Si no entra, esta opción cae y
se pasa a 2b/2c.

### 2b / 2c — alternativas

Sólo si 2a no pasa la medida. Están descritas con sus costes en
`widget/docs/PLAN.md`; ninguna de las dos se implementa en este repo sin que la
primitiva correspondiente exista ya publicada en `widget`.

## Etapa 3 — declarar la anidación

Independientemente de cuál gane, `targetlist` declara explícitamente la relación
que hoy sólo existe en el markup:

```go
Within(PartMenu, PartOptions, style.Flyout(style.SideStart))
```

Es lo que permite a `Validate()` comprobar la cadena en vez de adivinarla. Y si
alguien vuelve a meter un `Docked` en medio, ahora **falla en `Validate()`**, no
en el navegador de un usuario.

## Etapa 4 — borrar la regla memorizada

Quitar de `targetlist/css.go` el comentario *"No `Anchor()` aquí — `Docked` ya lo
hace contenedor de bloque, y los dos se pelean por `position`"*, y todo lo que
haya quedado explicando cómo esquivar el problema. Si la etapa A hizo su
trabajo, esa prosa ya no es información: es ruido que sobrevivió a su causa.

## Etapa 5 — `usermenu` gana los tests que nunca tuvo

```
?   github.com/tinywasm/components/usermenu   [no test files]
```

`usermenu` es **el caso que funciona**: `Root(Anchor())` → `PartPanel` con
`Flyout(SideEnd)`, sin nada posicionado en medio. Es la forma que el DSL
documenta, es la referencia contra la que se contrasta el bug de `targetlist`…
y no la vigila nadie.

Añadir `usermenu/anchor_contract_test.go` con la **misma** comprobación que
`targetlist` (los helpers de recorrido de markup se pueden compartir o duplicar;
son ~40 líneas). Debe pasar **desde el primer día, sin tocar `usermenu`**: si no
pasa, o el helper está mal o `usermenu` tiene el mismo bug sin que nadie lo haya
visto todavía. Las dos posibilidades merecen saberse.

---

## Orden de ejecución

| # | Etapa | Verde cuando |
|---|---|---|
| 1 | verificar dependencia A | `go doc … Within` responde |
| 2 | adoptar 2a (o 2b/2c) | `TestFlyoutHangsFromTheRowNotFromTheTrigger` pasa |
| 3 | `Within(PartMenu, PartOptions, …)` | `Validate()` sin errores |
| 4 | borrar el comentario memorizado | revisión |
| 5 | tests de `usermenu` | pasan sin tocar `usermenu` |

Verificación final **en navegador**, las dos que importan:

- **Desktop:** `options.top >= row.bottom` y
  `options.offsetParent === .targetlist__row`.
- **Móvil, tira desplazada al formulario:** el `⋮` sigue visible y pulsable
  dentro del *sliver*.

---

## Lo que este plan NO hace

- **No toca `usermenu` ni `fieldset` ni `selectsearch`.** `fieldset` usa
  `Anchor` + `Docked(Parent)` pero **no tiene `Flyout`**, así que la composición
  que rompe no existe ahí. Comprobado, no supuesto.
- **No cambia el comportamiento móvil de `PartOptions`.** El
  `Docked(Viewport, …)` de móvil es correcto y está justificado en su comentario
  (el disparador se mueve entre ~10px y ~370px según cómo esté desplazada la
  tira); no depende de esta cadena y no se toca.
- **No revisa el `Veil()` del `Backdrop`.** Ya está arreglado y verificado en
  navegador.
- **No persigue el solape residual de ~4.4px entre badge y botón flotante** en
  móvil. Es interno de `targetlist`, no tiene que ver con este contenedor de
  bloque, y sigue aparcado.
