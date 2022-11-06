package randstr

import (
	"bytes"
	"crypto/rand"
	"io"
	"sync"
)

// Random bytes / string generator.
type Generator struct {
	randMu         sync.Mutex
	randReader     io.Reader
	byteLen        uint
	encoderFactory func(w io.Writer) io.Writer
}

// New creates new *Generator with options applied.
//
// Default options are:
//   - rand byte generator: rand.Rand.
//   - length of generated rand value: 16.
//   - encoder of value: none.
func New(opts ...Option) *Generator {
	g := &Generator{
		randReader:     rand.Reader,
		byteLen:        16,
		encoderFactory: func(w io.Writer) io.Writer { return w },
	}

	for _, opt := range opts {
		g = opt(g)
	}

	return g
}

// BytesLen returns randomly generated slice of bytes.
//
// It generates random bytes of byteLen long, then processes bytes through an encoder
// which can be specified with a corresponding Option passed to New.
// Its output length may vary depending on the applied encoder
// (e.g. For hex encoder, output bytes has length of 2 * byteLen).
func (g *Generator) BytesLen(byteLen int64) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, byteLen))
	encoder := g.encoderFactory(buf)

	g.randMu.Lock()
	_, err := io.CopyN(encoder, g.randReader, byteLen)
	g.randMu.Unlock()

	if cl, ok := encoder.(io.Closer); ok {
		cl.Close()
	}

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// StringLen returns randomly generated string.
//
// It only converts the type of ByteLen.
func (g *Generator) StringLen(byteLen int64) (string, error) {
	buf, err := g.BytesLen(byteLen)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// Bytes generates randomly generated slice of bytes.
//
// It generates random bytes, then processes bytes through an encoder.
// Both of the length of the internal bytes and the encoder can be specified by passing
// corresponding Options to New.
// Its output length may vary depending on the applied encoder
// (e.g. For hex encoder, output bytes has length of 2 * byteLen).
func (g *Generator) Bytes() ([]byte, error) {
	return g.BytesLen(int64(g.byteLen))
}

// String returns randomly craeted string.
//
// It only converts the type of Bytes.
func (g *Generator) String() (string, error) {
	return g.StringLen(int64(g.byteLen))
}
