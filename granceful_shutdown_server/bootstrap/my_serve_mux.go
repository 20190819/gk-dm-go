package bootstrap

import (
	"net/http"
)

type MyServeMux struct {
	reject bool
	*http.ServeMux
}

func (mux *MyServeMux) ServerHTTP(w http.ResponseWriter, r *http.Request) {

	// 拒绝请求
	if mux.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("服务器关闭"))
		return
	}

	mux.ServeHTTP(w, r)
}
