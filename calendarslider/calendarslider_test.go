//go:build !wasm

package calendarslider

import (
	"strings"
	"testing"

	. "github.com/tinywasm/fmt"
	"github.com/tinywasm/widget"
)

// La aritmética pura de fechas (bisiestos, días por mes, día de la semana,
// suma de meses, claves YYYY-MM) vive y se prueba en github.com/tinywasm/date
// — reutilizable fuera de este componente. Ver TestBuildMonthPadsToSixWeeks
// más abajo para la única cobertura que sigue siendo de calendarslider: cómo
// usa esa aritmética para armar la grilla.

// TestBuildMonthAgosto2026 fija la geometría del mes de agosto 2026: comienza
// sábado (5 huecos iniciales), tiene 31 días y las celdas finales quedan en la
// última semana parcial (1 día). El mes lleva: fila de días de la semana +
// 6 semanas + etiqueta del mes.
func TestBuildMonthAgosto2026(t *testing.T) {
	c := &CalendarSlider{today: "2026-08-11"}
	m := c.buildMonth(2026, 8, "2026-07", "2026-09")
	if m == nil {
		t.Fatal("buildMonth returned nil")
	}
	children := m.Children()
	if len(children) != 10 {
		t.Fatalf("agosto 2026 debería tener fila de días + 6 semanas + etiqueta + prev + next = 10 hijos, tiene %d", len(children))
	}
	if children[0].String() != c.buildWeekdayRow().String() {
		t.Error("el primer hijo del mes debería ser la fila de días de la semana")
	}
	if !Contains(children[7].String(), "Agosto 2026") {
		t.Errorf("la etiqueta del mes debería decir 'Agosto 2026', dice %s", children[7].String())
	}
	if !Contains(children[8].String(), "href='#cs-m-2026-07'") {
		t.Errorf("el enlace anterior debería apuntar a julio, dice %s", children[8].String())
	}
	if !Contains(children[9].String(), "href='#cs-m-2026-09'") {
		t.Errorf("el enlace siguiente debería apuntar a septiembre, dice %s", children[9].String())
	}
}

// TestBuildMonthAlwaysHasBothLinks cubre buildMonth en sí: siempre arma los
// dos enlaces con lo que se le pasa, sin importar si el llamador (Render)
// los está usando como vecino real o como vuelta del bucle. Quien decide el
// bucle infinito es Render, no buildMonth — ver TestRenderWrapsAround.
func TestBuildMonthAlwaysHasBothLinks(t *testing.T) {
	c := &CalendarSlider{today: "2026-08-11"}
	m := c.buildMonth(2026, 8, "2026-07", "2026-09").String()
	if !Contains(m, "Mes anterior") || !Contains(m, "Mes siguiente") {
		t.Error("buildMonth debería incluir siempre ambos enlaces")
	}
}

// TestRenderWrapsAround cubre el bucle infinito: el ‹ del primer mes de la
// tira apunta al último y el › del último apunta al primero — igual que el
// deslizador original, para volver al inicio sin recorrer los N meses.
func TestRenderWrapsAround(t *testing.T) {
	c := &CalendarSlider{Start: "2026-08", NumMonths: 3}
	c.Init(nil)
	c.today = "2026-08-11"
	htmlOut := c.Render().String()

	if got := strings.Count(htmlOut, "calendarslider__prev"); got != 3 {
		t.Errorf("cada uno de los 3 meses debería tener un enlace 'prev', hay %d", got)
	}
	if got := strings.Count(htmlOut, "calendarslider__next"); got != 3 {
		t.Errorf("cada uno de los 3 meses debería tener un enlace 'next', hay %d", got)
	}

	// Agosto (primero) enlaza hacia atrás con octubre (último) — el bucle.
	augMonth := c.buildMonth(2026, 8, "2026-10", "2026-09").String()
	if !Contains(augMonth, "href='#cs-m-2026-10'") {
		t.Errorf("el 'prev' de agosto (primero) debería envolver a octubre (último), dice %s", augMonth)
	}
	// Octubre (último) enlaza hacia adelante con agosto (primero) — el bucle.
	octMonth := c.buildMonth(2026, 10, "2026-09", "2026-08").String()
	if !Contains(octMonth, "href='#cs-m-2026-08'") {
		t.Errorf("el 'next' de octubre (último) debería envolver a agosto (primero), dice %s", octMonth)
	}
}

// weeksPerMonth es siempre 6, ocupe o no el mes real sus 6 filas: sin este
// piso fijo, un mes de menos filas deja una tarjeta más baja y la etiqueta
// del mes (y el ‹ › que se ancla a la tarjeta) saltan verticalmente al
// deslizar entre meses. Febrero de 2021 empieza lunes y tiene 28 días — 4
// filas reales — el caso más corto posible.
func TestBuildMonthPadsToSixWeeks(t *testing.T) {
	c := &CalendarSlider{today: "2021-02-11"}
	m := c.buildMonth(2021, 2, "2021-01", "2021-03")
	children := m.Children()

	// fila de días + 6 semanas + etiqueta + prev + next, igual que un mes de
	// 6 filas reales (TestBuildMonthAgosto2026) — el total nunca cambia.
	if len(children) != 10 {
		t.Fatalf("febrero 2021 debería tener 10 hijos (igual que cualquier mes), tiene %d", len(children))
	}

	// Las 2 últimas semanas (índices 5 y 6 tras la fila de días de la
	// semana) son relleno: aria-hidden y sin ningún día real.
	for _, idx := range []int{5, 6} {
		week := children[idx]
		if !Contains(week.String(), "aria-hidden") {
			t.Errorf("la semana de relleno %d debería llevar aria-hidden", idx)
		}
		if Contains(week.String(), "data-date") {
			t.Errorf("la semana de relleno %d no debería contener ningún día real", idx)
		}
	}

	// La etiqueta sigue siendo el hijo 7 (después de las 6 semanas), igual
	// que en un mes con 6 filas reales — la posición no se mueve.
	if !Contains(children[7].String(), "Febrero 2021") {
		t.Errorf("la etiqueta debería seguir en la posición 7, dice %s", children[7].String())
	}
}

// TestBuildMonthCells valida el contenido de las celdas: huecos iniciales,
// marcadores de hoy/feriado/domingo, ocupación y seleccionabilidad.
func TestBuildMonthCells(t *testing.T) {
	c := &CalendarSlider{
		today:      "2026-08-11",
		Holidays:   []Holiday{{Date: "2026-08-15", Name: "Asunción de la Virgen"}},
		Occupation: []OccupationDay{{Date: "2026-08-11", Percent: 60}, {Date: "2026-08-02", Percent: 30}},
	}
	// 1 de agosto de 2026 = sábado → 5 huecos, el día 1 cae en la columna 6.
	firstWeek := c.buildMonth(2026, 8, "2026-07", "2026-09").Children()[1]
	cells := firstWeek.Children()
	if len(cells) != 7 {
		t.Fatalf("primera semana debería tener 7 celdas, tiene %d", len(cells))
	}
	for i := 0; i < 5; i++ {
		if !Contains(cells[i].String(), "aria-hidden") {
			t.Errorf("celda hueco %d debería ser una celda vacía aria-hidden", i)
		}
	}

	day1 := cells[5]
	if !Contains(day1.String(), "data-date='2026-08-01'") {
		t.Error("día 1 debería caer en la columna 6")
	}
	if !Contains(day1.String(), "day-off") || Contains(day1.String(), "day-red") {
		t.Error("1 de agosto (sábado) sin ocupación ni feriado debería ser un día normal (day-off), el rojo es solo domingo/feriado")
	}
	if Contains(day1.String(), "day-selectable") || Contains(day1.String(), "data-use") {
		t.Error("un sábado sin ocupación no debería ser seleccionable")
	}

	// Domingo 2 con ocupación: la ocupación gana sobre el domingo (fiel al
	// original) → seleccionable, no rojo.
	day2 := cells[6]
	if !Contains(day2.String(), "day-selectable") {
		t.Error("domingo 2 con ocupación debería ser seleccionable")
	}
	if Contains(day2.String(), "day-red") {
		t.Error("la ocupación debería ganarle al domingo")
	}
	if !Contains(day2.String(), "data-use='30'") || !Contains(day2.String(), "--meter-fill:30%") {
		t.Error("la barra de ocupación debería llevar el porcentaje")
	}

	// Martes 11: hoy + ocupación → marcador de hoy y seleccionable.
	day11 := c.buildDay(2026, 8, 11)
	if !Contains(day11.String(), "day-today") {
		t.Error("el día de hoy debería marcarse con day-today")
	}
	if !Contains(day11.String(), "day-selectable") {
		t.Error("el día de hoy con ocupación debería ser seleccionable")
	}
	if !Contains(day11.String(), "title='Hoy'") {
		t.Error("el día de hoy debería llevar title 'Hoy'")
	}

	// Sábado 15: feriado → rojo con nombre en el title, no seleccionable.
	day15 := c.buildDay(2026, 8, 15)
	if !Contains(day15.String(), "day-red") {
		t.Error("el feriado debería marcarse rojo")
	}
	if !Contains(day15.String(), "Asunción de la Virgen") {
		t.Error("el feriado debería llevar el nombre en el title")
	}
	if Contains(day15.String(), "day-selectable") {
		t.Error("un feriado no debería ser seleccionable")
	}

	// Viernes 14: día hábil sin ocupación → inactivo, no seleccionable.
	day14 := c.buildDay(2026, 8, 14)
	if !Contains(day14.String(), "day-off") {
		t.Error("día hábil sin ocupación debería marcarse como inactivo")
	}
	if Contains(day14.String(), "day-selectable") {
		t.Error("día sin ocupación no debería ser seleccionable")
	}
}

func TestClampOccupation(t *testing.T) {
	c := &CalendarSlider{Occupation: []OccupationDay{{Date: "2026-08-11", Percent: 150}, {Date: "2026-08-12", Percent: -5}}}
	if !Contains(c.buildDay(2026, 8, 11).String(), "data-use='100'") {
		t.Error("ocupación sobre 100 debería recortarse a 100")
	}
	if !Contains(c.buildDay(2026, 8, 12).String(), "data-use='0'") {
		t.Error("ocupación negativa debería recortarse a 0")
	}
}

func TestRenderStructure(t *testing.T) {
	c := &CalendarSlider{
		Start:      "2026-08",
		NumMonths:  3,
		Holidays:   []Holiday{{Date: "2026-08-15", Name: "Asunción"}},
		Occupation: []OccupationDay{{Date: "2026-08-11", Percent: 60}},
	}
	c.Init(nil)
	c.today = "2026-08-11"

	htmlOut := c.Render().String()

	// La tira [agosto, septiembre, octubre] es la esperada — Start es el
	// primer mes, no el centro; un elemento estático por mes, sin señal ni
	// reconstrucción.
	for _, id := range []string{"cs-m-2026-08", "cs-m-2026-09", "cs-m-2026-10"} {
		if !strings.Contains(htmlOut, "id='"+id+"'") {
			t.Errorf("la tira debería incluir el mes %q, no aparece en:\n%s", id, htmlOut)
		}
	}

	// Cada mes lleva sus dos enlaces (el bucle infinito se cubre en
	// TestRenderWrapsAround); agosto (Start) enlaza a septiembre como
	// siguiente.
	if !Contains(htmlOut, "href='#cs-m-2026-09'") {
		t.Error("agosto (Start) debería enlazar a septiembre como mes siguiente")
	}

	// Las flechas van superpuestas (EdgeStrip) en la hoja de estilos servida,
	// no como estilo inline del HTML renderizado.
	cssOut := c.RenderCSS().String()
	if !strings.Contains(cssOut, "position: absolute;") {
		t.Error("las flechas deberían ir superpuestas (position:absolute) en la hoja de estilos")
	}

	if !Contains(htmlOut, "Lun") || !Contains(htmlOut, "Dom") {
		t.Error("la fila de días de la semana debería estar dentro de cada mes")
	}
	if !strings.Contains(htmlOut, "calendarslider__week-row") {
		t.Error("las semanas deberían llevar la clase de grilla de 7 columnas")
	}
	if !Contains(htmlOut, "Agosto 2026") {
		t.Error("cada mes debería llevar su etiqueta (mes + año)")
	}

	// Selección reactiva: la señal ya está seteada antes del render y el
	// binding de estado se serializa en la celda del día.
	c.Selected.Set("2026-08-11")
	htmlOut = c.Render().String()
	if !strings.Contains(htmlOut, "data-selected='true'") {
		t.Error("las celdas deberían llevar el estado data-selected")
	}
}

// TestPairMarkupAndStylesheet: toda clase de la hoja CSS existe en el HTML
// renderizado. Los meses ahora son hijos estáticos (Child, no BindChildren),
// así que Render().String() ya los serializa completos.
func TestPairMarkupAndStylesheet(t *testing.T) {
	c := &CalendarSlider{Start: "2026-08", Holidays: []Holiday{{Date: "2026-08-15", Name: "x"}}, Occupation: []OccupationDay{{Date: "2026-08-11", Percent: 10}}}
	c.Init(nil)
	c.today = "2026-08-11"
	htmlOut := c.Render().String()
	cssOut := c.RenderCSS().String()

	extractClasses := func(hay, prefix string) map[string]bool {
		out := make(map[string]bool)
		rest := hay
		for {
			idx := strings.Index(rest, prefix)
			if idx < 0 {
				break
			}
			end := idx
			for end < len(rest) && rest[end] != ' ' && rest[end] != '{' && rest[end] != ',' && rest[end] != '}' && rest[end] != '\n' && rest[end] != '\r' && rest[end] != '\t' && rest[end] != ':' && rest[end] != '[' {
				end++
			}
			out[rest[idx:end]] = true
			rest = rest[end:]
		}
		return out
	}

	for cls := range extractClasses(cssOut, "calendarslider__") {
		if !strings.Contains(htmlOut, cls) {
			t.Errorf("clase CSS %q no existe en el HTML renderizado", cls)
		}
	}
}

// TestNumMonthsClampsToMax cubre el tope de la tira: un NumMonths por
// encima de maxMonths (12) se recorta, igual que el límite del calendario
// original — evita una tira de scroll-snap sin control.
func TestNumMonthsClampsToMax(t *testing.T) {
	c := &CalendarSlider{Start: "2026-08", NumMonths: 20}
	c.Init(nil)
	htmlOut := c.Render().String()

	if got := strings.Count(htmlOut, "id='cs-m-"); got != maxMonths {
		t.Fatalf("NumMonths=20 debería recortarse a %d meses, la tira tiene %d", maxMonths, got)
	}
	if !Contains(htmlOut, "id='cs-m-2026-08'") {
		t.Error("el primer mes de una tira de 12 empezando en agosto 2026 debería ser agosto 2026 (Start)")
	}
	if !Contains(htmlOut, "id='cs-m-2027-07'") {
		t.Error("el último mes de una tira de 12 empezando en agosto 2026 debería ser julio 2027")
	}
}

// TestNumMonthsDefaultsToThree cubre el default documentado (0 = 3).
func TestNumMonthsDefaultsToThree(t *testing.T) {
	c := &CalendarSlider{Start: "2026-08"}
	c.Init(nil)
	htmlOut := c.Render().String()
	if got := strings.Count(htmlOut, "id='cs-m-"); got != 3 {
		t.Errorf("NumMonths sin especificar debería dar 3 meses, la tira tiene %d", got)
	}
}

func TestCalendarSlider_SatisfiesFilterable(t *testing.T) {
	var c widget.Filterable = &CalendarSlider{}
	got := ""
	c.OnFilterChange(func(term string) { got = term })
	if got != "" {
		t.Fatalf("sink must not fire on registration, got %q", got)
	}
}
