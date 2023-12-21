package main

import (
	"MultiState-P2P/pkg/network"
	"fmt"
	"os"
)

func main() {
	// Initializing a node and connecting to the network
	ip, err := network.getLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		os.Exit(1)
	}

	// Define a port
	port := "8888"

	// Combine IP and port
	my_IP := ip + ":" + port
	peer_IP := "some ip address"

	myNode := network.NewNode(my_IP, "myAccessToken")
	conn, err := network.ConnectToNetwork(myNode, peer_IP) // IP of a known peer
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// The node would then listen for incoming requests and handle them accordingly
}
