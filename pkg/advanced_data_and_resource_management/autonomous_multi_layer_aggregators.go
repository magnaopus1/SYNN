package advanced_data_and_resource_management

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewAMLAManager initializes a new manager for multi-layer aggregation.
func NewAMLAManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *AMLAManager {
	return &AMLAManager{
		Aggregations:      make(map[string]*AMLAData),
		PendingRollups:    []*AMLARollup{},
		CompletedData:     []*AMLAData{},
		EncryptionService: encryptionService,
		Ledger:            ledgerInstance,
	}
}

// CreateRollup handles incoming rollup data from different layers, encrypts it, and logs it in the ledger.
func (amla *AMLAManager) CreateRollup(layer int, payload string, aggregatorNode string) (*AMLARollup, error) {
	amla.mu.Lock()
	defer amla.mu.Unlock()

	if layer <= 0 {
		return nil, errors.New("invalid layer: must be greater than 0")
	}
	if len(payload) == 0 {
		return nil, errors.New("payload cannot be empty")
	}
	if len(aggregatorNode) == 0 {
		return nil, errors.New("aggregator node cannot be empty")
	}

	rollupID := generateUniqueID()
	encryptedPayload, err := amla.EncryptionService.EncryptData("AES", []byte(payload), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt rollup payload: %w", err)
	}
	encryptedPayloadStr := base64.StdEncoding.EncodeToString(encryptedPayload)

	rollup := &AMLARollup{
		RollupID:         rollupID,
		Layer:            layer,
		DataPayload:      payload,
		AggregatorNode:   aggregatorNode,
		AggregationTime:  time.Now(),
		EncryptedPayload: encryptedPayloadStr,
	}

	amla.PendingRollups = append(amla.PendingRollups, rollup)
	ledgerMsg, err := amla.Ledger.RollupLedger.RecordRollupCreation(rollupID)
	if err != nil {
		return nil, fmt.Errorf("failed to log rollup in ledger: %w", err)
	}
	log.Printf("Rollup %s created for layer %d by node %s. Ledger Message: %s", rollupID, layer, aggregatorNode, ledgerMsg)
	return rollup, nil
}

// AggregateRollups aggregates all pending rollups into a single AMLAData instance.
func (amla *AMLAManager) AggregateRollups(validatorNode string) (*AMLAData, error) {
	amla.mu.Lock()
	defer amla.mu.Unlock()

	if len(amla.PendingRollups) == 0 {
		return nil, errors.New("no pending rollups available for aggregation")
	}
	if len(validatorNode) == 0 {
		return nil, errors.New("validator node cannot be empty")
	}

	dataID := generateUniqueID()
	rollups := amla.PendingRollups
	amla.PendingRollups = nil

	aggregation := &AMLAData{
		DataID:        dataID,
		Rollups:       rollups,
		ValidatorNode: validatorNode,
		Status:        "Pending",
	}

	amla.Aggregations[dataID] = aggregation
	log.Printf("Aggregated %d rollups into data ID %s, awaiting validation by node %s", len(rollups), dataID, validatorNode)
	return aggregation, nil
}

// ValidateAggregation validates an aggregation, encrypts its ID, and logs the validation in the ledger.
func (amla *AMLAManager) ValidateAggregation(dataID string, validatorNode string) error {
	amla.mu.Lock()
	defer amla.mu.Unlock()

	if len(dataID) == 0 || len(validatorNode) == 0 {
		return errors.New("data ID and validator node must not be empty")
	}

	aggregation, exists := amla.Aggregations[dataID]
	if !exists {
		return fmt.Errorf("aggregation %s not found", dataID)
	}

	encryptedDataID, err := amla.EncryptionService.EncryptData("AES", []byte(dataID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data ID: %w", err)
	}

	aggregation.ValidationTime = time.Now()
	aggregation.ValidatorNode = validatorNode
	aggregation.EncryptedDataID = base64.StdEncoding.EncodeToString(encryptedDataID)
	aggregation.Status = "Validated"

	ledgerMsg, err := amla.Ledger.DataManagementLedger.RecordAggregationValidation(dataID)
	if err != nil {
		return fmt.Errorf("failed to log aggregation validation in ledger: %w", err)
	}
	log.Printf("Aggregation %s validated by node %s. Ledger Message: %s", dataID, validatorNode, ledgerMsg)
	amla.CompletedData = append(amla.CompletedData, aggregation)
	delete(amla.Aggregations, dataID)
	return nil
}

// MonitorAggregationStatus retrieves the status of a specific aggregation.
func (amla *AMLAManager) MonitorAggregationStatus(dataID string) (string, error) {
	amla.mu.Lock()
	defer amla.mu.Unlock()

	aggregation, exists := amla.Aggregations[dataID]
	if !exists {
		return "", fmt.Errorf("aggregation %s not found", dataID)
	}

	return aggregation.Status, nil
}

// generateUniqueID creates a cryptographically secure unique ID.
func generateUniqueID() string {
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		log.Fatalf("failed to generate unique ID: %v", err)
	}
	return hex.EncodeToString(id)
}
