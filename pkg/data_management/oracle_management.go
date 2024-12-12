package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// RegisterNewOracle registers a new oracle source in the ledger
func RegisterNewOracle(oracleID, oracleSource string, updateInterval time.Duration, trustLevel int) error {
	oracle := common.OracleSource{
		OracleID:       oracleID,
		Source:         oracleSource,
		LastUpdate:     time.Now(),
		UpdateInterval: updateInterval,
		TrustLevel:     trustLevel,
		Status:         "active",
	}
	return common.SaveOracleSource(oracle)
}

// RemoveOracleSource deactivates an oracle source in the ledger
func RemoveOracleSource(oracleID string) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.Status = "inactive"
	oracle.RemovedAt = time.Now()
	return common.SaveOracleSource(oracle)
}

// QueryOracleLastUpdate retrieves the last update timestamp for a specified oracle
func QueryOracleLastUpdate(oracleID string) (time.Time, error) {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to fetch oracle: %v", err)
	}
	return oracle.LastUpdate, nil
}

// SetOracleUpdateInterval sets a new update interval for an oracle source
func SetOracleUpdateInterval(oracleID string, newInterval time.Duration) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.UpdateInterval = newInterval
	return common.SaveOracleSource(oracle)
}

// EnableOracleRetry enables retry for failed updates for a specified oracle
func EnableOracleRetry(oracleID string) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.RetryEnabled = true
	return common.SaveOracleSource(oracle)
}

// DisableOracleRetry disables retry for failed updates for a specified oracle
func DisableOracleRetry(oracleID string) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.RetryEnabled = false
	return common.SaveOracleSource(oracle)
}

// CheckDataFeedIntegrity verifies the integrity of the data feed from the oracle
func CheckDataFeedIntegrity(oracleID string, expectedHash string) (bool, error) {
	oracleData, err := common.FetchOracleData(oracleID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch oracle data: %v", err)
	}
	calculatedHash := common.CalculateHash(oracleData)
	return calculatedHash == expectedHash, nil
}

// SetOracleTrustLevel sets the trust level of an oracle to a specific value
func SetOracleTrustLevel(oracleID string, trustLevel int) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.TrustLevel = trustLevel
	return common.SaveOracleSource(oracle)
}

// FetchOracleMetadata retrieves metadata about an oracle source
func FetchOracleMetadata(oracleID string) (common.OracleMetadata, error) {
	metadata, err := common.FetchOracleMetadata(oracleID)
	if err != nil {
		return common.OracleMetadata{}, fmt.Errorf("failed to fetch oracle metadata: %v", err)
	}
	return metadata, nil
}

// LogDataFeedError logs an error encountered during data feed retrieval
func LogDataFeedError(oracleID, errorDetails string) error {
	errorLog := common.OracleErrorLog{
		OracleID:     oracleID,
		ErrorDetails: errorDetails,
		LoggedAt:     time.Now(),
	}
	return common.SaveOracleErrorLog(errorLog)
}
