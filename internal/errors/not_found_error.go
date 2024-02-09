package errors

func NewNotFoundError() *NotFoundError {
	return &NotFoundError{}
}

type NotFoundError struct{}

func (*NotFoundError) Error() string {
	return "Entity not found"
}
