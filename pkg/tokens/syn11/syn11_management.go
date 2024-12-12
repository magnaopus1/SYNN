package syn11

import (
	"errors"
	"fmt"
	"sync"
	"time"

  	"synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// InterestPaymentManager handles the payment of interest (coupons) to gilt holders.
type InterestPaymentManager struct {
	mutex           sync.Mutex
	Ledger          *ledger.Ledger                 // Ledger to track payments
	Consensus       *consensus.SynnergyConsensus   // Synnergy Consensus for validation
	Encryption      *encryption.EncryptionService  // Encryption service for securing payment details
	InterestHistory map[string][]InterestPayment   // Map of TokenID to interest payment history
}

// InterestPayment holds the details of an interest payment made to a holder.
type InterestPayment struct {
	TokenID     string    // Associated Token ID (Syn11)
	HolderID    string    // Holder receiving the interest payment
	Amount      uint64    // Amount of interest paid
	PaymentDate time.Time // Date when the interest was paid
	EncryptedData string  // Encrypted payment details for security
}

// NewInterestPaymentManager creates a new InterestPaymentManager.
func NewInterestPaymentManager(ledgerInstance *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *InterestPaymentManager {
	return &InterestPaymentManager{
		Ledger:          ledgerInstance,
		Consensus:       consensusEngine,
		Encryption:      encryptionService,
		InterestHistory: make(map[string][]InterestPayment),
	}
}

// PayInterest handles the payment of interest (coupons) for a specific Syn11 token.
func (ipm *InterestPaymentManager) PayInterest(tokenID string, holderID string, amount uint64) error {
	ipm.mutex.Lock()
	defer ipm.mutex.Unlock()

	// Validate the interest payment through Synnergy Consensus
	if err := ipm.Consensus.ValidateInterestPayment(tokenID, holderID, amount); err != nil {
		return fmt.Errorf("interest payment validation failed: %w", err)
	}

	// Encrypt the interest payment details
	encryptedData, err := ipm.Encryption.Encrypt([]byte(fmt.Sprintf("Payment of %d to %s", amount, holderID)))
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	// Record the interest payment in the ledger
	interestPayment := InterestPayment{
		TokenID:      tokenID,
		HolderID:     holderID,
		Amount:       amount,
		PaymentDate:  time.Now(),
		EncryptedData: string(encryptedData),
	}
	if err := ipm.Ledger.RecordInterestPayment(tokenID, interestPayment); err != nil {
		return fmt.Errorf("ledger update failed: %w", err)
	}

	// Append the interest payment to the history
	ipm.InterestHistory[tokenID] = append(ipm.InterestHistory[tokenID], interestPayment)

	return nil
}

// GetInterestPaymentHistory returns the interest payment history for a Syn11 token.
func (ipm *InterestPaymentManager) GetInterestPaymentHistory(tokenID string) ([]InterestPayment, error) {
	history, exists := ipm.InterestHistory[tokenID]
	if !exists {
		return nil, errors.New("no interest payment history found for the token")
	}
	return history, nil
}

// CouponManager handles coupon-related logic (interest coupons for gilt holders).
type CouponManager struct {
	mutex         sync.Mutex
	Ledger        *ledger.Ledger
	Consensus     *consensus.SynnergyConsensus
	Encryption    *encryption.EncryptionService
	CouponHistory map[string][]CouponRecord // Map of TokenID to coupon history
}

// CouponRecord represents a coupon issued for a Syn11 gilt.
type CouponRecord struct {
	TokenID     string    // Associated Token ID
	CouponRate  float64   // Coupon rate (interest rate)
	IssuedDate  time.Time // Date when the coupon was issued
	EncryptedData string  // Encrypted coupon details
}

// NewCouponManager creates a new CouponManager.
func NewCouponManager(ledgerInstance *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *CouponManager {
	return &CouponManager{
		Ledger:        ledgerInstance,
		Consensus:     consensusEngine,
		Encryption:    encryptionService,
		CouponHistory: make(map[string][]CouponRecord),
	}
}

// IssueCoupon handles the issuance of a coupon for a specific Syn11 gilt.
func (cm *CouponManager) IssueCoupon(tokenID string, couponRate float64) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Validate the coupon issuance through Synnergy Consensus
	if err := cm.Consensus.ValidateCouponIssuance(tokenID, couponRate); err != nil {
		return fmt.Errorf("coupon validation failed: %w", err)
	}

	// Encrypt the coupon details
	encryptedData, err := cm.Encryption.Encrypt([]byte(fmt.Sprintf("Coupon of %.2f issued for %s", couponRate, tokenID)))
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	// Record the coupon issuance in the ledger
	coupon := CouponRecord{
		TokenID:      tokenID,
		CouponRate:   couponRate,
		IssuedDate:   time.Now(),
		EncryptedData: string(encryptedData),
	}
	if err := cm.Ledger.RecordCoupon(tokenID, coupon); err != nil {
		return fmt.Errorf("ledger update failed: %w", err)
	}

	// Append the coupon record to the history
	cm.CouponHistory[tokenID] = append(cm.CouponHistory[tokenID], coupon)

	return nil
}

// GetCouponHistory retrieves the coupon issuance history for a Syn11 gilt.
func (cm *CouponManager) GetCouponHistory(tokenID string) ([]CouponRecord, error) {
	history, exists := cm.CouponHistory[tokenID]
	if !exists {
		return nil, errors.New("no coupon history found for the token")
	}
	return history, nil
}

// RedemptionManager handles the redemption of Syn11 tokens for fiat or other assets.
type RedemptionManager struct {
	mutex           sync.Mutex
	Ledger          *ledger.Ledger
	Consensus       *consensus.SynnergyConsensus
	Encryption      *encryption.EncryptionService
	RedemptionLog   map[string][]RedemptionRecord // Map of TokenID to redemption records
}

// RedemptionRecord holds the details of a redemption.
type RedemptionRecord struct {
	TokenID      string    // Associated Token ID
	HolderID     string    // Holder redeeming the tokens
	Amount       uint64    // Amount redeemed
	RedemptionDate time.Time // Date of the redemption
	EncryptedData string   // Encrypted redemption details
}

// NewRedemptionManager creates a new RedemptionManager.
func NewRedemptionManager(ledgerInstance *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *RedemptionManager {
	return &RedemptionManager{
		Ledger:        ledgerInstance,
		Consensus:     consensusEngine,
		Encryption:    encryptionService,
		RedemptionLog: make(map[string][]RedemptionRecord),
	}
}

// RedeemToken handles the redemption of a Syn11 token.
func (rm *RedemptionManager) RedeemToken(tokenID, holderID string, amount uint64) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Validate the redemption through Synnergy Consensus
	if err := rm.Consensus.ValidateRedemption(tokenID, holderID, amount); err != nil {
		return fmt.Errorf("redemption validation failed: %w", err)
	}

	// Encrypt the redemption details
	encryptedData, err := rm.Encryption.Encrypt([]byte(fmt.Sprintf("Redemption of %d tokens by %s", amount, holderID)))
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	// Record the redemption in the ledger
	redemption := RedemptionRecord{
		TokenID:      tokenID,
		HolderID:     holderID,
		Amount:       amount,
		RedemptionDate: time.Now(),
		EncryptedData: string(encryptedData),
	}
	if err := rm.Ledger.RecordRedemption(tokenID, redemption); err != nil {
		return fmt.Errorf("ledger update failed: %w", err)
	}

	// Append the redemption record to the log
	rm.RedemptionLog[tokenID] = append(rm.RedemptionLog[tokenID], redemption)

	return nil
}

// GetRedemptionLog retrieves the redemption log for a specific Syn11 token.
func (rm *RedemptionManager) GetRedemptionLog(tokenID string) ([]RedemptionRecord, error) {
	log, exists := rm.RedemptionLog[tokenID]
	if !exists {
		return nil, errors.New("no redemption log found for the token")
	}
	return log, nil
}


