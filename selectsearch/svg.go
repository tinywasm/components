package selectsearch

import "github.com/tinywasm/svg"

var (
	arrowDown = svg.Define("ss-arrow-down", "0 0 16 16",
		svg.Path("M1.5 4.5l6.5 7 6.5-7H1.5z"),
	)
)

func (c *SelectSearch) IconSvg() *svg.Sprite {
	return svg.NewSprite(arrowDown)
}
