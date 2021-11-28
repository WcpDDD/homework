package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/all-f-0/golang/homework/http_server/src/handles"
	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	LogBufferSize     = 100
	DefaultConfigFile = "./http-server.yaml"
)

type requestLog struct {
	ip   string
	code int
}

type serverAppConfig struct {
	Port int
}

type serverConfig struct {
	App serverAppConfig
}

type httpServer struct {
	server     *http.Server
	config     serverConfig
	mutex      sync.Mutex
	exitLogger chan bool
}

func (server *httpServer) stopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.server.Shutdown(ctx); err != nil {
		glog.Fatalln("服务器优雅终止失败")
	}
	// 关闭日志
	server.exitLogger <- true
}

func main() {
	configPath := flag.String("config-path", DefaultConfigFile, "配置文件路径")

	flag.Parse()
	defer func() {
		glog.Flush()
	}()

	serverConfig := loadServerConfig(configPath)

	server := startServer(serverConfig)
	stopWatch := make(chan bool)

	done := make(chan os.Signal, 1)

	if watcher, err := fsnotify.NewWatcher(); err != nil {
		glog.Warning("配置文件监听器创建失败")
	} else {
		watchConfigFile(watcher, func() {
			glog.V(5).Infoln("监听到配置文件发生修改")
			serverConfig = loadServerConfig(configPath)
			server = reloadConfig(serverConfig, server)
			glog.V(5).Infoln("配置文件已重新加载")
		}, stopWatch)
	}
	// 获取SIGTERM信号
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	// 关闭文件监听
	stopWatch <- true
	glog.V(5).Infoln("检测到服务器关闭信号")
	server.stopServer()
	glog.V(5).Infoln("服务器关闭成功")
}

// 判断配置文件更改是否需要重启服务器逻辑
func restartServerPredicate(config serverConfig, server *httpServer) bool {
	return config.App.Port != server.config.App.Port
}

func reloadConfig(config serverConfig, server *httpServer) *httpServer {
	server.mutex.Lock()
	defer func() {
		server.mutex.Unlock()
	}()
	// 如果App配置发生变化 则重启server 否则仅替换其中配置文件
	if restartServerPredicate(config, server) {
		glog.V(5).Infoln("服务器需要重启")
		server.stopServer()
		return startServer(config)
	} else {
		glog.V(5).Infoln("服务器不需要重启")
		server.config = config
		return server
	}
}

func watchConfigFile(watcher *fsnotify.Watcher, cb func(), stop chan bool) {
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				{
					switch event.Op {
					case fsnotify.Create:
						fallthrough
					case fsnotify.Write:
						cb()
					}
				}
			case err := <-watcher.Errors:
				{
					glog.Errorln("文件监听失败", err)
				}
			case <-stop:
				break
			}
		}
	}()
	watcher.Add(DefaultConfigFile)
}

func startServer(config serverConfig) *httpServer {
	mux := http.NewServeMux()

	logChan := make(chan requestLog, LogBufferSize)
	exitChan := make(chan bool, 1)
	go requestLogger(logChan, exitChan)
	registerHandle(handles.IndexHandle{}, mux, logChan)
	registerHandle(handles.Healthz{}, mux, logChan)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.App.Port),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			glog.V(5).Infof("http server 服务终止:%+v", err)
		}
	}()

	return &httpServer{
		server:     &server,
		config:     config,
		mutex:      sync.Mutex{},
		exitLogger: exitChan,
	}
}

// 加载服务器配置
func loadServerConfig(configPath *string) serverConfig {
	config := serverConfig{
		App: serverAppConfig{
			Port: 80,
		},
	}
	glog.V(5).Infof("加载配置文件")
	if file, err := ioutil.ReadFile(*configPath); err != nil {
		glog.Warning("配置文件加载失败，使用默认配置")
		config = serverConfig{
			App: serverAppConfig{
				Port: 80,
			},
		}
	} else {
		err := yaml.Unmarshal(file, &config)
		// 处理掉配置文件格式不正确的情况
		if err != nil {
			config = serverConfig{
				App: serverAppConfig{
					Port: 80,
				},
			}
		}
	}
	return config
}

// 记录用户访问信息
func requestLogger(ch chan requestLog, exitChan chan bool) {
	for {
		select {
		case rl := <-ch:
			glog.V(5).Infof("客户端地址:%s, 返回码:%d", rl.ip, rl.code)
		case <-exitChan:
			break
		}
	}
}

// 这里的io异常会不会导致和客户端的连接被挂起，直到超时？ 有没有什么处理方式
func sendResponse(statusCode int, body string, header http.Header, w http.ResponseWriter) {
	defer func() {
		if err := recover(); err != nil {
			glog.V(5).Infof("发送响应信息失败.")
		}
	}()
	for key, value := range header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(statusCode)
	if _, err := io.WriteString(w, body); err != nil {
		glog.Errorf("发送响应信息失败.")
	}
}

// 包装handle 处理异常及打印日志
func handleWrapper(h handles.Handle, ch chan requestLog) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusOK

		defer func() {
			if err := recover(); err != nil {
				glog.Errorf("请求处理失败:%+v\n", err)
				// 服务端异常
				sendResponse(http.StatusInternalServerError, "", http.Header{}, w)
			}
			ch <- requestLog{
				ip:   r.RemoteAddr,
				code: statusCode,
			}
		}()

		// 如果路径不匹配 则404
		if !strings.EqualFold(h.Path(), r.RequestURI) {
			statusCode = http.StatusNotFound
			sendResponse(statusCode, "", http.Header{}, w)
			return
		}

		method := r.Method
		if strings.EqualFold(method, h.Method()) {
			glog.V(1).Infof("处理请求:%s,%s\n", h.Path(), h.Method())
			h.Invoke(r, func(responseInfo handles.ResponseInfo, err error) {
				if err != nil {
					statusCode = http.StatusInternalServerError
					sendResponse(statusCode, "", responseInfo.Header, w)
					return
				}
				statusCode = http.StatusOK
				sendResponse(statusCode, responseInfo.Body, responseInfo.Header, w)
			})
		} else {
			// method不匹配
			statusCode = http.StatusMethodNotAllowed
			sendResponse(statusCode, "", http.Header{}, w)
		}
	}
}

func registerHandle(handle handles.Handle, mux *http.ServeMux, ch chan requestLog) {
	path := handle.Path()
	if len(path) > 0 {
		mux.HandleFunc(path, handleWrapper(handle, ch))
	}
}
