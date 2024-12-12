package dao

import (
    "errors"
    "fmt"
    "sync"
    "time"

    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/consensus"
)

// NewDAO initializes a new DAO with given parameters.
func NewDAO(name string, consensusEngine *common.SynnergyConsensus, ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) (*DAO, error) {
    if name == "" {
        return nil, errors.New("DAO name cannot be empty")
    }

    dao := &DAO{
        DAOID:            GenerateUniqueID(),
        Name:             name,
        CreationTime:     time.Now(),
        Members:          make(map[string]string),
        Proposals:        make(map[string]*DAOProposal),
        ProposalHistory:  make(map[string]*DAOProposal),
        ConsensusEngine:  consensusEngine,
        Ledger:           ledgerInstance,
        EncryptionService: encryptionService,
        ProposalDuration: 7 * 24 * time.Hour, // Default proposal voting period: 7 days
    }

    // Add the creator as an admin
    dao.Members[GetCreatorWallet()] = DAOAdmin

    // Record DAO creation in the ledger.
    err := ledgerInstance.DAOLedger.RecordDAOCreation(dao)
    if err != nil {
        return nil, fmt.Errorf("failed to record DAO creation in the ledger: %v", err)
    }

    fmt.Printf("DAO %s created successfully\n", dao.Name)
    return dao, nil
}

// AddMember adds a new member to the DAO.
func (d *DAO) AddMember(walletAddress, role string) error {
    d.mutex.Lock()
    defer d.mutex.Unlock()

    if role != DAOMember && role != DAOAdmin && role != DAOProposalAuthor {
        return errors.New("invalid role")
    }

    // Add the member to the DAO.
    d.Members[walletAddress] = role

    // Record the member addition in the ledger.
    err := d.ledger.DAOLedger.RecordMemberAddition(d.DAOID, walletAddress, role)
    if err != nil {
        return fmt.Errorf("failed to record member addition: %v", err)
    }

    fmt.Printf("Member %s added to DAO %s with role %s\n", walletAddress, d.Name, role)
    return nil
}

// CreateProposal allows a DAO member to create a new proposal.
func (d *DAO) CreateProposal(authorWallet, title, description string) (*DAOProposal, error) {
    d.mutex.Lock()
    defer d.mutex.Unlock()

    if _, exists := d.Members[authorWallet]; !exists {
        return nil, errors.New("author is not a member of the DAO")
    }

    proposal := &DAOProposal{
        ProposalID:   GenerateUniqueID(),
        Title:        title,
        Description:  description,
        Author:       authorWallet,
        CreationTime: time.Now(),
        Status:       "Pending",
    }

    d.Proposals[proposal.ProposalID] = proposal

    // Record the proposal creation in the ledger.
    err := d.Ledger.DAOLedger.RecordProposalCreation(d.DAOID, proposal)
    if err != nil {
        return nil, fmt.Errorf("failed to record proposal creation in the ledger: %v", err)
    }

    fmt.Printf("Proposal %s created by %s in DAO %s\n", title, authorWallet, d.Name)
    return proposal, nil
}

// VoteOnProposal allows members to vote on a DAO proposal.
func (d *DAO) VoteOnProposal(walletAddress, proposalID, vote string) error {
    d.mutex.Lock()
    defer d.mutex.Unlock()

    proposal, exists := d.Proposals[proposalID]
    if !exists {
        return errors.New("proposal not found")
    }

    if _, isMember := d.Members[walletAddress]; !isMember {
        return errors.New("voter is not a member of the DAO")
    }

    if proposal.Status != "Pending" {
        return errors.New("voting period for this proposal has ended")
    }

    switch vote {
    case "approve":
        proposal.ApproveCount++
    case "reject":
        proposal.RejectCount++
    default:
        return errors.New("invalid vote option")
    }

    proposal.VoteCount++

    // Check if voting has reached a conclusion.
    if proposal.ApproveCount >= d.ConsensusEngine.RequiredApprovals() {
        proposal.Status = "Approved"
    } else if proposal.RejectCount >= d.ConsensusEngine.RequiredRejections() {
        proposal.Status = "Rejected"
    }

    // Record the vote in the ledger.
    err := d.Ledger.DAOLedger.RecordProposalVote(d.DAOID, proposalID, walletAddress, vote)
    if err != nil {
        return fmt.Errorf("failed to record vote in ledger: %v", err)
    }

    fmt.Printf("Vote recorded for proposal %s by member %s\n", proposalID, walletAddress)
    return nil
}

// GetProposalStatus checks the status of a proposal.
func (d *DAO) GetProposalStatus(proposalID string) (string, error) {
    d.mutex.Lock()
    defer d.mutex.Unlock()

    proposal, exists := d.Proposals[proposalID]
    if !exists {
        return "", errors.New("proposal not found")
    }

    return proposal.Status, nil
}

// FinalizeProposal marks a proposal as completed after the voting period ends.
func (d *DAO) FinalizeProposal(proposalID string) error {
    d.mutex.Lock()
    defer d.mutex.Unlock()

    proposal, exists := d.Proposals[proposalID]
    if !exists {
        return errors.New("proposal not found")
    }

    if time.Since(proposal.CreationTime) > d.ProposalDuration {
        if proposal.ApproveCount >= d.ConsensusEngine.RequiredApprovals() {
            proposal.Status = "Approved"
        } else {
            proposal.Status = "Rejected"
        }

        // Move the proposal to history.
        d.ProposalHistory[proposalID] = proposal
        delete(d.Proposals, proposalID)

        // Record the final status in the ledger.
        err := d.Ledger.RecordProposalFinalization(d.DAOID, proposalID, proposal.Status)
        if err != nil {
            return fmt.Errorf("failed to record proposal finalization in ledger: %v", err)
        }

        fmt.Printf("Proposal %s has been finalized with status: %s\n", proposalID, proposal.Status)
        return nil
    }

    return errors.New("voting period has not ended yet")
}

// RemoveMember removes a member from the DAO.
func (d *DAO) RemoveMember(walletAddress string) error {
    d.mutex.Lock()
    defer d.mutex.Unlock()

    if _, exists := d.Members[walletAddress]; !exists {
        return errors.New("member not found")
    }

    // Remove the member from the DAO.
    delete(d.Members, walletAddress)

    // Record the member removal in the ledger.
    err := d.Ledger.RecordMemberRemoval(d.DAOID, walletAddress)
    if err != nil {
        return fmt.Errorf("failed to record member removal in the ledger: %v", err)
    }

    fmt.Printf("Member %s removed from DAO %s\n", walletAddress, d.Name)
    return nil
}
