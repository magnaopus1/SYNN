package common

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)

// SpaceTimeProof represents a proof of space and time for validating data storage and retrieval over time
type SpaceTimeProof struct {
	ProofID        string                         // Unique identifier for the proof
	DataHash       string                         // Hash of the data being proven
	SpaceUsed      uint64                         // Amount of space used (in bytes)
	TimeStored     time.Duration                  // Duration for which the data was stored
	IsValid        bool                           // Whether the proof has been validated
	Ledger         *ledger.Ledger                 // Reference to the ledger for recording events
	Encryption     *Encryption         // Encryption service for securing proof data
	Consensus      *SynnergyConsensus // Synnergy Consensus mechanism for validation
	mu             sync.Mutex                     // Mutex for concurrency handling
}

// NewSpaceTimeProof initializes a new space-time proof
func NewSpaceTimeProof(proofID string, dataHash string, spaceUsed uint64, timeStored time.Duration, ledgerInstance *ledger.Ledger, encryptionService *Encryption, consensus *SynnergyConsensus) *SpaceTimeProof {
	return &SpaceTimeProof{
		ProofID:    proofID,
		DataHash:   dataHash,
		SpaceUsed:  spaceUsed,
		TimeStored: timeStored,
		IsValid:    false,
		Ledger:     ledgerInstance,
		Encryption: encryptionService,
		Consensus:  consensus,
	}
}

// ValidateProof validates the space-time proof using Synnergy Consensus
func (stp *SpaceTimeProof) ValidateProof() error {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	if stp.IsValid {
		return errors.New("proof has already been validated")
	}

	// Simulate space-time validation using Synnergy Consensus
	err := stp.Consensus.ValidateSpaceTimeProof(stp.ProofID, stp.DataHash, stp.SpaceUsed, stp.TimeStored)
	if err != nil {
		return fmt.Errorf("proof validation failed: %v", err)
	}

	stp.IsValid = true

	// Log the proof validation in the ledger with the expected number of arguments (ProofID, Validator)
	validator := "synnergy_validator" // Use the actual validator ID here if you have it.
	err = stp.Ledger.StorageLedger.RecordProofValidation(stp.ProofID, validator)
	if err != nil {
		return fmt.Errorf("failed to log proof validation: %v", err)
	}

	// Additional logging can be done separately for SpaceUsed, TimeStored, and Timestamp
	fmt.Printf("Space-Time Proof %s validated successfully with SpaceUsed: %d, TimeStored: %s\n", 
		stp.ProofID, stp.SpaceUsed, stp.TimeStored)

	return nil
}


// RetrieveProofData retrieves the encrypted proof data for validation or auditing
func (stp *SpaceTimeProof) RetrieveProofData() (string, error) {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	if !stp.IsValid {
		return "", errors.New("proof is not yet validated")
	}

	// Define a random IV (ensure it's 16 bytes for AES encryption or adjust based on your encryption method)
	iv := []byte("randomIV-16bytes") // Example: make sure the IV fits the required length of the encryption method

	// Encrypt proof data for secure retrieval
	encryptedData, err := stp.Encryption.EncryptData("encryption_key_string", []byte(stp.DataHash), iv)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt proof data: %v", err)
	}

	fmt.Printf("Encrypted proof data for proof %s retrieved\n", stp.ProofID)
	return string(encryptedData), nil
}


// RecordStorageEvent records an event when data is stored in relation to a space-time proof
func (stp *SpaceTimeProof) RecordStorageEvent(storageDuration time.Duration, storageSpace uint64) error {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	// Update the proof's space and time data
	stp.TimeStored += storageDuration
	stp.SpaceUsed += storageSpace

	// Convert the non-string values to strings
	storageSpaceStr := fmt.Sprintf("%d", storageSpace)                      // Convert uint64 to string
	storageDurationStr := storageDuration.String()                          // Convert time.Duration to string
	timestampStr := time.Now().Format(time.RFC3339)                         // Format time.Time to string

	// Log the storage event in the ledger
	err := stp.Ledger.StorageLedger.RecordStorageEvent(stp.ProofID, storageSpaceStr, storageDurationStr, timestampStr)
	if err != nil {
		return fmt.Errorf("failed to log storage event: %v", err)
	}

	fmt.Printf("Storage event recorded for proof %s: %d bytes stored for %v\n", stp.ProofID, storageSpace, storageDuration)
	return nil
}

// RevalidateProof revalidates the space-time proof to ensure consistency over a long duration
func (stp *SpaceTimeProof) RevalidateProof() error {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	if !stp.IsValid {
		return errors.New("proof has not been validated before")
	}

	// Simulate revalidation of the proof
	err := stp.Consensus.RevalidateSpaceTimeProof(stp.ProofID, stp.DataHash, stp.SpaceUsed, stp.TimeStored)
	if err != nil {
		return fmt.Errorf("proof revalidation failed: %v", err)
	}

	// Convert the non-string values to strings for the ledger
	spaceUsedStr := fmt.Sprintf("%d", stp.SpaceUsed) // Convert uint64 to string

	// Log the revalidation event in the ledger (only two arguments: ProofID and spaceUsedStr)
	err = stp.Ledger.StorageLedger.RecordProofRevalidation(stp.ProofID, spaceUsedStr)
	if err != nil {
		return fmt.Errorf("failed to log proof revalidation: %v", err)
	}

	fmt.Printf("Space-Time Proof %s revalidated successfully\n", stp.ProofID)
	return nil
}


// InvalidateProof invalidates the proof if a violation or inconsistency is detected
func (stp *SpaceTimeProof) InvalidateProof() error {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	if !stp.IsValid {
		return errors.New("proof is not valid, cannot invalidate")
	}

	stp.IsValid = false

	// Convert the timestamp to string format
	timestampStr := time.Now().Format(time.RFC3339)          // Convert time.Time to string

	// Log the invalidation in the ledger
	err := stp.Ledger.StorageLedger.RecordProofInvalidation(stp.ProofID, timestampStr)
	if err != nil {
		return fmt.Errorf("failed to log proof invalidation: %v", err)
	}

	fmt.Printf("Space-Time Proof %s invalidated\n", stp.ProofID)
	return nil
}
