package routes

import "net/http"

func SigninHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		SigninPost(response, request)
	}
}

func SigninPost(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
}
