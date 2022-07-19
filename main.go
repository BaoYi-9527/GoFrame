package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	// 独立的 {} 相当于隔离出一块作用域
	v1 := r.Group("v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=Bob
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("v2")
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
