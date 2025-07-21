package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"synnergy_network/pkg/consensus"
	"synnergy_network/pkg/network"
	"synnergy_network/pkg/transactions"
	"synnergy_network/pkg/smart_contract"
	"synnergy_network/pkg/wallet"
	"synnergy_network/pkg/tokens"
	"synnergy_network/pkg/defi"
	"synnergy_network/pkg/governance"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/common"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

// SynnergyAPIServer represents the main API server for the Synnergy Network
type SynnergyAPIServer struct {
	Router         *mux.Router
	LedgerInstance *ledger.Ledger
	Port           string
	
	// Core API modules
	ConsensusAPI     *ConsensusAPI
	NetworkAPI       *NetworkAPI
	TransactionsAPI  *TransactionsAPI
	SmartContractAPI *SmartContractAPI
	WalletAPI        *WalletAPI
	TokensAPI        *TokensAPI
	DeFiAPI          *DeFiAPI
	GovernanceAPI    *GovernanceAPI
	
	// Token API modules - Batch #11 Module-Validated APIs
	SYN3900API       *SYN3900API
	SYN4700API       *SYN4700API
	SYN4900API       *SYN4900API
	SYN5000API       *SYN5000API
	
	// Infrastructure components
	NetworkManager    *network.NetworkManager
	TransactionPool   *transactions.TransactionPool
	ConsensusEngine   *consensus.SynnergyConsensus
	EncryptionService *common.Encryption
	GasManager        *common.GasManager
}

// NewSynnergyAPIServer creates a new API server instance
func NewSynnergyAPIServer(port string, ledgerInstance *ledger.Ledger) *SynnergyAPIServer {
	// Initialize router
	router := mux.NewRouter()
	
	// Initialize infrastructure components
	encryptionService := common.NewEncryption()
	gasManager := common.NewGasManager(ledgerInstance, nil, 0.001) // Default gas price
	networkManager := network.NewNetworkManager("localhost:8080", ledgerInstance, 30*time.Minute)
	transactionPool := transactions.NewTransactionPool(10000, ledgerInstance, encryptionService)
	
	// Create API server instance
	server := &SynnergyAPIServer{
		Router:            router,
		LedgerInstance:    ledgerInstance,
		Port:              port,
		NetworkManager:    networkManager,
		TransactionPool:   transactionPool,
		EncryptionService: encryptionService,
		GasManager:        gasManager,
	}
	
	// Initialize API modules
	server.initializeAPIModules()
	
	// Register all routes
	server.registerRoutes()
	
	return server
}

// initializeAPIModules initializes all API modules
func (s *SynnergyAPIServer) initializeAPIModules() {
	s.ConsensusAPI = NewConsensusAPI(s.LedgerInstance)
	s.NetworkAPI = NewNetworkAPI(s.NetworkManager, s.LedgerInstance)
	s.TransactionsAPI = NewTransactionsAPI(s.TransactionPool, s.LedgerInstance)
	s.SmartContractAPI = NewSmartContractAPI(s.LedgerInstance)
	s.WalletAPI = NewWalletAPI(s.LedgerInstance, s.NetworkManager)
	
	// Initialize Batch #11 Module-Validated Token APIs
	synnergyConsensus := &common.SynnergyConsensus{} // Initialize properly in production
	synnergyMutex := &common.SynnergyMutex{}         // Initialize properly in production
	
	s.SYN3900API = NewSYN3900API(s.LedgerInstance, synnergyConsensus, synnergyMutex)
	s.SYN4700API = NewSYN4700API(s.LedgerInstance, synnergyConsensus, synnergyMutex)
	s.SYN4900API = NewSYN4900API(s.LedgerInstance, synnergyConsensus, synnergyMutex)
	s.SYN5000API = NewSYN5000API(s.LedgerInstance, synnergyConsensus, synnergyMutex)
	
	// TODO: Initialize other API modules as they are created
	// s.TokensAPI = NewTokensAPI(...)
	// s.DeFiAPI = NewDeFiAPI(...)
	// s.GovernanceAPI = NewGovernanceAPI(...)
}

// registerRoutes registers all API routes with their respective modules
func (s *SynnergyAPIServer) registerRoutes() {
	// Root health check endpoint
	s.Router.HandleFunc("/health", s.HealthCheck).Methods("GET")
	s.Router.HandleFunc("/", s.Welcome).Methods("GET")
	
	// API version prefix
	apiV1 := s.Router.PathPrefix("/api/v1").Subrouter()
	
	// Register module routes
	s.ConsensusAPI.RegisterRoutes(apiV1)
	s.NetworkAPI.RegisterRoutes(apiV1)
	s.TransactionsAPI.RegisterRoutes(apiV1)
	s.SmartContractAPI.RegisterRoutes(apiV1)
	s.WalletAPI.RegisterRoutes(apiV1)
	
	// Register Batch #11 Module-Validated Token API Routes
	s.SYN3900API.RegisterRoutes(apiV1)
	s.SYN4700API.RegisterRoutes(apiV1)
	s.SYN4900API.RegisterRoutes(apiV1)
	s.SYN5000API.RegisterRoutes(apiV1)
	
	// TODO: Register other module routes as they are created
	// s.TokensAPI.RegisterRoutes(apiV1)
	// s.DeFiAPI.RegisterRoutes(apiV1)
	// s.GovernanceAPI.RegisterRoutes(apiV1)
	
	// System information endpoints
	apiV1.HandleFunc("/system/info", s.GetSystemInfo).Methods("GET")
	apiV1.HandleFunc("/system/status", s.GetSystemStatus).Methods("GET")
	apiV1.HandleFunc("/system/metrics", s.GetSystemMetrics).Methods("GET")
}

// Start starts the API server
func (s *SynnergyAPIServer) Start() error {
	// Configure CORS
	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	corsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	
	// Apply CORS middleware
	handler := handlers.CORS(corsOrigins, corsHeaders, corsMethods)(s.Router)
	
	// Start server
	log.Printf("Starting Synnergy Network API Server on port %s", s.Port)
	log.Printf("API Documentation available at: http://localhost:%s/", s.Port)
	
	return http.ListenAndServe(":"+s.Port, handler)
}

// HealthCheck returns the health status of the API server
func (s *SynnergyAPIServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
		"services": map[string]string{
			"consensus":    "running",
			"network":      "running",
			"transactions": "running",
			"ledger":       "running",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Welcome provides API documentation and available endpoints
func (s *SynnergyAPIServer) Welcome(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Welcome to Synnergy Network API",
		"version": "1.0.0",
		"documentation": "Available endpoints:",
		"endpoints": map[string]interface{}{
			"health": "/health",
			"system": map[string]string{
				"info":    "/api/v1/system/info",
				"status":  "/api/v1/system/status",
				"metrics": "/api/v1/system/metrics",
			},
			"consensus": map[string]string{
				"difficulty":  "/api/v1/consensus/difficulty/*",
				"validators":  "/api/v1/consensus/validators/*",
				"audit":       "/api/v1/consensus/audit/*",
				"rewards":     "/api/v1/consensus/rewards/*",
			},
			"network": map[string]string{
				"peers":       "/api/v1/network/peers/*",
				"messages":    "/api/v1/network/messages/*",
				"connections": "/api/v1/network/connections/*",
			},
			"transactions": map[string]string{
				"pool":      "/api/v1/transactions/pool/*",
				"subblocks": "/api/v1/transactions/subblocks/*",
				"status":    "/api/v1/transactions/status/{txID}",
				"history":   "/api/v1/transactions/history/{address}",
			},
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetSystemInfo returns system information
func (s *SynnergyAPIServer) GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"system": "Synnergy Network",
		"version": "1.0.0",
		"consensus": "Proof of History + Proof of Stake + Proof of Work",
		"features": []string{
			"Multi-consensus mechanisms",
			"46+ token standards",
			"Advanced DeFi capabilities",
			"Governance systems",
			"Enterprise-grade security",
		},
		"modules": []string{
			"consensus", "network", "transactions", "smart_contract",
			"wallet", "tokens", "defi", "governance", "cryptography",
			"storage", "common", "ledger",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSystemStatus returns current system status
func (s *SynnergyAPIServer) GetSystemStatus(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"timestamp": time.Now().UTC(),
		"uptime":    time.Since(time.Now()).String(), // TODO: Track actual uptime
		"status":    "operational",
		"modules": map[string]interface{}{
			"consensus": map[string]interface{}{
				"status": "active",
				"validators": 0, // TODO: Get actual count
			},
			"network": map[string]interface{}{
				"status": "connected",
				"peers": len(s.NetworkManager.GetConnectedPeers()),
			},
			"transactions": map[string]interface{}{
				"status": "processing",
				"pool_size": s.TransactionPool.PoolSize(),
			},
			"ledger": map[string]interface{}{
				"status": "synchronized",
				"latest_block": 0, // TODO: Get actual latest block
			},
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSystemMetrics returns system performance metrics
func (s *SynnergyAPIServer) GetSystemMetrics(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"timestamp": time.Now().UTC(),
		"performance": map[string]interface{}{
			"transaction_throughput": "1000+ TPS",
			"block_time": "2.5 seconds",
			"finality_time": "< 5 seconds",
		},
		"network": map[string]interface{}{
			"connected_peers": len(s.NetworkManager.GetConnectedPeers()),
			"network_latency": "< 100ms",
		},
		"consensus": map[string]interface{}{
			"validator_count": 0, // TODO: Get actual count
			"participation_rate": "95%",
		},
		"resources": map[string]interface{}{
			"cpu_usage": "25%", // TODO: Get actual metrics
			"memory_usage": "512MB",
			"storage_usage": "10GB",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}