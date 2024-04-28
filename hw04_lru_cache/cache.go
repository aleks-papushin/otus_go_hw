package hw04lrucache

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
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(k Key, v interface{}) bool {
	value := lruValue{
		key:   k,
		value: v,
	}

	if item, ok := c.items[k]; ok { // если элемент в словаре - обновить и переместить в начало очереди
		item.Value = value
		c.queue.MoveToFront(item)
		return ok
	}

	// если нет, то добавить и поместить в начало,
	newItem := c.queue.PushFront(value)
	c.items[k] = newItem

	// при этом, если размер очереди становится больше емкости кэша,
	// то удалить последний элемент из очереди и его значение из словаря
	if c.queue.Len() > c.capacity {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		delete(c.items, lastItem.Value.(lruValue).key)
	}
	return false
}

func (c *lruCache) Get(k Key) (interface{}, bool) {
	if item, ok := c.items[k]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(lruValue).value, ok
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}
