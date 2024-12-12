package transactions

import (
	"fmt"
	"time"

	"synnergy_network/pkg/ledger"
)

// NewTransactionMetricsManager initializes the transaction metrics manager.
func NewTransactionMetricsManager(ledgerInstance *ledger.Ledger) *TransactionMetricsManager {
	return &TransactionMetricsManager{
		ledger: ledgerInstance,
	}
}

// RecordTransactionMetrics records the metrics for a transaction, including gas used and fees.
func (tmm *TransactionMetricsManager) RecordTransactionMetrics(transactionID string, gasUsed int, gasLimit int, fee float64) error {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	// Increment total transactions
	tmm.totalTransactions++

	// Calculate gas efficiency (used gas vs gas limit)
	var gasEfficiency float64
	if gasLimit > 0 {
		gasEfficiency = float64(gasUsed) / float64(gasLimit)
	}
	tmm.gasEfficiency = gasEfficiency

	// Add gas consumption
	tmm.totalGasConsumed += gasUsed

	// Add to total fees collected
	tmm.totalFeesCollected += fee

	// Combine the transaction metrics into a single string
	metricsStr := fmt.Sprintf("GasUsed: %d, GasLimit: %d, Fee: %.2f, GasEfficiency: %.2f", gasUsed, gasLimit, fee, gasEfficiency)

	// Log metrics to the ledger for audit (pass only transactionID and the combined metrics string)
	err := tmm.ledger.RecordTransactionMetrics(transactionID, metricsStr)
	if err != nil {
		return fmt.Errorf("failed to record transaction metrics: %v", err)
	}

	return nil
}


// RecordSubBlockMetrics records metrics for each validated sub-block.
func (tmm *TransactionMetricsManager) RecordSubBlockMetrics(subBlockID string, transactionCount int, subBlockTime time.Duration, subBlockSize int64, parentBlockID string) error {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	// Increment total sub-blocks
	tmm.totalSubBlocks++

	// Update transaction throughput (transactions per second)
	tps := float64(transactionCount) / subBlockTime.Seconds()
	tmm.transactionThroughput = (tmm.transactionThroughput + tps) / 2 // Average TPS

	// Combine the sub-block metrics into a single string
	metricsStr := fmt.Sprintf("TransactionCount: %d, SubBlockTime: %s, TPS: %.2f", transactionCount, subBlockTime.String(), tps)

	// Record sub-block metrics in the ledger
	err := tmm.ledger.RecordSubBlockMetrics(subBlockID, metricsStr, subBlockSize, transactionCount, parentBlockID)
	if err != nil {
		return fmt.Errorf("failed to record sub-block metrics: %v", err)
	}

	return nil
}

// RecordBlockMetrics records metrics for each validated block.
func (tmm *TransactionMetricsManager) RecordBlockMetrics(blockID string, subBlocksCount int, blockTime time.Duration, blockSize int64, validatorID string) error {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	// Increment total blocks
	tmm.totalBlocks++

	// Combine the block metrics into a single string
	metricsStr := fmt.Sprintf("SubBlocks: %d, BlockTime: %s", subBlocksCount, blockTime.String())

	// Record block metrics in the ledger (pass blockSize, subBlocksCount, and validatorID)
	err := tmm.ledger.RecordBlockMetrics(blockID, metricsStr, blockSize, subBlocksCount, validatorID)
	if err != nil {
		return fmt.Errorf("failed to record block metrics: %v", err)
	}

	return nil
}




// GetTransactionThroughput returns the average transaction throughput (TPS).
func (tmm *TransactionMetricsManager) GetTransactionThroughput() float64 {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	return tmm.transactionThroughput
}

// GetTotalTransactions returns the total number of transactions recorded.
func (tmm *TransactionMetricsManager) GetTotalTransactions() int {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	return tmm.totalTransactions
}

// GetTotalGasConsumed returns the total gas consumed by all transactions.
func (tmm *TransactionMetricsManager) GetTotalGasConsumed() int {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	return tmm.totalGasConsumed
}

// GetTotalFeesCollected returns the total fees collected from all transactions.
func (tmm *TransactionMetricsManager) GetTotalFeesCollected() float64 {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	return tmm.totalFeesCollected
}

// GetTotalSubBlocks returns the total number of validated sub-blocks.
func (tmm *TransactionMetricsManager) GetTotalSubBlocks() int {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	return tmm.totalSubBlocks
}

// GetTotalBlocks returns the total number of validated blocks.
func (tmm *TransactionMetricsManager) GetTotalBlocks() int {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	return tmm.totalBlocks
}

// GetGasEfficiency returns the average gas efficiency (used gas vs gas limit).
func (tmm *TransactionMetricsManager) GetGasEfficiency() float64 {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	return tmm.gasEfficiency
}

// ResetMetrics resets all transaction metrics. Used for testing or system resets.
func (tmm *TransactionMetricsManager) ResetMetrics() {
	tmm.metricsLock.Lock()
	defer tmm.metricsLock.Unlock()

	tmm.totalTransactions = 0
	tmm.totalSubBlocks = 0
	tmm.totalBlocks = 0
	tmm.totalGasConsumed = 0
	tmm.totalFeesCollected = 0
	tmm.transactionThroughput = 0
	tmm.gasEfficiency = 0
}
