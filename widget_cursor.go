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

func (widget CursorWidget) drawAtPoint(cursor Cursor, x int, y int, pressure int, style *Style) (int, int) {
	x_pos := x
	y_pos := y

	cursorstyle := cursor.style(style)
	disabled := style.Sub("Disabled")

	if pressure < 5 || pressure == 6 {
		x_pos += style.StringOut("Cursor: ", x_pos, y_pos)
	}
	if cursor.mode == StringMode {
		x_pos += cursorstyle.StringOut("(t)ext", x_pos, y_pos)
	} else {
		x_pos += style.StringOut("(t)ext", x_pos, y_pos)
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
		x_pos += cursorstyle.StringOut("(p)attern", x_pos, y_pos)
	} else {
		x_pos += style.StringOut("(p)attern", x_pos, y_pos)
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
		x_pos += cursorstyle.StringOut("(i)nteger", x_pos, y_pos)
	} else {
		x_pos += style.StringOut("(i)nteger", x_pos, y_pos)
	}
	if pressure < 8 {
		x_pos++
	} else {
		x_pos = x
		y_pos++
	}
	if cursor.mode == FloatingPointMode {
		x_pos += cursorstyle.StringOut("(f)loat", x_pos, y_pos)
	} else {
		x_pos += style.StringOut("(f)loat", x_pos, y_pos)
	}
	x_pos = x
	y_pos++
	if pressure < 5 || pressure == 6 {
		if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
			x_pos += style.StringOut("Toggle: ", x_pos, y_pos)
		} else {
			x_pos += disabled.StringOut("Toggle: ", x_pos, y_pos)
		}
	}
	if pressure >= 8 {
		x_pos = x
	}
	if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
		if cursor.big_endian {
			x_pos += style.StringOut("(E)ndian", x_pos, y_pos)
		} else {
			x_pos += style.StringOut("(e)ndian", x_pos, y_pos)
		}
		x_pos++
	} else if cursor.mode == BitPatternMode || cursor.mode == StringMode {
		x_pos += disabled.StringOut("(e)ndian", x_pos, y_pos)
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
			x_pos += style.StringOut("(U)nsigned", x_pos, y_pos)
		} else {
			x_pos += style.StringOut("(u)nsigned", x_pos, y_pos)
		}
	} else {
		x_pos += disabled.StringOut("(u)nsigned", x_pos, y_pos)
	}
	if pressure < 6 {
		x_pos += 4
		if cursor.mode == IntegerMode || cursor.mode == FloatingPointMode {
			x_pos += style.StringOut("Size:", x_pos, y_pos)
			if cursor.length() > cursor.minimumLength() {
				x_pos += style.StringOut(" ←H", x_pos, y_pos)
			} else {
				x_pos += disabled.StringOut(" ←H", x_pos, y_pos)
			}
			if cursor.length() < cursor.maximumLength() {
				x_pos += style.StringOut(" →L", x_pos, y_pos)
			} else {
				x_pos += disabled.StringOut(" →L", x_pos, y_pos)
			}
		} else {
			x_pos += disabled.StringOut("Size: ←H →L", x_pos, y_pos)
		}
	}
	return x_pos - x, y_pos - y + 1
}
