//go:build !wasm

package targetlist

import (
	"regexp"
	"strings"
	"testing"
)

// tagRegex matches an opening or closing tag. Every element this component
// renders has an explicit closing tag (even <use href='...'></use>), so a
// plain stack is enough to reconstruct the nesting — no void-element table.
var tagRegex = regexp.MustCompile(`<(/?)([a-zA-Z0-9]+)([^>]*)>`)

var classAttrRegex = regexp.MustCompile(`class='([^']*)'`)

// ancestorClasses walks markup and returns the class lists of the elements
// enclosing the first element whose class list contains want, OUTERMOST first.
// It reads the real rendered markup rather than a hand-written expectation, so
// restructuring the row re-derives the chain instead of silently invalidating
// the assertions built on it.
func ancestorClasses(markup, want string) []string {
	var stack []string
	for _, m := range tagRegex.FindAllStringSubmatch(markup, -1) {
		closing, attrs := m[1] == "/", m[3]
		if closing {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			continue
		}
		cls := ""
		if c := classAttrRegex.FindStringSubmatch(attrs); c != nil {
			cls = c[1]
		}
		for _, f := range strings.Fields(cls) {
			if f == want {
				out := make([]string, len(stack))
				copy(out, stack)
				return out
			}
		}
		stack = append(stack, cls)
	}
	return nil
}

// baseRuleBlock returns the declaration block of the standalone `.class {` rule
// in @layer widgets, before any @media. Base only: the mobile override is a
// different mechanism (see the exemption noted in the test below), and the
// primitives layer groups shared flags into multi-selector rules that carry
// none of the positioning declarations this reads.
func baseRuleBlock(cssStr, class string) string {
	if i := strings.Index(cssStr, "@media"); i != -1 {
		cssStr = cssStr[:i]
	}
	if i := strings.Index(cssStr, "@layer widgets {"); i != -1 {
		cssStr = cssStr[i:]
	}
	i := strings.Index(cssStr, "."+class+" {")
	if i == -1 {
		return ""
	}
	body := cssStr[i+strings.Index(cssStr[i:], "{"):]
	end := strings.Index(body, "}")
	if end == -1 {
		return ""
	}
	return body[:end]
}

func declValue(block, prop string) string {
	for _, line := range strings.Split(block, ";") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prop+":") {
			return strings.TrimSpace(strings.TrimPrefix(line, prop+":"))
		}
	}
	return ""
}

// TestFlyoutHangsFromTheRowNotFromTheTrigger is the regression net for the
// desktop bug where the ⋮ dropdown painted OVER the row that opened it,
// covering its own label.
//
// Measured in the browser before this test existed (1440x900, first row):
//
//	row      top 113.2  bottom 163.2  (height 50)
//	summary  top 118.0  bottom 142.0  (height 24)
//	options  top 142.0                → 21.2px ABOVE its own row's bottom
//	options.offsetParent === .targetlist__menu    ← the smoking gun
//
// The mechanism is a containing-block theft, and it is invisible in the DSL:
//
//	PartRow     Anchor()                 → position: relative
//	PartMenu    Docked(Parent, ...)      → position: absolute   ← steals it
//	PartOptions Flyout(SideStart)        → position: absolute; inset-block-start: 100%
//
// Flyout's own doc says "the nearest Anchor() ancestor is what it hangs from",
// but CSS resolves `100%` against the nearest POSITIONED ancestor, which is the
// 24px-tall <details>, not the 50px row. The Anchor() on the row is dead code
// today: nothing ever hangs from it.
//
// The assertion is deliberately NOT "the nearest positioned ancestor must be
// the Anchor" — that phrasing would forbid a fix that keeps the <details>
// positioned but makes its box span the row. What every correct fix has in
// common is that the Flyout's containing block ENDS where the Anchor ends: an
// intermediate positioned element pinned on only one block edge (the other
// `auto`) is sized by its own content, so `100%` can never reach the row's
// bottom. That is the invariant, and it holds for both a spanning trigger and a
// restructured DOM.
//
// Mobile is exempt on purpose: there PartOptions is Docked(Viewport) →
// position: fixed, which resolves against the screen and never consults this
// chain at all.
func TestFlyoutHangsFromTheRowNotFromTheTrigger(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)

	markup := tl.buildRow(Item{ID: "7", Label: "Alpha", Description: "activo"}).String()
	cssStr := tl.RenderCSS().String()

	flyoutClass := string(clsMenuList)
	chain := ancestorClasses(markup, flyoutClass)
	if len(chain) == 0 {
		t.Fatalf("could not locate .%s in the rendered row markup:\n%s", flyoutClass, markup)
	}

	// targetlist has no Flyout any more — the options are an in-flow accordion
	// inside the row (see TestOptionsNeverLeaveTheFlow, which is the guard that
	// matters here now). Say so instead of passing in silence: a check whose
	// precondition has quietly vanished is worse than no check, because it
	// still reports green. The chain logic below stays live for usermenu, which
	// does hang a real Flyout off a real Anchor.
	if !strings.Contains(baseRuleBlock(cssStr, flyoutClass), "position: absolute") {
		t.Skipf("no Flyout on .%s: the options are in flow, so there is no "+
			"containing block to steal. TestOptionsNeverLeaveTheFlow is what "+
			"keeps it that way.", flyoutClass)
	}

	// Walk inward from the outermost ancestor; the anchor is the innermost
	// `position: relative` on the chain, and everything BETWEEN it and the
	// flyout has to be transparent to the 100% resolution.
	anchorAt := -1
	for i, classes := range chain {
		for _, c := range strings.Fields(classes) {
			if declValue(baseRuleBlock(cssStr, c), "position") == "relative" {
				anchorAt = i
			}
		}
	}
	if anchorAt == -1 {
		t.Fatalf("no Anchor() (position: relative) ancestor found for .%s; chain=%v",
			flyoutClass, chain)
	}

	for _, classes := range chain[anchorAt+1:] {
		for _, c := range strings.Fields(classes) {
			block := baseRuleBlock(cssStr, c)
			pos := declValue(block, "position")
			if pos != "absolute" && pos != "fixed" {
				continue // in-flow: not a containing block, harmless
			}
			start := declValue(block, "inset-block-start")
			end := declValue(block, "inset-block-end")
			if start == "auto" || end == "auto" || start == "" || end == "" {
				t.Errorf(""+
					".%s sits between the Anchor (.%s) and the Flyout (.%s) and is "+
					"position: %s pinned on one block edge only "+
					"(inset-block-start: %q, inset-block-end: %q).\n"+
					"It becomes the Flyout's containing block, so the Flyout's "+
					"inset-block-start: 100%% resolves against THIS element's own "+
					"content height instead of the Anchor's — the dropdown opens "+
					"inside its own row and covers the row's content.\n"+
					"Either pin it on both block edges so its box spans the Anchor, "+
					"or take it off the chain between the two.",
					c, chain[anchorAt], flyoutClass, pos, start, end)
			}
		}
	}
}
