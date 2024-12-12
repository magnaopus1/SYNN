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

// TransactionFeeRefundAutomation automates the refund of unused gas and transaction fees.
type TransactionFeeRefundAutomation struct {
	ledgerInstance  *ledger.Ledger
	feeManager      *transactions.FeeManager
	transactionPool *transactions.TransactionPool
	mutex           sync.Mutex
	stopChan        chan bool
}

// NewTransactionFeeRefundAutomation initializes the automation for unused gas and fee refunds.
func NewTransactionFeeRefundAutomation(ledgerInstance *ledger.Ledger, transactionPool *transactions.TransactionPool, feeManager *transactions.FeeManager) *TransactionFeeRefundAutomation {
	return &TransactionFeeRefundAutomation{
		ledgerInstance:  ledgerInstance,
		transactionPool: transactionPool,
		feeManager:      feeManager,
		stopChan:        make(chan bool),
	}
}

// Start begins the process of monitoring completed transactions and refunding unused gas and fees.
func (t *TransactionFeeRefundAutomation) Start() {
	go t.runRefundLoop()
	log.Println("Transaction Fee Refund Automation started.")
}

// Stop halts the fee refund automation.
func (t *TransactionFeeRefundAutomation) Stop() {
	t.stopChan <- true
	log.Println("Transaction Fee Refund Automation stopped.")
}

// runRefundLoop periodically checks for completed transactions and triggers gas/fee refunds.
func (t *TransactionFeeRefundAutomation) runRefundLoop() {
	ticker := time.NewTicker(500 * time.Millisecond) // Check every 500ms for completed transactions
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.processTransactionRefunds()
		case <-t.stopChan:
			return
		}
	}
}

// processTransactionRefunds handles refunding unused gas and fees for completed transactions.
func (t *TransactionFeeRefundAutomation) processTransactionRefunds() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Retrieve all transactions that are completed but have not yet been refunded
	completedTransactions, err := t.ledgerInstance.GetCompletedTransactionsPendingRefund()
	if err != nil {
		log.Printf("Failed to fetch completed transactions for refund: %v", err)
		return
	}

	for _, tx := range completedTransactions {
		err := t.refundUnusedGas(tx)
		if err != nil {
			log.Printf("Failed to refund gas for transaction %s: %v", tx.ID, err)
		}

		err = t.refundUnusedFees(tx)
		if err != nil {
			log.Printf("Failed to refund fees for transaction %s: %v", tx.ID, err)
		}
	}
}

// refundUnusedGas calculates and refunds unused gas after transaction execution.
func (t *TransactionFeeRefundAutomation) refundUnusedGas(tx common.Transaction) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Calculate the difference between gas used and gas limit
	if tx.GasUsed < tx.GasLimit {
		// Refund the difference in gas
		unusedGas := tx.GasLimit - tx.GasUsed
		refundAmount := float64(unusedGas) * tx.GasPricePerUnit

		// Log the refund
		log.Printf("Refunding %f gas for transaction %s", refundAmount, tx.ID)

		// Encrypt the refund details before logging to the ledger
		encryptedRefund, err := encryption.EncryptData(fmt.Sprintf("%f", refundAmount), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt gas refund: %v", err)
		}

		// Log the refund in the ledger
		err = t.ledgerInstance.RecordGasRefund(tx.ID, encryptedRefund)
		if err != nil {
			return fmt.Errorf("failed to record gas refund in ledger: %v", err)
		}
	}

	return nil
}

// refundUnusedFees handles refunding any unused transaction fees.
func (t *TransactionFeeRefundAutomation) refundUnusedFees(tx common.Transaction) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Retrieve actual fee after validation (for example, through consensus validation results)
	actualFee, err := t.ledgerInstance.GetTransactionFeeAfterValidation(tx.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve actual transaction fee: %v", err)
	}

	// If the actual fee is lower than the fee paid, refund the difference
	if actualFee < tx.Fee {
		feeRefund := tx.Fee - actualFee

		// Log the refund
		log.Printf("Refunding %f unused fees for transaction %s", feeRefund, tx.ID)

		// Encrypt the refund details before logging to the ledger
		encryptedRefund, err := encryption.EncryptData(fmt.Sprintf("%f", feeRefund), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt fee refund: %v", err)
		}

		// Log the refund in the ledger
		err = t.ledgerInstance.RecordFeeRefund(tx.ID, encryptedRefund)
		if err != nil {
			return fmt.Errorf("failed to record fee refund in ledger: %v", err)
		}
	}

	return nil
}
