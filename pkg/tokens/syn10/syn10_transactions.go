package syn10

import (
	"errors"
	"fmt"
	"sync"
	"time"

  	"synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)



// NewBatchTransferProcessor creates a new BatchTransferProcessor.
func NewBatchTransferProcessor(ledger *ledger.TokenLedger, validator *validators.TransferValidator, encryptionService *encryption.Service) *BatchTransferProcessor {
	return &BatchTransferProcessor{
		ledger:            ledger,
		validator:         validator,
		encryptionService: encryptionService,
	}
}

// NewOwnershipTransferProcessor creates a new OwnershipTransferProcessor.
func NewOwnershipTransferProcessor(ledger *ledger.TokenLedger, validator *validators.TransferValidator, encryptionService *encryption.Service) *OwnershipTransferProcessor {
	return &OwnershipTransferProcessor{
		ledger:            ledger,
		validator:         validator,
		encryptionService: encryptionService,
	}
}

// NewSaleHistoryProcessor creates a new SaleHistoryProcessor.
func NewSaleHistoryProcessor(ledger *ledger.TokenLedger, validator *validators.SaleValidator, encryptionService *encryption.Service) *SaleHistoryProcessor {
	return &SaleHistoryProcessor{
		ledger:            ledger,
		validator:         validator,
		encryptionService: encryptionService,
	}
}

// NewSyn10FeeFreeTransactionProcessor creates a new FeeFreeTransactionProcessor.
func NewSyn10FeeFreeTransactionProcessor(ledger *ledger.TokenLedger, validator *validators.TransactionValidator, encryptionService *encryption.Service) *Syn10FeeFreeTransactionProcessor {
	return &Syn10FeeFreeTransactionProcessor{
		ledger:            ledger,
		validator:         validator,
		encryptionService: encryptionService,
	}
}

// ProcessBatch processes a batch of transfers and ensures atomic validation.
func (b *BatchTransferProcessor) ProcessBatch(batch BatchTransfer) error {
	if len(batch.Transfers) == 0 {
		return errors.New("no transfers in the batch")
	}

	if err := b.validateBatch(batch); err != nil {
		return fmt.Errorf("batch validation failed: %v", err)
	}

	return b.processTransfers(batch)
}

// validateBatch validates all transfers in the batch.
func (b *BatchTransferProcessor) validateBatch(batch BatchTransfer) error {
	for _, transfer := range batch.Transfers {
		if err := b.validator.ValidateTransfer(batch.TokenID, batch.SenderAddress, transfer.ReceiverAddress, transfer.Amount); err != nil {
			return err
		}
	}
	return nil
}

// processTransfers executes the batch transfers.
func (b *BatchTransferProcessor) processTransfers(batch BatchTransfer) error {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	errors := make([]error, len(batch.Transfers))

	for i, transfer := range batch.Transfers {
		wg.Add(1)
		go func(i int, transfer TransferDetail) {
			defer wg.Done()
			if err := b.ledger.TransferTokens(batch.TokenID, batch.SenderAddress, transfer.ReceiverAddress, transfer.Amount); err != nil {
				mutex.Lock()
				errors[i] = err
				mutex.Unlock()
			}
		}(i, transfer)
	}

	wg.Wait()

	for _, err := range errors {
		if err != nil {
			return err
		}
	}

	// Encrypt and store batch details.
	encryptedVerificationID, err := b.encryptionService.Encrypt([]byte(batch.VerificationID))
	if err != nil {
		return fmt.Errorf("failed to encrypt verification ID: %v", err)
	}
	return b.ledger.StoreBatchDetails(batch.TokenID, encryptedVerificationID, batch.Transfers)
}

// TransferOwnership processes a token ownership transfer.
func (o *OwnershipTransferProcessor) TransferOwnership(transfer OwnershipTransfer) error {
	if err := o.validateTransfer(transfer); err != nil {
		return fmt.Errorf("ownership validation failed: %w", err)
	}

	encryptedVerificationID, err := o.encryptionService.Encrypt([]byte(transfer.VerificationID))
	if err != nil {
		return fmt.Errorf("failed to encrypt verification ID: %w", err)
	}

	if err := o.ledger.TransferTokens(transfer.TokenID, transfer.SenderAddress, transfer.ReceiverAddress, transfer.Amount); err != nil {
		return fmt.Errorf("ledger transfer failed: %w", err)
	}

	return o.ledger.LogOwnershipTransfer(transfer.TokenID, transfer.SenderAddress, transfer.ReceiverAddress, transfer.Amount, string(encryptedVerificationID))
}

// RecordSale logs a sale transaction to the ledger.
func (s *SaleHistoryProcessor) RecordSale(record SaleRecord) error {
	if err := s.validateSale(record); err != nil {
		return fmt.Errorf("sale validation failed: %w", err)
	}

	encryptedTransactionID, err := s.encryptionService.Encrypt([]byte(record.TransactionID))
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction ID: %w", err)
	}

	return s.ledger.LogSaleRecord(record.TokenID, record.SellerAddress, record.BuyerAddress, record.Amount, record.SalePrice, record.Timestamp, string(encryptedTransactionID))
}

// ProcessTransaction processes a fee-free token transaction.
func (p *Syn10FeeFreeTransactionProcessor) ProcessTransaction(tx TransactionFeeFree) error {
	if err := p.validateTransaction(tx); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	encryptedTransactionID, err := p.encryptionService.Encrypt([]byte(tx.TransactionID))
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction ID: %w", err)
	}

	return p.ledger.RecordTransaction(tx.TokenID, tx.FromAddress, tx.ToAddress, tx.Amount, tx.Timestamp, string(encryptedTransactionID), true)
}

// validateTransaction ensures the transaction meets all criteria.
func (p *Syn10FeeFreeTransactionProcessor) validateTransaction(tx TransactionFeeFree) error {
	if err := p.validator.ValidateSender(tx.FromAddress); err != nil {
		return fmt.Errorf("sender validation failed: %w", err)
	}
	if err := p.validator.ValidateReceiver(tx.ToAddress); err != nil {
		return fmt.Errorf("receiver validation failed: %w", err)
	}
	if err := p.validator.ValidateAmount(tx.Amount); err != nil {
		return fmt.Errorf("amount validation failed: %w", err)
	}
	if err := p.validator.ValidateTransactionID(tx.TransactionID); err != nil {
		return fmt.Errorf("transaction ID validation failed: %w", err)
	}
	return nil
}
