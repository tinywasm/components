# PLAN — `components` module: upgrade to `tinywasm/css v0.1.0`

> Single goal: replace the `New(...)` constructor call with `NewStylesheet(...)` in every `ssr.go` after `tinywasm/css v0.1.0` is published. No RawRule issues exist in this module — all 7 components are clean.

**Blocked on**: `github.com/tinywasm/css v0.1.0` published (see `tinywasm/css/docs/PLAN.md`).

---

## 1. Scope

7 files, 1 change each — all on line 10 of each `ssr.go`:

| File | Change |
| --- | --- |
| `actionbutton/ssr.go` | `New(` → `NewStylesheet(` |
| `contentcard/ssr.go` | `New(` → `NewStylesheet(` |
| `datatable/ssr.go` | `New(` → `NewStylesheet(` |
| `dialog/ssr.go` | `New(` → `NewStylesheet(` |
| `navbar/ssr.go` | `New(` → `NewStylesheet(` |
| `selectsearch/ssr.go` | `New(` → `NewStylesheet(` |
| `themetoggle/ssr.go` | `New(` → `NewStylesheet(` |

---

## 2. Steps

```bash
# 1. bump dependency
go get github.com/tinywasm/css@v0.1.0
go mod tidy

# 2. rename (all at once)
grep -rn "\bNew(" --include="ssr.go" .

# 3. build + test
go build ./...
go test ./...

# 4. verify no old constructor remains
grep -rn "\bNew(" --include="ssr.go" .   # must return empty
```

---

## 3. Acceptance criteria

- `go.mod` references `tinywasm/css v0.1.0`.
- `go build ./...` and `go test ./...` green.
- `grep -rn "\bNew(" --include="ssr.go" .` returns **empty**.
- No `RawRule` introduced (zero existed before; zero after).
