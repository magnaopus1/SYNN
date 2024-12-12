package syn1700

import (
	"errors"
	"time"
)

// SYN1700Token represents an event ticket token under the SYN1700 standard.
type SYN1700Token struct {
	TokenID               string              // Unique identifier for the token
	Owner                 string              // The owner of the event ticket token
	EventMetadata         EventMetadata       // Metadata related to the event for which the ticket is issued
	TicketMetadata        TicketMetadata      // Metadata related to the specific ticket (class, price, etc.)
	OwnershipHistory      []OwnershipRecord   // Historical ownership data of the event ticket
	TicketTransactions    []TicketTransaction // Logs of transactions (sales, transfers, etc.)
	ComplianceStatus      string              // Compliance status of the token (e.g., "Compliant", "Pending Audit", etc.)
	RevocationStatus      bool                // Indicates whether access has been revoked by the organizer
	AccessRights          AccessRights        // Access rights for this event (e.g., entry, VIP areas, etc.)
	ImmutableRecords      []ImmutableRecord   // Immutable records for traceability and compliance
	TimeLock              TimeLock            // Time-lock information for ticket entry and use
	RestrictedTransfers   bool                // Indicates if transfers are restricted (e.g., anti-scalping)
	ApprovalRequired      bool                // Whether transfers or sales require approval from an event organizer or third party
	EncryptedMetadata     []byte              // Encrypted metadata containing sensitive information
}



// EventMetadata contains detailed information about the event related to the ticket.
type EventMetadata struct {
	EventID       string    // Unique identifier for the event
	EventName     string    // Name of the event
	EventLocation string    // Location of the event
	EventDate     time.Time // Date and time of the event
	StartTime     time.Time // Event start time
	EndTime       time.Time // Event end time
	TicketSupply  int       // Total number of tickets available for the event
}

// TicketMetadata contains detailed information about the specific ticket.
type TicketMetadata struct {
	TicketID     string    // Unique identifier for the ticket
	EventID      string    // Associated event ID
	TicketClass  string    // Class of the ticket (e.g., "Standard", "VIP", etc.)
	TicketType   string    // Type of the ticket (e.g., "Early-bird", "Standard", "Late-release")
	TicketPrice  float64   // Price of the ticket
	Conditions   string    // Special conditions for the ticket (e.g., "Adult", "Disabled")
	PurchaseDate time.Time // Date the ticket was purchased
}

// OwnershipRecord captures ownership changes for the event ticket.
type OwnershipRecord struct {
	PreviousOwner string    // ID of the previous owner
	NewOwner      string    // ID of the new owner
	TransferDate  time.Time // Date of ownership transfer
	TransferType  string    // Type of transfer (e.g., "Primary Sale", "Secondary Market")
}

// TicketTransaction logs the transactions involving the ticket.
type TicketTransaction struct {
	TransactionID     string    // Unique identifier for the transaction
	Sender            string    // ID of the ticket sender or seller
	Recipient         string    // ID of the recipient or buyer
	TransactionDate   time.Time // Date of the transaction
	TransactionAmount float64   // Amount paid or transferred
	TransactionType   string    // Type of transaction (e.g., "Purchase", "Transfer")
}

// AccessRights defines the access permissions granted by the event ticket.
type AccessRights struct {
	OwnerID           string    // ID of the current owner of the ticket
	AccessType        string    // Type of access granted (e.g., "General Admission", "VIP Area")
	DelegatedAccessID string    // If access is delegated, the ID of the person with delegated access
	AccessGranted     bool      // Whether access is currently granted
	AccessStartTime   time.Time // Time when access begins
	AccessEndTime     time.Time // Time when access ends
}

// TimeLock defines the time-locked properties of the ticket, such as when it can be used.
type TimeLock struct {
	LockStartTime time.Time // Time when the time lock begins
	LockEndTime   time.Time // Time when the time lock ends
}

// ImmutableRecord stores immutable records for compliance and transparency.
type ImmutableRecord struct {
	RecordID    string    // Unique identifier for the record
	Description string    // Description of the record (e.g., "Ticket Issued", "Transfer")
	Timestamp   time.Time // When the record was created
}

// Factory to create a new SYN1700Token
type TokenFactory struct {
	ledgerInstance *ledger.Ledger
}

// NewTokenFactory creates a new instance of TokenFactory
func NewTokenFactory(ledger *ledger.Ledger) *TokenFactory {
	return &TokenFactory{ledgerInstance: ledger}
}

// CreateSYN1700Token creates a new SYN1700 event ticket token and stores it in the ledger
func (factory *TokenFactory) CreateSYN1700Token(owner string, eventMeta EventMetadata, ticketMeta TicketMetadata, restrictedTransfers bool, approvalRequired bool) (*SYN1700Token, error) {
	tokenID := common.GenerateUniqueID()

	// Create encrypted metadata for sensitive data
	encryptedMetadata, err := encryption.EncryptMetadata(ticketMeta)
	if err != nil {
		return nil, err
	}

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

	// Save token to ledger
	err = factory.ledgerInstance.StoreToken(newToken)
	if err != nil {
		return nil, err
	}

	// Log event in ledger
	factory.ledgerInstance.LogEvent("TokenCreated", newToken.TokenID, time.Now())

	return newToken, nil
}

// TransferOwnership transfers the ownership of a ticket to a new owner
func (token *SYN1700Token) TransferOwnership(newOwner string, transferType string, approvalRequired bool) error {
	if token.RestrictedTransfers && approvalRequired {
		return errors.New("Transfer restricted or approval required")
	}

	ownershipRecord := OwnershipRecord{
		PreviousOwner: token.Owner,
		NewOwner:      newOwner,
		TransferDate:  time.Now(),
		TransferType:  transferType,
	}

	token.OwnershipHistory = append(token.OwnershipHistory, ownershipRecord)
	token.Owner = newOwner

	// Log the transfer event in the ledger
	ledger.LogEvent("OwnershipTransferred", token.TokenID, time.Now())

	return nil
}

// VerifyAccess verifies if the token owner can access the event
func (token *SYN1700Token) VerifyAccess() (bool, error) {
	currentTime := time.Now()
	if token.TimeLock.LockStartTime.Before(currentTime) && token.TimeLock.LockEndTime.After(currentTime) {
		return true, nil
	}
	return false, errors.New("Access denied: ticket is not valid for entry at this time")
}

// RevokeAccess revokes access rights for the token
func (token *SYN1700Token) RevokeAccess() {
	token.RevocationStatus = true

	// Log revocation in ledger
	ledger.LogEvent("AccessRevoked", token.TokenID, time.Now())
}

// LogTransaction logs a transaction related to the ticket
func (token *SYN1700Token) LogTransaction(sender string, recipient string, amount float64, transactionType string) {
	transaction := TicketTransaction{
		TransactionID:     common.GenerateUniqueID(),
		Sender:            sender,
		Recipient:         recipient,
		TransactionDate:   time.Now(),
		TransactionAmount: amount,
		TransactionType:   transactionType,
	}

	token.TicketTransactions = append(token.TicketTransactions, transaction)

	// Log transaction in ledger
	ledger.LogEvent("TicketTransaction", transaction.TransactionID, time.Now())
}

// AddImmutableRecord adds an immutable record to the token's history
func (token *SYN1700Token) AddImmutableRecord(description string) {
	record := ImmutableRecord{
		RecordID:    common.GenerateUniqueID(),
		Description: description,
		Timestamp:   time.Now(),
	}

	token.ImmutableRecords = append(token.ImmutableRecords, record)

	// Log record addition in ledger
	ledger.LogEvent("ImmutableRecordAdded", record.RecordID, time.Now())
}
