package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewStateChannelPerformance initializes the performance tracker for a state channel
func NewStateChannelPerformance(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.StateChannelPerformance {
	return &common.StateChannelPerformance{
		ChannelID:        channelID,
		Participants:     participants,
		TransactionTimes: []time.Duration{},
		State:            make(map[string]interface{}),
		Ledger:           ledgerInstance,
		Encryption:       encryptionService,
	}
}

// TrackTransactionTime adds a new transaction time to the performance tracker
func (scp *common.StateChannelPerformance) TrackTransactionTime(startTime, endTime time.Time) error {
	scp.mu.Lock()
	defer scp.mu.Unlock()

	duration := endTime.Sub(startTime)
	scp.TransactionTimes = append(scp.TransactionTimes, duration)

	// Log the transaction time in the ledger
	err := scp.Ledger.RecordTransactionTime(scp.ChannelID, duration, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction time: %v", err)
	}

	fmt.Printf("Transaction time of %v tracked for channel %s\n", duration, scp.ChannelID)
	return nil
}

// CalculatePerformanceMetrics calculates throughput and latency for the state channel
func (scp *common.StateChannelPerformance) CalculatePerformanceMetrics() error {
	scp.mu.Lock()
	defer scp.mu.Unlock()

	if len(scp.TransactionTimes) == 0 {
		return errors.New("no transaction times available to calculate performance metrics")
	}

	// Calculate average latency
	var totalDuration time.Duration
	for _, t := range scp.TransactionTimes {
		totalDuration += t
	}
	scp.Latency = totalDuration / time.Duration(len(scp.TransactionTimes))

	// Calculate throughput (transactions per second)
	scp.Throughput = float64(len(scp.TransactionTimes)) / totalDuration.Seconds()

	// Log the performance metrics in the ledger
	err := scp.Ledger.RecordPerformanceMetrics(scp.ChannelID, scp.Latency, scp.Throughput, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log performance metrics: %v", err)
	}

	fmt.Printf("Performance metrics calculated for channel %s: Latency = %v, Throughput = %f tps\n", scp.ChannelID, scp.Latency, scp.Throughput)
	return nil
}

// OptimizePerformance optimizes the state channel's performance by adjusting internal parameters
func (scp *common.StateChannelPerformance) OptimizePerformance() error {
	scp.mu.Lock()
	defer scp.mu.Unlock()

	// Perform optimization logic based on current performance data
	if scp.Throughput < 10 {
		scp.State["optimized"] = true
		fmt.Printf("Optimization triggered for channel %s due to low throughput\n", scp.ChannelID)
	} else {
		scp.State["optimized"] = false
		fmt.Printf("No optimization needed for channel %s\n", scp.ChannelID)
	}

	// Log the optimization event in the ledger
	err := scp.Ledger.RecordOptimization(scp.ChannelID, time.Now(), scp.State["optimized"].(bool))
	if err != nil {
		return fmt.Errorf("failed to log optimization: %v", err)
	}

	return nil
}

// RetrievePerformanceMetrics retrieves the current performance metrics for the state channel
func (scp *common.StateChannelPerformance) RetrievePerformanceMetrics() (time.Duration, float64, error) {
	scp.mu.Lock()
	defer scp.mu.Unlock()

	if len(scp.TransactionTimes) == 0 {
		return 0, 0, errors.New("no performance metrics available")
	}

	fmt.Printf("Retrieved performance metrics for channel %s: Latency = %v, Throughput = %f tps\n", scp.ChannelID, scp.Latency, scp.Throughput)
	return scp.Latency, scp.Throughput, nil
}

// ResetPerformanceMetrics clears the performance data for the state channel
func (scp *common.StateChannelPerformance) ResetPerformanceMetrics() error {
	scp.mu.Lock()
	defer scp.mu.Unlock()

	// Clear the transaction times and reset metrics
	scp.TransactionTimes = []time.Duration{}
	scp.Latency = 0
	scp.Throughput = 0

	// Log the reset event in the ledger
	err := scp.Ledger.RecordPerformanceReset(scp.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log performance reset: %v", err)
	}

	fmt.Printf("Performance metrics reset for channel %s\n", scp.ChannelID)
	return nil
}
