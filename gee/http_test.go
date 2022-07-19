package gee

import (
	"fmt"
	"reflect"
	"testing"
)

/*
PS D:\Go\src\GoFrame\gee> go test
2022/07/19 13:53:57 Route  GET - /
2022/07/19 13:53:57 Route  GET - /hello/:name
2022/07/19 13:53:57 Route  GET - /hello/b/c
2022/07/19 13:53:57 Route  GET - /hi/:name
2022/07/19 13:53:57 Route  GET - /asset/*filepath
matched path: /hello/:name, params['name']: Bob
PASS
ok      gee     0.598s
*/

func newTestRoute() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/asset/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRoute()
	n, ps := r.getRoute("GET", "/hello/Bob")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "Bob" {
		t.Fatal("name should be equal to 'Bob' ")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
}
