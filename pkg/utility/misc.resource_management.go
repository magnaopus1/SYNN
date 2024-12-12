package utility

import (
	"errors"
	"sync"
	"synnergy_network/pkg/common"
	"time"
)

var resourcePriorityMap sync.Map
var scheduledEvents sync.Map
var resourceAllocations sync.Map

// SET_EXECUTION_PRIORITY: Sets the execution priority for a specific process
func SET_EXECUTION_PRIORITY(processID string, priority int) {
    resourcePriorityMap.Store(processID, priority)
    LogDiagnostic("Execution Priority", "Execution priority set for process: " + processID)
}

// GET_CONTRACT_UPTIME: Retrieves uptime for a specific contract
func GET_CONTRACT_UPTIME(contractID string) (time.Duration, error) {
    startTime, err := common.ledger.GetContractStartTime(contractID)
    if err != nil {
        return 0, err
    }
    uptime := time.Since(startTime)
    LogDiagnostic("Contract Uptime", "Contract " + contractID + " uptime: " + uptime.String())
    return uptime, nil
}

// CANCEL_SCHEDULED_EVENT: Cancels a previously scheduled event
func CANCEL_SCHEDULED_EVENT(eventID string) error {
    if _, exists := scheduledEvents.Load(eventID); !exists {
        return errors.New("event not found")
    }
    scheduledEvents.Delete(eventID)
    LogDiagnostic("Scheduled Event", "Event " + eventID + " cancelled")
    return nil
}

// EXECUTE_IF_CONDITION_MET: Executes a function if a specified condition is met
func EXECUTE_IF_CONDITION_MET(condition func() bool, fn func()) {
    if condition() {
        fn()
        LogDiagnostic("Conditional Execution", "Function executed as condition met")
    }
}

// RESERVE_RESOURCE: Reserves a specific resource for a process
func RESERVE_RESOURCE(resourceID string, processID string) error {
    if _, inUse := resourceAllocations.Load(resourceID); inUse {
        return errors.New("resource is already reserved")
    }
    resourceAllocations.Store(resourceID, processID)
    LogDiagnostic("Resource Reservation", "Resource " + resourceID + " reserved for process " + processID)
    return nil
}

// RELEASE_RESOURCE: Releases a reserved resource
func RELEASE_RESOURCE(resourceID string) error {
    if _, inUse := resourceAllocations.Load(resourceID); !inUse {
        return errors.New("resource not currently reserved")
    }
    resourceAllocations.Delete(resourceID)
    LogDiagnostic("Resource Release", "Resource " + resourceID + " released")
    return nil
}

// ALLOCATE_RESOURCE: Allocates a resource for a specific transaction
func ALLOCATE_RESOURCE(txID, resourceID string) error {
    err := common.ledger.AllocateResource(txID, resourceID)
    if err != nil {
        LogDiagnostic("Resource Allocation", "Failed to allocate resource " + resourceID + " for transaction " + txID)
        return err
    }
    LogDiagnostic("Resource Allocation", "Resource " + resourceID + " allocated to transaction " + txID)
    return nil
}

// DEALLOCATE_RESOURCE: Deallocates a resource from a transaction
func DEALLOCATE_RESOURCE(txID, resourceID string) error {
    err := common.ledger.DeallocateResource(txID, resourceID)
    if err != nil {
        LogDiagnostic("Resource Deallocation", "Failed to deallocate resource " + resourceID + " from transaction " + txID)
        return err
    }
    LogDiagnostic("Resource Deallocation", "Resource " + resourceID + " deallocated from transaction " + txID)
    return nil
}

// SET_RESOURCE_TIMEOUT: Sets a timeout for resource usage
func SET_RESOURCE_TIMEOUT(resourceID string, duration time.Duration) {
    go func() {
        time.Sleep(duration)
        RELEASE_RESOURCE(resourceID)
    }()
    LogDiagnostic("Resource Timeout", "Timeout set for resource " + resourceID)
}

// GET_RESOURCE_STATUS: Retrieves the status of a specific resource
func GET_RESOURCE_STATUS(resourceID string) (string, error) {
    status, err := common.ledger.GetResourceStatus(resourceID)
    if err != nil {
        LogDiagnostic("Resource Status", "Failed to retrieve status for resource " + resourceID)
        return "", err
    }
    LogDiagnostic("Resource Status", "Status for resource " + resourceID + ": " + status)
    return status, nil
}

// RECLAIM_MEMORY: Releases unused memory in the system
func RECLAIM_MEMORY() {
    // Hypothetical memory reclamation process
    common.system.ReclaimMemory()
    LogDiagnostic("Memory Reclamation", "Unused memory reclaimed")
}

// FORCE_GARBAGE_COLLECTION: Forces a garbage collection cycle to free memory
func FORCE_GARBAGE_COLLECTION() {
    common.system.ForceGC()
    LogDiagnostic("Garbage Collection", "Forced garbage collection executed")
}

// RESET_EXECUTION_STATE: Resets the execution state of the blockchain system
func RESET_EXECUTION_STATE() error {
    err := common.ledger.ResetExecutionState()
    if err != nil {
        LogDiagnostic("Execution State", "Failed to reset execution state")
        return err
    }
    LogDiagnostic("Execution State", "Execution state reset")
    return nil
}

// SCHEDULE_RECURRING_EVENT: Schedules an event to occur at regular intervals
func SCHEDULE_RECURRING_EVENT(eventID string, interval time.Duration, fn func()) {
    scheduledEvents.Store(eventID, true)
    go func() {
        for {
            if _, active := scheduledEvents.Load(eventID); !active {
                break
            }
            fn()
            time.Sleep(interval)
        }
    }()
    LogDiagnostic("Recurring Event", "Recurring event " + eventID + " scheduled")
}

// CANCEL_RECURRING_EVENT: Cancels a recurring event
func CANCEL_RECURRING_EVENT(eventID string) error {
    if _, exists := scheduledEvents.Load(eventID); !exists {
        return errors.New("recurring event not found")
    }
    scheduledEvents.Delete(eventID)
    LogDiagnostic("Recurring Event", "Recurring event " + eventID + " cancelled")
    return nil
}


