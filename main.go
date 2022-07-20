package main

import (
	"fmt"
	"gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.Default()

	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static") // 指定静态资源路径

	stu1 := &student{Name: "Bob", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}

	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})

	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	// 独立的 {} 相当于隔离出一块作用域
	v1 := r.Group("/v1")
	{
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
