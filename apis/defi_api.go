package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/defi"
	"synnergy_network/pkg/ledger"
)

// DeFiAPI exposes decentralized finance operations
type DeFiAPI struct {
	manager *defi.DeFiManagement
}

// NewDeFiAPI creates a DeFiAPI with a new management instance
func NewDeFiAPI(l *ledger.Ledger, enc *common.Encryption) *DeFiAPI {
	return &DeFiAPI{manager: defi.NewDeFiManagement(l, enc)}
}

// RegisterRoutes registers DeFi routes
func (api *DeFiAPI) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/defi/liquidity/create", api.CreateLiquidityPool).Methods("POST")
	r.HandleFunc("/defi/liquidity/{id}", api.GetLiquidityPool).Methods("GET")
	r.HandleFunc("/defi/farming/stake", api.ManageYield).Methods("POST")
	r.HandleFunc("/defi/rewards/{id}", api.DistributeRewards).Methods("POST")
	r.HandleFunc("/defi/health", api.Health).Methods("GET")
}

type liquidityRequest struct {
	PoolID string  `json:"pool_id"`
	Amount float64 `json:"amount"`
	Reward float64 `json:"reward"`
}

// CreateLiquidityPool creates a new liquidity pool
func (api *DeFiAPI) CreateLiquidityPool(w http.ResponseWriter, r *http.Request) {
	var req liquidityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	_, err := api.manager.CreateLiquidityPool(req.PoolID, req.Amount, req.Reward)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create pool: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

// GetLiquidityPool retrieves pool details
func (api *DeFiAPI) GetLiquidityPool(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pool, err := api.manager.GetLiquidityPoolDetails(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("pool not found: %v", err), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pool)
}

type farmingRequest struct {
	FarmingID string  `json:"farming_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	PoolID    string  `json:"pool_id"`
}

// ManageYield adds a yield farming record
func (api *DeFiAPI) ManageYield(w http.ResponseWriter, r *http.Request) {
	var req farmingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	_, err := api.manager.ManageYieldFarming(req.FarmingID, req.UserID, req.Amount, req.PoolID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

// DistributeRewards distributes rewards for farming record
func (api *DeFiAPI) DistributeRewards(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := api.manager.DistributeRewards(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("distribution failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func (api *DeFiAPI) Health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}
