package handles

import (
	"net/http"
	"net/url"
	"os"
	"strings"
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

func (IndexHandle) Invoke(r *http.Request, callback func(ResponseInfo, error)) {
	defer func() {
		if err := recover(); err != nil {
			callback(ResponseInfo{}, err.(error))
		}
	}()

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
