package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Syn20Token defines the structure of a SYN20 token contract.
type Syn20Token struct {
	mutex        sync.Mutex
	TokenName    string
	TokenSymbol  string
	TotalSupply  *big.Int
	Decimals     uint8
	BalanceSheet map[string]*big.Int             // Wallet address -> token balance
	Allowances   map[string]map[string]*big.Int  // Owner -> Spender -> Allowance
	Ledger       *ledger.Ledger
	Metadata     *SYN20Metadata
	Encryption   *encryption.Encryption
}

// NewSyn20Token initializes a new SYN20 token contract with the metadata.
func NewSyn20Token(tokenName, tokenSymbol string, decimals uint8, initialSupply *big.Int, ledgerInstance *ledger.Ledger, metadata *SYN20Metadata, encryptionService *encryption.Encryption) *Syn20Token {
	token := &Syn20Token{
		TokenName:    tokenName,
		TokenSymbol:  tokenSymbol,
		TotalSupply:  initialSupply,
		Decimals:     decimals,
		BalanceSheet: make(map[string]*big.Int),
		Allowances:   make(map[string]map[string]*big.Int),
		Ledger:       ledgerInstance,
		Metadata:     metadata,
		Encryption:   encryptionService,
	}

	// Allocate initial supply to the contract creator
	token.BalanceSheet[common.ContractOwnerAddress] = initialSupply

	// Store metadata
	token.Metadata.TokenName = tokenName
	token.Metadata.TokenSymbol = tokenSymbol
	token.Metadata.TotalSupply = initialSupply
	token.Metadata.Decimals = decimals

	// Log the token creation in the ledger
	token.Ledger.RecordTokenCreation(tokenName, tokenSymbol, initialSupply)

	return token
}

// TotalSupply returns the total supply of the token.
func (t *Syn20Token) TotalSupply() *big.Int {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.TotalSupply
}

// BalanceOf returns the token balance of a specific address.
func (t *Syn20Token) BalanceOf(address string) (*big.Int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if balance, exists := t.BalanceSheet[address]; exists {
		return balance, nil
	}

	return nil, errors.New("address not found")
}

// Transfer transfers tokens from one address to another.
func (t *Syn20Token) Transfer(sender, recipient string, amount *big.Int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if amount.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("amount must be greater than zero")
	}

	// Check if sender has enough tokens
	if balance, exists := t.BalanceSheet[sender]; !exists || balance.Cmp(amount) < 0 {
		return errors.New("insufficient balance")
	}

	// Perform the transfer
	t.BalanceSheet[sender].Sub(t.BalanceSheet[sender], amount)
	if _, exists := t.BalanceSheet[recipient]; !exists {
		t.BalanceSheet[recipient] = big.NewInt(0)
	}
	t.BalanceSheet[recipient].Add(t.BalanceSheet[recipient], amount)

	// Log the transfer in the ledger
	err := t.Ledger.RecordTransfer(sender, recipient, amount)
	if err != nil {
		return fmt.Errorf("failed to record transfer: %v", err)
	}

	fmt.Printf("Transferred %s tokens from %s to %s.\n", amount.String(), sender, recipient)
	return nil
}

// Approve allows a spender to spend tokens on behalf of the owner.
func (t *Syn20Token) Approve(owner, spender string, amount *big.Int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if amount.Cmp(big.NewInt(0)) < 0 {
		return errors.New("amount must be non-negative")
	}

	// Initialize the allowance map for the owner if not already present
	if _, exists := t.Allowances[owner]; !exists {
		t.Allowances[owner] = make(map[string]*big.Int)
	}

	// Set the allowance
	t.Allowances[owner][spender] = amount

	// Log the approval in the ledger
	err := t.Ledger.RecordApproval(owner, spender, amount)
	if err != nil {
		return fmt.Errorf("failed to record approval: %v", err)
	}

	fmt.Printf("Approved %s to spend %s tokens on behalf of %s.\n", spender, amount.String(), owner)
	return nil
}

// TransferFrom allows a spender to transfer tokens on behalf of an owner.
func (t *Syn20Token) TransferFrom(spender, owner, recipient string, amount *big.Int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if amount.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("amount must be greater than zero")
	}

	// Check the allowance
	if allowance, exists := t.Allowances[owner][spender]; !exists || allowance.Cmp(amount) < 0 {
		return errors.New("allowance exceeded")
	}

	// Check if owner has enough tokens
	if balance, exists := t.BalanceSheet[owner]; !exists || balance.Cmp(amount) < 0 {
		return errors.New("insufficient balance")
	}

	// Perform the transfer
	t.BalanceSheet[owner].Sub(t.BalanceSheet[owner], amount)
	if _, exists := t.BalanceSheet[recipient]; !exists {
		t.BalanceSheet[recipient] = big.NewInt(0)
	}
	t.BalanceSheet[recipient].Add(t.BalanceSheet[recipient], amount)

	// Subtract from the spender's allowance
	t.Allowances[owner][spender].Sub(t.Allowances[owner][spender], amount)

	// Log the transfer and the change in allowance in the ledger
	err := t.Ledger.RecordTransfer(owner, recipient, amount)
	if err != nil {
		return fmt.Errorf("failed to record transfer: %v", err)
	}
	t.Ledger.RecordAllowanceChange(owner, spender, t.Allowances[owner][spender])

	fmt.Printf("Transferred %s tokens from %s to %s on behalf of %s.\n", amount.String(), owner, recipient, spender)
	return nil
}

// Burn burns a specific amount of tokens from an address.
func (t *Syn20Token) Burn(owner string, amount *big.Int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if amount.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("amount must be greater than zero")
	}

	// Check if owner has enough tokens
	if balance, exists := t.BalanceSheet[owner]; !exists || balance.Cmp(amount) < 0 {
		return errors.New("insufficient balance")
	}

	// Burn the tokens
	t.BalanceSheet[owner].Sub(t.BalanceSheet[owner], amount)
	t.TotalSupply.Sub(t.TotalSupply, amount)

	// Log the burn in the ledger
	err := t.Ledger.RecordBurn(owner, amount)
	if err != nil {
		return fmt.Errorf("failed to record burn: %v", err)
	}

	fmt.Printf("Burned %s tokens from %s.\n", amount.String(), owner)
	return nil
}

// Mint mints new tokens to a specified address.
func (t *Syn20Token) Mint(to string, amount *big.Int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if amount.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("amount must be greater than zero")
	}

	// Mint the tokens
	if _, exists := t.BalanceSheet[to]; !exists {
		t.BalanceSheet[to] = big.NewInt(0)
	}
	t.BalanceSheet[to].Add(t.BalanceSheet[to], amount)
	t.TotalSupply.Add(t.TotalSupply, amount)

	// Log the minting in the ledger
	err := t.Ledger.RecordMint(to, amount)
	if err != nil {
		return fmt.Errorf("failed to record minting: %v", err)
	}

	fmt.Printf("Minted %s tokens to %s.\n", amount.String(), to)
	return nil
}

// Allowance returns the current allowance for a spender on behalf of an owner.
func (t *Syn20Token) Allowance(owner, spender string) (*big.Int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if allowance, exists := t.Allowances[owner][spender]; exists {
		return allowance, nil
	}

	return nil, errors.New("no allowance set")
}
