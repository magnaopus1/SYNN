package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// SetOraclePriority assigns a priority level to an oracle to determine its data usage precedence
func SetOraclePriority(oracleID string, priorityLevel int) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.Priority = priorityLevel
	return common.SaveOracleSource(oracle)
}

// CheckOracleTrustScore retrieves and returns the current trust score of an oracle
func CheckOracleTrustScore(oracleID string) (int, error) {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch oracle: %v", err)
	}
	return oracle.TrustScore, nil
}

// SetOracleDataExpiry sets an expiration period for the data fetched from an oracle
func SetOracleDataExpiry(oracleID string, expiryDuration time.Duration) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.DataExpiry = expiryDuration
	return common.SaveOracleSource(oracle)
}

// EnableDataFeedMonitoring activates monitoring for an oracle's data feed to track performance and errors
func EnableDataFeedMonitoring(oracleID string) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.MonitoringEnabled = true
	return common.SaveOracleSource(oracle)
}

// DisableDataFeedMonitoring deactivates monitoring for an oracle's data feed
func DisableDataFeedMonitoring(oracleID string) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.MonitoringEnabled = false
	return common.SaveOracleSource(oracle)
}

// NotifyDataFeedThresholdBreach sends a notification if a data feed exceeds a predefined threshold
func NotifyDataFeedThresholdBreach(oracleID string, thresholdType string, value float64) error {
	notification := common.OracleNotification{
		OracleID:       oracleID,
		ThresholdType:  thresholdType,
		BreachValue:    value,
		NotificationAt: time.Now(),
	}
	return common.SendNotification(notification)
}

// ResetOracleUsageMetrics resets usage metrics for an oracle, such as request counts and latency records
func ResetOracleUsageMetrics(oracleID string) error {
	metrics := common.OracleMetrics{
		OracleID:       oracleID,
		RequestCount:   0,
		AverageLatency: 0,
		LastReset:      time.Now(),
	}
	return common.SaveOracleMetrics(metrics)
}

// AnalyzeDataFeedLatency calculates the average latency of the oracle's data feed for performance assessment
func AnalyzeDataFeedLatency(oracleID string) (float64, error) {
	metrics, err := common.FetchOracleMetrics(oracleID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch oracle metrics: %v", err)
	}
	if metrics.RequestCount == 0 {
		return 0, errors.New("no requests recorded, unable to calculate latency")
	}
	return metrics.AverageLatency, nil
}

// RetrieveOracleSubscribers retrieves a list of subscribers to a specific oracle data feed
func RetrieveOracleSubscribers(oracleID string) ([]string, error) {
	subscribers, err := common.FetchOracleSubscribers(oracleID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscribers for oracle: %v", err)
	}
	return subscribers, nil
}

// SetMultipleOracleSupport enables an oracle to support multiple data types or sources simultaneously
func SetMultipleOracleSupport(oracleID string, dataTypes []string) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.SupportedDataTypes = dataTypes
	return common.SaveOracleSource(oracle)
}
