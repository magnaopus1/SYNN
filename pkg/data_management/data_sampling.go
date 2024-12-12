package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// DATA_OVERSAMPLING applies oversampling to minority class data
func DATA_OVERSAMPLING(data []common.DataRecord, targetClass string) []common.DataRecord {
	classSamples := []common.DataRecord{}
	for _, record := range data {
		if record.Class == targetClass {
			classSamples = append(classSamples, record)
		}
	}
	oversampledData := append(data, classSamples...)
	for len(oversampledData) < 2*len(data) {
		oversampledData = append(oversampledData, classSamples[rand.Intn(len(classSamples))])
	}
	return oversampledData
}

// DATA_UNDERSAMPLING applies undersampling to majority class data
func DATA_UNDERSAMPLING(data []common.DataRecord, targetClass string) []common.DataRecord {
	undersampledData := []common.DataRecord{}
	for _, record := range data {
		if record.Class == targetClass || rand.Float64() > 0.5 {
			undersampledData = append(undersampledData, record)
		}
	}
	return undersampledData
}

// APPLY_MINIMUM_TRANSFORM replaces each feature with the minimum value for scaling
func APPLY_MINIMUM_TRANSFORM(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	minVal := data[0]
	for _, val := range data {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

// APPLY_MAXIMUM_TRANSFORM replaces each feature with the maximum value for scaling
func APPLY_MAXIMUM_TRANSFORM(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	maxVal := data[0]
	for _, val := range data {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

// IDENTIFY_CORRELATED_FEATURES identifies highly correlated features for removal
func IDENTIFY_CORRELATED_FEATURES(data [][]float64, threshold float64) map[int][]int {
	correlatedFeatures := make(map[int][]int)
	for i := 0; i < len(data); i++ {
		for j := i + 1; j < len(data); j++ {
			if math.Abs(correlation(data[i], data[j])) > threshold {
				correlatedFeatures[i] = append(correlatedFeatures[i], j)
			}
		}
	}
	return correlatedFeatures
}

// REMOVE_CORRELATED_FEATURES removes highly correlated features based on identified pairs
func REMOVE_CORRELATED_FEATURES(data [][]float64, correlated map[int][]int) [][]float64 {
	removedIndices := map[int]bool{}
	for i := range correlated {
		for _, j := range correlated[i] {
			removedIndices[j] = true
		}
	}

	filteredData := [][]float64{}
	for i := 0; i < len(data); i++ {
		if !removedIndices[i] {
			filteredData = append(filteredData, data[i])
		}
	}
	return filteredData
}

// EXTRACT_KEY_PHRASES extracts key phrases from text data
func EXTRACT_KEY_PHRASES(text string, minWordLength int) []string {
	words := strings.Fields(text)
	keyPhrases := []string{}
	for _, word := range words {
		if len(word) >= minWordLength {
			keyPhrases = append(keyPhrases, word)
		}
	}
	return keyPhrases
}

// RECORD_TIMESTAMP logs the current timestamp for a specific data record
func RECORD_TIMESTAMP(dataID string) error {
	timestampLog := common.TimestampRecord{
		DataID:    dataID,
		Timestamp: time.Now(),
	}
	return common.SaveTimestampRecord(timestampLog)
}

// TRACK_DATA_CHANGE records changes in a data record, hashing it to ensure data integrity
func TRACK_DATA_CHANGE(dataID string, newData []byte) error {
	hash := sha256.Sum256(newData)
	dataChangeLog := common.DataChangeLog{
		DataID:     dataID,
		ChangeHash: hex.EncodeToString(hash[:]),
		ChangedAt:  time.Now(),
	}
	return common.SaveDataChangeLog(dataChangeLog)
}

// STORE_DATA_LINEAGE records the lineage of data, tracking its origin and transformation
func STORE_DATA_LINEAGE(dataID string, lineageInfo string) error {
	lineageLog := common.DataLineage{
		DataID:       dataID,
		LineageInfo:  lineageInfo,
		RecordedAt:   time.Now(),
	}
	return common.SaveDataLineage(lineageLog)
}

// VERIFY_DATA_INTEGRITY verifies the integrity of a data record by comparing hash values
func VERIFY_DATA_INTEGRITY(dataID string, data []byte) (bool, error) {
	storedRecord, err := common.FetchDataChangeLog(dataID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch data change log: %v", err)
	}
	currentHash := sha256.Sum256(data)
	return storedRecord.ChangeHash == hex.EncodeToString(currentHash[:]), nil
}

// GENERATE_AUDIT_TRAIL generates an audit trail for a data record, including all access and modifications
func GENERATE_AUDIT_TRAIL(dataID string) (*common.AuditTrail, error) {
	accessLogs, err := common.FetchDataAccessLogs(dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data access logs: %v", err)
	}

	changeLogs, err := common.FetchDataChangeLogs(dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data change logs: %v", err)
	}

	return &common.AuditTrail{
		DataID:      dataID,
		AccessLogs:  accessLogs,
		ChangeLogs:  changeLogs,
		GeneratedAt: time.Now(),
	}, nil
}

// LINK_DATA_PROVENANCE links a data record to its provenance information for tracking its entire lifecycle
func LINK_DATA_PROVENANCE(dataID, provenanceID string) error {
	provenanceLink := common.ProvenanceLink{
		DataID:      dataID,
		ProvenanceID: provenanceID,
		LinkedAt:    time.Now(),
	}
	return common.SaveProvenanceLink(provenanceLink)
}

// HISTORICAL_REVERSION reverts a data record to a previous state, restoring an earlier version
func HISTORICAL_REVERSION(dataID string, version int) (*common.DataRecord, error) {
	previousVersion, err := common.FetchHistoricalVersion(dataID, version)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve historical version: %v", err)
	}
	return previousVersion, nil
}

// Helper function: correlation calculates the correlation coefficient between two data slices
func correlation(x, y []float64) float64 {
	n := len(x)
	if n != len(y) || n == 0 {
		return 0
	}

	meanX, meanY := mean(x), mean(y)
	var numerator, denomX, denomY float64

	for i := 0; i < n; i++ {
		diffX := x[i] - meanX
		diffY := y[i] - meanY
		numerator += diffX * diffY
		denomX += diffX * diffX
		denomY += diffY * diffY
	}

	if denomX == 0 || denomY == 0 {
		return 0
	}
	return numerator / math.Sqrt(denomX*denomY)
}

// Helper function: mean calculates the mean of a data slice
func mean(data []float64) float64 {
	sum := 0.0
	for _, val := range data {
		sum += val
	}
	return sum / float64(len(data))
}
