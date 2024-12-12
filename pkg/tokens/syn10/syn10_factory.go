package syn10

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)



func NewSYN10Token(
    tokenName, tokenID string,
    issuer IssuerInfo,
    currencyCode string,
    initialSupply uint64,
    consensus *common.SynnergyConsensus,
    ledgerInstance *ledger.SYN10Ledger,
    encryptionService *common.Encryption,
    complianceService *SYN10ComplianceManager,
    exchangeRateAPI string, // URL to fetch exchange rates from
    centralBank string) (*SYN10Token, error) {

    // Step 1: Validate parameters
    if tokenID == "" || currencyCode == "" || initialSupply == 0 || tokenName == "" {
        return nil, fmt.Errorf("invalid token parameters: ensure all required fields are non-empty")
    }

    // Step 2: Fetch the initial exchange rate dynamically
    initialRate, err := fetchCurrentExchangeRate(exchangeRateAPI, currencyCode)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch initial exchange rate: %v", err)
    }

    // Step 3: Initialize metadata
    metadata := &SYN10Metadata{
        TokenID:           generateTokenID(tokenID, currencyCode),
        CurrencyCode:      currencyCode,
        Issuer:            issuer,
        ExchangeRate:      initialRate,
        CreationDate:      time.Now(),
        TotalSupply:       big.NewInt(int64(initialSupply)),
        CirculatingSupply: big.NewInt(int64(initialSupply)),
        PeggingMechanism:  PeggingInfo{},
        LegalCompliance:   LegalInfo{},
    }

    // Step 4: Encrypt metadata
    metadataString := fmt.Sprintf("%+v", metadata) // Convert metadata to string for encryption
    encryptedData, err := encryptionService.EncryptData(metadataString, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("error encrypting metadata: %v", err)
    }
    metadata.EncryptedMetadata = encryptedData

    // Step 5: Store metadata in the ledger
    if err := ledgerInstance.AddTokenMetadata(metadata); err != nil {
        return nil, fmt.Errorf("error storing token metadata in ledger: %v", err)
    }

    // Step 6: Validate token creation using Synnergy Consensus
    valid, err := consensus.ValidateTokenCreation(metadata)
    if err != nil {
        return nil, fmt.Errorf("error during token creation validation: %v", err)
    }
    if !valid {
        return nil, fmt.Errorf("token creation validation failed")
    }

    // Step 7: Initialize SYN10Token instance
    token := &SYN10Token{
        TokenName:        tokenName,
        Metadata:         metadata,
        Ledger:           ledgerInstance,
        Consensus:        consensus,
        Encryption:       encryptionService,
        Compliance:       complianceService,
        CentralBank:      centralBank,
        ExchangeRate:     initialRate,
        TransactionLimits: make(map[string]uint64),
        Allowances:       make(map[string]uint64),
        SecurityProtocols: make(map[string]string),
    }

    // Step 8: Log token creation event
    creationLog := fmt.Sprintf("Token %s created with ID %s by issuer %s at %v", tokenName, tokenID, issuer.Name, time.Now())
    if err := ledgerInstance.recordComplianceAuditLog(tokenID, creationLog); err != nil {
        return nil, fmt.Errorf("failed to log token creation event: %v", err)
    }

    return token, nil
}


// fetchCurrentExchangeRate fetches the current exchange rate for the given currency code from the provided API.
func fetchCurrentExchangeRate(apiURL string, currencyCode string) (float64, error) {
    response, err := http.Get(fmt.Sprintf("%s?currency=%s", apiURL, currencyCode))
    if err != nil {
        return 0, err
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("failed to fetch exchange rate")
    }

    var result map[string]interface{}
    if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
        return 0, err
    }

    rate, ok := result["rate"].(float64)
    if !ok {
        return 0, errors.New("invalid response format")
    }

    return rate, nil
}

// UpdateExchangeRate fetches the current exchange rate from the API and updates the metadata.
func (token *SYN10Token) UpdateExchangeRate() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Fetch the latest exchange rate using the API
    newRate, err := fetchCurrentExchangeRate(token.ExchangeRateAPI, token.Metadata.CurrencyCode)
    if err != nil {
        return fmt.Errorf("failed to update exchange rate: %v", err)
    }

    token.Metadata.ExchangeRate = newRate
    return token.updateLedger()
}

// Mint allows only the central bank to mint new tokens, updates the ledger, and validates the transaction.
func (token *SYN10Token) Mint(amount uint64, issuer string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Only the central bank can mint tokens
    if issuer != token.CentralBank {
        return fmt.Errorf("only the central bank can mint tokens")
    }

    if amount == 0 {
        return fmt.Errorf("amount must be greater than zero")
    }

    token.Metadata.TotalSupply += amount
    token.Metadata.CirculatingSupply += amount

    if err := token.Ledger.AddTokens(token.Metadata.Issuer.Name, amount); err != nil {
        return err
    }

    // Validate the minting process via consensus
    if valid, err := token.Consensus.ValidateMintTransaction(token.Metadata, amount); !valid || err != nil {
        return fmt.Errorf("error validating mint transaction: %v", err)
    }

    return nil
}

// Burn allows only the central bank to burn tokens, updates the ledger, and validates the transaction.
func (token *SYN10Token) Burn(amount uint64, issuer string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Only the central bank can burn tokens
    if issuer != token.CentralBank {
        return fmt.Errorf("only the central bank can burn tokens")
    }

    if amount == 0 || amount > token.Metadata.CirculatingSupply {
        return fmt.Errorf("invalid burn amount")
    }

    token.Metadata.CirculatingSupply -= amount
    token.Metadata.TotalSupply -= amount

    if err := token.Ledger.RemoveTokens(token.Metadata.Issuer.Name, amount); err != nil {
        return err
    }

    // Validate the burning process via consensus
    if valid, err := token.Consensus.ValidateBurnTransaction(token.Metadata, amount); !valid || err != nil {
        return fmt.Errorf("error validating burn transaction: %v", err)
    }

    return nil
}

// Transfer facilitates the transfer of tokens between accounts, updates the ledger, and validates the transaction.
func (token *SYN10Token) Transfer(from, to string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Run KYC/AML checks before transfer
    if err := token.Compliance.RunKYCAmlCheck(from, to); err != nil {
        return fmt.Errorf("transfer denied due to KYC/AML failure: %v", err)
    }

    if err := token.Consensus.ValidateSender(from); err != nil {
        return err
    }
    if err := token.Consensus.ValidateReceiver(to); err != nil {
        return err
    }
    if err := token.Consensus.ValidateAmount(amount); err != nil {
        return err
    }

    if err := token.Ledger.Transfer(from, to, amount); err != nil {
        return err
    }

    // Validate transfer through consensus
    if valid, err := token.Consensus.ValidateTransfer(from, to, amount); !valid || err != nil {
        return fmt.Errorf("error validating transfer: %v", err)
    }

    return nil
}

// GetBalance retrieves the balance of a given address from the ledger.
func (token *SYN10Token) GetBalance(address string) (uint64, error) {
    return token.Ledger.GetBalance(address)
}

// GetTokenDetails provides the metadata of the token.
func (token *SYN10Token) GetTokenDetails() *SYN10Metadata {
    return token.Metadata
}

// updateLedger encrypts metadata and updates it in the ledger.
func (token *SYN10Token) updateLedger() error {
    encryptedData, err := token.Encryption.EncryptData(fmt.Sprintf("%v", token.Metadata), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("error encrypting metadata: %v", err)
    }
    token.Metadata.EncryptedData = encryptedData

    if err := token.Ledger.UpdateTokenMetadata(token.Metadata); err != nil {
        return fmt.Errorf("error updating token metadata in ledger: %v", err)
    }

    return nil
}

// generateTokenID generates a unique token identifier using SHA-256 hashing.
func generateTokenID(tokenID, currencyCode string) string {
    hash := sha256.New()
    hash.Write([]byte(tokenID + currencyCode))
    return hex.EncodeToString(hash.Sum(nil))
}
