// readtext is collection of trivial text-reader functions.
//
// Some of code is copied from Go standard packages. So let's keep it creditted.
//
// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Full note can be also found in ./GO_LICENSE.
//
// Parts modified or written by me are governed by license that can be found in LICENSE.
package readtext

// EqualAsciiCaseInsensitiveStrict reports s1 and s2 is same string in case-insensitive way.
// It assumes both only contain ascii code.
func EqualAsciiCaseInsensitive(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	return EqualAsciiCaseInsensitiveStrict(s1, s2)
}

// EqualAsciiCaseInsensitiveStrict reports s1 and s2 is same string in case-insensitive way.
// It assumes both only contain ascii code. Both must be same len.
func EqualAsciiCaseInsensitiveStrict(s1, s2 string) bool {
	for i := 0; i < len(s1); i++ {
		c1 := s1[i]
		c2 := s2[i]
		if c1 != c2 {
			// Switch to lower-case; 'a'-'A' is known to be a single bit. (0b100000)
			c1 |= 'a' - 'A'
			c2 |= 'a' - 'A'
			if c1 != c2 || c1 < 'a' || c1 > 'z' {
				return false
			}
		}
	}
	return true
}

// IsDigit reports i-th is in range and i-th byte of val is digit of ascii code.
func IsDigit(val string, i int) bool {
	if len(val) <= i {
		return false
	}
	c := val[i]
	return '0' <= c && c <= '9'
}

// ReadNum2 cuts up to 2 digits from input val, returns parsed num as int and remaining string.
// It reports successful operation by returning true found.
// if shouldBePadded is true, 1st and 2nd byte of val must be digit (namely '0' <= chara <= '9').
func ReadNum2(val string, shouldBePadded bool) (num int, remaining string, found bool) {
	if !IsDigit(val, 0) {
		return 0, val, false
	}
	if !IsDigit(val, 1) {
		if shouldBePadded {
			return 0, val, false
		}
		return int(val[0] - '0'), val[1:], true
	}
	return int(val[0]-'0')*10 + int(val[1]-'0'), val[2:], true
}

// ReadNumN reads up to n digits from head of val, returns string parsed as number.
// remaining is rest of string.
// It reports successful operation by returning true found.
// if shouldBePadded is true, first n bytes must be digits (namely '0' <= chara <= '9').
func ReadNumN(val string, shouldBePadded bool, n int) (num int, remaining string, found bool) {
	var i int
	for ; i < n && IsDigit(val, i); i++ {
		num = num*10 + int(val[i]-'0')
	}
	if i == 0 || (shouldBePadded && i != n) {
		return 0, val, false
	}
	return num, val[i:], true
}

// ReadNumSpN reads up to n digits from head of val, returns string parsed as number.
// remaining is rest of string.
//
// A Space after first digits but before n digit is failure.
//
// It reports successful operation by returning true found.
// if shouldBePadded is true, first n bytes must be digits (namely '0' <= chara <= '9').
func ReadNumSpN(val string, n int) (num int, remaining string, found bool) {
	var i int
	for ; i < n && (val[i] == ' ' || IsDigit(val, i)); i++ {
		if val[i] != ' ' {
			found = true
			num = num*10 + int(val[i]-'0')
		}
	}
	if i == 0 || i != n || !found {
		return 0, val, false
	}
	return num, val[i:], true
}

// ReadMatched tests val from head against every element of tab and reports first matched index.
// It assumes val and elements of tab only contain ascii code.
//
// idx == -1 and remaining is val if no element matches.
func ReadMatchedCaseInsensitive(tab []string, val string) (idx int, remaining string) {
	for i, v := range tab {
		if len(val) >= len(v) && EqualAsciiCaseInsensitiveStrict(val[0:len(v)], v) {
			return i, val[len(v):]
		}
	}
	return -1, val
}
