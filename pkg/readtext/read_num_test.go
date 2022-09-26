package readtext_test

import (
	"testing"

	"github.com/ngicks/gommon/pkg/readtext"
	"github.com/stretchr/testify/assert"
)

func TestEqualAsciiCaseInsensitive(t *testing.T) {
	func() {
		defer func() {
			assert.NotNil(t, recover())
		}()
		assert.Equal(t, false, readtext.EqualAsciiCaseInsensitive("aaaa", "aaa"))
		_ = readtext.EqualAsciiCaseInsensitiveStrict("aaaa", "aaa")
	}()

	for _, v := range [][2]string{{"aaa", "aaa"}, {"ABC", "aBc"}, {"あいうえ", "あいうえ"}} {
		assert.Equal(t, true, readtext.EqualAsciiCaseInsensitive(v[0], v[1]))
		assert.Equal(t, true, readtext.EqualAsciiCaseInsensitiveStrict(v[0], v[1]))
	}
	for _, v := range [][2]string{{"aaa", "aaaa"}, {"ABC", "aBca"}} {
		assert.Equal(t, false, readtext.EqualAsciiCaseInsensitive(v[0], v[1]))
		// false report
		assert.Equal(t, true, readtext.EqualAsciiCaseInsensitiveStrict(v[0], v[1]))
	}

	for _, v := range [][2]string{{"aab", "aaa"}, {"abc", "acv"}, {"aBc", "abv"}} {
		assert.Equal(t, false, readtext.EqualAsciiCaseInsensitive(v[0], v[1]))
		assert.Equal(t, false, readtext.EqualAsciiCaseInsensitiveStrict(v[0], v[1]))
	}
}

func TestReadRuneN(t *testing.T) {
	input := `fooばあbaz😂漢字`

	var out, remaining string
	out, remaining = readtext.ReadRuneN(input, 1)
	assert.Equal(t, "f", out)
	assert.Equal(t, "ooばあbaz😂漢字", remaining)

	out, remaining = readtext.ReadRuneN(input, 4)
	assert.Equal(t, "fooば", out)
	assert.Equal(t, "あbaz😂漢字", remaining)

	out, remaining = readtext.ReadRuneN(input, 9)
	assert.Equal(t, "fooばあbaz😂", out)
	assert.Equal(t, "漢字", remaining)

	out, remaining = readtext.ReadRuneN(input, 11)
	assert.Equal(t, "fooばあbaz😂漢字", out)
	assert.Equal(t, "", remaining)

	out, remaining = readtext.ReadRuneN(input, 100)
	assert.Equal(t, "fooばあbaz😂漢字", out)
	assert.Equal(t, "", remaining)
}

type readNumNTestCase struct {
	input             string
	n                 int
	shouldBePadded    bool
	expectedFound     bool
	expectedNum       int
	expectedRemaining string
}

func TestReadNum2(t *testing.T) {
	cases := []readNumNTestCase{
		{
			input:             "123nnn",
			shouldBePadded:    true,
			expectedFound:     true,
			expectedNum:       12,
			expectedRemaining: "3nnn",
		},
		{
			input:             "1nnn",
			shouldBePadded:    false,
			expectedFound:     true,
			expectedNum:       1,
			expectedRemaining: "nnn",
		},
		{
			input:             "1nnn",
			shouldBePadded:    true,
			expectedFound:     false,
			expectedNum:       0,
			expectedRemaining: "1nnn",
		},
	}

	for _, testCase := range cases {
		num, remaining, found := readtext.ReadNum2(testCase.input, testCase.shouldBePadded)
		assert.Equal(t, testCase.expectedFound, found, "%+v", testCase)
		assert.Equal(t, testCase.expectedNum, num, "%+v", testCase)
		assert.Equal(t, testCase.expectedRemaining, remaining, "%+v", testCase)
	}
}

func TestReadNumN(t *testing.T) {
	cases := []readNumNTestCase{
		{input: "123nnn", n: 2, shouldBePadded: true, expectedNum: 12, expectedRemaining: "3nnn"},
		{input: "123nnn", n: 3, shouldBePadded: true, expectedNum: 123, expectedRemaining: "nnn"},
		{input: "123nnn", n: 4, shouldBePadded: false, expectedNum: 123, expectedRemaining: "nnn"},
		{input: "12345678nnn", n: 5, shouldBePadded: true, expectedNum: 12345, expectedRemaining: "678nnn"},
	}

	for _, testCase := range cases {
		num, remaining, found := readtext.ReadNumN(testCase.input, testCase.shouldBePadded, testCase.n)
		assert.Equal(t, found, true)
		assert.Equal(t, testCase.expectedNum, num)
		assert.Equal(t, testCase.expectedRemaining, remaining)
	}
}

func TestReadNumSpN(t *testing.T) {
	cases := []readNumNTestCase{
		{input: "123nnn", n: 2, expectedNum: 12, expectedFound: true, expectedRemaining: "3nnn"},
		{input: "123nnn", n: 3, expectedNum: 123, expectedFound: true, expectedRemaining: "nnn"},
		{input: "123nnn", n: 4, expectedNum: 0, expectedFound: false, expectedRemaining: "123nnn"},
		{input: "12345678nnn", n: 5, expectedNum: 12345, expectedFound: true, expectedRemaining: "678nnn"},
		{input: " 1234nnn", n: 4, expectedNum: 123, expectedFound: true, expectedRemaining: "4nnn"},
		{input: "  234nnn", n: 4, expectedNum: 23, expectedFound: true, expectedRemaining: "4nnn"},
		{input: "   34nnn", n: 4, expectedNum: 3, expectedFound: true, expectedRemaining: "4nnn"},
		{input: "    4nnn", n: 4, expectedNum: 0, expectedFound: false, expectedRemaining: "    4nnn"},
	}

	for _, testCase := range cases {
		num, remaining, found := readtext.ReadNumSpN(testCase.input, testCase.n)
		assert.Equal(t, testCase.expectedFound, found, "%+v", testCase)
		assert.Equal(t, testCase.expectedNum, num)
		assert.Equal(t, testCase.expectedRemaining, remaining)
	}
}
