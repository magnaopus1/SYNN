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
    VMSyncCheckInterval = 5000 * time.Millisecond // Interval for checking VM synchronization status
    SubBlocksPerBlock   = 1000                    // Number of sub-blocks in a block
)

// VirtualMachineSynchronizationAutomation automates the process of synchronizing virtual machines across nodes
type VirtualMachineSynchronizationAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance     *ledger.Ledger               // Ledger to record synchronization events
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    syncCheckCycle     int                          // Counter for synchronization check cycles
}

// NewVirtualMachineSynchronizationAutomation initializes the automation for VM synchronization
func NewVirtualMachineSynchronizationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VirtualMachineSynchronizationAutomation {
    return &VirtualMachineSynchronizationAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        syncCheckCycle:  0,
    }
}

// StartVMSyncCheck starts the continuous loop to monitor and enforce virtual machine synchronization across nodes
func (automation *VirtualMachineSynchronizationAutomation) StartVMSyncCheck() {
    ticker := time.NewTicker(VMSyncCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndSyncVMs()
        }
    }()
}

// checkAndSyncVMs checks the VM synchronization status across all nodes and triggers synchronization if necessary
func (automation *VirtualMachineSynchronizationAutomation) checkAndSyncVMs() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the synchronization status of all nodes' VMs
    syncStatus := automation.consensusSystem.GetVMSynchronizationStatus()

    outOfSyncNodes := automation.findOutOfSyncNodes(syncStatus)

    if len(outOfSyncNodes) > 0 {
        for _, node := range outOfSyncNodes {
            fmt.Printf("Node %s is out of sync. Triggering VM synchronization.\n", node.ID)
            automation.syncVMForNode(node)
        }
    } else {
        fmt.Println("All nodes are synchronized. No action required.")
    }

    automation.syncCheckCycle++
    fmt.Printf("VM synchronization check cycle #%d executed.\n", automation.syncCheckCycle)

    if automation.syncCheckCycle%SubBlocksPerBlock == 0 {
        automation.finalizeSyncCheckCycle()
    }
}

// findOutOfSyncNodes identifies nodes that are out of sync
func (automation *VirtualMachineSynchronizationAutomation) findOutOfSyncNodes(syncStatus []common.VMSyncStatus) []common.Node {
    var outOfSyncNodes []common.Node

    for _, status := range syncStatus {
        if !status.IsSynchronized {
            outOfSyncNodes = append(outOfSyncNodes, status.Node)
        }
    }

    return outOfSyncNodes
}

// syncVMForNode triggers synchronization of the virtual machine on the specified node
func (automation *VirtualMachineSynchronizationAutomation) syncVMForNode(node common.Node) {
    // Encrypt node data before initiating synchronization
    encryptedNodeData := automation.encryptNodeData(node)

    // Trigger VM synchronization through the Synnergy Consensus system
    syncSuccess := automation.consensusSystem.SynchronizeVMOnNode(encryptedNodeData)

    if syncSuccess {
        fmt.Printf("Virtual machine successfully synchronized on node %s.\n", node.ID)
        automation.logVMSyncEvent(node)
    } else {
        fmt.Printf("Error synchronizing virtual machine on node %s.\n", node.ID)
    }
}

// finalizeSyncCheckCycle finalizes the VM synchronization check cycle and logs the result in the ledger
func (automation *VirtualMachineSynchronizationAutomation) finalizeSyncCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeVMSyncCheckCycle()
    if success {
        fmt.Println("VM synchronization check cycle finalized successfully.")
        automation.logVMSyncCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing VM synchronization check cycle.")
    }
}

// logVMSyncEvent logs the VM synchronization event for a specific node into the ledger
func (automation *VirtualMachineSynchronizationAutomation) logVMSyncEvent(node common.Node) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-sync-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Synchronization",
        Status:    "Completed",
        Details:   fmt.Sprintf("VM successfully synchronized on node %s.", node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with VM synchronization event for node %s.\n", node.ID)
}

// logVMSyncCheckCycleFinalization logs the finalization of a VM synchronization check cycle into the ledger
func (automation *VirtualMachineSynchronizationAutomation) logVMSyncCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-sync-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "VM Sync Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with VM sync check cycle finalization.")
}

// encryptNodeData encrypts the node data before initiating VM synchronization
func (automation *VirtualMachineSynchronizationAutomation) encryptNodeData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node data successfully encrypted.")
    return node
}

// ensureVMSyncIntegrity checks the integrity of the synchronization process and re-triggers it if necessary
func (automation *VirtualMachineSynchronizationAutomation) ensureVMSyncIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateVMSyncIntegrity()
    if !integrityValid {
        fmt.Println("VM synchronization integrity breach detected. Re-triggering synchronization checks.")
        automation.checkAndSyncVMs()
    } else {
        fmt.Println("VM synchronization integrity is valid.")
    }
}
