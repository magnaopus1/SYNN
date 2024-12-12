package syn845

import (
	"errors"
	"sync"
	"time"

)

// Data structure definitions

type DebtMetadata struct {
	ID              string
	Owner           string
	OriginalAmount  float64
	InterestRate    float64
	RepaymentPeriod time.Duration
	PenaltyRate     float64
}

type PaymentRecord struct {
	Date      time.Time
	Amount    float64
	Interest  float64
	Principal float64
	Balance   float64
}

type StatusLog struct {
	Status string
	Date   time.Time
}

type CollateralRecord struct {
	AssetID string
	Value   float64
}

type EventLog struct {
	EventType string
	Date      time.Time
	Details   string
}

// Database struct with mutex for thread safety
type Database struct {
	debtMetadata      map[string]string   // map[DebtID]encryptedData
	paymentRecords    map[string][]string // map[DebtID][]encryptedData
	statusLogs        map[string][]string // map[DebtID][]encryptedData
	collateralRecords map[string]string   // map[AssetID]encryptedData
	eventLogs         map[string][]string // map[DebtID][]encryptedData
	mu                sync.RWMutex
	encryptionKey     [32]byte
	Ledger            *ledger.Ledger               // Ledger for recording storage actions
	ConsensusEngine   *consensus.SynnergyConsensus // Consensus for validating storage actions
}

// NewDatabase initializes a new database instance
func NewDatabase(encryptionKey string, ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus) *Database {
	return &Database{
		debtMetadata:      make(map[string]string),
		paymentRecords:    make(map[string][]string),
		statusLogs:        make(map[string][]string),
		collateralRecords: make(map[string]string),
		eventLogs:         make(map[string][]string),
		encryptionKey:     sha256.Sum256([]byte(encryptionKey)),
		Ledger:            ledger,
		ConsensusEngine:   consensusEngine,
	}
}

// Encryption and decryption methods
func (db *Database) encrypt(data []byte) (string, error) {
	block, err := aes.NewCipher(db.encryptionKey[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (db *Database) decrypt(encryptedData string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(db.encryptionKey[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// Methods to store and retrieve debt metadata with ledger and consensus validation
func (db *Database) StoreDebtMetadata(metadata DebtMetadata) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	encryptedData, err := db.encrypt(data)
	if err != nil {
		return err
	}

	// Validate the storage action via Synnergy Consensus
	if err := db.ConsensusEngine.ValidateStorageAction(metadata.ID, string(data)); err != nil {
		return errors.New("storage validation failed via Synnergy Consensus")
	}

	db.debtMetadata[metadata.ID] = encryptedData

	// Record the storage action in the ledger
	if err := db.Ledger.RecordStorageAction(metadata.ID, encryptedData); err != nil {
		return errors.New("failed to record debt metadata in the ledger")
	}

	return nil
}

func (db *Database) RetrieveDebtMetadata(debtID string) (*DebtMetadata, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	encryptedData, exists := db.debtMetadata[debtID]
	if !exists {
		return nil, errors.New("debt metadata not found")
	}

	data, err := db.decrypt(encryptedData)
	if err != nil {
		return nil, err
	}

	var metadata DebtMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// Methods to store and retrieve payment records with consensus and ledger integration
func (db *Database) StorePaymentRecord(debtID string, record PaymentRecord) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	encryptedData, err := db.encrypt(data)
	if err != nil {
		return err
	}

	// Validate the storage action via Synnergy Consensus
	if err := db.ConsensusEngine.ValidateStorageAction(debtID, string(data)); err != nil {
		return errors.New("storage validation failed via Synnergy Consensus")
	}

	db.paymentRecords[debtID] = append(db.paymentRecords[debtID], encryptedData)

	// Record the storage action in the ledger
	if err := db.Ledger.RecordStorageAction(debtID, encryptedData); err != nil {
		return errors.New("failed to record payment record in the ledger")
	}

	return nil
}

func (db *Database) RetrievePaymentRecords(debtID string) ([]PaymentRecord, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	encryptedRecords, exists := db.paymentRecords[debtID]
	if !exists {
		return nil, errors.New("payment records not found")
	}

	var records []PaymentRecord
	for _, encryptedData := range encryptedRecords {
		data, err := db.decrypt(encryptedData)
		if err != nil {
			return nil, err
		}

		var record PaymentRecord
		if err := json.Unmarshal(data, &record); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

// Methods to store and retrieve status logs with consensus and ledger integration
func (db *Database) StoreStatusLog(debtID string, log StatusLog) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(log)
	if err != nil {
		return err
	}

	encryptedData, err := db.encrypt(data)
	if err != nil {
		return err
	}

	// Validate the storage action via Synnergy Consensus
	if err := db.ConsensusEngine.ValidateStorageAction(debtID, string(data)); err != nil {
		return errors.New("storage validation failed via Synnergy Consensus")
	}

	db.statusLogs[debtID] = append(db.statusLogs[debtID], encryptedData)

	// Record the storage action in the ledger
	if err := db.Ledger.RecordStorageAction(debtID, encryptedData); err != nil {
		return errors.New("failed to record status log in the ledger")
	}

	return nil
}

func (db *Database) RetrieveStatusLogs(debtID string) ([]StatusLog, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	encryptedLogs, exists := db.statusLogs[debtID]
	if !exists {
		return nil, errors.New("status logs not found")
	}

	var logs []StatusLog
	for _, encryptedData := range encryptedLogs {
		data, err := db.decrypt(encryptedData)
		if err != nil {
			return nil, err
		}

		var log StatusLog
		if err := json.Unmarshal(data, &log); err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// Methods to store and retrieve collateral records with consensus and ledger integration
func (db *Database) StoreCollateralRecord(record CollateralRecord) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	encryptedData, err := db.encrypt(data)
	if err != nil {
		return err
	}

	// Validate the storage action via Synnergy Consensus
	if err := db.ConsensusEngine.ValidateStorageAction(record.AssetID, string(data)); err != nil {
		return errors.New("storage validation failed via Synnergy Consensus")
	}

	db.collateralRecords[record.AssetID] = encryptedData

	// Record the storage action in the ledger
	if err := db.Ledger.RecordStorageAction(record.AssetID, encryptedData); err != nil {
		return errors.New("failed to record collateral record in the ledger")
	}

	return nil
}

func (db *Database) RetrieveCollateralRecord(assetID string) (*CollateralRecord, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	encryptedData, exists := db.collateralRecords[assetID]
	if !exists {
		return nil, errors.New("collateral record not found")
	}

	data, err := db.decrypt(encryptedData)
	if err != nil {
		return nil, err
	}

	var record CollateralRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

// Methods to store and retrieve event logs with consensus and ledger integration
func (db *Database) StoreEventLog(debtID string, log EventLog) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(log)
	if err != nil {
		return err
	}

	encryptedData, err := db.encrypt(data)
	if err != nil {
		return err
	}

	// Validate the storage action via Synnergy Consensus
	if err := db.ConsensusEngine.ValidateStorageAction(debtID, string(data)); err != nil {
		return errors.New("storage validation failed via Synnergy Consensus")
	}

	db.eventLogs[debtID] = append(db.eventLogs[debtID], encryptedData)

	// Record the storage action in the ledger
	if err := db.Ledger.RecordStorageAction(debtID, encryptedData); err != nil {
		return errors.New("failed to record event log in the ledger")
	}

	return nil
}

func (db *Database) RetrieveEventLogs(debtID string) ([]EventLog, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	encryptedLogs, exists := db.eventLogs[debtID]
	if !exists {
		return nil, errors.New("event logs not found")
	}

	var logs []EventLog
	for _, encryptedData := range encryptedLogs {
		data, err := db.decrypt(encryptedData)
		if err != nil {
			return nil, err
		}

		var log EventLog
		if err := json.Unmarshal(data, &log); err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}
