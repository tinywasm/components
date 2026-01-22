//go:build !wasm

package button

func (b *Button) RenderCSS() string {
	return `.btn { 
		padding: 0.5rem 1rem; 
		cursor: pointer;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
	}`
}
