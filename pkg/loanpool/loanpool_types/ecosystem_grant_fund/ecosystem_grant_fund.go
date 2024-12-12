package loanpool

import (
	"math/big"
	"synnergy_network/pkg/ledger"
)


// NewEcosystemGrantFund initializes the Ecosystem Grant Fund.
func NewEcosystemGrantFund(ledgerInstance *ledger.Ledger) *EcosystemGrantFund {
	return &EcosystemGrantFund{
		TotalBalance:      big.NewInt(0),
		GrantsDistributed: big.NewInt(0),
		Ledger:            ledgerInstance,
	}
}

// ViewFundBalance returns the current balance available in the fund.
func (fund *EcosystemGrantFund) ViewFundBalance() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.TotalBalance)
}

// ViewGrantsDistributed returns the total amount of grants that have been distributed.
func (fund *EcosystemGrantFund) ViewGrantsDistributed() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.GrantsDistributed)
}
