package response

// CustomError custom error
type CustomError struct {
	httpStatus int
	Code       int
	Message    string
}

func (cErr CustomError) Error() string {
	return cErr.Message
}
