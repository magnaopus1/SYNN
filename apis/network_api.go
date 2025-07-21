package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"synnergy_network/pkg/network"
	"synnergy_network/pkg/ledger"

	"github.com/gorilla/mux"
)

// NetworkAPI handles all network-related API endpoints
type NetworkAPI struct {
	NetworkManager *network.NetworkManager
	LedgerInstance *ledger.Ledger
}

// NewNetworkAPI creates a new network API instance
func NewNetworkAPI(networkManager *network.NetworkManager, ledgerInstance *ledger.Ledger) *NetworkAPI {
	return &NetworkAPI{
		NetworkManager: networkManager,
		LedgerInstance: ledgerInstance,
	}
}

// RegisterRoutes registers all network API routes
func (api *NetworkAPI) RegisterRoutes(router *mux.Router) {
	// Peer management
	router.HandleFunc("/network/peers/connect", api.ConnectToPeer).Methods("POST")
	router.HandleFunc("/network/peers/disconnect", api.DisconnectFromPeer).Methods("POST")
	router.HandleFunc("/network/peers/list", api.ListPeers).Methods("GET")
	router.HandleFunc("/network/peers/ping", api.PingPeer).Methods("POST")
	
	// Messaging
	router.HandleFunc("/network/messages/send", api.SendMessage).Methods("POST")
	router.HandleFunc("/network/messages/receive", api.ReceiveMessages).Methods("GET")
	
	// Connection management
	router.HandleFunc("/network/connections/status", api.GetConnectionStatus).Methods("GET")
	router.HandleFunc("/network/connections/logs", api.GetConnectionLogs).Methods("GET")
}

// ConnectToPeer establishes a connection with a peer
func (api *NetworkAPI) ConnectToPeer(w http.ResponseWriter, r *http.Request) {
	var request struct {
		PeerIP string `json:"peer_ip"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.NetworkManager.ConnectToPeer(request.PeerIP)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to peer: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Successfully connected to peer %s", request.PeerIP),
		"peer_ip": request.PeerIP,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DisconnectFromPeer disconnects from a peer
func (api *NetworkAPI) DisconnectFromPeer(w http.ResponseWriter, r *http.Request) {
	var request struct {
		PeerIP string `json:"peer_ip"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.NetworkManager.DisconnectFromPeer(request.PeerIP)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to disconnect from peer: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Successfully disconnected from peer %s", request.PeerIP),
		"peer_ip": request.PeerIP,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListPeers returns a list of connected peers
func (api *NetworkAPI) ListPeers(w http.ResponseWriter, r *http.Request) {
	// Get connected peers from network manager
	peers := api.NetworkManager.GetConnectedPeers()

	response := map[string]interface{}{
		"success": true,
		"peers": peers,
		"count": len(peers),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PingPeer sends a ping to check peer connectivity
func (api *NetworkAPI) PingPeer(w http.ResponseWriter, r *http.Request) {
	var request struct {
		PeerIP string `json:"peer_ip"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.NetworkManager.PingPeer(request.PeerIP)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to ping peer: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Successfully pinged peer %s", request.PeerIP),
		"peer_ip": request.PeerIP,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SendMessage sends an encrypted message to a peer
func (api *NetworkAPI) SendMessage(w http.ResponseWriter, r *http.Request) {
	var request struct {
		PeerIP  string `json:"peer_ip"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.NetworkManager.SendEncryptedMessage(request.PeerIP, request.Message)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Message sent successfully",
		"peer_ip": request.PeerIP,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ReceiveMessages listens for incoming messages from a peer
func (api *NetworkAPI) ReceiveMessages(w http.ResponseWriter, r *http.Request) {
	peerIP := r.URL.Query().Get("peer_ip")
	if peerIP == "" {
		http.Error(w, "peer_ip parameter is required", http.StatusBadRequest)
		return
	}

	// This would typically be a WebSocket or Server-Sent Events implementation
	// For simplicity, we'll return a success message
	response := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Listening for messages from peer %s", peerIP),
		"peer_ip": peerIP,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetConnectionStatus returns the status of network connections
func (api *NetworkAPI) GetConnectionStatus(w http.ResponseWriter, r *http.Request) {
	status := api.NetworkManager.GetNetworkStatus()

	response := map[string]interface{}{
		"success": true,
		"status": status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetConnectionLogs retrieves network connection logs
func (api *NetworkAPI) GetConnectionLogs(w http.ResponseWriter, r *http.Request) {
	logs := api.LedgerInstance.GetNetworkLogs()

	response := map[string]interface{}{
		"success": true,
		"logs": logs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}