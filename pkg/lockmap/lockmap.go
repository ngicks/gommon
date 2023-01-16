package lockmap

import (
	"github.com/ngicks/locker"
	syncparam "github.com/ngicks/type-param-common/sync-param"
)

type LockMap[K comparable, V any] struct {
	locks *locker.Locker[K]
	m     syncparam.Map[K, V]
}

func New[K comparable, V any]() *LockMap[K, V] {
	return &LockMap[K, V]{
		locks: locker.New[K](),
	}
}

func (m *LockMap[K, V]) Get(k K) (v V, ok bool) {
	m.locks.Lock(k)
	defer m.locks.Unlock(k)
	return m.m.Load(k)
}

func (m *LockMap[K, V]) Set(k K, v V) {
	m.locks.Lock(k)
	defer m.locks.Unlock(k)
	m.m.Store(k, v)
}

func (m *LockMap[K, V]) Delete(k K) {
	m.locks.Lock(k)
	defer m.locks.Unlock(k)
	m.m.Delete(k)
}

func (m *LockMap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key K, value V) bool {
		m.locks.Lock(key)
		defer m.locks.Unlock(key)
		return f(key, value)
	})
}

func (m *LockMap[K, V]) RunWithinLock(k K, fn func(v V, has bool, set func(v V))) (ok bool) {
	m.locks.Lock(k)
	defer m.locks.Unlock(k)
	v, ok := m.m.Load(k)
	fn(v, ok, func(v V) { m.m.Store(k, v) })
	return true
}
