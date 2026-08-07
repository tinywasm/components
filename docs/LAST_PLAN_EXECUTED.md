---
PLAN: "refactor(components): adoptar el apilamiento declarado y retirar los parches de espacio"
EXECUTOR: jules
REVIEWER: none
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.
>
> Es la **etapa 3 de 4**. Orden obligatorio: **css → widget → components →
> layout**. Requiere `--chip-height` (`tinywasm/css`) y las etapas 1-5 de
> `tinywasm/widget` publicadas. No empezar antes: casi todo este plan consiste en
> **borrar** cosas que dejan de hacer falta.

# Plan — quitar los parches que el DSL ya no obliga a escribir

## 1. Por qué

Depurando una vista CRUD en móvil aparecieron cuatro defectos visuales. Tres se
resolvieron **en este repo, a mano**, porque el DSL no ofrecía otra cosa:

| Parche actual | Dónde | Qué compensaba |
|---|---|---|
| `PadEdge(EdgeBottom, Space12)` en `PartList` (solo móvil) | `targetlist/css.go` | Que un botón flotante de otro widget se solapaba con el badge de la última fila. |
| `PadInline(Space8)` en `PartLabel` | `targetlist/css.go` | Que el ⋮ `Docked` se pinta *dentro* de la caja y pisaba el texto. |
| Comentario "la leyenda debe medir lo mismo que el badge, por eso no lleva padding vertical" | `fieldset/css.go` | Que la altura de un chip era emergente y las dos coincidían por casualidad. |

Ninguno es un error de quien los escribió: eran la única salida. Con las etapas
de `tinywasm/widget` publicadas dejan de serlo, y **un parche que sobrevive a su
causa se convierte en deuda**: el siguiente que lo lea no sabrá si sigue haciendo
falta.

La regla del harness que aplica aquí es explícita: *un consumidor nunca recrea
localmente lo que falta arriba; si la librería no expone lo que necesitas, se
para y se reporta*. Esto es el paso de vuelta — una vez arreglado arriba, el
parche de abajo se retira.

## 2. Contexto del repo para un agente sin contexto previo

- Módulo: `github.com/tinywasm/components`. `docs/PLAN.md` va junto a `go.mod`.
- Un componente = un paquete. Los relevantes: `targetlist/`, `fieldset/`.
- **Separación SSR por extensión**: el CSS va en `css.go` con `//go:build !wasm`,
  los iconos en `svg.go`. Nunca CSS ni SVG en el `.go` principal.
- **Sin `front.go`**: la interactividad WASM va en el archivo del componente vía
  `OnMount()`.
- Todo el estilo sale del DSL de `tinywasm/widget/style`. **No hay escotilla de
  CSS crudo y no se debe añadir.**
- Nada de librería estándar en paquetes WASM: `tinywasm/fmt`.
- Empotrado por valor: `dom.Element` como valor, nunca `*dom.Element`.
- Prohibidas las cadenas repetidas en la lógica: constante con nombre.
- `conformance_test.go` en la raíz aplica reglas transversales a todos los
  componentes; leerlo antes de tocar nada, porque **falla por cosas que parecen
  correctas** (por ejemplo, exige que `Cue(Hover, X)` tenga su `Cue(Focus, X)`
  emparejado, para que quien navega con teclado no se quede sin la señal que sí
  recibe quien usa ratón).

## 3. Etapas

### Etapa 1 — `targetlist`: retirar la reserva de espacio manual

En `targetlist/css.go`, la regla móvil de `PartList` lleva hoy:

```go
On(css.Mobile, PartList,
    style.Pad(style.SpaceNone),
    style.PadEdge(style.EdgeBottom, style.Space12), // ← retirar
).
```

Con la etapa 5 de `widget`, `Scroll()` ya emite
`padding-block-end: var(--floating-bottom, 0px)`, y el host que flota el botón
declara cuánto ocupa. La reserva deja de ser asunto de este componente.

Retirar **solo** el `PadEdge`. El `Pad(SpaceNone)` se queda: responde a otra cosa
(recuperar milímetros para que el ⋮ quepa en la franja visible de la lista en
móvil) y está documentado en su comentario.

Actualizar el comentario para que no quede describiendo algo que ya no está.

**Aceptación:** `grep -n "PadEdge" targetlist/css.go` no devuelve nada; el badge
de la última fila deja de solaparse con un botón flotante del host cuando el host
declara su `FloatingChrome` (se verifica en el plan de `layout`).

### Etapa 2 — `targetlist`: revisar la holgura del texto tras `OnEdge` sin transform

`PartLabel` lleva `PadInline(Space8)` para librarse del ⋮ `Docked` en el borde
inicial. `Docked` mantiene el elemento **dentro** de la caja, así que esa holgura
sigue siendo necesaria y **no se retira**.

Lo que sí hay que volver a medir es `PartBadge`: con la etapa 2 de `widget`,
`OnEdge` ya no usa `transform`, así que el badge pasa a ocupar espacio real y
puede empujar la altura de la fila. Comprobar que la fila sigue midiendo
`--control-height` y que el badge sigue montado sobre la línea inferior.

Si la fila crece, la causa será que la mitad inferior del chip ahora cuenta:
compensarlo en `PartRow` y dejarlo escrito, **no** volviendo a meter un
`transform` por la puerta de atrás.

**Aceptación:** todas las filas miden lo mismo con y sin badge, y con etiquetas de
distinta longitud.

### Etapa 3 — `fieldset`: la leyenda deja de coincidir por casualidad

Con `--chip-height` publicado, la altura del chip es un valor declarado. El
comentario de `PartLabel` que explica que **no** lleva padding vertical para
igualar al badge de `targetlist` describe un acuerdo verbal entre dos repos:
reemplazarlo por la referencia al token, que es lo que ahora garantiza la
igualdad por construcción.

Verificar además, con la etapa 4 de `widget` (anillos con `box-shadow` en vez de
`outline`):
- el borde del estado `Locked` **ya no** pinta por encima de la leyenda;
- ese borde respeta el `border-radius` — el anillo cuadrado que se veía en un
  iPhone 7 (iOS 15, Safari sin `border-radius` en outlines) debe haber
  desaparecido;
- la caja **no** cambia de tamaño al entrar en `Locked` (era la razón de usar
  `outline`; `box-shadow` la conserva).

**Aceptación:** ninguna regla emitida por `fieldset` contiene `outline`; la
leyenda es legible en el estado `Locked`.

### Etapa 4 — retirar los z-index implícitos

Con la etapa 3 de `widget`, el apilamiento lo decide el DSL. Revisar que ningún
componente de este repo dependa de orden de DOM para quedar por encima de un
hermano. Si aparece alguno, se declara — no se reordena el marcado.

### Etapa 5 — tests

1. `targetlist`: la regla móvil de `PartList` **no** emite `padding-block-end`
   propio, y la regla base sigue emitiendo el `padding` normal.
2. `targetlist`: `PartBadge` no emite `transform`.
3. `fieldset`: el estado `Locked` no emite `outline`.
4. Ejecutar `conformance_test.go` completo: es el que atrapa las reglas
   transversales.

Ejecutar `go build ./... && go test ./... -count=1` en la raíz. Deben pasar
**todos** los paquetes, no solo los dos tocados.

| Etapa | Archivos | Puerta |
|---|---|---|
| 1 | `targetlist/css.go` | tras publicar `widget` |
| 2 | `targetlist/css.go` | tras 1 |
| 3 | `fieldset/css.go` | tras publicar `widget` |
| 4 | todo `*/css.go` | tras 1 y 3 |
| 5 | `*/*_test.go` | tras 4 |

## 4. Lo que este plan NO hace

- **No toca el desplegable del ⋮ anclado al viewport en móvil.** Esa decisión
  (anclar al viewport en vez de a la fila) responde a que el disparador se mueve
  ~340px entre los dos estados del scroll-snap y CSS no puede observarlo. Sigue
  siendo una elección explícita del consumidor y es correcta.
- No cambia el marcado de ningún componente, solo su hoja.
