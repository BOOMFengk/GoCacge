package lru

import "container/list"

type Cache struct {
	maxBytes int64                    //允许最大内存
	nbytes   int64                    //目前已使用的内存
	ll       *list.List               //双向链表
	cache    map[string]*list.Element //字典, 键是字符串，值是双向链表中对应节点的指针

	OnEvicted func(key string, value Value) //是某条记录被移除时的回调函数,可以为nil
}

// entry 双向链表结点的数据类型
type entry struct {
	key   string
	value Value
}

// Value 接口的任意类型
type Value interface {
	Len() int //用于返回值所占用的内存大小
}

// New 构造器，实例化 Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 从字典中找到对应的双向链表的结点，将该节点移动到队尾

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 移除最近最少访问的节点，因为前面Get是往前 所以这里假定了最前面是队首
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 插入 先检测有没有有就放头上，并更新值，没有放头上，超过范围删除尾部
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
