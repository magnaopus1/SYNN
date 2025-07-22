package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/dao"
	"synnergy_network/pkg/ledger"
)

// DAOAPI exposes DAO management functionality
type DAOAPI struct {
	manager *dao.DAOManagement
}

// NewDAOAPI creates a new DAOAPI
func NewDAOAPI(l *ledger.Ledger, enc *common.Encryption) *DAOAPI {
	return &DAOAPI{manager: dao.NewDAOManagement(l, enc, nil)}
}

// RegisterRoutes registers DAO routes
func (api *DAOAPI) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/dao/create", api.CreateDAO).Methods("POST")
	r.HandleFunc("/dao/{id}/member", api.AddMember).Methods("POST")
	r.HandleFunc("/dao/{id}", api.ViewDAO).Methods("GET")
	r.HandleFunc("/dao/health", api.Health).Methods("GET")
}

type createDAORequest struct {
	ID      string            `json:"id"`
	Creator string            `json:"creator"`
	Members map[string]string `json:"members"`
}

// CreateDAO handles DAO creation
func (api *DAOAPI) CreateDAO(w http.ResponseWriter, r *http.Request) {
	var req createDAORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	_, err := api.manager.CreateDAO(req.ID, req.Creator, req.Members)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create dao: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

type memberRequest struct {
	Wallet string `json:"wallet"`
	Role   string `json:"role"`
}

// AddMember adds a DAO member
func (api *DAOAPI) AddMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req memberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := api.manager.AddMember(id, req.Wallet, req.Role); err != nil {
		http.Error(w, fmt.Sprintf("failed to add member: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

// ViewDAO returns DAO info
func (api *DAOAPI) ViewDAO(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	daoObj, err := api.manager.ViewDAO(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("dao not found: %v", err), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(daoObj)
}

func (api *DAOAPI) Health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}
