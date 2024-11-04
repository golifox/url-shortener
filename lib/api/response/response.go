package response

type Response struct {
	Status      string `json:"status"`
	Error       string `json:"error,omitempty"`
	EncodedLink string `json:"encoded_link,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func OKWithEncodedLink(encodedLink string) Response {
	return Response{
		Status:      StatusOK,
		EncodedLink: encodedLink,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
