package api

import (
	"encoding/json"
	"net/http"
)

// handlerError is a custom error type.
type handlerError struct {
	Error   error
	Message string
	Code    int
}

type handler func(http.ResponseWriter, *http.Request) (interface{}, *handlerError)

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := fn(w, r)
	if err != nil {
		http.Error(w, err.Message, err.Code)
		return
	}

	data, err2 := json.Marshal(res)
	if err2 != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
