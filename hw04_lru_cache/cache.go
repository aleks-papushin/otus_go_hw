package hw04lrucache

type Key string

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
	if item, ok := c.items[k]; ok { // если элемент в словаре - обновить и переместить в начало очереди
		item.Value = v
		c.queue.MoveToFront(item)
		return ok
	} else { // если нет, то добавить и поместить в начало,
		newItem := c.queue.PushFront(v)
		c.items[k] = newItem

		// при этом, если размер очереди становится больше емкости кэша,
		// то удалить последний элемент из очереди и его значение из словаря
		if c.queue.Len() > c.capacity {
			lastItem := c.queue.Back()
			c.queue.Remove(lastItem)
			for k, v := range c.items {
				if v == lastItem {
					delete(c.items, k)
				}
			}
		}
		return false
	}
}

func (c *lruCache) Get(k Key) (interface{}, bool) {
	if item, ok := c.items[k]; ok {
		c.queue.MoveToFront(item)
		return item.Value, ok
	} else {
		return nil, false
	}
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}
