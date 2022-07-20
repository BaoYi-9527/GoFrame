package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc 定义一个 request Handler 的类型
type HandlerFunc func(c *Context)

// RouterGroup
// @Description: 路由分组
type RouterGroup struct {
	prefix      string        // 前缀
	middlewares []HandlerFunc // 中间件
	parent      *RouterGroup  // 支持嵌套
	engine      *Engine       // 所有的分组共享一个 Engine 示例
}

// Engine
// @Description: 实现了 ServeHTTP 的 Engine
type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup     // 存储所有的分组
	htmlTemplates *template.Template // HTML render
	funcMap       template.FuncMap   // HTML render
}

// New
// @Description: Engine 构造器
// @return *Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}             // 构造一个 Engine
	engine.RouterGroup = &RouterGroup{engine: engine}  // 构造一个路由分组 并且注入 当前 Engine
	engine.groups = []*RouterGroup{engine.RouterGroup} // 将当前 Engine 的路由分组 放入 Engine 的分组管理中
	return engine
}

// Default
// @Description: 默认使用 Logger 和 Recovery 中间件
// @return *Engine
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// Group
// @Description: 定义一个新的路由分组
// @PS: 所有路由分组共享一个 engine 所有分组都的 engine 都继承自父类 也就是都继承自根类
// @receiver group
// @param prefix
// @return *RouterGroup
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// createStaticHandler
// @Description: 静态资源处理器
// @receiver group
// @param relativePath
// @param fs
// @return HandlerFunc
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// 检查文件是否存在 以及 权限是否通过
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static
// @Description: 静态文件路由
// @receiver group
// @param relativePath
// @param root
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

// SetFuncMap
// @Description:
// @receiver engine
// @param funcMap
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob
// @Description:
// @receiver engine
// @param pattern
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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

// addRoute
// @Description: 分组新增路由
// @receiver group
// @param method
// @param comp
// @param handler
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Use
// @Description: 为路由组添加中间件
// @receiver group
// @param middlewares
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// ServeHTTP
// @Description: 实现 ServeHTTP
// @receiver engine
// @param w	 Response 返回
// @param req Request 请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc

	// 为当前路由的上下文 Context 添加中间件
	for _, group := range engine.groups {
		// 如果当前路由包含 路由组的前缀就将路由组的中间件添加到路由的 上下文中
		// 匹配前缀的时候是连带 / 的
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine // 注入 engine
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
