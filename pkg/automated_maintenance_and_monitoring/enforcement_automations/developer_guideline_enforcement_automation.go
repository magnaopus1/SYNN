package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/developer"
)

// Configuration for developer guideline enforcement automation
const (
	DeveloperCheckInterval        = 20 * time.Second // Interval to check compliance with developer guidelines
	MaxSubmissionViolations       = 3                // Maximum submission violations before restricting developer
	PerformanceBenchmarkThreshold = 80               // Minimum performance score required for submission
	SecurityComplianceRequired    = true             // Security checks must be passed for all submissions
)

// DeveloperGuidelineEnforcementAutomation monitors and enforces developer guideline compliance
type DeveloperGuidelineEnforcementAutomation struct {
	devManager        *developer.DeveloperManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	violationCount    map[string]int // Tracks guideline violations per developer
}

// NewDeveloperGuidelineEnforcementAutomation initializes the developer guideline enforcement automation
func NewDeveloperGuidelineEnforcementAutomation(devManager *developer.DeveloperManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *DeveloperGuidelineEnforcementAutomation {
	return &DeveloperGuidelineEnforcementAutomation{
		devManager:       devManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
		violationCount:   make(map[string]int),
	}
}

// StartGuidelineEnforcement begins continuous monitoring and enforcement of developer guideline compliance
func (automation *DeveloperGuidelineEnforcementAutomation) StartGuidelineEnforcement() {
	ticker := time.NewTicker(DeveloperCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkGuidelineCompliance()
		}
	}()
}

// checkGuidelineCompliance verifies each developer submission for compliance with coding, security, and performance standards
func (automation *DeveloperGuidelineEnforcementAutomation) checkGuidelineCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, submissionID := range automation.devManager.GetRecentSubmissions() {
		if !automation.validateSubmission(submissionID) {
			automation.handleSubmissionViolation(submissionID)
		}
	}
}

// validateSubmission checks if a submission meets developer guideline compliance
func (automation *DeveloperGuidelineEnforcementAutomation) validateSubmission(submissionID string) bool {
	compliant, err := automation.devManager.ValidateCodingStandards(submissionID)
	if err != nil || !compliant {
		automation.logGuidelineAction(submissionID, "Coding Standards Violation")
		return false
	}

	performanceScore := automation.devManager.GetPerformanceScore(submissionID)
	if performanceScore < PerformanceBenchmarkThreshold {
		automation.logGuidelineAction(submissionID, "Performance Benchmark Violation")
		return false
	}

	securityCompliant, err := automation.devManager.ValidateSecurityCompliance(submissionID)
	if err != nil || !securityCompliant {
		automation.logGuidelineAction(submissionID, "Security Compliance Violation")
		return false
	}

	// Consensus validation for additional enforcement
	if err := automation.consensusEngine.ValidateSubmission(submissionID); err != nil {
		automation.logGuidelineAction(submissionID, "Consensus Validation Failure")
		return false
	}

	return true
}

// handleSubmissionViolation restricts developers with repeated submission violations
func (automation *DeveloperGuidelineEnforcementAutomation) handleSubmissionViolation(submissionID string) {
	developerID := automation.devManager.GetDeveloperID(submissionID)
	automation.violationCount[developerID]++

	if automation.violationCount[developerID] >= MaxSubmissionViolations {
		err := automation.devManager.RestrictDeveloper(developerID)
		if err != nil {
			fmt.Printf("Failed to restrict developer %s due to repeated guideline violations: %v\n", developerID, err)
			automation.logGuidelineAction(developerID, "Failed Developer Restriction")
		} else {
			fmt.Printf("Developer %s restricted due to repeated guideline violations.\n", developerID)
			automation.logGuidelineAction(developerID, "Developer Restricted for Guideline Violations")
			automation.violationCount[developerID] = 0
		}
	} else {
		fmt.Printf("Guideline violation detected for submission %s by developer %s.\n", submissionID, developerID)
		automation.logGuidelineAction(submissionID, "Guideline Compliance Violation Detected")
	}
}

// logGuidelineAction securely logs actions related to developer guideline enforcement
func (automation *DeveloperGuidelineEnforcementAutomation) logGuidelineAction(entityID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Entity ID: %s", action, entityID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("developer-guideline-enforcement-%s-%d", entityID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Developer Guideline Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log developer guideline enforcement action for entity %s: %v\n", entityID, err)
	} else {
		fmt.Println("Developer guideline enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DeveloperGuidelineEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualGuidelineCheck allows administrators to manually check guideline compliance for a specific submission
func (automation *DeveloperGuidelineEnforcementAutomation) TriggerManualGuidelineCheck(submissionID string) {
	fmt.Printf("Manually triggering guideline compliance check for submission: %s\n", submissionID)

	if !automation.validateSubmission(submissionID) {
		automation.handleSubmissionViolation(submissionID)
	} else {
		fmt.Printf("Submission %s is compliant with developer guidelines.\n", submissionID)
		automation.logGuidelineAction(submissionID, "Manual Guideline Compliance Check Passed")
	}
}
