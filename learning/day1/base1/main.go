package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// HandleFunc 注册一个处理器函数 handler 和对应的模式 pattern（注册到 DefaultServeMux ）。

	// ServeMux 的文档解释了模式的匹配机制。
	// ServeMux 类型是 HTTP 请求的多路转接器。它会将每一个接收的请求的 URL 与一个注册模式的列表进行匹配，并调用和 URL 最匹配的模式的处理器。

	// Handler根据 r.Method、r.Host和 r.URL.Path 等数据，返回将用于处理该请求的 HTTP 处理器。它总是返回一个 非nil 的处理器。
	// 如果路径不是它的规范格式，将返回内建的用于重定向到等价的规范路径的处理器。
	// Handler 也会返回匹配该请求的的已注册模式；在内建重定向处理器的情况下，pattern 会在重定向后进行匹配。
	// 如果没有已注册模式可以应用于该请求，本方法将返回一个内建的 "404 page not found" 处理器和一个空字符串模式。
	http.HandleFunc("/", indexHandler) // indexHandler ->  ServeHTTP(w ResponseWriter, r *Request)
	http.HandleFunc("/hello", helloHandler)
	// ListenAndServe 监听 srv.Addr 指定的 TCP 地址，并且会调用 Serve 方法接收到的连接。
	// 如果 srv.Addr 为空字符串，会使用":http"。
	// Fatal 等价于 {Print(v...); os.Exit(1)}
	log.Fatal(http.ListenAndServe(":9999", nil))
}

// ResponseWriter 接口被 HTTP 处理器用于构造 HTTP 回复。
// Request 类型代表一个服务端接受到的或者客户端发送出去的 HTTP 请求。
func indexHandler(w http.ResponseWriter, req *http.Request) {
	// Fprintf根据 format 参数生成格式化的字符串并写入 w。返回写入的字节数和遇到的任何错误。
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
