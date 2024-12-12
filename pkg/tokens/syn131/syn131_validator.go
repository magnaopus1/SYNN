package syn131

import (
	"errors"

)

// ValidateSYN131ComplianceRecord validates compliance records for SYN131 tokens
func ValidateSYN131ComplianceRecord(record *ComplianceRecord) error {
	if record == nil {
		return errors.New("compliance record is nil")
	}
	if record.TokenID == "" {
		return errors.New("missing token ID in compliance record")
	}
	if record.Status == "" {
		return errors.New("missing status in compliance record")
	}
	if record.ComplianceDate.IsZero() {
		return errors.New("missing compliance date in compliance record")
	}
	if record.Fee < 0 {
		return errors.New("invalid fee in compliance record")
	}
	return nil
}

// ValidateSYN131Event validates events in the SYN131 system
func ValidateSYN131Event(event *Event) error {
	if event == nil {
		return errors.New("event is nil")
	}
	if event.ID == "" {
		return errors.New("missing event ID")
	}
	if event.Type == "" {
		return errors.New("missing event type")
	}
	if event.Timestamp.IsZero() {
		return errors.New("missing timestamp in event")
	}
	if event.Fee < 0 {
		return errors.New("invalid fee in event")
	}
	return nil
}

// ValidateSYN131TokenCreation validates token creation transactions
func ValidateSYN131TokenCreation(token *Syn131Token) error {
	if token == nil {
		return errors.New("token is nil")
	}
	if token.ID == "" {
		return errors.New("missing token ID")
	}
	if token.Owner == "" {
		return errors.New("missing token owner")
	}
	if token.Status == "" {
		return errors.New("missing token status")
	}
	if token.CreatedAt.IsZero() {
		return errors.New("missing creation time")
	}
	return nil
}

// ValidateSYN131TokenUpdate validates token update transactions
func ValidateSYN131TokenUpdate(token *Syn131Token) error {
	if token == nil {
		return errors.New("token is nil")
	}
	if token.ID == "" {
		return errors.New("missing token ID")
	}
	if token.UpdatedAt.IsZero() {
		return errors.New("missing update time")
	}
	return nil
}

// ValidateSYN131Transaction validates a generic SYN131 transaction
func ValidateSYN131Transaction(tx *SYN131Transaction) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}
	if tx.TransactionID == "" {
		return errors.New("missing transaction ID")
	}
	if tx.TokenID == "" {
		return errors.New("missing token ID")
	}
	if tx.Status == "" {
		return errors.New("missing transaction status")
	}
	if tx.Timestamp.IsZero() {
		return errors.New("missing transaction timestamp")
	}
	if tx.Fee < 0 {
		return errors.New("invalid transaction fee")
	}
	return nil
}

// ValidateSYN131TokenTransfer validates token transfer transactions
func ValidateSYN131TokenTransfer(tx *OwnershipTransaction) error {
	if tx == nil {
		return errors.New("ownership transaction is nil")
	}
	if tx.TransactionID == "" {
		return errors.New("missing transaction ID")
	}
	if tx.AssetID == "" {
		return errors.New("missing asset ID")
	}
	if tx.FromOwner == "" || tx.ToOwner == "" {
		return errors.New("missing ownership details")
	}
	if tx.Amount <= 0 {
		return errors.New("invalid transfer amount")
	}
	if tx.Fee < 0 {
		return errors.New("invalid transaction fee")
	}
	return nil
}

// ValidateShardedOwnershipTransaction validates sharded ownership transactions
func ValidateShardedOwnershipTransaction(tx *ShardedOwnershipTransaction) error {
	if tx == nil {
		return errors.New("sharded ownership transaction is nil")
	}
	if tx.TransactionID == "" {
		return errors.New("missing transaction ID")
	}
	if tx.AssetID == "" {
		return errors.New("missing asset ID")
	}
	if len(tx.FromOwners) == 0 || len(tx.ToOwners) == 0 {
		return errors.New("missing ownership details")
	}
	if tx.TotalAmount <= 0 {
		return errors.New("invalid total amount")
	}
	if tx.Fee < 0 {
		return errors.New("invalid transaction fee")
	}
	return nil
}

// ValidateOwnershipTransaction validates standard ownership transactions
func ValidateOwnershipTransaction(tx *OwnershipTransaction) error {
	return ValidateSYN131TokenTransfer(tx)
}

// ValidateRentalTransaction validates rental payment transactions
func ValidateRentalTransaction(tx *RentalTransaction) error {
	if tx == nil {
		return errors.New("rental transaction is nil")
	}
	if tx.TransactionID == "" {
		return errors.New("missing transaction ID")
	}
	if tx.RentalAgreementID == "" {
		return errors.New("missing rental agreement ID")
	}
	if tx.Amount <= 0 {
		return errors.New("invalid rental amount")
	}
	if tx.Fee < 0 {
		return errors.New("invalid transaction fee")
	}
	return nil
}

// ValidateLeaseTransaction validates lease payment transactions
func ValidateLeaseTransaction(tx *LeaseTransaction) error {
	if tx == nil {
		return errors.New("lease transaction is nil")
	}
	if tx.TransactionID == "" {
		return errors.New("missing transaction ID")
	}
	if tx.LeaseAgreementID == "" {
		return errors.New("missing lease agreement ID")
	}
	if tx.Amount <= 0 {
		return errors.New("invalid lease amount")
	}
	if tx.Fee < 0 {
		return errors.New("invalid transaction fee")
	}
	return nil
}

// ValidatePurchaseTransaction validates purchase transactions
func ValidatePurchaseTransaction(tx *PurchaseTransaction) error {
	if tx == nil {
		return errors.New("purchase transaction is nil")
	}
	if tx.TransactionID == "" {
		return errors.New("missing transaction ID")
	}
	if tx.AssetID == "" {
		return errors.New("missing asset ID")
	}
	if tx.Amount <= 0 {
		return errors.New("invalid purchase amount")
	}
	if tx.Fee < 0 {
		return errors.New("invalid transaction fee")
	}
	return nil
}
