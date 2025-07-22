package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/defi"
	"synnergy_network/pkg/ledger"
)

// DeFiAPI provides endpoints for decentralized finance operations
type DeFiAPI struct {
	management *defi.DeFiManagement
}

// NewDeFiAPI creates a new DeFiAPI instance
func NewDeFiAPI(ledgerInstance *ledger.Ledger) *DeFiAPI {
	enc := common.NewEncryption()
	mgmt := defi.NewDeFiManagement(ledgerInstance, enc)
	return &DeFiAPI{management: mgmt}
}

// RegisterRoutes registers defi API endpoints
func (api *DeFiAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/defi/pool/create", api.CreateLiquidityPool).Methods("POST")
	router.HandleFunc("/defi/pool/{id}", api.GetLiquidityPool).Methods("GET")
	router.HandleFunc("/defi/pool/{id}/close", api.CloseLiquidityPool).Methods("POST")
}

// CreateLiquidityPool creates a new liquidity pool
func (api *DeFiAPI) CreateLiquidityPool(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PoolID           string  `json:"pool_id"`
		InitialLiquidity float64 `json:"initial_liquidity"`
		RewardRate       float64 `json:"reward_rate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	pool, err := api.management.CreateLiquidityPool(req.PoolID, req.InitialLiquidity, req.RewardRate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(pool)
}

// GetLiquidityPool returns pool details
func (api *DeFiAPI) GetLiquidityPool(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pool, err := api.management.GetLiquidityPoolDetails(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pool)
}

// CloseLiquidityPool closes a pool
func (api *DeFiAPI) CloseLiquidityPool(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := api.management.CloseLiquidityPool(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "closed"})
}
