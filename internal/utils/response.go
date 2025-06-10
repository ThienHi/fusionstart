package utils

type CustomResponse struct {
	Code    int         `json:"code"`
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Response(code int, error bool, message string, data interface{}) *CustomResponse {
	return &CustomResponse{
		Code:    code,
		Error:   error,
		Message: message,
		Data:    data,
	}
}
