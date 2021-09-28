package handles

import (
	"github.com/all-f-0/golang/homework/util"
	"net/http"
)

type ResponseInfo struct {
	Header http.Header
	Body   string
}

func CreateResponseInfo(body string, header http.Header) ResponseInfo {
	headerOptional := util.NewOptional(header)
	return ResponseInfo{
		Body: body,
		Header: headerOptional.OrElseGet(func() interface{} {
			return http.Header{}
		}).(http.Header),
	}
}

type Handle interface {
	Invoke(r *http.Request, callback func(ResponseInfo, error))
	Path() string
	Method() string
}
