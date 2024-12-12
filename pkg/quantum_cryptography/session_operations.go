package quantum_cryptography

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// QRSession represents a quantum-resistant session with a generated session key
type QRSession struct {
    SessionID   string
    SessionKey  []byte
    CreatedAt   time.Time
    ExpiresAt   time.Time
}

var sessionStore = make(map[string]QRSession)
var sessionLock sync.Mutex

// QRSessionKeyGeneration: Generates a quantum-resistant session key with an expiration time
func QRSessionKeyGeneration(sessionID string, duration time.Duration) (*QRSession, error) {
    sessionLock.Lock()
    defer sessionLock.Unlock()

    if _, exists := sessionStore[sessionID]; exists {
        LogSessionOperation("QRSessionKeyGeneration", "Session already exists: "+sessionID)
        return nil, errors.New("session already exists")
    }

    // Generate a 256-bit quantum-resistant session key
    sessionKey := make([]byte, 32) // 256 bits
    _, err := rand.Read(sessionKey)
    if err != nil {
        LogSessionOperation("QRSessionKeyGeneration", "Session key generation failed")
        return nil, errors.New("failed to generate session key")
    }

    session := QRSession{
        SessionID:  sessionID,
        SessionKey: sessionKey,
        CreatedAt:  time.Now(),
        ExpiresAt:  time.Now().Add(duration),
    }
    sessionStore[sessionID] = session

    LogSessionOperation("QRSessionKeyGeneration", fmt.Sprintf("Session key generated for sessionID %s, expires at %s", sessionID, session.ExpiresAt))
    return &session, nil
}

// Helper Functions

// LogSessionOperation: Logs session operations with encryption
func LogSessionOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("SessionOperation", encryptedMessage)
}
