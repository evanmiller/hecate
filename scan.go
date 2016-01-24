package main

import (
	"fmt"
	"math"
	"time"
	"strings"
)

func scanEpoch(value string, epoch time.Time) time.Time {
	parsed_time, parse_error := time.Parse("1/2/2006", value)
	if parse_error != nil {
		return epoch
	}
	return parsed_time
}

func scanOffset(value string, file_pos int) int {
	var err error
	expression := ""
	scanned_file_pos := 0
	if n, _ := fmt.Sscanf(value, "+%s", &expression); n > 0 {
		if scanned_file_pos, err = evaluateExpression(expression); err == nil {
			return file_pos + scanned_file_pos
		}
		return -1
	}
	if scanned_file_pos, err = evaluateExpression(value); err == nil {
		if scanned_file_pos < 0 {
			if scanned_file_pos+file_pos < 0 {
				return 0
			}
			return scanned_file_pos + file_pos
		}
		return scanned_file_pos
	}
	return -1
}

func scanSearchString(value string, bytes []byte, cursor Cursor, quit <-chan bool, progress chan<- int) *Cursor {
	representations := make(map[string]*Cursor)

	var scanned_fp float64
	var rest_of_value string
	if n, _ := fmt.Sscanf(value, "%g%s", &scanned_fp, &rest_of_value); n > 0 && len(rest_of_value) == 0 {
		fp32_string := binaryStringForInteger32(math.Float32bits(float32(scanned_fp)), cursor.big_endian)
		fp32_cursor := Cursor{mode: FloatingPointMode, fp_length: 4, big_endian: cursor.big_endian}
		representations[fp32_string] = &fp32_cursor

		fp64_string := binaryStringForInteger64(math.Float64bits(scanned_fp), cursor.big_endian)
		fp64_cursor := Cursor{mode: FloatingPointMode, fp_length: 8, big_endian: cursor.big_endian}
		representations[fp64_string] = &fp64_cursor

		var scanned_int int64
		if n, _ := fmt.Sscanf(value, "%d%s", &scanned_int, &rest_of_value); n > 0 && scanned_fp == float64(scanned_int) && len(rest_of_value) == 0 {
			if scanned_int >= math.MinInt8 && scanned_int <= math.MaxUint8 {
				int8_string := binaryStringForInteger8(uint8(scanned_int))
				int8_cursor := Cursor{mode: IntegerMode, int_length: 1, unsigned: (scanned_int > math.MaxInt8)}
				representations[int8_string] = &int8_cursor
			}
			if scanned_int >= math.MinInt16 && scanned_int <= math.MaxUint16 {
				int16_string := binaryStringForInteger16(uint16(scanned_int), cursor.big_endian)
				int16_cursor := Cursor{mode: IntegerMode, int_length: 2, unsigned: (scanned_int > math.MaxInt16),
					big_endian: cursor.big_endian}
				representations[int16_string] = &int16_cursor
			}
			if scanned_int >= math.MinInt32 && scanned_int <= math.MaxUint32 {
				int32_string := binaryStringForInteger32(uint32(scanned_int), cursor.big_endian)
				int32_cursor := Cursor{mode: IntegerMode, int_length: 4, unsigned: (scanned_int > math.MaxInt32),
					big_endian: cursor.big_endian}
				representations[int32_string] = &int32_cursor
			}
			int64_string := binaryStringForInteger64(uint64(scanned_int), cursor.big_endian)
			int64_cursor := Cursor{mode: IntegerMode, int_length: 8, unsigned: (scanned_int > math.MaxInt64),
				big_endian: cursor.big_endian}
			representations[int64_string] = &int64_cursor
		}
	}
	text_cursor := Cursor{mode: StringMode}
	representations[value] = &text_cursor

	first_match := -1
	first_length := 1
	first_cursor := cursor

	for k, v := range representations {
		start_pos := cursor.pos + 1
		found_pos := -1
		if start_pos < len(bytes) {
			found_pos = interruptibleSearch(bytes[start_pos:], k, quit, progress)
			if found_pos >= 0 {
				found_pos += start_pos
			}
		}
		if found_pos == -1 {
			found_pos = interruptibleSearch(bytes[0:cursor.pos], k, quit, progress)
		}
		if found_pos == -2 {
			return nil
		}
		v.pos = found_pos
	}

	found_match := false

	for _, v := range representations {
		if v.pos != -1 {
			found_match = true
			ranked_pos := v.pos - cursor.pos
			if ranked_pos < 0 {
				ranked_pos += len(bytes)
			}
			if first_match == -1 || ranked_pos < first_match ||
				(ranked_pos == first_match && v.length() > first_length) {
				first_match = ranked_pos
				first_length = v.length()
				first_cursor.pos = v.pos
				if v.mode == FloatingPointMode {
					first_cursor.fp_length = v.fp_length
				} else if v.mode == IntegerMode {
					first_cursor.int_length = v.int_length
				}
				first_cursor.mode = v.mode
				first_cursor.unsigned = v.unsigned
			}
		}
	}
	if !found_match {
		return nil
	}
	return &first_cursor
}

func scanEditedContent (value string, cursor Cursor) string {
	if cursor.mode == IntegerMode {
		var scanned_int int64
		if n, _ := fmt.Sscanf(value, "%d", &scanned_int); n < 1 {
			return ""
		}
		if cursor.int_length == 1 {
			if cursor.unsigned {
				return binaryStringForInteger8(uint8(scanned_int))
			} else {
				return binaryStringForInteger8(uint8(scanned_int))
			}
		} else if cursor.int_length == 2 {
			if cursor.unsigned {
				return binaryStringForInteger16(uint16(scanned_int), cursor.big_endian)
			} else {
				return binaryStringForInteger16(uint16(scanned_int), cursor.big_endian)
			}
		} else if cursor.int_length == 4 {
			if cursor.unsigned {
				return binaryStringForInteger32(uint32(scanned_int), cursor.big_endian)
			} else {
				return binaryStringForInteger32(uint32(scanned_int), cursor.big_endian)
			}
		} else if cursor.int_length == 8 {
			if cursor.unsigned {
				return binaryStringForInteger64(uint64(scanned_int), cursor.big_endian)
			} else {
				return binaryStringForInteger64(uint64(scanned_int), cursor.big_endian)
			}
		}
	} else if cursor.mode == FloatingPointMode {
		var scanned_fp float64
		if n, _ := fmt.Sscanf(value, "%g", &scanned_fp); n < 1 {
			return ""
		}
		if cursor.fp_length == 4 {
			return binaryStringForInteger32(math.Float32bits(float32(scanned_fp)), cursor.big_endian)
		} else if cursor.fp_length == 8 {
			return binaryStringForInteger64(math.Float64bits(scanned_fp), cursor.big_endian)
		}
	} else if cursor.mode == BitPatternMode {
		var scanned_int int64
		scan_fmt := "%" + string('0' + (cursor.bit_length * 2)) + "x"
		if n, _ := fmt.Sscanf(strings.Replace(value, " ", "", -1), scan_fmt, &scanned_int); n < 1 {
			return ""
		}
		if cursor.bit_length == 1 {
			return binaryStringForInteger8(uint8(scanned_int))
		} else if cursor.bit_length == 2 {
			return binaryStringForInteger16(uint16(scanned_int), true)
		}
	}

	return value
}
