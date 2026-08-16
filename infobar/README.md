# infobar

`infobar` is a top contact/location bar component displaying contact information items (phone, email, address, working hours) with SVG icons.

## Usage

```go
import "github.com/tinywasm/components/infobar"

bar := &infobar.InfoBar{
    Items: []infobar.InfoItem{
        {Icon: infobar.IconPhone, Text: "+56 9 1234 5678", Href: "tel:+56912345678"},
        {Icon: infobar.IconMail, Text: "contacto@clinica.cl", Href: "mailto:contacto@clinica.cl"},
        {Icon: infobar.IconPin, Text: "Av. Providencia 1234"},
        {Icon: infobar.IconClock, Text: "Lun-Vie 8:00 - 20:00"},
    },
}
```
