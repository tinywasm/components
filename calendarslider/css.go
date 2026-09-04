//go:build !wasm

package calendarslider

import (
	"github.com/tinywasm/css"
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

// RenderCSS describe el aspecto del calendario con los tokens del tema. La
// puesta en escena es fiel al calendario original: un mes a la vez, centrado
// en la página (Center), con las secciones vecinas fuera de vista en una
// tira de scroll-snap (ScrollRow) — cada mes ocupa el 100% de la tira, así
// que solo uno es visible; el hoy se rellena con el color primario y las
// flechas ‹ › son tiras superpuestas en los bordes del mes (EdgeStrip), no
// ancladas a una esquina, para que cubran toda su altura como en el
// original. Cada mes lleva sus propias flechas — apuntan al vecino, así que
// solo hace falta ancla (Anchor) en la sección misma, no en la raíz.
//
// PartMonth lleva PadInline: las flechas se posicionan contra su propio
// padding-box (position:absolute no ve el padding de un ancestro salvo el
// propio), así que ese padding es el gutter que evita que se superpongan
// con la columna de días más externa.
//
// El mismo diseño de un mes a la vez ya funciona en mobile sin reglas
// aparte; solo el tamaño de la celda del día crece para un objetivo táctil
// más cómodo.
func (c *CalendarSlider) RenderCSS() *css.Stylesheet {
	return style.For(c).
		Root(
			style.Center(style.Compact),
		).
		Part(PartWeekRow,
			style.FixedGrid(7, style.SpaceNone),
		).
		Part(PartStrip,
			style.ScrollRow(style.SpaceNone),
		).
		Part(PartMonth,
			style.As(style.Panel),
			style.Pad(style.Space2),
			style.Anchor(),
			style.PadInline(style.Space6),
			style.Width(style.Full),
			style.Stack(style.Space1),
		).
		Part(PartWeekday,
			style.CenterContent(),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
		).
		Part(PartMonthName,
			style.CenterContent(),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
			style.Pad(style.Space1),
		).
		Part(PartDayStack,
			style.Stack(style.SpaceNone),
		).
		Part(PartDay,
			style.CenterContent(),
			style.IconBox(style.IconMd),
			style.CenterSelf(),
		).
		Part(PartDayNum,
			style.Grow(),
			style.CenterContent(),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
		).
		Part(PartDayUse,
			style.As(style.Inset),
			style.Round(style.RadiusFull),
			style.Meter(style.Space1),
			style.CenterSelf(),
		).
		Part(PartDaySelectable,
			style.Interactive(style.Inset),
		).
		Part(PartDayRed,
			style.Glyph(style.Danger),
		).
		Part(PartDayToday,
			style.As(style.Primary),
			style.Round(style.RadiusSm),
		).
		Part(PartPrev,
			style.As(style.Bare),
			style.CenterContent(),
			style.FontSize(style.Text2xl),
			style.FontWeight(style.WeightBold),
			style.Pad(style.Space1),
			style.Interactive(style.Subtle),
			style.Round(style.RadiusSm),
			style.EdgeStrip(style.Parent, style.SideStart),
		).
		Part(PartNext,
			style.As(style.Bare),
			style.CenterContent(),
			style.FontSize(style.Text2xl),
			style.FontWeight(style.WeightBold),
			style.Pad(style.Space1),
			style.Interactive(style.Subtle),
			style.Round(style.RadiusSm),
			style.EdgeStrip(style.Parent, style.SideEnd),
		).
		When(widget.Selected, PartDay,
			style.As(style.Accent),
			style.Raise(style.Raised),
		).
		On(css.Mobile, PartDay,
			style.IconBox(style.IconLg),
		).
		Stylesheet()
}
