package controllers

type Response[T any] struct {
	Data    T      `json:"data,omitempty"`
	Success bool   `json:"success"`
	Errors  string `json:"errors,omitempty"`
}
