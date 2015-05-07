package main

import (
	"github.com/nsf/termbox-go"
)

func handleSpecialKeys(key termbox.Key) {}

const outputMode = termbox.OutputNormal

func availableColors() []int {
	result := make([]int, int(termbox.ColorWhite))
	for i := termbox.ColorBlack; i <= termbox.ColorWhite; i++ {
		result[i-termbox.ColorBlack] = int(i)
	}
	return result
}

func defaultStyle() *Style {
	style, err := StyleFromJson(`
 {
 	"BG": 1,
 	"FG": 8,
	"Data": {
		"Hex": {
			"Highlight": {
				"FG": 6
			}
		},

		"Rune": {
			"FG": 4,

			"Code": {
				"FG": 8,

				"Highlight": {
					"FG": 4
				}
			},

			"Highlight": {
				"FG": 8
			}
		},

		"Bit": {
			"FG": 7
		},

		"Int": {
			"FG": 7
		},

		"Cursor": {
			"Int": {
				"BG": 7
			},
			"Bit": {
				"BG": 7
			},
			"Text": {
				"BG": 2
			},
			"Float": {
				"BG": 2
			}
		},

		"Disabled": {
			"FG": 8
		},

		"Edit": {
			"BG": 8,
 			"FG": 1
		}
	},
 	"About": {
 		"Logo": {
 			"BG": 2
 		}
 	},
 	"Palette": {}
 }`)

	if err != nil {
		panic(err)
	}

	return style
}
