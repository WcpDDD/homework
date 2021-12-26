package common

import (
	"context"
	"github.com/golang/glog"
	"net/http"
	"sync"
	"time"
)

type ServerAppConfig struct {
	Port int
}

type HandleConfig struct {
	NextService string `yaml:"next-service"`
}

type ServerConfig struct {
	App    ServerAppConfig
	Handle HandleConfig
}

type HttpServer struct {
	Server     *http.Server
	Config     ServerConfig
	Mutex      sync.Mutex
	ExitLogger chan bool
}

func (server *HttpServer) StopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.Server.Shutdown(ctx); err != nil {
		glog.Fatalln("服务器优雅终止失败")
	}
	// 关闭日志
	server.ExitLogger <- true
}
