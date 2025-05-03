package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Error(message string) Response {
	return Response{
		Status: "error",
		Error:  message,
	}
}
