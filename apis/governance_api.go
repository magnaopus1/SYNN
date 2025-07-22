package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/governance"
	"synnergy_network/pkg/ledger"
)

// GovernanceAPI exposes governance module operations via HTTP
type GovernanceAPI struct {
	ledger *ledger.Ledger
}

// NewGovernanceAPI creates a new GovernanceAPI instance
func NewGovernanceAPI(ledgerInstance *ledger.Ledger) *GovernanceAPI {
	return &GovernanceAPI{ledger: ledgerInstance}
}

// RegisterRoutes registers governance endpoints
func (api *GovernanceAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/governance/proposal/revert", api.RevertProposal).Methods("POST")
	router.HandleFunc("/governance/vote/escrow", api.EscrowVoteTokens).Methods("POST")
	router.HandleFunc("/governance/vote/release", api.ReleaseVoteTokens).Methods("POST")
}

type revertProposalRequest struct {
	ProposalID string `json:"proposal_id"`
}

// RevertProposal reverts a governance proposal
func (api *GovernanceAPI) RevertProposal(w http.ResponseWriter, r *http.Request) {
	var req revertProposalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := governance.GovernanceRevertProposal(req.ProposalID, api.ledger); err != nil {
		http.Error(w, fmt.Sprintf("revert failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"proposal_id": req.ProposalID,
		"timestamp":   time.Now(),
	})
}

type escrowVoteRequest struct {
	VoteID string `json:"vote_id"`
	Amount int64  `json:"amount"`
}

// EscrowVoteTokens places vote tokens in escrow
func (api *GovernanceAPI) EscrowVoteTokens(w http.ResponseWriter, r *http.Request) {
	var req escrowVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := governance.GovernanceEscrowVoteTokens(req.VoteID, req.Amount, api.ledger); err != nil {
		http.Error(w, fmt.Sprintf("escrow failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"vote_id":   req.VoteID,
		"amount":    req.Amount,
		"timestamp": time.Now(),
	})
}

type releaseVoteRequest struct {
	VoteID string `json:"vote_id"`
}

// ReleaseVoteTokens releases escrowed vote tokens
func (api *GovernanceAPI) ReleaseVoteTokens(w http.ResponseWriter, r *http.Request) {
	var req releaseVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := governance.GovernanceReleaseVoteTokens(req.VoteID, api.ledger); err != nil {
		http.Error(w, fmt.Sprintf("release failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"vote_id":   req.VoteID,
		"timestamp": time.Now(),
	})
}
