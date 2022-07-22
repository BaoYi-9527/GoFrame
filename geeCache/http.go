// Package geeCache
// @Description: 提供供其他节点访问的能力
package geeCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

// HTTPPool
// @Description:
type HTTPPool struct {
	self     string
	basePath string
}

// NewHTTPPool
// @Description: HTTPPool 实例化
// @param self	记录自己的地址，包括主机名/IP 和 端口
// @return *HTTPPool 节点间通讯地址前缀
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// SeverHTTP
// @Description:
// @receiver p
// @param w
// @param r
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// 分组与键名
	groupName := parts[0]
	key := parts[1]
	// 实例化分组
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group:"+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
