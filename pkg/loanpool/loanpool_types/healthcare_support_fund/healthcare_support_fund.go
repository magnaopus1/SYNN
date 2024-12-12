package loanpool

import (
	"math/big"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewHealthcareSupportFund initializes the Healthcare Support Fund with the provided ledger, consensus engine, and encryption service.
func NewHealthcareSupportFund(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption) *HealthcareSupportFund {
	return &HealthcareSupportFund{
		TotalBalance:      big.NewInt(0),
		GrantsDistributed: big.NewInt(0),
		Ledger:            ledgerInstance,
		Consensus:         consensusEngine,
		EncryptionService: encryptionService,
	}
}

// ViewFundBalance returns the current balance available in the healthcare support fund.
func (fund *HealthcareSupportFund) ViewFundBalance() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.TotalBalance)
}

// ViewGrantsDistributed returns the total amount of healthcare grants that have been distributed.
func (fund *HealthcareSupportFund) ViewGrantsDistributed() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.GrantsDistributed)
}
