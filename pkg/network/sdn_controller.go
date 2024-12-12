package network

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewSDNController initializes the SDN Controller
func NewSDNController(encryptionKey string, ledgerInstance *ledger.Ledger) *SDNController {
	return &SDNController{
		Nodes:          make(map[string]*SDNNode),
		LedgerInstance: ledgerInstance,
		EncryptionKey:  encryptionKey,
	}
}

// AddNode adds a new node to the network
func (controller *SDNController) AddNode(nodeID string, address net.IP) {
	controller.NodeLock.Lock()
	defer controller.NodeLock.Unlock()

	// Add the node to the network
	controller.Nodes[nodeID] = &SDNNode{
		NodeID:    nodeID,
		Address:   address,
		Status:    "active",
		LastCheck: time.Now(),
	}

	fmt.Printf("Node %s added to the network at %s\n", nodeID, address)
}

// RemoveNode removes a node from the network
func (controller *SDNController) RemoveNode(nodeID string) {
	controller.NodeLock.Lock()
	defer controller.NodeLock.Unlock()

	// Remove the node from the network
	delete(controller.Nodes, nodeID)

	fmt.Printf("Node %s removed from the network\n", nodeID)
}

// SendControlMessage sends an encrypted control message to a node
func (controller *SDNController) SendControlMessage(nodeID, message string) error {
	controller.NodeLock.Lock()
	defer controller.NodeLock.Unlock()

	node, exists := controller.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}

	// Create an encryption instance
	encryption := &common.Encryption{}

	// Encrypt the control message (using AES as an example)
	encryptedMessage, err := encryption.EncryptData("AES", []byte(message), []byte(controller.EncryptionKey)) // Convert EncryptionKey to []byte
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %v", err)
	}

	// Simulate sending the message over the network
	fmt.Printf("Sending encrypted message to node %s at %s: %x\n", nodeID, node.Address.String(), encryptedMessage)
	return nil
}



// ReceiveControlMessage receives and decrypts a control message from a node
func (controller *SDNController) ReceiveControlMessage(nodeID, encryptedMessage string) (string, error) {
	controller.NodeLock.Lock()
	defer controller.NodeLock.Unlock()

	node, exists := controller.Nodes[nodeID]
	if !exists {
		return "", fmt.Errorf("node %s not found", nodeID)
	}

	// Use the node variable for logging or other operations (to avoid the unused variable error)
	fmt.Printf("Receiving message from node %s at address %s\n", nodeID, node.Address.String())

	// Create an encryption instance
	encryption := &common.Encryption{}

	// Decrypt the control message
	decryptedMessage, err := encryption.DecryptData([]byte(encryptedMessage), []byte(controller.EncryptionKey)) // Only two arguments now
	if err != nil {
		return "", fmt.Errorf("failed to decrypt message: %v", err)
	}

	fmt.Printf("Received decrypted message from node %s: %s\n", nodeID, string(decryptedMessage))
	return string(decryptedMessage), nil
}




// GetNodeStatus returns the status of a node in the network
func (controller *SDNController) GetNodeStatus(nodeID string) (string, error) {
	controller.NodeLock.Lock()
	defer controller.NodeLock.Unlock()

	node, exists := controller.Nodes[nodeID]
	if !exists {
		return "", fmt.Errorf("node %s not found", nodeID)
	}

	return node.Status, nil
}

// UpdateNodeStatus updates the status of a node (e.g., during a health check)
func (controller *SDNController) UpdateNodeStatus(nodeID string, status string) error {
	controller.NodeLock.Lock()
	defer controller.NodeLock.Unlock()

	node, exists := controller.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}

	node.Status = status
	node.LastCheck = time.Now()

	fmt.Printf("Updated status of node %s to %s\n", nodeID, status)
	return nil
}

// MonitorNodeHealth performs a periodic health check of all nodes
func (controller *SDNController) MonitorNodeHealth(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			controller.NodeLock.Lock()
			for nodeID, node := range controller.Nodes {
				// Check if the node has responded within the expected time window
				if time.Since(node.LastCheck) > interval*2 {
					node.Status = "inactive"
					fmt.Printf("Node %s marked as inactive\n", nodeID)
				} else {
					fmt.Printf("Node %s is healthy\n", nodeID)
				}
			}
			controller.NodeLock.Unlock()
		}
	}
}

// SecureNetworkBootstrap initializes the secure network bootstrap process
func (controller *SDNController) SecureNetworkBootstrap(certFile, keyFile string) error {
	// Load TLS certificate and key for secure communication
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificates: %v", err)
	}

	// Create a secure listener for handling node connections
	listener, err := tls.Listen("tcp", ":8080", &tls.Config{Certificates: []tls.Certificate{cert}})
	if err != nil {
		return fmt.Errorf("failed to start TLS listener: %v", err)
	}

	fmt.Println("SDN Controller is securely listening for connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go controller.handleNodeConnection(conn)
	}
}

// handleNodeConnection processes incoming encrypted connections from nodes
func (controller *SDNController) handleNodeConnection(conn net.Conn) {
	defer conn.Close()

	// Read the incoming data
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	// Create an encryption instance
	encryption := &common.Encryption{}

	// Convert controller.EncryptionKey to []byte if it's a string
	encryptionKey := []byte(controller.EncryptionKey)

	// Decrypt the message received from the node
	decryptedMessage, err := encryption.DecryptData(buffer[:n], encryptionKey)
	if err != nil {
		fmt.Printf("Error decrypting node message: %v\n", err)
		return
	}

	fmt.Printf("Received secure message: %s\n", string(decryptedMessage))
}

