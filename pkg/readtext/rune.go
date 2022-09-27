package readtext

import "unicode/utf8"

// ReadRuneN cuts n utf-8 code points from val and returns as runes, returns rest as remaining.
func ReadRuneN(val string, n int) (runes string, remaining string) {
	var offset int
	for i := 0; i < n; i++ {
		_, size := utf8.DecodeRune([]byte(val[offset:]))
		if size == 0 {
			break
		}
		offset += size
	}
	return val[0:offset], val[offset:]
}
