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
    ProxyVoteCheckInterval       = 15 * time.Second  // Interval for checking proxy voting actions
    SubBlocksPerBlock            = 1000              // Number of sub-blocks in a block
    MaxProxyDelegationsPerMember = 3                 // Maximum number of allowed proxy delegations per DAO member
    MaxVoteAttempts              = 5                 // Maximum vote attempts for any proxy
)

// DAOProxyVotingRestrictionAutomation manages and enforces proxy voting rules within a DAO
type DAOProxyVotingRestrictionAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging proxy voting actions
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    voteDelegationMap map[string]int               // Tracks the number of proxies delegated by a member
    voteCycleCount    int                          // Counter for proxy voting check cycles
}

// NewDAOProxyVotingRestrictionAutomation initializes the automation for managing proxy voting restrictions
func NewDAOProxyVotingRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DAOProxyVotingRestrictionAutomation {
    return &DAOProxyVotingRestrictionAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        voteDelegationMap: make(map[string]int),
        voteCycleCount:    0,
    }
}

// StartProxyVotingMonitoring starts the continuous loop for enforcing proxy voting restrictions
func (automation *DAOProxyVotingRestrictionAutomation) StartProxyVotingMonitoring() {
    ticker := time.NewTicker(ProxyVoteCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndRestrictProxyVotes()
        }
    }()
}

// monitorAndRestrictProxyVotes checks all DAO proxy voting actions and applies restrictions
func (automation *DAOProxyVotingRestrictionAutomation) monitorAndRestrictProxyVotes() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of all members attempting proxy votes
    proxyVoteList := automation.consensusSystem.GetProxyVoteActions()

    if len(proxyVoteList) > 0 {
        for _, voteAction := range proxyVoteList {
            fmt.Printf("Checking proxy voting restrictions for member %s delegating to %s.\n", voteAction.DelegatorID, voteAction.ProxyID)
            automation.enforceProxyVotingRestriction(voteAction)
        }
    } else {
        fmt.Println("No proxy voting actions detected.")
    }

    automation.voteCycleCount++
    fmt.Printf("Proxy voting check cycle #%d executed.\n", automation.voteCycleCount)

    if automation.voteCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeVoteCycle()
    }
}

// enforceProxyVotingRestriction enforces the DAO's proxy voting restrictions for a member
func (automation *DAOProxyVotingRestrictionAutomation) enforceProxyVotingRestriction(voteAction common.ProxyVoteAction) {
    currentDelegations := automation.voteDelegationMap[voteAction.DelegatorID]

    if currentDelegations >= MaxProxyDelegationsPerMember {
        fmt.Printf("Proxy voting restriction violated by member %s. Max delegations: %d.\n", voteAction.DelegatorID, MaxProxyDelegationsPerMember)
        automation.logProxyViolation(voteAction.DelegatorID)
        automation.restrictProxyVote(voteAction)
    } else {
        automation.voteDelegationMap[voteAction.DelegatorID]++
        fmt.Printf("Proxy voting allowed for member %s. Current delegations: %d.\n", voteAction.DelegatorID, automation.voteDelegationMap[voteAction.DelegatorID])
        automation.logProxyVote(voteAction)
    }
}

// restrictProxyVote blocks a proxy voting action that violates DAO governance
func (automation *DAOProxyVotingRestrictionAutomation) restrictProxyVote(voteAction common.ProxyVoteAction) {
    restrictedVote := automation.consensusSystem.BlockProxyVote(voteAction)

    if restrictedVote {
        fmt.Printf("Proxy vote by %s restricted successfully.\n", voteAction.DelegatorID)
    } else {
        fmt.Printf("Error restricting proxy vote by %s.\n", voteAction.DelegatorID)
    }
}

// logProxyViolation logs a proxy voting restriction violation into the ledger
func (automation *DAOProxyVotingRestrictionAutomation) logProxyViolation(delegatorID string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("proxy-vote-violation-%s", delegatorID),
        Timestamp: time.Now().Unix(),
        Type:      "Proxy Vote Violation",
        Status:    "Violation",
        Details:   fmt.Sprintf("Member %s exceeded max proxy delegations.", delegatorID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with proxy voting violation for member %s.\n", delegatorID)
}

// logProxyVote logs a successful proxy vote action into the ledger
func (automation *DAOProxyVotingRestrictionAutomation) logProxyVote(voteAction common.ProxyVoteAction) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("proxy-vote-%s-to-%s", voteAction.DelegatorID, voteAction.ProxyID),
        Timestamp: time.Now().Unix(),
        Type:      "Proxy Vote",
        Status:    "Allowed",
        Details:   fmt.Sprintf("Member %s delegated proxy vote to %s.", voteAction.DelegatorID, voteAction.ProxyID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with proxy vote by member %s.\n", voteAction.DelegatorID)
}

// finalizeVoteCycle finalizes the proxy voting restriction cycle and logs it into the ledger
func (automation *DAOProxyVotingRestrictionAutomation) finalizeVoteCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeVoteCycle()
    if success {
        fmt.Println("Proxy voting restriction cycle finalized successfully.")
        automation.logVoteCycleFinalization()
    } else {
        fmt.Println("Error finalizing proxy voting restriction cycle.")
    }
}

// logVoteCycleFinalization logs the finalization of a proxy voting restriction cycle into the ledger
func (automation *DAOProxyVotingRestrictionAutomation) logVoteCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("proxy-vote-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Proxy Vote Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with proxy voting restriction cycle finalization.")
}

// encryptVoteData encrypts the proxy voting data before processing or logging
func (automation *DAOProxyVotingRestrictionAutomation) encryptVoteData(voteAction common.ProxyVoteAction) common.EncryptedProxyVoteData {
    encryptedData, err := encryption.EncryptData([]byte(voteAction.ProxyID))
    if err != nil {
        fmt.Println("Error encrypting proxy voting data:", err)
        return common.EncryptedProxyVoteData{}
    }

    fmt.Println("Proxy vote data successfully encrypted.")
    return common.EncryptedProxyVoteData{
        DelegatorID:   voteAction.DelegatorID,
        EncryptedProxy: encryptedData,
    }
}
