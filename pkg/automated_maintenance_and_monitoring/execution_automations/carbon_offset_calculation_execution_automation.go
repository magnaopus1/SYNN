package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/carbon"
)

const (
	OffsetCalculationInterval  = 30 * time.Minute  // Interval for running carbon offset calculations
	CarbonUnitPrice            = 10.0              // Price per carbon unit for calculation purposes (example)
	MaxAllowedCarbonEmission   = 1000.0            // Maximum allowed carbon emissions per node before offset is required
)

// CarbonOffsetExecutionAutomation automates the process of calculating and logging carbon offsets for the blockchain network
type CarbonOffsetExecutionAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validating offsets
	ledgerInstance    *ledger.Ledger                        // Ledger for recording offset actions
	carbonCalculator  *carbon.Calculator                    // Carbon emission and offset calculator
	stateMutex        *sync.RWMutex                         // Mutex for thread-safe operations
}

// NewCarbonOffsetExecutionAutomation initializes the carbon offset automation
func NewCarbonOffsetExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, carbonCalculator *carbon.Calculator, stateMutex *sync.RWMutex) *CarbonOffsetExecutionAutomation {
	return &CarbonOffsetExecutionAutomation{
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		carbonCalculator: carbonCalculator,
		stateMutex:       stateMutex,
	}
}

// StartOffsetCalculationMonitor starts monitoring for carbon offset calculations
func (automation *CarbonOffsetExecutionAutomation) StartOffsetCalculationMonitor() {
	ticker := time.NewTicker(OffsetCalculationInterval)

	go func() {
		for range ticker.C {
			automation.calculateAndLogOffsets()
		}
	}()
}

// calculateAndLogOffsets calculates carbon offsets and logs the result into the ledger
func (automation *CarbonOffsetExecutionAutomation) calculateAndLogOffsets() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch emission data from carbon calculator
	emissionData, err := automation.carbonCalculator.FetchEmissionData()
	if err != nil {
		fmt.Println("Error fetching carbon emission data:", err)
		return
	}

	// Calculate total offset required based on current emissions
	totalOffset := automation.calculateTotalOffset(emissionData)
	if totalOffset > 0 {
		automation.logOffsetCalculation(totalOffset)
	} else {
		fmt.Println("No carbon offset required at this time.")
	}
}

// calculateTotalOffset calculates the carbon offset required for emissions exceeding the threshold
func (automation *CarbonOffsetExecutionAutomation) calculateTotalOffset(emissionData carbon.EmissionData) float64 {
	excessEmission := emissionData.TotalEmission - MaxAllowedCarbonEmission
	if excessEmission > 0 {
		offsetRequired := excessEmission * CarbonUnitPrice
		fmt.Printf("Offset required for excess emissions: %.2f units\n", offsetRequired)
		return offsetRequired
	}
	return 0
}

// logOffsetCalculation securely logs the carbon offset calculation into the ledger
func (automation *CarbonOffsetExecutionAutomation) logOffsetCalculation(totalOffset float64) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("carbon-offset-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Carbon Offset Calculation",
		Status:    "Completed",
		Details:   fmt.Sprintf("Total offset required: %.2f carbon units.", totalOffset),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	if err := automation.ledgerInstance.AddEntry(entry); err != nil {
		fmt.Printf("Error logging carbon offset calculation: %v\n", err)
	} else {
		fmt.Println("Carbon offset calculation successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *CarbonOffsetExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
