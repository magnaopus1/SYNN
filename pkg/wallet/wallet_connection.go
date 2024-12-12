package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
	"time"
)

// WalletConnection manages the secure connection of a wallet to the blockchain network.
type WalletConnection struct {
	ledgerInstance *ledger.Ledger
	networkManager *network.NetworkManager  // Network manager to handle network operations
	connections    map[string]string        // Active connections mapped by wallet address and status
	mutex          sync.Mutex
}

// NewWalletConnection initializes a new WalletConnection instance.
func NewWalletConnection(ledgerInstance *ledger.Ledger, networkManager *network.NetworkManager) *WalletConnection {
	return &WalletConnection{
		ledgerInstance: ledgerInstance,
		networkManager: networkManager,  
		connections:    make(map[string]string), // Properly initialize the map with make
	}
}


// ConnectWallet connects a wallet to a node in the blockchain network securely.
func (wc *WalletConnection) ConnectWallet(walletAddress string, privateKey *ecdsa.PrivateKey) error {
    wc.mutex.Lock()
    defer wc.mutex.Unlock()

    // Ensure the wallet is not already connected
    if _, exists := wc.connections[walletAddress]; exists {
        return fmt.Errorf("wallet %s is already connected", walletAddress)
    }

    // Generate a secure token for wallet authentication (not used in current implementation)
    _, err := wc.generateAuthenticationToken(walletAddress, privateKey)
    if err != nil {
        return fmt.Errorf("failed to generate authentication token: %v", err)
    }

    // Use the network manager to establish a connection
    err = wc.networkManager.ConnectToPeer(walletAddress)
    if err != nil {
        return fmt.Errorf("failed to establish connection: %v", err)
    }

    // Create a connection event
    connectionEvent := ledger.ConnectionEvent{
        EventID:      fmt.Sprintf("conn-%d", time.Now().UnixNano()), // Unique Event ID
        ConnectionID: walletAddress, // Using walletAddress as ConnectionID
        WalletID:     walletAddress,
        EventType:    ledger.EventTypeConnected,
        EventTime:    time.Now(),
        Details:      "Wallet connected successfully",
    }

    // Log the connection event in the ledger
    if err := wc.ledgerInstance.RecordConnectionEvent(walletAddress, connectionEvent); err != nil {
        return fmt.Errorf("failed to log connection: %v", err)
    }

    // Store the connection status
    wc.connections[walletAddress] = "Connected"

    fmt.Printf("Wallet %s connected successfully.\n", walletAddress)
    return nil
}



// DisconnectWallet safely disconnects a wallet from the blockchain network.
func (wc *WalletConnection) DisconnectWallet(walletAddress string) error {
    wc.mutex.Lock()
    defer wc.mutex.Unlock()

    _, exists := wc.connections[walletAddress]
    if !exists {
        return errors.New("wallet is not connected")
    }

    // Create a connection event for disconnection
    connectionEvent := ledger.ConnectionEvent{
        EventID:      fmt.Sprintf("conn-%d", time.Now().UnixNano()), // Unique Event ID
        ConnectionID: walletAddress,                                 // Using walletAddress as ConnectionID
        WalletID:     walletAddress,
        EventType:    "DISCONNECTED",                               // Use a string constant or direct string
        EventTime:    time.Now(),
        Details:      "Wallet disconnected successfully",
    }

    // Log the disconnection event in the ledger
    if err := wc.ledgerInstance.RecordConnectionEvent(walletAddress, connectionEvent); err != nil {
        return fmt.Errorf("failed to log disconnection: %v", err)
    }

    // Remove the connection from the map
    delete(wc.connections, walletAddress)
    fmt.Printf("Wallet %s disconnected successfully.\n", walletAddress)

    return nil
}



// generateAuthenticationToken creates an encrypted token using the wallet's private key.
func (wc *WalletConnection) generateAuthenticationToken(walletAddress string, privateKey *ecdsa.PrivateKey) (string, error) {
    // Generate a random challenge
    challenge := make([]byte, 32)
    if _, err := rand.Read(challenge); err != nil {
        return "", fmt.Errorf("failed to generate random challenge: %v", err)
    }

    // Sign the challenge with the wallet's private key
    r, s, err := ecdsa.Sign(rand.Reader, privateKey, challenge)
    if err != nil {
        return "", fmt.Errorf("failed to sign challenge: %v", err)
    }

    // Combine the signature (r and s) with the challenge to form the auth token
    authToken := fmt.Sprintf("%s:%x:%x", walletAddress, r, s)

    // Create an encryption instance and encrypt the auth token
    encryptionInstance := &common.Encryption{} // Create encryption instance
    encryptedToken, err := encryptionInstance.EncryptData("AES", []byte(authToken), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt authentication token: %v", err)
    }

    return string(encryptedToken), nil
}

// VerifyWalletAuthenticationToken verifies an authentication token received from a wallet.
func (wc *WalletConnection) VerifyWalletAuthenticationToken(walletAddress string, authToken string, publicKey *ecdsa.PublicKey) (bool, error) {
    // Create an encryption instance and decrypt the auth token
    encryptionInstance := &common.Encryption{} // Create encryption instance
    decryptedToken, err := encryptionInstance.DecryptData([]byte(authToken), common.EncryptionKey)
    if err != nil {
        return false, fmt.Errorf("failed to decrypt authentication token: %v", err)
    }

    // Parse the token into wallet address, r, and s
    var r, s big.Int
    var addr string
    _, err = fmt.Sscanf(string(decryptedToken), "%s:%x:%x", &addr, &r, &s)
    if err != nil {
        return false, fmt.Errorf("invalid token format: %v", err)
    }

    // Ensure the wallet address matches
    if addr != walletAddress {
        return false, errors.New("wallet address mismatch in token")
    }

    // Recreate the challenge from the wallet address (here we assume it's deterministic)
    challenge := []byte(addr)

    // Verify the signature using the public key
    if ecdsa.Verify(publicKey, challenge, &r, &s) {
        return true, nil
    }

    return false, errors.New("invalid wallet authentication")
}

// IsWalletConnected checks whether a wallet is currently connected.
func (wc *WalletConnection) IsWalletConnected(walletAddress string) bool {
	wc.mutex.Lock()
	defer wc.mutex.Unlock()

	_, exists := wc.connections[walletAddress]
	return exists
}

// BroadcastTransaction sends a signed transaction from the wallet to the network.
func (wc *WalletConnection) BroadcastTransaction(walletAddress string, transaction []byte, amount float64) error {
    wc.mutex.Lock()
    defer wc.mutex.Unlock()

    // Ensure the wallet is connected
    _, exists := wc.connections[walletAddress]
    if !exists {
        return errors.New("wallet is not connected to the network")
    }

    // Convert the transaction to a string or hex format
    transactionStr := hex.EncodeToString(transaction)

    // Use the NetworkManager to send the encrypted transaction to the peer
    err := wc.networkManager.SendEncryptedMessage(walletAddress, transactionStr)
    if err != nil {
        return fmt.Errorf("failed to send transaction: %v", err)
    }

    // Log the transaction broadcast in the ledger
    err = wc.ledgerInstance.RecordTransaction(walletAddress, transactionStr, amount) // Remove currency argument
    if err != nil {
        return fmt.Errorf("failed to log transaction: %v", err)
    }

    return nil
}




// ListenForNetworkEvents listens for incoming network events (e.g., incoming transactions or messages) for the connected wallet.
func (wc *WalletConnection) ListenForNetworkEvents(walletAddress string) error {
	wc.mutex.Lock()
	defer wc.mutex.Unlock()

	// Check if the wallet is connected
	_, exists := wc.connections[walletAddress]
	if !exists {
		return errors.New("wallet is not connected to the network")
	}

	// Use the NetworkManager to listen for messages for the connected wallet
	err := wc.networkManager.ReceiveMessages(walletAddress)
	if err != nil {
		return fmt.Errorf("failed to listen for network events: %v", err)
	}

	fmt.Printf("Listening for network events for wallet: %s\n", walletAddress)
	return nil
}



// SecureWalletData securely transmits wallet data across the network.
func (wc *WalletConnection) SecureWalletData(walletAddress string, data []byte) error {
	wc.mutex.Lock()
	defer wc.mutex.Unlock()

	// Check if the wallet is connected
	_, exists := wc.connections[walletAddress]
	if !exists {
		return errors.New("wallet is not connected to the network")
	}

	// Create an instance of Encryption and encrypt the data
	encryptionInstance := &common.Encryption{}
	encryptedData, err := encryptionInstance.EncryptData("AES", data, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt wallet data: %v", err)
	}

	// Use the network manager to send the encrypted data to the network
	err = wc.networkManager.SendEncryptedMessage(walletAddress, string(encryptedData))
	if err != nil {
		return fmt.Errorf("failed to send encrypted data: %v", err)
	}

	fmt.Printf("Encrypted wallet data sent for wallet %s.\n", walletAddress)
	return nil
}

