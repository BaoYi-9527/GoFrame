package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// start timer
		t := time.Now()
		// Process request
		c.Next()
		// 记录日志
		log.Printf("[%d] %s in %v for group v2-logger", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
