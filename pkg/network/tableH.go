package network

// TableH represents the 'Who-Has-What' table in the network.
// It maps file names to a list of IPs that have the file.
type TableH struct {
    entries map[string][]string
}

// NewTableH creates a new 'Who-Has-What' table.
func NewTableH() *TableH {
    return &TableH{
        entries: make(map[string][]string),
    }
}

// AddEntry adds a new file and its location to the table.
func (t *TableH) AddEntry(fileName, nodeIP string) {
    t.entries[fileName] = append(t.entries[fileName], nodeIP)
}

// RemoveEntry removes a file entry from a specific node.
func (t *TableH) RemoveEntry(fileName, nodeIP string) {
    if nodes, ok := t.entries[fileName]; ok {
        for i, ip := range nodes {
            if ip == nodeIP {
                t.entries[fileName] = append(nodes[:i], nodes[i+1:]...)
                break
            }
        }
        // If the file is no longer available at any node, delete the file entry
        if len(t.entries[fileName]) == 0 {
            delete(t.entries, fileName)
        }
    }
}

// RemoveNode removes all entries for a given node IP.
func (t *TableH) RemoveNode(nodeIP string) {
    for fileName, nodes := range t.entries {
        for i, ip := range nodes {
            if ip == nodeIP {
                t.entries[fileName] = append(nodes[:i], nodes[i+1:]...)
            }
        }
        // Clean up any files that no longer have any nodes hosting them
        if len(t.entries[fileName]) == 0 {
            delete(t.entries, fileName)
        }
    }
}

// GetNodesWithFile returns a list of nodes that have the specified file.
func (t *TableH) GetNodesWithFile(fileName string) []string {
    return t.entries[fileName]
}
