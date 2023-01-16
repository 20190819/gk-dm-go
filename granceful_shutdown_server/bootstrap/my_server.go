package bootstrap

import (
	"context"
	"log"
	"net/http"
	"time"
)

type ServerInterface interface {
	Handle(pattern string, handler http.Handler)
	Start() error
	RejectRequest()
	Stop()
}

type MyServer struct {
	srv  *http.Server
	addr string
	name string
	mux  *MyServeMux
}

func NewMyServer(name, addr string) *MyServer {
	mux := &MyServeMux{
		ServeMux: http.NewServeMux(),
	}
	return &MyServer{
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		addr: addr,
		name: name,
		mux:  mux,
	}
}

func (ms *MyServer) Handle(pattern string, handler http.Handler) {
	ms.mux.Handle(pattern, handler)
}

func (ms *MyServer) Start() error {
	return ms.srv.ListenAndServe()
}

func (ms *MyServer) Stop(ctx context.Context) error {
	log.Printf("服务器%s关闭中", ms.name)
	return ms.srv.Shutdown(ctx)
}

func (ms *MyServer) RejectRequest() {
	ms.mux.reject = true
}

func (ms *MyServer) checkStartSuccess() {
	for {
		time.Sleep(time.Second)
		log.Printf("检查服务%s是否启动\n", ms.name)
		resp, err := http.Get("http://" + ms.addr + "/")
		if err != nil {
			log.Printf("服务%s启动失败\n", ms.name)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("服务%s启动失败\n", ms.name)
			continue
		}
		break
	}
	log.Printf("服务%s启动成功\n", ms.name)
}
