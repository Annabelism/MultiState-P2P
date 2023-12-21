package main

import (
	"MultiState-P2P/pkg/network"
	"fmt"
	"os"
	"strings"
	"sync"
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
	peer_IP := "some ip address"

	myNode := network.NewNode(my_IP, "myAccessToken")
	fmt.Printf("node built for peer A\n")
	conn, err := network.ConnectToNetwork(myNode, peer_IP)
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
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter 'request' to start an request: ")
			text, _ := reader.ReadString('\n')
			input := strings.TrimSuffix(strings.ToLower(text), "\n")
			terminalInput <- input
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
			network.MakeRequest(myNode)

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
