package network

import (
    "MultiState-P2P/pkg/protocol"
    "net"
    "fmt"
)

// ConnectToNetwork tries to connect to the network via an existing peer and returns the established TCP connection.
func ConnectToNetwork(n *Node, peerIP string) (net.Conn, error) {
    conn, err := net.Dial("tcp", peerIP)
    if err != nil {
        return nil, err
    }

    // Send a connection request to the peer
    req := protocol.ConnectionRequest{
        DestinationIP: n.IP,
        AccessToken:   n.AccessToken,
    }

    // Send the update request to the peer
    err = protocol.SendRequest(conn, req)
    if err != nil {
        conn.Close() // Close the connection in case of an error after establishing it
        return nil, fmt.Errorf("error @ connection request to peer %s: %w", peerIP, err)
    }

    // Return the established connection
    return conn, nil
}

// DisconnectFromNetwork handles the disconnection of a node from the network.
func DisconnectFromNetwork(n *Node) error {
    // Iterate through all stored connections to send update requests
    for peerIP, conn := range n.Connections {
        if conn != nil {
            // Create an update request to remove this node's entries
            updateReq := protocol.UpdateTuple{
                Action: protocol.Remove, // Assuming 'Remove' is a defined action in the 'Action' type
                Index:  n.IP,            // Assuming 'n.IP' is the IP of the current node
                Value:  "",              // Value might be empty or contain relevant data based on your implementation
            }

            // Send the update request to the peer
            err := protocol.SendRequest(conn, updateReq)
            if err != nil {
                return fmt.Errorf("error sending update request to peer %s: %w", peerIP, err)
            }
        }
    }

    // After sending updates, close all connections
    for peerIP, conn := range n.Connections {
        if conn != nil {
            err := conn.Close()
            if err != nil {
                return fmt.Errorf("error closing connection to peer %s: %w", peerIP, err)
            }
        }
    }

    // Optionally, clear the connections map after closing all connections
    n.Connections = make(map[string]net.Conn)

    return nil
}


// HandleConnectionRequest handles incoming connection requests.
func HandleConnectionRequest(n *Node, req protocol.ConnectionRequest) protocol.ConnectionResponse {
    // Check if the access token is correct
    if req.AccessToken != n.AccessToken {
        return protocol.ConnectionResponse{
            Status:  protocol.Unauthorized,
            Message: "Invalid access token",
        }
    }

    // Update Table H with new node's information
    n.TableH.AddEntry("", req.DestinationIP) // Example entry, adjust according to your logic

    // Propagate the updated Table H to all peers, including the new node
    // Get all peers from the node's TableH
    allPeers := n.TableH.GetAllPeers()

    // Iterate through all peers
    for _, peerIP := range allPeers {
        // Establish a TCP connection to the peer
        conn, err := net.Dial("tcp", peerIP)
        if err != nil {
            fmt.Printf("error connecting to peer %s: %v\n", peerIP, err)
        }
        defer conn.Close()

        // Create an update request to remove this node's entries
        updateReq := protocol.UpdateTuple{
            Action: protocol.Add, // Assuming 'Remove' is a defined action in the 'Action' type
            Index:  req.DestinationIP,   // Assuming 'n.IP' is the IP of the current node
            Value:  "",     // Value might be empty or contain relevant data based on your implementation
        }

        // Send the update request to the peer
        err = protocol.SendRequest(conn, updateReq)
        if err != nil {
            fmt.Printf("error connecting to peer %s: %v\n", peerIP, err)
        }
    }

    return protocol.ConnectionResponse{
        Status:  protocol.Success,
        Message: "Connection successful",
    }
}

