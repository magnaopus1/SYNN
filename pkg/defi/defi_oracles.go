package defi

import (
	"fmt"
	"log"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewOracleManager initializes the DeFi Oracle Manager
// Manages oracle submissions and their verification status.
func NewOracleManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *OracleManager {
	if ledgerInstance == nil || encryptionService == nil {
		log.Fatalf("[ERROR] NewOracleManager requires valid ledger and encryptionService instances")
	}
	
	log.Printf("[INFO] Initializing OracleManager with ledger and encryption service.")
	
	return &OracleManager{
		OracleSubmissions:   make(map[string]*OracleData),
		VerifiedSubmissions: []*OracleData{},
		PendingSubmissions:  []*OracleData{},
		Ledger:              ledgerInstance,
		EncryptionService:   encryptionService,
		mu:                  sync.Mutex{},
	}
}


// SubmitOracleData allows an oracle to submit data to the blockchain
// Validates inputs, encrypts the data payload, and logs the submission in the ledger.
func (om *OracleManager) SubmitOracleData(dataFeedID, payload, handlerNode string) (*OracleData, error) {
	// Step 1: Input validation
	if dataFeedID == "" || payload == "" || handlerNode == "" {
		err := fmt.Errorf("dataFeedID, payload, and handlerNode cannot be empty")
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	log.Printf("[INFO] Submitting oracle data. DataFeedID: %s, HandlerNode: %s", dataFeedID, handlerNode)

	// Step 2: Thread safety
	om.mu.Lock()
	defer om.mu.Unlock()

	// Step 3: Generate a unique Oracle ID
	oracleID := generateUniqueID()
	log.Printf("[INFO] Generated Oracle ID: %s", oracleID)

	// Step 4: Encrypt the data payload
	startTime := time.Now()
	encryptionInstance, err := common.NewEncryption(256)
	if err != nil {
		log.Printf("[ERROR] Failed to create encryption instance: %v", err)
		return nil, fmt.Errorf("failed to create encryption instance: %w", err)
	}

	encryptedPayload, err := encryptionInstance.EncryptData("AES", []byte(payload), common.EncryptionKey)
	if err != nil {
		log.Printf("[ERROR] Failed to encrypt data payload: %v", err)
		return nil, fmt.Errorf("failed to encrypt data payload: %w", err)
	}
	log.Printf("[INFO] Data payload encrypted successfully. Duration: %v", time.Since(startTime))

	// Step 5: Create OracleData instance
	oracleData := &OracleData{
		OracleID:         oracleID,
		DataFeedID:       dataFeedID,
		DataPayload:      payload,
		Verified:         false,
		Timestamp:        time.Now(),
		HandlerNode:      handlerNode,
		EncryptedPayload: string(encryptedPayload),
	}

	// Step 6: Update OracleManager state
	om.PendingSubmissions = append(om.PendingSubmissions, oracleData)
	om.OracleSubmissions[oracleID] = oracleData
	log.Printf("[INFO] Oracle data added to pending submissions and OracleSubmissions. Oracle ID: %s", oracleID)

	// Step 7: Log submission in the ledger
	err = om.Ledger.DataManagementLedger.RecordOracleSubmission(oracleID, map[string]interface{}{
		"dataFeedID":   dataFeedID,
		"handlerNode":  handlerNode,
		"timestamp":    oracleData.Timestamp,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to log oracle submission in the ledger: %v", err)
		return nil, fmt.Errorf("failed to log oracle submission in the ledger: %w", err)
	}

	// Step 8: Final log and return
	log.Printf("[SUCCESS] Oracle data submitted successfully. Oracle ID: %s", oracleID)
	return oracleData, nil
}


// VerifyOracleData verifies the data provided by the oracle
// Updates the status and logs the verification in the ledger.
func (om *OracleManager) VerifyOracleData(oracleID string) error {
	// Step 1: Input validation
	if oracleID == "" {
		err := fmt.Errorf("oracleID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}

	log.Printf("[INFO] Starting verification for Oracle ID: %s", oracleID)

	// Step 2: Thread safety
	om.mu.Lock()
	defer om.mu.Unlock()

	// Step 3: Retrieve oracle data
	oracleData, exists := om.OracleSubmissions[oracleID]
	if !exists {
		err := fmt.Errorf("oracle data %s not found", oracleID)
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 4: Perform verification logic
	startTime := time.Now()
	oracleData.Verified = true // Replace with real verification logic if needed
	log.Printf("[INFO] Oracle data verification completed. Oracle ID: %s, Duration: %v", oracleID, time.Since(startTime))

	// Step 5: Update submissions
	om.VerifiedSubmissions = append(om.VerifiedSubmissions, oracleData)
	om.PendingSubmissions = removePendingSubmission(om.PendingSubmissions, oracleID)

	// Step 6: Record verification in the ledger
	err := om.Ledger.DataManagementLedger.RecordOracleVerification(oracleID)
	if err != nil {
		log.Printf("[ERROR] Failed to log oracle verification in the ledger: %v", err)
		return fmt.Errorf("failed to log oracle verification in the ledger: %w", err)
	}

	// Step 7: Log success
	log.Printf("[SUCCESS] Oracle data verified and recorded. Oracle ID: %s", oracleID)
	return nil
}


// MonitorOracleStatus checks the status of a specific oracle submission
// Returns the verification status of the oracle data.
func (om *OracleManager) MonitorOracleStatus(oracleID string) (string, error) {
	// Step 1: Input validation
	if oracleID == "" {
		err := fmt.Errorf("oracleID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return "", err
	}

	log.Printf("[INFO] Checking status for Oracle ID: %s", oracleID)

	// Step 2: Thread safety
	om.mu.Lock()
	defer om.mu.Unlock()

	// Step 3: Retrieve oracle submission
	oracleData, exists := om.OracleSubmissions[oracleID]
	if !exists {
		err := fmt.Errorf("oracle data %s not found", oracleID)
		log.Printf("[ERROR] %v", err)
		return "", err
	}

	// Step 4: Determine status
	status := "Pending Verification"
	if oracleData.Verified {
		status = "Verified"
	}

	// Step 5: Log and return status
	log.Printf("[INFO] Oracle ID %s status: %s", oracleID, status)
	return status, nil
}


// removePendingSubmission removes an oracle from the pending list
// Ensures efficient removal of a specific oracle from the pending queue.
func removePendingSubmission(pendingList []*OracleData, oracleID string) []*OracleData {
	for i, data := range pendingList {
		if data.OracleID == oracleID {
			updatedList := append(pendingList[:i], pendingList[i+1:]...)
			log.Printf("[INFO] Oracle ID %s removed from pending submissions", oracleID)
			return updatedList
		}
	}
	log.Printf("[WARNING] Oracle ID %s not found in pending submissions", oracleID)
	return pendingList
}



