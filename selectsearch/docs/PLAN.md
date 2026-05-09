# PLAN: Dark mode fix — `tinywasm/components/selectsearch`

## Contexto

El componente `selectsearch` usa tokens CSS de `tinywasm/dom` (`--color-*`) en la
mayoría de sus reglas. Sin embargo, dos elementos no definen sus propiedades de
color/background explícitamente, lo que hace que hereden el default del browser
(fondo blanco, texto oscuro) y rompan visualmente en modo dark.

---

## Problemas identificados en `selectsearch.css`

### 1. `.ss-search` — sin background ni color

```css
/* ACTUAL — hereda browser default: fondo blanco, texto oscuro */
.ss-search {
    width: 100%;
    border: 0.2em solid var(--color-tertiary);
    border-radius: 0.4em 0.4em 0 0;
    padding: var(--mag-pri);
    font-size: 1rem;
    box-sizing: border-box;
}
```

En modo dark, el fondo sigue siendo blanco → rompe la coherencia visual.

**Fix:**
```css
.ss-search {
    width: 100%;
    border: 0.2em solid var(--color-tertiary);
    border-radius: 0.4em 0.4em 0 0;
    padding: var(--mag-pri);
    font-size: 1rem;
    box-sizing: border-box;
    background: var(--color-gray);   /* ← añadir */
    color: var(--color-primary);     /* ← añadir */
}
```

`--color-gray` en dark = `#0D1117` (superficie oscura).
`--color-primary` en dark = `#E6EDF3` (texto claro).

---

### 2. `.ss-header { color: var(--color-gray) }` — contraste pobre en dark

```css
/* ACTUAL */
.ss-header {
    background: var(--color-secondary);   /* cyan #00ADD8 siempre */
    color: var(--color-gray);             /* dark: #0D1117 = casi negro sobre cyan */
}
```

En dark mode `--color-gray` es `#0D1117` (casi negro). Sobre cyan es técnicamente
legible (ratio ~4.5:1) pero la intención semántica del token es "superficie" no
"texto". El resultado es texto muy oscuro sobre un fondo de color vivo, lo que
luce incongruente con el resto del componente en dark mode.

**Fix: usar `--color-primary` para el texto del header**

```css
.ss-header {
    background: var(--color-secondary);
    color: var(--color-primary);   /* light: #1C1C1E / dark: #E6EDF3 */
}
```

`--color-primary` está semánticamente definido como "texto principal" para cada
modo. En dark es `#E6EDF3` (casi blanco) sobre cyan → contraste correcto.

---

### 3. `.ss-option:hover { color: var(--color-gray) }` — revisión

```css
.ss-option:hover { background: var(--color-hover); color: var(--color-gray); }
```

En dark, `--color-hover = #F7DF1E` (amarillo JS) y `--color-gray = #0D1117` (casi negro).
Negro sobre amarillo tiene buen contraste (ratio ~12:1) y es visualmente coherente
con la identidad del design system. **No requiere cambio.**

---

## Cambios requeridos

Solo se modifica `selectsearch.css`:

```diff
 .ss-header {
     background: var(--color-secondary);
-    color: var(--color-gray);
+    color: var(--color-primary);
     padding: var(--mag-pri) calc(var(--mag-pri) * 2);
     cursor: pointer;
     border-radius: 0.4em;
     display: flex;
     justify-content: space-between;
     align-items: center;
 }

 .ss-search {
     width: 100%;
     border: 0.2em solid var(--color-tertiary);
     border-radius: 0.4em 0.4em 0 0;
     padding: var(--mag-pri);
     font-size: 1rem;
     box-sizing: border-box;
+    background: var(--color-gray);
+    color: var(--color-primary);
 }
```

---

## Verificación

Una vez implementados `tinywasm/dom` (tema API) y `tinywasm/components/themeswitch`:

```go
import (
    . "github.com/tinywasm/dom"
    "github.com/tinywasm/components/themeswitch"
)

func main() {
    // dom.init() restaura el tema automáticamente al importar dom
    Render("app", &App{})
    Append("body", &themeswitch.ThemeSwitch{})
    select {}
}
```

El botón flotante (esquina superior derecha) permite ciclar entre
`auto / dark / light` y verificar el componente sin salir del browser.

---

## Checklist de implementación

- [ ] En `selectsearch.css`: cambiar `color: var(--color-gray)` → `color: var(--color-primary)` en `.ss-header`
- [ ] En `selectsearch.css`: añadir `background: var(--color-gray)` y `color: var(--color-primary)` en `.ss-search`
- [ ] Verificar visualmente en modo dark con `ThemeSwitch`
- [ ] Verificar que el modo light sigue correcto
- [ ] Actualizar `web/client.go` para usar `ThemeSwitch` como demo
