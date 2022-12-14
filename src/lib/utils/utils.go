package utils

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(response http.ResponseWriter, httpStatus int, responseStruct any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(httpStatus)
	json.NewEncoder(response).Encode(responseStruct)
}
