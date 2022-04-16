package memgrf

import "sync"

type locker struct {
	mu sync.RWMutex
}

func (l *locker) rlock() func() {
	l.mu.RLock()
	return func() { l.mu.RUnlock() }
}

func (l *locker) lock() func() {
	l.mu.Lock()
	return func() { l.mu.Unlock() }
}
