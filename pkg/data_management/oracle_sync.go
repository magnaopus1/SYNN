package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// QueryOracleSyncStatus retrieves the current synchronization status of an oracle
func QueryOracleSyncStatus(oracleID string) (string, error) {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch oracle: %v", err)
	}
	return oracle.SyncStatus, nil
}

// SetDataFeedCacheDuration sets the duration for caching data from an oracle feed
func SetDataFeedCacheDuration(oracleID string, cacheDuration time.Duration) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.CacheDuration = cacheDuration
	return common.SaveOracleSource(oracle)
}

// CheckOracleEventHistory retrieves the event history for an oracle to review past actions and updates
func CheckOracleEventHistory(oracleID string, fromDate, toDate time.Time) ([]common.OracleEventLog, error) {
	events, err := common.FetchOracleEventsInRange(oracleID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch oracle event history: %v", err)
	}
	return events, nil
}

// UpdateOracleReliabilityScore adjusts the reliability score of an oracle based on recent performance metrics
func UpdateOracleReliabilityScore(oracleID string, newScore int) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	if newScore < 0 || newScore > 100 {
		return errors.New("reliability score must be between 0 and 100")
	}
	oracle.ReliabilityScore = newScore
	return common.SaveOracleSource(oracle)
}

// SetOracleRetryLimit sets the maximum number of retries allowed for an oracle in case of data retrieval failure
func SetOracleRetryLimit(oracleID string, retryLimit int) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	if retryLimit < 0 {
		return errors.New("retry limit must be a non-negative integer")
	}
	oracle.RetryLimit = retryLimit
	return common.SaveOracleSource(oracle)
}
