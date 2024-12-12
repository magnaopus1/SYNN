package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/defi/liquidity_pool"
)

const (
	LiquidityCheckInterval     = 5 * time.Minute // Interval for checking liquidity pool conditions
	MinimumLiquidityThreshold  = 10000.0         // Minimum liquidity required in the pool (example value)
	AutoLiquidityProvisionRate = 0.02            // Auto-provision rate, 2% of the pool to stabilize liquidity
	LiquidityLoggingInterval   = 30 * time.Minute // Interval for logging liquidity status to the ledger
)

// DeFiLiquidityExecutionAutomation manages liquidity in the DeFi space
type DeFiLiquidityExecutionAutomation struct {
	liquidityPool   *liquidity_pool.LiquidityPool  // Liquidity pool for the DeFi system
	consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validation and integration
	ledgerInstance  *ledger.Ledger                 // Ledger to track liquidity events
	liquidityMutex  *sync.RWMutex                  // Mutex for thread-safe operations
}

// NewDeFiLiquidityExecutionAutomation initializes a new DeFiLiquidityExecutionAutomation instance
func NewDeFiLiquidityExecutionAutomation(liquidityPool *liquidity_pool.LiquidityPool, consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, liquidityMutex *sync.RWMutex) *DeFiLiquidityExecutionAutomation {
	return &DeFiLiquidityExecutionAutomation{
		liquidityPool:  liquidityPool,
		consensusEngine: consensusEngine,
		ledgerInstance: ledgerInstance,
		liquidityMutex: liquidityMutex,
	}
}

// StartDeFiLiquidityMonitoring starts the automated monitoring and management of the liquidity pool
func (automation *DeFiLiquidityExecutionAutomation) StartDeFiLiquidityMonitoring() {
	ticker := time.NewTicker(LiquidityCheckInterval)
	go func() {
		for range ticker.C {
			automation.monitorLiquidityPool()
		}
	}()

	// Log liquidity pool status periodically
	logTicker := time.NewTicker(LiquidityLoggingInterval)
	go func() {
		for range logTicker.C {
			automation.logLiquidityStatus()
		}
	}()
}

// monitorLiquidityPool checks the liquidity pool and triggers automated liquidity provisioning if necessary
func (automation *DeFiLiquidityExecutionAutomation) monitorLiquidityPool() {
	automation.liquidityMutex.Lock()
	defer automation.liquidityMutex.Unlock()

	currentLiquidity := automation.liquidityPool.GetCurrentLiquidity()

	// If liquidity drops below the threshold, trigger automatic liquidity provisioning
	if currentLiquidity < MinimumLiquidityThreshold {
		automation.provideAutomaticLiquidity()
	} else {
		fmt.Printf("Current liquidity is sufficient: %.2f\n", currentLiquidity)
	}
}

// provideAutomaticLiquidity automatically provides liquidity to stabilize the pool
func (automation *DeFiLiquidityExecutionAutomation) provideAutomaticLiquidity() {
	amountToProvide := automation.liquidityPool.GetTotalPoolValue() * AutoLiquidityProvisionRate
	err := automation.liquidityPool.AddLiquidity(amountToProvide)
	if err != nil {
		fmt.Println("Failed to add automatic liquidity:", err)
		return
	}

	// Log the automatic liquidity provision in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("automatic-liquidity-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Liquidity Provision",
		Status:    "Success",
		Details:   fmt.Sprintf("Automatically provided %.2f liquidity to stabilize the pool.", amountToProvide),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log automatic liquidity provision:", err)
	} else {
		fmt.Println("Automatic liquidity provision successfully logged in the ledger.")
	}
}

// logLiquidityStatus logs the current status of the liquidity pool to the ledger
func (automation *DeFiLiquidityExecutionAutomation) logLiquidityStatus() {
	currentLiquidity := automation.liquidityPool.GetCurrentLiquidity()

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("liquidity-status-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Liquidity Status",
		Status:    "Success",
		Details:   fmt.Sprintf("Current liquidity: %.2f", currentLiquidity),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log liquidity status:", err)
	} else {
		fmt.Println("Liquidity status successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DeFiLiquidityExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualLiquidity allows for manual addition of liquidity to the pool by administrators
func (automation *DeFiLiquidityExecutionAutomation) TriggerManualLiquidity(amount float64) {
	automation.liquidityMutex.Lock()
	defer automation.liquidityMutex.Unlock()

	err := automation.liquidityPool.AddLiquidity(amount)
	if err != nil {
		fmt.Printf("Failed to manually add liquidity: %.2f\n", err)
		return
	}

	// Log the manual liquidity provision in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("manual-liquidity-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Manual Liquidity Provision",
		Status:    "Success",
		Details:   fmt.Sprintf("Manually provided %.2f liquidity to the pool.", amount),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log manual liquidity provision:", err)
	} else {
		fmt.Println("Manual liquidity provision successfully logged in the ledger.")
	}
}
