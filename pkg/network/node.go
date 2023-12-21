package network

import (
	"MultiState-P2P/pkg/protocol"
	"MultiState-P2P/pkg/util"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

// NodeState represents the state of a node in the network.
type NodeState string

const (
	Idle     NodeState = "Idle"
	Share    NodeState = "Share"
	Request  NodeState = "Request"
	Update   NodeState = "Update"
	Transmit NodeState = "Transmit"
	Dead     NodeState = "Dead" // Inactive state
)

// Node represents a peer in the network.
type Node struct {
	IP          string
	State       NodeState
	Buffer      []protocol.Request // Assuming Request is a struct that handles requests.
	TableH      *TableH            // Map of file names to a list of IPs that store them.
	AccessToken string
	Connections map[string]net.Conn
}

// NewNode creates a new node with the given IP and access token.
func NewNode(IP, accessToken string) *Node {
	return &Node{
		IP:          IP,
		State:       Idle,
		Buffer:      make([]protocol.Request, 0),
		TableH:      NewTableH(),
		AccessToken: accessToken,
		Connections: make(map[string]net.Conn),
	}
}

func (n *Node) HandleRequest(conn net.Conn) error {
	// Get the sender's address
	senderAddr := conn.RemoteAddr().String()

	// Split the address into IP and port
	senderIP, senderPort, err := net.SplitHostPort(senderAddr)
	if err != nil {
		return fmt.Errorf("failed to parse remote address: %w", err)
	}

	fmt.Printf("Received connection from %s:%s\n", senderIP, senderPort)
	// Add logic here to handle the request

	jsonData, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Errorf("error reading request data: %w", err)
	}

	// Unmarshal JSON data into a Request struct
	var req protocol.Request
	err = json.Unmarshal(jsonData, &req)
	if err != nil {
		return fmt.Errorf("error unmarshalling request JSON: %w", err)
	}

	// Process the request based on the content of the req object
	// Check the request type and process accordingly
	fmt.Println("REQUEST TYPE: ", req.Type)
	switch req.Type {
	case protocol.UpdateReqType:
		HandleUpdateRequest(n, req)

	case protocol.DownloadReqType:
		downloadReq, ok := req.Payload.(protocol.DownloadRequest)
		if !ok {
			return fmt.Errorf("payload is not a DownloadRequest")
		}
		filename := downloadReq.FileName
		downloadResp, err := protocol.CreateDownloadResponse("here is your file", filename)
		if err != nil {
			return err
		}
		err = protocol.SendResponse(conn, downloadResp)
		if err != nil {
			return err
		}

	case protocol.ConnectionReqType:
		connectionReq, ok := req.Payload.(protocol.ConnectionRequest)
		if !ok {
			return fmt.Errorf("payload is not a ConnectionRequest")
		}
		ConnectionResp := HandleConnectionRequest(n, connectionReq)
		fmt.Printf("%#v\n", ConnectionResp)
	}
	return nil
}

func HandleUpdateRequest(n *Node, req protocol.Request) error {
	fmt.Println("RECEIVED UPDATE REQ")
	updateReq, ok := req.Payload.(protocol.UpdateRequest)
	if !ok {
		return fmt.Errorf("payload is not an UpdateRequest")
	}
	// Process each UpdateTuple
	for _, updateTuple := range updateReq.Updates {
		fmt.Printf("Action: %s, Index: %s, Value: %s\n", updateTuple.Action, updateTuple.Index, updateTuple.Value)
		switch updateTuple.Action {
		case protocol.Add:
			fmt.Printf("Add Action - Index: %s\n", updateTuple.Index)
			n.TableH.AddEntry(updateTuple.Index, updateTuple.Value)

		case protocol.Delete:
			fmt.Printf("Delete Action - Index: %s\n", updateTuple.Index)
			n.TableH.RemoveEntry(updateTuple.Index, updateTuple.Value)

		case protocol.Remove:
			fmt.Printf("Remove Action - Index: %s\n", updateTuple.Index)
			n.TableH.RemoveNode(updateTuple.Index)

		default:
			// Handle unknown action
			fmt.Printf("Unknown Action: %s\n", updateTuple.Action)
		}
	}
	fmt.Println("TABLE H AFTER UPDATE: ")
	PrintTableH(n.TableH)
	return nil
}

func (n *Node) CreateRequest(input string) error {

	timeout := 5 * time.Second
	timeoutCh := time.After(timeout)

	// Perform operations in a loop until timeout
	for {
		select {
		case <-timeoutCh:
			// Timeout occurred, exit the loop
			//Need to implement: go back to idle state
			fmt.Println("Timeout reached. Exiting loop.")
			return util.TimeoutError("Timeout")
		default:
			req, err := MakeRequest(n, input)
			if err != nil {
				fmt.Println("Error making request.")
				continue
			}
			req_type := req.Type
			var peer_ip string
			switch req_type {
			case protocol.DownloadReqType:
				peer_ip = req.Payload.(protocol.DownloadRequest).DestinationIP
			case protocol.UpdateReqType:
				fmt.Println("ENTERED UPDATE REQ")

				peers := n.TableH.GetAllNodes() // Assuming GetAllPeers returns a slice of peer IP addresses
				fmt.Println("LEN: ", len(peers))
				for _, peer_ip := range peers {
					fmt.Println(peer_ip)
					conn, err := BuildConnections(n, peer_ip)
					if err != nil {
						fmt.Println("Error building connections for update.")
						continue
					}
					protocol.SendRequest(conn, req)

				}
				PrintTableH(n.TableH)
				return nil
			}
			conn, err := BuildConnections(n, peer_ip)
			if err != nil {
				fmt.Println("Error building connections.")
				continue
			}
			protocol.SendRequest(conn, req)
			fmt.Println("Request sent.")
		}
	}
}

func MakeRequest(n *Node, input_req string) (protocol.Request, error) {
	//read from console. From https://freshman.tech/snippets/go/read-console-input/
	//can be moved to node later

	switch input_req {
	case "download", "1":
		fmt.Print("Enter the file name you're requesting (or enter 'x' to cancel request): \n")
		// read the file name from user input
		input_file, err := ReadFromConsole()
		if err != nil {
			empty_request := &protocol.Request{
				Type:    protocol.EmptyReqType,
				Payload: nil,
			}
			return *empty_request, err
		}
		if input_file == "x" {
			fmt.Println("Request Canceled.")
			fmt.Print("Choose your request: download(1), update(2), cancel request(x)\n")
			empty_request := &protocol.Request{
				Type:    protocol.EmptyReqType,
				Payload: nil,
			}
			return *empty_request, util.CanceledRequestError("Request Canceled")
		}
		//change later. get one node that has the target file
		destination_ips := n.TableH.GetNodesWithFile(input_file)
		if len(destination_ips) == 0 {
			fmt.Println("File Not Found")
			empty_request := &protocol.Request{
				Type:    protocol.EmptyReqType,
				Payload: nil,
			}
			return *empty_request, util.FileNotFoundError("File not found")
		}
		req := protocol.CreateDownloadRequest(input_file, destination_ips[0])
		return req, nil
	case "update", "2":
		fmt.Print("Enter the action you want to make: add, delete, remove, cancel request(x)\n")
		// read the file name from user input
		input_action, err := ReadFromConsole()
		if err != nil {
			empty_request := &protocol.Request{
				Type:    protocol.EmptyReqType,
				Payload: nil,
			}
			return *empty_request, err
		}
		fmt.Print("Enter the file you want to update (enter 'x' to cancel request): \n")
		input_file, err := ReadFromConsole()
		if err != nil {
			empty_request := &protocol.Request{
				Type:    protocol.EmptyReqType,
				Payload: nil,
			}
			return *empty_request, err
		}
		if input_file == "x" {
			fmt.Println("Request Canceled.")
			fmt.Print("Choose your request: download(1), update(2), cancel request(x)\n")
			empty_request := &protocol.Request{
				Type:    protocol.EmptyReqType,
				Payload: nil,
			}
			return *empty_request, util.CanceledRequestError("Request Canceled")
		}
		err = CheckFile(n, input_file)
		if err != nil {
			fmt.Println("Check file failed.")
			empty_request := &protocol.Request{
				Type:    protocol.EmptyReqType,
				Payload: nil,
			}
			return *empty_request, util.FileNotFoundError("File not found")
		}
		fmt.Println("CREATING REQUEST: ", input_action, input_file, n.IP)
		req := protocol.CreateUpdateRequest(input_action, n.IP, input_file)
		fmt.Println("REQUEST CREATED: ", input_action, input_file, n.IP)
		return req, nil
	case "x":
		fmt.Println("Request Canceled.")
		fmt.Print("Choose your request: download(1), update(2), cancel request(x)\n")
		empty_request := &protocol.Request{
			Type:    protocol.EmptyReqType,
			Payload: nil,
		}
		return *empty_request, util.CanceledRequestError("Request Canceled")
	default:
		fmt.Println("Invalid input. Please try again.")
		fmt.Print("Choose your request: download(1), update(2), cancel request(x)\n")
		empty_request := &protocol.Request{
			Type:    protocol.EmptyReqType,
			Payload: nil,
		}
		return *empty_request, util.InvalidInputError("Invalid Input")
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

func CheckFile(n *Node, file_name string) error {
	file_path := "file/" + file_name

	if _, err := os.Stat(file_path); err != nil {
		fmt.Println("ERROR checking file:", err)
		return err
	}
	return nil
}
