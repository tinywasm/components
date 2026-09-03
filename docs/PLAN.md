---
PLAN: "feat: targethour list component + calendarslider widget.Filterable"
EXECUTOR: jules
REVIEWER: none
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Part of `RESERVATION_MODULE_MASTER_PLAN.md` (monorepo root), Phase A.
> **DRAFT** — the maintainer will refine it. Implement what is specified; where
> a detail is marked "maintainer will refine", implement the stated default.

# Plan — `targethour` + `calendarslider` as a filter control

Two independent changes in `github.com/tinywasm/components`, one commit:

1. **`calendarslider`** gains `widget.Filterable` so a `crudview` can use a
   month calendar as its filter control (pick a day → filter the list to that
   day), the same way `searchbar` and `selectsearch` already do.
2. **`targethour`** — a NEW selectable list component, sibling of `targetlist`
   and `targetdate`, whose rows lead with a prominent **hour** (`HH:MM`) and
   carry an optional per-row **status tint**.

Neither touches the WASM binary surface beyond what `targetlist`/`targetdate`
already do.

---

## Ecosystem rules (this repo — apply to every file you write)

`github.com/tinywasm/components` is a **public library**: code, comments,
identifiers and error messages are **in English**.

- **No Go stdlib in untagged or `//go:build wasm` files.** Use
  `github.com/tinywasm/fmt` instead of `errors`/`strconv`/`strings`
  (`fmt.Sprintf`, `fmt.Err`, `fmt.Errf`, `fmt.Contains`). `tinywasm/fmt` has
  **no** `Itoa`/`FormatInt` — intercalate an int with `fmt.Sprintf("%d", n)`.
  Stdlib IS allowed in `_test.go` files (`testing`, `strings`, `reflect`).
- **No `map[...]`** in untagged / `wasm` files (pulls TinyGo hashing +
  runtime). A set is a `[]string`; a linear scan over tens of entries is free.
  A `map` is fine in a `//go:build !wasm` file (`css.go`).
- **Embed `dom.Element` by value**, never `*dom.Element`.
- **SSR split by extension**: CSS in `css.go` under `//go:build !wasm` with the
  entry point named EXACTLY `RenderCSS` (a method, not a free function — the
  dot-imported `css` package already exports a free `RenderCSS`); SVG geometry
  in `svg.go` under `//go:build !wasm` with the entry point named EXACTLY
  `IconSvg` (also a method on the same receiver type). A CSS builder named
  anything else is silently never emitted — the component renders unstyled and
  nothing fails at build time.
- **No `front.go`.** WASM interactivity lives in the main component file.
- **Icons**: reference from an untagged file
  (`github.com/tinywasm/icons/trash`, `.../pencil` → `trash.Ref`,
  `pencil.Ref`); geometry from `svg.go` (`trash.Def()`, `pencil.Def()`). NEVER
  import `github.com/tinywasm/svg/sprite` from an untagged file.
- **Tests**: `gotest`, never `go test`. Stdlib `testing` only — no testify.

Mandatory pre-finish check (must print nothing):

```bash
GOOS=js GOARCH=wasm go list -deps ./... | grep tinywasm/svg/sprite
```

---

## Stage 1 — `calendarslider` implements `widget.Filterable`

**File: `calendarslider/calendarslider.go`** (untagged).

`widget.Filterable` is the whole contract (`github.com/tinywasm/widget`,
`capability.go`):

```go
type Filterable interface{ OnFilterChange(func(term string)) }
```

`searchbar` and `selectsearch` already implement it; `crudview` auto-wires any
`Filter` control that satisfies it (see
`https://github.com/tinywasm/layout/blob/main/crudview/crudview.go`, the
`if src, ok := v.Filter.(widget.Filterable); ok` block in `Init`).

Changes:

1. Add an unexported field to the `CalendarSlider` struct, next to `OnSelect`:

   ```go
   onFilter func(term string) // set via OnFilterChange — satisfies widget.Filterable
   ```

2. Add the compile proof and the method (place near `WidgetName`/`WidgetKind`):

   ```go
   var _ widget.Filterable = (*CalendarSlider)(nil)

   // OnFilterChange implements widget.Filterable: it registers the sink called
   // with the picked day ("YYYY-MM-DD") on every day selection. The signature
   // is fixed by widget.Filterable — do not add a parameter, do not rename.
   func (c *CalendarSlider) OnFilterChange(fn func(term string)) { c.onFilter = fn }
   ```

3. In the day-cell click handler (the `li.On("click", func(Event) { ... })`
   block, currently around line 372), fire the sink after `OnSelect`:

   ```go
   li.On("click", func(Event) {
       c.Selected.Set(dateStr)
       if c.OnSelect != nil {
           c.OnSelect(dateStr)
       }
       if c.onFilter != nil {
           c.onFilter(dateStr)
       }
   })
   ```

   `OnSelect` and `onFilter` are BOTH fired — a host may want either or both
   (mirrors `selectsearch`, which exposes its own richer callback alongside
   `OnFilterChange`). Do not remove `OnSelect`.

**Do NOT** change `Render`, the strip-building logic, the `‹ ›` navigation, or
any CSS. This stage is additive only.

### Stage 1 tests — `calendarslider/calendarslider_test.go` (`//go:build !wasm`)

Add one test:

```go
func TestCalendarSlider_SatisfiesFilterable(t *testing.T) {
	var c widget.Filterable = &CalendarSlider{}
	got := ""
	c.OnFilterChange(func(term string) { got = term })
	// the sink is stored; the click that fires it is covered by the wasm test
	if got != "" {
		t.Fatalf("sink must not fire on registration, got %q", got)
	}
}
```

In `calendarslider/calendarslider_wasm_test.go` (`//go:build wasm`), extend the
existing day-click test (or add one) to assert that clicking a selectable day
cell fires the `OnFilterChange` sink with that cell's `data-date` value.

---

## Stage 2 — new component `components/targethour/`

A flat package (no subfolders). Mirror `components/targetdate/` almost exactly
— same multi-selection mechanics (it assembles `components/listselect`), same
`view.Item` alias, same `crudview.ListView` shape — the ONLY differences are
(a) the lead renders a prominent **hour** instead of a stacked date badge, and
(b) an optional per-row **status tint**.

Reference implementation to copy structure from (verbatim where noted):
`https://github.com/tinywasm/components/blob/main/targetdate/targetdate.go`
and its `css.go`, `svg.go`, `dangermode_test.go`, `selectmode_test.go`,
`row_test.go`. The full current `targetdate.go` is appended at the end of this
plan as **Appendix A** — recycle its `buildRow` selection-state logic
(`isEditCheckSig`, `isDangerSig`, `isSelSig`, `isMarkedSig`, the `check` span
with `BindState(widget.Selected/Invalid, ...)`, the row `click` handler)
**unchanged** — that logic is shared and must not drift between the three
list components.

### Files to create

| File | Build tag | Contents |
|---|---|---|
| `targethour/targethour.go` | *(none)* | `TargetHour` struct, `Name*`/`Part*`/`cls*` constants, `Item` alias, `Init`, `SetSelectMode`, `SetDanger`, `OnCheckedChange`, `CheckedIDs`, `SetItems`, `Items`, `Count`, `Render`, `buildRow`, the `Status` enum + `StatusOf` field. |
| `targethour/css.go` | `//go:build !wasm` | `func (t *TargetHour) RenderCSS() *css.Stylesheet` + unexported `sheet()`. |
| `targethour/svg.go` | `//go:build !wasm` | `func (t *TargetHour) IconSvg() *sprite.Sprite` returning `sprite.NewSprite(trash.Def(), pencil.Def())`. |
| `targethour/row_test.go` | `//go:build !wasm` | markup assertions — mirror `targetdate/row_test.go`. |
| `targethour/dangermode_test.go` | `//go:build !wasm` | mirror `targetdate/dangermode_test.go` (selectors become `.targethour__*`). |
| `targethour/selectmode_test.go` | `//go:build !wasm` (`package targethour_test`) | mirror `targetdate/selectmode_test.go`. |
| `targethour/statustint_test.go` | `//go:build !wasm` | the new status-tint behaviour (see below). |

### `targethour/targethour.go` — spec

Package doc:

```go
// Package targethour is targetlist's sibling for a day's booked slots: each
// row leads with a prominent hour (HH:MM) and may carry a status tint
// (pending / confirmed / attended). Same multi-selection mechanics as
// targetlist/targetdate — it assembles components/listselect, it does not
// re-declare it.
package targethour
```

Identity:

```go
const NameTargetHour = widget.Name("targethour")

const (
	PartRow         = widget.Part("row")
	PartContent     = widget.Part("content")
	PartCheck       = widget.Part("check")
	PartCheckTrash  = widget.Part("check-trash")
	PartCheckPencil = widget.Part("check-pencil")
	PartHour        = widget.Part("hour")
	PartLabel       = widget.Part("label")
	PartBadge       = widget.Part("badge")
	PartList        = widget.Part("list")
)
```

`cls*` vars mirror `targetdate` (`NameTargetHour.Root()`, `.Class(Part*)`).

```go
type Item = view.Item // same alias reason as targetlist.Item

// Status is a row's booking state — drives the tint only. The zero value
// (StatusPending) paints no tint, exactly like a plain targetlist row.
type Status uint8

const (
	StatusPending Status = iota // no tint
	StatusConfirmed             // confirmed by reception
	StatusAttended             // patient already attended
)

type TargetHour struct {
	Element

	Selected *SignalString
	OnSelect func(it Item)

	// StatusOf maps a row to its booking state for the tint. Optional — nil
	// means every row is StatusPending (no tint). The host owns the mapping
	// from its own model / view.Item to this typed enum, so the library holds
	// no localized status strings.
	StatusOf func(it Item) Status

	items []Item
	rows  *SignalNodes
	sel   listselect.Mode
}
```

`WidgetName` → `NameTargetHour`; `WidgetKind` → `widget.Combobox` (same as
targetdate). `Init`, `ensure`, `SetSelectMode`, `SetDanger`,
`OnCheckedChange`, `CheckedIDs`, `SetItems`, `Items`, `Count`, `Render`:
**copy verbatim** from `targetdate` (Appendix A), renaming `TargetDate`→
`TargetHour`, `t.sel`, `clsListWrap`/`clsList`, key prefix `"td-"`→`"th-"`.

`buildRow(it Item) *Element`:

- Copy the selection-state block (`isEditCheckSig` … `isMarkedSig`), the `row`
  element with its `BindState`/`BindAttrBool`/`click` handler, and the `check`
  span **verbatim** from Appendix A.
- The lead is a single prominent hour, NOT a stacked badge:

  ```go
  hour := Span().Set(clsHour.AsAttr()).Text(it.LeadMain)
  ```

  (`it.LeadMain` carries `"HH:MM"` — the host sets it. `LeadTop`/`LeadBottom`
  are ignored by this component.)
- Content: hour, then check, then label, then optional description chip —
  mirror `targetdate`'s `content` assembly:

  ```go
  content := Div().Set(clsContent.AsAttr()).
      Child(hour).
      Child(check).
      Child(Span().Set(clsLabel.AsAttr()).Text(it.Label))
  if it.Description != "" {
      content.Child(Span().Set(clsBadge.AsAttr()).
          Attr("title", it.Description).
          Text(fmt.Convert(it.Description).Truncate(badgeChars).String()))
  }
  row.Child(content)
  ```

  Keep `const badgeChars = 16` (same as targetdate).
- Status tint: bind a widget state on the ROW from `StatusOf`. Use
  `widget.Current` for `StatusConfirmed` and `widget.Busy` for
  `StatusAttended` (both are existing `widget.State` values — see
  `https://github.com/tinywasm/widget/blob/main/state.go`; do NOT invent a new
  state, that is a `widget` change out of scope):

  ```go
  st := StatusPending
  if t.StatusOf != nil {
      st = t.StatusOf(it)
  }
  row.BindStateFunc(widget.Current, func() bool { return st == StatusConfirmed })
  row.BindStateFunc(widget.Busy, func() bool { return st == StatusAttended })
  ```

  (`st` is captured per row at build time — the host rebuilds rows via
  `SetItems` on every reload, same as `targetdate`, so a plain capture is
  correct; do NOT add a signal.)

### `targethour/css.go` — spec

Copy `targetdate/css.go` structure. `sheet()`:

- `style.For(t).Root(style.Fill(), style.Stack(style.SpaceNone))`.
- `listgap.Apply(s, PartList)` + `s.On(css.Mobile, PartList, listgap.MobileOpts()...)`.
- `listselect.Apply(s, PartCheck, PartCheckTrash, PartCheckPencil, PartRow)`.
- `Part(PartRow, ...)` — same as targetdate's row.
- `Part(PartContent, style.Row(style.Space2), style.Grow(), style.Pad(style.Space2))`.
- `Part(PartHour, style.FontSize(style.TextLg), style.FontWeight(style.WeightBold), style.KeepSize(), style.CenterContent(), style.Divider(style.SideEnd), style.Pad(style.Space2))`
  — the hour is the prominent lead; `Divider(SideEnd)` gives the same hairline
  targetdate's `PartLead` uses.
- `Part(PartLabel, style.FontWeight(style.WeightBold), style.Grow())`.
- `Part(PartBadge, ...)` — copy targetdate's `PartBadge` verbatim.
- Selection cues on `PartRow` — copy targetdate's `When(widget.Selected, ...)`
  / `Cue(widget.Hover/Focus/Press, ...)` verbatim.
- **Status tint** — two overlay rules, low-key washes so they never fight the
  blue selection highlight:

  ```go
  When(widget.Current, PartRow, style.As(style.AccentWash)).  // confirmed
  When(widget.Busy, PartRow, style.As(style.Subtle)).         // attended (muted)
  ```

  (maintainer will refine the exact tokens; ship these.)

### `targethour/svg.go` — spec

```go
//go:build !wasm

package targethour

import (
	"github.com/tinywasm/icons/pencil"
	"github.com/tinywasm/icons/trash"
	"github.com/tinywasm/svg/sprite"
)

func (t *TargetHour) IconSvg() *sprite.Sprite {
	return sprite.NewSprite(trash.Def(), pencil.Def())
}
```

### `targethour/statustint_test.go` — spec

```go
//go:build !wasm

package targethour

import (
	"strings"
	"testing"

	"github.com/tinywasm/view"
	"github.com/tinywasm/widget"
)

// A row with no StatusOf carries no tint state.
func TestTargetHour_NoStatusMapperNoTint(t *testing.T) {
	th := &TargetHour{}
	th.Init(nil)
	html := th.buildRow(Item{ID: "1", Label: "x", LeadMain: "09:00"}).String()
	if strings.Contains(html, string(widget.Current.Attr().Key())+"='true'") ||
		strings.Contains(html, string(widget.Busy.Attr().Key())+"='true'") {
		t.Errorf("a row without a StatusOf mapper must carry no tint state\n%s", html)
	}
}

// StatusConfirmed -> data-current='true'; StatusAttended -> data-busy='true';
// StatusPending -> neither.
func TestTargetHour_StatusDrivesTheTintState(t *testing.T) {
	cases := []struct {
		st        Status
		wantKey   string
		wantOther string
	}{
		{StatusConfirmed, widget.Current.Attr().Key(), widget.Busy.Attr().Key()},
		{StatusAttended, widget.Busy.Attr().Key(), widget.Current.Attr().Key()},
	}
	for _, c := range cases {
		th := &TargetHour{StatusOf: func(view.Item) Status { return c.st }}
		th.Init(nil)
		html := th.buildRow(Item{ID: "1", Label: "x", LeadMain: "09:00"}).String()
		if !strings.Contains(html, c.wantKey+"='true'") {
			t.Errorf("status %d must set %s='true'\n%s", c.st, c.wantKey, html)
		}
		if strings.Contains(html, c.wantOther+"='true'") {
			t.Errorf("status %d must NOT set %s='true'\n%s", c.st, c.wantOther, html)
		}
	}
}

// The lead renders it.LeadMain as the prominent hour.
func TestTargetHour_LeadRendersTheHour(t *testing.T) {
	th := &TargetHour{}
	th.Init(nil)
	html := th.buildRow(Item{ID: "1", Label: "Ana", LeadMain: "14:30"}).String()
	if !strings.Contains(html, "targethour__hour") || !strings.Contains(html, "14:30") {
		t.Errorf("the lead must render LeadMain as the hour\n%s", html)
	}
}
```

The `dangermode_test.go` / `selectmode_test.go` / `row_test.go` mirrors: take
`targetdate`'s files, replace `targetdate`→`targethour`, `TargetDate`→
`TargetHour`, `td-`→`th-`, and the lead assertions (targetdate checks
`lead-top`/`lead-main`/`lead-bottom`; targethour checks `targethour__hour`
carries the hour). Keep every check-box / selection-mode / danger-mode
assertion identical — that behaviour is shared and must not regress.

---

## Stage 3 — documentation (verify against the code you wrote)

- **`components/docs/CATALOG.md`** — add a `targethour` entry next to
  `targetdate` and `targetlist`: one line on what it is (a day's booked-slot
  list, hour lead, status tint) and its `crudview.ListView` role. Note
  `calendarslider` now satisfies `widget.Filterable`.
- **`components/docs/SKILL.md`** — if it enumerates the list components or the
  filter controls, add `targethour` / the `calendarslider` Filterable line.
- **`components/README.md`** — must index every file under `docs/`; if it lists
  components, add `targethour`.
- Do NOT reference this `docs/PLAN.md` from any permanent doc — it is deleted
  when the PR merges.

---

## Acceptance criteria

```bash
go build ./...                                              # clean
gotest                                                      # vet + tests + race + wasm, all green
GOOS=js GOARCH=wasm go list -deps ./... | grep tinywasm/svg/sprite   # prints NOTHING
GOOS=js GOARCH=wasm go list -deps ./targethour | grep -E 'errors|strconv|strings'  # prints NOTHING (no stdlib in wasm graph)
grep -rn "package targethour" targethour/                   # the new package exists
grep -n "widget.Filterable" calendarslider/calendarslider.go # the interface proof is present
```

- `targethour` satisfies `crudview.ListView` — verify by adding, in
  `selectmode_test.go`, a compile-time assertion via a local interface that
  lists the `ListView` methods (do NOT import `layout` — `components` must not
  depend on `layout`; declare the method-set locally and assign
  `var _ theInterface = (*targethour.TargetHour)(nil)`).
- Every `targetdate` check-box / danger-mode / select-mode test has a green
  `targethour` mirror.

## Stages table

| # | File(s) | Action |
|---|---|---|
| 1 | `calendarslider/calendarslider.go`, `calendarslider/calendarslider_test.go`, `calendarslider/calendarslider_wasm_test.go` | add `widget.Filterable` (field + `var _` + method + fire sink in day click) + tests |
| 2 | `targethour/{targethour.go,css.go,svg.go,row_test.go,dangermode_test.go,selectmode_test.go,statustint_test.go}` | new component, structure copied from `targetdate`, lead = hour, `StatusOf` tint |
| 3 | `components/docs/CATALOG.md`, `components/docs/SKILL.md`, `components/README.md` | add `targethour` + `calendarslider` Filterable; verify against code |

---

## Appendix A — current `targetdate/targetdate.go` (recycle its buildRow logic)

```go
// Package targetdate is targetlist's sibling for rows that need a prominent
// leading badge — an hour, a day, anything view.Item.LeadTop/Main/Bottom
// carries — instead of a plain label. Same multi-selection mechanics as
// targetlist.
package targetdate

import (
	. "github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/view"
	"github.com/tinywasm/widget"

	"github.com/tinywasm/components/listselect"
	"github.com/tinywasm/icons/pencil"
	"github.com/tinywasm/icons/trash"
)

const badgeChars = 16

const NameTargetDate = widget.Name("targetdate")

const (
	PartRow         = widget.Part("row")
	PartContent     = widget.Part("content")
	PartCheck       = widget.Part("check")
	PartCheckTrash  = widget.Part("check-trash")
	PartCheckPencil = widget.Part("check-pencil")
	PartBadge       = widget.Part("badge")
	PartLabel       = widget.Part("label")
	PartList        = widget.Part("list")
	PartLead        = widget.Part("lead")
	PartLeadStack   = widget.Part("lead-stack")
	PartLeadTop     = widget.Part("lead-top")
	PartLeadMain    = widget.Part("lead-main")
	PartLeadBottom  = widget.Part("lead-bottom")
)

// (cls* vars: NameTargetDate.Root() and NameTargetDate.Class(Part*))

type Item = view.Item

type TargetDate struct {
	Element
	Selected *SignalString
	OnSelect func(it Item)
	items    []Item
	rows     *SignalNodes
	sel      listselect.Mode
}

func (t *TargetDate) WidgetName() widget.Name { return NameTargetDate }
func (t *TargetDate) WidgetKind() widget.Kind { return widget.Combobox }

func (t *TargetDate) ensure() {
	if t.rows == nil {
		t.rows = NewNodes()
	}
	if t.Selected == nil {
		t.Selected = NewString("")
	}
}

func (t *TargetDate) Init(_ Ctx) { t.ensure() }

func (t *TargetDate) SetSelectMode(on bool)        { t.sel.SetOn(on) }
func (t *TargetDate) SetDanger(on bool)            { t.sel.SetDanger(on) }
func (t *TargetDate) OnCheckedChange(fn func(int)) { t.sel.OnChange = fn }

func (t *TargetDate) CheckedIDs() []string {
	ids := make([]string, len(t.items))
	for i, it := range t.items {
		ids[i] = it.ID
	}
	return t.sel.CheckedIDs(ids)
}

func (t *TargetDate) SetItems(items []Item) {
	t.ensure()
	t.items = items
	nodes := make([]*Element, 0, len(items))
	for _, it := range items {
		nodes = append(nodes, t.buildRow(it))
	}
	t.rows.Set(nodes)
}

func (t *TargetDate) Items() []Item { return t.items }
func (t *TargetDate) Count() int    { return len(t.items) }

func (t *TargetDate) Render() *Element {
	list := Ul().Set(clsList.AsAttr()).Attr("role", "listbox").BindChildren(t.rows)
	return Div().Set(clsListWrap.AsAttr()).
		BindStateFunc(widget.Open, func() bool { return t.sel.On().Get() }).
		Child(list)
}

func (t *TargetDate) buildRow(it Item) *Element {
	id := it.ID
	key := "td-" + id

	isEditCheckSig := DeriveBool(func() bool {
		_ = t.sel.Changed().Get()
		return t.sel.On().Get() && !t.sel.Danger().Get() && t.sel.IsChecked(id)
	})
	isDangerSig := DeriveBool(func() bool {
		_ = t.sel.Changed().Get()
		return t.sel.On().Get() && t.sel.Danger().Get() && t.sel.IsChecked(id)
	})
	isSelSig := DeriveBool(func() bool {
		_ = t.sel.Changed().Get()
		if t.sel.On().Get() {
			return isEditCheckSig.Get()
		}
		return t.Selected.Get() == id
	})
	isMarkedSig := DeriveBool(func() bool { return isSelSig.Get() || isDangerSig.Get() })

	row := Li().Set(clsRow.AsAttr()).
		ID(key).
		Key(key).
		Attr("role", "option").
		BindAttrBool("aria-selected", isMarkedSig).
		BindState(widget.Selected, isSelSig).
		BindState(widget.Invalid, isDangerSig)

	row.On("click", func(Event) {
		if t.sel.On().Get() {
			t.sel.Toggle(id)
			return
		}
		if t.OnSelect != nil {
			t.OnSelect(it)
		}
	})

	lead := Div().Set(clsLead.AsAttr()).Child(
		Div().Set(clsLeadStack.AsAttr()).Child(
			Span().Set(clsLeadTop.AsAttr()).Text(it.LeadTop),
			Span().Set(clsLeadMain.AsAttr()).Text(it.LeadMain),
			Span().Set(clsLeadBottom.AsAttr()).Text(it.LeadBottom),
		),
	)

	check := Span().Set(clsCheck.AsAttr()).
		BindState(widget.Selected, isEditCheckSig).
		BindState(widget.Invalid, isDangerSig).
		Child(trash.Ref.Render(string(clsCheckTrash))).
		Child(pencil.Ref.Render(string(clsCheckPencil)))

	content := Div().Set(clsContent.AsAttr()).
		Child(check).
		Child(Span().Set(clsLabel.AsAttr()).Text(it.Label))
	if it.Description != "" {
		content.Child(Span().Set(clsBadge.AsAttr()).
			Attr("title", it.Description).
			Text(fmt.Convert(it.Description).Truncate(badgeChars).String()))
	}

	row.Child(lead)
	row.Child(content)
	return row
}
```

For `targethour`, the lead becomes `Span().Set(clsHour.AsAttr()).Text(it.LeadMain)`
placed as the FIRST child of `content` (before `check`), and the `StatusOf`
tint binds `widget.Current` / `widget.Busy` on `row` as described in Stage 2.
