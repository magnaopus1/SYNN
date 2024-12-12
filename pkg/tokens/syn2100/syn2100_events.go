package syn2100

import (
	"errors"
	"time"

)

// RecordInvoiceTokenizationEvent logs the event of an invoice being tokenized and updates the ledger.
func RecordInvoiceTokenizationEvent(token *common.SYN2100Token, issuer string, invoiceID string) error {
	// Create event log for tokenization
	tokenizationEvent := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Invoice Tokenization",
		Description: "Invoice " + invoiceID + " has been tokenized.",
		PerformedBy: issuer,
		EventDate:   time.Now(),
	}

	// Add event log to token's event history
	token.EventLogs = append(token.EventLogs, tokenizationEvent)

	// Update the ledger to reflect this tokenization event
	err := ledger.RecordEvent(token.TokenID, "Invoice Tokenization", tokenizationEvent)
	if err != nil {
		return errors.New("failed to log tokenization event in ledger")
	}

	return nil
}

// RecordTransferEvent logs the event of a token transfer between parties and updates the ledger.
func RecordTransferEvent(token *common.SYN2100Token, sender, recipient string, amount float64) error {
	// Create event log for token transfer
	transferEvent := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Token Transfer",
		Description: "Transferred token to " + recipient + " for amount: " + formatAmount(amount),
		PerformedBy: sender,
		EventDate:   time.Now(),
	}

	// Add event log to token's event history
	token.EventLogs = append(token.EventLogs, transferEvent)

	// Update the ledger to reflect the transfer event
	err := ledger.RecordEvent(token.TokenID, "Token Transfer", transferEvent)
	if err != nil {
		return errors.New("failed to log transfer event in ledger")
	}

	return nil
}

// RecordDynamicDiscountEvent logs the event of applying dynamic discounting and updates the ledger.
func RecordDynamicDiscountEvent(token *common.SYN2100Token, issuer string, originalAmount, discountedAmount float64) error {
	// Create event log for dynamic discounting
	discountEvent := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Dynamic Discount Applied",
		Description: "Original amount: " + formatAmount(originalAmount) + ", discounted to: " + formatAmount(discountedAmount),
		PerformedBy: issuer,
		EventDate:   time.Now(),
	}

	// Add event log to token's event history
	token.EventLogs = append(token.EventLogs, discountEvent)

	// Update the ledger with the discounting event
	err := ledger.RecordEvent(token.TokenID, "Dynamic Discount Applied", discountEvent)
	if err != nil {
		return errors.New("failed to log dynamic discount event in ledger")
	}

	return nil
}

// RecordSettlementEvent logs the settlement of a financial document and updates the ledger.
func RecordSettlementEvent(token *common.SYN2100Token, settledBy string, settlementAmount float64) error {
	// Create event log for settlement
	settlementEvent := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Document Settled",
		Description: "Settled for amount: " + formatAmount(settlementAmount),
		PerformedBy: settledBy,
		EventDate:   time.Now(),
	}

	// Add event log to token's event history
	token.EventLogs = append(token.EventLogs, settlementEvent)

	// Update the ledger with the settlement event
	err := ledger.RecordEvent(token.TokenID, "Document Settled", settlementEvent)
	if err != nil {
		return errors.New("failed to log settlement event in ledger")
	}

	return nil
}

// RecordVerificationEvent logs the verification of a financial document and updates the ledger.
func RecordVerificationEvent(token *common.SYN2100Token, verifier string, documentID string) error {
	// Encrypt document ID for security
	encryptedDocumentID, err := encryption.Encrypt([]byte(documentID))
	if err != nil {
		return errors.New("failed to encrypt document ID")
	}

	// Create event log for verification
	verificationEvent := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Document Verification",
		Description: "Document " + documentID + " has been verified.",
		PerformedBy: verifier,
		EventDate:   time.Now(),
	}

	// Add event log to token's event history
	token.EventLogs = append(token.EventLogs, verificationEvent)

	// Update the ledger to reflect the verification event
	err = ledger.RecordEvent(token.TokenID, "Document Verification", verificationEvent)
	if err != nil {
		return errors.New("failed to log verification event in ledger")
	}

	// Store encrypted document ID in the ledger
	err = ledger.StoreEncryptedDocumentID(token.TokenID, encryptedDocumentID)
	if err != nil {
		return errors.New("failed to store encrypted document ID in ledger")
	}

	return nil
}

// RecordLiquidityEvent logs liquidity provision for a tokenized invoice and updates the ledger.
func RecordLiquidityEvent(token *common.SYN2100Token, provider string, liquidityAmount float64) error {
	// Create event log for liquidity provision
	liquidityEvent := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Liquidity Provided",
		Description: "Liquidity of " + formatAmount(liquidityAmount) + " provided by " + provider,
		PerformedBy: provider,
		EventDate:   time.Now(),
	}

	// Add event log to token's event history
	token.EventLogs = append(token.EventLogs, liquidityEvent)

	// Update the ledger to reflect the liquidity provision event
	err := ledger.RecordEvent(token.TokenID, "Liquidity Provided", liquidityEvent)
	if err != nil {
		return errors.New("failed to log liquidity provision event in ledger")
	}

	return nil
}

// Utility function to generate a unique ID (can be replaced with a real implementation)
func generateUniqueID() string {
	// This is a placeholder for unique ID generation
	return "some-unique-id"
}

// Utility function to format amounts (for better readability)
func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
