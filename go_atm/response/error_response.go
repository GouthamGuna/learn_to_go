package response

type ErrorResponse struct {
	Error     string `json:"error"`
	ErrorType error  `json:"errorType"`
}
