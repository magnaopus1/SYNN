package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// RequestOracleData initiates a request for data from an oracle source and logs the request
func RequestOracleData(oracleID string) (string, error) {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch oracle: %v", err)
	}
	requestID := generateRequestID()
	request := common.OracleRequest{
		RequestID:  requestID,
		OracleID:   oracleID,
		RequestedAt: time.Now(),
		Status:     "pending",
	}
	if err := common.SaveOracleRequest(request); err != nil {
		return "", fmt.Errorf("failed to log oracle request: %v", err)
	}
	return requestID, nil
}

// VerifyOracleData verifies the authenticity and integrity of data returned from an oracle
func VerifyOracleData(oracleID string, data []byte, expectedHash string) (bool, error) {
	calculatedHash := generateHash(data)
	if calculatedHash != expectedHash {
		return false, errors.New("data integrity verification failed")
	}
	return true, nil
}

// FetchExternalData retrieves data from an external source and integrates it with the oracle system
func FetchExternalData(sourceURL string) ([]byte, error) {
	data, err := common.RetrieveDataFromURL(sourceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch external data: %v", err)
	}
	return data, nil
}

// UpdateOracleData updates the data stored for a specific oracle and refreshes the last update timestamp
func UpdateOracleData(oracleID string, data []byte) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.LastUpdate = time.Now()
	oracle.DataHash = generateHash(data)
	if err := common.SaveOracleSource(oracle); err != nil {
		return fmt.Errorf("failed to update oracle data: %v", err)
	}
	return nil
}

// SetDataFeedFrequency sets the update frequency for an oracle's data feed
func SetDataFeedFrequency(oracleID string, frequency time.Duration) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	oracle.UpdateFrequency = frequency
	return common.SaveOracleSource(oracle)
}

// TriggerDataFeed manually triggers a data feed update for an oracle
func TriggerDataFeed(oracleID string) error {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return fmt.Errorf("failed to fetch oracle: %v", err)
	}
	data, err := FetchExternalData(oracle.Source)
	if err != nil {
		return fmt.Errorf("failed to fetch data during trigger: %v", err)
	}
	return UpdateOracleData(oracleID, data)
}

// GetOracleResponse retrieves the latest response data from an oracle
func GetOracleResponse(oracleID string) ([]byte, error) {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch oracle: %v", err)
	}
	data, err := common.FetchOracleData(oracleID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve oracle data: %v", err)
	}
	return data, nil
}

// LogOracleEvent logs an event related to an oracle's operation, such as data updates or errors
func LogOracleEvent(oracleID, eventDetails string) error {
	eventLog := common.OracleEventLog{
		OracleID:     oracleID,
		EventDetails: eventDetails,
		LoggedAt:     time.Now(),
	}
	return common.SaveOracleEventLog(eventLog)
}

// CheckOracleStatus checks the current operational status of an oracle
func CheckOracleStatus(oracleID string) (string, error) {
	oracle, err := common.FetchOracleSource(oracleID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch oracle: %v", err)
	}
	return oracle.Status, nil
}

// ValidateDataFeedSignature verifies the cryptographic signature of data received from an oracle feed
func ValidateDataFeedSignature(data []byte, signature string, publicKey []byte) (bool, error) {
	isValid, err := common.VerifySignature(data, signature, publicKey)
	if err != nil {
		return false, fmt.Errorf("signature validation failed: %v", err)
	}
	return isValid, nil
}

// Helper function: generateRequestID generates a unique request ID for an oracle data request
func generateRequestID() string {
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d", timestamp)))
	return hex.EncodeToString(hash[:])
}

// Helper function: generateHash generates a SHA-256 hash for given data
func generateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
