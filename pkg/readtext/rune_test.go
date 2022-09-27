package readtext_test

import (
	"testing"

	"github.com/ngicks/gommon/pkg/readtext"
	"github.com/stretchr/testify/assert"
)

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
