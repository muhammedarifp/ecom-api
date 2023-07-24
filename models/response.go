package models

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error"`
}
