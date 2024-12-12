package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewGovernanceTracking initializes a new GovernanceTracking system
func NewGovernanceTracking(ledgerInstance *ledger.Ledger) *GovernanceTracking {
    return &GovernanceTracking{
        ProposalHistory: make(map[string]*GovernanceProposalStatus),
        LedgerInstance:  ledgerInstance,
    }
}

// TrackProposal creates an entry to track the governance proposal
func (gt *GovernanceTracking) TrackProposal(proposalID, encryptedDetails string) error {
    gt.mutex.Lock()
    defer gt.mutex.Unlock()

    if _, exists := gt.ProposalHistory[proposalID]; exists {
        return fmt.Errorf("proposal %s is already being tracked", proposalID)
    }

    newStatus := &GovernanceProposalStatus{
        ProposalID:      proposalID,
        Status:          "Pending",
        Timestamps:      []time.Time{time.Now()},
        EncryptedDetails: encryptedDetails,
    }

    gt.ProposalHistory[proposalID] = newStatus
    gt.logProposalStatusToLedger(proposalID, newStatus)

    fmt.Printf("Tracking started for proposal %s.\n", proposalID)
    return nil
}

// UpdateProposalStatus updates the status of a governance proposal and logs the update
func (gt *GovernanceTracking) UpdateProposalStatus(proposalID, newStatus string) error {
    gt.mutex.Lock()
    defer gt.mutex.Unlock()

    status, exists := gt.ProposalHistory[proposalID]
    if !exists {
        return fmt.Errorf("proposal %s not found", proposalID)
    }

    status.Status = newStatus
    status.Timestamps = append(status.Timestamps, time.Now())
    gt.logProposalStatusToLedger(proposalID, status)

    fmt.Printf("Updated proposal %s status to: %s\n", proposalID, newStatus)
    return nil
}

// GenerateGovernanceReport generates an encrypted report of all governance activities
func (gt *GovernanceTracking) GenerateGovernanceReport() (string, error) {
    gt.mutex.Lock()
    defer gt.mutex.Unlock()

    report := "Governance Report:\n"
    for proposalID, status := range gt.ProposalHistory {
        report += fmt.Sprintf("Proposal %s: Status: %s\n", proposalID, status.Status)
        for i, timestamp := range status.Timestamps {
            report += fmt.Sprintf("  Timestamp %d: %s\n", i+1, timestamp.String())
        }
    }

    // Create an Encryption instance
    encryptionInstance := &common.Encryption{}

    // Encrypt the report
    encryptedReport, err := encryptionInstance.EncryptData("AES", []byte(report), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt governance report: %v", err)
    }

    fmt.Println("Governance report generated and encrypted.")
    return string(encryptedReport), nil // Convert the encrypted []byte to string
}


// logProposalStatusToLedger logs the proposal status to the ledger
func (gt *GovernanceTracking) logProposalStatusToLedger(proposalID string, status *GovernanceProposalStatus) error {
    logEntry := fmt.Sprintf("Proposal %s status: %s", proposalID, status.Status)

    // Create an Encryption instance
    encryptionInstance := &common.Encryption{}

    // Encrypt the log entry
    encryptedLog, err := encryptionInstance.EncryptData("AES", []byte(logEntry), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt proposal status log: %v", err)
    }

    hash := gt.generateProposalHash(proposalID)

    // Log the encrypted status to the ledger, providing a dummy creationFee (e.g., 0.0)
    err = gt.LedgerInstance.RecordProposal(hash, string(encryptedLog), proposalID, 0.0) // Added creationFee as 0.0
    if err != nil {
        return fmt.Errorf("failed to log proposal status to ledger: %v", err)
    }

    fmt.Printf("Status of proposal %s logged to ledger.\n", proposalID)
    return nil
}



// QueryProposalStatus allows querying the current status of a proposal
func (gt *GovernanceTracking) QueryProposalStatus(proposalID string) (*GovernanceProposalStatus, error) {
    gt.mutex.Lock()
    defer gt.mutex.Unlock()

    status, exists := gt.ProposalHistory[proposalID]
    if !exists {
        return nil, fmt.Errorf("proposal %s not found", proposalID)
    }

    return status, nil
}

// generateProposalHash generates a SHA-256 hash for a proposal based on its ID.
func (gt *GovernanceTracking) generateProposalHash(proposalID string) string {
    // Step 1: Create a new SHA-256 hasher
    hasher := sha256.New()

    // Step 2: Write the proposal ID into the hasher (as bytes)
    hasher.Write([]byte(proposalID))

    // Step 3: Compute the hash
    hashBytes := hasher.Sum(nil)

    // Step 4: Convert the hash bytes to a hex string
    return hex.EncodeToString(hashBytes)
}