package sustainability

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewCarbonOffsetMatchingSystem initializes the carbon offset matching system
func NewCarbonOffsetMatchingSystem(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *CarbonOffsetMatchingSystem {
	return &CarbonOffsetMatchingSystem{
		Requests:         make(map[string]*OffsetRequest),
		Matches:          make(map[string]*CarbonOffsetMatch),
		AvailableCredits: make(map[string]*Syn700Token),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// AddOffsetRequest allows a user to request carbon offsets for their emissions
func (coms *CarbonOffsetMatchingSystem) AddOffsetRequest(requestID, requester string, amount float64) (*OffsetRequest, error) {
	coms.mu.Lock()
	defer coms.mu.Unlock()

	// Encrypt request data
	requestData := fmt.Sprintf("RequestID: %s, Requester: %s, OffsetAmount: %f", requestID, requester, amount)
	encryptedData, err := coms.EncryptionService.EncryptData([]byte(requestData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt offset request data: %v", err)
	}

	// Create the offset request
	request := &OffsetRequest{
		RequestID:     requestID,
		Requester:     requester,
		OffsetAmount:  amount,
		RequestedTime: time.Now(),
		IsFulfilled:   false,
	}

	// Add the request to the system
	coms.Requests[requestID] = request

	// Log the offset request in the ledger
	err = coms.Ledger.RecordOffsetRequest(requestID, requester, amount, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log offset request: %v", err)
	}

	fmt.Printf("Offset request %s made by %s for %f tons of CO2\n", requestID, requester, amount)
	return request, nil
}

// AddAvailableCredit adds an available carbon credit to the system for matching
func (coms *CarbonOffsetMatchingSystem) AddAvailableCredit(creditID string, credit *Syn700Token) error {
	coms.mu.Lock()
	defer coms.mu.Unlock()

	// Ensure the carbon credit is not retired
	if credit.IsRetired {
		return fmt.Errorf("carbon credit %s has already been retired", creditID)
	}

	// Add the credit to the available credits pool
	coms.AvailableCredits[creditID] = credit

	// Log the addition of available credits in the ledger
	err := coms.Ledger.RecordAvailableCredit(creditID, credit.TokenID, credit.Owner, credit.IssuedAmount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log available carbon credit: %v", err)
	}

	fmt.Printf("Carbon credit %s added as available for matching\n", creditID)
	return nil
}

// MatchOffsetRequest attempts to match an offset request with available carbon credits
func (coms *CarbonOffsetMatchingSystem) MatchOffsetRequest(requestID string) (*CarbonOffsetMatch, error) {
	coms.mu.Lock()
	defer coms.mu.Unlock()

	// Retrieve the offset request
	request, exists := coms.Requests[requestID]
	if !exists || request.IsFulfilled {
		return nil, fmt.Errorf("offset request %s not found or already fulfilled", requestID)
	}

	// Attempt to match with available credits
	var matchedAmount float64
	var matchedCredit *Syn700Token

	for creditID, credit := range coms.AvailableCredits {
		// Check if the credit can fulfill the request
		if credit.IssuedAmount >= request.OffsetAmount && !credit.IsRetired {
			matchedAmount = request.OffsetAmount
			matchedCredit = credit

			// Mark the credit as retired
			credit.IsRetired = true
			credit.RetiredTime = time.Now()

			// Remove the credit from available credits
			delete(coms.AvailableCredits, creditID)
			break
		}
	}

	if matchedCredit == nil {
		return nil, fmt.Errorf("no available carbon credit could match the offset request %s", requestID)
	}

	// Mark the request as fulfilled
	request.IsFulfilled = true
	request.FulfilledTime = time.Now()

	// Create and log the offset match
	match := &CarbonOffsetMatch{
		MatchID:       generateUniqueID(),
		RequestID:     requestID,
		CreditID:      matchedCredit.CreditID,
		MatchedAmount: matchedAmount,
		MatchedTime:   time.Now(),
	}
	coms.Matches[match.MatchID] = match

	// Log the match in the ledger
	err := coms.Ledger.RecordOffsetMatch(match.MatchID, requestID, matchedCredit.CreditID, matchedAmount, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log carbon offset match: %v", err)
	}

	fmt.Printf("Offset request %s fulfilled with carbon credit %s for %f tons of CO2\n", requestID, matchedCredit.CreditID, matchedAmount)
	return match, nil
}

// ViewOffsetRequest allows viewing the details of a specific offset request
func (coms *CarbonOffsetMatchingSystem) ViewOffsetRequest(requestID string) (*OffsetRequest, error) {
	coms.mu.Lock()
	defer coms.mu.Unlock()

	// Retrieve the offset request
	request, exists := coms.Requests[requestID]
	if !exists {
		return nil, fmt.Errorf("offset request %s not found", requestID)
	}

	return request, nil
}

