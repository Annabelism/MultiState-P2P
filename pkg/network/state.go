package network

import (
	"MultiState-P2P/pkg/protocol"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
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

	timeout := 5 * time.Second
	timeoutCh := time.After(timeout)

	// Perform operations in a loop until timeout
	for {
		select {
		case <-timeoutCh:
			// Timeout occurred, exit the loop
			//Need to implement: go back to idle state
			fmt.Println("Timeout reached. Exiting loop.")
			return
		default:
			//get peerIP, return connection
			conn, err := ConnectToNetwork(n, peerIP)
			// Continue performing your operations here
			fmt.Println("Performing operation...")
			// Simulate some work
			req, err := MakeRequest(n)
			if err != nil {
				continue
			}
			//pass in peer node
			protocol.SendRequest(conn, req)
		}
	}
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

func MakeRequest(n *Node) (interface{}, error) {
	//read from console. From https://freshman.tech/snippets/go/read-console-input/
	//can be moved to node later
	fmt.Print("Choose your operation: download(1), update(2)")
	input_req, err := ReadFromConsole()
	if err != nil {
		fmt.Print("Choose your operation: Download(1), Update(2)")
		return nil, err
	}
	switch input_req {
	case "download", "1":
		fmt.Print("Enter the file name you're requesting: add(1), delete(2), remove(3)")
		// read the file name from user input
		input_file, err := ReadFromConsole()
		if err != nil {
			return nil, err
		}
		//change later. get one node that has the target file
		destination_ip := n.TableH.GetNodesWithFile(input_file)[0]
		req := protocol.CreateDownloadRequest(input_file, destination_ip)
		return req, nil
	case "update", "2":
		fmt.Print("Enter the action you want to make: ")
		// read the file name from user input
		input_action, err := ReadFromConsole()
		if err != nil {
			return nil, err
		}
		fmt.Print("Enter the file you want to update: ")
		input_index, err := ReadFromConsole()
		if err != nil {
			return nil, err
		}
		req := protocol.CreateUpdateRequest(input_action, input_index, n.IP)
		return req, nil
	default:
		fmt.Println("Invalid input. Please try again.")
		return nil, err
	}
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
