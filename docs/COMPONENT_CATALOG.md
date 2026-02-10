# Components Catalog

This catalog documents the available reusable components in `tinywasm/components`.

## Overview

All components follow the [Component Creation Guide](./COMPONENT_CREATION.md).

### Usage

```go
import "github.com/tinywasm/components/button"

// Instantiate
btn := &button.Button{
    Text: "Click me",
    Variant: "primary",
}

// Render
dom.Append("app", btn)
```

---

## 1. Button

Versatile button component with variant support.

### Import
`"github.com/tinywasm/components/button"`

### Usage

```go
btn := &button.Button{
    Text: "Save Changes",
    Variant: "primary",
    OnClick: func(e dom.Event) {
        fmt.Println("Saved!")
    },
}
```

### Properties
- `Text` (string): The button label.
- `Variant` (string): Visual style (`primary`, `secondary`, `danger`).
- `OnClick` (func(dom.Event)): Click handler.

### CSS Variables
- `--color-primary`
- `--color-secondary`
- `--color-error`

---

## 2. Card

Container with header, body, and footer sections.

### Import
`"github.com/tinywasm/components/card"`

### Usage

```go
c := &card.Card{
    Header: dom.Tag("h3").Text("Card Title").ToComponent(),
    Body:   dom.Tag("p").Text("Content goes here...").ToComponent(),
    Footer: dom.Tag("small").Text("Last updated: today").ToComponent(),
}
```

### Properties
- `Header` (dom.Component): Content for header section.
- `Body` (dom.Component): Main content.
- `Footer` (dom.Component): Footer content.

---

## 3. Input

Text input with label and validation support.

### Import
`"github.com/tinywasm/components/input"`

### Usage

```go
inp := &input.Input{
    Label: "Username",
    Placeholder: "Enter username",
    Type: "text",
    Required: true,
    OnInput: func(e dom.Event) {
        // Handle input
    },
}
```

### Properties
- `Label` (string): Label text.
- `Placeholder` (string): Placeholder text.
- `Type` (string): HTML input type (`text`, `email`, `password`, etc).
- `Value` (string): Initial value.
- `Required` (bool): Whether input is required.
- `OnInput` (func(dom.Event)): Input handler.

---

## 4. Nav

Navigation menu with support for icons.

### Import
`"github.com/tinywasm/components/nav"`

### Usage

```go
n := &nav.Nav{
    Items: []nav.NavItem{
        {Label: "Dashboard", Route: "dashboard", Icon: "icon-home"},
        {Label: "Settings", Route: "settings", Icon: "icon-cog"},
    },
}
```

### Properties
- `Items` ([]NavItem): List of navigation items.
    - `Label` (string): Display text.
    - `Route` (string): Target route/hash.
    - `Icon` (string): SVG sprite ID for icon (optional).

---

## 5. Modal

Dialog overlay with backdrop and close button.

### Import
`"github.com/tinywasm/components/modal"`

### Usage

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

### Properties
- `Title` (string): Modal title.
- `Content` (dom.Component): Modal body content.
- `Visible` (bool): Initial visibility state.

### Methods
- `Open()`: Shows the modal.
- `Close(e dom.Event)`: Hides the modal.

---

## 6. Table

Data table.

### Import
`"github.com/tinywasm/components/table"`

### Usage

```go
t := &table.Table{
    Headers: []string{"ID", "Name", "Role"},
    Rows: [][]string{
        {"1", "Alice", "Admin"},
        {"2", "Bob", "User"},
    },
}
```

### Properties
- `Headers` ([]string): Column headers.
- `Rows` ([][]string): Data rows.

---

## 7. Form

Wrapper for form fields.

### Import
`"github.com/tinywasm/components/form"`

### Usage

```go
f := &form.Form{
    Fields: []dom.Component{
        &input.Input{Label: "Email"},
        &button.Button{Text: "Submit"},
    },
    OnSubmit: func(e dom.Event) {
        e.PreventDefault()
        // Handle submit
    },
}
```

### Properties
- `Fields` ([]dom.Component): List of form components.
- `OnSubmit` (func(dom.Event)): Submit handler.
