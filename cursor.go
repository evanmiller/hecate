package main

import (
	"fmt"
	"unsafe"

	"github.com/nsf/termbox-go"
)

type CursorMode int

const MAX_INTEGER_WIDTH = 8
const MIN_INTEGER_WIDTH = 1
const MAX_FLOATING_POINT_WIDTH = 8
const MIN_FLOATING_POINT_WIDTH = 4

const (
	StringMode CursorMode = iota + 1
	BitPatternMode
	IntegerMode
	FloatingPointMode
)

type ByteRange struct {
	pos    int
	length int
}

type Cursor struct {
	pos        int
	int_length int
	fp_length  int
	mode       CursorMode
	unsigned   bool
	big_endian bool
	hex_mode   bool
}

func (cursor Cursor) c_type() string {
	if cursor.mode == IntegerMode {
		if cursor.unsigned {
			return fmt.Sprintf("uint%d_t", cursor.int_length*8)
		} else {
			return fmt.Sprintf(" int%d_t", cursor.int_length*8)
		}
	} else if cursor.mode == FloatingPointMode {
		if cursor.fp_length == 4 {
			return " float"
		} else if cursor.fp_length == 8 {
			return " double"
		}
	} else if cursor.mode == BitPatternMode {
		return " char"
	} else if cursor.mode == StringMode {
		return " char *"
	}
	return ""
}

func (cursor Cursor) length() int {
	if cursor.mode == IntegerMode {
		return cursor.int_length
	}
	if cursor.mode == FloatingPointMode {
		return cursor.fp_length
	}
	return 1
}

func (cursor Cursor) maximumLength() int {
	if cursor.mode == IntegerMode {
		return MAX_INTEGER_WIDTH
	}
	if cursor.mode == FloatingPointMode {
		return MAX_FLOATING_POINT_WIDTH
	}
	return 1
}

func (cursor Cursor) minimumLength() int {
	if cursor.mode == IntegerMode {
		return MIN_INTEGER_WIDTH
	}
	if cursor.mode == FloatingPointMode {
		return MIN_FLOATING_POINT_WIDTH
	}
	return 1
}

func (cursor Cursor) color(style Style) termbox.Attribute {
	if cursor.mode == IntegerMode {
		return style.int_cursor_hex_bg
	}
	if cursor.mode == FloatingPointMode {
		return style.fp_cursor_hex_bg
	}
	if cursor.mode == BitPatternMode {
		return style.bit_cursor_hex_bg
	}
	return style.text_cursor_hex_bg
}

func (cursor Cursor) highlightRange(data []byte) ByteRange {
	var hilite ByteRange
	if cursor.mode != StringMode || !isPrintable(data[cursor.pos]) {
		return hilite
	}
	hilite.pos = cursor.pos
	for ; hilite.pos > 0 && isPrintable(data[hilite.pos-1]); hilite.pos-- {
	}
	for ; hilite.pos+hilite.length < len(data) && isPrintable(data[hilite.pos+hilite.length]); hilite.length++ {
	}
	return hilite
}

func (cursor Cursor) formatBytesAsNumber(data []byte) string {
	str := ""
	var integer uint64
	if cursor.big_endian {
		for i := 0; i < len(data); i++ {
			integer = (integer * 256) + uint64(data[i])
		}
	} else {
		for i := len(data) - 1; i >= 0; i-- {
			integer = (integer * 256) + uint64(data[i])
		}
	}
	if cursor.mode == IntegerMode {
		if cursor.int_length == 1 {
			if cursor.unsigned {
				str = fmt.Sprintf("%d", uint8(integer))
			} else {
				str = fmt.Sprintf("%d", int8(integer))
			}
		} else if cursor.int_length == 2 {
			if cursor.unsigned {
				str = fmt.Sprintf("%d", uint16(integer))
			} else {
				str = fmt.Sprintf("%d", int16(integer))
			}
		} else if cursor.int_length == 4 {
			if cursor.unsigned {
				str = fmt.Sprintf("%d", uint32(integer))
			} else {
				str = fmt.Sprintf("%d", int32(integer))
			}
		} else if cursor.int_length == 8 {
			if cursor.unsigned {
				str = fmt.Sprintf("%d", uint64(integer))
			} else {
				str = fmt.Sprintf("%d", int64(integer))
			}
		}
	} else if cursor.mode == FloatingPointMode {
		if cursor.fp_length == 4 {
			var integer32 uint32 = uint32(integer)
			str = fmt.Sprintf("%.5g", *(*float32)(unsafe.Pointer(&integer32)))
		} else if cursor.fp_length == 8 {
			str = fmt.Sprintf("%g", *(*float64)(unsafe.Pointer(&integer)))
		}
	}
	return str
}
