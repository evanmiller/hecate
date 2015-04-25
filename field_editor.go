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
}

func (field_editor *FieldEditor) handleKeyEvent(event termbox.Event, file_pos int) int {
	new_file_pos := -1
	if event.Key == termbox.KeyEnter {
		if len(field_editor.value) > 0 {
			scanned_file_pos := 0
			if n, _ := fmt.Sscanf(string(field_editor.value), "+%v", &scanned_file_pos); n > 0 {
				new_file_pos = file_pos + scanned_file_pos
			} else if n, _ := fmt.Sscanf(string(field_editor.value), "%v", &scanned_file_pos); n > 0 {
				if scanned_file_pos < 0 {
					if scanned_file_pos+file_pos < 0 {
						new_file_pos = 0
					} else {
						new_file_pos = scanned_file_pos + file_pos
					}
				} else {
					new_file_pos = scanned_file_pos
				}
			}
		} else {
			new_file_pos = file_pos
		}
	} else if event.Key == termbox.KeyEsc {
		new_file_pos = file_pos
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
	}
	return new_file_pos
}
