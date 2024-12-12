package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// GovernanceProposal represents a governance proposal in the network
type GovernanceProposal struct {
	ProposalID      string         // Unique ID of the proposal
	Title           string         // Title of the proposal
	Description     string         // Description of the proposal
	Creator         string         // Address of the proposer
	CreatedAt       time.Time      // Timestamp of creation
	Status          ProposalStatus 
	VotesFor        int            // Number of votes in favor
	VotesAgainst    int            // Number of votes against
	ExpirationTime  time.Time      // Proposal expiration timestamp
	EncryptedDetails string        // Encrypted proposal details
	CreationFee     float64        // Fee charged for proposal creation
}

// ProposalManager manages the lifecycle of governance proposals
type ProposalManager struct {
	Proposals      map[string]*GovernanceProposal // Map of proposal ID to proposals
	mutex          sync.Mutex                     // Mutex for thread-safe operations
	LedgerInstance *ledger.Ledger                 // Ledger instance for tracking proposals
	FeePercentage  float64                        // Fee percentage based on transaction fees (0.25%)
}

// ProposalStatus represents the status of a governance proposal
type ProposalStatus string

const (
    Pending  ProposalStatus = "Pending"
    Approved ProposalStatus = "Approved"
    Rejected ProposalStatus = "Rejected"
)


// NewProposalManager initializes a new ProposalManager with a 0.25% fee
func NewProposalManager(ledgerInstance *ledger.Ledger) *ProposalManager {
    return &ProposalManager{
        Proposals:      make(map[string]*GovernanceProposal),
        LedgerInstance: ledgerInstance,
        FeePercentage:  0.0025, // 0.25% fee
    }
}


// CreateProposal allows a user to submit a new proposal for governance with a dynamic fee based on transaction fees
func (pm *ProposalManager) CreateProposal(creator string, title string, description string, expirationDuration time.Duration, syn900Token *SYN900Token, ledgerInstance *ledger.Ledger) (string, error) {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    proposalID := pm.generateProposalID(creator, title)

    // Check if a proposal with the same ID already exists
    if _, exists := pm.Proposals[proposalID]; exists {
        return "", fmt.Errorf("proposal with ID %s already exists", proposalID)
    }

    // Verify Syn-900 token and stamp the creator's ID onto the proposal
    if !pm.verifySyn900Token(syn900Token, creator) {
        return "", errors.New("invalid Syn-900 token")
    }

    // Calculate the proposal creation fee based on the last 500 blocks' transaction fees
    creationFee, err := pm.calculateCreationFee()
    if err != nil {
        return "", fmt.Errorf("failed to calculate proposal creation fee: %v", err)
    }

    // Charge the creator's account the creation fee (include ledger instance)
    err = pm.chargeFee(creator, creationFee, ledgerInstance)
    if err != nil {
        return "", fmt.Errorf("failed to charge creation fee: %v", err)
    }

    // Create an instance of Encryption
    encryptionInstance := &Encryption{}

    // Encrypt proposal details using AES algorithm
    encryptedDetails, err := encryptionInstance.EncryptData("AES", []byte(description), EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt proposal details: %v", err)
    }

    // Create and store the new proposal
    newProposal := &GovernanceProposal{
        ProposalID:      proposalID,
        Title:           title,
        Description:     description,
        Creator:         creator,
        CreatedAt:       time.Now(),
        Status:          Pending,  // Correctly use the ProposalStatus type
        VotesFor:        0,
        VotesAgainst:    0,
        ExpirationTime:  time.Now().Add(expirationDuration),
        EncryptedDetails: string(encryptedDetails), // Store the encrypted details
        CreationFee:     creationFee,
    }

    pm.Proposals[proposalID] = newProposal

    // Log the proposal creation in the ledger
    err = pm.logProposalToLedger(newProposal)
    if err != nil {
        return "", fmt.Errorf("failed to log proposal to ledger: %v", err)
    }

    // Destroy the Syn-900 token after successful proposal creation (pass by reference)
    err = pm.destroySyn900Token(syn900Token)
    if err != nil {
        return "", fmt.Errorf("failed to destroy Syn-900 token: %v", err)
    }

    fmt.Printf("Proposal %s created by %s with a fee of %.2f SYNN.\n", proposalID, creator, creationFee)
    return proposalID, nil
}




// calculateCreationFee calculates the average transaction fee for the last 500 blocks and applies a 0.25% fee
func (pm *ProposalManager) calculateCreationFee() (float64, error) {
    totalFees, err := pm.LedgerInstance.BlockchainConsensusCoinLedger.GetTotalTransactionFeesForLastBlocks(500)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve transaction fees: %v", err)
    }

    averageFee := totalFees / 500.0
    creationFee := averageFee * pm.FeePercentage
    return creationFee, nil
}

// chargeFee deducts the proposal creation fee from the creator's account and distributes it
func (pm *ProposalManager) chargeFee(creator string, feeAmount float64, ledgerInstance *ledger.Ledger) error {
    // Add the transaction to the ledger (deduct fee)
    err := ledgerInstance.BlockchainConsensusCoinLedger.AddTransaction(creator, "", feeAmount)
    if err != nil {
        return fmt.Errorf("failed to record transaction: %v", err)
    }

    // Distribute the fee across the different pools
    tdm := NewTransactionDistributionManager(ledgerInstance)
    err = tdm.DistributeRewards("blockID_placeholder", feeAmount) // Use actual blockID if available
    if err != nil {
        return fmt.Errorf("failed to distribute rewards: %v", err)
    }

    fmt.Printf("Charged %s a fee of %.2f SYNN for proposal creation and distributed the fee.\n", creator, feeAmount)
    return nil
}





// verifySyn900Token ensures the Syn-900 ID token is valid for the creator
func (pm *ProposalManager) verifySyn900Token(token *SYN900Token, creator string) bool {
    // Verify if the token's owner matches the creator and if the status is 'active'
    return token.Owner == creator && token.Status == "active"
}



// destroySyn900Token destroys the Syn-900 ID token after proposal creation
func (pm *ProposalManager) destroySyn900Token(token *SYN900Token) error {
    // Lock the token's mutex to ensure thread safety
    token.mutex.Lock()
    defer token.mutex.Unlock() // Ensure the mutex is unlocked after the operation

    // Set the token's status to 'revoked'
    token.Status = "revoked"
    fmt.Printf("Syn-900 ID token for creator %s destroyed after proposal creation.\n", token.Owner)
    return nil
}




// logProposalToLedger logs the proposal details to the ledger
func (pm *ProposalManager) logProposalToLedger(proposal *GovernanceProposal) error {
    // Encrypt the proposal details
    encryptionInstance := &Encryption{}
    encryptedDetails, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%+v", proposal)), EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt proposal details: %v", err)
    }

    // Generate a hash for the proposal
    hash := pm.generateProposalHash(proposal)

    // Log the proposal to the ledger, using CreationFee as the float64 argument
    err = pm.LedgerInstance.GovernanceLedger.RecordProposal(hash, string(encryptedDetails), "pending", proposal.CreationFee)
    if err != nil {
        return fmt.Errorf("failed to log proposal to ledger: %v", err)
    }

    fmt.Printf("Proposal %s logged to the ledger.\n", proposal.ProposalID)
    return nil
}




// generateProposalID creates a unique ID for a proposal based on the creator and title
func (pm *ProposalManager) generateProposalID(creator, title string) string {
    hashInput := fmt.Sprintf("%s%s%d", creator, title, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// generateProposalHash creates a hash for a proposal
func (pm *ProposalManager) generateProposalHash(proposal *GovernanceProposal) string {
    hashInput := fmt.Sprintf("%s%d", proposal.ProposalID, proposal.CreatedAt.UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}
