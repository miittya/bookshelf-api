package handler

type Response struct {
	Error string `json:"error,omitempty"`
}

type Status struct {
	Status string `json:"status"`
}

func OK() Status {
	return Status{Status: "OK"}
}

func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}
