---
PLAN: "fix(selectsearch): shape, icon turn, focus clipping and list parity"
TAG: v0.6.0
EXECUTOR: jules
REVIEWER: none
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.

# Plan — `selectsearch`: forma, giro del icono, foco recortado y paridad con el resto del chasis

> **Idioma:** este documento está en español porque lo pidió el autor.
> **El código, los comentarios de código y los nombres de símbolos van SIEMPRE en
> inglés** — `tinywasm/*` es librería pública. No traduzcas identificadores ni
> escribas comentarios en español dentro de los `.go`.

## Prerrequisito (ejecutar primero)

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

Se ejecuta `gotest` (nunca `go test`) desde la raíz del repo. No invoques
`gopush` ni `codejob`.

## ⛔ Gate: este plan depende de `tinywasm/widget`

Las etapas 2 y 6 consumen `style.Rotate` / `style.TurnHalf`, que **todavía no
existen**. Se añaden en
<https://github.com/tinywasm/widget/blob/main/docs/PLAN.md>. No empieces hasta
que ese plan esté publicado y `go.mod` apunte a la versión que los trae.
Comprobación:

```bash
grep -rn "func Rotate" $(go env GOMODCACHE)/github.com/tinywasm/widget*/style/except.go
```

---

## 1. Contexto

`components/selectsearch` es un combobox: una cabecera con un chevron que
despliega un panel con un buscador y una lista de opciones. El autor reportó
seis defectos. Todos comparten una causa de fondo: **el componente reinventó
piezas que el chasis ya tenía resueltas** en lugar de ensamblar las existentes
(`searchbar` para la cabecera, `listgap` + `targetlist` para la lista).

Eso viola la doctrina del ecosistema (`CONSTRUCTION_HARNESS.md`):

> *"Applications and server implementations do not re-implement anything; they
> **assemble pieces**."*
> *"The glue is written once, in the library that owns it."*

Archivos implicados (todos dentro de `selectsearch/`):

| Archivo | Qué contiene hoy |
|---|---|
| `css.go` | 135 líneas: el `style.Sheet` completo del componente |
| `selectsearch.go` | 244 líneas: markup, señales, `buildRows`, `selectOption` |
| `svg.go` | el icono `ss-arrow-down` |
| `selectsearch_test.go`, `css_test.go`, `selectsearch_contract_test.go` | tests |

## 2. Los seis defectos, con su causa exacta

### A. Pierde la forma al estrechar la ventana

`style.Row()` emite `flex-wrap: wrap` (verificable en
`widget/style/emit_primitives.go`, bloque `rowSel`). En `css.go`:

- `PartHeader` es `Row(SpaceNone)` **sin** `ControlBox()` ni `KeepSize()`.
- `PartIcon` tiene `ControlBox()` pero **sin** `KeepSize()`, así que encoge.

Al faltar espacio, el texto salta a una segunda línea debajo del cuadro azul y
la cabecera deja de ser una barra.

**`searchbar/css.go` ya resuelve exactamente esto** y es el patrón a copiar:

```go
Root(Row(SpaceNone), Round(RadiusMd), HideOverflow(), ControlBox(), KeepSize())
Part(PartIcon, As(Primary), MediaBox(AspectSquare), ControlBox(), KeepSize())
Part(PartInput, As(Inset), Pad(Space2), Grow(), ControlBox())
```

### B. El chevron no gira al desplegar

No hay ninguna regla de rotación, y el DOM no escribe `data-open` en ninguna
parte que el CSS pueda seleccionar. `c.isOpen` sólo alimenta el `checked` del
checkbox oculto y el `Show()` del panel.

### C. El buscador aparece con las esquinas recortadas al enfocarse

`PartDropdown` es `As(Panel)` sin `Round()` explícito → hereda
`border-radius: var(--radius-md)` por defecto (ver `Surface.defaultRadius()` en
`widget/style/surface.go`) **y** lleva `HideOverflow()`. `PartDropdown` no
tiene padding, así que `PartSearch` queda a ras del borde superior redondeado y
sus esquinas —y el anillo de foco ámbar, que se dibuja por dentro— quedan
cortadas por el clip del padre.

### D. El hover no es el del resto de la aplicación

`PartOption` usa `Interactive(style.Panel)`, cuyo hover es una mezcla gris
genérica. El lenguaje del chasis para "vas camino de seleccionar esto" es el
ámbar: `targetlist/css.go` lo declara así y lo documenta —

```go
Cue(widget.Hover, PartRow, style.As(style.AccentWash))
Cue(widget.Focus, PartRow, style.As(style.AccentWash))
Cue(widget.Press, PartRow, style.As(style.Accent))
When(widget.Selected, PartRow, style.As(style.Accent))
```

### E. El listado se ve distinto al resto de listados

`PartOptions` declara su contenedor a mano (`Stack(SpaceNone)` + `Scroll()`),
mientras que `targetlist` y `targetdate` ensamblan la pieza compartida
`components/listgap`. Y `PartOption` no tiene la forma de fila del chasis
(`Round`, `Pad(Space3)`, `ControlBox`, `Interactive(Page)`).

`listgap` existe justo para esto — su doc de paquete:

> *"It is a lego piece — the list components assemble it, they do not
> re-declare its shape or its spacing."*

### F. (adicional) IDs fijos: dos instancias en una página se pisan

`selectsearch.go` escribe literales `"ss-toggle"`, `"ss-search"`, `"ss-options"`
y `"ss-opt-"+id`. El `<label for="ss-toggle">` de una segunda instancia abre el
desplegable de la primera, y `Get("ss-search")` enfoca el buscador equivocado.
No es visible hoy porque la demo monta una sola, pero es el mismo tipo de
defecto silencioso que el resto del plan corrige.

## 3. Cambios exactos

### 3.1 `selectsearch.go` — estado `Open` en la raíz

En `Render()`, el `Div` raíz debe publicar el estado que el CSS necesita.
Sustituir el return final:

```go
	return Div().Set(ClsSsBox.AsAttr()).
		Child(toggle).
		Child(header).
		Child(Show(c.isOpen, dropdown))
```

por:

```go
	// BindState, not a class toggled by hand: data-open is the single value the
	// stylesheet selects on, so markup and CSS cannot disagree. It is what lets
	// the chevron turn be a CSS state rule instead of a second source of truth
	// in Go.
	return Div().Set(ClsSsBox.AsAttr()).
		BindState(widget.Open, c.isOpen).
		Child(toggle).
		Child(header).
		Child(Show(c.isOpen, dropdown))
```

`widget` ya está importado en el archivo. `BindState` acepta `widget.Open`
directamente (`dom.StateAttr` es una interfaz `Key()/Value()` que `widget.State`
satisface).

### 3.2 `selectsearch.go` — IDs por instancia (defecto F)

Añadir un campo no exportado y un helper. `Init` genera el prefijo una vez:

```go
	uid string // per-instance id prefix; two pickers on one page must not collide
```

En `Init`, **antes** de crear las señales:

```go
	// A page may mount more than one picker. The label's `for`, the focus
	// lookup and every option id are derived from this prefix so instance B's
	// header cannot toggle instance A's checkbox — the failure a fixed
	// "ss-toggle" guarantees the moment a second one appears.
	c.uid = fmt.Sprintf("%s-%d", string(NameSelectSearch), nextSelectSearchID())
```

> **Anti-footgun verificado:** `github.com/tinywasm/fmt` **no** tiene `Itoa` ni
> `FormatInt`. El helper sancionado para intercalar un entero es `fmt.Sprintf`,
> que ya se usa en código WASM de este mismo repo
> (`calendarslider/calendarslider.go:264`). No importes `strconv`.

Y a nivel de paquete:

```go
var selectSearchSeq int

func nextSelectSearchID() int {
	selectSearchSeq++
	return selectSearchSeq
}
```

Reemplazar los cuatro literales por derivados de `c.uid`:

| Antes | Después |
|---|---|
| `ID("ss-toggle")` | `ID(c.uid + "-toggle")` |
| `Attr("for", "ss-toggle")` | `Attr("for", c.uid+"-toggle")` |
| `Get("ss-search")` | `Get(c.uid + "-search")` |
| `ID("ss-search")` | `ID(c.uid + "-search")` |
| `ID("ss-options")` | `ID(c.uid + "-options")` |
| `Attr("aria-controls", "ss-options")` | `Attr("aria-controls", c.uid+"-options")` |
| `ID("ss-opt-"+opt.ID)` | `ID(c.uid + "-opt-" + opt.ID)` |

> **Anti-footgun:** `selectsearch.go` compila a WASM. Usa `fmt.Sprintf` de
> `github.com/tinywasm/fmt` (ya importado en el archivo), **nunca** `strconv`
> ni `strings` del stdlib.

### 3.3 `css.go` — la hoja completa

Reescribir `sheet()` con estos cambios. Mantén el resto de comentarios que ya
existen y que siguen siendo ciertos; los que cito abajo son nuevos o
reemplazan a uno obsoleto.

**Imports:** añadir `"github.com/tinywasm/components/listgap"` y
`"github.com/tinywasm/widget"`.

**(a) `PartHeader` y `PartIcon` — defecto A.** Mismo par que searchbar:

```go
		// ControlBox+KeepSize on BOTH the bar and its cap, exactly as
		// searchbar/css.go declares them: Row() carries flex-wrap: wrap, so a
		// narrow viewport wrapped the text onto its own line under the square
		// and the header stopped being a bar. KeepSize is what forbids the
		// wrap; ControlBox is what keeps the two halves the same height.
		Part(PartHeader,
			style.Row(style.SpaceNone),
			style.Round(style.RadiusMd),
			style.HideOverflow(),
			style.ControlBox(),
			style.KeepSize(),
		).
		Part(PartIcon,
			style.As(style.Primary),
			style.MediaBox(style.AspectSquare),
			style.ControlBox(),
			style.KeepSize(),
		).
```

**(b) `PartGlyph` — defecto B.** El que gira es el `<svg>`, nunca el cuadro
azul:

```go
		// The GLYPH turns, not PartIcon: rotating the cap would spin the whole
		// filled square. TurnNone is the resting rule Animate() transitions
		// from — without a base value there is no start state to move off.
		Part(PartGlyph,
			style.IconBox(style.IconSm),
			style.Rotate(style.TurnNone),
			style.Animate(style.MotionBase),
		).
```

y, junto a las demás reglas de estado:

```go
		// The chevron IS the open state. WhenWithin from the root, because
		// data-open lives on the root (see selectsearch.go's BindState) and the
		// glyph is two levels down inside the header.
		WhenWithin(widget.Open, "", PartGlyph,
			style.Rotate(style.TurnHalf),
		).
```

**(c) `PartDropdown` y `PartSearch` — defecto C:**

```go
		// Pad(Space1): the dropdown is a rounded, clipping panel
		// (As(Panel) brings RadiusMd by default; HideOverflow clips to it), so
		// a child flush against its top edge loses its corners — and with them
		// the inset focus ring, which is what made the search box look sawn
		// off. The padding is what keeps every child inside the rounded area.
		Part(PartDropdown,
			style.Stack(style.Space1),
			style.As(style.Panel),
			style.Pad(style.Space1),
			style.HideOverflow(),
		).
		// A radius of its own so the box reads as a control, not as a slab
		// filling the panel's width.
		Part(PartSearch,
			style.As(style.Inset),
			style.Pad(style.Space2),
			style.Round(style.RadiusSm),
			style.ControlBox(),
		).
```

**(d) `PartOptions` — defecto E.** Borrar el `Part(PartOptions, ...)` a mano y
ensamblar la pieza compartida. Como `sheet()` hoy es una sola cadena de
llamadas, hay que romperla igual que hace `targetlist/css.go`:

```go
func (c *SelectSearch) sheet() *style.Sheet {
	s := style.For(c).
		Root(
			style.Anchor(),
			style.Stack(style.Space1),
			style.As(style.Panel),
			style.Round(style.RadiusMd),
		)
	// The same list container targetlist and targetdate assemble — one lego
	// piece, so the dropdown's rows breathe with the same rhythm as every
	// other list in the app instead of inventing a third spacing.
	listgap.Apply(s, PartOptions)
	s.On(css.Mobile, PartOptions, listgap.MobileOpts()...)
	return s.
		Part(PartToggle, style.Hide()).
		// … el resto de la cadena, tal cual, empezando por PartDropdown
}
```

**(e) `PartOption` — defectos D y E.** La forma y los estados de una fila de
`targetlist`, literalmente:

```go
		// The row recipe targetlist/css.go declares for PartRow: same surface,
		// same box, same radius. A dropdown option and a list row are the same
		// object to the user; they must not be two different shapes.
		Part(PartOption,
			style.Row(style.Space2),
			style.KeepSize(),
			style.ControlBox(),
			style.Interactive(style.Page),
			style.Pad(style.Space3),
			style.Round(style.RadiusMd),
		).
```

y los cuatro estados, con los mismos colores y por las mismas razones que
`targetlist`:

```go
		// Amber is the chassis' one "where I am" statement — the rail's current
		// nav item and a selected list row wear it too. AccentWash on hover and
		// focus reads as "on the way to selected"; a grey mix reads as chrome
		// with no relation to what clicking does. Focus repeats Hover on
		// purpose: a keyboard user gets no :hover and must not be stranded on a
		// different colour.
		When(widget.Selected, PartOption, style.As(style.Accent)).
		Cue(widget.Hover, PartOption, style.As(style.AccentWash)).
		Cue(widget.Focus, PartOption, style.As(style.AccentWash)).
		Cue(widget.Press, PartOption, style.As(style.Accent)).
```

**(f) `PartDesc`:** dejarlo como está (`As(Inset)`, `FontSize(TextXs)`).
**No** lo conviertas en `OnEdge`: el `rowGap` de `listgap` está calibrado para
el desbordamiento del badge de `targetlist` y aquí no hay tal desbordamiento.

### 3.4 `selectsearch.go` — marcar la opción elegida

Para que `When(widget.Selected, PartOption, …)` tenga algo que seleccionar,
`buildRows` debe escribir el estado en la fila elegida. Añadir un
`*SignalString` con el id seleccionado y, en `buildRows`, sobre cada `item`:

```go
		item.BindStateFunc(widget.Selected, func() bool { return c.selectedID.Get() == o.ID })
```

Declarar `selectedID *SignalString` en el struct, inicializarlo en `Init`
(`c.selectedID = NewString("")`) y fijarlo en `selectOption`
(`c.selectedID.Set(o.ID)`), junto a `c.selectedLabel.Set(o.Label)`.

## 4. Reglas de calidad obligatorias

- **Sin strings sueltos en la lógica.** Los prefijos de id salen de `c.uid`;
  los sufijos (`-toggle`, `-search`, `-options`, `-opt-`) van en constantes de
  paquete sin exportar, no repetidos inline.
- **Sin librería estándar en código que compila a WASM.** `selectsearch.go`
  usa `github.com/tinywasm/fmt` — nunca `strconv`, `strings` ni `errors` del
  stdlib. *Anti-footgun:* `css.go` y `svg.go` llevan `//go:build !wasm` y ahí
  el stdlib sí es legítimo; no "arregles" esos imports.
- **Embebido por valor.** `SelectSearch` embebe `Element` por valor, nunca
  `*Element` (restricción de heap de TinyGo). No lo cambies.
- **Superficie mínima.** `uid`, `selectedID`, `selectSearchSeq` y
  `nextSelectSearchID` quedan sin exportar.
- **No dupliques `listgap`.** Si necesitas un valor de espaciado suyo, se
  asamblea llamando a `Apply`/`MobileOpts`; no copies sus constantes.

## 5. Tests

Molde *consumer-shaped*, como el resto del repo. Añadir a
`selectsearch/css_test.go` (o un archivo nuevo `selectsearch_visual_test.go`):

```go
func TestHeaderKeepsItsShapeWhenNarrow(t *testing.T) {
	// Row() carries flex-wrap: wrap. Without KeepSize on both the bar and its
	// cap a narrow viewport wrapped the text under the square. searchbar
	// already answers this; the header must answer it identically.
	s := (&SelectSearch{}).RenderCSS().String()
	for _, want := range []string{
		".selectsearch__header",
		"flex-shrink: 0;",
		"min-height: var(--control-height",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("expected %q in the header/icon rules, got:\n%s", want, s)
		}
	}
}

func TestChevronTurnsOnOpen(t *testing.T) {
	s := (&SelectSearch{}).RenderCSS().String()
	if !strings.Contains(s, "transform: rotate(0deg);") {
		t.Errorf("expected the chevron's resting rotation, got:\n%s", s)
	}
	if !strings.Contains(s, `.selectsearch[data-open="true"] .selectsearch__glyph`) {
		t.Errorf("expected the open state to reach the glyph, got:\n%s", s)
	}
	if !strings.Contains(s, "transform: rotate(180deg);") {
		t.Errorf("expected the open rotation, got:\n%s", s)
	}
}

func TestDropdownPadsItsClippedCorners(t *testing.T) {
	// As(Panel) brings RadiusMd and HideOverflow clips to it: a child flush
	// against the edge loses its corners AND its inset focus ring.
	s := (&SelectSearch{}).RenderCSS().String()
	i := strings.Index(s, ".selectsearch__dropdown {")
	if i < 0 {
		t.Fatal("expected a dropdown rule")
	}
	b := s[i:]
	if e := strings.Index(b, "}"); e > 0 {
		b = b[:e]
	}
	if !strings.Contains(b, "padding:") {
		t.Errorf("the clipping dropdown must pad its children off the rounded corners, got:\n%s", b)
	}
}

func TestOptionsWearTheChassisHoverAndSelection(t *testing.T) {
	// The same amber language targetlist uses. A grey hover made the dropdown
	// read as a foreign piece bolted onto the app.
	s := (&SelectSearch{}).RenderCSS().String()
	for _, want := range []string{
		"--color-accent-wash",              // hover + focus
		`.selectsearch__option[data-selected="true"]`,
		"--color-accent",                    // selected + press
	} {
		if !strings.Contains(s, want) {
			t.Errorf("expected %q, got:\n%s", want, s)
		}
	}
}

func TestOptionsReuseTheSharedListContainer(t *testing.T) {
	// listgap owns the list rhythm for the whole component set. Re-declaring
	// it here is what let the dropdown drift into a third spacing.
	s := (&SelectSearch{}).RenderCSS().String()
	if !strings.Contains(s, ".selectsearch__options") || !strings.Contains(s, "overflow-y: auto;") {
		t.Errorf("expected the shared list container on the options, got:\n%s", s)
	}
}
```

Y para el defecto F, en `selectsearch_test.go`:

```go
func TestTwoPickersDoNotShareIDs(t *testing.T) {
	a, b := &SelectSearch{}, &SelectSearch{}
	a.Init(nil)
	b.Init(nil)
	ha, hb := a.Render().String(), b.Render().String()
	if strings.Contains(ha, "ss-toggle") || strings.Contains(hb, "ss-toggle") {
		t.Error("ids must be per-instance, not the fixed ss-* literals")
	}
	if ha == hb {
		t.Error("two pickers rendered identical markup — their ids collide")
	}
}
```

## 6. Criterios de aceptación (verificables)

```bash
gotest                                                     # vet ✅ race ✅ tests ✅ wasm ✅
grep -rn "ss-toggle\|ss-search\|ss-options" selectsearch/*.go | grep -v _test   # → vacío
grep -rn "listgap.Apply" selectsearch/css.go               # 1 resultado
grep -rn "Interactive(style.Panel)" selectsearch/css.go    # → vacío (era el hover gris)
grep -rn "AccentWash" selectsearch/css.go                  # hover + focus
grep -rn "style.Rotate" selectsearch/css.go                # 2 resultados (base + estado)
grep -rn "BindState(widget.Open" selectsearch/selectsearch.go  # 1 resultado
```

Además: `components/conformance_test.go` debe seguir verde sin tocarlo — si
falla, es que `RenderCSS()` emite una clase nueva no declarada, o que
`Validate()` rechaza la hoja.

## 7. Verificación manual (la hace el desarrollador, no el agente)

Sólo el resultado visual se comprueba a mano; todo lo funcional está cubierto
arriba. Lista para el revisor humano:

1. Estrechar la ventana hasta móvil: la cabecera sigue siendo **una barra**.
2. Abrir: el chevron gira 180° **con transición**; cerrar: vuelve.
3. Enfocar el buscador: anillo ámbar **completo**, sin esquinas cortadas.
4. Pasar el ratón por una opción: ámbar claro, **igual** que una fila de la
   lista de fichas.
5. La opción elegida queda en ámbar sólido.

## 8. Fuera de alcance (NO hacer)

- No toques `searchbar`, `targetlist`, `targetdate` ni `listgap`: son la
  referencia, no el objetivo. Si crees que a `listgap` le falta algo, **para y
  repórtalo** (regla del harness: un consumidor no recrea lo que falta arriba).
- No añadas un `internal/` ni un wrapper sobre `listgap`.
- El vaciado del formulario al cambiar de paciente **no** es de este plan: va
  en <https://github.com/tinywasm/layout/blob/main/docs/PLAN.md>.

## 9. Etapas

| # | Etapa | Defecto | Archivos | Cierra cuando |
|---|---|---|---|---|
| 1 | Forma de la cabecera | A | `css.go` | `TestHeaderKeepsItsShapeWhenNarrow` pasa |
| 2 | Giro del chevron | B | `css.go`, `selectsearch.go` | `TestChevronTurnsOnOpen` pasa |
| 3 | Esquinas del buscador | C | `css.go` | `TestDropdownPadsItsClippedCorners` pasa |
| 4 | Contenedor de lista compartido | E | `css.go` | `TestOptionsReuseTheSharedListContainer` pasa |
| 5 | Hover/selección del chasis | D | `css.go`, `selectsearch.go` | `TestOptionsWearTheChassisHoverAndSelection` pasa |
| 6 | IDs por instancia | F | `selectsearch.go` | `TestTwoPickersDoNotShareIDs` pasa; grep de `ss-toggle` vacío |
| 7 | Suite completa | — | — | `gotest` verde, conformance intacto |

La etapa 4 es **gate** de la 5 (la 5 estiliza filas dentro del contenedor que
la 4 ensambla). Las etapas 1, 2, 3 y 6 son paralelas entre sí.
