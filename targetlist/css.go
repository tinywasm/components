//go:build !wasm

package targetlist

import (
	. "github.com/tinywasm/css"
)

// Selection tokens referenced by name (not as css-lib constants) so the component
// compiles against any css version; the active theme supplies the values.
const (
	selectionBg = "var(--color-selection, #f5a623)"
	selectionFg = "var(--color-on-selection, #1c1c1e)"
)

// RenderCSS owns the targetlist look only: the scroll list, the row cards, the
// selected highlight, the description badge and the ⋮ options menu. Consumes
// theme tokens; the host layout positions the list, it does not style it.
func (t *TargetList) RenderCSS() *Stylesheet {
	return NewStylesheet(
		// Wrapper: passes flex sizing through from the host to the <ul> — the
		// backdrop (see below) needs to sit outside the <ul> since that element's
		// children are fully owned by the keyed reconcile (BindChildren), which
		// would fight a statically-added sibling.
		Rule(clsListWrap,
			Display(Flex_),
			FlexDirection(Column),
			FlexGrow(Str("1")),
			MinHeight(Str("0")),
			Height(Pct(100)),
		),

		// Full-page click-catcher, shown only while a row's ⋮ menu is open, so a
		// click anywhere outside a menu closes it (native <details> has no such
		// behavior on its own).
		Rule(clsMenuBackdrop,
			Display(None),
			Position(Str("fixed")),
			Top(Zero),
			Left(Zero),
			Right(Zero),
			Bottom(Zero),
			ZIndex(Str("4")), // below clsMenuList (5), above the rows
		),
		Rule(Selector("."+string(clsListWrap)+":has(."+string(clsMenu)+"[open]) ."+string(clsMenuBackdrop)),
			Display(Str("block")),
		),

		Rule(clsList,
			Display(Flex_),
			FlexDirection(Column),
			RawRule("gap: .5rem"),
			FlexGrow(Str("1")),
			MinHeight(Str("0")),
			Overflow(Auto),
			Padding(Str(".2rem")),
			ListStyle(None),
			Margin(Zero),
		),

		// Row card.
		Rule(clsRow,
			Position(Relative),
			Display(Flex_),
			FlexDirection(Column),
			JustifyContent(Center),
			MinHeight(Str("56px")),
			PaddingTop(Str(".5em")),
			PaddingRight(Str("2em")), // room for the ⋮ button
			PaddingBottom(Str(".5em")),
			PaddingLeft(Str(".7em")),
			Cursor(Pointer),
			FontWeight(Str("600")),
			Background(ColorBackground),
			Border(Str("1px solid "+ColorMuted.Var())),
			BorderRadius(Em(0.4)),
			Color(ColorOnSurface),
			Transition(Str("box-shadow .2s ease, border-color .2s ease")),
		),
		Rule(clsRow.Hover(),
			BorderColor(ColorPrimary),
			RawRule("box-shadow: 0 1px 4px rgba(0,0,0,.12)"),
		),
		// --color-selection / --color-on-selection are theme tokens (bright orange
		// on light, muted on dark; text flips dark→white). Referenced by name so
		// the component does not need a css-lib version that exports the constant.
		Rule(clsRowOn,
			Background(Str(selectionBg)),
			Color(Str(selectionFg)),
			BorderColor(Str(selectionBg)),
		),

		// Description badge, bottom-right.
		Rule(clsBadge,
			Position(Absolute),
			Right(Px(6)),
			Bottom(Px(6)),
			Background(ColorSurface),
			Color(ColorMuted),
			FontSize(Pct(75)),
			FontWeight(Str("500")),
			Padding(Str(".15em .5em")),
			BorderRadius(Em(0.4)),
		),

		// ⋮ menu (native <details>), top-right.
		Rule(clsMenu,
			Position(Absolute),
			Top(Px(4)),
			Right(Px(4)),
		),
		Rule(clsMenuBtn,
			ListStyle(None),
			Cursor(Pointer),
			Display(Flex_),
			AlignItems(Center),
			JustifyContent(Center),
			Width(Em(1.6)),
			Height(Em(1.6)),
			BorderRadius(Em(0.3)),
			Color(ColorMuted),
		),
		Rule(Selector("."+string(clsMenuBtn)+"::-webkit-details-marker"),
			Display(None),
		),
		Rule(clsMenuBtn.Hover(),
			Background(ColorSurface),
			Color(ColorOnSurface),
		),
		Rule(clsMenuIcon,
			Width(Px(14)),
			Height(Px(14)),
			Decl{Prop: "fill", Val: "currentColor"},
		),

		// Dropdown list shown while the <details> is open.
		Rule(clsMenuList,
			Position(Absolute),
			Right(Zero),
			RawRule("top: calc(100% + 2px)"),
			MinWidth(Em(7)),
			Display(Flex_),
			FlexDirection(Column),
			Background(ColorBackground),
			Border(Str("1px solid "+ColorMuted.Var())),
			BorderRadius(Em(0.4)),
			Overflow(Hidden),
			ZIndex(Str("5")),
			RawRule("box-shadow: 0 4px 12px rgba(0,0,0,.18)"),
		),
		// The list container clips overflow (it must, to scroll) — so the last
		// row's dropdown, which has no room below it, gets cut off opening
		// downward. Flip it to open upward instead, where the earlier rows
		// leave room.
		Rule(Selector("."+string(clsRow)+":last-child ."+string(clsMenuList)),
			RawRule("top: auto; bottom: calc(100% + 2px)"),
		),

		Rule(clsMenuItem,
			Display(Block),
			Width(Pct(100)),
			TextAlign(Str("left")),
			Padding(Str(".4em .7em")),
			Background(None),
			Border(None),
			Cursor(Pointer),
			FontSize(Rem(0.9)),
			Color(ColorOnSurface),
		),
		Rule(clsMenuItem.Hover(),
			Background(ColorSurface),
		),
	)
}
