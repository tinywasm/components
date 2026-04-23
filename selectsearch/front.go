//go:build wasm

package selectsearch

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
)

func (c *SelectSearch) OnMount() {
	id := c.GetID()

	// Wire search filtering
	if searchEl, ok := dom.Get(id + "-search"); ok {
		searchEl.On("input", func(e dom.Event) {
			term := e.TargetValue()
			c.filterTerm = term

			// Check if we need to call OnSearch
			allHidden := true
			lowerTerm := fmt.Convert(term).ToLower().String()
			for _, opt := range c.Options {
				if c.matches(opt, lowerTerm) {
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

			c.Update()

			// After update, we need to restore focus to search input and put cursor at the end
			if newSearchEl, ok := dom.Get(id + "-search"); ok {
				newSearchEl.Focus()
			}
		})
	}

	// Wire option click — close dropdown and call OnSelect
	if optionsEl, ok := dom.Get(id + "-options"); ok {
		optionsEl.On("click", func(e dom.Event) {
			target, ok := dom.Get(e.TargetID())
			if !ok {
				return
			}
			optID := target.GetAttr("data-id")
			optDesc := target.GetAttr("data-description")

			if optID != "" {
				// Find the option in c.Options to get the label
				for _, opt := range c.Options {
					if opt.ID == optID {
						c.selectedLabel = opt.Label
						break
					}
				}

				if c.OnSelect != nil {
					c.OnSelect(optID, optDesc)
				}

				c.filterTerm = ""
				c.Update()
			}
		})
	}
}
