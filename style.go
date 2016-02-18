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

	selected_option_bg termbox.Attribute
	search_progress_fg termbox.Attribute

	text_cursor_hex_bg termbox.Attribute
	bit_cursor_hex_bg  termbox.Attribute
	int_cursor_hex_bg  termbox.Attribute
	fp_cursor_hex_bg   termbox.Attribute

	hilite_hex_fg  termbox.Attribute
	hilite_rune_fg termbox.Attribute

	field_editor_bg termbox.Attribute
	field_editor_fg termbox.Attribute

	field_editor_last_bg termbox.Attribute
	field_editor_last_fg termbox.Attribute

	field_editor_invalid_bg termbox.Attribute
	field_editor_invalid_fg termbox.Attribute

	about_logo_bg termbox.Attribute

	filled_bit_rune rune
	empty_bit_rune rune
	space_rune rune
}
