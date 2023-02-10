package demo

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

type HandleFunc func(ctx *context.Context)

type Server interface {
	http.Handler
	Start(addr string) error
	AddRoute(method, path string, handler HandleFunc)
}

type HTTPServer struct {
	req  *http.Request
	resp http.ResponseWriter
	router
}

type Context struct {
	Req    *http.Request
	Resp   http.ResponseWriter
	Params map[string]string
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:    request,
		Resp:   writer,
		Params: make(map[string]string),
	}

	s.Next(ctx)
}

func (s *HTTPServer) Next(ctx *Context) {
	// 找路由

	// 执行业务逻辑
}

func (s *HTTPServer) AddRoute(method, path string, handler HandleFunc) {
	s.addRoute(method, path, handler)
}

func (s *HTTPServer) Start(addr string) error {

	// 方式01
	// 阻塞方法后面无法做些什么
	//return http.ListenAndServe(addr, s)

	// 方式02
	// 启动端口
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// 此处可以做一些其他事情，如 服务注册
	// ...
	log.Println("启动服务中")
	return http.Serve(listen, s)

}

func StartWeb(addr string) {
	var s Server = &HTTPServer{}

	s.AddRoute(http.MethodGet, "/user", func(ctx *context.Context) {
		fmt.Println("/user")
	})
	s.AddRoute(http.MethodGet, "/order/*", func(ctx *context.Context) {
		fmt.Println("/orders/*")
	})

	err := s.Start(addr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
