# herobanner

`herobanner` is a prominent header component for site landing pages, featuring background images, main title, subtitle, and call-to-action buttons.

## Usage

```go
import (
    "github.com/tinywasm/components/actionbutton"
    "github.com/tinywasm/components/herobanner"
)

hero := &herobanner.HeroBanner{
    Title:    "Salud y Cuidado Integral",
    Subtitle: "Especialistas dedicados a tu bienestar y el de tu familia.",
    Images:   []string{"/images/hero1.jpg", "/images/hero2.jpg"},
    Actions: []dom.Component{
        &actionbutton.ActionButton{Text: "Reserva tu hora", Variant: "primary"},
        &actionbutton.ActionButton{Text: "Ver especialidades", Variant: "secondary"},
    },
}
```
