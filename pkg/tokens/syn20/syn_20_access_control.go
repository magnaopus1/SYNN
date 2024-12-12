package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// AccessControlManager manages role-based access control for SYN20 contracts.
type AccessControlManager struct {
	mutex        sync.Mutex
	RoleMappings map[string]map[string]string // contractID -> (address -> role)
	Ledger       *ledger.Ledger               // Ledger for recording role changes and actions
	Encryption   *encryption.Encryption       // Encryption service to secure role-based data
}

// RoleType defines possible roles in the SYN20 access control system.
type RoleType string

const (
	RoleAdmin   RoleType = "admin"
	RoleMinter  RoleType = "minter"
	RoleBurner  RoleType = "burner"
	RoleHolder  RoleType = "holder"
	RoleDefault RoleType = "default"
)

// NewAccessControlManager initializes a new AccessControlManager.
func NewAccessControlManager(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *AccessControlManager {
	return &AccessControlManager{
		RoleMappings: make(map[string]map[string]string),
		Ledger:       ledgerInstance,
		Encryption:   encryptionService,
	}
}

// AssignRole assigns a role to a specific address for a contract.
func (acm *AccessControlManager) AssignRole(contractID string, address string, role RoleType) error {
	acm.mutex.Lock()
	defer acm.mutex.Unlock()

	if _, exists := acm.RoleMappings[contractID]; !exists {
		acm.RoleMappings[contractID] = make(map[string]string)
	}

	// Encrypt role assignment data
	encryptedRole, err := acm.Encryption.EncryptData(string(role), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting role assignment for address %s: %v", address, err)
	}

	// Log role assignment to the ledger
	err = acm.Ledger.RecordRoleAssignment(contractID, address, encryptedRole)
	if err != nil {
		return fmt.Errorf("error logging role assignment for address %s: %v", address, err)
	}

	// Assign the role
	acm.RoleMappings[contractID][address] = string(role)
	fmt.Printf("Assigned role %s to address %s for contract %s.\n", role, address, contractID)
	return nil
}

// RemoveRole removes a role from a specific address for a contract.
func (acm *AccessControlManager) RemoveRole(contractID string, address string) error {
	acm.mutex.Lock()
	defer acm.mutex.Unlock()

	if _, exists := acm.RoleMappings[contractID]; !exists {
		return errors.New("contract does not exist")
	}

	if _, exists := acm.RoleMappings[contractID][address]; !exists {
		return errors.New("role not assigned to the address")
	}

	// Encrypt role removal data
	encryptedRemoval, err := acm.Encryption.EncryptData("role_removed", common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting role removal for address %s: %v", address, err)
	}

	// Log role removal to the ledger
	err = acm.Ledger.RecordRoleRemoval(contractID, address, encryptedRemoval)
	if err != nil {
		return fmt.Errorf("error logging role removal for address %s: %v", address, err)
	}

	delete(acm.RoleMappings[contractID], address)
	fmt.Printf("Removed role from address %s for contract %s.\n", address, contractID)
	return nil
}

// CheckRole checks if an address has the specified role for a contract.
func (acm *AccessControlManager) CheckRole(contractID string, address string, role RoleType) (bool, error) {
	acm.mutex.Lock()
	defer acm.mutex.Unlock()

	if _, exists := acm.RoleMappings[contractID]; !exists {
		return false, errors.New("contract does not exist")
	}

	if assignedRole, exists := acm.RoleMappings[contractID][address]; exists && assignedRole == string(role) {
		return true, nil
	}

	return false, nil
}

// GetRolesForContract retrieves all role assignments for a specific contract.
func (acm *AccessControlManager) GetRolesForContract(contractID string) (map[string]string, error) {
	acm.mutex.Lock()
	defer acm.mutex.Unlock()

	if roles, exists := acm.RoleMappings[contractID]; exists {
		return roles, nil
	}

	return nil, errors.New("contract not found")
}

// GetRoleForAddress retrieves the role assigned to an address in a contract.
func (acm *AccessControlManager) GetRoleForAddress(contractID string, address string) (string, error) {
	acm.mutex.Lock()
	defer acm.mutex.Unlock()

	if _, exists := acm.RoleMappings[contractID]; !exists {
		return "", errors.New("contract does not exist")
	}

	if role, exists := acm.RoleMappings[contractID][address]; exists {
		return role, nil
	}

	return "", errors.New("role not assigned to address")
}

// ValidateAdmin validates if an address is an admin for a contract.
func (acm *AccessControlManager) ValidateAdmin(contractID string, address string) (bool, error) {
	return acm.CheckRole(contractID, address, RoleAdmin)
}

// ValidateMinter validates if an address is a minter for a contract.
func (acm *AccessControlManager) ValidateMinter(contractID string, address string) (bool, error) {
	return acm.CheckRole(contractID, address, RoleMinter)
}

// ValidateBurner validates if an address is a burner for a contract.
func (acm *AccessControlManager) ValidateBurner(contractID string, address string) (bool, error) {
	return acm.CheckRole(contractID, address, RoleBurner)
}
