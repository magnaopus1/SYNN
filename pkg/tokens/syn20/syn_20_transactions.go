package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SYN20TransactionManager handles transactions for SYN20 tokens.
type SYN20TransactionManager struct {
	mutex       sync.Mutex                 // For thread-safe operations
	Ledger      *ledger.Ledger             // Reference to the blockchain ledger
	Consensus   *synnergy_consensus.Engine // Synnergy Consensus engine
	Storage     *SYN20Storage              // Storage for balances and allowances
	Encryption  *encryption.Encryption     // Encryption service
}

// NewSYN20TransactionManager initializes a new SYN20 transaction manager.
func NewSYN20TransactionManager(ledgerInstance *ledger.Ledger, consensus *synnergy_consensus.Engine, storage *SYN20Storage, encryptionService *encryption.Encryption) *SYN20TransactionManager {
	return &SYN20TransactionManager{
		Ledger:     ledgerInstance,
		Consensus:  consensus,
		Storage:    storage,
		Encryption: encryptionService,
	}
}

// Transfer handles the transfer of SYN20 tokens between two addresses.
func (tm *SYN20TransactionManager) Transfer(fromAddress, toAddress string, amount float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate transaction through Synnergy Consensus
	valid, err := tm.Consensus.ValidateTransaction(fromAddress, toAddress, amount)
	if !valid || err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Perform the transfer by updating the balances
	if err := tm.Storage.Transfer(fromAddress, toAddress, amount); err != nil {
		return fmt.Errorf("error performing token transfer: %v", err)
	}

	// Log the transaction in the ledger
	transactionID := common.GenerateTransactionID()
	if err := tm.Ledger.RecordTokenTransfer(transactionID, fromAddress, toAddress, amount); err != nil {
		return fmt.Errorf("error recording transaction in ledger: %v", err)
	}

	fmt.Printf("Transaction %s: %f SYN transferred from %s to %s\n", transactionID, amount, fromAddress, toAddress)
	return nil
}

// TransferFrom allows a spender to transfer tokens on behalf of an owner.
func (tm *SYN20TransactionManager) TransferFrom(ownerAddress, spenderAddress, recipientAddress string, amount float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate transaction through Synnergy Consensus
	valid, err := tm.Consensus.ValidateTransaction(spenderAddress, recipientAddress, amount)
	if !valid || err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Transfer tokens using allowance mechanism
	if err := tm.Storage.TransferFrom(ownerAddress, spenderAddress, recipientAddress, amount); err != nil {
		return fmt.Errorf("error performing token transfer: %v", err)
	}

	// Log the transaction in the ledger
	transactionID := common.GenerateTransactionID()
	if err := tm.Ledger.RecordTokenTransfer(transactionID, ownerAddress, recipientAddress, amount); err != nil {
		return fmt.Errorf("error recording transaction in ledger: %v", err)
	}

	fmt.Printf("Transaction %s: %f SYN transferred from %s (via %s) to %s\n", transactionID, amount, ownerAddress, spenderAddress, recipientAddress)
	return nil
}

// Approve allows a spender to spend tokens on behalf of the owner up to a specified amount.
func (tm *SYN20TransactionManager) Approve(ownerAddress, spenderAddress string, amount float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Approve the spender to spend on behalf of the owner
	if err := tm.Storage.Approve(ownerAddress, spenderAddress, amount); err != nil {
		return fmt.Errorf("error approving allowance: %v", err)
	}

	// Log the approval in the ledger
	approvalID := common.GenerateTransactionID()
	if err := tm.Ledger.RecordTokenApproval(approvalID, ownerAddress, spenderAddress, amount); err != nil {
		return fmt.Errorf("error recording approval in ledger: %v", err)
	}

	fmt.Printf("Approval %s: %s approved %s to spend %f SYN\n", approvalID, ownerAddress, spenderAddress, amount)
	return nil
}

// Burn allows a user to destroy a specified amount of tokens from their balance, reducing the total supply.
func (tm *SYN20TransactionManager) Burn(fromAddress string, amount float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate that the burn is acceptable
	if err := tm.Storage.Burn(fromAddress, amount); err != nil {
		return fmt.Errorf("error burning tokens: %v", err)
	}

	// Log the burn operation in the ledger
	burnID := common.GenerateTransactionID()
	if err := tm.Ledger.RecordTokenBurn(burnID, fromAddress, amount); err != nil {
		return fmt.Errorf("error recording burn operation in ledger: %v", err)
	}

	fmt.Printf("Burn %s: %f SYN burned from %s\n", burnID, amount, fromAddress)
	return nil
}

// Mint allows a token owner to mint new tokens to a specified address, increasing the total supply.
func (tm *SYN20TransactionManager) Mint(toAddress string, amount float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate that the mint is acceptable and doesn't exceed the maximum supply
	if err := tm.Storage.Mint(toAddress, amount); err != nil {
		return fmt.Errorf("error minting tokens: %v", err)
	}

	// Log the mint operation in the ledger
	mintID := common.GenerateTransactionID()
	if err := tm.Ledger.RecordTokenMint(mintID, toAddress, amount); err != nil {
		return fmt.Errorf("error recording mint operation in ledger: %v", err)
	}

	fmt.Printf("Mint %s: %f SYN minted to %s\n", mintID, amount, toAddress)
	return nil
}

// GetBalance retrieves the balance of a specific address.
func (tm *SYN20TransactionManager) GetBalance(address string) (float64, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the balance from storage
	balance, err := tm.Storage.GetBalance(address)
	if err != nil {
		return 0, fmt.Errorf("error retrieving balance for address %s: %v", address, err)
	}

	return balance, nil
}

// GetAllowance retrieves the amount a spender is allowed to spend on behalf of an owner.
func (tm *SYN20TransactionManager) GetAllowance(ownerAddress, spenderAddress string) (float64, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the allowance from storage
	allowance, err := tm.Storage.GetAllowance(ownerAddress, spenderAddress)
	if err != nil {
		return 0, fmt.Errorf("error retrieving allowance for spender %s on behalf of owner %s: %v", spenderAddress, ownerAddress, err)
	}

	return allowance, nil
}

// ValidateAndExecuteTransaction validates the transaction using Synnergy Consensus and executes it.
func (tm *SYN20TransactionManager) ValidateAndExecuteTransaction(fromAddress, toAddress string, amount float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate transaction through Synnergy Consensus
	valid, err := tm.Consensus.ValidateTransaction(fromAddress, toAddress, amount)
	if !valid || err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Perform the token transfer
	if err := tm.Storage.Transfer(fromAddress, toAddress, amount); err != nil {
		return fmt.Errorf("error performing token transfer: %v", err)
	}

	// Log the transaction in the ledger
	transactionID := common.GenerateTransactionID()
	if err := tm.Ledger.RecordTokenTransfer(transactionID, fromAddress, toAddress, amount); err != nil {
		return fmt.Errorf("error recording transaction in ledger: %v", err)
	}

	fmt.Printf("Transaction %s: %f SYN successfully transferred from %s to %s\n", transactionID, amount, fromAddress, toAddress)
	return nil
}
