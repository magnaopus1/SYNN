package environment_and_system_core

import (
	"fmt"
	"log"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// Mutex for managing concurrent access to event trigger operations
var eventTriggerLock sync.Mutex

// createEventTrigger creates a new event trigger with encrypted conditions and specified actions.
func createEventTrigger(ledgerInstance *ledger.Ledger, triggerID string, conditions map[string]interface{}, actions []string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("createEventTrigger: ledger instance cannot be nil")
	}
	if triggerID == "" {
		return fmt.Errorf("createEventTrigger: triggerID cannot be empty")
	}
	if len(conditions) == 0 {
		return fmt.Errorf("createEventTrigger: conditions cannot be empty")
	}
	if len(actions) == 0 {
		return fmt.Errorf("createEventTrigger: actions cannot be empty")
	}

	// Lock for concurrency safety
	eventTriggerLock.Lock()
	defer eventTriggerLock.Unlock()

	// Encrypt the conditions
	encryptedConditions, err := common.EncryptData(conditions)
	if err != nil {
		return fmt.Errorf("createEventTrigger: failed to encrypt conditions: %w", err)
	}

	// Record the event trigger in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.EventManager.RecordEventTrigger(triggerID, encryptedConditions, actions); err != nil {
		return fmt.Errorf("createEventTrigger: failed to create event trigger %s: %w", triggerID, err)
	}

	// Log success
	log.Printf("Event trigger %s created successfully.", triggerID)
	return nil
}


// removeEventTrigger removes an existing event trigger.
func removeEventTrigger(ledgerInstance *ledger.Ledger, triggerID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("removeEventTrigger: ledger instance cannot be nil")
	}
	if triggerID == "" {
		return fmt.Errorf("removeEventTrigger: triggerID cannot be empty")
	}

	// Lock for concurrency safety
	eventTriggerLock.Lock()
	defer eventTriggerLock.Unlock()

	// Delete the event trigger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.EventManager.DeleteEventTrigger(triggerID); err != nil {
		return fmt.Errorf("removeEventTrigger: failed to remove event trigger %s: %w", triggerID, err)
	}

	// Log success
	log.Printf("Event trigger %s removed successfully.", triggerID)
	return nil
}


// triggerEvent manually triggers an event based on its ID.
func triggerEvent(ledgerInstance *ledger.Ledger, triggerID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("triggerEvent: ledger instance cannot be nil")
	}
	if triggerID == "" {
		return fmt.Errorf("triggerEvent: triggerID cannot be empty")
	}

	// Lock for concurrency safety
	eventTriggerLock.Lock()
	defer eventTriggerLock.Unlock()

	// Activate the trigger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.EventManager.ActivateTrigger(triggerID); err != nil {
		return fmt.Errorf("triggerEvent: failed to trigger event %s: %w", triggerID, err)
	}

	// Log success
	log.Printf("Event %s triggered successfully.", triggerID)
	return nil
}


// checkEventStatus retrieves the current status of an event trigger.
func checkEventStatus(ledgerInstance *ledger.Ledger, triggerID string) (ledger.EventStatus, error) {
	// Validate inputs
	if ledgerInstance == nil {
		return ledger.EventStatus{}, fmt.Errorf("checkEventStatus: ledger instance cannot be nil")
	}
	if triggerID == "" {
		return ledger.EventStatus{}, fmt.Errorf("checkEventStatus: triggerID cannot be empty")
	}

	// Get the event status
	status, err := ledgerInstance.EnvironmentSystemCoreLedger.EventManager.GetEventStatus(triggerID)
	if err != nil {
		return ledger.EventStatus{}, fmt.Errorf("checkEventStatus: failed to retrieve status for event %s: %w", triggerID, err)
	}

	// Return the status
	return status, nil
}

// scheduleAutomationTask schedules an automation task with specified details.
func scheduleAutomationTask(ledgerInstance *ledger.Ledger, taskID string, scheduledTime time.Time, taskDetails map[string]interface{}) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("scheduleAutomationTask: ledger instance cannot be nil")
	}
	if taskID == "" {
		return fmt.Errorf("scheduleAutomationTask: taskID cannot be empty")
	}
	if taskDetails == nil {
		return fmt.Errorf("scheduleAutomationTask: taskDetails cannot be nil")
	}

	// Lock for concurrency safety
	eventTriggerLock.Lock()
	defer eventTriggerLock.Unlock()

	// Schedule the task in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.TaskManager.ScheduleTask(taskID, scheduledTime, taskDetails); err != nil {
		return fmt.Errorf("scheduleAutomationTask: failed to schedule task %s: %w", taskID, err)
	}

	// Log success
	log.Printf("Automation task %s scheduled for %v.", taskID, scheduledTime)
	return nil
}


// cancelScheduledTask cancels a previously scheduled automation task.
func cancelScheduledTask(ledgerInstance *ledger.Ledger, taskID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("cancelScheduledTask: ledger instance cannot be nil")
	}
	if taskID == "" {
		return fmt.Errorf("cancelScheduledTask: taskID cannot be empty")
	}

	// Lock for concurrency safety
	eventTriggerLock.Lock()
	defer eventTriggerLock.Unlock()

	// Cancel the task in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.TaskManager.CancelTask(taskID); err != nil {
		return fmt.Errorf("cancelScheduledTask: failed to cancel task %s: %w", taskID, err)
	}

	// Log success
	log.Printf("Scheduled task %s canceled successfully.", taskID)
	return nil
}


// monitorCondition continuously evaluates a condition and triggers an event if met.
func monitorCondition(ledgerInstance *ledger.Ledger, conditionID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("monitorCondition: ledger instance cannot be nil")
	}
	if conditionID == "" {
		return fmt.Errorf("monitorCondition: conditionID cannot be empty")
	}

	// Evaluate the condition
	result, err := ledgerInstance.EnvironmentSystemCoreLedger.ConditionManager.EvaluateCondition(conditionID)
	if err != nil {
		return fmt.Errorf("monitorCondition: failed to evaluate condition %s: %w", conditionID, err)
	}

	// Trigger the event if the condition is met
	if result.Met {
		if err := triggerEvent(ledgerInstance, conditionID); err != nil {
			return fmt.Errorf("monitorCondition: failed to trigger event for condition %s: %w", conditionID, err)
		}
	}

	// Log evaluation result
	log.Printf("Condition %s evaluated: Met=%v.", conditionID, result.Met)
	return nil
}

// logAutomationEvent logs an automation-related event in the ledger.
func logAutomationEvent(ledgerInstance *ledger.Ledger, event ledger.AutomationEvent) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("logAutomationEvent: ledger instance cannot be nil")
	}

	// Record the automation event
	if err := ledgerInstance.EnvironmentSystemCoreLedger.EventManager.RecordAutomationLog(event); err != nil {
		return fmt.Errorf("logAutomationEvent: failed to log automation event: %w", err)
	}

	// Log success
	log.Printf("Automation event %s logged successfully.", event.EventID)
	return nil
}

// setEventPriority updates the priority of an event trigger.
func setEventPriority(ledgerInstance *ledger.Ledger, eventID string, priority int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("setEventPriority: ledger instance cannot be nil")
	}
	if eventID == "" {
		return fmt.Errorf("setEventPriority: eventID cannot be empty")
	}
	if priority < 0 {
		return fmt.Errorf("setEventPriority: priority cannot be negative")
	}

	// Update the event priority in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.EventManager.UpdateEventPriority(eventID, priority); err != nil {
		return fmt.Errorf("setEventPriority: failed to set priority for event %s: %w", eventID, err)
	}

	// Log success
	log.Printf("Event %s priority set to %d.", eventID, priority)
	return nil
}


// getNextScheduledEvent retrieves the next scheduled event from the ledger.
func getNextScheduledEvent(ledgerInstance *ledger.Ledger) (ledger.ScheduledEvent, error) {
	// Validate inputs
	if ledgerInstance == nil {
		return ledger.ScheduledEvent{}, fmt.Errorf("getNextScheduledEvent: ledger instance cannot be nil")
	}

	// Fetch the next scheduled event from the ledger
	nextEvent, err := ledgerInstance.EnvironmentSystemCoreLedger.EventScheduler.FetchNextScheduledEvent()
	if err != nil {
		return ledger.ScheduledEvent{}, fmt.Errorf("getNextScheduledEvent: failed to retrieve the next scheduled event: %w", err)
	}

	// Log success
	log.Printf("Next scheduled event retrieved successfully: EventID=%s, ScheduledTime=%v.", nextEvent.EventID, nextEvent.ScheduledTime)
	return nextEvent, nil
}

// registerEventListener registers a listener to an event by its ID.
func registerEventListener(ledgerInstance *ledger.Ledger, eventID, listenerID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("registerEventListener: ledger instance cannot be nil")
	}
	if eventID == "" || listenerID == "" {
		return fmt.Errorf("registerEventListener: eventID and listenerID cannot be empty")
	}

	// Add the event listener in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.EventManager.AddEventListener(eventID, listenerID); err != nil {
		return fmt.Errorf("registerEventListener: failed to register listener %s for event %s: %w", listenerID, eventID, err)
	}

	// Log success
	log.Printf("Listener %s registered to event %s successfully.", listenerID, eventID)
	return nil
}


// deregisterEventListener removes a listener from an event by its ID.
func deregisterEventListener(ledgerInstance *ledger.Ledger, eventID, listenerID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("deregisterEventListener: ledger instance cannot be nil")
	}
	if eventID == "" || listenerID == "" {
		return fmt.Errorf("deregisterEventListener: eventID and listenerID cannot be empty")
	}

	// Remove the event listener in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.EventManager.RemoveEventListener(eventID, listenerID); err != nil {
		return fmt.Errorf("deregisterEventListener: failed to deregister listener %s for event %s: %w", listenerID, eventID, err)
	}

	// Log success
	log.Printf("Listener %s deregistered from event %s successfully.", listenerID, eventID)
	return nil
}


// activateAutomationSequence activates an automation sequence by its ID.
func activateAutomationSequence(ledgerInstance *ledger.Ledger, sequenceID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("activateAutomationSequence: ledger instance cannot be nil")
	}
	if sequenceID == "" {
		return fmt.Errorf("activateAutomationSequence: sequenceID cannot be empty")
	}

	// Start the automation sequence in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.AutomationManager.StartAutomationSequence(sequenceID); err != nil {
		return fmt.Errorf("activateAutomationSequence: failed to activate automation sequence %s: %w", sequenceID, err)
	}

	// Log success
	log.Printf("Automation sequence %s activated successfully.", sequenceID)
	return nil
}


// deactivateAutomationSequence deactivates an automation sequence by its ID.
func deactivateAutomationSequence(ledgerInstance *ledger.Ledger, sequenceID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("deactivateAutomationSequence: ledger instance cannot be nil")
	}
	if sequenceID == "" {
		return fmt.Errorf("deactivateAutomationSequence: sequenceID cannot be empty")
	}

	// Stop the automation sequence in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.AutomationManager.StopAutomationSequence(sequenceID); err != nil {
		return fmt.Errorf("deactivateAutomationSequence: failed to deactivate automation sequence %s: %w", sequenceID, err)
	}

	// Log success
	log.Printf("Automation sequence %s deactivated successfully.", sequenceID)
	return nil
}
