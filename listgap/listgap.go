//go:build !wasm

// Package listgap holds the row-to-row gap for any target* list component
// (targetlist, targethour, …), reused for that same list's own lateral
// padding so a row sits equidistant from its neighbors and from the list's
// own edge in every direction, at both breakpoints.
//
// A leaf package on purpose, sibling to targetlist/targethour rather than a
// const in the components root: the root package's own conformance_test.go
// already imports targetlist, so a css.go importing the root back would be
// components → targetlist → components — a cycle in the test build. Nothing
// imports listgap but the list components (and, if a conformance check ever
// needs it, the root test file), so there is no path back.
package listgap

import "github.com/tinywasm/widget/style"

const (
	Gap       = style.Space1
	GapMobile = style.Space2
)
