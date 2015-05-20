package main

import (
	"fmt"
)

type OffsetWidget int

func (widget OffsetWidget) layoutUnderPressure(pressure int) Size {
	if pressure > 3 {
		return Size{0, 0}
	}
	return Size{18, 2}
}

func (widget OffsetWidget) drawAtPoint(cursor Cursor, point Point, pressure int, style Style, mode EditMode) Size {
	if pressure > 3 {
		return Size{0, 0}
	}
	fg := style.default_fg
	bg := style.default_bg
	if mode == EditingSearch {
		drawStringAtPoint("Search(/)", point.x, point.y, fg, bg)
	} else if cursor.hex_mode {
		drawStringAtPoint(fmt.Sprintf("Offset(:)  0x%x", cursor.pos), point.x, point.y, fg, bg)
	} else {
		drawStringAtPoint(fmt.Sprintf("Offset(:)  %d", cursor.pos), point.x, point.y, fg, bg)
	}
	return Size{drawStringAtPoint(fmt.Sprintf("  Type :  %s", cursor.c_type()), point.x, point.y+1, fg, bg), 2}
}
