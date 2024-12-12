package network

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewWebRTCManager initializes a new WebRTCManager
func NewWebRTCManager() *WebRTCManager {
	return &WebRTCManager{
		Peers: make(map[string]*PeerConnection),
	}
}

// AddPeerConnection adds a new peer connection to the WebRTCManager
func (wm *WebRTCManager) AddPeerConnection(peerID string, conn WebRTCConn) error {
    wm.lock.Lock()
    defer wm.lock.Unlock()

    if _, exists := wm.Peers[peerID]; exists {
        return fmt.Errorf("Peer %s is already connected", peerID)
    }

    // Use WebRTCConn as it is and assign it to the custom ConnectionPool
    connectionPool := &ConnectionPool{
        // WebRTCConn doesn't implement net.Conn, so we can't treat it as net.Conn directly
        ActiveConns: []net.Conn{}, // Assuming WebRTCConn isn't treated as net.Conn
    }

    peerConn := &PeerConnection{
        NodeID:       peerID,
        Connection:   connectionPool,
        LastPingTime: time.Now(),
        IsAlive:      true,                
        EncryptionKey: conn.EncryptionKey, // Assuming WebRTCConn now has EncryptionKey
    }

    wm.Peers[peerID] = peerConn

    fmt.Printf("Peer %s connected via WebRTC\n", peerID)

    // Adjusted for correct fields in WebRTCConnection struct
    webRTCConnection := ledger.WebRTCConnection{
        ConnectionID: conn.ConnectionID,
        PeerID:       peerID,
        Timestamp:    time.Now(), // Assuming this is the correct field instead of "ConnectedAt"
    }

    err := wm.LedgerInstance.RecordWebRTCConnection(webRTCConnection)
    if err != nil {
        return fmt.Errorf("Failed to record WebRTC connection for peer %s in the ledger: %v", peerID, err)
    }

    return nil
}


// RemovePeerConnection removes a peer connection from the WebRTCManager
func (wm *WebRTCManager) RemovePeerConnection(peerID string) error {
	wm.lock.Lock()
	defer wm.lock.Unlock()

	if _, exists := wm.Peers[peerID]; !exists {
		return fmt.Errorf("Peer %s not found", peerID)
	}

	delete(wm.Peers, peerID)
	fmt.Printf("Peer %s disconnected from WebRTC\n", peerID)

	// Remove connection details from the ledger
	err := wm.LedgerInstance.RemoveWebRTCConnection(peerID)
	if err != nil {
		return fmt.Errorf("Failed to remove WebRTC connection for peer %s from the ledger: %v", peerID, err)
	}

	return nil
}

// SendEncryptedMessage sends an encrypted message over a WebRTC connection
func (wm *WebRTCManager) SendEncryptedMessage(peerID string, message []byte) error {
    wm.lock.Lock()
    defer wm.lock.Unlock()

    peerConn, exists := wm.Peers[peerID]
    if !exists || !peerConn.IsAlive {  // Use IsAlive instead of IsActive
        return fmt.Errorf("Peer %s is not connected or inactive", peerID)
    }

    // Create an encryption instance
    encryption := &common.Encryption{}

    // Encrypt the message before sending
    encryptedMessage, err := encryption.EncryptData("AES", message, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("Failed to encrypt message: %v", err)
    }

    // Assuming you have a method for sending encrypted messages in your connection pool
    err = peerConn.Connection.Send(encryptedMessage)  // Adjust to the correct method for sending
    if err != nil {
        return fmt.Errorf("Failed to send encrypted message to peer %s: %v", peerID, err)
    }

    fmt.Printf("Encrypted message sent to peer %s\n", peerID)
    return nil
}

// ReceiveEncryptedMessage listens for an encrypted message from a peer and decrypts it
func (wm *WebRTCManager) ReceiveEncryptedMessage(peerID string) ([]byte, error) {
    wm.lock.Lock()
    defer wm.lock.Unlock()

    peerConn, exists := wm.Peers[peerID]
    if !exists || !peerConn.IsAlive {  // Use IsAlive instead of IsActive
        return nil, fmt.Errorf("Peer %s is not connected or inactive", peerID)
    }

    // Receive the encrypted message over the connection pool
    encryptedMessage, err := peerConn.Connection.Receive()  // Receive through the connection pool
    if err != nil {
        return nil, fmt.Errorf("Failed to receive message from peer %s: %v", peerID, err)
    }

    // Create an encryption instance
    encryption := &common.Encryption{}

    // Decrypt the received message
    decryptedMessage, err := encryption.DecryptData(encryptedMessage, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("Failed to decrypt message from peer %s: %v", peerID, err)
    }

    fmt.Printf("Message received from peer %s and decrypted\n", peerID)
    return decryptedMessage, nil
}


// BroadcastToAllPeers sends an encrypted message to all connected peers
func (wm *WebRTCManager) BroadcastToAllPeers(message []byte) error {
    wm.lock.Lock()
    defer wm.lock.Unlock()

    // Create an encryption instance
    encryption := &common.Encryption{}

    // Encrypt the message
    encryptedMessage, err := encryption.EncryptData("AES", message, common.EncryptionKey)  // Corrected to include encryption method as the first argument
    if err != nil {
        return fmt.Errorf("Failed to encrypt message for broadcast: %v", err)
    }

    for peerID, peerConn := range wm.Peers {
        if peerConn.IsAlive { // Use IsAlive instead of IsActive
            err := peerConn.Connection.Send(encryptedMessage) // Use ConnectionPool.Send method
            if err != nil {
                fmt.Printf("Failed to send message to peer %s: %v\n", peerID, err)
            } else {
                fmt.Printf("Message broadcasted to peer %s\n", peerID)
            }
        }
    }

    return nil
}


// ValidateConnectionHealth checks the status of all peer connections
func (wm *WebRTCManager) ValidateConnectionHealth() {
    wm.lock.Lock()
    defer wm.lock.Unlock()

    for peerID, peerConn := range wm.Peers {
        if !peerConn.Connection.IsAlive() { // Use IsAlive method on ConnectionPool
            fmt.Printf("Peer %s connection is inactive, removing from WebRTC\n", peerID)
            peerConn.IsAlive = false
            wm.RemovePeerConnection(peerID)
        }
    }
}

// SyncWithLedger ensures all WebRTC connections are synced with the ledger for auditing
func (wm *WebRTCManager) SyncWithLedger() {
    for peerID, peerConn := range wm.Peers {
        if peerConn == nil {
            continue // Skip if the peerConn is nil
        }

        // Generate or use an existing ConnectionID for the peer connection
        connectionID := fmt.Sprintf("conn-%s-%s", peerID, time.Now().Format("20060102150405"))

        // Create a WebRTCConnection struct for the ledger
        webRTCConnection := ledger.WebRTCConnection{
            ConnectionID: connectionID,  // Use generated or managed ConnectionID
            PeerID:       peerID,
            Timestamp:    time.Now(),    // Use Timestamp instead of ConnectedAt
        }

        err := wm.LedgerInstance.RecordWebRTCConnection(webRTCConnection)
        if err != nil {
            fmt.Printf("Failed to sync WebRTC connection for peer %s with ledger: %v\n", peerID, err)
        }
    }
}




// generateMessageHash creates a SHA-256 hash for verifying message integrity
func generateMessageHash(message []byte) string {
	hash := sha256.New()
	hash.Write(message)
	return hex.EncodeToString(hash.Sum(nil))
}


