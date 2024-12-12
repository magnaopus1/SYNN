package defi

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"time"
)

// CrowdfundingCreateCampaign creates a new crowdfunding campaign.
// Validates inputs, encrypts sensitive data, and stores the campaign in the ledger.
func CrowdfundingCreateCampaign(campaignID, title, description string, goalAmount float64, endTime time.Time, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Starting campaign creation process for Campaign ID: %s", campaignID)

	// Step 1: Validate input parameters
	if campaignID == "" {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format")
	}
	if title == "" {
		return fmt.Errorf("campaign title cannot be empty")
	}
	if description == "" {
		return fmt.Errorf("campaign description cannot be empty")
	}
	if goalAmount <= 0 {
		return fmt.Errorf("goal amount must be greater than zero")
	}
	if time.Now().After(endTime) {
		return fmt.Errorf("end time must be in the future")
	}

	// Step 2: Encrypt sensitive data
	log.Printf("[INFO] Encrypting sensitive campaign data...")
	encryptedTitle := encryption.EncryptString(title)
	encryptedDescription := encryption.EncryptString(description)

	// Step 3: Create a campaign object
	campaign := ledger.CrowdfundingCampaign{
		CampaignID:  campaignID,
		Title:       encryptedTitle,
		Description: encryptedDescription,
		GoalAmount:  goalAmount,
		EndTime:     endTime,
		Status:      "Active",
	}

	// Step 4: Store the campaign in the ledger
	log.Printf("[INFO] Storing campaign in ledger...")
	if err := ledgerInstance.DeFiLedger.CreateCampaign(campaign); err != nil {
		log.Printf("[ERROR] Failed to create campaign with ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to create campaign: %w", err)
	}

	// Step 5: Log and confirm success
	log.Printf("[SUCCESS] Crowdfunding campaign created successfully. Campaign ID: %s", campaignID)
	return nil
}


// CrowdfundingContribute allows a user to contribute to a crowdfunding campaign.
// Validates inputs, encrypts the user ID, and records the contribution in the ledger.
func CrowdfundingContribute(campaignID, userID string, amount float64, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Starting contribution process for Campaign ID: %s, User ID: %s", campaignID, userID)

	// Step 1: Validate input parameters
	if campaignID == "" {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format")
	}
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if amount <= 0 {
		return fmt.Errorf("contribution amount must be greater than zero")
	}

	// Step 2: Check if the campaign exists and is active
	log.Printf("[INFO] Validating campaign status...")
	campaign, err := ledgerInstance.DeFiLedger.GetCampaignByID(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch campaign with ID %s: %v", campaignID, err)
		return fmt.Errorf("campaign not found: %w", err)
	}
	if campaign.Status != "Active" {
		log.Printf("[ERROR] Contribution rejected: Campaign ID %s is not active", campaignID)
		return fmt.Errorf("campaign is not active")
	}

	// Step 3: Encrypt sensitive data
	log.Printf("[INFO] Encrypting user ID for security...")
	encryptedUserID := encryption.EncryptString(userID)

	// Step 4: Record the contribution in the ledger
	log.Printf("[INFO] Recording contribution in ledger...")
	if err := ledgerInstance.DeFiLedger.Contribute(campaignID, encryptedUserID, amount); err != nil {
		log.Printf("[ERROR] Failed to record contribution for Campaign ID %s by User ID %s: %v", campaignID, userID, err)
		return fmt.Errorf("failed to contribute to campaign: %w", err)
	}

	// Step 5: Log and confirm success
	log.Printf("[SUCCESS] Contribution recorded successfully. Campaign ID: %s, User ID: %s, Amount: %.2f", campaignID, userID, amount)
	return nil
}

// CrowdfundingRefundContributors refunds all contributors of a specific crowdfunding campaign.
// Validates the campaign ID, checks campaign state, and processes refunds through the ledger.
func CrowdfundingRefundContributors(campaignID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Starting refund process for Campaign ID: %s", campaignID)

	// Step 1: Validate campaign ID
	if campaignID == "" {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format")
	}

	// Step 2: Check campaign existence and state
	log.Printf("[INFO] Validating campaign state for refunds...")
	campaign, err := ledgerInstance.DeFiLedger.GetCampaignByID(campaignID)
	if err != nil {
		log.Printf("[ERROR] Campaign with ID %s not found: %v", campaignID, err)
		return fmt.Errorf("campaign not found: %w", err)
	}
	if campaign.Status != "Failed" {
		log.Printf("[ERROR] Refunds cannot be processed: Campaign ID %s is not in 'Failed' state", campaignID)
		return fmt.Errorf("refunds are only allowed for failed campaigns")
	}

	// Step 3: Process refunds in the ledger
	log.Printf("[INFO] Processing refunds for contributors of Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.RefundContributors(campaignID); err != nil {
		log.Printf("[ERROR] Failed to process refunds for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to refund contributors: %w", err)
	}

	// Step 4: Update campaign status to 'Refunded'
	log.Printf("[INFO] Updating campaign status to 'Refunded' for Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.UpdateCampaignStatus(campaignID, "Refunded"); err != nil {
		log.Printf("[ERROR] Failed to update campaign status for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to update campaign status: %w", err)
	}

	// Step 5: Log and confirm success
	log.Printf("[SUCCESS] All contributors refunded successfully for Campaign ID: %s", campaignID)
	return nil
}


// CrowdfundingDistributeFunds distributes the raised funds to the campaign creator.
// Validates the campaign ID, ensures the campaign succeeded, and processes the distribution through the ledger.
func CrowdfundingDistributeFunds(campaignID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Starting fund distribution for Campaign ID: %s", campaignID)

	// Step 1: Validate campaign ID
	if campaignID == "" {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format")
	}

	// Step 2: Check campaign existence and state
	log.Printf("[INFO] Validating campaign state for fund distribution...")
	campaign, err := ledgerInstance.DeFiLedger.GetCampaignByID(campaignID)
	if err != nil {
		log.Printf("[ERROR] Campaign with ID %s not found: %v", campaignID, err)
		return fmt.Errorf("campaign not found: %w", err)
	}
	if campaign.Status != "Successful" {
		log.Printf("[ERROR] Funds cannot be distributed: Campaign ID %s is not in 'Successful' state", campaignID)
		return fmt.Errorf("funds can only be distributed for successful campaigns")
	}

	// Step 3: Process fund distribution in the ledger
	log.Printf("[INFO] Distributing funds to campaign creator for Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.DistributeFunds(campaignID); err != nil {
		log.Printf("[ERROR] Failed to distribute funds for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to distribute funds: %w", err)
	}

	// Step 4: Update campaign status to 'Distributed'
	log.Printf("[INFO] Updating campaign status to 'Distributed' for Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.UpdateCampaignStatus(campaignID, "Distributed"); err != nil {
		log.Printf("[ERROR] Failed to update campaign status for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to update campaign status: %w", err)
	}

	// Step 5: Log and confirm success
	log.Printf("[SUCCESS] Funds distributed successfully to the creator of Campaign ID: %s", campaignID)
	return nil
}


// isValidCampaignIDFormat checks if the campaign ID matches the expected format.
func isValidCampaignIDFormat(campaignID string) bool {
	if len(campaignID) == 0 {
		log.Printf("[ERROR] Campaign ID cannot be empty.")
		return false
	}

	// Define a regex for a valid campaign ID (alphanumeric with hyphens).
	regex := `^[a-zA-Z0-9\-]+$`
	matched, err := regexp.MatchString(regex, campaignID)
	if err != nil {
		log.Printf("[ERROR] Regex match failed for campaign ID %s: %v", campaignID, err)
		return false
	}

	if !matched {
		log.Printf("[ERROR] Campaign ID %s is not in the valid format.", campaignID)
	}

	return matched
}

// CrowdfundingAuditCampaign audits a specific crowdfunding campaign.
// Validates the campaign ID, ensures the campaign exists, and performs an audit via the ledger.
func CrowdfundingAuditCampaign(campaignID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Initiating audit for Campaign ID: %s", campaignID)

	// Step 1: Validate campaign ID
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format for ID: %s", campaignID)
	}

	// Step 2: Check if the campaign exists
	campaign, err := ledgerInstance.DeFiLedger.GetCampaignByID(campaignID)
	if err != nil {
		log.Printf("[ERROR] Campaign with ID %s not found: %v", campaignID, err)
		return fmt.Errorf("campaign not found: %w", err)
	}

	// Step 3: Ensure campaign is in an auditable state
	if campaign.Status != "Active" && campaign.Status != "Completed" {
		log.Printf("[ERROR] Campaign ID %s is not in an auditable state. Status: %s", campaignID, campaign.Status)
		return fmt.Errorf("campaign must be in 'Active' or 'Completed' state for auditing")
	}

	// Step 4: Perform the audit via the ledger
	log.Printf("[INFO] Performing audit for Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.AuditCampaign(campaignID); err != nil {
		log.Printf("[ERROR] Failed to audit Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("audit failed: %w", err)
	}

	// Step 5: Log and confirm success
	log.Printf("[SUCCESS] Campaign audited successfully. Campaign ID: %s", campaignID)
	return nil
}


// CrowdfundingTrackCampaignProgress tracks the progress of a specific crowdfunding campaign.
// Validates the campaign ID, ensures the campaign exists, and updates progress via the ledger.
func CrowdfundingTrackCampaignProgress(campaignID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Initiating progress tracking for Campaign ID: %s", campaignID)

	// Step 1: Validate campaign ID
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format for ID: %s", campaignID)
	}

	// Step 2: Check if the campaign exists
	campaign, err := ledgerInstance.DeFiLedger.GetCampaignByID(campaignID)
	if err != nil {
		log.Printf("[ERROR] Campaign with ID %s not found: %v", campaignID, err)
		return fmt.Errorf("campaign not found: %w", err)
	}

	// Step 3: Ensure campaign is in a trackable state
	if campaign.Status != "Active" {
		log.Printf("[ERROR] Campaign ID %s is not in a trackable state. Status: %s", campaignID, campaign.Status)
		return fmt.Errorf("campaign must be in 'Active' state for progress tracking")
	}

	// Step 4: Track the campaign progress via the ledger
	log.Printf("[INFO] Tracking progress for Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.TrackCampaignProgress(campaignID); err != nil {
		log.Printf("[ERROR] Failed to track progress for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("progress tracking failed: %w", err)
	}

	// Step 5: Log and confirm success
	log.Printf("[SUCCESS] Campaign progress tracked successfully. Campaign ID: %s", campaignID)
	return nil
}


// CrowdfundingFetchCampaignDetails retrieves the details of a specific crowdfunding campaign.
// Validates the campaign ID, ensures the campaign exists, and fetches details from the ledger.
func CrowdfundingFetchCampaignDetails(campaignID string, ledgerInstance *ledger.Ledger) (ledger.CrowdfundingCampaign, error) {
	log.Printf("[INFO] Fetching details for Campaign ID: %s", campaignID)

	// Step 1: Validate campaign ID
	if len(campaignID) == 0 {
		return ledger.CrowdfundingCampaign{}, fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return ledger.CrowdfundingCampaign{}, fmt.Errorf("invalid campaign ID format for ID: %s", campaignID)
	}

	// Step 2: Fetch campaign details from the ledger
	log.Printf("[INFO] Retrieving campaign details from ledger for Campaign ID: %s", campaignID)
	campaign, err := ledgerInstance.DeFiLedger.FetchCampaignDetails(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch details for Campaign ID %s: %v", campaignID, err)
		return ledger.CrowdfundingCampaign{}, fmt.Errorf("failed to fetch campaign details: %w", err)
	}

	// Step 3: Log success and return the campaign details
	log.Printf("[SUCCESS] Campaign details retrieved successfully. Campaign ID: %s", campaignID)
	return campaign, nil
}


// CrowdfundingCloseCampaign closes a specific crowdfunding campaign.
// Validates the campaign ID, ensures the campaign is eligible for closure, and performs the operation via the ledger.
func CrowdfundingCloseCampaign(campaignID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Initiating closure for Campaign ID: %s", campaignID)

	// Step 1: Validate campaign ID
	if len(campaignID) == 0 {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format for ID: %s", campaignID)
	}

	// Step 2: Fetch the campaign to verify its state
	log.Printf("[INFO] Fetching campaign to verify closure eligibility. Campaign ID: %s", campaignID)
	campaign, err := ledgerInstance.DeFiLedger.FetchCampaignDetails(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve campaign details for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to fetch campaign details: %w", err)
	}

	// Step 3: Ensure the campaign is eligible for closure
	if campaign.Status != "Active" && campaign.Status != "Completed" {
		log.Printf("[ERROR] Campaign ID %s is not in a closeable state. Status: %s", campaignID, campaign.Status)
		return fmt.Errorf("campaign must be in 'Active' or 'Completed' state for closure")
	}

	// Step 4: Perform the closure operation in the ledger
	log.Printf("[INFO] Closing campaign in the ledger. Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.CloseCampaign(campaignID); err != nil {
		log.Printf("[ERROR] Failed to close Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to close campaign: %w", err)
	}

	// Step 5: Log and confirm success
	log.Printf("[SUCCESS] Campaign closed successfully. Campaign ID: %s", campaignID)
	return nil
}


// CrowdfundingLockFunds locks a specific amount of funds for a crowdfunding campaign.
// Validates inputs and securely locks funds in the ledger.
func CrowdfundingLockFunds(campaignID string, amount float64, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Initiating fund lock for Campaign ID: %s, Amount: %.2f", campaignID, amount)

	// Step 1: Validate Campaign ID
	if len(campaignID) == 0 {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format for ID: %s", campaignID)
	}

	// Step 2: Validate Amount
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero, received: %.2f", amount)
	}

	// Step 3: Ensure the campaign exists and is active
	log.Printf("[INFO] Validating campaign state for Campaign ID: %s", campaignID)
	campaign, err := ledgerInstance.DeFiLedger.FetchCampaignDetails(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch campaign details for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("campaign does not exist or could not be retrieved: %w", err)
	}
	if campaign.Status != "Active" {
		log.Printf("[ERROR] Campaign ID %s is not in an active state. Current Status: %s", campaignID, campaign.Status)
		return fmt.Errorf("funds can only be locked for active campaigns")
	}

	// Step 4: Lock Funds in the Ledger
	log.Printf("[INFO] Locking funds in the ledger for Campaign ID: %s, Amount: %.2f", campaignID, amount)
	if err := ledgerInstance.DeFiLedger.LockFunds(campaignID, amount); err != nil {
		log.Printf("[ERROR] Failed to lock funds for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to lock funds: %w", err)
	}

	// Step 5: Log and Confirm Success
	log.Printf("[SUCCESS] Funds locked successfully for Campaign ID: %s, Amount: %.2f", campaignID, amount)
	return nil
}


// CrowdfundingUnlockFunds unlocks all funds for a crowdfunding campaign.
// Validates the campaign ID, ensures the campaign is eligible, and unlocks funds via the ledger.
func CrowdfundingUnlockFunds(campaignID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Initiating fund unlock for Campaign ID: %s", campaignID)

	// Step 1: Validate Campaign ID
	if len(campaignID) == 0 {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format for ID: %s", campaignID)
	}

	// Step 2: Fetch the Campaign to Ensure it is Eligible
	log.Printf("[INFO] Validating campaign state for unlocking funds. Campaign ID: %s", campaignID)
	campaign, err := ledgerInstance.DeFiLedger.FetchCampaignDetails(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve campaign details for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to retrieve campaign details: %w", err)
	}
	if campaign.Status != "Completed" && campaign.Status != "Cancelled" {
		log.Printf("[ERROR] Campaign ID %s is not eligible for fund unlocking. Current Status: %s", campaignID, campaign.Status)
		return fmt.Errorf("funds can only be unlocked for completed or cancelled campaigns")
	}

	// Step 3: Unlock Funds in the Ledger
	log.Printf("[INFO] Unlocking funds in the ledger for Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.UnlockFunds(campaignID); err != nil {
		log.Printf("[ERROR] Failed to unlock funds for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to unlock funds: %w", err)
	}

	// Step 4: Log and Confirm Success
	log.Printf("[SUCCESS] Funds unlocked successfully for Campaign ID: %s", campaignID)
	return nil
}


// CrowdfundingEscrowFunds securely escrows a specific amount of funds for a crowdfunding campaign.
// Validates inputs, ensures campaign eligibility, and records the escrow operation in the ledger.
func CrowdfundingEscrowFunds(campaignID string, amount float64, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Starting escrow process for Campaign ID: %s, Amount: %.2f", campaignID, amount)

	// Step 1: Validate Campaign ID
	if len(campaignID) == 0 {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format: %s", campaignID)
	}

	// Step 2: Validate Amount
	if amount <= 0 {
		return fmt.Errorf("escrow amount must be greater than zero, received: %.2f", amount)
	}

	// Step 3: Verify Campaign Status
	log.Printf("[INFO] Validating campaign eligibility for Campaign ID: %s", campaignID)
	campaign, err := ledgerInstance.DeFiLedger.FetchCampaignDetails(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch campaign details for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("unable to fetch campaign details: %w", err)
	}
	if campaign.Status != "Active" {
		log.Printf("[ERROR] Campaign ID %s is not in an active state. Current Status: %s", campaignID, campaign.Status)
		return fmt.Errorf("escrow is only allowed for active campaigns")
	}

	// Step 4: Escrow Funds in the Ledger
	log.Printf("[INFO] Escrowing funds in the ledger for Campaign ID: %s, Amount: %.2f", campaignID, amount)
	if err := ledgerInstance.DeFiLedger.EscrowFunds(campaignID, amount); err != nil {
		log.Printf("[ERROR] Failed to escrow funds for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to escrow funds: %w", err)
	}

	// Step 5: Log Success
	log.Printf("[SUCCESS] Funds escrowed successfully for Campaign ID: %s, Amount: %.2f", campaignID, amount)
	return nil
}


// CrowdfundingExtendCampaignDuration extends the duration of a crowdfunding campaign.
// Validates inputs, ensures campaign eligibility, and updates the end time in the ledger.
func CrowdfundingExtendCampaignDuration(campaignID string, additionalTime time.Duration, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Initiating campaign duration extension for Campaign ID: %s, Additional Time: %s", campaignID, additionalTime)

	// Step 1: Validate Campaign ID
	if len(campaignID) == 0 {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format: %s", campaignID)
	}

	// Step 2: Validate Additional Time
	if additionalTime <= 0 {
		return fmt.Errorf("additional time must be greater than zero, received: %s", additionalTime)
	}

	// Step 3: Fetch Campaign Details and Verify Eligibility
	log.Printf("[INFO] Fetching campaign details for validation. Campaign ID: %s", campaignID)
	campaign, err := ledgerInstance.DeFiLedger.FetchCampaignDetails(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch campaign details for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("unable to fetch campaign details: %w", err)
	}
	if campaign.Status != "Active" {
		log.Printf("[ERROR] Campaign ID %s is not in an active state. Current Status: %s", campaignID, campaign.Status)
		return fmt.Errorf("duration can only be extended for active campaigns")
	}

	// Step 4: Calculate the New End Time
	newEndTime := campaign.EndTime.Add(additionalTime)
	if time.Now().After(newEndTime) {
		log.Printf("[ERROR] Calculated new end time is in the past. New End Time: %v", newEndTime)
		return fmt.Errorf("calculated new end time must be in the future")
	}

	// Step 5: Update the Campaign Duration in the Ledger
	log.Printf("[INFO] Updating campaign duration in the ledger. Campaign ID: %s, New End Time: %v", campaignID, newEndTime)
	if err := ledgerInstance.DeFiLedger.ExtendCampaignDuration(campaignID, newEndTime); err != nil {
		log.Printf("[ERROR] Failed to extend campaign duration for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to update campaign duration: %w", err)
	}

	// Step 6: Log Success
	log.Printf("[SUCCESS] Campaign duration extended successfully. Campaign ID: %s, New End Time: %v", campaignID, newEndTime)
	return nil
}


// CrowdfundingFetchContributionHistory retrieves the contribution history of a crowdfunding campaign.
// Validates the campaign ID, fetches the history from the ledger, and ensures secure data handling.
func CrowdfundingFetchContributionHistory(campaignID string, ledgerInstance *ledger.Ledger) ([]ledger.ContributionRecord, error) {
	log.Printf("[INFO] Fetching contribution history for Campaign ID: %s", campaignID)

	// Step 1: Validate Campaign ID
	if len(campaignID) == 0 {
		return nil, fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return nil, fmt.Errorf("invalid campaign ID format: %s", campaignID)
	}

	// Step 2: Fetch Contribution History
	log.Printf("[INFO] Retrieving contribution history from ledger for Campaign ID: %s", campaignID)
	history, err := ledgerInstance.DeFiLedger.GetContributionHistory(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch contribution history for Campaign ID %s: %v", campaignID, err)
		return nil, fmt.Errorf("unable to fetch contribution history: %w", err)
	}

	// Step 3: Log Success
	log.Printf("[SUCCESS] Contribution history retrieved successfully for Campaign ID: %s", campaignID)
	return history, nil
}


// CrowdfundingSetContributionLimits sets the minimum and maximum contribution limits for a crowdfunding campaign.
// Validates inputs, ensures campaign eligibility, and updates the limits in the ledger.
func CrowdfundingSetContributionLimits(campaignID string, minLimit, maxLimit float64, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Setting contribution limits for Campaign ID: %s, Min Limit: %.2f, Max Limit: %.2f", campaignID, minLimit, maxLimit)

	// Step 1: Validate Campaign ID
	if len(campaignID) == 0 {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format: %s", campaignID)
	}

	// Step 2: Validate Contribution Limits
	if minLimit < 0 || maxLimit <= 0 || minLimit > maxLimit {
		return fmt.Errorf("invalid contribution limits: minLimit must be >= 0, maxLimit must be > 0, and minLimit <= maxLimit")
	}

	// Step 3: Verify Campaign Status
	log.Printf("[INFO] Verifying campaign status for Campaign ID: %s", campaignID)
	campaign, err := ledgerInstance.DeFiLedger.FetchCampaignDetails(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch campaign details for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("unable to fetch campaign details: %w", err)
	}
	if campaign.Status != "Active" {
		log.Printf("[ERROR] Campaign ID %s is not active. Current Status: %s", campaignID, campaign.Status)
		return fmt.Errorf("contribution limits can only be set for active campaigns")
	}

	// Step 4: Update Contribution Limits in the Ledger
	log.Printf("[INFO] Updating contribution limits in the ledger for Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.SetContributionLimits(campaignID, minLimit, maxLimit); err != nil {
		log.Printf("[ERROR] Failed to set contribution limits for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to set contribution limits: %w", err)
	}

	// Step 5: Log Success
	log.Printf("[SUCCESS] Contribution limits set successfully for Campaign ID: %s, Min Limit: %.2f, Max Limit: %.2f", campaignID, minLimit, maxLimit)
	return nil
}


// CrowdfundingFetchContributionLimits fetches the minimum and maximum contribution limits for a campaign.
// Validates the campaign ID, retrieves the limits from the ledger, and logs results dynamically.
func CrowdfundingFetchContributionLimits(campaignID string, ledgerInstance *ledger.Ledger) (float64, float64, error) {
	log.Printf("[INFO] Fetching contribution limits for Campaign ID: %s", campaignID)

	// Step 1: Validate Campaign ID
	if len(campaignID) == 0 {
		return 0, 0, fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return 0, 0, fmt.Errorf("invalid campaign ID format: %s", campaignID)
	}

	// Step 2: Fetch Contribution Limits
	log.Printf("[INFO] Retrieving contribution limits from ledger for Campaign ID: %s", campaignID)
	minLimit, maxLimit, err := ledgerInstance.DeFiLedger.GetContributionLimits(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch contribution limits for Campaign ID %s: %v", campaignID, err)
		return 0, 0, fmt.Errorf("unable to fetch contribution limits: %w", err)
	}

	// Step 3: Log Success
	log.Printf("[SUCCESS] Contribution limits retrieved successfully for Campaign ID: %s, Min Limit: %.2f, Max Limit: %.2f", campaignID, minLimit, maxLimit)
	return minLimit, maxLimit, nil
}


// CrowdfundingPauseCampaign pauses a crowdfunding campaign.
// Validates the campaign ID, checks the campaign status, and updates the status to paused in the ledger.
func CrowdfundingPauseCampaign(campaignID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Pausing crowdfunding campaign with Campaign ID: %s", campaignID)

	// Step 1: Validate Campaign ID
	if len(campaignID) == 0 {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format: %s", campaignID)
	}

	// Step 2: Fetch Campaign Details and Validate Status
	log.Printf("[INFO] Verifying campaign status for Campaign ID: %s", campaignID)
	campaign, err := ledgerInstance.DeFiLedger.FetchCampaignDetails(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch campaign details for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("unable to fetch campaign details: %w", err)
	}
	if campaign.Status != "Active" {
		log.Printf("[ERROR] Campaign ID %s cannot be paused as it is not active. Current Status: %s", campaignID, campaign.Status)
		return fmt.Errorf("campaign must be active to pause")
	}

	// Step 3: Update Campaign Status to Paused
	log.Printf("[INFO] Updating campaign status to 'Paused' for Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.PauseCampaign(campaignID); err != nil {
		log.Printf("[ERROR] Failed to pause campaign with Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to pause campaign: %w", err)
	}

	// Step 4: Log Success
	log.Printf("[SUCCESS] Campaign paused successfully. Campaign ID: %s", campaignID)
	return nil
}


// CrowdfundingResumeCampaign resumes a paused crowdfunding campaign.
// Validates the campaign ID, checks the campaign's current status, and updates the status to active in the ledger.
func CrowdfundingResumeCampaign(campaignID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Attempting to resume crowdfunding campaign. Campaign ID: %s", campaignID)

	// Step 1: Validate Campaign ID
	if len(campaignID) == 0 {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format: %s", campaignID)
	}

	// Step 2: Fetch Campaign Details and Validate Status
	log.Printf("[INFO] Fetching campaign details for verification. Campaign ID: %s", campaignID)
	campaign, err := ledgerInstance.DeFiLedger.FetchCampaignDetails(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch campaign details for Campaign ID %s: %v", campaignID, err)
		return fmt.Errorf("failed to fetch campaign details: %w", err)
	}
	if campaign.Status != "Paused" {
		log.Printf("[ERROR] Campaign ID %s cannot be resumed. Current Status: %s", campaignID, campaign.Status)
		return fmt.Errorf("campaign must be paused to resume")
	}

	// Step 3: Update Campaign Status to Active
	log.Printf("[INFO] Resuming campaign. Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.ResumeCampaign(campaignID); err != nil {
		log.Printf("[ERROR] Failed to resume campaign. Campaign ID: %s, Error: %v", campaignID, err)
		return fmt.Errorf("failed to resume campaign: %w", err)
	}

	// Step 4: Log Success
	log.Printf("[SUCCESS] Campaign resumed successfully. Campaign ID: %s", campaignID)
	return nil
}


// CrowdfundingMonitorContributionFlow monitors the contribution flow for a crowdfunding campaign.
// Validates the campaign ID and enables contribution flow monitoring in the ledger.
func CrowdfundingMonitorContributionFlow(campaignID string, ledgerInstance *ledger.Ledger) error {
	log.Printf("[INFO] Enabling contribution flow monitoring. Campaign ID: %s", campaignID)

	// Step 1: Validate Campaign ID
	if len(campaignID) == 0 {
		return fmt.Errorf("campaign ID cannot be empty")
	}
	if !isValidCampaignIDFormat(campaignID) {
		return fmt.Errorf("invalid campaign ID format: %s", campaignID)
	}

	// Step 2: Verify Campaign Exists
	log.Printf("[INFO] Validating campaign existence for Campaign ID: %s", campaignID)
	_, err := ledgerInstance.DeFiLedger.FetchCampaignDetails(campaignID)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch campaign details for monitoring. Campaign ID: %s, Error: %v", campaignID, err)
		return fmt.Errorf("failed to verify campaign for monitoring: %w", err)
	}

	// Step 3: Enable Contribution Flow Monitoring
	log.Printf("[INFO] Activating contribution flow monitoring for Campaign ID: %s", campaignID)
	if err := ledgerInstance.DeFiLedger.MonitorContributionFlow(campaignID); err != nil {
		log.Printf("[ERROR] Failed to enable contribution flow monitoring. Campaign ID: %s, Error: %v", campaignID, err)
		return fmt.Errorf("failed to monitor contribution flow: %w", err)
	}

	// Step 4: Log Success
	log.Printf("[SUCCESS] Contribution flow monitoring enabled successfully for Campaign ID: %s", campaignID)
	return nil
}


