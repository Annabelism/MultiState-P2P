package network

import (
	"MultiState-P2P/pkg/protocol"
	"bufio"
	"fmt"
	"os"
	"strings"
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
	Type EventType
	// ... other fields
}

// EventType represents the type of an event.
type EventType string

// Define your event types here, e.g., ReceivedUpdateRequest, ErrorOccurred, etc.
const (
	ReceivedUpdateRequest EventType = "ReceivedUpdateRequest"
	ErrorOccurred         EventType = "ErrorOccurred"
	// ... other event types
)

// These handle*State functions implement the logic for each state.
func (n *Node) handleIdleState(event Event) {
	// Implement logic for when the node is in the Idle state
}

func (n *Node) handleShareState(event Event) {
	// Implement logic for when the node is in the Share state
}

func (n *Node) handleRequestState(event Event) {
	// Implement logic for when the node is in the Request state

}

func (n *Node) handleUpdateState(event Event) {
	// Implement logic for when the node is in the Update state
}

func (n *Node) handleTransmitState(event Event) {
	// Implement logic for when the node is in the Transmit state
}

func (n *Node) handleDeadState(event Event) {
	// Implement logic for when the node is in the Dead state
}

func MakeRequest(n *Node) interface{} {
	//read from console. From https://freshman.tech/snippets/go/read-console-input/
	//can be moved to node later
	var req protocol.Request
	fmt.Print("Choose your operation: Download(1), Update(2)")
	input_req, err := ReadFromConsole()
	if err != nil {
		fmt.Print("Choose your operation: Download(1), Update(2)")
		return nil
	} else {
		switch input_req {
		case "download", "1":
			req.Type = protocol.DownloadReqType
			fmt.Print("Enter the file name you're requesting: ")
			// ReadString will block until the delimiter is entered
			input_file, err := ReadFromConsole()
			if err != nil {
				return nil
			} else {
				//temp. get one node that has the target file
				destination_ip := n.TableH.GetNodesWithFile(input_file)[0]
				req.Payload = protocol.DownloadRequest{
					DestinationIP: destination_ip,
					FileName:      input_file,
				}
			}

		case "update", "2":
			req.Type = protocol.UpdateReqType
			return nil

		default:
			fmt.Println()
			return nil
		}
	}
	return nil
}

func ReadFromConsole() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return "", err
	}
	// remove the delimeter from the string
	input = strings.TrimSuffix(strings.ToLower(input), "\n")
	return input, nil
}
