//go:build !wasm

package form

func (f *Form) RenderCSS() string {
	return ""
}
