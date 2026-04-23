//go:build wasm

package selectsearch

import (
	"github.com/tinywasm/dom"
	"strings"
)

func (c *SelectSearch) OnMount() {
	id := c.GetID()

	// Wire search filtering
	if searchEl, ok := dom.Get(id + "-search"); ok {
		searchEl.On("input", func(e dom.Event) {
			term := e.TargetValue()
			c.FilterTerm = term

			// Check if we need to call OnSearch
			allHidden := true
			lowerTerm := strings.ToLower(term)
			for _, opt := range c.Options {
				if strings.Contains(strings.ToLower(opt.Label), lowerTerm) ||
				   strings.Contains(strings.ToLower(opt.Description), lowerTerm) {
					allHidden = false
					break
				}
			}

			if allHidden && c.OnSearch != nil {
				newOptions := c.OnSearch(term)
				if len(newOptions) > 0 {
					c.Options = append(c.Options, newOptions...)
				}
			}

			dom.Update(c)

			// After update, we need to restore focus to search input and put cursor at the end
			if newSearchEl, ok := dom.Get(id + "-search"); ok {
				newSearchEl.Focus()
			}
		})
	}
}
