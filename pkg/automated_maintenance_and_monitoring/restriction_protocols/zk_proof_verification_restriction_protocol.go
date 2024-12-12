package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/zkproof"
)

// Constants for ZK Proof verification
const (
	ZKProofVerificationInterval = 5 * time.Second // Interval to check for new ZK Proof verifications
	MaxPendingVerifications     = 100             // Maximum number of pending verifications allowed at once
)

// ZKProofVerificationAutomation manages the verification process for ZK Proofs
type ZKProofVerificationAutomation struct {
	consensusSystem  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	transactionMutex *sync.RWMutex
}

// NewZKProofVerificationAutomation initializes the automation for ZK Proof verification restrictions
func NewZKProofVerificationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, transactionMutex *sync.RWMutex) *ZKProofVerificationAutomation {
	return &ZKProofVerificationAutomation{
		consensusSystem:  consensusSystem,
		ledgerInstance:   ledgerInstance,
		transactionMutex: transactionMutex,
	}
}

// StartMonitoring continuously checks ZK Proof verifications across the network
func (automation *ZKProofVerificationAutomation) StartMonitoring() {
	ticker := time.NewTicker(ZKProofVerificationInterval)

	go func() {
		for range ticker.C {
			automation.evaluatePendingZKProofVerifications()
		}
	}()
}

// evaluatePendingZKProofVerifications checks pending ZK Proof verifications and restricts if limits are exceeded
func (automation *ZKProofVerificationAutomation) evaluatePendingZKProofVerifications() {
	automation.transactionMutex.Lock()
	defer automation.transactionMutex.Unlock()

	pendingVerifications := automation.consensusSystem.GetPendingZKProofVerifications()

	if len(pendingVerifications) > MaxPendingVerifications {
		automation.enforceVerificationLimit(len(pendingVerifications))
	} else {
		automation.processPendingVerifications(pendingVerifications)
	}
}

// enforceVerificationLimit restricts ZK Proof verifications if the pending list exceeds the allowed limit
func (automation *ZKProofVerificationAutomation) enforceVerificationLimit(pendingCount int) {
	err := automation.consensusSystem.RestrictZKProofVerifications()
	if err != nil {
		fmt.Println("Failed to enforce ZK Proof verification restriction:", err)
		return
	}

	// Log the restriction into the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("zkproof-verification-limit-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "ZK Proof Verification Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("ZK Proof verification restricted due to pending verifications exceeding limit (%d)", pendingCount),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log ZK Proof verification restriction:", err)
	} else {
		fmt.Println("ZK Proof verification restriction logged in the ledger.")
	}
}

// processPendingVerifications processes all pending ZK Proof verifications
func (automation *ZKProofVerificationAutomation) processPendingVerifications(verifications []zkproof.ZKProof) {
	for _, proof := range verifications {
		if !automation.verifyZKProof(proof) {
			automation.logFailedVerification(proof)
		} else {
			automation.logSuccessfulVerification(proof)
		}
	}
}

// verifyZKProof performs the actual ZK Proof verification process
func (automation *ZKProofVerificationAutomation) verifyZKProof(proof zkproof.ZKProof) bool {
	// Leverage the consensus system for verification
	isValid, err := automation.consensusSystem.VerifyZKProof(proof)
	if err != nil {
		fmt.Println("Error verifying ZK Proof:", err)
		return false
	}
	return isValid
}

// logFailedVerification logs any failed ZK Proof verifications into the ledger
func (automation *ZKProofVerificationAutomation) logFailedVerification(proof zkproof.ZKProof) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("zkproof-failed-verification-%s-%d", proof.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "ZK Proof Verification Failed",
		Status:    "Failed",
		Details:   fmt.Sprintf("ZK Proof verification failed for proof ID %s.", proof.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log failed ZK Proof verification:", err)
	} else {
		fmt.Println("Failed ZK Proof verification logged for proof ID:", proof.ID)
	}
}

// logSuccessfulVerification logs successful ZK Proof verifications into the ledger
func (automation *ZKProofVerificationAutomation) logSuccessfulVerification(proof zkproof.ZKProof) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("zkproof-successful-verification-%s-%d", proof.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "ZK Proof Verification Successful",
		Status:    "Verified",
		Details:   fmt.Sprintf("ZK Proof verification succeeded for proof ID %s.", proof.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log successful ZK Proof verification:", err)
	} else {
		fmt.Println("Successful ZK Proof verification logged for proof ID:", proof.ID)
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ZKProofVerificationAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
