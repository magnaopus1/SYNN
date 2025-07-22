package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/dao"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/tokens/syn900"
)

// DAOAPI exposes DAO management functions over HTTP
type DAOAPI struct {
	management *dao.DAOManagement
}

// NewDAOAPI constructs a DAOAPI
func NewDAOAPI(ledgerInstance *ledger.Ledger) *DAOAPI {
	enc := &common.Encryption{}
	verifier := &syn900.Verifier{}
	mgmt := dao.NewDAOManagement(ledgerInstance, enc, verifier)
	return &DAOAPI{management: mgmt}
}

// RegisterRoutes registers DAO API routes
func (api *DAOAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/dao/create", api.CreateDAO).Methods("POST")
	router.HandleFunc("/dao/{id}/member", api.AddMember).Methods("POST")
	router.HandleFunc("/dao/{id}/proposal", api.SubmitProposal).Methods("POST")
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
	daoObj, err := api.management.CreateDAO(req.Name, req.Creator, req.Members)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(daoObj)
}

// AddMember adds a member to a DAO
func (api *DAOAPI) AddMember(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Wallet string `json:"wallet"`
		Role   string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	daoID := mux.Vars(r)["id"]
	if err := api.management.AddMember(daoID, req.Wallet, req.Role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "member added"})
}

// SubmitProposal submits a proposal to a DAO
func (api *DAOAPI) SubmitProposal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Proposal  string `json:"proposal"`
		Submitter string `json:"submitter"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	daoID := mux.Vars(r)["id"]
	propID, err := api.management.SubmitProposal(daoID, req.Proposal, req.Submitter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"proposal_id": propID})
}
