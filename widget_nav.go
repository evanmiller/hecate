package main

type NavigationWidget int

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
