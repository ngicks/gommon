package logger

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var fooKey *int
var barKey *int

func init() {
	var foo, bar int
	fooKey = &foo
	barKey = &bar
}

func TestZapLogger(t *testing.T) {
	var ctx context.Context
	ctx = context.WithValue(context.Background(), fooKey, "foooooo")
	ctx = context.WithValue(ctx, barKey, "barrrrr")

	zz, buf := bufLogger()
	var logger Logger = NewZapLogger(zz.Sugar())

	var ll Logger
	{
		ll = logger.SetContextKey(
			ContextField{Label: "foo", ContextKey: fooKey},
			ContextField{Label: "bar", ContextKey: barKey},
		)
		ll = ll.WithContext(ctx)
		ll.Debug("aaaaaahhhhhhhhhh!")
		logLine := unmarshalLastLine(buf)
		assert.Equal(t, "foooooo", logLine["foo"])
		assert.Equal(t, "barrrrr", logLine["bar"])
	}
	{
		ll = logger.SetContextKey(
			ContextField{Label: "foofoo", ContextKey: fooKey},
			ContextField{Label: "barbar", ContextKey: barKey},
		)
		ll = ll.WithContext(ctx)
		ll.Debug("aaaaaahhhhhhhhhh!")
		logLine := unmarshalLastLine(buf)
		assert.Equal(t, nil, logLine["foo"])
		assert.Equal(t, nil, logLine["bar"])
		assert.Equal(t, "foooooo", logLine["foofoo"])
		assert.Equal(t, "barrrrr", logLine["barbar"])
	}
	{
		ll = logger.SetContextKey(
			ContextField{Label: "foofoo", ContextKey: fooKey},
		)
		ll = ll.WithContext(ctx)
		ll.Debug("aaaaaahhhhhhhhhh!")
		logLine := unmarshalLastLine(buf)
		assert.Equal(t, "foooooo", logLine["foofoo"])
		assert.Equal(t, nil, logLine["barbar"])
	}
}

func bufLogger() (*zap.Logger, *bytes.Buffer) {
	buf := bytes.NewBuffer([]byte{})
	conf := zap.NewProductionEncoderConfig()

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(conf),
		zapcore.AddSync(buf),
		zap.DebugLevel,
	)

	return zap.New(core), buf
}

func unmarshalLastLine(buf *bytes.Buffer) map[string]interface{} {
	bin := buf.Bytes()
	cloned := make([]byte, len(bin))

	copy(cloned, bin)

	scanner := bufio.NewScanner(bytes.NewBuffer(cloned))

	var lastLine []byte
	for scanner.Scan() {
		lastLine = scanner.Bytes()
	}

	if lastLine == nil {
		panic("empty buf")
	}

	unmarshaled := map[string]interface{}{}
	err := json.Unmarshal(lastLine, &unmarshaled)
	if err != nil {
		panic(err)
	}

	return unmarshaled
}
