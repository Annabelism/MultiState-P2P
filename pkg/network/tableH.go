package network

import (
	"encoding/json"
	"net"
)

type TableH struct {
	Entries map[string][]string
}

// Share sends the TableH data structure over a network connection
func ShareTableH(table *TableH, peerIP string, initial bool) error { // Encode the TableH structure to JSON
	// Establish a TCP connection
	conn, err := net.Dial("tcp", peerIP)
	if err != nil {
		return err
	}

	// Construct the JSON data using a map
	updateMsg := map[string]interface{}{
		"type":   "update",
		"tableH": table.Entries,
	}

	if initial {
		updateMsg["type"] = "update-i"
	}

	jsonData, err := json.Marshal(updateMsg)
	if err != nil {
		conn.Close()
		return err
	}

	// Send the JSON data
	_, err = conn.Write(jsonData)
	conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// NewTableH creates a new table.
func NewTableH() *TableH {
	return &TableH{
		Entries: make(map[string][]string),
	}
}

// AddEntry adds a new filename to the list of files hosted by a specific node (identified by IP:port).
func (t *TableH) AddEntry(nodeAddr string, fileName string) {
	t.Entries[nodeAddr] = append(t.Entries[nodeAddr], fileName)
}

// RemoveEntry removes a filename from a specific node's entry (IP:port).
func (t *TableH) RemoveEntry(nodeAddr string, fileName string) {
	if files, ok := t.Entries[nodeAddr]; ok {
		for i, f := range files {
			if f == fileName {
				t.Entries[nodeAddr] = append(files[:i], files[i+1:]...)
				break
			}
		}
		// If the node no longer hosts any files, delete the node entry
		if len(t.Entries[nodeAddr]) == 0 {
			delete(t.Entries, nodeAddr)
		}
	}
}

func (t *TableH) AddNode(nodeAddr string) {
	t.Entries[nodeAddr] = make([]string, 0)
}

// RemoveNode removes all Entries for a given node (IP:port).
func (t *TableH) RemoveNode(nodeAddr string) {
	delete(t.Entries, nodeAddr)
}

// GetFilesByNode returns a list of filenames hosted by a specific node (IP:port).
func (t *TableH) GetFilesByNode(nodeAddr string) []string {
	return t.Entries[nodeAddr]
}

// GetAllNodes returns a list of all unique node addresses (IP:port) in the TableH.
func (t *TableH) GetAllNodes() []string {
	var nodes []string
	for nodeAddr := range t.Entries {
		nodes = append(nodes, nodeAddr)
	}
	return nodes
}

// GetNodesWithFile returns a list of nodes (IP:port) that have the specified file.
func (t *TableH) GetNodesWithFile(fileName string) []string {
	var nodesWithFile []string
	for nodeAddr, files := range t.Entries {
		for _, f := range files {
			if f == fileName {
				nodesWithFile = append(nodesWithFile, nodeAddr)
				break // Once found, no need to check other files for the same node
			}
		}
	}
	return nodesWithFile
}
