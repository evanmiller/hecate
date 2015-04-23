package main

import (
	"github.com/nsf/termbox-go"
)

type Style struct {
	default_bg termbox.Attribute
	default_fg termbox.Attribute

	rune_fg       termbox.Attribute
	space_rune_fg termbox.Attribute
	int_fg        termbox.Attribute
	bit_fg        termbox.Attribute

	text_cursor_hex_bg termbox.Attribute
	bit_cursor_hex_bg  termbox.Attribute
	int_cursor_hex_bg  termbox.Attribute
	fp_cursor_hex_bg   termbox.Attribute

	hilite_hex_fg  termbox.Attribute
	hilite_rune_fg termbox.Attribute

	field_editor_bg termbox.Attribute
	field_editor_fg termbox.Attribute
}

func defaultStyle() Style {
	var style Style
	style.default_bg = termbox.Attribute(1)
	style.default_fg = termbox.Attribute(256)
	style.rune_fg = termbox.Attribute(248)
	style.int_fg = termbox.Attribute(154)
	style.bit_fg = termbox.Attribute(154)
	style.space_rune_fg = termbox.Attribute(240)

	style.text_cursor_hex_bg = termbox.Attribute(167)
	style.bit_cursor_hex_bg = termbox.Attribute(26)
	style.int_cursor_hex_bg = termbox.Attribute(63)
	style.fp_cursor_hex_bg = termbox.Attribute(127)

	style.hilite_hex_fg = termbox.Attribute(231)
	style.hilite_rune_fg = termbox.Attribute(256)

	style.field_editor_bg = style.default_fg
	style.field_editor_fg = style.default_bg

	return style
}
