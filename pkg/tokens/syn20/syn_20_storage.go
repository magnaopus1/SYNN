package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SYN20Storage handles storage for SYN20 token balances and allowances.
type SYN20Storage struct {
	mutex       sync.Mutex                 // For thread-safe operations
	Balances    map[string]float64         // Mapping of address to balance
	Allowances  map[string]map[string]float64 // Mapping of owner to spender and the amount allowed
	Ledger      *ledger.Ledger             // Reference to the ledger
	Consensus   *synnergy_consensus.Engine // Synnergy Consensus engine
	Encryption  *encryption.Encryption     // Encryption service for sensitive data
}

// NewSYN20Storage initializes a new SYN20 storage manager.
func NewSYN20Storage(ledgerInstance *ledger.Ledger, consensus *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN20Storage {
	return &SYN20Storage{
		Balances:   make(map[string]float64),
		Allowances: make(map[string]map[string]float64),
		Ledger:     ledgerInstance,
		Consensus:  consensus,
		Encryption: encryptionService,
	}
}

// GetBalance retrieves the balance for a specific address.
func (ss *SYN20Storage) GetBalance(address string) (float64, error) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Get balance from the ledger
	balance, err := ss.Ledger.GetBalance(address)
	if err != nil {
		return 0, fmt.Errorf("error retrieving balance for address %s: %v", address, err)
	}

	// Decrypt the balance
	decryptedBalance, err := ss.Encryption.DecryptData(balance, common.EncryptionKey)
	if err != nil {
		return 0, fmt.Errorf("error decrypting balance for address %s: %v", address, err)
	}

	return decryptedBalance, nil
}

// SetBalance sets the balance for a specific address, storing it securely in the ledger.
func (ss *SYN20Storage) SetBalance(address string, amount float64) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Encrypt the balance
	encryptedBalance, err := ss.Encryption.EncryptData(amount, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting balance for address %s: %v", address, err)
	}

	// Store encrypted balance in the ledger
	if err := ss.Ledger.UpdateBalance(address, encryptedBalance); err != nil {
		return fmt.Errorf("error storing balance for address %s: %v", address, err)
	}

	return nil
}

// Transfer transfers tokens from one address to another and updates the ledger accordingly.
func (ss *SYN20Storage) Transfer(fromAddress, toAddress string, amount float64) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Get the sender's balance
	fromBalance, err := ss.GetBalance(fromAddress)
	if err != nil {
		return err
	}

	// Check if sender has enough balance
	if fromBalance < amount {
		return errors.New("insufficient balance")
	}

	// Get the recipient's balance
	toBalance, err := ss.GetBalance(toAddress)
	if err != nil {
		return err
	}

	// Perform the transfer by deducting and adding balances
	newFromBalance := fromBalance - amount
	newToBalance := toBalance + amount

	// Encrypt and store the updated balances
	if err := ss.SetBalance(fromAddress, newFromBalance); err != nil {
		return err
	}
	if err := ss.SetBalance(toAddress, newToBalance); err != nil {
		return err
	}

	// Log transfer in the ledger
	if err := ss.Ledger.RecordTransfer(fromAddress, toAddress, amount); err != nil {
		return fmt.Errorf("error logging transfer in the ledger: %v", err)
	}

	return nil
}

// Approve sets an allowance for a spender on behalf of the token owner.
func (ss *SYN20Storage) Approve(owner, spender string, amount float64) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Encrypt the allowance
	encryptedAllowance, err := ss.Encryption.EncryptData(amount, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting allowance: %v", err)
	}

	// Store the allowance
	if ss.Allowances[owner] == nil {
		ss.Allowances[owner] = make(map[string]float64)
	}
	ss.Allowances[owner][spender] = encryptedAllowance

	// Log the approval in the ledger
	if err := ss.Ledger.RecordAllowance(owner, spender, encryptedAllowance); err != nil {
		return fmt.Errorf("error logging allowance in ledger: %v", err)
	}

	return nil
}

// GetAllowance retrieves the allowance set by the owner for a spender.
func (ss *SYN20Storage) GetAllowance(owner, spender string) (float64, error) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Get the encrypted allowance from the ledger
	allowance, err := ss.Ledger.GetAllowance(owner, spender)
	if err != nil {
		return 0, fmt.Errorf("error retrieving allowance for %s by %s: %v", spender, owner, err)
	}

	// Decrypt the allowance
	decryptedAllowance, err := ss.Encryption.DecryptData(allowance, common.EncryptionKey)
	if err != nil {
		return 0, fmt.Errorf("error decrypting allowance for %s by %s: %v", spender, owner, err)
	}

	return decryptedAllowance, nil
}

// TransferFrom transfers tokens from an owner's address to a recipient on behalf of a spender.
func (ss *SYN20Storage) TransferFrom(owner, spender, recipient string, amount float64) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Get the current allowance for the spender
	allowance, err := ss.GetAllowance(owner, spender)
	if err != nil {
		return err
	}

	// Check if allowance is sufficient
	if allowance < amount {
		return errors.New("allowance exceeded")
	}

	// Perform the transfer
	if err := ss.Transfer(owner, recipient, amount); err != nil {
		return err
	}

	// Decrease the allowance by the amount transferred
	newAllowance := allowance - amount
	if err := ss.Approve(owner, spender, newAllowance); err != nil {
		return err
	}

	// Log the transfer in the ledger
	if err := ss.Ledger.RecordTransfer(owner, recipient, amount); err != nil {
		return fmt.Errorf("error logging transfer in ledger: %v", err)
	}

	return nil
}

// Burn burns tokens from an address, reducing the total supply.
func (ss *SYN20Storage) Burn(address string, amount float64) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Get the balance of the address
	balance, err := ss.GetBalance(address)
	if err != nil {
		return err
	}

	// Check if there are enough tokens to burn
	if balance < amount {
		return errors.New("insufficient balance to burn")
	}

	// Reduce the balance and update the ledger
	newBalance := balance - amount
	if err := ss.SetBalance(address, newBalance); err != nil {
		return err
	}

	// Log the burn operation in the ledger
	if err := ss.Ledger.RecordBurn(address, amount); err != nil {
		return fmt.Errorf("error logging burn operation in ledger: %v", err)
	}

	return nil
}

// Mint mints new tokens to a specified address, increasing the total supply.
func (ss *SYN20Storage) Mint(address string, amount float64) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Get the current balance of the address
	balance, err := ss.GetBalance(address)
	if err != nil {
		return err
	}

	// Increase the balance and update the ledger
	newBalance := balance + amount
	if err := ss.SetBalance(address, newBalance); err != nil {
		return err
	}

	// Log the mint operation in the ledger
	if err := ss.Ledger.RecordMint(address, amount); err != nil {
		return fmt.Errorf("error logging mint operation in ledger: %v", err)
	}

	return nil
}
