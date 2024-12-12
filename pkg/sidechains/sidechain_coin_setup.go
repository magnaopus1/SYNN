package sidechains

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	
)


// NewSidechainCoinSetup initializes a new sidechain coin setup
func NewSidechainCoinSetup(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus) *common.SidechainCoinSetup {
	return &common.SidechainCoinSetup{
		Coins:        make(map[string]*common.SidechainCoin),
		Balances:     make(map[string]map[string]float64),
		Transactions: make([]*common.Transaction, 0),
		Ledger:       ledgerInstance,
		Encryption:   encryptionService,
		Consensus:    common.SynnergyConsensus,
	}
}

// CreateCoin creates a new coin on the sidechain
func (sc *common.SidechainCoinSetup) CreateCoin(coinID, name, symbol string, totalSupply float64, decimals int) (*common.SidechainCoin, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if _, exists := sc.Coins[coinID]; exists {
		return nil, errors.New("coin with this ID already exists")
	}

	// Create the coin
	coin := &common.SidechainCoin{
		CoinID:      coinID,
		Name:        name,
		Symbol:      symbol,
		TotalSupply: totalSupply,
		Decimals:    decimals,
	}

	// Set initial supply to creator
	sc.Balances[coinID] = make(map[string]float64)
	sc.Balances[coinID]["creator"] = totalSupply

	// Record the coin creation in the ledger
	err := sc.Ledger.RecordCoinCreation(coinID, name, symbol, totalSupply, decimals, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log coin creation: %v", err)
	}

	sc.Coins[coinID] = coin
	fmt.Printf("Coin %s (%s) created with total supply of %f\n", name, symbol, totalSupply)
	return coin, nil
}

// TransferCoin transfers a specific amount of coins from one user to another
func (sc *common.SidechainCoinSetup) TransferCoin(coinID, senderID, receiverID string, amount float64) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	coin, exists := sc.Coins[coinID]
	if !exists {
		return errors.New("coin does not exist")
	}

	// Check sender's balance
	if sc.Balances[coinID][senderID] < amount {
		return errors.New("insufficient balance")
	}

	// Update balances
	sc.Balances[coinID][senderID] -= amount
	sc.Balances[coinID][receiverID] += amount

	// Create a transaction
	tx := &common.Transaction{
		TxID:       common.GenerateTransactionID(),
		From:       senderID,
		To:         receiverID,
		Amount:     amount,
		CoinID:     coinID,
		Timestamp:  time.Now(),
	}

	// Encrypt transaction details
	encryptedTx, err := sc.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction data: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add transaction to the list
	sc.Transactions = append(sc.Transactions, tx)

	// Log the transaction in the ledger
	err = sc.Ledger.RecordTransaction(tx.TxID, senderID, receiverID, coinID, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction: %v", err)
	}

	fmt.Printf("Transferred %f %s from %s to %s\n", amount, coin.Symbol, senderID, receiverID)
	return nil
}

// GetBalance retrieves the balance of a specific user for a specific coin
func (sc *common.SidechainCoinSetup) GetBalance(coinID, userID string) (float64, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	_, exists := sc.Coins[coinID]
	if !exists {
		return 0, errors.New("coin does not exist")
	}

	balance, exists := sc.Balances[coinID][userID]
	if !exists {
		return 0, errors.New("user does not have any balance for this coin")
	}

	fmt.Printf("Balance of user %s for coin %s: %f\n", userID, coinID, balance)
	return balance, nil
}

// MintCoins allows the coin issuer to mint additional coins
func (sc *common.SidechainCoinSetup) MintCoins(coinID, issuerID string, amount float64) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	coin, exists := sc.Coins[coinID]
	if !exists {
		return errors.New("coin does not exist")
	}

	// Add new coins to the issuer's balance
	sc.Balances[coinID][issuerID] += amount
	coin.TotalSupply += amount

	// Record the minting event in the ledger
	err := sc.Ledger.RecordCoinMinting(coinID, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log coin minting: %v", err)
	}

	fmt.Printf("Minted %f new coins for %s\n", amount, issuerID)
	return nil
}

// BurnCoins allows the coin issuer to burn (destroy) coins, reducing total supply
func (sc *common.SidechainCoinSetup) BurnCoins(coinID, burnerID string, amount float64) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	coin, exists := sc.Coins[coinID]
	if !exists {
		return errors.New("coin does not exist")
	}

	// Check burner's balance
	if sc.Balances[coinID][burnerID] < amount {
		return errors.New("insufficient balance to burn")
	}

	// Reduce balance and total supply
	sc.Balances[coinID][burnerID] -= amount
	coin.TotalSupply -= amount

	// Record the burn event in the ledger
	err := sc.Ledger.RecordCoinBurning(coinID, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log coin burning: %v", err)
	}

	fmt.Printf("Burned %f coins from %s\n", amount, burnerID)
	return nil
}

// ValidateTransaction uses Synnergy Consensus to validate transactions
func (sc *common.SidechainCoinSetup) ValidateTransaction(txID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Find the transaction
	for _, tx := range sc.Transactions {
		if tx.TxID == txID {
			// Use consensus to validate
			err := sc.Consensus.consensus.ValidateTransaction(tx)
			if err != nil {
				return fmt.Errorf("transaction validation failed: %v", err)
			}

			fmt.Printf("Transaction %s validated\n", txID)
			return nil
		}
	}

	return fmt.Errorf("transaction %s not found", txID)
}

// GetTotalSupply retrieves the total supply of a sidechain coin
func (sc *common.SidechainCoinSetup) GetTotalSupply(coinID string) (float64, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	coin, exists := sc.Coins[coinID]
	if !exists {
		return 0, errors.New("coin does not exist")
	}

	fmt.Printf("Total supply of coin %s: %f\n", coinID, coin.TotalSupply)
	return coin.TotalSupply, nil
}

// RetrieveTransaction retrieves a transaction by its ID
func (sc *common.SidechainCoinSetup) RetrieveTransaction(txID string) (*common.Transaction, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for _, tx := range sc.Transactions {
		if tx.TxID == txID {
			fmt.Printf("Transaction %s retrieved\n", txID)
			return tx, nil
		}
	}

	return nil, fmt.Errorf("transaction %s not found", txID)
}
