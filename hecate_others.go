// +build !windows

package main

import (
	"os"
	"syscall"

	"github.com/nsf/termbox-go"
)

func handleSpecialKeys(key termbox.Key) {
	if key == termbox.KeyCtrlZ {
		process, _ := os.FindProcess(os.Getpid())
		termbox.Close()
		process.Signal(syscall.SIGSTOP)
		termbox.Init()
	}
}

func availableColors() []int {
	result := make([]int, 256)
	for i := 1; i <= 256; i++ {
		result[i-1] = int(i)
	}
	return result
}

const outputMode = termbox.Output256

func defaultStyle() *Style {
	style, err := StyleFromJson(`
 {
 	"BG": 1,
 	"FG": 256,
	"Data": {
		"Hex": {
			"Highlight": {
				"FG": 231
			}
		},

		"Rune": {
			"FG": 248,

			"Code": {
				"FG": 240,

				"Highlight": {
					"FG": 248
				}
			},

			"Highlight": {
				"FG": 256
			}
		},

		"Bit": {
			"FG": 154
		},

		"Int": {
			"FG": 154
		},

		"Cursor": {
			"Int": {
				"BG": 63
			},
			"Bit": {
				"BG": 26
			},
			"Text": {
				"BG": 167
			},
			"Float": {
				"BG": 127
			}
		},

		"Disabled": {
			"FG": 240
		},

		"Edit": {
			"BG": 256,
 			"FG": 1
		}
	},
 	"About": {
 		"Logo": {
 			"BG": 125
 		}
 	},
 	"Palette": {}
 }`)

	if err != nil {
		panic(err)
	}

	return style
}
