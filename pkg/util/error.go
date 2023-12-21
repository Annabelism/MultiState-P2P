package util

// MyError is a custom error type
type CanceledRequestError string
type InvalidInputError string
type TimeoutError string
type FileNotFoundError string

// Error implements the error interface for MyError
func (e CanceledRequestError) Error() string {
	return string(e)
}

// Error implements the error interface for MyError
func (e InvalidInputError) Error() string {
	return string(e)
}

func (e TimeoutError) Error() string {
	return string(e)
}

func (e FileNotFoundError) Error() string {
	return string(e)
}
