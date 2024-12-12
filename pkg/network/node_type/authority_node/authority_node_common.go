package authority_node

import (
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


// Permission defines the structure for specific permissions granted to an authority node
type Permission struct {
	PermissionID string    // Unique ID for the permission
	Description  string    // Description of the permission granted
	GrantedAt    time.Time // Time when the permission was granted
	GrantedBy    string    // ID of the entity granting the permission
}

// AuthorityNodeType defines different types of authority nodes in the network
type AuthorityNodeTypes string

// Enumeration of Authority Node Types
const (
	ElectedAuthorityNode AuthorityNodeType = "Elected Authority Node"
	MilitaryNode         AuthorityNodeType = "Military Node"
	BankingNode          AuthorityNodeType = "Banking Node"
	CentralBankNode      AuthorityNodeType = "Central Bank Node"
	ExchangeNode         AuthorityNodeType = "Exchange Node"
	GovernmentNode       AuthorityNodeType = "Government Node"
	RegulatorNode        AuthorityNodeType = "Regulator Node"
)


