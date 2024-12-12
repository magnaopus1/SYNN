package ledger

import (
	"errors"
	"fmt"
	"time"
)

// GetIdentityLogs retrieves identity-related logs, filtered by NodeID or Action.
func (l *IdentityLedger ) GetIdentityLogs(nodeID, actionFilter string) []IdentityLog {
	l.Lock()
	defer l.Unlock()

	var filteredLogs []IdentityLog
	for _, log := range l.IdentityLogs {
		if (nodeID == "" || log.NodeID == nodeID) && (actionFilter == "" || log.Action == actionFilter) {
			filteredLogs = append(filteredLogs, log)
		}
	}
	return filteredLogs
}



// PurgeOldIdentityLogs removes identity logs older than the specified duration.
func (l *IdentityLedger ) PurgeOldIdentityLogs(duration time.Duration) {
	l.Lock()
	defer l.Unlock()

	cutoff := time.Now().Add(-duration)
	var recentLogs []IdentityLog
	for _, log := range l.IdentityLogs {
		if log.Timestamp.After(cutoff) {
			recentLogs = append(recentLogs, log)
		}
	}
	l.IdentityLogs = recentLogs
	fmt.Printf("Purged identity logs older than %v.\n", duration)
}

// RecordIdentityCreation records the creation of a new identity in the ledger.
func (l *IdentityLedger ) RecordIdentityCreation(identityID string, identity Identity) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.Identities[identityID]; exists {
		return fmt.Errorf("identity %s already exists", identityID)
	}

	l.Identities[identityID] = &identity
	log := IdentityLog{
		NodeID:    identityID,
		Action:    "Identity created",
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Owner: %s, Type: %s", identity.Owner, identity.IdentityType),
	}
	l.IdentityLogs = append(l.IdentityLogs, log)
	fmt.Printf("Identity %s created.\n", identityID)
	return nil
}


// RecordIdentityVerification marks an identity as verified in the ledger.
func (l *IdentityLedger) RecordIdentityVerification(identityID string) error {
	l.Lock()
	defer l.Unlock()

	identity, exists := l.Identities[identityID]
	if !exists {
		return fmt.Errorf("identity %s does not exist", identityID)
	}

	identity.IsVerified = true
	l.Identities[identityID] = identity

	log := IdentityLog{
		NodeID:    identityID,
		Action:    "Identity verified",
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Identity %s verified", identityID),
	}
	l.IdentityLogs = append(l.IdentityLogs, log)
	fmt.Printf("Identity %s verified.\n", identityID)
	return nil
}


// RecordAccessControlChange logs changes in access control for a specific node or user.
func (l *IdentityLedger ) RecordAccessControlChange(nodeID, action string) {
	l.Lock()
	defer l.Unlock()

	accessLog := IdentityLog{
		NodeID:    nodeID,
		Action:    action,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Access control changed: %s", action),
	}
	l.IdentityLogs = append(l.IdentityLogs, accessLog)
	fmt.Printf("Access control change recorded for node %s: %s\n", nodeID, action)
}


// RecordPrivacyChange records a change in privacy settings for a user.
func (l *IdentityLedger ) RecordPrivacyChange(userID string, newSettings PrivacySettings) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.PrivacyManager.PrivacyRecords[userID]; !exists {
		return fmt.Errorf("user %s does not exist", userID)
	}

	l.PrivacyManager.PrivacyRecords[userID] = &newSettings
	log := IdentityLog{
		NodeID:    userID,
		Action:    "Privacy settings changed",
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("New privacy settings: Encryption=%v, SharePermission=%v", newSettings.DataEncryption, newSettings.PermissionToShare),
	}
	l.IdentityLogs = append(l.IdentityLogs, log)
	fmt.Printf("Privacy settings updated for user %s.\n", userID)
	return nil
}


// GetPrivacySettings retrieves the privacy settings for a user.
func (l *IdentityLedger) GetPrivacySettings(userID string) (*PrivacySettings, error) {
	l.Lock()
	defer l.Unlock()

	settings, exists := l.PrivacyManager.PrivacyRecords[userID]
	if !exists {
		return nil, fmt.Errorf("privacy settings for user %s not found", userID)
	}

	fmt.Printf("Privacy settings retrieved for user %s.\n", userID)
	return settings, nil
}

// RecordIdentityProofSubmission logs the submission of identity proof.
func (l *IdentityLedger) RecordIdentityProofSubmission(participant, proofHash string) error {
	l.Lock()
	defer l.Unlock()

	// Log the identity proof submission
	proof := IdentityProof{
		Participant: participant,
		ProofHash:   proofHash,
		Timestamp:   time.Now(),
		Verified:    false, // Initial state of the proof is unverified
	}

	// Store the proof in the ledger
	l.IdentityProofs[participant] = proof
	return nil
}


// RecordIdentityProofVerification logs the verification of an identity proof.
func (l *IdentityLedger) RecordIdentityProofVerification(proofID, verifierID string) error {
	l.Lock()
	defer l.Unlock()

	// Verify the identity proof
	proof, exists := l.IdentityProofs[proofID]
	if !exists {
		return errors.New("proof not found")
	}

	// Mark the proof as verified
	proof.Verified = true
	proof.Timestamp = time.Now()

	// Update the proof in the ledger
	l.IdentityProofs[proofID] = proof
	return nil
}




// RecordAccessChange logs any changes made to access control for nodes or accounts.
func (l *IdentityLedger) RecordAccessChange(nodeID, action string) {
	l.Lock()
	defer l.Unlock()

	accessLog := IdentityLog{
		NodeID:    nodeID,
		Action:    action,
		Timestamp: time.Now(),
	}
	l.IdentityLogs = append(l.IdentityLogs, accessLog)
	fmt.Printf("Access change recorded for node %s: %s\n", nodeID, action)
}


// RecordIdentityAction logs identity actions such as changes in node identity or identity verification.
func (l *IdentityLedger) RecordIdentityAction(nodeID, action string) {
	l.Lock()
	defer l.Unlock()

	identityLog := IdentityLog{
		NodeID:    nodeID,
		Action:    fmt.Sprintf("Identity action: %s", action),
		Timestamp: time.Now(),
	}
	l.IdentityLogs = append(l.IdentityLogs, identityLog)
	fmt.Printf("Identity action recorded for node %s: %s\n", nodeID, action)
}

// RecordPrivacyAction logs actions taken that affect the privacy settings or data privacy of nodes or users.
func (l *IdentityLedger) RecordPrivacyAction(nodeID, action string) {
	l.Lock()
	defer l.Unlock()

	privacyLog := IdentityLog{
		NodeID:    nodeID,
		Action:    fmt.Sprintf("Privacy action: %s", action),
		Timestamp: time.Now(),
	}
	l.IdentityLogs = append(l.IdentityLogs, privacyLog)
	fmt.Printf("Privacy action recorded for node %s: %s\n", nodeID, action)
}