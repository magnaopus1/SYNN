package network

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	"synnergy_network/pkg/common"
)

// NewTLSHandshakeManager initializes a new TLS handshake manager
func NewTLSHandshakeManager(certFile, keyFile, caFile string) *TLSHandshakeManager {
	return &TLSHandshakeManager{
		CertFile:    certFile,
		KeyFile:     keyFile,
		CAFile:      caFile,
		Connections: make(map[string]net.Conn),
	}
}

// InitiateTLSHandshake starts a TLS handshake with a remote node for secure communication
func (manager *TLSHandshakeManager) InitiateTLSHandshake(address string) (net.Conn, error) {
	// Load CA certificate
	caCert, err := ioutil.ReadFile(manager.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Load client certificates
	cert, err := tls.LoadX509KeyPair(manager.CertFile, manager.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate and key: %v", err)
	}

	// Configure TLS settings
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true, // For demo purposes, skip certificate validation (not recommended for production)
	}

	// Establish a secure connection using TLS
	conn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate TLS connection: %v", err)
	}

	fmt.Printf("TLS handshake initiated with node at %s\n", address)

	// Store the connection in the active connection pool
	manager.lock.Lock()
	manager.Connections[address] = conn
	manager.lock.Unlock()

	return conn, nil
}

// HandleIncomingTLSHandshake handles incoming TLS handshakes from remote nodes
func (manager *TLSHandshakeManager) HandleIncomingTLSHandshake(listener net.Listener) error {
	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		// Perform the TLS handshake
		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			fmt.Println("Failed to cast incoming connection to TLS")
			conn.Close()
			continue
		}

		err = tlsConn.Handshake()
		if err != nil {
			fmt.Printf("TLS handshake failed: %v\n", err)
			tlsConn.Close()
			continue
		}

		clientAddr := tlsConn.RemoteAddr().String()
		fmt.Printf("Secure TLS handshake completed with %s\n", clientAddr)

		// Add the connection to the active connection pool
		manager.lock.Lock()
		manager.Connections[clientAddr] = tlsConn
		manager.lock.Unlock()

		// Handle encrypted communication with the client
		go manager.handleTLSConnection(tlsConn)
	}
}

// handleTLSConnection handles incoming encrypted messages after the handshake
func (manager *TLSHandshakeManager) handleTLSConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	buffer := make([]byte, 4096)

	// Create an instance of Encryption
	encryption := &common.Encryption{}

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Failed to read from %s: %v\n", clientAddr, err)
			break
		}

		encryptedData := buffer[:n]

		// Decrypt the received data using the EncryptionKey from common package
		decryptedData, err := encryption.DecryptData(encryptedData, common.EncryptionKey)
		if err != nil {
			fmt.Printf("Failed to decrypt message from %s: %v\n", clientAddr, err)
			continue
		}

		// Process the decrypted message
		manager.processMessage(decryptedData, conn)
	}
}



// processMessage handles the incoming decrypted message and triggers appropriate actions
func (manager *TLSHandshakeManager) processMessage(message []byte, conn net.Conn) {
	var request NetworkRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		fmt.Printf("Failed to unmarshal request: %v\n", err)
		return
	}

	switch request.Type {
	case "SUBMIT_TRANSACTION":
		// Handle transaction submission (integrating with the ledger)
		manager.handleTransactionSubmission(request, conn)
	default:
		fmt.Printf("Unknown request type: %s\n", request.Type)
	}
}

// handleTransactionSubmission processes a transaction submission and stores it in the ledger
func (manager *TLSHandshakeManager) handleTransactionSubmission(request NetworkRequest, conn net.Conn) {
	var tx common.Transaction
	err := json.Unmarshal(request.Data, &tx)
	if err != nil {
		fmt.Printf("Failed to unmarshal transaction: %v\n", err)
		manager.sendErrorResponse(conn, "Invalid transaction format")
		return
	}

	// Assuming the Transaction struct contains FromAddress, ToAddress, and Amount fields
	sender := tx.FromAddress
	recipient := tx.ToAddress
	amount := tx.Amount

	// Store the transaction in the ledger with the correct arguments
	err = manager.Ledger.AddTransaction(sender, recipient, amount)
	if err != nil {
		fmt.Printf("Failed to add transaction to ledger: %v\n", err)
		manager.sendErrorResponse(conn, "Failed to add transaction to ledger")
		return
	}

	// Respond with success
	manager.sendSuccessResponse(conn, "Transaction successfully submitted and stored in ledger")
}



// sendErrorResponse sends an error response back to the client
func (manager *TLSHandshakeManager) sendErrorResponse(conn net.Conn, errorMsg string) {
	response := NetworkResponse{
		Type:    "ERROR",
		Message: errorMsg,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Failed to marshal error response: %v\n", err)
		return
	}

	// Create an instance of Encryption
	encryption := &common.Encryption{}

	// Encrypt the response before sending, using EncryptionKey from common package
	encryptedResponse, err := encryption.EncryptData("AES", responseData, common.EncryptionKey)
	if err != nil {
		fmt.Printf("Failed to encrypt error response: %v\n", err)
		return
	}

	_, err = conn.Write(encryptedResponse)
	if err != nil {
		fmt.Printf("Failed to send error response: %v\n", err)
	}
}



// sendSuccessResponse sends a success response back to the client
func (manager *TLSHandshakeManager) sendSuccessResponse(conn net.Conn, successMsg string) {
	response := NetworkResponse{
		Type:    "SUCCESS",
		Message: successMsg,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Failed to marshal success response: %v\n", err)
		return
	}

	// Create an instance of Encryption
	encryption := &common.Encryption{}

	// Encrypt the response before sending, using EncryptionKey from common package
	encryptedResponse, err := encryption.EncryptData("AES", responseData, common.EncryptionKey)
	if err != nil {
		fmt.Printf("Failed to encrypt success response: %v\n", err)
		return
	}

	_, err = conn.Write(encryptedResponse)
	if err != nil {
		fmt.Printf("Failed to send success response: %v\n", err)
	}
}



// CloseAllTLSConnections gracefully closes all active TLS connections
func (manager *TLSHandshakeManager) CloseAllTLSConnections() {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	for addr, conn := range manager.Connections {
		conn.Close()
		fmt.Printf("Closed TLS connection with %s\n", addr)
		delete(manager.Connections, addr)
	}
}
