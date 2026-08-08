//go:build !wasm

package usermenu

import (
	"regexp"
	"strings"
	"testing"
)

// tagRegex matches an opening or closing tag. Every element this component
// renders has an explicit closing tag, so a plain stack is enough to
// reconstruct the nesting — no void-element table.
var tagRegex = regexp.MustCompile(`<(/?)([a-zA-Z0-9]+)([^>]*)>`)

var classAttrRegex = regexp.MustCompile(`class='([^']*)'`)

// ancestorClasses walks markup and returns the class lists of the elements
// enclosing the first element whose class list contains want, OUTERMOST first.
// It reads the real rendered markup rather than a hand-written expectation, so
// restructuring the tree re-derives the chain instead of silently invalidating
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
// in @layer widgets, before any @media. Base only: the mobile accordion is a
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
		line = strings.TrimPrefix(line, "{")
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prop+":") {
			return strings.TrimSpace(strings.TrimPrefix(line, prop+":"))
		}
	}
	return ""
}

// TestPanelHangsFromTheRootAnchor is the reference case the targetlist anchor
// contract test contrasts against: Root(Anchor()) → PartPanel with
// Flyout(SideEnd), with NOTHING positioned in between. It is the shape the
// style DSL documents as the legal one — and, until this test existed, nobody
// watched it: usermenu shipped with no test files at all.
//
// The panel's containing block is the root itself: its inset resolves against
// the Anchor, so the flyout hangs exactly where the trigger's own box ends.
// If a Docked/OnEdge/Backdrop part ever lands between the two, this test turns
// red — the same interposition that broke targetlist's row menus.
//
// Mobile is exempt on purpose: there the panel is an in-flow accordion inside
// a drawer (no Flyout at all), so there is no chain to hang from — and the
// test only reads base rules, outside @media.
func TestPanelHangsFromTheRootAnchor(t *testing.T) {
	m := &UserMenu{}
	m.Init(nil)
	m.Name = "Ada"
	m.Avatar = "https://example.invalid/ada.png"

	markup := m.Render().String()
	cssStr := m.RenderCSS().String()

	flyoutClass := string(clsPanel)
	chain := ancestorClasses(markup, flyoutClass)
	if len(chain) == 0 {
		t.Fatalf("could not locate .%s in the rendered markup:\n%s", flyoutClass, markup)
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
					"content height instead of the Anchor's.\n"+
					"Either pin it on both block edges so its box spans the Anchor, "+
					"or take it off the chain between the two.",
					c, chain[anchorAt], flyoutClass, pos, start, end)
			}
		}
	}
}
