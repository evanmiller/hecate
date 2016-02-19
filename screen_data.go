package main

import "github.com/nsf/termbox-go"

const NAME_PADDING = 2
const TAB_MARGIN = 3

const THICK_LINE = "━"
const RIGHT_JOINT = "┕"
const LEFT_JOINT = "┙"
const DOUBLE_JOINT = "┷"

/*
const THICK_LINE = "═"
const RIGHT_JOINT = "╘"
const LEFT_JOINT = "╛"
const DOUBLE_JOINT = "╧"
*/

type TabListViewPort struct {
	width  int
	offset int
}

type DataScreen struct {
	tabs          []*DataTab
	tab_view_port TabListViewPort
	active_tab    int
	show_tabs     bool
}

func (screen *DataScreen) initializeWithFiles(files []FileInfo) {
	var tabs []*DataTab
	for _, file := range files {
		tab := NewDataTab(file)
		tabs = append(tabs, &tab)
	}

	screen.tabs = tabs
	screen.show_tabs = true
}

func (screen *DataScreen) receiveEvents(input <-chan termbox.Event, output chan<- interface{}, quit <-chan bool) {
	for _, t := range screen.tabs {
		go func(tab *DataTab) {
			tab.receiveEvents(output)
		}(t)
	}

	for {
		do_quit := false
		select {
		case event := <-input:
			output <- ScreenIndex(screen.handleKeyEvent(event, output))
		case <-quit:
			do_quit = true
		}
		if do_quit {
			for _, t := range screen.tabs {
				t.quit_channel <- true
			}
			break
		}
	}
}

func (screen *DataScreen) handleKeyEvent(event termbox.Event, output chan<- interface{}) int {
	active_tab := screen.tabs[screen.active_tab]
	if active_tab.field_editor != nil {
		return active_tab.handleKeyEvent(event)
	} else if event.Key == termbox.KeyCtrlLsqBracket { // color palette
		return PALETTE_SCREEN_INDEX
	} else if event.Ch == '?' { // about
		return ABOUT_SCREEN_INDEX
	} else if event.Ch == 'S' {
		screen.show_tabs = true
		return DATA_SCREEN_INDEX
	} else if event.Ch == 'W' {
		screen.show_tabs = false
		return DATA_SCREEN_INDEX
	} else if event.Key == termbox.KeyCtrlT {
		var new_tabs []*DataTab
		for index, old_tab := range screen.tabs {
			new_tabs = append(new_tabs, old_tab)
			if old_tab == active_tab {
				tab_copy := NewDataTab(old_tab.file_info)
				tab_copy.cursor = old_tab.cursor
				tab_copy.view_port = old_tab.view_port
				tab_copy.cursor.pos = tab_copy.view_port.first_row * tab_copy.view_port.bytes_per_row
				tab_copy.cursor.mode = StringMode
				new_tabs = append(new_tabs, &tab_copy)
				screen.active_tab = index + 1
				go func() {
					(&tab_copy).receiveEvents(output)
				}()
			}
		}
		screen.tabs = new_tabs
		screen.show_tabs = true
		return DATA_SCREEN_INDEX
	} else if event.Key == termbox.KeyCtrlW {
		active_tab.quit_channel <- true

		if len(screen.tabs) == 1 {
			return EXIT_SCREEN_INDEX
		}

		var new_tabs []*DataTab
		for _, old_tab := range screen.tabs {
			if old_tab != active_tab {
				new_tabs = append(new_tabs, old_tab)
			}
		}
		screen.tabs = new_tabs
		if screen.active_tab >= len(new_tabs) {
			screen.active_tab = len(new_tabs) - 1
		}
		return DATA_SCREEN_INDEX
	} else if event.Ch == 'D' && screen.show_tabs {
		if screen.active_tab < len(screen.tabs)-1 {
			screen.active_tab++
		}
		return DATA_SCREEN_INDEX
	} else if event.Ch == 'A' && screen.show_tabs {
		if screen.active_tab > 0 {
			screen.active_tab--
		}
		return DATA_SCREEN_INDEX
	}
	return active_tab.handleKeyEvent(event)
}

func (screen *DataScreen) performLayout() {
	width, height := termbox.Size()

	for _, tab := range screen.tabs {
		if screen.show_tabs {
			tab.performLayout(width, height-3)
		} else {
			tab.performLayout(width, height)
		}
	}

	screen.tab_view_port.width = width
	tab_pos := TAB_MARGIN
	active_tab_start_pos := 0
	active_tab_name_len := 0
	for index, tab := range screen.tabs {
		if index == screen.active_tab {
			active_tab_start_pos = tab_pos
			for _ = range tab.file_info.baseName() {
				active_tab_name_len++
			}
		}
		tab_pos += 2 + 2*NAME_PADDING
		for _ = range tab.file_info.baseName() {
			tab_pos++
		}
	}
	active_tab_end_pos := active_tab_start_pos + 2 + 2*NAME_PADDING + active_tab_name_len
	if tab_pos+TAB_MARGIN < screen.tab_view_port.offset+width {
		if tab_pos+TAB_MARGIN > width {
			screen.tab_view_port.offset = tab_pos + TAB_MARGIN - width
		} else {
			screen.tab_view_port.offset = 0
		}
	}
	if screen.tab_view_port.offset > active_tab_start_pos-TAB_MARGIN {
		screen.tab_view_port.offset = active_tab_start_pos - TAB_MARGIN
	}
	if screen.tab_view_port.offset+width < active_tab_end_pos+TAB_MARGIN {
		screen.tab_view_port.offset = active_tab_end_pos - width + TAB_MARGIN
	}
}

func (screen *DataScreen) drawScreen(style Style) {
	width, _ := termbox.Size()
	active_tab := screen.tabs[screen.active_tab]
	if screen.show_tabs {
		fg := style.default_fg
		bg := style.default_bg
		x_pos := -screen.tab_view_port.offset
		for i := 0; i < TAB_MARGIN; i++ {
			drawStringAtPoint(THICK_LINE, x_pos, 2, fg, bg)
			x_pos++
		}
		for _, tab := range screen.tabs {
			name_fg := fg
			if tab != active_tab {
				name_fg = style.rune_fg
			}
			drawStringAtPoint("╭", x_pos, 0, fg, bg)
			drawStringAtPoint("│", x_pos, 1, fg, bg)
			if tab == active_tab {
				drawStringAtPoint(LEFT_JOINT, x_pos, 2, fg, bg)
			} else {
				drawStringAtPoint(DOUBLE_JOINT, x_pos, 2, fg, bg)
			}
			x_pos++

			nameLength := drawStringAtPoint(tab.file_info.baseName(), x_pos+NAME_PADDING, 1, name_fg, bg)
			for i := 0; i < 2*NAME_PADDING+nameLength; i++ {
				drawStringAtPoint("─", x_pos, 0, fg, bg)
				if tab != active_tab {
					drawStringAtPoint(THICK_LINE, x_pos, 2, fg, bg)
				}
				x_pos++
			}
			drawStringAtPoint("╮", x_pos, 0, fg, bg)
			drawStringAtPoint("│", x_pos, 1, fg, bg)
			if tab == active_tab {
				drawStringAtPoint(RIGHT_JOINT, x_pos, 2, fg, bg)
			} else {
				drawStringAtPoint(DOUBLE_JOINT, x_pos, 2, fg, bg)
			}
			x_pos++
		}
		if width-x_pos > 22 {
			drawStringAtPoint("(?) help", width-20, 1, fg, bg)
		}
		if width-x_pos > 12 {
			drawStringAtPoint("(q)uit", width-10, 1, fg, bg)
		}
		for x_pos < width {
			drawStringAtPoint(THICK_LINE, x_pos, 2, fg, bg)
			x_pos++
		}
		active_tab.drawTab(style, 3)
	} else {
		active_tab.drawTab(style, 0)
	}
}
