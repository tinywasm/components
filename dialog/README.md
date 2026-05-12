# Modal Component

Dialog overlay with backdrop and close button.

## Import
`"github.com/tinywasm/components/modal"`

## Usage

```go
m := &modal.Modal{
    Title: "Confirm Action",
    Content: dom.Tag("p").Text("Are you sure?").ToComponent(),
    Visible: false,
}

// To open
m.Open()

// To close (programmatically)
m.Close(dom.Event{})
```

## Properties
- `Title` (string): Modal title.
- `Content` (dom.Component): Modal body content.
- `Visible` (bool): Initial visibility state.

## Methods
- `Open()`: Shows the modal.
- `Close(e dom.Event)`: Hides the modal.

---
[Back to Catalog](../docs/CATALOG.md)
