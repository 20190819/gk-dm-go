package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type shutdownCallbackFunc func(ctx context.Context)
type option func(*MyApp)

func WithCallbackOption(cb shutdownCallbackFunc) option {

	return func(app *MyApp) {
		app.CallBacks = append(app.CallBacks, cb)
	}
}

type MyApp struct {
	Servers         []*MyServer
	ShutdownTimeout time.Duration
	WaitTimeout     time.Duration
	CallbackTimeout time.Duration
	CallBacks       []shutdownCallbackFunc
}

func NewMyApp(servers []*MyServer, ops ...option) *MyApp {
	app := &MyApp{
		Servers:         servers,
		ShutdownTimeout: 30 * time.Second,
		WaitTimeout:     10 * time.Second,
		CallbackTimeout: 3 * time.Second,
	}

	for _, option := range ops {
		option(app)
	}
	return app
}

func (app *MyApp) StartAndService() {

	for _, server := range app.Servers {
		srv := server

		go func() {
			log.Printf("准备启动%s服务\n", srv.name)
			err := srv.Start() // 阻塞后续执行
			if err != nil {
				if err == http.ErrServerClosed {
					log.Printf("服务器%s已关闭", srv.name)
				} else {
					log.Printf("服务器%s异常关闭: %v\n", srv.name, err.Error())
				}
			}
		}()

		// 开启另外的协程 ping 服务地址，ping 通则成功，否则失败
		go func() {
			srv.checkStartSuccess()
		}()
	}

	// 监听系统信号
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, SignalWin...)
	<-ch // 阻塞
	go func() {
		select {
		case <-ch:
			log.Println("强制退出")
			os.Exit(1)
		case <-time.After(app.ShutdownTimeout):
			os.Exit(1)
		}
	}()
	app.Shutdown()
}

func (app *MyApp) Shutdown() {
	log.Println("开始关闭应用，停止接收新请求")
	for _, server := range app.Servers {
		server.RejectRequest()
	}

	log.Println("等待进行中的请求完结")
	time.Sleep(app.WaitTimeout * time.Second)

	log.Println("===开始关闭服务器===")
	// 并发关闭服务器
	// 协调所有 server 关闭后进行下一步
	var wg sync.WaitGroup
	wg.Add(len(app.Servers))
	ctx := context.Background()
	for _, server := range app.Servers {
		srvCopy := server
		go func() {
			if err := srvCopy.Stop(ctx); err != nil {
				log.Printf("关闭服务器失败 %s \n", srvCopy.name)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("开始执行自定义回调")
	wg.Add(len(app.CallBacks))
	for _, backFunc := range app.CallBacks {
		backFuncCopy := backFunc
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), app.CallbackTimeout)
			backFuncCopy(ctx)
			cancel()
			wg.Done()
		}()
	}
	wg.Wait()

	// 释放资源
	log.Println("释放资源...")
	app.Close()

}

func (app *MyApp) Close() {
	// 释放资源
	time.Sleep(time.Second)
	log.Println("应用已经关闭 Done")
}
