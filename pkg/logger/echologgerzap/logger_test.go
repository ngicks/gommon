package echologgerzap_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ngicks/gommon/pkg/logger/echologgerzap"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	logger := echologgerzap.Default()
	buf := bytes.NewBuffer([]byte{})

	t.Run("Output / SetOutput", func(t *testing.T) {
		logger.SetOutput(buf)
		sizeBefore := buf.Len()
		logger.Debug("foobarbaz")
		sizeAfter := buf.Len()

		require.Greater(t, sizeAfter, sizeBefore)

		require.Equal(t, buf, logger.Output())
	})

	t.Run("Prefix", func(t *testing.T) {
		prefix := "foobar"
		logger.SetPrefix(prefix)

		require.Equal(t, prefix, logger.Prefix())

		logger.Info("quux")

		lastLine := unmarshalLastLine(buf)
		require.Contains(t, lastLine, "prefix")
		require.Equal(t, prefix, lastLine["prefix"])
	})

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
