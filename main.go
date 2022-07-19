package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func main() {
	r := gee.New()
	// 全局中间件
	r.Use(gee.Logger())
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	// 独立的 {} 相当于隔离出一块作用域
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=Bob
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	// TODO::要注意 路由组前缀一定要加 / 后续做中间件处理的时候使用路由匹配的时候会带上 / 否则匹配不到
	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 局部中间件
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/Bob
			c.String(http.StatusOK, "hello %s, you're at %s \n", c.Param("name"), c.Path)
		})

		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start Timer
		t := time.Now()
		// 若是服务端异常报 500？
		//c.Fail(http.StatusInternalServerError, "Internal Server Error")
		// 计算处理时间
		log.Printf("[%d] %s in %v for group v2....", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
