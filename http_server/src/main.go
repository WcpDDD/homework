package main

import (
	"github.com/all-f-0/golang/homework/http_server/src/handles"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	LogBufferSize = 100
)

type requestLog struct {
	ip   string
	code int
}

func main() {
	logChan := make(chan requestLog, LogBufferSize)
	exitChan := make(chan bool, 1)
	go requestLogger(logChan, exitChan)
	registerHandle(handles.IndexHandle{}, logChan)
	registerHandle(handles.Healthz{}, logChan)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("http server 启动失败:%+v", err)
	}
}

// 记录用户访问信息
func requestLogger(ch chan requestLog, exitChan chan bool) {
	for {
		select {
		case rl := <-ch:
			log.Printf("客户端地址:%s, 返回码:%d", rl.ip, rl.code)
		case <-exitChan:
			break
		}
	}
}

// 这里的io异常会不会导致和客户端的连接被挂起，直到超时？ 有没有什么处理方式
func sendResponse(statusCode int, body string, header http.Header, w http.ResponseWriter) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("发送响应信息失败.")
		}
	}()
	for key, value := range header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(statusCode)
	if _, err := io.WriteString(w, body); err != nil {
		log.Printf("发送响应信息失败.")
	}
}

// 包装handle 处理异常及打印日志
func handleWrapper(h handles.Handle, ch chan requestLog) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusOK

		defer func() {
			if err := recover(); err != nil {
				log.Printf("请求处理失败:%+v\n", err)
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
			log.Printf("处理请求:%s,%s\n", h.Path(), h.Method())
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

func registerHandle(handle handles.Handle, ch chan requestLog) {
	path := handle.Path()
	if len(path) > 0 {
		http.HandleFunc(path, handleWrapper(handle, ch))
	}
}
