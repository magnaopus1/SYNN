package syn3400

import (
	"errors"
	"sync"
	"time"

)

// ForexPolicy represents a policy for managing Forex trading rules and compliance
type ForexPolicy struct {
	PolicyID      string    `json:"policy_id"`
	Description   string    `json:"description"`
	Rules         []string  `json:"rules"`
	EffectiveDate time.Time `json:"effective_date"`
	ExpiryDate    time.Time `json:"expiry_date"`
	Encrypted     bool      `json:"encrypted"` // Indicate if the policy is encrypted
}

// PolicyManager manages Forex trading policies for SYN3400 tokens
type PolicyManager struct {
	policies  map[string]ForexPolicy
	ledger    *ledger.Ledger
	encryptor *encryption.Encryptor
	consensus *consensus.SynnergyConsensus
	mutex     sync.Mutex
}

// NewPolicyManager creates a new instance of PolicyManager
func NewPolicyManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *PolicyManager {
	return &PolicyManager{
		policies:  make(map[string]ForexPolicy),
		ledger:    ledger,
		encryptor: encryptor,
		consensus: consensus,
	}
}

// AddPolicy adds a new Forex trading policy, encrypts it, and logs it in the ledger
func (pm *PolicyManager) AddPolicy(policy *ForexPolicy) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Validate inputs
	if policy == nil || policy.PolicyID == "" {
		return errors.New("invalid policy: policy is nil or has empty PolicyID")
	}

	// Encrypt the policy if needed
	if !policy.Encrypted {
		encryptedPolicy, err := pm.encryptor.EncryptData(policy)
		if err != nil {
			return err
		}
		policy = encryptedPolicy.(*ForexPolicy)
		policy.Encrypted = true
	}

	// Store the policy
	pm.policies[policy.PolicyID] = *policy

	// Log the policy addition in the ledger
	pm.ledger.LogEvent("PolicyAdded", time.Now(), policy.PolicyID)

	// Validate the policy addition with consensus
	return pm.consensus.ValidateSubBlock(policy.PolicyID)
}

// UpdatePolicy updates an existing Forex trading policy and ensures compliance with Synnergy Consensus
func (pm *PolicyManager) UpdatePolicy(policy *ForexPolicy) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Validate inputs
	if policy == nil || policy.PolicyID == "" {
		return errors.New("invalid policy: policy is nil or has empty PolicyID")
	}

	// Check if the policy exists
	if _, exists := pm.policies[policy.PolicyID]; !exists {
		return errors.New("policy does not exist")
	}

	// Encrypt the policy if needed
	if !policy.Encrypted {
		encryptedPolicy, err := pm.encryptor.EncryptData(policy)
		if err != nil {
			return err
		}
		policy = encryptedPolicy.(*ForexPolicy)
		policy.Encrypted = true
	}

	// Update the policy
	pm.policies[policy.PolicyID] = *policy

	// Log the policy update in the ledger
	pm.ledger.LogEvent("PolicyUpdated", time.Now(), policy.PolicyID)

	// Validate the policy update with consensus
	return pm.consensus.ValidateSubBlock(policy.PolicyID)
}

// GetPolicy retrieves a Forex trading policy by its ID
func (pm *PolicyManager) GetPolicy(policyID string) (*ForexPolicy, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Check if the policy exists
	policy, exists := pm.policies[policyID]
	if !exists {
		return nil, errors.New("policy not found")
	}

	// Return the policy
	return &policy, nil
}

// RemovePolicy removes an existing Forex trading policy, logs the removal in the ledger, and syncs with consensus
func (pm *PolicyManager) RemovePolicy(policyID string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Check if the policy exists
	if _, exists := pm.policies[policyID]; !exists {
		return errors.New("policy does not exist")
	}

	// Remove the policy
	delete(pm.policies, policyID)

	// Log the policy removal in the ledger
	pm.ledger.LogEvent("PolicyRemoved", time.Now(), policyID)

	// Validate the removal with consensus
	return pm.consensus.ValidateSubBlock(policyID)
}

// ListPolicies lists all active policies managed by the PolicyManager
func (pm *PolicyManager) ListPolicies() ([]ForexPolicy, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	var activePolicies []ForexPolicy
	currentTime := time.Now()

	for _, policy := range pm.policies {
		if policy.EffectiveDate.Before(currentTime) && policy.ExpiryDate.After(currentTime) {
			activePolicies = append(activePolicies, policy)
		}
	}

	return activePolicies, nil
}

// ValidatePolicy checks whether a policy is active based on its EffectiveDate and ExpiryDate
func (pm *PolicyManager) ValidatePolicy(policyID string) (bool, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Retrieve the policy
	policy, exists := pm.policies[policyID]
	if !exists {
		return false, errors.New("policy not found")
	}

	// Check if the policy is currently active
	currentTime := time.Now()
	if policy.EffectiveDate.Before(currentTime) && policy.ExpiryDate.After(currentTime) {
		return true, nil
	}

	return false, errors.New("policy is not active")
}
