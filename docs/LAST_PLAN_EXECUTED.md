---
PLAN: "feat(calendarslider): arrow layout, infinite loop, mobile collapse"
REVIEWER: none
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.

# PLAN — index for 3 independent CalendarSlider improvements

Three separate requests against `components/calendarslider`, different in
nature (a layout/CSS fix, a navigation-behavior fix, a new interaction) —
split into 3 stage files on purpose, so each can be assigned to whatever
executor its own complexity calls for, and dispatched on its own schedule.
They do **not** need to run as one sequential pass through this index the
way a single-topic multi-stage plan would; **pick a stage file, promote it
to the thing you dispatch (or point your executor straight at it), and
repeat for the others independently.**

| Stage | File | What | Blocked on |
|---|---|---|---|
| 1 | [PLAN_STAGE_1_ARROW_HOVER_LAYOUT.md](PLAN_STAGE_1_ARROW_HOVER_LAYOUT.md) | Day grid reclaims the width the ‹ › gutter used to reserve permanently; arrows become a compact row flanking the month name (always visible, any device) and additionally grow into full-height edge overlays on mouse-hover/keyboard-focus only | Nothing — can run first, alone, anytime (it also fixes the month label's now-stale raw-Spanish rendering, see its Stage 1.2, but that fix is safe against either version of `webtyp/date`) |
| 2 | [PLAN_STAGE_2_INFINITE_LOOP.md](PLAN_STAGE_2_INFINITE_LOOP.md) | The wrap edge (last→first, first→last) jumps instantly instead of smooth-scrolling backwards across every month in between | `webtyp/dom`'s `PLAN.md` (below) must be merged + published + version-bumped in `components/go.mod` first |
| 3 | [PLAN_STAGE_3_MOBILE_COLLAPSE.md](PLAN_STAGE_3_MOBILE_COLLAPSE.md) | On mobile only, picking a day collapses the calendar to a compact chip ("Monday 18 August 2026" by default); tapping it re-expands | `webtyp/date`'s `PLAN.md` (below) must be merged + published + version-bumped in `components/go.mod` first |

## External prerequisites (separate repos, separate plans, already written)

Two of the three stages need one small, additive capability each from a
foundational repo `components` depends on. Both are self-contained plans,
already written, sitting in their own repo:

- `https://github.com/webtyp/dom/blob/main/docs/PLAN.md` — adds
  `Reference.ScrollIntoViewInstant()`. Pure addition. Gates Stage 2.
- `https://github.com/webtyp/date/blob/main/docs/PLAN.md` — adds
  `WeekdayName` and `ParseDateKey`, **and changes `MonthName`'s existing
  return values from Spanish to English** — a breaking change to
  already-shipped behavior, deliberate (see "Translation registration"
  below for why), not a pure addition like the `dom` one. Gates Stage 3 —
  and Stage 3's own prerequisite section lists the 3 existing tests this
  breaks and must fix in the same commit as the version bump.

Dispatch and merge these two first (they have no dependency on each other
or on any calendarslider stage), then bump `components/go.mod` (`go get -u
github.com/webtyp/dom@latest github.com/webtyp/date@latest && go mod
tidy`) before dispatching Stage 2 or Stage 3. Stage 1 needs neither and has
no reason to wait for them.

## Translation registration — no stage in `components` does this, on purpose

`date.MonthName`/`date.WeekdayName` return English — the canonical,
untranslated form. `calendarslider` (Stages 1 and 3) renders them through
`lang.Translate(...)` (`webtyp.com/fmt/lang`) but registers **no**
dictionary itself: *"the library never registers words, so you decide the
language your users see"* —
`https://github.com/webtyp/layout/blob/main/docs/DICTIONARY.md`,
the same rule `layout/crudview` already follows for its own strings
("Confirm", "Cancel", "Delete", …). With no dictionary registered anywhere
in the running app, everything renders in English — that is the correct
default after these 3 stages land, not a bug.

A real application that wants Spanish (or any of the other 7 languages
`fmt/lang` supports) registers these 19 words itself, once, the same way
`layout/docs/DICTIONARY.md` shows for crudview's words:

```go
import "webtyp.com/fmt/lang"

func init() {
	lang.RegisterWords([]lang.DictEntry{
		{EN: "January", ES: "Enero"}, {EN: "February", ES: "Febrero"},
		{EN: "March", ES: "Marzo"}, {EN: "April", ES: "Abril"},
		{EN: "May", ES: "Mayo"}, {EN: "June", ES: "Junio"},
		{EN: "July", ES: "Julio"}, {EN: "August", ES: "Agosto"},
		{EN: "September", ES: "Septiembre"}, {EN: "October", ES: "Octubre"},
		{EN: "November", ES: "Noviembre"}, {EN: "December", ES: "Diciembre"},
		{EN: "Sunday", ES: "Domingo"}, {EN: "Monday", ES: "Lunes"},
		{EN: "Tuesday", ES: "Martes"}, {EN: "Wednesday", ES: "Miércoles"},
		{EN: "Thursday", ES: "Jueves"}, {EN: "Friday", ES: "Viernes"},
		{EN: "Saturday", ES: "Sábado"},
	})
}

func main() {
	lang.OutLang(lang.ES) // or lang.OutLang() to auto-detect
	// ...
}
```

This registration is **out of scope for all 3 stages above** — it belongs
to whichever real application embeds `calendarslider` (`app-demo` included,
if its own demo is meant to keep reading in Spanish once these stages and
their prerequisites land). None of the 3 stage files add it, and none
should.

## The one real ordering constraint: Stages 2 and 3 both touch `Render()`

Stage 1 only touches `buildMonth`'s internals (markup) and `css.go` — it
cannot conflict with the other two. Stages 2 and 3 both edit
`CalendarSlider.Render()` in `calendarslider.go` (Stage 2 adds two bool
params threaded through the per-month loop; Stage 3 wraps the strip and adds
a sibling). They are logically independent of each other, but dispatching
both as separate PRs from the same starting commit risks a merge conflict
in that one function. Finish and merge one (2 or 3, whichever order — no
preference) before starting the other's PR; do not run them concurrently
from the same base commit.

## Verification once all 3 (plus both prerequisites) have landed

```
gotest -tinygo
```

green, from `components/`, plus the manual per-stage checks each stage file
lists in its own "Acceptance criteria" section (hover-reveal on a real
desktop browser, wrap-jump feel, mobile collapse/expand) — these are visual/
interaction checks a test suite cannot fully cover on its own.

## Execution log (2026-09-05, single sequential pass: Stage 1 → 2 → 3)

- Both external prerequisites were already published at dispatch time:
  `dom` v0.13.9 ships `ScrollIntoViewInstant`; `date` v0.0.5 ships
  `WeekdayName`/`ParseDateKey` with English `MonthName`. `components/go.mod`
  bumped `date` v0.0.2 → v0.0.5 (`go get -u + go mod tidy`) before Stage 3.
- Stage 1 + prerequisite test fixes: `TestBuildMonthPadsToSixWeeks`
  ("Febrero 2021" → "February 2021") needed the same English fix as the 3
  assertions Stage 3's prerequisite lists — the stage file omits it, but the
  old-Spanish assertion fails identically under `date` v0.0.5.
- **Deviation from Stage 3 as written:** the stage prescribes
  `BindState(widget.Open, …)` + `RevealedBy(widget.Open)` for the strip and
  the collapsed chip. `CalendarSlider.WidgetKind()` is `Grid`, and
  `widget.Kind.Allows` permits `Grid` only `Selected`/`Current` — the root
  `conformance_test.go` rejects `Open` on a grid (`gotest` full suite red).
  Implemented with `widget.Current` instead (both `BindState`s and both
  `RevealedBy`s); each element's `data-current` stays independent (strip
  follows `expanded`, chip its negation), same mechanics as specified.
- Verification: `gotest` (vet ✅ race ✅ tests ✅ wasm ✅) and
  `gotest -tinygo` green from `components/`; sprite-leak check
  (`go list -deps ./... | grep webtyp/svg/sprite`) empty. Manual
  browser checks (hover reveal, wrap-jump feel, mobile collapse/expand)
  still pending — they need a real viewport, not covered here.
