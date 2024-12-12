package authorization

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// AuthorizationLog represents a record of an authorization event within the ledger.
type AuthorizationLog struct {
    LogID       string    // Unique ID for the log entry
    OperationID string    // ID of the operation being approved or queried
    SignerID    string    // ID of the signer or user involved in the action
    Action      string    // Type of action, e.g., "Approved", "Rejected"
    Timestamp   time.Time // Time the action was recorded
    IPAddress   string    // IP address from where the action was initiated
}


// AuditTrail records details about an audit trail for traceability.
type AuditTrail struct {
	TrailID        string    // Unique identifier for the audit trail entry
	EventType      string    // Type of event (e.g., "transaction", "access", "authorization")
	UserID         string    // User ID associated with the action
	NodeID         string    // ID of the node where the event occurred
	Timestamp      time.Time // Timestamp of the event
	ActionDetails  string    // Detailed description of the action
	Status         string    // Current status of the audit trail (e.g., "Enabled", "Disabled")
	OperationID    string    // Operation ID associated with the audit trail
}

// UnauthorizedAccess logs unauthorized access attempts for security monitoring.
type UnauthorizedAccess struct {
    AccessID       string    // Unique identifier for the unauthorized access attempt
    UserID         string    // User ID if identifiable, otherwise could be "unknown"
    NodeID         string    // ID of the node where the unauthorized attempt occurred
    AttemptedAccessLevel string // Access level attempted, e.g., "admin", "write"
    Timestamp      time.Time // Timestamp of the unauthorized access attempt
    IPAddress      string    // IP address from which the attempt originated
    ActionTaken    string    // Action taken in response, e.g., "blocked", "alerted"
    Notes          string    // Additional notes or reason for unauthorized access
}

// QueryAuthorizationLog retrieves authorization logs for a specific operation or within a time period.
func QueryAuthorizationLog(ledgerInstance *ledger.Ledger, logID string, startTime, endTime time.Time) ([]ledger.AuthorizationLog, error) {
    // Fetch logs within the given time range
    logs, err := ledgerInstance.AuthorizationLedger.FetchAuthorizationLogs(logID, startTime, endTime)
    if err != nil {
        return nil, fmt.Errorf("error fetching authorization logs: %v", err)
    }

    return logs, nil
}

// getNodeIPAddress connects briefly to a network access point to determine the outbound IP address.
func getNodeIPAddress(accessPoint string, port string) (string, error) {
    // Formulate the address for the access point (usually a node or central network entry)
    endpoint := net.JoinHostPort(accessPoint, port)

    // Dial a UDP connection to the endpoint, without actually sending data.
    conn, err := net.DialTimeout("udp", endpoint, 2*time.Second)
    if err != nil {
        return "", fmt.Errorf("failed to dial network access point: %v", err)
    }
    defer conn.Close()

    // Extract the local address of the connection, which will contain the node’s outbound IP
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP.String(), nil
}


// DenyOperationWithSignature records a denial in the ledger, verifying the signature and logging the denial event.
func DenyOperationWithSignature(consensus *common.SynnergyConsensus, ledgerInstance *ledger.Ledger, operationID, signerID string, signature []byte) (bool, error) {
	// Verify the signature using Synnergy Consensus
	verified, err := consensus.VerifySignature(ledgerInstance, operationID, signerID, signature)
	if err != nil || !verified {
		return false, errors.New("signature verification failed")
	}

	// Create the denial record using ledger.AuthorizationLog struct
	denialRecord := ledger.AuthorizationLog{
		OperationID: operationID,
		SignerID:    signerID,
		Action:      "Denied",
		Timestamp:   time.Now(),
	}

	// Record the denial action in the ledger
	if err := ledgerInstance.AuthorizationLedger.RecordAuthorizationAction(denialRecord); err != nil {
		return false, errors.New("failed to record denial in ledger")
	}

	return true, nil
}


// EnableAuditTrail initiates an audit trail for the specified operation, logging all authorization attempts.
func EnableAuditTrail(ledgerInstance *ledger.Ledger, operationID string) error {
	auditRecord := ledger.AuditTrail{
		TrailID:     generateTrailID(), // Assume generateTrailID() generates a unique ID for each audit trail
		EventType:   "Authorization",
		UserID:      "system",
		NodeID:      "node1",
		Timestamp:   time.Now(),
		Status:      "Enabled",
		ActionDetails: "Audit trail initiated for authorization attempts",
		OperationID: operationID,
	}

	if err := ledgerInstance.ComplianceLedger.RecordAuditTrail(auditRecord); err != nil {
		return errors.New("failed to enable audit trail")
	}

	return nil
}

// DisableAuditTrail stops the audit trail for a specific operation, marking the end time and storing it in the ledger.
func DisableAuditTrail(ledgerInstance *ledger.Ledger, operationID string) error {
	auditRecord := ledger.AuditTrail{
		TrailID:      generateTrailID(),
		EventType:    "Authorization",
		UserID:       "system",
		NodeID:       "node1",
		Timestamp:    time.Now(),
		Status:       "Disabled",
		ActionDetails: "Audit trail disabled for authorization attempts",
		OperationID:   operationID,
	}

	if err := ledgerInstance.ComplianceLedger.RecordAuditTrail(auditRecord); err != nil {
		return errors.New("failed to disable audit trail")
	}

	return nil
}

// generateTrailID generates a unique ID for an audit trail entry.
func generateTrailID() string {
	// Create a 16-byte random ID
	id := make([]byte, 16)
	_, err := rand.Read(id)
	if err != nil {
		fmt.Println("Error generating trail ID:", err)
		return ""
	}
	// Convert the random bytes to a hex string for readability
	return hex.EncodeToString(id)
}

// TrackSignerActivity logs the activity of a specific signer, recording all actions and storing them in the ledger.
func TrackSignerActivity(ledgerInstance *ledger.Ledger, signerID string, activity ledger.AuthorizationLog) error {
	activity.Timestamp = time.Now()
	if err := ledgerInstance.AuthorizationLedger.RecordSignerActivity(signerID, activity); err != nil {
		return fmt.Errorf("failed to record signer activity in ledger: %v", err)
	}
	return nil
}

// ReviewAuthorizationHistory allows the review of authorization actions for an operation within a specific period.
func ReviewAuthorizationHistory(ledgerInstance *ledger.Ledger, operationID string, startTime, endTime time.Time) ([]ledger.AuthorizationLog, error) {
	history, err := ledgerInstance.AuthorizationLedger.FetchAuthorizationHistory(operationID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve authorization history: %v", err)
	}
	return history, nil
}

// FlagUnauthorizedAccess flags any unauthorized access attempts for an operation, storing them securely in the ledger.
func FlagUnauthorizedAccess(ledgerInstance *ledger.Ledger, operationID, signerID, details string) error {
	unauthorizedAccessRecord := ledger.UnauthorizedAccess{
		OperationID: operationID,
		SignerID:    signerID,
		Details:     details,
		Timestamp:   time.Now(),
	}
	if err := ledgerInstance.AuthorizationLedger.RecordUnauthorizedAccess(unauthorizedAccessRecord); err != nil {
		return fmt.Errorf("failed to flag unauthorized access in ledger: %v", err)
	}
	return nil
}

