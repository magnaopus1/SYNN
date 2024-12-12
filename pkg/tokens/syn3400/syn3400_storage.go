package syn3400

import (
	"errors"
	"sync"
	"time"

)

// ForexStorage manages secure storage of Forex tokens and related data
type ForexStorage struct {
	forexPairs    map[string]*ForexPair // In-memory storage for Forex pairs
	positions     map[string]*Position  // In-memory storage for speculative positions
	hedging       map[string]*HedgingPosition // In-memory storage for hedging positions
	storageMutex  sync.Mutex            // Mutex to ensure thread-safe access
	ledger        *ledger.Ledger        // Ledger integration for logging events
	encryptor     *encryption.Encryptor // Encryptor for data security
	consensus     *consensus.SynnergyConsensus // Consensus mechanism for validating storage actions
}

// NewForexStorage creates a new instance of ForexStorage
func NewForexStorage(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *ForexStorage {
	return &ForexStorage{
		forexPairs:   make(map[string]*ForexPair),
		positions:    make(map[string]*Position),
		hedging:      make(map[string]*HedgingPosition),
		ledger:       ledger,
		encryptor:    encryptor,
		consensus:    consensus,
	}
}

// StoreForexPair securely stores a Forex pair in the storage and logs the action
func (fs *ForexStorage) StoreForexPair(pair *ForexPair) error {
	fs.storageMutex.Lock()
	defer fs.storageMutex.Unlock()

	// Encrypt the Forex pair before storing
	encryptedPair, err := fs.encryptor.EncryptData(pair)
	if err != nil {
		return err
	}

	fs.forexPairs[pair.PairID] = encryptedPair.(*ForexPair)

	// Log the event in the ledger
	fs.ledger.LogEvent("ForexPairStored", time.Now(), pair.PairID)

	// Validate the storage action using consensus
	return fs.consensus.ValidateSubBlock(pair.PairID)
}

// GetForexPair retrieves a Forex pair from storage securely by decrypting it
func (fs *ForexStorage) GetForexPair(pairID string) (*ForexPair, error) {
	fs.storageMutex.Lock()
	defer fs.storageMutex.Unlock()

	pair, exists := fs.forexPairs[pairID]
	if !exists {
		return nil, errors.New("forex pair not found")
	}

	// Decrypt the Forex pair before returning
	decryptedPair, err := fs.encryptor.DecryptData(pair)
	if err != nil {
		return nil, err
	}

	return decryptedPair.(*ForexPair), nil
}

// StorePosition securely stores a speculative Forex position and logs the action
func (fs *ForexStorage) StorePosition(position *Position) error {
	fs.storageMutex.Lock()
	defer fs.storageMutex.Unlock()

	// Encrypt the position data before storing
	encryptedPosition, err := fs.encryptor.EncryptData(position)
	if err != nil {
		return err
	}

	fs.positions[position.PositionID] = encryptedPosition.(*Position)

	// Log the event in the ledger
	fs.ledger.LogEvent("PositionStored", time.Now(), position.PositionID)

	// Validate the storage action using consensus
	return fs.consensus.ValidateSubBlock(position.PositionID)
}

// GetPosition retrieves a speculative position from storage securely by decrypting it
func (fs *ForexStorage) GetPosition(positionID string) (*Position, error) {
	fs.storageMutex.Lock()
	defer fs.storageMutex.Unlock()

	position, exists := fs.positions[positionID]
	if !exists {
		return nil, errors.New("position not found")
	}

	// Decrypt the position data before returning
	decryptedPosition, err := fs.encryptor.DecryptData(position)
	if err != nil {
		return nil, err
	}

	return decryptedPosition.(*Position), nil
}

// StoreHedgingPosition securely stores a hedging position in the storage and logs the action
func (fs *ForexStorage) StoreHedgingPosition(position *HedgingPosition) error {
	fs.storageMutex.Lock()
	defer fs.storageMutex.Unlock()

	// Encrypt the hedging position before storing
	encryptedPosition, err := fs.encryptor.EncryptData(position)
	if err != nil {
		return err
	}

	fs.hedging[position.PositionID] = encryptedPosition.(*HedgingPosition)

	// Log the event in the ledger
	fs.ledger.LogEvent("HedgingPositionStored", time.Now(), position.PositionID)

	// Validate the storage action using consensus
	return fs.consensus.ValidateSubBlock(position.PositionID)
}

// GetHedgingPosition retrieves a hedging position from storage securely by decrypting it
func (fs *ForexStorage) GetHedgingPosition(positionID string) (*HedgingPosition, error) {
	fs.storageMutex.Lock()
	defer fs.storageMutex.Unlock()

	position, exists := fs.hedging[positionID]
	if !exists {
		return nil, errors.New("hedging position not found")
	}

	// Decrypt the hedging position before returning
	decryptedPosition, err := fs.encryptor.DecryptData(position)
	if err != nil {
		return nil, err
	}

	return decryptedPosition.(*HedgingPosition), nil
}

// DeleteForexPair deletes a Forex pair from the storage and logs the action
func (fs *ForexStorage) DeleteForexPair(pairID string) error {
	fs.storageMutex.Lock()
	defer fs.storageMutex.Unlock()

	if _, exists := fs.forexPairs[pairID]; !exists {
		return errors.New("forex pair not found")
	}

	delete(fs.forexPairs, pairID)

	// Log the event in the ledger
	fs.ledger.LogEvent("ForexPairDeleted", time.Now(), pairID)

	// Validate the deletion using consensus
	return fs.consensus.ValidateSubBlock(pairID)
}

// DeletePosition deletes a speculative position from the storage and logs the action
func (fs *ForexStorage) DeletePosition(positionID string) error {
	fs.storageMutex.Lock()
	defer fs.storageMutex.Unlock()

	if _, exists := fs.positions[positionID]; !exists {
		return errors.New("position not found")
	}

	delete(fs.positions, positionID)

	// Log the event in the ledger
	fs.ledger.LogEvent("PositionDeleted", time.Now(), positionID)

	// Validate the deletion using consensus
	return fs.consensus.ValidateSubBlock(positionID)
}

// DeleteHedgingPosition deletes a hedging position from the storage and logs the action
func (fs *ForexStorage) DeleteHedgingPosition(positionID string) error {
	fs.storageMutex.Lock()
	defer fs.storageMutex.Unlock()

	if _, exists := fs.hedging[positionID]; !exists {
		return errors.New("hedging position not found")
	}

	delete(fs.hedging, positionID)

	// Log the event in the ledger
	fs.ledger.LogEvent("HedgingPositionDeleted", time.Now(), positionID)

	// Validate the deletion using consensus
	return fs.consensus.ValidateSubBlock(positionID)
}
