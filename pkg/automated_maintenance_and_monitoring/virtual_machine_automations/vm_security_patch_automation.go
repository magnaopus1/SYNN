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
    SecurityPatchCheckInterval = 6000 * time.Millisecond // Interval for checking security patches
    SubBlocksPerBlock          = 1000                    // Number of sub-blocks in a block
)

// VMSecurityPatchAutomation automates the process of applying security patches to VMs
type VMSecurityPatchAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging security patches
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    patchCheckCount    int                          // Counter for patch check cycles
}

// NewVMSecurityPatchAutomation initializes the automation for applying security patches to VMs
func NewVMSecurityPatchAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMSecurityPatchAutomation {
    return &VMSecurityPatchAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        patchCheckCount: 0,
    }
}

// StartPatchCheck starts the continuous loop for monitoring and applying security patches to VMs
func (automation *VMSecurityPatchAutomation) StartPatchCheck() {
    ticker := time.NewTicker(SecurityPatchCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndApplyPatches()
        }
    }()
}

// monitorAndApplyPatches checks for available security patches and applies them if necessary
func (automation *VMSecurityPatchAutomation) monitorAndApplyPatches() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the available security patches from the system
    availablePatches := automation.consensusSystem.GetAvailableSecurityPatches()

    if len(availablePatches) > 0 {
        for _, patch := range availablePatches {
            fmt.Printf("Applying security patch %s to VM %s.\n", patch.ID, patch.VM.ID)
            automation.applySecurityPatch(patch.VM, patch)
        }
    } else {
        fmt.Println("No new security patches available for VMs.")
    }

    automation.patchCheckCount++
    fmt.Printf("Security patch check cycle #%d executed.\n", automation.patchCheckCount)

    if automation.patchCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizePatchCycle()
    }
}

// applySecurityPatch applies the security patch to the specified VM
func (automation *VMSecurityPatchAutomation) applySecurityPatch(vm common.VirtualMachine, patch common.SecurityPatch) {
    // Encrypt patch data before applying it
    encryptedPatch := automation.encryptPatchData(patch)

    // Apply the patch through the Synnergy Consensus system
    patchSuccess := automation.consensusSystem.ApplySecurityPatch(vm, encryptedPatch)

    if patchSuccess {
        fmt.Printf("Security patch %s successfully applied to VM %s.\n", patch.ID, vm.ID)
        automation.logPatchEvent(vm, patch)
    } else {
        fmt.Printf("Error applying security patch %s to VM %s.\n", patch.ID, vm.ID)
    }
}

// finalizePatchCycle finalizes the patch check cycle and logs the result in the ledger
func (automation *VMSecurityPatchAutomation) finalizePatchCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizePatchCycle()
    if success {
        fmt.Println("Security patch check cycle finalized successfully.")
        automation.logPatchCycleFinalization()
    } else {
        fmt.Println("Error finalizing security patch check cycle.")
    }
}

// logPatchEvent logs the patch action for a specific VM into the ledger
func (automation *VMSecurityPatchAutomation) logPatchEvent(vm common.VirtualMachine, patch common.SecurityPatch) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("security-patch-%s-%s", patch.ID, vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Security Patch",
        Status:    "Applied",
        Details:   fmt.Sprintf("Security patch %s successfully applied to VM %s.", patch.ID, vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with security patch event for VM %s.\n", vm.ID)
}

// logPatchCycleFinalization logs the finalization of a patch check cycle into the ledger
func (automation *VMSecurityPatchAutomation) logPatchCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("security-patch-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Security Patch Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with security patch cycle finalization.")
}

// encryptPatchData encrypts the patch data before applying it
func (automation *VMSecurityPatchAutomation) encryptPatchData(patch common.SecurityPatch) common.SecurityPatch {
    encryptedData, err := encryption.EncryptData(patch)
    if err != nil {
        fmt.Println("Error encrypting patch data:", err)
        return patch
    }

    patch.EncryptedData = encryptedData
    fmt.Println("Security patch data successfully encrypted.")
    return patch
}

// ensurePatchIntegrity checks the integrity of the patch data and re-applies the patch if necessary
func (automation *VMSecurityPatchAutomation) ensurePatchIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidatePatchIntegrity()
    if !integrityValid {
        fmt.Println("Security patch integrity breach detected. Re-triggering patch application.")
        automation.monitorAndApplyPatches()
    } else {
        fmt.Println("Security patch integrity is valid.")
    }
}
