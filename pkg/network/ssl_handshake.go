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


// NewSSLHandshakeManager initializes a new SSL handshake manager
func NewSSLHandshakeManager(certFile, keyFile, caFile string) *SSLHandshakeManager {
	return &SSLHandshakeManager{
		CertFile:    certFile,
		KeyFile:     keyFile,
		CAFile:      caFile,
		Connections: make(map[string]net.Conn),
	}
}

// StartHandshake initializes an SSL handshake with a remote node for secure communication
func (manager *SSLHandshakeManager) StartHandshake(address string) (net.Conn, error) {
	// Load CA cert for client-side validation
	caCert, err := ioutil.ReadFile(manager.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA file: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Load the client certificates
	cert, err := tls.LoadX509KeyPair(manager.CertFile, manager.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %v", err)
	}

	// Configure TLS client settings
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// Dial the remote server using TLS
	conn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to establish TLS connection: %v", err)
	}

	fmt.Printf("TLS handshake completed with %s\n", address)

	// Add the connection to the active connection pool
	manager.lock.Lock()
	manager.Connections[address] = conn
	manager.lock.Unlock()

	return conn, nil
}

// HandleIncomingHandshake handles incoming secure TLS handshakes
func (manager *SSLHandshakeManager) HandleIncomingHandshake(listener net.Listener) error {
	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			fmt.Println("Failed to cast connection to TLS")
			conn.Close()
			continue
		}

		// Perform TLS handshake
		err = tlsConn.Handshake()
		if err != nil {
			fmt.Printf("TLS handshake failed: %v\n", err)
			tlsConn.Close()
			continue
		}

		clientAddr := tlsConn.RemoteAddr().String()
		fmt.Printf("TLS handshake successful with %s\n", clientAddr)

		// Add the connection to the active pool
		manager.lock.Lock()
		manager.Connections[clientAddr] = tlsConn
		manager.lock.Unlock()

		// Start handling incoming data after handshake
		go manager.handleConnection(tlsConn)
	}
}

// handleConnection handles incoming messages after SSL handshake
func (manager *SSLHandshakeManager) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	buffer := make([]byte, 4096)

	// Create an encryption instance
	encryption := &common.Encryption{}

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Failed to read from %s: %v\n", clientAddr, err)
			break
		}

		encryptedMessage := buffer[:n]

		// Decrypt the message (without specifying the algorithm, as AES is implied by the key)
		decryptedMessage, err := encryption.DecryptData(encryptedMessage, common.EncryptionKey) // Only two arguments
		if err != nil {
			fmt.Printf("Failed to decrypt message from %s: %v\n", clientAddr, err)
			continue
		}

		// Process the decrypted message
		manager.processMessage(decryptedMessage, conn)
	}
}



// processMessage processes the decrypted message and performs the necessary action
func (manager *SSLHandshakeManager) processMessage(message []byte, conn net.Conn) {
	var request NetworkRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		fmt.Printf("Invalid request format: %v\n", err)
		return
	}

	switch request.Type {
	case "PING":
		fmt.Println("Received PING request, sending PONG.")
		manager.sendResponse(conn, "PONG", "PONG message received.")
	default:
		fmt.Printf("Unknown request type: %s\n", request.Type)
	}
}

// sendResponse sends an encrypted response after processing the message
func (manager *SSLHandshakeManager) sendResponse(conn net.Conn, responseType, message string) {
	response := NetworkResponse{
		Type:    responseType,
		Message: message,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Failed to marshal response: %v\n", err)
		return
	}

	// Create an encryption instance
	encryption := &common.Encryption{}

	// Encrypt the response data
	encryptedResponse, err := encryption.EncryptData("AES", responseData, []byte(common.EncryptionKey))
	if err != nil {
		fmt.Printf("Failed to encrypt response: %v\n", err)
		return
	}

	// Send the encrypted response back to the client
	_, err = conn.Write(encryptedResponse)
	if err != nil {
		fmt.Printf("Failed to send response: %v\n", err)
	}
}


// CloseConnection closes an active connection and removes it from the pool
func (manager *SSLHandshakeManager) CloseConnection(address string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if conn, exists := manager.Connections[address]; exists {
		conn.Close()
		fmt.Printf("Closed connection with %s\n", address)
		delete(manager.Connections, address)
	}
}

// CloseAllConnections gracefully closes all active connections
func (manager *SSLHandshakeManager) CloseAllConnections() {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	for addr, conn := range manager.Connections {
		conn.Close()
		fmt.Printf("Closed connection with %s\n", addr)
		delete(manager.Connections, addr)
	}
}
