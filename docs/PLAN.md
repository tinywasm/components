---
PLAN: "feat(herobanner): srcset en las imagenes del hero"
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 4439125634798363599
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.
>
> **PUERTA: no despachar hasta que `github.com/tinywasm/image` publique
> `image.Responsive`.** Este plan la usa; sin ella no compila.

# Plan — `herobanner`: `srcset` en el hero

Hacer que las imágenes del hero declaren las tres variantes que el pipeline ya
genera, en vez de servir la de escritorio a todo el mundo.

## Por qué el hero primero

La imagen del hero es el elemento **LCP** (Largest Contentful Paint) de un sitio
de vitrina: es la que Google cronometra. Servirla en 1920 px a un teléfono de
400 px es la forma más cara de perder ese número.

Medido sobre `veltylabs/mjosefa-website` en producción: el HTML no tiene **ni un**
`srcset`, y las fotos del hero pesan 47–110 KB donde a un teléfono le bastaban
~25 KB.

## Lo que NO se toca

- **`loadingFor` es correcto y se queda como está**: capa 0 `eager`, el resto
  `lazy`. El hero es el LCP y debe ser la única imagen `eager`. **No la cambies.**
- **`autoRotateLayers`** y toda la mecánica de rotación.
- El resto de los componentes. `usermenu` emite un avatar y `sitenav` un logo
  SVG: ninguno tiene variantes y ninguno las necesita. **No los toques.**

---

## 0. Reglas de desarrollo

`herobanner.Render()` **compila para wasm**. Por lo tanto:

- **Sin `fmt`, `errors`, `strconv`, `strings`, `log`** de la stdlib. Usa
  `github.com/tinywasm/fmt`.
- **Sin `map[K]V`**, sin `reflect`, sin `encoding/json`.
- **No importes `widget/style`** desde `Render()`: es `!wasm`. Es exactamente por
  eso que `autoRotateLayers` está duplicado con un comentario que lo explica —
  no lo "arregles" importando el paquete.
- Los productores SSR (CSS, SVG, JS, HTML pesado) viven en archivos con nombre de
  extensión y `//go:build !wasm`. Nunca en un `ssr.go`.
- Código en inglés; documentación y comentarios de prosa en español.
- Sin strings mágicos: todo string repetido es una constante nombrada.

---

## 1. El cambio

En `herobanner.go`, dentro de `Render()`, hoy:

```go
img := NewElement("img").
	Attr("src", imgSrc).
	Attr("alt", "").
	Attr("loading", loadingFor(i)).
	NoCloseTag()
mediaLayer.Child(img)
```

Pasa a construirse con el constructor responsivo:

```go
img := image.Responsive(imgSrc, "").
	Attr("loading", loadingFor(i)).
	AsElement()
mediaLayer.Child(img)
```

`image.Responsive` ya emite `NoCloseTag()`, `src`, `srcset`, `sizes` y `alt`.

**`sizes` se queda en el `100vw` por defecto**: el hero ocupa el ancho completo
de la ventana, que es justo lo que ese valor declara. No lo sobrescribas.

`HeroBanner.Images` sigue siendo `[]string` y sigue recibiendo **rutas base**
(`/img/foto.jpg`, sin sufijo de variante). La firma pública no cambia: para el
consumidor esto es puramente aditivo.

**Anti-footgun:** `image.Responsive` recibe la ruta **base**. Si un consumidor
venía pasando `/img/foto.M.jpg` (con el sufijo escrito a mano), el resultado
serán rutas `/img/foto.M.S.jpg` que no existen. Eso se corrige **en las apps**,
no aquí, y no es parte de este plan — pero documéntalo en el `README.md`.

---

## 2. `alt` vacío

El `alt=""` actual es **correcto y se mantiene**: estas imágenes son decorativas
(un fondo rotatorio detrás del titular), y un `alt` descriptivo en seis capas
haría que un lector de pantalla leyera seis veces lo mismo. `alt=""` es lo que le
dice que las ignore.

---

## 3. Tests

| # | Caso | Espera |
|---|---|---|
| 1 | `HeroBanner` con una imagen | el `<img>` trae `srcset` con las tres variantes |
| 2 | la capa 0 | `loading="eager"` |
| 3 | las capas 1..5 | `loading="lazy"` |
| 4 | `sizes` | `100vw` |
| 5 | `alt` | vacío, en todas las capas |
| 6 | `HeroBanner` sin imágenes | no emite la capa de medios, sin pánico |
| 7 | la salida sigue siendo markup válido | `<img>` sin etiqueta de cierre |

Los casos 2 y 3 son de **regresión**: verifican que este cambio no tocó el
`loading`, que ya estaba bien.

---

## 4. Documentación

- `README.md` de `herobanner` — que `Images` recibe **rutas base sin sufijo de
  variante**, y que de ahí salen las tres.
- Si escribes diagramas: **nunca uses `subgraph`** (rompe el TUI). `flowchart TD`
  y `<br/>` para los saltos.

---

## 5. Criterios de aceptación

- [ ] `go vet ./...` limpio; tests en verde con los 7 casos.
- [ ] `GOOS=js GOARCH=wasm go build ./...` sin errores.
- [ ] `grep -n "NewElement(\"img\")" herobanner/herobanner.go` → vacío.
- [ ] `grep -n "func loadingFor" -A 5 herobanner/herobanner.go` → **sin cambios**.
- [ ] `git diff --stat` toca **sólo** `herobanner/`: ni `usermenu/`, ni
      `sitenav/`, ni ningún otro componente.
- [ ] La firma de `HeroBanner` no cambia: `Images` sigue siendo `[]string`.
- [ ] `grep -n "\"fmt\"\|\"strings\"\|\"errors\"\|\"strconv\"" herobanner/*.go` → vacío.

## 6. Fuera de alcance

`landing` (fase C, su propio repo), las apps consumidoras, y cualquier otro
componente. Tampoco AVIF ni art direction.
