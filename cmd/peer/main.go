package main

import (
	"MultiState-P2P/pkg/network"
	"fmt"
	"os"
)

func main() {
	// Use GetLocalIP from network package
	ip, err := network.GetLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		os.Exit(1)
	}

	port := "8888"
	my_IP := ip + ":" + port
	peer_IP := "10.205.0.199:6666"

	myNode := network.NewNode(my_IP, "myAccessToken")
	fmt.Printf("node built for peer A\n")
	conn, err := network.ConnectToNetwork(myNode, peer_IP)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return // Stop execution if there's an error
	}
	fmt.Print("established connection\n")
	network.PrintTableH(myNode.TableH)
	defer conn.Close()

	// Now you can use BuildConnections to establish connections with peers
	err = network.BuildConnections(myNode)
	if err != nil {
		fmt.Println("Error building connections:", err)
		return
	}

	// Rest of your main function...

	// On program exit, disconnect from the network
	defer network.DisconnectFromNetwork(myNode)
}
