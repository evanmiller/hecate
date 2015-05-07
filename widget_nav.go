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

func (widget NavigationWidget) drawAtPoint(cursor Cursor, x int, y int, pressure int, style *Style) (int, int) {
	x_pos := x
	if pressure == 0 {
		x_pos += StringOut("Navigate: ←h ↓j ↑k →l", x_pos, y, style)
		x_pos = x + 10
		x_pos += StringOut("←←←←b w→→→→", x_pos, y+1, style)
	} else if pressure == 1 {
		x_pos += StringOut("Navigate: ←h ↓j", x_pos, y, style)
		x_pos = x + 10
		x_pos += StringOut("↑k →l", x_pos, y+1, style)
	}
	return x_pos - x, 2
}
