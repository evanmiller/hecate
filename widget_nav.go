package main

type NavigationWidget int

func (widget NavigationWidget) sizeForLayout(layout Layout) Size {
	if layout.pressure > 1 {
		return Size{0, 0}
	}
	layouts := map[int]string{
		0: "Navigate: ←h ↓j ↑k →l",
		1: "Navigate: ←h ↓j",
	}
	runeCount := 0
	for _, _ = range layouts[layout.pressure] {
		runeCount++
	}
	return Size{runeCount, 2}
}

func (widget NavigationWidget) drawAtPoint(screen *DataScreen, layout Layout, point Point, style Style) Size {
	fg := style.default_fg
	bg := style.default_bg
	x_pos := point.x
	if layout.pressure == 0 {
		x_pos += drawStringAtPoint("Navigate: ←h ↓j ↑k →l", x_pos, point.y, fg, bg)
		x_pos = point.x + 10
		x_pos += drawStringAtPoint("←←←←b w→→→→", x_pos, point.y+1, fg, bg)
	} else if layout.pressure == 1 {
		x_pos += drawStringAtPoint("Navigate: ←h ↓j", x_pos, point.y, fg, bg)
		x_pos = point.x + 10
		x_pos += drawStringAtPoint("↑k →l", x_pos, point.y+1, fg, bg)
	}
	return Size{x_pos - point.x, 2}
}
