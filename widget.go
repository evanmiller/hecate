package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type Widget interface {
	layoutUnderPressure(pressure int) (int, int)
	drawAtPoint(cursor Cursor, x int, y int, pressure int, style Style) (int, int)
}

type NavigationWidget int
type CursorWidget int
type OffsetWidget int

func (widget NavigationWidget) layoutUnderPressure(pressure int) (int, int) {
	if pressure > 1 {
		return 0, 0
	}
	layouts := map[int]string{
		0: "Navigate: ←h ↓j ↑k →l",
		1: "Navigate: ←h ↓j",
	}
	runeCount := 0
	for _, _ = range layouts[pressure] {
		runeCount++
	}
	return runeCount, 2
}

func (widget NavigationWidget) drawAtPoint(cursor Cursor, x int, y int, pressure int, style Style) (int, int) {
	fg := style.default_fg
	bg := style.default_bg
	x_pos := x
	if pressure == 0 {
		x_pos += drawStringAtPoint("Navigate: ←h ↓j ↑k →l", x_pos, y, fg, bg)
		x_pos = x + 10
		x_pos += drawStringAtPoint("←←←←b w→→→→", x_pos, y+1, fg, bg)
	} else if pressure == 1 {
		x_pos += drawStringAtPoint("Navigate: ←h ↓j", x_pos, y, fg, bg)
		x_pos = x + 10
		x_pos += drawStringAtPoint("↑k →l", x_pos, y+1, fg, bg)
	}
	return x_pos - x, 2
}

func (widget CursorWidget) layoutUnderPressure(pressure int) (int, int) {
	runeCount := 0
	height := 2
	if pressure < 5 || pressure == 6 {
		for _, _ = range "Cursor: " {
			runeCount++
		}
	}
	if pressure < 6 {
		for _, _ = range "(t)ext (p)attern (i)nteger (f)loat" {
			runeCount++
		}
	} else if pressure < 8 {
		for _, _ = range "(t)ext    (p)attern" {
			/*            (i)nteger (f)loat
			              (e)ndian  (u)nsigned */
			runeCount++
		}
		height = 3
	} else {
		for _, _ = range "(u)nsigned" {
			runeCount++
		}
		height = 6
	}
	return runeCount, height
}

func (widget CursorWidget) drawAtPoint(cursor Cursor, x int, y int, pressure int, style Style) (int, int) {
	fg := style.default_fg
	bg := style.default_bg
	x_pos := x
	y_pos := y
	if pressure < 5 || pressure == 6 {
		x_pos += drawStringAtPoint("Cursor: ", x_pos, y_pos, fg, bg)
	}
	if cursor.mode == StringMode {
		x_pos += drawStringAtPoint("(t)ext", x_pos, y_pos, fg, cursor.color(style))
	} else {
		x_pos += drawStringAtPoint("(t)ext", x_pos, y_pos, fg, bg)
	}
	if pressure < 6 {
		x_pos++
	} else if pressure < 8 {
		x_pos += 4
	} else {
		x_pos = x
		y_pos++
	}
	if cursor.mode == BitPatternMode {
		x_pos += drawStringAtPoint("(p)attern", x_pos, y_pos, fg, cursor.color(style))
	} else {
		x_pos += drawStringAtPoint("(p)attern", x_pos, y_pos, fg, bg)
	}
	if pressure < 6 {
		x_pos++
	} else if pressure < 7 {
		x_pos = x + 8
		y_pos++
	} else {
		x_pos = x
		y_pos++
	}
	if cursor.mode == IntegerMode {
		x_pos += drawStringAtPoint("(i)nteger", x_pos, y_pos, fg, cursor.color(style))
	} else {
		x_pos += drawStringAtPoint("(i)nteger", x_pos, y_pos, fg, bg)
	}
	if pressure < 8 {
		x_pos++
	} else {
		x_pos = x
		y_pos++
	}
	if cursor.mode == FloatingPointMode {
		x_pos += drawStringAtPoint("(f)loat", x_pos, y_pos, fg, cursor.color(style))
	} else {
		x_pos += drawStringAtPoint("(f)loat", x_pos, y_pos, fg, bg)
	}
	x_pos = x
	y_pos++
	if pressure < 5 || pressure == 6 {
		if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
			x_pos += drawStringAtPoint("Toggle: ", x_pos, y_pos, fg, bg)
		} else {
			x_pos += drawStringAtPoint("Toggle: ", x_pos, y_pos, style.space_rune_fg, bg)
		}
	}
	if pressure >= 8 {
		x_pos = x
	}
	if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
		if cursor.big_endian {
			x_pos += drawStringAtPoint("(E)ndian", x_pos, y_pos, fg, bg)
		} else {
			x_pos += drawStringAtPoint("(e)ndian", x_pos, y_pos, fg, bg)
		}
		x_pos++
	} else if cursor.mode == BitPatternMode || cursor.mode == StringMode {
		x_pos += drawStringAtPoint("(e)ndian", x_pos, y_pos, style.space_rune_fg, bg)
		x_pos++
	}
	if pressure >= 8 {
		x_pos = x
		y_pos++
	} else if pressure >= 6 {
		x_pos++
	}
	if cursor.mode == IntegerMode {
		if cursor.unsigned {
			x_pos += drawStringAtPoint("(U)nsigned", x_pos, y_pos, fg, bg)
		} else {
			x_pos += drawStringAtPoint("(u)nsigned", x_pos, y_pos, fg, bg)
		}
	} else {
		x_pos += drawStringAtPoint("(u)nsigned", x_pos, y_pos, style.space_rune_fg, bg)
	}
	if pressure < 6 {
		x_pos += 4
		if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
			x_pos += drawStringAtPoint("Size:", x_pos, y_pos, fg, bg)
			if cursor.length() > cursor.minimumLength() {
				x_pos += drawStringAtPoint(" ←H", x_pos, y_pos, fg, bg)
			} else {
				x_pos += drawStringAtPoint(" ←H", x_pos, y_pos, style.space_rune_fg, bg)
			}
			if cursor.length() < cursor.maximumLength() {
				x_pos += drawStringAtPoint(" →L", x_pos, y_pos, fg, bg)
			} else {
				x_pos += drawStringAtPoint(" →L", x_pos, y_pos, style.space_rune_fg, bg)
			}
		} else {
			x_pos += drawStringAtPoint("Size: ←H →L", x_pos, y_pos, style.space_rune_fg, bg)
		}
	}
	return x_pos - x, y_pos - y + 1
}

func (widget OffsetWidget) layoutUnderPressure(pressure int) (int, int) {
	if pressure > 3 {
		return 0, 0
	}
	return 16, 2
}

func (widget OffsetWidget) drawAtPoint(cursor Cursor, x int, y int, pressure int, style Style) (int, int) {
	if pressure > 3 {
		return 0, 0
	}
	fg := style.default_fg
	bg := style.default_bg
	drawStringAtPoint(fmt.Sprintf("Offset:  %d", cursor.pos), x, y, fg, bg)
	return drawStringAtPoint(fmt.Sprintf("  Type: %s", cursor.c_type()), x, y+1, fg, bg), 2
}

func sizeOfWidgets(widgets []Widget, pressure int) (int, int) {
	total_widget_width := 0
	max_widget_height := 0
	for _, widget := range widgets {
		widget_width, widget_height := widget.layoutUnderPressure(pressure)
		total_widget_width += widget_width
		if widget_height > max_widget_height {
			max_widget_height = widget_height
		}
	}
	return total_widget_width, max_widget_height
}

func numberOfVisibleWidgets(widgets []Widget, pressure int) int {
	count := 0
	for _, widget := range widgets {
		widget_width, _ := widget.layoutUnderPressure(pressure)
		if widget_width > 0 {
			count++
		}
	}
	return count
}

func drawWidgets(cursor Cursor, style Style) int {
	all_widgets := [...]Widget{
		NavigationWidget(0),
		CursorWidget(0),
		OffsetWidget(0),
	}
	widgets := all_widgets[:]

	width, height := termbox.Size()
	max_widget_height := 2
	spacing := 4
	padding := 2
	total_widget_width := 0
	pressure := 0

	for ; pressure < 10; pressure++ {
		spacing = 4
		total_widget_width, _ = sizeOfWidgets(widgets, pressure)
		num_spaces := numberOfVisibleWidgets(widgets, pressure) - 1
		for ; total_widget_width+num_spaces*spacing > (width-2*padding) && spacing > 2; spacing-- {
		}
		if total_widget_width+num_spaces*spacing <= (width - 2*padding) {
			break
		}
	}
	total_widget_width, max_widget_height = sizeOfWidgets(widgets, pressure)
	num_spaces := numberOfVisibleWidgets(widgets, pressure) - 1
	start_x, start_y := (width-2*padding-(total_widget_width+num_spaces*spacing))/2+padding, height-max_widget_height
	x, y := start_x, start_y
	for _, widget := range widgets {
		widget_width, _ := widget.drawAtPoint(cursor, x, y, pressure, style)
		x += widget_width
		if widget_width > 0 {
			x += spacing
		}
	}

	return max_widget_height
}
