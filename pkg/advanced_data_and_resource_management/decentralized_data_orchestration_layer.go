package advanced_data_and_resource_management

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

func NewDDOLManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *DDOLManager {
	if ledgerInstance == nil {
		panic("Ledger instance is required for DDOLManager initialization")
	}
	if encryptionService == nil {
		panic("Encryption service is required for DDOLManager initialization")
	}

	return &DDOLManager{
		Orchestrations:    make(map[string]*OrchestratedData),
		CompletedData:     []*OrchestratedData{},
		PendingQueue:      []string{},
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}



func (ddol *DDOLManager) RequestOrchestration(appID, payload, handlerNode string) (*OrchestratedData, error) {
	ddol.mu.Lock()
	defer ddol.mu.Unlock()

	if appID == "" || handlerNode == "" {
		return nil, errors.New("application ID and handler node must be provided")
	}
	if len(payload) == 0 {
		return nil, errors.New("payload cannot be empty")
	}

	// Generate a unique ID for the request
	dataID := generateUniqueID()

	// Encrypt the payload
	encryptedPayload, err := ddol.EncryptionService.EncryptData("AES", []byte(payload), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	encryptedPayloadStr := base64.StdEncoding.EncodeToString(encryptedPayload)

	orchestratedData := &OrchestratedData{
		DataID:            dataID,
		AppID:             appID,
		DataPayload:       payload,
		OrchestrationTime: time.Now(),
		Status:            "Pending",
		EncryptedPayload:  encryptedPayloadStr,
		HandlerNode:       handlerNode,
	}

	// Add the request to the pending queue
	ddol.PendingQueue = append(ddol.PendingQueue, dataID)
	ddol.Orchestrations[dataID] = orchestratedData

	// Log the orchestration request in the ledger
	logEntry := map[string]interface{}{
		"appID":            appID,
		"handlerNode":      handlerNode,
		"orchestrationTime": orchestratedData.OrchestrationTime,
		"encryptedPayload":  encryptedPayloadStr,
	}
	if _, err := ddol.Ledger.UtilityLedger.RecordOrchestrationRequest(dataID, logEntry); err != nil {
		return nil, fmt.Errorf("failed to log orchestration request: %w", err)
	}

	log.Printf("Orchestration request created for appID %s with data ID %s", appID, dataID)
	return orchestratedData, nil
}




func (ddol *DDOLManager) ProcessOrchestration() error {
	ddol.mu.Lock()
	defer ddol.mu.Unlock()

	if len(ddol.PendingQueue) == 0 {
		return errors.New("no pending orchestration requests")
	}

	dataID := ddol.PendingQueue[0]
	ddol.PendingQueue = ddol.PendingQueue[1:]

	orchestratedData, exists := ddol.Orchestrations[dataID]
	if !exists {
		return fmt.Errorf("orchestration data %s not found", dataID)
	}

	// Simulate processing by fetching off-chain data
	externalData, err := ddol.fetchExternalData(orchestratedData.HandlerNode)
	if err != nil {
		return fmt.Errorf("failed to fetch external data for data ID %s: %w", dataID, err)
	}

	// Validate and process data
	processedData, err := ddol.validateAndTransform(orchestratedData.DataPayload, externalData)
	if err != nil {
		return fmt.Errorf("data processing failed for data ID %s: %w", dataID, err)
	}

	orchestratedData.Status = "Processed"
	orchestratedData.Result = processedData
	orchestratedData.OrchestrationTime = time.Now()

	ddol.CompletedData = append(ddol.CompletedData, orchestratedData)

	// Log the completion in the ledger
	completionDetails := map[string]interface{}{
		"handlerNode":        orchestratedData.HandlerNode,
		"orchestrationTime":  orchestratedData.OrchestrationTime,
		"status":             orchestratedData.Status,
		"result":             orchestratedData.Result,
	}
	if err := ddol.Ledger.UtilityLedger.RecordOrchestrationCompletion(dataID, completionDetails); err != nil {
		return fmt.Errorf("failed to log orchestration completion: %w", err)
	}

	log.Printf("Orchestration completed for data ID %s", dataID)
	return nil
}

func (ddol *DDOLManager) fetchExternalData(handlerNode string) (string, error) {
	// Replace with actual logic to fetch external data
	return fmt.Sprintf("ExternalDataFrom-%s", handlerNode), nil
}

func (ddol *DDOLManager) validateAndTransform(payload, externalData string) (string, error) {
	if len(payload) < 10 {
		return "", errors.New("payload too short to process")
	}
	if len(externalData) < 10 {
		return "", errors.New("external data too short to process")
	}
	return fmt.Sprintf("Processed-%s-%s", payload, externalData), nil
}


func (ddol *DDOLManager) MonitorOrchestrationStatus(dataID string) (string, error) {
	ddol.mu.Lock()
	defer ddol.mu.Unlock()

	orchestratedData, exists := ddol.Orchestrations[dataID]
	if !exists {
		return "", fmt.Errorf("orchestration with data ID %s not found", dataID)
	}

	log.Printf("Orchestration status for data ID %s: %s", dataID, orchestratedData.Status)
	return orchestratedData.Status, nil
}


// ExternalAPIResponse represents the structure of the JSON response from the external API
type ExternalAPIResponse struct {
	Data string `json:"data"`
}
