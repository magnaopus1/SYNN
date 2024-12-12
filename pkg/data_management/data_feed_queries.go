package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// GetLatestDataFeedValue retrieves the most recent value from a specified data feed
func GetLatestDataFeedValue(feedID string) (common.DataFeedValue, error) {
	latestValue, err := common.FetchLatestFeedValue(feedID)
	if err != nil {
		return common.DataFeedValue{}, fmt.Errorf("failed to retrieve latest data feed value: %v", err)
	}
	return latestValue, nil
}

// QueryOracleFeedUpdates retrieves recent updates from an oracle feed
func QueryOracleFeedUpdates(oracleID string, limit int) ([]common.OracleFeedUpdate, error) {
	updates, err := common.FetchOracleFeedUpdates(oracleID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query oracle feed updates: %v", err)
	}
	return updates, nil
}

// SetFeedDataCaching enables or configures caching for a data feed
func SetFeedDataCaching(feedID string, cachingEnabled bool, duration time.Duration) error {
	cacheConfig := common.FeedCacheConfig{
		FeedID:         feedID,
		CachingEnabled: cachingEnabled,
		CacheDuration:  duration,
	}
	return common.SaveFeedCacheConfig(cacheConfig)
}

// ClearFeedDataCache clears cached data for a specified data feed
func ClearFeedDataCache(feedID string) error {
	return common.RemoveFeedDataCache(feedID)
}

// SynchronizeOracleData synchronizes oracle data across the network to ensure consistency
func SynchronizeOracleData(oracleID string) error {
	data, err := common.FetchOracleFeedUpdates(oracleID, 0) // Fetch all data for sync
	if err != nil {
		return fmt.Errorf("failed to fetch oracle data for synchronization: %v", err)
	}

	for _, record := range data {
		err = common.BroadcastDataToNodes(record) // Synchronize data across nodes
		if err != nil {
			return fmt.Errorf("failed to broadcast data to nodes: %v", err)
		}
	}
	return nil
}

// SetExternalDataSource registers an external data source for a specific data feed
func SetExternalDataSource(feedID, sourceName, sourceURL string) error {
	source := common.ExternalDataSource{
		FeedID:     feedID,
		SourceName: sourceName,
		SourceURL:  sourceURL,
		AddedAt:    time.Now(),
	}
	return common.SaveExternalDataSource(source)
}

// RemoveExternalDataSource removes a previously registered external data source
func RemoveExternalDataSource(feedID, sourceName string) error {
	return common.DeleteExternalDataSource(feedID, sourceName)
}

// RetrieveDataFeedErrorLogs fetches error logs associated with a specific data feed
func RetrieveDataFeedErrorLogs(feedID string) ([]common.DataFeedErrorLog, error) {
	logs, err := common.FetchDataFeedErrorLogs(feedID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data feed error logs: %v", err)
	}
	return logs, nil
}

// ValidateExternalDataSource checks the validity and reliability of an external data source
func ValidateExternalDataSource(sourceURL string) (bool, error) {
	// Example validation by attempting to reach and confirm the response from the data source
	valid, err := common.ValidateDataSourceConnection(sourceURL)
	if err != nil {
		return false, fmt.Errorf("validation of data source failed: %v", err)
	}
	return valid, nil
}

// TrackOracleResponseTime logs and calculates the response time of oracle data feed
func TrackOracleResponseTime(oracleID string) (time.Duration, error) {
	startTime := time.Now()
	
	// Simulate data retrieval for response time calculation
	_, err := common.FetchLatestFeedValue(oracleID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch oracle feed value: %v", err)
	}

	responseTime := time.Since(startTime)
	log := common.OracleResponseTimeLog{
		OracleID:      oracleID,
		ResponseTime:  responseTime,
		Timestamp:     time.Now(),
	}
	err = common.SaveOracleResponseTimeLog(log)
	if err != nil {
		return 0, fmt.Errorf("failed to log oracle response time: %v", err)
	}

	return responseTime, nil
}
