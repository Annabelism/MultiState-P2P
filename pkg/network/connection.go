package network

import (
	"encoding/json"
	"errors"

	"fmt"
	"net"
)

func ConnectToNetwork(n *Node, peerIP string) (string, error) {
	if peerIP == n.IP {
		n.TableH.AddEntry(n.IP, "")
		return "", errors.New("no connection is made, peerIP = node.IP")
	}

	// Establish a TCP connection

	conn, err := net.Dial("tcp", peerIP)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	fmt.Println("DIALED TO ", peerIP)

	// Construct the request data
	requestData := map[string]string{
		"type":        "connection",
		"accessToken": n.AccessToken,
		"myIP":        n.IP,
	}

	// Marshal the request data into JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// Send the request
	_, err = conn.Write(jsonData)
	if err != nil {
		return "", err
	}

	fmt.Println("FINISHED SENDING, WAITING RESPONSE")

	// Wait for the response
	//responseData, err := io.ReadAll(conn)
	buffer := make([]byte, 1024)

	r, err := conn.Read(buffer)

	if err != nil {
		return "", err
	}
	responseData := buffer[:r]

	fmt.Println("RESPONSE RECEIVED: ", responseData)

	return string(responseData), nil
}
