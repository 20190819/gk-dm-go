package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type ShutdownCallback func(ctx context.Context)
type Option func(*App)

func WithShutdownCallbacks(cbs ...ShutdownCallback) Option {

	return func(app *App) {
		app.cbs = cbs
	}
}

type MyServer struct {
	name string
	srv  *http.Server
	mux  *serverMux
}

func NewMyServer(name, addr string) *MyServer {

	mux := &serverMux{
		ServeMux: http.NewServeMux(),
	}
	return &MyServer{
		name: name,
		mux:  mux,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *MyServer) Start() error {
	return s.srv.ListenAndServe()
}

func (s *MyServer) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *MyServer) RejectRequest() {
	s.mux.reject = true
}

func (s *MyServer) stop(ctx context.Context) error {
	log.Printf("服务器%s关闭中", s.name)
	return s.srv.Shutdown(ctx)
}

// serverMux 可以看作装饰器模式或代理模式
type serverMux struct {
	reject bool
	*http.ServeMux
}

func (s *serverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if s.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("服务器已关闭"))
		return
	}
	s.ServeMux.ServeHTTP(w, r)
}

type App struct {
	servers         []*MyServer
	shutdownTimeout time.Duration // 优雅退出整个超时时间 默认 30 s
	waitTime        time.Duration // 等待处理已有请求,默认 10 s
	cbTimeout       time.Duration // 自定义回调超时,默认 3 s
	cbs             []ShutdownCallback
}

func NewApp(servers []*MyServer, opts ...Option) *App {
	app := &App{
		servers:         servers,
		shutdownTimeout: 30,
		waitTime:        10,
		cbTimeout:       3,
	}

	for _, opt := range opts {
		opt(app)
	}
	return app
}

func (app *App) StartAndService() {

	for _, server := range app.servers {
		srv := server
		go func() {
			err := srv.Start()
			if err != nil {
				if err == http.ErrServerClosed {
					log.Printf("服务器%s关闭中\r", srv.name)
				} else {
					log.Printf("服务器%s异常退出\n", srv.name)
				}
			}
		}()
		log.Printf("服务器 %s 启动成功\n", srv.name)
	}

	// 监听系统信号
	// 定义 channel 接收系统信号
	// 定义要监听的系统信号量 signals []os.Signal
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, signalWindows...)
	<-ch
	log.Println("hello")
	go func() {
		select {
		case <-ch:
			log.Println("强制退出")
			os.Exit(1)
		case <-time.After(app.shutdownTimeout):
			log.Println("超时强制退出")
			os.Exit(1)
		}
	}()

	app.ShunDown()
}

func (app *App) ShunDown() {
	log.Println("开始关闭应用，停止接收新请求")
	for _, server := range app.servers {
		server.RejectRequest()
	}

	log.Println("等待正在执行请求完结")
	time.Sleep(app.waitTime)

	log.Println("开始关闭服务器")
	// 并发关闭服务器，同时要注意协调所有的 server 都关闭之后才能步入下一个阶段
	var wg sync.WaitGroup
	wg.Add(len(app.servers))
	ctx := context.Background()
	for _, server := range app.servers {
		srvCopy := server
		go func() {
			if err := srvCopy.stop(ctx); err != nil {
				log.Printf("关闭服务器失败 %s \n", srvCopy.name)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("开始执行自定义回调")
	// 并发执行回调，要注意协调所有的回调都执行完才会步入下一个阶段
	wg.Add(len(app.cbs))
	for _, cb := range app.cbs {
		callbackFunc := cb
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), app.cbTimeout)
			callbackFunc(ctx)
			cancel()
			wg.Done()
		}()
	}
	wg.Wait()

	// 释放资源
	log.Println("开始释放资源")
	app.Close()
}

func (app *App) Close() {
	// 释放可能的一些资源
	time.Sleep(time.Second)
	log.Println("应用关闭")
}
