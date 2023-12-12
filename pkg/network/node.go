package network

import (
	"MultiState-P2P/pkg/protocol"
)

// NodeState represents the state of a node in the network.
type NodeState string

const (
	Idle    NodeState = "Idle"
	Share   NodeState = "Share"
	Request NodeState = "Request"
	Update  NodeState = "Update"
	Transmit NodeState = "Transmit"
	Dead    NodeState = "Dead" // Inactive state
)

// Node represents a peer in the network.
type Node struct {
	IP         string
	State      NodeState
	Buffer     []protocol.Request // Assuming Request is a struct that handles requests.
	TableH     map[string][]string // Map of file names to a list of IPs that store them.
	AccessToken string
}

// NewNode creates a new node with the given IP and access token.
func NewNode(IP, accessToken string) *Node {
	return &Node{
		IP:         IP,
		State:      Idle,
		Buffer:     make([]protocol.Request, 0),
		TableH:     make(map[string][]string),
		AccessToken: accessToken,
	}
}

// UpdateTableH updates the node's Table H based on the action required.
func (n *Node) UpdateTableH(action protocol.Action, file string, nodeIP string) {
	// Implementation of updating Table H goes here
	// Depending on the action, it will add, delete, or remove an entry
}
