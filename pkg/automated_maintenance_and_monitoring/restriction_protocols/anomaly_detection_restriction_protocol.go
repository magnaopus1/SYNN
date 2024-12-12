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
	AnomalyDetectionInterval       = 1 * time.Second  // Check for anomalies every second
	TransactionSpikeThreshold      = 5000             // Maximum transactions per minute before flagging as anomalous
	UnusualValidatorActivityLimit  = 10               // Maximum validator activity deviation before detection
)

// AnomalyDetectionAutomation handles real-time anomaly detection across the network
type AnomalyDetectionAutomation struct {
	consensusSystem          *consensus.SynnergyConsensus
	ledgerInstance           *ledger.Ledger
	stateMutex               *sync.RWMutex
	anomalousBehaviorTracker map[string]int // Tracks anomalies detected by user or validator
}

// NewAnomalyDetectionAutomation initializes and returns an instance of AnomalyDetectionAutomation
func NewAnomalyDetectionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AnomalyDetectionAutomation {
	return &AnomalyDetectionAutomation{
		consensusSystem:          consensusSystem,
		ledgerInstance:           ledgerInstance,
		stateMutex:               stateMutex,
		anomalousBehaviorTracker: make(map[string]int),
	}
}

// StartAnomalyDetectionMonitoring starts real-time monitoring for network anomalies every second
func (automation *AnomalyDetectionAutomation) StartAnomalyDetectionMonitoring() {
	ticker := time.NewTicker(AnomalyDetectionInterval)

	go func() {
		for range ticker.C {
			automation.monitorAnomalies()
		}
	}()
}

// monitorAnomalies continuously checks network activity to detect anomalies in transactions and validator performance
func (automation *AnomalyDetectionAutomation) monitorAnomalies() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Retrieve relevant network metrics
	transactionRate := automation.consensusSystem.GetTransactionRate()
	validatorActivity := automation.consensusSystem.GetValidatorActivityMetrics()

	// Check for transaction rate anomalies
	if transactionRate > TransactionSpikeThreshold {
		automation.flagAnomalousTransactionRate(transactionRate)
	}

	// Check for anomalies in validator behavior
	for validator, activityRate := range validatorActivity {
		if activityRate > UnusualValidatorActivityLimit {
			automation.flagAnomalousValidator(validator, activityRate)
		}
	}
}

// flagAnomalousTransactionRate handles the logging and response to anomalous transaction spikes
func (automation *AnomalyDetectionAutomation) flagAnomalousTransactionRate(transactionRate float64) {
	fmt.Printf("Anomalous transaction rate detected: %.2f transactions per minute.\n", transactionRate)

	// Encrypt and log the anomaly in the ledger
	encryptedData := automation.encryptAnomalyData(fmt.Sprintf("Anomalous transaction rate: %.2f tpm", transactionRate))
	automation.logAnomalyToLedger("Anomalous Transaction Rate", encryptedData, fmt.Sprintf("Transaction rate of %.2f tpm detected.", transactionRate))
}

// flagAnomalousValidator handles the detection and flagging of unusual validator activity
func (automation *AnomalyDetectionAutomation) flagAnomalousValidator(validator string, activityRate float64) {
	fmt.Printf("Anomalous validator activity detected. Validator: %s, Activity rate: %.2f\n", validator, activityRate)

	// Track suspicious validator activity
	automation.anomalousBehaviorTracker[validator]++
	if automation.anomalousBehaviorTracker[validator] >= UnusualValidatorActivityLimit {
		fmt.Printf("Validator %s flagged for repeated anomalous activity.\n", validator)
		automation.logAnomalyToLedger("Validator Anomalous Activity", automation.encryptAnomalyData(validator), fmt.Sprintf("Validator %s exceeded activity rate of %.2f", validator, activityRate))
	}
}

// encryptAnomalyData encrypts anomaly data before logging it to the ledger
func (automation *AnomalyDetectionAutomation) encryptAnomalyData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting anomaly data:", err)
		return data
	}
	return string(encryptedData)
}

// logAnomalyToLedger logs the detected anomaly into the ledger with detailed information
func (automation *AnomalyDetectionAutomation) logAnomalyToLedger(anomalyType, anomalyData, additionalDetails string) {
	// Create a ledger entry with detailed anomaly information
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("anomaly-detection-%s-%d", anomalyType, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Anomaly Detected",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Anomaly (%s) detected: %s. Details: %s", anomalyType, anomalyData, additionalDetails),
	}

	// Add entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log anomaly into ledger: %v\n", err)
	} else {
		fmt.Printf("Anomaly logged: %s\n", anomalyType)
	}
}
