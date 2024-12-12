package execution_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    LiquidityCheckInterval     = 5 * time.Second  // Interval for checking liquidity levels in AMM pools
    PriceUpdateInterval        = 10 * time.Second // Interval for updating token prices
    PriceImpactThreshold       = 0.05             // Threshold for price impact triggering market maker actions (5%)
    LiquidityProvisionThreshold = 1000            // Minimum liquidity level before triggering new liquidity provision
)

// AutomaticMarketMakerExecutionAutomation handles automated market-making operations
type AutomaticMarketMakerExecutionAutomation struct {
    consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine
    ledgerInstance    *ledger.Ledger                         // Ledger instance for recording market actions
    stateMutex        *sync.RWMutex                          // Mutex for thread-safe operations
}

// NewAutomaticMarketMakerExecutionAutomation initializes the market maker automation
func NewAutomaticMarketMakerExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AutomaticMarketMakerExecutionAutomation {
    return &AutomaticMarketMakerExecutionAutomation{
        consensusEngine:  consensusEngine,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
    }
}

// StartMarketMakerMonitor begins the automatic market maker execution and liquidity monitoring
func (automation *AutomaticMarketMakerExecutionAutomation) StartMarketMakerMonitor() {
    liquidityTicker := time.NewTicker(LiquidityCheckInterval)
    priceTicker := time.NewTicker(PriceUpdateInterval)

    go func() {
        for range liquidityTicker.C {
            automation.checkLiquidityLevels()
        }
    }()

    go func() {
        for range priceTicker.C {
            automation.updateTokenPrices()
        }
    }()
}

// checkLiquidityLevels monitors AMM pools for liquidity and triggers rebalancing if necessary
func (automation *AutomaticMarketMakerExecutionAutomation) checkLiquidityLevels() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    liquidityPools := automation.consensusEngine.GetLiquidityPools()

    for _, pool := range liquidityPools {
        if pool.Liquidity < LiquidityProvisionThreshold {
            fmt.Printf("Low liquidity detected in pool %s. Triggering liquidity provision.\n", pool.ID)
            automation.provideLiquidity(pool)
        }
    }
}

// provideLiquidity adds liquidity to a pool and logs the event to the ledger
func (automation *AutomaticMarketMakerExecutionAutomation) provideLiquidity(pool common.LiquidityPool) {
    success := automation.consensusEngine.AddLiquidity(pool.ID, LiquidityProvisionThreshold)

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("liquidity-provision-%s", pool.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Liquidity Provision",
        Status:    "Success",
        Details:   fmt.Sprintf("Added liquidity to pool %s.", pool.ID),
    }

    if !success {
        entry.Status = "Failure"
        entry.Details = fmt.Sprintf("Failed to add liquidity to pool %s.", pool.ID)
    }

    encryptedEntry := automation.encryptData(entry.Details)
    entry.Details = encryptedEntry

    if err := automation.ledgerInstance.AddEntry(entry); err != nil {
        fmt.Printf("Error logging liquidity provision: %v\n", err)
    } else {
        fmt.Printf("Liquidity provision logged successfully for pool %s.\n", pool.ID)
    }
}

// updateTokenPrices fetches the latest market data and updates token prices accordingly
func (automation *AutomaticMarketMakerExecutionAutomation) updateTokenPrices() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    pools := automation.consensusEngine.GetLiquidityPools()

    for _, pool := range pools {
        priceImpact := automation.consensusEngine.CalculatePriceImpact(pool.ID)

        if priceImpact > PriceImpactThreshold {
            fmt.Printf("High price impact detected in pool %s. Rebalancing market.\n", pool.ID)
            automation.rebalancePool(pool)
        }
    }
}

// rebalancePool rebalances a liquidity pool based on market conditions and logs the action to the ledger
func (automation *AutomaticMarketMakerExecutionAutomation) rebalancePool(pool common.LiquidityPool) {
    success := automation.consensusEngine.RebalancePool(pool.ID)

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("pool-rebalance-%s", pool.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Pool Rebalancing",
        Status:    "Success",
        Details:   fmt.Sprintf("Rebalanced liquidity pool %s due to high price impact.", pool.ID),
    }

    if !success {
        entry.Status = "Failure"
        entry.Details = fmt.Sprintf("Failed to rebalance liquidity pool %s.", pool.ID)
    }

    encryptedEntry := automation.encryptData(entry.Details)
    entry.Details = encryptedEntry

    if err := automation.ledgerInstance.AddEntry(entry); err != nil {
        fmt.Printf("Error logging pool rebalance: %v\n", err)
    } else {
        fmt.Printf("Pool rebalance successfully logged for pool %s.\n", pool.ID)
    }
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *AutomaticMarketMakerExecutionAutomation) encryptData(data string) string {
    encryptedData, err := encryption.EncryptData([]byte(data))
    if err != nil {
        fmt.Println("Error encrypting data:", err)
        return data
    }
    return string(encryptedData)
}
