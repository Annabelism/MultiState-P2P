package network

import (
    "fmt"
)

// TableH represents a table mapping from an "IP:port" string to a list of filenames.
type TableH struct {
    entries map[string][]string
}

// NewTableH creates a new table.
func NewTableH() *TableH {
    return &TableH{
        entries: make(map[string][]string),
    }
}

// AddEntry adds a new filename to the list of files hosted by a specific node (identified by IP:port).
func (t *TableH) AddEntry(nodeAddr string, fileName string) {
    t.entries[nodeAddr] = append(t.entries[nodeAddr], fileName)
}

// RemoveEntry removes a filename from a specific node's entry (IP:port).
func (t *TableH) RemoveEntry(nodeAddr string, fileName string) {
    if files, ok := t.entries[nodeAddr]; ok {
        for i, f := range files {
            if f == fileName {
                t.entries[nodeAddr] = append(files[:i], files[i+1:]...)
                break
            }
        }
        // If the node no longer hosts any files, delete the node entry
        if len(t.entries[nodeAddr]) == 0 {
            delete(t.entries, nodeAddr)
        }
    }
}

// RemoveNode removes all entries for a given node (IP:port).
func (t *TableH) RemoveNode(nodeAddr string) {
    delete(t.entries, nodeAddr)
}

// GetFilesByNode returns a list of filenames hosted by a specific node (IP:port).
func (t *TableH) GetFilesByNode(nodeAddr string) []string {
    return t.entries[nodeAddr]
}

// GetAllNodes returns a list of all unique node addresses (IP:port) in the TableH.
func (t *TableH) GetAllNodes() []string {
    var nodes []string
    for nodeAddr := range t.entries {
        nodes = append(nodes, nodeAddr)
    }
    return nodes
}

// GetNodesWithFile returns a list of nodes (IP:port) that have the specified file.
func (t *TableH) GetNodesWithFile(fileName string) []string {
    var nodesWithFile []string
    for nodeAddr, files := range t.entries {
        for _, f := range files {
            if f == fileName {
                nodesWithFile = append(nodesWithFile, nodeAddr)
                break // Once found, no need to check other files for the same node
            }
        }
    }
    return nodesWithFile
}

func PrintTableH(table *TableH) {
    for key, values := range table.entries {
        fmt.Printf("%s: %v\n", key, values)
    }
}