package integrated_charity_management

import (
	"encoding/base64"
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewInternalCharityPool initializes the internal charity pool
func NewInternalCharityPool(ownerAddress string, ledgerInstance *ledger.Ledger) *InternalCharityPool {
    icp := &InternalCharityPool{
        PoolBalance:    0, // Initial balance is 0
        WalletAddresses: make(map[string]float64),
        OwnerAddress:   ownerAddress,
        LedgerInstance: ledgerInstance,
        stopChan:       make(chan bool),
    }

    // Start the 24-hour distribution cycle
    go icp.start24HrDistribution()

    return icp
}

// AddWalletAddress allows the blockchain owner to add a wallet address for internal charities
func (icp *InternalCharityPool) AddWalletAddress(walletAddress string) error {
    icp.mutex.Lock()
    defer icp.mutex.Unlock()

    if walletAddress == "" {
        return errors.New("wallet address cannot be empty")
    }

    // Check if the wallet address already exists
    if _, exists := icp.WalletAddresses[walletAddress]; exists {
        return fmt.Errorf("wallet address %s already exists", walletAddress)
    }

    initialBalance := 0.0 // Initialize with 0 balance
    icp.WalletAddresses[walletAddress] = initialBalance
    fmt.Printf("Wallet address %s added to the internal charity pool.\n", walletAddress)

    // Log the addition to the ledger with wallet address and initial balance
    err := icp.logWalletAdditionToLedger(walletAddress, initialBalance)
    if err != nil {
        return fmt.Errorf("failed to log wallet address addition to ledger: %v", err)
    }

    return nil
}

// UpdatePoolBalance updates the internal charity pool balance by receiving funds from the charity pool
func (icp *InternalCharityPool) UpdatePoolBalance(amount float64) error {
    icp.mutex.Lock()
    defer icp.mutex.Unlock()

    if amount <= 0 {
        return errors.New("invalid amount: must be greater than zero")
    }

    icp.PoolBalance += amount
    fmt.Printf("Internal charity pool balance updated by %.2f SYNN. Current balance: %.2f SYNN.\n", amount, icp.PoolBalance)

    // Log the pool update to the ledger
    return icp.logPoolBalanceUpdateToLedger(amount)
}

// start24HrDistribution handles fund distribution every 24 hours
func (icp *InternalCharityPool) start24HrDistribution() {
    for {
        select {
        case <-icp.stopChan:
            return
        case <-time.After(24 * time.Hour): // Wait for 24 hours
            fmt.Println("Starting 24-hour distribution cycle for the internal charity pool.")
            icp.DistributeFunds()
        }
    }
}

// DistributeFunds distributes available pool funds equally among registered wallet addresses every 24 hours
func (icp *InternalCharityPool) DistributeFunds() error {
    icp.mutex.Lock()
    defer icp.mutex.Unlock()

    if icp.PoolBalance <= 0 {
        return errors.New("insufficient funds in the internal charity pool")
    }

    numWallets := len(icp.WalletAddresses)
    if numWallets == 0 {
        return errors.New("no wallet addresses available for distribution")
    }

    // Calculate the amount each wallet will receive
    distributionAmount := icp.PoolBalance / float64(numWallets)

    // Distribute funds and reset the pool balance
    for wallet := range icp.WalletAddresses {
        icp.WalletAddresses[wallet] += distributionAmount
        fmt.Printf("Distributed %.2f SYNN to wallet %s.\n", distributionAmount, wallet)

        // Log the distribution to the ledger
        err := icp.logFundDistributionToLedger(wallet, distributionAmount)
        if err != nil {
            return fmt.Errorf("failed to log fund distribution to ledger for wallet %s: %v", wallet, err)
        }
    }

    // Reset the pool balance after distribution
    icp.PoolBalance = 0
    fmt.Println("Internal charity pool balance reset to 0 after fund distribution.")

    return nil
}

// Stop24HrDistribution stops the 24-hour distribution cycle
func (icp *InternalCharityPool) Stop24HrDistribution() {
    icp.stopChan <- true
}

// GetWalletBalance retrieves the balance of a specific wallet in the internal charity pool
func (icp *InternalCharityPool) GetWalletBalance(walletAddress string) (float64, error) {
    icp.mutex.Lock()
    defer icp.mutex.Unlock()

    balance, exists := icp.WalletAddresses[walletAddress]
    if !exists {
        return 0, fmt.Errorf("wallet address %s not found", walletAddress)
    }

    return balance, nil
}

// logWalletAdditionToLedger logs the addition of a wallet address to the ledger
func (icp *InternalCharityPool) logWalletAdditionToLedger(walletAddress string, balance float64) error {
    logData := fmt.Sprintf("Added wallet address: %s to internal charity pool with balance %.2f.", walletAddress, balance)

    // Create an encryption instance
    encryptionInstance, err := common.NewEncryption(256) // Assuming AES 256-bit key
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt the log data
    encryptedLog, err := encryptionInstance.EncryptData("AES", []byte(logData), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt wallet address log: %v", err)
    }

    // Convert encrypted log to string if needed (optional logging)
    encryptedLogString := base64.StdEncoding.EncodeToString(encryptedLog)
    fmt.Printf("Encrypted Wallet Addition Log: %s\n", encryptedLogString)

    // Only pass the balance to the ledger
    return icp.LedgerInstance.RecordInternalCharityWalletAddition(balance)
}


// logPoolBalanceUpdateToLedger logs the update of the pool balance to the ledger
func (icp *InternalCharityPool) logPoolBalanceUpdateToLedger(amount float64) error {
    logData := fmt.Sprintf("Internal charity pool updated by %.2f SYNN.", amount)

    // Create an encryption instance
    encryptionInstance, err := common.NewEncryption(256) // Assuming AES 256-bit key
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt the log data
    encryptedLog, err := encryptionInstance.EncryptData("AES", []byte(logData), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt pool balance update log: %v", err)
    }

    // Convert encrypted log to string if needed (optional logging)
    encryptedLogString := base64.StdEncoding.EncodeToString(encryptedLog)
    fmt.Printf("Encrypted Pool Balance Update Log: %s\n", encryptedLogString)

    // Only pass the amount to the ledger
    return icp.LedgerInstance.RecordInternalCharityPoolUpdate(amount)
}



// logFundDistributionToLedger logs the distribution of funds to a wallet in the ledger
func (icp *InternalCharityPool) logFundDistributionToLedger(walletAddress string, amount float64) error {
    logData := fmt.Sprintf("Distributed %.2f SYNN to wallet address: %s from internal charity pool.", amount, walletAddress)

    // Create an encryption instance
    encryptionInstance, err := common.NewEncryption(256) // Assuming AES 256-bit key
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt the log data
    encryptedLog, err := encryptionInstance.EncryptData("AES", []byte(logData), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt fund distribution log: %v", err)
    }

    // Convert encrypted log to string if needed (optional logging)
    encryptedLogString := base64.StdEncoding.EncodeToString(encryptedLog)
    fmt.Printf("Encrypted Fund Distribution Log: %s\n", encryptedLogString)

    // Only pass the amount to the ledger
    return icp.LedgerInstance.RecordInternalCharityFundDistribution(amount)
}


