package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/security"
)

const (
	AttackCheckInterval    = 30 * time.Second // Interval to check for potential attacks
	BlacklistingThreshold  = 5                // Number of suspicious activities before blacklisting
	RateLimitThreshold     = 100              // Number of transactions from one node before applying rate limiting
)

// AttackPreventionEnforcementAutomation handles detecting and preventing attacks in the blockchain network
type AttackPreventionEnforcementAutomation struct {
	securityManager    *security.SecurityManager
	consensusEngine    *consensus.SynnergyConsensus
	ledgerInstance     *ledger.Ledger
	enforcementMutex   *sync.RWMutex
	blacklistedNodes   map[string]bool
	nodeTransactionMap map[string]int
}

// NewAttackPreventionEnforcementAutomation initializes attack prevention enforcement automation
func NewAttackPreventionEnforcementAutomation(securityManager *security.SecurityManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *AttackPreventionEnforcementAutomation {
	return &AttackPreventionEnforcementAutomation{
		securityManager:    securityManager,
		consensusEngine:    consensusEngine,
		ledgerInstance:     ledgerInstance,
		enforcementMutex:   enforcementMutex,
		blacklistedNodes:   make(map[string]bool),
		nodeTransactionMap: make(map[string]int),
	}
}

// StartAttackPreventionEnforcement begins continuous monitoring of the network for attacks
func (automation *AttackPreventionEnforcementAutomation) StartAttackPreventionEnforcement() {
	ticker := time.NewTicker(AttackCheckInterval)

	go func() {
		for range ticker.C {
			automation.detectAndPreventAttacks()
		}
	}()
}

// detectAndPreventAttacks identifies suspicious activities and prevents further attacks
func (automation *AttackPreventionEnforcementAutomation) detectAndPreventAttacks() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	// Retrieve the list of nodes with suspicious activity from the security manager
	suspiciousNodes := automation.securityManager.GetSuspiciousNodes()

	for _, nodeID := range suspiciousNodes {
		if automation.blacklistedNodes[nodeID] {
			// Node is already blacklisted, enforce prevention
			automation.preventAttackFromBlacklistedNode(nodeID)
		} else {
			// Check transaction count and decide if rate limiting or blacklisting is needed
			automation.nodeTransactionMap[nodeID]++

			if automation.nodeTransactionMap[nodeID] >= RateLimitThreshold {
				automation.applyRateLimiting(nodeID)
			}

			if automation.nodeTransactionMap[nodeID] >= BlacklistingThreshold {
				automation.blacklistNode(nodeID)
			}
		}
	}
}

// blacklistNode blacklists a node from participating in the network
func (automation *AttackPreventionEnforcementAutomation) blacklistNode(nodeID string) {
	automation.blacklistedNodes[nodeID] = true
	automation.securityManager.BlockNode(nodeID)

	fmt.Printf("Node %s has been blacklisted due to suspicious activities.\n", nodeID)
	automation.logAttackPreventionAction(nodeID, "Blacklisted")
}

// applyRateLimiting applies rate limiting to nodes exhibiting suspicious behavior
func (automation *AttackPreventionEnforcementAutomation) applyRateLimiting(nodeID string) {
	err := automation.securityManager.RateLimitNode(nodeID)
	if err != nil {
		fmt.Printf("Failed to apply rate limiting to node %s: %v\n", nodeID, err)
		return
	}

	fmt.Printf("Rate limiting applied to node %s.\n", nodeID)
	automation.logAttackPreventionAction(nodeID, "Rate Limited")
}

// preventAttackFromBlacklistedNode ensures no transactions are processed from blacklisted nodes
func (automation *AttackPreventionEnforcementAutomation) preventAttackFromBlacklistedNode(nodeID string) {
	err := automation.consensusEngine.BlockTransactionsFromNode(nodeID)
	if err != nil {
		fmt.Printf("Failed to block transactions from blacklisted node %s: %v\n", nodeID, err)
		return
	}

	fmt.Printf("Blocked all transactions from blacklisted node %s.\n", nodeID)
	automation.logAttackPreventionAction(nodeID, "Transactions Blocked")
}

// logAttackPreventionAction securely logs actions taken to prevent attacks
func (automation *AttackPreventionEnforcementAutomation) logAttackPreventionAction(nodeID string, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Node: %s", action, nodeID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("attack-prevention-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Attack Prevention",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log attack prevention action for node %s: %v\n", nodeID, err)
	} else {
		fmt.Println("Attack prevention action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *AttackPreventionEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualBlacklisting allows administrators to manually blacklist a node
func (automation *AttackPreventionEnforcementAutomation) TriggerManualBlacklisting(nodeID string) {
	fmt.Printf("Manually blacklisting node: %s\n", nodeID)

	automation.blacklistedNodes[nodeID] = true
	automation.securityManager.BlockNode(nodeID)

	automation.logAttackPreventionAction(nodeID, "Manually Blacklisted")
}
