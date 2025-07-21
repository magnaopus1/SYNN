package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"
	"github.com/gorilla/mux"
)

type SYN2400API struct{}

func NewSYN2400API() *SYN2400API { return &SYN2400API{} }

func (api *SYN2400API) RegisterRoutes(router *mux.Router) {
	// Core social media token management (15 endpoints)
	router.HandleFunc("/syn2400/tokens", api.CreateSocialToken).Methods("POST")
	router.HandleFunc("/syn2400/tokens/{tokenID}", api.GetToken).Methods("GET")
	router.HandleFunc("/syn2400/tokens", api.ListTokens).Methods("GET")
	router.HandleFunc("/syn2400/tokens/{tokenID}/transfer", api.TransferTokens).Methods("POST")
	router.HandleFunc("/syn2400/tokens/{tokenID}/tip", api.TipCreator).Methods("POST")
	router.HandleFunc("/syn2400/tokens/balance/{address}", api.GetBalance).Methods("GET")
	router.HandleFunc("/syn2400/tokens/{tokenID}/stake", api.StakeTokens).Methods("POST")
	router.HandleFunc("/syn2400/tokens/{tokenID}/rewards", api.ClaimRewards).Methods("POST")
	router.HandleFunc("/syn2400/tokens/{tokenID}/burn", api.BurnTokens).Methods("POST")
	router.HandleFunc("/syn2400/tokens/{tokenID}/mint", api.MintTokens).Methods("POST")
	router.HandleFunc("/syn2400/tokens/{tokenID}/subscription", api.SubscribeToCreator).Methods("POST")
	router.HandleFunc("/syn2400/tokens/{tokenID}/fan-token", api.CreateFanToken).Methods("POST")
	router.HandleFunc("/syn2400/tokens/{tokenID}/governance", api.GetGovernanceRights).Methods("GET")
	router.HandleFunc("/syn2400/tokens/{tokenID}/metadata", api.UpdateMetadata).Methods("PUT")
	router.HandleFunc("/syn2400/tokens/{tokenID}/history", api.GetTokenHistory).Methods("GET")

	// Content monetization (20 endpoints)
	router.HandleFunc("/syn2400/content", api.CreateContent).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}", api.GetContent).Methods("GET")
	router.HandleFunc("/syn2400/content/{contentID}/monetize", api.MonetizeContent).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}/unlock", api.UnlockContent).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}/like", api.LikeContent).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}/share", api.ShareContent).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}/comment", api.CommentOnContent).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}/revenue", api.GetContentRevenue).Methods("GET")
	router.HandleFunc("/syn2400/content/{contentID}/analytics", api.GetContentAnalytics).Methods("GET")
	router.HandleFunc("/syn2400/content/{contentID}/nft", api.CreateContentNFT).Methods("POST")
	router.HandleFunc("/syn2400/content/feed", api.GetPersonalizedFeed).Methods("GET")
	router.HandleFunc("/syn2400/content/trending", api.GetTrendingContent).Methods("GET")
	router.HandleFunc("/syn2400/content/search", api.SearchContent).Methods("GET")
	router.HandleFunc("/syn2400/content/{contentID}/collaborators", api.ManageCollaborators).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}/licensing", api.SetContentLicense).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}/copyright", api.ManageCopyright).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}/moderation", api.ModerateContent).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}/boost", api.BoostContent).Methods("POST")
	router.HandleFunc("/syn2400/content/{contentID}/schedule", api.ScheduleContent).Methods("POST")
	router.HandleFunc("/syn2400/content/batch-upload", api.BatchUploadContent).Methods("POST")

	// Creator economy (15 endpoints)
	router.HandleFunc("/syn2400/creators", api.RegisterCreator).Methods("POST")
	router.HandleFunc("/syn2400/creators/{creatorID}", api.GetCreator).Methods("GET")
	router.HandleFunc("/syn2400/creators/{creatorID}/verify", api.VerifyCreator).Methods("POST")
	router.HandleFunc("/syn2400/creators/{creatorID}/followers", api.GetFollowers).Methods("GET")
	router.HandleFunc("/syn2400/creators/{creatorID}/earnings", api.GetCreatorEarnings).Methods("GET")
	router.HandleFunc("/syn2400/creators/{creatorID}/fan-fund", api.CreateFanFund).Methods("POST")
	router.HandleFunc("/syn2400/creators/{creatorID}/merchandise", api.ManageMerchandise).Methods("POST")
	router.HandleFunc("/syn2400/creators/{creatorID}/events", api.CreateEvent).Methods("POST")
	router.HandleFunc("/syn2400/creators/{creatorID}/courses", api.CreateCourse).Methods("POST")
	router.HandleFunc("/syn2400/creators/{creatorID}/memberships", api.ManageMemberships).Methods("POST")
	router.HandleFunc("/syn2400/creators/{creatorID}/brand-deals", api.ManageBrandDeals).Methods("POST")
	router.HandleFunc("/syn2400/creators/{creatorID}/analytics", api.GetCreatorAnalytics).Methods("GET")
	router.HandleFunc("/syn2400/creators/{creatorID}/reputation", api.GetCreatorReputation).Methods("GET")
	router.HandleFunc("/syn2400/creators/leaderboard", api.GetCreatorLeaderboard).Methods("GET")
	router.HandleFunc("/syn2400/creators/{creatorID}/sponsor", api.SponsorCreator).Methods("POST")

	// Social features (10 endpoints)
	router.HandleFunc("/syn2400/social/follow", api.FollowUser).Methods("POST")
	router.HandleFunc("/syn2400/social/unfollow", api.UnfollowUser).Methods("POST")
	router.HandleFunc("/syn2400/social/friends", api.GetFriends).Methods("GET")
	router.HandleFunc("/syn2400/social/groups", api.ManageGroups).Methods("POST")
	router.HandleFunc("/syn2400/social/communities", api.JoinCommunity).Methods("POST")
	router.HandleFunc("/syn2400/social/messages", api.SendMessage).Methods("POST")
	router.HandleFunc("/syn2400/social/notifications", api.GetNotifications).Methods("GET")
	router.HandleFunc("/syn2400/social/mentions", api.GetMentions).Methods("GET")
	router.HandleFunc("/syn2400/social/hashtags", api.GetHashtagTrends).Methods("GET")
	router.HandleFunc("/syn2400/social/influence", api.GetInfluenceScore).Methods("GET")

	// Analytics and reporting (10 endpoints)
	router.HandleFunc("/syn2400/analytics/engagement", api.GetEngagementAnalytics).Methods("GET")
	router.HandleFunc("/syn2400/analytics/revenue", api.GetRevenueAnalytics).Methods("GET")
	router.HandleFunc("/syn2400/analytics/growth", api.GetGrowthMetrics).Methods("GET")
	router.HandleFunc("/syn2400/analytics/demographics", api.GetDemographics).Methods("GET")
	router.HandleFunc("/syn2400/analytics/performance", api.GetPerformanceMetrics).Methods("GET")
	router.HandleFunc("/syn2400/reports/creator", api.GenerateCreatorReport).Methods("GET")
	router.HandleFunc("/syn2400/reports/platform", api.GeneratePlatformReport).Methods("GET")
	router.HandleFunc("/syn2400/analytics/trends", api.GetTrendAnalysis).Methods("GET")
	router.HandleFunc("/syn2400/analytics/roi", api.GetROIAnalytics).Methods("GET")
	router.HandleFunc("/syn2400/analytics/retention", api.GetRetentionMetrics).Methods("GET")

	// Administrative (5 endpoints)
	router.HandleFunc("/syn2400/admin/settings", api.UpdateSettings).Methods("PUT")
	router.HandleFunc("/syn2400/admin/moderation", api.SetModerationPolicies).Methods("POST")
	router.HandleFunc("/syn2400/admin/health", api.GetSystemHealth).Methods("GET")
	router.HandleFunc("/syn2400/admin/logs", api.GetSystemLogs).Methods("GET")
	router.HandleFunc("/syn2400/admin/backup", api.CreateBackup).Methods("POST")
}

func (api *SYN2400API) CreateSocialToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("Creating social media token - Request from %s", r.RemoteAddr)
	
	var request struct {
		Name         string  `json:"name" validate:"required"`
		Symbol       string  `json:"symbol" validate:"required"`
		TotalSupply  float64 `json:"total_supply" validate:"required,min=1"`
		CreatorID    string  `json:"creator_id" validate:"required"`
		Platform     string  `json:"platform" validate:"required"`
		TokenType    string  `json:"token_type" validate:"required,oneof=creator fan utility governance"`
		Utility      string  `json:"utility" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error":"Invalid JSON format","code":"INVALID_JSON"}`, http.StatusBadRequest)
		return
	}

	tokenID := fmt.Sprintf("SOCIAL_%d", time.Now().UnixNano())
	
	response := map[string]interface{}{
		"success":      true,
		"token_id":     tokenID,
		"name":         request.Name,
		"symbol":       request.Symbol,
		"total_supply": request.TotalSupply,
		"creator_id":   request.CreatorID,
		"platform":     request.Platform,
		"token_type":   request.TokenType,
		"utility":      request.Utility,
		"created_at":   time.Now().Format(time.RFC3339),
		"status":       "active",
		"network":      "synnergy",
		"decimals":     18,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	log.Printf("Social media token created: %s", tokenID)
}

func (api *SYN2400API) GetToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]
	response := map[string]interface{}{"success": true, "token_id": tokenID, "token_type": "creator", "status": "active"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *SYN2400API) GetEngagementAnalytics(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"success": true, "engagement": map[string]interface{}{"likes": 15000, "shares": 2500, "comments": 5000, "rate": 8.5}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}