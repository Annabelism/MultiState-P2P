package main

import (
	"MultiState-P2P/pkg/network"
	"MultiState-P2P/pkg/protocol"
	"fmt"
)

func main() {
	// Initializing a node and connecting to the network
	my_IP := "my ip address"
	peer_IP := "some ip address"

	myNode := network.NewNode(my_IP, "myAccessToken")
	err := network.ConnectToNetwork(myNode, peer_IP) // IP of a known peer
	if err != nil {
        fmt.Printf("error: %v\n", err)
	}

	// The node would then listen for incoming requests and handle them accordingly
}
