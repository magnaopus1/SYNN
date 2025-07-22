package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/defi"
	"synnergy_network/pkg/ledger"
)

// DeFiAPI exposes decentralized finance operations
type DeFiAPI struct {
	manager *defi.DeFiManagement
}

// NewDeFiAPI initializes the DeFi API
func NewDeFiAPI(ledgerInstance *ledger.Ledger, enc *common.Encryption) *DeFiAPI {
	return &DeFiAPI{manager: defi.NewDeFiManagement(ledgerInstance, enc)}
}

// RegisterRoutes registers defi endpoints
func (api *DeFiAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/defi/liquidity/create", api.CreateLiquidityPool).Methods("POST")
	router.HandleFunc("/defi/liquidity/{id}", api.GetLiquidityPool).Methods("GET")
	router.HandleFunc("/defi/rewards/{id}/distribute", api.DistributeRewards).Methods("POST")
}

type liquidityReq struct {
	PoolID           string  `json:"pool_id"`
	InitialLiquidity float64 `json:"initial_liquidity"`
	RewardRate       float64 `json:"reward_rate"`
}

// CreateLiquidityPool creates a new liquidity pool
func (api *DeFiAPI) CreateLiquidityPool(w http.ResponseWriter, r *http.Request) {
	var req liquidityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	pool, err := api.manager.CreateLiquidityPool(req.PoolID, req.InitialLiquidity, req.RewardRate)
	if err != nil {
		http.Error(w, fmt.Sprintf("creation failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pool)
}

// GetLiquidityPool returns details about a pool
func (api *DeFiAPI) GetLiquidityPool(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pool, err := api.manager.GetLiquidityPoolDetails(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("not found: %v", err), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pool)
}

// DistributeRewards distributes rewards for a farming record
func (api *DeFiAPI) DistributeRewards(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	reward, err := api.manager.DistributeRewards(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("distribution failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"farming_id": id,
		"reward":     reward,
		"timestamp":  time.Now(),
	})
}
