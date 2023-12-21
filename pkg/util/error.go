package util

// MyError is a custom error type
type CanceledRequestError string
type InvalidInputError string

// Error implements the error interface for MyError
func (e CanceledRequestError) Error() string {
	return string(e)
}

// Error implements the error interface for MyError
func (e InvalidInputError) Error() string {
	return string(e)
}
