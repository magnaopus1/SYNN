package sidechains

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/consensus"
	"synnergy_network/pkg/network"
)

// NewSidechainSecurityManager initializes the security manager
func NewSidechainSecurityManager(consensus *common.SynnergyConsensus, encryptionService *encryption.Encryption, ledgerInstance *ledger.Ledger, sidechainNetwork *common.SidechainNetwork) *common.SidechainSecurityManager {
	return &common.SidechainSecurityManager{
		Consensus:        consensus,
		Encryption:       encryptionService,
		Ledger:           ledgerInstance,
		SidechainNetwork: sidechainNetwork,
	}
}

// ValidateTransaction ensures that a transaction follows the consensus rules and is securely transmitted
func (ssm *common.SidechainSecurityManager) ValidateTransaction(tx *common.Transaction) error {
	ssm.mu.Lock()
	defer ssm.mu.Unlock()

	// Ensure the transaction is valid according to Synnergy Consensus
	err := ssm.Consensus.ValidateTransaction(tx)
	if err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Log the successful validation in the ledger
	err = ssm.Ledger.RecordTransactionValidation(tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction validation: %v", err)
	}

	fmt.Printf("Transaction %s validated successfully\n", tx.TxID)
	return nil
}

// SecureTransaction encrypts the transaction data before broadcasting to the network
func (ssm *common.SidechainSecurityManager) SecureTransaction(tx *common.Transaction) error {
	ssm.mu.Lock()
	defer ssm.mu.Unlock()

	// Encrypt the transaction data
	encryptedData, err := ssm.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedData)

	// Log the encryption event in the ledger
	err = ssm.Ledger.RecordEncryption(tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log encryption event: %v", err)
	}

	fmt.Printf("Transaction %s encrypted and secured\n", tx.TxID)
	return nil
}

// BroadcastSecureData broadcasts encrypted transaction data to the sidechain network
func (ssm *common.SidechainSecurityManager) BroadcastSecureData(tx *common.Transaction, nodeID string) error {
	ssm.mu.Lock()
	defer ssm.mu.Unlock()

	// Encrypt the transaction data for secure transmission
	encryptedData, err := ssm.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}

	// Broadcast the encrypted transaction data to the target node
	node, err := ssm.SidechainNetwork.GetNode(nodeID)
	if err != nil {
		return fmt.Errorf("failed to find node: %v", err)
	}

	err = ssm.SidechainNetwork.BroadcastData(node, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction to node %s: %v", nodeID, err)
	}

	// Log the secure broadcast event in the ledger
	err = ssm.Ledger.RecordSecureBroadcast(tx.TxID, nodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log secure broadcast: %v", err)
	}

	fmt.Printf("Secure data broadcast for transaction %s to node %s\n", tx.TxID, nodeID)
	return nil
}

// MonitorSecurity ensures the security of the sidechain network through regular checks
func (ssm *common.SidechainSecurityManager) MonitorSecurity() error {
	ssm.mu.Lock()
	defer ssm.mu.Unlock()

	// Perform security checks (e.g., network integrity, consensus performance, encryption validity)
	fmt.Println("Running sidechain security monitoring...")

	// Check the consensus mechanism for any performance issues
	err := ssm.Consensus.CheckIntegrity()
	if err != nil {
		return fmt.Errorf("consensus integrity check failed: %v", err)
	}

	// Check if all encryption keys and methods are functioning as expected
	err = ssm.Encryption.CheckEncryptionHealth()
	if err != nil {
		return fmt.Errorf("encryption health check failed: %v", err)
	}

	// Perform node security checks via the sidechain network
	err = ssm.SidechainNetwork.CheckNodeSecurity()
	if err != nil {
		return fmt.Errorf("node security check failed: %v", err)
	}

	// Log the security monitoring event in the ledger
	err = ssm.Ledger.RecordSecurityCheck(time.Now())
	if err != nil {
		return fmt.Errorf("failed to log security monitoring: %v", err)
	}

	fmt.Println("Sidechain security monitoring completed successfully")
	return nil
}

// HandleSecurityIncident handles and logs a detected security incident on the sidechain
func (ssm *common.SidechainSecurityManager) HandleSecurityIncident(incidentDescription string) error {
	ssm.mu.Lock()
	defer ssm.mu.Unlock()

	// Log the security incident in the ledger
	err := ssm.Ledger.RecordSecurityIncident(incidentDescription, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log security incident: %v", err)
	}

	fmt.Printf("Security incident handled: %s\n", incidentDescription)
	return nil
}

// RotateEncryptionKeys rotates encryption keys to enhance security
func (ssm *common.SidechainSecurityManager) RotateEncryptionKeys() error {
	ssm.mu.Lock()
	defer ssm.mu.Unlock()

	// Rotate encryption keys via the encryption service
	err := ssm.Encryption.RotateKeys()
	if err != nil {
		return fmt.Errorf("failed to rotate encryption keys: %v", err)
	}

	// Log the key rotation event
	err = ssm.Ledger.RecordKeyRotation(time.Now())
	if err != nil {
		return fmt.Errorf("failed to log encryption key rotation: %v", err)
	}

	fmt.Println("Encryption keys rotated successfully")
	return nil
}
