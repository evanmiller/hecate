package main

type NavigationWidget int

func (widget NavigationWidget) layoutUnderPressure(pressure int) Size {
	if pressure > 1 {
		return Size{0, 0}
	}
	layouts := map[int]string{
		0: "Navigate: ←h ↓j ↑k →l",
		1: "Navigate: ←h ↓j",
	}
	runeCount := 0
	for _, _ = range layouts[pressure] {
		runeCount++
	}
	return Size{runeCount, 2}
}

func (widget NavigationWidget) drawAtPoint(cursor Cursor, point Point, pressure int, style Style, mode EditMode) Size {
	fg := style.default_fg
	bg := style.default_bg
	x_pos := point.x
	if pressure == 0 {
		x_pos += drawStringAtPoint("Navigate: ←h ↓j ↑k →l", x_pos, point.y, fg, bg)
		x_pos = point.x + 10
		x_pos += drawStringAtPoint("←←←←b w→→→→", x_pos, point.y+1, fg, bg)
	} else if pressure == 1 {
		x_pos += drawStringAtPoint("Navigate: ←h ↓j", x_pos, point.y, fg, bg)
		x_pos = point.x + 10
		x_pos += drawStringAtPoint("↑k →l", x_pos, point.y+1, fg, bg)
	}
	return Size{x_pos - point.x, 2}
}
