package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/governance"
	"synnergy_network/pkg/ledger"
)

// GovernanceAPI provides endpoints for on-chain governance
type GovernanceAPI struct {
	contract *governance.GovernanceContract
}

// NewGovernanceAPI creates a new GovernanceAPI instance
func NewGovernanceAPI(ledgerInstance *ledger.Ledger) *GovernanceAPI {
	contract := governance.NewGovernanceContract("gov_contract", ledgerInstance)
	return &GovernanceAPI{contract: contract}
}

// RegisterRoutes registers governance API endpoints
func (api *GovernanceAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/governance/vote", api.CastVote).Methods("POST")
	router.HandleFunc("/governance/tally", api.TallyVotes).Methods("POST")
	router.HandleFunc("/governance/execute", api.ExecuteProposal).Methods("POST")
}

// CastVote allows a user to vote on a proposal
func (api *GovernanceAPI) CastVote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProposalID string `json:"proposal_id"`
		VoterID    string `json:"voter_id"`
		Vote       string `json:"vote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := api.contract.CastVote(req.ProposalID, req.VoterID, req.Vote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "vote recorded"})
}

// TallyVotes tallies votes for a proposal
func (api *GovernanceAPI) TallyVotes(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProposalID string `json:"proposal_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := api.contract.TallyVotes(req.ProposalID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "tallied"})
}

// ExecuteProposal executes an approved proposal
func (api *GovernanceAPI) ExecuteProposal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProposalID string                 `json:"proposal_id"`
		Params     map[string]interface{} `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	vm, _ := common.NewVirtualMachine(1, api.contract.LedgerInstance, nil, false)
	enc := &common.Encryption{}
	if err := api.contract.ExecuteProposal(vm, req.ProposalID, req.Params, enc, common.EncryptionKey); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "executed"})
}
