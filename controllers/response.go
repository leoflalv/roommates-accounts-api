package controllers

type Response[T any] struct {
	Data    T      `json:"data"`
	Success bool   `json:"success"`
	Errors  string `json:"errors"`
}
