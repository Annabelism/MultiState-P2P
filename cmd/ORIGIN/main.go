package main

import (
	"MultiState-P2P/pkg/network"
	"fmt"
	"os"
	"net"
)

func BuildORIGIN() {
	// Retrieve the local IP address
	ip, err := network.GetLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		os.Exit(1)
	}

	// Define the port on which the ORIGIN peer will listen
	port := "6666"
	my_IP := ip + ":" + port

	// Initialize the ORIGIN node with no initial connections
	myNode := network.NewNode(my_IP, "myAccessToken")

	// Listen on the specified port for incoming connections
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error setting up listener on port %s: %v\n", port, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("ORIGIN node initialized and listening on %s\n", my_IP)

	// Main loop to handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting incoming connection: %v\n", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go myNode.HandleRequest(conn)
	}

	// Additional logic for graceful shutdown, signal handling, etc., can be added here
}

func main() {
	BuildORIGIN()
}