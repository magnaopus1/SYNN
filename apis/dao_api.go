package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/dao"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// DAOAPI provides endpoints for DAO management
type DAOAPI struct {
	LedgerInstance    *ledger.Ledger
	EncryptionService *encryption.Encryption
}

// NewDAOAPI creates a new DAOAPI instance
func NewDAOAPI(ledgerInstance *ledger.Ledger) *DAOAPI {
	return &DAOAPI{LedgerInstance: ledgerInstance, EncryptionService: encryption.NewEncryption()}
}

// RegisterRoutes registers DAO API routes
func (api *DAOAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/dao", api.CreateDAO).Methods("POST")
	router.HandleFunc("/dao/{daoId}/member", api.AddMember).Methods("POST")
	router.HandleFunc("/dao/{daoId}/proposal", api.SubmitProposal).Methods("POST")
	router.HandleFunc("/dao/{daoId}/proposal/{proposalId}/vote", api.VoteProposal).Methods("POST")
}

// CreateDAO handles DAO creation
func (api *DAOAPI) CreateDAO(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string            `json:"name"`
		Creator string            `json:"creator"`
		Members map[string]string `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	_, err := dao.NewDAOManagement(api.LedgerInstance, api.EncryptionService, nil).CreateDAO(req.Name, req.Creator, req.Members)
	if err != nil {
		http.Error(w, fmt.Sprintf("create failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// AddMember adds a member to DAO
func (api *DAOAPI) AddMember(w http.ResponseWriter, r *http.Request) {
	daoID := mux.Vars(r)["daoId"]
	var req struct {
		Wallet string `json:"wallet"`
		Role   string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err := dao.NewDAOManagement(api.LedgerInstance, api.EncryptionService, nil).AddMember(daoID, req.Wallet, req.Role)
	if err != nil {
		http.Error(w, fmt.Sprintf("add failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// SubmitProposal submits a DAO proposal
func (api *DAOAPI) SubmitProposal(w http.ResponseWriter, r *http.Request) {
	daoID := mux.Vars(r)["daoId"]
	var req struct {
		Proposal  string `json:"proposal"`
		Submitter string `json:"submitter"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	_, err := dao.NewDAOManagement(api.LedgerInstance, api.EncryptionService, nil).SubmitProposal(daoID, req.Proposal, req.Submitter)
	if err != nil {
		http.Error(w, fmt.Sprintf("submit failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// VoteProposal casts a vote on a proposal
func (api *DAOAPI) VoteProposal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	daoID := vars["daoId"]
	proposalID := vars["proposalId"]
	var req struct {
		Voter string `json:"voter"`
		Vote  string `json:"vote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err := dao.NewDAOManagement(api.LedgerInstance, api.EncryptionService, nil).VoteOnProposal(daoID, proposalID, req.Voter, req.Vote)
	if err != nil {
		http.Error(w, fmt.Sprintf("vote failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
