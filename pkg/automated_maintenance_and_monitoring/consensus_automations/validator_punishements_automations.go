package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    PunishmentCheckInterval     = 10 * time.Minute   // Interval for checking validator punishments
    PoWPunishmentThreshold      = 5.0                // Threshold for PoW validator punishment (e.g., 5 failed block attempts)
    PoSPunishmentThreshold      = 24.0               // Threshold for PoS validator punishment (e.g., 24 hours of downtime)
    PoHPunishmentThreshold      = 1.0                // Threshold for PoH validator punishment (e.g., 1 missed participation cycle)
)

// ValidatorPunishmentsAutomation tracks validator behavior and enforces punishments for PoS, PoW, and PoH validators.
type ValidatorPunishmentsAutomation struct {
    ledgerInstance     *ledger.Ledger                        // Blockchain ledger for tracking punishment events
    consensusEngine    *synnergy_consensus.SynnergyConsensus  // Synnergy Consensus engine for monitoring validator behavior
    punishmentManager  *synnergy_consensus.PunishmentManager // Punishment Manager for handling penalties
    rewardManager      *synnergy_consensus.RewardManager     // Reward Manager for invoking punishments when necessary
    stateMutex         *sync.RWMutex                        // Mutex for thread-safe operations
}

// NewValidatorPunishmentsAutomation initializes the punishment automation for validators.
func NewValidatorPunishmentsAutomation(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, punishmentManager *synnergy_consensus.PunishmentManager, rewardManager *synnergy_consensus.RewardManager, stateMutex *sync.RWMutex) *ValidatorPunishmentsAutomation {
    return &ValidatorPunishmentsAutomation{
        ledgerInstance:    ledgerInstance,
        consensusEngine:   consensusEngine,
        punishmentManager: punishmentManager,
        rewardManager:     rewardManager,
        stateMutex:        stateMutex,
    }
}

// StartPunishmentMonitoring begins the continuous monitoring of validator behavior for punishment enforcement.
func (automation *ValidatorPunishmentsAutomation) StartPunishmentMonitoring() {
    ticker := time.NewTicker(PunishmentCheckInterval)

    go func() {
        for range ticker.C {
            fmt.Println("Checking validator behavior for punishments...")
            automation.checkAndEnforcePunishments()
        }
    }()
}

// checkAndEnforcePunishments checks validators' behavior and enforces punishments if thresholds are violated.
func (automation *ValidatorPunishmentsAutomation) checkAndEnforcePunishments() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Check for PoS validator violations (e.g., downtime, performance issues)
    automation.checkPoSValidators()

    // Check for PoW validator violations (e.g., failed mining attempts)
    automation.checkPoWValidators()

    // Check for PoH participation violations (e.g., missed cycles)
    automation.checkPoHValidators()
}

// checkPoSValidators checks and enforces punishments for PoS validators based on downtime or performance issues.
func (automation *ValidatorPunishmentsAutomation) checkPoSValidators() {
    violations := make(map[string]float64)
    
    posValidators := automation.consensusEngine.GetActivePoSValidators()
    for _, validator := range posValidators {
        downtime := automation.consensusEngine.GetPoSValidatorDowntime(validator)
        if downtime >= PoSPunishmentThreshold {
            violations[validator.Address] = downtime
            fmt.Printf("PoS validator %s has violated the threshold with %.2f hours of downtime.\n", validator.Address, downtime)
        }
    }

    // Enforce punishments if there are any violations
    if len(violations) > 0 {
        automation.punishmentManager.EnforcePunishments(violations, "PoS")
    }
}

// checkPoWValidators checks and enforces punishments for PoW validators based on failed block mining attempts.
func (automation *ValidatorPunishmentsAutomation) checkPoWValidators() {
    violations := make(map[string]float64)
    
    powValidators := automation.consensusEngine.GetActivePoWValidators()
    for _, miner := range powValidators {
        failedAttempts := automation.consensusEngine.GetPoWFailedAttempts(miner)
        if failedAttempts >= PoWPunishmentThreshold {
            violations[miner.Address] = failedAttempts
            fmt.Printf("PoW miner %s has violated the threshold with %.2f failed block attempts.\n", miner.Address, failedAttempts)
        }
    }

    // Enforce punishments if there are any violations
    if len(violations) > 0 {
        automation.punishmentManager.EnforcePunishments(violations, "PoW")
    }
}

// checkPoHValidators checks and enforces punishments for PoH validators based on missed participation cycles.
func (automation *ValidatorPunishmentsAutomation) checkPoHValidators() {
    violations := make(map[string]float64)

    pohParticipants := automation.consensusEngine.GetActivePoHParticipants()
    for _, participant := range pohParticipants {
        missedCycles := automation.consensusEngine.GetPoHParticipationMisses(participant)
        if missedCycles >= PoHPunishmentThreshold {
            violations[participant.Address] = missedCycles
            fmt.Printf("PoH participant %s has violated the threshold with %.2f missed participation cycles.\n", participant.Address, missedCycles)
        }
    }

    // Enforce punishments if there are any violations
    if len(violations) > 0 {
        automation.punishmentManager.EnforcePunishments(violations, "PoH")
    }
}

// ResetPunishments triggers a reset of older punishments after a specified interval (e.g., 90 days).
func (automation *ValidatorPunishmentsAutomation) ResetPunishments() {
    fmt.Println("Resetting older punishments...")
    automation.punishmentManager.ResetPunishments()
}

// ManualOverridePunishment allows for the manual enforcement of punishment on a specific validator.
func (automation *ValidatorPunishmentsAutomation) ManualOverridePunishment(validatorAddress string, punishmentAmount float64, category string) {
    fmt.Printf("Manually overriding punishment for validator %s with %.2f SYNN for %s violation.\n", validatorAddress, punishmentAmount, category)

    violations := map[string]float64{
        validatorAddress: punishmentAmount,
    }

    // Enforce manual punishment immediately
    automation.punishmentManager.EnforcePunishments(violations, category)
}
