---
PLAN: "fix: calendarslider nav must not touch location.hash; move list-selection chrome into listselect"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 8064228250952934631
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Do NOT run `gopush` or `codejob`; they are developer tools outside this repo.

# PLAN — `tinywasm/components`

Two independent fixes, both found testing the `#reservation` view of the demo
app in the browser. Do **Part 1** and **Part 2** in order; each is
self-contained; finish one (its acceptance criteria green) before the next.
Never mix changes between them.

Root-cause write-up (context, not instructions):
`../RESERVATION_VIEW_FIXES_MASTER_PLAN.md` in the monorepo — but everything you
need is inline here.

---

## Repo rules you MUST follow (from `AGENTS.md` + `CONSTRUCTION_HARNESS.md`)

- **Public library → everything in English**: code, comments, identifiers,
  error messages. (Only Spanish that already exists in test comments may stay.)
- **No Go stdlib in files that compile to WASM.** Use `github.com/tinywasm/fmt`
  (`fmt.Sprintf`, `fmt.Errf`). No `errors`, `strconv`, `strings`. `tinywasm/fmt`
  has **no** `Itoa` — use `fmt.Sprintf("%d", n)`.
- **No `map[...]` in WASM-compiled code.** Linear scan over a slice.
- **Embed `dom.Element` by value**, never `*dom.Element`.
- **Build-tag split by file**: runtime/`Render()` code is untagged; CSS lives in
  `css.go` under `//go:build !wasm`; SVG geometry in `svg.go` under
  `//go:build !wasm`. Never create `front.go`.
- **Component contract is `Render()` + optional `Init(ctx)` only.** No
  `OnMount`/lifecycle hooks. One-time work goes in `Init`.
- **Reactivity is signals only** (`SignalString`/`SignalBool`/`SignalNodes`,
  `DeriveBool`, `Bind*`/`Bind*Func`). No generics anywhere.
- **A shared piece is assembled, never re-declared.** If two components would
  contain the same block, it belongs in the lego piece that owns the concern
  (`listselect`, `listgap`), exposed as a typed helper — not pasted twice.
- Tests: `gotest`, never `go test`. Stdlib `testing` only. WASM tests run
  against a real DOM. Every new/changed public helper needs a
  **consumer-shaped test in this repo**.
- Pre-publish sprite-leak check must stay clean:
  `GOOS=js GOARCH=wasm go list -deps ./... | grep tinywasm/svg/sprite` → empty.

---

# PART 1 — `calendarslider`: the ‹ › controls must not mutate `location.hash`

## The bug

`calendarslider/calendarslider.go` `buildMonth` renders the month-navigation
controls as page anchors:

```go
monthEl.Child(A("#cs-m-" + prevKey).Set(clsPrev.AsAttr())...)
monthEl.Child(A("#cs-m-" + nextKey).Set(clsNext.AsAttr())...)
```

The package doc (lines 1–9) states this is deliberate: *"the ‹ › controls are
plain same-page anchor links (`<a href="#cs-m-...">`) … the browser's native
scroll-snap does the sliding — no click handler"*.

That assumes the URL hash is free real estate. In `platformd` (the shell the
demo runs in) **the hash is the router**. Clicking ‹ or › sets
`location.hash = "#cs-m-2026-10"`; the router finds no module by that name and
**blanks the whole stage**. Confirmed live: pressing › left `<main>` empty with
`location.hash` = `#cs-m-2026-10`.

## What is already correct — do NOT change it

- `NumMonths` (field, `calendarslider.go:100`) already makes the month count
  configurable. `numMonths()` (`:153`) already applies `0 → 3` and caps at
  `maxMonths` (12). **The demo's "3 months" is already the default. Add no new
  configurability.**
- The wrap-around ("rotate between them") is already built: `Render` (`:198`)
  links `keys[(i-1+n)%n]` as prev and `keys[(i+1)%n]` as next, so the first
  month's ‹ points at the last and the last month's › points at the first.
  **Keep this behaviour exactly.**

The ONLY defect is that the control is an `<a href="#...">` instead of a button
that scrolls.

## Fix — file by file

### `calendarslider/calendarslider.go`

1. **`buildMonth`** — replace the two `A("#cs-m-"+key)` children with buttons.
   Keep every other attribute (`clsPrev`/`clsNext` class, `aria-label`,
   `title`, the `‹` / `›` text). Add `data-target` carrying the destination id
   (needed by the SSR test to assert the wrap-around, and it documents intent):

   ```go
   prev := Button().Set(clsPrev.AsAttr()).
       Attr("type", "button").
       Attr("aria-label", "Mes anterior").
       Attr("title", "Mes anterior").
       Attr("data-target", "cs-m-"+prevKey).
       Text("‹")
   prev.On("click", func(Event) { slideToMonth(prevKey) })
   monthEl.Child(prev)

   next := Button().Set(clsNext.AsAttr()).
       Attr("type", "button").
       Attr("aria-label", "Mes siguiente").
       Attr("title", "Mes siguiente").
       Attr("data-target", "cs-m-"+nextKey).
       Text("›")
   next.On("click", func(Event) { slideToMonth(nextKey) })
   monthEl.Child(next)
   ```

   `Button` is already available (`. "github.com/tinywasm/html"` dot-import).
   `Event` and `Get` are already available (`. "github.com/tinywasm/dom"`).

2. **Add the helper** (same file, package-level, near `buildMonth`):

   ```go
   // slideToMonth jumps the scroll-snap strip to the month card carrying the
   // given key. The ‹ › controls call this instead of being
   // <a href="#cs-m-..."> anchors: an anchor mutates location.hash, and a
   // hash-routed shell (platformd) reads that as a route change, not a scroll.
   // A <button> plus this handler keeps the slide entirely inside the widget;
   // the browser still animates it via scroll-snap. ScrollIntoView is a no-op
   // on the server build.
   func slideToMonth(key string) {
       if ref, ok := Get("cs-m-" + key); ok {
           ref.ScrollIntoView()
       }
   }
   ```

   `dom.Get(id) (Reference, bool)` and `Reference.ScrollIntoView()` already
   exist and are documented for exactly this (jumping a horizontal scroll-snap
   container to another panel) — see `dom/dom.go` and `dom/element_wasm.go`.

3. **Package doc (lines 1–9)** — replace the sentence claiming the controls are
   anchors with no click handler. New wording, same paragraph:

   > the ‹ › controls are `<button>`s; a small click handler
   > (`slideToMonth`) jumps the scroll-snap strip to the neighbouring month
   > card with `ScrollIntoView`, and the browser's scroll-snap animates it.
   > They are deliberately **not** `<a href="#cs-m-...">`: an anchor mutates
   > `location.hash`, which a hash-routed shell (`platformd`) reads as a route
   > change, blanking the view.

### `calendarslider/README.md`

Find every place that describes the ‹ › as anchor links / "no JavaScript" for
navigation and correct it to match the new package doc. The day-click handler
was always there, so "pure Go, zero JavaScript" as a general claim is already
loose — keep the spirit ("no hand-written JS, no infinite-slider DOM
recycling") but do not claim the nav has no handler.

### `calendarslider/css.go`

`clsPrev` / `clsNext` styling: a `<button>` carries a UA background, border and
font. Add whatever is needed so the buttons render visually identical to the
old anchors (`As(Bare)` or equivalent to strip the UA chrome, keep the existing
size / position / `‹`/`›` glyph styling). Do not change their `OnEdge`
placement on the month card.

### Tests

`calendarslider/calendarslider_test.go` (SSR / backend build):

- `TestBuildMonthAgosto2026` (`:39`, `:42`): the two `href='#cs-m-...'`
  assertions → assert the child is a `<button>` and its `data-target` is
  `cs-m-2026-07` / `cs-m-2026-09` respectively (e.g.
  `Contains(children[8].String(), `data-target='cs-m-2026-07'`)` and
  `Contains(children[8].String(), "<button")`). Keep the child-count and
  label assertions.
- `TestRenderWrapsAround` (`:77`, `:82`): `href='#cs-m-2026-10'` /
  `href='#cs-m-2026-08'` → `data-target='cs-m-2026-10'` /
  `data-target='cs-m-2026-08'`. Keep the `calendarslider__prev` /
  `calendarslider__next` count assertions (3 each).
- `:234` (`href='#cs-m-2026-09'`): switch to the `data-target` form.
- `TestBuildMonthAlwaysHasBothLinks`: `"Mes anterior"` / `"Mes siguiente"`
  still present — unchanged.

`calendarslider/calendarslider_wasm_test.go` (real DOM):

- Rename `TestNavLinksTargetTheNeighbor` → `TestNavButtonsSlideToNeighbor`.
  Rewrite its body:
  - Assert `.calendarslider__prev` / `.calendarslider__next` have
    `tagName == "BUTTON"`.
  - Record `js.Global().Get("location").Get("hash").String()` before, click
    `#cs-m-2026-08 .calendarslider__next`, assert the hash is **unchanged**
    (still the same string — do NOT assert it becomes `#cs-m-2026-09`).
  - Keep the "first month has a `prev`, last month has a `next`" existence
    checks (they read class presence, not `href`).
  - Keep the final loop asserting no `#cs-m-...` month element was removed from
    the DOM.
- `TestExternalSelectedHighlights`: after the `.calendarslider__next` click,
  keep "selection survives"; add an assertion that `location.hash` did not
  change.
- `calendarslider_contract_test.go`: read it; if it asserts anchor `href`s,
  move it to the `data-target` / `<button>` form. If it does not touch nav,
  leave it.

## Part 1 acceptance

- `grep -rn 'A("#cs-m-' calendarslider/` → empty.
- `grep -rn "href='#cs-m-\|href=\"#cs-m-" calendarslider/` → empty (source
  **and** tests).
- `grep -n "ScrollIntoView" calendarslider/calendarslider.go` → present once,
  in `slideToMonth`.
- `gotest ./calendarslider/...` green (SSR and WASM).
- `NumMonths` field and `numMonths()` unchanged (`git diff` shows no edit to
  lines 100 or 152–162 beyond context).

---

# PART 2 — one owner for list-selection chrome: move it into `listselect`

## The bug (and the debt behind it)

Testing `#reservation` in delete mode: the top-right corner of the list card is
a pile-up — a dark trash-glyph fragment, the row-count bubble, and the text
`3 / 4`, overlapping the first row.

Two **absolutely-positioned** overlays claim the same corner of the same box
(`.crudview__list`, which is `position: relative`):

| Overlay | Rendered by | Placement | Shows |
|---|---|---|---|
| row-count bubble (`countbadge`) | `layout/crudview` | `OnEdge(EdgeTop, SideEnd)` | total (`4`) |
| select-all master check | THIS repo — `targethour.buildMasterCheck()` etc. | `OnEdge(EdgeTop, SideEnd, …, Space2)` | glyph + `k / N` |

They were built in two separate loops and never reconciled. On top of that:

- **`buildMasterCheck()` is byte-identical** in `targetlist.go`, `targetdate.go`
  and `targethour.go` (~30 lines each). Copy #1, #2, #3.
- The **per-row selection check** (4 `DeriveBool`s + the `check` span with two
  glyph children) is near-identical across the same three files. Copy #1, #2,
  #3.
- The master check's box uses `As(DangerWash)` / `As(AccentWash)` — a pale fill
  with dark `currentColor`, so its glyph reads near-black (the row check solved
  this exact thing with **solid** `As(Danger)` + `As(AccentInverse)` → white
  glyph). And it has no width cap, no solid background, no `z-index`, so its
  content spills onto row 1.

`CONSTRUCTION_HARNESS.md`: *"The glue is written once, in the library that owns
it."* `listselect` owns "the multi-selection mode a record list enters" — so it
owns the **UI** of that mode too, not just the `Mode` state and the CSS
options. Today it half-owns it (`Apply` / `ApplyMaster` are CSS-only; the DOM
is pasted into each widget).

## The fix — `listselect` builds the selection chrome; the three widgets assemble it

No absolute positioning anywhere. The select-all control + the count live in a
**normal in-flow header strip** that `listselect` builds and each widget
`Child()`s above its `<ul>`. The strip is `display:none` in normal mode
(revealed by the widget root's `Open` state, exactly as the check boxes are
today), so **normal mode looks identical to now** — no header, no count. In
selection mode the strip appears: `[select-all box] [k / N]`.

This also deletes `crudview`'s row-count `countbadge` (that is Part of the
`layout` plan — `../../layout/docs/PLAN.md`; do not touch `layout` from here).

### New public API in `components/listselect/`

Add to `listselect.go` (untagged — compiles to WASM). It will need new imports:
`"github.com/tinywasm/fmt"`, `"github.com/tinywasm/widget"`,
`"github.com/tinywasm/icons/trash"`, `"github.com/tinywasm/icons/pencil"`
(all already in `components/go.mod`).

```go
// Row is the per-row selection wiring listselect hands a target* widget so the
// three widgets stop hand-rolling identical derives. Build it once per row in
// buildRow.
//
//   - Check is the glyph box: place it in the row. It reveals a trash glyph
//     when the row is marked while the danger tone is armed, a pencil when
//     marked while it is not, and is invisible otherwise.
//   - Edit   is "marked, danger tone OFF" — the widget ORs this with its own
//     "this is the loaded record" highlight for the row's Selected state.
//   - Danger is "marked, danger tone ON" — bind the row's Invalid state to it.
type Row struct {
    Check  *Element
    Edit   *SignalBool
    Danger *SignalBool
}

// RowOf builds the selection wiring for one row id. name is the host widget's
// WidgetName() — listselect namespaces its parts under it so the CSS
// (ApplyRow) and the element agree.
func RowOf(m *Mode, id string, name widget.Name) Row

// Header builds the in-flow selection header strip: a select-all / deselect-all
// box (trash glyph while the danger tone is armed, pencil otherwise) and a
// count that reads "k / N". Child() it above the list <ul>. It is hidden in
// normal mode (ApplyHeader reveals it on the root's Open state), so the host
// shows no header and no count until selection mode starts.
//
// total reports the row count the widget currently renders (pass the widget's
// Count method value, e.g. t.Count). name is the host's WidgetName().
func Header(m *Mode, total func() int, name widget.Name) *Element
```

Implementation notes (keep them faithful to today's behaviour):

- `RowOf`'s `Edit` / `Danger` are `DeriveBool`s that read `m.Changed().Get()`
  first (so a tap repaints — see the `Mode.changed` field doc), then
  `m.On() && !m.Danger() && m.IsChecked(id)` and
  `m.On() && m.Danger() && m.IsChecked(id)` respectively.
- `RowOf`'s `Check` is a `Span` with class `name.Class(partCheck)`,
  `BindState(widget.Selected, r.Edit)`, `BindState(widget.Invalid, r.Danger)`,
  and two glyph children: `trash.Ref.Render(string(name.Class(partCheckTrash)))`
  and `pencil.Ref.Render(string(name.Class(partCheckPencil)))`.
- `Header`'s select-all box: `Span` class `name.Class(partAll)`,
  `Attr("role", "checkbox")`, `BindAttrBool("aria-checked", allChecked)` where
  `allChecked := DeriveBool(func() bool { _ = m.Changed().Get(); n := m.Count(); return n > 0 && n == total() })`,
  `BindState(widget.Invalid, DeriveBool(func() bool { return m.On().Get() && m.Danger().Get() }))`,
  `BindState(widget.Selected, DeriveBool(func() bool { return m.On().Get() && !m.Danger().Get() }))`,
  two glyph children (`partAllTrash` / `partAllPencil`), and a count child
  `Span` class `name.Class(partCount)` with
  `BindTextFunc(func() string { _ = m.Changed().Get(); return fmt.Sprintf("%d / %d", m.Count(), total()) })`.
- `Header`'s click handler: `if n := m.Count(); n > 0 && n == total() { m.Clear(); return }` else `m.CheckAll(idsInRenderOrder)`.
  **Problem:** `CheckAll` needs the ids in render order, which the widget owns,
  not `listselect`. Fix by having `Header` take them lazily: change the
  signature to `Header(m *Mode, ids func() []string, name widget.Name) *Element`
  and derive `total()` as `len(ids())`. The widget passes a closure returning
  `t.itemIDs()`. (Update the doc comment accordingly — `ids` returns the
  current rows in render order; count is its length.)
- Part name constants are **unexported in `listselect`** (`partCheck`,
  `partCheckTrash`, `partCheckPencil`, `partAll`, `partAllTrash`,
  `partAllPencil`, `partAllCount`, each `widget.Part("sel-...")`). Nothing
  outside `listselect` needs them — the widgets pass only `name`.

### Rewrite `components/listselect/css.go`

Replace `Apply` and `ApplyMaster` with:

```go
// ApplyRow adds the per-row selection-check skin to s (the host widget's
// sheet). The host calls this instead of hand-writing the block. row is the
// host's own row part (its class differs per widget), used only for the
// danger wash under a marked row.
func ApplyRow(s *style.Sheet, row widget.Part) *style.Sheet

// ApplyHeader adds the in-flow selection-header skin to s: the strip is hidden
// until the host root's Open state, then a centred flex row carrying the
// select-all box and the count.
func ApplyHeader(s *style.Sheet) *style.Sheet
```

- `ApplyRow` body = today's `Apply` body but with `listselect`'s own
  `partCheck` / `partCheckTrash` / `partCheckPencil` instead of the passed-in
  parts. Unchanged rules: base `Hide()` + `OnEdge(EdgeTop, SideEnd, SpaceNone,
  SpaceNone)` + `IconBox(IconMd)` + `KeepSize()` + `Round(RadiusSm)` +
  `As(Inset)` + `Animate(MotionFast)` on the box; `Hide()` + `IconBox(IconSm)`
  on each glyph; `WhenWithin(Open, "", partCheck, Show(), Row(SpaceNone),
  CenterContent())`; `WhenWithin(Invalid, partCheck, partCheckTrash, Show())`
  and the Selected/pencil pair; `When(Invalid, partCheck, As(Danger))` and
  `When(Selected, partCheck, As(AccentInverse))`; `When(Invalid, row,
  As(DangerWash))`.
- `ApplyHeader` body — the strip is in flow, not `OnEdge`:
  - `s.Part(partHeader, style.Hide())`
  - `s.Part(partAll, style.IconBox(style.IconMd), style.KeepSize(),
    style.Round(style.RadiusSm), style.As(style.Inset),
    style.Animate(style.MotionFast))` — resting box, no `OnEdge`.
  - `s.Part(partAllTrash, style.Hide(), style.IconBox(style.IconSm))`; same for
    `partAllPencil`.
  - `s.Part(partAllCount, style.FontSize(style.TextXs),
    style.FontWeight(style.WeightBold))`.
  - `s.WhenWithin(widget.Open, "", partHeader, style.Show(),
    style.Row(style.Space2), style.CenterContent())` — appears as a flex strip.
  - `s.WhenWithin(widget.Invalid, partAll, partAllTrash, style.Show())` and the
    `Selected` / `partAllPencil` pair.
  - **Solid fills** (this is the A4 fix — white glyph, opaque box):
    `s.When(widget.Invalid, partAll, style.As(style.Danger))` and
    `s.When(widget.Selected, partAll, style.As(style.Accent))`.

### Update the three widgets

`targetlist/targetlist.go`, `targetdate/targetdate.go`, `targethour/targethour.go`
— the same edit in each:

1. **Delete** `buildMasterCheck()` entirely.
2. **Delete** the local part constants and `cls*` vars for the check /
   check-all (`PartCheck`, `PartCheckTrash`, `PartCheckPencil`, `PartCheckAll`,
   `PartCheckAllTrash`, `PartCheckAllPencil`, `PartCheckAllCount` and their
   `clsCheck*` / `clsCheckAll*`). Keep `PartRow`, `PartList`, `PartLabel`,
   `PartBadge`, and the widget's own lead/hour/content parts.
3. **Drop** the `icons/trash` + `icons/pencil` imports if nothing else in the
   file uses them (the row check moves into `listselect`).
4. `Render()`: `Child(t.buildMasterCheck())` → `Child(listselect.Header(&t.sel,
   t.itemIDs, t.WidgetName()))`.
5. `buildRow(it)`: replace the block of 4 `DeriveBool`s + the `check` span
   with:

   ```go
   r := listselect.RowOf(&t.sel, id, t.WidgetName())
   isSel := DeriveBool(func() bool {
       _ = t.sel.Changed().Get()
       if t.sel.On().Get() {
           return r.Edit.Get()
       }
       return t.Selected.Get() == id
   })
   row.BindState(widget.Selected, isSel).
       BindState(widget.Invalid, r.Danger).
       BindAttrBool("aria-selected", DeriveBool(func() bool { return isSel.Get() || r.Danger.Get() }))
   ```

   and use `r.Check` where the old `check` span was placed in the row content
   (`targetlist`: after the row, before the label; `targetdate` /
   `targethour`: inside `content`, in the same slot as before).
6. Each widget's `css.go` (`//go:build !wasm`): replace
   `listselect.Apply(s, PartCheck, PartCheckTrash, PartCheckPencil, PartRow)`
   with `listselect.ApplyRow(s, PartRow)`, and
   `listselect.ApplyMaster(s, PartCheckAll, …)` with `listselect.ApplyHeader(s)`.

### Tests

- **`components/listselect/`** gets a new WASM test file
  (`listselect_ui_wasm_test.go`, `//go:build wasm`) — the consumer-shaped proof
  the harness requires. Using a real `Mode` and a fixed
  `ids := []string{"a","b","c"}`:
  - `Header`: render it; assert the count text is `0 / 3`; click the select-all
    box; assert `m.Count() == 3` and the text is `3 / 3`; click again; assert
    `m.Count() == 0`.
  - Arm the danger tone (`m.SetDanger(true)`, `m.SetOn(true)`); assert the
    trash glyph is the revealed one (pencil hidden) — read
    `data-state`/class or computed `display` on the two glyph parts.
  - `RowOf`: with `m.SetOn(true)` and `m.Toggle("b")`, assert `Row.Edit` for
    `"b"` is true and `Row.Danger` false; flip `SetDanger(true)` and assert the
    reverse.
- **`targetlist` / `targetdate` / `targethour`** WASM tests: wherever they
  drove `buildMasterCheck` output or the old `clsCheckAll*` classes, retarget
  to the header the widget now renders (`.<name>__sel-header`,
  `.<name>__sel-all`, `.<name>__sel-count`). The row-check assertions move from
  `.<name>__check` to the same class emitted by `RowOf` (keep `partCheck` =
  `widget.Part("sel-check")` so the class is `.<name>__sel-check` — update the
  test selectors to match). Behavioural assertions (enter selection mode →
  boxes appear; tap a row → it marks; normal mode → no check, no header) stay
  and must pass.
- `components/conformance_test.go` (imports `targetlist`): if it asserts the
  old check classes, update; otherwise leave.

## Part 2 acceptance

- `grep -rn "buildMasterCheck" components/` → empty.
- `grep -rn "func Apply\b\|func ApplyMaster\b" components/listselect/` → empty
  (replaced by `ApplyRow` / `ApplyHeader`).
- `grep -rn "OnEdge" components/listselect/css.go` → only in `ApplyRow`'s
  per-row box rule (the header has **no** `OnEdge`).
- `grep -rn "DangerWash\|AccentWash" components/listselect/css.go` → empty (the
  header fills are solid `Danger` / `Accent`).
- `grep -rn "PartCheckAll" components/targetlist/ components/targetdate/ components/targethour/` → empty.
- The per-row derive block (`isEditCheckSig`, `isDangerSig` …) appears **once**
  — inside `listselect`, not in the three widgets:
  `grep -rn "isEditCheckSig" components/target*/` → empty.
- `gotest ./...` green.
- `GOOS=js GOARCH=wasm go list -deps ./... | grep tinywasm/svg/sprite` → empty.
- Manual (daemon hot-reloads `components`): open the demo `#reservation`, pick a
  day, press 🗑 → a single clean header strip above the list (`[trash box] 0 /
  4`), white glyph, nothing overlapping row 1; tap rows → count climbs, box
  goes solid red; press select-all → `4 / 4`, all rows red; leave mode → header
  gone, list identical to before. Repeat with ✏ → same, blue/`Accent`.

---

## Docs to update before finishing (both parts)

- `components/docs/CATALOG.md` — `calendarslider` nav description;
  `listselect` entry gains `Header` / `RowOf` / `ApplyRow` / `ApplyHeader`.
- `components/docs/ARCHITECTURE.md` (or `DESIGN.md` if that is where component
  contracts live) — a short "listselect owns the selection chrome (header +
  row check), the target\* widgets assemble it" note; the calendarslider nav
  mechanism (button + `slideToMonth`, not anchors, and why).
- `components/README.md` — must still index every file in `docs/`.
- `components/AGENTS.md` — the line about `listselect.Apply` being "the check
  mark's whole visual contract" → update to `ApplyRow` / `ApplyHeader` +
  `Header` / `RowOf` (the piece now owns the DOM, not only the skin).

## Stages

| # | Scope | Files | Done when |
|---|---|---|---|
| 1 | calendarslider nav | `calendarslider/calendarslider.go`, `css.go`, `README.md`, `calendarslider_test.go`, `calendarslider_wasm_test.go`, `calendarslider_contract_test.go` | Part 1 acceptance green |
| 2 | listselect owns the chrome | `listselect/listselect.go`, `listselect/css.go`, new `listselect/listselect_ui_wasm_test.go` | `Header`/`RowOf`/`ApplyRow`/`ApplyHeader` exist + tested |
| 3 | assemble in the widgets | `targetlist/{targetlist.go,css.go,*_wasm_test.go}`, `targetdate/…`, `targethour/…`, `components/conformance_test.go` | Part 2 acceptance green |
| 4 | docs | `docs/CATALOG.md`, `docs/ARCHITECTURE.md`/`DESIGN.md`, `README.md`, `AGENTS.md` | indexed + accurate |

Final: `gotest ./...` green, WASM build clean, sprite-leak check empty.
