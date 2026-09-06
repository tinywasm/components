---
PLAN: "feat(calendarslider): reclaim day-grid width, hover-reveal arrows on desktop, compact flanking arrows always visible on mobile"
REVIEWER: none
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Part of a 3-stage calendarslider improvement — see `docs/PLAN.md` for the
> index. This stage has no PUBLISHED-version dependency on the other two —
> it can be dispatched and merged before, after, or independent of them.
> It does, however, touch the SAME month-label line Stage 1.2 below
> describes regardless of order; read that section even though the rest of
> this stage is unrelated to language.

# Plan — CalendarSlider Stage 1: arrow layout

## Context

`PartMonth` (`components/calendarslider/css.go`) reserves
`style.PadInline(style.Space6)` permanently on both sides, so the ‹ › nav
buttons (rendered as full-height `EdgeStrip` overlays) never overlap the
day grid. That gutter shrinks the day grid's real width **at all times**,
mouse or touch, whether the arrows are being looked at or not.

The fix: arrows become a normal, compact, always-visible row flanking the
month name ("‹ September 2026 ›" by default — see Stage 1.2 for why it's
English) — no permanent gutter needed,
the day grid reclaims the full width. On a fine (mouse) pointer only, the
SAME two buttons additionally grow into full-height edge overlays while the
pointer hovers the month card (or while a button itself holds focus, for
keyboard users) — `widget/style` already has the exact primitives for this
gated correctly (`(hover: hover)` media, never touch): `Cue`, `CueWithinHover`.
Verified against a throwaway probe program before writing this plan — the
combination compiles and emits exactly the expected CSS (see "Verified
mechanism" below).

## Verified mechanism (do not re-derive, use as given)

```go
style.For(w).
    Part(partPrev, /* base, always-visible options */).
    Cue(widget.Focus, partPrev, style.EdgeStrip(style.Parent, style.SideStart)).
    CueWithinHover(widget.Hover, partMonth, partPrev, style.EdgeStrip(style.Parent, style.SideStart)).
    Stylesheet()
```

emits (confirmed by running it):

```css
@layer states {
.w__prev:focus-visible { position: absolute; inset-block: 0; inset-inline-start: 0; z-index: 1; }
}
@media (hover: hover) {
@layer states {
.w__card:hover .w__prev { position: absolute; inset-block: 0; inset-inline-start: 0; z-index: 1; }
}
}
```

`:focus-visible`'s rule is NOT inside the media query — keyboard focus reveals
the edge overlay on every device, mouse or not; that is correct and intended
(the base state already has the button visible and reachable everywhere, so
this only additionally *repositions/enlarges* it while focused — an instant
snap, not an animated move, since `position` cannot be transitioned). Only
the hover-triggered rule is gated to `(hover: hover)`, so it never fires on a
touch tap. This is the DELIBERATE reason to use `Cue`/`CueWithinHover`, not a
hand-rolled `:hover` — no risk of a "stuck hover" on mobile.

## Stage 1.1 — restructure `buildMonth`'s markup

**File: `calendarslider.go`**

Add a new part constant, next to the existing ones (around line 26-42):

```go
	PartMonthNav      = widget.Part("month-nav")
```

Add its class var, next to the others (around line 44-61):

```go
	clsMonthNav = NameCalendarSlider.Class(PartMonthNav)
```

In `buildMonth` (currently lines 241-297), replace the three separate
`monthEl.Child(...)` calls for the month-name label, `prev`, and `next` —
today they are three independent children of `monthEl` appended in this
order: label, then `prev`, then `next`. Wrap all three in one new
`PartMonthNav` container instead, in the order `prev, label, next` (so the
label sits between the two arrows once they're a compact flanking row):

```go
	monthName := Div().Set(clsMonthNm.AsAttr()).
		Text(lang.Translate(date.MonthName(month), year).String())

	prev := Button().Set(clsPrev.AsAttr()).
		Attr("type", "button").
		Attr("aria-label", "Mes anterior").
		Attr("title", "Mes anterior").
		Attr("data-target", "cs-m-"+prevKey).
		Text("‹")
	prev.On("click", func(Event) { slideToMonth(prevKey) })

	next := Button().Set(clsNext.AsAttr()).
		Attr("type", "button").
		Attr("aria-label", "Mes siguiente").
		Attr("title", "Mes siguiente").
		Attr("data-target", "cs-m-"+nextKey).
		Text("›")
	next.On("click", func(Event) { slideToMonth(nextKey) })

	monthEl.Child(Div().Set(clsMonthNav.AsAttr()).
		Child(prev).
		Child(monthName).
		Child(next))
```

Do NOT touch `slideToMonth`'s signature or body in this stage — that is
Stage 2's job (`docs/PLAN_STAGE_2_INFINITE_LOOP.md`), which depends on a
`webtyp/dom` addition not yet published. Keep the click handlers exactly
as they are today, just relocated into this new wrapper.

`monthEl`'s own children are now, in order: the weekday row, the (up to 6)
week rows, then this one `PartMonthNav` div — 8 children total for a normal
month card (weekday row + 6 week rows + 1 nav), where today it is 10
(weekday row + 6 week rows + label + prev + next as three separate
children). This changes two existing tests — see Stage 1.4.

## Stage 1.2 — the month label goes through `fmt/lang`, not raw

`date.MonthName` used to return the Spanish name directly ("Agosto") —
`components/docs/PLAN_STAGE_2_INFINITE_LOOP.md`'s sibling repo plan
(`https://github.com/webtyp/date/blob/main/docs/PLAN.md`) changes that to
the English canonical name ("August"): **this library never hardcodes a
human language, the consumer decides** — the same rule
`layout/crudview` already follows for every string it renders (see
`https://github.com/webtyp/layout/blob/main/docs/DICTIONARY.md`). The
snippet above already reflects the fix: `lang.Translate(date.MonthName(month),
year).String()` instead of `fmt.Sprintf("%s %d", date.MonthName(month), year)`.

Add the import to `calendarslider.go`'s import block:

```go
"webtyp.com/fmt/lang"
```

This call is **safe to make regardless of dispatch order** against the
`date` plan or `components/go.mod`'s version of it:

- Old `date` (still Spanish) → `lang.Translate("Agosto", 2026)` — "Agosto"
  is an unrecognized dictionary key, so it passes through unchanged, exactly
  as `fmt.Sprintf` would have rendered it. No behavior change yet.
- New `date` (English) with NO consumer dictionary registered anywhere →
  renders "August 2026" — English, per `fmt/lang`'s own documented default
  ("without any dictionary, everything renders in English").
  `calendarslider` itself registers nothing — per the same rule, a
  reusable component is not the place to decide the end user's language,
  same as `crudview` never registers Spanish for "Confirm"/"Cancel" itself.
  A REAL application embedding `calendarslider` that wants "Agosto"
  registers `lang.RegisterWords([]lang.DictEntry{{EN: "August", ES:
  "Agosto"}, ...})` for the 12 months once, itself — exactly the pattern
  `layout/docs/DICTIONARY.md` documents for crudview's own strings. That
  registration is OUT OF SCOPE for this stage (it belongs to whichever real
  app wants Spanish output, `app-demo` included, as its own follow-up) —
  do not add it here.
- New `date` WITH a consumer dictionary registered somewhere in the running
  app → renders "Agosto 2026" (or whatever language that app activated via
  `lang.OutLang`), with no further code change needed here.

Verified against a real `go run` before writing this plan — `lang.Translate`
mixing a translated word with a raw `int` year produces exactly the
expected spacing ("Lunes 18 Agosto 2026"-shaped output; see Stage 3's plan
for the full natural-language date, which uses the identical mechanism).

## Stage 1.3 — CSS

**File: `css.go`**

1. Remove `style.PadInline(style.Space6)` from `Part(PartMonth, ...)`. Leave
   every other option on that Part call untouched
   (`As(Panel)`, `Pad(Space2)`, `Anchor()`, `Width(Full)`, `Stack(Space1)`).
   `Anchor()` MUST stay — it is still the containing block the hover-revealed
   `EdgeStrip` positions against.

2. Add a new `Part(PartMonthNav, ...)` block — a centered, compact row:

   ```go
   Part(PartMonthNav,
       style.Row(style.Space2),
       style.CenterContent(),
   ).
   ```

3. Change `Part(PartPrev, ...)` and `Part(PartNext, ...)`: remove
   `style.EdgeStrip(style.Parent, style.SideStart)` /
   `style.EdgeStrip(style.Parent, style.SideEnd)` from these base blocks —
   that positioning now belongs ONLY in the reveal rules added in step 4.
   Do NOT add a transition/`Animate` here: a `position` change (static →
   absolute) is not something CSS can interpolate, so there is nothing for a
   transition to smooth — the reveal is, and stays, an instant snap on
   hover/focus. That is expected and fine at hover speed; do not try to
   animate it.
   Keep every other option on both (`As(Bare)`, `CenterContent()`,
   `FontSize(Text2xl)`, `FontWeight(WeightBold)`, `Pad(Space1)`,
   `Interactive(Subtle)`, `Round(RadiusSm)`) exactly as they are — this is
   now their permanent, always-visible, in-flow appearance (mobile and
   desktop-before-hover alike).

4. Add the reveal rules (new, at the end of the builder chain, before
   `.Stylesheet()`):

   ```go
   Cue(widget.Focus, PartPrev, style.EdgeStrip(style.Parent, style.SideStart)).
   Cue(widget.Focus, PartNext, style.EdgeStrip(style.Parent, style.SideEnd)).
   CueWithinHover(widget.Hover, PartMonth, PartPrev, style.EdgeStrip(style.Parent, style.SideStart)).
   CueWithinHover(widget.Hover, PartMonth, PartNext, style.EdgeStrip(style.Parent, style.SideEnd)).
   ```

   Place them right after the existing `When(widget.Selected, PartDay, ...)`
   block and before the `On(css.Mobile, PartDay, ...)` block, or anywhere
   else in the chain — order among independent `Cue`/`When`/`On` calls does
   not matter to the emitted output.

Net result: the day grid (weekday row + week rows) always gets `PartMonth`'s
full content width — no side gutter, on any device. The nav row is compact
and always visible everywhere. Only on a fine pointer does hovering the
card (or focusing a button) additionally turn that button into a full-height
edge overlay.

## Stage 1.4 — fix the two tests this markup change breaks

**File: `calendarslider_test.go`**

Both `TestBuildMonthAgosto2026` and `TestBuildMonthPadsToSixWeeks` assert
`len(children) != 10` (must become `!= 8`) and reach into `children[8]` /
`children[9]` for the prev/next buttons, which no longer exist as separate
indices — both buttons are now inside `children[7]` (the new
`PartMonthNav` div) alongside the label.

In `TestBuildMonthAgosto2026` (currently lines 23-40): change
`len(children) != 10` to `!= 8`; the fatal message's "10" → "8" and its
"etiqueta + prev + next" → "etiqueta + prev + next (en una sola fila)" or
similar. Change the two `children[8].String()` / `children[9].String()`
checks to both read `children[7].String()` instead (delete one of the two
now-duplicate `Contains(..., "<button")` checks or keep both against the
same string — they both still hold, since `children[7]` now serializes prev,
the label, and next together). The `children[7]` "Agosto 2026" check on line
35 needs no change — that string is still a substring of `children[7]`'s
serialization once nav is index 7, just update its comment ("el mes lleva la
etiqueta" → "el mes lleva la fila de navegación, con la etiqueta en el
medio") so it does not claim children[7] is only the label.

In `TestBuildMonthPadsToSixWeeks` (currently lines 92-120): same `10`→`8`
child-count fix. Its `children[7]` "Febrero 2021" check (line 117) needs no
code change (same reasoning), only the preceding comment
("La etiqueta sigue siendo el hijo 7") updated to say the nav row (label +
both arrows) is hijo 7, not the label alone.

**File: `calendarslider_wasm_test.go`, `calendarslider_contract_test.go`,
`calendarslider_test.go`'s other tests** (`TestBuildMonthCells`,
`TestClampOccupation`, `TestRenderStructure`, `TestPairMarkupAndStylesheet`,
`TestNumMonthsClampsToMax`, `TestNumMonthsDefaultsToThree`,
`TestCalendarSlider_SatisfiesFilterable`, everything in
`calendarslider_wasm_test.go`, everything in
`calendarslider_contract_test.go`): read through them, but expect NO changes
— they query by CSS class (`.calendarslider__prev`, `#cs-m-... .calendarslider__month-name`,
etc.) or by content substring, never by child index beyond `Children()[1]`
for the first week row (index 1, still unaffected — only indices 7+ moved).
Confirm this by actually running the suite (next section) rather than
skipping a file because this plan says so.

One exception worth double-checking by eye, not by assumption: in
`TestRenderStructure` (`calendarslider_test.go`), the comment above the
`position: absolute` check says "Las flechas van superpuestas (EdgeStrip) en
la hoja de estilos servida" — true before this stage, now only true while
hovering/focused. The assertion itself
(`strings.Contains(cssOut, "position: absolute;")`) still passes unchanged,
because the `:focus-visible` reveal rule emits `position: absolute;`
unconditionally (outside any `@media` block — see "Verified mechanism"
above) — but update the comment so it says the overlay is conditional
(hover/focus), not permanent.

## Acceptance criteria

- `grep -n "PartMonthNav" calendarslider.go css.go` → matches in both files.
- `grep -n "EdgeStrip" css.go` → 0 matches inside the base `Part(PartPrev`/
  `Part(PartNext` blocks, exactly 4 matches total (2 `Cue(Focus,...)` + 2
  `CueWithinHover(Hover,...)`).
- `grep -n "PadInline" css.go` → 0 matches.
- `grep -n "fmt/lang" calendarslider.go` → one match (the import); `grep -n "date.MonthName" calendarslider.go` → still one match, now wrapped inside a `lang.Translate(...)` call, never concatenated raw via `fmt.Sprintf` again.
- `gotest` (stdlib **and** `-tinygo`) green, all packages in `components`.
- Manually confirm in a real browser (or the project's `run` skill) once
  built: on a desktop-width viewport, hovering a month card reveals full-height
  ‹ › at its edges; on a narrow/touch viewport, ‹ › sit permanently flanking
  the month name, no gutter around the day grid either way.

## Stages

| Stage | File(s) | Done when |
|---|---|---|
| 1.1 | `calendarslider.go` | `buildMonth` emits one `PartMonthNav` div wrapping prev/label/next |
| 1.2 | `calendarslider.go` | month label renders via `lang.Translate`, not raw `date.MonthName` |
| 1.3 | `css.go` | gutter gone, nav row compact by default, edge-overlay only on hover/focus |
| 1.4 | `calendarslider_test.go` | full `gotest` suite green |
