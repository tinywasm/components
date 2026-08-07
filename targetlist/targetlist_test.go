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

func TestTargetList_RowHasLabelBadgeAndMenu(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)

	html := tl.buildRow(Item{ID: "7", Label: "Alpha", Description: "192.168.0.7"}).String()

	for _, want := range []string{"targetlist__row", "Alpha", "targetlist__badge", "192.168.0.7", "targetlist__menu", "Editar", "Eliminar"} {
		if !strings.Contains(html, want) {
			t.Errorf("buildRow output missing %q\ngot: %s", want, html)
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

func TestTargetList_MenuOpenStateBackdrop(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	tl.SetItems([]Item{{ID: "1", Label: "One"}})

	// Initially, with no open menus, backdrop shouldn't have data-open="true"
	htmlInit := tl.Render().String()
	t.Logf("htmlInit: %s", htmlInit)
	if strings.Contains(htmlInit, `data-open='true'`) {
		t.Error("expected backdrop NOT to have data-open='true' initially")
	}

	// Mocking menu open
	tl.menuOpen.Set(true)
	htmlOpen := tl.Render().String()
	t.Logf("htmlOpen: %s", htmlOpen)
	if !strings.Contains(htmlOpen, `data-open='true'`) {
		t.Error("expected backdrop to have data-open='true' when a menu is open")
	}

	// Verify closeAllMenus clears it
	tl.closeAllMenus()
	if tl.menuOpen.Get() {
		t.Error("expected menuOpen to be false after closeAllMenus")
	}
}

func TestTargetList_CSSDoesNotContainHas(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	css := tl.RenderCSS().String()

	if strings.Contains(css, ":has(") {
		t.Error("expected CSS not to contain forbidden :has( selector")
	}
	if !strings.Contains(css, "display: none") {
		t.Error("expected CSS to have display: none under normal condition for backdrop")
	}
	if !strings.Contains(css, "[data-open=\"true\"]") {
		t.Error("expected CSS to contain selector matching [data-open=\"true\"]")
	}
	if !strings.Contains(css, "display: block") {
		t.Error("expected CSS to have display: block under open condition")
	}
}

func TestTargetList_SelectionUsesAccent(t *testing.T) {
	// PLAN v0.2.0 item 3: the selected row wears the amber Accent surface —
	// the same "where I am" statement the rail's current nav item makes —
	// never the 15% blue Highlight wash, which was close to invisible.
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

// TestTargetList_ListGutterDoesNotClobberTheScrollSeam is the net for
// PLAN.md Stage 1: the list's own gutter is inline-only, on every breakpoint.
// A `padding:` shorthand would land in the widgets layer and override the
// primitives-layer seam (padding-block-end: var(--floating-bottom, 0px)),
// silently putting the last row's badge back under a FloatingChrome host's
// button. The seam must keep owning the block edges.
func TestTargetList_ListGutterDoesNotClobberTheScrollSeam(t *testing.T) {
	css := (&TargetList{}).RenderCSS().String()

	// The seam itself must be present (Scroll() emits it in primitives).
	if !strings.Contains(css, "padding-block-end: var(--floating-bottom, 0px);") {
		t.Errorf("expected Scroll() to emit the floating-bottom seam, got:\n%s", css)
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
	if strings.Contains(baseBody, "padding: ") || strings.Contains(baseBody, "padding-block") {
		t.Errorf("the base rule must not clobber the seam's block edges, block:\n%s", baseBody)
	}

	// Mobile media rule: the flush reclaims the sliver inline, never the seam.
	mediaIdx := strings.Index(css, "@media (max-width")
	if mediaIdx == -1 {
		t.Fatal("expected a mobile media query")
	}
	mobileRegion := css[mediaIdx:]
	if next := strings.Index(mobileRegion[1:], "@media"); next != -1 {
		mobileRegion = mobileRegion[:next+1]
	}
	mi := strings.Index(mobileRegion, "\n.targetlist__list {")
	if mi == -1 {
		t.Fatal("expected a mobile rule for .targetlist__list")
	}
	mobileBody := mobileRegion[mi:]
	if end := strings.Index(mobileBody, "}"); end != -1 {
		mobileBody = mobileBody[:end]
	}
	if !strings.Contains(mobileBody, "padding-inline: 0;") {
		t.Errorf("expected the mobile flush to be inline-only, block:\n%s", mobileBody)
	}
	if strings.Contains(mobileBody, "padding: ") || strings.Contains(mobileBody, "padding-block") {
		t.Errorf("the mobile rule must not clobber the seam's block edges, block:\n%s", mobileBody)
	}
}

// TestTargetList_BadgeStraddlesWithoutTransform is the net for PLAN.md
// Stage 2: OnEdge straddles with half a --chip-height of negative margin, not
// a transform — the badge's box must stay visible to scroll-size calculations
// so a host's FloatingChrome reservation applies to where it really paints.
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

// TestMenuDocksToLeadingEdge is the net for PLAN.md Stage A1: on the mobile
// master-detail strip, selecting a row leaves only a sliver of the list's
// LEADING edge visible. A row menu anchored to the trailing edge is the one
// part of the row that sliver can never show, stranding the only control that
// unlocks the now-read-only form on the panel the user just left tapping the
// row navigated away from.
func TestMenuDocksToLeadingEdge(t *testing.T) {
	css := (&TargetList{}).RenderCSS().String()
	// Leading newline: without it this also matches the substring
	// ".targetlist__menu {" inside the primitives-layer combined selector
	// ".targetlist__badge, .targetlist__menu {", which appears earlier in the
	// sheet and carries no inset declarations at all.
	i := strings.Index(css, "\n.targetlist__menu {")
	if i == -1 {
		t.Fatal("expected a standalone rule for .targetlist__menu")
	}
	body := css[i:]
	end := strings.Index(body, "}")
	if end == -1 {
		t.Fatal("malformed rule block")
	}
	b := body[:end]
	if !strings.Contains(b, "inset-inline-start:") {
		t.Errorf("expected the menu docked to the leading edge (inset-inline-start), block:\n%s", b)
	}
	if !strings.Contains(b, "inset-inline-end: auto;") {
		t.Errorf("expected the trailing edge left auto, block:\n%s", b)
	}
}

// TestOptionsPanelStaysOnScreenOnMobile is the net for PLAN.md Stage A2: once
// the menu trigger sits at the row's leading edge (~372px into a 375-400px
// viewport, itself already only a 10% sliver in from the panel edge), a
// row-anchored Flyout for the dropdown overflows the screen by roughly its
// own width. Mobile must anchor it to the viewport instead; desktop keeps the
// row-anchored Flyout, now matching the trigger's own leading-edge side.
func TestOptionsPanelStaysOnScreenOnMobile(t *testing.T) {
	css := (&TargetList{}).RenderCSS().String()

	// Desktop (ungated) rule: a row-anchored Flyout, SideStart to match the
	// trigger — not the SideEnd it used before Stage A1 moved the trigger.
	// Anchored to start searching after "@layer widgets {": the primitives
	// layer emits its OWN standalone ".targetlist__options {" rule (just
	// HideOverflow's "overflow: hidden;"), which a bare or newline-anchored
	// search finds first and which carries no position/inset declarations at
	// all.
	widgetsIdx := strings.Index(css, "@layer widgets {")
	if widgetsIdx == -1 {
		t.Fatal("expected an @layer widgets block")
	}
	di := strings.Index(css[widgetsIdx:], "\n.targetlist__options {")
	if di == -1 {
		t.Fatal("expected a widgets-layer rule for .targetlist__options")
	}
	di += widgetsIdx
	desktopBody := css[di:]
	if end := strings.Index(desktopBody, "}"); end != -1 {
		desktopBody = desktopBody[:end]
	}
	if !strings.Contains(desktopBody, "position: absolute;") {
		t.Errorf("expected the desktop options panel to stay a row-anchored Flyout, block:\n%s", desktopBody)
	}
	if !strings.Contains(desktopBody, "inset-inline-start: 0;") {
		t.Errorf("expected the desktop Flyout on the leading (start) side, block:\n%s", desktopBody)
	}

	// Mobile media query: viewport-docked, not row-anchored. The device rule
	// emits TWO consecutive blocks for the same selector — one for Stack's
	// flow decls (display/flex-direction/gap), one for Docked's position/inset
	// decls — so LastIndex, not the first match, is the one carrying
	// "position: fixed".
	mediaIdx := strings.Index(css, "@media (max-width")
	if mediaIdx == -1 {
		t.Fatal("expected a mobile media query")
	}
	mobileRegion := css[mediaIdx:]
	if next := strings.Index(mobileRegion[1:], "@media"); next != -1 {
		mobileRegion = mobileRegion[:next+1]
	}
	mi := strings.LastIndex(mobileRegion, "\n.targetlist__options {")
	if mi == -1 {
		t.Fatal("expected a mobile rule for .targetlist__options")
	}
	mobileBody := mobileRegion[mi:]
	if end := strings.Index(mobileBody, "}"); end != -1 {
		mobileBody = mobileBody[:end]
	}
	if !strings.Contains(mobileBody, "position: fixed;") {
		t.Errorf("expected the mobile options panel docked to the viewport (position: fixed), block:\n%s", mobileBody)
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
