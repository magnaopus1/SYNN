package identity_services

import (
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network/node_type/authority_node"
	"time"
)

// Permission types for various authority nodes
type Permission string

const (
    TransactionCancellationApproval Permission = "TransactionCancellationApproval"
    EvidenceViewing                 Permission = "EvidenceViewing"
    TransactionReversalApproval     Permission = "TransactionReversalApproval"
    GlobalPermission                Permission = "GlobalPermission"
)

// AuthorityNodeType defines different types of authority nodes in the network
type AuthorityNodeTypes string

// Enumeration of Authority Node Types
const (
	ElectedAuthorityNode authority_node.AuthorityNodeTypes = "Elected Authority Node"
	MilitaryNode         authority_node.AuthorityNodeTypes = "Military Node"
	BankingNode          authority_node.AuthorityNodeTypes = "Banking Node"
	CentralBankNode      authority_node.AuthorityNodeTypes = "Central Bank Node"
	ExchangeNode         authority_node.AuthorityNodeTypes = "Exchange Node"
	GovernmentNode       authority_node.AuthorityNodeTypes = "Government Node"
	RegulatorNode        authority_node.AuthorityNodeTypes = "Regulator Node"
)

// Assuming this struct exists in the authority_node package
type AuthorityNodeStruct struct {
    NodeID            string
    AuthorityNodeType authority_node.AuthorityNodeTypes
    SecretKey         string
    EncryptedKey      string
    CreatedAt         time.Time
    Permissions       []Permission  // Add the Permissions field

}



// NewAccessControlManager initializes a new AccessControlManager
func NewAccessControlManager(ledgerInstance *ledger.Ledger, ownerSecret string) *AccessControlManager {
    return &AccessControlManager{
        AuthorityNodes:   make(map[string]*AuthorityNodeStruct),  // Corrected map type
        UserAccess:       make(map[string]string),
        StakeholderAccess: make(map[string]string),
        OwnerAccess:      ownerSecret,
        LedgerInstance:   ledgerInstance,
        Encryption:       &common.Encryption{},  // Initialize the encryption instance
    }
}

// AddAuthorityNode adds a new authority node to the system with an encrypted secret key
func (acm *AccessControlManager) AddAuthorityNode(nodeID string, nodeType authority_node.AuthorityNodeTypes, secretKey string) error {
    acm.mutex.Lock()
    defer acm.mutex.Unlock()

    // Check if the node already exists
    if _, exists := acm.AuthorityNodes[nodeID]; exists {
        return fmt.Errorf("authority node with ID %s already exists", nodeID)
    }

    // Encrypt the secret key using the acm.Encryption instance
    encryptedKey, err := acm.Encryption.EncryptData("AES", []byte(secretKey), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt node secret key: %v", err)
    }

    // Create a new AuthorityNodeStruct instance (assuming this struct represents the authority node)
    newNode := &AuthorityNodeStruct{
        NodeID:            nodeID,
        AuthorityNodeType: nodeType,  // Use the enum or string directly
        SecretKey:         secretKey,
        CreatedAt:         time.Now(),
        EncryptedKey:      string(encryptedKey),  // Convert []byte to string
    }

    // Add the new node to the AuthorityNodes map
    acm.AuthorityNodes[nodeID] = newNode

    // Record the node addition in the ledger
    acm.LedgerInstance.RecordAccessChange(fmt.Sprintf("Authority Node %s added with type %s", nodeID, nodeType), string(encryptedKey))

    fmt.Printf("Authority Node %s of type %s added successfully.\n", nodeID, nodeType)
    return nil
}





// VerifyNodeAccess verifies if a node can access the system with the given secret key
func (acm *AccessControlManager) VerifyNodeAccess(nodeID, secretKey string) bool {
    acm.mutex.Lock()
    defer acm.mutex.Unlock()

    node, exists := acm.AuthorityNodes[nodeID]
    if !exists {
        return false
    }

    return node.SecretKey == secretKey
}

// GrantUserAccess grants user access by encrypting the access level
func (acm *AccessControlManager) GrantUserAccess(userID, accessLevel string) error {
    acm.mutex.Lock()
    defer acm.mutex.Unlock()

    // Encrypt the access level using the Encryption instance
    encryptedAccess, err := acm.Encryption.EncryptData("AES", []byte(accessLevel), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt user access level: %v", err)
    }

    // Store the encrypted access level
    acm.UserAccess[userID] = string(encryptedAccess)

    // Record user access grant in the ledger (do not assign it to err)
    acm.LedgerInstance.RecordAccessChange(userID, fmt.Sprintf("User %s granted access level %s", userID, accessLevel))

    fmt.Printf("User %s granted access level: %s.\n", userID, accessLevel)
    return nil
}

// GrantStakeholderAccess grants access to a stakeholder by encrypting their access level
func (acm *AccessControlManager) GrantStakeholderAccess(stakeholderID, accessLevel string) error {
    acm.mutex.Lock()
    defer acm.mutex.Unlock()

    // Encrypt the access level using the Encryption instance
    encryptedAccess, err := acm.Encryption.EncryptData("AES", []byte(accessLevel), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt stakeholder access level: %v", err)
    }

    // Store the encrypted access level
    acm.StakeholderAccess[stakeholderID] = string(encryptedAccess)

    // Record stakeholder access grant in the ledger (do not assign it to err)
    acm.LedgerInstance.RecordAccessChange(stakeholderID, fmt.Sprintf("Stakeholder %s granted access level %s", stakeholderID, accessLevel))

    fmt.Printf("Stakeholder %s granted access level: %s.\n", stakeholderID, accessLevel)
    return nil
}


// VerifyOwnerAccess verifies if the owner can access the system with the given secret key
func (acm *AccessControlManager) VerifyOwnerAccess(secretKey string) bool {
    acm.mutex.Lock()
    defer acm.mutex.Unlock()

    return acm.OwnerAccess == secretKey
}

// ModifyAuthorityNode modifies the permissions of an authority node
func (acm *AccessControlManager) ModifyAuthorityNode(nodeID string, newPermissions []Permission) error {
    acm.mutex.Lock()
    defer acm.mutex.Unlock()

    node, exists := acm.AuthorityNodes[nodeID]
    if !exists {
        return fmt.Errorf("authority node %s not found", nodeID)
    }

    // Update the permissions
    node.Permissions = newPermissions

    // Log modification in the ledger (no assignment needed)
    acm.LedgerInstance.RecordAccessChange(nodeID, fmt.Sprintf("Authority Node %s modified permissions", nodeID))

    fmt.Printf("Authority Node %s permissions modified successfully.\n", nodeID)
    return nil
}

// RevokeAuthorityNodeAccess revokes the access of an authority node
func (acm *AccessControlManager) RevokeAuthorityNodeAccess(nodeID string) error {
    acm.mutex.Lock()
    defer acm.mutex.Unlock()

    if _, exists := acm.AuthorityNodes[nodeID]; !exists {
        return fmt.Errorf("authority node %s not found", nodeID)
    }

    // Remove the node from the map
    delete(acm.AuthorityNodes, nodeID)

    // Log the revocation in the ledger (no assignment needed)
    acm.LedgerInstance.RecordAccessChange(nodeID, fmt.Sprintf("Authority Node %s access revoked", nodeID))

    fmt.Printf("Authority Node %s access revoked successfully.\n", nodeID)
    return nil
}

// ViewPermissions returns the permissions for an authority node
func (acm *AccessControlManager) ViewPermissions(nodeID string) ([]Permission, error) {
    acm.mutex.Lock()
    defer acm.mutex.Unlock()

    node, exists := acm.AuthorityNodes[nodeID]
    if !exists {
        return nil, fmt.Errorf("authority node %s not found", nodeID)
    }

    return node.Permissions, nil
}

