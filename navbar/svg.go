//go:build !wasm

package navbar

func (n *NavBar) IconSvg() map[string]string {
	return nil
}
