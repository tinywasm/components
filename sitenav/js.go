//go:build !wasm

package sitenav

import "webtyp.com/js"

// RenderJS returns the vanilla JavaScript snippet responsible for the mobile
// menu toggle behavior and link click auto-close.
//
// The open state is written as data-open on the MENU, which is the attribute
// widget.Open resolves to and the one RenderCSS selects on. It used to toggle
// an "is-open" class instead — a name no stylesheet in the ecosystem has ever
// matched, so the button flipped a class nothing read and the menu never
// opened. aria-expanded stays on the button: that is the assistive-technology
// contract, and it lives on the control, not on what it controls.
func (sn *SiteNav) RenderJS() []*js.Script {
	return []*js.Script{{Content: `(function() {
	if (window.__sitenavInit) return;
	window.__sitenavInit = true;

	document.addEventListener('click', function(e) {
		var toggle = e.target.closest('[aria-controls="sitenav-menu"]');
		if (toggle) {
			var isExpanded = toggle.getAttribute('aria-expanded') === 'true';
			var nextState = !isExpanded;
			toggle.setAttribute('aria-expanded', String(nextState));
			toggle.setAttribute('aria-label', nextState ? 'Cerrar menú de navegación' : 'Abrir menú de navegación');
			// The button carries the state too, not just the menu: the glyph it
			// swaps between lives inside the button, and a descendant rule can
			// only reach it through an ancestor that holds the attribute.
			if (nextState) {
				toggle.setAttribute('data-open', 'true');
			} else {
				toggle.removeAttribute('data-open');
			}
			var menu = document.getElementById('sitenav-menu');
			if (menu) {
				if (nextState) {
					menu.setAttribute('data-open', 'true');
				} else {
					menu.removeAttribute('data-open');
				}
			}
			return;
		}

		var menu = document.getElementById('sitenav-menu');
		if (menu && menu.contains(e.target)) {
			var link = e.target.closest('a');
			if (link) {
				var btn = document.querySelector('[aria-controls="sitenav-menu"]');
				if (btn) {
					btn.setAttribute('aria-expanded', 'false');
					btn.setAttribute('aria-label', 'Abrir menú de navegación');
					btn.removeAttribute('data-open');
				}
				menu.removeAttribute('data-open');
			}
		}
	});
})();`}}
}
