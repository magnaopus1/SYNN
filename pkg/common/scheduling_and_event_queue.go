package common

import (
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	
	"time"
)

// ScheduledEvent represents a task scheduled for future execution.
type ScheduledEvent struct {
	EventID       string      // Unique identifier for the event
	ExecutionTime time.Time   // Time at which the event should be executed
	Opcode        byte        // Opcode associated with the event
	Operands      []int       // Operands for the opcode
	Status        string      // Current status (e.g., "scheduled", "executed", "failed")
}

// EventQueue manages a queue of scheduled events with synchronization support.
type EventQueue struct {
	events        []ScheduledEvent
	ledgerInstance *ledger.Ledger
	mutex         sync.Mutex
}

// NewEventQueue initializes a new event queue.
func NewEventQueue(ledgerInstance *ledger.Ledger) *EventQueue {
	return &EventQueue{
		events:         make([]ScheduledEvent, 0),
		ledgerInstance: ledgerInstance,
	}
}

// ScheduleEvent schedules a new event for a future time.
func (eq *EventQueue) ScheduleEvent(opcode byte, operands []int, delay time.Duration) (string, error) {
    eq.mutex.Lock()
    defer eq.mutex.Unlock()

    eventID := generateUniqueEventID()
    event := ScheduledEvent{
        EventID:       eventID,
        ExecutionTime: time.Now().Add(delay),
        Opcode:        opcode,
        Operands:      operands,
        Status:        "scheduled",
    }

    eq.events = append(eq.events, event)
    eq.ledgerInstance.VirtualMachineLedger.RecordScheduledEvent(event.EventID, event.ExecutionTime, opcode) // Log event scheduling in the ledger

    return eventID, nil
}


// ExecuteDueEvents checks the event queue and executes all events whose time has come.
func (eq *EventQueue) ExecuteDueEvents(vm *VirtualMachine) {
    eq.mutex.Lock()
    defer eq.mutex.Unlock()

    currentTime := time.Now()
    for i, event := range eq.events {
        if event.ExecutionTime.Before(currentTime) && event.Status == "scheduled" {
            err := eq.executeEvent(vm.vm, &event) // Pass vm.vm which is of type VMInterface
            if err != nil {
                event.Status = "failed"
                eq.ledgerInstance.VirtualMachineLedger.RecordFailedEvent(event.EventID, err.Error()) // Pass the event ID and error message separately
            } else {
                event.Status = "executed"
                eq.ledgerInstance.VirtualMachineLedger.RecordExecutedEvent(event.EventID)
            }
            // Update event in queue
            eq.events[i] = event
        }
    }
}


// executeEvent executes a scheduled event using its opcode and operands.
func (eq *EventQueue) executeEvent(vm VMInterface, event *ScheduledEvent) error {
	// Convert Opcode from byte to string and set up the instruction payload
	instruction := Instruction{
		ID:        event.EventID,
		Opcode:    string(event.Opcode), // Convert byte Opcode to string
		Payload:   event.Operands,       // Store operands in the Payload field
		Timestamp: event.ExecutionTime,
	}

	// Use ExecuteInstruction from VMInterface
	return vm.ExecuteInstruction(instruction)
}


// RemoveExecutedEvents clears executed events from the queue to free resources.
func (eq *EventQueue) RemoveExecutedEvents() {
	eq.mutex.Lock()
	defer eq.mutex.Unlock()

	filteredEvents := eq.events[:0]
	for _, event := range eq.events {
		if event.Status != "executed" {
			filteredEvents = append(filteredEvents, event)
		}
	}
	eq.events = filteredEvents
}

// generateUniqueEventID generates a unique event ID (you could use UUIDs or a similar approach).
func generateUniqueEventID() string {
	return fmt.Sprintf("event-%d", time.Now().UnixNano())
}

// EventScheduler manages scheduling and periodic execution checks.
type EventScheduler struct {
	eventQueue    *EventQueue
	ticker        *time.Ticker
	stopChannel   chan bool
}

// NewEventScheduler initializes an EventScheduler with a specific interval for checking due events.
func NewEventScheduler(eventQueue *EventQueue, interval time.Duration) *EventScheduler {
	return &EventScheduler{
		eventQueue:  eventQueue,
		ticker:      time.NewTicker(interval),
		stopChannel: make(chan bool),
	}
}

// Start initiates the periodic event checking and execution process.
func (es *EventScheduler) Start(vm *VirtualMachine) {
	go func() {
		for {
			select {
			case <-es.ticker.C:
				es.eventQueue.ExecuteDueEvents(vm)
			case <-es.stopChannel:
				es.ticker.Stop()
				return
			}
		}
	}()
}

// Stop halts the scheduler's periodic event execution.
func (es *EventScheduler) Stop() {
	es.stopChannel <- true
	close(es.stopChannel)
}
