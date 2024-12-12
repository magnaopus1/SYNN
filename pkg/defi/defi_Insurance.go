package defi

import (
	"errors"
	"fmt"
	"log"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"time"
)

// InsuranceCreatePolicy creates a new insurance policy with specified terms.
// Encrypts sensitive data, validates inputs, and records the policy in the ledger.
func InsuranceCreatePolicy(
    policyID, insuredEntity string, premium float64, terms string, 
    coverageAmount float64, duration time.Duration, ledgerInstance *ledger.Ledger,
) error {
    // Log the start of the operation
    log.Printf("[INFO] Starting InsuranceCreatePolicy. Policy ID: %s, Insured Entity: %s", policyID, insuredEntity)

    // Step 1: Validate input parameters
    if policyID == "" || insuredEntity == "" || terms == "" {
        err := errors.New("policyID, insuredEntity, and terms cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if premium <= 0 {
        err := errors.New("premium must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if coverageAmount <= 0 {
        err := errors.New("coverage amount must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if duration <= 0 {
        err := errors.New("duration must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt the terms of the policy
    encryptedTerms, err := encryption.EncryptString(terms)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt policy terms for policy ID %s: %v", policyID, err)
        return fmt.Errorf("failed to encrypt policy terms: %w", err)
    }

    // Step 3: Create the policy object
    policy := InsurancePolicy{
        PolicyID:       policyID,
        InsuredEntity:  insuredEntity,
        Premium:        premium,
        Terms:          encryptedTerms,
        CoverageAmount: coverageAmount,
        Status:         "Pending",
        Duration:       duration,
        CreatedAt:      time.Now(),
    }

    // Step 4: Record the policy in the ledger
    if err := ledgerInstance.DeFiLedger.CreatePolicy(policy); err != nil {
        log.Printf("[ERROR] Failed to create policy %s: %v", policyID, err)
        return fmt.Errorf("failed to create policy: %w", err)
    }

    // Log success
    log.Printf("[SUCCESS] Insurance policy created successfully. Policy ID: %s", policyID)
    return nil
}


// InsuranceActivatePolicy activates a pending insurance policy.
// Validates the policy ID and updates the ledger.
func InsuranceActivatePolicy(policyID string, ledgerInstance *ledger.Ledger) error {
    // Log the start of the operation
    log.Printf("[INFO] Starting InsuranceActivatePolicy. Policy ID: %s", policyID)

    // Step 1: Validate the policy ID
    if policyID == "" {
        err := errors.New("policyID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Activate the policy in the ledger
    if err := ledgerInstance.DeFiLedger.ActivatePolicy(policyID); err != nil {
        log.Printf("[ERROR] Failed to activate policy %s: %v", policyID, err)
        return fmt.Errorf("failed to activate policy: %w", err)
    }

    // Log success
    log.Printf("[SUCCESS] Policy activated successfully. Policy ID: %s", policyID)
    return nil
}


// InsuranceClaimPolicy processes a claim on an insurance policy.
// Encrypts claimant data, validates inputs, and records the claim in the ledger.
func InsuranceClaimPolicy(policyID, claimantID string, claimAmount float64, ledgerInstance *ledger.Ledger) error {
    // Log the start of the operation
    log.Printf("[INFO] Starting InsuranceClaimPolicy. Policy ID: %s, Claimant ID: %s", policyID, claimantID)

    // Step 1: Validate input parameters
    if policyID == "" || claimantID == "" {
        err := errors.New("policyID and claimantID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if claimAmount <= 0 {
        err := errors.New("claim amount must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt the claimant ID
    encryptedClaimantID, err := encryption.EncryptString(claimantID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt claimant ID for policy %s: %v", policyID, err)
        return fmt.Errorf("failed to encrypt claimant ID: %w", err)
    }

    // Step 3: Record the claim in the ledger
    if err := ledgerInstance.DeFiLedger.ClaimPolicy(policyID, encryptedClaimantID, claimAmount); err != nil {
        log.Printf("[ERROR] Failed to process claim for policy %s: %v", policyID, err)
        return fmt.Errorf("failed to process claim: %w", err)
    }

    // Log success
    log.Printf("[SUCCESS] Claim processed successfully. Policy ID: %s, Claim Amount: %.2f", policyID, claimAmount)
    return nil
}


// InsuranceSetPremium updates the premium for an insurance policy.
// Validates inputs and records the change in the ledger.
func InsuranceSetPremium(policyID string, premium float64, ledgerInstance *ledger.Ledger) error {
    // Log the start of the operation
    log.Printf("[INFO] Starting InsuranceSetPremium. Policy ID: %s, New Premium: %.2f", policyID, premium)

    // Step 1: Validate input parameters
    if policyID == "" {
        err := errors.New("policyID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if premium <= 0 {
        err := errors.New("premium must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Update the premium in the ledger
    if err := ledgerInstance.DeFiLedger.SetPremium(policyID, premium); err != nil {
        log.Printf("[ERROR] Failed to set premium for policy %s: %v", policyID, err)
        return fmt.Errorf("failed to set premium: %w", err)
    }

    // Log success
    log.Printf("[SUCCESS] Premium updated successfully for policy %s. New Premium: %.2f", policyID, premium)
    return nil
}


// InsuranceCalculatePayout calculates the payout for a claim on a policy.
// Validates inputs and fetches the payout from the ledger.
func InsuranceCalculatePayout(policyID string, claimAmount float64, ledgerInstance *ledger.Ledger) (float64, error) {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceCalculatePayout. Policy ID: %s, Claim Amount: %.2f", policyID, claimAmount)

	// Step 1: Validate input parameters
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return 0, err
	}
	if claimAmount <= 0 {
		err := errors.New("claim amount must be greater than zero")
		log.Printf("[ERROR] %v", err)
		return 0, err
	}

	// Step 2: Calculate the payout in the ledger
	payout, err := ledgerInstance.DeFiLedger.CalculatePayout(policyID, claimAmount)
	if err != nil {
		log.Printf("[ERROR] Failed to calculate payout for policy %s: %v", policyID, err)
		return 0, fmt.Errorf("failed to calculate payout: %w", err)
	}

	// Log success
	log.Printf("[SUCCESS] Payout calculated successfully. Policy ID: %s, Payout: %.2f", policyID, payout)
	return payout, nil
}


// InsuranceVerifyClaim verifies the validity of a claim on an insurance policy.
// Validates inputs and verifies the claim in the ledger.
func InsuranceVerifyClaim(policyID, claimantID string, ledgerInstance *ledger.Ledger) (bool, error) {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceVerifyClaim. Policy ID: %s, Claimant ID: %s", policyID, claimantID)

	// Step 1: Validate input parameters
	if policyID == "" || claimantID == "" {
		err := errors.New("policyID and claimantID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return false, err
	}

	// Step 2: Verify the claim in the ledger
	verified, err := ledgerInstance.DeFiLedger.VerifyClaim(policyID, claimantID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify claim for policy %s: %v", policyID, err)
		return false, fmt.Errorf("failed to verify claim: %w", err)
	}

	// Log success
	log.Printf("[SUCCESS] Claim verification completed. Policy ID: %s, Claimant ID: %s, Verified: %t", policyID, claimantID, verified)
	return verified, nil
}


// InsuranceDistributePayout distributes the payout for a claim on a policy.
// Encrypts claimant data, validates inputs, and records the distribution in the ledger.
func InsuranceDistributePayout(policyID, claimantID string, payoutAmount float64, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceDistributePayout. Policy ID: %s, Claimant ID: %s, Payout Amount: %.2f", policyID, claimantID, payoutAmount)

	// Step 1: Validate input parameters
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}
	if claimantID == "" {
		err := errors.New("claimantID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}
	if payoutAmount <= 0 {
		err := errors.New("payout amount must be greater than zero")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Encrypt the claimant ID
	encryptedClaimantID, err := encryption.EncryptString(claimantID)
	if err != nil {
		log.Printf("[ERROR] Failed to encrypt claimant ID for policy %s: %v", policyID, err)
		return fmt.Errorf("failed to encrypt claimant ID: %w", err)
	}

	// Step 3: Record the payout distribution in the ledger
	if err := ledgerInstance.DeFiLedger.DistributePayout(policyID, encryptedClaimantID, payoutAmount); err != nil {
		log.Printf("[ERROR] Failed to distribute payout for policy %s: %v", policyID, err)
		return fmt.Errorf("failed to distribute payout: %w", err)
	}

	// Log success
	log.Printf("[SUCCESS] Payout distributed successfully. Policy ID: %s, Claimant ID: %s, Payout Amount: %.2f", policyID, claimantID, payoutAmount)
	return nil
}




// InsuranceTrackPolicyStatus tracks the status of a policy in the ledger.
// Validates input, retrieves the status, and logs the result.
func InsuranceTrackPolicyStatus(policyID string, ledgerInstance *ledger.Ledger) (string, error) {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceTrackPolicyStatus. Policy ID: %s", policyID)

	// Step 1: Validate the policy ID
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return "", err
	}

	// Step 2: Retrieve the policy status from the ledger
	status, err := ledgerInstance.DeFiLedger.TrackPolicyStatus(policyID)
	if err != nil {
		log.Printf("[ERROR] Failed to track status for policy %s: %v", policyID, err)
		return "", fmt.Errorf("failed to track policy status: %w", err)
	}

	// Log success
	log.Printf("[SUCCESS] Policy status retrieved successfully. Policy ID: %s, Status: %s", policyID, status)
	return status, nil
}


// InsuranceAuditPolicy audits a specific policy in the ledger.
// Validates input and logs the audit operation.
func InsuranceAuditPolicy(policyID string, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceAuditPolicy. Policy ID: %s", policyID)

	// Step 1: Validate the policy ID
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Perform the policy audit in the ledger
	startTime := time.Now()
	if err := ledgerInstance.DeFiLedger.AuditPolicy(policyID); err != nil {
		log.Printf("[ERROR] Failed to audit policy %s: %v", policyID, err)
		return fmt.Errorf("failed to audit policy: %w", err)
	}

	// Step 3: Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Policy audited successfully. Policy ID: %s, Duration: %v", policyID, duration)
	return nil
}


// InsuranceEscrowFunds escrows a specified amount of funds for a policy.
// Validates inputs, updates the ledger, and logs the operation.
func InsuranceEscrowFunds(policyID string, amount float64, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceEscrowFunds. Policy ID: %s, Amount: %.2f", policyID, amount)

	// Step 1: Validate input parameters
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}
	if amount <= 0 {
		err := errors.New("amount must be greater than zero")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Perform the escrow operation in the ledger
	startTime := time.Now()
	if err := ledgerInstance.DeFiLedger.EscrowFunds(policyID, amount); err != nil {
		log.Printf("[ERROR] Failed to escrow funds for policy %s: %v", policyID, err)
		return fmt.Errorf("failed to escrow funds: %w", err)
	}

	// Step 3: Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Funds escrowed successfully. Policy ID: %s, Amount: %.2f, Duration: %v", policyID, amount, duration)
	return nil
}


// InsuranceReleaseEscrow releases escrowed funds for a policy.
// Validates the input and updates the ledger.
func InsuranceReleaseEscrow(policyID string, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceReleaseEscrow. Policy ID: %s", policyID)

	// Step 1: Validate the policy ID
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Attempt to release escrowed funds
	startTime := time.Now()
	if err := ledgerInstance.DeFiLedger.ReleaseEscrow(policyID); err != nil {
		log.Printf("[ERROR] Failed to release escrow for policy %s: %v", policyID, err)
		return fmt.Errorf("failed to release escrow: %w", err)
	}

	// Step 3: Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Escrow funds released successfully. Policy ID: %s, Duration: %v", policyID, duration)
	return nil
}


// InsuranceCancelPolicy cancels a specific insurance policy.
// Validates the input and updates the ledger.
func InsuranceCancelPolicy(policyID string, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceCancelPolicy. Policy ID: %s", policyID)

	// Step 1: Validate the policy ID
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Attempt to cancel the policy in the ledger
	startTime := time.Now()
	if err := ledgerInstance.DeFiLedger.CancelPolicy(policyID); err != nil {
		log.Printf("[ERROR] Failed to cancel policy %s: %v", policyID, err)
		return fmt.Errorf("failed to cancel policy: %w", err)
	}

	// Step 3: Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Policy canceled successfully. Policy ID: %s, Duration: %v", policyID, duration)
	return nil
}


// InsuranceLockPolicy locks a specific insurance policy to prevent changes or claims.
// Validates the input and updates the ledger.
func InsuranceLockPolicy(policyID string, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceLockPolicy. Policy ID: %s", policyID)

	// Step 1: Validate the policy ID
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Attempt to lock the policy in the ledger
	startTime := time.Now()
	if err := ledgerInstance.DeFiLedger.LockPolicy(policyID); err != nil {
		log.Printf("[ERROR] Failed to lock policy %s: %v", policyID, err)
		return fmt.Errorf("failed to lock policy: %w", err)
	}

	// Step 3: Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Policy locked successfully. Policy ID: %s, Duration: %v", policyID, duration)
	return nil
}


// InsuranceUnlockPolicy unlocks a policy to allow further modifications or claims.
// Validates input and updates the ledger.
func InsuranceUnlockPolicy(policyID string, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceUnlockPolicy. Policy ID: %s", policyID)

	// Step 1: Validate the policy ID
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Attempt to unlock the policy in the ledger
	startTime := time.Now()
	if err := ledgerInstance.DeFiLedger.UnlockPolicy(policyID); err != nil {
		log.Printf("[ERROR] Failed to unlock policy %s: %v", policyID, err)
		return fmt.Errorf("failed to unlock policy: %w", err)
	}

	// Step 3: Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Policy unlocked successfully. Policy ID: %s, Duration: %v", policyID, duration)
	return nil
}


// InsuranceFetchPolicyTerms retrieves the terms of a policy.
// Validates input, fetches terms from the ledger, and decrypts them.
func InsuranceFetchPolicyTerms(policyID string, ledgerInstance *ledger.Ledger) (string, error) {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceFetchPolicyTerms. Policy ID: %s", policyID)

	// Step 1: Validate the input
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return "", err
	}

	// Step 2: Fetch policy terms from the ledger
	startTime := time.Now()
	terms, err := ledgerInstance.DeFiLedger.FetchPolicyTerms(policyID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch terms for policy %s: %v", policyID, err)
		return "", fmt.Errorf("failed to fetch policy terms: %w", err)
	}
	log.Printf("[INFO] Successfully fetched policy terms from ledger. Policy ID: %s", policyID)

	// Step 3: Decrypt the policy terms
	decryptedTerms, decryptionErr := encryption.DecryptString(terms)
	if decryptionErr != nil {
		log.Printf("[ERROR] Failed to decrypt terms for policy %s: %v", policyID, decryptionErr)
		return "", fmt.Errorf("failed to decrypt policy terms: %w", decryptionErr)
	}

	// Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Policy terms fetched and decrypted successfully. Policy ID: %s, Duration: %v", policyID, duration)

	return decryptedTerms, nil
}


// InsuranceUpdatePolicyTerms updates the terms of a policy.
// Encrypts the new terms, validates input, and updates the ledger.
func InsuranceUpdatePolicyTerms(policyID, newTerms string, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceUpdatePolicyTerms. Policy ID: %s", policyID)

	// Step 1: Validate inputs
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}
	if newTerms == "" {
		err := errors.New("new terms cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Encrypt the new policy terms
	encryptedTerms, encryptionErr := encryption.EncryptString(newTerms)
	if encryptionErr != nil {
		log.Printf("[ERROR] Failed to encrypt new terms for policy %s: %v", policyID, encryptionErr)
		return fmt.Errorf("failed to encrypt policy terms: %w", encryptionErr)
	}

	// Step 3: Update the policy terms in the ledger
	startTime := time.Now()
	if err := ledgerInstance.DeFiLedger.UpdatePolicyTerms(policyID, encryptedTerms); err != nil {
		log.Printf("[ERROR] Failed to update terms for policy %s: %v", policyID, err)
		return fmt.Errorf("failed to update policy terms: %w", err)
	}

	// Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Policy terms updated successfully. Policy ID: %s, Duration: %v", policyID, duration)
	return nil
}


// InsuranceFreezePolicy freezes a policy to temporarily block modifications or claims.
// Validates input and updates the ledger.
func InsuranceFreezePolicy(policyID string, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceFreezePolicy. Policy ID: %s", policyID)

	// Step 1: Validate the input
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Freeze the policy in the ledger
	startTime := time.Now()
	if err := ledgerInstance.DeFiLedger.FreezePolicy(policyID); err != nil {
		log.Printf("[ERROR] Failed to freeze policy %s: %v", policyID, err)
		return fmt.Errorf("failed to freeze policy: %w", err)
	}

	// Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Policy frozen successfully. Policy ID: %s, Duration: %v", policyID, duration)
	return nil
}


// InsuranceFetchClaimHistory retrieves the claim history for a policy.
// Validates input and fetches claim records from the ledger.
func InsuranceFetchClaimHistory(policyID string, ledgerInstance *ledger.Ledger) ([]ledger.ClaimRecord, error) {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceFetchClaimHistory. Policy ID: %s", policyID)

	// Step 1: Validate the input
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// Step 2: Fetch claim history from the ledger
	startTime := time.Now()
	history, err := ledgerInstance.DeFiLedger.GetClaimHistory(policyID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch claim history for policy %s: %v", policyID, err)
		return nil, fmt.Errorf("failed to fetch claim history: %w", err)
	}

	// Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Claim history fetched successfully. Policy ID: %s, Records: %d, Duration: %v", policyID, len(history), duration)

	return history, nil
}


// InsuranceSetClaimProcessingFee sets the processing fee for claims on a policy.
// Validates inputs and updates the ledger.
func InsuranceSetClaimProcessingFee(policyID string, fee float64, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceSetClaimProcessingFee. Policy ID: %s, Fee: %.2f", policyID, fee)

	// Step 1: Validate the inputs
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}
	if fee < 0 {
		err := errors.New("fee cannot be negative")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Update the claim processing fee in the ledger
	startTime := time.Now()
	if err := ledgerInstance.DeFiLedger.SetClaimProcessingFee(policyID, fee); err != nil {
		log.Printf("[ERROR] Failed to set claim processing fee for policy %s: %v", policyID, err)
		return fmt.Errorf("failed to set claim processing fee: %w", err)
	}

	// Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Claim processing fee set successfully. Policy ID: %s, Fee: %.2f, Duration: %v", policyID, fee, duration)
	return nil
}


// InsuranceFetchClaimProcessingFee retrieves the processing fee for claims on a policy.
// Validates input and fetches the fee from the ledger.
func InsuranceFetchClaimProcessingFee(policyID string, ledgerInstance *ledger.Ledger) (float64, error) {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceFetchClaimProcessingFee. Policy ID: %s", policyID)

	// Step 1: Validate the input
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return 0, err
	}

	// Step 2: Fetch the claim processing fee from the ledger
	startTime := time.Now()
	fee, err := ledgerInstance.DeFiLedger.GetClaimProcessingFee(policyID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch claim processing fee for policy %s: %v", policyID, err)
		return 0, fmt.Errorf("failed to fetch claim processing fee: %w", err)
	}

	// Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Claim processing fee fetched successfully. Policy ID: %s, Fee: %.2f, Duration: %v", policyID, fee, duration)
	return fee, nil
}


// InsuranceAutoRenewPolicy enables auto-renewal for a policy.
// Validates input and updates the ledger.
func InsuranceAutoRenewPolicy(policyID string, ledgerInstance *ledger.Ledger) error {
	// Log the start of the operation
	log.Printf("[INFO] Starting InsuranceAutoRenewPolicy. Policy ID: %s", policyID)

	// Step 1: Validate the input
	if policyID == "" {
		err := errors.New("policyID cannot be empty")
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Step 2: Enable auto-renewal in the ledger
	startTime := time.Now()
	if err := ledgerInstance.DeFiLedger.AutoRenewPolicy(policyID); err != nil {
		log.Printf("[ERROR] Failed to enable auto-renewal for policy %s: %v", policyID, err)
		return fmt.Errorf("failed to set up auto-renewal: %w", err)
	}

	// Step 3: Log success
	duration := time.Since(startTime)
	log.Printf("[SUCCESS] Auto-renewal enabled successfully. Policy ID: %s, Duration: %v", policyID, duration)

	return nil
}

