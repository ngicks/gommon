package timing

import "sync"

// CreateWaiterFn returns a waiter, where caller can be blocked until passed fn's return.
func CreateWaiterFn(fn ...func()) (waiter func()) {
	var wg sync.WaitGroup

	sw := make(chan struct{})
	for _, f := range fn {
		wg.Add(1)
		go func(fn func()) {
			<-sw
			defer wg.Done()
			fn()
		}(f)
	}

	for i := 0; i < len(fn); i++ {
		sw <- struct{}{}
	}

	return wg.Wait
}

func CreateRepeatedWaiterFn(fn func(), repeat int) (waiter func()) {
	repeated := make([]func(), repeat)
	for i := 0; i < repeat; i++ {
		repeated[i] = fn
	}

	return CreateWaiterFn(repeated...)
}

func createWaiterCh(waiterMaker func() func()) (waiter <-chan struct{}) {
	waiterFn := waiterMaker()

	waiterCh := make(chan struct{})

	go func() {
		waiterFn()
		close(waiterCh)
	}()

	return waiterCh
}

func CreateWaiterCh(fn ...func()) (waiter <-chan struct{}) {
	return createWaiterCh(func() func() { return CreateWaiterFn(fn...) })
}

func CreateRepeatedWaiterCh(fn func(), repeat int) (waiter <-chan struct{}) {
	return createWaiterCh(func() func() { return CreateRepeatedWaiterFn(fn, repeat) })
}
