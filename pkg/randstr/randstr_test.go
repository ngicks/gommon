package randstr_test

import (
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ngicks/gommon/pkg/randstr"
	"github.com/ngicks/type-param-common/set"
	"github.com/ngicks/type-param-common/util"
	"github.com/stretchr/testify/require"
)

// making it shorter to call
func must[T any](v T, err error) T {
	return util.Must(v, err)
}

// Tests that generator surley generates random values...
// Generate random arbitrary number of values and ensure all values are unique.
func TestRandomStrGen(t *testing.T) {

	// byteLen must be large enough to avoid value conflictions.
	byteLen := uint(128)
	generator := randstr.New(
		randstr.RandBytelen(byteLen),
		randstr.EncoderFactory(hex.NewEncoder),
	)
	genStrSet := set.New[string]()

	for i := 0; i < 1000; i++ {
		s, err := generator.String()
		if err != nil {
			t.Fatal(err)
		}
		require.Len(t, s, 128*2)
		genStrSet.Add(s)
	}
	require.Equal(t, genStrSet.Len(), 1000)
}

// Tests if default generator is valid.
func TestDefault(t *testing.T) {
	require := require.New(t)

	defaultRandStr := randstr.New()

	str, err := defaultRandStr.String()
	require.NoError(err)
	require.Len(str, 16)
}

// Tests if set options are correctly respected.
func TestOption(t *testing.T) {
	require := require.New(t)

	// set rand reader is used.
	seed := 123
	gen1 := randstr.New(randstr.RandReader(rand.New(rand.NewSource(int64(seed)))))
	gen2 := randstr.New(randstr.RandReader(rand.New(rand.NewSource(int64(seed)))))

	for i := 0; i < 100; i++ {
		str1 := string(must(gen1.Bytes()))
		str2 := must(gen2.String())
		require.Equal(str1, str2)
	}

	byteLen := 16
	gen1 = randstr.New(
		randstr.RandReader(rand.New(rand.NewSource(int64(seed)))),
		randstr.RandBytelen(uint(byteLen)),
		randstr.Hex(),
	)
	gen2 = randstr.New(
		randstr.RandReader(rand.New(rand.NewSource(int64(seed)))),
		randstr.RandBytelen(uint(byteLen)),
		randstr.Base64(),
	)

	str1 := must(gen1.String())
	str2 := must(gen2.String())

	// set len is used. hex transforms 1 bytes to 2 char. set bytelen(16) * 2 = 32
	require.Len(str1, byteLen*2)

	require.Condition(
		func() bool {
			// set encoder is surely used.
			decoded1 := must(hex.DecodeString(str1))
			decoded2 := must(base64.StdEncoding.DecodeString(str2))
			return cmp.Equal(decoded1, decoded2)
		},
	)
}

// Tests if an internal lock prevents race conditions.
func TestRace(t *testing.T) {
	require := require.New(t)

	seed := time.Now().UnixMicro()
	gen := randstr.New(randstr.RandReader(rand.New(rand.NewSource(seed))))

	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			_, err := gen.Bytes()
			require.NoError(err)
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestLen(t *testing.T) {
	require := require.New(t)

	var gen *randstr.Generator
	gen = randstr.New(randstr.RandBytelen(128))

	require.Len(must(gen.Bytes()), 128)
	require.Len(must(gen.BytesLen(256)), 256)
	require.Len(must(gen.String()), 128)
	require.Len(must(gen.StringLen(256)), 256)

	gen = randstr.New(randstr.RandBytelen(128), randstr.Hex())

	require.Len(must(gen.Bytes()), 256)
	require.Len(must(gen.BytesLen(256)), 512)
	require.Len(must(gen.String()), 256)
	require.Len(must(gen.StringLen(256)), 512)

	gen = randstr.New()

	for i := 0; i < 1024; i++ {
		require.Len(must(gen.BytesLen(int64(i))), i)
	}
}
