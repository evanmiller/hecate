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

type SwitchScreen interface {
	screenIndex() int
}

type ScreenIndex int
func (idx ScreenIndex) screenIndex() int {
	return int(idx)
}

type FileInfo struct {
	filename string
	bytes    []byte
	rw       bool
}

func (file_info *FileInfo) baseName () string {
	suffix := map[bool]string {
		true: " *",
		false: "",
	}
	return path.Base(file_info.filename) + suffix[file_info.rw]
}

func (file_info *FileInfo) reopen (read_write bool) error {
	new_file, err := openFile(file_info.filename, read_write)
	if err == nil {
		*file_info = *new_file
	}

	return err
}

func mainLoop(files []FileInfo, style Style) {
	screens := defaultScreensForFiles(files)
	active_idx := DATA_SCREEN_INDEX

	var screen_key_channels []chan termbox.Event
	var screen_quit_channels []chan bool
	command_channel := make(chan interface{})
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
			screen.receiveEvents(screen_key_channels[index], command_channel,
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
		case cmd := <-command_channel:
			switch cmd := cmd.(type) {
			case SwitchScreen:
				new_screen_index := cmd.screenIndex()
				if new_screen_index < len(screens) {
					active_idx = new_screen_index
					layoutAndDrawScreen(screens[active_idx], style)
				} else {
					do_quit = true
				}
			default:
				fmt.Printf("unknown command: %T\n", cmd)
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

func openFile (filename string, read_write bool) (*FileInfo, error) {
	file_mode := os.O_RDONLY
	mmap_mode := mmap.RDONLY
	if read_write {
		file_mode = os.O_RDWR
		mmap_mode = mmap.RDWR
	}

	file, err := os.OpenFile(filename, file_mode, 0)
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

	mm, err := mmap.Map(file, mmap_mode, 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error mmap'ing file: %q\n", err.Error()))
	}

	return &FileInfo{ filename: filename, bytes: mm, rw: read_write }, nil
}

func main() {
	var err error

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename> [...]\n", PROGRAM_NAME)
		os.Exit(1)
	}
	var files []FileInfo
	for i := 1; i < len(os.Args); i++ {
		file_info, err := openFile(os.Args[i], false)
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
