package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/defi"
	"synnergy_network/pkg/ledger"
)

// DeFiAPI provides endpoints for decentralized finance operations
type DeFiAPI struct {
	LedgerInstance *ledger.Ledger
}

// NewDeFiAPI creates a new DeFiAPI instance
func NewDeFiAPI(ledgerInstance *ledger.Ledger) *DeFiAPI {
	return &DeFiAPI{LedgerInstance: ledgerInstance}
}

// RegisterRoutes registers DeFi API routes
func (api *DeFiAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/defi/policy", api.CreatePolicy).Methods("POST")
	router.HandleFunc("/defi/policy/{policyId}/claim", api.ClaimPolicy).Methods("POST")
	router.HandleFunc("/defi/policy/{policyId}/status", api.GetPolicyStatus).Methods("GET")
	router.HandleFunc("/defi/policy/{policyId}/payout", api.DistributePayout).Methods("POST")
}

// CreatePolicy creates a new insurance policy
func (api *DeFiAPI) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PolicyID string  `json:"policy_id"`
		Holder   string  `json:"holder"`
		Amount   float64 `json:"amount"`
		Premium  float64 `json:"premium"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	_, err := defi.InsuranceCreatePolicy(req.PolicyID, req.Holder, req.Amount, req.Premium, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("creation failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// ClaimPolicy submits a claim on a policy
func (api *DeFiAPI) ClaimPolicy(w http.ResponseWriter, r *http.Request) {
	policyID := mux.Vars(r)["policyId"]
	var req struct {
		Claimant string  `json:"claimant"`
		Amount   float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err := defi.InsuranceClaimPolicy(policyID, req.Claimant, req.Amount, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("claim failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// GetPolicyStatus returns the status of a policy
func (api *DeFiAPI) GetPolicyStatus(w http.ResponseWriter, r *http.Request) {
	policyID := mux.Vars(r)["policyId"]
	status, err := defi.InsuranceTrackPolicyStatus(policyID, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("status failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"status": status})
}

// DistributePayout distributes payout for a claim
func (api *DeFiAPI) DistributePayout(w http.ResponseWriter, r *http.Request) {
	policyID := mux.Vars(r)["policyId"]
	var req struct {
		Claimant string  `json:"claimant"`
		Amount   float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err := defi.InsuranceDistributePayout(policyID, req.Claimant, req.Amount, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("payout failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
