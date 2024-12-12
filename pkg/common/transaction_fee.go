package common

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"synnergy_network/pkg/ledger"
)

type FeeType string

const (
	BaseFee    FeeType = "BaseFee"
	VariableFee FeeType = "VariableFee"
	PriorityFee FeeType = "PriorityFee"
)

// TransactionFee represents the breakdown of fees for a single transaction.
type TransactionFee struct {
	TransactionID string   // Unique identifier for the transaction.
	BaseFee       float64  // The base fee, determined by the median network fee.
	VariableFee   float64  // The variable fee, based on gas units and gas price.
	PriorityFee   float64  // An optional tip to speed up the transaction processing.
	TotalFee      float64  // The total calculated fee for the transaction.
	FeeType       FeeType  // The type of fee (BaseFee, VariableFee, etc.)
}

// NetworkFeeManager handles transaction fee calculation, distribution, and refunding.
type NetworkFeeManager struct {
    feeLock         sync.Mutex      // For thread-safe fee management
    ledgerInstance  *ledger.Ledger  // Reference to the ledger for recording transactions
    medianFees      []float64       // Stores the last 1000 blocks' fees for base fee calculation
}

// NewNetworkFeeManager initializes the fee manager.
func NewNetworkFeeManager(ledgerInstance *ledger.Ledger) *NetworkFeeManager {
	return &NetworkFeeManager{
		ledgerInstance: ledgerInstance,
		medianFees:     make([]float64, 1000), // Stores the last 1000 blocks' fees
	}
}

// CalculateBaseFee calculates the base fee using the median fee of the last 1000 blocks.
func (nfm *NetworkFeeManager) CalculateBaseFee() (float64, error) {
	nfm.feeLock.Lock()
	defer nfm.feeLock.Unlock()

	// Median calculation from last 1000 blocks
	medianFee := calculateMedian(nfm.medianFees)

	// Adjust the base fee dynamically based on network conditions
	adjustmentFactor := nfm.calculateAdjustmentFactor()
	baseFee := medianFee * (1 + adjustmentFactor)

	return baseFee, nil
}

// calculateMedian calculates the median of an array of fees.
func calculateMedian(fees []float64) float64 {
	sortedFees := make([]float64, len(fees))
	copy(sortedFees, fees)
	sort.Float64s(sortedFees)

	mid := len(sortedFees) / 2
	if len(sortedFees)%2 == 0 {
		return (sortedFees[mid-1] + sortedFees[mid]) / 2
	}
	return sortedFees[mid]
}

// CalculateVariableFee calculates the variable fee based on gas units and gas price.
func (nfm *NetworkFeeManager) CalculateVariableFee(gasUnits int, gasPricePerUnit float64) float64 {
	return float64(gasUnits) * gasPricePerUnit
}

// CalculatePriorityFee allows the user to specify a tip for faster transaction processing.
func (nfm *NetworkFeeManager) CalculatePriorityFee(userSpecifiedTip float64) float64 {
	return userSpecifiedTip
}

// CalculateTotalFee calculates the total fee based on the base fee, variable fee, and optional priority fee.
func (nfm *NetworkFeeManager) CalculateTotalFee(baseFee, variableFee, priorityFee float64) float64 {
	return baseFee + variableFee + priorityFee
}

// RecordMedianFee updates the median fee array with the latest block fee.
func (nfm *NetworkFeeManager) RecordMedianFee(fee float64) {
	nfm.feeLock.Lock()
	defer nfm.feeLock.Unlock()

	nfm.medianFees = append(nfm.medianFees[1:], fee)
}

// calculateAdjustmentFactor dynamically adjusts the base fee based on network congestion.
func (nfm *NetworkFeeManager) calculateAdjustmentFactor() float64 {
	// Network conditions can adjust this value (e.g., based on block capacity utilization)
	// In a real-world implementation, this should reflect network congestion and capacity.
	congestionFactor := 0.1 // Example static factor
	return congestionFactor
}

// ProcessTransactionFee processes the transaction fee by calculating and applying the base, variable, and priority fees.
func (nfm *NetworkFeeManager) ProcessTransactionFee(transactionID string, gasUnits int, gasPricePerUnit float64, userTip float64) (*TransactionFee, error) {
	baseFee, err := nfm.CalculateBaseFee()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate base fee: %v", err)
	}

	variableFee := nfm.CalculateVariableFee(gasUnits, gasPricePerUnit)
	priorityFee := nfm.CalculatePriorityFee(userTip)
	totalFee := nfm.CalculateTotalFee(baseFee, variableFee, priorityFee)

	// Convert transactionID from string to uint64 for recording in the ledger
	convertedTransactionID, err := convertTransactionIDToUint64(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert transaction ID to uint64: %v", err)
	}

	// Convert uint64 to string for the ledger method
	convertedTransactionIDStr := strconv.FormatUint(convertedTransactionID, 10)

	// Convert totalFee from float64 to uint64 (you might want to round it)
	totalFeeUint64 := uint64(totalFee)

	// Record the transaction fee in the ledger (for auditing purposes)
	err = nfm.ledgerInstance.BlockchainConsensusCoinLedger.RecordTransactionFee(convertedTransactionIDStr, totalFeeUint64)
	if err != nil {
		return nil, fmt.Errorf("failed to record transaction fee: %v", err)
	}

	// Return the structured fee details
	return &TransactionFee{
		BaseFee:     baseFee,
		VariableFee: variableFee,
		PriorityFee: priorityFee,
		TotalFee:    totalFee,
	}, nil
}

// Helper function to convert transaction ID from string to uint64
func convertTransactionIDToUint64(transactionID string) (uint64, error) {
	convertedID, err := strconv.ParseUint(transactionID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid transaction ID format: %v", err)
	}
	return convertedID, nil
}


// Operation represents an action in a transaction that consumes gas.
type Operation struct {
	Type  string      // Type of operation (e.g., "Transfer", "SmartContractExecution", "DataStorage", etc.)
	Value interface{} // Additional details about the operation (could be amount, contract ID, or other context)
}

// GetGasCostForOperation returns the gas cost for a given operation type.
func GetGasCostForOperation(op Operation) int {
	switch op.Type {
	case "Transfer":
		return 5000 // Reduced gas cost for a basic transfer
	case "SmartContractExecution":
		return 10000 // Lowered gas cost for smart contract execution
	case "DataStorage":
		return 25000 // Reduced gas cost for data storage on-chain
	case "GovernanceVote":
		return 3000 // Gas cost for governance-related operations like voting
	case "TokenMinting":
		return 12000 // Gas cost for minting new tokens
	default:
		return 1000 // Minimum gas cost for unrecognized operation types
	}
}

// TransactionGasEstimator estimates the total gas required for a transaction based on its operations.
func TransactionGasEstimator(operations []Operation) int {
	totalGasUnits := 0
	for _, op := range operations {
		totalGasUnits += GetGasCostForOperation(op)
	}
	return totalGasUnits
}

// CalculateGasFee calculates the gas fee as a percentage of the transaction value, with a ceiling of 0.25% of the transaction value.
func CalculateGasFee(transactionValue float64, totalGasUnits int, gasPricePerUnit float64) float64 {
	// Calculate the raw gas fee
	rawGasFee := float64(totalGasUnits) * gasPricePerUnit

	// Set a ceiling for gas fees at 0.25% of the transaction value
	feeCeiling := transactionValue * 0.0025

	// Apply the ceiling if the raw gas fee exceeds it
	if rawGasFee > feeCeiling {
		return feeCeiling
	}
	return rawGasFee
}

// RefundUnusedGas refunds unused gas if the transaction consumes less than the gas limit.
func (nfm *NetworkFeeManager) RefundUnusedGas(transactionID string, gasUsed int, gasLimit int, gasPricePerUnit float64) error {
	if gasUsed < gasLimit {
		refundAmount := float64(gasLimit-gasUsed) * gasPricePerUnit
		return nfm.ledgerInstance.BlockchainConsensusCoinLedger.RefundTransactionGas(transactionID, uint64(refundAmount)) // Cast float64 to uint64
	}
	return nil
}

// DistributeTransactionFees handles the distribution of fees to validators, miners, and other pool allocations.
func (nfm *NetworkFeeManager) DistributeTransactionFees(fee *TransactionFee) error {
	// Define the fee distribution percentages
	validatorShare := 0.70 * fee.TotalFee
	charityPoolShare := 0.10 * fee.TotalFee
	loanPoolShare := 0.05 * fee.TotalFee
	passiveIncomeShare := 0.05 * fee.TotalFee
	authorityNodeShare := 0.05 * fee.TotalFee
	devPoolShare := 0.05 * fee.TotalFee

	// Create a map of fees to distribute
	feeDistribution := map[string]uint64{
		"validator":      uint64(validatorShare),
		"charityPool":    uint64(charityPoolShare),
		"loanPool":       uint64(loanPoolShare),
		"passiveIncome":  uint64(passiveIncomeShare),
		"authorityNode":  uint64(authorityNodeShare),
		"devPool":        uint64(devPoolShare),
	}

	// Distribute fees
	err := nfm.ledgerInstance.BlockchainConsensusCoinLedger.DistributeFees(fee.TransactionID, feeDistribution)
	if err != nil {
		return fmt.Errorf("error distributing transaction fees: %v", err)
	}
	return nil
}
