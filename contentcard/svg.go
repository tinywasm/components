//go:build !wasm

package contentcard

func (c *ContentCard) IconSvg() map[string]string {
	return nil
}
