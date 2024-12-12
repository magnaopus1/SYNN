package community_engagement

import (
	"errors"
	"fmt"
	"synnergy_network/pkg/ledger"
	"time"
)

// Event struct to store information about an event in the blockchain network
type Event struct {
	ID              string    // Unique identifier for the event
	Name            string    // Name of the event
	Description     string    // Description of the event
	Location        string    // Location of the event
	Date            time.Time // Date and time of the event
	CreatorID       string    // User ID of the event creator
	MaxParticipants int       // Maximum number of participants allowed
	Participants    []string  // List of participant IDs
	CreatedAt       time.Time // Timestamp when the event was created
}

// createEvent allows users to create an event within the blockchain network
func CreateEvent(creatorID, eventName, description, location string, eventDate time.Time, maxParticipants int) (string, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Generate a unique ID for the event
	eventID := generateUniqueID(creatorID + eventName)

	// Create the event struct
	event := ledger.Event{
		ID:              eventID,
		Name:            eventName,
		Description:     description,
		Location:        location,
		Date:            eventDate,
		CreatorID:       creatorID,
		MaxParticipants: maxParticipants,
		Participants:    []string{},
		CreatedAt:       time.Now(),
	}

	// Record the event in the ledger
	if err := l.CommunityEngagementLedger.RecordEvent(event); err != nil {
		return "", fmt.Errorf("failed to record event: %v", err)
	}

	return eventID, nil
}

// joinEvent registers a user to participate in an event, ensuring participant limit compliance
func JoinEvent(userID, eventID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the event
	event, err := l.CommunityEngagementLedger.FetchEvent(eventID)
	if err != nil {
		return fmt.Errorf("failed to retrieve event: %v", err)
	}

	// Check if the event has reached its maximum number of participants
	if len(event.Participants) >= event.MaxParticipants {
		return errors.New("event is full")
	}

	// Check if the user is already a participant
	for _, participant := range event.Participants {
		if participant == userID {
			return errors.New("user is already registered for this event")
		}
	}

	// Add the user to the event's participant list and update the ledger
	event.Participants = append(event.Participants, userID)
	if err := l.CommunityEngagementLedger.UpdateEvent(event); err != nil {
		return fmt.Errorf("failed to update event participants: %v", err)
	}

	return nil
}

// leaveEvent removes a user from an event's participant list
func LeaveEvent(userID, eventID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the event
	event, err := l.CommunityEngagementLedger.FetchEvent(eventID)
	if err != nil {
		return fmt.Errorf("failed to retrieve event: %v", err)
	}

	// Check if the user is registered as a participant
	participantFound := false
	for i, participant := range event.Participants {
		if participant == userID {
			// Remove the user from the participants list
			event.Participants = append(event.Participants[:i], event.Participants[i+1:]...)
			participantFound = true
			break
		}
	}

	if !participantFound {
		return errors.New("user is not registered for this event")
	}

	// Update the event in the ledger
	if err := l.CommunityEngagementLedger.UpdateEvent(event); err != nil {
		return fmt.Errorf("failed to update event participants: %v", err)
	}

	return nil
}


// listAllEvents retrieves all events created within the network
func ListAllEvents() ([]ledger.Event, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve all events
	events, err := l.CommunityEngagementLedger.ListAllEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to list all events: %v", err)
	}
	return events, nil
}

// viewEventDetails provides detailed information about a specific event
func ViewEventDetails(eventID string) (ledger.Event, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the specific event by event ID
	event, err := l.CommunityEngagementLedger.FetchEvent(eventID)
	if err != nil {
		return ledger.Event{}, fmt.Errorf("failed to retrieve event details: %v", err)
	}
	return event, nil
}