package herobanner

import (
	. "github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
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
			img := NewElement("img").
				Attr("src", imgSrc).
				Attr("alt", "").
				Attr("loading", "eager")
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
