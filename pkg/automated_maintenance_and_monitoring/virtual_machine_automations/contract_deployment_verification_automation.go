package automations

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

const (
    DeploymentCheckInterval = 3000 * time.Millisecond // Interval for checking contract deployments
    SubBlocksPerBlock       = 1000                    // Number of sub-blocks in a block
)

// ContractDeploymentVerificationAutomation automates the process of verifying contract deployments
type ContractDeploymentVerificationAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus
    ledgerInstance       *ledger.Ledger
    stateMutex           *sync.RWMutex
    deploymentCheckCount int
}

// NewContractDeploymentVerificationAutomation initializes the automation for contract deployment verification
func NewContractDeploymentVerificationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractDeploymentVerificationAutomation {
    return &ContractDeploymentVerificationAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        deploymentCheckCount: 0,
    }
}

// StartDeploymentCheck starts the continuous loop for verifying contract deployments
func (automation *ContractDeploymentVerificationAutomation) StartDeploymentCheck() {
    ticker := time.NewTicker(DeploymentCheckInterval)

    go func() {
        for range ticker.C {
            automation.verifyAndLogContractDeployment()
        }
    }()
}

// verifyAndLogContractDeployment verifies contract deployments and logs the results
func (automation *ContractDeploymentVerificationAutomation) verifyAndLogContractDeployment() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch pending deployments that need to be verified
    pendingDeployments := automation.consensusSystem.GetPendingContractDeployments()

    for _, deployment := range pendingDeployments {
        fmt.Printf("Verifying contract deployment: %s\n", deployment.ContractID)

        verified := automation.verifyContract(deployment)
        if verified {
            fmt.Printf("Contract %s successfully verified.\n", deployment.ContractID)
            automation.logContractDeployment(deployment)
        } else {
            fmt.Printf("Contract %s failed verification.\n", deployment.ContractID)
        }
    }

    automation.deploymentCheckCount++
    if automation.deploymentCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeDeploymentCycle()
    }
}

// verifyContract verifies a specific contract deployment through encryption and consensus
func (automation *ContractDeploymentVerificationAutomation) verifyContract(deployment common.ContractDeployment) bool {
    // Encrypt contract data
    encryptedData := encryption.EncryptContractData(deployment)

    // Verify the contract using the Synnergy Consensus
    return automation.consensusSystem.VerifyContractDeployment(encryptedData)
}

// logContractDeployment logs the contract deployment verification in the ledger
func (automation *ContractDeploymentVerificationAutomation) logContractDeployment(deployment common.ContractDeployment) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("contract-deployment-%s", deployment.ContractID),
        Timestamp: time.Now().Unix(),
        Type:      "Contract Deployment",
        Status:    "Verified",
        Details:   fmt.Sprintf("Contract %s successfully verified and deployed.", deployment.ContractID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with contract deployment verification for contract %s.\n", deployment.ContractID)
}

// finalizeDeploymentCycle finalizes the deployment verification cycle and logs it in the ledger
func (automation *ContractDeploymentVerificationAutomation) finalizeDeploymentCycle() {
    success := automation.consensusSystem.FinalizeDeploymentCycle()
    if success {
        fmt.Println("Deployment verification cycle finalized successfully.")
        automation.logDeploymentCycleFinalization()
    } else {
        fmt.Println("Error finalizing deployment verification cycle.")
    }
}

// logDeploymentCycleFinalization logs the finalization of the deployment verification cycle
func (automation *ContractDeploymentVerificationAutomation) logDeploymentCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("deployment-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Deployment Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with deployment cycle finalization.")
}
