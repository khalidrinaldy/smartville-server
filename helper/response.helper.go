package helper

type Result struct {
	Error   bool                     `json:"error"`
	Message string                   `json:"message"`
	Data    interface{} `json:"data"`
}

func ResultResponse(error bool, message string, data interface{}) Result {
	var result Result
	result.Error = error
	result.Message = message
	result.Data = data
	return result
}