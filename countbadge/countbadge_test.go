//go:build !wasm

package countbadge

import (
	"strings"
	"testing"

	. "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
)

func newBadge(n string, visible bool) *CountBadge {
	count := NewString(n)
	show := NewBool(visible)
	return &CountBadge{Count: count, Visible: show}
}

func TestCountBadge_Render(t *testing.T) {
	b := newBadge("3", true)

	html := b.Render().String()

	if !strings.HasPrefix(html, "<span") {
		t.Error("expected span tag, got: " + html)
	}
	if !strings.Contains(html, "countbadge__badge") {
		t.Error("expected countbadge__badge class, got: " + html)
	}
	if !strings.Contains(html, ">3<") {
		t.Error("expected the count as text, got: " + html)
	}
}

// TestCountBadge_VisibilityIsAState proves the bubble paints only when told:
// the stylesheet hides it by default and reveals it on the same Open state
// the markup binds (data-open is only written in a live DOM, so the stdlib
// render cannot show it — the states layer is where the contract lives).
// Hidden at zero (Visible=false), shown above zero.
func TestCountBadge_VisibilityIsAState(t *testing.T) {
	cssStr := (&CountBadge{}).RenderCSS().String()

	if !strings.Contains(cssStr, `.countbadge__badge[data-open="true"]`) {
		t.Error("expected the stylesheet to reveal the badge on data-open, got:\n" + cssStr)
	}
}

// TestCountBadge_BadgeIsOutOfFlow is the regression net for the defect that
// motivated this component: an in-flow counter stretches its button the
// moment it appears. The bubble must be absolute, so the host box never
// depends on the count.
func TestCountBadge_BadgeIsOutOfFlow(t *testing.T) {
	cssStr := (&CountBadge{}).RenderCSS().String()

	mark := ".countbadge__badge {"
	i := strings.Index(cssStr, mark)
	if i == -1 {
		t.Fatal("expected a rule for .countbadge__badge")
	}
	body := cssStr[i:]
	if end := strings.Index(body, "}"); end != -1 {
		body = body[:end]
	}
	if !fmt.Contains(body, "position: absolute") {
		t.Errorf("the badge must be out of flow (position: absolute), block:\n%s", body)
	}
}

func TestPairMarkupAndStylesheet(t *testing.T) {
	cssStr := (&CountBadge{}).RenderCSS().String()
	if !strings.Contains(cssStr, ".countbadge__badge") {
		t.Error("expected .countbadge__badge in the stylesheet")
	}
	html := newBadge("3", true).Render().String()
	if !strings.Contains(html, "countbadge__badge") {
		t.Error("expected countbadge__badge in the markup")
	}
	if strings.Contains(cssStr, "!important") {
		t.Error("stylesheet contains forbidden !important directive")
	}
}
