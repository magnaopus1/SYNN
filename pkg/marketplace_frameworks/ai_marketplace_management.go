package marketplace

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

func scheduleAIUsage(moduleID string, usageSchedule AIUsageSchedule, ledgerInstance *Ledger) error {
    usageSchedule.ModuleID = moduleID
    usageSchedule.ScheduledAt = time.Now()
    return ledgerInstance.recordAIUsageSchedule(usageSchedule)
}

func enableAIUsageRestrictions(moduleID string, ledgerInstance *Ledger) error {
    module, err := ledgerInstance.getAIModule(moduleID)
    if err != nil {
        return fmt.Errorf("failed to retrieve AI module: %v", err)
    }
    module.UsageRestrictionsEnabled = true
    return ledgerInstance.updateAIModule(module)
}

func disableAIUsageRestrictions(moduleID string, ledgerInstance *Ledger) error {
    module, err := ledgerInstance.getAIModule(moduleID)
    if err != nil {
        return fmt.Errorf("failed to retrieve AI module: %v", err)
    }
    module.UsageRestrictionsEnabled = false
    return ledgerInstance.updateAIModule(module)
}

func logAIEvent(moduleID, eventType, description string, ledgerInstance *Ledger) error {
    encryptedDescription, err := encryption.EncryptData(description)
    if err != nil {
        return fmt.Errorf("failed to encrypt event description: %v", err)
    }
    event := AIEventLog{
        ModuleID:    moduleID,
        EventType:   eventType,
        Description: encryptedDescription,
        Timestamp:   time.Now(),
    }
    return ledgerInstance.recordAIEvent(event)
}

func approveAINetworkUsage(requestID string, ledgerInstance *Ledger) error {
    return ledgerInstance.updateNetworkUsageRequest(requestID, "Approved")
}

func denyAINetworkUsage(requestID string, ledgerInstance *Ledger) error {
    return ledgerInstance.updateNetworkUsageRequest(requestID, "Denied")
}

func setAIResourceAllocation(moduleID string, allocation AIResourceAllocation, ledgerInstance *Ledger) error {
    allocation.ModuleID = moduleID
    allocation.Timestamp = time.Now()
    return ledgerInstance.recordAIResourceAllocation(allocation)
}

func trackAIModelPerformance(moduleID string, metrics AIMetrics, ledgerInstance *Ledger) error {
    metrics.ModuleID = moduleID
    metrics.Timestamp = time.Now()
    return ledgerInstance.recordAIModelPerformance(metrics)
}

func retrieveAIModelMetrics(moduleID string, ledgerInstance *Ledger) (AIMetrics, error) {
    return ledgerInstance.getAIModelMetrics(moduleID)
}

func updateAITrainingData(moduleID string, newData AITrainingData, ledgerInstance *Ledger) error {
    newData.ModuleID = moduleID
    newData.UpdatedAt = time.Now()
    return ledgerInstance.updateAITrainingData(newData)
}

func manageAIModelVersioning(moduleID, version string, changes AIModelChanges, ledgerInstance *Ledger) error {
    versionRecord := AIModelVersion{
        ModuleID:   moduleID,
        Version:    version,
        Changes:    changes,
        Timestamp:  time.Now(),
    }
    return ledgerInstance.recordAIModelVersion(versionRecord)
}

func archiveAIModule(moduleID string, ledgerInstance *Ledger) error {
    module, err := ledgerInstance.getAIModule(moduleID)
    if err != nil {
        return fmt.Errorf("failed to retrieve AI module: %v", err)
    }
    module.Archived = true
    return ledgerInstance.updateAIModule(module)
}

func activateAIModule(moduleID string, ledgerInstance *Ledger) error {
    module, err := ledgerInstance.getAIModule(moduleID)
    if err != nil {
        return fmt.Errorf("failed to retrieve AI module: %v", err)
    }
    module.Archived = false
    return ledgerInstance.updateAIModule(module)
}


func deactivateAIModule(moduleID string, ledgerInstance *Ledger) error {
    module, err := ledgerInstance.getAIModule(moduleID)
    if err != nil {
        return fmt.Errorf("failed to retrieve AI module: %v", err)
    }
    module.Active = false
    return ledgerInstance.updateAIModule(module)
}

func assignAITask(moduleID string, task AITask, ledgerInstance *Ledger) error {
    task.ModuleID = moduleID
    task.AssignedAt = time.Now()
    task.Completed = false
    return ledgerInstance.recordAITask(task)
}

func verifyAITaskCompletion(taskID string, ledgerInstance *Ledger) (bool, error) {
    task, err := ledgerInstance.getAITask(taskID)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve AI task: %v", err)
    }
    return task.Completed, nil
}

func rewardAIModuleUsage(moduleID string, reward AIReward, ledgerInstance *Ledger) error {
    reward.ModuleID = moduleID
    reward.IssuedAt = time.Now()
    return ledgerInstance.recordAIReward(reward)
}

func penaltyAIModuleMisuse(moduleID string, penalty AIPenalty, ledgerInstance *Ledger) error {
    penalty.ModuleID = moduleID
    penalty.IssuedAt = time.Now()
    return ledgerInstance.recordAIPenalty(penalty)
}

func linkAIModelToDataset(moduleID, datasetID string, ledgerInstance *Ledger) error {
    association := AIDatasetLink{
        ModuleID:  moduleID,
        DatasetID: datasetID,
        LinkedAt:  time.Now(),
    }
    return ledgerInstance.recordAIDatasetLink(association)
}

func removeAIModelDatasetLink(moduleID, datasetID string, ledgerInstance *Ledger) error {
    return ledgerInstance.removeAIDatasetLink(moduleID, datasetID)
}
