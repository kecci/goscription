package models

// BaseResponse base response
type BaseResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   []string    `json:"error"`
}
