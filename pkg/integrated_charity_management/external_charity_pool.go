package integrated_charity_management

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

const (
    CharityProposalDuration    = 14 * 24 * time.Hour // 14 days for proposals
    VotingDuration             = 60 * 24 * time.Hour // 60 days for voting
    MinimumCharityEntries      = 30                  // Minimum number of charities to start voting
    CharityCycleDuration       = 90 * 24 * time.Hour // 90-day charity cycle
    InitialExternalPoolBalance = 0.0                 // Initial external charity pool balance
)


// NewExternalCharityPoolManager initializes a new ExternalCharityPoolManager
func NewExternalCharityPoolManager(ledgerInstance *ledger.Ledger) *ExternalCharityPoolManager {
    return &ExternalCharityPoolManager{
        CharityEntries:      make(map[string]*CharityProposal),
        LedgerInstance:      ledgerInstance,
        ExternalPoolBalance: InitialExternalPoolBalance, // Initialize with 0 balance
    }
}

// StartCharityProposalPeriod initiates the charity proposal process for the external charity pool
func (cpm *ExternalCharityPoolManager) StartCharityProposalPeriod() {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    cpm.ProposalStart = time.Now()
    fmt.Println("External charity proposal period started.")
}

// SubmitCharityProposal allows charities to submit their proposals
func (cpm *ExternalCharityPoolManager) SubmitCharityProposal(name, charityNumber, description, website string, addresses []string) (string, error) {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    if time.Since(cpm.ProposalStart) > CharityProposalDuration && len(cpm.CharityEntries) >= MinimumCharityEntries {
        return "", errors.New("charity proposal period has ended")
    }

    // Ensure no duplicate charity names
    for _, proposal := range cpm.CharityEntries {
        if proposal.Name == name {
            return "", errors.New("charity with the same name already exists")
        }
    }

    // Step 1: Create an encryption instance
    encryptionInstance, err := common.NewEncryption(256) // Assuming NewEncryption creates AES with a 256-bit key
    if err != nil {
        return "", fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Step 2: Encrypt charity addresses using the encryption instance
    encryptedAddresses, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%v", addresses)), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt charity addresses: %v", err)
    }

    charityID := cpm.generateCharityID(name, charityNumber)

    // Step 3: Store the encrypted addresses as a string in CharityProposal
    charity := &CharityProposal{
        CharityID:     charityID,
        Name:          name,
        CharityNumber: charityNumber,
        Description:   description,
        Website:       website,
        Addresses:     []string{base64.StdEncoding.EncodeToString(encryptedAddresses)}, // Store encrypted addresses as base64-encoded string
        CreatedAt:     time.Now(),
        IsValid:       true,
    }

    cpm.CharityEntries[charityID] = charity
    fmt.Printf("Charity proposal submitted: %s\n", name)

    // Log charity proposal submission in the ledger
    err = cpm.logCharityProposalToLedger(charity)
    if err != nil {
        return "", fmt.Errorf("failed to log charity proposal to ledger: %v", err)
    }

    return charityID, nil
}


// VoteForCharity allows users to vote for a charity
func (cpm *ExternalCharityPoolManager) VoteForCharity(charityID string) error {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    charity, exists := cpm.CharityEntries[charityID]
    if !exists {
        return errors.New("charity not found")
    }

    if time.Since(cpm.ProposalStart) > VotingDuration {
        return errors.New("voting period has ended")
    }

    charity.VoteCount++
    fmt.Printf("Vote cast for charity: %s\n", charity.Name)

    // Check if charity is reported as fake (if more than 25 votes cast)
    if charity.VoteCount > 25 && !charity.IsValid {
        cpm.removeFakeCharity(charityID)
    }

    return nil
}

// removeFakeCharity removes a charity flagged as fake
func (cpm *ExternalCharityPoolManager) removeFakeCharity(charityID string) {
    charity := cpm.CharityEntries[charityID]
    delete(cpm.CharityEntries, charityID)
    fmt.Printf("Removed fake charity: %s\n", charity.Name)

    // Log removal of fake charity to ledger
    _ = cpm.logCharityRemovalToLedger(charity)
}

// SelectTopCharities finalizes the voting process and selects the top 20 charities
func (cpm *ExternalCharityPoolManager) SelectTopCharities() ([]*CharityProposal, error) {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    if len(cpm.CharityEntries) < MinimumCharityEntries {
        return nil, errors.New("not enough charity entries to proceed")
    }

    if time.Since(cpm.ProposalStart) < VotingDuration {
        return nil, errors.New("voting period not yet over")
    }

    // Sort charities by vote count
    var charities []*CharityProposal
    for _, charity := range cpm.CharityEntries {
        charities = append(charities, charity)
    }

    // Sort based on vote count, descending
    SortByVotes(charities)

    // Select top 20 charities
    cpm.CurrentCycle = charities[:20]
    cpm.VotingEnd = time.Now()
    fmt.Println("Top 20 charities selected.")

    // Log the top charities selection to ledger
    err := cpm.logTopCharitiesToLedger()
    if err != nil {
        return nil, err
    }

    return cpm.CurrentCycle, nil
}

// SortByVotes sorts the list of CharityProposals in descending order by vote count
func SortByVotes(charities []*CharityProposal) {
    sort.Slice(charities, func(i, j int) bool {
        return charities[i].VoteCount > charities[j].VoteCount
    })
}

// DistributeCharityFunds distributes funds to the top 20 selected charities every 24 hours
func (cpm *ExternalCharityPoolManager) DistributeCharityFunds() error {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    if len(cpm.CurrentCycle) == 0 {
        return errors.New("no charities selected in the current cycle")
    }

    if cpm.ExternalPoolBalance <= 0 {
        return errors.New("insufficient funds in the external charity pool")
    }

    // Convert CharityCycleDuration to hours as a float64 for correct division
    durationInHours := float64(CharityCycleDuration.Hours())
    dailyFunds := cpm.ExternalPoolBalance / float64(len(cpm.CurrentCycle)) / (durationInHours / 24)

    // Iterate over the current cycle and distribute funds individually
    for _, charity := range cpm.CurrentCycle {
        // Distribute funds via the LedgerInstance for each charity
        err := cpm.LedgerInstance.DistributeFunds(dailyFunds) // Assuming DistributeFunds expects a float64
        if err != nil {
            return fmt.Errorf("failed to distribute funds for charity %s: %v", charity.CharityID, err)
        }
        fmt.Printf("Distributed %.2f SYNN to charity: %s\n", dailyFunds, charity.Name)
    }

    // Reduce external pool balance after each distribution
    cpm.ExternalPoolBalance -= dailyFunds * float64(len(cpm.CurrentCycle))

    return nil
}



// UpdateExternalPoolBalance adds transaction fees to the external charity pool balance
func (cpm *ExternalCharityPoolManager) UpdateExternalPoolBalance(amount float64) {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    cpm.ExternalPoolBalance += amount
    fmt.Printf("External charity pool balance updated: %.2f SYNN\n", cpm.ExternalPoolBalance)
}

// GetExternalPoolBalance returns the current external charity pool balance
func (cpm *ExternalCharityPoolManager) GetExternalPoolBalance() float64 {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    return cpm.ExternalPoolBalance
}

// generateCharityID generates a unique ID for a charity based on name and charity number
func (cpm *ExternalCharityPoolManager) generateCharityID(name, charityNumber string) string {
    hashInput := fmt.Sprintf("%s%s%d", name, charityNumber, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// logCharityProposalToLedger logs the charity proposal to the ledger
func (cpm *ExternalCharityPoolManager) logCharityProposalToLedger(charity *CharityProposal) error {
    // Step 1: Create an encryption instance
    encryptionInstance, err := common.NewEncryption(256) // Assuming NewEncryption creates AES with a 256-bit key
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Step 2: Encrypt charity details using the encryption instance
    charityDetails := fmt.Sprintf("%+v", charity)
    encryptedDetails, err := encryptionInstance.EncryptData("AES", []byte(charityDetails), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt charity details: %v", err)
    }

    // Step 3: Convert the encrypted details to a base64 string
    encryptedDetailsStr := base64.StdEncoding.EncodeToString(encryptedDetails)

    // Step 4: Record the encrypted charity details in the ledger
    return cpm.LedgerInstance.RecordCharityProposal(charity.CharityID, encryptedDetailsStr)
}



// logCharityRemovalToLedger logs the removal of a fake charity to the ledger
func (cpm *ExternalCharityPoolManager) logCharityRemovalToLedger(charity *CharityProposal) error {
    // Call RecordCharityRemoval without passing charityID
    return cpm.LedgerInstance.RecordCharityRemoval()
}



// logTopCharitiesToLedger logs the top 20 selected charities to the ledger
func (cpm *ExternalCharityPoolManager) logTopCharitiesToLedger() error {
    // Step 1: Create an encryption instance
    encryptionInstance, err := common.NewEncryption(256) // Assuming AES with a 256-bit key
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Step 2: Format the top charities data and encrypt it
    topCharities := fmt.Sprintf("%+v", cpm.CurrentCycle)
    encryptedTopCharities, err := encryptionInstance.EncryptData("AES", []byte(topCharities), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt top charities: %v", err)
    }

    // Convert the encrypted data to a string and wrap it in a slice of strings
    encryptedDataAsString := string(encryptedTopCharities)
    encryptedDataSlice := []string{encryptedDataAsString} // Wrapping in a slice of strings

    // Step 3: Record the encrypted top charities in the ledger
    err = cpm.LedgerInstance.RecordTopCharities(encryptedDataSlice) // Passing as []string
    if err != nil {
        return fmt.Errorf("failed to record top charities to ledger: %v", err)
    }

    fmt.Println("Top charities logged to the ledger.")
    return nil
}


