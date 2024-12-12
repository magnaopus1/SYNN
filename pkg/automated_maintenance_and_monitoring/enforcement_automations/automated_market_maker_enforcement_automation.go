package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/market"
)

// Configurable constants
const (
	LiquidityCheckInterval  = 15 * time.Second
	PriceManipulationMargin = 0.05 // 5% allowed price deviation before triggering manipulation check
	LiquidityThreshold      = 1000 // Minimum liquidity in the pool
)

// AutomatedMarketMakerEnforcementAutomation manages AMM compliance and manipulation prevention
type AutomatedMarketMakerEnforcementAutomation struct {
	marketManager      *market.MarketManager
	consensusEngine    *consensus.SynnergyConsensus
	ledgerInstance     *ledger.Ledger
	enforcementMutex   *sync.RWMutex
	liquidityPools     map[string]float64 // Track liquidity in each pool
}

// NewAutomatedMarketMakerEnforcementAutomation initializes AMM enforcement automation
func NewAutomatedMarketMakerEnforcementAutomation(marketManager *market.MarketManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *AutomatedMarketMakerEnforcementAutomation {
	return &AutomatedMarketMakerEnforcementAutomation{
		marketManager:      marketManager,
		consensusEngine:    consensusEngine,
		ledgerInstance:     ledgerInstance,
		enforcementMutex:   enforcementMutex,
		liquidityPools:     make(map[string]float64),
	}
}

// StartAMMEnforcement begins continuous monitoring and enforcement for AMM operations
func (automation *AutomatedMarketMakerEnforcementAutomation) StartAMMEnforcement() {
	ticker := time.NewTicker(LiquidityCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkLiquidityThresholds()
			automation.detectPriceManipulation()
		}
	}()
}

// checkLiquidityThresholds enforces minimum liquidity in each pool
func (automation *AutomatedMarketMakerEnforcementAutomation) checkLiquidityThresholds() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for poolID, liquidity := range automation.marketManager.GetAllLiquidityPools() {
		automation.liquidityPools[poolID] = liquidity

		if liquidity < LiquidityThreshold {
			automation.enforceLiquidity(poolID)
		}
	}
}

// enforceLiquidity adds liquidity to a pool below the threshold
func (automation *AutomatedMarketMakerEnforcementAutomation) enforceLiquidity(poolID string) {
	err := automation.marketManager.AddLiquidity(poolID, LiquidityThreshold - automation.liquidityPools[poolID])
	if err != nil {
		fmt.Printf("Failed to enforce liquidity for pool %s: %v\n", poolID, err)
		return
	}

	fmt.Printf("Enforced liquidity for pool %s. Liquidity adjusted to threshold.\n", poolID)
	automation.logAMMEnforcementAction(poolID, "Liquidity Enforced")
}

// detectPriceManipulation checks for abnormal price changes in AMM pools
func (automation *AutomatedMarketMakerEnforcementAutomation) detectPriceManipulation() {
	for poolID, priceChange := range automation.marketManager.GetRecentPriceChanges() {
		if priceChange > PriceManipulationMargin {
			automation.preventPriceManipulation(poolID)
		}
	}
}

// preventPriceManipulation applies measures to stabilize prices in case of manipulation
func (automation *AutomatedMarketMakerEnforcementAutomation) preventPriceManipulation(poolID string) {
	err := automation.marketManager.StabilizePrice(poolID)
	if err != nil {
		fmt.Printf("Failed to prevent price manipulation for pool %s: %v\n", poolID, err)
		return
	}

	fmt.Printf("Price manipulation prevented in pool %s.\n", poolID)
	automation.logAMMEnforcementAction(poolID, "Price Manipulation Prevented")
}

// logAMMEnforcementAction securely logs actions taken to enforce AMM compliance
func (automation *AutomatedMarketMakerEnforcementAutomation) logAMMEnforcementAction(poolID string, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Pool: %s", action, poolID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("amm-enforcement-%s-%d", poolID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "AMM Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log AMM enforcement action for pool %s: %v\n", poolID, err)
	} else {
		fmt.Println("AMM enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *AutomatedMarketMakerEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualLiquidityAdjustment allows administrators to manually adjust liquidity in a pool
func (automation *AutomatedMarketMakerEnforcementAutomation) TriggerManualLiquidityAdjustment(poolID string, amount float64) {
	fmt.Printf("Manually adjusting liquidity for pool: %s by %f\n", poolID, amount)

	err := automation.marketManager.AddLiquidity(poolID, amount)
	if err != nil {
		fmt.Printf("Failed to manually adjust liquidity for pool %s: %v\n", poolID, err)
		return
	}

	automation.logAMMEnforcementAction(poolID, "Manually Adjusted Liquidity")
}
