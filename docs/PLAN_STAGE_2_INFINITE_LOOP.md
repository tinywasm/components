---
PLAN: "feat(calendarslider): instant jump at the wrap boundary instead of a long smooth scroll the wrong way"
REVIEWER: none
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Part of a 3-stage calendarslider improvement — see `docs/PLAN.md` for the
> index.

## ⚠️ BLOCKING PREREQUISITE — do not dispatch this stage until it is satisfied

This stage calls `dom.Reference.ScrollIntoViewInstant()`, which does not
exist in any published `webtyp/dom` version yet. It is added by
`https://github.com/webtyp/dom/blob/main/docs/PLAN.md` (a separate,
already-written, self-contained plan in that repo). Before dispatching THIS
stage:

1. That `dom` plan must be dispatched, merged, and published (new `dom`
   version tagged).
2. `components/go.mod` must be bumped to that new version:
   `go get -u webtyp.com/dom@latest` from `components/`, then
   `go mod tidy`.
3. Confirm `grep -n "ScrollIntoViewInstant" $(go env GOPATH)/pkg/mod/webtyp.com/dom@*/reference.go`
   finds it before writing a single line of this stage.

If dispatched before that, every step below fails to compile — this is not
a design choice the executor can route around; stop and wait for the
prerequisite.

# Plan — CalendarSlider Stage 2: instant jump at the wrap edge

## Context

The strip is already logically circular: the ‹ of the first month's card
targets the key of the last month, and the › of the last targets the first
(see `Render`'s `prevKey := keys[(i-1+n)%n]` / `nextKey := keys[(i+1)%n]`).
The bug is not the data, it is the ANIMATION: `slideToMonth` always calls
`ScrollIntoView()`, which is always smooth. At the wrap boundary that means
scrolling smoothly across every month card in between, in the direction
OPPOSITE to what the pressed arrow implies (pressing ‹ on the first month
scrolls the strip rightward, through every other month, to reach the last
one) — that long, backwards-reading scroll is what reads as "reloading from
the start".

Every OTHER navigation (adjacent months) already behaves correctly and stays
untouched. The fix is narrowly targeted: only the two wrap edges (first
month's ‹, last month's ›) jump instantly; everything else keeps the smooth
`ScrollIntoView()` it already has.

This does not change how many months are in the DOM, does not add any new
elements, and does not touch `css.go` at all — it is a pure behavior change
in three functions.

## Stage 2.1 — thread "is this a wrap" through to the click handlers

**File: `calendarslider.go`**

1. Change `slideToMonth`'s signature (currently a single-line function
   calling `Get` + `ScrollIntoView`):

   ```go
   // slideToMonth jumps the scroll-snap strip to the month card carrying the
   // given key. instant selects ScrollIntoViewInstant over the normal smooth
   // ScrollIntoView — reserved for the two wrap edges (first month's ‹, last
   // month's ›), where a smooth scroll would visibly travel across every
   // month in between in the wrong apparent direction. Every adjacent-month
   // navigation keeps calling this with instant=false.
   func slideToMonth(key string, instant bool) {
       ref, ok := Get("cs-m-" + key)
       if !ok {
           return
       }
       if instant {
           ref.ScrollIntoViewInstant()
           return
       }
       ref.ScrollIntoView()
   }
   ```

2. Change `buildMonth`'s signature to accept the two wrap flags, and use them
   in the two `On("click", ...)` handlers already there:

   ```go
   func (c *CalendarSlider) buildMonth(year, month int, prevKey, nextKey string, prevWraps, nextWraps bool) *Element {
   ```

   ```go
   prev.On("click", func(Event) { slideToMonth(prevKey, prevWraps) })
   ```

   ```go
   next.On("click", func(Event) { slideToMonth(nextKey, nextWraps) })
   ```

   Nothing else in `buildMonth` changes — this is additive parameters
   threaded into the two existing closures, not new UI.

3. In `Render`, compute the two flags per iteration and pass them:

   ```go
   strip := Div().Set(clsStrip.AsAttr())
   for i, key := range keys {
       y, m := date.ParseMonthKey(key)
       prevKey := keys[(i-1+n)%n]
       nextKey := keys[(i+1)%n]
       prevWraps := i == 0
       nextWraps := i == n-1
       strip.Child(c.buildMonth(y, m, prevKey, nextKey, prevWraps, nextWraps))
   }
   ```

   When `n == 1` both flags are true for the single month (its own prev AND
   next target itself) — that already falls out of `i==0 && i==n-1` with no
   special-casing needed; clicking either just re-jumps to the same card,
   harmless.

## Stage 2.2 — update the 6 existing direct calls to `buildMonth`

`buildMonth` is called directly (not through `Render`) in 6 places across
`calendarslider_test.go` — all need the two new trailing bool arguments.
None of them are testing wrap timing (that is new, covered in Stage 2.3), so
pass whatever reflects what each test is actually about:

| Line (today) | Call | New trailing args |
|---|---|---|
| 25 (`TestBuildMonthAgosto2026`) | `c.buildMonth(2026, 8, "2026-07", "2026-09")` | `, false, false` |
| 53 (`TestBuildMonthAlwaysHasBothLinks`) | `c.buildMonth(2026, 8, "2026-07", "2026-09")` | `, false, false` |
| 76 (`TestRenderWrapsAround`, `augMonth`) | `c.buildMonth(2026, 8, "2026-10", "2026-09")` | `, true, false` — this call specifically represents agosto (first month), whose prev DOES wrap |
| 81 (`TestRenderWrapsAround`, `octMonth`) | `c.buildMonth(2026, 10, "2026-09", "2026-08")` | `, false, true` — represents octubre (last month), whose next DOES wrap |
| 94 (`TestBuildMonthPadsToSixWeeks`) | `c.buildMonth(2021, 2, "2021-01", "2021-03")` | `, false, false` |
| 131 (`TestBuildMonthCells`) | `c.buildMonth(2026, 8, "2026-07", "2026-09")` | `, false, false` |

Re-run `grep -n "buildMonth(" calendarslider*.go` after editing — every call
site must now pass 6 arguments; if any still passes 4, the package will not
compile (which is a fine, loud way to confirm you found them all).

## Stage 2.3 — new test proving the actual browser call, not just the wiring

`calendarslider_test.go`'s `TestRenderWrapsAround` only proves the `href`/
`data-target` data is a loop — it says nothing about the animation, because
that is a live-browser behavior invisible in a `.String()` HTML dump. Add a
new WASM test that intercepts the real `Element.prototype.scrollIntoView`
call and asserts the `behavior` option differs between a wrap click and an
adjacent click.

**File: `calendarslider_wasm_test.go`** — add:

```go
// TestWrapNavigationJumpsInstantly cubre el bug real: el ‹ del primer mes
// (envuelve al último) y el › del último (envuelve al primero) deben saltar
// sin animación — un scroll suave ahí viaja visualmente en la dirección
// contraria a través de todos los meses intermedios, lo que se lee como un
// reinicio. La navegación entre meses vecinos sigue siendo suave.
func TestWrapNavigationJumpsInstantly(t *testing.T) {
	c := &CalendarSlider{Start: "2026-08", NumMonths: 3}
	c.Init(nil)
	Render("app", c.Render())

	var lastBehavior string
	proto := js.Global().Get("Element").Get("prototype")
	original := proto.Get("scrollIntoView")
	spy := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) > 0 {
			lastBehavior = args[0].Get("behavior").String()
		}
		return nil
	})
	proto.Set("scrollIntoView", spy)
	defer func() {
		proto.Set("scrollIntoView", original)
		spy.Release()
	}()

	query(t, "#cs-m-2026-08 .calendarslider__prev").Call("click")
	if lastBehavior != "instant" {
		t.Errorf("wrap prev (agosto -> octubre) behavior = %q, want instant", lastBehavior)
	}

	query(t, "#cs-m-2026-08 .calendarslider__next").Call("click")
	if lastBehavior != "smooth" {
		t.Errorf("adjacent next (agosto -> septiembre) behavior = %q, want smooth", lastBehavior)
	}

	query(t, "#cs-m-2026-10 .calendarslider__next").Call("click")
	if lastBehavior != "instant" {
		t.Errorf("wrap next (octubre -> agosto) behavior = %q, want instant", lastBehavior)
	}
}
```

This uses the same `query` helper already defined in this file (top of
`calendarslider_wasm_test.go`) — no new test scaffolding needed.

## Acceptance criteria

- `grep -n "func slideToMonth" calendarslider.go` shows the 2-argument form.
- `grep -c "buildMonth(" calendarslider.go calendarslider_test.go` — every
  call site (7 total: 1 in `Render`, 6 in tests) passes 6 arguments; `go
  vet ./...` catches any that don't.
- `gotest -tinygo` green, including the new `TestWrapNavigationJumpsInstantly`.
- Manual check once built: rapidly clicking › through all months and past
  the last one shows N-1 smooth slides then one instant cut back to the
  first month — never a long scroll backwards through everything.

## Stages

| Stage | File(s) | Done when |
|---|---|---|
| 2.1 | `calendarslider.go` | `slideToMonth`/`buildMonth`/`Render` thread the wrap flags |
| 2.2 | `calendarslider_test.go` | all 6 direct `buildMonth` calls updated, package compiles |
| 2.3 | `calendarslider_wasm_test.go` | `TestWrapNavigationJumpsInstantly` passes under `gotest -tinygo` |
