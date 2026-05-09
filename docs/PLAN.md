# PLAN: Migración de tokens CSS v2 — `tinywasm/components`

## Contexto

Migración consecuente al breaking change de `tinywasm/dom` que adopta la
convención Material Design 3 / Bootstrap 5 con patrón two-layer variables.
Todos los archivos CSS de este módulo referencian los tokens viejos y deben
actualizarse a los nuevos.

**Prerequisito bloqueante:** completar `tinywasm/dom` → `dom/docs/PLAN.md`
antes de ejecutar este plan.

---

## Tabla de sustitución

| Token viejo | Token nuevo | Cuándo |
|---|---|---|
| `--color-primary` (como texto) | `--color-on-surface` | texto sobre cualquier superficie |
| `--color-primary` (como texto sobre cyan) | `--color-on-primary` | texto cuyo fondo es `--color-primary` |
| `--color-secondary` | `--color-primary` | fondo cyan / acento de marca |
| `--color-selection` | `--color-secondary` | fondo purple / acento interactivo |
| `--color-gray` (como texto sobre purple) | `--color-on-secondary` | texto cuyo fondo es `--color-secondary` |
| `--color-tertiary` | `--color-muted` | texto apagado / bordes sutiles |
| `--color-quaternary` | `--color-surface` | fondo de paneles / cards |
| `--color-gray` (como fondo) | `--color-background` | fondo de página / inputs |
| `--color-hover` | `--color-hover` | sin cambio |

**Regla clave:** el token `on-X` se usa cuando el elemento tiene
`background: var(--color-X)`. En cualquier otro contexto con texto: `--color-on-surface`.

---

## Archivos afectados y cambios requeridos

### `button/button.css`

| Selector | Propiedad | Viejo | Nuevo |
|---|---|---|---|
| `button[name*="btn"]` | `background-color` | `--color-secondary` | `--color-primary` |
| `button[name*="btn"]` | `color` | `--color-primary` | `--color-on-primary` |
| `button[name*="btn"]:disabled` | `background-color` | `--color-tertiary` | `--color-muted` |
| `button[name*="btn"]:hover` | `background-color` | `--color-selection` | `--color-secondary` |
| `.contebuton` | `background` | `--color-primary` | `--color-surface` |
| `.btn-url, .btn-login` | `background` | `--color-secondary` | `--color-primary` |
| `.btn-url-down` | `background` | `--color-gray` | `--color-background` |
| `.btn-url-down svg path` | `fill` | `--color-secondary` | `--color-primary` |
| `.btn-url-disable` | `background` | `--color-tertiary` | `--color-muted` |
| `.btn-selected` | `background` | `--color-selection` | `--color-secondary` |
| `.btn-url, .btn-selected, ...` | `box-shadow` | `--color-quaternary` | `--color-surface` |
| `.btn-url-pulse` | `box-shadow` | `--color-quaternary` | `--color-surface` |
| `@keyframes pulse-url 0%` | `box-shadow` | `--color-selection` | `--color-secondary` |
| `@keyframes pulse-url 100%` | `background` | `--color-selection` | `--color-secondary` |
| `input[type=submit]` | `background-color` | `--color-secondary` | `--color-primary` |
| `input[type=submit]` | `color` | `--color-primary` | `--color-on-primary` |
| `.btn-primary` | `background` | `--color-secondary` | `--color-primary` |
| `.btn-primary` | `color` | `--color-gray` | `--color-on-primary` |
| `.btn-secondary` | `background` | `--color-tertiary` | `--color-muted` |
| `.btn-secondary` | `color` | `--color-primary` | `--color-on-surface` |
| `.btn-danger` | `color` | `--color-primary` | `--color-on-surface` |

---

### `card/card.css`

| Selector | Propiedad | Viejo | Nuevo |
|---|---|---|---|
| `.card` | `background-color` | `--color-gray` | `--color-background` |
| `.card-header, .card-body, .card-footer` | `color` | `--color-primary` | `--color-on-surface` |
| `.card, .card-header, .card-footer` | `border` | `--color-tertiary` | `--color-muted` |
| `.card-footer` | `background-color` | `--color-quaternary` | `--color-surface` |

---

### `modal/modal.css`

| Selector | Propiedad | Viejo | Nuevo |
|---|---|---|---|
| `.modal-backdrop` | `background-color` | `--color-quaternary` | `--color-surface` |
| `.modal-content` | `background-color` | `--color-gray` | `--color-background` |
| `.modal-content, .modal-close` | `color` | `--color-primary` | `--color-on-surface` |
| `.modal-header` | `border` | `--color-tertiary` | `--color-muted` |

---

### `nav/nav.css`

| Selector | Propiedad | Viejo | Nuevo |
|---|---|---|---|
| `.nav` | `background-color` | `--color-quaternary` | `--color-surface` |
| `.nav` | `color` | `--color-primary` | `--color-on-surface` |
| `.nav-item a:hover` | `color` | `--color-secondary` | `--color-primary` |

---

### `selectsearch/selectsearch.css`

| Selector | Propiedad | Viejo | Nuevo |
|---|---|---|---|
| `.ss-header` | `background` | `--color-secondary` | `--color-primary` |
| `.ss-header` | `color` | `--color-primary` | `--color-on-primary` |
| `.ss-search` | `background` | `--color-gray` | `--color-background` |
| `.ss-search` | `color` | `--color-primary` | `--color-on-surface` |
| `.ss-search` | `border` | `--color-tertiary` | `--color-muted` |
| `.ss-options` | `background` | `--color-quaternary` | `--color-surface` |
| `.ss-option` | `color` | `--color-primary` | `--color-on-surface` |
| `.ss-option` | `border-bottom` | `--color-tertiary` | `--color-muted` |
| `.ss-desc` | `background` | `--color-quaternary` | `--color-surface` |
| `.ss-desc` | `color` | `--color-primary` | `--color-on-surface` |
| `.ss-options scrollbar-thumb` | `background` | `--color-secondary` | `--color-primary` |

---

### `table/table.css`

| Selector | Propiedad | Viejo | Nuevo |
|---|---|---|---|
| `.table` | `color` | `--color-primary` | `--color-on-surface` |
| `.table th` | `background-color` | `--color-quaternary` | `--color-surface` |
| `.table th, .table td` | `border-bottom` | `--color-tertiary` | `--color-muted` |

---

### `themeswitch/themeswitch.css`

**Los bloques `[data-theme]` se eliminan de este archivo** — ahora viven
exclusivamente en `dom/theme.css`. Este archivo queda solo con estilos del botón.

| Selector | Propiedad | Viejo | Nuevo |
|---|---|---|---|
| `.ts-btn` | `background` | `--color-secondary` | `--color-primary` |
| `.ts-btn` | `color` | `--color-primary` | `--color-on-primary` |
| `.ts-btn:hover` | `background` | `--color-selection` | `--color-secondary` |
| `.ts-btn:hover` | `color` | `--color-gray` | `--color-on-secondary` |

**CSS resultante de `themeswitch.css` después de la migración:**

```css
.ts-btn {
    position: fixed;
    top: 1rem;
    right: 1rem;
    z-index: 9999;
    width: 2.4rem;
    height: 2.4rem;
    padding: 0;
    border-radius: 50%;
    border: none;
    cursor: pointer;
    font-size: 1.1rem;
    line-height: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-primary);
    color: var(--color-on-primary);
    opacity: 0.85;
    transition: opacity 0.2s, transform 0.2s;
}

.ts-btn:hover {
    opacity: 1;
    transform: scale(1.12);
    background: var(--color-secondary);
    color: var(--color-on-secondary);
}
```

---

## Checklist de implementación

**Prerequisito:** `tinywasm/dom` → `dom/docs/PLAN.md` completado y publicado.

- [ ] `button/button.css` — aplicar tabla de cambios
- [ ] `card/card.css` — aplicar tabla de cambios
- [ ] `modal/modal.css` — aplicar tabla de cambios
- [ ] `nav/nav.css` — aplicar tabla de cambios
- [ ] `selectsearch/selectsearch.css` — aplicar tabla de cambios
- [ ] `table/table.css` — aplicar tabla de cambios
- [ ] `themeswitch/themeswitch.css` — aplicar tabla + eliminar bloques `[data-theme]`
- [ ] Ejecutar `gotest` en cada componente con demo browser para validar visualmente
- [ ] `gopush 'feat(components)!: migrate CSS tokens to v2 (Material/Bootstrap convention)'`
