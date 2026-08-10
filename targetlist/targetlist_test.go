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

// selectorBodies concatenates the body of every top-level "selector { ... }"
// block found in region, in order. A style.Sheet can legally emit the same
// selector more than once within one @layer (e.g. a flow shape block ahead
// of a value-carrying one — see On()'s handling of a breakpoint override
// that sets its own flow primitive), so a test asserting on "the" rule for a
// selector needs every block, not just the first.
func selectorBodies(region, selector string) string {
	var out strings.Builder
	needle := "\n" + selector + " {"
	rest := region
	for {
		i := strings.Index(rest, needle)
		if i == -1 {
			break
		}
		rest = rest[i+len(needle):]
		end := strings.Index(rest, "}")
		if end == -1 {
			break
		}
		out.WriteString(rest[:end])
		out.WriteString("\n")
		rest = rest[end+1:]
	}
	return out.String()
}

func TestTargetList_RowHasLabelBadgeAndMenu(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)

	html := tl.buildRow(Item{ID: "7", Label: "Alpha", Description: "192.168.0.7"}).String()

	for _, want := range []string{"targetlist__row", "Alpha", "targetlist__badge", "192.168.0.7", "targetlist__button", "targetlist__options", "Eliminar"} {
		if !strings.Contains(html, want) {
			t.Errorf("buildRow output missing %q\ngot: %s", want, html)
		}
	}
	// Editar left the menu with the lock it existed to undo: the ⋮ opens a
	// single-option accordion now, and tapping a row already leaves the form
	// editable. Neither the label nor the pencil glyph may come back.
	if strings.Contains(html, "Editar") {
		t.Errorf("buildRow must not render an Editar option (the lock is gone)\ngot: %s", html)
	}
	if strings.Contains(html, "tl-edit") {
		t.Errorf("buildRow must not reference the tl-edit icon\ngot: %s", html)
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

// The accordion is exclusive by construction: openMenu holds ONE id, so there
// is no state in which two rows are expanded. This replaced a native
// <details name="…"> group, whose open state lived in the DOM and had to be
// read back out of it — see the openMenu comment in targetlist.go.
func TestTargetList_OnlyOneRowExpandsAtATime(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	tl.SetItems([]Item{{ID: "1", Label: "One"}, {ID: "2", Label: "Two"}})

	// buildRow, not Render(): the rows are cached in a SignalNodes that the
	// wrapper only reconciles in the browser, but each row's own state binding
	// evaluates on String(), which is what this is about.
	expanded := func() []string {
		var open []string
		for _, it := range tl.Items() {
			if strings.Contains(tl.buildRow(it).String(), `data-open='true'`) {
				open = append(open, it.ID)
			}
		}
		return open
	}

	if got := expanded(); len(got) != 0 {
		t.Errorf("no row may be expanded initially, got %v", got)
	}

	tl.openMenu.Set("1")
	if got := expanded(); len(got) != 1 || got[0] != "1" {
		t.Errorf("expanding row 1 must expand exactly row 1, got %v", got)
	}

	tl.openMenu.Set("2")
	if got := expanded(); len(got) != 1 || got[0] != "2" {
		t.Errorf("expanding row 2 must collapse row 1, got %v", got)
	}

	tl.closeAllMenus()
	if tl.openMenu.Get() != "" {
		t.Errorf("closeAllMenus must clear openMenu, got %q", tl.openMenu.Get())
	}
	if got := expanded(); len(got) != 0 {
		t.Errorf("closeAllMenus must collapse every row, got %v", got)
	}
}

func TestTargetList_CSSDoesNotContainHas(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	css := tl.RenderCSS().String()

	if strings.Contains(css, ":has(") {
		t.Error("expected CSS not to contain forbidden :has( selector")
	}
	// The options are hidden until their own row's id is the one in openMenu:
	// RevealedBy(Open) is the whole open/close mechanism now that no native
	// <details> is doing it.
	if !strings.Contains(css, "display: none") {
		t.Error("expected CSS to hide the options by default")
	}
	if !strings.Contains(css, "[data-open=\"true\"]") {
		t.Error("expected CSS to contain selector matching [data-open=\"true\"]")
	}
	if !strings.Contains(css, "display: flex") {
		t.Error("expected the revealed options to get a real flow back")
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

	// Mobile media rule: the inline gutter stays inline on mobile too — the
	// two-column indent budget (see crudview's cardInset) — never the seam.
	mediaIdx := strings.Index(css, "@media (max-width")
	if mediaIdx == -1 {
		t.Fatal("expected a mobile media query")
	}
	mobileRegion := css[mediaIdx:]
	if next := strings.Index(mobileRegion[1:], "@media"); next != -1 {
		mobileRegion = mobileRegion[:next+1]
	}
	// Mobile now also overrides Stack's gap (to match PadInline's own bump —
	// see listgap.GapMobile), so On() emits the part's full flow shape
	// (display/flex-direction/gap:var(--gap)/min-height) as its own block
	// ahead of the value-carrying one (--gap/padding-inline) — the same
	// two-block pattern platformd's own On(Mobile, …, Stack(...)) parts
	// already produce. Collect every .targetlist__list block in the region,
	// not just the first, or this assertion would inspect the shape block
	// and never see padding-inline at all.
	mobileBody := selectorBodies(mobileRegion, ".targetlist__list")
	if mobileBody == "" {
		t.Fatal("expected a mobile rule for .targetlist__list")
	}
	if !strings.Contains(mobileBody, "padding-inline: var(--space-2") {
		t.Errorf("expected the mobile gutter to be the Space2 inline pad, block:\n%s", mobileBody)
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

// TestMenuTrailsTheRow supersedes an earlier net that required the OPPOSITE
// placement (the ⋮ trigger as the row's FIRST child, leading edge, for the
// mobile master-detail sliver's sake — see git history for that version's
// full rationale). The trigger is now DOM-LAST, on purpose: PartLabel's
// Grow() already claims 100% of the row's free space during flex
// resolution, so a margin-auto push on a leading trigger had nothing left
// to distribute — trailing placement is the only one that actually lands
// at the row's trailing edge (see css.go's PartButton comment). The
// trade-off — the ⋮ is no longer reachable from the mobile sliver — is
// accepted, not overlooked.
func TestMenuTrailsTheRow(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	markup := tl.buildRow(Item{ID: "7", Label: "Alpha"}).String()

	// The row's LAST in-flow child (before the options panel) must be the
	// ⋮ trigger: in a Row() flow, order is the position.
	last, bestIdx := "", -1
	if i := strings.Index(markup, "targetlist__row"); i != -1 {
		if j := strings.Index(markup[i:], ">"); j != -1 {
			rest := markup[i+j+1:]
			for _, c := range []string{"targetlist__button", "targetlist__label", "targetlist__badge"} {
				if k := strings.Index(rest, c); k != -1 && k > bestIdx {
					bestIdx = k
					last = c
				}
			}
		}
	}
	if last != "targetlist__button" {
		t.Errorf("expected the ⋮ trigger to be the row's LAST child (trailing edge by flex order), got last=%q\nmarkup: %s", last, markup)
	}

	// And it must be in flow: no position declarations on its own base rule.
	cssStr := tl.RenderCSS().String()
	block := baseRuleBlock(cssStr, string(clsMenuBtn))
	if block == "" {
		t.Fatal("expected a base rule for .targetlist__button")
	}
	for _, prop := range []string{"position", "inset-block-start", "inset-block-end", "inset-inline-start"} {
		if declValue(block, prop) != "" {
			t.Errorf("expected .targetlist__button to be in flow (no %s), block:\n%s", prop, block)
		}
	}
}

func TestSheetValidates(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	if errs := tl.sheet().Validate(); len(errs) > 0 {
		t.Errorf("targetlist sheet must validate, got:\n%v", errs)
	}
}

// TestOptionsNeverLeaveTheFlow is the regression net for BOTH shipped bugs of
// this panel, which had one root: an out-of-flow options panel cannot coexist
// with the list's own Scroll() region.
//
//  1. As a Flyout it resolved inset-block-start: 100% against the nearest
//     POSITIONED ancestor. Measured at 1440x900: the panel opened 21.2px inside
//     its own row and covered 8.2px of the row's label.
//  2. Anchored correctly to the row, it was then clipped by the scroller: on
//     the last row, 10px of an 84.8px panel survived.
//  3. The mobile Docked(Viewport, …) escape hatch from (2) landed the buttons
//     502px from their row, over two unrelated rows, which then needed a
//     Veil()'d backdrop that blurred the very row being acted on.
//
// In flow there is no clipper to escape and nothing to disambiguate. Measured
// after: the same last row shows 41.6px of 41.6px.
//
// So: no positioning on the options, on ANY device. A future device override
// that reaches for absolute/fixed to "get more room" is re-entering the loop.
func TestOptionsNeverLeaveTheFlow(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	cssStr := tl.RenderCSS().String()

	for _, blk := range strings.Split(cssStr, "}") {
		if !strings.Contains(blk, ".targetlist__options") {
			continue
		}
		for _, banned := range []string{"position: absolute;", "position: fixed;"} {
			if strings.Contains(blk, banned) {
				t.Errorf("the options must stay in flow (found %q); an out-of-flow "+
					"panel is clipped by the list's Scroll() region, and escaping "+
					"that clipper is what detached it from its row.\nblock:%s",
					banned, blk)
			}
		}
	}

	// And the row must be able to grow to hold them: inside a Scroll() column a
	// flex item defaults to flex-shrink: 1, which pins the row at its
	// min-height and lets the options paint outside its box.
	if b := baseRuleBlock(cssStr, string(clsRow)); !strings.Contains(b, "flex-shrink: 0") {
		if p := ruleContaining(cssStr, ".targetlist__row", "flex-shrink"); !strings.Contains(p, "flex-shrink: 0") {
			t.Errorf("the row must not shrink, or it cannot grow to contain the expanded options; got:\n%s%s", b, p)
		}
	}
}

// ruleContaining returns the first rule block whose selector list mentions sel
// and whose body mentions prop. The primitives layer groups shared flags across
// selectors, so a per-part lookup can legitimately miss them.
func ruleContaining(cssStr, sel, prop string) string {
	for _, blk := range strings.Split(cssStr, "}") {
		if strings.Contains(blk, sel) && strings.Contains(blk, prop) {
			return blk
		}
	}
	return ""
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
