package handlers

type APIError struct {
	Message string `json:"message"`
}

func makeApiError(err error) *APIError {
	return &APIError{
		Message: err.Error(),
	}
}
