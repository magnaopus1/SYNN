package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/encryption"
)

const (
    VMRunsCheckInterval = 3000 * time.Millisecond // Interval for checking if VMs are running on every node
    SubBlocksPerBlock   = 1000                    // Number of sub-blocks in a block
)

// VirtualMachineRunsOnEveryNodeAutomation ensures that the virtual machine runs on every node in the network
type VirtualMachineRunsOnEveryNodeAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance      *ledger.Ledger               // Ledger for recording events
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    vmRunsCheckCycle    int                          // Counter for VM runs check cycles
}

// NewVirtualMachineRunsOnEveryNodeAutomation initializes the automation that checks if VMs run on every node
func NewVirtualMachineRunsOnEveryNodeAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VirtualMachineRunsOnEveryNodeAutomation {
    return &VirtualMachineRunsOnEveryNodeAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        vmRunsCheckCycle:  0,
    }
}

// StartVMRunsCheck starts the continuous loop that monitors whether virtual machines are running on every node
func (automation *VirtualMachineRunsOnEveryNodeAutomation) StartVMRunsCheck() {
    ticker := time.NewTicker(VMRunsCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkVMRunsOnAllNodes()
        }
    }()
}

// checkVMRunsOnAllNodes checks if virtual machines are running on all the network nodes
func (automation *VirtualMachineRunsOnEveryNodeAutomation) checkVMRunsOnAllNodes() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Get the list of all active nodes in the network
    activeNodes := automation.consensusSystem.GetActiveNodes()

    for _, node := range activeNodes {
        if !automation.isVMRunningOnNode(node) {
            fmt.Printf("VM is not running on node %s. Initiating VM deployment.\n", node.ID)
            automation.deployVMOnNode(node)
        } else {
            fmt.Printf("VM is running on node %s.\n", node.ID)
        }
    }

    automation.vmRunsCheckCycle++
    fmt.Printf("VM runs check cycle #%d executed.\n", automation.vmRunsCheckCycle)

    if automation.vmRunsCheckCycle%SubBlocksPerBlock == 0 {
        automation.finalizeVMRunsCheckCycle()
    }
}

// isVMRunningOnNode checks if the virtual machine is running on the specified node
func (automation *VirtualMachineRunsOnEveryNodeAutomation) isVMRunningOnNode(node common.Node) bool {
    return automation.consensusSystem.IsVMRunningOnNode(node)
}

// deployVMOnNode triggers the deployment of a virtual machine on the specified node
func (automation *VirtualMachineRunsOnEveryNodeAutomation) deployVMOnNode(node common.Node) {
    // Encrypt the node data before deploying the VM
    encryptedNodeData := automation.encryptNodeData(node)

    // Trigger VM deployment through the Synnergy Consensus system
    deploymentSuccess := automation.consensusSystem.DeployVMOnNode(encryptedNodeData)

    if deploymentSuccess {
        fmt.Printf("Virtual machine successfully deployed on node %s.\n", node.ID)
        automation.logVMDeploymentEvent(node)
    } else {
        fmt.Printf("Error deploying virtual machine on node %s.\n", node.ID)
    }
}

// finalizeVMRunsCheckCycle finalizes the VM runs check cycle and logs the result in the ledger
func (automation *VirtualMachineRunsOnEveryNodeAutomation) finalizeVMRunsCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeVMRunsCheckCycle()
    if success {
        fmt.Println("VM runs check cycle finalized successfully.")
        automation.logVMRunsCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing VM runs check cycle.")
    }
}

// logVMDeploymentEvent logs the VM deployment event for a specific node into the ledger
func (automation *VirtualMachineRunsOnEveryNodeAutomation) logVMDeploymentEvent(node common.Node) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-deployment-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Deployment",
        Status:    "Completed",
        Details:   fmt.Sprintf("Virtual machine successfully deployed on node %s.", node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with VM deployment event for node %s.\n", node.ID)
}

// logVMRunsCheckCycleFinalization logs the finalization of a VM runs check cycle into the ledger
func (automation *VirtualMachineRunsOnEveryNodeAutomation) logVMRunsCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-runs-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "VM Runs Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with VM runs check cycle finalization.")
}

// encryptNodeData encrypts the node data before initiating VM deployment
func (automation *VirtualMachineRunsOnEveryNodeAutomation) encryptNodeData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node data successfully encrypted.")
    return node
}

// ensureVMIntegrity checks the integrity of the VM deployment and operation and triggers redeployment if necessary
func (automation *VirtualMachineRunsOnEveryNodeAutomation) ensureVMIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateVMIntegrity()
    if !integrityValid {
        fmt.Println("VM integrity breach detected. Re-triggering VM deployment checks.")
        automation.checkVMRunsOnAllNodes()
    } else {
        fmt.Println("VM integrity is valid.")
    }
}
