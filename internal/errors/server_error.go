package errors

func NewServerError(s string) *ServerError {
	return &ServerError{
		Message: s,
	}
}

type ServerError struct {
	Message string
}

func (e *ServerError) Error() string {
	return e.Message
}
