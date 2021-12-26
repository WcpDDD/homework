package handles

import (
	"github.com/all-f-0/golang/homework/http_server/src/common"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type IndexHandle struct{}

/*
获取环境变量中version相关的配置 并添加进headers
*/
func addVersionHeaders(header *http.Header) {
	for _, envItem := range os.Environ() {
		eqIndex := strings.Index(envItem, "=")
		name := envItem[0:eqIndex]
		value := envItem[eqIndex+1:]
		if flag := strings.Index(strings.ToUpper(name), "VERSION"); flag >= 0 {
			header.Add(url.QueryEscape(name), url.QueryEscape(value))
		}
	}
}

/*
将request中的headers添加进response
*/
func copyHeaders(r *http.Request, header *http.Header) {
	for key, value := range r.Header {
		for _, v := range value {
			header.Add(key, v)
		}
	}
}
func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func (IndexHandle) Invoke(r *http.Request, server *common.HttpServer, callback func(ResponseInfo, error)) {
	defer func() {
		if err := recover(); err != nil {
			callback(ResponseInfo{}, err.(error))
		}
	}()

	// 添加延时
	time.Sleep(time.Duration(randInt(500, 3000)) * time.Millisecond)

	header := http.Header{}
	copyHeaders(r, &header)
	addVersionHeaders(&header)

	response := CreateResponseInfo("", header)
	callback(response, nil)
	//callback(response, fmt.Errorf("xxx错误"))
}

func (IndexHandle) Path() (url string) {
	url = "/"
	return
}

func (IndexHandle) Method() (method string) {
	method = "Get"
	return
}
