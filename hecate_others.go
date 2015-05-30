// +build !windows

package main

import (
	"os"
	"syscall"

	"github.com/nsf/termbox-go"
)

func handleSpecialKeys(key termbox.Key) {
	if key == termbox.KeyCtrlZ {
		process, _ := os.FindProcess(os.Getpid())
		termbox.Close()
		process.Signal(syscall.SIGSTOP)
		termbox.Init()
	}
}

const outputMode = termbox.Output256

func defaultStyle() Style {
	var style Style
	style.default_bg = termbox.Attribute(1)
	style.default_fg = termbox.Attribute(256)
	style.rune_fg = termbox.Attribute(248)
	style.int_fg = termbox.Attribute(154)
	style.bit_fg = termbox.Attribute(154)
	style.space_rune_fg = termbox.Attribute(240)
	style.selected_option_bg = termbox.Attribute(240)

	style.text_cursor_hex_bg = termbox.Attribute(167)
	style.bit_cursor_hex_bg = termbox.Attribute(26)
	style.int_cursor_hex_bg = termbox.Attribute(63)
	style.fp_cursor_hex_bg = termbox.Attribute(127)

	style.hilite_hex_fg = termbox.Attribute(231)
	style.hilite_rune_fg = termbox.Attribute(256)

	style.about_logo_bg = termbox.Attribute(125)

	style.field_editor_bg = style.default_fg
	style.field_editor_fg = style.default_bg

	return style
}
