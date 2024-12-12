package dao

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/tokens/syn900"
)



// NewDAOManagement initializes a new DAO management system.
func NewDAOManagement(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, syn900Verifier *syn900.Verifier) *DAOManagement {
	return &DAOManagement{
		DAOs:              make(map[string]*DAO),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
		Syn900Verifier:    syn900Verifier,
	}
}

// CreateDAO creates a new DAO with initial members and roles.
func (dm *DAOManagement) CreateDAO(name string, creatorWallet string, initialMembers map[string]string) (*DAO, error) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Generate a new DAO ID
	daoID := GenerateUniqueID()

	// Create DAO with creator and initial members
	members := make(map[string]*DAOMember)
	for wallet, role := range initialMembers {
		members[wallet] = &DAOMember{
			WalletAddress: wallet,
			Role:          role,
			VotingPower:   1,
			IsAuthorized:  true,
		}
	}

	// Create the DAO object
	dao := &DAO{
		DAOID:           daoID,
		Name:            name,
		CreatorWallet:   creatorWallet,
		CreatedAt:       time.Now(),
		Members:         members,
		FundsVault:      NewDAOFundVault(daoID, 0, dm.Ledger, dm.EncryptionService, dm.Syn900Verifier), // Fund vault with 0 initial balance
		VotingThreshold: 3, // Default voting threshold
		IsActive:        true,
	}

	// Store the DAO in memory
	dm.DAOs[daoID] = dao

	// Record the DAO creation in the ledger
	err := dm.Ledger.DAOLedger.RecordDAOCreation(dao)
	if err != nil {
		return nil, fmt.Errorf("failed to record DAO creation in ledger: %v", err)
	}

	fmt.Printf("DAO %s created successfully by %s\n", name, creatorWallet)
	return dao, nil
}

// AddMember adds a new member to the DAO with a specified role.
func (dm *DAOManagement) AddMember(daoID string, walletAddress string, role string) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dao, exists := dm.DAOs[daoID]
	if !exists {
		return errors.New("DAO not found")
	}

	// Check if the member already exists
	if _, exists := dao.Members[walletAddress]; exists {
		return errors.New("member already exists in the DAO")
	}

	// Add the new member
	dao.Members[walletAddress] = &DAOMember{
		WalletAddress: walletAddress,
		Role:          role,
		VotingPower:   1,
		IsAuthorized:  true,
	}

	// Record the addition of the member in the ledger
	err := dm.Ledger.DAOLedger.RecordMemberAddition(daoID, walletAddress, role)
	if err != nil {
		return fmt.Errorf("failed to record member addition in ledger: %v", err)
	}

	fmt.Printf("Member %s added to DAO %s with role %s\n", walletAddress, daoID, role)
	return nil
}

// RemoveMember removes a member from the DAO.
func (dm *DAOManagement) RemoveMember(daoID string, walletAddress string) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dao, exists := dm.DAOs[daoID]
	if !exists {
		return errors.New("DAO not found")
	}

	// Check if the member exists
	if _, exists := dao.Members[walletAddress]; !exists {
		return errors.New("member not found in the DAO")
	}

	// Remove the member
	delete(dao.Members, walletAddress)

	// Record the removal of the member in the ledger
	err := dm.Ledger.DAOLedger.RecordMemberRemoval(daoID, walletAddress)
	if err != nil {
		return fmt.Errorf("failed to record member removal in ledger: %v", err)
	}

	fmt.Printf("Member %s removed from DAO %s\n", walletAddress, daoID)
	return nil
}

// ChangeVotingPower adjusts the voting power of a member in the DAO.
func (dm *DAOManagement) ChangeVotingPower(daoID string, walletAddress string, newVotingPower int) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dao, exists := dm.DAOs[daoID]
	if !exists {
		return errors.New("DAO not found")
	}

	// Check if the member exists
	member, exists := dao.Members[walletAddress]
	if !exists {
		return errors.New("member not found in the DAO")
	}

	// Change the voting power
	member.VotingPower = newVotingPower

	// Record the change in the ledger
	err := dm.Ledger.DAOLedger.RecordVotingPowerChange(daoID, walletAddress, newVotingPower)
	if err != nil {
		return fmt.Errorf("failed to record voting power change in ledger: %v", err)
	}

	fmt.Printf("Voting power of member %s in DAO %s changed to %d\n", walletAddress, daoID, newVotingPower)
	return nil
}

// SubmitProposal submits a proposal to the DAO for voting.
func (dm *DAOManagement) SubmitProposal(daoID string, proposal string, submittedBy string) (string, error) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dao, exists := dm.DAOs[daoID]
	if !exists {
		return "", errors.New("DAO not found")
	}

	// Check if the submitter is an authorized member
	member, exists := dao.Members[submittedBy]
	if !exists || !member.IsAuthorized {
		return "", errors.New("unauthorized member")
	}

	// Generate a unique proposal ID
	proposalID := generateUniqueID()

	// Record the proposal submission in the ledger
	err := dm.Ledger.DAOLedger.RecordDAOProposal(daoID, proposalID, proposal, submittedBy)
	if err != nil {
		return "", fmt.Errorf("failed to record proposal in ledger: %v", err)
	}

	fmt.Printf("Proposal %s submitted to DAO %s by %s\n", proposalID, daoID, submittedBy)
	return proposalID, nil
}

// VoteOnProposal allows a DAO member to vote on a submitted proposal.
func (dm *DAOManagement) VoteOnProposal(daoID string, proposalID string, voterWallet string, vote string) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dao, exists := dm.DAOs[daoID]
	if !exists {
		return errors.New("DAO not found")
	}

	// Check if the voter is an authorized member
	member, exists := dao.Members[voterWallet]
	if !exists || !member.IsAuthorized {
		return errors.New("unauthorized voter")
	}

	// Record the vote in the ledger
	err := dm.Ledger.DAOLedger.RecordDAOVote(daoID, proposalID, voterWallet, vote)
	if err != nil {
		return fmt.Errorf("failed to record vote in ledger: %v", err)
	}

	fmt.Printf("Member %s voted on proposal %s in DAO %s\n", voterWallet, proposalID, daoID)
	return nil
}

// DeactivateDAO deactivates the DAO, freezing its functions.
func (dm *DAOManagement) DeactivateDAO(daoID string, requesterWallet string) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dao, exists := dm.DAOs[daoID]
	if !exists {
		return errors.New("DAO not found")
	}

	// Check if the requester is the DAO creator or an admin
	if dao.CreatorWallet != requesterWallet && dao.Members[requesterWallet].Role != "Admin" {
		return errors.New("unauthorized request")
	}

	// Deactivate the DAO
	dao.IsActive = false

	// Record the deactivation in the ledger
	err := dm.Ledger.DAOLedger.RecordDAODeactivation(daoID)
	if err != nil {
		return fmt.Errorf("failed to record DAO deactivation in ledger: %v", err)
	}

	fmt.Printf("DAO %s deactivated by %s\n", daoID, requesterWallet)
	return nil
}

// ViewDAO returns the details of the specified DAO.
func (dm *DAOManagement) ViewDAO(daoID string) (*DAO, error) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dao, exists := dm.DAOs[daoID]
	if !exists {
		return nil, errors.New("DAO not found")
	}

	return dao, nil
}
