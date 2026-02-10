package all

import (
	"github.com/tinywasm/components/button"
	"github.com/tinywasm/components/card"
	"github.com/tinywasm/components/form"
	"github.com/tinywasm/components/input"
	"github.com/tinywasm/components/modal"
	"github.com/tinywasm/components/nav"
	"github.com/tinywasm/components/table"
	"github.com/tinywasm/dom"
)

// Components returns a list of all component factories or instances.
func Components() []dom.Component {
	return []dom.Component{
		&button.Button{},
		&card.Card{},
		&form.Form{},
		&input.Input{},
		&modal.Modal{},
		&nav.Nav{},
		&table.Table{},
	}
}
