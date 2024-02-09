package errors

func NewApplicationError(s string) *ApplicationError {
	return &ApplicationError{
		Message: s,
	}
}

type ApplicationError struct {
	Message string
}

func (e *ApplicationError) Error() string {
	return e.Message
}
