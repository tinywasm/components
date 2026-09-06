//go:build !wasm

package listgap_test

import (
	"strings"
	"testing"

	"webtyp.com/components/listgap"
	"webtyp.com/css"
	"webtyp.com/widget"
	"webtyp.com/widget/style"
)

// demoList is a consumer-shaped stand-in for targetlist / targetdate: a real
// widget, the real style.Sheet, listgap through it exactly as a css.go would,
// then assert the CSS it emits.
type demoList struct{}

func (demoList) WidgetName() widget.Name { return widget.Name("demolist") }
func (demoList) WidgetKind() widget.Kind { return widget.Combobox }

const demoPartList = widget.Part("list")

func render() string {
	s := style.For(demoList{}).Root(style.Fill(), style.Stack(style.SpaceNone))
	listgap.Apply(s, demoPartList)
	s.On(css.Mobile, demoPartList, listgap.MobileOpts()...)
	return s.Stylesheet().String()
}

// listRule returns the body of the LAST `.demolist__list {` block inside region
// — the value-carrying one (the primitives/shape blocks come first).
func listRule(t *testing.T, region string) string {
	t.Helper()
	const sel = ".demolist__list {"
	i := strings.LastIndex(region, sel)
	if i == -1 {
		t.Fatalf("no %q in region:\n%s", sel, region)
	}
	body := region[i+len(sel):]
	if end := strings.Index(body, "}"); end != -1 {
		body = body[:end]
	}
	return body
}

func TestApply_VerticalRhythmClearsTheBadgeStraddle(t *testing.T) {
	css := render()
	mediaAt := strings.Index(css, "@media (max-width")
	if mediaAt == -1 {
		t.Fatalf("expected a mobile media query, got:\n%s", css)
	}
	base, mobile := css[:mediaAt], css[mediaAt:]

	// PartBadge straddles calc(-0.5 * --chip-height) below its own row — 10px
	// at the default --chip-height (1.25rem). The row rhythm must clear that
	// overshoot or the badge paints on top of the next row. Space4 (16px) and
	// Space6 (24px) each leave real margin; Space2/Space3 measured 2px of
	// actual overlap in the browser and must never come back here.
	if got := listRule(t, base); !strings.Contains(got, "--gap: var(--space-4") {
		t.Errorf("desktop row rhythm must be Space4 (clears the 10px badge overshoot), block:\n%s", got)
	}
	if got := listRule(t, mobile); !strings.Contains(got, "--gap: var(--space-6") {
		t.Errorf("mobile row rhythm must be Space6, block:\n%s", got)
	}
}

func TestApply_LateralInsetStaysOnTheCrudviewBudget(t *testing.T) {
	css := render()
	mediaAt := strings.Index(css, "@media (max-width")
	base, mobile := css[:mediaAt], css[mediaAt:]

	// The inset is deliberately NOT bumped with the vertical rhythm: it is a
	// summand of crudview's 16px master-detail indent budget.
	if got := listRule(t, base); !strings.Contains(got, "padding-inline: var(--space-1") {
		t.Errorf("desktop lateral inset must stay Space1, block:\n%s", got)
	}
	if got := listRule(t, mobile); !strings.Contains(got, "padding-inline: var(--space-2") {
		t.Errorf("mobile lateral inset must stay Space2, block:\n%s", got)
	}
}

func TestApply_TopBottomGutterMatchesTheSidesAndAddsToTheSeam(t *testing.T) {
	out := render()
	mediaAt := strings.Index(out, "@media (max-width")
	base, mobile := out[:mediaAt], out[mediaAt:]

	// ScrollGutter folds the ambient top/bottom gutter into the SAME calc as
	// the FloatingChrome reservation — additive, never a plain override — so
	// the primitives layer carries both in one declaration.
	for _, want := range []string{
		"padding-block-start: calc(var(--floating-top, 0px) + " + css.Space1.Var() + ");",
		"padding-block-end: calc(var(--floating-bottom, 0px) + " + css.Space1.Var() + ");",
	} {
		if !strings.Contains(base, want) {
			t.Errorf("expected the desktop gutter to add to the seam with %q, got:\n%s", want, base)
		}
	}
	for _, want := range []string{
		"padding-block-start: calc(var(--floating-top, 0px) + " + css.Space2.Var() + ");",
		"padding-block-end: calc(var(--floating-bottom, 0px) + " + css.Space2.Var() + ");",
	} {
		if !strings.Contains(mobile, want) {
			t.Errorf("expected the mobile gutter to add to the seam with %q, got:\n%s", want, mobile)
		}
	}

	// The WIDGETS-layer list rule itself must still carry no padding-block of
	// its own: ScrollGutter's decls live in the PRIMITIVES layer precisely so
	// a later widgets-layer declaration can never plainly override — and
	// therefore erase — the seam's reservation.
	for _, region := range []string{base, mobile} {
		body := listRule(t, region)
		if strings.Contains(body, "padding-block") || strings.Contains(body, "padding: ") {
			t.Errorf("the widgets-layer list rule must not touch the block edges (the seam owns that declaration), block:\n%s", body)
		}
	}
}
