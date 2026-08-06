module github.com/tinywasm/components

go 1.25.2

require (
	github.com/tinywasm/css v0.4.7
	github.com/tinywasm/dom v0.13.2
	github.com/tinywasm/fmt v0.25.5
	github.com/tinywasm/html v0.0.12
	github.com/tinywasm/svg v0.1.8
	github.com/tinywasm/widget v0.5.8
)

require (
	github.com/tinywasm/font v0.0.4 // indirect
	github.com/tinywasm/json v0.5.17 // indirect
	github.com/tinywasm/model v0.1.2 // indirect
)

// ── replaces de desarrollo local ─────────────────────────────────────────────
// PLAN v0.2.0 — the widget/style EdgeToEdge fix is exercised here before it is
// published. Revert when the work lands.
