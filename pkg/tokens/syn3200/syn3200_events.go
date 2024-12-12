package syn3200

import (
	"time"
	"errors"
	"sync"

)

// BillPaymentEvent represents an event where a payment is made on a bill.
type BillPaymentEvent struct {
	EventID     string    `json:"event_id"`
	BillID      string    `json:"bill_id"`
	Payer       string    `json:"payer"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	Status      string    `json:"status"`   // Completed, Failed
	Timestamp   time.Time `json:"timestamp"`
}

// AutomatedPaymentEvent represents an event where an automated payment is executed.
type AutomatedPaymentEvent struct {
	EventID     string    `json:"event_id"`
	PaymentID   string    `json:"payment_id"`
	BillID      string    `json:"bill_id"`
	Payer       string    `json:"payer"`
	Amount      float64   `json:"amount"`
	Schedule    string    `json:"schedule"` // Daily, Weekly, Monthly
	Status      string    `json:"status"`   // Pending, Completed, Cancelled
	ExecutedAt  time.Time `json:"executed_at"`
	Timestamp   time.Time `json:"timestamp"`
}

// ConditionalEnforcementEvent represents an event when a conditional bill enforcement is executed.
type ConditionalEnforcementEvent struct {
	EventID        string    `json:"event_id"`
	EnforcementID  string    `json:"enforcement_id"`
	BillID         string    `json:"bill_id"`
	Condition      string    `json:"condition"`
	Enforced       bool      `json:"enforced"`
	EnforcementDate time.Time `json:"enforcement_date"`
	Timestamp      time.Time `json:"timestamp"`
}

// FairAllocationEvent represents an event when a fair allocation is performed.
type FairAllocationEvent struct {
	EventID       string    `json:"event_id"`
	AllocationID  string    `json:"allocation_id"`
	BillID        string    `json:"bill_id"`
	AllocatedTo   string    `json:"allocated_to"`
	AllocatedAmount float64 `json:"allocated_amount"`
	AllocationDate time.Time `json:"allocation_date"`
	Timestamp     time.Time `json:"timestamp"`
}

// EventManager manages the creation and tracking of all events for the SYN3200 standard.
type EventManager struct {
	ledgerService     *ledger.Ledger
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewEventManager creates a new instance of EventManager.
func NewEventManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *EventManager {
	return &EventManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// LogBillPayment logs a bill payment event.
func (em *EventManager) LogBillPayment(billID, payer string, amount float64, status string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	event := &BillPaymentEvent{
		EventID:     billID + "-" + payer + "-payment-event",
		BillID:      billID,
		Payer:       payer,
		Amount:      amount,
		PaymentDate: time.Now(),
		Status:      status,
		Timestamp:   time.Now(),
	}

	// Encrypt the event.
	encryptedEvent, err := em.encryptionService.EncryptData(event)
	if err != nil {
		return err
	}

	// Log the event in the ledger.
	err = em.ledgerService.LogEvent("BillPaymentEvent", time.Now(), event.EventID)
	if err != nil {
		return err
	}

	// Validate the event via consensus.
	err = em.consensusService.ValidateSubBlock(event.EventID)
	if err != nil {
		return err
	}

	// Store the event.
	err = em.ledgerService.StoreEvent(event.EventID, encryptedEvent)
	if err != nil {
		return err
	}

	return nil
}

// LogAutomatedPayment logs an automated payment event.
func (em *EventManager) LogAutomatedPayment(paymentID, billID, payer string, amount float64, schedule, status string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	event := &AutomatedPaymentEvent{
		EventID:    billID + "-" + paymentID + "-automated-payment-event",
		PaymentID:  paymentID,
		BillID:     billID,
		Payer:      payer,
		Amount:     amount,
		Schedule:   schedule,
		Status:     status,
		ExecutedAt: time.Now(),
		Timestamp:  time.Now(),
	}

	// Encrypt the event.
	encryptedEvent, err := em.encryptionService.EncryptData(event)
	if err != nil {
		return err
	}

	// Log the event in the ledger.
	err = em.ledgerService.LogEvent("AutomatedPaymentEvent", time.Now(), event.EventID)
	if err != nil {
		return err
	}

	// Validate the event via consensus.
	err = em.consensusService.ValidateSubBlock(event.EventID)
	if err != nil {
		return err
	}

	// Store the event.
	err = em.ledgerService.StoreEvent(event.EventID, encryptedEvent)
	if err != nil {
		return err
	}

	return nil
}

// LogConditionalEnforcement logs a conditional enforcement event.
func (em *EventManager) LogConditionalEnforcement(enforcementID, billID, condition string, enforced bool) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	event := &ConditionalEnforcementEvent{
		EventID:        enforcementID + "-conditional-enforcement-event",
		EnforcementID:  enforcementID,
		BillID:         billID,
		Condition:      condition,
		Enforced:       enforced,
		EnforcementDate: time.Now(),
		Timestamp:      time.Now(),
	}

	// Encrypt the event.
	encryptedEvent, err := em.encryptionService.EncryptData(event)
	if err != nil {
		return err
	}

	// Log the event in the ledger.
	err = em.ledgerService.LogEvent("ConditionalEnforcementEvent", time.Now(), event.EventID)
	if err != nil {
		return err
	}

	// Validate the event via consensus.
	err = em.consensusService.ValidateSubBlock(event.EventID)
	if err != nil {
		return err
	}

	// Store the event.
	err = em.ledgerService.StoreEvent(event.EventID, encryptedEvent)
	if err != nil {
		return err
	}

	return nil
}

// LogFairAllocation logs a fair allocation event.
func (em *EventManager) LogFairAllocation(allocationID, billID, allocatedTo string, allocatedAmount float64) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	event := &FairAllocationEvent{
		EventID:         allocationID + "-fair-allocation-event",
		AllocationID:    allocationID,
		BillID:          billID,
		AllocatedTo:     allocatedTo,
		AllocatedAmount: allocatedAmount,
		AllocationDate:  time.Now(),
		Timestamp:       time.Now(),
	}

	// Encrypt the event.
	encryptedEvent, err := em.encryptionService.EncryptData(event)
	if err != nil {
		return err
	}

	// Log the event in the ledger.
	err = em.ledgerService.LogEvent("FairAllocationEvent", time.Now(), event.EventID)
	if err != nil {
		return err
	}

	// Validate the event via consensus.
	err = em.consensusService.ValidateSubBlock(event.EventID)
	if err != nil {
		return err
	}

	// Store the event.
	err = em.ledgerService.StoreEvent(event.EventID, encryptedEvent)
	if err != nil {
		return err
	}

	return nil
}
