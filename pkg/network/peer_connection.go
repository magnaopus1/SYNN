package network

import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// PeerConnectionManager manages all active peer connections in the network
type PeerConnectionManager struct {
	Connections map[string]*PeerConnection // Active peer connections
	Mutex       sync.Mutex                 // Thread-safe access to peer connections
	Ledger      *ledger.Ledger             // Pointer to the ledger for recording transactions
}

// PeerConnection represents a connection to a peer node, using a ConnectionPool
type PeerConnection struct {
    NodeID       string
    Connection   *ConnectionPool  // Custom connection pool instead of net.Conn
    LastPingTime time.Time        // Last time the peer was pinged
    IsAlive      bool             // Status of the connection
    EncryptionKey []byte          // Encryption key used for secure communication
}


// NewPeerConnectionManager initializes a new peer connection manager
func NewPeerConnectionManager(ledger *ledger.Ledger) *PeerConnectionManager {
	return &PeerConnectionManager{
		Connections: make(map[string]*PeerConnection),
		Ledger:      ledger,
	}
}

// ConnectToPeer establishes a connection to a peer using a ConnectionPool
func (pcm *PeerConnectionManager) ConnectToPeer(address, publicKey string) error {
    pcm.Mutex.Lock()
    defer pcm.Mutex.Unlock()

    // Step 1: Establish a TCP connection to the peer
    conn, err := net.Dial("tcp", address)
    if err != nil {
        return fmt.Errorf("failed to connect to peer %s: %v", address, err)
    }

    // Step 2: Add the TCP connection to the ConnectionPool's ActiveConns
    connectionPool := &ConnectionPool{
        ActiveConns: []net.Conn{conn},  // Add the new connection to the ActiveConns list
    }

    // Step 3: Perform key exchange using the net.Conn from the pool
    encryptionKey, err := pcm.performKeyExchange(conn, publicKey)  // Use conn directly
    if err != nil {
        conn.Close()  // Close the connection if the key exchange fails
        return fmt.Errorf("key exchange failed with peer %s: %v", address, err)
    }

    // Step 4: Create a new peer connection with the ConnectionPool
    peerConn := &PeerConnection{
        NodeID:       address,            // Use the peer's address as the NodeID
        Connection:   connectionPool,     // Store the ConnectionPool
        LastPingTime: time.Now(),         // Set the current time for LastPingTime
        IsAlive:      true,               // Mark the connection as alive
        EncryptionKey: encryptionKey,     // Store the encryption key
    }

    // Step 5: Store the connection in the PeerConnectionManager's connections map
    pcm.Connections[address] = peerConn
    fmt.Printf("Connected to peer %s\n", address)

    return nil
}


// DisconnectFromPeer closes the connection to a peer
func (pcm *PeerConnectionManager) DisconnectFromPeer(address string) error {
	pcm.Mutex.Lock()
	defer pcm.Mutex.Unlock()

	peerConn, exists := pcm.Connections[address]
	if !exists {
		return fmt.Errorf("peer %s not found", address)
	}

	// Close all active connections in the ConnectionPool
	err := peerConn.Connection.Close()
	if err != nil {
		return fmt.Errorf("failed to disconnect from peer %s: %v", address, err)
	}

	delete(pcm.Connections, address)
	fmt.Printf("Disconnected from peer %s\n", address)
	return nil
}






// performKeyExchange handles the Diffie-Hellman key exchange to establish a shared encryption key
func (pcm *PeerConnectionManager) performKeyExchange(conn net.Conn, peerPublicKey string) ([]byte, error) {
	// Generate a local private/public key pair
	localPrivateKey, localPublicKey, err := GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %v", err)
	}

	// Send local public key to the peer
	_, err = conn.Write(localPublicKey) // no need to convert to []byte since it already is
	if err != nil {
		return nil, fmt.Errorf("failed to send public key: %v", err)
	}

	// Receive peer's public key
	peerKeyBuffer := make([]byte, len(peerPublicKey))
	_, err = conn.Read(peerKeyBuffer)
	if err != nil {
		return nil, fmt.Errorf("failed to receive peer public key: %v", err)
	}

	// Derive shared encryption key
	sharedKey, err := DeriveSharedKey(localPrivateKey, peerKeyBuffer) // Pass peerKeyBuffer directly as []byte
	if err != nil {
		return nil, fmt.Errorf("failed to derive shared key: %v", err)
	}

	return sharedKey, nil
}


// GenerateKeyPair generates a private and public key pair using elliptic curve cryptography
func GenerateKeyPair() ([]byte, []byte, error) {
	priv, x, y, err := elliptic.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("key generation failed: %v", err)
	}
	pub := elliptic.Marshal(elliptic.P256(), x, y)
	return priv, pub, nil
}


// DeriveSharedKey derives a shared encryption key using the private and peer public key
func DeriveSharedKey(priv []byte, peerPubKey []byte) ([]byte, error) {
	x, y := elliptic.Unmarshal(elliptic.P256(), peerPubKey)
	if x == nil || y == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	sharedX, _ := elliptic.P256().ScalarMult(x, y, priv)
	return sharedX.Bytes(), nil
}


// SendMessage encrypts and sends a message to a connected peer
func (pcm *PeerConnectionManager) SendMessage(address string, message string) error {
	pcm.Mutex.Lock()
	defer pcm.Mutex.Unlock()

	peerConn, exists := pcm.Connections[address]
	if !exists {
		return fmt.Errorf("peer %s not found", address)
	}

	// Create an encryption instance and encrypt the message using the shared encryption key
	encryption := &common.Encryption{}
	encryptedMessage, err := encryption.EncryptData("AES", []byte(message), peerConn.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt message for peer %s: %v", address, err)
	}

	// Use an active connection to send the encrypted message
	if len(peerConn.Connection.ActiveConns) == 0 {
		return fmt.Errorf("no active connections for peer %s", address)
	}
	activeConn := peerConn.Connection.ActiveConns[0]

	_, err = activeConn.Write(encryptedMessage)
	if err != nil {
		return fmt.Errorf("failed to send message to peer %s: %v", address, err)
	}

	peerConn.LastPingTime = time.Now() // Update the LastPingTime instead of LastActive
	fmt.Printf("Sent encrypted message to peer %s\n", address)
	return nil
}

// ReceiveMessage receives and decrypts a message from a peer
func (pcm *PeerConnectionManager) ReceiveMessage(address string) (string, error) {
	pcm.Mutex.Lock()
	defer pcm.Mutex.Unlock()

	peerConn, exists := pcm.Connections[address]
	if !exists {
		return "", fmt.Errorf("peer %s not found", address)
	}

	// Use an active connection to receive the encrypted message
	if len(peerConn.Connection.ActiveConns) == 0 {
		return "", fmt.Errorf("no active connections for peer %s", address)
	}
	activeConn := peerConn.Connection.ActiveConns[0]

	messageBuffer := make([]byte, 1024)
	n, err := activeConn.Read(messageBuffer)
	if err != nil {
		return "", fmt.Errorf("failed to receive message from peer %s: %v", address, err)
	}

	// Decrypt the message using the shared encryption key
	encryption := &common.Encryption{}
	decryptedMessage, err := encryption.DecryptData(messageBuffer[:n], peerConn.EncryptionKey) // Corrected to use only two arguments
	if err != nil {
		return "", fmt.Errorf("failed to decrypt message from peer %s: %v", address, err)
	}

	peerConn.LastPingTime = time.Now() // Update the LastPingTime instead of LastActive
	return string(decryptedMessage), nil
}



// HandleIncomingConnections listens for and handles incoming peer connections
func (pcm *PeerConnectionManager) HandleIncomingConnections(listener net.Listener) {
	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept incoming connection: %v\n", err)
			continue
		}

		// Handle the incoming peer connection in a separate goroutine
		go pcm.handleNewPeerConnection(conn)
	}
}

// handleNewPeerConnection handles the setup of a new incoming peer connection
func (pcm *PeerConnectionManager) handleNewPeerConnection(conn net.Conn) {
	// Perform key exchange and establish encryption
	peerAddress := conn.RemoteAddr().String()
	peerPublicKeyBuffer := make([]byte, 1024)

	_, err := conn.Read(peerPublicKeyBuffer)
	if err != nil {
		fmt.Printf("Failed to read public key from peer %s: %v\n", peerAddress, err)
		conn.Close()
		return
	}

	// Clean up the peerPublicKeyBuffer to ensure no trailing bytes are included
	peerPublicKey := string(bytes.Trim(peerPublicKeyBuffer, "\x00"))

	// Perform key exchange to establish a shared encryption key
	encryptionKey, err := pcm.performKeyExchange(conn, peerPublicKey)
	if err != nil {
		fmt.Printf("Failed to perform key exchange with peer %s: %v\n", peerAddress, err)
		conn.Close()
		return
	}

	// Add the connection to the ConnectionPool
	pool := &ConnectionPool{} // Assuming you have a proper initialization for ConnectionPool
	pool.AddConnection(conn)  // Add net.Conn to your custom connection pool

	// Add the new peer connection
	pcm.Mutex.Lock()
	pcm.Connections[peerAddress] = &PeerConnection{
		NodeID:       peerAddress,    // Use NodeID to store peer's address
		Connection:   pool,           // Assign the ConnectionPool instead of net.Conn
		LastPingTime: time.Now(),
		IsAlive:      true,
		EncryptionKey: encryptionKey, // Store the encryption key established in key exchange
	}
	pcm.Mutex.Unlock()

	fmt.Printf("New peer connected: %s\n", peerAddress)
}
