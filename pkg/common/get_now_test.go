package common_test

import (
	"testing"
	"time"

	"github.com/ngicks/gommon/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestGetNow(t *testing.T) {
	getNow := common.GetNowImpl{}

	now := time.Now()
	gNow := getNow.GetNow()

	// Oooh Does this really mean something? It should at least make coverage.
	require.InDelta(t, now.UnixNano(), gNow.UnixNano(), float64(time.Millisecond))
}
