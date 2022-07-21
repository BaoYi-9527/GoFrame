package geeCache

import (
	"fmt"
	"log"
	"sync"
)

// Getter
// @Description: 定义接口 Getter 和 回调函数 Get
type Getter interface {
	Get(key string) ([]byte, error)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

// GetterFunc 定义函数类型 并实现 Getter
type GetterFunc func(key string) ([]byte, error)

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// Get
// @Description:
// @receiver f
// @param key
// @return []byte
// @return error
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// NewGroup
// @Description: 缓存的命名空间
// @param name
// @param cacheBytes
// @param getter	缓存未命中时调用的回调函数
// @return *Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup
// @Description:
// @param name
// @return *Group
func GetGroup(name string) *Group {
	mu.RLock() // 只读锁
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	// 先从 mainCache zho中查找
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	// 缓存取不到 调用回调函数
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	// 获取源数据的时候讲数据放入缓存供下次读取
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
