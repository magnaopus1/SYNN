package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// LOG_ACCESS_RESTRICTIONS logs any access restrictions applied to a data record
func LOG_ACCESS_RESTRICTIONS(dataID string, restrictionDetails string) error {
	accessLog := common.AccessRestrictionLog{
		DataID:             dataID,
		RestrictionDetails: restrictionDetails,
		LoggedAt:           time.Now(),
	}
	return common.SaveAccessRestrictionLog(accessLog)
}

// QUERY_OWNERSHIP_HISTORY retrieves the ownership history for a specific data record
func QUERY_OWNERSHIP_HISTORY(dataID string) ([]common.OwnershipRecord, error) {
	history, err := common.FetchOwnershipHistory(dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ownership history: %v", err)
	}
	return history, nil
}

// TAG_CRITICAL_DATA tags a data record as critical for added security and monitoring
func TAG_CRITICAL_DATA(dataID string) error {
	record, err := common.FetchDataRecord(dataID)
	if err != nil {
		return fmt.Errorf("failed to fetch data record: %v", err)
	}
	record.IsCritical = true
	return common.SaveDataRecord(record)
}

// UNTAG_CRITICAL_DATA removes the critical tag from a data record
func UNTAG_CRITICAL_DATA(dataID string) error {
	record, err := common.FetchDataRecord(dataID)
	if err != nil {
		return fmt.Errorf("failed to fetch data record: %v", err)
	}
	record.IsCritical = false
	return common.SaveDataRecord(record)
}

// CERTIFY_DATA_INTEGRITY certifies the integrity of a data record by generating a cryptographic hash
func CERTIFY_DATA_INTEGRITY(dataID string, data []byte) error {
	hash := sha256.Sum256(data)
	certification := common.DataCertification{
		DataID:        dataID,
		CertHash:      hex.EncodeToString(hash[:]),
		CertifiedAt:   time.Now(),
	}
	return common.SaveDataCertification(certification)
}

// VALIDATE_CERTIFICATION validates the certification hash of a data record to confirm integrity
func VALIDATE_CERTIFICATION(dataID string, currentData []byte) (bool, error) {
	certification, err := common.FetchDataCertification(dataID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch certification record: %v", err)
	}
	currentHash := sha256.Sum256(currentData)
	return certification.CertHash == hex.EncodeToString(currentHash[:]), nil
}

// GENERATE_CERTIFICATE_LOG creates a log entry for data certification events
func GENERATE_CERTIFICATE_LOG(dataID string, certHash string) error {
	certLog := common.CertificationLog{
		DataID:     dataID,
		CertHash:   certHash,
		LoggedAt:   time.Now(),
	}
	return common.SaveCertificationLog(certLog)
}

// CONFIRM_DATA_REVOCATION logs the revocation of access or status for a specific data record
func CONFIRM_DATA_REVOCATION(dataID string, reason string) error {
	revocation := common.DataRevocation{
		DataID:      dataID,
		Reason:      reason,
		RevokedAt:   time.Now(),
	}
	return common.SaveDataRevocation(revocation)
}

// REGISTER_DATA_ASSET registers a new data asset within the ledger for tracking and auditing
func REGISTER_DATA_ASSET(dataID string, assetDetails common.DataAssetDetails) error {
	asset := common.DataAsset{
		DataID:       dataID,
		Details:      assetDetails,
		RegisteredAt: time.Now(),
	}
	return common.SaveDataAsset(asset)
}

// DEREGISTER_DATA_ASSET removes a data asset from active registration within the ledger
func DEREGISTER_DATA_ASSET(dataID string) error {
	return common.RemoveDataAsset(dataID)
}

// MAP_DATA_DEPENDENCIES maps dependencies for a data asset, tracking associated records
func MAP_DATA_DEPENDENCIES(dataID string, dependencies []string) error {
	dependencyRecord := common.DataDependency{
		DataID:       dataID,
		Dependencies: dependencies,
		LoggedAt:     time.Now(),
	}
	return common.SaveDataDependency(dependencyRecord)
}

// LOG_DEPENDENCY_ACCESS logs access events for data dependencies, providing traceability
func LOG_DEPENDENCY_ACCESS(dataID string, dependencyID string) error {
	dependencyAccess := common.DependencyAccessLog{
		DataID:         dataID,
		DependencyID:   dependencyID,
		AccessedAt:     time.Now(),
	}
	return common.SaveDependencyAccessLog(dependencyAccess)
}

// SET_PROVENANCE_ATTRIBUTES sets provenance attributes for a data record, such as origin and certification
func SET_PROVENANCE_ATTRIBUTES(dataID string, attributes common.ProvenanceAttributes) error {
	provenance := common.DataProvenance{
		DataID:      dataID,
		Attributes:  attributes,
		AssignedAt:  time.Now(),
	}
	return common.SaveDataProvenance(provenance)
}

// CLEAR_PROVENANCE_ATTRIBUTES clears provenance attributes from a data record
func CLEAR_PROVENANCE_ATTRIBUTES(dataID string) error {
	return common.RemoveDataProvenance(dataID)
}
