# Agent Guide — `tinywasm/components`

Constraints for agents adding or modifying UI components. Read this before any change.

---

## Construction Harness — typed & explicit (the TinyWasm approach)

This library is part of TinyWasm's **construction harness**: the typed, explicit API is what keeps an
agent that doesn't know the library from building wrong code. The compiler must reject mistakes; what
it can't catch becomes a `devMode` warning — never a silent failure.

- **Typed over `any`** — no generic slots; typed builder methods (like `tinywasm/json`), reusing `fmt` types. Anything reactive goes only through a signal binding (`BindText`/`Bind*`), which requires a signal.
- **Explicit names** — `Text` (static) vs `BindText` (reactive); reading the call states intent.
- **Illegal states unrepresentable** — dynamic content has ONE path, typed to require a signal.
- **Minimal public surface** — export only what the author types; engine plumbing stays unexported.
- **Docs are minimal "how" instructions, not long skills** — if a rule must be *remembered*, close it
  with types, not prose.

(Ecosystem rationale: `tinywasm/docs/ARNES_DE_CONSTRUCCION.md`.)

---

## Component Contract — ONE way (signals)

A component implements **only**:

- `Render() *dom.Element` — describes the structure ONCE; dynamic parts are signal bindings.
- `Init(ctx dom.Ctx)` — optional; runs ONCE before first render (load storage, fetch, subscribe).

There is **NO** `OnMount`/`OnUpdate`/`OnUnmount` and **NO** manual `Update()` (it is unexported in
`dom`). Forgetting a lifecycle call is impossible because there are no lifecycle calls.

State the UI shows lives in **typed signals**; changing one patches only the bound DOM node — never
re-render a whole element, never use a Virtual DOM.

```go
type ThemeToggle struct {
    dom.Element                 // value-embed, never *dom.Element
    theme *dom.SignalString     // UI state is a signal, not a plain field
}
func (t *ThemeToggle) Init(_ dom.Ctx) { t.theme = dom.NewString("") /* load saved */ }
func (t *ThemeToggle) Render() *dom.Element {
    return dom.Button().
        BindTextFunc(func() string { return iconFor(t.theme.Get()) }).   // auto-tracked; no deps list
        On("click", func(dom.Event) { t.theme.Set(next(t.theme.Get())) })
}
```

## No Generics

Zero generic functions in this ecosystem (follow the `tinywasm/fmt` codec rule "cero any, cero map").
The DOM boundary is `string`/`bool`, so use concrete typed signals — never `Signal[T]`:

- `SignalString`/`SignalBool`/`SignalNodes`; `NewString`/`NewBool`/`NewNodes`; `Get`/`Set`; `Toggle()` on bool
- Bindings (raw signal): `BindText`, `BindAttr`, `BindClass`, `BindAttrBool`, `Bind` (two-way input),
  `BindChildren` (keyed list — build `[]*Element` in a normal loop, set with `.Key(id)`), `Autofocus`
- Bindings (computed, **auto-tracked, no deps list**): `BindTextFunc`/`BindAttrFunc`/`BindClassFunc`/`BindAttrBoolFunc`;
  `DeriveString`/`DeriveBool` for a named shared computed value
- `Show(boolSig, render)` for conditional subtrees.

## Minimal Public Surface

Export only what a component user types. **Unexport any symbol only this package uses** — e.g. theme
constants (`ThemeDark`/`ThemeLight`/`ThemeAuto`), helpers like `cycle`/`icon`/`label`. Struct fields
stay unexported; expose behavior through methods.

## WASM / TinyGo

- Lifecycle/reactive code goes in `//go:build wasm` files; provide `!wasm` stubs for any function
  called from tag-less code (e.g. inside `Render()`), returning no-ops / `""`.
- **No Go stdlib** (`fmt`/`strings`/`errors`): use `github.com/tinywasm/fmt`. DOM only via
  `github.com/tinywasm/dom`, never `syscall/js`.
- `switch`, not `map`. No `defer/recover` (a no-op in TinyGo WASM) — use O(1) guards.

## Testing

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest   # external agents have no global gotest
gotest
```

- `gotest`, never `go test`. Stdlib assertions only (`testing`/`reflect`, no testify).
- Dual WASM/stdlib via build tags sharing one runner; WASM tests run against a real DOM.
- Cover the **frequent use cases**: load-on-init (no flash), two-way input (IME/cursor safe), derived
  value, conditional `Show`, keyed list, and the no-recursion regression.
- Publish with `gopush 'message'`, never raw `git commit`/`push`.

## Documentation First

Update docs **before** code and before `gopush`: `docs/ARCHITECTURE.md` (what/why),
`docs/DESIGN.md` (decisions), `docs/CATALOG.md`, and the component standards in the `components`
skill (`devflow/skills/components/SKILL.md` — run the `llmskill` sync after editing). `README.md`
must index every file in `docs/`. Diagrams: `flowchart TD`, no `subgraph`, `<br/>` for line breaks.
