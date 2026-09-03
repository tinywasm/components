# CountBadge

Notification bubble: a count riding the top-end corner of the button that
owns it. Out of the flow (`OnEdge`), so the host keeps its box whether the
bubble shows "1", "99" or nothing at all.

## Import
`"github.com/tinywasm/components/countbadge"`

## Usage

```go
import "github.com/tinywasm/components/countbadge"

count := dom.NewString("0")
visible := dom.NewBool(false)

span := (&countbadge.CountBadge{Count: count, Visible: visible}).Render()
```

## Properties
- `Count` (*dom.SignalString): The number to show, as text. A signal
  because the count changes without re-rendering the host.
- `Visible` (*dom.SignalBool): Whether the bubble paints. Showing "0" is
  noise — wire something false at zero (e.g. `n > 0`).

## Host contract
The host button must be an `Anchor` (`style.Anchor()` — `position:
relative`, zero visual change) so the bubble resolves against the button,
not the page.
