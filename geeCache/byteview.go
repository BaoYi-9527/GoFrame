// Package geeCache
// @Description: 缓存值的抽象与封装
package geeCache

// ByteView
// @Description: 只读数据结构 用于存储缓存值
type ByteView struct {
	b []byte // 存储真实的缓存值
}

// Len
// @Description: 获取缓存值的长度
// @receiver v
// @return int
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice
// @Description: 获取数据拷贝 防止缓存值被外部程序修改
// @receiver v
// @return []byte
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
