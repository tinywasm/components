//go:build !wasm

package card

func (c *Card) RenderCSS() string {
	return `.card { 
		padding: 1rem; 
		border: 1px solid #ccc; 
		border-radius: 8px;
		margin: 1rem 0;
	}`
}
