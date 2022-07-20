package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// 原始对象
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息
	Path   string
	Method string
	Params map[string]string
	// 响应信息
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int // index 记录当前执行到的中间件的索引
	// engine pointer
	engine *Engine
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next
// @Description: 调用中间件
// @receiver c
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// PostForm
// @Description:  获取 POST 请求参数
// @receiver c
// @param key
// @return string
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query
// @Description: 获取 Query String 中的参数
// @receiver c
// @param key
// @return string
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status
// @Description: 设置响应状态码 code
// @receiver c
// @param code
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader
// @Description: 设置请求头
// @receiver c
// @param key
// @param value
func (c *Context) SetHeader(key string, value string) {
	// Header 返回响应的头部，该值会被 WriterHeader 方法发送
	// 在 WriterHeader 或 Writer 方法后再改变该对象是没有意义的
	c.Writer.Header().Set(key, value)
}

// String
// @Description: 向HTTP链接中写入回复数据
// @receiver c
// @param code
// @param format
// @param values...interface{}
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	// Write 向连接中写入作为 HTTP 的一部分回复的数据。
	// 如果被调用时还未调用 WriteHeader，本方法会先调用 WriteHeader(http.StatusOK)
	// 如果 Header 中没有 "Content-Type" 键，
	// 本方法会使用包函数 DetectContentType 检查数据的前 512 字节，将返回值作为该键的值。

	// Sprintf 根据 format 参数生成格式化的字符串并返回该字符串。
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON
// @Description: 返回 JSON 类型数据
// @receiver c
// @param code
// @param obj
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "text/json")
	c.Status(code)
	// NewDecoder 创建一个从 c.Writer 读取并解码 json 对象的 *Decoder，解码器有自己的缓冲，并可能超前读取部分 json 数据。
	encoder := json.NewEncoder(c.Writer)
	// Encode将 obj 的 json 编码写入输出流，并会写入一个换行符
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
