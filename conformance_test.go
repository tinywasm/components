package components

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/tinywasm/components/actionbutton"
	"github.com/tinywasm/components/contentcard"
	"github.com/tinywasm/components/datatable"
	"github.com/tinywasm/components/fieldset"
	"github.com/tinywasm/components/modaldialog"
	"github.com/tinywasm/components/selectsearch"
	"github.com/tinywasm/components/targetlist"
	"github.com/tinywasm/components/themetoggle"
)

// Allowed CSS variables in github.com/tinywasm/css v0.2.0 (and widget/style)
var allowedVars = map[string]bool{
	"--color-primary":      true,
	"--color-on-primary":   true,
	"--color-secondary":    true,
	"--color-on-secondary": true,
	"--color-success":      true,
	"--color-on-success":   true,
	"--color-error":        true,
	"--color-on-error":     true,

	"--color-background":     true,
	"--color-surface":        true,
	"--color-surface-sunken": true,
	"--color-on-surface":     true,
	"--color-outline":        true,
	"--color-muted":          true,
	"--color-hover":          true,
	"--color-selection":      true,
	"--color-on-selection":   true,
	"--color-disabled":       true,
	"--color-on-disabled":    true,

	"--color-background-light":      true,
	"--color-background-dark":       true,
	"--color-surface-light":         true,
	"--color-surface-dark":          true,
	"--color-surface-sunken-light":  true,
	"--color-surface-sunken-dark":   true,
	"--color-on-surface-light":      true,
	"--color-on-surface-dark":       true,
	"--color-outline-light":         true,
	"--color-outline-dark":          true,
	"--color-muted-light":           true,
	"--color-muted-dark":            true,
	"--color-hover-light":           true,
	"--color-hover-dark":            true,
	"--color-selection-light":       true,
	"--color-selection-dark":        true,
	"--color-on-selection-light":    true,
	"--color-on-selection-dark":     true,
	"--color-disabled-light":        true,
	"--color-disabled-dark":         true,
	"--color-on-disabled-light":     true,
	"--color-on-disabled-dark":      true,

	"--color-background-hover": true,
	"--color-surface-hover":    true,
	"--color-primary-hover":    true,
	"--color-secondary-hover":  true,
	"--color-selection-hover":  true,
	"--color-success-hover":    true,
	"--color-error-hover":      true,
	"--color-muted-hover":      true,

	"--color-background-focus": true,
	"--color-surface-focus":    true,
	"--color-primary-focus":    true,
	"--color-secondary-focus":  true,
	"--color-selection-focus":  true,

	"--color-background-press": true,
	"--color-surface-press":    true,
	"--color-primary-press":    true,
	"--color-secondary-press":  true,
	"--color-selection-press":  true,

	"--text-xs":   true,
	"--text-sm":   true,
	"--text-base": true,
	"--text-lg":   true,
	"--text-xl":   true,
	"--text-2xl":  true,

	"--leading-tight":         true,
	"--leading-normal":        true,
	"--leading-relaxed":       true,
	"--font-weight-regular":   true,
	"--font-weight-medium":    true,
	"--font-weight-bold":      true,
	"--tracking-tight":        true,
	"--tracking-normal":       true,
	"--tracking-wide":         true,
	"--space-0":               true,
	"--space-1":               true,
	"--space-2":               true,
	"--space-3":               true,
	"--space-4":               true,
	"--space-6":               true,
	"--space-8":               true,
	"--space-12":              true,
	"--radius-sm":             true,
	"--radius-md":             true,
	"--radius-lg":             true,
	"--radius-full":           true,
	"--shadow-sm":             true,
	"--shadow-md":             true,
	"--shadow-lg":             true,
	"--shadow-xl":             true,
	"--duration-fast":         true,
	"--duration-base":         true,
	"--duration-slow":         true,
	"--ease-in":               true,
	"--ease-out":              true,
	"--ease-in-out":           true,
	"--z-base":                true,
	"--z-dropdown":            true,
	"--z-sticky":              true,
	"--z-modal":               true,
	"--z-toast":               true,
	"--z-tooltip":             true,
	"--bp-sm":                 true,
	"--bp-md":                 true,
	"--bp-lg":                 true,
	"--bp-xl":                 true,
	"--max-w-prose":           true,
	"--max-w-content":         true,
	"--max-w-screen":          true,
}

var (
	hexColorRegex   = regexp.MustCompile(`#[0-9a-fA-F]{3,6}\b`)
	rgbColorRegex   = regexp.MustCompile(`(rgb|rgba|hsl|hsla)\(`)
	viewportRegex   = regexp.MustCompile(`\d+(vw|vh|vmin|vmax)`)
	varRegex        = regexp.MustCompile(`var\((--[a-zA-Z0-9-]+)`)
	classRegex      = regexp.MustCompile(`\.([a-zA-Z0-9_-]+)`)
	htmlClassRegex  = regexp.MustCompile(`class='([^']*)'`)
	htmlClassRegex2 = regexp.MustCompile(`class="([^"]*)"`)
)

func TestConformance(t *testing.T) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "docs" || info.Name() == ".git" || info.Name() == "web" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			if n == nil {
				return true
			}

			// 1. Cero literales de color / Viewport / Vars / RawRule / Media en BasicLit (strings)
			if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				val := lit.Value
				// Strip quotes
				if len(val) >= 2 {
					val = val[1 : len(val)-1]
				}

				// Skip inline SVGs or Base64 strings (these contain '#' and are not color literals)
				if !strings.Contains(val, "data:image/") && !strings.Contains(val, "PHN2Z") {
					// Check color literals
					if hexColorRegex.MatchString(val) {
						t.Errorf("%s: contains hex color literal %q", path, val)
					}
					if rgbColorRegex.MatchString(val) {
						t.Errorf("%s: contains rgb/hsl color literal %q", path, val)
					}
					// Check viewport units
					if viewportRegex.MatchString(val) {
						t.Errorf("%s: contains viewport unit %q", path, val)
					}
					// Check literal color names (e.g., exact match for simple color names)
					for _, colorName := range []string{"white", "black", "red", "blue", "green", "yellow"} {
						if val == colorName {
							t.Errorf("%s: contains forbidden color name %q", path, val)
						}
					}
				}

				// Check var(--...) exists in allowed vars
				matches := varRegex.FindAllStringSubmatch(val, -1)
				for _, match := range matches {
					if len(match) > 1 {
						varName := match[1]
						if !allowedVars[varName] {
							t.Errorf("%s: references non-catalog variable %q", path, varName)
						}
					}
				}
			}

			// 2. Cero RawRule / Media / RootCSS / Theme / Declare identifiers
			if ident, ok := n.(*ast.Ident); ok {
				if ident.Name == "RawRule" {
					t.Errorf("%s: uses forbidden RawRule", path)
				}
				if ident.Name == "Media" {
					t.Errorf("%s: uses forbidden Media", path)
				}
				if ident.Name == "RootCSS" {
					t.Errorf("%s: uses forbidden RootCSS (only reserved for app / tinywasm/dom)", path)
				}
				if ident.Name == "Theme" {
					t.Errorf("%s: uses forbidden Theme (only reserved for app / tinywasm/dom)", path)
				}
				if ident.Name == "Declare" {
					t.Errorf("%s: uses forbidden Declare (only reserved for app / tinywasm/dom)", path)
				}
			}

			// 3. Cero constantes de clase escritas a mano (Class = "...")
			if valueSpec, ok := n.(*ast.ValueSpec); ok {
				isClassType := false
				if valueSpec.Type != nil {
					if ident, ok := valueSpec.Type.(*ast.Ident); ok && ident.Name == "Class" {
						isClassType = true
					} else if sel, ok := valueSpec.Type.(*ast.SelectorExpr); ok {
						if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "css" && sel.Sel.Name == "Class" {
							isClassType = true
						}
					}
				}
				if isClassType {
					for _, valExpr := range valueSpec.Values {
						if _, ok := valExpr.(*ast.BasicLit); ok {
							// Forbid hand-written Class assignment
							t.Errorf("%s: hand-written Class constant %q detected", path, valueSpec.Names[0].Name)
						}
					}
				}
			}

			return true
		})

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}

// TestPairMarkupAndStylesheet extracts classes from the compiled CSS stylesheet
// and ensures they map exactly to either the rendered HTML output of the components,
// or their exported/known classes.
func TestPairMarkupAndStylesheet(t *testing.T) {
	// Helper to extract classes from CSS string
	extractCSSClasses := func(css string) map[string]bool {
		classes := make(map[string]bool)
		matches := classRegex.FindAllStringSubmatch(css, -1)
		for _, m := range matches {
			if len(m) > 1 {
				classes[m[1]] = true
			}
		}
		return classes
	}

	// Helper to extract classes from HTML string
	extractHTMLClasses := func(html string) map[string]bool {
		classes := make(map[string]bool)
		matches1 := htmlClassRegex.FindAllStringSubmatch(html, -1)
		for _, m := range matches1 {
			if len(m) > 1 {
				for _, cls := range strings.Fields(m[1]) {
					classes[cls] = true
				}
			}
		}
		matches2 := htmlClassRegex2.FindAllStringSubmatch(html, -1)
		for _, m := range matches2 {
			if len(m) > 1 {
				for _, cls := range strings.Fields(m[1]) {
					classes[cls] = true
				}
			}
		}
		return classes
	}

	// Helper to filter classes by prefix (ignores global layout/utility classes)
	filterClasses := func(classes map[string]bool, prefix string) map[string]bool {
		filtered := make(map[string]bool)
		for cls := range classes {
			if strings.HasPrefix(cls, prefix) {
				filtered[cls] = true
			}
		}
		return filtered
	}

	// 1. ActionButton
	{
		btn := &actionbutton.ActionButton{Text: "Button", Variant: "primary"}
		html := btn.Render().String()
		css := btn.Style().Stylesheet().String()

		htmlClasses := filterClasses(extractHTMLClasses(html), "actionbutton")
		cssClasses := filterClasses(extractCSSClasses(css), "actionbutton")

		// Every class in CSS should correspond to a valid actionbutton class
		for cls := range cssClasses {
			if cls != "actionbutton" && cls != "actionbutton__primary" && cls != "actionbutton__secondary" && cls != "actionbutton__danger" {
				t.Errorf("ActionButton CSS contains unexpected class %q", cls)
			}
		}
		// Root class should be in HTML
		if !htmlClasses["actionbutton"] {
			t.Errorf("ActionButton HTML missing root class 'actionbutton'")
		}
	}

	// 2. ContentCard
	{
		c := &contentcard.ContentCard{
			Header: &contentcard.ContentCard{},
			Body:   &contentcard.ContentCard{},
			Footer: &contentcard.ContentCard{},
		}
		html := c.Render().String()
		css := c.Style().Stylesheet().String()

		htmlClasses := filterClasses(extractHTMLClasses(html), "contentcard")
		cssClasses := filterClasses(extractCSSClasses(css), "contentcard")

		for cls := range cssClasses {
			if !htmlClasses[cls] {
				t.Errorf("ContentCard: CSS class %q does not exist in rendered HTML", cls)
			}
		}
		for cls := range htmlClasses {
			if _, ok := cssClasses[cls]; !ok {
				t.Errorf("ContentCard: HTML class %q is unstyled in CSS", cls)
			}
		}
	}

	// 3. DataTable
	{
		dt := &datatable.DataTable{Headers: []string{"Col"}, Rows: [][]string{{"Val"}}}
		dt.Init(nil)
		html := dt.Render().String()
		css := dt.Style().Stylesheet().String()

		htmlClasses := filterClasses(extractHTMLClasses(html), "datatable")
		cssClasses := filterClasses(extractCSSClasses(css), "datatable")

		if !htmlClasses["datatable"] {
			t.Errorf("DataTable HTML missing root class 'datatable'")
		}
		if !cssClasses["datatable"] {
			t.Errorf("DataTable CSS missing root class 'datatable'")
		}
	}

	// 4. Fieldset (uses classes from tinywasm/form v0.3.0)
	{
		f := &fieldset.Fieldset{}
		css := f.Style().Stylesheet().String()
		cssClasses := filterClasses(extractCSSClasses(css), "tw-field")

		expectedFormClasses := map[string]bool{
			"tw-field":              true,
			"tw-field__label":       true,
			"tw-field__input":       true,
			"tw-field__error":       true,
			"tw-field__radio-group": true,
		}

		for cls := range cssClasses {
			if !expectedFormClasses[cls] {
				t.Errorf("Fieldset CSS contains unexpected class %q which form v0.3.0 does not emit", cls)
			}
		}
		for cls := range expectedFormClasses {
			if !cssClasses[cls] {
				t.Errorf("Fieldset CSS missing style rule for form v0.3.0 class %q", cls)
			}
		}
	}

	// 5. ModalDialog
	{
		md := &modaldialog.ModalDialog{Title: "Title", Content: &modaldialog.ModalDialog{}}
		md.Init(nil)
		md.Open()
		html := md.Render().String()
		css := md.Style().Stylesheet().String()

		htmlClasses := filterClasses(extractHTMLClasses(html), "modaldialog")
		cssClasses := filterClasses(extractCSSClasses(css), "modaldialog")

		for cls := range cssClasses {
			// Close button close does not have custom styling, but verify styling classes appear in HTML
			if !htmlClasses[cls] {
				t.Errorf("ModalDialog CSS class %q does not exist in rendered HTML", cls)
			}
		}
	}

	// 6. SelectSearch
	{
		ss := &selectsearch.SelectSearch{}
		css := ss.Style().Stylesheet().String()
		cssClasses := filterClasses(extractCSSClasses(css), "selectsearch")

		expectedSelectSearchClasses := map[string]bool{
			"selectsearch":           true,
			"selectsearch__toggle":   true,
			"selectsearch__dropdown": true,
			"selectsearch__header":   true,
			"selectsearch__icon":     true,
			"selectsearch__search":   true,
			"selectsearch__options":  true,
			"selectsearch__option":   true,
			"selectsearch__label":    true,
			"selectsearch__desc":     true,
		}

		for cls := range cssClasses {
			if !expectedSelectSearchClasses[cls] {
				t.Errorf("SelectSearch CSS contains unexpected class %q", cls)
			}
		}
		for cls := range expectedSelectSearchClasses {
			// selectsearch__label and selectsearch__desc might not have direct CSS rules if styled via parts,
			// but verify their core parts exist in stylesheet
			if cls == "selectsearch" || cls == "selectsearch__dropdown" || cls == "selectsearch__header" || cls == "selectsearch__search" || cls == "selectsearch__options" || cls == "selectsearch__option" {
				if !cssClasses[cls] {
					t.Errorf("SelectSearch CSS missing style rule for %q", cls)
				}
			}
		}
	}

	// 7. TargetList
	{
		tl := &targetlist.TargetList{}
		tl.Init(nil)
		css := tl.Style().Stylesheet().String()
		cssClasses := filterClasses(extractCSSClasses(css), "targetlist")

		expectedTargetListClasses := map[string]bool{
			"targetlist":           true,
			"targetlist__list":     true,
			"targetlist__backdrop": true,
			"targetlist__row":      true,
			"targetlist__label":    true,
			"targetlist__badge":    true,
			"targetlist__menu":     true,
			"targetlist__button":   true,
			"targetlist__icon":     true,
			"targetlist__options":  true,
			"targetlist__item":     true,
		}

		for cls := range cssClasses {
			if !expectedTargetListClasses[cls] {
				t.Errorf("TargetList CSS contains unexpected class %q", cls)
			}
		}
		for cls := range expectedTargetListClasses {
			if cls != "targetlist__label" && cls != "targetlist__button" && cls != "targetlist__icon" && cls != "targetlist__item" {
				if !cssClasses[cls] {
					t.Errorf("TargetList CSS missing style rule for expected class %q", cls)
				}
			}
		}
	}

	// 8. ThemeToggle
	{
		tt := &themetoggle.ThemeToggle{}
		tt.Init(nil)
		html := tt.Render().String()
		css := tt.Style().Stylesheet().String()

		htmlClasses := filterClasses(extractHTMLClasses(html), "themetoggle")
		cssClasses := filterClasses(extractCSSClasses(css), "themetoggle")

		if !htmlClasses["themetoggle"] {
			t.Errorf("ThemeToggle HTML missing class 'themetoggle'")
		}
		if !cssClasses["themetoggle"] {
			t.Errorf("ThemeToggle CSS missing class 'themetoggle'")
		}
	}
}
