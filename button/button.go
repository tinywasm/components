package button

import "github.com/tinywasm/fmt"

type Button struct {
	Label string
}

func (b *Button) HandlerName() string {
	return "btn-" + b.Label
}

func (b *Button) RenderHTML() string {
	return fmt.Html(`<button id="%s" class="btn">%s</button>`, b.HandlerName(), b.Label).String()
}
