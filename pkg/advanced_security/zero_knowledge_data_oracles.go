package advanced_security

import (
	"fmt"
	"log"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// ZKOracleData represents the data provided by the oracle with zk-proofs
type ZKOracleData struct {
	OracleID        string    // Unique ID of the oracle
	DataFeedID      string    // Data feed ID for the specific oracle
	ZKProof         string    // Zero-Knowledge Proof of the data validity
	DataPayload     string    // Encrypted data payload
	Verified        bool      // Whether the zk-proof has been verified
	Timestamp       time.Time // Timestamp of the oracle data submission
	HandlerNode     string    // Node handling the oracle submission
}

// ZKDOManager manages zero-knowledge data oracles (ZKDO)
type ZKDOManager struct {
	OracleSubmissions    map[string]*ZKOracleData  // Active oracle submissions
	CompletedSubmissions []*ZKOracleData           // Log of verified submissions
	PendingVerifications []*ZKOracleData           // Queue of pending zk-proof verifications
	Ledger               *ledger.Ledger            // Ledger instance for recording zk-proof data
	EncryptionService    *common.Encryption    // Encryption service for secure data handling
	mu                   sync.Mutex                // Mutex for concurrent operations
	Consensus 			common.SynnergyConsensus
}


// NewZKDOManager initializes the Zero-Knowledge Data Oracles manager
func NewZKDOManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *ZKDOManager {
	return &ZKDOManager{
		OracleSubmissions:   make(map[string]*ZKOracleData),
		CompletedSubmissions: []*ZKOracleData{},
		PendingVerifications: []*ZKOracleData{},
		Ledger:               ledgerInstance,
		EncryptionService:    encryptionService,
	}
}

// SubmitOracleData allows an oracle to submit data with a zk-proof
func (zkdo *ZKDOManager) SubmitOracleData(dataFeedID, payload, handlerNode string) (*ZKOracleData, error) {
    zkdo.mu.Lock()
    defer zkdo.mu.Unlock()

    // Generate a unique Oracle ID
    oracleID := generateUniqueID()

    // Encrypt the data payload using AES
    encryptedPayload, err := zkdo.EncryptionService.EncryptData("AES", []byte(payload), common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt payload: %v", err)
    }

    // Create an instance of ZkProof to generate the proof
    zkProofInstance := &common.ZkProof{
        ProofID:    generateUniqueID(),
        ProverID:   handlerNode,
        Ledger:     zkdo.Ledger,             // Pass the ledger to the zkProof instance
        Encryption: zkdo.EncryptionService,  // Use the same encryption service
    }

    // Generate a zero-knowledge proof for the payload
    err = zkProofInstance.GenerateProof([]byte(payload))
    if err != nil {
        return nil, fmt.Errorf("failed to generate zero-knowledge proof: %v", err)
    }

    // Convert encryptedPayload (which is []byte) to string for DataPayload
    encryptedPayloadStr := string(encryptedPayload)

    // Convert zkProofInstance.ProofData from []byte to string
    zkProofStr := string(zkProofInstance.ProofData)

    // Create ZKOracleData instance
    oracleData := &ZKOracleData{
        OracleID:    oracleID,
        DataFeedID:  dataFeedID,
        ZKProof:     zkProofStr,              // Use the string version of the zk-proof data
        DataPayload: encryptedPayloadStr,     // Encrypted data as string
        Verified:    false,
        Timestamp:   time.Now(),
        HandlerNode: handlerNode,
    }

    // Add to pending verifications and oracle submissions
    zkdo.PendingVerifications = append(zkdo.PendingVerifications, oracleData)
    zkdo.OracleSubmissions[oracleID] = oracleData

    // Log the oracle submission in the ledger
    dataMap := map[string]interface{}{
        "DataFeedID":   dataFeedID,
        "HandlerNode":  handlerNode,
        "Timestamp":    time.Now(),
    }
    successMsg, err := zkdo.Ledger.DataManagementLedger.RecordOracleSubmission(oracleID, dataMap)
    if err != nil {
        return nil, fmt.Errorf("failed to log oracle submission in the ledger: %v", err)
    }

    // Optionally log the success message from the ledger
    log.Printf("Ledger Record Success: %s", successMsg)

    fmt.Printf("Oracle data submitted for data feed %s by node %s\n", dataFeedID, handlerNode)
    return oracleData, nil
}



// VerifyZKProof verifies the zk-proof and marks the oracle data as verified
func (zkdo *ZKDOManager) VerifyZKProof(oracleID string) error {
    zkdo.mu.Lock()
    defer zkdo.mu.Unlock()

    // Find the oracle data
    oracleData, exists := zkdo.OracleSubmissions[oracleID]
    if !exists {
        return fmt.Errorf("oracle data %s not found", oracleID)
    }

    // Create or retrieve the zk-proof instance (convert ZKProof from string to []byte)
    zkProofInstance := &common.ZkProof{
        ProofID:    oracleData.OracleID,
        ProofData:  []byte(oracleData.ZKProof), // Convert string to []byte
        Ledger:     zkdo.Ledger,                // Reference to ledger
        Encryption: zkdo.EncryptionService,     // Reference to encryption service
        Consensus:  &zkdo.Consensus,            // Pass a pointer to the SynnergyConsensus instance
        VerifierID: oracleData.HandlerNode,     // Assuming handlerNode is the verifier
    }

    // Verify the zk-proof (real-world proof verification)
    err := zkProofInstance.VerifyProof([]byte(oracleData.DataPayload)) // Pass the expected data (DataPayload)
    if err != nil {
        oracleData.Verified = false
        return fmt.Errorf("zero-knowledge proof verification failed for oracle %s: %v", oracleID, err)
    }

    // Mark as verified
    oracleData.Verified = true

    // Move from pending to completed submissions
    zkdo.CompletedSubmissions = append(zkdo.CompletedSubmissions, oracleData)
    zkdo.PendingVerifications = removePendingVerification(zkdo.PendingVerifications, oracleID)

    // Log the verified submission in the ledger (pass only oracleID)
    err = zkdo.Ledger.DataManagementLedger.RecordOracleVerification(oracleID)
    if err != nil {
        return fmt.Errorf("failed to log oracle verification in the ledger: %v", err)
    }

    fmt.Printf("Oracle data %s verified with zk-proof\n", oracleID)
    return nil
}


// MonitorOracleStatus allows the system to check the verification status of a specific oracle submission
func (zkdo *ZKDOManager) MonitorOracleStatus(oracleID string) (string, error) {
	zkdo.mu.Lock()
	defer zkdo.mu.Unlock()

	// Retrieve the oracle submission by ID
	oracleData, exists := zkdo.OracleSubmissions[oracleID]
	if !exists {
		return "", fmt.Errorf("oracle data %s not found", oracleID)
	}

	if oracleData.Verified {
		return "Verified", nil
	}
	return "Pending Verification", nil
}


// removePendingVerification removes an oracle from the pending list
func removePendingVerification(pendingList []*ZKOracleData, oracleID string) []*ZKOracleData {
	for i, data := range pendingList {
		if data.OracleID == oracleID {
			return append(pendingList[:i], pendingList[i+1:]...)
		}
	}
	return pendingList
}
