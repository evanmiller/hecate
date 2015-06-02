package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

type DataScreen struct {
	tabs       []*DataTab
	active_tab int
	show_tabs  bool
}

func (screen *DataScreen) initializeWithFiles(files []FileInfo) {
	cursor := Cursor{int_length: 4, fp_length: 4, mode: StringMode,
		epoch_unit: SecondsSinceEpoch, epoch_time: time.Unix(0, 0).UTC()}

	var tabs []*DataTab
	for _, file := range files {
		tab := DataTab{
			search_result_channel:   make(chan *Cursor),
			search_quit_channel:     make(chan bool),
			search_progress_channel: make(chan int),
			bytes:     file.bytes,
			filename:  file.filename,
			cursor:    cursor,
			hilite:    cursor.highlightRange(file.bytes),
			prev_mode: cursor.mode,
		}
		tabs = append(tabs, &tab)
	}

	screen.tabs = tabs
	screen.show_tabs = len(tabs) > 1
}

func (screen *DataScreen) receiveEvents(input <-chan termbox.Event, output chan<- int, quit <-chan bool) {
	var tab_quit_channels []chan bool
	for _ = range screen.tabs {
		quit_channel := make(chan bool, 10)
		tab_quit_channels = append(tab_quit_channels, quit_channel)
	}
	for i, t := range screen.tabs {
		go func(index int) {
			t.receiveEvents(output, tab_quit_channels[index])
		}(i)
	}

	for {
		do_quit := false
		select {
		case event := <-input:
			output <- screen.handleKeyEvent(event)
		case <-quit:
			do_quit = true
		}
		if do_quit {
			for _, c := range tab_quit_channels {
				c <- true
			}
			break
		}
	}
}

func (screen *DataScreen) handleKeyEvent(event termbox.Event) int {
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
	} else if event.Key == termbox.KeyCtrlW {
		if len(screen.tabs) > 1 {
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
		}
	} else if event.Key == termbox.KeyTab && screen.show_tabs {
		screen.active_tab = (screen.active_tab + 1) % len(screen.tabs)
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
		for x_pos < width {
			drawStringAtPoint("━", x_pos, 2, fg, bg)
			x_pos++
		}
		active_tab.drawTab(style, 3)
	} else {
		active_tab.drawTab(style, 0)
	}
}
