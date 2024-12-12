package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// QUERY_PROVENANCE_HISTORY retrieves the full history of changes for a data record
func QUERY_PROVENANCE_HISTORY(dataID string) ([]common.ProvenanceRecord, error) {
	history, err := common.FetchProvenanceHistory(dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve provenance history: %v", err)
	}
	return history, nil
}

// LOG_DATA_ACCESS logs each access event of a data record for tracking
func LOG_DATA_ACCESS(dataID, accessorID string) error {
	accessLog := common.DataAccessLog{
		DataID:      dataID,
		AccessorID:  accessorID,
		AccessTime:  time.Now(),
	}
	return common.SaveDataAccessLog(accessLog)
}

// VALIDATE_DATA_ORIGIN verifies the original source of a data record for authenticity
func VALIDATE_DATA_ORIGIN(dataID, originalSource string) (bool, error) {
	provenance, err := common.FetchProvenanceRecord(dataID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch provenance record: %v", err)
	}
	return provenance.OriginalSource == originalSource, nil
}

// STORE_ORIGINAL_HASH stores the hash of the original data for future integrity checks
func STORE_ORIGINAL_HASH(dataID string, data []byte) error {
	hash := sha256.Sum256(data)
	originalHashRecord := common.DataHashRecord{
		DataID:        dataID,
		OriginalHash:  hex.EncodeToString(hash[:]),
		RecordedAt:    time.Now(),
	}
	return common.SaveDataHashRecord(originalHashRecord)
}

// COMPARE_WITH_ORIGINAL_HASH compares current data hash with the original to check for modifications
func COMPARE_WITH_ORIGINAL_HASH(dataID string, currentData []byte) (bool, error) {
	originalRecord, err := common.FetchDataHashRecord(dataID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch original hash record: %v", err)
	}
	currentHash := sha256.Sum256(currentData)
	return originalRecord.OriginalHash == hex.EncodeToString(currentHash[:]), nil
}

// ARCHIVE_DATA_VERSION archives a specific version of data for long-term storage
func ARCHIVE_DATA_VERSION(dataID string, version int) error {
	data, err := common.FetchDataVersion(dataID, version)
	if err != nil {
		return fmt.Errorf("failed to fetch data version: %v", err)
	}
	data.IsArchived = true
	data.ArchivedAt = time.Now()
	return common.SaveArchivedDataVersion(data)
}

// RESTORE_ARCHIVED_VERSION restores a previously archived data version for active use
func RESTORE_ARCHIVED_VERSION(dataID string, version int) (*common.DataRecord, error) {
	archivedData, err := common.FetchArchivedDataVersion(dataID, version)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch archived data version: %v", err)
	}
	if !archivedData.IsArchived {
		return nil, errors.New("specified data version is not archived")
	}
	return archivedData, nil
}

// IDENTIFY_DATA_ANOMALIES scans for anomalies in data patterns, indicating potential tampering
func IDENTIFY_DATA_ANOMALIES(dataID string) ([]common.AnomalyRecord, error) {
	anomalies, err := common.ScanDataForAnomalies(dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to identify data anomalies: %v", err)
	}
	return anomalies, nil
}

// ASSIGN_DATA_OWNER assigns an owner to a data record for ownership tracking
func ASSIGN_DATA_OWNER(dataID, ownerID string) error {
	ownershipRecord := common.DataOwnership{
		DataID:    dataID,
		OwnerID:   ownerID,
		AssignedAt: time.Now(),
	}
	return common.SaveDataOwnership(ownershipRecord)
}

// VERIFY_OWNERSHIP checks if a specified user is the current owner of the data record
func VERIFY_OWNERSHIP(dataID, userID string) (bool, error) {
	ownershipRecord, err := common.FetchDataOwnership(dataID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch data ownership record: %v", err)
	}
	return ownershipRecord.OwnerID == userID, nil
}

// TRACK_DATA_RETENTION logs data retention events to track data's lifecycle within retention policies
func TRACK_DATA_RETENTION(dataID string, eventType string) error {
	retentionLog := common.DataRetentionLog{
		DataID:     dataID,
		EventType:  eventType,
		Timestamp:  time.Now(),
	}
	return common.SaveDataRetentionLog(retentionLog)
}

// RECORD_DATA_LIFECYCLE_EVENT logs significant lifecycle events of a data record
func RECORD_DATA_LIFECYCLE_EVENT(dataID, eventType string, description string) error {
	lifecycleEvent := common.DataLifecycleEvent{
		DataID:      dataID,
		EventType:   eventType,
		Description: description,
		Timestamp:   time.Now(),
	}
	return common.SaveDataLifecycleEvent(lifecycleEvent)
}

// SET_DATA_EXPIRATION sets an expiration date for a data record for automatic deletion or archival
func SET_DATA_EXPIRATION(dataID string, expirationDate time.Time) error {
	expiration := common.DataExpiration{
		DataID:         dataID,
		ExpirationDate: expirationDate,
		SetAt:          time.Now(),
	}
	return common.SaveDataExpiration(expiration)
}

// REVOKE_DATA_ACCESS revokes access to a data record for a specified user
func REVOKE_DATA_ACCESS(dataID, userID string) error {
	revokeRecord := common.AccessRevocation{
		DataID:    dataID,
		UserID:    userID,
		RevokedAt: time.Now(),
	}
	return common.SaveAccessRevocation(revokeRecord)
}
