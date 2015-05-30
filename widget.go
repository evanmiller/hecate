package main

import (
	"github.com/nsf/termbox-go"
)

type Widget interface {
	sizeForLayout(layout Layout) Size
	drawAtPoint(cursor Cursor, data []byte, point Point, layout Layout, style Style, mode EditMode) Size
}

type WidgetSlice []Widget

func (widgets WidgetSlice) sizeForLayout(layout Layout) Size {
	total_widget_width := 0
	max_widget_height := 0
	for _, widget := range widgets {
		widget_size := widget.sizeForLayout(layout)
		total_widget_width += widget_size.width
		if widget_size.height > max_widget_height {
			max_widget_height = widget_size.height
		}
	}
	return Size{total_widget_width, max_widget_height}
}

func (widgets WidgetSlice) numberVisibleForLayout(layout Layout) int {
	count := 0
	for _, widget := range widgets {
		widget_size := widget.sizeForLayout(layout)
		if widget_size.width > 0 {
			count++
		}
	}
	return count
}

func (widgets WidgetSlice) layout(show_date bool) Layout {
	width, _ := termbox.Size()
	layout := Layout{pressure: 0, spacing: 4, num_spaces: 0, widget_size: Size{0, 0}, show_date: show_date}
	padding := 2
	for ; layout.pressure < 10; layout.pressure++ {
		layout.spacing = 4
		layout.widget_size = widgets.sizeForLayout(layout)
		layout.num_spaces = widgets.numberVisibleForLayout(layout) - 1
		for ; layout.width() > (width-2*padding) && layout.spacing > 2; layout.spacing-- {
		}
		if layout.width() <= (width - 2*padding) {
			break
		}
	}
	return layout
}

func listOfWidgets() WidgetSlice {
	all_widgets := [...]Widget{
		NavigationWidget(0),
		CursorWidget(0),
		OffsetWidget(0),
	}

	return all_widgets[:]
}

func heightOfWidgets(show_date bool) int {
	widgets := listOfWidgets()
	layout := widgets.layout(show_date)
	return layout.widget_size.height
}

func drawWidgets(cursor Cursor, data []byte, style Style, mode EditMode, show_date bool) Layout {
	widgets := listOfWidgets()

	width, height := termbox.Size()
	padding := 2
	layout := widgets.layout(show_date)
	start_x := (width-2*padding-layout.width())/2 + padding
	start_y := height - layout.widget_size.height
	point := Point{start_x, start_y}
	for _, widget := range widgets {
		widget_size := widget.drawAtPoint(cursor, data, point, layout, style, mode)
		point.x += widget_size.width
		if widget_size.width > 0 {
			point.x += layout.spacing
		}
	}

	return layout
}
