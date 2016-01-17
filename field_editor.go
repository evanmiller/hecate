package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

type FieldEditor struct {
	value      []byte
	cursor_pos int
	last_value string
	width      int
	fixed      int
	overwrite  bool
}

func (field_editor *FieldEditor) handleKeyEvent(event termbox.Event) (string, bool) {
	is_done := false

	if event.Ch == 0 && utf8.RuneCount(field_editor.value) == 0 {
		field_editor.value = []byte(field_editor.last_value)
	}

	if event.Key == termbox.KeyEnter {
		is_done = true
	} else if event.Key == termbox.KeyEsc {
		is_done = true
		field_editor.value = nil
	} else if event.Key == termbox.KeyArrowLeft {
		if field_editor.cursor_pos > 0 {
			field_editor.cursor_pos--
		}
	} else if event.Key == termbox.KeyArrowUp || event.Key == termbox.KeyCtrlA {
		field_editor.cursor_pos = 0
	} else if event.Key == termbox.KeyArrowRight {
		if field_editor.cursor_pos < utf8.RuneCount(field_editor.value) {
			field_editor.cursor_pos++
		}
	} else if event.Key == termbox.KeyArrowDown || event.Key == termbox.KeyCtrlE {
		field_editor.cursor_pos = utf8.RuneCount(field_editor.value)
	} else if event.Key == termbox.KeyCtrlH || event.Key == termbox.KeyBackspace {
		if field_editor.cursor_pos > 0 {
			field_editor.value = removeRuneAtIndex(field_editor.value, field_editor.cursor_pos-1)
			field_editor.cursor_pos--
		}
	} else if event.Key == termbox.KeyCtrlK {
		field_editor.cursor_pos = 0
		field_editor.value = make([]byte, 0)
	} else if unicode.IsPrint(event.Ch) {
		field_editor.insert(event.Ch)
	} else if event.Key == termbox.KeySpace {
		field_editor.insert(' ')
	}
	return string(field_editor.value), is_done
}

func (field_editor *FieldEditor) insert(r rune) {
	if field_editor.overwrite && field_editor.cursor_pos < utf8.RuneCount(field_editor.value) {
		field_editor.value[field_editor.cursor_pos] = byte(r)
	} else {
		if field_editor.fixed > 0 && field_editor.cursor_pos == field_editor.fixed {
			return
		}
		field_editor.value = insertRuneAtIndex(field_editor.value, field_editor.cursor_pos, r)
	}
	field_editor.cursor_pos++
}

func (field_editor *FieldEditor) drawFieldValueAtPoint(style Style, x, y int) int {
	termbox.SetCursor(x+1+field_editor.cursor_pos, y)
	if utf8.RuneCount(field_editor.value) > 0 || len(field_editor.last_value) == 0 {
		return drawStringAtPoint(fmt.Sprintf(" %-*s ", field_editor.width, field_editor.value), x, y,
			style.field_editor_fg, style.field_editor_bg)
	} else {
		return drawStringAtPoint(fmt.Sprintf(" %-*s ", field_editor.width, field_editor.last_value), x, y,
			style.field_editor_last_fg, style.field_editor_last_bg)
	}
}
