package syn1700

import (
	"errors"
	"time"
)

// SYN1700TransactionManager handles all transaction operations related to SYN1700 tokens (event tickets).
type SYN1700TransactionManager struct {
	ledgerInstance *ledger.Ledger
}

// NewSYN1700TransactionManager initializes a new SYN1700TransactionManager with ledger integration.
func NewSYN1700TransactionManager(ledger *ledger.Ledger) *SYN1700TransactionManager {
	return &SYN1700TransactionManager{
		ledgerInstance: ledger,
	}
}

// TransferTicket facilitates the transfer of ownership for a SYN1700 ticket from one user to another.
func (tm *SYN1700TransactionManager) TransferTicket(tokenID, senderID, recipientID string, amount float64) error {
	// Retrieve the token from the ledger
	token, err := tm.ledgerInstance.GetToken(tokenID)
	if err != nil {
		return errors.New("token retrieval failed: " + err.Error())
	}

	// Verify sender's ownership
	if token.Owner != senderID {
		return errors.New("sender does not own this ticket")
	}

	// Validate transaction with Synnergy Consensus sub-block validation
	err = tm.validateSubBlockTransaction(tokenID, senderID, recipientID)
	if err != nil {
		return err
	}

	// Update ownership and log the transfer
	newOwnershipRecord := OwnershipRecord{
		PreviousOwner: senderID,
		NewOwner:      recipientID,
		TransferDate:  time.Now(),
		TransferType:  "Transfer",
	}
	token.OwnershipHistory = append(token.OwnershipHistory, newOwnershipRecord)
	token.Owner = recipientID

	// Log transaction details
	transactionLog := TicketTransaction{
		TransactionID:     common.GenerateUniqueID(),
		Sender:            senderID,
		Recipient:         recipientID,
		TransactionDate:   time.Now(),
		TransactionAmount: amount,
		TransactionType:   "Transfer",
	}
	token.TicketTransactions = append(token.TicketTransactions, transactionLog)

	// Encrypt sensitive metadata before updating
	encryptedMetadata, err := encryption.EncryptData(token.EncryptedMetadata)
	if err != nil {
		return err
	}
	token.EncryptedMetadata = encryptedMetadata

	// Update the token in the ledger
	err = tm.ledgerInstance.UpdateToken(tokenID, token)
	if err != nil {
		return err
	}

	// Log the transaction in the ledger
	tm.ledgerInstance.LogEvent("TicketTransferred", tokenID, time.Now())

	return nil
}

// ValidateTicket ensures that the ticket is valid for entry at the event.
func (tm *SYN1700TransactionManager) ValidateTicket(tokenID, ownerID string) error {
	// Retrieve the token from the ledger
	token, err := tm.ledgerInstance.GetToken(tokenID)
	if err != nil {
		return errors.New("token retrieval failed: " + err.Error())
	}

	// Ensure the owner matches the ticket holder
	if token.Owner != ownerID {
		return errors.New("the ticket does not belong to this owner")
	}

	// Check the revocation status
	if token.RevocationStatus {
		return errors.New("ticket access has been revoked")
	}

	// Validate time lock (if ticket is time-locked)
	if time.Now().Before(token.TimeLock.LockStartTime) || time.Now().After(token.TimeLock.LockEndTime) {
		return errors.New("ticket cannot be used outside of the designated time")
	}

	// Log the validation event
	tm.ledgerInstance.LogEvent("TicketValidated", tokenID, time.Now())

	return nil
}

// IssueTicket handles issuing a new SYN1700 token and storing it in the ledger.
func (tm *SYN1700TransactionManager) IssueTicket(eventMetadata EventMetadata, ticketMetadata TicketMetadata, ownerID string) (string, error) {
	tokenID := common.GenerateUniqueID()

	// Create a new SYN1700Token
	newToken := SYN1700Token{
		TokenID:          tokenID,
		Owner:            ownerID,
		EventMetadata:    eventMetadata,
		TicketMetadata:   ticketMetadata,
		OwnershipHistory: []OwnershipRecord{},
		TicketTransactions: []TicketTransaction{
			{
				TransactionID:     common.GenerateUniqueID(),
				Sender:            "EventOrganizer",
				Recipient:         ownerID,
				TransactionDate:   time.Now(),
				TransactionAmount: ticketMetadata.TicketPrice,
				TransactionType:   "Issue",
			},
		},
		ComplianceStatus:    "Compliant",
		RevocationStatus:    false,
		AccessRights:        AccessRights{OwnerID: ownerID, AccessType: "General Admission"},
		ImmutableRecords:    []ImmutableRecord{},
		TimeLock:            TimeLock{LockStartTime: eventMetadata.EventDate.Add(-time.Hour), LockEndTime: eventMetadata.EndTime},
		RestrictedTransfers: false,
		ApprovalRequired:    false,
	}

	// Encrypt sensitive metadata before storing
	encryptedMetadata, err := encryption.EncryptData(newToken.EncryptedMetadata)
	if err != nil {
		return "", err
	}
	newToken.EncryptedMetadata = encryptedMetadata

	// Store the token in the ledger
	err = tm.ledgerInstance.StoreToken(newToken.TokenID, newToken)
	if err != nil {
		return "", err
	}

	// Log the issuance event in the ledger
	tm.ledgerInstance.LogEvent("TicketIssued", newToken.TokenID, time.Now())

	return newToken.TokenID, nil
}

// RevokeTicket revokes a ticket, preventing further use for event access.
func (tm *SYN1700TransactionManager) RevokeTicket(tokenID, organizerID string) error {
	// Retrieve the token from the ledger
	token, err := tm.ledgerInstance.GetToken(tokenID)
	if err != nil {
		return errors.New("token retrieval failed: " + err.Error())
	}

	// Verify that the event organizer is authorized to revoke the ticket
	if organizerID != token.EventMetadata.EventID {
		return errors.New("unauthorized to revoke this ticket")
	}

	// Mark the token's access as revoked
	token.RevocationStatus = true

	// Update the token in the ledger
	err = tm.ledgerInstance.UpdateToken(tokenID, token)
	if err != nil {
		return err
	}

	// Log the revocation event in the ledger
	tm.ledgerInstance.LogEvent("TicketRevoked", tokenID, time.Now())

	return nil
}

// validateSubBlockTransaction validates the transaction using Synnergy Consensus by validating sub-blocks before the transfer is confirmed.
func (tm *SYN1700TransactionManager) validateSubBlockTransaction(tokenID, senderID, recipientID string) error {
	// Generate sub-blocks for the transaction
	subBlocks := common.GenerateSubBlocks(tokenID, 1000)

	// Validate each sub-block using Synnergy Consensus
	for _, subBlock := range subBlocks {
		err := common.ValidateSubBlock(subBlock)
		if err != nil {
			return errors.New("sub-block validation failed: " + err.Error())
		}
	}

	// Store the validated sub-blocks in the ledger
	for _, subBlock := range subBlocks {
		err := tm.ledgerInstance.StoreSubBlock(subBlock)
		if err != nil {
			return errors.New("sub-block storage failed: " + err.Error())
		}
	}

	return nil
}

// GetTransactionHistory retrieves the transaction history of a SYN1700 token.
func (tm *SYN1700TransactionManager) GetTransactionHistory(tokenID string) ([]TicketTransaction, error) {
	// Retrieve the token from the ledger
	token, err := tm.ledgerInstance.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("token retrieval failed: " + err.Error())
	}

	// Return the transaction history
	return token.TicketTransactions, nil
}

