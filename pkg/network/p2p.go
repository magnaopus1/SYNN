package network

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// Peer represents a connected peer in the network
type Peer struct {
    Address     string    // Peer address (IP/Port)
    Conn        net.Conn  // Peer connection
	PublicKey string    // Public key of the peer (as a string, representing the AES key)
    LastSeen    time.Time // Last time the peer was seen (used for timeouts)
}

// P2PNetwork manages the peer-to-peer network and communication
type P2PNetwork struct {
	Address          string                // The node's own address
	MessageQueue     chan P2PMessage       // Queue to store outgoing messages
	PendingBlocks    []common.Block        // Blocks pending validation
	Peers            map[string]*Peer      // Connected peers
	IncomingMessages chan P2PMessage       // Incoming message channel
	OutgoingMessages chan P2PMessage       // Outgoing message channel
	PeerLock         sync.Mutex            // Synchronizes access to peers
	NodeKey          *common.NodeKey       // Node's public-private keypair for encryption
	LedgerInstance   *ledger.Ledger        // Pointer to the ledger for transaction and block management
	ConsensusEngine  *common.SynnergyConsensus // Pointer to the Synnergy Consensus engine
}

// P2PMessage represents the structure of a peer-to-peer message
type P2PMessage struct {
	Sender        string    // Sender's public key (encoded as string)
	Recipient     string    // Recipient peer address
	EncryptedKey  string    // Base64 encoded AES key
	Content       string    // Base64 encoded encrypted message content
	Timestamp     time.Time // Timestamp of the message
}


// NewP2PNetwork initializes a new P2P network with encryption, ledger, and consensus integration
func NewP2PNetwork(nodeKey *common.NodeKey, ledger *ledger.Ledger, consensusEngine *common.SynnergyConsensus, address string) *P2PNetwork {
	return &P2PNetwork{
		Peers:            make(map[string]*Peer),
		IncomingMessages: make(chan P2PMessage, 100),
		OutgoingMessages: make(chan P2PMessage, 100),
		NodeKey:          nodeKey,
		LedgerInstance:   ledger,
		ConsensusEngine:  consensusEngine,
		Address:          address,  // Use the address passed as an argument
		MessageQueue:     make(chan P2PMessage, 100),
	}
}


// Start starts the P2P network and begins accepting connections
func (p2p *P2PNetwork) Start(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting network listener:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("P2P Network started on port %s\n", port)

	go p2p.handleIncomingMessages()
	go p2p.handleOutgoingMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go p2p.handlePeerConnection(conn)
	}
}

// handlePeerConnection manages an incoming peer connection
func (p2p *P2PNetwork) handlePeerConnection(conn net.Conn) {
	defer conn.Close()

	peerAddress := conn.RemoteAddr().String()
	fmt.Printf("Connected to peer: %s\n", peerAddress)

	p2p.PeerLock.Lock()
	p2p.Peers[peerAddress] = &Peer{
		Address:   peerAddress,
		Conn:      conn,            // Use Conn instead of Connection
		PublicKey: "",              // Initialize as an empty string instead of nil
		LastSeen:  time.Now(),      // Use LastSeen instead of LastContact
	}
	p2p.PeerLock.Unlock()

	// Perform a handshake with the peer to exchange public keys
	err := p2p.performHandshake(peerAddress)
	if err != nil {
		fmt.Println("Handshake failed with peer:", peerAddress, err)

		// Remove the peer if handshake fails
		p2p.PeerLock.Lock()
		delete(p2p.Peers, peerAddress)
		p2p.PeerLock.Unlock()
		return
	}

	// Listen for messages from the peer
	p2p.listenForMessages(peerAddress)
}



// performHandshake exchanges AES-encrypted messages for secure communication with a peer
func (p2p *P2PNetwork) performHandshake(peerAddress string) error {
	peer := p2p.Peers[peerAddress]
	conn := peer.Conn // Use the peer's connection

	// Encode the node's address as a message and encrypt it using AES
	encodedMsg, err := common.EncodeMessage(p2p.Address)
	if err != nil {
		return fmt.Errorf("failed to encode address: %v", err)
	}

	// Send the encrypted message over the connection
	_, err = conn.Write([]byte(encodedMsg))
	if err != nil {
		return fmt.Errorf("failed to send encrypted address: %v", err)
	}

	// Receive the peer's encrypted message (address or public key)
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read encrypted peer message: %v", err)
	}

	// Decode the peer's encrypted message
	peerDecodedMsg, err := common.DecodeMessage(string(buffer[:n]))
	if err != nil {
		return fmt.Errorf("failed to decode peer message: %v", err)
	}

	// Store the decoded peer's address
	peer.Address = peerDecodedMsg

	fmt.Printf("Handshake successful with peer: %s\n", peerAddress)
	return nil
}

// listenForMessages listens for incoming messages from a peer
func (p2p *P2PNetwork) listenForMessages(peerAddress string) {
	peer := p2p.Peers[peerAddress]
	conn := peer.Conn  // Corrected: Conn, not Connection

	for {
		// Create a buffer to store the incoming data
		buffer := make([]byte, 4096)

		// Read the message from the connection
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading message from %s: %v\n", peerAddress, err)
			return
		}

		// Decode the message using AES decryption
		decryptedMessage, err := common.DecodeMessage(string(buffer[:n]))
		if err != nil {
			fmt.Printf("Error decoding message from %s: %v\n", peerAddress, err)
			continue
		}

		// Construct the P2PMessage
		message := P2PMessage{
			Sender:  peerAddress,
			Content: decryptedMessage,  // Assuming the decrypted message is in the correct format
		}

		// Process the message and push it to the incoming message queue
		fmt.Printf("Received encrypted message from %s: %s\n", peerAddress, decryptedMessage)
		p2p.IncomingMessages <- message
	}
}


// SendMessage sends an encrypted message to a specific peer
func (p2p *P2PNetwork) SendMessage(recipient, content string) error {
	// Check if the recipient peer is connected
	peer, exists := p2p.Peers[recipient]
	if !exists {
		return fmt.Errorf("peer %s not connected", recipient)
	}

	// Step 1: Create an encryption instance
	encryptionInstance, err := common.NewEncryption(256) // Assuming NewEncryption creates an AES instance with a 256-bit key
	if err != nil {
		return fmt.Errorf("failed to create encryption instance: %v", err)
	}

	// Step 2: Generate a random AES key
	aesKey := make([]byte, 32) // 32 bytes = 256-bit AES key
	if _, err := rand.Read(aesKey); err != nil {
		return fmt.Errorf("failed to generate AES key: %v", err)
	}

	// Step 3: Encrypt the message content using AES
	encryptedContent, err := encryptionInstance.EncryptData("AES", aesKey, []byte(content)) // Encrypting the message content
	if err != nil {
		return fmt.Errorf("failed to encrypt message content: %v", err)
	}

	// Step 4: Get the recipient's public key
	recipientPublicKey, err := encryptionInstance.DecodePublicKey(peer.PublicKey) // Decode recipient's public key from string to *rsa.PublicKey
	if err != nil {
		return fmt.Errorf("failed to decode recipient's public key: %v", err)
	}

	// Step 5: Encrypt the AES key using the recipient's public key
	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, recipientPublicKey, aesKey, nil)
	if err != nil {
		return fmt.Errorf("failed to encrypt AES key: %v", err)
	}

	// Step 6: Encode the encrypted AES key into base64
	encodedAESKey := base64.StdEncoding.EncodeToString(encryptedAESKey)

	// Step 7: Encode the sender's public key for inclusion in the message
	encodedSenderPublicKey, err := encryptionInstance.EncodePublicKey(p2p.NodeKey.PublicKey) // Encode sender's public key
	if err != nil {
		return fmt.Errorf("failed to encode sender's public key: %v", err)
	}

	// Step 8: Create the P2PMessage with the encrypted AES key and message content
	message := P2PMessage{
		Sender:        encodedSenderPublicKey,  // Sender's public key (encoded as string)
		Recipient:     recipient,               // Recipient peer address
		EncryptedKey:  encodedAESKey,           // Base64 encoded encrypted AES key
		Content:       base64.StdEncoding.EncodeToString(encryptedContent), // Base64 encoded encrypted message content
		Timestamp:     time.Now(),              // Current timestamp
	}

	// Step 9: Add the message to the outgoing queue
	p2p.OutgoingMessages <- message
	fmt.Printf("Encrypted message sent to %s\n", recipient)

	return nil
}

// handleIncomingMessages processes incoming messages and acts upon them
func (p2p *P2PNetwork) handleIncomingMessages() {
	for message := range p2p.IncomingMessages {
		// Process the message content, e.g., validating a block or transaction
		fmt.Printf("Processing incoming message from %s\n", message.Sender)
		// Add further logic to process block and transaction validation
	}
}

// handleOutgoingMessages sends queued messages to the appropriate peers
func (p2p *P2PNetwork) handleOutgoingMessages() {
	for message := range p2p.OutgoingMessages {
		peer, exists := p2p.Peers[message.Recipient]
		if !exists {
			fmt.Printf("Peer %s not connected\n", message.Recipient)
			continue
		}

		// Correctly handle the return values from EncodeMessage
		encodedMessage, err := common.EncodeMessage(message.Content)
		if err != nil {
			fmt.Printf("Error encoding message for %s: %v\n", message.Recipient, err)
			continue
		}

		// Send the encoded message over the connection
		_, err = peer.Conn.Write([]byte(encodedMessage))
		if err != nil {
			fmt.Printf("Error sending message to %s: %v\n", message.Recipient, err)
		} else {
			fmt.Printf("Message sent to %s\n", message.Recipient)
		}
	}
}


// encryptMessage encrypts a message using AES with the recipient's public key (as a string here, representing AES key)
func (p2p *P2PNetwork) encryptMessage(message string, recipientPublicKey string) (string, error) {
	// Create an encryption instance (assuming NewEncryption creates an AES-based encryption with the specified key size)
	encryptionInstance, err := common.NewEncryption(256) // Assuming AES 256-bit encryption
	if err != nil {
		return "", fmt.Errorf("failed to create encryption instance: %v", err)
	}

	// Here recipientPublicKey represents the AES key as a string, so we decode it into bytes
	aesKey := []byte(recipientPublicKey) // Convert the string into []byte key

	// Encrypt the message using the AES key
	encryptedContent, err := encryptionInstance.EncryptData("AES", aesKey, []byte(message)) // Use the AES key and message
	if err != nil {
		return "", fmt.Errorf("encryption failed: %v", err)
	}

	// Convert the encrypted message to a base64 encoded string
	encryptedMessage := base64.StdEncoding.EncodeToString(encryptedContent)

	return encryptedMessage, nil
}



// decryptMessage decrypts a message using AES with the node's private key (RSA to decrypt AES key, then AES for the message)
func (p2p *P2PNetwork) decryptMessage(encryptedMessage, encryptedAESKey string) (string, error) {
	// Step 1: Decode the encrypted AES key (base64-encoded) and decrypt it using RSA private key
	encryptedAESKeyBytes, err := base64.StdEncoding.DecodeString(encryptedAESKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted AES key: %v", err)
	}

	// Decrypt the AES key using the node's RSA private key
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, p2p.NodeKey.PrivateKey, encryptedAESKeyBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt AES key: %v", err)
	}

	// Step 2: Decode the encrypted message (base64-encoded) to get the ciphertext
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted message: %v", err)
	}

	// Step 3: Create an encryption instance (assuming it's AES)
	encryptionInstance, err := common.NewEncryption(256) // AES 256-bit encryption instance
	if err != nil {
		return "", fmt.Errorf("failed to create encryption instance: %v", err)
	}

	// Step 4: Decrypt the message using the decrypted AES key
	decrypted, err := encryptionInstance.DecryptData(encryptedBytes, aesKey) // Use the decrypted AES key
	if err != nil {
		return "", fmt.Errorf("decryption failed: %v", err)
	}

	// Step 5: Return the decrypted message as a string
	return string(decrypted), nil
}



// BroadcastBlock broadcasts a block to all connected peers
func (p2p *P2PNetwork) BroadcastBlock(block common.Block) error {
	blockHash := generateBlockHash(block)
	content := fmt.Sprintf("Broadcasting Block #%d: %s", block.Index, blockHash)

	for peerAddress := range p2p.Peers {
		err := p2p.SendMessage(peerAddress, content)
		if err != nil {
			fmt.Printf("Failed to broadcast block to %s: %v\n", peerAddress, err)
		}
	}

	return nil
}

// generateBlockHash generates a SHA-256 hash for a block
func generateBlockHash(block common.Block) string {
	hashInput := fmt.Sprintf("%d%s%s%d", block.Index, block.Timestamp.String(), block.PrevHash, block.Nonce)
	hash := sha256.New()
	hash.Write([]byte(hashInput))
	return hex.EncodeToString(hash.Sum(nil))
}


// StartNetwork starts listening for incoming peer connections and message passing
func (network *P2PNetwork) StartNetwork(port string) {
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        fmt.Println("Error starting network listener:", err)
        return
    }
    fmt.Printf("Node %s started listening on port %s\n", network.Address, port)

    go network.handleIncomingConnections(listener)
    go network.processOutgoingMessages()
}

// handleIncomingConnections handles incoming peer connections and spawns goroutines to listen for messages
func (network *P2PNetwork) handleIncomingConnections(listener net.Listener) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        go network.handlePeerConnection(conn)
    }
}

// processOutgoingMessages processes the message queue and sends messages to the corresponding peers
func (network *P2PNetwork) processOutgoingMessages() {
    for message := range network.MessageQueue {
        peerConn, exists := network.Peers[message.Recipient]
        if !exists {
            fmt.Printf("Peer %s is not connected. Dropping message.\n", message.Recipient)
            continue
        }

        // Step 1: Encode the message using AES encryption (assuming common.EncodeMessage does this)
        encodedMessage, err := common.EncodeMessage(message.Content) // Encode the message content
        if err != nil {
            fmt.Printf("Failed to encode message for %s: %v\n", message.Recipient, err)
            continue
        }

        // Step 2: Send the encoded message over the peer's connection
        _, err = peerConn.Conn.Write([]byte(encodedMessage)) // Assuming peerConn has a `Conn` field of type `net.Conn`
        if err != nil {
            fmt.Printf("Failed to send message to %s: %v\n", message.Recipient, err)
        } else {
            fmt.Printf("Encrypted message sent to %s\n", message.Recipient)
        }
    }
}

// EncryptMessage encrypts a message using AES
func (network *P2PNetwork) EncryptMessage(message string, aesKey []byte) (string, error) {
    // Step 1: Create an encryption instance with the provided AES key
    encryptionInstance, err := common.NewEncryption(256) // Assuming 256-bit encryption
    if err != nil {
        return "", fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Step 2: Encrypt the message using AES and the provided key
    encryptedMessageBytes, err := encryptionInstance.EncryptData("AES", []byte(message), aesKey) // Pass aesKey here
    if err != nil {
        return "", fmt.Errorf("encryption failed: %v", err)
    }

    // Step 3: Convert the encrypted bytes to a base64 string for readability
    encryptedMessage := base64.StdEncoding.EncodeToString(encryptedMessageBytes)

    return encryptedMessage, nil
}

// DecryptMessage decrypts a received encrypted message using AES
func (network *P2PNetwork) DecryptMessage(encryptedMessage string, aesKey []byte) (string, error) {
    // Step 1: Create an encryption instance with the provided AES key
    encryptionInstance, err := common.NewEncryption(256) // Assuming 256-bit encryption
    if err != nil {
        return "", fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Step 2: Decode the base64-encoded encrypted message
    encryptedMessageBytes, err := base64.StdEncoding.DecodeString(encryptedMessage)
    if err != nil {
        return "", fmt.Errorf("failed to decode encrypted message: %v", err)
    }

    // Step 3: Decrypt the message using AES and the provided key
    decryptedMessageBytes, err := encryptionInstance.DecryptData(encryptedMessageBytes, aesKey) // Pass aesKey here
    if err != nil {
        return "", fmt.Errorf("decryption failed: %v", err)
    }

    // Step 4: Convert the decrypted byte slice back to a string
    decryptedMessage := string(decryptedMessageBytes)

    return decryptedMessage, nil
}




// generateBlockHash creates a unique hash for a block
func (network *P2PNetwork) generateBlockHash(block common.Block) string {
    hashInput := fmt.Sprintf("%d%s%s%d", block.Index, block.Timestamp.String(), block.PrevHash, block.Nonce)
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}
