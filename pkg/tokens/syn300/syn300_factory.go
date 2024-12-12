package syn300

import (
	"errors"
	"sync"
	"time"
)

// Token represents a governance token in the Synthron ecosystem, facilitating complex governance functionalities.
type SYN300Token struct {
	ID           string
	BalanceOf    map[string]uint64
	VotingPower  map[string]uint64
	Delegations  map[string]string  // Maps delegator's address to a delegate's address.
	Votes        map[string]map[string]uint64 // Maps proposal IDs to voter addresses and their vote weights.
	Metadata     SYN300TokenMetadata
	mutex        sync.RWMutex
	Ledger       *ledger.Ledger
}

// TokenMetadata contains important details about the governance token.
type SYN300TokenMetadata struct {
	Name        string
	Symbol      string
	Decimals    int
	TotalSupply uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewToken creates a new governance token with the given name, symbol, and total supply.
func NewToken(name, symbol string, initialSupply uint64, ledger *ledger.Ledger) (*Token, error) {
	token := &Token{
		ID:          generateUniqueID(),
		BalanceOf:   make(map[string]uint64),
		VotingPower: make(map[string]uint64),
		Delegations: make(map[string]string),
		Votes:       make(map[string]map[string]uint64),
		Metadata: TokenMetadata{
			Name:        name,
			Symbol:      symbol,
			Decimals:    18,
			TotalSupply: initialSupply,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Ledger: ledger,
	}

	// Record token creation in the ledger
	err := ledger.RecordTokenCreation(token.Metadata.Name, token.Metadata.Symbol, token.Metadata.TotalSupply)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Transfer allows the transfer of tokens between addresses.
func (t *Token) Transfer(from, to string, amount uint64) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if amount == 0 || from == to {
		return errors.New("invalid transfer")
	}

	fromBalance := t.BalanceOf[from]
	if fromBalance < amount {
		return errors.New("insufficient balance")
	}

	t.BalanceOf[from] -= amount
	t.BalanceOf[to] += amount

	// Record the transfer in the ledger
	err := t.Ledger.RecordTransferEvent(t.ID, from, to, amount, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// Delegate allows a user to delegate their voting power to another address.
func (t *Token) Delegate(delegator, delegatee string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if delegator == delegatee {
		return errors.New("cannot delegate to self")
	}

	// Ensure delegator has balance
	if t.BalanceOf[delegator] == 0 {
		return errors.New("no balance to delegate")
	}

	// Record the delegation
	t.Delegations[delegator] = delegatee
	t.VotingPower[delegatee] += t.BalanceOf[delegator]

	// Record the delegation in the ledger
	err := t.Ledger.RecordDelegationEvent(t.ID, delegator, delegatee, t.BalanceOf[delegator], time.Now())
	if err != nil {
		return err
	}

	return nil
}

// CastVote allows a user to cast their vote on a proposal.
func (t *Token) CastVote(voter, proposalID string, voteWeight uint64) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exists := t.Votes[proposalID]; !exists {
		t.Votes[proposalID] = make(map[string]uint64)
	}

	t.Votes[proposalID][voter] = voteWeight

	// Record the vote in the ledger
	err := t.Ledger.RecordVoteEvent(t.ID, voter, proposalID, voteWeight, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// generateUniqueID creates a unique identifier for each token.
func generateUniqueID() string {
	// Implementation for generating a unique token ID
	// This can use encryption or a UUID generator
	return encryption.GenerateUUID()
}
