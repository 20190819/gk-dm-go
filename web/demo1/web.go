package demo1

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type HandleFunc func(ctx *context.Context)

type Server interface {
	http.Handler
	Start(addr string) error
	addRoute(method, path string, handler HandleFunc)
}

type HTTPServer struct {
	req  *http.Request
	resp http.ResponseWriter
}

type Context struct {
	Req    *http.Request
	Writer http.ResponseWriter
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}

func (s *HTTPServer) addRoute(method, path string, handler HandleFunc) {
	//
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

func Start(addr string) {
	var s Server = &HTTPServer{}
	var h1 HandleFunc = func(ctx *context.Context) {
		fmt.Println("step 1 ")
		time.Sleep(time.Second)
	}

	var h2 HandleFunc = func(ctx *context.Context) {
		fmt.Println("step 2 ")
		time.Sleep(time.Second)
	}

	s.addRoute(http.MethodPost, "/user", func(ctx *context.Context) {
		h1(ctx)
		h2(ctx)
	})

	err := s.Start(addr)
	if err != nil {
		panic("http server start failed")
	}
}
