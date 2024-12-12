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

const (
	AirdropCheckInterval      = 60 * time.Second // Interval for checking restricted airdrop eligibility
	MaxAllowedAirdropAttempts = 3                // Maximum allowed failed airdrop attempts before restriction
)

// RestrictedAirdropAutomation enforces airdrop restrictions based on network criteria
type RestrictedAirdropAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	airdropViolationCount  map[string]int // Tracks failed airdrop attempts per user
}

// NewRestrictedAirdropAutomation initializes and returns an instance of RestrictedAirdropAutomation
func NewRestrictedAirdropAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedAirdropAutomation {
	return &RestrictedAirdropAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		airdropViolationCount: make(map[string]int),
	}
}

// StartAirdropMonitoring starts continuous monitoring of restricted airdrop eligibility
func (automation *RestrictedAirdropAutomation) StartAirdropMonitoring() {
	ticker := time.NewTicker(AirdropCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorAirdropEligibility()
		}
	}()
}

// monitorAirdropEligibility checks for failed airdrop attempts and enforces restrictions if necessary
func (automation *RestrictedAirdropAutomation) monitorAirdropEligibility() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch airdrop data from Synnergy Consensus
	airdropData := automation.consensusSystem.GetAirdropData()

	for userID, airdropStatus := range airdropData {
		// Check if the user has violated airdrop eligibility rules
		if airdropStatus == "failed" {
			automation.flagAirdropViolation(userID, "Airdrop eligibility requirements not met")
		}
	}
}

// flagAirdropViolation flags a user's failed airdrop attempt and logs it in the ledger
func (automation *RestrictedAirdropAutomation) flagAirdropViolation(userID string, reason string) {
	fmt.Printf("Airdrop violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.airdropViolationCount[userID]++

	// Log the violation in the ledger
	automation.logAirdropViolation(userID, reason)

	// Check if the user has exceeded the allowed number of airdrop attempts
	if automation.airdropViolationCount[userID] >= MaxAllowedAirdropAttempts {
		automation.restrictAirdropAccess(userID)
	}
}

// logAirdropViolation logs the flagged airdrop violation into the ledger with full details
func (automation *RestrictedAirdropAutomation) logAirdropViolation(userID string, violationReason string) {
	// Create a ledger entry for airdrop violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("airdrop-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Airdrop Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated airdrop eligibility rules. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptAirdropData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log airdrop violation:", err)
	} else {
		fmt.Println("Airdrop violation logged.")
	}
}

// restrictAirdropAccess restricts a user's access to future airdrops after exceeding allowed attempts
func (automation *RestrictedAirdropAutomation) restrictAirdropAccess(userID string) {
	// Add restriction details to the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("airdrop-access-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Airdrop Access Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from participating in airdrops due to repeated violations.", userID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptAirdropData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log airdrop access restriction:", err)
	} else {
		fmt.Println("Airdrop access restriction applied.")
	}
}

// encryptAirdropData encrypts the airdrop data before logging for security
func (automation *RestrictedAirdropAutomation) encryptAirdropData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting airdrop data:", err)
		return data
	}
	return string(encryptedData)
}
