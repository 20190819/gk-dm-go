package main

import (
	"context"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"gk-dm-go/ctx/graceful_shutdown/service"
	"net/http"
)

func main() {
	s1 := service.NewMyServer("frontend", ":8070")
	s2 := service.NewMyServer("backend", ":8071")

	s1.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("graceful server " + carbon.Now().ToDateTimeString()))
	}))

	app := service.NewApp([]*service.MyServer{s1, s2}, service.WithShutdownCallbacks(CacheToDataBase))
	app.StartAndService()
}

func CacheToDataBase(ctx context.Context) {

	signal := make(chan struct{}, 1)
	go func() {
		fmt.Println("刷新缓存中...")
		signal <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("刷新缓存超时")
	case <-signal:
		fmt.Println("写入数据库成功")
	}
}
