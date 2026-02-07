package components

import (
	"fmt"
	"reflect"
	"strings"
)

// CssVars defines the structure for the site's CSS variables.
type CssVars struct {
	Primary            string `css:"--color-primary"`
	Secondary          string `css:"--color-secondary"`
	Tertiary           string `css:"--color-tertiary"`
	Quaternary         string `css:"--color-quaternary"`
	Gray               string `css:"--color-gray"`
	Selection          string `css:"--color-selection"`
	Hover              string `css:"--color-hover"`
	Success            string `css:"--color-success"`
	Error              string `css:"--color-error"`
	MenuWidthCollapsed string `css:"--menu-width-collapsed"`
	MenuWidthExpanded  string `css:"--menu-width-expanded"`
}

// GetDefaultCssVars returns the default CSS variables.
func GetDefaultCssVars() CssVars {
	return CssVars{
		Primary:            "#ffffff",
		Secondary:          "#7c3aed",
		Tertiary:           "#94a3b8",
		Quaternary:         "#1e293b",
		Gray:               "#f8fafc",
		Selection:          "#a78bfa",
		Hover:              "#6d28d9",
		Success:            "#10b981",
		Error:              "#ef4444",
		MenuWidthCollapsed: "64px",
		MenuWidthExpanded:  "250px",
	}
}

// Render returns the CSS variables as a string (e.g. key: value;).
func (c CssVars) Render() string {
	var sb strings.Builder
	v := reflect.ValueOf(c)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("css")
		value := v.Field(i).String()

		if tag != "" && value != "" {
			sb.WriteString(fmt.Sprintf("%s: %s;\n", tag, value))
		}
	}
	return sb.String()
}
