package util

// MyError is a custom error type
type CanceledRequestError string

// Error implements the error interface for MyError
func (e CanceledRequestError) Error() string {
	return string(e)
}

func CanceledRequest() error {
	return CanceledRequestError("Request Canceled")
}
