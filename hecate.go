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
	active_idx := DATA_SCREEN_INDEX

	var screen_key_channels []chan termbox.Event
	var screen_quit_channels []chan bool
	switch_channel := make(chan int)
	main_key_channel := make(chan termbox.Event, 10)

	layoutAndDrawScreen(screens[active_idx], style)

	for _ = range screens {
		key_channel := make(chan termbox.Event, 10)
		screen_key_channels = append(screen_key_channels, key_channel)

		quit_channel := make(chan bool)
		screen_quit_channels = append(screen_quit_channels, quit_channel)
	}

	for i, s := range screens {
		go func(index int) {
			s.receiveEvents(screen_key_channels[index], switch_channel,
				screen_quit_channels[index])
		}(i)
	}

	go func() {
		for {
			main_key_channel <- termbox.PollEvent()
		}
	}()

	for {
		do_quit := false
		select {
		case event := <-main_key_channel:
			if event.Type == termbox.EventKey {
				handleSpecialKeys(event.Key)

				screen_key_channels[active_idx] <- event
			}
			if event.Type == termbox.EventResize {
				layoutAndDrawScreen(screens[active_idx], style)
			}
		case new_screen_index := <-switch_channel:
			if new_screen_index < len(screens) {
				active_idx = new_screen_index
				layoutAndDrawScreen(screens[active_idx], style)
			} else {
				do_quit = true
			}
		}
		if do_quit {
			for _, c := range screen_quit_channels {
				c <- true
			}
			break
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
