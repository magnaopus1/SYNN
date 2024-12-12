package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/defi"
)

// Configuration for liquidity pool enforcement automation
const (
	LiquidityPoolCheckInterval        = 30 * time.Second // Interval to check liquidity pool compliance
	MinLiquidityThreshold             = 1000.0           // Minimum liquidity threshold per pool
	MaxTransactionLimitPerPool        = 100000.0         // Maximum transaction limit per liquidity pool
	HighVolumeTransactionThreshold    = 50000.0          // Threshold for high-volume transactions in pools
)

// LiquidityPoolEnforcementAutomation monitors and enforces compliance within liquidity pools
type LiquidityPoolEnforcementAutomation struct {
	defiManager          *defi.DeFiManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	liquidityStatusMap   map[string]float64 // Tracks liquidity status of each pool
	highVolumeTxnCount   map[string]int     // Tracks high-volume transactions per pool
}

// NewLiquidityPoolEnforcementAutomation initializes the liquidity pool enforcement automation
func NewLiquidityPoolEnforcementAutomation(defiManager *defi.DeFiManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *LiquidityPoolEnforcementAutomation {
	return &LiquidityPoolEnforcementAutomation{
		defiManager:         defiManager,
		consensusEngine:     consensusEngine,
		ledgerInstance:      ledgerInstance,
		enforcementMutex:    enforcementMutex,
		liquidityStatusMap:  make(map[string]float64),
		highVolumeTxnCount:  make(map[string]int),
	}
}

// StartLiquidityPoolEnforcement begins continuous monitoring and enforcement of liquidity pool compliance
func (automation *LiquidityPoolEnforcementAutomation) StartLiquidityPoolEnforcement() {
	ticker := time.NewTicker(LiquidityPoolCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkLiquidityPoolCompliance()
		}
	}()
}

// checkLiquidityPoolCompliance monitors each liquidity pool's status and restricts non-compliant pools
func (automation *LiquidityPoolEnforcementAutomation) checkLiquidityPoolCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyLiquidityLevels()
	automation.enforceTransactionLimits()
}

// verifyLiquidityLevels ensures that each liquidity pool meets the minimum liquidity threshold
func (automation *LiquidityPoolEnforcementAutomation) verifyLiquidityLevels() {
	for _, poolID := range automation.defiManager.GetAllPools() {
		liquidity := automation.defiManager.GetPoolLiquidity(poolID)
		automation.liquidityStatusMap[poolID] = liquidity

		if liquidity < MinLiquidityThreshold {
			fmt.Printf("Liquidity enforcement triggered for pool %s due to low liquidity.\n", poolID)
			automation.applyRestriction(poolID, "Insufficient Liquidity")
		}
	}
}

// enforceTransactionLimits monitors and restricts high-volume transactions in liquidity pools
func (automation *LiquidityPoolEnforcementAutomation) enforceTransactionLimits() {
	for _, poolID := range automation.defiManager.GetAllPools() {
		transactionVolume := automation.defiManager.GetPoolTransactionVolume(poolID)

		if transactionVolume > HighVolumeTransactionThreshold {
			automation.highVolumeTxnCount[poolID]++
			if automation.highVolumeTxnCount[poolID] > 1 {
				fmt.Printf("High-volume transaction limit enforcement triggered for pool %s.\n", poolID)
				automation.applyRestriction(poolID, "High Transaction Volume Limit Exceeded")
			}
		} else {
			automation.highVolumeTxnCount[poolID] = 0 // Reset count if within limits
		}
	}
}

// applyRestriction restricts operations for liquidity pools that fail to meet compliance standards
func (automation *LiquidityPoolEnforcementAutomation) applyRestriction(poolID, reason string) {
	err := automation.defiManager.RestrictPool(poolID)
	if err != nil {
		fmt.Printf("Failed to restrict non-compliant liquidity pool %s: %v\n", poolID, err)
		automation.logPoolAction(poolID, "Restriction Failed", reason)
	} else {
		fmt.Printf("Liquidity pool %s restricted due to %s.\n", poolID, reason)
		automation.logPoolAction(poolID, "Restricted", reason)
	}
}

// logPoolAction securely logs actions related to liquidity pool enforcement
func (automation *LiquidityPoolEnforcementAutomation) logPoolAction(poolID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, Pool ID: %s, Reason: %s", action, poolID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("liquidity-pool-enforcement-%s-%d", poolID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Liquidity Pool Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log liquidity pool enforcement action for pool %s: %v\n", poolID, err)
	} else {
		fmt.Println("Liquidity pool enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *LiquidityPoolEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualLiquidityCheck allows administrators to manually enforce compliance for a specific liquidity pool
func (automation *LiquidityPoolEnforcementAutomation) TriggerManualLiquidityCheck(poolID string) {
	fmt.Printf("Manually triggering compliance check for liquidity pool: %s\n", poolID)

	liquidity := automation.defiManager.GetPoolLiquidity(poolID)
	if liquidity < MinLiquidityThreshold {
		automation.applyRestriction(poolID, "Manual Trigger: Insufficient Liquidity")
	} else {
		fmt.Printf("Liquidity pool %s meets compliance standards.\n", poolID)
		automation.logPoolAction(poolID, "Manual Compliance Check Passed", "Liquidity Compliance Verified")
	}
}
