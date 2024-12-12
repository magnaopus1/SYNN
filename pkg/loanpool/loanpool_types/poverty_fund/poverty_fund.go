package loanpool

import (
	"math/big"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewPovertyFund initializes the Poverty Fund with the provided ledger, consensus engine, and encryption service.
func NewPovertyFund(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption) *PovertyFund {
	return &PovertyFund{
		TotalBalance:      big.NewInt(0),
		GrantsDistributed: big.NewInt(0),
		Ledger:            ledgerInstance,
		Consensus:         consensusEngine,
		EncryptionService: encryptionService,
	}
}

// ViewFundBalance returns the current balance available in the Poverty Fund.
func (fund *PovertyFund) ViewFundBalance() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.TotalBalance)
}

// ViewGrantsDistributed returns the total amount of grants that have been distributed from the Poverty Fund.
func (fund *PovertyFund) ViewGrantsDistributed() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.GrantsDistributed)
}
