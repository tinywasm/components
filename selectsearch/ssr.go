//go:build !wasm

package selectsearch

import _ "embed"

//go:embed selectsearch.css
var css string

func (c *SelectSearch) RenderCSS() string {
	return css
}

func (c *SelectSearch) IconSvg() map[string]string {
	return map[string]string{
		// Arrow down icon for the dropdown header
		// viewBox 0 0 16 16 (default)
		"ss-arrow-down": `<path fill-rule="evenodd" d="M1.5 4.5l6.5 7 6.5-7H1.5z"/>`,
	}
}
