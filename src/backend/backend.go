package backend

import "net/http"

func RouteHandler(httpMethod string, handler func(response http.ResponseWriter, request *http.Request)) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		if request.Method == httpMethod {
			handler(response, request)
		}
	}
}
