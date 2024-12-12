package automations

import (
	"log"
	"sync"
	"time"

	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/transactions"
)

// TransactionMetricsAutomation automates the collection and reporting of transaction metrics.
type TransactionMetricsAutomation struct {
	ledgerInstance    *ledger.Ledger
	metricsManager    *transactions.TransactionMetricsManager
	mutex             sync.Mutex
	stopChan          chan bool
	collectionPeriod  time.Duration
	lastCollectedTime time.Time
}

// NewTransactionMetricsAutomation initializes a new TransactionMetricsAutomation.
func NewTransactionMetricsAutomation(ledgerInstance *ledger.Ledger, metricsManager *transactions.TransactionMetricsManager) *TransactionMetricsAutomation {
	return &TransactionMetricsAutomation{
		ledgerInstance:    ledgerInstance,
		metricsManager:    metricsManager,
		stopChan:          make(chan bool),
		collectionPeriod:  60 * time.Second, // Metrics collected and reset every 60 seconds
		lastCollectedTime: time.Now(),
	}
}

// Start begins the process of periodically collecting and resetting transaction metrics.
func (t *TransactionMetricsAutomation) Start() {
	go t.runMetricsCollectionLoop()
	log.Println("Transaction Metrics Automation started.")
}

// Stop halts the metrics collection automation.
func (t *TransactionMetricsAutomation) Stop() {
	t.stopChan <- true
	log.Println("Transaction Metrics Automation stopped.")
}

// runMetricsCollectionLoop continuously collects and resets transaction metrics every collection period.
func (t *TransactionMetricsAutomation) runMetricsCollectionLoop() {
	ticker := time.NewTicker(t.collectionPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.collectAndResetMetrics()
		case <-t.stopChan:
			return
		}
	}
}

// collectAndResetMetrics gathers transaction metrics, logs them, and resets the data for the next period.
func (t *TransactionMetricsAutomation) collectAndResetMetrics() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Calculate the time range for which we are collecting metrics
	periodStartTime := t.lastCollectedTime
	t.lastCollectedTime = time.Now()

	// Collect metrics from the transaction pool and ledger
	totalTransactions, err := t.metricsManager.GetTotalTransactions()
	if err != nil {
		log.Printf("Failed to retrieve total transactions: %v", err)
		return
	}

	totalGasConsumed, err := t.metricsManager.GetTotalGasConsumed()
	if err != nil {
		log.Printf("Failed to retrieve total gas consumed: %v", err)
		return
	}

	totalFeesCollected, err := t.metricsManager.GetTotalFeesCollected()
	if err != nil {
		log.Printf("Failed to retrieve total fees collected: %v", err)
		return
	}

	transactionThroughput, err := t.metricsManager.GetTransactionThroughput()
	if err != nil {
		log.Printf("Failed to calculate transaction throughput: %v", err)
		return
	}

	finalizationTime, err := t.metricsManager.GetFinalizationTime()
	if err != nil {
		log.Printf("Failed to retrieve block finalization time: %v", err)
		return
	}

	// Log metrics for auditing and further processing
	err = t.logMetrics(totalTransactions, totalGasConsumed, totalFeesCollected, transactionThroughput, finalizationTime)
	if err != nil {
		log.Printf("Failed to log transaction metrics: %v", err)
	}

	// Reset the metrics after collection
	t.metricsManager.ResetMetrics()
}

// logMetrics securely logs transaction metrics into the ledger with encryption.
func (t *TransactionMetricsAutomation) logMetrics(totalTransactions, totalGasConsumed, totalFeesCollected int, throughput float64, finalizationTime float64) error {
	logEntry := common.TransactionMetrics{
		TotalTransactions: totalTransactions,
		TotalGasConsumed:  totalGasConsumed,
		TotalFeesCollected: totalFeesCollected,
		TransactionThroughput: throughput,
		FinalizationTime: finalizationTime,
		Timestamp:        time.Now(),
	}

	// Encrypt the log entry for secure storage
	encryptedLog, err := encryption.EncryptData(logEntry.ToString(), common.EncryptionKey)
	if err != nil {
		return err
	}

	// Store the encrypted log in the ledger
	err = t.ledgerInstance.RecordMetricsLog(encryptedLog)
	if err != nil {
		return err
	}

	log.Printf("Transaction metrics collected and logged successfully.")
	return nil
}

// Utility functions for metrics logging

// RetrieveMetrics allows manual retrieval of metrics for the current period
func (t *TransactionMetricsAutomation) RetrieveMetrics() (common.TransactionMetrics, error) {
	metrics := common.TransactionMetrics{
		TotalTransactions: t.metricsManager.GetTotalTransactions(),
		TotalGasConsumed:  t.metricsManager.GetTotalGasConsumed(),
		TotalFeesCollected: t.metricsManager.GetTotalFeesCollected(),
		TransactionThroughput: t.metricsManager.GetTransactionThroughput(),
		FinalizationTime: t.metricsManager.GetFinalizationTime(),
		Timestamp: time.Now(),
	}

	// Log and Encrypt the collected metrics
	encryptedLog, err := encryption.EncryptData(metrics.ToString(), common.EncryptionKey)
	if err != nil {
		return metrics, err
	}

	// Store the encrypted metrics log into the ledger for future access
	err = t.ledgerInstance.RecordMetricsLog(encryptedLog)
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}
