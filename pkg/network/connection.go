package network

import (
    "MultiState-P2P/pkg/protocol"
    "net"
    "fmt"
)

// ConnectToNetwork tries to connect to the network via an existing peer.
func ConnectToNetwork(n *Node, peerIP string) error {
    conn, err := net.Dial("tcp", peerIP)
    if err != nil {
        return err
    }
    defer conn.Close()

    // Send a connection request to the peer
    req := protocol.ConnectionRequest{
        DestinationIP: n.IP,
        AccessToken:   n.AccessToken,
    }

    // Send the update request to the peer
    err = protocol.SendRequest(conn, req)
    if err != nil {
        return fmt.Errorf("error @ connection request to peer %s: %w", peerIP, err)
    }

    return nil
}

// DisconnectFromNetwork handles the disconnection of a node from the network.
func DisconnectFromNetwork(n *Node) error {
    // Get all peers from the node's TableH
    allPeers := n.TableH.GetAllPeers()

    // Iterate through all peers
    for _, peerIP := range allPeers {
        // Establish a TCP connection to the peer
        conn, err := net.Dial("tcp", peerIP)
        if err != nil {
            return fmt.Errorf("error connecting to peer %s: %w", peerIP, err)
        }
        defer conn.Close()

        // Create an update request to remove this node's entries
        updateReq := protocol.UpdateTuple{
            Action: protocol.Remove, // Assuming 'Remove' is a defined action in the 'Action' type
            Index:  n.IP,   // Assuming 'n.IP' is the IP of the current node
            Value:  "",     // Value might be empty or contain relevant data based on your implementation
        }

        // Send the update request to the peer
        err = protocol.SendRequest(conn, updateReq)
        if err != nil {
            return fmt.Errorf("error sending update request to peer %s: %w", peerIP, err)
        }
    }

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

