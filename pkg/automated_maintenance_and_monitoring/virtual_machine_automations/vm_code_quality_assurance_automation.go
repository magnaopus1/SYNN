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
    CodeQualityCheckInterval = 3000 * time.Millisecond // Interval for checking code quality across VMs
    SubBlocksPerBlock        = 1000                    // Number of sub-blocks in a block
)

// VMCodeQualityAssuranceAutomation automates the process of validating and maintaining code quality in VM executions
type VMCodeQualityAssuranceAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance    *ledger.Ledger               // Ledger to record code quality checks
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    qualityCheckCycle int                          // Counter for quality assurance cycles
}

// NewVMCodeQualityAssuranceAutomation initializes the automation for VM code quality assurance
func NewVMCodeQualityAssuranceAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMCodeQualityAssuranceAutomation {
    return &VMCodeQualityAssuranceAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        qualityCheckCycle: 0,
    }
}

// StartCodeQualityCheck starts the continuous loop to monitor and validate code quality in VMs
func (automation *VMCodeQualityAssuranceAutomation) StartCodeQualityCheck() {
    ticker := time.NewTicker(CodeQualityCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndValidateCodeQuality()
        }
    }()
}

// monitorAndValidateCodeQuality checks the code quality across VMs and ensures all code meets quality standards
func (automation *VMCodeQualityAssuranceAutomation) monitorAndValidateCodeQuality() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of active VMs
    vmList := automation.consensusSystem.GetActiveVMs()

    for _, vm := range vmList {
        fmt.Printf("Validating code quality for VM %s.\n", vm.ID)
        automation.validateVMCodeQuality(vm)
    }

    automation.qualityCheckCycle++
    fmt.Printf("Code quality check cycle #%d executed.\n", automation.qualityCheckCycle)

    if automation.qualityCheckCycle%SubBlocksPerBlock == 0 {
        automation.finalizeCodeQualityCheckCycle()
    }
}

// validateVMCodeQuality performs code quality checks for a specific virtual machine
func (automation *VMCodeQualityAssuranceAutomation) validateVMCodeQuality(vm common.VirtualMachine) {
    // Encrypt VM data before checking code quality
    encryptedVMData := automation.encryptVMData(vm)

    // Validate the code quality of the VM through the Synnergy Consensus system
    qualityValid := automation.consensusSystem.ValidateVMCodeQuality(encryptedVMData)

    if qualityValid {
        fmt.Printf("Code quality for VM %s is valid.\n", vm.ID)
        automation.logCodeQualityEvent(vm)
    } else {
        fmt.Printf("Code quality for VM %s failed validation.\n", vm.ID)
        automation.triggerCodeCorrection(vm)
    }
}

// triggerCodeCorrection triggers the necessary actions to correct code quality issues for a VM
func (automation *VMCodeQualityAssuranceAutomation) triggerCodeCorrection(vm common.VirtualMachine) {
    // Initiating automatic code correction
    success := automation.consensusSystem.InitiateCodeCorrection(vm.ID)

    if success {
        fmt.Printf("Code correction successfully initiated for VM %s.\n", vm.ID)
        automation.logCodeCorrectionEvent(vm)
    } else {
        fmt.Printf("Failed to initiate code correction for VM %s.\n", vm.ID)
    }
}

// finalizeCodeQualityCheckCycle finalizes the code quality check cycle and logs the result in the ledger
func (automation *VMCodeQualityAssuranceAutomation) finalizeCodeQualityCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCodeQualityCheckCycle()
    if success {
        fmt.Println("Code quality check cycle finalized successfully.")
        automation.logCodeQualityCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing code quality check cycle.")
    }
}

// logCodeQualityEvent logs the code quality validation for a specific VM into the ledger
func (automation *VMCodeQualityAssuranceAutomation) logCodeQualityEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("code-quality-vm-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Code Quality",
        Status:    "Validated",
        Details:   fmt.Sprintf("Code quality validated for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with code quality validation event for VM %s.\n", vm.ID)
}

// logCodeCorrectionEvent logs the code correction action into the ledger
func (automation *VMCodeQualityAssuranceAutomation) logCodeCorrectionEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("code-correction-vm-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Code Correction",
        Status:    "Triggered",
        Details:   fmt.Sprintf("Code correction successfully triggered for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with code correction event for VM %s.\n", vm.ID)
}

// logCodeQualityCheckCycleFinalization logs the finalization of a code quality check cycle into the ledger
func (automation *VMCodeQualityAssuranceAutomation) logCodeQualityCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("code-quality-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Code Quality Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with code quality check cycle finalization.")
}

// encryptVMData encrypts the VM data before performing code quality checks
func (automation *VMCodeQualityAssuranceAutomation) encryptVMData(vm common.VirtualMachine) common.VirtualMachine {
    encryptedData, err := encryption.EncryptData(vm)
    if err != nil {
        fmt.Println("Error encrypting VM data:", err)
        return vm
    }
    vm.EncryptedData = encryptedData
    fmt.Println("VM data successfully encrypted.")
    return vm
}

// ensureCodeQualityIntegrity checks the integrity of the code quality validation process and triggers re-validation if necessary
func (automation *VMCodeQualityAssuranceAutomation) ensureCodeQualityIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateCodeQualityIntegrity()
    if !integrityValid {
        fmt.Println("Code quality validation integrity breach detected. Re-triggering validation checks.")
        automation.monitorAndValidateCodeQuality()
    } else {
        fmt.Println("Code quality validation integrity is valid.")
    }
}
