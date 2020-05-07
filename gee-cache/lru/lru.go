package lru

import "container/list"

type Cache struct {
	maxBytes int64                    //maximum of the cache
	nBytes   int64                    //used Cache amount
	ll       *list.List               // linked List
	cacheMap map[string]*list.Element // cache Map of the Nodes
	//triggered after delete
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        &list.List{},
		cacheMap:  make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cacheMap[key]; ok {
		c.ll.MoveToFront(ele) //push to front  update linked list
		entry := ele.Value.(*entry)
		return entry.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() // pop from back
	if ele != nil {
		c.ll.Remove(ele)
		entry := ele.Value.(*entry)
		delete(c.cacheMap, entry.key)
		c.nBytes -= int64(len(entry.key)) + int64(entry.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(entry.key, entry.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cacheMap[key]; ok { // exists, update value on existing key
		c.ll.MoveToFront(ele)
		entry := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(entry.value.Len())
		ele.Value = value
	} else {
		entry := &entry{
			key:   key,
			value: value,
		}
		ele := c.ll.PushFront(entry) // return a Element
		c.cacheMap[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	if c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
