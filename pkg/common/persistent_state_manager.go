package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"synnergy_network/pkg/ledger"
)

// VMState represents the state of the Virtual Machine during execution.
type VMState struct {
	ProgramCounter int             // Tracks the current instruction position
	StackPointer   int             // Points to the top of the stack
	Registers      map[string]int  // Stores register values
	Flags          map[string]bool // Condition flags (e.g., Zero, Carry)
	Memory         []byte          // VM memory for instructions and data
	Accounts       map[string]int64 // Stores account balances or other account-related data
}

// NewVMState initializes a new VMState with default values.
func NewVMState(memorySize int) *VMState {
	return &VMState{
		ProgramCounter: 0,
		StackPointer:   0,
		Registers:      make(map[string]int),
		Flags:          make(map[string]bool),
		Memory:         make([]byte, memorySize),
		Accounts:       make(map[string]int64),
	}
}

// PersistentStateManager manages the VMState using the ledger with enhanced control.
type PersistentStateManager struct {
	state              *VMState
	ledger             *ledger.Ledger
	mutex              sync.RWMutex
	logger             *LogManager
	stateFile          string
	stateRebuildNeeded bool // Flag to check if a rebuild is necessary
	rebuildAttempted   bool // Flag to track if a rebuild was attempted recently
}

var (
	persistentStateManagerInstance *PersistentStateManager
	persistentStateManagerOnce     sync.Once
)

// GetPersistentStateManager returns the singleton instance of PersistentStateManager.
func GetPersistentStateManager(ledgerInstance *ledger.Ledger, stateFile string, logger *LogManager) (*PersistentStateManager, error) {
	var err error
	persistentStateManagerOnce.Do(func() {
		persistentStateManagerInstance, err = newPersistentStateManager(ledgerInstance, stateFile, logger)
	})
	return persistentStateManagerInstance, err
}

// newPersistentStateManager initializes a new PersistentStateManager.
func newPersistentStateManager(ledgerInstance *ledger.Ledger, stateFile string, logger *LogManager) (*PersistentStateManager, error) {
	psm := &PersistentStateManager{
		state:              NewVMState(1024),
		ledger:             ledgerInstance,
		logger:             logger,
		stateFile:          stateFile,
		stateRebuildNeeded: true, // Initial rebuild/load required
		rebuildAttempted:   false,
	}

	// Register ledger block listener
	ledgerInstance.BlockchainConsensusCoinLedger.RegisterBlockListener(func(block ledger.Block) {
		psm.mutex.Lock()
		defer psm.mutex.Unlock()
		psm.stateRebuildNeeded = true
		psm.rebuildAttempted = false
	})

	// Load or rebuild the state based on file existence
	if err := psm.loadOrRebuildState(); err != nil {
		return nil, fmt.Errorf("failed to initialize state: %v", err)
	}

	return psm, nil
}

// loadOrRebuildState tries to load state from file or rebuild if necessary.
func (psm *PersistentStateManager) loadOrRebuildState() error {
	psm.mutex.Lock()
	defer psm.mutex.Unlock()

	if _, err := os.Stat(psm.stateFile); err == nil {
		if err := psm.loadState(); err == nil {
			psm.stateRebuildNeeded = false // Mark rebuild as unnecessary after successful load
			psm.rebuildAttempted = false   // Reset rebuild attempt flag
			return nil
		}
		psm.logger.Error("Failed to load state from file; will attempt rebuild", nil)
	}

	// Trigger a rebuild only if the flag allows it
	if psm.stateRebuildNeeded && !psm.rebuildAttempted {
		psm.rebuildAttempted = true
		return psm.rebuildState()
	}

	psm.logger.Info("State rebuild not needed or already attempted; bypassing unnecessary rebuild", nil)
	return nil
}

// RebuildState rebuilds the state from the ledger, skips if no blocks.
func (psm *PersistentStateManager) rebuildState() error {
	psm.mutex.Lock()
	defer psm.mutex.Unlock()

	// Reset state and check blocks to avoid redundant rebuild attempts
	psm.state = NewVMState(1024) 
	blocks := psm.ledger.BlockchainConsensusCoinLedger.GetBlocks()
	if len(blocks) == 0 {
		psm.logger.Info("No blocks found in ledger; state rebuild skipped", nil)
		psm.stateRebuildNeeded = false // Prevent further attempts
		return nil
	}

	psm.logger.Info("Rebuilding state from ledger", nil)
	for _, block := range blocks {
		for _, subBlock := range block.SubBlocks {
			for _, tx := range subBlock.Transactions {
				if err := psm.applyTransactionNoLock(tx); err != nil {
					psm.logger.Error("Failed to apply transaction during state rebuild", map[string]interface{}{
						"txID":  tx.TransactionID,
						"error": err,
					})
					return err
				}
			}
		}
	}
	psm.stateRebuildNeeded = false
	psm.logger.Info("State successfully rebuilt from ledger", nil)
	return nil
}


// applyTransactionNoLock applies a transaction to the state without acquiring a lock.
func (psm *PersistentStateManager) applyTransactionNoLock(tx ledger.Transaction) error {
	if tx.FromAddress == "" || tx.ToAddress == "" || tx.Amount <= 0 {
		return errors.New("invalid transaction parameters")
	}

	fromBalance, fromExists := psm.state.Accounts[tx.FromAddress]
	if !fromExists {
		psm.state.Accounts[tx.FromAddress] = 0 // Initialize account with zero balance
		fromBalance = 0
	}
	if fromBalance < int64(tx.Amount) {
		psm.logger.Error("Insufficient funds for transaction", map[string]interface{}{
			"sender":   tx.FromAddress,
			"balance":  fromBalance,
			"required": tx.Amount,
		})
		return fmt.Errorf("insufficient funds in account '%s'", tx.FromAddress)
	}

	if _, toExists := psm.state.Accounts[tx.ToAddress]; !toExists {
		psm.state.Accounts[tx.ToAddress] = 0
		psm.logger.Info("Created new recipient account", map[string]interface{}{"recipient": tx.ToAddress})
	}

	psm.state.Accounts[tx.FromAddress] -= int64(tx.Amount)
	psm.state.Accounts[tx.ToAddress] += int64(tx.Amount)

	psm.logger.Info("Transaction applied successfully", map[string]interface{}{
		"from":    tx.FromAddress,
		"to":      tx.ToAddress,
		"amount":  tx.Amount,
		"newFrom": psm.state.Accounts[tx.FromAddress],
		"newTo":   psm.state.Accounts[tx.ToAddress],
	})

	return nil
}

// CommitState serializes and saves the current state to a file.
func (psm *PersistentStateManager) CommitState() error {
	psm.mutex.Lock()
	defer psm.mutex.Unlock()

	err := psm.saveState()
	if err != nil {
		psm.logger.Error("Failed to save state", map[string]interface{}{"error": err})
		return err
	}

	psm.logger.Info("State committed", nil)
	return nil
}

// saveState serializes the state to a file, handling errors.
func (psm *PersistentStateManager) saveState() error {
	data, err := json.Marshal(psm.state)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(psm.stateFile, data, 0644)
}

// loadState loads the VMState from a file with error handling.
func (psm *PersistentStateManager) loadState() error {
	data, err := ioutil.ReadFile(psm.stateFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, psm.state)
}

// Close gracefully closes the PersistentStateManager, committing the state.
func (psm *PersistentStateManager) Close() error {
	if err := psm.CommitState(); err != nil {
		return err
	}
	psm.logger.Info("PersistentStateManager closed", nil)
	return nil
}
