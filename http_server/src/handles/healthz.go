package handles

import (
	"github.com/all-f-0/golang/homework/http_server/src/common"
	"net/http"
)

type Healthz struct{}

func (Healthz) Invoke(r *http.Request, server *common.HttpServer, callback func(ResponseInfo, error)) {
	response := CreateResponseInfo("", http.Header{})
	callback(response, nil)
}

func (Healthz) Method() (method string) {
	method = "Get"
	return
}

func (Healthz) Path() (path string) {
	path = "/healthz"
	return
}
