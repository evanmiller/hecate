package main

type CursorWidget int

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
