package network

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewPeerDiscoveryManager initializes a new peer discovery manager
func NewPeerDiscoveryManager(localAddress, localPublicKey string, ledgerInstance *ledger.Ledger, knownNetworks []string) *PeerDiscoveryManager {
	return &PeerDiscoveryManager{
		Peers:         make(map[string]*Peer),
		Ledger:        ledgerInstance,
		LocalAddress:  localAddress,
		LocalPublicKey: localPublicKey,
		KnownNetworks: knownNetworks,
	}
}

// DiscoverPeers discovers new peers on the network by sending a "HELLO" handshake message
func (pdm *PeerDiscoveryManager) DiscoverPeers() {
	for _, network := range pdm.KnownNetworks {
		fmt.Printf("Discovering peers in network: %s\n", network)
		
		conn, err := net.Dial("udp", network)
		if err != nil {
			fmt.Printf("Failed to connect to network %s: %v\n", network, err)
			continue
		}

		// Send a "HELLO" handshake message encrypted with our local public key
		helloMessage := fmt.Sprintf("HELLO from %s", pdm.LocalAddress)

		// Create an instance of Encryption
		encryption := &common.Encryption{}

		// Encrypt the hello message using our public key
		encryptedHelloMessage, err := encryption.EncryptData("RSA", []byte(helloMessage), []byte(pdm.LocalPublicKey))
		if err != nil {
			fmt.Printf("Failed to encrypt handshake message: %v\n", err)
			conn.Close()  // Ensure the connection is closed on error
			continue
		}

		_, err = conn.Write(encryptedHelloMessage)
		if err != nil {
			fmt.Printf("Failed to send handshake to network %s: %v\n", network, err)
		}

		conn.Close()
	}
}



// HandleIncomingHandshake handles incoming handshake messages from other peers
func (pdm *PeerDiscoveryManager) HandleIncomingHandshake(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Failed to read incoming handshake: %v\n", err)
		return
	}

	// Create an encryption instance
	encryption := &common.Encryption{}

	// Decrypt the message using the local private key (assuming RSA encryption)
	encryptedMessage := buffer[:n]
	decryptedMessage, err := encryption.DecryptData(encryptedMessage, []byte(pdm.LocalPrivateKey)) // DecryptData now takes only two arguments
	if err != nil {
		fmt.Printf("Failed to decrypt handshake message: %v\n", err)
		return
	}

	// Parse the decrypted handshake
	message := string(decryptedMessage)
	if !strings.HasPrefix(message, "HELLO from ") {
		fmt.Printf("Invalid handshake message received: %s\n", message)
		return
	}

	peerAddress := strings.TrimPrefix(message, "HELLO from ")
	peerPublicKey := pdm.extractPublicKeyFromMessage(message)

	// Add the peer to the discovered list
	pdm.addPeer(peerAddress, peerPublicKey)
}



// addPeer adds a new peer to the list of discovered peers
func (pdm *PeerDiscoveryManager) addPeer(address, publicKey string) {
	if _, exists := pdm.Peers[address]; !exists {
		pdm.Peers[address] = &Peer{
			Address:   address,
			PublicKey: publicKey,
		}
		fmt.Printf("Discovered new peer: %s\n", address)

		// Create PeerInfo struct and pass it to RecordPeerDiscovery
		peerInfo := ledger.PeerInfo{
			Address:   address,
			PublicKey: publicKey,
		}
		pdm.Ledger.RecordPeerDiscovery(peerInfo)
	}
}


// extractPublicKeyFromMessage extracts the public key from the handshake message
func (pdm *PeerDiscoveryManager) extractPublicKeyFromMessage(message string) string {
	// Simulate extracting public key, actual implementation would depend on the message format
	hash := sha256.New()
	hash.Write([]byte(message))
	return hex.EncodeToString(hash.Sum(nil))
}

// BroadcastNewBlock broadcasts a newly mined block to all discovered peers
func (pdm *PeerDiscoveryManager) BroadcastNewBlock(block common.Block) {
	for _, peer := range pdm.Peers {
		conn, err := net.Dial("tcp", peer.Address)
		if err != nil {
			fmt.Printf("Failed to connect to peer %s: %v\n", peer.Address, err)
			continue
		}

		// Create an encryption instance
		encryption := &common.Encryption{}

		// Encrypt the block data with the peer's public key
		blockData := fmt.Sprintf("NEW BLOCK: %s", block.Hash)
		encryptedBlockData, err := encryption.EncryptData("RSA", []byte(blockData), []byte(peer.PublicKey))  // Assuming RSA encryption and PEM encoded key
		if err != nil {
			fmt.Printf("Failed to encrypt block data for peer %s: %v\n", peer.Address, err)
			conn.Close()
			continue
		}

		_, err = conn.Write(encryptedBlockData)
		if err != nil {
			fmt.Printf("Failed to send block data to peer %s: %v\n", peer.Address, err)
		}

		conn.Close()
	}
}


// ListenForIncomingPeers listens on a specific port for incoming peer handshakes
func (pdm *PeerDiscoveryManager) ListenForIncomingPeers(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Failed to start peer discovery listener on port %s: %v\n", port, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening for incoming peers on port %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept incoming connection: %v\n", err)
			continue
		}

		go pdm.HandleIncomingHandshake(conn)
	}
}
