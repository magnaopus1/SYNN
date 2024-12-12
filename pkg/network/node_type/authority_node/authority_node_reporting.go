package authority_node

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/syn900"
)

// ReportingThresholds defines the report thresholds for penalties.
const (
	Syn900Report1Threshold        = 100000  // 100,000 syn900 reports within 30 days: 50% reward slash for 30 days
	Syn900Report2Threshold        = 500000  // 500,000 syn900 reports within 30 days: 80% reward slash for 60 days
	Syn900Report3Threshold        = 1000000 // 1,000,000 syn900 reports within 30 days: Revocation of status pending review

	AuthorityReport1Threshold     = 3  // 3 authority node reports within 30 days: 50% reward slash for 30 days
	AuthorityReport2Threshold     = 7  // 7 authority node reports within 30 days: 80% reward slash for 60 days
	AuthorityReport3Threshold     = 12 // 12 authority node reports within 30 days: Immediate revocation of authority node status
)

// AuthorityNodeReport represents a report filed against an authority node.
type AuthorityNodeReport struct {
	NodeID           string    // Node ID being reported
	ReportType       string    // "syn900" or "authority" report
	ReporterID       string    // Reporter ID (syn900 token or authority node ID)
	ReportTimestamp  time.Time // Timestamp when the report was filed
}

// AuthorityNodeReportManager manages the reporting system for authority nodes.
type AuthorityNodeReportManager struct {
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	Ledger            *ledger.Ledger                // Reference to the ledger for storing reports
	EncryptionService *encryption.Encryption        // Encryption service for securing reports
	Reports           map[string][]*AuthorityNodeReport // Map of reports against authority nodes by NodeID
}

// NewAuthorityNodeReportManager initializes a new AuthorityNodeReportManager.
func NewAuthorityNodeReportManager(ledger *ledger.Ledger, encryptionService *encryption.Encryption) *AuthorityNodeReportManager {
	return &AuthorityNodeReportManager{
		Ledger:            ledger,
		EncryptionService: encryptionService,
		Reports:           make(map[string][]*AuthorityNodeReport),
	}
}

// FileSyn900Report files a syn900 token report against an authority node.
func (rm *AuthorityNodeReportManager) FileSyn900Report(nodeID, reporterID string) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if nodeID == "" || reporterID == "" {
		return errors.New("invalid nodeID or reporterID")
	}

	report := &AuthorityNodeReport{
		NodeID:          nodeID,
		ReportType:      "syn900",
		ReporterID:      reporterID,
		ReportTimestamp: time.Now(),
	}

	// Add the report to the node's report list.
	rm.Reports[nodeID] = append(rm.Reports[nodeID], report)

	// Encrypt and store the report in the ledger.
	encryptedReport, err := rm.EncryptionService.EncryptData([]byte(fmt.Sprintf("%v", report)), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt report: %v", err)
	}

	err = rm.Ledger.RecordNodeReport(nodeID, encryptedReport)
	if err != nil {
		return fmt.Errorf("failed to record report in ledger: %v", err)
	}

	// Check for penalties.
	return rm.evaluatePenalties(nodeID, "syn900")
}

// FileAuthorityNodeReport files an authority node report against another authority node.
func (rm *AuthorityNodeReportManager) FileAuthorityNodeReport(nodeID, reporterID string) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if nodeID == "" || reporterID == "" {
		return errors.New("invalid nodeID or reporterID")
	}

	report := &AuthorityNodeReport{
		NodeID:          nodeID,
		ReportType:      "authority",
		ReporterID:      reporterID,
		ReportTimestamp: time.Now(),
	}

	// Add the report to the node's report list.
	rm.Reports[nodeID] = append(rm.Reports[nodeID], report)

	// Encrypt and store the report in the ledger.
	encryptedReport, err := rm.EncryptionService.EncryptData([]byte(fmt.Sprintf("%v", report)), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt report: %v", err)
	}

	err = rm.Ledger.RecordNodeReport(nodeID, encryptedReport)
	if err != nil {
		return fmt.Errorf("failed to record report in ledger: %v", err)
	}

	// Check for penalties.
	return rm.evaluatePenalties(nodeID, "authority")
}

// evaluatePenalties evaluates penalties for an authority node based on the reports.
func (rm *AuthorityNodeReportManager) evaluatePenalties(nodeID, reportType string) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Count reports within the last 30 days.
	reportCount := 0
	now := time.Now()
	for _, report := range rm.Reports[nodeID] {
		if report.ReportType == reportType && now.Sub(report.ReportTimestamp) <= 30*24*time.Hour {
			reportCount++
		}
	}

	// Evaluate penalties for syn900 reports.
	if reportType == "syn900" {
		switch {
		case reportCount >= Syn900Report3Threshold:
			err := rm.revokeNodeStatusPendingReview(nodeID)
			if err != nil {
				return fmt.Errorf("failed to revoke node status pending review: %v", err)
			}
		case reportCount >= Syn900Report2Threshold:
			err := rm.applyRewardSlash(nodeID, 80, 60)
			if err != nil {
				return fmt.Errorf("failed to apply 80%% reward slash: %v", err)
			}
		case reportCount >= Syn900Report1Threshold:
			err := rm.applyRewardSlash(nodeID, 50, 30)
			if err != nil {
				return fmt.Errorf("failed to apply 50%% reward slash: %v", err)
			}
		}
	}

	// Evaluate penalties for authority node reports.
	if reportType == "authority" {
		switch {
		case reportCount >= AuthorityReport3Threshold:
			err := rm.revokeNodeStatusImmediately(nodeID)
			if err != nil {
				return fmt.Errorf("failed to revoke node status: %v", err)
			}
		case reportCount >= AuthorityReport2Threshold:
			err := rm.applyRewardSlash(nodeID, 80, 60)
			if err != nil {
				return fmt.Errorf("failed to apply 80%% reward slash: %v", err)
			}
		case reportCount >= AuthorityReport1Threshold:
			err := rm.applyRewardSlash(nodeID, 50, 30)
			if err != nil {
				return fmt.Errorf("failed to apply 50%% reward slash: %v", err)
			}
		}
	}

	return nil
}

// applyRewardSlash applies a reward slash to an authority node.
func (rm *AuthorityNodeReportManager) applyRewardSlash(nodeID string, percentage int, durationDays int) error {
	// Apply the reward slash in the ledger.
	err := rm.Ledger.ApplyRewardSlash(nodeID, percentage, durationDays)
	if err != nil {
		return fmt.Errorf("failed to apply reward slash for node %s: %v", nodeID, err)
	}

	fmt.Printf("Applied %d%% reward slash to node %s for %d days.\n", percentage, nodeID, durationDays)
	return nil
}

// revokeNodeStatusPendingReview revokes the status of a node pending review by other authority nodes.
func (rm *AuthorityNodeReportManager) revokeNodeStatusPendingReview(nodeID string) error {
	// Mark the node status as revoked pending review in the ledger.
	err := rm.Ledger.RevokeNodeStatusPendingReview(nodeID)
	if err != nil {
		return fmt.Errorf("failed to revoke node status pending review for node %s: %v", nodeID, err)
	}

	// Trigger a review process that requires 7 authority node confirmations or rejections.
	err = rm.triggerReviewProcess(nodeID)
	if err != nil {
		return fmt.Errorf("failed to trigger review process for node %s: %v", nodeID, err)
	}

	fmt.Printf("Node %s status revoked pending review.\n", nodeID)
	return nil
}

// revokeNodeStatusImmediately revokes the status of a node immediately.
func (rm *AuthorityNodeReportManager) revokeNodeStatusImmediately(nodeID string) error {
	// Mark the node status as immediately revoked in the ledger.
	err := rm.Ledger.RevokeNodeStatus(nodeID)
	if err != nil {
		return fmt.Errorf("failed to revoke node status for node %s: %v", nodeID, err)
	}

	fmt.Printf("Node %s status immediately revoked.\n", nodeID)
	return nil
}

// triggerReviewProcess triggers a review process for a node with 7 authority node confirmations or rejections.
func (rm *AuthorityNodeReportManager) triggerReviewProcess(nodeID string) error {
	// Logic to randomly select 7 authority nodes for the review process.
	selectedNodes, err := rm.selectRandomAuthorityNodes(7)
	if err != nil {
		return fmt.Errorf("failed to select authority nodes for review process: %v", err)
	}

	for _, authorityNode := range selectedNodes {
		// Send the node review request to selected authority nodes.
		err = rm.NetworkManager.RequestNodeReview(authorityNode, nodeID)
		if err != nil {
			return fmt.Errorf("failed to send review request to node %s: %v", authorityNode, err)
		}
	}

	return nil
}

// selectRandomAuthorityNodes randomly selects a given number of authority nodes for review.
func (rm *AuthorityNodeReportManager) selectRandomAuthorityNodes(count int) ([]string, error) {
	authorityNodes := rm.NetworkManager.GetActiveAuthorityNodes()
	if len(authorityNodes) < count {
		return nil, errors.New("not enough authority nodes available")
	}

	selectedNodes := []string{}
	for len(selectedNodes) < count {
		randomNode := authorityNodes[common.RandomInt(len(authorityNodes))]
		if !contains(selectedNodes, randomNode) {
			selectedNodes = append(selectedNodes, randomNode)
		}
	}

	return selectedNodes, nil
}

// Helper function to check if an item is in a slice.
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
