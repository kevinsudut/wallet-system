package response

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func WriteErrorResponse(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, content interface{}) {
	bJson, err := jsoniter.Marshal(content)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(bJson)
}
