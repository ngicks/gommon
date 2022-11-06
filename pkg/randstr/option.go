package randstr

import (
	"encoding/base64"
	"encoding/hex"
	"io"
)

type Option func(g *Generator) *Generator

// RandReader creates an Option that sets the source of random bytes.
func RandReader(r io.Reader) Option {
	return func(g *Generator) *Generator {
		g.randReader = r
		return g
	}
}

// RandByteLen creates an Option that sets length of the internal random bytes.
func RandBytelen(l uint) Option {
	return func(g *Generator) *Generator {
		g.byteLen = l
		return g
	}
}

// EncoderFactory creates an Option that sets encoder factory of a generator.
func EncoderFactory(factory func(w io.Writer) io.Writer) Option {
	return func(g *Generator) *Generator {
		g.encoderFactory = factory
		return g
	}
}

// Hex creates an Option that sets the encoder factory hex.NewEncoder.
func Hex() Option {
	return EncoderFactory(hex.NewEncoder)
}

// Base64 creates an Option that sets the encoder factory base64.NewEncoder using base64.StdEncoding.
func Base64() Option {
	return EncoderFactory(
		func(w io.Writer) io.Writer { return base64.NewEncoder(base64.StdEncoding, w) },
	)
}
