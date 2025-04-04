package lru2

import "time"

type Item struct {
	key   string
	value any
	time  time.Time
}

type Cache struct {
	capacity int
	items    map[string]*Item
}

func New(capacity int) *Cache {

	return &Cache{
		capacity: capacity,
		items:    make(map[string]*Item),
	}
}

func (l *Cache) Get(key string) any {
	item, ok := l.items[key]
	if !ok {
		return -1
	}

	l.items[key].time = time.Now()

	return item.value
}

func (l *Cache) Set(key string, value any) {
	if len(l.items) == l.capacity {
		oldestTime := time.Now()
		oldestKey := ""
		for i := range l.items {
			if l.items[i].time.Before(oldestTime) {
				oldestTime = l.items[i].time
				oldestKey = l.items[i].key
			}
		}

		delete(l.items, oldestKey)
	}

	l.items[key] = &Item{key: key, value: value, time: time.Time{}}
}
