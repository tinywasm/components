//go:build wasm

package themeswitch

import "github.com/tinywasm/dom"

// OnMount restaura el tema guardado en localStorage al montarse en el DOM.
// Si el storage no está disponible o la entrada está corrupta, sale limpiamente
// sin modificar el tema (modo auto por defecto).
func (t *ThemeSwitch) OnMount() {
	if !dom.LocalStorageAvailable() {
		return
	}
	saved, err := dom.LocalStorageGet(storageKey)
	if err != nil || saved == "" {
		return
	}
	theme := Theme(saved)
	if !valid(theme) {
		dom.LocalStorageDel(storageKey) // best-effort cleanup, error ignorado
		return
	}
	dom.SetDocumentAttr("data-theme", string(theme))
	t.Update()
}

func (t *ThemeSwitch) onClick(dom.Event) {
	current := Theme(dom.GetDocumentAttr("data-theme"))
	next := cycle(current)
	dom.SetDocumentAttr("data-theme", string(next)) // "" elimina el atributo para ThemeAuto
	// Persistencia best-effort: el tema se aplica aunque el storage falle.
	if next == ThemeAuto {
		dom.LocalStorageDel(storageKey) // error ignorado — tema ya aplicado
	} else {
		dom.LocalStorageSet(storageKey, string(next)) // error ignorado — tema ya aplicado
	}
	t.Update()
}
