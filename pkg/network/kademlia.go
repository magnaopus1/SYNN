package network

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewKademliaDHT initializes a new Kademlia Distributed Hash Table for the network
func NewKademliaDHT(nodeID string, ledger *ledger.Ledger) *KademliaDHT {
    return &KademliaDHT{
        NodeID:  nodeID,
        KBucket: make(map[string]KademliaNode),
        ledger:  ledger,
    }
}

// CalculateNodeID generates a unique node ID using SHA-256 hashing
func CalculateNodeID(input string) string {
    hash := sha256.New()
    hash.Write([]byte(input))
    return hex.EncodeToString(hash.Sum(nil))
}

// AddNode adds a node to the Kademlia DHT
func (k *KademliaDHT) AddNode(nodeID, ip string, location GeoLocation) {
    k.lock.Lock()
    defer k.lock.Unlock()

    node := KademliaNode{
        ID:         nodeID,
        IPAddress:  ip,
        LastActive: time.Now(),
        Location:   location,
    }
    k.KBucket[nodeID] = node
    fmt.Printf("Node %s added to Kademlia DHT at IP: %s\n", nodeID, ip)

    // Log the event with only the required arguments (event type and node ID)
    k.ledger.LogNodeEvent("NodeAdded", nodeID)
}


// FindClosestNodes finds the closest nodes to a given target ID in the Kademlia DHT
func (k *KademliaDHT) FindClosestNodes(targetID string, count int) []KademliaNode {
    k.lock.Lock()
    defer k.lock.Unlock()

    distances := make(map[string]*big.Int)
    targetHash := new(big.Int)
    targetHash.SetString(targetID, 16)

    for nodeID := range k.KBucket {
        nodeHash := new(big.Int)
        nodeHash.SetString(nodeID, 16)

        xorResult := new(big.Int).Xor(targetHash, nodeHash)
        distances[nodeID] = xorResult
    }

    // Sort nodes by XOR distance (closest nodes first)
    sortedNodes := sortNodesByDistance(k.KBucket, distances)

    // Return the closest nodes up to the requested count
    if len(sortedNodes) > count {
        return sortedNodes[:count]
    }
    return sortedNodes
}

// RemoveNode removes a node from the Kademlia DHT
func (k *KademliaDHT) RemoveNode(nodeID string) {
    k.lock.Lock()
    defer k.lock.Unlock()

    if _, exists := k.KBucket[nodeID]; exists {
        delete(k.KBucket, nodeID)
        fmt.Printf("Node %s removed from Kademlia DHT.\n", nodeID)

        // Log the event with only the required arguments (event type and node ID)
        k.ledger.LogNodeEvent("NodeRemoved", nodeID)
    }
}


// EncryptMessage encrypts a message to be sent between nodes using a shared public key
func (k *KademliaDHT) EncryptMessage(message string, pubKey *common.PublicKey) (string, error) {
    encryption := &common.Encryption{} // Instantiate the Encryption object

    // Use AES encryption (or whatever algorithm you are using)
    encryptedMsg, err := encryption.EncryptData("AES", []byte(message), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt message: %v", err)
    }

    // Return the encrypted message as a hex-encoded string
    return hex.EncodeToString(encryptedMsg), nil
}


// DecryptMessage decrypts a message received from another node using RSA private key
func (k *KademliaDHT) DecryptMessage(encryptedMsg string, privKey *rsa.PrivateKey) (string, error) {
    // Decode the hex-encoded encrypted message
    msgBytes, err := hex.DecodeString(encryptedMsg)
    if err != nil {
        return "", fmt.Errorf("failed to decode encrypted message: %v", err)
    }

    // Decrypt the message using the RSA private key
    decryptedMsg, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, msgBytes)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt message: %v", err)
    }

    // Return the decrypted message as a string
    return string(decryptedMsg), nil
}

// PingNode pings a node in the network to check if it is active
func (k *KademliaDHT) PingNode(nodeID string) error {
    node, exists := k.KBucket[nodeID]
    if !exists {
        return fmt.Errorf("node %s not found in the Kademlia DHT", nodeID)
    }

    // Simulate sending a ping message to the node's IP address
    fmt.Printf("Pinging node %s at IP: %s\n", nodeID, node.IPAddress)
    node.LastActive = time.Now()
    k.KBucket[nodeID] = node

    // Log the event with only the required arguments
    k.ledger.LogNodeEvent("NodePinged", nodeID)
    
    return nil
}


// sortNodesByDistance sorts the nodes by their XOR distance from the target
func sortNodesByDistance(nodes map[string]KademliaNode, distances map[string]*big.Int) []KademliaNode {
    type distanceNodePair struct {
        distance *big.Int
        node     KademliaNode
    }

    // Create a slice of distance-node pairs
    var distanceNodePairs []distanceNodePair
    for nodeID, dist := range distances {
        distanceNodePairs = append(distanceNodePairs, distanceNodePair{distance: dist, node: nodes[nodeID]})
    }

    // Sort by distance (ascending)
    sort.Slice(distanceNodePairs, func(i, j int) bool {
        return distanceNodePairs[i].distance.Cmp(distanceNodePairs[j].distance) == -1
    })

    // Extract the sorted nodes
    var sortedNodes []KademliaNode
    for _, pair := range distanceNodePairs {
        sortedNodes = append(sortedNodes, pair.node)
    }

    return sortedNodes
}
