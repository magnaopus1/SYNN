package network

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewNetworkManager initializes a new network manager
func NewNetworkManager(nodeAddress string, ledger *ledger.Ledger, maxIdleTime time.Duration) *NetworkManager {
	return &NetworkManager{
		nodeAddress:    nodeAddress,
		peers:          make(map[string]*PeerConnection),
		ledgerInstance: ledger,
		connectionPool: NewConnectionPool(maxIdleTime), // Pass maxIdleTime as an argument
	}
}


// ConnectToPeer establishes a connection with a peer and adds it to the peer list
func (nm *NetworkManager) ConnectToPeer(peerIP string) error {
    nm.lock.Lock()
    defer nm.lock.Unlock()

    if _, exists := nm.peers[peerIP]; exists {
        return fmt.Errorf("already connected to peer %s", peerIP)
    }

    // Remove the unused conn variable as it's not necessary
    _, err := net.Dial("tcp", peerIP)
    if err != nil {
        return fmt.Errorf("failed to connect to peer %s: %v", peerIP, err)
    }

    nm.peers[peerIP] = &PeerConnection{
        NodeID:        peerIP,            // Use NodeID instead of PublicIP
        Connection:    nm.connectionPool, // Use the connection pool
        IsAlive:       true,              // Use IsAlive instead of IsActive
        LastPingTime:  time.Now(),        // Use the correct field name LastPingTime
        EncryptionKey: []byte{},          // Initialize with an encryption key
    }

    // Log the connection event in the ledger
    nm.ledgerInstance.LogNetworkEvent(fmt.Sprintf("PeerConnected from %s to %s", nm.nodeAddress, peerIP))
    fmt.Printf("Connected to peer: %s\n", peerIP)
    return nil
}

// DisconnectFromPeer disconnects from a peer and removes it from the peer list
func (nm *NetworkManager) DisconnectFromPeer(peerIP string) error {
	nm.lock.Lock()
	defer nm.lock.Unlock()

	peer, exists := nm.peers[peerIP]
	if !exists {
		return fmt.Errorf("no active connection with peer %s", peerIP)
	}

	// Close all active connections in the connection pool
	nm.connectionPool.mutex.Lock()
	for _, conn := range peer.Connection.ActiveConns {
		if err := conn.Close(); err != nil {
			nm.connectionPool.mutex.Unlock()
			return fmt.Errorf("failed to close connection with peer %s: %v", peerIP, err)
		}
	}
	nm.connectionPool.mutex.Unlock()

	// Remove the peer from the list
	delete(nm.peers, peerIP)

	// Log the event with the correct number of arguments
	event := fmt.Sprintf("PeerDisconnected: %s from %s", peerIP, nm.nodeAddress)
	nm.ledgerInstance.LogNetworkEvent(event)
	
	fmt.Printf("Disconnected from peer: %s\n", peerIP)
	return nil
}


// SendEncryptedMessage sends an encrypted message to a connected peer
func (nm *NetworkManager) SendEncryptedMessage(peerIP, message string) error {
	nm.lock.Lock()
	defer nm.lock.Unlock()

	peer, exists := nm.peers[peerIP]
	if !exists {
		return fmt.Errorf("no active connection with peer %s", peerIP)
	}

	// Encrypt the message
	encryptedMessage, err := nm.encryptMessage(message)
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %v", err)
	}

	// Lock the connection pool and write the message to the first available connection
	nm.connectionPool.mutex.Lock()
	defer nm.connectionPool.mutex.Unlock()

	if len(peer.Connection.ActiveConns) == 0 {
		return fmt.Errorf("no active connections available for peer %s", peerIP)
	}

	// Use the first active connection for sending the message
	activeConn := peer.Connection.ActiveConns[0]
	_, err = activeConn.Write([]byte(encryptedMessage))
	if err != nil {
		return fmt.Errorf("failed to send message to peer %s: %v", peerIP, err)
	}

	// Log the event in the ledger
	event := fmt.Sprintf("MessageSent from %s to %s at %v", nm.nodeAddress, peerIP, time.Now())
	nm.ledgerInstance.LogNetworkEvent(event)

	fmt.Printf("Encrypted message sent to peer %s\n", peerIP)
	return nil
}


// ReceiveMessages listens for incoming messages from peers
func (nm *NetworkManager) ReceiveMessages(peerIP string) error {
	peer, exists := nm.peers[peerIP]
	if !exists {
		return fmt.Errorf("no active connection with peer %s", peerIP)
	}

	// Retrieve the first active connection from the connection pool
	nm.connectionPool.mutex.Lock()
	if len(peer.Connection.ActiveConns) == 0 {
		nm.connectionPool.mutex.Unlock()
		return fmt.Errorf("no active connections available for peer %s", peerIP)
	}
	activeConn := peer.Connection.ActiveConns[0]
	nm.connectionPool.mutex.Unlock()

	buf := make([]byte, 1024)
	for {
		n, err := activeConn.Read(buf)
		if err != nil {
			fmt.Printf("Error reading from peer %s: %v\n", peerIP, err)
			return err
		}

		encryptedMessage := string(buf[:n])
		decryptedMessage, err := nm.decryptMessage(encryptedMessage)
		if err != nil {
			fmt.Printf("Failed to decrypt message from peer %s: %v\n", peerIP, err)
			continue
		}

		// Log the message reception in the ledger
		event := fmt.Sprintf("MessageReceived from %s to %s at %v", peerIP, nm.nodeAddress, time.Now())
		nm.ledgerInstance.LogNetworkEvent(event)

		fmt.Printf("Decrypted message from peer %s: %s\n", peerIP, decryptedMessage)
	}
}

// PingPeer sends a ping to a peer to check if the connection is still alive
func (nm *NetworkManager) PingPeer(peerIP string) error {
	peer, exists := nm.peers[peerIP]
	if !exists {
		return fmt.Errorf("no active connection with peer %s", peerIP)
	}

	// Retrieve the first active connection from the connection pool
	nm.connectionPool.mutex.Lock()
	if len(peer.Connection.ActiveConns) == 0 {
		nm.connectionPool.mutex.Unlock()
		peer.IsAlive = false
		return fmt.Errorf("no active connections available for peer %s", peerIP)
	}
	activeConn := peer.Connection.ActiveConns[0]
	nm.connectionPool.mutex.Unlock()

	_, err := activeConn.Write([]byte("ping"))
	if err != nil {
		peer.IsAlive = false
		return fmt.Errorf("failed to ping peer %s: %v", peerIP, err)
	}

	peer.LastPingTime = time.Now()
	fmt.Printf("Pinged peer %s\n", peerIP)
	return nil
}


// encryptMessage encrypts a message using a public key (peer-specific)
func (nm *NetworkManager) encryptMessage(message string) (string, error) {
    encryption := &common.Encryption{} // Create an instance of the Encryption struct
    pubKey := GetNetworkPublicKey()     // Get the RSA public key

    // Convert the public key to bytes (if using RSA)
    pubKeyBytes := x509.MarshalPKCS1PublicKey(pubKey)

    // Encrypt using RSA and the public key
    encryptedContent, err := encryption.EncryptData("RSA", []byte(message), pubKeyBytes) // Using "RSA" as the encryption method
    if err != nil {
        return "", fmt.Errorf("encryption failed: %v", err)
    }
    return hex.EncodeToString(encryptedContent), nil
}

// decryptMessage decrypts an incoming encrypted message using the node's private key
func (nm *NetworkManager) decryptMessage(encryptedMessage string) (string, error) {
    encryption := &common.Encryption{} // Create an instance of the Encryption struct
    privKey := GetNodePrivateKey()     // Get the RSA private key

    contentBytes, err := hex.DecodeString(encryptedMessage)
    if err != nil {
        return "", fmt.Errorf("failed to decode encrypted message: %v", err)
    }

    // Convert the private key to bytes
    privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)

    // Decrypt the message using the private key (Remove "RSA" as it's unnecessary)
    decryptedContent, err := encryption.DecryptData(contentBytes, privKeyBytes) // Just two arguments: contentBytes and privKeyBytes
    if err != nil {
        return "", fmt.Errorf("decryption failed: %v", err)
    }

    return string(decryptedContent), nil
}


// LogConnection logs the connection event into the ledger
func (nm *NetworkManager) LogConnection(eventType, fromIP, toIP string) {
    logMessage := fmt.Sprintf("%s: from %s to %s at %s", eventType, fromIP, toIP, time.Now().Format(time.RFC3339))
    nm.ledgerInstance.LogNetworkEvent(logMessage)
}


// GenerateConnectionID generates a unique connection ID based on node addresses and timestamp
func GenerateConnectionID(fromIP, toIP string) string {
	hashInput := fmt.Sprintf("%s_%s_%d", fromIP, toIP, time.Now().UnixNano())
	hash := sha256.New()
	hash.Write([]byte(hashInput))
	return hex.EncodeToString(hash.Sum(nil))
}
