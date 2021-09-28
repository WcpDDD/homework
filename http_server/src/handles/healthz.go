package handles

import (
	"net/http"
)

type Healthz struct{}

func (Healthz) Invoke(r *http.Request, callback func(ResponseInfo, error)) {
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
