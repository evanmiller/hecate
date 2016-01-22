package main

import (
	"fmt"
	"os"
	"path"
	"errors"

	mmap "github.com/edsrzf/mmap-go"

	"github.com/nsf/termbox-go"
)

const PROGRAM_NAME = "hecate"

type FileInfo struct {
	filename string
	bytes    []byte
}

func mainLoop(files []FileInfo, style Style) {
	screens := defaultScreensForFiles(files)
	active_idx := DATA_SCREEN_INDEX

	var screen_key_channels []chan termbox.Event
	var screen_quit_channels []chan bool
	switch_channel := make(chan int)
	main_key_channel := make(chan termbox.Event, 10)

	layoutAndDrawScreen(screens[active_idx], style)

	for _ = range screens {
		key_channel := make(chan termbox.Event, 10)
		screen_key_channels = append(screen_key_channels, key_channel)

		quit_channel := make(chan bool, 10)
		screen_quit_channels = append(screen_quit_channels, quit_channel)
	}

	for i, s := range screens {
		go func(index int, screen Screen) {
			screen.receiveEvents(screen_key_channels[index], switch_channel,
				screen_quit_channels[index])
		}(i, s)
	}

	go func() {
		for {
			event := termbox.PollEvent()
			if event.Type == termbox.EventInterrupt {
				break
			} else {
				main_key_channel <- event
			}
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
			termbox.Interrupt()
			break
		}
	}
}

func openFile (filename string) (*FileInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening file: %q\n", err.Error()))
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error stat'ing file: %q\n", err.Error()))
	}

	if fi.Size() < 8 {
		return nil, errors.New(fmt.Sprintf("File %s is too short to be edited\n", filename))
	}

	mm, err := mmap.Map(file, mmap.RDONLY, 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error mmap'ing file: %q\n", err.Error()))
	}

	return &FileInfo{filename: path.Base(filename), bytes: mm}, nil
}

func main() {
	var err error

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename> [...]\n", PROGRAM_NAME)
		os.Exit(1)
	}
	var files []FileInfo
	for i := 1; i < len(os.Args); i++ {
		file_info, err := openFile(os.Args[i])
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		files = append(files, *file_info)
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	style := defaultStyle()
	termbox.SetOutputMode(outputMode)

	mainLoop(files, style)
}
