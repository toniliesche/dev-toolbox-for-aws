package model

type ErrorResponse struct {
	Type    string `json:"__type"`
	Message string `json:"Message"`
}
