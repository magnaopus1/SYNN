package defi

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewInsuranceManager initializes the DeFi Insurance Manager with a ledger and encryption service.
// Ensures the manager is ready for policy and claim operations.
func NewInsuranceManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *InsuranceManager {
	log.Printf("[INFO] Initializing Insurance Manager...")
	if ledgerInstance == nil {
		log.Fatalf("[ERROR] Ledger instance cannot be nil")
	}
	if encryptionService == nil {
		log.Fatalf("[ERROR] Encryption service cannot be nil")
	}

	manager := &InsuranceManager{
		Policies:          make(map[string]*InsurancePolicy),
		Claims:            make(map[string]*InsuranceClaim),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
		mu:                sync.Mutex{},
	}

	log.Printf("[SUCCESS] Insurance Manager initialized successfully.")
	return manager
}


// CreatePolicy creates a new insurance policy for a user.
// Validates inputs, encrypts sensitive data, and logs the policy creation in the ledger.
func (im *InsuranceManager) CreatePolicy(holder string, insuredAmount, premium float64, duration time.Duration) (*InsurancePolicy, error) {
	log.Printf("[INFO] Creating new insurance policy. Holder: %s", holder)

	// Step 1: Lock the manager to ensure thread-safe access.
	im.mu.Lock()
	defer im.mu.Unlock()

	// Step 2: Validate input parameters.
	if holder == "" {
		return nil, fmt.Errorf("holder address cannot be empty")
	}
	if insuredAmount <= 0 {
		return nil, fmt.Errorf("insured amount must be greater than zero")
	}
	if premium <= 0 {
		return nil, fmt.Errorf("premium must be greater than zero")
	}
	if duration <= 0 {
		return nil, fmt.Errorf("policy duration must be greater than zero")
	}

	// Step 3: Generate a unique policy ID.
	policyID := generateUniqueID()
	log.Printf("[INFO] Generated policy ID: %s", policyID)

	// Step 4: Encrypt the holder's address.
	encryptedHolder, err := im.EncryptionService.EncryptData("AES", []byte(holder), common.EncryptionKey)
	if err != nil {
		log.Printf("[ERROR] Failed to encrypt holder address for policy creation: %v", err)
		return nil, fmt.Errorf("failed to encrypt holder address: %w", err)
	}
	encryptedHolderStr := string(encryptedHolder)

	// Step 5: Create the insurance policy.
	startDate := time.Now()
	policy := &InsurancePolicy{
		PolicyID:        policyID,
		Holder:          holder,
		InsuredAmount:   insuredAmount,
		Premium:         premium,
		PolicyDuration:  duration,
		StartDate:       startDate,
		ExpiryDate:      startDate.Add(duration),
		Status:          "Active",
		EncryptedHolder: encryptedHolderStr,
	}
	log.Printf("[INFO] Policy created in memory. Policy ID: %s", policyID)

	// Step 6: Add the policy to the manager's local state.
	im.Policies[policyID] = policy

	// Step 7: Record the policy creation in the ledger.
	log.Printf("[INFO] Recording policy creation in ledger. Policy ID: %s", policyID)
	err = im.Ledger.DeFiLedger.RecordPolicyCreation(policyID, holder, insuredAmount, premium, duration)
	if err != nil {
		delete(im.Policies, policyID) // Rollback if ledger operation fails.
		log.Printf("[ERROR] Failed to log policy creation in the ledger for policy %s: %v", policyID, err)
		return nil, fmt.Errorf("failed to log policy creation in the ledger: %w", err)
	}

	// Step 8: Log Success.
	log.Printf("[SUCCESS] Insurance policy created successfully. Policy ID: %s, Holder: %s", policyID, holder)
	return policy, nil
}

// ClaimPolicy allows a user to make a claim on an active insurance policy.
// Validates the policy and claim amount, encrypts sensitive data, and logs the claim in the ledger.
func (im *InsuranceManager) ClaimPolicy(policyID string, claimAmount float64) (*InsuranceClaim, error) {
	log.Printf("[INFO] Processing claim request. Policy ID: %s, Claim Amount: %.2f", policyID, claimAmount)

	// Step 1: Lock the manager to ensure thread safety.
	im.mu.Lock()
	defer im.mu.Unlock()

	// Step 2: Validate input parameters.
	if policyID == "" {
		return nil, fmt.Errorf("policy ID cannot be empty")
	}
	if claimAmount <= 0 {
		return nil, fmt.Errorf("claim amount must be greater than zero")
	}

	// Step 3: Retrieve the policy.
	policy, exists := im.Policies[policyID]
	if !exists {
		log.Printf("[ERROR] Policy not found: %s", policyID)
		return nil, fmt.Errorf("policy %s not found", policyID)
	}

	// Step 4: Validate policy status.
	if policy.Status != "Active" {
		log.Printf("[ERROR] Policy %s is not active. Current status: %s", policyID, policy.Status)
		return nil, fmt.Errorf("policy %s is not active", policyID)
	}
	if time.Now().After(policy.ExpiryDate) {
		log.Printf("[ERROR] Policy %s is expired. Expiry Date: %v", policyID, policy.ExpiryDate)
		return nil, fmt.Errorf("policy %s is expired", policyID)
	}

	// Step 5: Validate claim amount.
	if claimAmount > policy.InsuredAmount {
		log.Printf("[ERROR] Claim amount exceeds insured amount for policy %s: Claim=%.2f, Insured=%.2f", policyID, claimAmount, policy.InsuredAmount)
		return nil, fmt.Errorf("claim amount exceeds insured amount for policy %s", policyID)
	}

	// Step 6: Generate unique claim ID and encrypt claim data.
	claimID := generateUniqueID()
	claimData := fmt.Sprintf("PolicyID: %s, ClaimAmount: %.2f", policyID, claimAmount)
	encryptedClaimData, err := im.EncryptionService.EncryptData("AES", []byte(claimData), common.EncryptionKey)
	if err != nil {
		log.Printf("[ERROR] Failed to encrypt claim data for policy %s: %v", policyID, err)
		return nil, fmt.Errorf("failed to encrypt claim data: %w", err)
	}

	// Step 7: Create and store the claim.
	claim := &InsuranceClaim{
		ClaimID:       claimID,
		PolicyID:      policyID,
		ClaimAmount:   claimAmount,
		ClaimDate:     time.Now(),
		ClaimStatus:   "Pending",
		EncryptedData: string(encryptedClaimData),
	}
	im.Claims[claimID] = claim

	// Step 8: Record claim in the ledger.
	log.Printf("[INFO] Recording claim submission in ledger. Claim ID: %s", claimID)
	err = im.Ledger.DeFiLedger.RecordClaimSubmission(claimID, claimAmount)
	if err != nil {
		delete(im.Claims, claimID) // Rollback claim creation in case of failure.
		log.Printf("[ERROR] Failed to record claim in ledger for claim %s: %v", claimID, err)
		return nil, fmt.Errorf("failed to record claim in ledger: %w", err)
	}

	// Step 9: Log success and return the claim.
	log.Printf("[SUCCESS] Claim submitted successfully. Claim ID: %s, Policy ID: %s", claimID, policyID)
	return claim, nil
}


// ApproveClaim approves a pending insurance claim and updates the associated policy status.
// Validates the claim and policy, updates their statuses, and logs the approval in the ledger.
func (im *InsuranceManager) ApproveClaim(claimID string) error {
	log.Printf("[INFO] Approving claim. Claim ID: %s", claimID)

	// Step 1: Lock the manager to ensure thread safety.
	im.mu.Lock()
	defer im.mu.Unlock()

	// Step 2: Validate claim ID.
	if claimID == "" {
		return fmt.Errorf("claim ID cannot be empty")
	}

	// Step 3: Retrieve the claim.
	claim, exists := im.Claims[claimID]
	if !exists {
		log.Printf("[ERROR] Claim not found: %s", claimID)
		return fmt.Errorf("claim %s not found", claimID)
	}

	// Step 4: Validate claim status.
	if claim.ClaimStatus != "Pending" {
		log.Printf("[ERROR] Claim %s is not in a pending state. Current status: %s", claimID, claim.ClaimStatus)
		return fmt.Errorf("claim %s is not in a pending state", claimID)
	}

	// Step 5: Retrieve the associated policy.
	policy, exists := im.Policies[claim.PolicyID]
	if !exists {
		log.Printf("[ERROR] Policy not found for claim %s: %s", claimID, claim.PolicyID)
		return fmt.Errorf("policy %s not found for claim %s", claim.PolicyID, claimID)
	}

	// Step 6: Approve the claim and update policy status.
	log.Printf("[INFO] Approving claim and updating policy. Claim ID: %s, Policy ID: %s", claimID, policy.PolicyID)
	claim.ClaimStatus = "Approved"
	policy.Status = "Claimed"

	// Step 7: Record the claim approval in the ledger.
	err := im.Ledger.DeFiLedger.RecordClaimApproval(claimID)
	if err != nil {
		log.Printf("[ERROR] Failed to record claim approval in ledger for claim %s: %v", claimID, err)
		claim.ClaimStatus = "Pending" // Rollback claim status.
		policy.Status = "Active"      // Rollback policy status.
		return fmt.Errorf("failed to record claim approval in ledger: %w", err)
	}

	// Step 8: Log success.
	log.Printf("[SUCCESS] Claim approved successfully. Claim ID: %s, Policy ID: %s", claimID, policy.PolicyID)
	return nil
}


// RejectClaim rejects a pending insurance claim.
// Validates the claim, updates its status, and logs the rejection in the ledger.
func (im *InsuranceManager) RejectClaim(claimID string) error {
	log.Printf("[INFO] Rejecting claim. Claim ID: %s", claimID)

	// Step 1: Lock the InsuranceManager to ensure thread safety.
	im.mu.Lock()
	defer im.mu.Unlock()

	// Step 2: Validate the claim ID.
	if claimID == "" {
		return fmt.Errorf("claim ID cannot be empty")
	}

	// Step 3: Retrieve the claim.
	claim, exists := im.Claims[claimID]
	if !exists {
		log.Printf("[ERROR] Claim not found: %s", claimID)
		return fmt.Errorf("claim %s not found", claimID)
	}

	// Step 4: Validate the claim status.
	if claim.ClaimStatus != "Pending" {
		log.Printf("[ERROR] Claim %s is not in a pending state. Current status: %s", claimID, claim.ClaimStatus)
		return fmt.Errorf("claim %s is not in a pending state", claimID)
	}

	// Step 5: Update the claim status to "Rejected".
	log.Printf("[INFO] Rejecting claim. Claim ID: %s", claimID)
	claim.ClaimStatus = "Rejected"

	// Step 6: Record the rejection in the ledger.
	err := im.Ledger.DeFiLedger.RecordClaimRejection(claimID)
	if err != nil {
		log.Printf("[ERROR] Failed to log claim rejection in ledger for claim %s: %v", claimID, err)
		claim.ClaimStatus = "Pending" // Rollback the claim status to its original state in case of failure.
		return fmt.Errorf("failed to log claim rejection in ledger: %w", err)
	}

	// Step 7: Log success and return.
	log.Printf("[SUCCESS] Claim rejected successfully. Claim ID: %s", claimID)
	return nil
}

