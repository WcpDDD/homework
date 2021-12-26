package handles

import (
	"fmt"
	"github.com/all-f-0/golang/homework/http_server/src/common"
	"github.com/golang/glog"
	"net/http"
	"strings"
)

type TraceHandle struct{}

func (TraceHandle) Invoke(r *http.Request, server *common.HttpServer, callback func(ResponseInfo, error)) {
	defer func() {
		if err := recover(); err != nil {
			callback(ResponseInfo{}, err.(error))
		}
	}()

	header := http.Header{}
	copyHeaders(r, &header)
	addVersionHeaders(&header)

	header.Add("canary", "canary")

	response := CreateResponseInfo("", header)

	// 如果下一条服务链路不为空 则执行访问
	if len(server.Config.Handle.NextService) > 0 {
		if err := remote(server.Config.Handle.NextService, header); err != nil {
			callback(response, err)
			return
		}
	}

	callback(response, nil)
}

func (TraceHandle) Path() (url string) {
	url = "/trace"
	return
}

func (TraceHandle) Method() (method string) {
	method = "Get"
	return
}

func remote(url string, header http.Header) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("%s", err)
	}
	lowerCaseHeader := make(http.Header)
	for key, value := range header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	glog.Info("headers:", lowerCaseHeader)
	req.Header = lowerCaseHeader
	client := &http.Client{}
	if _, err := client.Do(req); err != nil {
		return err
	}
	return nil
}
