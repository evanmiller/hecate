package main

import (
	"fmt"
)

type OffsetWidget int

func (widget OffsetWidget) layoutUnderPressure(pressure int) (int, int) {
	if pressure > 3 {
		return 0, 0
	}
	return 18, 2
}

func (widget OffsetWidget) drawAtPoint(cursor Cursor, x int, y int, pressure int, style *Style) (int, int) {
	if pressure > 3 {
		return 0, 0
	}
	if cursor.hex_mode {
		style.StringOut(fmt.Sprintf("Offset(:)  0x%x", cursor.pos), x, y)
	} else {
		style.StringOut(fmt.Sprintf("Offset(:)  %d", cursor.pos), x, y)
	}
	return style.StringOut(fmt.Sprintf("  Type :  %s", cursor.c_type()), x, y+1), 2
}
