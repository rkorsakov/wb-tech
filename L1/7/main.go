package main

import (
	"errors"
	"sync"
)

type ConcurMap[T comparable, V any] struct {
	mu        sync.Mutex
	concurMap map[T]V
}

func New[T comparable, V any]() *ConcurMap[T, V] {
	return &ConcurMap[T, V]{
		concurMap: make(map[T]V),
	}
}

func (m *ConcurMap[T, V]) Put(key T, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.concurMap[key] = value
}

func (m *ConcurMap[T, V]) Get(key T) (V, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if value, ok := m.concurMap[key]; ok {
		return value, nil
	}
	return *new(V), errors.New("key not found")
}

func main() {
	coolMap := New[int, string]()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			coolMap.Put(i, "hello")
		}(i)
	}
	wg.Wait()
}
