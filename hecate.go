package main

import (
	"fmt"
	"os"

	mmap "github.com/edsrzf/mmap-go"

	"github.com/nsf/termbox-go"
)

const PROGRAM_NAME = "hecate"

func mainLoop(bytes []byte, style Style) {
	screens := defaultScreensForData(bytes)
	display_screen := screens[DATA_SCREEN_INDEX]
	layoutAndDrawScreen(display_screen, style)
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			handleSpecialKeys(event.Key)

			new_screen_index := display_screen.handleKeyEvent(event)
			if new_screen_index < len(screens) {
				display_screen = screens[new_screen_index]
				layoutAndDrawScreen(display_screen, style)
			} else {
				break
			}
		}
		if event.Type == termbox.EventResize {
			layoutAndDrawScreen(display_screen, style)
		}
	}
}

func main() {
	var err error

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <filename>\n", PROGRAM_NAME)
		os.Exit(1)
	}
	path := os.Args[1]

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %q\n", err.Error())
		os.Exit(1)
	}

	fi, err := file.Stat()
	if err != nil {
		fmt.Printf("Error stat'ing file: %q\n", err.Error())
		os.Exit(1)
	}

	if fi.Size() < 8 {
		fmt.Printf("File %s is too short to be edited\n", path)
		os.Exit(1)
	}

	mm, err := mmap.Map(file, mmap.RDONLY, 0)
	if err != nil {
		fmt.Printf("Error mmap'ing file: %q\n", err.Error())
		os.Exit(1)
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	style := defaultStyle()
	termbox.SetOutputMode(outputMode)

	mainLoop(mm, style)
}
