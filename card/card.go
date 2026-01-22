package card

import "github.com/tinywasm/fmt"

type Card struct {
	Title   string
	Content string
}

func (c *Card) HandlerName() string {
	return "card-" + c.Title
}

func (c *Card) RenderHTML() string {
	return fmt.Html(`<div id="%s" class="card">
		<h3>%s</h3>
		<p>%s</p>
	</div>`, c.HandlerName(), c.Title, c.Content).String()
}
