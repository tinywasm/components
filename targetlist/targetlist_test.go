//go:build !wasm

package targetlist

import (
	"regexp"
	"strings"
	"testing"
)

var (
	classRegex      = regexp.MustCompile(`\.([a-zA-Z0-9_-]+)`)
	htmlClassRegex  = regexp.MustCompile(`class='([^']*)'`)
	htmlClassRegex2 = regexp.MustCompile(`class="([^"]*)"`)
)

func TestTargetList_RowHasLabelBadgeAndCheck(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)

	html := tl.buildRow(Item{ID: "7", Label: "Alpha", Description: "192.168.0.7"}).String()

	for _, want := range []string{"targetlist__row", "Alpha", "targetlist__badge", "192.168.0.7", "targetlist__check"} {
		if !strings.Contains(html, want) {
			t.Errorf("buildRow output missing %q\ngot: %s", want, html)
		}
	}
	for _, unwanted := range []string{"targetlist__button", "targetlist__options", "Eliminar", "Editar"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("buildRow must not render per-row menu artifact %q\ngot: %s", unwanted, html)
		}
	}
}

func TestTargetList_SetItemsPopulatesRows(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	tl.SetItems([]Item{{ID: "1", Label: "One"}, {ID: "2", Label: "Two"}})

	if got := len(tl.rows.Get()); got != 2 {
		t.Fatalf("expected 2 rows, got %d", got)
	}
}

func TestTargetList_CSSDoesNotContainHas(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	css := tl.RenderCSS().String()

	if strings.Contains(css, ":has(") {
		t.Error("expected CSS not to contain forbidden :has( selector")
	}
}

func TestTargetList_SelectionUsesAccent(t *testing.T) {
	css := (&TargetList{}).RenderCSS().String()
	i := strings.Index(css, `.targetlist__row[data-selected="true"] {`)
	if i == -1 {
		t.Fatal("expected a rule for the selected row state")
	}
	body := css[i:]
	end := strings.Index(body, "}")
	if end == -1 {
		t.Fatal("malformed rule block")
	}
	b := body[:end]
	if !strings.Contains(b, "--color-accent") {
		t.Errorf("selected row must use the Accent surface, block:\n%s", b)
	}
	if strings.Contains(b, "--color-selection") {
		t.Errorf("selected row must not use the Highlight surface, block:\n%s", b)
	}
}

// Restored: neither of these two guards a menu, and the PR that removed the
// per-row menu took them out with it. Both protect bugs this package has
// already been bitten by once — see the comments inside them.
func TestTargetList_ListGutterDoesNotClobberTheScrollSeam(t *testing.T) {
	css := (&TargetList{}).RenderCSS().String()

	// The seam is present in the primitives layer, now with the ambient gutter
	// added to it (ScrollGutter folds both into one calc so neither can be lost
	// to a later plain override).
	if !strings.Contains(css, "padding-block-end: calc(var(--floating-bottom, 0px) + var(--space-1") {
		t.Errorf("expected Scroll()+ScrollGutter to emit the additive floating-bottom seam, got:\n%s", css)
	}
	if !strings.Contains(css, "padding-block-start: calc(var(--floating-top, 0px) + var(--space-1") {
		t.Errorf("expected a symmetric additive top seam, got:\n%s", css)
	}

	// Base widgets-layer rule: inline gutter only, no shorthand, no block pad.
	widgetsIdx := strings.Index(css, "@layer widgets {")
	if widgetsIdx == -1 {
		t.Fatal("expected an @layer widgets block")
	}
	bi := strings.Index(css[widgetsIdx:], "\n.targetlist__list {")
	if bi == -1 {
		t.Fatal("expected a widgets-layer rule for .targetlist__list")
	}
	bi += widgetsIdx
	baseBody := css[bi:]
	if end := strings.Index(baseBody, "}"); end != -1 {
		baseBody = baseBody[:end]
	}
	if !strings.Contains(baseBody, "padding-inline: var(--space-1") {
		t.Errorf("expected the base rule to keep its inline gutter, block:\n%s", baseBody)
	}
	// The widgets-layer rule carries the inline gutter only. The block gutter
	// is the primitives-layer additive calc (asserted above) — never a plain
	// padding-block here, which would replace it outright, and never a
	// `padding:` shorthand.
	if strings.Contains(baseBody, "padding: ") || strings.Contains(baseBody, "padding-block") {
		t.Errorf("the widgets-layer base rule must not carry a block padding of its own, block:\n%s", baseBody)
	}

	// Mobile: the additive seam is re-emitted for this breakpoint with the
	// Space2 inset (matching the mobile lateral inset), and the widgets-layer
	// mobile rule still carries the inline gutter only.
	mediaIdx := strings.Index(css, "@media (max-width")
	if mediaIdx == -1 {
		t.Fatal("expected a mobile media query")
	}
	mobileRegion := css[mediaIdx:]
	if next := strings.Index(mobileRegion[1:], "@media"); next != -1 {
		mobileRegion = mobileRegion[:next+1]
	}
	if !strings.Contains(mobileRegion, "padding-block-end: calc(var(--floating-bottom, 0px) + var(--space-2") {
		t.Errorf("expected the mobile additive seam with the Space2 gutter, region:\n%s", mobileRegion)
	}
	// The value-carrying widgets-layer mobile block: the one with --gap. It
	// must hold the inline gutter and nothing on the block axis.
	widgetsMobile := ruleContaining(mobileRegion, ".targetlist__list", "--gap:")
	if widgetsMobile == "" {
		t.Fatal("expected a widgets-layer mobile rule for .targetlist__list carrying --gap")
	}
	if !strings.Contains(widgetsMobile, "padding-inline: var(--space-2") {
		t.Errorf("expected the mobile widgets rule to keep the Space2 inline pad, block:\n%s", widgetsMobile)
	}
	if strings.Contains(widgetsMobile, "padding: ") || strings.Contains(widgetsMobile, "padding-block") {
		t.Errorf("the mobile widgets rule must not carry a block padding of its own, block:\n%s", widgetsMobile)
	}
}

func TestTargetList_BadgeStraddlesWithoutTransform(t *testing.T) {
	css := (&TargetList{}).RenderCSS().String()
	i := strings.Index(css, "\n.targetlist__badge {")
	if i == -1 {
		t.Fatal("expected a rule for .targetlist__badge")
	}
	body := css[i:]
	if end := strings.Index(body, "}"); end != -1 {
		body = body[:end]
	}
	if !strings.Contains(body, "margin-block-end: calc(-0.5 * var(--chip-height") {
		t.Errorf("expected the badge to straddle the row's bottom line by half a chip, block:\n%s", body)
	}
	if strings.Contains(body, "transform") {
		t.Errorf("the badge must not use a transform to straddle, block:\n%s", body)
	}
}

func TestSheetValidates(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	if errs := tl.sheet().Validate(); len(errs) > 0 {
		t.Errorf("targetlist sheet must validate, got:\n%v", errs)
	}
}

func TestPairMarkupAndStylesheet(t *testing.T) {
	extractCSSClasses := func(css string) map[string]bool {
		classes := make(map[string]bool)
		matches := classRegex.FindAllStringSubmatch(css, -1)
		for _, m := range matches {
			if len(m) > 1 {
				classes[m[1]] = true
			}
		}
		return classes
	}

	extractHTMLClasses := func(html string) map[string]bool {
		classes := make(map[string]bool)
		matches1 := htmlClassRegex.FindAllStringSubmatch(html, -1)
		for _, m := range matches1 {
			if len(m) > 1 {
				for _, cls := range strings.Fields(m[1]) {
					classes[cls] = true
				}
			}
		}
		matches2 := htmlClassRegex2.FindAllStringSubmatch(html, -1)
		for _, m := range matches2 {
			if len(m) > 1 {
				for _, cls := range strings.Fields(m[1]) {
					classes[cls] = true
				}
			}
		}
		return classes
	}

	filterClasses := func(classes map[string]bool, prefix string) map[string]bool {
		filtered := make(map[string]bool)
		for cls := range classes {
			if strings.HasPrefix(cls, prefix) {
				filtered[cls] = true
			}
		}
		return filtered
	}

	tl := &TargetList{}
	tl.Init(nil)
	html := tl.Render().String() + tl.buildRow(Item{ID: "1", Label: "A", Description: "B"}).String()
	css := tl.RenderCSS().String()

	htmlClasses := filterClasses(extractHTMLClasses(html), "targetlist")
	cssClasses := filterClasses(extractCSSClasses(css), "targetlist")

	for cls := range cssClasses {
		if !htmlClasses[cls] {
			t.Errorf("TargetList CSS class %q does not exist in rendered HTML", cls)
		}
	}
	for cls := range htmlClasses {
		if !cssClasses[cls] {
			t.Errorf("TargetList HTML class %q is unstyled in CSS", cls)
		}
	}
}

// ruleContaining returns the first CSS block that mentions both sel and prop.
// Used by TestTargetList_BadgeStraddlesWithoutTransform to read one rule out of
// the emitted sheet without parsing it.
func ruleContaining(cssStr, sel, prop string) string {
	for _, blk := range strings.Split(cssStr, "}") {
		if strings.Contains(blk, sel) && strings.Contains(blk, prop) {
			return blk
		}
	}
	return ""
}
