package defi

import (
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"time"
)

// BettingSetOdds sets the odds for a specific bet in a DeFi betting ledger.
// This function validates the bet ID, ensures positive odds, and updates the ledger securely.
func BettingSetOdds(betID string, odds float64, ledgerInstance *ledger.Ledger) error {
	// Ensure thread safety with a mutex for ledger operations
	mutex.Lock()
	defer mutex.Unlock()

	log.Printf("[INFO] Setting odds for bet ID: %s with odds: %.2f", betID, odds)

	// Step 1: Validate inputs
	if betID == "" {
		return fmt.Errorf("bet ID cannot be empty")
	}

	if odds <= 0 {
		return fmt.Errorf("odds must be greater than zero")
	}

	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance provided")
	}

	// Step 2: Check if the bet exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Error checking existence of bet ID %s: %v", betID, err)
		return fmt.Errorf("error validating bet ID: %w", err)
	}

	if !exists {
		return fmt.Errorf("bet ID %s does not exist", betID)
	}

	// Step 3: Update the odds in the ledger
	if err := ledgerInstance.DeFiLedger.SetOdds(betID, odds); err != nil {
		log.Printf("[ERROR] Failed to set odds for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to update odds: %w", err)
	}

	// Step 4: Log and return success
	log.Printf("[SUCCESS] Odds set successfully for bet ID: %s, Odds: %.2f", betID, odds)
	return nil
}

// BettingCalculatePayout calculates the potential payout for a given bet amount and odds.
// Ensures inputs are valid and rounds the result to two decimal places.
func BettingCalculatePayout(amount, odds float64) (float64, error) {
	log.Printf("[INFO] Calculating payout. Amount: %.2f, Odds: %.2f", amount, odds)

	// Step 1: Validate inputs
	if amount <= 0 {
		return 0, fmt.Errorf("bet amount must be greater than zero")
	}

	if odds <= 0 {
		return 0, fmt.Errorf("odds must be greater than zero")
	}

	// Step 2: Calculate the payout
	payout := amount * odds

	// Step 3: Round the payout to two decimal places
	roundedPayout := math.Round(payout*100) / 100

	log.Printf("[SUCCESS] Payout calculated. Amount: %.2f, Odds: %.2f, Payout: %.2f", amount, odds, roundedPayout)
	return roundedPayout, nil
}


// BettingTrackBet tracks a bet for updates or status in the ledger.
// Ensures thread safety, validates inputs, and dynamically handles tracking.
func BettingTrackBet(betID string, ledgerInstance *ledger.Ledger) error {
	// Mutex to ensure thread-safe operations on the ledger
	mutex.Lock()
	defer mutex.Unlock()

	log.Printf("[INFO] Initiating bet tracking for bet ID: %s", betID)

	// Step 1: Validate inputs
	if betID == "" {
		return fmt.Errorf("bet ID cannot be empty")
	}

	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance provided")
	}

	// Step 2: Verify bet ID exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to verify bet ID: %w", err)
	}

	if !exists {
		return fmt.Errorf("bet ID %s does not exist", betID)
	}

	// Step 3: Track the bet in the ledger
	if err := ledgerInstance.DeFiLedger.TrackBet(betID); err != nil {
		log.Printf("[ERROR] Failed to track bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to initiate tracking for bet ID: %w", err)
	}

	// Step 4: Log and return success
	log.Printf("[SUCCESS] Bet tracking initiated successfully for bet ID: %s", betID)
	return nil
}


// BettingFetchBetStatus fetches the current status of a bet from the ledger.
// Validates inputs and dynamically retrieves the status.
func BettingFetchBetStatus(betID string, ledgerInstance *ledger.Ledger) (string, error) {
	log.Printf("[INFO] Fetching status for bet ID: %s", betID)

	// Step 1: Validate inputs
	if betID == "" {
		return "", fmt.Errorf("bet ID cannot be empty")
	}

	if ledgerInstance == nil {
		return "", fmt.Errorf("invalid ledger instance provided")
	}

	// Step 2: Verify bet ID exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify bet ID %s: %v", betID, err)
		return "", fmt.Errorf("failed to verify bet ID: %w", err)
	}

	if !exists {
		return "", fmt.Errorf("bet ID %s does not exist", betID)
	}

	// Step 3: Retrieve the bet status from the ledger
	status, err := ledgerInstance.DeFiLedger.GetBetStatus(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch status for bet ID %s: %v", betID, err)
		return "", fmt.Errorf("failed to fetch bet status: %w", err)
	}

	// Step 4: Log and return the status
	log.Printf("[SUCCESS] Retrieved status for bet ID: %s, Status: %s", betID, status)
	return status, nil
}


// BettingDistributeWinnings distributes winnings for a specific bet.
// Validates the bet ID, fetches bet details, and performs the distribution in the ledger.
func BettingDistributeWinnings(betID string, ledgerInstance *ledger.Ledger) error {
	// Log the start of the process
	log.Printf("[INFO] Starting winnings distribution for bet ID: %s", betID)

	// Step 1: Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance: cannot process winnings distribution")
	}
	if err := validateBetID(betID); err != nil {
		return fmt.Errorf("invalid bet ID: %w", err)
	}

	// Step 2: Check if the bet exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to check existence of bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to verify bet existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("bet ID %s does not exist in the ledger", betID)
	}

	// Step 3: Fetch bet details (ensure the bet is resolved and has winners)
	betDetails, err := ledgerInstance.DeFiLedger.GetBetDetails(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch details for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to fetch bet details: %w", err)
	}
	if !betDetails.IsResolved {
		return fmt.Errorf("bet ID %s is not resolved and cannot distribute winnings", betID)
	}
	if len(betDetails.Winners) == 0 {
		log.Printf("[INFO] No winners for bet ID %s. No distribution required.", betID)
		return nil
	}

	// Step 4: Distribute winnings
	if err := ledgerInstance.DeFiLedger.DistributeWinnings(betID); err != nil {
		log.Printf("[ERROR] Failed to distribute winnings for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to distribute winnings: %w", err)
	}

	// Step 5: Log success
	log.Printf("[SUCCESS] Winnings successfully distributed for bet ID: %s", betID)
	return nil
}


// validateBetID validates the format and ensures the bet ID follows business rules.
func validateBetID(betID string) error {
	log.Printf("[INFO] Validating bet ID: %s", betID)

	// Check if the bet ID is empty
	if len(betID) == 0 {
		return fmt.Errorf("bet ID cannot be empty")
	}

	// Validate the format of the bet ID (e.g., alphanumeric, length constraints)
	if !isValidBetIDFormat(betID) {
		return fmt.Errorf("bet ID format is invalid")
	}

	// Log success
	log.Printf("[SUCCESS] Bet ID validated successfully: %s", betID)
	return nil
}

// isValidBetIDFormat checks the format of the bet ID based on predefined rules.
func isValidBetIDFormat(betID string) bool {
	// Example: Ensure the bet ID is alphanumeric and within a specific length range
	const minLength, maxLength = 8, 64
	if len(betID) < minLength || len(betID) > maxLength {
		return false
	}
	for _, char := range betID {
		if !('a' <= char && char <= 'z') && !('A' <= char && char <= 'Z') && !('0' <= char && char <= '9') {
			return false
		}
	}
	return true
}


// BettingEscrowBetFunds escrows funds for a specific bet.
// Validates inputs and ensures funds are securely escrowed in the ledger.
func BettingEscrowBetFunds(betID string, amount float64, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Starting funds escrow for bet ID: %s, Amount: %.2f", betID, amount)

	// Step 1: Input validation
	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance: cannot escrow funds")
	}
	if err := validateBetID(betID); err != nil {
		return fmt.Errorf("invalid bet ID: %w", err)
	}
	if amount <= 0 {
		return fmt.Errorf("invalid amount: must be greater than zero")
	}

	// Step 2: Verify the bet exists and is in a valid state for escrow
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify bet existence for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to verify bet existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("bet ID %s does not exist in the ledger", betID)
	}

	isValidState, err := ledgerInstance.DeFiLedger.IsBetStateValidForEscrow(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify bet state for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to verify bet state: %w", err)
	}
	if !isValidState {
		return fmt.Errorf("bet ID %s is not in a valid state for escrow", betID)
	}

	// Step 3: Escrow funds in the ledger
	if err := ledgerInstance.DeFiLedger.EscrowFunds(betID, amount); err != nil {
		log.Printf("[ERROR] Failed to escrow funds for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to escrow bet funds: %w", err)
	}

	// Step 4: Log success
	log.Printf("[SUCCESS] Funds escrowed successfully for bet ID: %s, Amount: %.2f", betID, amount)
	return nil
}

// BettingReleaseBetFunds releases escrowed funds for a specific bet.
// Validates the bet ID and ensures the funds are properly released in the ledger.
func BettingReleaseBetFunds(betID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Starting release of escrowed funds for bet ID: %s", betID)

	// Step 1: Input validation
	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance: cannot release escrowed funds")
	}
	if err := validateBetID(betID); err != nil {
		return fmt.Errorf("invalid bet ID: %w", err)
	}

	// Step 2: Verify the bet exists and is in a valid state for fund release
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify bet existence for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to verify bet existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("bet ID %s does not exist in the ledger", betID)
	}

	isValidState, err := ledgerInstance.DeFiLedger.IsBetStateValidForRelease(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify bet state for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to verify bet state: %w", err)
	}
	if !isValidState {
		return fmt.Errorf("bet ID %s is not in a valid state for fund release", betID)
	}

	// Step 3: Release funds in the ledger
	if err := ledgerInstance.DeFiLedger.ReleaseEscrowedFunds(betID); err != nil {
		log.Printf("[ERROR] Failed to release escrowed funds for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to release escrowed funds: %w", err)
	}

	// Step 4: Log success
	log.Printf("[SUCCESS] Escrowed funds released successfully for bet ID: %s", betID)
	return nil
}


// BettingAuditBet audits a specific bet for compliance and accuracy.
// Validates the bet ID and ensures a comprehensive audit is performed in the ledger.
func BettingAuditBet(betID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Starting audit for bet ID: %s", betID)

	// Step 1: Validate the bet ID
	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance: cannot perform bet audit")
	}
	if err := validateBetID(betID); err != nil {
		return fmt.Errorf("invalid bet ID: %w", err)
	}

	// Step 2: Check if the bet exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify existence of bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to verify bet existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("bet ID %s does not exist in the ledger", betID)
	}

	// Step 3: Perform the audit
	log.Printf("[INFO] Performing audit for bet ID: %s", betID)
	if err := ledgerInstance.DeFiLedger.AuditBet(betID); err != nil {
		log.Printf("[ERROR] Audit failed for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to audit bet: %w", err)
	}

	// Step 4: Log audit success
	log.Printf("[SUCCESS] Audit completed successfully for bet ID: %s", betID)
	return nil
}

// BettingMonitorOdds enables monitoring of odds for a specific bet.
// Validates the bet ID and initiates odds monitoring in the ledger.
func BettingMonitorOdds(betID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Starting odds monitoring for bet ID: %s", betID)

	// Step 1: Validate the bet ID
	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance: cannot monitor odds")
	}
	if err := validateBetID(betID); err != nil {
		return fmt.Errorf("invalid bet ID: %w", err)
	}

	// Step 2: Check if the bet exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify existence of bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to verify bet existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("bet ID %s does not exist in the ledger", betID)
	}

	// Step 3: Enable odds monitoring
	log.Printf("[INFO] Enabling odds monitoring for bet ID: %s", betID)
	if err := ledgerInstance.DeFiLedger.MonitorOdds(betID); err != nil {
		log.Printf("[ERROR] Failed to enable odds monitoring for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to monitor odds: %w", err)
	}

	// Step 4: Log success
	log.Printf("[SUCCESS] Odds monitoring enabled successfully for bet ID: %s", betID)
	return nil
}


// BettingFetchBetHistory retrieves the transaction history for a specific bet.
// Validates the bet ID and fetches the history records from the ledger.
func BettingFetchBetHistory(betID string, ledgerInstance *ledger.Ledger) ([]ledger.BetHistoryRecord, error) {
	log.Printf("[INFO] Fetching bet history for bet ID: %s", betID)

	// Step 1: Validate the ledger instance and bet ID
	if ledgerInstance == nil {
		return nil, fmt.Errorf("invalid ledger instance: cannot fetch bet history")
	}
	if err := validateBetID(betID); err != nil {
		return nil, fmt.Errorf("invalid bet ID: %w", err)
	}

	// Step 2: Check if the bet exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify existence of bet ID %s: %v", betID, err)
		return nil, fmt.Errorf("failed to verify bet existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("bet ID %s does not exist in the ledger", betID)
	}

	// Step 3: Fetch the history from the ledger
	log.Printf("[INFO] Fetching transaction history for bet ID: %s", betID)
	history, err := ledgerInstance.DeFiLedger.GetBetHistory(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch bet history for bet ID %s: %v", betID, err)
		return nil, fmt.Errorf("failed to fetch bet history: %w", err)
	}

	// Step 4: Log success and return the history
	log.Printf("[SUCCESS] Bet history retrieved successfully for bet ID: %s", betID)
	return history, nil
}


// BettingSetMaximumBet sets the maximum allowable bet amount for a specific bet.
// Validates the bet ID and ensures the maximum amount is updated in the ledger.
func BettingSetMaximumBet(betID string, maxAmount float64, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Setting maximum bet amount for bet ID: %s", betID)

	// Step 1: Validate the ledger instance and bet ID
	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance: cannot set maximum bet")
	}
	if err := validateBetID(betID); err != nil {
		return fmt.Errorf("invalid bet ID: %w", err)
	}

	// Step 2: Validate the maximum amount
	if maxAmount <= 0 {
		return fmt.Errorf("maximum bet amount must be greater than zero")
	}

	// Step 3: Check if the bet exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify existence of bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to verify bet existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("bet ID %s does not exist in the ledger", betID)
	}

	// Step 4: Update the maximum bet amount in the ledger
	log.Printf("[INFO] Updating maximum bet amount for bet ID: %s to %.2f", betID, maxAmount)
	if err := ledgerInstance.DeFiLedger.SetMaxBet(betID, maxAmount); err != nil {
		log.Printf("[ERROR] Failed to set maximum bet amount for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to set maximum bet amount: %w", err)
	}

	// Step 5: Log success
	log.Printf("[SUCCESS] Maximum bet amount set successfully for bet ID: %s, Maximum Amount: %.2f", betID, maxAmount)
	return nil
}


// BettingFetchMaximumBet retrieves the maximum allowable bet amount for a specific bet.
// Validates the bet ID and fetches the value from the ledger.
func BettingFetchMaximumBet(betID string, ledgerInstance *ledger.Ledger) (float64, error) {
	log.Printf("[INFO] Fetching maximum bet amount for bet ID: %s", betID)

	// Step 1: Validate inputs
	if ledgerInstance == nil {
		return 0, fmt.Errorf("invalid ledger instance: cannot fetch maximum bet amount")
	}
	if err := validateBetID(betID); err != nil {
		return 0, fmt.Errorf("invalid bet ID: %w", err)
	}

	// Step 2: Check if the bet exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify existence of bet ID %s: %v", betID, err)
		return 0, fmt.Errorf("failed to verify bet existence: %w", err)
	}
	if !exists {
		return 0, fmt.Errorf("bet ID %s does not exist in the ledger", betID)
	}

	// Step 3: Retrieve the maximum bet amount from the ledger
	log.Printf("[INFO] Retrieving maximum bet amount for bet ID: %s", betID)
	maxBet, err := ledgerInstance.DeFiLedger.GetMaxBet(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch maximum bet amount for bet ID %s: %v", betID, err)
		return 0, fmt.Errorf("failed to fetch maximum bet amount: %w", err)
	}

	// Step 4: Log success and return the value
	log.Printf("[SUCCESS] Maximum bet amount retrieved for bet ID: %s, Maximum Amount: %.2f", betID, maxBet)
	return maxBet, nil
}


// BettingSetBetExpiration sets an expiration time for a specific bet.
// Validates the bet ID and expiration time, and updates the ledger.
func BettingSetBetExpiration(betID string, expiration time.Time, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Setting expiration for bet ID: %s", betID)

	// Step 1: Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance: cannot set bet expiration")
	}
	if err := validateBetID(betID); err != nil {
		return fmt.Errorf("invalid bet ID: %w", err)
	}
	if time.Now().After(expiration) {
		return fmt.Errorf("expiration time must be in the future")
	}

	// Step 2: Check if the bet exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify existence of bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to verify bet existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("bet ID %s does not exist in the ledger", betID)
	}

	// Step 3: Set the expiration time in the ledger
	log.Printf("[INFO] Setting expiration time for bet ID: %s to %v", betID, expiration)
	if err := ledgerInstance.DeFiLedger.SetBetExpiration(betID, expiration); err != nil {
		log.Printf("[ERROR] Failed to set expiration for bet ID %s: %v", betID, err)
		return fmt.Errorf("failed to set bet expiration: %w", err)
	}

	// Step 4: Log success
	log.Printf("[SUCCESS] Bet expiration set successfully for bet ID: %s, Expiration: %v", betID, expiration)
	return nil
}

// BettingFetchBetExpiration retrieves the expiration time for a specific bet.
// Validates the bet ID and fetches the expiration time from the ledger.
func BettingFetchBetExpiration(betID string, ledgerInstance *ledger.Ledger) (time.Time, error) {
	log.Printf("[INFO] Fetching expiration time for bet ID: %s", betID)

	// Step 1: Validate inputs
	if ledgerInstance == nil {
		return time.Time{}, fmt.Errorf("invalid ledger instance: cannot fetch bet expiration")
	}
	if err := validateBetID(betID); err != nil {
		return time.Time{}, fmt.Errorf("invalid bet ID: %w", err)
	}

	// Step 2: Check if the bet exists in the ledger
	exists, err := ledgerInstance.DeFiLedger.DoesBetExist(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify existence of bet ID %s: %v", betID, err)
		return time.Time{}, fmt.Errorf("failed to verify bet existence: %w", err)
	}
	if !exists {
		return time.Time{}, fmt.Errorf("bet ID %s does not exist in the ledger", betID)
	}

	// Step 3: Fetch the expiration time from the ledger
	expiration, err := ledgerInstance.DeFiLedger.GetBetExpiration(betID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch expiration time for bet ID %s: %v", betID, err)
		return time.Time{}, fmt.Errorf("failed to fetch bet expiration: %w", err)
	}

	// Step 4: Log success and return the expiration time
	log.Printf("[SUCCESS] Bet expiration retrieved for bet ID: %s, Expiration: %v", betID, expiration)
	return expiration, nil
}


// BettingPauseBetting pauses all betting operations on the platform.
// Updates the configuration in the ledger to reflect the paused state.
func BettingPauseBetting(ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Initiating pause for all betting operations.")

	// Step 1: Validate the ledger instance
	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance: cannot pause betting")
	}

	// Step 2: Update the betting configuration in the ledger
	log.Printf("[INFO] Updating betting configuration to paused state.")
	if err := ledgerInstance.DeFiLedger.UpdateConfig("BettingPaused", true); err != nil {
		log.Printf("[ERROR] Failed to update betting configuration: %v", err)
		return fmt.Errorf("failed to pause betting: %w", err)
	}

	// Step 3: Verify the update
	paused, err := ledgerInstance.DeFiLedger.GetConfig("BettingPaused")
	if err != nil || paused != true {
		log.Printf("[ERROR] Verification failed for betting pause state: %v", err)
		return fmt.Errorf("failed to verify betting pause state: %w", err)
	}

	// Step 4: Log success
	log.Printf("[SUCCESS] Betting operations paused successfully.")
	return nil
}


// BettingResumeBetting resumes all betting operations on the platform.
// Updates the configuration in the ledger to reflect the active state.
func BettingResumeBetting(ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Initiating resume of all betting operations.")

	// Step 1: Validate the ledger instance
	if ledgerInstance == nil {
		return fmt.Errorf("invalid ledger instance: cannot resume betting")
	}

	// Step 2: Fetch the current betting configuration from the ledger
	log.Printf("[INFO] Checking current betting configuration...")
	bettingPaused, err := ledgerInstance.DeFiLedger.GetConfig("BettingPaused")
	if err != nil {
		log.Printf("[ERROR] Failed to fetch current betting configuration: %v", err)
		return fmt.Errorf("failed to fetch betting configuration: %w", err)
	}

	// Step 3: Check if betting is already active
	if !bettingPaused {
		log.Printf("[INFO] Betting is already active. No action required.")
		return nil
	}

	// Step 4: Update the betting configuration to active
	log.Printf("[INFO] Updating betting configuration to active state...")
	if err := ledgerInstance.DeFiLedger.UpdateConfig("BettingPaused", false); err != nil {
		log.Printf("[ERROR] Failed to update betting configuration to active state: %v", err)
		return fmt.Errorf("failed to resume betting: %w", err)
	}

	// Step 5: Verify the configuration update
	log.Printf("[INFO] Verifying betting configuration update...")
	newConfig, err := ledgerInstance.DeFiLedger.GetConfig("BettingPaused")
	if err != nil || newConfig != false {
		log.Printf("[ERROR] Betting configuration verification failed: %v", err)
		return fmt.Errorf("failed to verify betting resumption: %w", err)
	}

	// Step 6: Log success and exit
	log.Printf("[SUCCESS] Betting operations resumed successfully.")
	return nil
}

