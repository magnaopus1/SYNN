package interoperability

import (
	"fmt"
	"time"
	"sync"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// CrossChainDisputeManager handles the management of cross-chain disputes
type CrossChainDisputeManager struct {
	consensusEngine *common.SynnergyConsensus
	ledgerInstance  *ledger.Ledger
	disputeMutex    *sync.RWMutex
}

// NewCrossChainDisputeManager initializes the CrossChainDisputeManager
func NewCrossChainDisputeManager(consensusEngine *common.SynnergyConsensus, ledgerInstance *ledger.Ledger, disputeMutex *sync.RWMutex) *CrossChainDisputeManager {
	return &CrossChainDisputeManager{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		disputeMutex:    disputeMutex,
	}
}

// InitiateDispute initiates a new cross-chain dispute and logs it in the ledger
func (manager *CrossChainDisputeManager) initiateDispute(disputeID, initiatorID, reason string) error {
    manager.disputeMutex.Lock()
    defer manager.disputeMutex.Unlock()

    dispute := Dispute{
        DisputeID:   disputeID,
        InitiatorID: initiatorID,
        Reason:      reason,
        Status:      "Initiated",
        CreatedAt:   time.Now(),
    }

    if err := manager.ledgerInstance.addDispute(dispute); err != nil {
        return fmt.Errorf("failed to initiate dispute %s: %v", disputeID, err)
    }

    manager.logDisputeEvent(disputeID, "Initiated", fmt.Sprintf("Dispute initiated by %s for reason: %s", initiatorID, reason))
    return nil
}

// ResolveDispute resolves an existing dispute, recording the outcome in the ledger
func (manager *CrossChainDisputeManager) resolveDispute(disputeID, resolution string) error {
    manager.disputeMutex.Lock()
    defer manager.disputeMutex.Unlock()

    if err := manager.ledgerInstance.resolveDispute(disputeID, resolution); err != nil {
        return fmt.Errorf("failed to resolve dispute %s: %v", disputeID, err)
    }

    manager.logDisputeEvent(disputeID, "Resolved", fmt.Sprintf("Dispute resolved with outcome: %s", resolution))
    return nil
}

// RegisterDisputeHandler assigns a handler to manage the specified dispute
func (manager *CrossChainDisputeManager) registerDisputeHandler(disputeID, handlerID string) error {
    manager.logDisputeEvent(disputeID, "Handler Registered", fmt.Sprintf("Handler %s assigned to dispute", handlerID))
    return nil
}

// RecordDisputeResolution securely records the details of the dispute resolution
func (manager *CrossChainDisputeManager) recordDisputeResolution(disputeID, resolutionDetails string) error {
    encryptedDetails := manager.encryptData(resolutionDetails)

    event := DisputeEvent{
        DisputeID: disputeID,
        EventType: "Resolution Recorded",
        Details:   encryptedDetails,
        Timestamp: time.Now(),
    }

    if err := manager.ledgerInstance.logDisputeEvent(event); err != nil {
        return fmt.Errorf("failed to record dispute resolution for dispute %s: %v", disputeID, err)
    }

    return nil
}

// EscalateCrossChainDispute escalates a dispute to a higher authority for resolution
func (manager *CrossChainDisputeManager) escalateCrossChainDispute(disputeID, reason string) error {
    manager.logDisputeEvent(disputeID, "Escalated", fmt.Sprintf("Dispute escalated due to: %s", reason))
    return nil
}

// AssignDisputeMediator assigns a mediator to the dispute for arbitration
func (manager *CrossChainDisputeManager) assignDisputeMediator(disputeID, mediatorID string) error {
    mediator := MediatorAssignment{
        DisputeID:  disputeID,
        MediatorID: mediatorID,
        AssignedAt: time.Now(),
    }

    if err := manager.ledgerInstance.assignMediator(mediator); err != nil {
        return fmt.Errorf("failed to assign mediator %s to dispute %s: %v", mediatorID, disputeID, err)
    }

    manager.logDisputeEvent(disputeID, "Mediator Assigned", fmt.Sprintf("Mediator %s assigned to dispute", mediatorID))
    return nil
}

// UnassignDisputeMediator removes an assigned mediator from the dispute
func (manager *CrossChainDisputeManager) unassignDisputeMediator(disputeID, mediatorID string) error {
    if err := manager.ledgerInstance.unassignMediator(disputeID); err != nil {
        return fmt.Errorf("failed to unassign mediator %s from dispute %s: %v", mediatorID, disputeID, err)
    }

    manager.logDisputeEvent(disputeID, "Mediator Unassigned", fmt.Sprintf("Mediator %s unassigned from dispute", mediatorID))
    return nil
}

// LogDisputeEvent logs events related to disputes for traceability
func (manager *CrossChainDisputeManager) logDisputeEvent(disputeID, eventType, details string) {
    eventDetails := fmt.Sprintf("Event: %s, Dispute ID: %s, Details: %s", eventType, disputeID, details)
    encryptedDetails := manager.encryptData(eventDetails)

    event := DisputeEvent{
        DisputeID: disputeID,
        EventType: eventType,
        Details:   encryptedDetails,
        Timestamp: time.Now(),
    }

    if err := manager.ledgerInstance.logDisputeEvent(event); err != nil {
        fmt.Printf("Failed to log dispute event for dispute ID %s: %v\n", disputeID, err)
    }
}

// encryptData encrypts data for secure logging and storage
func (manager *CrossChainDisputeManager) encryptData(data string) string {
    encryptedData, err := encryption.EncryptData([]byte(data))
    if err != nil {
        fmt.Println("Error encrypting data:", err)
        return data
    }
    return string(encryptedData)
}


// validateDisputeEvidence validates evidence provided for a dispute using the consensus mechanism
func (manager *CrossChainDisputeManager) validateDisputeEvidence(disputeID, evidence string) (bool, error) {
    manager.disputeMutex.Lock()
    defer manager.disputeMutex.Unlock()

    evidenceID := fmt.Sprintf("evidence-%s-%d", disputeID, time.Now().Unix())
    disputeEvidence := DisputeEvidence{
        EvidenceID:  evidenceID,
        DisputeID:   disputeID,
        Content:     evidence,
        Validated:   false,
        SubmittedAt: time.Now(),
    }

    if err := manager.ledgerInstance.addDisputeEvidence(disputeEvidence); err != nil {
        return false, fmt.Errorf("failed to add evidence for dispute %s: %v", disputeID, err)
    }

    isValid, err := manager.consensusEngine.ValidateEvidence(evidence)
    if err != nil || !isValid {
        manager.logDisputeEvent(disputeID, "Evidence Validation Failed", "Failed to validate dispute evidence")
        return false, fmt.Errorf("evidence validation failed for dispute %s", disputeID)
    }

    if err := manager.ledgerInstance.validateDisputeEvidence(evidenceID); err != nil {
        return false, fmt.Errorf("failed to mark evidence %s as validated: %v", evidenceID, err)
    }

    manager.logDisputeEvent(disputeID, "Evidence Validated", "Evidence validation successful")
    return true, nil
}

// generateArbitrationSummary generates a summary of the arbitration process for the dispute
func (manager *CrossChainDisputeManager) generateArbitrationSummary(disputeID, summaryDetails string) error {
    manager.disputeMutex.Lock()
    defer manager.disputeMutex.Unlock()

    summaryID := fmt.Sprintf("summary-%s-%d", disputeID, time.Now().Unix())
    encryptedSummary := manager.encryptData(summaryDetails)

    arbitrationSummary := ArbitrationSummary{
        SummaryID:   summaryID,
        DisputeID:   disputeID,
        Summary:     encryptedSummary,
        GeneratedAt: time.Now(),
    }

    if err := manager.ledgerInstance.addArbitrationSummary(arbitrationSummary); err != nil {
        return fmt.Errorf("failed to generate arbitration summary for dispute %s: %v", disputeID, err)
    }

    manager.logDisputeEvent(disputeID, "Arbitration Summary Generated", "Summary generated successfully")
    return nil
}

// logDisputeEvent logs dispute-related events in the ledger for traceability
func (manager *CrossChainDisputeManager) logDisputeEvent(disputeID, eventType, details string) {
    eventDetails := fmt.Sprintf("Event: %s, Dispute ID: %s, Details: %s", eventType, disputeID, details)
    encryptedDetails := manager.encryptData(eventDetails)

    if err := manager.ledgerInstance.logDisputeEvent(disputeID, eventType, encryptedDetails); err != nil {
        fmt.Printf("Failed to log dispute event for dispute ID %s: %v\n", disputeID, err)
    }
}

// encryptData encrypts data for secure logging and storage
func (manager *CrossChainDisputeManager) encryptData(data string) string {
    encryptedData, err := encryption.EncryptData([]byte(data))
    if err != nil {
        fmt.Println("Error encrypting data:", err)
        return data
    }
    return string(encryptedData)
}

// decryptData decrypts stored encrypted data
func (manager *CrossChainDisputeManager) decryptData(encryptedData string) string {
    decryptedData, err := encryption.DecryptData([]byte(encryptedData))
    if err != nil {
        fmt.Println("Error decrypting data:", err)
        return encryptedData
    }
    return string(decryptedData)
}
