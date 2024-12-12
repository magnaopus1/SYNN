package syn3300

import (
	"sync"
	"time"

)

// Syn3300Token represents a SYN3300 token with full details.
type Syn3300Token struct {
	ID               string               // Unique ID of the SYN3300 token
	Name             string               // Name of the ETF
	TotalSupply      float64              // Total supply of the ETF
	Value            float64              // Current value of the ETF
	ETFMetadata      ETFMetadata          // Metadata for the ETF
	ETFPortfolio     ETFPortfolioDetails  // Details of the ETF's portfolio
	mutex            sync.Mutex           // Mutex for thread-safe operations
	ledgerService    *ledger.Ledger       // Ledger integration
	encryptionService *encryption.Encryptor // Encryption service
	consensusService *consensus.SynnergyConsensus // Consensus service
}

// NewSyn3300Token creates a new SYN3300 token with all required services.
func NewSyn3300Token(id, name string, totalSupply, value float64, ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *Syn3300Token {
	return &Syn3300Token{
		ID:               id,
		Name:             name,
		TotalSupply:      totalSupply,
		Value:            value,
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService: consensus,
	}
}

// ETFLink represents the linking information of an ETF share.
type ETFLink struct {
	ETFID        string    `json:"etf_id"`         // ETF identifier
	ShareTokenID string    `json:"share_token_id"` // Associated share token
	Owner        string    `json:"owner"`          // Owner of the share
	Timestamp    time.Time `json:"timestamp"`      // Timestamp of the link
}

// ETFLinkingService provides methods to link ETF shares to specific ETFs.
type ETFLinkingService struct {
	ledgerService    *ledger.Ledger           // Ledger service for logging
	encryptionService *encryption.Encryptor    // Encryption service for securing the link
	consensusService *consensus.SynnergyConsensus // Consensus service for validating link updates
	mutex            sync.Mutex               // Mutex for thread-safe operations
}

// NewETFLinkingService creates a new ETF linking service.
func NewETFLinkingService(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *ETFLinkingService {
	return &ETFLinkingService{
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService: consensus,
	}
}

// LinkETFShare links a share token to an ETF.
func (els *ETFLinkingService) LinkETFShare(etfID, shareTokenID, owner string) (*ETFLink, error) {
	els.mutex.Lock()
	defer els.mutex.Unlock()

	// Create a new ETFLink.
	link := &ETFLink{
		ETFID:        etfID,
		ShareTokenID: shareTokenID,
		Owner:        owner,
		Timestamp:    time.Now(),
	}

	// Encrypt the link for security.
	encryptedLink, err := els.encryptionService.EncryptData(link)
	if err != nil {
		return nil, err
	}

	// Log the link in the ledger.
	els.ledgerService.LogEvent("ETFShareLinked", time.Now(), shareTokenID)

	// Validate the linking using consensus.
	err = els.consensusService.ValidateSubBlock(shareTokenID)
	if err != nil {
		return nil, err
	}

	return encryptedLink.(*ETFLink), nil
}

// ETFMetadata represents the metadata of an ETF.
type ETFMetadata struct {
	ETFID          string    `json:"etf_id"`         // ETF identifier
	Name           string    `json:"name"`           // ETF name
	Symbol         string    `json:"symbol"`         // ETF symbol
	TotalShares    int       `json:"total_shares"`   // Total number of shares
	AvailableShares int      `json:"available_shares"` // Available shares for trading
	CurrentPrice   float64   `json:"current_price"`  // Current price of the ETF
	Timestamp      time.Time `json:"timestamp"`      // Last updated time
}

// ETFPricePegManager manages pegging the ETF price to external indices.
type ETFPricePegManager struct {
	pricePegData map[string]float64  // External price data for pegging
	mutex        sync.Mutex          // Mutex for thread-safe operations
}

// NewETFPricePegManager creates a new instance of ETFPricePegManager.
func NewETFPricePegManager() *ETFPricePegManager {
	return &ETFPricePegManager{
		pricePegData: make(map[string]float64),
	}
}

// AddPricePeg adds or updates the pegged price for an ETF.
func (ppm *ETFPricePegManager) AddPricePeg(etfID string, peggedPrice float64) {
	ppm.mutex.Lock()
	defer ppm.mutex.Unlock()

	ppm.pricePegData[etfID] = peggedPrice
}

// GetPeggedPrice retrieves the pegged price for an ETF.
func (ppm *ETFPricePegManager) GetPeggedPrice(etfID string) (float64, error) {
	ppm.mutex.Lock()
	defer ppm.mutex.Unlock()

	price, exists := ppm.pricePegData[etfID]
	if !exists {
		return 0, errors.New("pegged price not found for ETF")
	}

	return price, nil
}

// ETFPortfolioDetails represents the detailed information of an ETF portfolio.
type ETFPortfolioDetails struct {
	ETFID           string    `json:"etf_id"`           // ETF identifier
	Name            string    `json:"name"`             // ETF name
	TotalShares     int       `json:"total_shares"`     // Total shares in the portfolio
	AvailableShares int       `json:"available_shares"` // Available shares in the portfolio
	CurrentPrice    float64   `json:"current_price"`    // Current price of the portfolio
	Holdings        []Holding `json:"holdings"`         // Holdings of the ETF
	CreatedAt       time.Time `json:"created_at"`       // Creation time of the portfolio
	UpdatedAt       time.Time `json:"updated_at"`       // Last update time
}

// Holding represents a single holding in the ETF portfolio.
type Holding struct {
	AssetID   string  `json:"asset_id"`   // Asset identifier
	AssetName string  `json:"asset_name"` // Asset name
	Quantity  int     `json:"quantity"`   // Quantity of the asset
	Value     float64 `json:"value"`      // Current value of the asset
}

// ETFPortfolioService provides methods to manage ETF portfolio details.
type ETFPortfolioService struct {
	ledgerService     *ledger.Ledger        // Ledger service for logging portfolio updates
	encryptionService *encryption.Encryptor // Encryption service for securing portfolio details
	mutex             sync.Mutex            // Mutex for thread-safe operations
}

// NewETFPortfolioService creates a new ETFPortfolioService.
func NewETFPortfolioService(ledger *ledger.Ledger, encryptor *encryption.Encryptor) *ETFPortfolioService {
	return &ETFPortfolioService{
		ledgerService:     ledger,
		encryptionService: encryptor,
	}
}
// AddHolding adds a new holding to the ETF portfolio.
func (eps *ETFPortfolioService) AddHolding(etfID string, holding *Holding) error {
	eps.mutex.Lock()
	defer eps.mutex.Unlock()

	// Retrieve the portfolio to update.
	portfolio, err := eps.retrievePortfolio(etfID)
	if err != nil {
		return err
	}

	// Add the new holding to the portfolio.
	portfolio.Holdings = append(portfolio.Holdings, *holding)

	// Update the portfolio's total shares and value based on the new holding.
	portfolio.TotalShares += holding.Quantity
	portfolio.CurrentPrice += holding.Value

	// Encrypt the updated portfolio.
	encryptedPortfolio, err := eps.encryptionService.EncryptData(portfolio)
	if err != nil {
		return err
	}

	// Log the addition in the ledger.
	if err := eps.ledgerService.LogEvent("HoldingAdded", time.Now(), etfID); err != nil {
		return err
	}

	// Store the updated portfolio in the ledger.
	if err := eps.storePortfolio(etfID, encryptedPortfolio.(*ETFPortfolioDetails)); err != nil {
		return err
	}

	return nil
}

// RemoveHolding removes a holding from the ETF portfolio.
func (eps *ETFPortfolioService) RemoveHolding(etfID, assetID string) error {
	eps.mutex.Lock()
	defer eps.mutex.Unlock()

	// Retrieve the portfolio to update.
	portfolio, err := eps.retrievePortfolio(etfID)
	if err != nil {
		return err
	}

	// Find and remove the holding from the portfolio.
	found := false
	for i, holding := range portfolio.Holdings {
		if holding.AssetID == assetID {
			// Subtract the holding value and quantity from portfolio totals.
			portfolio.TotalShares -= holding.Quantity
			portfolio.CurrentPrice -= holding.Value

			// Remove the holding from the list.
			portfolio.Holdings = append(portfolio.Holdings[:i], portfolio.Holdings[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("holding not found in the portfolio")
	}

	// Encrypt the updated portfolio.
	encryptedPortfolio, err := eps.encryptionService.EncryptData(portfolio)
	if err != nil {
		return err
	}

	// Log the removal in the ledger.
	if err := eps.ledgerService.LogEvent("HoldingRemoved", time.Now(), etfID); err != nil {
		return err
	}

	// Store the updated portfolio in the ledger.
	if err := eps.storePortfolio(etfID, encryptedPortfolio.(*ETFPortfolioDetails)); err != nil {
		return err
	}

	return nil
}

// UpdatePortfolioValue updates the total value of the ETF portfolio.
func (eps *ETFPortfolioService) UpdatePortfolioValue(etfID string, newValue float64) error {
	eps.mutex.Lock()
	defer eps.mutex.Unlock()

	// Retrieve the portfolio to update.
	portfolio, err := eps.retrievePortfolio(etfID)
	if err != nil {
		return err
	}

	// Update the portfolio's current value.
	portfolio.CurrentPrice = newValue

	// Encrypt the updated portfolio.
	encryptedPortfolio, err := eps.encryptionService.EncryptData(portfolio)
	if err != nil {
		return err
	}

	// Log the update in the ledger.
	if err := eps.ledgerService.LogEvent("PortfolioValueUpdated", time.Now(), etfID); err != nil {
		return err
	}

	// Store the updated portfolio in the ledger.
	if err := eps.storePortfolio(etfID, encryptedPortfolio.(*ETFPortfolioDetails)); err != nil {
		return err
	}

	return nil
}

// GetPortfolioValue returns the current value of the portfolio.
func (eps *ETFPortfolioService) GetPortfolioValue(etfID string) (float64, error) {
	// Retrieve the latest value of the portfolio.
	portfolio, err := eps.retrievePortfolio(etfID)
	if err != nil {
		return 0, err
	}

	return portfolio.CurrentPrice, nil
}

// retrievePortfolio retrieves the portfolio details from ledger storage.
func (eps *ETFPortfolioService) retrievePortfolio(etfID string) (*ETFPortfolioDetails, error) {
	// Retrieve the portfolio data from the ledger storage.
	data, err := eps.ledgerService.RetrievePortfolio(etfID) 
	if err != nil {
		return nil, err
	}

	// Decrypt the portfolio before returning.
	decryptedPortfolio, err := eps.encryptionService.DecryptData(data)
	if err != nil {
		return nil, err
	}

	return decryptedPortfolio.(*ETFPortfolioDetails), nil
}

// storePortfolio stores the updated portfolio details in the ledger storage.
func (eps *ETFPortfolioService) storePortfolio(etfID string, portfolio *ETFPortfolioDetails) error {
	// Encrypt the portfolio for storage.
	encryptedData, err := eps.encryptionService.EncryptData(portfolio)
	if err != nil {
		return err
	}

	// Store the encrypted portfolio into the ledger's persistent storage.
	if err := eps.ledgerService.StorePortfolio(etfID, encryptedData); err != nil {
		return err
	}

	return nil
}
