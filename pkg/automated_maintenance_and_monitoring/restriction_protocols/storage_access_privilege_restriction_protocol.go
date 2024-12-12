package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

// Constants
const (
	StorageAccessCheckInterval    = 10 * time.Second // Interval to check access attempts
	MaxUnauthorizedAttempts       = 3
	UnauthorizedStorageAccess     = "Unauthorized storage access detected"
	ExcessiveStorageAccessWarning = "Excessive storage access attempts detected"
)

// StorageAccessRestrictionAutomation handles storage access privilege restrictions
type StorageAccessRestrictionAutomation struct {
	consensusSystem     *consensus.SynnergyConsensus
	ledgerInstance      *ledger.Ledger
	accessMutex         *sync.RWMutex
	accessAttempts      map[string]int    // Track access attempts per wallet
	restrictedWallets   map[string]bool   // Track restricted wallets
}

// NewStorageAccessRestrictionAutomation initializes the storage access restriction automation
func NewStorageAccessRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, accessMutex *sync.RWMutex) *StorageAccessRestrictionAutomation {
	return &StorageAccessRestrictionAutomation{
		consensusSystem:   consensusSystem,
		ledgerInstance:    ledgerInstance,
		accessMutex:       accessMutex,
		accessAttempts:    make(map[string]int),
		restrictedWallets: make(map[string]bool),
	}
}

// StartAccessMonitoring starts continuous monitoring of storage access privileges
func (automation *StorageAccessRestrictionAutomation) StartAccessMonitoring() {
	ticker := time.NewTicker(StorageAccessCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateStorageAccess()
		}
	}()
}

// evaluateStorageAccess checks for unauthorized or excessive storage access attempts and enforces restrictions
func (automation *StorageAccessRestrictionAutomation) evaluateStorageAccess() {
	automation.accessMutex.Lock()
	defer automation.accessMutex.Unlock()

	accessData := automation.consensusSystem.GetStorageAccessData()

	for walletID, accessCount := range accessData {
		if accessCount > MaxUnauthorizedAttempts {
			automation.logExcessiveAccessAttempts(walletID, accessCount)
			automation.restrictWallet(walletID)
		} else if automation.consensusSystem.IsUnauthorizedStorageAccess(walletID) {
			automation.logUnauthorizedAccess(walletID)
			automation.restrictWallet(walletID)
		}
	}
}

// logExcessiveAccessAttempts logs instances of excessive storage access attempts in the ledger
func (automation *StorageAccessRestrictionAutomation) logExcessiveAccessAttempts(walletID string, accessCount int) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("excessive-storage-access-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Excessive Storage Access",
		Status:    "Warning",
		Details:   fmt.Sprintf("Wallet %s attempted to access storage %d times, exceeding the limit.", walletID, accessCount),
	}

	// Encrypt the details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log excessive storage access attempts:", err)
	} else {
		fmt.Println("Excessive storage access logged for wallet:", walletID)
	}
}

// logUnauthorizedAccess logs unauthorized storage access in the ledger
func (automation *StorageAccessRestrictionAutomation) logUnauthorizedAccess(walletID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("unauthorized-storage-access-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Unauthorized Storage Access",
		Status:    "Critical",
		Details:   fmt.Sprintf("Unauthorized storage access detected for wallet %s.", walletID),
	}

	// Encrypt the unauthorized access details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log unauthorized storage access:", err)
	} else {
		fmt.Println("Unauthorized storage access logged for wallet:", walletID)
	}
}

// restrictWallet restricts storage access for wallets that breach access policies
func (automation *StorageAccessRestrictionAutomation) restrictWallet(walletID string) {
	fmt.Printf("Wallet %s has been restricted from storage access due to violations.\n", walletID)

	// Log the restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-restriction-storage-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Storage Access Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Wallet %s restricted from accessing storage due to policy violations.", walletID),
	}

	// Encrypt the restriction details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log wallet storage restriction:", err)
	} else {
		fmt.Println("Storage access restriction logged for:", walletID)
	}

	// Update the consensus system to restrict storage access for the wallet
	automation.consensusSystem.RestrictStorageAccess(walletID)
}

// encryptData encrypts sensitive information before storing in the ledger
func (automation *StorageAccessRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting access data:", err)
		return data
	}
	return string(encryptedData)
}
