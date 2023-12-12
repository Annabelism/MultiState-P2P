package protocol

// Define standard response statuses
type Status string

const (
	Success      Status = "success"
	Fail         Status = "fail"
	Unauthorized Status = "unauthorized"
)

// ConnectionResponse is the response to a ConnectionRequest
type ConnectionResponse struct {
	Status  Status
	Message string
}

// DownloadResponse is the response to a DownloadRequest
type DownloadResponse struct {
	Status  Status
	Message string
	File    []byte // The file data should be sent as bytes
}

// UpdateResponse is the response to an UpdateRequest
type UpdateResponse struct {
	Status  Status
	Message string
}

// CreateSuccessResponse creates a success response with a custom message
func CreateSuccessResponse(message string) *ConnectionResponse {
	return &ConnectionResponse{
		Status:  Success,
		Message: message,
	}
}

// CreateFailResponse creates a fail response with a custom message
func CreateFailResponse(message string) *ConnectionResponse {
	return &ConnectionResponse{
		Status:  Fail,
		Message: message,
	}
}

// CreateUnauthorizedResponse creates an unauthorized response
func CreateUnauthorizedResponse() *ConnectionResponse {
	return &ConnectionResponse{
		Status:  Unauthorized,
		Message: "Access token is invalid or expired.",
	}
}
