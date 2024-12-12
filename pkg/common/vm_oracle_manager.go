package common

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Oracle represents a registered oracle.
type Oracle struct {
    ID          string
    URL         string
    PublicKey   []byte
    LastActive  time.Time
    Description string
}

// OracleRequest represents a request for external data.
type OracleRequest struct {
    RequestID    string
    DataSource   string
    Parameters   map[string]interface{}
    Timestamp    time.Time
    ResponseChan chan OracleResponse
}

// OracleResponse represents a response from an oracle.
type OracleResponse struct {
    RequestID  string
    Data       interface{}
    Signature  []byte
    Timestamp  time.Time
    OracleID   string
    Verified   bool
    Error      error
}

// OracleManager manages oracle interactions.
type OracleManager struct {
    oracles           map[string]*Oracle           // Registered oracles
    requests          map[string]*OracleRequest    // Active requests
    responses         map[string]OracleResponse    // Received responses
    cache             map[string]OracleResponse    // Cached responses
    cacheDuration     time.Duration                // Duration to keep responses in cache
    mutex             sync.RWMutex                 // Mutex for thread safety
    requestTimeout    time.Duration                // Timeout for oracle requests
    maxRetries        int                          // Maximum retries for failed requests
    trustedOracleIDs  []string                     // List of trusted oracle IDs
    verificationFunc  func(OracleResponse) bool    // Function to verify responses
    processingQueue   chan OracleRequest           // Queue for processing requests
    stopChan          chan struct{}                // Channel to signal stopping
    processingWorkers int                          // Number of worker goroutines
	logger           *LogManager // Adding logger to OracleManager

}

// NewOracleManager initializes a new OracleManager.
func NewOracleManager(cacheDuration, requestTimeout time.Duration, maxRetries, processingWorkers int, trustedOracleIDs []string) *OracleManager {
    manager := &OracleManager{
        oracles:           make(map[string]*Oracle),
        requests:          make(map[string]*OracleRequest),
        responses:         make(map[string]OracleResponse),
        cache:             make(map[string]OracleResponse),
        cacheDuration:     cacheDuration,
        requestTimeout:    requestTimeout,
        maxRetries:        maxRetries,
        trustedOracleIDs:  trustedOracleIDs,
        verificationFunc:  defaultVerificationFunc,
        processingQueue:   make(chan OracleRequest, 1000),
        stopChan:          make(chan struct{}),
        processingWorkers: processingWorkers,
    }

    // Start processing workers
    for i := 0; i < processingWorkers; i++ {
        go manager.requestProcessor()
    }

    // Start cache cleaner
    go manager.cacheCleaner()

    return manager
}

// defaultVerificationFunc is the default function to verify oracle responses.
func defaultVerificationFunc(response OracleResponse) bool {
    // Implement cryptographic verification here
    // For example, verify the signature using the oracle's public key
    return true // Placeholder implementation
}

// RegisterOracle registers a new oracle with the manager.
func (manager *OracleManager) RegisterOracle(oracle Oracle) error {
    manager.mutex.Lock()
    defer manager.mutex.Unlock()

    if _, exists := manager.oracles[oracle.ID]; exists {
        return errors.New("oracle already registered")
    }

    manager.oracles[oracle.ID] = &oracle
    log.Printf("Oracle %s registered.\n", oracle.ID)
    return nil
}

// UnregisterOracle removes an oracle from the manager.
func (manager *OracleManager) UnregisterOracle(oracleID string) error {
    manager.mutex.Lock()
    defer manager.mutex.Unlock()

    if _, exists := manager.oracles[oracleID]; !exists {
        return errors.New("oracle not found")
    }

    delete(manager.oracles, oracleID)
    log.Printf("Oracle %s unregistered.\n", oracleID)
    return nil
}

// CreateRequest creates a new oracle request.
func (manager *OracleManager) CreateRequest(dataSource string, parameters map[string]interface{}) (string, error) {
    requestID := generateRequestID(dataSource, parameters)
    request := &OracleRequest{
        RequestID:    requestID,
        DataSource:   dataSource,
        Parameters:   parameters,
        Timestamp:    time.Now(),
        ResponseChan: make(chan OracleResponse, 1),
    }

    manager.mutex.Lock()
    manager.requests[requestID] = request
    manager.mutex.Unlock()

    // Add request to processing queue
    manager.processingQueue <- *request

    return requestID, nil
}

// GetResponse retrieves the response for a given request ID.
func (manager *OracleManager) GetResponse(requestID string) (OracleResponse, error) {
    // Check cache first
    manager.mutex.RLock()
    if response, exists := manager.cache[requestID]; exists {
        manager.mutex.RUnlock()
        return response, nil
    }
    manager.mutex.RUnlock()

    // Wait for response
    manager.mutex.RLock()
    request, exists := manager.requests[requestID]
    manager.mutex.RUnlock()

    if !exists {
        return OracleResponse{}, errors.New("request not found")
    }

    select {
    case response := <-request.ResponseChan:
        if response.Error != nil {
            return response, response.Error
        }
        // Cache the response
        manager.cacheResponse(requestID, response)
        return response, nil
    case <-time.After(manager.requestTimeout):
        return OracleResponse{}, errors.New("request timed out")
    }
}

// requestProcessor processes oracle requests.
func (manager *OracleManager) requestProcessor() {
    for {
        select {
        case request := <-manager.processingQueue:
            manager.processRequest(request)
        case <-manager.stopChan:
            return
        }
    }
}

// processRequest sends the request to registered oracles and handles the response.
func (manager *OracleManager) processRequest(request OracleRequest) {
    var wg sync.WaitGroup
    responsesChan := make(chan OracleResponse, len(manager.trustedOracleIDs))

    manager.mutex.RLock()
    oracles := manager.getTrustedOracles()
    manager.mutex.RUnlock()

    for _, oracle := range oracles {
        wg.Add(1)
        go func(oracle *Oracle) {
            defer wg.Done()
            response, err := manager.sendRequestToOracle(request, oracle)
            if err != nil {
                log.Printf("Error receiving response from oracle %s: %v\n", oracle.ID, err)
                return
            }
            responsesChan <- response
        }(oracle)
    }

    wg.Wait()
    close(responsesChan)

    // Collect and verify responses
    verifiedResponse, err := manager.collectAndVerifyResponses(request.RequestID, responsesChan)
    if err != nil {
        manager.mutex.RLock()
        reqPtr, exists := manager.requests[request.RequestID]
        manager.mutex.RUnlock()
        if exists && reqPtr != nil {
            reqPtr.ResponseChan <- OracleResponse{
                RequestID: reqPtr.RequestID,
                Error:     err,
            }
        }
        return
    }

    // Send response back to requester
    manager.mutex.RLock()
    reqPtr, exists := manager.requests[request.RequestID]
    manager.mutex.RUnlock()
    if exists && reqPtr != nil {
        reqPtr.ResponseChan <- verifiedResponse
    }
}


// sendRequestToOracle sends a request to a specific oracle and waits for a response over HTTP.
func (manager *OracleManager) sendRequestToOracle(request OracleRequest, oracle *Oracle) (OracleResponse, error) {
	// Prepare request data
	requestData := []byte(fmt.Sprintf("requestID=%s&dataSource=%s", request.RequestID, request.DataSource))
	req, err := http.NewRequest("POST", oracle.URL, bytes.NewBuffer(requestData))
	if err != nil {
		return OracleResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform HTTP request to oracle
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return OracleResponse{}, fmt.Errorf("failed to send request to oracle %s: %v", oracle.ID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return OracleResponse{}, fmt.Errorf("oracle %s returned non-200 status: %d", oracle.ID, resp.StatusCode)
	}

	// Parse oracle response
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return OracleResponse{}, fmt.Errorf("failed to read response from oracle %s: %v", oracle.ID, err)
	}

	// Create OracleResponse
	response := OracleResponse{
		RequestID: request.RequestID,
		Data:      string(responseData),
		Timestamp: time.Now(),
		OracleID:  oracle.ID,
	}

	// Verify the response's signature using oracle's public key
	err = manager.verifyResponseSignature(response, oracle.PublicKey)
	if err != nil {
		return response, fmt.Errorf("failed to verify response from oracle %s: %v", oracle.ID, err)
	}
	response.Verified = true

	return response, nil
}

// verifyResponseSignature verifies the signature on an oracle response.
func (manager *OracleManager) verifyResponseSignature(response OracleResponse, publicKey []byte) error {
	// Decode public key
	block, _ := pem.Decode(publicKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return errors.New("failed to decode oracle public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("invalid public key: %v", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return errors.New("not an RSA public key")
	}

	// Verify response data and timestamp
	message := fmt.Sprintf("%s%s%s", response.RequestID, response.Data, response.Timestamp)
	hashed := sha256.Sum256([]byte(message))

	// Verify the signature (assuming `Signature` is a valid RSA signature)
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], response.Signature)
	if err != nil {
		return errors.New("signature verification failed")
	}

	return nil
}

// collectAndVerifyResponses collects verified responses and aggregates them to form the final result.
func (manager *OracleManager) collectAndVerifyResponses(requestID string, responsesChan chan OracleResponse) (OracleResponse, error) {
	var verifiedResponses []OracleResponse
	responseCounts := make(map[string]int) // Track the frequency of each unique response
	highestCount := 0
	var finalResponse OracleResponse

	// Collect all verified responses
	for response := range responsesChan {
		if response.Verified {
			verifiedResponses = append(verifiedResponses, response)

			// Ensure `response.Data` is a string for map indexing
			dataStr, ok := response.Data.(string)
			if !ok {
				return OracleResponse{}, errors.New("response data is not a string")
			}
			responseCounts[dataStr]++

			// Determine the majority response
			if responseCounts[dataStr] > highestCount {
				highestCount = responseCounts[dataStr]
				finalResponse = response
			}
		}
	}

	if len(verifiedResponses) == 0 {
		return OracleResponse{}, errors.New("no verified responses received from oracles")
	}

	// Log final aggregation details
	if manager.logger != nil {
		manager.logger.Info("Oracle responses aggregated", map[string]interface{}{
			"requestID":      requestID,
			"totalResponses": len(verifiedResponses),
			"finalResponse":  finalResponse.Data,
		})
	}

	return finalResponse, nil
}

// cacheResponse stores the response in the cache.
func (manager *OracleManager) cacheResponse(requestID string, response OracleResponse) {
    manager.mutex.Lock()
    defer manager.mutex.Unlock()
    manager.cache[requestID] = response
}

// cacheCleaner periodically cleans up expired cache entries.
func (manager *OracleManager) cacheCleaner() {
    ticker := time.NewTicker(manager.cacheDuration)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            manager.mutex.Lock()
            for key, response := range manager.cache {
                if time.Since(response.Timestamp) > manager.cacheDuration {
                    delete(manager.cache, key)
                }
            }
            manager.mutex.Unlock()
        case <-manager.stopChan:
            return
        }
    }
}

// getTrustedOracles returns a list of trusted oracles.
func (manager *OracleManager) getTrustedOracles() []*Oracle {
    oracles := []*Oracle{}
    for _, oracleID := range manager.trustedOracleIDs {
        if oracle, exists := manager.oracles[oracleID]; exists {
            oracles = append(oracles, oracle)
        }
    }
    return oracles
}

// generateRequestID generates a unique ID for a request.
func generateRequestID(dataSource string, parameters map[string]interface{}) string {
    data := fmt.Sprintf("%s:%v:%d", dataSource, parameters, time.Now().UnixNano())
    hash := sha256.Sum256([]byte(data))
    return fmt.Sprintf("%x", hash)
}

// Stop gracefully stops the OracleManager.
func (manager *OracleManager) Stop() {
    close(manager.stopChan)
    log.Println("OracleManager stopped.")
}
