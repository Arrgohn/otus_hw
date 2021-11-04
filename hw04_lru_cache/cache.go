package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	mu *sync.Mutex
	items    map[Key]*ListItem
}

func (l lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.items[key]; ok {
		l.items[key].Value = cacheItem{key: string(key), value: value}

		l.queue.MoveToFront(l.items[key])

		return true
	}

	el := l.queue.PushFront(cacheItem{key: string(key), value: value})
	l.items[key] = el

	if len(l.items) > l.capacity{
		el := l.queue.Back()

		l.queue.Remove(el)

		key := el.Value.(cacheItem).key
		delete(l.items, Key(key))
	}

	return false
}

func (l lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.items[key]; ok {
		el := l.items[key]

		l.queue.MoveToFront(el)
		l.items[key] = l.queue.Front()

		return el.Value.(cacheItem).value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		mu: &sync.Mutex{},
		items:    make(map[Key]*ListItem, capacity),
	}
}
