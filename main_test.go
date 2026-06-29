package main

import (
	"testing"
)

func BenchmarkMutexConcurrentGet(b *testing.B) {
	state := &MutexState{
		db: make(map[string]string),
	}

	dispatchCommandMutex([]string{"SET", "testKey", "testValue"}, state)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			dispatchCommandMutex([]string{"GET", "testKey"}, state)
		}
	})
}

func BenchmarkChannelConcurrentGet(b *testing.B) {
	state := &ChanState{
		db:      make(map[string]string),
		actions: make(chan func()),
	}
	go runBackendManager(state)
	dispatchCommandChan([]string{"SET", "testKey", "myValue"}, state)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			dispatchCommandChan([]string{"GET", "testKey"}, state)
		}
	})
}
