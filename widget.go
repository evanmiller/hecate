package main

import (
	"github.com/nsf/termbox-go"
)

type Widget interface {
	layoutUnderPressure(pressure int) Size
	drawAtPoint(cursor Cursor, point Point, pressure int, style Style, mode EditMode) Size
}

type WidgetSlice []Widget

func (widgets WidgetSlice) sizeAtPressure(pressure int) Size {
	total_widget_width := 0
	max_widget_height := 0
	for _, widget := range widgets {
		widget_size := widget.layoutUnderPressure(pressure)
		total_widget_width += widget_size.width
		if widget_size.height > max_widget_height {
			max_widget_height = widget_size.height
		}
	}
	return Size{total_widget_width, max_widget_height}
}

func (widgets WidgetSlice) numberVisibleAtPressure(pressure int) int {
	count := 0
	for _, widget := range widgets {
		widget_size := widget.layoutUnderPressure(pressure)
		if widget_size.width > 0 {
			count++
		}
	}
	return count
}

func (widgets WidgetSlice) layout() (int, int) {
	width, _ := termbox.Size()
	pressure := 0
	spacing := 4
	padding := 2
	for ; pressure < 10; pressure++ {
		spacing = 4
		total_widget_size := widgets.sizeAtPressure(pressure)
		num_spaces := widgets.numberVisibleAtPressure(pressure) - 1
		for ; total_widget_size.width+num_spaces*spacing > (width-2*padding) && spacing > 2; spacing-- {
		}
		if total_widget_size.width+num_spaces*spacing <= (width - 2*padding) {
			break
		}
	}
	return pressure, spacing
}

func listOfWidgets() WidgetSlice {
	all_widgets := [...]Widget{
		NavigationWidget(0),
		CursorWidget(0),
		OffsetWidget(0),
	}

	return all_widgets[:]
}

func heightOfWidgets() int {
	widgets := listOfWidgets()
	pressure, _ := widgets.layout()
	widget_size := widgets.sizeAtPressure(pressure)
	return widget_size.height
}

func drawWidgets(cursor Cursor, style Style, mode EditMode) Size {
	widgets := listOfWidgets()

	width, height := termbox.Size()
	spacing := 4
	padding := 2
	pressure, spacing := widgets.layout()
	total_widget_size := widgets.sizeAtPressure(pressure)
	num_spaces := widgets.numberVisibleAtPressure(pressure) - 1
	start_x := (width-2*padding-(total_widget_size.width+num_spaces*spacing))/2 + padding
	start_y := height - total_widget_size.height
	point := Point{start_x, start_y}
	for _, widget := range widgets {
		widget_size := widget.drawAtPoint(cursor, point, pressure, style, mode)
		point.x += widget_size.width
		if widget_size.width > 0 {
			point.x += spacing
		}
	}

	return Size{total_widget_size.width + num_spaces*spacing, total_widget_size.height}
}
