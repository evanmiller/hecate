package main

import (
	"github.com/nsf/termbox-go"
	"unicode/utf8"
)

func isASCII(val byte) bool {
	return (val >= 0x20 && val < 0x7f)
}

func isCode(val byte) bool {
	return val == 0x09 || val == 0x0A || val == 0x0D
}

func isPrintable(val byte) bool {
	return isASCII(val) || isCode(val)
}

func removeRuneAtIndex(value []byte, index int) []byte {
	var runes []rune
	new_string := make([]byte, utf8.UTFMax*(len(value)+1))
	pos := 0
	for _, runeValue := range string(value) {
		if pos != index {
			runes = append(runes, runeValue)
		}
		pos++
	}
	pos = 0
	for _, runeValue := range runes {
		pos += utf8.EncodeRune(new_string[pos:], runeValue)
	}
	return new_string[0:pos]
}

func insertRuneAtIndex(value []byte, index int, newRuneValue rune) []byte {
	var runes []rune
	new_string := make([]byte, utf8.UTFMax*(len(value)+1))
	pos := 0
	for _, runeValue := range string(value) {
		if pos == index {
			runes = append(runes, newRuneValue)
		}
		runes = append(runes, runeValue)
		pos++
	}
	if index == pos {
		runes = append(runes, newRuneValue)
	}
	pos = 0
	for _, runeValue := range runes {
		pos += utf8.EncodeRune(new_string[pos:], runeValue)
	}
	return new_string[0:pos]
}

func SetCell(x, y int, ch rune, s *Style) {
	termbox.SetCell(x, y, ch, s.Fg(), s.Bg())
}

func StringOut(str string, x, y int, s *Style) int {
	x_pos := x
	fg, bg := s.Fg(), s.Bg()
	for _, runeValue := range str {
		termbox.SetCell(x_pos, y, runeValue, fg, bg)
		x_pos++
	}
	return x_pos - x
}
