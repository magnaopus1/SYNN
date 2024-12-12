package syn2500

import (
	"time"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rsa"
	"crypto/rand"
)

// DAOStorage is the main struct responsible for storing and managing DAO token-related data.
type DAOStorage struct {
	TokenStore map[string]common.DAOToken  // Stores DAO tokens by Token ID
	ProposalStore map[string]common.DAOProposal // Stores DAO proposals by Proposal ID
}

// NewDAOStorage initializes the DAO storage system.
func NewDAOStorage() *DAOStorage {
	return &DAOStorage{
		TokenStore:    make(map[string]common.DAOToken),
		ProposalStore: make(map[string]common.DAOProposal),
	}
}

// StoreDAOToken securely stores a DAO token in the ledger and storage.
func (ds *DAOStorage) StoreDAOToken(token common.DAOToken) error {
	// Generate token hash for unique identification and verification
	tokenHash := ds.generateHash(token)
	token.TokenHash = tokenHash

	// Validate and store the token using Synnergy Consensus
	err := synconsensus.ValidateSubBlock(token.TokenHash)
	if err != nil {
		return errors.New("token validation failed through Synnergy Consensus")
	}

	// Store token in ledger for immutability
	err = ledger.StoreDAOToken(token)
	if err != nil {
		return errors.New("failed to store DAO token in ledger")
	}

	// Store token in local storage
	ds.TokenStore[token.TokenID] = token

	return nil
}

// RetrieveDAOToken retrieves a DAO token from storage by its Token ID.
func (ds *DAOStorage) RetrieveDAOToken(tokenID string) (common.DAOToken, error) {
	token, exists := ds.TokenStore[tokenID]
	if !exists {
		return common.DAOToken{}, errors.New("token not found in storage")
	}
	return token, nil
}

// StoreDAOProposal stores a DAO proposal securely and manages its metadata.
func (ds *DAOStorage) StoreDAOProposal(proposal common.DAOProposal) error {
	// Generate hash for the proposal to ensure its integrity
	proposalHash := ds.generateHash(proposal)
	proposal.ProposalHash = proposalHash

	// Validate and store the proposal using Synnergy Consensus
	err := synconsensus.ValidateSubBlock(proposal.ProposalHash)
	if err != nil {
		return errors.New("proposal validation failed through Synnergy Consensus")
	}

	// Store proposal in the ledger for immutability
	err = ledger.StoreDAOProposal(proposal)
	if err != nil {
		return errors.New("failed to store DAO proposal in ledger")
	}

	// Store proposal in local storage
	ds.ProposalStore[proposal.ProposalID] = proposal

	return nil
}

// RetrieveDAOProposal retrieves a DAO proposal by its Proposal ID.
func (ds *DAOStorage) RetrieveDAOProposal(proposalID string) (common.DAOProposal, error) {
	proposal, exists := ds.ProposalStore[proposalID]
	if !exists {
		return common.DAOProposal{}, errors.New("proposal not found in storage")
	}
	return proposal, nil
}

// generateHash generates a unique hash for a DAO token or proposal using SHA-256.
func (ds *DAOStorage) generateHash(data interface{}) string {
	hashInput := ""
	switch v := data.(type) {
	case common.DAOToken:
		hashInput = v.TokenID + v.Owner + v.Timestamp.String()
	case common.DAOProposal:
		hashInput = v.ProposalID + v.Creator + v.Timestamp.String()
	}

	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// UpdateDAOTokenStatus updates the status of a DAO token (e.g., Active, Inactive).
func (ds *DAOStorage) UpdateDAOTokenStatus(tokenID string, status string) error {
	token, exists := ds.TokenStore[tokenID]
	if !exists {
		return errors.New("token not found in storage")
	}

	token.Status = status
	token.Timestamp = time.Now()

	// Validate the updated token using Synnergy Consensus
	err := synconsensus.ValidateSubBlock(token.TokenHash)
	if err != nil {
		return errors.New("token validation failed for update through Synnergy Consensus")
	}

	// Store updated token in ledger
	err = ledger.StoreDAOToken(token)
	if err != nil {
		return errors.New("failed to update DAO token in ledger")
	}

	// Update the token in local storage
	ds.TokenStore[token.TokenID] = token

	return nil
}

// DeleteDAOToken securely deletes a DAO token from the storage and ledger.
func (ds *DAOStorage) DeleteDAOToken(tokenID string) error {
	_, exists := ds.TokenStore[tokenID]
	if !exists {
		return errors.New("token not found in storage")
	}

	// Remove the token from ledger
	err := ledger.DeleteDAOToken(tokenID)
	if err != nil {
		return errors.New("failed to delete DAO token from ledger")
	}

	// Remove the token from local storage
	delete(ds.TokenStore, tokenID)

	return nil
}

// DeleteDAOProposal deletes a DAO proposal from the storage and ledger.
func (ds *DAOStorage) DeleteDAOProposal(proposalID string) error {
	_, exists := ds.ProposalStore[proposalID]
	if !exists {
		return errors.New("proposal not found in storage")
	}

	// Remove the proposal from ledger
	err := ledger.DeleteDAOProposal(proposalID)
	if err != nil {
		return errors.New("failed to delete DAO proposal from ledger")
	}

	// Remove the proposal from local storage
	delete(ds.ProposalStore, proposalID)

	return nil
}

// ListAllDAOTokens returns a list of all stored DAO tokens.
func (ds *DAOStorage) ListAllDAOTokens() []common.DAOToken {
	tokens := []common.DAOToken{}
	for _, token := range ds.TokenStore {
		tokens = append(tokens, token)
	}
	return tokens
}

// ListAllDAOProposals returns a list of all stored DAO proposals.
func (ds *DAOStorage) ListAllDAOProposals() []common.DAOProposal {
	proposals := []common.DAOProposal{}
	for _, proposal := range ds.ProposalStore {
		proposals = append(proposals, proposal)
	}
	return proposals
}
