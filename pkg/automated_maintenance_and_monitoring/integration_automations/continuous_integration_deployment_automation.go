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
    CIDeploymentCheckInterval = 2500 * time.Millisecond // Interval for checking CI/CD process
    SubBlocksPerBlock         = 1000                    // Number of sub-blocks in a block
)

// ContinuousIntegrationDeploymentAutomation automates continuous integration and deployment for new code changes and features
type ContinuousIntegrationDeploymentAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance   *ledger.Ledger               // Ledger to store CI/CD logs and actions
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    deploymentCount  int                          // Counter for deployment cycles
}

// NewContinuousIntegrationDeploymentAutomation initializes the CI/CD automation for deployment
func NewContinuousIntegrationDeploymentAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContinuousIntegrationDeploymentAutomation {
    return &ContinuousIntegrationDeploymentAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        deploymentCount: 0,
    }
}

// StartCIDeployment starts the continuous loop for integration and deployment monitoring
func (automation *ContinuousIntegrationDeploymentAutomation) StartCIDeployment() {
    ticker := time.NewTicker(CIDeploymentCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkForCIDeployment()
        }
    }()
}

// checkForCIDeployment checks if new code changes or features are ready for integration and deployment
func (automation *ContinuousIntegrationDeploymentAutomation) checkForCIDeployment() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    newFeatures := automation.consensusSystem.GetNewCodeChanges() // Fetch new code changes or features to be deployed

    for _, feature := range newFeatures {
        fmt.Printf("New feature detected: %s. Initiating deployment.\n", feature.Name)
        deploymentSuccess := automation.deployFeature(feature)

        if deploymentSuccess {
            fmt.Printf("Feature %s deployed successfully.\n", feature.Name)
            automation.logDeploymentEvent(feature.Name, "Success")
        } else {
            fmt.Printf("Error deploying feature %s.\n", feature.Name)
            automation.logDeploymentEvent(feature.Name, "Failed")
        }
    }

    automation.deploymentCount++
    fmt.Printf("CI/CD deployment cycle #%d completed.\n", automation.deploymentCount)

    if automation.deploymentCount%SubBlocksPerBlock == 0 {
        automation.finalizeDeploymentCycle()
    }
}

// deployFeature performs the deployment of the new feature with encryption and logs the action
func (automation *ContinuousIntegrationDeploymentAutomation) deployFeature(feature common.Feature) bool {
    fmt.Printf("Encrypting deployment data for feature: %s\n", feature.Name)

    // Encrypt feature data before deploying it
    encryptedFeatureData, err := encryption.EncryptData(feature)
    if err != nil {
        fmt.Printf("Error encrypting feature data for %s: %s\n", feature.Name, err.Error())
        return false
    }

    feature.EncryptedData = encryptedFeatureData
    fmt.Printf("Feature data for %s encrypted successfully.\n", feature.Name)

    // Deploy the feature to the system
    return automation.consensusSystem.DeployFeature(feature)
}

// logDeploymentEvent logs the deployment process into the ledger for traceability
func (automation *ContinuousIntegrationDeploymentAutomation) logDeploymentEvent(featureName string, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("ci-deployment-%s", featureName),
        Timestamp: time.Now().Unix(),
        Type:      "CI/CD Deployment",
        Status:    result,
        Details:   fmt.Sprintf("Deployment result for feature %s: %s", featureName, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with CI/CD deployment result for feature %s: %s.\n", featureName, result)
}

// finalizeDeploymentCycle finalizes the deployment check cycle and logs the result in the ledger
func (automation *ContinuousIntegrationDeploymentAutomation) finalizeDeploymentCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeDeploymentCycle()
    if success {
        fmt.Println("CI/CD deployment cycle finalized successfully.")
        automation.logDeploymentCycleFinalization()
    } else {
        fmt.Println("Error finalizing CI/CD deployment cycle.")
    }
}

// logDeploymentCycleFinalization logs the finalization of a deployment check cycle into the ledger
func (automation *ContinuousIntegrationDeploymentAutomation) logDeploymentCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("ci-deployment-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "CI/CD Deployment Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with CI/CD deployment cycle finalization.")
}
