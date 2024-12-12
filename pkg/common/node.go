package common

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"time"
)

// AuthorityNode represents an authority node with unique permissions
type AuthorityNodeVersion struct {
	NodeID            string    // Unique ID for the authority node
	SecretKey         string    // Secret key for node access control
	CreatedAt         time.Time // Timestamp of node creation
	EncryptedKey      string    // Encrypted form of the secret key
	AuthorityNodeType AuthorityNodeTypes  // Assuming 'nodeType' is a valid type or enum defined elsewhere
}

// AuthorityNodeType defines different types of authority nodes in the network
type AuthorityNodeTypes string

// NodeCategory represents the category of the node (e.g., authority, standard).
type NodeCategory string

const (
    AuthorityCategory NodeCategory = "Authority"
    StandardCategory  NodeCategory = "Standard"
)

// NodeType represents the specific type of the node within a category.
type NodeType string

// Authority node types
const (
    GovernmentNode NodeType = "Government"
    BankingNode    NodeType = "Banking"
    RegulatorNode  NodeType = "Regulator"
)

// Standard node types
const (
    LightningNode NodeType = "Lightning"
    LightNode     NodeType = "Light"
    FullNode      NodeType = "Full"
)

// Node represents a node in the Synnergy Network.
type Node struct {
    Address      string          // The public address of the node
    Name         string          // Name of the node
    NodeCategory NodeCategory    // Category of the node (e.g., authority, standard)
    NodeType     NodeType        // Specific type of the node (e.g., government, light, full)
    NodeKey      *NodeKey  // Public-private key pair for encryption
	IsActive     bool            // Whether the node is active or not
    HasLightVM bool   // Whether the node is running a LightVM
	HasHeavyVM bool   // Whether the node is running a HeavyVM
    
}


// NewNode creates and initializes a new node with both category and type.
func NewNode(address, name string, nodeCategory NodeCategory, nodeType NodeType) (*Node, error) {
    nodeKey, err := GenerateNodeKeyPair()
    if err != nil {
        return nil, fmt.Errorf("failed to generate node key pair: %v", err)
    }

    return &Node{
        Address:      address,
        Name:         name,
        NodeCategory: nodeCategory,
        NodeType:     nodeType,
        NodeKey:      nodeKey,
    }, nil
}

// NodeDescription returns a description of the node with its category and type.
func (n *Node) NodeDescription() string {
    return fmt.Sprintf("Node %s (%s) is of category %s and type %s", n.Name, n.Address, n.NodeCategory, n.NodeType)
}

// GenerateNodeHash creates a unique hash for the node based on its address and public key.
func (n *Node) GenerateNodeHash() (string, error) {
    // Serialize the public key into a DER-encoded byte slice
    pubKeyBytes, err := x509.MarshalPKIXPublicKey(n.NodeKey.PublicKey)
    if err != nil {
        return "", fmt.Errorf("failed to marshal public key: %v", err)
    }

    // Create a unique string to hash using the address and the serialized public key
    hashInput := fmt.Sprintf("%s%s", n.Address, hex.EncodeToString(pubKeyBytes))

    // Generate the hash
    hash := sha256.New()
    hash.Write([]byte(hashInput))

    return hex.EncodeToString(hash.Sum(nil)), nil
}


// ReceiveEncryptedMessage allows the node to handle incoming encrypted messages
func (n *Node) ReceiveEncryptedMessage(message EncryptedMessage) error {
	// Decrypt and process the message here
	fmt.Printf("Node %s received an encrypted message\n", n.Name)

	// Implement the logic to decrypt and handle the message
	return nil
}

// EncryptedMessage represents an encrypted payload to be sent between nodes.
type EncryptedMessage struct {
    CipherText  []byte    // The encrypted message content
    Hash        string    // Hash of the original message for integrity check
    CreatedAt   time.Time // When the message was encrypted (Renamed from Timestamp to CreatedAt)
}


// Implementing the NetworkNode methods for Node
func (n *Node) GetAddress() string {
    return n.Address
}

func (n *Node) IsNodeActive() bool {
    return n.IsActive
}
