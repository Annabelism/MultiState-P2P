package network

import (
	"MultiState-P2P/pkg/protocol"
	"net"
    "encoding/json"
    "fmt"
    "io"
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
	TableH     *TableH // Map of file names to a list of IPs that store them.
	AccessToken string
	Connections map[string]net.Conn
}

// NewNode creates a new node with the given IP and access token.
func NewNode(IP, accessToken string) *Node {
	return &Node{
		IP:         IP,
		State:      Idle,
		Buffer:     make([]protocol.Request, 0),
		TableH:     NewTableH(),
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
	switch req.Type {
    case protocol.UpdateReqType:
		updateReq := req.Payload
	
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
	
	case protocol.DownloadReqType:
		downloadReq := req.Payload
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
		
		ConnectionResp := HandleConnectionRequest(n, req.Payload)
		fmt.Printf("%#v\n", ConnectionResp)
    }
    return nil
}