package geeCache

import "container/list"

type Cache struct {
	maxBytes  int64                         // 允许使用的最大内存
	nbytes    int64                         // 当前已使用的内存
	ll        *list.List                    // 标准包实现的双向链表
	cache     map[string]*list.Element      // 缓存字典
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数 可以为 nil
}

// entry
// @Description: 键值对 双向链表节点的数据类型
type entry struct {
	key   string // 保存 key 是为了在淘汰 key 时方便从字典中删除对应的映射
	value Value
}

// Value
// @Description: 值
type Value interface {
	Len() int
}

// New
// @Description: 实例化 Cache
// @param maxBytes
// @param onEvicted
// @return *Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get
// @Description: 查找
// @receiver c
// @param key
// @return value
// @return ok
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest
// @Description: 删除 淘汰策略
// @receiver c
func (c *Cache) RemoveOldest() {
	// 获取队首节点
	ele := c.ll.Back()
	if ele != nil {
		// 从淘汰队列中删除该节点
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 从字典中删除该节点的映射关系
		delete(c.cache, kv.key)
		// 维护占用内存
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 如果回调函数不为 nil 则调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add
// @Description: 新增或更新键
// @receiver c
// @param key
// @param value
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 更新
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 新增
		ele := c.ll.PushFront(&entry{key, value}) // 淘汰队列添加新的节点
		c.cache[key] = ele                        // 节点添加至映射
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// 若占用内存超出了内存限制 调用淘汰策略
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len
// @Description: 当前节点数
// @receiver c
// @return int
func (c *Cache) Len() int {
	return c.ll.Len()
}
