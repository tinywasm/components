package herobanner

import (
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/image"
	"github.com/tinywasm/widget"
)

// NameHeroBanner is the widget name for herobanner.
const NameHeroBanner = widget.Name("herobanner")

const (
	PartMedia    = widget.Part("media")
	PartContent  = widget.Part("content")
	PartTitle    = widget.Part("title")
	PartSubtitle = widget.Part("subtitle")
	PartActions  = widget.Part("actions")
)

var (
	clsHero         = NameHeroBanner.Root()
	clsHeroMedia    = NameHeroBanner.Class(PartMedia)
	clsHeroContent  = NameHeroBanner.Class(PartContent)
	clsHeroTitle    = NameHeroBanner.Class(PartTitle)
	clsHeroSubtitle = NameHeroBanner.Class(PartSubtitle)
	clsHeroActions  = NameHeroBanner.Class(PartActions)
)

// loadingFor decides the loading hint per slide. Only the first layer is
// painted when the page opens — the rest are revealed by the rotation
// keyframes seconds later — so marking every one eager makes the browser
// fetch six full-size photographs before first paint, competing with the one
// slide the visitor can actually see. The first stays eager because it IS the
// largest contentful paint; the others can arrive while the first is on
// screen.
func loadingFor(layer int) string {
	if layer == 0 {
		return "eager"
	}
	return "lazy"
}

// autoRotateLayers mirrors style.AutoRotateLayers (widget/style is !wasm-only
// and Render() must compile for wasm, so it cannot import that package) —
// keep this in sync if the widget-side constant ever changes.
const autoRotateLayers = 6

// HeroBanner displays a prominent header banner with text, background image(s), and actions.
type HeroBanner struct {
	Element
	Title    string
	Subtitle string
	Images   []string
	Actions  []Component
}

func (h *HeroBanner) WidgetName() widget.Name { return NameHeroBanner }
func (h *HeroBanner) WidgetKind() widget.Kind { return widget.Region }

func (h *HeroBanner) Render() *Element {
	banner := Section().Set(clsHero.AsAttr())

	if len(h.Images) > 0 {
		mediaLayer := Div().Set(clsHeroMedia.AsAttr())
		for i := 0; i < autoRotateLayers; i++ {
			imgSrc := h.Images[i%len(h.Images)]
			img := image.Responsive(imgSrc, "").
				Attr("loading", loadingFor(i)).
				AsElement()
			mediaLayer.Child(img)
		}
		banner.Child(mediaLayer)
	}

	content := Div().Set(clsHeroContent.AsAttr())

	if h.Title != "" {
		content.Child(H1().Set(clsHeroTitle.AsAttr()).Text(h.Title))
	}

	if h.Subtitle != "" {
		content.Child(P().Set(clsHeroSubtitle.AsAttr()).Text(h.Subtitle))
	}

	if len(h.Actions) > 0 {
		actionsRow := Div().Set(clsHeroActions.AsAttr())
		for _, act := range h.Actions {
			if act != nil {
				actionsRow.Child(act)
			}
		}
		content.Child(actionsRow)
	}

	banner.Child(content)
	return banner
}
