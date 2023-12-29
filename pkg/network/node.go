package network

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// Node represents a peer in the network.
type Node struct {
	IP          string
	TableH      *TableH // Map of file names to a list of IPs that store them.
	AccessToken string
}

// NewNode creates a new node with the given IP and access token.
func NewNode(IP, accessToken string) *Node {
	return &Node{
		IP:          IP,
		TableH:      NewTableH(),
		AccessToken: accessToken,
	}
}

func (n *Node) Broadcast() error {
	peers := n.TableH.GetAllNodes()
	for _, peerIP := range peers {
		if peerIP != n.IP {
			err := ShareTableH(n.TableH, peerIP, false)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (myNode *Node) HandleRequest(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	r, err := conn.Read(buffer)

	if err != nil {
		return
	}
	requestData := buffer[:r]

	fmt.Println("REQUEST DATA RECEIVED")

	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	// Unmarshal the JSON data into a map.
	var request map[string]interface{}
	err = json.Unmarshal(requestData, &request)
	if err != nil {
		fmt.Println("Error unmarshaling:", err.Error())
		return
	}

	// Type assertion for the "type" field
	requestType, ok := request["type"].(string)
	if !ok {
		fmt.Println("Error: type field is not a string")
		return
	}
	fmt.Println(requestType)
	if requestType == "connection" {
		fmt.Println("RECEIVED CONNECTION REQUEST")
		// Type assertion for the "accessToken" field
		accessToken, ok := request["accessToken"].(string)
		if !ok {
			fmt.Println("Error: accessToken field is not a string")
			return
		}

		if accessToken == "myAccessToken" {
			fmt.Println("TOKEN VARIFIED")
			conn.Write([]byte("Access granted\n"))

			// Type assertion for the "myIP" field
			myIP, ok := request["myIP"].(string)
			if !ok {
				fmt.Println("Error: myIP field is not a string")
				return
			}

			ShareTableH(myNode.TableH, myIP, true)
		} else {
			conn.Write([]byte("Access denied"))
		}
	}

	if requestType == "update-i" {
		fmt.Printf("TYPE FOR TABLEH: %T\n", request["tableH"])
		// Type assertion for the "tableH" field
		// tableH, ok := request["tableH"].(map[string]interface{})
		// if !ok {
		// 	fmt.Println("Error: tableH field is not the expected type")
		// 	return
		// }

		// Convert entries to the proper type (map[string][]string)
		// convertedEntries := make(map[string][]string)
		// for key, value := range tableH {
		// 	fmt.Printf("TYPE FOR VALUE: %T\n", value)
		// 	convertedRows := make([]string, 0) // Initialize an empty string slice

		// 	for val := range value {
		// 		tinyVal, ok := val.(string)
		// 		if !ok {
		// 			// Handle error or perform additional type assertions if necessary
		// 			fmt.Println("Error: Tiny tableH mismatch-i")
		// 			return
		// 		}
		// 		convertedRows = append(convertedRows, tinyVal)
		// 	}
		// 	convertedEntries[key] = convertedRows
		// }
		entries, ok := request["tableH"].(map[string]interface{})
		if !ok {
			fmt.Println("Error: tableH field is not a map[string]interface{}")
			return
		}

		convertedEntries := make(map[string][]string)
		for key, value := range entries {
			// Assert 'value' to []interface{}
			slice, ok := value.([]interface{})
			if !ok {
				fmt.Println("Error: value is not a slice")
				continue
			}

			// Now you can iterate over 'slice'
			for _, val := range slice {
				strVal, ok := val.(string)
				if !ok {
					fmt.Println("Error: value in the slice is not a string")
					continue
				}
				convertedEntries[key] = append(convertedEntries[key], strVal)
			}
		}

		// Now, 'convertedEntries' is a map[string][]string
		// Now, convertedEntries is a map[string][]string that you can use to update TableH

		fmt.Println("Received tableH update (INITIALIZATION):", convertedEntries)
		myNode.TableH.Entries = convertedEntries
		myNode.TableH.AddEntry(myNode.IP, "")
		myNode.Broadcast()
	}

	if requestType == "update" {
		fmt.Printf("TYPE FOR TABLEH: %T\n", request["tableH"])
		// Type assertion for the "tableH" field
		entries, ok := request["tableH"].(map[string]interface{})
		if !ok {
			fmt.Println("Error: tableH field is not a map[string]interface{}")
			return
		}

		convertedEntries := make(map[string][]string)
		for key, value := range entries {
			// Assert 'value' to []interface{}
			slice, ok := value.([]interface{})
			if !ok {
				fmt.Println("Error: value is not a slice")
				continue
			}

			// Now you can iterate over 'slice'
			for _, val := range slice {
				strVal, ok := val.(string)
				if !ok {
					fmt.Println("Error: value in the slice is not a string")
					continue
				}
				convertedEntries[key] = append(convertedEntries[key], strVal)
			}
		}

		// Now, 'convertedEntries' is a map[string][]string
		// Now, convertedEntries is a map[string][]string that you can use to update TableH

		fmt.Println("Received tableH update (UPDATE):", convertedEntries)
		myNode.TableH.Entries = convertedEntries
	}

	if requestType == "download" {
		filename, ok := request["filename"].(string)
		if !ok {
			fmt.Println("Error: filename field is not the expected type")
			return
		}

		filedir := "file/"
		filePath := filepath.Join(filedir, filename)

		fmt.Println("FILE PATH: ", filePath)
		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Fprintf(conn, "File does not exist: %s\n", filePath)
			return
		}

		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(conn, "Failed to open file: %s\n", err)
			return
		}
		defer file.Close()

		fmt.Println("FILE OPENED")
		// Send the file over the TCP connection
		if _, err := io.Copy(conn, file); err != nil {
			fmt.Fprintf(conn, "Failed to send file: %s\n", err)
		}

		fmt.Println("FILE SENT")
	}
}

func MakeRequest(n *Node, input_req string) error {
	//read from console. From https://freshman.tech/snippets/go/read-console-input/

	switch input_req {
	case "download", "1":
		fmt.Print("Enter the file name you're requesting (or enter 'x' to cancel request): \n")
		// read the file name from user input
		input_file, err := ReadFromConsole()
		if err != nil {
			return err
		}
		if input_file == "x" {
			fmt.Println("Request Canceled.")
			fmt.Print("Choose your request: download(1), update(2), leave(3), or cancel request(x)\n")
		}

		peerIP := n.TableH.GetNodesWithFile(input_file)
		if len(peerIP) == 0 {
			return errors.New("file Does Not Exist")
		}

		conn, err := net.Dial("tcp", peerIP[0])
		if err != nil {
			return fmt.Errorf("failed to connect to server: %s", err)
		}
		defer conn.Close()

		requestData := map[string]string{
			"type":     "download",
			"filename": input_file,
		}

		// Marshal the request data into JSON
		jsonData, err := json.Marshal(requestData)
		if err != nil {
			return err
		}

		// Send the request
		_, err = conn.Write(jsonData)
		if err != nil {
			return err
		}

		// Create or truncate the local file
		file, err := os.Create("file_received/" + input_file)
		if err != nil {
			return fmt.Errorf("failed to create file: %s", err)
		}
		fmt.Printf("FILE CREATED")
		defer file.Close()

		// Read the response (file data) and write it to the file
		_, err = io.Copy(file, conn)
		if err != nil {
			return fmt.Errorf("failed to receive and write file: %s", err)
		}
		fmt.Printf("FILE COPIED")

	case "update", "2":
		fmt.Print("Enter the action you want to make: add, delete, or cancel request(x)\n")
		input_action, err := ReadFromConsole()
		if err != nil {
			return err
		}
		if input_action == "x" {
			fmt.Println("Request Canceled.")
			fmt.Print("Choose your request: download(1), update(2), leave(3), or cancel request(x)\n")

			return errors.New("request Canceled")
		}
		fmt.Print("Enter the file you want to update (enter 'x' to cancel request): \n")
		input_file, err := ReadFromConsole()
		if err != nil {
			return err
		}
		if input_file == "x" {
			fmt.Println("Request Canceled.")
			fmt.Print("Choose your request: download(1), update(2), leave(3), or cancel request(x)\n")

			return errors.New("request Canceled")
		}
		err = CheckFile(n, input_file)
		if err != nil {
			fmt.Println("Check file failed.")
			return errors.New("file not found")
		}
		if input_action == "add" {
			n.TableH.AddEntry(n.IP, input_file)
			err := n.Broadcast()
			if err != nil {
				return errors.New("Broadcast failed. ")
			}
		}
		if input_action == "delete" {
			n.TableH.RemoveEntry(n.IP, input_file)
			err := n.Broadcast()
			if err != nil {
				return errors.New("Broadcast failed. ")
			}
		}
		return nil

	case "leave", "3":
		n.TableH.RemoveNode(n.IP)
		err := n.Broadcast()
		if err != nil {
			return errors.New("Broadcast failed. ")
		}
	case "x":
		fmt.Println("Request Canceled.")
		fmt.Print("Choose your request: download(1), update(2), leave(3), or cancel request(x)\n")
		return errors.New("request canceled")

	default:
		fmt.Println("Invalid input. Please try again.")
		fmt.Print("Choose your request: download(1), update(2), leave(3), or cancel request(x)\n")
		return errors.New("invalid input")
	}
	return nil
}

func ReadFromConsole() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return "", err
	}
	// remove the delimeter from the string
	input = strings.TrimSuffix(strings.ToLower(input), "\n")
	return input, nil
}

func CheckFile(n *Node, file_name string) error {
	file_path := "file/" + file_name

	if _, err := os.Stat(file_path); err != nil {
		fmt.Println("ERROR checking file:", err)
		return err
	}
	return nil
}
