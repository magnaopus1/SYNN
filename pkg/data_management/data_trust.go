package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// ESTABLISH_PROVENANCE_LINK links a data record to its provenance record for lifecycle tracking
func ESTABLISH_PROVENANCE_LINK(dataID, provenanceID string) error {
	link := common.ProvenanceLink{
		DataID:       dataID,
		ProvenanceID: provenanceID,
		LinkedAt:     time.Now(),
	}
	return common.SaveProvenanceLink(link)
}

// ENCRYPT_DATA encrypts data using AES encryption and returns the encrypted data as a byte array
func ENCRYPT_DATA(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %v", err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

// DECRYPT_DATA decrypts AES-encrypted data using the provided key and returns the original data
func DECRYPT_DATA(encryptedData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}
	if len(encryptedData) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := encryptedData[:aes.BlockSize]
	data := encryptedData[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return data, nil
}

// DATA_ACCESS_CONTROL manages access permissions for a data record
func DATA_ACCESS_CONTROL(dataID string, userID string, permissionLevel string) error {
	accessControl := common.DataAccessControl{
		DataID:         dataID,
		UserID:         userID,
		PermissionLevel: permissionLevel,
		SetAt:          time.Now(),
	}
	return common.SaveDataAccessControl(accessControl)
}

// VERIFY_PROVENANCE checks if the data's current state matches its provenance record for integrity
func VERIFY_PROVENANCE(dataID string, currentHash string) (bool, error) {
	provenanceRecord, err := common.FetchProvenanceRecord(dataID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch provenance record: %v", err)
	}
	return provenanceRecord.DataHash == currentHash, nil
}

// SECURE_DATA_TRANSFER facilitates secure transfer of data by creating an encrypted transmission record
func SECURE_DATA_TRANSFER(data []byte, key []byte, recipientID string) ([]byte, error) {
	encryptedData, err := ENCRYPT_DATA(data, key)
	if err != nil {
		return nil, err
	}
	transferLog := common.SecureTransferLog{
		RecipientID:  recipientID,
		TransferredAt: time.Now(),
		DataHash:      hex.EncodeToString(sha256.New().Sum(encryptedData)),
	}
	if err := common.SaveSecureTransferLog(transferLog); err != nil {
		return nil, fmt.Errorf("failed to log data transfer: %v", err)
	}
	return encryptedData, nil
}

// ANONYMIZE_DATA replaces identifiable data with anonymized placeholders
func ANONYMIZE_DATA(data map[string]string, fieldsToAnonymize []string) map[string]string {
	anonymizedData := make(map[string]string)
	for key, val := range data {
		if contains(fieldsToAnonymize, key) {
			anonymizedData[key] = "ANONYMIZED"
		} else {
			anonymizedData[key] = val
		}
	}
	return anonymizedData
}

// GENERATE_DATA_HASH creates a SHA-256 hash for a given data payload
func GENERATE_DATA_HASH(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// REVOKE_DATA_ACCESS revokes access for a specific user to a data record
func REVOKE_DATA_ACCESS(dataID string, userID string) error {
	revokeLog := common.AccessRevocation{
		DataID:    dataID,
		UserID:    userID,
		RevokedAt: time.Now(),
	}
	return common.SaveAccessRevocation(revokeLog)
}

// PROVENANCE_AUDIT_CHECK performs an audit check on the data’s provenance to verify integrity
func PROVENANCE_AUDIT_CHECK(dataID string) ([]common.AuditRecord, error) {
	auditRecords, err := common.FetchProvenanceAuditLogs(dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch audit logs: %v", err)
	}
	return auditRecords, nil
}

// SET_DATA_ACCESS_EXPIRATION sets an expiration date for data access, after which access will be revoked
func SET_DATA_ACCESS_EXPIRATION(dataID string, userID string, expirationDate time.Time) error {
	expiration := common.DataAccessExpiration{
		DataID:         dataID,
		UserID:         userID,
		ExpirationDate: expirationDate,
		SetAt:          time.Now(),
	}
	return common.SaveDataAccessExpiration(expiration)
}

// MASK_SENSITIVE_INFORMATION masks sensitive data fields with placeholder text
func MASK_SENSITIVE_INFORMATION(data map[string]string, fieldsToMask []string) map[string]string {
	maskedData := make(map[string]string)
	for key, val := range data {
		if contains(fieldsToMask, key) {
			maskedData[key] = "MASKED"
		} else {
			maskedData[key] = val
		}
	}
	return maskedData
}

// UNMASK_SENSITIVE_INFORMATION retrieves the original data by accessing secure storage or logs
func UNMASK_SENSITIVE_INFORMATION(dataID string, fieldsToUnmask []string) (map[string]string, error) {
	originalData, err := common.FetchOriginalDataRecord(dataID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve original data: %v", err)
	}
	unmaskedData := make(map[string]string)
	for key, val := range originalData {
		if contains(fieldsToUnmask, key) {
			unmaskedData[key] = val
		} else {
			unmaskedData[key] = "MASKED"
		}
	}
	return unmaskedData, nil
}

// REGISTER_ACCESS_LOG logs access events for auditing and compliance purposes
func REGISTER_ACCESS_LOG(dataID string, userID string) error {
	accessLog := common.AccessLog{
		DataID:    dataID,
		UserID:    userID,
		AccessedAt: time.Now(),
	}
	return common.SaveAccessLog(accessLog)
}

// VALIDATE_ACCESS_POLICY checks if a user’s access meets specified policy requirements
func VALIDATE_ACCESS_POLICY(dataID string, userID string, requiredPermission string) (bool, error) {
	accessControl, err := common.FetchDataAccessControl(dataID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch access control record: %v", err)
	}
	return accessControl.PermissionLevel == requiredPermission, nil
}

// Helper function: contains checks if a slice contains a specified element
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}
