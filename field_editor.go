package main

import (
	"fmt"
	"unicode"

	"github.com/nsf/termbox-go"
)

type FieldEditor struct {
	value      []rune
	cursor_pos int
	last_value string
	init_value string
	width      int
	fixed      int
	overwrite  bool
	at_bol     bool
	at_eol     bool
}

func (field_editor *FieldEditor) handleKeyEvent(event termbox.Event) (string, bool) {
	is_done := false

	if len(field_editor.value) == 0 {
		if event.Ch == 0 {
			field_editor.setValue([]rune(field_editor.last_value))
		} else if len(field_editor.init_value) > 0 {
			field_editor.setValue([]rune(field_editor.init_value))
		}
	}

	if field_editor.fixed > 0 {
		if len(field_editor.value) > field_editor.fixed {
			field_editor.setValue(field_editor.value[:field_editor.fixed])
		}
	}

	if event.Key == termbox.KeyEnter {
		is_done = true
	} else if event.Key == termbox.KeyEsc {
		is_done = true
		field_editor.value = nil
	} else if event.Key == termbox.KeyArrowLeft {
		field_editor.moveCursor(-1)
	} else if event.Key == termbox.KeyArrowUp || event.Key == termbox.KeyCtrlA {
		field_editor.at_bol = event.Key == termbox.KeyArrowUp
		field_editor.setCursorPos(0)
	} else if event.Key == termbox.KeyArrowRight {
		field_editor.moveCursor(1)
	} else if event.Key == termbox.KeyArrowDown || event.Key == termbox.KeyCtrlE {
		field_editor.at_eol = event.Key == termbox.KeyArrowDown
		field_editor.setCursorPos(len(field_editor.value))
	} else if event.Key == termbox.KeyCtrlH || event.Key == termbox.KeyBackspace {
		field_editor.delete_back()
	} else if event.Key == termbox.KeyCtrlD || event.Key == termbox.KeyDelete {
		field_editor.delete_front()
	} else if event.Key == termbox.KeyCtrlK {
		field_editor.setValue(make([]rune, 0))
	} else if unicode.IsPrint(event.Ch) {
		field_editor.insert(event.Ch)
	} else if event.Key == termbox.KeySpace {
		field_editor.insert(' ')
	}
	return string(field_editor.value), is_done
}

func (field_editor *FieldEditor) setValue (value []rune) {
	field_editor.value = value

	field_editor.at_bol = false
	field_editor.at_eol = false

	if field_editor.cursor_pos > len(field_editor.value) {
		field_editor.setCursorPos(len(field_editor.value))
	}
}

func (field_editor *FieldEditor) setCursorPos (pos int) {
	bol := field_editor.at_bol
	eol := field_editor.at_eol

	if pos < 0 {
		pos = 0
		bol = true
		eol = false
	} else if field_editor.fixed > 0 && pos >= field_editor.fixed {
		pos = field_editor.fixed - 1
		bol = false
		eol = true
	}
	if pos > len(field_editor.value) {
		pos = len(field_editor.value)
		bol = false
		eol = true
	}

	field_editor.cursor_pos = pos
	field_editor.at_bol = bol
	field_editor.at_eol = eol
}

func (field_editor *FieldEditor) moveCursor (delta int) {
	field_editor.setCursorPos(field_editor.cursor_pos + delta)
}

func (field_editor *FieldEditor) insert(r rune) {
	if field_editor.overwrite && field_editor.cursor_pos < len(field_editor.value) {
		field_editor.value[field_editor.cursor_pos] = r
	} else {
		if field_editor.fixed > 0 && field_editor.cursor_pos == field_editor.fixed {
			return
		}
		pos := field_editor.cursor_pos
		field_editor.value = append(field_editor.value[:pos], append([]rune{r}, field_editor.value[pos:]...)...)
	}

	field_editor.moveCursor(1)
}

func (field_editor *FieldEditor) delete_front() {
	pos := field_editor.cursor_pos
	if pos < len(field_editor.value) {
		field_editor.value = append(field_editor.value[:pos], field_editor.value[pos+1:]...)
	}
}

func (field_editor *FieldEditor) delete_back() {
	pos := field_editor.cursor_pos
	if pos == 0 {
		return
	} else if pos < len(field_editor.value) {
		field_editor.value = append(field_editor.value[:pos-1], field_editor.value[pos:]...)
	} else {
		field_editor.value = field_editor.value[:pos-1]
	}

	field_editor.moveCursor(-1)
}

func (field_editor *FieldEditor) drawFieldValueAtPoint(style Style, x, y int) int {
	termbox.SetCursor(x+1+field_editor.cursor_pos, y)
	if len(field_editor.value) > 0 || len(field_editor.last_value) == 0 {
		return drawStringAtPoint(fmt.Sprintf(" %-*s ", field_editor.width, string(field_editor.value)), x, y,
			style.field_editor_fg, style.field_editor_bg)
	} else {
		return drawStringAtPoint(fmt.Sprintf(" %-*s ", field_editor.width, field_editor.last_value), x, y,
			style.field_editor_last_fg, style.field_editor_last_bg)
	}
}
