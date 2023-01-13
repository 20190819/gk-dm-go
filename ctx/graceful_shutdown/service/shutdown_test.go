package service

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

// TestGraceFulShutdown 测试优雅退出
func TestGraceFulShutdown(t *testing.T) {
	s1 := NewMyServer("biz", ":8070")
	s2 := NewMyServer("biz", ":8071")

	s1.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("graceful server 1"))
	}))

	app := NewApp([]*MyServer{s1, s2}, WithShutdownCallbacks(CacheToDataBase))
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
