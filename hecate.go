package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

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
			if event.Key == termbox.KeyCtrlZ {
				process, _ := os.FindProcess(os.Getpid())
				termbox.Close()
				process.Signal(syscall.SIGSTOP)
				termbox.Init()
			}
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

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %q\n", err.Error())
		os.Exit(1)
	}
	if len(bytes) < 8 {
		fmt.Printf("File %s is too short to be edited\n", path)
		os.Exit(1)
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	style := defaultStyle()
	termbox.SetOutputMode(termbox.Output256)

	mainLoop(bytes, style)
}
