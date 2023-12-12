package network

import (
    "MultiState-P2P/pkg/protocol"
    "net"
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
    // Assume we have a method to send the request over the connection
    protocol.SendRequest(conn, req)

    // Further logic to handle the response goes here

    return nil
}

// DisconnectFromNetwork handles the logic for a node to gracefully leave the network.
func DisconnectFromNetwork(n *Node) error {
    // Send an update request to all peers to remove this node's entries from Table H
    // Assume we have a method to broadcast the request to all peers
    // BroadcastUpdateRequest(n)

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
    // BroadcastUpdateToPeers(n, req.DestinationIP)

    return protocol.ConnectionResponse{
        Status:  protocol.Success,
        Message: "Connection successful",
    }
}

// BroadcastUpdateToPeers broadcasts the updated Table H to all peers.
func BroadcastUpdateToPeers(n *Node, newPeerIP string) {
    // Logic to broadcast the update to all peers
    // Iterate over all known peers and send the updated Table H
    // This might involve establishing a TCP connection to each peer and sending the data
}

// BroadcastUpdateRequest sends an update request to all peers when disconnecting.
func BroadcastUpdateRequest(n *Node) {
    // Logic to send an update request to all peers
    // This indicates to peers that they should remove entries related to this node from their Table H
}

