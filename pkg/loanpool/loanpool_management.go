package loanpool

import (
	"fmt"
	"math/big"

	"synnergy_network/pkg/ledger"
)


// NewLoanPoolManager initializes a new LoanPoolManager instance for viewing loan pool balances.
func NewLoanPoolManager(loanPool *LoanPool, ledgerInstance *ledger.Ledger) *LoanPoolManager {
	return &LoanPoolManager{
		LoanPool: loanPool,
		Ledger:   ledgerInstance,
	}
}

// ViewMainLoanPoolBalance returns the balance of the main loan pool.
func (lpm *LoanPoolManager) ViewMainLoanPoolBalance() *big.Int {
	lpm.mutex.Lock()
	defer lpm.mutex.Unlock()

	return lpm.LoanPool.MainFund
}

// ViewLoanPoolSubFundBalances returns the balances of all the sub-funds.
func (lpm *LoanPoolManager) ViewLoanPoolSubFundBalances() map[string]*big.Int {
	lpm.mutex.Lock()
	defer lpm.mutex.Unlock()

	// Return the current balances of each sub-fund
	return map[string]*big.Int{
		"PersonalGrantFund":    lpm.LoanPool.PersonalGrantFund,
		"EcosystemGrantFund":   lpm.LoanPool.EcosystemGrantFund,
		"EducationFund":        lpm.LoanPool.EducationFund,
		"HealthcareSupportFund": lpm.LoanPool.HealthcareSupportFund,
		"PovertyFund":          lpm.LoanPool.PovertyFund,
		"SecuredFund":          lpm.LoanPool.SecuredFund,
		"BusinessGrantFund":    lpm.LoanPool.BusinessGrantFund,
		"UnsecuredLoanFund":    lpm.LoanPool.UnsecuredLoanFund,
		"EnvironmentalFund":    lpm.LoanPool.EnvironmentalFund,
	}
}

// PrintLoanPoolBalances prints the balances of the main fund and all sub-funds to the console.
func (lpm *LoanPoolManager) PrintLoanPoolBalances() {
	lpm.mutex.Lock()
	defer lpm.mutex.Unlock()

	fmt.Printf("Main Loan Pool Balance: %s\n", lpm.LoanPool.MainFund.String())
	fmt.Printf("Personal Grant Fund Balance: %s\n", lpm.LoanPool.PersonalGrantFund.String())
	fmt.Printf("Ecosystem Grant Fund Balance: %s\n", lpm.LoanPool.EcosystemGrantFund.String())
	fmt.Printf("Education Fund Balance: %s\n", lpm.LoanPool.EducationFund.String())
	fmt.Printf("Healthcare Support Fund Balance: %s\n", lpm.LoanPool.HealthcareSupportFund.String())
	fmt.Printf("Poverty Fund Balance: %s\n", lpm.LoanPool.PovertyFund.String())
	fmt.Printf("Secured Fund Balance: %s\n", lpm.LoanPool.SecuredFund.String())
	fmt.Printf("Business Grant Fund Balance: %s\n", lpm.LoanPool.BusinessGrantFund.String())
	fmt.Printf("Unsecured Loan Fund Balance: %s\n", lpm.LoanPool.UnsecuredLoanFund.String())
	fmt.Printf("Environmental Fund Balance: %s\n", lpm.LoanPool.EnvironmentalFund.String())
}
