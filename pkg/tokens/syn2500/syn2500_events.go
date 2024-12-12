package syn2500

import (
	"time"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rsa"
	"crypto/rand"
)

// SYN2500Events handles the registration, storage, and processing of DAO-related events for SYN2500 tokens.
type SYN2500Events struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewSYN2500Events initializes the event management for SYN2500 DAO tokens with RSA keys.
func NewSYN2500Events() (*SYN2500Events, error) {
	// Generate encryption keys
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &SYN2500Events{
		privateKey: privKey,
		publicKey:  &privKey.PublicKey,
	}, nil
}

// RegisterEvent registers a new event for a DAO, such as proposals, votes, or membership changes.
func (events *SYN2500Events) RegisterEvent(event common.DAOEvent, token *common.SYN2500Token) error {
	// Generate a unique event ID
	event.EventID = generateUniqueID()
	event.Timestamp = time.Now()

	// Encrypt the event data for secure storage
	encryptedEvent, err := events.encryptEvent(event)
	if err != nil {
		return err
	}

	// Store the encrypted event in the ledger with sub-block validation using Synnergy Consensus
	err = ledger.StoreDAOEvent(encryptedEvent, synconsensus.SubBlockValidation)
	if err != nil {
		return err
	}

	// Append the event to the DAO token's event log
	token.EventLog = append(token.EventLog, event)

	return nil
}

// GetEvent retrieves a specific event from the DAO event log by event ID.
func (events *SYN2500Events) GetEvent(eventID string, token *common.SYN2500Token) (common.DAOEvent, error) {
	for _, event := range token.EventLog {
		if event.EventID == eventID {
			// Return the found event
			return event, nil
		}
	}
	return common.DAOEvent{}, errors.New("event not found")
}

// ValidateEvent validates the integrity and authenticity of an event using its hash and stored metadata.
func (events *SYN2500Events) ValidateEvent(event common.DAOEvent) (bool, error) {
	eventBytes := serializeEvent(event)
	hash := sha256.Sum256(eventBytes)

	// Verify the event hash using the public key
	err := rsa.VerifyPKCS1v15(events.publicKey, crypto.SHA256, hash[:], event.Signature)
	if err != nil {
		return false, errors.New("event validation failed: invalid signature")
	}

	return true, nil
}

// encryptEvent encrypts the event data before storing it in the ledger
func (events *SYN2500Events) encryptEvent(event common.DAOEvent) ([]byte, error) {
	eventBytes := serializeEvent(event)
	hashed := sha512.Sum512(eventBytes)
	encryptedBytes, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, events.publicKey, hashed[:], nil)
	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

// serializeEvent serializes the DAO event into a byte slice for encryption and storage
func serializeEvent(event common.DAOEvent) []byte {
	eventBytes, _ := json.Marshal(event)
	return eventBytes
}

// generateUniqueID creates a unique ID for DAO events
func generateUniqueID() string {
	timestamp := time.Now().UnixNano()
	hash := sha256.New()
	hash.Write([]byte(string(timestamp)))
	return hex.EncodeToString(hash.Sum(nil))
}

// ListEvents returns a list of all events related to a specific DAO token
func (events *SYN2500Events) ListEvents(token *common.SYN2500Token) ([]common.DAOEvent, error) {
	if len(token.EventLog) == 0 {
		return nil, errors.New("no events found for this DAO token")
	}

	return token.EventLog, nil
}

// SignEvent signs the event data using the private key to ensure authenticity
func (events *SYN2500Events) SignEvent(event *common.DAOEvent) error {
	eventBytes := serializeEvent(*event)
	hashed := sha256.Sum256(eventBytes)

	signature, err := rsa.SignPKCS1v15(rand.Reader, events.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return err
	}

	event.Signature = signature
	return nil
}
