package loanpool

import (
	"math/big"
	"synnergy_network/pkg/ledger"
)


// NewBusinessPersonalGrantFund initializes the Business Personal Grant Fund.
func NewBusinessPersonalGrantFund(ledgerInstance *ledger.Ledger) *BusinessPersonalGrantFund {
	return &BusinessPersonalGrantFund{
		TotalBalance:      big.NewInt(0),
		GrantsDistributed: big.NewInt(0),
		Ledger:            ledgerInstance,
	}
}

// ViewFundBalance returns the current balance available in the fund.
func (fund *BusinessPersonalGrantFund) ViewFundBalance() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.TotalBalance)
}

// ViewGrantsDistributed returns the total amount of grants that have been distributed.
func (fund *BusinessPersonalGrantFund) ViewGrantsDistributed() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.GrantsDistributed)
}
