//go:build !wasm

package sitenav

// RenderJS returns the vanilla JavaScript snippet responsible for the mobile
// menu toggle behavior and link click auto-close.
func (sn *SiteNav) RenderJS() string {
	return `(function() {
	if (window.__sitenavInit) return;
	window.__sitenavInit = true;

	document.addEventListener('click', function(e) {
		var toggle = e.target.closest('[aria-controls="sitenav-menu"]');
		if (toggle) {
			var isExpanded = toggle.getAttribute('aria-expanded') === 'true';
			var nextState = !isExpanded;
			toggle.setAttribute('aria-expanded', String(nextState));
			toggle.setAttribute('aria-label', nextState ? 'Cerrar menú de navegación' : 'Abrir menú de navegación');
			var menu = document.getElementById('sitenav-menu');
			if (menu) {
				if (nextState) {
					menu.classList.add('is-open');
				} else {
					menu.classList.remove('is-open');
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
				}
				menu.classList.remove('is-open');
			}
		}
	});
})();`
}
