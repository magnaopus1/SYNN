package loanpool

import (
	"math/big"

	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/common"
)

// NewSmallBusinessGrantFund initializes the Small Business Grant Fund.
func NewSmallBusinessGrantFund(ledgerInstance *ledger.Ledger) *common.SmallBusinessGrantFund {
	return &common.SmallBusinessGrantFund{
		TotalBalance:      big.NewInt(0),
		GrantsDistributed: big.NewInt(0),
		Ledger:            ledgerInstance,
	}
}

// ViewFundBalance returns the current balance available in the fund.
func (fund *common.SmallBusinessGrantFund) ViewFundBalance() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.TotalBalance)
}

// ViewGrantsDistributed returns the total amount of grants that have been distributed.
func (fund *common.SmallBusinessGrantFund) ViewGrantsDistributed() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.GrantsDistributed)
}
