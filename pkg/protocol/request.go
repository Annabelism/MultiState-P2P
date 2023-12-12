package protocol

import (
	"fmt"
	"net"
    "encoding/json"
)

// Define the types of actions that can be used in update requests
type Action string

const (
	Add    Action = "add"
	Delete Action = "delete"
	Remove Action = "remove"
)

// Define the types of requests that can be handled
type RequestType string

const (
	ConnectionReqType RequestType = "Connection"
	DownloadReqType   RequestType = "Download"
	UpdateReqType     RequestType = "Update"
)

// Request is a wrapper for different types of requests in the network
type Request struct {
	Type    RequestType
	Payload interface{} // Payload can be any of the defined request types
}

// ConnectionRequest is used when a new device wants to join the network
type ConnectionRequest struct {
	DestinationIP string
	AccessToken   string
}

// DownloadRequest is used to request a file from another node
type DownloadRequest struct {
	DestinationIP string
	FileName      string
}

// UpdateTuple represents the changes to be made in the update request
type UpdateTuple struct {
	Action Action
	Index  string
	Value  string
}

// UpdateRequest is used to modify Table H
type UpdateRequest struct {
	Updates []UpdateTuple
}

// IsValid checks the validity of the action in UpdateTuple
func (ut *UpdateTuple) IsValid() bool {
	switch ut.Action {
	case Add, Delete, Remove:
		return true
	default:
		return false
	}
}

// SendRequest sends a request over a TCP connection
func SendRequest(conn net.Conn, req interface{}) error {
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