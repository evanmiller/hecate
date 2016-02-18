// +build windows

package main

import (
	"github.com/nsf/termbox-go"
)

func handleSpecialKeys(key termbox.Key) {}

const outputMode = termbox.OutputNormal

func defaultStyle() Style {
	var style Style
	style.default_bg = termbox.ColorBlack
	style.default_fg = termbox.ColorWhite
	style.rune_fg = termbox.ColorYellow
	style.int_fg = termbox.ColorCyan
	style.bit_fg = termbox.ColorCyan
	style.space_rune_fg = termbox.ColorWhite
	style.selected_option_bg = termbox.ColorBlue
	style.search_progress_fg = termbox.ColorCyan

	style.text_cursor_hex_bg = termbox.ColorRed
	style.bit_cursor_hex_bg = termbox.ColorCyan
	style.int_cursor_hex_bg = termbox.ColorCyan
	style.fp_cursor_hex_bg = termbox.ColorRed

	style.hilite_hex_fg = termbox.ColorMagenta
	style.hilite_rune_fg = termbox.ColorMagenta

	style.about_logo_bg = termbox.ColorRed

	style.field_editor_bg = style.default_fg
	style.field_editor_fg = style.default_bg

	style.field_editor_last_bg = style.rune_fg
	style.field_editor_last_fg = style.default_fg

	style.field_editor_invalid_bg = termbox.ColorRed
	style.field_editor_invalid_fg = style.rune_fg

	style.space_rune = '•'
	style.filled_bit_rune = '●'
	style.empty_bit_rune = '○'

	return style
}
