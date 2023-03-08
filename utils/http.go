package utils

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Data    T      `json:"data,omitempty"`
	Success bool   `json:"success"`
	Errors  string `json:"errors,omitempty"`
}

func HttpError(w http.ResponseWriter, status int, message string) {
	resp := Response[struct{}]{Success: false, Errors: message}
	w.WriteHeader(status)
	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

func HttpSuccess[T any](w http.ResponseWriter, status int, data *T) {
	var resp Response[T]

	if data != nil {
		resp = Response[T]{Success: true, Data: *data}
	} else {
		resp = Response[T]{Success: true}
	}

	w.WriteHeader(status)
	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}
