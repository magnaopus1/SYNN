// blockchain_cache_and_capacity_management.go

package main

import (
	"fmt"
	"sync"

	"synnergy_network/pkg/ledger"
)

var mutex sync.Mutex

// blockchainSetTransactionRetryLimit sets the maximum retry limit for transactions.
func blockchainSetTransactionRetryLimit(limit int, ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetTransactionRetryLimit(limit); err != nil {
        return fmt.Errorf("failed to set transaction retry limit: %v", err)
    }
    fmt.Printf("Transaction retry limit set to %d.\n", limit)
    return nil
}

// blockchainGetTransactionRetryLimit retrieves the current retry limit for transactions.
func blockchainGetTransactionRetryLimit(ledgerInstance *ledger.Ledger) (int, error) {
    mutex.Lock()
    defer mutex.Unlock()
    limit, err := ledgerInstance.BlockchainConsensusCoinLedger.GetTransactionRetryLimit()
    if err != nil {
        return 0, fmt.Errorf("failed to get transaction retry limit: %v", err)
    }
    fmt.Printf("Transaction retry limit is %d.\n", limit)
    return limit, nil
}

// blockchainMonitorTransactionRetries monitors and logs retry counts of failed transactions.
func blockchainMonitorTransactionRetries(ledgerInstance *ledger.Ledger) {
    retries := ledgerInstance.BlockchainConsensusCoinLedger.GetTransactionRetryData()
    for txID, count := range retries {
        fmt.Printf("Transaction %s has retried %d times.\n", txID, count)
    }
}

// blockchainAuditTransactionRetries performs an audit on the retry data of transactions.
func blockchainAuditTransactionRetries(ledgerInstance *ledger.Ledger) {
    err := ledgerInstance.BlockchainConsensusCoinLedger.AuditTransactionRetries()
    if err != nil {
        fmt.Println("Failed to audit transaction retries:", err)
    } else {
        fmt.Println("Transaction retry audit completed.")
    }
}


// blockchainTrackTransactionRetries tracks retry counts for each transaction in the ledger.
func blockchainTrackTransactionRetries(txID string, ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()
    err := ledgerInstance.BlockchainConsensusCoinLedger.IncrementRetryCount(txID)
    if err != nil {
        fmt.Printf("Failed to track retries for transaction %s: %v\n", txID, err)
    }
}

// blockchainEnableAuditLogging enables transaction audit logging.
func blockchainEnableAuditLogging(ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()
    ledgerInstance.ComplianceLedger.SetAuditLogging(true)
    fmt.Println("Audit logging enabled.")
}

// blockchainDisableAuditLogging disables transaction audit logging.
func blockchainDisableAuditLogging(ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()
    ledgerInstance.ComplianceLedger.SetAuditLogging(false)
    fmt.Println("Audit logging disabled.")
}

// blockchainSetSubblockCapacity sets the maximum capacity for sub-blocks.
func blockchainSetSubblockCapacity(capacity int, ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetSubblockCapacity(capacity); err != nil {
        return fmt.Errorf("failed to set subblock capacity: %v", err)
    }
    fmt.Printf("Subblock capacity set to %d.\n", capacity)
    return nil
}

// blockchainGetSubblockCapacity retrieves the current sub-block capacity.
func blockchainGetSubblockCapacity(ledgerInstance *ledger.Ledger) (int, error) {
    mutex.Lock()
    defer mutex.Unlock()
    capacity, err := ledgerInstance.BlockchainConsensusCoinLedger.GetSubblockCapacity()
    if err != nil {
        return 0, fmt.Errorf("failed to get subblock capacity: %v", err)
    }
    fmt.Printf("Subblock capacity is %d.\n", capacity)
    return capacity, nil
}

// blockchainMonitorSubblockCapacity monitors and logs the capacity status of sub-blocks.
func blockchainMonitorSubblockCapacity(ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()
    capacities := ledgerInstance.BlockchainConsensusCoinLedger.GetSubblockCapacityData()
    for subblockID, capacity := range capacities {
        fmt.Printf("Subblock %s capacity usage: %d.\n", subblockID, capacity)
    }
}

// blockchainAuditSubblockCapacity audits the capacity usage across sub-blocks.
func blockchainAuditSubblockCapacity(ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()
    err := ledgerInstance.BlockchainConsensusCoinLedger.AuditSubblockCapacity()
    if err != nil {
        fmt.Println("Failed to audit subblock capacity:", err)
    } else {
        fmt.Println("Subblock capacity audit completed.")
    }
}

// blockchainLogSubblockCapacity logs capacity usage for a sub-block to the ledger.
func blockchainLogSubblockCapacity(subblockID string, capacity int, ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()
    err := ledgerInstance.BlockchainConsensusCoinLedger.LogSubblockCapacity(subblockID, capacity)
    if err != nil {
        fmt.Printf("Failed to log capacity for subblock %s: %v\n", subblockID, err)
    }
}

// blockchainFetchSubblockCapacityHistory retrieves the capacity history for sub-blocks.
func blockchainFetchSubblockCapacityHistory(subblockID string, ledgerInstance *ledger.Ledger) ([]int, error) {
    mutex.Lock()
    defer mutex.Unlock()
    history, err := ledgerInstance.BlockchainConsensusCoinLedger.GetSubblockCapacityHistory(subblockID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch subblock capacity history: %v", err)
    }
    return history, nil
}

// blockchainSetSynthronCoinDenomination sets the denomination of Synthron Coin.
func blockchainSetSynthronCoinDenomination(denomination string, ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()
    err := ledgerInstance.BlockchainConsensusCoinLedger.SetSynthronCoinDenomination(denomination)
    if err != nil {
        return fmt.Errorf("failed to set Synthron Coin denomination: %v", err)
    }
    fmt.Printf("Synthron Coin denomination set to %s.\n", denomination)
    return nil
}


// blockchainGetSynthronCoinDenomination retrieves the current Synthron Coin denomination.
func blockchainGetSynthronCoinDenomination(ledgerInstance *ledger.Ledger) (string, error) {
    mutex.Lock()
    defer mutex.Unlock()
    denomination, err := ledgerInstance.BlockchainConsensusCoinLedger.GetSynthronCoinDenomination()
    if err != nil {
        return "", fmt.Errorf("failed to get Synthron Coin denomination: %v", err)
    }
    fmt.Printf("Synthron Coin denomination is %s.\n", denomination)
    return denomination, nil
}

// blockchainAuditCoinDenominations audits the denomination settings of Synthron Coin.
func blockchainAuditCoinDenominations(ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()
    err := ledgerInstance.BlockchainConsensusCoinLedger.AuditCoinDenominations()
    if err != nil {
        fmt.Println("Failed to audit Synthron Coin denominations:", err)
    } else {
        fmt.Println("Synthron Coin denomination audit completed.")
    }
}

// blockchainTrackCoinDenominations tracks denomination changes in Synthron Coin.
func blockchainTrackCoinDenominations(change string, ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()
    err := ledgerInstance.BlockchainConsensusCoinLedger.TrackCoinDenominationChange(change)
    if err != nil {
        fmt.Printf("Failed to track Synthron Coin denomination change: %v\n", err)
    }
}

// blockchainMonitorCoinDenominations monitors and logs changes in Synthron Coin denominations.
func blockchainMonitorCoinDenominations(ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()
    denominations := ledgerInstance.BlockchainConsensusCoinLedger.GetCoinDenominationData()
    for date, denomination := range denominations {
        fmt.Printf("On %s, Synthron Coin denomination was %s.\n", date, denomination)
    }
}

// blockchainSetSubblockCacheLimit sets the cache limit for sub-blocks.
func blockchainSetSubblockCacheLimit(limit int, ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()
    err := ledgerInstance.BlockchainConsensusCoinLedger.SetSubblockCacheLimit(limit)
    if err != nil {
        return fmt.Errorf("failed to set subblock cache limit: %v", err)
    }
    fmt.Printf("Subblock cache limit set to %d.\n", limit)
    return nil
}


// blockchainGetSubblockCacheLimit retrieves the cache limit for sub-blocks.
func blockchainGetSubblockCacheLimit(ledgerInstance *ledger.Ledger) (int, error) {
    mutex.Lock()
    defer mutex.Unlock()

    limit, err := ledgerInstance.BlockchainConsensusCoinLedger.GetSubblockCacheLimit()
    if err != nil {
        return 0, fmt.Errorf("failed to get subblock cache limit: %v", err)
    }
    fmt.Printf("Subblock cache limit is %d.\n", limit)
    return limit, nil
}

// blockchainEnableCacheMonitoring enables monitoring of cache usage in the ledger.
func blockchainEnableCacheMonitoring(ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()

    ledgerInstance.StorageLedger.SetCacheMonitoring(true)
    fmt.Println("Cache monitoring enabled.")
}

// blockchainDisableCacheMonitoring disables monitoring of cache usage in the ledger.
func blockchainDisableCacheMonitoring(ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()

    ledgerInstance.StorageLedger.SetCacheMonitoring(false)
    fmt.Println("Cache monitoring disabled.")
}

// blockchainAuditCacheUsage audits cache usage in the system.
func blockchainAuditCacheUsage(ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()

    err := ledgerInstance.ComplianceLedger.AuditCacheUsage()
    if err != nil {
        fmt.Println("Failed to audit cache usage:", err)
    } else {
        fmt.Println("Cache usage audit completed.")
    }
}

// blockchainTrackCacheUsage tracks cache usage over time.
func blockchainTrackCacheUsage(cacheData map[string]int, ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()

    for cacheID, usage := range cacheData {
        err := ledgerInstance.StorageLedger.RecordCacheUsage(cacheID, usage)
        if err != nil {
            fmt.Printf("Failed to track cache usage for %s: %v\n", cacheID, err)
        }
    }
}

// blockchainFetchCacheUsageHistory retrieves the historical cache usage data.
func blockchainFetchCacheUsageHistory(ledgerInstance *ledger.Ledger) (map[string]int, error) {
    mutex.Lock()
    defer mutex.Unlock()

    history, err := ledgerInstance.StorageLedger.GetCacheUsageHistory()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch cache usage history: %v", err)
    }
    return history, nil
}