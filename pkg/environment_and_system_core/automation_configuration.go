package environment_and_system_core

import (
	"fmt"
	"log"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// Lock for managing concurrent access to event configuration
var automationLock sync.Mutex

func defineEventAction(ledgerInstance *ledger.Ledger, eventID string, action string) error {
    automationLock.Lock()
    defer automationLock.Unlock()

    if ledgerInstance == nil {
        return fmt.Errorf("defineEventAction: ledger instance cannot be nil")
    }
    if eventID == "" {
        return fmt.Errorf("defineEventAction: eventID cannot be empty")
    }
    if action == "" {
        return fmt.Errorf("defineEventAction: action cannot be empty")
    }

    if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordEventAction(eventID, action); err != nil {
        return fmt.Errorf("defineEventAction: failed to define action for event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Event action defined for event %s: %s", eventID, action)
    return nil
}

func executeOnEventTrigger(ledgerInstance *ledger.Ledger, eventID string) error {
    if ledgerInstance == nil {
        return fmt.Errorf("executeOnEventTrigger: ledger instance cannot be nil")
    }
    if eventID == "" {
        return fmt.Errorf("executeOnEventTrigger: eventID cannot be empty")
    }

    action, err := ledgerInstance.EnvironmentSystemCoreLedger.GetEventAction(eventID)
    if err != nil {
        return fmt.Errorf("executeOnEventTrigger: failed to retrieve action for event %s: %w", eventID, err)
    }

    if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecuteAction(action); err != nil {
        return fmt.Errorf("executeOnEventTrigger: failed to execute action for event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Event triggered and action executed for event %s.", eventID)
    return nil
}


func validateTriggerCondition(ledgerInstance *ledger.Ledger, eventID string) (bool, error) {
    if ledgerInstance == nil {
        return false, fmt.Errorf("validateTriggerCondition: ledger instance cannot be nil")
    }
    if eventID == "" {
        return false, fmt.Errorf("validateTriggerCondition: eventID cannot be empty")
    }

    condition, err := ledgerInstance.EnvironmentSystemCoreLedger.GetTriggerCondition(eventID)
    if err != nil {
        return false, fmt.Errorf("validateTriggerCondition: failed to retrieve condition for event %s: %w", eventID, err)
    }

    isValid := condition.IsValid()
    log.Printf("[INFO] Trigger condition validation for event %s: %t", eventID, isValid)
    return isValid, nil
}


func setTriggerThreshold(ledgerInstance *ledger.Ledger, eventID string, threshold int) error {
    if ledgerInstance == nil {
        return fmt.Errorf("setTriggerThreshold: ledger instance cannot be nil")
    }
    if eventID == "" {
        return fmt.Errorf("setTriggerThreshold: eventID cannot be empty")
    }
    if threshold <= 0 {
        return fmt.Errorf("setTriggerThreshold: threshold must be positive")
    }

    condition, err := ledgerInstance.EnvironmentSystemCoreLedger.GetTriggerCondition(eventID)
    if err != nil {
        return fmt.Errorf("setTriggerThreshold: failed to retrieve trigger condition for event %s: %w", eventID, err)
    }

    condition.Details["threshold"] = threshold
    if err := ledgerInstance.EnvironmentSystemCoreLedger.SetTriggerCondition(eventID, condition); err != nil {
        return fmt.Errorf("setTriggerThreshold: failed to set trigger threshold for event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Trigger threshold set for event %s: %d", eventID, threshold)
    return nil
}


func retrieveEventHistory(ledgerInstance *ledger.Ledger, eventID string) ([]ledger.EventLog, error) {
    if ledgerInstance == nil {
        return nil, fmt.Errorf("retrieveEventHistory: ledger instance cannot be nil")
    }
    if eventID == "" {
        return nil, fmt.Errorf("retrieveEventHistory: eventID cannot be empty")
    }

    history, err := ledgerInstance.EnvironmentSystemCoreLedger.GetEventHistory(eventID)
    if err != nil {
        return nil, fmt.Errorf("retrieveEventHistory: failed to retrieve history for event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Event history retrieved for event %s.", eventID)
    return history, nil
}


func logEventFailure(ledgerInstance *ledger.Ledger, eventID, reason string) error {
    if err := ledgerInstance.EnvironmentSystemCoreLedger.LogFailure(eventID, reason); err != nil {
        return fmt.Errorf("failed to log event failure: %v", err)
    }
    return nil
}

func notifyOnEventTrigger(ledgerInstance *ledger.Ledger, eventID string, recipients []string) error {
    if ledgerInstance == nil {
        return fmt.Errorf("notifyOnEventTrigger: ledger instance cannot be nil")
    }
    if eventID == "" {
        return fmt.Errorf("notifyOnEventTrigger: eventID cannot be empty")
    }
    if len(recipients) == 0 {
        return fmt.Errorf("notifyOnEventTrigger: recipients list cannot be empty")
    }

    notification := ledgerInstance.EnvironmentSystemCoreLedger.CreateNotification(eventID, recipients)
    if err := ledgerInstance.EnvironmentSystemCoreLedger.DispatchNotification(notification); err != nil {
        return fmt.Errorf("notifyOnEventTrigger: failed to dispatch notification for event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Notification dispatched for event %s to recipients: %v", eventID, recipients)
    return nil
}


func delayEventExecution(ledgerInstance *ledger.Ledger, eventID string, delay time.Duration) error {
    if ledgerInstance == nil {
        return fmt.Errorf("delayEventExecution: ledger instance cannot be nil")
    }
    if eventID == "" {
        return fmt.Errorf("delayEventExecution: eventID cannot be empty")
    }
    if delay <= 0 {
        return fmt.Errorf("delayEventExecution: delay must be a positive duration")
    }

    if err := ledgerInstance.EnvironmentSystemCoreLedger.DelayEvent(eventID, delay); err != nil {
        return fmt.Errorf("delayEventExecution: failed to delay event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Event %s delayed by %v", eventID, delay)
    return nil
}


func resumeScheduledEvent(ledgerInstance *ledger.Ledger, eventID string) error {
    if ledgerInstance == nil {
        return fmt.Errorf("resumeScheduledEvent: ledger instance cannot be nil")
    }
    if eventID == "" {
        return fmt.Errorf("resumeScheduledEvent: eventID cannot be empty")
    }

    if err := ledgerInstance.EnvironmentSystemCoreLedger.ResumeEvent(eventID); err != nil {
        return fmt.Errorf("resumeScheduledEvent: failed to resume event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Event %s resumed successfully", eventID)
    return nil
}

func queryEventProperties(ledgerInstance *ledger.Ledger, eventID string) (map[string]interface{}, error) {
    if ledgerInstance == nil {
        return nil, fmt.Errorf("queryEventProperties: ledger instance cannot be nil")
    }
    if eventID == "" {
        return nil, fmt.Errorf("queryEventProperties: eventID cannot be empty")
    }

    properties, err := ledgerInstance.EnvironmentSystemCoreLedger.GetEventProperties(eventID)
    if err != nil {
        return nil, fmt.Errorf("queryEventProperties: failed to retrieve properties for event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Retrieved properties for event %s", eventID)
    return properties, nil
}


func setRecurringEvent(ledgerInstance *ledger.Ledger, eventID string, interval time.Duration) error {
    if ledgerInstance == nil {
        return fmt.Errorf("setRecurringEvent: ledger instance cannot be nil")
    }
    if eventID == "" {
        return fmt.Errorf("setRecurringEvent: eventID cannot be empty")
    }
    if interval <= 0 {
        return fmt.Errorf("setRecurringEvent: interval must be positive")
    }

    if err := ledgerInstance.EnvironmentSystemCoreLedger.SetRecurringEvent(eventID, interval); err != nil {
        return fmt.Errorf("setRecurringEvent: failed to set recurring event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Recurring event %s set with interval %v", eventID, interval)
    return nil
}


func checkRecurringEventStatus(ledgerInstance *ledger.Ledger, eventID string) (bool, error) {
    if ledgerInstance == nil {
        return false, fmt.Errorf("checkRecurringEventStatus: ledger instance cannot be nil")
    }
    if eventID == "" {
        return false, fmt.Errorf("checkRecurringEventStatus: eventID cannot be empty")
    }

    status, err := ledgerInstance.EnvironmentSystemCoreLedger.IsRecurringEvent(eventID)
    if err != nil {
        return false, fmt.Errorf("checkRecurringEventStatus: failed to check status for recurring event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Recurring status for event %s: %t", eventID, status)
    return status, nil
}


func updateEventTriggerConditions(ledgerInstance *ledger.Ledger, eventID string, newConditions map[string]interface{}) error {
    if ledgerInstance == nil {
        return fmt.Errorf("updateEventTriggerConditions: ledger instance cannot be nil")
    }
    if eventID == "" {
        return fmt.Errorf("updateEventTriggerConditions: eventID cannot be empty")
    }
    if len(newConditions) == 0 {
        return fmt.Errorf("updateEventTriggerConditions: newConditions cannot be empty")
    }

    if err := ledgerInstance.EnvironmentSystemCoreLedger.UpdateTriggerConditions(eventID, newConditions); err != nil {
        return fmt.Errorf("updateEventTriggerConditions: failed to update conditions for event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Trigger conditions updated for event %s", eventID)
    return nil
}


func manageEventDependencies(ledgerInstance *ledger.Ledger, eventID string, dependencies []string) error {
    if ledgerInstance == nil {
        return fmt.Errorf("manageEventDependencies: ledger instance cannot be nil")
    }
    if eventID == "" {
        return fmt.Errorf("manageEventDependencies: eventID cannot be empty")
    }
    if len(dependencies) == 0 {
        return fmt.Errorf("manageEventDependencies: dependencies cannot be empty")
    }

    if err := ledgerInstance.EnvironmentSystemCoreLedger.SetEventDependencies(eventID, dependencies); err != nil {
        return fmt.Errorf("manageEventDependencies: failed to set dependencies for event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Dependencies set for event %s: %v", eventID, dependencies)
    return nil
}


func resetEventSequence(ledgerInstance *ledger.Ledger, eventID string) error {
    if ledgerInstance == nil {
        return fmt.Errorf("resetEventSequence: ledger instance cannot be nil")
    }
    if eventID == "" {
        return fmt.Errorf("resetEventSequence: eventID cannot be empty")
    }

    if err := ledgerInstance.EnvironmentSystemCoreLedger.ResetEventSequence(eventID); err != nil {
        return fmt.Errorf("resetEventSequence: failed to reset sequence for event %s: %w", eventID, err)
    }

    log.Printf("[INFO] Event sequence reset for event %s", eventID)
    return nil
}


func validateAutomationSchedule(ledgerInstance *ledger.Ledger) (bool, error) {
    if ledgerInstance == nil {
        return false, fmt.Errorf("validateAutomationSchedule: ledger instance cannot be nil")
    }

    schedule, err := ledgerInstance.EnvironmentSystemCoreLedger.GetAutomationSchedule()
    if err != nil {
        return false, fmt.Errorf("validateAutomationSchedule: failed to retrieve automation schedule: %w", err)
    }

    isValid := schedule.IsValid()
    log.Printf("[INFO] Automation schedule validation: %t", isValid)
    return isValid, nil
}

