---
PLAN: "feat: componentes de sitio público — infobar, sitenav, herobanner, statgrid"
STATUS: review
SESSION: 15129426777814571399
PR: https://github.com/tinywasm/components/pull/20
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.

# Plan — faltan las cuatro piezas de un sitio público

## Contexto

Este catálogo cubre bien las aplicaciones (CRUD, formularios, listas,
diálogos), pero **nunca se construyó un sitio público** con él. Al preparar uno
real —una clínica: cabecera con menú móvil, hero con imágenes, cifras de
impacto, barra de contacto— aparecen cuatro huecos.

Dos piezas ya existen y **se reutilizan tal cual, no se duplican**:
`contentcard` para las tarjetas de servicio y `fieldset` + `tinywasm/form`
para los formularios.

Lo que este plan **no** hace: componer la página. Ordenar secciones, anclas y
rutas es responsabilidad del layout, que tiene su propio plan. Aquí solo se
entregan piezas reutilizables que no saben en qué página viven.

## Los cuatro componentes

Todos siguen el patrón ya establecido en este repo —mira `contentcard` como
referencia canónica: `widget.Name` + constantes `widget.Part`, `Render()
*dom.Element`, `RenderCSS()` en un `css.go` con `//go:build !wasm` usando
`style.For(c)`, y un test de contrato.

### 1. `infobar` — barra de datos de contacto

La franja superior: dirección, teléfono, correo, cada uno con su ícono.

- Contenido inyectado, no hardcodeado: recibe una lista de ítems
  `{Icon, Text, Href}`. `Href` opcional — un teléfono es `tel:`, un correo
  `mailto:`, una dirección puede no enlazar.
- Aporta `IconSvg()` con los íconos que use (patrón de `item_catalog/svg.go`
  o `targetlist/svg.go`).
- En móvil la franja debe poder colapsar a lo esencial sin romper el layout;
  resuélvelo con CSS, no con JS.

### 2. `sitenav` — cabecera con navegación y menú móvil

Logo + enlaces + botón hamburguesa que despliega en móvil.

- **Es el único de los cuatro que necesita `RenderJS()`.** El comportamiento
  es abrir/cerrar y cerrar al pulsar un enlace. Que sea JS y no WASM es
  deliberado: esta barra tiene que funcionar antes y sin que cargue ningún
  binario, en la primera pintura, en un móvil con red mala.
- Accesibilidad no opcional: el botón lleva `aria-controls`, `aria-expanded`
  actualizado al abrir/cerrar, y `aria-label`. Es una cabecera de sitio
  público, la va a leer un lector de pantalla.
- El logo acepta dos fuentes (ancha y compacta) para conmutar por viewport con
  `<picture>`; usa `tinywasm/image` para construirlo, no armes el markup a
  mano.

### 3. `herobanner` — portada con imágenes y llamada a la acción

Imagen (o varias, rotando) de fondo, titular, bajada y botones.

- **La rotación es CSS, no JS.** Un carrusel de fondo es
  `@keyframes` + `animation-delay` escalonado; no hay estado que gestionar y
  no vale la pena un byte de script. Si el número de imágenes es variable, los
  retardos se calculan al emitir el CSS.
- **Respeta `prefers-reduced-motion`**: con la preferencia activa, se muestra
  la primera imagen fija y no se anima. No es un extra, es lo que evita marear
  a quien lo pidió explícitamente.
- El texto va sobre la imagen: garantiza contraste con una capa de velo
  tokenizada, no con un color a ojo.

### 4. `statgrid` — cifras destacadas

Rejilla de pares valor/etiqueta ("80+" / "Años de historia").

- Puramente presentacional, sin estado ni JS.
- El número de columnas se adapta al viewport; que la rejilla no dependa de
  cuántas cifras reciba.

## Restricciones

- **Nada de valores literales de estilo.** Ni un hex, ni un `rgba()`, ni un
  `px` de espaciado suelto: todo sale de tokens de `tinywasm/css` vía
  `widget/style`. Este repo ya tiene un test que lo vigila
  (`style.TestNoInventedValues` en `tinywasm/widget`) y la revisión lo va a
  mirar.
- Sin carpetas `internal/`.
- Cada componente en su propio paquete, con su `README.md`, igual que los
  existentes.
- No dupliques `contentcard` ni `fieldset` — si les falta algo para servir
  aquí, **dilo en el PR**; se arregla en ellos, no se copia.
- El JS de `sitenav` es vanilla, sin dependencias, y debe ser idempotente:
  `RenderJS()` puede acabar incluido más de una vez en una página y no debe
  registrar el manejador dos veces.

## Verificación

- Test de contrato por componente, con la forma de
  `contentcard/card_contract_test.go`: implementa `dom.Component` y `Render()`
  es idempotente.
- `RenderCSS()` de cada uno emite en las capas correctas y sin literales
  (apóyate en cómo lo comprueban los tests existentes de `widget/style`).
- `sitenav`: un test que afirme que el markup lleva `aria-expanded` y
  `aria-controls`, y que el JS emitido no asume que el DOM ya existe (debe
  tolerar ejecutarse en `<head>` con `defer`).
- `herobanner`: un test que afirme que el CSS emitido contiene una regla
  `@media (prefers-reduced-motion: reduce)` que desactiva la animación.
- `go build ./... && go vet ./... && go test ./...` verde.

## Etapas

| # | Alcance | Aceptación |
|---|---|---|
| 1 | `statgrid` | el más simple; fija el patrón de rejilla adaptativa y el test de contrato |
| 2 | `infobar` + sus íconos | ítems inyectados; `IconSvg()` descubierto por el escaneo SSR |
| 3 | `herobanner` | rotación por CSS; `prefers-reduced-motion` probado |
| 4 | `sitenav` | menú accesible; `RenderJS()` idempotente y sin dependencias |
