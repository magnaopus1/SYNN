package automations

import (
	"fmt"
	"log"
	"sync"
	"time"
	"synnergy_network_demo/transactions"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
)

// TransactionFeeDistributionAutomation automates the distribution of transaction fees
// across various pools every time a sub-block is validated.
type TransactionFeeDistributionAutomation struct {
	ledgerInstance      *ledger.Ledger
	distributionManager *transactions.TransactionDistributionManager
	mutex               sync.Mutex
	stopChan            chan bool
}

// NewTransactionFeeDistributionAutomation initializes a new automation instance.
func NewTransactionFeeDistributionAutomation(ledgerInstance *ledger.Ledger) *TransactionFeeDistributionAutomation {
	return &TransactionFeeDistributionAutomation{
		ledgerInstance:      ledgerInstance,
		distributionManager: transactions.NewTransactionDistributionManager(ledgerInstance),
		stopChan:            make(chan bool),
	}
}

// Start begins the continuous monitoring and distribution process, triggered by sub-block validation.
func (t *TransactionFeeDistributionAutomation) Start() {
	go t.runAutomationLoop()
	log.Println("Transaction Fee Distribution Automation started.")
}

// Stop stops the continuous automation.
func (t *TransactionFeeDistributionAutomation) Stop() {
	t.stopChan <- true
	log.Println("Transaction Fee Distribution Automation stopped.")
}

// runAutomationLoop continuously checks for new validated sub-blocks and distributes fees.
func (t *TransactionFeeDistributionAutomation) runAutomationLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.processSubBlockFeeDistribution()
		case <-t.stopChan:
			return
		}
	}
}

// processSubBlockFeeDistribution processes fee distribution for newly validated sub-blocks.
func (t *TransactionFeeDistributionAutomation) processSubBlockFeeDistribution() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Fetch the list of recently validated sub-blocks from the ledger
	subBlocks, err := t.ledgerInstance.GetValidatedSubBlocks()
	if err != nil {
		log.Printf("Failed to fetch validated sub-blocks: %v", err)
		return
	}

	// Iterate through each sub-block and process fee distribution
	for _, subBlock := range subBlocks {
		totalTransactionFees, err := t.calculateTotalFees(subBlock)
		if err != nil {
			log.Printf("Error calculating total transaction fees for sub-block %s: %v", subBlock.ID, err)
			continue
		}

		err = t.distributionManager.DistributeRewards(subBlock.ID, totalTransactionFees)
		if err != nil {
			log.Printf("Error distributing rewards for sub-block %s: %v", subBlock.ID, err)
		} else {
			log.Printf("Successfully distributed rewards for sub-block %s", subBlock.ID)
		}
	}
}

// calculateTotalFees calculates the total fees for a validated sub-block.
func (t *TransactionFeeDistributionAutomation) calculateTotalFees(subBlock common.SubBlock) (float64, error) {
	var totalFees float64
	for _, tx := range subBlock.Transactions {
		totalFees += tx.Fee
	}
	return totalFees, nil
}

// SetupTriggers sets up the trigger to automatically process fees when a sub-block is validated.
func (t *TransactionFeeDistributionAutomation) SetupTriggers() {
	// Trigger every time a sub-block is validated
	t.ledgerInstance.OnSubBlockValidated(t.processSubBlockFeeDistribution)
	log.Println("Transaction fee distribution trigger setup.")
}

// recordSubBlockTransaction records the fee distribution for the sub-block in the ledger with encryption.
func (t *TransactionFeeDistributionAutomation) recordSubBlockTransaction(subBlockID, description string, amount float64) error {
	encryptedAmount, err := encryption.EncryptData(fmt.Sprintf("%f", amount), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt fee distribution for sub-block: %v", err)
	}
	return t.ledgerInstance.RecordSubBlockTransaction(subBlockID, description, encryptedAmount)
}
