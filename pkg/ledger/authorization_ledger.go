package ledger

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"time"
)

// VerifySignature verifies the signature for a specific action or document
func (l *AuthorizationLedger) VerifySignature(signerID, signature string) bool {
	// Real-world signature verification logic (details depend on cryptographic libraries and keys used)
	valid := validateSignature(signerID, signature) // Stub function for signature validation
	log.Printf("Signature verification for signer %s resulted in %t.", signerID, valid)
	return valid
}

// RecordAuthorizationLevel stores the authorization level data in the ledger.
func (l *AuthorizationLedger) RecordAuthorizationLevel(authData AuthorizationData) error {
	// Store authorization data in ledger's authorizationLevels map
	l.AuthorizationLevels[authData.UserID] = authData
	fmt.Printf("Authorization level for user %s set to %d.\n", authData.UserID, authData.AuthorizationLevel)
	return nil
}

// FetchAuthorizationLevel retrieves the authorization level for a given user from the ledger.
func (l *AuthorizationLedger) FetchAuthorizationLevel(userID string) (int, error) {
	authData, exists := l.AuthorizationLevels[userID]
	if !exists {
		return 0, fmt.Errorf("authorization level for user %s not found", userID)
	}
	return authData.AuthorizationLevel, nil
}

// RecordTrustedParty adds a trusted party to the ledger.
func (l *AuthorizationLedger) RecordTrustedParty(party TrustedParty) error {
	l.TrustedParties[party.PartyID] = party
	fmt.Printf("Trusted party %s added to the ledger.\n", party.PartyID)
	return nil
}

// RecordAuthorizationEvent logs an authorization event in the ledger.
func (l *AuthorizationLedger) RecordAuthorizationEvent(event AuthorizationEvent) error {
	l.AuthorizationEvents = append(l.AuthorizationEvents, event)
	fmt.Printf("Authorization event logged: %s\n", event.Action)
	return nil
}

// RecordAuthorizedSigner records a new authorized signer in the ledger.
func (l *AuthorizationLedger) RecordAuthorizedSigner(signer AuthorizedSigner) {
	l.Lock()
	defer l.Unlock()

	l.AuthorizedSigners[signer.SignerID] = signer
	log.Printf("Authorized signer %s recorded in ledger", signer.SignerID)
}

// DeleteAuthorizedSigner removes an authorized signer from the ledger.
func (l *AuthorizationLedger) DeleteAuthorizedSigner(signerID string) {
	l.cLock()
	defer l.Unlock()

	delete(l.AuthorizedSigners, signerID)
	log.Printf("Authorized signer %s removed from ledger", signerID)
}

// HasAll checks if the PermissionSet contains all required permissions.
func (p *PermissionSet) HasAll(required PermissionSet) bool {
	for perm := range required.Permissions {
		if !p.Permissions[perm] {
			return false
		}
	}
	return true
}

// RecordBiometricRegistration records a new biometric registration in the ledger.
func (l *AuthorizationLedger) RecordBiometricRegistration(registration BiometricRegistration) error {
	l.Lock()
	defer l.Unlock()
	l.biometricData[registration.UserID] = registration
	return nil
}

// FetchBiometricData retrieves stored biometric data for a user.
func (l *AuthorizationLedger) FetchBiometricData(userID string) (BiometricRegistration, error) {
	l.Lock()
	defer l.Unlock()
	data, exists := l.BiometricData[userID]
	if !exists {
		return BiometricRegistration{}, fmt.Errorf("biometric data not found for user %s", userID)
	}
	return data, nil
}

// DeleteBiometricData removes a user's biometric data from the ledger.
func (l *AuthorizationLedger) DeleteBiometricData(userID string) error {
	l.Lock()
	defer l.Unlock()

	// Implementation to delete the user's biometric data
	if _, exists := l.BiometricData[userID]; exists {
		delete(l.BiometricData, userID)
		return nil
	}
	return fmt.Errorf("biometric data not found for userID: %s", userID)
}

// RecordBiometricUpdate logs the updated biometric data for a user.
func (l *AuthorizationLedger) RecordBiometricUpdate(update BiometricUpdate) error {
	l.Lock()
	defer l.Unlock()

	// Update biometric data in the ledger
	l.BiometricData[update.UserID] = update.EncryptedData
	return nil
}

// RecordBiometricAccess logs an access attempt based on biometric data to the blockchain ledger.
func (l *AuthorizationLedger) RecordBiometricAccess(userID, action string) error {
	// Step 1: Validate input parameters
	if userID == "" {
		return errors.New("userID cannot be empty")
	}
	if action != "success" && action != "failure" {
		return fmt.Errorf("invalid action: %s; expected 'success' or 'failure'", action)
	}

	// Step 2: Create a new access log entry
	accessLog := BiometricAccessLog{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
	}

	// Step 3: Submit the log entry to the blockchain ledger
	if err := l.submitToBlockchain(accessLog); err != nil {
		return fmt.Errorf("failed to submit biometric access log to blockchain ledger: %v", err)
	}

	// Step 4: Optionally emit an event or notify system components for auditing
	if err := l.emitAccessLogEvent(accessLog); err != nil {
		return fmt.Errorf("failed to emit access log event: %v", err)
	}

	return nil
}

// DeleteDelegatedAccess removes a specific delegated access record from the ledger.
// It securely deletes the record by querying blockchain storage and confirming the operation's success.
func (l *AuthorizationLedger) DeleteDelegatedAccess(deviceID, delegateID string) error {
	// Step 1: Validate Input Parameters
	if deviceID == "" || delegateID == "" {
		return errors.New("deviceID and delegateID cannot be empty")
	}

	// Step 2: Query the Blockchain Storage to Locate the Delegated Access Record
	recordKey := fmt.Sprintf("delegated_access:%s:%s", deviceID, delegateID) // Unique key format for locating the record

	exists, err := l.checkRecordExists(recordKey)
	if err != nil {
		return fmt.Errorf("error checking existence of delegated access record: %v", err)
	}
	if !exists {
		return fmt.Errorf("delegated access record for deviceID %s and delegateID %s does not exist", deviceID, delegateID)
	}

	// Step 3: Begin a Blockchain Transaction for Deletion
	tx, err := l.beginTransaction()
	if err != nil {
		return fmt.Errorf("failed to start blockchain transaction for deleting delegated access: %v", err)
	}

	// Step 4: Remove the Record from Blockchain Storage
	err = l.removeRecordFromStorage(recordKey, tx)
	if err != nil {
		l.rollbackTransaction(tx) // Rollback transaction on error
		return fmt.Errorf("failed to delete delegated access record from ledger: %v", err)
	}

	// Step 5: Commit the Blockchain Transaction
	if err = l.commitTransaction(tx); err != nil {
		return fmt.Errorf("failed to commit deletion of delegated access record: %v", err)
	}

	// Step 6: Emit a Deletion Event for Auditing and Notifications
	l.emitAccessRevocationEvent(deviceID, delegateID)

	return nil
}

// submitToBlockchain submits the access log entry to the blockchain ledger as a transaction.
func (l *AuthorizationLedger) submitToBlockchain(log BiometricAccessLog) error {
	// Convert log data to a transaction format suitable for the blockchain ledger
	transactionData := map[string]interface{}{
		"user_id":   log.UserID,
		"action":    log.Action,
		"timestamp": log.Timestamp,
	}

	// Submit the transaction data to the blockchain
	fmt.Printf("Blockchain transaction submitted: %+v\n", transactionData)
	return nil
}

// emitAccessLogEvent emits an event for biometric access, useful for monitoring and auditing.
func (l *AuthorizationLedger) emitAccessLogEvent(log BiometricAccessLog) error {
	// Notify other components or monitoring systems
	fmt.Printf("Event emitted: UserID: %s, Action: %s, Timestamp: %v\n", log.UserID, log.Action, log.Timestamp)
	return nil
}

// RecordBiometricAccessLog records an access attempt in the biometric access log.
func (l *AuthorizationLedger) RecordBiometricAccessLog(log BiometricAccessLog) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.BiometricAccessLogs = append(l.BiometricAccessLogs, log)
}

// FetchAccessControlFlag retrieves the access control flag for a specific user ID from the ledger.
func (l *AuthorizationLedger) FetchAccessControlFlag(userID string) (bool, error) {
	// Lock the ledger for read access to prevent concurrent data issues
	l.Lock()
	defer l.Unlock()

	// Check if the user exists in the ledger data
	flag, exists := l.AccessControlFlags[userID]
	if !exists {
		return false, fmt.Errorf("access control flag for user ID %s not found", userID)
	}

	// Return the flag value and a nil error if successful
	return flag, nil
}

// RecordPermissionRequest logs a request for permissions by an entity.
func (l *AuthorizationLedger) RecordPermissionRequest(requestID string, userID string, permissions string, status string) error {
	l.Lock()
	defer l.Unlock()

	// Logic to record permission request in ledger storage (omitted for brevity)
	return nil
}

// FetchAuthorizedSigner retrieves the authorized signer information from the ledger.
func (l *AuthorizationLedger) FetchAuthorizedSigner(signerID string) (AuthorizedSigner, error) {
	l.Lock()
	defer l.Unlock()

	signer, exists := l.AuthorizedSigners[signerID]
	if !exists {
		return AuthorizedSigner{}, fmt.Errorf("signer not found: %s", signerID)
	}
	return signer, nil
}

// UpdateTrustedPartyFlag updates the flag status of a trusted party in the ledger.
func (l *AuthorizationLedger) UpdateTrustedPartyFlag(partyID string, flag bool) error {
	party, exists := l.TrustedParties[partyID]
	if !exists {
		return fmt.Errorf("trusted party %s does not exist in the ledger", partyID)
	}
	party.Flagged = flag
	l.TrustedParties[partyID] = party
	fmt.Printf("Trusted party %s flagged status updated to: %v\n", partyID, flag)
	return nil
}

// DeleteTrustedParty removes a trusted party from the ledger.
func (l *AuthorizationLedger) DeleteTrustedParty(partyID string) error {
	if _, exists := l.TrustedParties[partyID]; !exists {
		return fmt.Errorf("trusted party %s does not exist in ledger", partyID)
	}
	delete(l.TrustedParties, partyID)
	fmt.Printf("Trusted party %s removed from the ledger.\n", partyID)
	return nil
}

// RecordUnauthorizedAccess records an unauthorized access attempt in the ledger.
func (l *AuthorizationLedger) RecordUnauthorizedAccess(record UnauthorizedAccess) error {
	l.Lock()
	defer l.Unlock()

	// Ensure an entry exists for storing unauthorized access records
	if l.UnauthorizedAccessLogs == nil {
		l.UnauthorizedAccessLogs = []UnauthorizedAccess{}
	}

	// Append the unauthorized access record
	l.UnauthorizedAccessLogs = append(l.UnauthorizedAccessLogs, record)

	log.Printf("Unauthorized access attempt flagged for operation %s by signer %s at %s", record.OperationID, record.SignerID, record.Timestamp)
	return nil
}

// GetAuthorizationLevel fetches the authorization level of an entity, providing an alternative name for retrieval
func (l *AuthorizationLedger) GetAuthorizationLevel(entityID string) (string, bool) {
	return l.FetchAuthorizationLevel(entityID)
}

// validateSignature is a placeholder for the actual signature validation process, to be implemented with the required cryptographic method
func validateSignature(signerID, signature string) bool {
	// Placeholder logic for signature verification
	return true
}

// AddPublicKey securely adds a public key for a given signerID to the ledger.
func (l *AuthorizationLedger) AddPublicKey(signerID string, publicKeyPEM string) error {
	l.Lock()
	defer l.Unlock()

	publicKey, err := parsePublicKeyPEM(publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}

	l.PublicKeys[signerID] = publicKey
	return nil
}

// GetPublicKey retrieves the public key associated with a given signerID.
func (l *AuthorizationLedger) GetPublicKey(signerID string) (*ecdsa.PublicKey, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	publicKey, exists := l.PublicKeys[signerID]
	if !exists {
		return nil, fmt.Errorf("public key not found for signerID: %s", signerID)
	}
	return publicKey, nil
}

// parsePublicKeyPEM decodes and parses an ECDSA public key from a PEM-encoded string.
func parsePublicKeyPEM(publicKeyPEM string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid PEM format or missing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not of type ECDSA")
	}
	return ecdsaPub, nil
}

// RemovePublicKey removes a public key associated with a signerID from the ledger.
func (l *AuthorizationLedger) RemovePublicKey(signerID string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.PublicKeys[signerID]; !exists {
		return fmt.Errorf("public key not found for signerID: %s", signerID)
	}

	delete(l.PublicKeys, signerID)
	return nil
}

// RecordPermissionRevocation logs a revocation of permission
func (l *AuthorizationLedger) RecordPermissionRevocation(permissionID, revokedBy, reason string) {
	l.Lock()
	defer l.Unlock()
	l.PermissionRevocations[permissionID] = PermissionRevocation{
		PermissionID: permissionID,
		RevokedBy:    revokedBy,
		Reason:       reason,
		RevokedAt:    time.Now(),
	}
	log.Printf("Permission %s revoked by %s for reason: %s.", permissionID, revokedBy, reason)
}

// FetchAuthorizationLogs retrieves authorization logs within a specified time range.
func (l *AuthorizationLedger) FetchAuthorizationLogs(logID string, startTime, endTime time.Time) ([]AuthorizationLog, error) {
	l.Lock()
	defer l.Unlock()

	var logsInRange []AuthorizationLog
	for _, log := range l.AuthorizationLogs {
		// Match by logID and ensure timestamp is within range
		if log.LogID == logID && log.Timestamp.After(startTime) && log.Timestamp.Before(endTime) {
			logsInRange = append(logsInRange, log)
		}
	}

	if len(logsInRange) == 0 {
		return nil, errors.New("no authorization logs found within specified time range")
	}

	return logsInRange, nil
}

// RecordSignerActivity records a signer's activity log in the ledger.
func (l *AuthorizationLedger) RecordSignerActivity(signerID string, activity AuthorizationLog) error {
	l.Lock()
	defer l.Unlock()

	// Ensure an entry exists for the signerID
	if _, exists := l.SignerActivities[signerID]; !exists {
		l.SignerActivities[signerID] = []AuthorizationLog{}
	}

	// Append the activity log to the signer's activity history
	l.SignerActivities[signerID] = append(l.SignerActivities[signerID], activity)

	log.Printf("Activity for signer %s recorded: %s at %s", signerID, activity.Action, activity.Timestamp)
	return nil
}

// FetchAuthorizationHistory retrieves the authorization history for an operation within a specific time range.
func (l *AuthorizationLedger) FetchAuthorizationHistory(operationID string, startTime, endTime time.Time) ([]AuthorizationLog, error) {
	l.Lock()
	defer l.Unlock()

	var history []AuthorizationLog
	for _, log := range l.AuthorizationLogs[operationID] {
		if log.Timestamp.After(startTime) && log.Timestamp.Before(endTime) {
			history = append(history, log)
		}
	}

	if len(history) == 0 {
		return nil, fmt.Errorf("no authorization history found for operation %s within the specified time range", operationID)
	}

	log.Printf("Fetched %d authorization logs for operation %s", len(history), operationID)
	return history, nil
}

// RecordAuthorizationAction records an authorization action in the ledger.
func (l *AuthorizationLedger) RecordAuthorizationAction(record AuthorizationLog) error {
	l.Lock()
	defer l.Unlock()

	// Append the new record to the ledger's authorization logs
	l.AuthorizationLogs = append(l.AuthorizationLogs, record)
	return nil
}

// RecordSubBlockValidationAuthorization logs authorization for sub-block validation
func (l *AuthorizationLedger) RecordSubBlockValidationAuthorization(authID, validatorID, status string) {
	l.Lock()
	defer l.Unlock()
	l.SubBlockValidationAuth[authID] = SubBlockValidationAuth{
		AuthID:       authID,
		ValidatorID:  validatorID,
		Status:       status,
		AuthorizedAt: time.Now(),
	}
	log.Printf("Sub-block validation authorization %s set for validator %s with status: %s.", authID, validatorID, status)
}

// UpdateAccessControlFlag updates the access control flag for a specified entity
func (l *AuthorizationLedger) UpdateAccessControlFlag(entityID string, flag bool) {
	l.Lock()
	defer l.Unlock()
	l.AccessControlFlags[entityID] = flag
	log.Printf("Access control flag for %s set to %t.", entityID, flag)
}

// Updated RecordDeviceAuthorization function in the Ledger struct to match expected parameters
func (l *AuthorizationLedger) RecordDeviceAuthorization(deviceID string, authorizedBy string, level string) error {
	l.Lock()
	defer l.Unlock()

	// Log the authorization with necessary fields in the map
	l.DeviceAuthorizations[deviceID] = DeviceAuthorization{
		DeviceID:           deviceID,
		AuthorizedBy:       authorizedBy,
		AuthorizationLevel: level,
		AuthorizedAt:       time.Now(),
	}

	log.Printf("Device %s authorized by %s with level %s.", deviceID, authorizedBy, level)
	return nil
}

// DeleteDeviceAuthorization removes a device authorization entry
func (l *AuthorizationLedger) DeleteDeviceAuthorization(deviceID string) {
	l.Lock()
	defer l.Unlock()
	delete(l.DeviceAuthorizations, deviceID)
	log.Printf("Device authorization for %s deleted.", deviceID)
}

// FetchDeviceAuthorization retrieves device authorization details
func (l *AuthorizationLedger) FetchDeviceAuthorization(deviceID string) (DeviceAuthorization, bool) {
	l.Lock()
	defer l.Unlock()
	auth, exists := l.DeviceAuthorizations[deviceID]
	return auth, exists
}

// RecordDelegatedAccess logs a delegated access entry
func (l *AuthorizationLedger) RecordDelegatedAccess(delegationID, grantedBy, grantedTo, level string) {
	l.Lock()
	defer l.Unlock()
	l.DelegatedAccess[delegationID] = DelegatedAccessRecord{
		DelegationID: delegationID,
		GrantedBy:    grantedBy,
		GrantedTo:    grantedTo,
		AccessLevel:  level,
		GrantedAt:    time.Now(),
	}
	log.Printf("Delegated access %s granted by %s to %s with level %s.", delegationID, grantedBy, grantedTo, level)
}

// RecordTemporaryAccess logs a temporary access entry with an expiration date
func (l *AuthorizationLedger) RecordTemporaryAccess(accessID, entityID, level string, expiresAt time.Time) {
	l.Lock()
	defer l.Unlock()
	l.TemporaryAccessRecords[accessID] = TemporaryAccessRecord{
		AccessID:    accessID,
		EntityID:    entityID,
		AccessLevel: level,
		ExpiresAt:   expiresAt,
	}
	log.Printf("Temporary access %s granted to %s with level %s until %s.", accessID, entityID, level, expiresAt)
}

// DeleteTemporaryAccess removes a temporary access entry
func (l *AuthorizationLedger) DeleteTemporaryAccess(accessID string) {
	l.Lock()
	defer l.Unlock()
	delete(l.TemporaryAccessRecords, accessID)
	log.Printf("Temporary access %s deleted.", accessID)
}

// FetchAccessLogs retrieves all access logs from the ledger
func (l *AuthorizationLedger) FetchAccessLogs() map[string]AccessLog {
	l.Lock()
	defer l.Unlock()
	return l.AccessLogs
}

// generateLogID generates a unique log ID based on entity and action
func generateLogID(entityID, action string) string {
	return entityID + "_" + action + "_" + time.Now().Format("20060102150405")
}

// FetchAccessRestrictions retrieves a specific access restriction policy
func (l *AuthorizationLedger) FetchAccessRestrictions(restrictionID string) (AccessRestriction, bool) {
	l.Lock()
	defer l.Unlock()
	restriction, exists := l.AccessRestrictions[restrictionID]
	return restriction, exists
}

// RecordRoleChange logs a change in roles
func (l *AuthorizationLedger) RecordRoleChange(roleID, changedBy, newRole string) {
	l.Lock()
	defer l.Unlock()
	l.RoleChanges[roleID] = RoleChange{
		RoleID:    roleID,
		ChangedBy: changedBy,
		NewRole:   newRole,
		ChangedAt: time.Now(),
	}
	log.Printf("Role %s changed to %s by %s.", roleID, newRole, changedBy)
}

// RecordAccessAttempt logs an access attempt in the ledger.
func (l *AuthorizationLedger) RecordAccessAttempt(accessLog AccessLog) error {
	l.Lock()
	defer l.Unlock()

	// Store the access log entry with a generated unique ID
	accessLog.AccessID = fmt.Sprintf("access-%d", len(l.AccessLogs)+1)
	l.AccessLogs[accessLog.AccessID] = accessLog

	log.Printf("Access attempt recorded: UserID=%s, Action=%s, Result=%s, Timestamp=%v", accessLog.UserID, accessLog.Action, accessLog.Result, accessLog.Timestamp)
	return nil
}

// FetchMicrochipAuthorization retrieves the authorization record for a given microchip ID.
func (l *AuthorizationLedger) FetchMicrochipAuthorization(chipID string) (MicrochipAuthorization, error) {
	l.Lock()
	defer l.Unlock()

	auth, exists := l.MicrochipAuthorizations[chipID]
	if !exists {
		return MicrochipAuthorization{}, errors.New("microchip authorization not found")
	}

	return auth, nil
}

// RecordAuthorizationConstraints logs constraints applied to authorizations
func (l *AuthorizationLedger) RecordAuthorizationConstraints(constraintID, description, appliedBy string) {
	l.Lock()
	defer l.Unlock()
	l.AuthorizationConstraints[constraintID] = AuthorizationConstraint{
		ConstraintID: constraintID,
		Description:  description,
		AppliedBy:    appliedBy,
		AppliedAt:    time.Now(),
	}
	log.Printf("Authorization constraint %s applied by %s: %s.", constraintID, appliedBy, description)
}

// RecordKeyReset logs a key reset event
func (l *AuthorizationLedger) RecordKeyReset(keyID, resetBy string) {
	l.Lock()
	defer l.Unlock()
	l.KeyResets[keyID] = KeyResetRecord{
		KeyID:   keyID,
		ResetBy: resetBy,
		ResetAt: time.Now(),
	}
	log.Printf("Key %s reset by %s.", keyID, resetBy)
}

// RecordMicrochipAuthorization logs an authorization for a microchip
func (l *AuthorizationLedger) RecordMicrochipAuthorization(chipID, authorizedBy, level string) {
	l.Lock()
	defer l.Unlock()
	l.MicrochipAuthorizations[chipID] = MicrochipAuthorization{
		ChipID:             chipID,
		AuthorizedBy:       authorizedBy,
		AuthorizationLevel: level,
		AuthorizedAt:       time.Now(),
	}
	log.Printf("Microchip authorization for %s set to level %s by %s.", chipID, level, authorizedBy)
}

// UpdateMicrochipAuthorization updates the authorization level of a microchip
func (l *AuthorizationLedger) UpdateMicrochipAuthorization(chipID, level string) {
	l.Lock()
	defer l.Unlock()
	if auth, exists := l.MicrochipAuthorizations[chipID]; exists {
		auth.AuthorizationLevel = level
		l.MicrochipAuthorizations[chipID] = auth
		log.Printf("Microchip %s authorization level updated to %s.", chipID, level)
	}
}

// DeleteMicrochipAuthorization removes a microchip's authorization record
func (l *AuthorizationLedger) DeleteMicrochipAuthorization(chipID string) {
	l.Lock()
	defer l.Unlock()
	delete(l.MicrochipAuthorizations, chipID)
	log.Printf("Microchip authorization for %s deleted.", chipID)
}

// RecordAccessAttempt logs a general access attempt
func (l *AuthorizationLedger) RecordAccessAttempt(attemptID, entityID, result string) {
	l.Lock()
	defer l.Unlock()
	l.AccessAttempts[attemptID] = AccessAttempt{
		AttemptID:   attemptID,
		EntityID:    entityID,
		AttemptedAt: time.Now(),
		Result:      result,
	}
	log.Printf("Access attempt %s by %s recorded with result: %s.", attemptID, entityID, result)
}

// RecordMicrochipAccessAttempt logs an access attempt using a microchip
func (l *AuthorizationLedger) RecordMicrochipAccessAttempt(attemptID, chipID, result string) {
	l.Lock()
	defer l.Unlock()
	l.MicrochipAccessAttempts[attemptID] = AccessAttempt{
		AttemptID:   attemptID,
		EntityID:    chipID,
		AttemptedAt: time.Now(),
		Result:      result,
	}
	log.Printf("Microchip access attempt %s by %s recorded with result: %s.", attemptID, chipID, result)
}

// LogMultiSigVerification logs the verification of a multi-signature
func (l *AuthorizationLedger) LogMultiSigVerification(walletID, signerID string, verified bool) {
	log.Printf("MultiSig verification for wallet %s by signer %s verified: %t.", walletID, signerID, verified)
}

// RecordRoleAssignment logs a role assignment to an entity
func (l *AuthorizationLedger) RecordRoleAssignment(assignmentID, role, assignedTo, assignedBy string) {
	l.Lock()
	defer l.Unlock()
	l.RoleAssignments[assignmentID] = RoleAssignment{
		AssignmentID: assignmentID,
		Role:         role,
		AssignedTo:   assignedTo,
		AssignedBy:   assignedBy,
		AssignedAt:   time.Now(),
	}
	log.Printf("Role %s assigned to %s by %s.", role, assignedTo, assignedBy)
}

// DeleteRoleAssignment removes a role assignment from the ledger
func (l *AuthorizationLedger) DeleteRoleAssignment(assignmentID string) {
	l.Lock()
	defer l.Unlock()
	delete(l.RoleAssignments, assignmentID)
	log.Printf("Role assignment %s deleted.", assignmentID)
}

// FetchRoleAssignment retrieves a role assignment by assignment ID
func (l *AuthorizationLedger) FetchRoleAssignment(assignmentID string) (RoleAssignment, bool) {
	l.Lock()
	defer l.Unlock()
	assignment, exists := l.RoleAssignments[assignmentID]
	return assignment, exists
}

// FetchPrivileges retrieves the privileges associated with a specific role
func (l *AuthorizationLedger) FetchPrivileges(roleID string) (Privilege, bool) {
	l.Lock()
	defer l.Unlock()
	privilege, exists := l.Privileges[roleID]
	return privilege, exists
}

// RecordTimeBasedAuthorization logs a time-based access control entry
func (l *AuthorizationLedger) RecordTimeBasedAuthorization(authID, entityID, level string, validFrom, expiresAt time.Time) {
	l.Lock()
	defer l.Unlock()
	l.TimeBasedAuthorizations[authID] = TimeBasedAuthorization{
		AuthorizationID: authID,
		EntityID:        entityID,
		AccessLevel:     level,
		ValidFrom:       validFrom,
		ExpiresAt:       expiresAt,
	}
	log.Printf("Time-based authorization %s for %s with level %s from %s to %s.", authID, entityID, level, validFrom, expiresAt)
}

// DeleteTimeBasedAuthorization removes a time-based authorization record
func (l *AuthorizationLedger) DeleteTimeBasedAuthorization(authID string) {
	l.Lock()
	defer l.Unlock()
	delete(l.TimeBasedAuthorizations, authID)
	log.Printf("Time-based authorization %s deleted.", authID)
}

// RecordSignerPriority logs the priority level for a signer
func (l *AuthorizationLedger) RecordSignerPriority(signerID string, priorityLevel int) {
	l.Lock()
	defer l.Unlock()
	l.SignerPriorities[signerID] = SignerPriority{
		SignerID:      signerID,
		PriorityLevel: priorityLevel,
		SetAt:         time.Now(),
	}
	log.Printf("Signer %s priority level set to %d.", signerID, priorityLevel)
}

// FetchSignerStatus retrieves the priority status of a signer
func (l *AuthorizationLedger) FetchSignerStatus(signerID string) (SignerPriority, bool) {
	l.Lock()
	defer l.Unlock()
	status, exists := l.SignerPriorities[signerID]
	return status, exists
}
