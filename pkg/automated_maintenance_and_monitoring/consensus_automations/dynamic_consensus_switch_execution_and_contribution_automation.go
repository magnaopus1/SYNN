package consensus_automations

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    DynamicSwitchInterval        = 5 * time.Second  // Interval for dynamic validator allocation
    ContributionLoggingInterval  = 15 * time.Second // Interval for logging validator contributions
    SwitchExecutionKey           = "dynamic_switch_execution_key" // Encryption key for consensus switch and contribution logs
    ContributionTrackingKey      = "contribution_tracking_key"    // Encryption key for tracking validator contributions
)

// DynamicConsensusSwitchExecutionAndContributionAutomation automates consensus switching and contribution management in Synnergy Consensus
type DynamicConsensusSwitchExecutionAndContributionAutomation struct {
    ledgerInstance   *ledger.Ledger                   // Blockchain ledger for tracking consensus actions
    consensusEngine  *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine
    stateMutex       *sync.RWMutex                    // Mutex for thread-safe ledger access
    apiURL           string                           // API URL for consensus operations
}

// NewDynamicConsensusSwitchExecutionAndContributionAutomation initializes the automation for dynamic validator allocation and contribution tracking
func NewDynamicConsensusSwitchExecutionAndContributionAutomation(apiURL string, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *DynamicConsensusSwitchExecutionAndContributionAutomation {
    return &DynamicConsensusSwitchExecutionAndContributionAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        apiURL:          apiURL,
    }
}

// StartExecutionMonitoring initiates continuous monitoring for dynamic validator allocation, contribution, and efficiency
func (automation *DynamicConsensusSwitchExecutionAndContributionAutomation) StartExecutionMonitoring() {
    ticker := time.NewTicker(DynamicSwitchInterval)
    go func() {
        for range ticker.C {
            automation.allocateValidators()
            automation.trackMultiValidatorContributions()
        }
    }()

    contributionTicker := time.NewTicker(ContributionLoggingInterval)
    go func() {
        for range contributionTicker.C {
            automation.logValidatorContributions()
        }
    }()
}

// allocateValidators dynamically allocates validators to PoH, PoS, or PoW based on system performance and load
func (automation *DynamicConsensusSwitchExecutionAndContributionAutomation) allocateValidators() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Retrieve system load and validator performance metrics
    systemLoad := automation.consensusEngine.MeasureSystemLoad()
    validatorStake := automation.consensusEngine.GetTotalStake()
    networkCongestion := automation.consensusEngine.MeasureNetworkCongestion()

    // Adjust allocation based on system load and network conditions
    if systemLoad > 0.8 || networkCongestion > 0.7 {
        fmt.Println("High load detected, increasing PoH validator allocation.")
        automation.adjustValidatorAllocation("PoH", validatorStake)
    } else if systemLoad > 0.5 {
        fmt.Println("Moderate load detected, increasing PoS validator allocation.")
        automation.adjustValidatorAllocation("PoS", validatorStake)
    } else {
        fmt.Println("Low load detected, favoring PoW block mining.")
        automation.adjustValidatorAllocation("PoW", validatorStake)
    }
}

// adjustValidatorAllocation dynamically adjusts validator allocations to PoH, PoS, or PoW
func (automation *DynamicConsensusSwitchExecutionAndContributionAutomation) adjustValidatorAllocation(consensusType string, totalStake float64) {
    allocationURL := fmt.Sprintf("%s/api/consensus/%s/add-stake", automation.apiURL, consensusType)
    allocationPayload := map[string]float64{"stake": totalStake * 0.1} // Example allocation logic

    reqBody, _ := json.Marshal(allocationPayload)
    resp, err := http.Post(allocationURL, "application/json", bytes.NewBuffer(reqBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error adjusting validator allocation for %s: %v\n", consensusType, err)
    } else {
        fmt.Printf("Validator allocation adjusted for %s.\n", consensusType)
    }
}

// trackMultiValidatorContributions handles contributions from multiple validators in sub-block validation and block mining
func (automation *DynamicConsensusSwitchExecutionAndContributionAutomation) trackMultiValidatorContributions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Retrieve contributions from validators
    contributions := automation.consensusEngine.GetValidatorContributions()

    for validatorID, contribution := range contributions {
        // Register each validator's contribution to the validation/mining process
        automation.registerValidatorContribution(validatorID, contribution)
    }
}

// registerValidatorContribution tracks the contribution of a validator
func (automation *DynamicConsensusSwitchExecutionAndContributionAutomation) registerValidatorContribution(validatorID string, contribution float64) {
    fmt.Printf("Registering contribution for validator %s: %f\n", validatorID, contribution)

    contributionEntry := common.ContributionEntry{
        ValidatorID: validatorID,
        Contribution: contribution,
        Timestamp: time.Now().Unix(),
    }

    // Encrypt the contribution entry for security purposes
    encryptedEntry, err := encryption.EncryptContributionEntry(contributionEntry, []byte(ContributionTrackingKey))
    if err != nil {
        fmt.Printf("Error encrypting contribution entry for validator %s: %v\n", validatorID, err)
        return
    }

    automation.consensusEngine.RecordValidatorContribution(encryptedEntry)
    fmt.Printf("Validator contribution registered: %s - %f.\n", validatorID, contribution)
}

// logValidatorContributions logs validator contributions to the ledger
func (automation *DynamicConsensusSwitchExecutionAndContributionAutomation) logValidatorContributions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    contributionLogs := automation.consensusEngine.GetContributionLogs()

    for _, logEntry := range contributionLogs {
        ledgerEntry := common.LedgerEntry{
            ID:        fmt.Sprintf("contribution-log-%d", time.Now().Unix()),
            Timestamp: time.Now().Unix(),
            Type:      "Validator Contribution",
            Status:    "Logged",
            Details:   logEntry,
        }

        // Encrypt the ledger entry for security purposes
        encryptedEntry, err := encryption.EncryptLedgerEntry(ledgerEntry, []byte(SwitchExecutionKey))
        if err != nil {
            fmt.Printf("Error encrypting contribution log: %v\n", err)
            continue
        }

        automation.ledgerInstance.AddEntry(encryptedEntry)
        fmt.Printf("Ledger updated with contribution log: %s\n", logEntry)
    }
}
