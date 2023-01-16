package lockmap_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ngicks/gommon/pkg/lockmap"
	"github.com/ngicks/gommon/pkg/timing"
	"github.com/stretchr/testify/assert"
)

func TestLockMap(t *testing.T) {
	assert := assert.New(t)

	someMap := lockmap.New[string, string]()

	someMap.Set("1", "1")

	done := make(chan struct{})
	waiter := timing.CreateWaiterCh(func() {
		someMap.RunWithinLock("1", func(v string, set func(v string)) {
			<-done
			set("123")
		})
	})

	someMap.Set("2", "2")
	v, ok := someMap.Get("2")
	assert.True(ok)
	assert.Equal("2", v)
	v, ok = someMap.Get("3")
	assert.False(ok)
	assert.Equal("", v)

	blockings := [](<-chan struct{}){
		timing.CreateWaiterCh(func() { someMap.Set("1", "15") }),
		timing.CreateWaiterCh(func() { someMap.Get("1") }),
		timing.CreateWaiterCh(func() { someMap.Delete("1") }),
		timing.CreateWaiterCh(func() {
			someMap.Range(func(key, value string) bool {
				return true
			})
		}),
	}

	for _, ch := range blockings {
		select {
		case <-ch:
			t.Errorf("key \"%s\" must be locked by RunWithinLock", "1")
		case <-time.After(time.Millisecond):
		}
	}

	close(done)
	<-waiter

	for _, ch := range blockings {
		<-ch
	}
}

func TestLockMap_race(t *testing.T) {
	someMap := lockmap.New[string, string]()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k, v := range map[string]string{
				"1": "1",
				"2": "2",
				"3": "3",
			} {
				someMap.Set(k, v)
				someMap.Get(k)
				someMap.Delete(k)
				someMap.Range(func(key, value string) bool {
					return true
				})
				someMap.RunWithinLock(k, func(v string, set func(v string)) {
					if v == "2" {
						set("15")
					}
				})
			}
		}()
	}

	wg.Wait()
}
