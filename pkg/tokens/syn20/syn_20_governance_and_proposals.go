package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// ENABLE_GOVERNANCE initializes governance features for the SYN20 token.
func (token *Syn20Token) ENABLE_GOVERNANCE() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.GovernanceEnabled = true
    return token.Ledger.RecordLog("GovernanceEnabled", "Governance enabled for SYN20 token")
}

// SUBMIT_PROPOSAL allows token holders to submit a proposal for governance.
func (token *Syn20Token) SUBMIT_PROPOSAL(proposalID, description string, proposer string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    proposal := ledger.GovernanceProposal{
        ProposalID:  proposalID,
        Title:       "New Proposal",
        Description: description,
        Creator:     proposer,
        CreatedAt:   time.Now(),
        Status:      "Pending",
    }
    token.Ledger.AddProposal(proposal)
    return token.Ledger.RecordLog("ProposalSubmitted", fmt.Sprintf("Proposal %s submitted by %s", proposalID, proposer))
}

// VOTE_ON_PROPOSAL allows token holders to vote on active proposals.
func (token *Syn20Token) VOTE_ON_PROPOSAL(proposalID string, voter string, voteFor bool) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.CastVote(proposalID, voter, voteFor)
    if err != nil {
        return fmt.Errorf("failed to cast vote: %v", err)
    }
    return token.Ledger.RecordLog("VoteCast", fmt.Sprintf("Vote cast on proposal %s by %s", proposalID, voter))
}

// DELEGATE_VOTING delegates voting power to another address.
func (token *Syn20Token) DELEGATE_VOTING(delegator, delegatee string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.DelegateVote(delegator, delegatee)
    if err != nil {
        return fmt.Errorf("failed to delegate voting: %v", err)
    }
    return token.Ledger.RecordLog("VoteDelegated", fmt.Sprintf("Vote delegated from %s to %s", delegator, delegatee))
}

// COLLECT_TAX collects a transaction tax from token transfers.
func (token *Syn20Token) COLLECT_TAX(account string, amount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    taxRate := token.Metadata.TaxRate
    taxAmount := new(big.Int).Div(new(big.Int).Mul(amount, big.NewInt(int64(taxRate*100))), big.NewInt(10000))
    token.BalanceSheet[account].Sub(token.BalanceSheet[account], taxAmount)

    return token.Ledger.RecordTransaction("TaxCollected", account, "TaxAccount", taxAmount)
}

// ENABLE_AUTO_TAX enables automatic tax collection on each transaction.
func (token *Syn20Token) ENABLE_AUTO_TAX() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.AutoTaxEnabled = true
    return token.Ledger.RecordLog("AutoTaxEnabled", "Automatic tax collection enabled")
}

// DISABLE_AUTO_TAX disables automatic tax collection on transactions.
func (token *Syn20Token) DISABLE_AUTO_TAX() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.AutoTaxEnabled = false
    return token.Ledger.RecordLog("AutoTaxDisabled", "Automatic tax collection disabled")
}

// SET_TRANSACTION_FEE sets a transaction fee for all SYN20 transfers.
func (token *Syn20Token) SET_TRANSACTION_FEE(feePercentage float64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.TransactionFee = feePercentage
    return token.Ledger.RecordLog("TransactionFeeSet", fmt.Sprintf("Transaction fee set to %.2f%%", feePercentage))
}

// CHECK_TRANSACTION_FEE retrieves the current transaction fee.
func (token *Syn20Token) CHECK_TRANSACTION_FEE() float64 {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Metadata.TransactionFee
}

// ADJUST_TRANSACTION_FEE adjusts the transaction fee to a new percentage.
func (token *Syn20Token) ADJUST_TRANSACTION_FEE(newFeePercentage float64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.TransactionFee = newFeePercentage
    return token.Ledger.RecordLog("TransactionFeeAdjusted", fmt.Sprintf("Transaction fee adjusted to %.2f%%", newFeePercentage))
}

// ENABLE_MULTI_SIG enables multi-signature requirements for high-value transactions.
func (token *Syn20Token) ENABLE_MULTI_SIG() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.MultiSigEnabled = true
    return token.Ledger.RecordLog("MultiSigEnabled", "Multi-signature requirement enabled")
}

// DISABLE_MULTI_SIG disables multi-signature requirements.
func (token *Syn20Token) DISABLE_MULTI_SIG() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.MultiSigEnabled = false
    return token.Ledger.RecordLog("MultiSigDisabled", "Multi-signature requirement disabled")
}

// SET_FEE_RECIPIENT sets the recipient for collected transaction fees.
func (token *Syn20Token) SET_FEE_RECIPIENT(recipient string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.FeeRecipient = recipient
    return token.Ledger.RecordLog("FeeRecipientSet", fmt.Sprintf("Fee recipient set to %s", recipient))
}

// FETCH_FEE_RECIPIENT retrieves the current fee recipient address.
func (token *Syn20Token) FETCH_FEE_RECIPIENT() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Metadata.FeeRecipient
}

// VALIDATE_MULTISIG validates if a transaction meets multi-signature requirements.
func (token *Syn20Token) VALIDATE_MULTISIG(transactionID string, signatures []string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if !token.Metadata.MultiSigEnabled {
        return true, nil // Multi-signature not required if disabled
    }

    valid, err := token.Ledger.ValidateMultiSig(transactionID, signatures)
    if err != nil {
        return false, fmt.Errorf("multi-signature validation failed: %v", err)
    }
    return valid, nil
}

// SET_FROZEN_ACCOUNT_LIST freezes or unfreezes accounts for regulatory purposes.
func (token *Syn20Token) SET_FROZEN_ACCOUNT_LIST(accounts []string, freeze bool) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for _, account := range accounts {
        token.Ledger.SetAccountStatus(account, freeze)
    }

    status := "frozen"
    if !freeze {
        status = "unfrozen"
    }
    return token.Ledger.RecordLog("AccountStatusChange", fmt.Sprintf("Accounts %v %s", accounts, status))
}
