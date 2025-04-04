package lru1

import "time"

type Item struct {
	key   string
	value any
	time  time.Time
}

type LRUCache struct {
	capacity int
	items    []Item
}

func New(capacity int) *LRUCache {

	return &LRUCache{
		capacity: capacity,
	}
}

func (l *LRUCache) Get(key string) any {
	for i := range l.items {
		if l.items[i].key == key {
			l.items[i].time = time.Now()
			return l.items[i].value
		}
	}
	return -1
}

func (l *LRUCache) Set(key string, value any) {
	if len(l.items) == l.capacity {
		oldestKey := -1
		oldestTime := time.Now()
		for i := range l.items {
			if oldestTime.After(l.items[i].time) {
				oldestTime = l.items[i].time
				oldestKey = i
			}
		}

		l.items[oldestKey].value = value
		l.items[oldestKey].key = key
		l.items[oldestKey].time = time.Time{}

		return
	}

	l.items = append(l.items, Item{key: key, value: value, time: time.Time{}})
}
