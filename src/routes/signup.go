package routes

import "net/http"

func SignupHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		SignupPost(response, request)
	}
}

func SignupPost(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
}
