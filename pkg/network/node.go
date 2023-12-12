package network

import (
	"sync"
)

// Node represents a single node in the P2P network.
type Node struct {
	IPAddress     string            // IP address of the node
	SharedFiles   []string          // List of files the node is sharing
	Buffer        []NetworkRequest  // Buffer to store incoming network requests
	TableH        map[string][]string // 'Who-Has-What' table (Table H)
	TableHLock    sync.RWMutex      // Read/Write mutex for thread-safe access to Table H
}

// NetworkRequest represents a network request.
type NetworkRequest struct {
	// Define the structure of a network request here.
	// This could include fields like request type, payload, sender IP, etc.
}