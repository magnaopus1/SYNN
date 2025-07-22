package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/governance"
	"synnergy_network/pkg/ledger"
)

// GovernanceAPI provides endpoints for governance operations
type GovernanceAPI struct {
	LedgerInstance *ledger.Ledger
}

// NewGovernanceAPI creates a new GovernanceAPI instance
func NewGovernanceAPI(ledgerInstance *ledger.Ledger) *GovernanceAPI {
	return &GovernanceAPI{LedgerInstance: ledgerInstance}
}

// RegisterRoutes registers governance routes
func (api *GovernanceAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/governance/proposal/{id}/validate", api.ValidateProposal).Methods("POST")
	router.HandleFunc("/governance/proposal/{id}/confirm", api.ConfirmProposal).Methods("POST")
	router.HandleFunc("/governance/delegation/threshold", api.SetDelegationThreshold).Methods("POST")
	router.HandleFunc("/governance/compliance/history", api.FetchComplianceHistory).Methods("GET")
}

// ValidateProposal checks proposal compliance
func (api *GovernanceAPI) ValidateProposal(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := governance.GovernanceValidateProposalCompliance(id, api.LedgerInstance); err != nil {
		http.Error(w, fmt.Sprintf("validation failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// ConfirmProposal confirms proposal compliance
func (api *GovernanceAPI) ConfirmProposal(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := governance.GovernanceConfirmProposalCompliance(id, api.LedgerInstance); err != nil {
		http.Error(w, fmt.Sprintf("confirm failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// SetDelegationThreshold sets delegation threshold
func (api *GovernanceAPI) SetDelegationThreshold(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Threshold int `json:"threshold"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := governance.GovernanceSetDelegationThreshold(req.Threshold, api.LedgerInstance); err != nil {
		http.Error(w, fmt.Sprintf("set failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// FetchComplianceHistory returns compliance history
func (api *GovernanceAPI) FetchComplianceHistory(w http.ResponseWriter, r *http.Request) {
	proposalID := r.URL.Query().Get("proposal_id")
	history, err := governance.GovernanceFetchComplianceHistory(proposalID, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("fetch failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(history)
}
