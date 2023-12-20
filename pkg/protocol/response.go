package protocol

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

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

func CreateDownloadResponse(message string, file_name string) (*DownloadResponse, error) {
	file_path := "../../files/" + file_name
	file, err := os.ReadFile(file_path)
	if err != nil {
		return nil, err
	}
	return &DownloadResponse{
		Status:  Success,
		Message: message,
		File:    file,
	}, nil
}

func CreateUpdateResponse(message string) *UpdateResponse {
	return &UpdateResponse{
		Status:  Success,
		Message: message,
	}
}

// SendResponse sends a response over a TCP connection
func SendResponse(conn net.Conn, req interface{}) error {
	// Convert the request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshalling request: %w", err)
	}

	// Send the JSON data
	_, err = conn.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	return nil
}
