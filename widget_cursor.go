package main

import "fmt"

type CursorWidget int

func (widget CursorWidget) sizeForLayout(layout Layout) Size {
	runeCount := 0
	height := 2
	if layout.show_date {
		height = 4
	}
	if layout.pressure < 5 || layout.pressure == 6 {
		for _, _ = range "Cursor: " {
			runeCount++
		}
	}
	if layout.pressure < 6 {
		for _, _ = range "(t)ext (p)attern (i)nteger (f)loat" {
			runeCount++
		}
	} else if layout.pressure < 8 {
		for _, _ = range "(t)ext    (p)attern" {
			/*            (i)nteger (f)loat
			              (e)ndian  (u)nsigned */
			runeCount++
		}
		if height < 3 {
			height = 3
		}
	} else {
		for _, _ = range "(u)nsigned" {
			runeCount++
		}
		height = 6
	}
	return Size{runeCount, height}
}

func (widget CursorWidget) drawAtPoint(cursor Cursor, data []byte, point Point, layout Layout, style Style, mode EditMode) Size {
	fg := style.default_fg
	bg := style.default_bg
	x_pos := point.x
	y_pos := point.y
	max_x_pos := x_pos
	pressure := layout.pressure

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
		x_pos = point.x
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
		x_pos = point.x + 8
		y_pos++
	} else {
		x_pos = point.x
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
		x_pos = point.x
		y_pos++
	}
	if cursor.mode == FloatingPointMode {
		x_pos += drawStringAtPoint("(f)loat", x_pos, y_pos, fg, cursor.color(style))
	} else {
		x_pos += drawStringAtPoint("(f)loat", x_pos, y_pos, fg, bg)
	}
	x_pos = point.x
	y_pos++
	if pressure < 5 || pressure == 6 {
		if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
			x_pos += drawStringAtPoint("Toggle: ", x_pos, y_pos, fg, bg)
		} else {
			x_pos += drawStringAtPoint("Toggle: ", x_pos, y_pos, style.space_rune_fg, bg)
		}
	}
	if pressure >= 8 {
		x_pos = point.x
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
		x_pos = point.x
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
	max_x_pos = x_pos
	x_pos = point.x
	if layout.show_date {
		y_pos++
		if pressure < 6 {
			y_pos++
		}
		date_fg := style.space_rune_fg
		if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
			date_fg = fg
		}
		if pressure < 5 || pressure == 6 {
			x_pos += drawStringAtPoint("(D)ate: ", x_pos, y_pos, date_fg, bg)
		}
		x_pos += drawStringAtPoint(fmt.Sprintf("%10s", cursor.interpretBytesAsTime(data).Format("1/2/2006")), x_pos, y_pos, date_fg, bg)
		x_pos++
		x_pos += drawStringAtPoint(fmt.Sprintf("%8s", cursor.interpretBytesAsTime(data).Format("3:04 PM")), x_pos, y_pos, date_fg, bg)
		if pressure < 6 {
			x_pos += 2
			if date_fg == fg && cursor.epoch_unit == SecondsSinceEpoch {
				x_pos += drawStringAtPoint("(s)ecs", x_pos, y_pos, date_fg, style.selected_option_bg)
			} else {
				x_pos += drawStringAtPoint("(s)ecs", x_pos, y_pos, date_fg, bg)
			}
			x_pos++
			if date_fg == fg && cursor.epoch_unit == DaysSinceEpoch {
				x_pos += drawStringAtPoint("(d)ays", x_pos, y_pos, date_fg, style.selected_option_bg)
			} else {
				x_pos += drawStringAtPoint("(d)ays", x_pos, y_pos, date_fg, bg)
			}
		}
	}
	return Size{max_x_pos - point.x, y_pos - point.y + 1}
}
