package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// SecurityManager handles all security-related operations for the blockchain network
type SecurityManager struct {
	LedgerInstance      *ledger.Ledger      // Ledger for recording security events
	ActiveSessions      map[string]*Session // Active sessions mapped by session ID
	mutex               sync.Mutex          // Mutex for thread-safe operations
	WhitelistedIPs      map[string]bool     // Whitelisted IP addresses
	BlacklistedIPs      map[string]bool     // Blacklisted IP addresses
	NodeKeys            map[string]string   // Holds public keys for authorized nodes
	MasterSecurityKey   string              // Master key for system-level security management
	EncryptionService   *Encryption // Encryption service for securing sensitive data
	FailedLoginAttempts map[string]int      // Tracks failed login attempts by IP or UserID
	AlertThreshold      int                 // Threshold for raising alerts (e.g., failed login attempts)
}

// Session represents an active session within the blockchain network
type Session struct {
	SessionID    string    // Unique ID for the session
	UserID       string    // User ID associated with the session
	StartTime    time.Time // Session start time
	LastActivity time.Time // Last activity timestamp
	IPAddress    string    // IP address associated with the session
	EncryptionKey string   // Encryption key for securing session data
}

// SecurityEvent represents a security-related event logged within the blockchain network
type SecurityEvent struct {
	EventID      string    // Unique ID for the event
	EventType    string    // Type of event (e.g., "LOGIN", "ACCESS_DENIED", "FAILED_LOGIN")
	Timestamp    time.Time // Time when the event occurred
	UserID       string    // User ID associated with the event
	IPAddress    string    // IP address involved in the event
	Details      string    // Additional details regarding the event
}

// LoginAttempt represents an attempt to log into the network
type LoginAttempt struct {
	AttemptID    string    // Unique ID for the login attempt
	UserID       string    // User ID making the login attempt
	Timestamp    time.Time // Time of the login attempt
	IPAddress    string    // IP address from which the attempt was made
	Successful   bool      // Whether the login attempt was successful
}

// IntrusionDetectionSystem handles the detection and prevention of suspicious activities in the network
type IntrusionDetectionSystem struct {
	SuspiciousIPs     map[string]int         // Tracks suspicious IP addresses and their activity count
	BlacklistedIPs    map[string]bool        // Blacklisted IP addresses
	Threshold         int                    // Threshold for flagging an IP as suspicious
	LedgerInstance    *ledger.Ledger         // Ledger for logging detected intrusions
	EncryptionService *Encryption // Encryption service for securing logs
	mutex             sync.Mutex             // Mutex for thread-safe access
}

// SecurityAlert represents a security alert raised within the network
type SecurityAlert struct {
	AlertID       string    // Unique ID for the alert
	AlertType     string    // Type of alert (e.g., "INTRUSION", "FAILED_LOGIN")
	Timestamp     time.Time // Time the alert was raised
	Severity      string    // Severity level of the alert (e.g., "LOW", "MEDIUM", "HIGH")
	Description   string    // Description of the alert
	Resolved      bool      // Whether the alert has been resolved
}

// EncryptionManager handles all encryption and decryption operations in the system
type EncryptionManager struct {
	EncryptionService *Encryption // Encryption service for managing cryptographic operations
	LedgerInstance    *ledger.Ledger         // Ledger for logging encryption-related events
}

// NodeSecurityManager manages node-specific security settings and access control
type NodeSecurityManager struct {
	NodeID            string            // Unique identifier for the node
	NodePublicKey     string            // Public key of the node for secure communication
	Whitelist         map[string]bool   // Whitelisted IP addresses for the node
	Blacklist         map[string]bool   // Blacklisted IP addresses for the node
	EncryptionService *Encryption // Encryption service for securing node data
	LedgerInstance    *ledger.Ledger    // Ledger for recording node security events
}


// NewSecurityManager initializes the security manager with a new ledger instance and master key
func NewSecurityManager(ledgerInstance *ledger.Ledger, masterKey string) *SecurityManager {
    return &SecurityManager{
        LedgerInstance:    ledgerInstance,
        ActiveSessions:    make(map[string]*Session),
        WhitelistedIPs:    make(map[string]bool),
        BlacklistedIPs:    make(map[string]bool),
        NodeKeys:          make(map[string]string),
        MasterSecurityKey: masterKey,
    }
}

// CreateSession starts a new user session and stores it securely
func (sm *SecurityManager) CreateSession(userID string) (string, error) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    sessionID := sm.generateSessionID(userID)
    session := &Session{
        SessionID:    sessionID,
        UserID:       userID,
        StartTime:    time.Now(),
        LastActivity: time.Now(),
    }

    sm.ActiveSessions[sessionID] = session

    fmt.Printf("Session created for user %s with session ID %s.\n", userID, sessionID)

    // Convert session.StartTime to string using a standard format (RFC3339 for example)
    startTimeStr := session.StartTime.Format(time.RFC3339)

    // Pass the formatted time string to RecordSessionStart
    err := sm.LedgerInstance.AdvancedSecurityLedger.RecordSessionStart(sessionID, userID, startTimeStr)
    if err != nil {
        return "", fmt.Errorf("failed to log session to ledger: %v", err)
    }

    return sessionID, nil
}


// ValidateSession checks if a session is active and refreshes its last activity time
func (sm *SecurityManager) ValidateSession(sessionID string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    session, exists := sm.ActiveSessions[sessionID]
    if !exists {
        return errors.New("invalid session ID")
    }

    // Update the last activity timestamp
    session.LastActivity = time.Now()

    fmt.Printf("Session %s validated and activity updated.\n", sessionID)

    return nil
}

// TerminateSession ends an active session and removes it from the session pool
func (sm *SecurityManager) TerminateSession(sessionID string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    session, exists := sm.ActiveSessions[sessionID]
    if !exists {
        return errors.New("session not found")
    }

    // Remove the session from the active session pool
    delete(sm.ActiveSessions, sessionID)

    // Log the session end with the correct number of arguments (assuming sessionID only)
    err := sm.LedgerInstance.AdvancedSecurityLedger.RecordSessionEnd(sessionID) // Adjusted call
    if err != nil {
        return fmt.Errorf("failed to log session termination: %v", err)
    }

    // Optionally use the session for logging the user ID
    fmt.Printf("Session %s for user %s terminated.\n", sessionID, session.UserID)

    return nil
}



// AddNodeKey adds a new node's public key to the list of authorized nodes
func (sm *SecurityManager) AddNodeKey(nodeID, publicKey string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if _, exists := sm.NodeKeys[nodeID]; exists {
        return fmt.Errorf("node %s is already authorized", nodeID)
    }

    sm.NodeKeys[nodeID] = publicKey

    fmt.Printf("Node %s added with public key %s.\n", nodeID, publicKey)

    return nil
}

// RemoveNodeKey removes a node's public key from the list of authorized nodes
func (sm *SecurityManager) RemoveNodeKey(nodeID string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if _, exists := sm.NodeKeys[nodeID]; !exists {
        return fmt.Errorf("node %s is not authorized", nodeID)
    }

    delete(sm.NodeKeys, nodeID)

    fmt.Printf("Node %s removed from the list of authorized nodes.\n", nodeID)

    return nil
}

// AddIPToWhitelist adds an IP address to the whitelist
func (sm *SecurityManager) AddIPToWhitelist(ip string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if sm.WhitelistedIPs[ip] {
        return fmt.Errorf("IP %s is already whitelisted", ip)
    }

    sm.WhitelistedIPs[ip] = true
    fmt.Printf("IP %s added to the whitelist.\n", ip)

    return nil
}

// RemoveIPFromWhitelist removes an IP address from the whitelist
func (sm *SecurityManager) RemoveIPFromWhitelist(ip string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if !sm.WhitelistedIPs[ip] {
        return fmt.Errorf("IP %s is not in the whitelist", ip)
    }

    delete(sm.WhitelistedIPs, ip)
    fmt.Printf("IP %s removed from the whitelist.\n", ip)

    return nil
}

// AddIPToBlacklist adds an IP address to the blacklist
func (sm *SecurityManager) AddIPToBlacklist(ip string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if sm.BlacklistedIPs[ip] {
        return fmt.Errorf("IP %s is already blacklisted", ip)
    }

    sm.BlacklistedIPs[ip] = true
    fmt.Printf("IP %s added to the blacklist.\n", ip)

    return nil
}

// RemoveIPFromBlacklist removes an IP address from the blacklist
func (sm *SecurityManager) RemoveIPFromBlacklist(ip string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if !sm.BlacklistedIPs[ip] {
        return fmt.Errorf("IP %s is not in the blacklist", ip)
    }

    delete(sm.BlacklistedIPs, ip)
    fmt.Printf("IP %s removed from the blacklist.\n", ip)

    return nil
}

// CheckAccess verifies if a node or user is authorized based on IP and key
func (sm *SecurityManager) CheckAccess(ip string, publicKey string) bool {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    // Check if the IP is blacklisted
    if sm.BlacklistedIPs[ip] {
        fmt.Printf("Access denied for IP %s: blacklisted.\n", ip)
        return false
    }

    // Check if the IP is whitelisted
    if sm.WhitelistedIPs[ip] {
        fmt.Printf("Access granted for whitelisted IP %s.\n", ip)
        return true
    }

    // Check if the public key matches any authorized node
    for _, key := range sm.NodeKeys {
        if key == publicKey {
            fmt.Printf("Access granted for node with public key %s.\n", publicKey)
            return true
        }
    }

    fmt.Printf("Access denied for IP %s with public key %s.\n", ip, publicKey)
    return false
}

// generateSessionID generates a unique session ID for a user
func (sm *SecurityManager) generateSessionID(userID string) string {
    hashInput := fmt.Sprintf("%s%d", userID, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// EncryptData securely encrypts the provided data using the system's encryption mechanism
func (sm *SecurityManager) EncryptData(data string) (string, error) {
    // Initialize the Encryption object
    encryption := &Encryption{}

    // Convert data to []byte
    encryptedBytes, err := encryption.EncryptData("AES", []byte(data), EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt data: %v", err)
    }

    // Convert encrypted bytes to string for returning
    return string(encryptedBytes), nil
}

// DecryptData securely decrypts the provided data
func (sm *SecurityManager) DecryptData(encryptedData string) (string, error) {
    // Initialize the Encryption object
    encryption := &Encryption{}

    // Convert encrypted data from string to []byte
    decryptedBytes, err := encryption.DecryptData([]byte(encryptedData), EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt data: %v", err)
    }

    // Convert decrypted bytes back to string for returning
    return string(decryptedBytes), nil
}

