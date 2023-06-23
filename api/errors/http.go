package errors

type HttpError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func UnauthorizedError() HttpError {
	return HttpError{
		401,
		"Unauthorized",
		"You are not authorized to access this resource",
	}
}

func TokenNotActive() HttpError {
	return HttpError{
		401,
		"Unauthorized",
		"Token is not active",
	}
}

func NotFoundError() *HttpError {
	return &HttpError{
		404,
		"Not found",
		"The requested resource was not found",
	}
}

func BadRequestError(message string) *HttpError {
	return &HttpError{
		400,
		"Bad Request",
		message,
	}
}
