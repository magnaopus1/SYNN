package data_management

import (
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// ConfigureDataFeedTimeout sets a timeout for data feed connections and retrievals
func ConfigureDataFeedTimeout(feedID string, timeout time.Duration) error {
	dataFeedConfig, err := ledger.DataManagementLedger.FetchDataFeedConfig(feedID)
	if err != nil {
		return fmt.Errorf("failed to fetch data feed config: %v", err)
	}

	dataFeedConfig.Timeout = timeout
	return ledger.SaveDataFeedConfig(dataFeedConfig)
}

// RegisterDataFeedSubscriber registers a subscriber to receive updates from a specific data feed
func RegisterDataFeedSubscriber(feedID, subscriberID string) error {
	subscriber := DataFeedSubscriber{
		FeedID:       feedID,
		SubscriberID: subscriberID,
		RegisteredAt: time.Now(),
	}

	return ledger.AddDataFeedSubscriber(subscriber)
}

// UnregisterDataFeedSubscriber removes a subscriber from a data feed subscription
func UnregisterDataFeedSubscriber(feedID, subscriberID string) error {
	return common.RemoveDataFeedSubscriber(feedID, subscriberID)
}

// QueryDataFeedFrequency retrieves the update frequency of a specific data feed
func QueryDataFeedFrequency(feedID string) (time.Duration, error) {
	config, err := ledger.FetchDataFeedConfig(feedID)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve data feed frequency: %v", err)
	}
	return config.UpdateFrequency, nil
}

// RetrieveOracleHistory fetches the historical data of oracle feed responses
func RetrieveOracleHistory(oracleID string, startTime, endTime time.Time) ([]common.OracleRecord, error) {
	history, err := common.FetchOracleHistory(oracleID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve oracle history: %v", err)
	}
	return history, nil
}

// AnalyzeOraclePerformance assesses the performance metrics of an oracle data feed
func AnalyzeOraclePerformance(oracleID string) (*common.OraclePerformanceReport, error) {
	records, err := common.FetchOracleHistory(oracleID, time.Time{}, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch oracle history for performance analysis: %v", err)
	}

	performance := &common.OraclePerformanceReport{
		OracleID: oracleID,
		TotalRequests: len(records),
	}

	var totalResponseTime time.Duration
	for _, record := range records {
		performance.SuccessfulResponses += record.IsSuccess
		performance.FailedResponses += !record.IsSuccess
		totalResponseTime += record.ResponseTime
	}

	if performance.TotalRequests > 0 {
		performance.AverageResponseTime = totalResponseTime / time.Duration(performance.TotalRequests)
	}
	return performance, nil
}

// SetDataFeedAlertThreshold configures alert thresholds for a data feed
func SetDataFeedAlertThreshold(feedID string, threshold common.AlertThreshold) error {
	alertConfig := common.DataFeedAlert{
		FeedID:        feedID,
		ThresholdType: threshold.Type,
		ThresholdValue: threshold.Value,
	}

	return common.SaveDataFeedAlert(alertConfig)
}

// ClearDataFeedAlert clears existing alerts set for a data feed
func ClearDataFeedAlert(feedID string) error {
	return common.RemoveDataFeedAlert(feedID)
}

// EnableOracleEventLogging enables logging for oracle data feed events
func EnableOracleEventLogging(oracleID string) error {
	logConfig, err := common.FetchOracleLogConfig(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch log config: %v", err)
	}

	logConfig.LoggingEnabled = true
	return common.SaveOracleLogConfig(logConfig)
}

// DisableOracleEventLogging disables logging for oracle data feed events
func DisableOracleEventLogging(oracleID string) error {
	logConfig, err := common.FetchOracleLogConfig(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch log config: %v", err)
	}

	logConfig.LoggingEnabled = false
	return common.SaveOracleLogConfig(logConfig)
}
