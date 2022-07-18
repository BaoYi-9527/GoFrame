package main

import (
	"fmt"
	"log"
	"net/http"
)

// Engine
// @Description: Engine 是一个独立的用于处理所有请求的引擎
type Engine struct {
}

// ServeHTTP
// @Description: 实现 Handler 接口
// 在 Go 语言中只要实现了一个接口中的所有方法就相当于实现了这个接口【非侵入式接口】
// Handler接口地址： net/http/server.go
// type Handler interface {
//	ServeHTTP(ResponseWriter, *Request)
// }
// @receiver engine
// @param w
// @param req
func (engine Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 拦截所有的 HTTP 请求 解析 URL Path 分发给对应的处理方法
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.PATH = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	engine := new(Engine)
	// 启动 http 监听服务
	log.Fatal(http.ListenAndServe(":9999", engine))
}
