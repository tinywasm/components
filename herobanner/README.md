# herobanner

`herobanner` is a prominent header component for site landing pages, featuring background images, main title, subtitle, and call-to-action buttons.

## Usage

```go
import (
    "webtyp.com/components/actionbutton"
    "webtyp.com/components/herobanner"
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

> **Nota sobre `Images`**: El campo `Images` recibe **rutas base sin sufijo de variante** (por ejemplo, `/images/hero1.jpg`). La renderización genera automáticamente los atributos `srcset` con las tres variantes (`.S.jpg`, `.M.jpg`, `.L.jpg`) y `sizes="100vw"`.
