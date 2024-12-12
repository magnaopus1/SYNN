package network_test

import (
	"encoding/json"
	"net"
	"os"
	"sync"
	"testing"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)

// Mock implementations of common and ledger interfaces if required, e.g., mock SynnergyConsensus, Encryption, etc.
// Replace with actual implementations if available.

func TestNewServerInitialization(t *testing.T) {
	// Initialize the ledger instance
	ledgerInstance := ledger.NewLedger() // Use actual NewLedger function

	// Initialize consensus engine (using a simple mock if needed)
	consensusEngine := &common.SynnergyConsensus{}

	// Server address and encryption key
	address := "localhost:8000"
	encryptionKey := "testEncryptionKey"

	// Initialize server instance
	server := network.NewServer(address, encryptionKey, ledgerInstance, consensusEngine)
	if server == nil {
		t.Fatal("Failed to initialize server instance")
	}

	if server.Address != address {
		t.Errorf("Expected server address %s, got %s", address, server.Address)
	}
}

func TestServerStartStop(t *testing.T) {
	// Create a self-signed TLS certificate for testing purposes
	certFile := "test_cert.pem"
	keyFile := "test_key.pem"
	err := network.CreateSelfSignedCert(certFile, keyFile)
	if err != nil {
		t.Fatalf("Failed to create test certificate: %v", err)
	}
	defer os.Remove(certFile)
	defer os.Remove(keyFile)

	// Initialize the ledger instance and server
	ledgerInstance := ledger.NewLedger()
	consensusEngine := &common.SynnergyConsensus{}
	server := network.NewServer("localhost:8001", "testEncryptionKey", ledgerInstance, consensusEngine)

	// Run server in a separate goroutine to enable testing
	go func() {
		err := server.StartServer(certFile, keyFile)
		if err != nil {
			t.Errorf("Failed to start server: %v", err)
		}
	}()
	time.Sleep(1 * time.Second) // Allow time for server to start

	// Stop the server
	server.StopServer()
}

func TestHandleConnection(t *testing.T) {
	// Prepare test server and listener
	ledgerInstance := ledger.NewLedger()
	consensusEngine := &common.SynnergyConsensus{}
	server := network.NewServer("localhost:8002", "testEncryptionKey", ledgerInstance, consensusEngine)

	// Setup a connection to test handleConnection function
	listener, err := net.Listen("tcp", server.Address)
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Errorf("Failed to accept connection: %v", err)
		}
		defer conn.Close()

		// Test server handling connection
		server.HandleConnection(conn)
		wg.Done()
	}()

	// Simulate client connection
	clientConn, err := net.Dial("tcp", server.Address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer clientConn.Close()

	// Send test message to server
	_, err = clientConn.Write([]byte("Test message"))
	if err != nil {
		t.Errorf("Failed to write to server: %v", err)
	}

	wg.Wait() // Wait for server to process the connection
}

func TestProcessMessage(t *testing.T) {
	// Prepare server with ledger and consensus
	ledgerInstance := ledger.NewLedger()
	consensusEngine := &common.SynnergyConsensus{}
	server := network.NewServer("localhost:8003", "testEncryptionKey", ledgerInstance, consensusEngine)

	// Create test network request
	request := network.NetworkRequest{
		Type:    "SUB_BLOCK",
		Payload: json.RawMessage(`{"SubBlockID": "sub1", "Index": 1}`),
	}
	message, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Setup a fake connection to test message processing
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	go func() {
		err := server.ProcessMessage(message, serverConn)
		if err != nil {
			t.Errorf("Failed to process message: %v", err)
		}
	}()

	// Validate server response
	buffer := make([]byte, 1024)
	_, err = clientConn.Read(buffer)
	if err != nil {
		t.Errorf("Failed to read server response: %v", err)
	}
}

