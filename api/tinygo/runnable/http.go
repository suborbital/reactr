package runnable

// #include "reactr.h"
import "C"
import (
	"github.com/suborbital/reactr/api/tinygo/runnable/method"
)

func GET(url string, headers map[string]string) ([]byte, error) {
	return doRequest(method.GET, url, nil, headers)
}

func HEAD(url string, headers map[string]string) ([]byte, error) {
	return doRequest(method.HEAD, url, nil, headers)
}

func OPTIONS(url string, headers map[string]string) ([]byte, error) {
	return doRequest(method.OPTIONS, url, nil, headers)
}

func POST(url string, body []byte, headers map[string]string) ([]byte, error) {
	return doRequest(method.POST, url, body, headers)
}

func PUT(url string, body []byte, headers map[string]string) ([]byte, error) {
	return doRequest(method.PUT, url, body, headers)
}

func PATCH(url string, body []byte, headers map[string]string) ([]byte, error) {
	return doRequest(method.PATCH, url, body, headers)
}

func DELETE(url string, headers map[string]string) ([]byte, error) {
	return doRequest(method.DELETE, url, nil, headers)
}

// Remark: The URL gets encoded with headers added on the end, seperated by ::
// eg. https://google.com/somepage::authorization:bearer qdouwrnvgoquwnrg::anotherheader:nicetomeetyou
func doRequest(method method.MethodType, url string, body []byte, headers map[string]string) ([]byte, error) {
	urlStr := url

	if headers != nil {
		headerStr := renderHeaderString(headers)
		if headerStr != "" {
			urlStr += "::" + headerStr
		}
	}

	urlPtr, urlSize := rawSlicePointer([]byte(urlStr))
	bodyPtr, bodySize := rawSlicePointer(body)

	size := C.fetch_url(int32(method), urlPtr, urlSize, bodyPtr, bodySize, ident())

	return result(size)
}

func renderHeaderString(headers map[string]string) string {
	out := ""

	for key, value := range headers {
		out += key + ":" + value
		out += "::"
	}

	return out[:len(out)-2]
}
