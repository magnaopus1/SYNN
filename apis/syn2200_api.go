package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"
	"github.com/gorilla/mux"
)

type SYN2200API struct{}

func NewSYN2200API() *SYN2200API { return &SYN2200API{} }

func (api *SYN2200API) RegisterRoutes(router *mux.Router) {
	// Core prediction token management (15 endpoints)
	router.HandleFunc("/syn2200/tokens", api.CreatePredictionToken).Methods("POST")
	router.HandleFunc("/syn2200/tokens/{tokenID}", api.GetToken).Methods("GET")
	router.HandleFunc("/syn2200/tokens", api.ListTokens).Methods("GET")
	router.HandleFunc("/syn2200/tokens/{tokenID}/transfer", api.TransferTokens).Methods("POST")
	router.HandleFunc("/syn2200/tokens/{tokenID}/stake", api.StakeTokens).Methods("POST")
	router.HandleFunc("/syn2200/tokens/balance/{address}", api.GetBalance).Methods("GET")
	router.HandleFunc("/syn2200/tokens/{tokenID}/metadata", api.UpdateMetadata).Methods("PUT")
	router.HandleFunc("/syn2200/tokens/{tokenID}/burn", api.BurnTokens).Methods("POST")
	router.HandleFunc("/syn2200/tokens/{tokenID}/mint", api.MintTokens).Methods("POST")
	router.HandleFunc("/syn2200/tokens/{tokenID}/freeze", api.FreezeToken).Methods("POST")
	router.HandleFunc("/syn2200/tokens/{tokenID}/history", api.GetTokenHistory).Methods("GET")
	router.HandleFunc("/syn2200/tokens/{tokenID}/oracle", api.GetOracleData).Methods("GET")
	router.HandleFunc("/syn2200/tokens/{tokenID}/fees", api.GetFeeStructure).Methods("GET")
	router.HandleFunc("/syn2200/tokens/{tokenID}/reputation", api.GetReputationScore).Methods("GET")
	router.HandleFunc("/syn2200/tokens/{tokenID}/governance", api.GetGovernanceRights).Methods("GET")

	// Market creation and management (20 endpoints)
	router.HandleFunc("/syn2200/markets", api.CreateMarket).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}", api.GetMarket).Methods("GET")
	router.HandleFunc("/syn2200/markets", api.ListMarkets).Methods("GET")
	router.HandleFunc("/syn2200/markets/{marketID}/bet", api.PlaceBet).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/odds", api.GetOdds).Methods("GET")
	router.HandleFunc("/syn2200/markets/{marketID}/volume", api.GetMarketVolume).Methods("GET")
	router.HandleFunc("/syn2200/markets/{marketID}/resolve", api.ResolveMarket).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/dispute", api.DisputeResolution).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/liquidity", api.AddLiquidity).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/withdraw-liquidity", api.WithdrawLiquidity).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/pause", api.PauseMarket).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/resume", api.ResumeMarket).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/cancel", api.CancelMarket).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/extend", api.ExtendMarket).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/split", api.SplitMarket).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/merge", api.MergeMarkets).Methods("POST")
	router.HandleFunc("/syn2200/markets/{marketID}/arbitrage", api.GetArbitrageOpportunities).Methods("GET")
	router.HandleFunc("/syn2200/markets/{marketID}/participants", api.GetMarketParticipants).Methods("GET")
	router.HandleFunc("/syn2200/markets/{marketID}/sentiment", api.GetMarketSentiment).Methods("GET")
	router.HandleFunc("/syn2200/markets/trending", api.GetTrendingMarkets).Methods("GET")

	// Betting and positions (15 endpoints)
	router.HandleFunc("/syn2200/bets", api.GetUserBets).Methods("GET")
	router.HandleFunc("/syn2200/bets/{betID}", api.GetBetDetails).Methods("GET")
	router.HandleFunc("/syn2200/bets/{betID}/cancel", api.CancelBet).Methods("POST")
	router.HandleFunc("/syn2200/bets/{betID}/hedge", api.HedgeBet).Methods("POST")
	router.HandleFunc("/syn2200/bets/portfolio", api.GetBettingPortfolio).Methods("GET")
	router.HandleFunc("/syn2200/bets/pnl", api.GetProfitAndLoss).Methods("GET")
	router.HandleFunc("/syn2200/positions", api.GetPositions).Methods("GET")
	router.HandleFunc("/syn2200/positions/{positionID}", api.GetPositionDetails).Methods("GET")
	router.HandleFunc("/syn2200/positions/{positionID}/close", api.ClosePosition).Methods("POST")
	router.HandleFunc("/syn2200/positions/{positionID}/modify", api.ModifyPosition).Methods("PUT")
	router.HandleFunc("/syn2200/positions/risk", api.GetPositionRisk).Methods("GET")
	router.HandleFunc("/syn2200/positions/exposure", api.GetMarketExposure).Methods("GET")
	router.HandleFunc("/syn2200/payouts", api.ClaimPayouts).Methods("POST")
	router.HandleFunc("/syn2200/payouts/history", api.GetPayoutHistory).Methods("GET")
	router.HandleFunc("/syn2200/strategy", api.GetBettingStrategy).Methods("GET")

	// Oracle and data feeds (10 endpoints)
	router.HandleFunc("/syn2200/oracles", api.RegisterOracle).Methods("POST")
	router.HandleFunc("/syn2200/oracles/{oracleID}", api.GetOracleInfo).Methods("GET")
	router.HandleFunc("/syn2200/oracles/{oracleID}/data", api.SubmitOracleData).Methods("POST")
	router.HandleFunc("/syn2200/oracles/{oracleID}/reputation", api.GetOracleReputation).Methods("GET")
	router.HandleFunc("/syn2200/oracles/consensus", api.GetOracleConsensus).Methods("GET")
	router.HandleFunc("/syn2200/feeds", api.GetDataFeeds).Methods("GET")
	router.HandleFunc("/syn2200/feeds/{feedID}", api.GetFeedData).Methods("GET")
	router.HandleFunc("/syn2200/feeds/{feedID}/subscribe", api.SubscribeToFeed).Methods("POST")
	router.HandleFunc("/syn2200/feeds/prices", api.GetPriceFeeds).Methods("GET")
	router.HandleFunc("/syn2200/feeds/events", api.GetEventFeeds).Methods("GET")

	// Analytics and reporting (10 endpoints)
	router.HandleFunc("/syn2200/analytics/volume", api.GetTradingVolume).Methods("GET")
	router.HandleFunc("/syn2200/analytics/accuracy", api.GetPredictionAccuracy).Methods("GET")
	router.HandleFunc("/syn2200/analytics/performance", api.GetPerformanceMetrics).Methods("GET")
	router.HandleFunc("/syn2200/analytics/trends", api.GetMarketTrends).Methods("GET")
	router.HandleFunc("/syn2200/analytics/sentiment", api.GetSentimentAnalysis).Methods("GET")
	router.HandleFunc("/syn2200/reports/market", api.GenerateMarketReport).Methods("GET")
	router.HandleFunc("/syn2200/reports/user", api.GenerateUserReport).Methods("GET")
	router.HandleFunc("/syn2200/reports/oracle", api.GenerateOracleReport).Methods("GET")
	router.HandleFunc("/syn2200/analytics/arbitrage", api.GetArbitrageAnalytics).Methods("GET")
	router.HandleFunc("/syn2200/analytics/correlation", api.GetMarketCorrelations).Methods("GET")

	// Administrative (5 endpoints)
	router.HandleFunc("/syn2200/admin/settings", api.UpdateSettings).Methods("PUT")
	router.HandleFunc("/syn2200/admin/emergency", api.EmergencyPause).Methods("POST")
	router.HandleFunc("/syn2200/admin/health", api.GetSystemHealth).Methods("GET")
	router.HandleFunc("/syn2200/admin/logs", api.GetSystemLogs).Methods("GET")
	router.HandleFunc("/syn2200/admin/backup", api.CreateBackup).Methods("POST")
}

func (api *SYN2200API) CreatePredictionToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("Creating prediction market token - Request from %s", r.RemoteAddr)
	
	var request struct {
		Name         string  `json:"name" validate:"required"`
		Symbol       string  `json:"symbol" validate:"required"`
		TotalSupply  float64 `json:"total_supply" validate:"required,min=1"`
		MarketType   string  `json:"market_type" validate:"required,oneof=binary categorical scalar"`
		Category     string  `json:"category" validate:"required"`
		ExpiryDate   string  `json:"expiry_date" validate:"required"`
		OracleSource string  `json:"oracle_source" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error":"Invalid JSON format","code":"INVALID_JSON"}`, http.StatusBadRequest)
		return
	}

	tokenID := fmt.Sprintf("PRED_%d", time.Now().UnixNano())
	
	response := map[string]interface{}{
		"success":       true,
		"token_id":      tokenID,
		"name":          request.Name,
		"symbol":        request.Symbol,
		"total_supply":  request.TotalSupply,
		"market_type":   request.MarketType,
		"category":      request.Category,
		"expiry_date":   request.ExpiryDate,
		"oracle_source": request.OracleSource,
		"created_at":    time.Now().Format(time.RFC3339),
		"status":        "active",
		"network":       "synnergy",
		"decimals":      18,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	log.Printf("Prediction market token created: %s", tokenID)
}

func (api *SYN2200API) GetToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]
	response := map[string]interface{}{"success": true, "token_id": tokenID, "market_type": "binary", "status": "active"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *SYN2200API) GetTradingVolume(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"success": true, "volume": map[string]interface{}{"total_24h": 2500000.0, "markets": 150, "bets": 5000}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}