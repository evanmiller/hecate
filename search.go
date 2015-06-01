package main

import "bytes"

// primeRK is the prime base used in Rabin-Karp algorithm.
const primeRK = 16777619

// hashStr returns the hash and the appropriate multiplicative
// factor for use in Rabin-Karp algorithm.
func hashStr(sep string) (uint32, uint32) {
	hash := uint32(0)
	for i := 0; i < len(sep); i++ {
		hash = hash*primeRK + uint32(sep[i])
	}
	var pow, sq uint32 = 1, primeRK
	for i := len(sep); i > 0; i >>= 1 {
		if i&1 != 0 {
			pow *= sq
		}
		sq *= sq
	}
	return hash, pow
}

// Index returns the index of the first instance of sep in s, or -1 if sep is not present in s.
func interruptibleSearch(s []byte, sep string, quit chan bool, progress chan int) int {
	n := len(sep)
	switch {
	case n == 0:
		return 0
	case n == len(s):
		if bytes.Equal([]byte(sep), s) {
			return 0
		}
		return -1
	case n > len(s):
		return -1
	}
	hashsep, pow := hashStr(sep)
	var h uint32
	for i := 0; i < n; i++ {
		h = h*primeRK + uint32(s[i])
	}
	if h == hashsep && bytes.Equal(s[:n], []byte(sep)) {
		return 0
	}
	for i := n; i < len(s); {
		h *= primeRK
		h += uint32(s[i])
		h -= pow * uint32(s[i-n])
		i++
		if h == hashsep && bytes.Equal(s[i-n:i], []byte(sep)) {
			return i - n
		}
		select {
		case <-quit:
			return -2
		default:
		}
		if i%10000 == 0 {
			progress <- 10000
		}
	}
	progress <- len(s) % 10000
	return -1
}
