package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"
	"github.com/gorilla/mux"
)

type SYN2500API struct{}

func NewSYN2500API() *SYN2500API { return &SYN2500API{} }

func (api *SYN2500API) RegisterRoutes(router *mux.Router) {
	// Core music royalty token management (15 endpoints)
	router.HandleFunc("/syn2500/tokens", api.CreateMusicToken).Methods("POST")
	router.HandleFunc("/syn2500/tokens/{tokenID}", api.GetToken).Methods("GET")
	router.HandleFunc("/syn2500/tokens", api.ListTokens).Methods("GET")
	router.HandleFunc("/syn2500/tokens/{tokenID}/transfer", api.TransferTokens).Methods("POST")
	router.HandleFunc("/syn2500/tokens/{tokenID}/royalties", api.DistributeRoyalties).Methods("POST")
	router.HandleFunc("/syn2500/tokens/balance/{address}", api.GetBalance).Methods("GET")
	router.HandleFunc("/syn2500/tokens/{tokenID}/metadata", api.UpdateMetadata).Methods("PUT")
	router.HandleFunc("/syn2500/tokens/{tokenID}/fractionalize", api.FractionalizeRights).Methods("POST")
	router.HandleFunc("/syn2500/tokens/{tokenID}/license", api.CreateLicense).Methods("POST")
	router.HandleFunc("/syn2500/tokens/{tokenID}/stream-data", api.UpdateStreamData).Methods("POST")
	router.HandleFunc("/syn2500/tokens/{tokenID}/copyright", api.ManageCopyright).Methods("POST")
	router.HandleFunc("/syn2500/tokens/{tokenID}/collaboration", api.ManageCollaboration).Methods("POST")
	router.HandleFunc("/syn2500/tokens/{tokenID}/history", api.GetRoyaltyHistory).Methods("GET")
	router.HandleFunc("/syn2500/tokens/{tokenID}/analytics", api.GetTokenAnalytics).Methods("GET")
	router.HandleFunc("/syn2500/tokens/{tokenID}/verify", api.VerifyOwnership).Methods("POST")

	// Music catalog management (20 endpoints)
	router.HandleFunc("/syn2500/catalog/songs", api.RegisterSong).Methods("POST")
	router.HandleFunc("/syn2500/catalog/songs/{songID}", api.GetSong).Methods("GET")
	router.HandleFunc("/syn2500/catalog/songs", api.ListSongs).Methods("GET")
	router.HandleFunc("/syn2500/catalog/albums", api.CreateAlbum).Methods("POST")
	router.HandleFunc("/syn2500/catalog/albums/{albumID}", api.GetAlbum).Methods("GET")
	router.HandleFunc("/syn2500/catalog/artists", api.RegisterArtist).Methods("POST")
	router.HandleFunc("/syn2500/catalog/artists/{artistID}", api.GetArtist).Methods("GET")
	router.HandleFunc("/syn2500/catalog/labels", api.RegisterLabel).Methods("POST")
	router.HandleFunc("/syn2500/catalog/publishers", api.RegisterPublisher).Methods("POST")
	router.HandleFunc("/syn2500/catalog/genres", api.ManageGenres).Methods("POST")
	router.HandleFunc("/syn2500/catalog/isrc", api.RegisterISRC).Methods("POST")
	router.HandleFunc("/syn2500/catalog/upc", api.RegisterUPC).Methods("POST")
	router.HandleFunc("/syn2500/catalog/composers", api.RegisterComposer).Methods("POST")
	router.HandleFunc("/syn2500/catalog/producers", api.RegisterProducer).Methods("POST")
	router.HandleFunc("/syn2500/catalog/search", api.SearchCatalog).Methods("GET")
	router.HandleFunc("/syn2500/catalog/trending", api.GetTrendingMusic).Methods("GET")
	router.HandleFunc("/syn2500/catalog/recommendations", api.GetRecommendations).Methods("GET")
	router.HandleFunc("/syn2500/catalog/playlists", api.ManagePlaylists).Methods("POST")
	router.HandleFunc("/syn2500/catalog/charts", api.GetCharts).Methods("GET")
	router.HandleFunc("/syn2500/catalog/backup", api.BackupCatalog).Methods("POST")

	// Royalty distribution (15 endpoints)
	router.HandleFunc("/syn2500/royalties/calculate", api.CalculateRoyalties).Methods("POST")
	router.HandleFunc("/syn2500/royalties/distribute", api.DistributePayments).Methods("POST")
	router.HandleFunc("/syn2500/royalties/statements", api.GenerateStatements).Methods("GET")
	router.HandleFunc("/syn2500/royalties/splits", api.ManageRoyaltySplits).Methods("POST")
	router.HandleFunc("/syn2500/royalties/advance", api.ProcessAdvance).Methods("POST")
	router.HandleFunc("/syn2500/royalties/escrow", api.ManageEscrow).Methods("POST")
	router.HandleFunc("/syn2500/royalties/audit", api.AuditRoyalties).Methods("POST")
	router.HandleFunc("/syn2500/royalties/disputes", api.HandleDisputes).Methods("POST")
	router.HandleFunc("/syn2500/royalties/collections", api.GetCollectionSocieties).Methods("GET")
	router.HandleFunc("/syn2500/royalties/mechanical", api.ProcessMechanicalRoyalties).Methods("POST")
	router.HandleFunc("/syn2500/royalties/performance", api.ProcessPerformanceRoyalties).Methods("POST")
	router.HandleFunc("/syn2500/royalties/sync", api.ProcessSyncRoyalties).Methods("POST")
	router.HandleFunc("/syn2500/royalties/digital", api.ProcessDigitalRoyalties).Methods("POST")
	router.HandleFunc("/syn2500/royalties/international", api.ProcessInternationalRoyalties).Methods("POST")
	router.HandleFunc("/syn2500/royalties/tax", api.HandleTaxReporting).Methods("POST")

	// Licensing and rights (10 endpoints)
	router.HandleFunc("/syn2500/licensing/sync", api.CreateSyncLicense).Methods("POST")
	router.HandleFunc("/syn2500/licensing/mechanical", api.CreateMechanicalLicense).Methods("POST")
	router.HandleFunc("/syn2500/licensing/performance", api.CreatePerformanceLicense).Methods("POST")
	router.HandleFunc("/syn2500/licensing/master", api.CreateMasterLicense).Methods("POST")
	router.HandleFunc("/syn2500/licensing/publishing", api.CreatePublishingLicense).Methods("POST")
	router.HandleFunc("/syn2500/licensing/sample", api.CreateSampleLicense).Methods("POST")
	router.HandleFunc("/syn2500/licensing/cover", api.CreateCoverLicense).Methods("POST")
	router.HandleFunc("/syn2500/licensing/remix", api.CreateRemixLicense).Methods("POST")
	router.HandleFunc("/syn2500/licensing/broadcast", api.CreateBroadcastLicense).Methods("POST")
	router.HandleFunc("/syn2500/licensing/streaming", api.CreateStreamingLicense).Methods("POST")

	// Analytics and reporting (10 endpoints)
	router.HandleFunc("/syn2500/analytics/streams", api.GetStreamingAnalytics).Methods("GET")
	router.HandleFunc("/syn2500/analytics/revenue", api.GetRevenueAnalytics).Methods("GET")
	router.HandleFunc("/syn2500/analytics/demographics", api.GetDemographics).Methods("GET")
	router.HandleFunc("/syn2500/analytics/geographic", api.GetGeographicAnalytics).Methods("GET")
	router.HandleFunc("/syn2500/analytics/platform", api.GetPlatformAnalytics).Methods("GET")
	router.HandleFunc("/syn2500/reports/royalty", api.GenerateRoyaltyReport).Methods("GET")
	router.HandleFunc("/syn2500/reports/performance", api.GeneratePerformanceReport).Methods("GET")
	router.HandleFunc("/syn2500/analytics/trends", api.GetTrendAnalysis).Methods("GET")
	router.HandleFunc("/syn2500/analytics/market", api.GetMarketAnalytics).Methods("GET")
	router.HandleFunc("/syn2500/analytics/forecast", api.GetRevenueForecast).Methods("GET")

	// Administrative (5 endpoints)
	router.HandleFunc("/syn2500/admin/settings", api.UpdateSettings).Methods("PUT")
	router.HandleFunc("/syn2500/admin/compliance", api.ManageCompliance).Methods("POST")
	router.HandleFunc("/syn2500/admin/health", api.GetSystemHealth).Methods("GET")
	router.HandleFunc("/syn2500/admin/logs", api.GetSystemLogs).Methods("GET")
	router.HandleFunc("/syn2500/admin/backup", api.CreateBackup).Methods("POST")
}

func (api *SYN2500API) CreateMusicToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("Creating music royalty token - Request from %s", r.RemoteAddr)
	
	var request struct {
		Title        string   `json:"title" validate:"required"`
		Artist       string   `json:"artist" validate:"required"`
		Album        string   `json:"album"`
		ISRC         string   `json:"isrc" validate:"required"`
		Duration     int      `json:"duration" validate:"required,min=1"`
		Genre        string   `json:"genre" validate:"required"`
		ReleaseDate  string   `json:"release_date" validate:"required"`
		Rights       []string `json:"rights" validate:"required"`
		RoyaltyRate  float64  `json:"royalty_rate" validate:"required,min=0,max=100"`
		Collaborators []map[string]interface{} `json:"collaborators"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error":"Invalid JSON format","code":"INVALID_JSON"}`, http.StatusBadRequest)
		return
	}

	tokenID := fmt.Sprintf("MUSIC_%d", time.Now().UnixNano())
	
	response := map[string]interface{}{
		"success":      true,
		"token_id":     tokenID,
		"title":        request.Title,
		"artist":       request.Artist,
		"album":        request.Album,
		"isrc":         request.ISRC,
		"duration":     request.Duration,
		"genre":        request.Genre,
		"release_date": request.ReleaseDate,
		"rights":       request.Rights,
		"royalty_rate": request.RoyaltyRate,
		"collaborators": request.Collaborators,
		"created_at":   time.Now().Format(time.RFC3339),
		"status":       "active",
		"network":      "synnergy",
		"contract_address": fmt.Sprintf("0x%s", tokenID),
		"copyright_status": "registered",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	log.Printf("Music royalty token created: %s", tokenID)
}

func (api *SYN2500API) GetToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]
	response := map[string]interface{}{"success": true, "token_id": tokenID, "type": "music_royalty", "status": "active"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *SYN2500API) GetStreamingAnalytics(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"success": true, "analytics": map[string]interface{}{"total_streams": 1500000, "revenue": 12500.0, "platforms": 15}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}