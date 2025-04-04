package lru3

type Item struct {
	key   string
	value any
	next  *Item
	prev  *Item
}

type Cache struct {
	capacity int
	dict     map[string]*Item
	items    *Item
	head     *Item
}

func New(capacity int) *Cache {

	return &Cache{
		capacity: capacity,
		dict:     make(map[string]*Item),
	}
}

func (l *Cache) Get(key string) any {
	item, ok := l.dict[key]
	if !ok {
		return -1
	}

	if item == l.head {
		return item.value
	}

	if item.prev != nil {
		item.prev.next = item.next
	}

	if item.next != nil {
		item.next.prev = item.prev
	}

	if item == l.items {
		l.items = item.next
	}

	l.head.next = item
	item.prev = l.head
	l.head = item
	l.head.next = nil

	return item.value
}

func (l *Cache) Set(key string, value any) {
	if item, ok := l.dict[key]; ok {
		// remove old item with the same key
		if item.prev != nil {
			item.prev.next = item.next
		}

		if item.next != nil {
			item.next.prev = item.prev
		}

		if item == l.items {
			l.items = l.items.next
		}

		if item == l.head {
			l.head = l.head.prev
		}
		delete(l.dict, key)
	}

	if len(l.dict) == l.capacity {
		leastItem := l.items

		delete(l.dict, leastItem.key)
		l.items.key = key
		l.items.value = value
		l.dict[key] = l.items
		return
	}

	newItem := &Item{key: key, value: value}
	newItem.next = l.items
	if l.items != nil {
		l.items.prev = newItem
	}
	l.items = newItem
	l.dict[key] = newItem

	if l.items.next == nil {
		l.head = l.items
	}
}
