// Package listselect owns one concern: the multi-selection mode a record list
// enters when its host is about to act on several rows at once. It is a lego
// piece — targetlist and targetdate assemble it, they do not re-declare it.
//
// Sibling of listgap, and for the same reason: two lists that must stay
// visually and behaviourally interchangeable (crudview swaps one for the
// other) cannot each own a private copy of the rule.
package listselect

import (
	. "github.com/tinywasm/dom"
)

// Mode is the selection state of one list. The zero value is a usable list in
// normal mode — a list is a list until its host says otherwise.
type Mode struct {
	on *SignalBool

	// checked is a SLICE, never a map. A map pulls TinyGo's hashing and
	// runtime machinery into the WASM binary, and this type compiles to the
	// browser. A record list holds tens of rows, so a linear scan costs
	// microseconds — the map was never buying anything here.
	//
	// []string and not []fmt.KeyValue: this is a SET, so there is no value to
	// carry. fmt.KeyValue is the map-free shape for an actual key→value pair;
	// inventing a Value of "true" would be a stringly-typed bool.
	checked []string

	// OnChange fires after every toggle with the current count, so a host can
	// label its commit button ("🗑 3") and disable it at zero.
	OnChange func(n int)
}

func (m *Mode) ensure() {
	if m.on == nil {
		m.on = NewBool(false)
	}
}

// On reports the signal a component binds its root state to, so the stylesheet
// can reveal the checks. Never a bool: the skin has to react.
func (m *Mode) On() *SignalBool {
	m.ensure()
	return m.on
}

// SetOn enters or leaves selection mode. Leaving ALWAYS clears the marks:
// a mode the user cancelled must not leave a hidden selection behind for the
// next entry to inherit silently.
func (m *Mode) SetOn(on bool) {
	m.ensure()
	if !on && len(m.checked) > 0 {
		m.checked = nil
		if m.OnChange != nil {
			m.OnChange(0)
		}
	}
	m.on.Set(on)
}

// Toggle marks or unmarks one id and fires OnChange.
func (m *Mode) Toggle(id string) {
	m.ensure()
	found := -1
	for i, c := range m.checked {
		if c == id {
			found = i
			break
		}
	}

	if found >= 0 {
		// remove
		m.checked = append(m.checked[:found], m.checked[found+1:]...)
	} else {
		// add
		m.checked = append(m.checked, id)
	}

	if m.OnChange != nil {
		m.OnChange(len(m.checked))
	}
}

// IsChecked answers for one id — what a row binds its check state to.
func (m *Mode) IsChecked(id string) bool {
	for _, c := range m.checked {
		if c == id {
			return true
		}
	}
	return false
}

// CheckedIDs returns the marked ids in the order given by ids, which the
// caller passes as its current render order.
//
// Ordering is NOT optional and NOT the caller's problem to remember: checked
// accumulates in TAP order, and a host building a confirmation message from
// tap order would list rows in an order that matches nothing on screen.
// Taking the render order as a parameter is what makes the wrong version
// unwritable — there is no accessor that returns the raw slice.
func (m *Mode) CheckedIDs(ids []string) []string {
	if len(m.checked) == 0 {
		return nil
	}
	var res []string
	for _, id := range ids {
		if m.IsChecked(id) {
			res = append(res, id)
		}
	}
	return res
}
