package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    NodeTypeDifferentiationCheckInterval = 3000 * time.Millisecond // Interval for checking node type roles and enforcement
    SubBlocksPerBlock                    = 1000                    // Number of sub-blocks in a block
)

// EnforceNodeTypeDifferentiationAutomation automates the process of enforcing node type differentiation across the network
type EnforceNodeTypeDifferentiationAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger to store node type differentiation actions
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    nodeTypeCheckCount    int                          // Counter for node type differentiation check cycles
}

// NewEnforceNodeTypeDifferentiationAutomation initializes the automation for enforcing node type differentiation across nodes
func NewEnforceNodeTypeDifferentiationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *EnforceNodeTypeDifferentiationAutomation {
    return &EnforceNodeTypeDifferentiationAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        nodeTypeCheckCount: 0,
    }
}

// StartNodeTypeDifferentiationCheck starts the continuous loop for monitoring and enforcing node type differentiation across the network
func (automation *EnforceNodeTypeDifferentiationAutomation) StartNodeTypeDifferentiationCheck() {
    ticker := time.NewTicker(NodeTypeDifferentiationCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceNodeTypeDifferentiation()
        }
    }()
}

// monitorAndEnforceNodeTypeDifferentiation checks the assigned roles of nodes in the network and ensures that node roles are enforced correctly
func (automation *EnforceNodeTypeDifferentiationAutomation) monitorAndEnforceNodeTypeDifferentiation() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of nodes and their roles from Synnergy Consensus
    nodeRoles := automation.consensusSystem.GetNodeRoles()

    for _, nodeRole := range nodeRoles {
        if automation.isNodeRoleValid(nodeRole) {
            fmt.Printf("Node %s is correctly assigned as %s.\n", nodeRole.Node.ID, nodeRole.Role)
        } else {
            fmt.Printf("Node %s has an invalid or misassigned role. Enforcing the correct role.\n", nodeRole.Node.ID)
            automation.enforceCorrectNodeRole(nodeRole)
        }
    }

    automation.nodeTypeCheckCount++
    fmt.Printf("Node type differentiation check cycle #%d executed.\n", automation.nodeTypeCheckCount)

    if automation.nodeTypeCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeNodeTypeDifferentiationCycle()
    }
}

// isNodeRoleValid validates whether a node has the correct assigned role in the network
func (automation *EnforceNodeTypeDifferentiationAutomation) isNodeRoleValid(nodeRole common.NodeRole) bool {
    expectedRole := automation.consensusSystem.GetExpectedRoleForNode(nodeRole.Node)
    return nodeRole.Role == expectedRole
}

// enforceCorrectNodeRole enforces the correct role for nodes that have invalid or misassigned roles
func (automation *EnforceNodeTypeDifferentiationAutomation) enforceCorrectNodeRole(nodeRole common.NodeRole) {
    // Encrypt node data before enforcement
    encryptedNodeData := automation.AddEncryptionToNodeData(nodeRole.Node)

    // Enforce the correct role assignment through the Synnergy Consensus
    enforcementSuccess := automation.consensusSystem.EnforceNodeRole(encryptedNodeData, nodeRole)

    if enforcementSuccess {
        fmt.Printf("Correct role enforced for node %s as %s.\n", nodeRole.Node.ID, nodeRole.Role)
        automation.logNodeRoleEnforcement(nodeRole.Node, nodeRole.Role)
    } else {
        fmt.Printf("Error enforcing role for node %s.\n", nodeRole.Node.ID)
    }
}

// finalizeNodeTypeDifferentiationCycle finalizes the node type differentiation check cycle and logs the result in the ledger
func (automation *EnforceNodeTypeDifferentiationAutomation) finalizeNodeTypeDifferentiationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeNodeRoleEnforcementCycle()
    if success {
        fmt.Println("Node type differentiation check cycle finalized successfully.")
        automation.logNodeTypeDifferentiationCycleFinalization()
    } else {
        fmt.Println("Error finalizing node type differentiation check cycle.")
    }
}

// logNodeRoleEnforcement logs the enforcement of the correct role for a specific node into the ledger for traceability
func (automation *EnforceNodeTypeDifferentiationAutomation) logNodeRoleEnforcement(node common.Node, role string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("node-role-enforcement-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Node Role Enforcement",
        Status:    "Enforced",
        Details:   fmt.Sprintf("Correct role %s enforced for node %s.", role, node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with node role enforcement event for node %s.\n", node.ID)
}

// logNodeTypeDifferentiationCycleFinalization logs the finalization of a node type differentiation check cycle into the ledger
func (automation *EnforceNodeTypeDifferentiationAutomation) logNodeTypeDifferentiationCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("node-type-differentiation-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Node Type Differentiation Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with node type differentiation cycle finalization.")
}

// AddEncryptionToNodeData encrypts the node data before enforcing node roles
func (automation *EnforceNodeTypeDifferentiationAutomation) AddEncryptionToNodeData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node data successfully encrypted.")
    return node
}

// ensureNodeTypeDifferentiationIntegrity checks the integrity of node role assignments and triggers enforcement if necessary
func (automation *EnforceNodeTypeDifferentiationAutomation) ensureNodeTypeDifferentiationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateNodeRoleAssignments()
    if !integrityValid {
        fmt.Println("Node type differentiation integrity breach detected. Re-triggering role enforcement.")
        automation.monitorAndEnforceNodeTypeDifferentiation()
    } else {
        fmt.Println("Node role assignment integrity is valid.")
    }
}
