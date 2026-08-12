//go:build wasm

package calendarslider

import (
	"testing"

	. "github.com/tinywasm/dom"
	"syscall/js"
)

func TestMain(m *testing.M) {
	app := js.Global().Get("document").Call("createElement", "div")
	app.Set("id", "app")
	js.Global().Get("document").Get("body").Call("appendChild", app)
	m.Run()
}

func domDoc() js.Value { return js.Global().Get("document") }

func query(t *testing.T, sel string) js.Value {
	t.Helper()
	el := domDoc().Call("querySelector", sel)
	if el.IsNull() || el.IsUndefined() {
		t.Fatalf("no se encontró %q", sel)
	}
	return el
}

func exists(sel string) bool {
	el := domDoc().Call("querySelector", sel)
	return !el.IsNull() && !el.IsUndefined()
}

// TestDayClickSelects cubre el uso principal: clic en un día ocupable lo
// selecciona (señal + callback) y un día no ocupable no.
func TestDayClickSelects(t *testing.T) {
	var got string
	c := &CalendarSlider{
		Start:      "2026-08",
		Holidays:   []Holiday{{Date: "2026-08-15", Name: "Asunción"}},
		Occupation: []OccupationDay{{Date: "2026-08-11", Percent: 60}},
		OnSelect:   func(date string) { got = date },
	}
	c.Init(nil)
	Render("app", c.Render())

	bookable := query(t, "#cs-d-2026-08-11")
	bookable.Call("click")
	if c.Selected.Get() != "2026-08-11" {
		t.Errorf("Selected = %q, want 2026-08-11", c.Selected.Get())
	}
	if got != "2026-08-11" {
		t.Errorf("OnSelect = %q, want 2026-08-11", got)
	}
	if query(t, "#cs-d-2026-08-11").Call("getAttribute", "data-selected").String() != "true" {
		t.Error("el día seleccionado debería llevar data-selected=true")
	}

	// Un sábado sin ocupación no es seleccionable: ni señal ni callback.
	c.Selected.Set("")
	got = ""
	query(t, "#cs-d-2026-08-01").Call("click")
	if c.Selected.Get() != "" || got != "" {
		t.Error("un día no ocupable no debería seleccionarse")
	}
}

// TestAllMonthsAlwaysInDOM cubre el cambio de arquitectura frente al viejo
// deslizador infinito: los NumMonths meses son hijos estáticos, ninguno se
// desmonta ni se recicla — el slide es puramente visual (scroll-snap).
func TestAllMonthsAlwaysInDOM(t *testing.T) {
	c := &CalendarSlider{Start: "2026-08"}
	c.Init(nil)
	Render("app", c.Render())

	for _, id := range []string{"#cs-m-2026-08", "#cs-m-2026-09", "#cs-m-2026-10"} {
		if !exists(id) {
			t.Errorf("%s debería existir en el DOM", id)
		}
	}
	if label := query(t, "#cs-m-2026-08 .calendarslider__month-name").Get("textContent").String(); label != "Agosto 2026" {
		t.Fatalf("etiqueta de agosto = %q, want Agosto 2026", label)
	}
	if label := query(t, "#cs-m-2026-10 .calendarslider__month-name").Get("textContent").String(); label != "Octubre 2026" {
		t.Fatalf("etiqueta de octubre = %q, want Octubre 2026", label)
	}
}

// TestNavLinksTargetTheNeighbor cubre la navegación: ‹ › son enlaces de
// página (<a href="#cs-m-...">) hacia el mes vecino — sin handler de clic,
// el propio navegador desplaza el scroll-snap al destino. Se verifica el
// destino a través del hash de la URL, que un enlace de ancla actualiza de
// forma nativa.
func TestNavLinksTargetTheNeighbor(t *testing.T) {
	c := &CalendarSlider{Start: "2026-08"}
	c.Init(nil)
	Render("app", c.Render())

	// Bucle infinito: agosto (primero) y octubre (último) también tienen
	// enlace hacia el otro extremo — nunca hay que recorrer los N meses en
	// orden para volver al principio.
	if !exists("#cs-m-2026-08 .calendarslider__prev") {
		t.Error("el primer mes debería tener enlace 'prev' (envuelve al último)")
	}
	if !exists("#cs-m-2026-10 .calendarslider__next") {
		t.Error("el último mes debería tener enlace 'next' (envuelve al primero)")
	}

	query(t, "#cs-m-2026-08 .calendarslider__next").Call("click")
	if hash := js.Global().Get("location").Get("hash").String(); hash != "#cs-m-2026-09" {
		t.Errorf("el next de agosto debería llevar a #cs-m-2026-09, el hash quedó en %q", hash)
	}

	query(t, "#cs-m-2026-09 .calendarslider__prev").Call("click")
	if hash := js.Global().Get("location").Get("hash").String(); hash != "#cs-m-2026-08" {
		t.Errorf("el prev de septiembre debería volver a #cs-m-2026-08, el hash quedó en %q", hash)
	}

	// El bucle: el prev de agosto (primero) envuelve directo a octubre
	// (último), sin pasar por septiembre.
	query(t, "#cs-m-2026-08 .calendarslider__prev").Call("click")
	if hash := js.Global().Get("location").Get("hash").String(); hash != "#cs-m-2026-10" {
		t.Errorf("el prev de agosto debería envolver a #cs-m-2026-10, el hash quedó en %q", hash)
	}

	// Ningún mes se desmontó durante la navegación.
	for _, id := range []string{"#cs-m-2026-08", "#cs-m-2026-09", "#cs-m-2026-10"} {
		if !exists(id) {
			t.Errorf("%s debería seguir existiendo tras navegar", id)
		}
	}
}

// TestExternalSelectedHighlights cubre el enlace inverso: el host escribe la
// señal y el DOM se pinta sin tocar el calendario.
func TestExternalSelectedHighlights(t *testing.T) {
	c := &CalendarSlider{Start: "2026-08", Occupation: []OccupationDay{{Date: "2026-08-11", Percent: 60}}}
	c.Init(nil)
	Render("app", c.Render())

	c.Selected.Set("2026-08-11")
	if query(t, "#cs-d-2026-08-11").Call("getAttribute", "data-selected").String() != "true" {
		t.Error("escribir Selected debería pintar data-selected en el DOM")
	}

	// Al navegar (el enlace ‹ › es un <a href>, no reconstruye nada), la
	// selección externa sobrevive.
	query(t, ".calendarslider__next").Call("click")
	if query(t, "#cs-d-2026-08-11").Call("getAttribute", "data-selected").String() != "true" {
		t.Error("la selección debería sobrevivir a la navegación")
	}
}
