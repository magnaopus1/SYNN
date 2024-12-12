package syn4200

import (
	"errors"
	"time"
	"sync"
)

// Syn4200EventManager manages events related to SYN4200 tokens, such as donations, fundraising milestones, and social impact tracking.
type Syn4200EventManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewSyn4200EventManager creates a new event manager for SYN4200 tokens.
func NewSyn4200EventManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *Syn4200EventManager {
	return &Syn4200EventManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// LogDonationEvent logs a donation event to the ledger, including the donor, amount, and campaign.
func (em *Syn4200EventManager) LogDonationEvent(tokenID, donor string, amount float64, campaign string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	eventDetails := "Donation of " + fmt.Sprintf("%.2f", amount) + " made by " + donor + " to campaign " + campaign
	eventTimestamp := time.Now()

	// Log the event in the ledger
	if err := em.ledgerService.LogEvent("DonationEvent", eventTimestamp, tokenID, eventDetails); err != nil {
		return err
	}

	// Encrypt and store the event details securely
	encryptedEventDetails, err := em.encryptionService.EncryptData(eventDetails)
	if err != nil {
		return err
	}

	// Validate the donation event using Synnergy Consensus
	if err := em.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	// Store encrypted event details in the ledger
	return em.ledgerService.StoreEvent(tokenID, encryptedEventDetails)
}

// LogMilestoneEvent logs a fundraising milestone event, marking progress towards a goal.
func (em *Syn4200EventManager) LogMilestoneEvent(tokenID, milestoneDescription string, progress float64) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	eventDetails := "Milestone reached: " + milestoneDescription + " | Progress: " + fmt.Sprintf("%.2f%%", progress*100)
	eventTimestamp := time.Now()

	// Log the event in the ledger
	if err := em.ledgerService.LogEvent("MilestoneEvent", eventTimestamp, tokenID, eventDetails); err != nil {
		return err
	}

	// Encrypt and store the event details securely
	encryptedEventDetails, err := em.encryptionService.EncryptData(eventDetails)
	if err != nil {
		return err
	}

	// Validate the milestone event using Synnergy Consensus
	if err := em.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	// Store encrypted event details in the ledger
	return em.ledgerService.StoreEvent(tokenID, encryptedEventDetails)
}

// LogImpactEvent logs an event tracking the social impact achieved through donations, including details about the funded projects and outcomes.
func (em *Syn4200EventManager) LogImpactEvent(tokenID, impactDetails string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	eventTimestamp := time.Now()

	// Log the impact event in the ledger
	if err := em.ledgerService.LogEvent("ImpactEvent", eventTimestamp, tokenID, impactDetails); err != nil {
		return err
	}

	// Encrypt and store the event details securely
	encryptedImpactDetails, err := em.encryptionService.EncryptData(impactDetails)
	if err != nil {
		return err
	}

	// Validate the impact event using Synnergy Consensus
	if err := em.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	// Store encrypted event details in the ledger
	return em.ledgerService.StoreEvent(tokenID, encryptedImpactDetails)
}

// GetEventHistory retrieves the event history for a specific token, decrypting each event.
func (em *Syn4200EventManager) GetEventHistory(tokenID string) ([]string, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Retrieve encrypted event history from the ledger
	encryptedEvents, err := em.ledgerService.RetrieveEvents(tokenID)
	if err != nil {
		return nil, err
	}

	// Decrypt the events
	var eventHistory []string
	for _, encryptedEvent := range encryptedEvents {
		decryptedEvent, err := em.encryptionService.DecryptData(encryptedEvent)
		if err != nil {
			return nil, err
		}
		eventHistory = append(eventHistory, decryptedEvent.(string))
	}

	return eventHistory, nil
}

// TrackFundraisingProgress triggers an event when fundraising progress hits a certain threshold.
func (em *Syn4200EventManager) TrackFundraisingProgress(tokenID string, goalAmount, raisedAmount float64) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	percentageAchieved := (raisedAmount / goalAmount) * 100
	if percentageAchieved >= 100 {
		return em.LogMilestoneEvent(tokenID, "Fundraising Goal Reached", 1.0)
	}

	// If not 100% yet, log the progress as a milestone event.
	return em.LogMilestoneEvent(tokenID, "Fundraising In Progress", percentageAchieved/100)
}

// GenerateUniqueEventID generates a unique identifier for events.
func generateUniqueEventID() string {
	// Implement a real-world unique event ID generation logic (e.g., UUID or timestamp).
	return "event-id-" + time.Now().Format("20060102150405")
}
