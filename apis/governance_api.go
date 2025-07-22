package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/ledger"
)

type GovernanceAPI struct {
	ledger *ledger.Ledger
}

// NewGovernanceAPI creates a GovernanceAPI instance
func NewGovernanceAPI(l *ledger.Ledger) *GovernanceAPI {
	return &GovernanceAPI{ledger: l}
}

func (api *GovernanceAPI) RegisterRoutes(r *mux.Router) {
}

type createProposalRequest struct {
	ProposalID string  `json:"proposal_id"`
	Proposer   string  `json:"proposer"`
	Details    string  `json:"details"`
	Fee        float64 `json:"fee"`
}

func (api *GovernanceAPI) CreateProposal(w http.ResponseWriter, r *http.Request) {
	var req createProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := api.ledger.GovernanceLedger.RecordProposal(
		&api.ledger.AccountsWalletLedger,
		req.ProposalID,
		req.Proposer,
		req.Details,
		req.Fee,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create proposal: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

type voteRequest struct {
	ProposalID string `json:"proposal_id"`
	Voter      string `json:"voter"`
	Vote       string `json:"vote"`
}

// CastVote records a vote on a proposal
func (api *GovernanceAPI) CastVote(w http.ResponseWriter, r *http.Request) {
	var req voteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := api.ledger.GovernanceLedger.RecordVote(&api.ledger.AccountsWalletLedger, req.ProposalID, req.Voter, req.Vote)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to cast vote: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

// TallyVotes finalizes a proposal by recording its execution
func (api *GovernanceAPI) TallyVotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["id"]
	if proposalID == "" {
		http.Error(w, "missing proposal id", http.StatusBadRequest)
		return
	}

	if err := api.ledger.GovernanceLedger.RecordExecution(proposalID); err != nil {
		http.Error(w, fmt.Sprintf("failed to tally votes: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

// GetProposal returns proposal details from the ledger
func (api *GovernanceAPI) GetProposal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["id"]
	record, ok := api.ledger.GovernanceLedger.GovernanceRecords[proposalID]
	if !ok {
		http.Error(w, "proposal not found", http.StatusNotFound)
		return
	}
	proposal := record.Proposals[proposalID]
	json.NewEncoder(w).Encode(proposal)
}

// Health simple health endpoint
func (api *GovernanceAPI) Health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}
