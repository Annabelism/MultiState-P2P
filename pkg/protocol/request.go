package protocol

// Define the types of actions that can be used in update requests
type Action string

const (
	Add    Action = "add"
	Delete Action = "delete"
	Remove Action = "remove"
)

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

// Authenticate validates the access token in the ConnectionRequest
func (cr *ConnectionRequest) Authenticate() bool {
	// This should be replaced with actual authentication logic
	return cr.AccessToken == "validToken"
}