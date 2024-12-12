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
    VMProvisioningCheckInterval = 5000 * time.Millisecond // Interval for checking VM provisioning requirements
    SubBlocksPerBlock           = 1000                    // Number of sub-blocks in a block
)

// VirtualMachineProvisioningAutomation handles the automation for provisioning new virtual machines as required
type VirtualMachineProvisioningAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance        *ledger.Ledger               // Ledger for recording provisioning actions
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    provisioningCycleCount int                         // Counter for provisioning check cycles
}

// NewVirtualMachineProvisioningAutomation initializes the automation for virtual machine provisioning
func NewVirtualMachineProvisioningAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VirtualMachineProvisioningAutomation {
    return &VirtualMachineProvisioningAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        provisioningCycleCount: 0,
    }
}

// StartProvisioningCheck starts the continuous loop for monitoring and provisioning virtual machines
func (automation *VirtualMachineProvisioningAutomation) StartProvisioningCheck() {
    ticker := time.NewTicker(VMProvisioningCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndProvisionVMs()
        }
    }()
}

// checkAndProvisionVMs checks whether new virtual machines are needed and provisions them
func (automation *VirtualMachineProvisioningAutomation) checkAndProvisionVMs() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch VM resource requirements from the consensus system
    vmRequirements := automation.consensusSystem.GetVMProvisioningRequirements()

    if len(vmRequirements) > 0 {
        for _, req := range vmRequirements {
            fmt.Printf("Provisioning new virtual machine for requirement: %s.\n", req.Description)
            automation.provisionNewVM(req)
        }
    } else {
        fmt.Println("No new virtual machine provisioning required at the moment.")
    }

    automation.provisioningCycleCount++
    fmt.Printf("Provisioning check cycle #%d executed.\n", automation.provisioningCycleCount)

    if automation.provisioningCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeProvisioningCycle()
    }
}

// provisionNewVM provisions a new virtual machine based on the given requirement
func (automation *VirtualMachineProvisioningAutomation) provisionNewVM(requirement common.VMProvisioningRequirement) {
    // Encrypt the VM provisioning details
    encryptedProvisioningDetails := automation.encryptProvisioningDetails(requirement)

    // Trigger VM provisioning through Synnergy Consensus
    provisioningSuccess := automation.consensusSystem.TriggerVMProvisioning(encryptedProvisioningDetails)

    if provisioningSuccess {
        fmt.Printf("New virtual machine successfully provisioned for requirement: %s.\n", requirement.Description)
        automation.logProvisioningEvent(requirement)
    } else {
        fmt.Printf("Error provisioning virtual machine for requirement: %s.\n", requirement.Description)
    }
}

// finalizeProvisioningCycle finalizes the provisioning check cycle and logs the result in the ledger
func (automation *VirtualMachineProvisioningAutomation) finalizeProvisioningCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeProvisioningCycle()
    if success {
        fmt.Println("Provisioning check cycle finalized successfully.")
        automation.logProvisioningCycleFinalization()
    } else {
        fmt.Println("Error finalizing provisioning check cycle.")
    }
}

// logProvisioningEvent logs the provisioning event for a specific virtual machine into the ledger
func (automation *VirtualMachineProvisioningAutomation) logProvisioningEvent(requirement common.VMProvisioningRequirement) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-provisioning-%s", requirement.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Provisioning",
        Status:    "Completed",
        Details:   fmt.Sprintf("Virtual machine successfully provisioned for requirement: %s.", requirement.Description),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with VM provisioning event for requirement: %s.\n", requirement.Description)
}

// logProvisioningCycleFinalization logs the finalization of a provisioning cycle into the ledger
func (automation *VirtualMachineProvisioningAutomation) logProvisioningCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-provisioning-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "VM Provisioning Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with virtual machine provisioning cycle finalization.")
}

// encryptProvisioningDetails encrypts the virtual machine provisioning details before the provisioning process
func (automation *VirtualMachineProvisioningAutomation) encryptProvisioningDetails(requirement common.VMProvisioningRequirement) common.VMProvisioningRequirement {
    encryptedData, err := encryption.EncryptData(requirement)
    if err != nil {
        fmt.Println("Error encrypting VM provisioning data:", err)
        return requirement
    }
    requirement.EncryptedData = encryptedData
    fmt.Println("Virtual machine provisioning data successfully encrypted.")
    return requirement
}

// ensureProvisioningIntegrity checks the integrity of the provisioning process and triggers a recheck if necessary
func (automation *VirtualMachineProvisioningAutomation) ensureProvisioningIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateProvisioningIntegrity()
    if !integrityValid {
        fmt.Println("Provisioning integrity breach detected. Re-triggering provisioning checks.")
        automation.checkAndProvisionVMs()
    } else {
        fmt.Println("Provisioning integrity is valid.")
    }
}
