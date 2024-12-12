package network

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewNATTraversalManager initializes a new NAT Traversal manager
func NewNATTraversalManager(publicIP, privateIP string, ledger *ledger.Ledger) *NATTraversalManager {
	return &NATTraversalManager{
		publicIP:       publicIP,
		privateIP:      privateIP,
		peerMap:        make(map[string]string),
		ledgerInstance: ledger,
		webrtcManager: NewWebRTCManager(),
	}
}

// RegisterPeer registers a peer's public-private IP mapping
func (nat *NATTraversalManager) RegisterPeer(publicIP, privateIP string) {
	nat.lock.Lock()
	defer nat.lock.Unlock()

	nat.peerMap[publicIP] = privateIP
	nat.ledgerInstance.LogNATEvent("PeerRegistered: " + publicIP)  // Corrected to a single argument
	fmt.Printf("Registered peer with public IP: %s, private IP: %s\n", publicIP, privateIP)
}

// UnregisterPeer removes a peer from the mapping
func (nat *NATTraversalManager) UnregisterPeer(publicIP string) {
	nat.lock.Lock()
	defer nat.lock.Unlock()

	delete(nat.peerMap, publicIP)
	nat.ledgerInstance.LogNATEvent("PeerUnregistered: " + publicIP)  // Corrected to a single argument
	fmt.Printf("Unregistered peer with public IP: %s\n", publicIP)
}

// SendMessageThroughNAT sends an encrypted message through NAT traversal
func (nat *NATTraversalManager) SendMessageThroughNAT(fromPublicIP, toPublicIP, message string) error {
	nat.lock.Lock()
	defer nat.lock.Unlock()

	privateIP, exists := nat.peerMap[toPublicIP]
	if !exists {
		return fmt.Errorf("peer with public IP %s not found", toPublicIP)
	}

	// Encrypt the message
	encryptedMessage, err := nat.encryptMessage(message)
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %v", err)
	}

	// Send the encrypted message to the peer
	if err := nat.sendEncryptedMessage(fromPublicIP, privateIP, encryptedMessage); err != nil {
		return err
	}

	// Log the NAT event (corrected to a single argument)
	nat.ledgerInstance.LogNATEvent("MessageSentThroughNAT: " + toPublicIP)
	fmt.Printf("Encrypted message sent from %s to %s\n", fromPublicIP, toPublicIP)
	return nil
}

// encryptMessage encrypts a message using the RSA public key of the recipient
func (nat *NATTraversalManager) encryptMessage(message string) (string, error) {
	// Fetch the public key for encryption (ensure the function returns *rsa.PublicKey)
	pubKey := GetNetworkPublicKey()

	// Hash the message with SHA-256 (optional but recommended)
	label := []byte("") // Optional label, can be left empty
	hash := sha256.New()

	// Encrypt the message using RSA-OAEP with the public key
	encryptedContent, err := rsa.EncryptOAEP(hash, rand.Reader, pubKey, []byte(message), label)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %v", err)
	}

	// Return the encrypted content as a hex string
	return hex.EncodeToString(encryptedContent), nil
}

// sendEncryptedMessage sends an encrypted message to the recipient via NAT traversal using WebRTC.
func (nat *NATTraversalManager) sendEncryptedMessage(fromPublicIP, toPrivateIP, message string) error {
	// Retrieve the WebRTC connection to the recipient
	peerConn, exists := nat.webrtcManager.Peers[toPrivateIP] // Use the private IP to identify the peer
	if !exists || !peerConn.IsAlive { // Corrected to use IsAlive
		return fmt.Errorf("peer at private IP %s is not connected or inactive", toPrivateIP)
	}

	// Create an encryption instance
	encryption := &common.Encryption{}

	// Encrypt the message before sending
	encryptedMessage, err := encryption.EncryptData("AES", []byte(message), peerConn.EncryptionKey) // Use EncryptionKey from peerConn
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %v", err)
	}

	// Retrieve an active connection from connectionPool
	nat.connectionPool.mutex.Lock() // Reference nat.connectionPool (lowercase 'c')
	defer nat.connectionPool.mutex.Unlock()

	// Assuming we are using the first available active connection for this example
	if len(peerConn.Connection.ActiveConns) == 0 {
		return fmt.Errorf("no active connections available for peer at %s", toPrivateIP)
	}
	activeConn := peerConn.Connection.ActiveConns[0]

	// Use the connection's Write method to send the encrypted message
	_, err = activeConn.Write(encryptedMessage) // Use Write method from net.Conn
	if err != nil {
		return fmt.Errorf("failed to send encrypted message to peer at %s: %v", toPrivateIP, err)
	}

	fmt.Printf("Encrypted message sent from public IP %s to private IP %s\n", fromPublicIP, toPrivateIP)
	return nil
}

// HandleIncomingConnection handles an incoming connection from a peer through NAT
func (nat *NATTraversalManager) HandleIncomingConnection(conn net.Conn) {
	defer conn.Close()

	// Read the incoming message
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to read message:", err)
		return
	}
	message := string(buf[:n])

	// Verify and decrypt the message
	decryptedMessage, err := nat.decryptMessage(message)
	if err != nil {
		fmt.Println("Failed to decrypt message:", err)
		return
	}

	fmt.Printf("Received decrypted message: %s\n", decryptedMessage)
}

// decryptMessage decrypts an encrypted message using the node's private key (RSA decryption)
func (nat *NATTraversalManager) decryptMessage(encryptedMessage string) (string, error) {
    // Fetch the node's private key
    privKey := GetNodePrivateKey()

    // Decode the hex-encoded encrypted message
    contentBytes, err := hex.DecodeString(encryptedMessage)
    if err != nil {
        return "", fmt.Errorf("failed to decode encrypted message: %v", err)
    }

    // Decrypt the content using RSA (PKCS1v15)
    decryptedContent, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, contentBytes)
    if err != nil {
        return "", fmt.Errorf("decryption failed: %v", err)
    }

    // Return the decrypted content as a string
    return string(decryptedContent), nil
}


// LogConnection logs the NAT traversal event into the ledger
func (nat *NATTraversalManager) LogConnectionEvent(eventType, fromIP, toIP string) {
	// Create a formatted log message for the event
	eventMessage := fmt.Sprintf("EventType: %s, From: %s, To: %s, Time: %v", eventType, fromIP, toIP, time.Now())

	// Log the message using the ledger's LogNATEvent method
	nat.ledgerInstance.LogNATEvent(eventMessage)

	fmt.Printf("NAT event logged: %s, from %s to %s\n", eventType, fromIP, toIP)
}

