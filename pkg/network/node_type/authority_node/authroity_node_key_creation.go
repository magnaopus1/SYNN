package authority_node

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
)

// NodeType represents the types of authority nodes.
const (
	AuthorityNode      = "AuthorityNode"
	BankNode           = "BankNode"
	CentralBankNode    = "CentralBankNode"
	CreditProviderNode = "CreditProviderNode"
	ElectedAuthority   = "ElectedAuthorityNode"
	ExchangeNode       = "ExchangeNode"
	GovernmentNode     = "GovernmentNode"
	MilitaryNode       = "MilitaryNode"
	RegulatorNode      = "RegulatorNode"
)

// AuthorityNodeKey represents a generated key for a specific node type.
type AuthorityNodeKey struct {
	KeyID          string    // Unique key identifier
	NodeType       string    // Type of node the key is associated with
	ExpirationDate time.Time // Key expiration date (30 months from creation)
	MaxNodes       int       // Maximum number of nodes this key can start (fixed at 50)
	UsedNodes      int       // Number of nodes started with this key
	IsExpired      bool      // Whether the key has expired
}

// KeyManager manages authority node key creation, usage, and expiration.
type KeyManager struct {
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	Ledger            *ledger.Ledger                // Reference to the ledger for storing key details
	EncryptionService *encryption.Encryption        // Encryption service for secure key management
	GeneratedKeys     map[string]*AuthorityNodeKey  // Map of generated keys by key ID
	GenesisCreated    bool                          // Flag to ensure genesis key creation happens only once
	KeyDisburser      *KeyDisburser                 // Handles disbursement and tracking of keys
}

// NewKeyManager initializes a new KeyManager.
func NewKeyManager(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, keyDisburser *KeyDisburser) *KeyManager {
	return &KeyManager{
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
		GeneratedKeys:     make(map[string]*AuthorityNodeKey),
		GenesisCreated:    false,
		KeyDisburser:      keyDisburser,
	}
}

// CreateAuthorityNodeKey generates a new key for the specified node type.
func (km *KeyManager) CreateAuthorityNodeKey(nodeType string) (*AuthorityNodeKey, error) {
	km.mutex.Lock()
	defer km.mutex.Unlock()

	if nodeType == "" {
		return nil, errors.New("invalid node type")
	}

	// Create a new key with a 30-month expiration period.
	key := &AuthorityNodeKey{
		KeyID:          common.GenerateUniqueID(),
		NodeType:       nodeType,
		ExpirationDate: time.Now().Add(30 * 24 * time.Hour * 30), // 30 months
		MaxNodes:       50,
		UsedNodes:      0,
		IsExpired:      false,
	}

	// Store the key in memory and record it in the ledger.
	km.GeneratedKeys[key.KeyID] = key
	err := km.Ledger.RecordKeyCreation(key)
	if err != nil {
		return nil, fmt.Errorf("failed to record key in ledger: %v", err)
	}

	return key, nil
}

// GenesisKeyCreation generates the initial set of authority node keys during blockchain genesis.
func (km *KeyManager) GenesisKeyCreation(ownerCSVPath string) error {
	km.mutex.Lock()
	defer km.mutex.Unlock()

	// Ensure this function is only used once.
	if km.GenesisCreated {
		return errors.New("genesis key creation has already been completed")
	}

	// Define the node types and quantities for genesis creation.
	nodeTypes := map[string]int{
		AuthorityNode:      5,
		BankNode:           10,
		CentralBankNode:    5,
		CreditProviderNode: 10,
		ElectedAuthority:   10,
		ExchangeNode:       5,
		GovernmentNode:     5,
		MilitaryNode:       0, // Not generating keys for military nodes.
		RegulatorNode:      5,
	}

	// Create a CSV file to store the keys.
	file, err := os.Create(ownerCSVPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	// Write the CSV headers.
	csvWriter.Write([]string{"KeyID", "NodeType", "ExpirationDate", "MaxNodes", "UsedNodes"})

	// Generate keys for each node type.
	for nodeType, count := range nodeTypes {
		for i := 0; i < count; i++ {
			key, err := km.CreateAuthorityNodeKey(nodeType)
			if err != nil {
				return fmt.Errorf("failed to create key for node type %s: %v", nodeType, err)
			}

			// Write the key details to the CSV file.
			csvWriter.Write([]string{
				key.KeyID,
				key.NodeType,
				key.ExpirationDate.Format(time.RFC3339),
				fmt.Sprintf("%d", key.MaxNodes),
				fmt.Sprintf("%d", key.UsedNodes),
			})
		}
	}

	// Mark genesis as created to prevent future use.
	km.GenesisCreated = true
	return nil
}

// UseAuthorityNodeKey marks a key as used for starting a new node.
func (km *KeyManager) UseAuthorityNodeKey(keyID string) error {
	km.mutex.Lock()
	defer km.mutex.Unlock()

	key, exists := km.GeneratedKeys[keyID]
	if !exists {
		return errors.New("key not found")
	}

	// Check if the key is expired.
	if time.Now().After(key.ExpirationDate) {
		key.IsExpired = true
		return errors.New("key has expired")
	}

	// Check if the key has already been used to the maximum limit.
	if key.UsedNodes >= key.MaxNodes {
		return errors.New("key has already been used to the maximum number of nodes")
	}

	// Increment the usage count.
	key.UsedNodes++

	// Update the ledger with the new key usage count.
	err := km.Ledger.UpdateKeyUsage(key)
	if err != nil {
		return fmt.Errorf("failed to update key usage in ledger: %v", err)
	}

	return nil
}

// ViewAuthorityNodeKey allows the user to view the details of a specific key.
func (km *KeyManager) ViewAuthorityNodeKey(keyID string) (*AuthorityNodeKey, error) {
	km.mutex.Lock()
	defer km.mutex.Unlock()

	key, exists := km.GeneratedKeys[keyID]
	if !exists {
		return nil, errors.New("key not found")
	}

	return key, nil
}

// RefreshAuthorityNodeKey refreshes a key after 30 months, allowing the owner to continue using it.
func (km *KeyManager) RefreshAuthorityNodeKey(keyID string) error {
	km.mutex.Lock()
	defer km.mutex.Unlock()

	key, exists := km.GeneratedKeys[keyID]
	if !exists {
		return errors.New("key not found")
	}

	// Extend the key's expiration date by another 30 months.
	key.ExpirationDate = time.Now().Add(30 * 24 * time.Hour * 30)
	key.IsExpired = false

	// Update the ledger with the new expiration date.
	err := km.Ledger.UpdateKeyExpiration(key)
	if err != nil {
		return fmt.Errorf("failed to update key expiration in ledger: %v", err)
	}

	fmt.Printf("Key %s has been successfully refreshed.\n", keyID)
	return nil
}

// DisburseKey disburses a new key for an accepted authority node proposal.
func (km *KeyManager) DisburseKey(proposal interface{}) error {
	km.mutex.Lock()
	defer km.mutex.Unlock()

	var nodeType, applicantWallet string

	switch p := proposal.(type) {
	case *AuthorityNodeProposal:
		nodeType = AuthorityNode
		applicantWallet = p.ApplicantWallet
	case *BankNodeProposal:
		nodeType = BankNode
		applicantWallet = p.ApplicantWallet
	case *CentralBankNodeProposal:
		nodeType = CentralBankNode
		applicantWallet = p.ApplicantWallet
	case *CreditProviderNodeProposal:
		nodeType = CreditProviderNode
		applicantWallet = p.ApplicantWallet
	case *ElectedAuthorityNodeProposal:
		nodeType = ElectedAuthority
		applicantWallet = p.ApplicantWallet
	case *ExchangeNodeProposal:
		nodeType = ExchangeNode
		applicantWallet = p.ApplicantWallet
	case *GovernmentNodeProposal:
		nodeType = GovernmentNode
		applicantWallet = p.ApplicantWallet
	case *MilitaryNodeProposal:
		nodeType = MilitaryNode
		applicantWallet = p.ApplicantWallet
	case *RegulatorNodeProposal:
		nodeType = RegulatorNode
		applicantWallet = p.ApplicantWallet
	default:
		return errors.New("invalid proposal type for key disbursement")
	}

	// Create a new key for the node type.
	key, err := km.CreateAuthorityNodeKey(nodeType)
	if err != nil {
		return fmt.Errorf("failed to create key for %s: %v", nodeType, err)
	}

	// Record the key disbursement in the ledger.
	err = km.Ledger.RecordKeyDisbursement(applicantWallet, key)
	if err != nil {
		return fmt.Errorf("failed to record key disbursement: %v", err)
	}

	fmt.Printf("Key successfully disbursed to %s for node type %s.\n", applicantWallet, nodeType)
	return nil
}
