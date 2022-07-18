package gee

import (
	"net/http"
)

// HandlerFunc 定义一个 request Handler 的类型
type HandlerFunc func(c *Context)

// Engine
// @Description: 实现了 ServeHTTP 的 Engine
type Engine struct {
	router *router
}

// New
// @Description: Engine 构造器
// @return *Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

// addRoute
// @Description: 新增一个路由
// @receiver engine
// @param method	请求方法
// @param pattern	路由 pattern
// @param handler	路由处理器
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// ServeHTTP
// @Description: 实现 ServeHTTP
// @receiver engine
// @param w	 Response 返回
// @param req Request 请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
