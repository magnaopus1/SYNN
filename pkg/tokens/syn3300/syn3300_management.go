package syn3300

import (
	"sync"
	"time"

)

// FractionalShare represents a fractional share of an ETF (SYN3300) token
type FractionalShare struct {
	ETFID        string    `json:"etf_id"`
	ShareTokenID string    `json:"share_token_id"`
	Owner        string    `json:"owner"`
	Fraction     float64   `json:"fraction"` // Fraction of the total ETF token (e.g., 0.25 for 1/4)
	Timestamp    time.Time `json:"timestamp"`
}

// FractionalShareService provides methods to manage fractional shares of ETFs
type FractionalShareService struct {
	ledgerService     *ledger.Ledger
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewFractionalShareService creates a new instance of FractionalShareService
func NewFractionalShareService(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *FractionalShareService {
	return &FractionalShareService{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// AddFractionalShare adds a new fractional share to an ETF.
func (fs *FractionalShareService) AddFractionalShare(etfID, shareTokenID, owner string, fraction float64) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	share := &FractionalShare{
		ETFID:        etfID,
		ShareTokenID: shareTokenID,
		Owner:        owner,
		Fraction:     fraction,
		Timestamp:    time.Now(),
	}

	// Encrypt the share for secure storage
	encryptedShare, err := fs.encryptionService.EncryptData(share)
	if err != nil {
		return err
	}

	// Log the addition of the fractional share in the ledger
	if err := fs.ledgerService.LogEvent("FractionalShareAdded", time.Now(), shareTokenID); err != nil {
		return err
	}

	// Validate the fractional share using consensus
	if err := fs.consensusService.ValidateSubBlock(shareTokenID); err != nil {
		return err
	}

	// Store the encrypted share in the ledger storage
	if err := fs.storeFractionalShare(etfID, encryptedShare.(*FractionalShare)); err != nil {
		return err
	}

	return nil
}

// RetrieveFractionalShare retrieves a fractional share based on the shareTokenID.
func (fs *FractionalShareService) RetrieveFractionalShare(shareTokenID string) (*FractionalShare, error) {
	// Retrieve the encrypted share from the ledger
	data, err := fs.ledgerService.RetrieveShare(shareTokenID)
	if err != nil {
		return nil, err
	}

	// Decrypt the share before returning
	decryptedShare, err := fs.encryptionService.DecryptData(data)
	if err != nil {
		return nil, err
	}

	return decryptedShare.(*FractionalShare), nil
}

// storeFractionalShare is a helper function to store the fractional share securely in the ledger.
func (fs *FractionalShareService) storeFractionalShare(etfID string, share *FractionalShare) error {
	encryptedData, err := fs.encryptionService.EncryptData(share)
	if err != nil {
		return err
	}

	// Store the encrypted share in the ledger
	return fs.ledgerService.StoreShare(etfID, encryptedData)
}

// ShareTrackingRecord represents a record of fractional share tracking.
type ShareTrackingRecord struct {
	RecordID      string    `json:"record_id"`
	ETFID         string    `json:"etf_id"`
	ShareTokenID  string    `json:"share_token_id"`
	Owner         string    `json:"owner"`
	Fraction      float64   `json:"fraction"`
	LastUpdated   time.Time `json:"last_updated"`
	TransactionID string    `json:"transaction_id"`
}

// TrackShare updates and tracks fractional shares over time.
func (fs *FractionalShareService) TrackShare(shareTokenID string, fraction float64, transactionID string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	share, err := fs.RetrieveFractionalShare(shareTokenID)
	if err != nil {
		return err
	}

	record := &ShareTrackingRecord{
		RecordID:      shareTokenID,
		ETFID:         share.ETFID,
		ShareTokenID:  shareTokenID,
		Owner:         share.Owner,
		Fraction:      fraction,
		LastUpdated:   time.Now(),
		TransactionID: transactionID,
	}

	// Encrypt the tracking record for secure storage.
	encryptedRecord, err := fs.encryptionService.EncryptData(record)
	if err != nil {
		return err
	}

	// Log the update in the ledger.
	if err := fs.ledgerService.LogEvent("ShareTracked", time.Now(), shareTokenID); err != nil {
		return err
	}

	// Store the encrypted tracking record in the ledger.
	return fs.storeTrackingRecord(share.ETFID, encryptedRecord.(*ShareTrackingRecord))
}

// storeTrackingRecord stores a share tracking record securely.
func (fs *FractionalShareService) storeTrackingRecord(etfID string, record *ShareTrackingRecord) error {
	encryptedData, err := fs.encryptionService.EncryptData(record)
	if err != nil {
		return err
	}

	// Store the encrypted tracking record in the ledger.
	return fs.ledgerService.StoreTrackingRecord(etfID, encryptedData)
}

// ShareDistributionRecord represents a record of share distribution.
type ShareDistributionRecord struct {
	RecordID     string    `json:"record_id"`
	ETFID        string    `json:"etf_id"`
	ShareTokenID string    `json:"share_token_id"`
	Recipient    string    `json:"recipient"`
	Fraction     float64   `json:"fraction"`
	Timestamp    time.Time `json:"timestamp"`
}

// DistributeShare distributes fractional shares to recipients.
func (fs *FractionalShareService) DistributeShare(etfID, shareTokenID, recipient string, fraction float64) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	distribution := &ShareDistributionRecord{
		RecordID:     shareTokenID,
		ETFID:        etfID,
		ShareTokenID: shareTokenID,
		Recipient:    recipient,
		Fraction:     fraction,
		Timestamp:    time.Now(),
	}

	// Encrypt the distribution record.
	encryptedRecord, err := fs.encryptionService.EncryptData(distribution)
	if err != nil {
		return err
	}

	// Log the distribution in the ledger.
	if err := fs.ledgerService.LogEvent("ShareDistributed", time.Now(), shareTokenID); err != nil {
		return err
	}

	// Store the encrypted distribution record in the ledger.
	return fs.storeDistributionRecord(etfID, encryptedRecord.(*ShareDistributionRecord))
}

// storeDistributionRecord stores a share distribution record securely.
func (fs *FractionalShareService) storeDistributionRecord(etfID string, record *ShareDistributionRecord) error {
	encryptedData, err := fs.encryptionService.EncryptData(record)
	if err != nil {
		return err
	}

	// Store the encrypted distribution record in the ledger.
	return fs.ledgerService.StoreDistributionRecord(etfID, encryptedData)
}

// InvestmentOption represents an investment option for fractional ownership of an ETF.
type InvestmentOption struct {
	OptionID      string    `json:"option_id"`
	ETFID         string    `json:"etf_id"`
	Description   string    `json:"description"`
	MinInvestment float64   `json:"min_investment"`
	MaxInvestment float64   `json:"max_investment"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AddInvestmentOption adds a new investment option for fractional ownership.
func (fs *FractionalShareService) AddInvestmentOption(option *InvestmentOption) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// Encrypt the investment option data.
	encryptedOption, err := fs.encryptionService.EncryptData(option)
	if err != nil {
		return err
	}

	// Log the new investment option in the ledger.
	if err := fs.ledgerService.LogEvent("InvestmentOptionAdded", time.Now(), option.OptionID); err != nil {
		return err
	}

	// Store the encrypted investment option in the ledger.
	return fs.storeInvestmentOption(option.ETFID, encryptedOption.(*InvestmentOption))
}

// storeInvestmentOption securely stores an investment option.
func (fs *FractionalShareService) storeInvestmentOption(etfID string, option *InvestmentOption) error {
	encryptedData, err := fs.encryptionService.EncryptData(option)
	if err != nil {
		return err
	}

	// Store the encrypted investment option in the ledger.
	return fs.ledgerService.StoreInvestmentOption(etfID, encryptedData)
}
