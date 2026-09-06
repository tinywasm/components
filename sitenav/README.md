# sitenav

`sitenav` provides a responsive header navigation bar for public websites, complete with responsive `<picture>` branding logos, navigation links, call-to-action components, and an accessible mobile hamburger toggle.

## Usage

```go
import (
    "webtyp.com/components/actionbutton"
    "webtyp.com/components/sitenav"
)

nav := &sitenav.SiteNav{
    WideLogoSrc:    "/assets/logo-wide.svg",
    CompactLogoSrc: "/assets/logo-compact.svg",
    LogoAlt:        "Clínica San José",
    Links: []sitenav.NavItem{
        {Label: "Inicio", Href: "/", Active: true},
        {Label: "Especialidades", Href: "/especialidades"},
        {Label: "Equipo Médico", Href: "/equipo"},
        {Label: "Contacto", Href: "/contacto"},
    },
    Actions: []dom.Component{
        &actionbutton.ActionButton{Text: "Reservar Cita", Variant: "primary"},
    },
}
```
