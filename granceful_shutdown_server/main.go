package main

import (
	"context"
	"github.com/golang-module/carbon/v2"
	"gk-dm-go/granceful_shutdown_server/bootstrap"
	"log"
	"net/http"
)

func main() {
	s1 := bootstrap.NewMyServer("users", ":8005")
	s2 := bootstrap.NewMyServer("posts", ":8006")

	s1.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("test graceful shutdown my server 1 \n" + carbon.Now().ToDateTimeString()))
	}))
	s2.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("test graceful shutdown my server 2 \n" + carbon.Now().ToDateTimeString()))
	}))

	app := bootstrap.NewMyApp([]*bootstrap.MyServer{s1, s2}, bootstrap.WithCallbackOption(cache))
	app.StartAndService()
}

func cache(context.Context) {
	// todo something
	log.Println("自定义缓存回调")
}
