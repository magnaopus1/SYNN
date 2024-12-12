// blockchain_compression_and_encryption.go

package main

import (
    "fmt"
    "synnergy_network/pkg/ledger"
)

// blockchainEnableBlockCompression enables compression for blocks.
func blockchainEnableBlockCompression(ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()

    err := ledgerInstance.BlockchainConsensusCoinLedger.EnableBlockCompression()
    if err != nil {
        return fmt.Errorf("failed to enable block compression: %v", err)
    }
    fmt.Println("Block compression enabled.")
    return nil
}

// blockchainDisableBlockCompression disables compression for blocks.
func blockchainDisableBlockCompression(ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()

    err := ledgerInstance.BlockchainConsensusCoinLedger.DisableBlockCompression()
    if err != nil {
        return fmt.Errorf("failed to disable block compression: %v", err)
    }
    fmt.Println("Block compression disabled.")
    return nil
}

// blockchainFetchCompressionStatus retrieves the current compression status.
func blockchainFetchCompressionStatus(ledgerInstance *ledger.Ledger) (bool, error) {
    mutex.Lock()
    defer mutex.Unlock()

    status, err := ledgerInstance.BlockchainConsensusCoinLedger.FetchCompressionStatus()
    if err != nil {
        return false, fmt.Errorf("failed to fetch compression status: %v", err)
    }
    fmt.Printf("Block compression status: %v\n", status)
    return status, nil
}

// blockchainSetBlockCompressionLevel sets the compression level for blocks.
func blockchainSetBlockCompressionLevel(level int, ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()

    err := ledgerInstance.BlockchainConsensusCoinLedger.SetBlockCompressionLevel(level)
    if err != nil {
        return fmt.Errorf("failed to set block compression level: %v", err)
    }
    fmt.Printf("Block compression level set to %d.\n", level)
    return nil
}

// blockchainGetBlockCompressionLevel retrieves the current block compression level.
func blockchainGetBlockCompressionLevel(ledgerInstance *ledger.Ledger) (int, error) {
    mutex.Lock()
    defer mutex.Unlock()

    level, err := ledgerInstance.BlockchainConsensusCoinLedger.GetBlockCompressionLevel()
    if err != nil {
        return 0, fmt.Errorf("failed to get block compression level: %v", err)
    }
    fmt.Printf("Block compression level is %d.\n", level)
    return level, nil
}

// blockchainEnableEncryption enables encryption for blockchain blocks.
func blockchainEnableEncryption(ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()

    err := ledgerInstance.BlockchainConsensusCoinLedger.EnableEncryption()
    if err != nil {
        return fmt.Errorf("failed to enable encryption: %v", err)
    }
    fmt.Println("Encryption enabled.")
    return nil
}

// blockchainDisableEncryption disables encryption for blockchain blocks.
func blockchainDisableEncryption(ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()

    err := ledgerInstance.BlockchainConsensusCoinLedger.DisableEncryption()
    if err != nil {
        return fmt.Errorf("failed to disable encryption: %v", err)
    }
    fmt.Println("Encryption disabled.")
    return nil
}

// blockchainSetEncryptionKey sets the encryption key for the blockchain.
func blockchainSetEncryptionKey(key string, ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()

    err := ledgerInstance.BlockchainConsensusCoinLedger.SetEncryptionKey(key)
    if err != nil {
        return fmt.Errorf("failed to set encryption key: %v", err)
    }
    fmt.Println("Encryption key set.")
    return nil
}

// blockchainGetEncryptionKey retrieves the current encryption key.
func blockchainGetEncryptionKey(ledgerInstance *ledger.Ledger) (string, error) {
    mutex.Lock()
    defer mutex.Unlock()

    key, err := ledgerInstance.BlockchainConsensusCoinLedger.GetEncryptionKey()
    if err != nil {
        return "", fmt.Errorf("failed to get encryption key: %v", err)
    }
    fmt.Println("Encryption key retrieved.")
    return key, nil
}

// blockchainVerifyEncryptionStatus checks if encryption is enabled.
func blockchainVerifyEncryptionStatus(ledgerInstance *ledger.Ledger) (bool, error) {
    mutex.Lock()
    defer mutex.Unlock()

    status, err := ledgerInstance.BlockchainConsensusCoinLedger.VerifyEncryptionStatus()
    if err != nil {
        return false, fmt.Errorf("failed to verify encryption status: %v", err)
    }
    fmt.Printf("Encryption status: %v\n", status)
    return status, nil
}

// blockchainEncryptBlock encrypts a block using AES encryption.
func blockchainEncryptBlock(blockData []byte, ledgerInstance *ledger.Ledger) ([]byte, error) {
    encryptedData, err := ledgerInstance.BlockchainConsensusCoinLedger.EncryptBlock(blockData)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt block: %v", err)
    }
    fmt.Println("Block encrypted.")
    return encryptedData, nil
}

// blockchainDecryptBlock decrypts an encrypted block.
func blockchainDecryptBlock(encryptedData []byte, ledgerInstance *ledger.Ledger) ([]byte, error) {
    decryptedData, err := ledgerInstance.BlockchainConsensusCoinLedger.DecryptBlock(encryptedData)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt block: %v", err)
    }
    fmt.Println("Block decrypted.")
    return decryptedData, nil
}

// blockchainSetSubblockValidationCriteria sets the validation criteria for sub-blocks.
func blockchainSetSubblockValidationCriteria(criteria string, ledgerInstance *ledger.Ledger) error {
    err := ledgerInstance.BlockchainConsensusCoinLedger.SetSubblockValidationCriteria(criteria)
    if err != nil {
        return fmt.Errorf("failed to set subblock validation criteria: %v", err)
    }
    fmt.Printf("Subblock validation criteria set to %s.\n", criteria)
    return nil
}

// blockchainFetchSubblockValidationCriteria retrieves the validation criteria for sub-blocks.
func blockchainFetchSubblockValidationCriteria(ledgerInstance *ledger.Ledger) (string, error) {
    criteria, err := ledgerInstance.BlockchainConsensusCoinLedger.GetSubblockValidationCriteria()
    if err != nil {
        return "", fmt.Errorf("failed to fetch subblock validation criteria: %v", err)
    }
    fmt.Printf("Subblock validation criteria: %s\n", criteria)
    return criteria, nil
}

// blockchainEnableSubblockCompression enables compression for sub-blocks.
func blockchainEnableSubblockCompression(ledgerInstance *ledger.Ledger) error {
    err := ledgerInstance.BlockchainConsensusCoinLedger.EnableSubblockCompression()
    if err != nil {
        return fmt.Errorf("failed to enable subblock compression: %v", err)
    }
    fmt.Println("Subblock compression enabled.")
    return nil
}

// blockchainDisableSubblockCompression disables compression for sub-blocks.
func blockchainDisableSubblockCompression(ledgerInstance *ledger.Ledger) error {
    err := ledgerInstance.BlockchainConsensusCoinLedger.DisableSubblockCompression()
    if err != nil {
        return fmt.Errorf("failed to disable subblock compression: %v", err)
    }
    fmt.Println("Subblock compression disabled.")
    return nil
}

// blockchainSetValidationInterval sets the validation interval for blocks and sub-blocks.
func blockchainSetValidationInterval(interval int, ledgerInstance *ledger.Ledger) error {
    err := ledgerInstance.BlockchainConsensusCoinLedger.SetValidationInterval(interval)
    if err != nil {
        return fmt.Errorf("failed to set validation interval: %v", err)
    }
    fmt.Printf("Validation interval set to %d.\n", interval)
    return nil
}

// blockchainGetValidationInterval retrieves the current validation interval.
func blockchainGetValidationInterval(ledgerInstance *ledger.Ledger) (int, error) {
    interval, err := ledgerInstance.BlockchainConsensusCoinLedger.GetValidationInterval()
    if err != nil {
        return 0, fmt.Errorf("failed to get validation interval: %v", err)
    }
    fmt.Printf("Validation interval is %d.\n", interval)
    return interval, nil
}

// blockchainSetTransactionLimit sets the transaction limit for blocks.
func blockchainSetTransactionLimit(limit int, ledgerInstance *ledger.Ledger) error {
    err := ledgerInstance.BlockchainConsensusCoinLedger.SetTransactionLimit(limit)
    if err != nil {
        return fmt.Errorf("failed to set transaction limit: %v", err)
    }
    fmt.Printf("Transaction limit set to %d.\n", limit)
    return nil
}

// blockchainGetTransactionLimit retrieves the current transaction limit.
func blockchainGetTransactionLimit(ledgerInstance *ledger.Ledger) (int, error) {
    limit, err := ledgerInstance.BlockchainConsensusCoinLedger.GetTransactionLimit()
    if err != nil {
        return 0, fmt.Errorf("failed to get transaction limit: %v", err)
    }
    fmt.Printf("Transaction limit is %d.\n", limit)
    return limit, nil
}

// blockchainEnableTransactionTracking enables tracking of transactions within the blockchain.
func blockchainEnableTransactionTracking(ledgerInstance *ledger.Ledger) error {
    err := ledgerInstance.BlockchainConsensusCoinLedger.EnableTransactionTracking()
    if err != nil {
        return fmt.Errorf("failed to enable transaction tracking: %v", err)
    }
    fmt.Println("Transaction tracking enabled.")
    return nil
}

// blockchainDisableTransactionTracking disables tracking of transactions within the blockchain.
func blockchainDisableTransactionTracking(ledgerInstance *ledger.Ledger) error {
    err := ledgerInstance.BlockchainConsensusCoinLedger.DisableTransactionTracking()
    if err != nil {
        return fmt.Errorf("failed to disable transaction tracking: %v", err)
    }
    fmt.Println("Transaction tracking disabled.")
    return nil
}