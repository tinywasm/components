# tinywasm/components — Plan: Signal-Driven Components

> **Master:** tinywasm/docs/PLAN.md · **Engine:** tinywasm/dom/docs/PLAN.md
> **Module:** `github.com/tinywasm/components`
> **Type:** Breaking-aligned migration. Removes the footgun that started this effort.

---

## Prerequisites

```bash
# Canonical test runner (WASM tests run against a real DOM). Required: external agents have no global gotest.
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

## Development Rules

- **Documentation First:** update the `components` skill / standards to teach the single contract
  (`Render()` once + `Signal` bindings; optional `Init(ctx)`; no `Update`/`OnMount`) before code.
- **WASM only:** reactive code in `//go:build wasm` files; keep backend stubs compiling.
- **TinyGo idioms:** `switch` not `map`; embed `Element` by value.
- **Tests:** `gotest` (never `go test`); stdlib only; dual WASM/stdlib. Publish with `gopush 'msg'`.
- **Minimal public API:** export only what a component *user* types; unexport anything only this package uses (helpers, field models, single-use constants). State lives in unexported fields exposed via signals.

## Signals API recap (from the dom engine — self-contained)

```go
s := dom.NewString("v"); s.Get(); s.Set("x")          // observable string cell; Set patches bound nodes
b := dom.NewBool(false); b.Toggle()                   // observable bool cell (Toggle flips it)
el.BindText(s); el.BindClass("on", b)                 // raw bindings (typed, no generics)
el.BindAttr("title", s); el.BindAttrBool("disabled", b)
el.BindTextFunc(func() string { ... })                // computed binding, AUTO-TRACKED (no deps list)
in.Bind(s); in.Autofocus()                            // two-way input; focus-on-appear
dom.Show(b, func() *Element { ... })                  // mount/unmount subtree
ul.BindChildren(rows)                                 // rows := dom.NewNodes(...); keyed list, surgical rows
// dom.DeriveString(func() string {...}) → named shared computed (auto-tracked); + DeriveBool
```

State the UI shows lives in a typed `Signal` (**no generics**). No `Update()`. `Init(ctx dom.Ctx)`
runs once before render.

---

## Change 1 — `themetoggle`: `theme` signal + `Init`

The original bug was `Update()` in `OnMount` recursing. With signals it is structurally impossible.

In themetoggle/themeswitch_wasm.go and `themeswitch.go`:

- Replace the implicit `data-theme`-as-state with an explicit signal on the struct:

```go
type ThemeToggle struct {
	dom.Element
	theme *dom.SignalString // "", "dark", "light"
}

func (t *ThemeToggle) Init(_ dom.Ctx) {
	t.theme = dom.NewString("") // ThemeAuto
	if dom.LocalStorageAvailable() {
		if s, err := dom.LocalStorageGet(storageKey); err == nil && valid(Theme(s)) {
			t.theme.Set(s) // value ready before first paint → correct icon, no flash
		}
	}
}

func (t *ThemeToggle) Render() *Element {
	// labelSig is used twice (title + aria-label) → a named shared computed. Auto-tracked: no deps list.
	labelSig := dom.DeriveString(func() string { return label(Theme(t.theme.Get())) })
	return Button().
		BindTextFunc(func() string { return icon(Theme(t.theme.Get())) }). // computed icon, auto-tracked
		BindAttr("title", labelSig).BindAttr("aria-label", labelSig).
		Set(clsTsBtn.AsAttr()).                                       // typed builder: Set(...fmt.KeyValue), no Add(...any)
		On("click", func(Event) {
			next := cycle(Theme(t.theme.Get()))
			dom.SetDocumentAttr("data-theme", string(next))    // applies the theme
			if next == ThemeAuto { dom.LocalStorageDel(storageKey) } else { dom.LocalStorageSet(storageKey, string(next)) }
			t.theme.Set(string(next))                          // patches icon + labels surgically
		})
}
```

- Delete `OnMount` and the `t.Update()` call. Keep `cycle`/`icon`/`label`/`valid` (switch-based).

## Change 2 — `selectsearch`: signals + `Show` + `.Autofocus()`

In selectsearch/selectsearch.go:

- Move state to signals: `isOpen *SignalBool`, `query *SignalString`, `selected *SignalString`
  (init in `Init`). Options can stay a field if static.
- Delete all three `c.Update()` calls (lines 60, 101, 122) and `OnMount` (lines 142-148).
- Wrap the dropdown in `Show(c.isOpen, …)`; mark the search input `.Bind(c.query).Autofocus()` so it
  focuses when the dropdown mounts and keeps focus + IME while typing (node never replaced).
- Handlers only `Set` signals: toggle → `c.isOpen.Set(true/false)`; item click → set `selected`,
  `c.isOpen.Set(false)`, call `c.OnSelect`.
- Render the filtered list with `ul.BindChildren(c.rows)`, where `c.rows *dom.SignalNodes` is rebuilt
  (`c.rows.Set(buildRows(options, c.query.Get()))`) inside the input handler, so typing patches only
  changed rows (keyed reconcile).

---

## Documentation (do FIRST)

- **`devflow/skills/components/SKILL.md`**: replace any `OnMount`/`Update()` guidance with the single
  contract (`Render()` once + typed signals; `Init(ctx)`; `BindChildren`/`Show`; `.Autofocus()`).
  These two components are the **canonical examples** — reference them. After editing, run the
  `llmskill` sync so Claude/Gemini copies update.
- **`docs/CATALOG.md`** / **`docs/SKILL.md`**: update the themetoggle & selectsearch entries to the
  signal-based API.
- Each component's `README.md`/doc comment: show the `Signal`-based usage. Re-index `README.md`.

## Tests — frequent use cases (`gotest`)

Stdlib assertions only; dual WASM/stdlib. These cover the everyday component patterns and serve as
living examples:

- **stdlib:** `cycle`/`icon`/`label`/`valid`; `selectsearch` filter logic over `query`.
- **wasm (real DOM):**
  - **themetoggle — load-on-init + derived:** after `Init` with a saved theme, the first-rendered
    button text == correct icon (no flash); click cycles, only the button's text/attrs patch (capture
    node ref, assert identity); **regression: no recursion** (the original bug).
  - **selectsearch — two-way input + conditional + list:** typing keeps focus/cursor (IME-safe);
    `Show` mounts/unmounts the dropdown once; `BindChildren` patches only changed rows on filter.
- **In-browser (tinywasm MCP):** themetoggle `data-theme` matches `.ts-btn` icon on load and cycles;
  `browser_get_errors` shows **no** `Maximum call stack size exceeded`; selectsearch focus retained.

## Done When

- Neither component implements `OnMount`/`OnUpdate` nor calls `Update()`; state is in signals.
- themetoggle icon correct on load and cycles; selectsearch focus + IME survive keystrokes.
- **Docs:** `components` SKILL + CATALOG updated and synced (`llmskill`); READMEs re-indexed.
  **Tests:** the use-case tests above pass under `gotest`.
