package hw04lrucache

import "sync"

type Key string

// Будем везде записывать в значение и сам ключ, чтобы иметь возможность удалять из мапы за O(1)
// в противном случае, удаляя элемент из очереди, пришлось бы по значению искать его в мапе полным обходом мапы.
type lruValue struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(k Key, v interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	value := lruValue{
		key:   k,
		value: v,
	}

	if item, ok := c.items[k]; ok { // если элемент в мапе - обновить и переместить в начало очереди
		item.Value = value
		c.queue.MoveToFront(item)
		return ok
	}

	// если размер очереди после добавления нового элемента превысит емкость кэша,
	// то удалить последний элемент из очереди и его значение из словаря
	if c.queue.Len() >= c.capacity {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		delete(c.items, lastItem.Value.(lruValue).key)
	}

	// если элемента не было в мапе, то добавить в начало очереди и в мапу
	newItem := c.queue.PushFront(value)
	c.items[k] = newItem

	return false
}

func (c *lruCache) Get(k Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[k]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(lruValue).value, ok
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem)

	c.mu.Unlock()
}
