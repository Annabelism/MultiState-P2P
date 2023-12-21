package main

import (
	"MultiState-P2P/pkg/network"
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {
	// Initializing a node and connecting to the network
	ip, err := network.GetLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		os.Exit(1)
	}

	// Define a port
	port := "8880"

	// Combine IP and port
	my_IP := ip + ":" + port
	peer_IP := "localhost:8888"

	myNode := network.NewNode(my_IP, "myAccessToken")
	_, err = network.ConnectToNetwork(myNode, peer_IP) // IP of a known peer
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// The node would then listen for incoming requests and terminal input, and handle them accordingly

	// Channel for terminal input
	terminalInput := make(chan string)

	// Channel for incoming data from TCP connections
	tcpInput := make(chan string)

	// WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup

	// Start a goroutine to listen for terminal input
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Print("Choose your request: download(1), update(2), cancel request(x)\n")
		for {
			input_req, _ := network.ReadFromConsole()
			terminalInput <- input_req
		}
	}()

	// Start a goroutine to simulate incoming data from TCP connections
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Simulate opening TCP connections (replace with your actual code)
		network.BuildConnections(myNode)

		listener, err := net.Listen("tcp", ":8888")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			go myNode.HandleRequest(conn)
		}
	}()

	// Main loop with select statement for handling multiple channels
	for {
		select {
		case input := <-terminalInput:
			// Handle terminal input
			fmt.Println("Received from terminal:", input)

			//create and handle request
			network.MakeRequest(myNode, input)

			// // Check if the input is 'exit' to quit the program
			// if input == "exit\n" {
			// 	fmt.Println("Exiting program.")
			// 	// Close the TCP listener to stop accepting new connections
			// 	// (You may need additional logic to gracefully close existing connections)
			// 	os.Exit(0)
			// }

		case tcpData := <-tcpInput:
			// Handle data received from TCP connections
			fmt.Println("Received from TCP connection:", tcpData)

		default:
			// Do other work or continue looping
			continue
		}
	}
}

// Handle incoming data from a TCP connection
func handleTCPConnection(conn net.Conn, tcpInput chan<- string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from TCP connection:", err)
			return
		}
		tcpInput <- data
	}
}
