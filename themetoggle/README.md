# ThemeSwitch

Componente visual para alternar el tema de la aplicación (light / dark / auto).

## Características

- 3 estados: `auto` (preferencia del OS), `dark` y `light`.
- Persistencia automática en `localStorage`.
- Botón flotante fijo en la esquina superior derecha.
- Diseñado para pruebas de desarrollo de componentes en diferentes temas.

## Uso

```go
import (
	"github.com/tinywasm/components/themeswitch"
	"github.com/tinywasm/dom"
)

func main() {
	ts := &themeswitch.ThemeSwitch{}
	dom.Append("body", ts)
	select {}
}
```
