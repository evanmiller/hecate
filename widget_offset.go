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
		StringOut(fmt.Sprintf("Offset(:)  0x%x", cursor.pos), x, y, style)
	} else {
		StringOut(fmt.Sprintf("Offset(:)  %d", cursor.pos), x, y, style)
	}
	return StringOut(fmt.Sprintf("  Type :  %s", cursor.c_type()), x, y+1, style), 2
}
