//go:build !wasm

package selectsearch

func (c *SelectSearch) IconSvg() map[string]string {
	return map[string]string{
		"ss-arrow-down": `<path fill="currentColor" d="M1.5 4.5l6.5 7 6.5-7H1.5z"/>`,
	}
}
