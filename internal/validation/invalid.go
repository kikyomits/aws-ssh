package validation

type InvalidInput struct {
	Message string
}

func (e *InvalidInput) Error() string {
	return e.Message
}

func NewInvalidInputError(reason string) error {
	return &InvalidInput{
		Message: "invalid input: " + reason,
	}
}
