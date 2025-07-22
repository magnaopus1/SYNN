package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/dao"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/tokens/syn900"
)

// DAOAPI exposes DAO management operations
type DAOAPI struct {
	manager *dao.DAOManagement
}

// NewDAOAPI initializes DAO management
func NewDAOAPI(ledgerInstance *ledger.Ledger, enc *encryption.Encryption, verifier *syn900.Verifier) *DAOAPI {
	return &DAOAPI{manager: dao.NewDAOManagement(ledgerInstance, enc, verifier)}
}

// RegisterRoutes registers DAO endpoints
func (api *DAOAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/dao/create", api.CreateDAO).Methods("POST")
	router.HandleFunc("/dao/{id}/proposal", api.SubmitProposal).Methods("POST")
	router.HandleFunc("/dao/{id}/vote", api.VoteOnProposal).Methods("POST")
}

type createDAORequest struct {
	Name          string            `json:"name"`
	CreatorWallet string            `json:"creator_wallet"`
	Members       map[string]string `json:"members"`
}

// CreateDAO handles DAO creation
func (api *DAOAPI) CreateDAO(w http.ResponseWriter, r *http.Request) {
	var req createDAORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	daoObj, err := api.manager.CreateDAO(req.Name, req.CreatorWallet, req.Members)
	if err != nil {
		http.Error(w, fmt.Sprintf("creation failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(daoObj)
}

type proposalReq struct {
	Proposal    string `json:"proposal"`
	SubmittedBy string `json:"submitted_by"`
}

// SubmitProposal submits a proposal for voting
func (api *DAOAPI) SubmitProposal(w http.ResponseWriter, r *http.Request) {
	var req proposalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	daoID := mux.Vars(r)["id"]
	proposalID, err := api.manager.SubmitProposal(daoID, req.Proposal, req.SubmittedBy)
	if err != nil {
		http.Error(w, fmt.Sprintf("submission failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"dao_id":      daoID,
		"proposal_id": proposalID,
		"timestamp":   time.Now(),
	})
}

type voteReq struct {
	ProposalID string `json:"proposal_id"`
	Voter      string `json:"voter"`
	Vote       string `json:"vote"`
}

// VoteOnProposal records a vote on a DAO proposal
func (api *DAOAPI) VoteOnProposal(w http.ResponseWriter, r *http.Request) {
	var req voteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	daoID := mux.Vars(r)["id"]
	if err := api.manager.VoteOnProposal(daoID, req.ProposalID, req.Voter, req.Vote); err != nil {
		http.Error(w, fmt.Sprintf("vote failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"dao_id":      daoID,
		"proposal_id": req.ProposalID,
		"timestamp":   time.Now(),
	})
}
