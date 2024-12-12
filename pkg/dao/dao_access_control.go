package dao

import (
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Roles within the DAO for access control
const (
	DAOMember         = "Member"
	DAOAdmin          = "Admin"
	DAOProposalAuthor = "ProposalAuthor"
	DAOAuditor        = "Auditor"

)

// Permissions define the actions allowed for each role
var rolePermissions = map[string][]string{
	DAOMember: {
		"view_proposals",
		"vote_proposals",
	},
	DAOAdmin: {
		"view_proposals",
		"vote_proposals",
		"create_proposals",
		"manage_members",
		"finalize_proposals",
	},
	DAOProposalAuthor: {
		"view_proposals",
		"vote_proposals",
		"create_proposals",
	},
	DAOAuditor: {
		"view_proposals",
		"audit_transactions",
	},
}



// NewAccessControl initializes a new AccessControl instance.
func NewAccessControl(daoID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.AccessControl {
	return &AccessControl{
		DAOID:             daoID,
		Members:           make(map[string]string),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// AssignRole assigns a specific role to a member.
func (ac *AccessControl) AssignRole(walletAddress, role string) error {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	// Validate role
	if _, exists := rolePermissions[role]; !exists {
		return errors.New("invalid role")
	}

	// Assign role
	ac.Members[walletAddress] = role

	// Encrypt the role information before storing in the ledger
	encryptedRole, err := ac.EncryptionService.EncryptData([]byte(role), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt role: %v", err)
	}

	// Store the role assignment in the ledger
	err = ac.Ledger.DAOLedger.RecordRoleAssignment(ac.DAOID, walletAddress, encryptedRole)
	if err != nil {
		return fmt.Errorf("failed to record role assignment in the ledger: %v", err)
	}

	fmt.Printf("Role %s assigned to member %s in DAO %s\n", role, walletAddress, ac.DAOID)
	return nil
}

// RevokeRole revokes the role from a DAO member.
func (ac *AccessControl) RevokeDAORole(walletAddress string) error {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	// Check if the member exists
	if _, exists := ac.Members[walletAddress]; !exists {
		return errors.New("member not found")
	}

	// Remove the member's role
	delete(ac.Members, walletAddress)

	// Record the revocation in the ledger
	err := ac.Ledger.DAOLedger.RecordRoleRevocation(ac.DAOID, walletAddress)
	if err != nil {
		return fmt.Errorf("failed to record role revocation in ledger: %v", err)
	}

	fmt.Printf("Role revoked from member %s in DAO %s\n", walletAddress, ac.DAOID)
	return nil
}

// CheckPermission checks if a wallet address has the permission to perform an action.
func (ac *AccessControl) CheckDAOPermission(walletAddress, action string) error {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	role, exists := ac.Members[walletAddress]
	if !exists {
		return errors.New("member not found")
	}

	permissions, roleExists := rolePermissions[role]
	if !roleExists {
		return errors.New("role not found")
	}

	// Check if the action is allowed for the given role
	for _, permission := range permissions {
		if permission == action {
			return nil
		}
	}

	return fmt.Errorf("permission denied: %s is not allowed to perform %s", role, action)
}

// GetMemberRole retrieves the role assigned to a member.
func (ac *AccessControl) GetDAOMemberRole(walletAddress string) (string, error) {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	role, exists := ac.Members[walletAddress]
	if !exists {
		return "", errors.New("member not found")
	}

	return role, nil
}

// ListRoles returns the available roles and their permissions.
func (ac *AccessControl) ListDAORoles() map[string][]string {
	return rolePermissions
}

// AuditTransaction allows an auditor to audit a transaction in the DAO.
func (ac *AccessControl) AuditTransaction(auditorAddress, transactionID string) error {
	// Check if the auditor has the permission to audit
	err := ac.CheckPermission(auditorAddress, "audit_transactions")
	if err != nil {
		return err
	}

	// Perform audit using the ledger
	transaction, err := ac.Ledger.BlockchainConsensusCoinLedger.GetTransactionByID(ac.DAOID, transactionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction: %v", err)
	}

	fmt.Printf("Transaction %s audited by %s in DAO %s\n", transactionID, auditorAddress, ac.DAOID)
	return nil
}

// EncryptMemberRoles encrypts all member roles and stores them securely in the ledger.
func (ac *AccessControl) EncryptDAOMemberRoles() error {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	for walletAddress, role := range ac.Members {
		encryptedRole, err := ac.EncryptionService.EncryptData([]byte(role), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt role: %v", err)
		}

		// Update the ledger with the encrypted role
		err = ac.Ledger.DAOLedger.UpdateEncryptedRole(ac.DAOID, walletAddress, encryptedRole)
		if err != nil {
			return fmt.Errorf("failed to update encrypted role in ledger: %v", err)
		}
	}

	fmt.Printf("All roles for DAO %s have been encrypted and stored in the ledger\n", ac.DAOID)
	return nil
}
