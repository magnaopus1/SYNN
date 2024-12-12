package loanpool

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// NewLoanPool initializes a new loan pool with zero balances.
func NewLoanPool(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption) *LoanPool {
	return &LoanPool{
		MainFund:              big.NewInt(0),
		PersonalGrantFund:     big.NewInt(0),
		EcosystemGrantFund:    big.NewInt(0),
		EducationFund:         big.NewInt(0),
		HealthcareSupportFund: big.NewInt(0),
		PovertyFund:           big.NewInt(0),
		SecuredFund:           big.NewInt(0),
		BusinessGrantFund:     big.NewInt(0),
		UnsecuredLoanFund:     big.NewInt(0),
		EnvironmentalFund:     big.NewInt(0),
		Ledger:                ledgerInstance,
		Consensus:             consensusEngine,
		Encryption:            encryptionService,
	}
}

// AddToMainFund adds funds to the main loan pool.
func (lp *LoanPool) AddToMainFund(amount *big.Int) error {
	lp.mutex.Lock()
	defer lp.mutex.Unlock()

	// Add the amount to the Main Fund.
	lp.MainFund.Add(lp.MainFund, amount)

	// Log the addition of funds in the ledger.
	err := lp.Ledger.RecordLoanPoolTransaction("MainFund", amount, "Added to Main Fund")
	if err != nil {
		return fmt.Errorf("failed to record transaction in ledger: %v", err)
	}

	fmt.Printf("Added %s to the Main Loan Fund. Total in Main Fund: %s\n", amount.String(), lp.MainFund.String())
	return nil
}

// DistributeFunds redistributes funds from the Main Fund to the other 9 funds based on their allocations.
func (lp *LoanPool) DistributeFunds() error {
	lp.mutex.Lock()
	defer lp.mutex.Unlock()

	// Ensure the Main Fund has a balance.
	if lp.MainFund.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("main fund has no available balance for distribution")
	}

	// Allocate the funds based on their percentages.
	totalFunds := new(big.Int).Set(lp.MainFund)
	lp.PersonalGrantFund.Add(lp.PersonalGrantFund, calculatePercentage(totalFunds, 25))
	lp.EcosystemGrantFund.Add(lp.EcosystemGrantFund, calculatePercentage(totalFunds, 25))
	lp.EducationFund.Add(lp.EducationFund, calculatePercentage(totalFunds, 5))
	lp.HealthcareSupportFund.Add(lp.HealthcareSupportFund, calculatePercentage(totalFunds, 5))
	lp.PovertyFund.Add(lp.PovertyFund, calculatePercentage(totalFunds, 5))
	lp.SecuredFund.Add(lp.SecuredFund, calculatePercentage(totalFunds, 15))
	lp.BusinessGrantFund.Add(lp.BusinessGrantFund, calculatePercentage(totalFunds, 25))
	lp.UnsecuredLoanFund.Add(lp.UnsecuredLoanFund, calculatePercentage(totalFunds, 15))
	lp.EnvironmentalFund.Add(lp.EnvironmentalFund, calculatePercentage(totalFunds, 5))

	// Empty the main fund after distribution.
	lp.MainFund.Set(big.NewInt(0))

	// Log the distribution in the ledger.
	err := lp.Ledger.RecordLoanPoolTransaction("DistributeFunds", totalFunds, "Funds Distributed to Loan Pools")
	if err != nil {
		return fmt.Errorf("failed to record distribution transaction in ledger: %v", err)
	}

	fmt.Println("Funds successfully distributed among the loan pools.")
	return nil
}

// calculatePercentage calculates the specific percentage of a total amount.
func calculatePercentage(total *big.Int, percent int) *big.Int {
	percentage := new(big.Int).Mul(total, big.NewInt(int64(percent)))
	return percentage.Div(percentage, big.NewInt(100))
}

// GetFundBalances retrieves the current balance of all loan pool funds.
func (lp *LoanPool) GetFundBalances() map[string]*big.Int {
	lp.mutex.Lock()
	defer lp.mutex.Unlock()

	return map[string]*big.Int{
		"MainFund":              lp.MainFund,
		"PersonalGrantFund":     lp.PersonalGrantFund,
		"EcosystemGrantFund":    lp.EcosystemGrantFund,
		"EducationFund":         lp.EducationFund,
		"HealthcareSupportFund": lp.HealthcareSupportFund,
		"PovertyFund":           lp.PovertyFund,
		"SecuredFund":           lp.SecuredFund,
		"BusinessGrantFund":     lp.BusinessGrantFund,
		"UnsecuredLoanFund":     lp.UnsecuredLoanFund,
		"EnvironmentalFund":     lp.EnvironmentalFund,
	}
}

// SecureFunds ensures that the funds are securely managed and inaccessible for any unauthorized operations.
func (lp *LoanPool) SecureFunds() error {
	lp.mutex.Lock()
	defer lp.mutex.Unlock()

	// Encrypt the entire fund pool data to protect it.
	encryptedData, err := lp.Encryption.EncryptData(fmt.Sprintf("%v", lp), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt loan pool data: %v", err)
	}

	// Store the encrypted data for security.
	err = lp.Ledger.StoreEncryptedLoanPoolData(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to store encrypted loan pool data: %v", err)
	}

	fmt.Println("Loan pool funds are securely encrypted and protected.")
	return nil
}

// RebalanceLoanPool runs automatically to rebalance the loan pool fund every 24 hours.
func (lp *LoanPool) RebalanceLoanPool() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		// Attempt to redistribute the funds every 24 hours.
		err := lp.DistributeFunds()
		if err != nil {
			fmt.Printf("Error redistributing funds: %v\n", err)
		} else {
			fmt.Println("Loan pool rebalanced successfully.")
		}
	}
}
