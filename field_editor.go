package main

import (
	"unicode"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

type FieldEditor struct {
	value      []byte
	cursor_pos int
}

func (field_editor *FieldEditor) handleKeyEvent(event termbox.Event) (string, bool) {
	is_done := false
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
	} else if unicode.IsPrint(event.Ch) {
		field_editor.value = insertRuneAtIndex(field_editor.value, field_editor.cursor_pos, event.Ch)
		field_editor.cursor_pos++
	} else if event.Key == termbox.KeySpace {
		field_editor.value = insertRuneAtIndex(field_editor.value, field_editor.cursor_pos, ' ')
		field_editor.cursor_pos++
	}
	return string(field_editor.value), is_done
}
