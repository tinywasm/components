---
PLAN: "feat: master select-all check for target* record lists"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 17971103714245978648
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> **DRAFT** — the maintainer will refine the visual state mapping (which
> widget.State drives which look). Implement the stated defaults; the wiring
> and API must be exactly as specified.

# Plan — a master "select all" check for `targetlist` / `targetdate` / `targethour`

Today a record list in selection mode is checked **one row at a time**. Add a
**master check** at the top-end corner of the list that selects / deselects
every row at once, and shows a `n / total` count. It only appears in selection
mode (root `data-open`), reuses the same trash/pencil glyph box the row checks
use, and lives entirely inside `components` — `crudview` needs no change (its
`🗑 N` footer badge already tracks the count through the existing
`OnCheckedChange` callback).

Two layers:

1. `components/listselect` — `Mode` gains `CheckAll` / `Clear` / `Count`, and a
   new `ApplyMaster` skin function (sibling of the existing `Apply`).
2. `components/{targetlist,targetdate,targethour}` — each renders the master
   check as the first child of its list root, wired to the tri-state toggle.

---

## Ecosystem rules (this repo — apply to every file you write)

`github.com/tinywasm/components` is a **public library**: code, comments,
identifiers, error messages **in English**.

- **No Go stdlib in untagged / `//go:build wasm` files.** Use
  `github.com/tinywasm/fmt` (`fmt.Sprintf`, `fmt.Err`). `tinywasm/fmt` has **no**
  `Itoa`/`FormatInt` — an int goes through `fmt.Sprintf("%d", n)`. Stdlib IS
  allowed in `_test.go`.
- **No `map[...]`** in untagged / `wasm` files — a set is a `[]string`, a linear
  scan over tens of rows is free. `map` is fine only in a `//go:build !wasm`
  file (`css.go`).
- **Embed `dom.Element` by value.**
- **SSR split**: CSS in `css.go` under `//go:build !wasm`, entry point named
  EXACTLY `RenderCSS` (a method); SVG in `svg.go` under `//go:build !wasm`,
  entry point `IconSvg` (a method on the same receiver). A misnamed builder is
  silently never emitted.
- **Icons**: reference (`trash.Ref`, `pencil.Ref`) from untagged files; geometry
  (`trash.Def()`, `pencil.Def()`) only from `svg.go`. NEVER import
  `github.com/tinywasm/svg/sprite` from an untagged file.
- **Tests**: `gotest`, never `go test`. Stdlib `testing` only.

Mandatory pre-finish check (must print nothing):

```bash
GOOS=js GOARCH=wasm go list -deps ./... | grep tinywasm/svg/sprite
```

---

## Stage 1 — `listselect/listselect.go` (`Mode` API)

`Mode` currently exposes `On`, `SetOn`, `SetDanger`, `Danger`, `Changed`,
`Toggle`, `IsChecked`, `CheckedIDs`, `OnChange`. Add three methods (place them
right after `Toggle`):

```go
// CheckAll marks every id in ids — the caller's CURRENT render order — and
// replaces any previous selection. It owns a fresh backing array (never
// aliases the caller's slice). Fires Changed() and OnChange with the new
// count. This is the master check's "select all" action.
func (m *Mode) CheckAll(ids []string) {
	m.ensure()
	m.checked = append([]string(nil), ids...)
	m.Changed().Toggle()
	if m.OnChange != nil {
		m.OnChange(len(m.checked))
	}
}

// Clear unmarks every row WITHOUT leaving selection mode — unlike
// SetOn(false), which also exits the mode. The master check's "deselect all".
// A no-op (no signal churn) when nothing is marked.
func (m *Mode) Clear() {
	m.ensure()
	if len(m.checked) == 0 {
		return
	}
	m.checked = nil
	m.Changed().Toggle()
	if m.OnChange != nil {
		m.OnChange(0)
	}
}

// Count reports how many rows are currently marked. The master check reads it
// to decide its tri-state (none / some / all) and to render "n / total".
func (m *Mode) Count() int { return len(m.checked) }
```

Do NOT change `Toggle`, `SetOn`, `CheckedIDs` or any existing method.

### Stage 1 tests — `listselect/listselect_test.go`

Add (the file is `//go:build !wasm`, `package listselect_test`):

```go
func TestCheckAllMarksEveryIDInOrder(t *testing.T) {
	var m listselect.Mode
	m.SetOn(true)
	m.CheckAll([]string{"a", "b", "c"})
	if m.Count() != 3 {
		t.Fatalf("Count = %d, want 3", m.Count())
	}
	got := m.CheckedIDs([]string{"a", "b", "c"})
	if len(got) != 3 || got[0] != "a" || got[2] != "c" {
		t.Errorf("CheckedIDs = %v, want [a b c]", got)
	}
}

func TestClearUnmarksButStaysInSelectionMode(t *testing.T) {
	var m listselect.Mode
	m.SetOn(true)
	m.CheckAll([]string{"a", "b"})
	m.Clear()
	if m.Count() != 0 {
		t.Errorf("Clear must unmark everything, Count = %d", m.Count())
	}
	if !m.On().Get() {
		t.Errorf("Clear must NOT leave selection mode")
	}
}

func TestCheckAllOwnsItsBackingArray(t *testing.T) {
	var m listselect.Mode
	m.SetOn(true)
	src := []string{"a", "b"}
	m.CheckAll(src)
	src[0] = "MUTATED"
	if m.IsChecked("MUTATED") || !m.IsChecked("a") {
		t.Errorf("CheckAll must copy, not alias, the caller's slice")
	}
}

func TestCheckAllFiresOnChangeWithTotal(t *testing.T) {
	var m listselect.Mode
	var last int
	m.OnChange = func(n int) { last = n }
	m.SetOn(true)
	m.CheckAll([]string{"a", "b", "c"})
	if last != 3 {
		t.Errorf("OnChange got %d, want 3", last)
	}
	m.Clear()
	if last != 0 {
		t.Errorf("OnChange after Clear got %d, want 0", last)
	}
}
```

---

## Stage 2 — `listselect/css.go` (`ApplyMaster` skin)

`Apply(s, check, trashIcon, pencilIcon, row)` styles the per-row check. Add a
sibling `ApplyMaster` for the master check — same DRY reason as `Apply`: the
three target lists must not each own a private copy.

```go
// ApplyMaster adds the master (select-all) check's skin to s. targetlist,
// targetdate and targethour call this instead of hand-writing the block.
//
// Same shape as Apply's per-row check: hidden until selection mode opens
// (Open on the list root), then a centred flex box carrying one glyph — trash
// while the danger tone is armed (Invalid on the box), pencil otherwise
// (Selected on the box). The box fills solid when every row is marked (Locked
// on the box); it shows a lighter "some marked" wash otherwise (Busy). count
// is a small "n / total" label beside the glyph.
//
// MAINTAINER: the exact tokens for the all/some/none looks are a first pass —
// refine As(...) below.
func ApplyMaster(s *style.Sheet, checkAll, trashIcon, pencilIcon, count widget.Part) *style.Sheet {
	return s.
		Part(checkAll,
			style.Hide(),
			style.IconBox(style.IconMd),
			style.KeepSize(),
			style.Round(style.RadiusSm),
			style.As(style.Inset),
			style.OnEdge(style.EdgeTop, style.SideEnd, style.SpaceNone, style.Space2),
			style.Animate(style.MotionFast),
		).
		Part(trashIcon, style.Hide(), style.IconBox(style.IconSm)).
		Part(pencilIcon, style.Hide(), style.IconBox(style.IconSm)).
		Part(count,
			style.Hide(),
			style.FontSize(style.TextXs),
			style.FontWeight(style.WeightBold),
		).
		WhenWithin(widget.Open, "", checkAll,
			style.Show(),
			style.Row(style.Space1),
			style.CenterContent(),
		).
		WhenWithin(widget.Open, "", count,
			style.Show(),
		).
		WhenWithin(widget.Invalid, checkAll, trashIcon, style.Show()).
		WhenWithin(widget.Selected, checkAll, pencilIcon, style.Show()).
		When(widget.Locked, checkAll, style.As(style.Danger)).       // all rows marked
		When(widget.Busy, checkAll, style.As(style.DangerWash))      // some rows marked
}

// Ensure compile check for css import
var _ = css.Mobile
```

(`css.go` already imports `css`, `widget`, `widget/style`.)

---

## Stage 3 — the three list components

Do **`targetlist` first**, then apply the identical pattern (renamed) to
`targetdate` and `targethour`. The full current `targetlist.go` structure is
your reference — you are ADDING a master check, changing nothing that exists.

### 3a. New parts + `cls*` vars

In each component's `const (...)` part block and `var (...)` cls block, add
(shown for `targetlist` — `targetdate` uses `NameTargetDate`, `targethour`
`NameTargetHour`):

```go
	PartCheckAll       = widget.Part("check-all")
	PartCheckAllTrash  = widget.Part("check-all-trash")
	PartCheckAllPencil = widget.Part("check-all-pencil")
	PartCheckAllCount  = widget.Part("check-all-count")
```
```go
	clsCheckAll       = NameTargetList.Class(PartCheckAll)
	clsCheckAllTrash  = NameTargetList.Class(PartCheckAllTrash)
	clsCheckAllPencil = NameTargetList.Class(PartCheckAllPencil)
	clsCheckAllCount  = NameTargetList.Class(PartCheckAllCount)
```

### 3b. `itemIDs()` helper

`CheckedIDs()` already builds `[]string{it.ID ...}` inline. Extract it so the
master check reuses it:

```go
func (t *TargetList) itemIDs() []string {
	ids := make([]string, len(t.items))
	for i, it := range t.items {
		ids[i] = it.ID
	}
	return ids
}
```

and change `CheckedIDs()` to `return t.sel.CheckedIDs(t.itemIDs())`.

### 3c. `buildMasterCheck()`

```go
func (t *TargetList) buildMasterCheck() *Element {
	allChecked := DeriveBool(func() bool {
		_ = t.sel.Changed().Get() // re-read after every toggle (see Mode.Changed)
		n := t.sel.Count()
		return n > 0 && n == len(t.items)
	})
	someChecked := DeriveBool(func() bool {
		_ = t.sel.Changed().Get()
		n := t.sel.Count()
		return n > 0 && n < len(t.items)
	})

	m := Span().Set(clsCheckAll.AsAttr()).
		Attr("role", "checkbox").
		BindAttrBool("aria-checked", allChecked).
		// glyph selector: trash while the danger tone is armed, pencil
		// otherwise — mirrors the row check's Invalid/Selected split, but
		// keyed on the MODE (sel.Danger), not on "this row is checked".
		BindState(widget.Invalid, DeriveBool(func() bool { return t.sel.On().Get() && t.sel.Danger().Get() })).
		BindState(widget.Selected, DeriveBool(func() bool { return t.sel.On().Get() && !t.sel.Danger().Get() })).
		// fill: solid when ALL marked, lighter wash when SOME marked.
		BindState(widget.Locked, allChecked).
		BindState(widget.Busy, someChecked).
		Child(trash.Ref.Render(string(clsCheckAllTrash))).
		Child(pencil.Ref.Render(string(clsCheckAllPencil))).
		Child(Span().Set(clsCheckAllCount.AsAttr()).
			BindTextFunc(func() string {
				_ = t.sel.Changed().Get()
				return fmt.Sprintf("%d / %d", t.sel.Count(), len(t.items))
			}))

	m.On("click", func(Event) {
		if n := t.sel.Count(); n > 0 && n == len(t.items) {
			t.sel.Clear()
			return
		}
		t.sel.CheckAll(t.itemIDs())
	})
	return m
}
```

`targetdate` / `targethour`: same body, `TargetDate`/`TargetHour` receiver,
`clsCheckAll*` from their own `cls` vars. `fmt` is already imported in all
three.

### 3d. `Render()` — master check as first child of the root

```go
func (t *TargetList) Render() *Element {
	list := Ul().Set(clsList.AsAttr()).Attr("role", "listbox").BindChildren(t.rows)

	return Div().Set(clsListWrap.AsAttr()).
		BindStateFunc(widget.Open, func() bool { return t.sel.On().Get() }).
		Child(t.buildMasterCheck()).
		Child(list)
}
```

(`targetdate` / `targethour` `Render()` are the same shape — add the
`.Child(t.buildMasterCheck())` before `.Child(list)`.)

### 3e. `css.go` — call `ApplyMaster`

In each `sheet()`, right after the existing `listselect.Apply(...)` line:

```go
	listselect.ApplyMaster(s, PartCheckAll, PartCheckAllTrash, PartCheckAllPencil, PartCheckAllCount)
```

Nothing else in `css.go` changes.

### 3f. `svg.go` — unchanged

The master check reuses `trash.Ref` / `pencil.Ref`; `IconSvg()` already ships
`trash.Def()` + `pencil.Def()`. Do NOT add anything.

---

## Stage 4 — tests for the three components

For **each** of `targetlist`, `targetdate`, `targethour`, add
`mastercheck_test.go` (`//go:build !wasm`, `package <name>`):

```go
//go:build !wasm

package targetlist

import (
	"strings"
	"testing"
)

// The master check is the first element in the root, hidden until selection
// mode opens.
func TestTargetList_MasterCheckHiddenUntilSelectionMode(t *testing.T) {
	css := (&TargetList{}).RenderCSS().String()
	i := strings.Index(css, ".targetlist__check-all {")
	if i == -1 {
		t.Fatal("expected a rule for .targetlist__check-all")
	}
	body := css[i:]
	if e := strings.Index(body, "}"); e != -1 {
		body = body[:e]
	}
	if !strings.Contains(body, "display: none") {
		t.Errorf("the master check must be hidden by default, block:\n%s", body)
	}
	if !strings.Contains(css, `.targetlist[data-open="true"] .targetlist__check-all {`) {
		t.Errorf("the master check must be revealed by the list's open state")
	}
}

// Tapping it with nothing / some marked selects every row; tapping it with
// all marked clears.
func TestTargetList_MasterCheckTogglesAll(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	tl.SetItems([]Item{{ID: "1"}, {ID: "2"}, {ID: "3"}})
	tl.SetSelectMode(true)

	m := tl.buildMasterCheck()
	fire := func() { m = tl.buildMasterCheck() } // rebuild to re-read state
	_ = fire

	// simulate the click handler directly (no DOM under SSR)
	if n := tl.sel.Count(); n > 0 && n == 3 {
		tl.sel.Clear()
	} else {
		tl.sel.CheckAll(tl.itemIDs())
	}
	if tl.sel.Count() != 3 {
		t.Fatalf("first tap must select all, Count = %d", tl.sel.Count())
	}
	if n := tl.sel.Count(); n > 0 && n == 3 {
		tl.sel.Clear()
	} else {
		tl.sel.CheckAll(tl.itemIDs())
	}
	if tl.sel.Count() != 0 {
		t.Fatalf("second tap must clear, Count = %d", tl.sel.Count())
	}
}

// The count label renders "n / total".
func TestTargetList_MasterCheckShowsNOfTotal(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	tl.SetItems([]Item{{ID: "1"}, {ID: "2"}})
	tl.SetSelectMode(true)
	tl.sel.CheckAll([]string{"1"})

	html := tl.buildMasterCheck().String()
	if !strings.Contains(html, "1 / 2") {
		t.Errorf("master check must show \"1 / 2\", got:\n%s", html)
	}
}

// The master reuses the shared glyphs, not a bespoke tick.
func TestTargetList_MasterCheckUsesSharedGlyphs(t *testing.T) {
	tl := &TargetList{}
	tl.Init(nil)
	html := tl.buildMasterCheck().String()
	if !strings.Contains(html, `href='#trash'`) || !strings.Contains(html, `href='#pencil'`) {
		t.Errorf("master check must reference the shared trash/pencil glyphs:\n%s", html)
	}
}
```

(`targetdate` / `targethour`: rename `targetlist`→`targetdate`/`targethour`,
`TargetList`→`TargetDate`/`TargetHour`, `.targetlist__`→`.targetdate__`/`.targethour__`.)

Also verify `conformance_test.go`'s `TestKindAllowsEveryState` still passes —
the master check binds `Selected`, `Invalid`, `Locked`, `Busy`; the list root's
`WidgetKind` is `widget.Combobox`, whose `Allows` covers `Selected`, `Invalid`
and the universal `Locked`/`Busy` (`widget/kind.go`). If a bind you add is not
allowed for `Combobox`, that test fails — do NOT add `widget.Current` (it is
NOT allowed for `Combobox`).

---

## Stage 5 — documentation (verify against your code)

- **`targetlist` / `targetdate` / `targethour` READMEs** — add a "Master check
  (select all)" paragraph: appears top-end of the list in selection mode,
  tri-state (none / some / all), reuses the mode glyph, shows `n / total`.
  (`targethour/README.md` exists; `targetlist`/`targetdate` have no folder
  README — do NOT create one for those, just skip.)
- **`components/docs/CATALOG.md`** — under the target* entries (or the
  `listselect` mention), one line that the record lists now carry a master
  select-all check.
- Do NOT reference this `docs/PLAN.md` from any permanent doc.

---

## Acceptance criteria

```bash
go build ./...
gotest                                                          # vet + tests + race + wasm, all green
GOOS=js GOARCH=wasm go list -deps ./... | grep tinywasm/svg/sprite   # prints NOTHING
grep -n "func (m \*Mode) CheckAll" listselect/listselect.go     # present
grep -n "func (m \*Mode) Clear"    listselect/listselect.go     # present
grep -rn "buildMasterCheck" targetlist/ targetdate/ targethour/ # present in all three
grep -n "ApplyMaster" listselect/css.go                         # present
grep -rn "ApplyMaster" targetlist/css.go targetdate/css.go targethour/css.go  # called by all three
```

- Every one of the four `mastercheck_test.go` tests green for all three
  components.
- `listselect_test.go`'s four new `CheckAll`/`Clear` tests green.
- `crudview` (in `tinywasm/layout`) is NOT touched and NOT imported.

## Stages table

| # | File(s) | Action |
|---|---|---|
| 1 | `listselect/listselect.go`, `listselect/listselect_test.go` | `Mode.CheckAll` / `Clear` / `Count` + tests |
| 2 | `listselect/css.go` | `ApplyMaster` skin function |
| 3 | `targetlist/{targetlist.go,css.go}`, `targetdate/{targetdate.go,css.go}`, `targethour/{targethour.go,css.go}` | new parts, `itemIDs()`, `buildMasterCheck()`, render it, call `ApplyMaster` |
| 4 | `targetlist/mastercheck_test.go`, `targetdate/mastercheck_test.go`, `targethour/mastercheck_test.go` | per-component tests |
| 5 | `targethour/README.md`, `docs/CATALOG.md` | document the master check |
