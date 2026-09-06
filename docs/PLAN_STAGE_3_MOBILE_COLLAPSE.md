---
PLAN: "feat(calendarslider): collapse to a compact chip on mobile once a day is picked"
REVIEWER: none
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Part of a 3-stage calendarslider improvement — see `docs/PLAN.md` for the
> index.

## ⚠️ BLOCKING PREREQUISITE — do not dispatch this stage until it is satisfied

This stage calls `date.WeekdayName` and `date.ParseDateKey`, neither of
which exists in any published `webtyp/date` version yet. Both are added by
`https://github.com/webtyp/date/blob/main/docs/PLAN.md` (a separate,
already-written, self-contained plan in that repo) — which ALSO changes
`date.MonthName`'s existing return values from Spanish to English (a
breaking change to already-shipped behavior, not just an addition: see that
plan's Stage 1). Before dispatching THIS stage:

1. That `date` plan must be dispatched, merged, and published (new `date`
   version tagged).
2. `components/go.mod` must be bumped to that new version:
   `go get -u webtyp.com/date@latest` from `components/`, then
   `go mod tidy`.
3. Confirm `grep -n "func WeekdayName\|func ParseDateKey" $(go env GOPATH)/pkg/mod/webtyp.com/date@*/date.go`
   finds both before writing a single line of this stage.
4. **The `MonthName` breaking change ripples into 3 EXISTING tests that
   assert its old Spanish output — fix these in the SAME commit as the
   version bump, before touching anything else in this stage:**
   - `calendarslider_test.go`, `TestBuildMonthAgosto2026`: the check
     `Contains(children[7].String(), "Agosto 2026")` → `"August 2026"`.
   - `calendarslider_test.go`, `TestRenderStructure`: the check
     `Contains(htmlOut, "Agosto 2026")` → `"August 2026"`.
   - `calendarslider_wasm_test.go`, `TestAllMonthsAlwaysInDOM`: both
     `label != "Agosto 2026"` and `label != "Octubre 2026"` →
     `"August 2026"` / `"October 2026"`.

   These 3 assertions now hold because, with no dictionary registered
   anywhere in the test binary, `lang.Translate("August", 2026).String()`
   (Stage 1.2 of `PLAN_STAGE_1_ARROW_HOVER_LAYOUT.md`) passes the English
   words through unchanged — see that stage for why this is the correct,
   intentional new default, not a regression to paper over. Do this fix
   regardless of whether Stage 1 has already merged — the string these
   tests check comes from `date.MonthName` either way, and this prerequisite
   step is what makes it English from here on.

Independent of Stage 1 and Stage 2 otherwise — no ordering requirement
against those beyond this repo's own prerequisite and the test fix above.

# Plan — CalendarSlider Stage 3: mobile collapse to a compact chip

## Context

On mobile, once a day is picked, the full calendar keeps occupying the same
vertical space it always did — the caller (e.g. an hour-list next to it)
gets no more room even though the user is done choosing a date for now.
The fix: on mobile ONLY, once a day is selected, the calendar collapses to
one compact row — an icon and the date in natural language, e.g. "Monday 18
August 2026" by default, or "Lunes 18 Agosto 2026" for a host that has
registered a Spanish dictionary (see Stage 3.3's `buildCollapsed` below for
exactly what that means and why this component never does that
registration itself) — exactly the shape `components/selectsearch` already
uses for
its own closed state (`PartHeader`: icon left, content right). Tapping the
collapsed row re-expands the full calendar to change the date; picking a new
day collapses it again. Desktop is unaffected at every step — the full
calendar always stays visible there, regardless of selection.

This reuses `selectsearch`'s own open/close mechanism verbatim: a hidden
`<input type="checkbox">` plus a `<label for="...">` wrapping the visible
row. Clicking the label toggles the checkbox NATIVELY (no JS needed for that
part); a `"change"` listener syncs the checkbox into a Go signal, exactly
like `selectsearch.go`'s `toggle`/`c.isOpen` pair. Read that file's `toggle`
and `header` construction (`selectsearch.go` lines ~213-261) before writing
this stage — this plan does not re-explain the idiom, it names the exact
lines to copy the shape from.

The mobile/desktop split, and "hidden until picked, shown while collapsed",
both use `widget/style`'s `On`/`OnlyOn` + `RevealedBy` — verified against a
throwaway probe program before writing this plan (see "Verified mechanism").

**Scope limit, deliberate, do not extend it**: the calendar collapses ONLY
from the day-click handler (the user tapping a bookable day). A host that
sets `Selected` programmatically (`cal.Selected.Set("2026-08-11")`, already
documented in `README.md`) does NOT trigger a collapse — that path is
untouched by this stage. If that gap matters later, it is a follow-up plan,
not an implicit addition here.

## Verified mechanism (do not re-derive, use as given)

```go
style.For(w).
    Part(partStrip, /* existing base options, unchanged */).
    On(css.Mobile, partStrip, style.RevealedBy(widget.Open)).
    OnlyOn(css.Mobile, partCollapsed, style.RevealedBy(widget.Open)).
    Stylesheet()
```

emits (confirmed by running it) — outside any media query, `.w__strip`
keeps its existing unconditional base rule (always visible, exactly as
today) and `.w__collapsed` is `display: none` unconditionally; ONLY inside
`@media (max-width: ...)` do both become conditional on their OWN
`data-open` attribute:

```css
@media (max-width: 639.98px) {
  .w__collapsed { display: none; }
  .w__collapsed[data-open="true"] { display: block; }
  .w__strip { display: none; }
  .w__strip[data-open="true"] { display: flex; }
}
```

This is exactly the split needed: desktop never hides the strip and never
shows the collapsed row, no matter what; mobile shows exactly one of the two,
picked by each element's OWN `data-open`.

## Stage 3.1 — a local icon (no new shared `webtyp/icons` package)

The collapsed chip's icon is specific to this one component's own chrome —
the same situation `selectsearch` was in for its dropdown chevron
(`iconArrowDown`, defined locally, not in `webtyp/icons`). Follow that
exact precedent, not the shared-icon-per-package convention `webtyp/icons`
uses for cross-component action glyphs (trash/pencil/etc.) — a local icon is
correct here, do not create a new `webtyp/icons` subpackage for this.

**New file: `svg.go`** (this component has no icons today, so no such file
exists yet — this is the SSR-split convention every icon-bearing component
follows: geometry behind `//go:build !wasm`, never in the main file):

```go
//go:build !wasm

package calendarslider

import "webtyp.com/svg/sprite"

// IconSvg ships the collapsed chip's calendar glyph. FontAwesome Free 6
// "calendar-day" (solid), viewBox 0 0 448 512.
func (c *CalendarSlider) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(
		sprite.Define(iconCalendar, "0 0 448 512",
			sprite.Path("M128 0c17.7 0 32 14.3 32 32l0 32 128 0 0-32c0-17.7 14.3-32 32-32s32 14.3 32 32l0 32 48 0c26.5 0 48 21.5 48 48l0 48L0 160l0-48C0 85.5 21.5 64 48 64l48 0 0-32c0-17.7 14.3-32 32-32zM0 192l448 0 0 272c0 26.5-21.5 48-48 48L48 512c-26.5 0-48-21.5-48-48L0 192zm80 64c-8.8 0-16 7.2-16 16l0 96c0 8.8 7.2 16 16 16l96 0c8.8 0 16-7.2 16-16l0-96c0-8.8-7.2-16-16-16l-96 0z"),
		),
	)
}
```

**File: `calendarslider.go`** — add the icon's id constant next to
`maxMonths` (around line 69):

```go
const iconCalendar = svg.Icon("cs-calendar")
```

Add `"webtyp.com/svg"` to this file's import block (not present
today). Also add `"webtyp.com/fmt/lang"`, needed by
`buildCollapsed` in Stage 3.3 below — if Stage 1 already merged, the import
is there already from its Stage 1.2; do not add it twice.

## Stage 3.2 — per-instance id (two calendars on one page must not collide)

**File: `calendarslider.go`** — mirror `selectsearch.go`'s `uid` pattern
exactly (its comment at lines ~153-158 explains why: two pickers whose
checkbox/label share a fixed id would toggle each other).

Add a package-level counter next to `maxMonths`:

```go
var calendarSliderSeq int

func nextCalendarSliderID() int {
	calendarSliderSeq++
	return calendarSliderSeq
}

const suffixCollapsedToggle = "-collapsed-toggle"
```

Add two fields to the `CalendarSlider` struct, next to `today`:

```go
	uid      string     // per-instance id prefix; two calendars on one page must not collide
	expanded *SignalBool // true = full calendar showing; false = mobile collapsed chip showing
```

In `Init`, set both — guarded the same way `Selected` already is, so a
second `Init()` call (the existing `TestCalendarSlider_InitTwiceSafe`
contract) does not hand out a fresh `uid` (which would desync the
checkbox's `id`/`for` pair from whatever the DOM already has) or silently
re-expand a calendar the user had collapsed:

```go
func (c *CalendarSlider) Init(_ Ctx) {
	if c.Selected == nil {
		c.Selected = NewString("")
	}
	if c.expanded == nil {
		c.expanded = NewBool(true)
	}
	if c.uid == "" {
		c.uid = fmt.Sprintf("%s-%d", string(NameCalendarSlider), nextCalendarSliderID())
	}
	c.today = time.FormatDate(time.Now())
}
```

(`fmt` is already imported in this file.)

## Stage 3.3 — the collapsed chip

**File: `calendarslider.go`** — add the four new part constants next to the
existing ones:

```go
	PartCollapsed       = widget.Part("collapsed")
	PartCollapsedToggle = widget.Part("collapsed-toggle")
	PartCollapsedIcon   = widget.Part("collapsed-icon")
	PartCollapsedText   = widget.Part("collapsed-text")
```

and their class vars next to the existing ones:

```go
	clsCollapsed       = NameCalendarSlider.Class(PartCollapsed)
	clsCollapsedToggle = NameCalendarSlider.Class(PartCollapsedToggle)
	clsCollapsedIcon   = NameCalendarSlider.Class(PartCollapsedIcon)
	clsCollapsedText   = NameCalendarSlider.Class(PartCollapsedText)
```

Add a new method, near `buildMonth`:

```go
// buildCollapsed builds the mobile-only compact chip shown once a day is
// picked: a hidden checkbox + <label> — selectsearch's own open/close
// idiom (see its PartHeader/toggle, selectsearch.go ~213-261) — toggles
// c.expanded with the same two-way sync, no JS beyond that one listener.
// Visible only within css.Mobile and only while c.expanded is false — see
// ApplyCollapse-equivalent rules in css.go (Stage 3.4).
//
// The date text goes through lang.Translate, exactly like the month label
// in PLAN_STAGE_1_ARROW_HOVER_LAYOUT.md's Stage 1.2 and for the identical
// reason: date.WeekdayName/date.MonthName return English canonical names,
// never a hardcoded language — this component registers no dictionary
// itself (see that stage for why), it only asks to translate. With no
// dictionary registered anywhere in the running app this reads "Monday 18
// August 2026"; a host that has registered Spanish (see the "Translation
// registration" note in docs/PLAN.md) sees "Lunes 18 Agosto 2026" with no
// further code change here.
func (c *CalendarSlider) buildCollapsed() *Element {
	toggle := Input("checkbox").Set(clsCollapsedToggle.AsAttr()).
		ID(c.uid + suffixCollapsedToggle).
		BindAttrBool("checked", c.expanded).
		On("change", func(e Event) { c.expanded.Set(e.TargetChecked()) })

	text := Span().Set(clsCollapsedText.AsAttr()).
		BindTextFunc(func() string {
			y, m, d := date.ParseDateKey(c.Selected.Get())
			if y == 0 {
				return ""
			}
			weekday := date.WeekdayName(date.Weekday(y, m, d))
			return lang.Translate(weekday, d, date.MonthName(m), y).String()
		})

	label := Label().Set(clsCollapsed.AsAttr()).
		Attr("for", c.uid+suffixCollapsedToggle).
		BindState(widget.Open, DeriveBool(func() bool { return !c.expanded.Get() })).
		Child(iconCalendar.Render(string(clsCollapsedIcon))).
		Child(text)

	return Div().Child(toggle).Child(label)
}
```

`label`'s own `widget.Open` state means "I, the collapsed row, should be
showing" — the INVERSE of `c.expanded` — bound here as its own derived
boolean; this does not conflict with the strip's `BindState(widget.Open,
c.expanded)` added in Stage 3.4 below, because each element's `data-open`
attribute is independent — reusing the state NAME `Open` on two different
elements with two different (here, inverse) source booleans is the same
pattern `selectsearch` itself is not put in a position to need, but nothing
in this framework requires the same physical bool behind every element that
uses the state name `Open`.

## Stage 3.4 — wire it into `Render` and `buildDay`, and the CSS

**File: `calendarslider.go`**

In `Render`, give the strip the `Open` state and add the collapsed chip
as the root's first child:

```go
func (c *CalendarSlider) Render() *Element {
	n := c.numMonths()
	sy, sm := c.startYearMonth()

	keys := make([]string, n)
	for i := 0; i < n; i++ {
		y, m := date.AddMonths(sy, sm, i)
		keys[i] = date.MonthKey(y, m)
	}

	strip := Div().Set(clsStrip.AsAttr()).
		Attr("role", "grid").
		Attr("aria-label", "Calendario").
		BindState(widget.Open, c.expanded)
	for i, key := range keys {
		y, m := date.ParseMonthKey(key)
		prevKey := keys[(i-1+n)%n]
		nextKey := keys[(i+1)%n]
		strip.Child(c.buildMonth(y, m, prevKey, nextKey /* plus Stage 2's two bools, if that stage already landed — see docs/PLAN.md for the merge order the human executing these picked */))
	}

	return Div().Set(clsRoot.AsAttr()).
		Child(c.buildCollapsed()).
		Child(strip)
}
```

`role="grid"` / `aria-label="Calendario"` move from the outer root div onto
`strip` — the root itself no longer always represents a grid once it can
show a plain label/chip instead; `strip` is the element that is actually the
grid whenever it is showing.

In `buildDay`, the existing "selectable" click handler gains one line
(shown in context — only the added line is new):

```go
	if selectable {
		li.On("click", func(Event) {
			c.Selected.Set(dateStr)
			c.expanded.Set(false)
			if c.OnSelect != nil {
				c.OnSelect(dateStr)
			}
			if c.onFilter != nil {
				c.onFilter(dateStr)
			}
		})
	}
```

**File: `css.go`** — add:

```go
Part(PartCollapsed,
    style.Row(style.Space2),
    style.CenterContent(),
    style.Pad(style.Space2),
    style.As(style.Panel),
    style.Round(style.RadiusSm),
    style.Interactive(style.Subtle),
).
Part(PartCollapsedToggle,
    style.Hide(),
).
Part(PartCollapsedIcon,
    style.IconBox(style.IconMd),
).
Part(PartCollapsedText,
    style.Grow(),
    style.FontSize(style.TextSm),
    style.FontWeight(style.WeightBold),
).
On(css.Mobile, PartStrip,
    style.RevealedBy(widget.Open),
).
OnlyOn(css.Mobile, PartCollapsed,
    style.RevealedBy(widget.Open),
).
```

`PartCollapsedToggle`'s `Hide()` is unconditional (matches
`selectsearch.go`'s own hidden-native-checkbox treatment — the checkbox
itself never paints, the `<label>` is the whole visible/clickable surface).

## Acceptance criteria

- `grep -n "PartCollapsed\b" calendarslider.go css.go` → matches in both.
- `grep -n "func.*IconSvg" svg.go` → the new file exists and defines it.
- `gotest -tinygo` green.
- Manual check once built, at a mobile viewport width: pick a day → the
  calendar collapses to one row reading e.g. "Monday 18 August 2026" (or,
  in whatever app has registered a Spanish dictionary and activated it via
  `lang.OutLang(lang.ES)` — see `docs/PLAN.md`'s translation note — "Lunes
  18 Agosto 2026") with a calendar glyph on its left; tapping that row
  re-expands the full calendar; picking a different day collapses it again.
  At a desktop viewport width, none of this is visible at any point — the
  full calendar stays exactly as it is today regardless of selection.

## Stages

| Stage | File(s) | Done when |
|---|---|---|
| 3.1 | `svg.go` (new) | `IconSvg` compiles, ships the calendar glyph |
| 3.2 | `calendarslider.go` | `uid`/`expanded` set in `Init`, no collisions across instances |
| 3.3 | `calendarslider.go` | `buildCollapsed` compiles, checkbox+label toggles `expanded` |
| 3.4 | `calendarslider.go`, `css.go` | full `gotest -tinygo` green; manual mobile/desktop check above passes |
