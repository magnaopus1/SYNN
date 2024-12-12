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
    AccessControlCheckInterval  = 10 * time.Second  // Interval for checking DAO access control
    SubBlocksPerBlock           = 1000              // Number of sub-blocks in a block
)

// DAOAccessControlAutomation manages DAO access control, ensuring permissions and roles are enforced
type DAOAccessControlAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging access control actions
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    accessCycleCount   int                          // Counter for access control cycles
    accessPermissions  map[string]common.DAORoles   // DAO roles and permissions for members
}

// NewDAOAccessControlAutomation initializes the automation for managing DAO access control
func NewDAOAccessControlAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DAOAccessControlAutomation {
    return &DAOAccessControlAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        accessPermissions: make(map[string]common.DAORoles),
    }
}

// StartAccessControlMonitoring starts the continuous loop for checking DAO access control
func (automation *DAOAccessControlAutomation) StartAccessControlMonitoring() {
    ticker := time.NewTicker(AccessControlCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAccessControl()
        }
    }()
}

// monitorAccessControl checks DAO roles and permissions to enforce access control
func (automation *DAOAccessControlAutomation) monitorAccessControl() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of DAO members and their permissions
    daoMembers := automation.consensusSystem.GetDAOMembers()

    for _, member := range daoMembers {
        role := automation.getMemberRole(member)
        fmt.Printf("Checking access control for member %s, role: %s.\n", member.ID, role)

        automation.enforceAccessControl(member, role)
    }

    automation.accessCycleCount++
    fmt.Printf("DAO access control cycle #%d executed.\n", automation.accessCycleCount)

    if automation.accessCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeAccessCycle()
    }
}

// getMemberRole retrieves the role of a DAO member
func (automation *DAOAccessControlAutomation) getMemberRole(member common.DAOMember) common.DAORoles {
    if role, exists := automation.accessPermissions[member.ID]; exists {
        return role
    }
    return common.DefaultRole // Return a default role if none is assigned
}

// enforceAccessControl enforces access control based on the member's role
func (automation *DAOAccessControlAutomation) enforceAccessControl(member common.DAOMember, role common.DAORoles) {
    // Encrypt sensitive role data
    encryptedRoleData := automation.encryptRoleData(member, role)

    // Validate access through the Synnergy Consensus system
    accessGranted := automation.consensusSystem.ValidateDAOAccess(member, encryptedRoleData)

    if accessGranted {
        fmt.Printf("Access granted to DAO member %s with role: %s.\n", member.ID, role)
        automation.logAccessEvent(member, "Access Granted", role)
    } else {
        fmt.Printf("Access denied to DAO member %s with role: %s.\n", member.ID, role)
        automation.logAccessEvent(member, "Access Denied", role)
    }
}

// finalizeAccessCycle finalizes the DAO access control cycle and logs the result in the ledger
func (automation *DAOAccessControlAutomation) finalizeAccessCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeAccessControlCycle()
    if success {
        fmt.Println("DAO access control cycle finalized successfully.")
        automation.logAccessCycleFinalization()
    } else {
        fmt.Println("Error finalizing DAO access control cycle.")
    }
}

// logAccessEvent logs an access control event into the ledger
func (automation *DAOAccessControlAutomation) logAccessEvent(member common.DAOMember, eventType string, role common.DAORoles) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("dao-access-%s-%s", member.ID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "DAO Access Event",
        Status:    eventType,
        Details:   fmt.Sprintf("DAO member %s access %s with role %s.", member.ID, eventType, role),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with access control event for member %s.\n", member.ID)
}

// logAccessCycleFinalization logs the finalization of an access control cycle into the ledger
func (automation *DAOAccessControlAutomation) logAccessCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("dao-access-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "DAO Access Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with DAO access control cycle finalization.")
}

// assignRole assigns a role to a DAO member
func (automation *DAOAccessControlAutomation) assignRole(member common.DAOMember, role common.DAORoles) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    automation.accessPermissions[member.ID] = role
    fmt.Printf("Role %s assigned to DAO member %s.\n", role, member.ID)
    automation.logRoleAssignment(member, role)
}

// revokeRole revokes a role from a DAO member
func (automation *DAOAccessControlAutomation) revokeRole(member common.DAOMember) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    delete(automation.accessPermissions, member.ID)
    fmt.Printf("Role revoked from DAO member %s.\n", member.ID)
    automation.logRoleRevocation(member)
}

// logRoleAssignment logs the role assignment event into the ledger
func (automation *DAOAccessControlAutomation) logRoleAssignment(member common.DAOMember, role common.DAORoles) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("dao-role-assignment-%s", member.ID),
        Timestamp: time.Now().Unix(),
        Type:      "DAO Role Assignment",
        Status:    "Assigned",
        Details:   fmt.Sprintf("Role %s assigned to DAO member %s.", role, member.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with role assignment for member %s.\n", member.ID)
}

// logRoleRevocation logs the role revocation event into the ledger
func (automation *DAOAccessControlAutomation) logRoleRevocation(member common.DAOMember) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("dao-role-revocation-%s", member.ID),
        Timestamp: time.Now().Unix(),
        Type:      "DAO Role Revocation",
        Status:    "Revoked",
        Details:   fmt.Sprintf("Role revoked from DAO member %s.", member.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with role revocation for member %s.\n", member.ID)
}

// encryptRoleData encrypts the DAO member role data before processing access control
func (automation *DAOAccessControlAutomation) encryptRoleData(member common.DAOMember, role common.DAORoles) common.EncryptedRoleData {
    encryptedData, err := encryption.EncryptData([]byte(role))
    if err != nil {
        fmt.Println("Error encrypting role data:", err)
        return common.EncryptedRoleData{}
    }

    fmt.Println("Role data successfully encrypted.")
    return common.EncryptedRoleData{
        MemberID:     member.ID,
        EncryptedRole: encryptedData,
    }
}
