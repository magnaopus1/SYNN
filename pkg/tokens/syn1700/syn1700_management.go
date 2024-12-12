package syn1700

import (
	"errors"
	"time"
)

// SYN1700Management is responsible for managing SYN1700 tokens.
type SYN1700Management struct {
	ledgerInstance *ledger.Ledger
}

// NewSYN1700Management creates a new instance for managing SYN1700 tokens.
func NewSYN1700Management(ledger *ledger.Ledger) *SYN1700Management {
	return &SYN1700Management{
		ledgerInstance: ledger,
	}
}

// CreateTicket creates a new event ticket under the SYN1700 standard and stores it in the ledger.
func (m *SYN1700Management) CreateTicket(owner string, eventMeta EventMetadata, ticketMeta TicketMetadata, restrictedTransfers bool, approvalRequired bool) (*SYN1700Token, error) {
	tokenID := common.GenerateUniqueID()

	// Encrypt sensitive metadata
	encryptedMetadata, err := encryption.EncryptMetadata(ticketMeta)
	if err != nil {
		return nil, err
	}

	// Create new token
	newToken := &SYN1700Token{
		TokenID:           tokenID,
		Owner:             owner,
		EventMetadata:     eventMeta,
		TicketMetadata:    ticketMeta,
		OwnershipHistory:  []OwnershipRecord{},
		TicketTransactions: []TicketTransaction{},
		ComplianceStatus:  "Compliant",
		RevocationStatus:  false,
		AccessRights:      AccessRights{},
		ImmutableRecords:  []ImmutableRecord{},
		TimeLock:          TimeLock{},
		RestrictedTransfers: restrictedTransfers,
		ApprovalRequired:    approvalRequired,
		EncryptedMetadata:   encryptedMetadata,
	}

	// Validate and process the token into sub-blocks using Synnergy Consensus
	if err := m.validateAndProcessTicket(newToken); err != nil {
		return nil, err
	}

	// Store token into the ledger
	err = m.ledgerInstance.StoreToken(newToken)
	if err != nil {
		return nil, err
	}

	// Log the creation of the new token in the ledger
	m.ledgerInstance.LogEvent("TicketCreated", newToken.TokenID, time.Now())

	return newToken, nil
}

// TransferTicket transfers ownership of a ticket to another user and logs it.
func (m *SYN1700Management) TransferTicket(token *SYN1700Token, newOwner string, transferType string, approvalRequired bool) error {
	if token.RestrictedTransfers && approvalRequired {
		return errors.New("transfer restricted or approval required")
	}

	// Ownership record for the transfer
	ownershipRecord := OwnershipRecord{
		PreviousOwner: token.Owner,
		NewOwner:      newOwner,
		TransferDate:  time.Now(),
		TransferType:  transferType,
	}

	token.OwnershipHistory = append(token.OwnershipHistory, ownershipRecord)
	token.Owner = newOwner

	// Log the transfer in the ledger
	m.ledgerInstance.LogEvent("TicketTransferred", token.TokenID, time.Now())

	return nil
}

// RevokeTicket revokes access rights to a ticket.
func (m *SYN1700Management) RevokeTicket(token *SYN1700Token) error {
	token.RevocationStatus = true

	// Log revocation in the ledger
	m.ledgerInstance.LogEvent("TicketRevoked", token.TokenID, time.Now())

	return nil
}

// LogTicketTransaction logs a ticket-related transaction (e.g., sale or transfer).
func (m *SYN1700Management) LogTicketTransaction(token *SYN1700Token, sender string, recipient string, amount float64, transactionType string) error {
	transaction := TicketTransaction{
		TransactionID:     common.GenerateUniqueID(),
		Sender:            sender,
		Recipient:         recipient,
		TransactionDate:   time.Now(),
		TransactionAmount: amount,
		TransactionType:   transactionType,
	}

	token.TicketTransactions = append(token.TicketTransactions, transaction)

	// Log transaction in the ledger
	m.ledgerInstance.LogEvent("TicketTransactionLogged", transaction.TransactionID, time.Now())

	return nil
}

// ValidateTicket checks if the ticket is valid for the event and can be used for entry.
func (m *SYN1700Management) ValidateTicket(token *SYN1700Token) (bool, error) {
	currentTime := time.Now()
	if token.RevocationStatus {
		return false, errors.New("ticket access revoked")
	}

	if token.TimeLock.LockStartTime.Before(currentTime) && token.TimeLock.LockEndTime.After(currentTime) {
		return true, nil
	}

	return false, errors.New("ticket is not valid for entry at this time")
}

// AddImmutableRecord adds an immutable record to the ticket for traceability.
func (m *SYN1700Management) AddImmutableRecord(token *SYN1700Token, description string) error {
	record := ImmutableRecord{
		RecordID:    common.GenerateUniqueID(),
		Description: description,
		Timestamp:   time.Now(),
	}

	token.ImmutableRecords = append(token.ImmutableRecords, record)

	// Log the immutable record addition in the ledger
	m.ledgerInstance.LogEvent("ImmutableRecordAdded", record.RecordID, time.Now())

	return nil
}

// validateAndProcessTicket validates the ticket through Synnergy Consensus and processes it into sub-blocks.
func (m *SYN1700Management) validateAndProcessTicket(token *SYN1700Token) error {
	// Synnergy Consensus: Validate transactions into sub-blocks (1000 sub-blocks per block)
	subBlocks := common.GenerateSubBlocks(token.TokenID, 1000)

	// Process validation logic (e.g., running consensus checks on each sub-block)
	for _, subBlock := range subBlocks {
		if err := common.ValidateSubBlock(subBlock); err != nil {
			return err
		}
	}

	// Store each validated sub-block into the ledger
	for _, subBlock := range subBlocks {
		if err := m.ledgerInstance.StoreSubBlock(subBlock); err != nil {
			return err
		}
	}

	// Log the validation in the ledger
	m.ledgerInstance.LogEvent("TicketValidated", token.TokenID, time.Now())

	return nil
}
