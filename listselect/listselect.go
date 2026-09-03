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

	// changed flips on every Toggle. checked is a plain slice — reading it
	// inside a DeriveBool tracks nothing, so without this signal a row's
	// bound states would go stale the moment after selection mode opens and
	// no tap would ever repaint. Rows read Changed() first in their
	// derives; the flip re-runs them and IsChecked is re-read live.
	changed *SignalBool

	// danger arms the danger tone: while set, a checked row means "marked
	// for a destructive action" and the skin paints it red instead of blue.
	// Owned by the host (crudview sets it from its mode); a plain list that
	// never arms it keeps the single Accent selection language it always had.
	danger *SignalBool

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

// SetDanger arms or disarms the danger tone. Additive: a host that never
// calls it gets no red anywhere, whatever the selection does.
func (m *Mode) SetDanger(on bool) {
	m.ensure()
	if m.danger == nil {
		m.danger = NewBool(false)
	}
	m.danger.Set(on)
}

// Danger reports the tone signal a row binds its Invalid state to. Never a
// bool: the skin has to react.
func (m *Mode) Danger() *SignalBool {
	m.ensure()
	if m.danger == nil {
		m.danger = NewBool(false)
	}
	return m.danger
}

// Changed flips on every Toggle. Rows read it first in their state derives
// (see the changed field); without that read a tap would update nothing on
// screen.
func (m *Mode) Changed() *SignalBool {
	m.ensure()
	if m.changed == nil {
		m.changed = NewBool(false)
	}
	return m.changed
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
	m.Changed().Toggle()

	if m.OnChange != nil {
		m.OnChange(len(m.checked))
	}
}

// CheckAll marks every id in ids — the caller's CURRENT render order — and
// replaces any previous selection. It owns a fresh backing array (never
// aliases the caller's slice). Fires Changed() and OnChange with the new
// count. This is the master check's "select all" action.
func (m *Mode) CheckAll(ids []string) {
	m.ensure()
	m.checked = append([]string(nil), ids...)
	m.Changed().Toggle()
	if m.OnChange != nil {
		m.OnChange(len(m.checked))
	}
}

// Clear unmarks every row WITHOUT leaving selection mode — unlike
// SetOn(false), which also exits the mode. The master check's "deselect all".
// A no-op (no signal churn) when nothing is marked.
func (m *Mode) Clear() {
	m.ensure()
	if len(m.checked) == 0 {
		return
	}
	m.checked = nil
	m.Changed().Toggle()
	if m.OnChange != nil {
		m.OnChange(0)
	}
}

// Count reports how many rows are currently marked. The master check reads it
// to decide its tri-state (none / some / all) and to render "n / total".
func (m *Mode) Count() int { return len(m.checked) }

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
