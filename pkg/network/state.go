package network

import (
	"MultiState-P2P/pkg/protocol"
	"fmt"
)

// Transition handles state transitions for a node.
func (n *Node) Transition(event Event) {
	switch n.State {
	case Idle:
		n.handleIdleState(event)
	case Share:
		n.handleShareState(event)
	case Request:
		n.handleRequestState(event)
	case Update:
		n.handleUpdateState(event)
	case Transmit:
		n.handleTransmitState(event)
	case Dead:
		n.handleDeadState(event)
	}
}

// Event represents an event that can trigger a state transition.
type Event struct {
	Type        EventType
	File        string
	SourceIP    string
	AccessToken string
	Error       string
}

// EventType represents the type of an event.
type EventType string

// Define your event types here, e.g., ReceivedUpdateRequest, ErrorOccurred, etc.
const (
	ReceivedUpdateRequest EventType = "ReceivedUpdateRequest"
	ConnectionRequest     EventType = "ConnectionRequest"
	ReceivedFileRequest   EventType = "ReceivedFileRequest"
	ErrorOccurred         EventType = "ErrorOccurred"
	// ... other event types
)

// These handle*State functions implement the logic for each state.
func (n *Node) handleIdleState(event Event) {
	switch event.Type {
	case ReceivedUpdateRequest:
		n.UpdateTableH(protocol.Add, event.File, event.SourceIP)
	case ConnectionRequest:
		if event.AccessToken == n.AccessToken {
			n.TableH.AddEntry("", event.SourceIP)
		} else {
			fmt.Println("Invalid access token for connection request.")
		}
	case ErrorOccurred:
		// Handle any errors that may have occurred.
		fmt.Println("Error occurred in Idle state:", event.Error)
	default:
		fmt.Println("Unhandled event type received in Idle state.")
	}
}

func (n *Node) handleShareState(event Event) {
	switch event.Type {
	case ReceivedUpdateRequest:
		// Node receives an update request while sharing; it updates TableH with the new file information.
		n.UpdateTableH(protocol.Add, event.File, event.SourceIP)
	case ErrorOccurred:
		fmt.Println("Error occurred in Share state:", event.Error)
	default:
		fmt.Println("Unhandled event type in Share state.")
	}
}

func (n *Node) handleRequestState(event Event) {
	switch event.Type {
	case ReceivedFileRequest:
		// Check if the requested file is available in the node's TableH.
		nodesWithFile := n.TableH.GetNodesWithFile(event.File)
		if len(nodesWithFile) > 0 {
			downloadRequest := protocol.DownloadRequest{
				DestinationIP: nodesWithFile[0], // Selecting the first node for simplicity.
				FileName:      event.File,
			}
			// Send the download request to the selected node.
			protocol.SendRequest(conn, downloadRequest)
		} else {
			// File not found in the network, handle accordingly.
			fmt.Printf("File %s not found in the network.\n", event.File)
		}
	case ErrorOccurred:
		fmt.Println("Error occurred in Request state:", event.Error)
	default:
		fmt.Println("Unhandled event type in Request state.")
	}
}

func (n *Node) handleUpdateState(event Event) {
	switch event.Type {
	case ReceivedUpdateRequest:
		n.UpdateTableH(protocol.Add, event.File, event.SourceIP)
	case ErrorOccurred:
		fmt.Println("Error occurred in Update state:", event.Error)
	default:
		fmt.Println("Unhandled event type in Update state.")
	}
}

func (n *Node) handleTransmitState(event Event) {
	switch event.Type {
	case ConnectionRequest:
		// Logic for handling file transmission to another node
		// This can include sending the file data to the requester
	case ErrorOccurred:
		fmt.Println("Error occurred in Transmit state:", event.Error)
	default:
		fmt.Println("Unhandled event type in Transmit state.")
	}
}

func (n *Node) handleDeadState(event Event) {
	// Typically, no actions are taken in the Dead state.
	fmt.Println("Node is in Dead state, no actions taken.")
}
