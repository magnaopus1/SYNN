package plasma

import (

	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)


// NewPlasmaClient initializes a new PlasmaClient
func NewPlasmaClient(clientID, walletAddress string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, plasmaChain *common.PlasmaChain, networkManager *common.NetworkManager) *common.PlasmaClient {
	return &common.PlasmaClient{
		ClientID:         clientID,
		WalletAddress:    walletAddress,
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
		PlasmaChain:      plasmaChain,
		NetworkManager:   networkManager,
	}
}

// CreateTransaction creates a new UTXO transaction on the Plasma childchain
func (pc *common.PlasmaClient) CreateTransaction(receiver string, amount float64) (*common.PlasmaUTXO, error) {
	// Create a unique transaction ID
	txID := common.GenerateUniqueID()

	// Create the UTXO transaction
	tx := &common.PlasmaUTXO{
		UTXOID:   txID,
		Owner:    pc.WalletAddress,
		Amount:   amount,
		IsSpent:  false,
	}

	// Log the transaction in the ledger
	err := pc.Ledger.RecordTransaction(txID, pc.WalletAddress, receiver, amount, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log transaction in ledger: %v", err)
	}

	fmt.Printf("Client %s created a transaction with ID %s to send %f tokens to %s\n", pc.ClientID, txID, amount, receiver)
	return tx, nil
}

// SubmitTransaction submits a transaction to a Plasma sub-block for validation
func (pc *common.PlasmaClient) SubmitTransaction(subBlockID string, tx *common.PlasmaUTXO) error {
	// Add the transaction to the specified sub-block in the Plasma childchain
	err := pc.PlasmaChain.AddTransaction(subBlockID, tx)
	if err != nil {
		return fmt.Errorf("failed to submit transaction to sub-block %s: %v", subBlockID, err)
	}

	// Encrypt the transaction for security
	encryptedTx, err := pc.EncryptionService.EncryptData([]byte(tx.UTXOID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}

	// Log the encrypted transaction submission in the ledger
	err = pc.Ledger.RecordTransactionSubmission(subBlockID, encryptedTx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction submission: %v", err)
	}

	fmt.Printf("Transaction %s submitted to sub-block %s\n", tx.UTXOID, subBlockID)
	return nil
}

// ValidateSubBlock validates all transactions in a Plasma sub-block
func (pc *common.PlasmaClient) ValidateSubBlock(subBlockID string) error {
	subBlock, err := pc.PlasmaChain.RetrieveSubBlock(subBlockID)
	if err != nil {
		return fmt.Errorf("failed to retrieve sub-block: %v", err)
	}

	// Validate each transaction in the sub-block
	for _, tx := range subBlock.Transactions {
		err := pc.PlasmaChain.ValidateTransaction(tx)
		if err != nil {
			return fmt.Errorf("validation failed for transaction %s: %v", tx.UTXOID, err)
		}
	}

	// Log the sub-block validation in the ledger
	err = pc.Ledger.RecordSubBlockValidation(subBlockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sub-block validation: %v", err)
	}

	fmt.Printf("Sub-block %s validated\n", subBlockID)
	return nil
}

// RetrieveBlockData retrieves the details of a Plasma block
func (pc *common.PlasmaClient) RetrieveBlockData(blockID string) (*common.PlasmaBlock, error) {
	block, err := pc.PlasmaChain.RetrieveBlock(blockID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block %s: %v", blockID, err)
	}

	fmt.Printf("Retrieved block %s with %d sub-blocks\n", blockID, len(block.SubBlocks))
	return block, nil
}

// SpendUTXO spends a UTXO on the Plasma childchain
func (pc *common.PlasmaClient) SpendUTXO(tx *common.PlasmaUTXO) error {
	err := pc.PlasmaChain.SpendTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to spend UTXO: %v", err)
	}

	fmt.Printf("UTXO %s spent by client %s\n", tx.UTXOID, pc.ClientID)
	return nil
}

// SyncWithNetwork syncs the Plasma client with the current state of the Plasma childchain from network nodes
func (pc *common.PlasmaClient) SyncWithNetwork() error {
	err := pc.NetworkManager.SyncChildChainState()
	if err != nil {
		return fmt.Errorf("failed to sync Plasma childchain state with network: %v", err)
	}

	fmt.Printf("Client %s synced with Plasma childchain network\n", pc.ClientID)
	return nil
}

// MonitorSubBlock monitors the status of a specific Plasma sub-block
func (pc *common.PlasmaClient) MonitorSubBlock(subBlockID string) error {
	subBlock, err := pc.PlasmaChain.RetrieveSubBlock(subBlockID)
	if err != nil {
		return fmt.Errorf("failed to monitor sub-block %s: %v", subBlockID, err)
	}

	fmt.Printf("Monitoring sub-block %s: %d transactions\n", subBlockID, len(subBlock.Transactions))
	return nil
}
