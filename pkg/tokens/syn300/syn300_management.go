package syn300

import (
	"errors"
	"sync"
	"time"
)

// syn300Management manages token balances and authorizes tokens for use in governance proposals.
type syn300Management struct {
	Ledger            *ledger.Ledger
	Balances          map[string]uint64           // Token balances of addresses
	Delegations       map[string]string           // Delegator to Delegatee mapping
	VotingPower       map[string]uint64           // Tracks the voting power per address
	AuthorizedTokens  map[string]map[string]uint64 // Maps ProposalID -> (Address -> Authorized Token Amount)
	mutex             sync.RWMutex
}

// NewSyn300Management initializes the syn300Management system for managing token governance.
func NewSyn300Management(ledger *ledger.Ledger) *syn300Management {
	return &syn300Management{
		Ledger:           ledger,
		Balances:         make(map[string]uint64),
		Delegations:      make(map[string]string),
		VotingPower:      make(map[string]uint64),
		AuthorizedTokens: make(map[string]map[string]uint64),
	}
}

// BalanceOf returns the token balance of the given address.
func (m *syn300Management) BalanceOf(address string) (uint64, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	balance, exists := m.Balances[address]
	if !exists {
		return 0, errors.New("address not found")
	}
	return balance, nil
}

// Delegate allows a user to delegate their voting power to another address.
func (m *syn300Management) Delegate(delegator, delegatee string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.Balances[delegator]; !exists {
		return errors.New("delegator address not found")
	}
	m.Delegations[delegator] = delegatee
	return nil
}

// GetDelegation retrieves the delegatee for a given delegator.
func (m *syn300Management) GetDelegation(delegator string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	delegatee, exists := m.Delegations[delegator]
	if !exists {
		return "", errors.New("no delegation found for the address")
	}
	return delegatee, nil
}

// AuthorizeTokensForProposal authorizes tokens to be used in a specific proposal for voting.
func (m *syn300Management) AuthorizeTokensForProposal(proposalID, address string, tokenAmount uint64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	balance, exists := m.Balances[address]
	if !exists {
		return errors.New("address not found")
	}

	if tokenAmount > balance {
		return errors.New("insufficient token balance for authorization")
	}

	if _, exists := m.AuthorizedTokens[proposalID]; !exists {
		m.AuthorizedTokens[proposalID] = make(map[string]uint64)
	}

	m.AuthorizedTokens[proposalID][address] = tokenAmount
	return nil
}

// GetAuthorizedTokens retrieves the amount of tokens authorized for a specific proposal by an address.
func (m *syn300Management) GetAuthorizedTokens(proposalID, address string) (uint64, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	proposalTokens, exists := m.AuthorizedTokens[proposalID]
	if !exists {
		return 0, errors.New("proposal ID not found")
	}

	tokenAmount, exists := proposalTokens[address]
	if !exists {
		return 0, errors.New("no tokens authorized for this address on the proposal")
	}

	return tokenAmount, nil
}

// RevokeTokensForProposal revokes the tokens authorized for voting on a specific proposal.
func (m *syn300Management) RevokeTokensForProposal(proposalID, address string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	proposalTokens, exists := m.AuthorizedTokens[proposalID]
	if !exists {
		return errors.New("proposal ID not found")
	}

	if _, exists := proposalTokens[address]; !exists {
		return errors.New("no tokens authorized for this address on the proposal")
	}

	delete(proposalTokens, address)
	return nil
}

