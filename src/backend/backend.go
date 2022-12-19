package backend

import "net/http"

type MethodMap struct {
	Get     func(response http.ResponseWriter, request *http.Request)
	Post    func(response http.ResponseWriter, request *http.Request)
	Patch   func(response http.ResponseWriter, request *http.Request)
	Put     func(response http.ResponseWriter, request *http.Request)
	Delete  func(response http.ResponseWriter, request *http.Request)
	Trace   func(response http.ResponseWriter, request *http.Request)
	Options func(response http.ResponseWriter, request *http.Request)
	Connect func(response http.ResponseWriter, request *http.Request)
	Head    func(response http.ResponseWriter, request *http.Request)
}

var Routes map[string]string = map[string]string{}

func Route(key string, path string, methodMap MethodMap) {
	Routes[key] = path

	http.HandleFunc(path, func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			if methodMap.Get != nil {
				methodMap.Get(response, request)
			}
		case http.MethodPost:
			if methodMap.Post != nil {
				methodMap.Post(response, request)
			}
		case http.MethodPatch:
			if methodMap.Patch != nil {
				methodMap.Patch(response, request)
			}
		case http.MethodPut:
			if methodMap.Put != nil {
				methodMap.Put(response, request)
			}
		case http.MethodDelete:
			if methodMap.Delete != nil {
				methodMap.Delete(response, request)
			}
		case http.MethodTrace:
			if methodMap.Trace != nil {
				methodMap.Trace(response, request)
			}
		case http.MethodOptions:
			if methodMap.Options != nil {
				methodMap.Options(response, request)
			}
		case http.MethodConnect:
			if methodMap.Connect != nil {
				methodMap.Connect(response, request)
			}
		case http.MethodHead:
			if methodMap.Head != nil {
				methodMap.Head(response, request)
			}
		}
	})
}
