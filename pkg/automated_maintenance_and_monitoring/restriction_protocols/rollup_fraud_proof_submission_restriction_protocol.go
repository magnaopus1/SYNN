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
	FraudProofCheckInterval      = 5 * time.Second  // Interval for checking fraud proof submissions
	MaxFraudProofSubmissionCount = 5                // Maximum allowed fraudulent proof submissions
)

// RollupFraudProofRestrictionAutomation monitors and restricts fraudulent proof submissions in rollups
type RollupFraudProofRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	fraudProofSubmissions  map[string]int // Tracks fraudulent proof submissions per user
}

// NewRollupFraudProofRestrictionAutomation initializes RollupFraudProofRestrictionAutomation
func NewRollupFraudProofRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RollupFraudProofRestrictionAutomation {
	return &RollupFraudProofRestrictionAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		fraudProofSubmissions: make(map[string]int),
	}
}

// StartFraudProofMonitoring begins the continuous monitoring of fraud proof submissions
func (automation *RollupFraudProofRestrictionAutomation) StartFraudProofMonitoring() {
	ticker := time.NewTicker(FraudProofCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorFraudProofSubmissions()
		}
	}()
}

// monitorFraudProofSubmissions checks for fraudulent proof submissions in rollup processes
func (automation *RollupFraudProofRestrictionAutomation) monitorFraudProofSubmissions() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch fraud proof submission data from Synnergy Consensus
	fraudProofData := automation.consensusSystem.GetFraudProofSubmissions()

	for userID, fraudStatus := range fraudProofData {
		// Check if a fraudulent proof has been submitted
		if fraudStatus == "fraudulent" {
			automation.flagFraudProofViolation(userID, "Fraudulent proof submission detected")
		}
	}
}

// flagFraudProofViolation flags a fraudulent proof submission and logs it in the ledger
func (automation *RollupFraudProofRestrictionAutomation) flagFraudProofViolation(userID string, reason string) {
	fmt.Printf("Fraudulent proof submission: User ID %s, Reason: %s\n", userID, reason)

	// Increment the fraudulent submission count for the user
	automation.fraudProofSubmissions[userID]++

	// Log the violation in the ledger
	automation.logFraudProofViolation(userID, reason)

	// Check if the user has exceeded the allowed number of fraudulent proof submissions
	if automation.fraudProofSubmissions[userID] >= MaxFraudProofSubmissionCount {
		automation.restrictFraudProofSubmission(userID)
	}
}

// logFraudProofViolation logs the fraudulent proof submission violation into the ledger
func (automation *RollupFraudProofRestrictionAutomation) logFraudProofViolation(userID string, violationReason string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("fraud-proof-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Fraud Proof Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated fraud proof restrictions. Reason: %s", userID, violationReason),
	}

	// Encrypt the violation log before adding it to the ledger
	encryptedDetails := automation.encryptFraudProofData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log fraud proof violation:", err)
	} else {
		fmt.Println("Fraud proof violation logged.")
	}
}

// restrictFraudProofSubmission restricts a user from submitting further fraud proofs after repeated violations
func (automation *RollupFraudProofRestrictionAutomation) restrictFraudProofSubmission(userID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("fraud-proof-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Fraud Proof Submission Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from submitting fraud proofs due to repeated violations.", userID),
	}

	// Encrypt the restriction details before logging it to the ledger
	encryptedDetails := automation.encryptFraudProofData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log fraud proof submission restriction:", err)
	} else {
		fmt.Println("Fraud proof submission restriction applied.")
	}
}

// encryptFraudProofData encrypts the fraud proof submission data before logging it to ensure privacy and security
func (automation *RollupFraudProofRestrictionAutomation) encryptFraudProofData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting fraud proof data:", err)
		return data
	}
	return string(encryptedData)
}
