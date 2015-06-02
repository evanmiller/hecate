package main

import "github.com/nsf/termbox-go"

type DataScreen struct {
	tabs       []*DataTab
	active_tab int
	show_tabs  bool
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

func (screen *DataScreen) receiveEvents(input <-chan termbox.Event, output chan<- int, quit <-chan bool) {
	for _, t := range screen.tabs {
		go func(tab *DataTab) {
			tab.receiveEvents(output)
		}(t)
	}

	for {
		do_quit := false
		select {
		case event := <-input:
			output <- screen.handleKeyEvent(event, output)
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

func (screen *DataScreen) handleKeyEvent(event termbox.Event, output chan<- int) int {
	active_tab := screen.tabs[screen.active_tab]
	if active_tab.field_editor != nil {
		return active_tab.handleKeyEvent(event)
	} else if event.Key == termbox.KeyCtrlP { // color palette
		return PALETTE_SCREEN_INDEX
	} else if event.Ch == '?' { // about
		return ABOUT_SCREEN_INDEX
	} else if event.Ch == 'T' {
		screen.show_tabs = !screen.show_tabs
		return DATA_SCREEN_INDEX
	} else if event.Key == termbox.KeyCtrlT {
		var new_tabs []*DataTab
		for index, old_tab := range screen.tabs {
			new_tabs = append(new_tabs, old_tab)
			if old_tab == active_tab {
				tab_copy := NewDataTab(FileInfo{filename: old_tab.filename, bytes: old_tab.bytes})
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
		if len(screen.tabs) > 1 {
			var new_tabs []*DataTab
			for _, old_tab := range screen.tabs {
				if old_tab != active_tab {
					new_tabs = append(new_tabs, old_tab)
				}
			}
			active_tab.quit_channel <- true
			screen.tabs = new_tabs
			if screen.active_tab >= len(new_tabs) {
				screen.active_tab = len(new_tabs) - 1
			}
			return DATA_SCREEN_INDEX
		}
	} else if event.Key == termbox.KeyTab && screen.show_tabs {
		screen.active_tab = (screen.active_tab + 1) % len(screen.tabs)
		return DATA_SCREEN_INDEX
	} else if event.Ch == '`' {
		screen.active_tab = (screen.active_tab + len(screen.tabs) - 1) % len(screen.tabs)
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
}

func (screen *DataScreen) drawScreen(style Style) {
	width, _ := termbox.Size()
	active_tab := screen.tabs[screen.active_tab]
	if screen.show_tabs {
		fg := style.default_fg
		bg := style.default_bg
		x_pos := 0
		for i := 0; i < 4; i++ {
			drawStringAtPoint("━", x_pos, 2, fg, bg)
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
				drawStringAtPoint("┙", x_pos, 2, fg, bg)
			} else {
				drawStringAtPoint("┷", x_pos, 2, fg, bg)
			}
			x_pos++

			nameLength := drawStringAtPoint(tab.filename, x_pos+2, 1, name_fg, bg)
			for i := 0; i < 2+nameLength+2; i++ {
				drawStringAtPoint("─", x_pos, 0, fg, bg)
				if tab != active_tab {
					drawStringAtPoint("━", x_pos, 2, fg, bg)
				}
				x_pos++
			}
			drawStringAtPoint("╮", x_pos, 0, fg, bg)
			drawStringAtPoint("│", x_pos, 1, fg, bg)
			if tab == active_tab {
				drawStringAtPoint("┕", x_pos, 2, fg, bg)
			} else {
				drawStringAtPoint("┷", x_pos, 2, fg, bg)
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
			drawStringAtPoint("━", x_pos, 2, fg, bg)
			x_pos++
		}
		active_tab.drawTab(style, 3)
	} else {
		active_tab.drawTab(style, 0)
	}
}
