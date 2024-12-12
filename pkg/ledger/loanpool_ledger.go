package ledger

import (
	"fmt"
	"math/big"
)

// RecordLoanPoolTransaction records a transaction in a loan or grant pool.
func (ledger *Ledger) RecordLoanPoolTransaction(poolID string, tx TransactionRecord) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Helper function to convert float64 to *big.Int (assume 18 decimal places)
    amountBigInt := new(big.Int)
    amountBigInt.SetString(fmt.Sprintf("%.0f", tx.Amount*1e18), 10) // Convert float64 to string and then to *big.Int

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        pool.LoanRecords[tx.ID].Transactions = append(pool.LoanRecords[tx.ID].Transactions, tx)
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Transaction recorded in secured loan pool %s: %s\n", poolID, tx.ID)
        return nil
    }

    // Unsecured loan pool
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        pool.LoanRecords[tx.ID].Transactions = append(pool.LoanRecords[tx.ID].Transactions, tx)
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Transaction recorded in unsecured loan pool %s: %s\n", poolID, tx.ID)
        return nil
    }

    // Business grant pool
    if pool, exists := ledger.businessGrantPools[poolID]; exists {
        pool.GrantsDistributed = pool.GrantsDistributed.Add(pool.GrantsDistributed, amountBigInt)
        ledger.businessGrantPools[poolID] = pool
        fmt.Printf("Transaction recorded in business grant pool %s: %s\n", poolID, tx.ID)
        return nil
    }

    // Ecosystem grant pool
    if pool, exists := ledger.ecosystemGrantPools[poolID]; exists {
        pool.GrantsDistributed = pool.GrantsDistributed.Add(pool.GrantsDistributed, amountBigInt)
        ledger.ecosystemGrantPools[poolID] = pool
        fmt.Printf("Transaction recorded in ecosystem grant pool %s: %s\n", poolID, tx.ID)
        return nil
    }

    // Education fund pool
    if pool, exists := ledger.educationFundPools[poolID]; exists {
        pool.GrantsDistributed = pool.GrantsDistributed.Add(pool.GrantsDistributed, amountBigInt)
        ledger.educationFundPools[poolID] = pool
        fmt.Printf("Transaction recorded in education fund pool %s: %s\n", poolID, tx.ID)
        return nil
    }

    // Healthcare support fund pool
    if pool, exists := ledger.healthcareSupportPools[poolID]; exists {
        pool.GrantsDistributed = pool.GrantsDistributed.Add(pool.GrantsDistributed, amountBigInt)
        ledger.healthcareSupportPools[poolID] = pool
        fmt.Printf("Transaction recorded in healthcare support fund pool %s: %s\n", poolID, tx.ID)
        return nil
    }

    // Poverty relief fund pool
    if pool, exists := ledger.povertyReliefPools[poolID]; exists {
        pool.GrantsDistributed = pool.GrantsDistributed.Add(pool.GrantsDistributed, amountBigInt)
        ledger.povertyReliefPools[poolID] = pool
        fmt.Printf("Transaction recorded in poverty relief fund pool %s: %s\n", poolID, tx.ID)
        return nil
    }

    return fmt.Errorf("loan or grant pool with ID %s does not exist", poolID)
}

// StoreEncryptedLoanPoolData stores encrypted data for loan or grant pools.
func (ledger *Ledger) StoreEncryptedLoanPoolData(poolID, encryptedData string) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Check all types of pools for encrypted data storage
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        pool.EncryptedData = encryptedData
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Encrypted data stored for secured loan pool %s.\n", poolID)
        return nil
    }

    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        pool.EncryptedData = encryptedData
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Encrypted data stored for unsecured loan pool %s.\n", poolID)
        return nil
    }

    if pool, exists := ledger.businessGrantPools[poolID]; exists {
        pool.EncryptedData = encryptedData
        ledger.businessGrantPools[poolID] = pool
        fmt.Printf("Encrypted data stored for business grant pool %s.\n", poolID)
        return nil
    }

    if pool, exists := ledger.ecosystemGrantPools[poolID]; exists {
        pool.EncryptedData = encryptedData
        ledger.ecosystemGrantPools[poolID] = pool
        fmt.Printf("Encrypted data stored for ecosystem grant pool %s.\n", poolID)
        return nil
    }

    if pool, exists := ledger.educationFundPools[poolID]; exists {
        pool.EncryptedData = encryptedData
        ledger.educationFundPools[poolID] = pool
        fmt.Printf("Encrypted data stored for education fund pool %s.\n", poolID)
        return nil
    }

    if pool, exists := ledger.healthcareSupportPools[poolID]; exists {
        pool.EncryptedData = encryptedData
        ledger.healthcareSupportPools[poolID] = pool
        fmt.Printf("Encrypted data stored for healthcare support fund pool %s.\n", poolID)
        return nil
    }

    if pool, exists := ledger.povertyReliefPools[poolID]; exists {
        pool.EncryptedData = encryptedData
        ledger.povertyReliefPools[poolID] = pool
        fmt.Printf("Encrypted data stored for poverty relief fund pool %s.\n", poolID)
        return nil
    }

    return fmt.Errorf("loan or grant pool with ID %s does not exist", poolID)
}

// RecordBusinessPersonalGrantProposal records a new business personal grant proposal in the ledger.
func (ledger *Ledger) RecordBusinessPersonalGrantProposal(proposal BusinessPersonalGrantProposal) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Check if the proposal already exists
	if _, exists := ledger.State.BusinessGrantProposals[proposal.RegistrationNumber]; exists {
		return fmt.Errorf("proposal for business %s already exists", proposal.BusinessName)
	}

	// Record the proposal
	ledger.State.BusinessGrantProposals[proposal.RegistrationNumber] = &proposal
	fmt.Printf("Business personal grant proposal for %s recorded successfully.\n", proposal.BusinessName)
	return nil
}

// RecordEcosystemGrantProposal records a new ecosystem grant proposal in the ledger.
func (ledger *Ledger) RecordEcosystemGrantProposal(proposal EcosystemGrantProposal) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Check if the proposal already exists
	if _, exists := ledger.State.EcosystemGrantProposals[proposal.RegistrationNumber]; exists {
		return fmt.Errorf("proposal for business %s already exists", proposal.BusinessName)
	}

	// Record the proposal
	ledger.State.EcosystemGrantProposals[proposal.RegistrationNumber] = &proposal
	fmt.Printf("Ecosystem grant proposal for %s recorded successfully.\n", proposal.BusinessName)
	return nil
}

// RecordEducationFundProposal records a new education fund proposal in the ledger.
func (ledger *Ledger) RecordEducationFundProposal(proposal EducationFundProposal) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Check if the proposal already exists
	if _, exists := ledger.State.EducationFundProposals[proposal.WalletAddress]; exists {
		return fmt.Errorf("proposal for applicant %s already exists", proposal.ApplicantName)
	}

	// Record the proposal
	ledger.State.EducationFundProposals[proposal.WalletAddress] = &proposal
	fmt.Printf("Education fund proposal for %s recorded successfully.\n", proposal.ApplicantName)
	return nil
}


// RecordHealthcareSupportFundProposal records a new healthcare support fund proposal in the ledger.
func (ledger *Ledger) RecordHealthcareSupportFundProposal(proposal HealthcareSupportFundProposal) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Check if the proposal already exists
	if _, exists := ledger.State.HealthcareSupportProposals[proposal.WalletAddress]; exists {
		return fmt.Errorf("proposal for applicant %s already exists", proposal.ApplicantName)
	}

	// Record the proposal
	ledger.State.HealthcareSupportProposals[proposal.WalletAddress] = &proposal
	fmt.Printf("Healthcare support fund proposal for %s recorded successfully.\n", proposal.ApplicantName)
	return nil
}


// RecordPovertyFundProposal records a new poverty fund proposal in the ledger.
func (ledger *Ledger) RecordPovertyFundProposal(proposal PovertyFundProposal) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Check if the proposal already exists
	if _, exists := ledger.State.PovertyFundProposals[proposal.WalletAddress]; exists {
		return fmt.Errorf("proposal for applicant %s already exists", proposal.ApplicantName)
	}

	// Record the proposal
	ledger.State.PovertyFundProposals[proposal.WalletAddress] = &proposal
	fmt.Printf("Poverty fund proposal for %s recorded successfully.\n", proposal.ApplicantName)
	return nil
}


// RecordSecuredLoanProposal records a new secured loan proposal in the ledger.
func (ledger *Ledger) RecordSecuredLoanProposal(proposal SecuredLoanProposal) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Check if the proposal already exists
	if _, exists := ledger.State.SecuredLoanProposals[proposal.LoanID]; exists {
		return fmt.Errorf("loan proposal with ID %s already exists", proposal.LoanID)
	}

	// Record the proposal
	ledger.State.SecuredLoanProposals[proposal.LoanID] = &proposal
	fmt.Printf("Secured loan proposal for loan ID %s recorded successfully.\n", proposal.LoanID)
	return nil
}



// RecordUnsecuredLoanProposal records a new unsecured loan proposal in the ledger.
func (ledger *Ledger) RecordUnsecuredLoanProposal(proposal UnsecuredLoanProposal) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Check if the proposal already exists
	if _, exists := ledger.State.UnsecuredLoanProposals[proposal.LoanID]; exists {
		return fmt.Errorf("loan proposal with ID %s already exists", proposal.LoanID)
	}

	// Record the proposal
	ledger.State.UnsecuredLoanProposals[proposal.LoanID] = &proposal
	fmt.Printf("Unsecured loan proposal for loan ID %s recorded successfully.\n", proposal.LoanID)
	return nil
}


// RecordSmallBusinessGrantProposal records a new small business grant proposal in the ledger.
func (ledger *Ledger) RecordSmallBusinessGrantProposal(proposal SmallBusinessGrantProposal) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Check if the proposal already exists
	if _, exists := ledger.State.SmallBusinessGrantProposals[proposal.RegistrationNumber]; exists {
		return fmt.Errorf("proposal for business %s already exists", proposal.BusinessName)
	}

	// Record the proposal
	ledger.State.SmallBusinessGrantProposals[proposal.RegistrationNumber] = &proposal
	fmt.Printf("Small business grant proposal for %s recorded successfully.\n", proposal.BusinessName)
	return nil
}

// RecordProposalApproval records the approval of a loan or grant proposal.
func (ledger *Ledger) RecordProposalApproval(poolID, proposalID string) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan proposal approval
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        proposal.Status = "approved"
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Proposal %s approved in secured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    // Unsecured loan proposal approval
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in unsecured loan pool %s", proposalID, poolID)
        }
        proposal.Status = "approved"
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Proposal %s approved in unsecured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    // Business grant proposal approval
    if pool, exists := ledger.businessGrantPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in business grant pool %s", proposalID, poolID)
        }
        proposal.Status = "approved"
        ledger.businessGrantPools[poolID] = pool
        fmt.Printf("Proposal %s approved in business grant pool %s.\n", proposalID, poolID)
        return nil
    }

    // Ecosystem grant proposal approval
    if pool, exists := ledger.ecosystemGrantPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in ecosystem grant pool %s", proposalID, poolID)
        }
        proposal.Status = "approved"
        ledger.ecosystemGrantPools[poolID] = pool
        fmt.Printf("Proposal %s approved in ecosystem grant pool %s.\n", proposalID, poolID)
        return nil
    }

    // Education fund proposal approval
    if pool, exists := ledger.educationFundPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in education fund pool %s", proposalID, poolID)
        }
        proposal.Status = "approved"
        ledger.educationFundPools[poolID] = pool
        fmt.Printf("Proposal %s approved in education fund pool %s.\n", proposalID, poolID)
        return nil
    }

    // Healthcare support fund proposal approval
    if pool, exists := ledger.healthcareSupportPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in healthcare support fund pool %s", proposalID, poolID)
        }
        proposal.Status = "approved"
        ledger.healthcareSupportPools[poolID] = pool
        fmt.Printf("Proposal %s approved in healthcare support fund pool %s.\n", proposalID, poolID)
        return nil
    }

    // Poverty relief fund proposal approval
    if pool, exists := ledger.povertyReliefPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in poverty relief fund pool %s", proposalID, poolID)
        }
        proposal.Status = "approved"
        ledger.povertyReliefPools[poolID] = pool
        fmt.Printf("Proposal %s approved in poverty relief fund pool %s.\n", proposalID, poolID)
        return nil
    }

    return fmt.Errorf("loan or grant pool with ID %s does not exist", poolID)
}

// RecordDisbursement records the disbursement of funds for a loan or grant.
func (ledger *Ledger) RecordDisbursement(poolID, proposalID string, amount float64) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        if pool.FundsAvailable < amount {
            return fmt.Errorf("insufficient funds in secured loan pool %s", poolID)
        }
        proposal.FundsDisbursed += amount
        pool.FundsAvailable -= amount
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Disbursed %f from secured loan pool %s for proposal %s.\n", amount, poolID, proposalID)
        return nil
    }

    // Unsecured loan pool
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in unsecured loan pool %s", proposalID, poolID)
        }
        if pool.FundsAvailable < amount {
            return fmt.Errorf("insufficient funds in unsecured loan pool %s", poolID)
        }
        proposal.FundsDisbursed += amount
        pool.FundsAvailable -= amount
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Disbursed %f from unsecured loan pool %s for proposal %s.\n", amount, poolID, proposalID)
        return nil
    }

    // Business grant pool
    if pool, exists := ledger.businessGrantPools[poolID]; exists {
        if pool.TotalBalance.Cmp(big.NewInt(int64(amount))) < 0 {
            return fmt.Errorf("insufficient funds in business grant pool %s", poolID)
        }
        pool.GrantsDistributed = pool.GrantsDistributed.Add(pool.GrantsDistributed, big.NewInt(int64(amount)))
        ledger.businessGrantPools[poolID] = pool
        fmt.Printf("Disbursed %f from business grant pool %s for proposal %s.\n", amount, poolID, proposalID)
        return nil
    }

    // Ecosystem grant pool
    if pool, exists := ledger.ecosystemGrantPools[poolID]; exists {
        if pool.TotalBalance.Cmp(big.NewInt(int64(amount))) < 0 {
            return fmt.Errorf("insufficient funds in ecosystem grant pool %s", poolID)
        }
        pool.GrantsDistributed = pool.GrantsDistributed.Add(pool.GrantsDistributed, big.NewInt(int64(amount)))
        ledger.ecosystemGrantPools[poolID] = pool
        fmt.Printf("Disbursed %f from ecosystem grant pool %s for proposal %s.\n", amount, poolID, proposalID)
        return nil
    }

    // Education fund pool
    if pool, exists := ledger.educationFundPools[poolID]; exists {
        if pool.TotalBalance.Cmp(big.NewInt(int64(amount))) < 0 {
            return fmt.Errorf("insufficient funds in education fund pool %s", poolID)
        }
        pool.GrantsDistributed = pool.GrantsDistributed.Add(pool.GrantsDistributed, big.NewInt(int64(amount)))
        ledger.educationFundPools[poolID] = pool
        fmt.Printf("Disbursed %f from education fund pool %s for proposal %s.\n", amount, poolID, proposalID)
        return nil
    }

    // Healthcare support fund pool
    if pool, exists := ledger.healthcareSupportPools[poolID]; exists {
        if pool.TotalBalance.Cmp(big.NewInt(int64(amount))) < 0 {
            return fmt.Errorf("insufficient funds in healthcare support fund pool %s", poolID)
        }
        pool.GrantsDistributed = pool.GrantsDistributed.Add(pool.GrantsDistributed, big.NewInt(int64(amount)))
        ledger.healthcareSupportPools[poolID] = pool
        fmt.Printf("Disbursed %f from healthcare support fund pool %s for proposal %s.\n", amount, poolID, proposalID)
        return nil
    }

    // Poverty relief fund pool
    if pool, exists := ledger.povertyReliefPools[poolID]; exists {
        if pool.TotalBalance.Cmp(big.NewInt(int64(amount))) < 0 {
            return fmt.Errorf("insufficient funds in poverty relief fund pool %s", poolID)
        }
        pool.GrantsDistributed = pool.GrantsDistributed.Add(pool.GrantsDistributed, big.NewInt(int64(amount)))
        ledger.povertyReliefPools[poolID] = pool
        fmt.Printf("Disbursed %f from poverty relief fund pool %s for proposal %s.\n", amount, poolID, proposalID)
        return nil
    }

    return fmt.Errorf("loan or grant pool with ID %s does not exist", poolID)
}

// RecordCollateralSubmission records a collateral submission for a loan proposal.
func (ledger *Ledger) RecordCollateralSubmission(poolID, proposalID string, collateral Collateral) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        proposal.Collateral = append(proposal.Collateral, collateral)
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Collateral submitted for proposal %s in secured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    return fmt.Errorf("loan pool with ID %s does not accept collateral submissions", poolID)
}

// RecordRepayment records a repayment transaction for a loan.
func (ledger *Ledger) RecordRepayment(poolID, proposalID string, repayment Repayment) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        proposal.Repayments = append(proposal.Repayments, repayment)
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Repayment recorded for proposal %s in secured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    // Unsecured loan pool
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in unsecured loan pool %s", proposalID, poolID)
        }
        proposal.Repayments = append(proposal.Repayments, repayment)
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Repayment recorded for proposal %s in unsecured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    return fmt.Errorf("loan pool with ID %s does not exist", poolID)
}

// RecordLoanTerms records the loan terms for a specific loan proposal.
func (ledger *Ledger) RecordLoanTerms(poolID, proposalID string, loanTerms LoanTerms) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        proposal.LoanTerms = loanTerms
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Loan terms recorded for proposal %s in secured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    // Unsecured loan pool
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in unsecured loan pool %s", proposalID, poolID)
        }
        proposal.LoanTerms = loanTerms
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Loan terms recorded for proposal %s in unsecured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    return fmt.Errorf("loan pool with ID %s does not exist", poolID)
}

// UpdateLoanTerms updates the loan terms for a specific loan proposal.
func (ledger *Ledger) UpdateLoanTerms(poolID, proposalID string, updatedTerms LoanTerms) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        proposal.LoanTerms = updatedTerms
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Loan terms updated for proposal %s in secured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    // Unsecured loan pool
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in unsecured loan pool %s", proposalID, poolID)
        }
        proposal.LoanTerms = updatedTerms
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Loan terms updated for proposal %s in unsecured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    return fmt.Errorf("loan pool with ID %s does not exist", poolID)
}

// RecordLoanDefault records a loan default for a specific loan proposal.
func (ledger *Ledger) RecordLoanDefault(poolID, proposalID string) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        proposal.Status = "defaulted"
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Loan default recorded for proposal %s in secured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    // Unsecured loan pool
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in unsecured loan pool %s", proposalID, poolID)
        }
        proposal.Status = "defaulted"
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Loan default recorded for proposal %s in unsecured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    return fmt.Errorf("loan pool with ID %s does not exist", poolID)
}

// RecordInterestPayment records an interest payment for a specific loan.
func (ledger *Ledger) RecordInterestPayment(poolID, proposalID string, interestPayment InterestPayment) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        proposal.InterestPayments = append(proposal.InterestPayments, interestPayment)
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Interest payment recorded for proposal %s in secured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    // Unsecured loan pool
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in unsecured loan pool %s", proposalID, poolID)
        }
        proposal.InterestPayments = append(proposal.InterestPayments, interestPayment)
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Interest payment recorded for proposal %s in unsecured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    return fmt.Errorf("loan pool with ID %s does not exist", poolID)
}

// FinalizeProposal finalizes a proposal in the loan or grant pool.
func (ledger *Ledger) FinalizeProposal(poolID, proposalID string) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        proposal.Status = "finalized"
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Proposal %s finalized in secured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    // Unsecured loan pool
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in unsecured loan pool %s", proposalID, poolID)
        }
        proposal.Status = "finalized"
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Proposal %s finalized in unsecured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    return fmt.Errorf("loan pool with ID %s does not exist", poolID)
}

// UpdateAffordabilityApproval records the approval of an affordability check for a loan proposal.
func (ledger *Ledger) UpdateAffordabilityApproval(poolID, proposalID string) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        proposal.AffordabilityCheck.Status = "approved"
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Affordability check approved for proposal %s in secured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    // Unsecured loan pool
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in unsecured loan pool %s", proposalID, poolID)
        }
        proposal.AffordabilityCheck.Status = "approved"
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Affordability check approved for proposal %s in unsecured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    return fmt.Errorf("loan pool with ID %s does not exist", poolID)
}

// UpdateAffordabilityRejection records the rejection of an affordability check for a loan proposal.
func (ledger *Ledger) UpdateAffordabilityRejection(poolID, proposalID string) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Secured loan pool
    if pool, exists := ledger.securedLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in secured loan pool %s", proposalID, poolID)
        }
        proposal.AffordabilityCheck.Status = "rejected"
        ledger.securedLoanPools[poolID] = pool
        fmt.Printf("Affordability check rejected for proposal %s in secured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    // Unsecured loan pool
    if pool, exists := ledger.unsecuredLoanPools[poolID]; exists {
        proposal, ok := pool.Proposals[proposalID]
        if !ok {
            return fmt.Errorf("proposal with ID %s does not exist in unsecured loan pool %s", proposalID, poolID)
        }
        proposal.AffordabilityCheck.Status = "rejected"
        ledger.unsecuredLoanPools[poolID] = pool
        fmt.Printf("Affordability check rejected for proposal %s in unsecured loan pool %s.\n", proposalID, poolID)
        return nil
    }

    return fmt.Errorf("loan pool with ID %s does not exist", poolID)
}

// RecordCharityFeeDistribution records the distribution of fees to the charity pool.
func (l *Ledger) RecordCharityFeeDistribution(fee float64) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Update the charity pool with the fee distribution
	l.charityPool.totalBalance += fee
	l.charityPool.LedgerInstance.State.TransactionHistory = append(l.charityPool.LedgerInstance.State.TransactionHistory, TransactionRecord{
		Amount:  fee,
		From:    "system_fee_pool",
		To:      "charity_pool",
		Hash:    generateHash("system_fee_pool" + "charity_pool" + strconv.FormatFloat(fee, 'f', 6, 64)),
		Status:  "confirmed",
	})

	return nil
}

// RecordCharityPoolWithdrawal records a withdrawal from the charity pool.
func (l *Ledger) RecordCharityPoolWithdrawal(amount float64) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Ensure the charity pool has enough funds for the withdrawal
	if l.charityPool.totalBalance < amount {
		return fmt.Errorf("insufficient funds in charity pool")
	}

	// Update the pool by subtracting the withdrawal amount
	l.charityPool.totalBalance -= amount

	// Record the transaction in the ledger state
	l.State.TransactionHistory = append(l.State.TransactionHistory, TransactionRecord{
		Amount:  amount,
		From:    "charity_pool",
		To:      "withdrawal",
		Hash:    generateHash("charity_pool" + "withdrawal" + strconv.FormatFloat(amount, 'f', 6, 64)),
		Status:  "confirmed",
	})

	return nil
}

// DistributeFunds distributes funds from the system to the charity pool.
func (l *Ledger) DistributeFunds(amount float64) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Update the charity pool
	l.charityPool.totalBalance += amount

	// Record the transaction in the ledger state
	l.State.TransactionHistory = append(l.State.TransactionHistory, TransactionRecord{
		Amount:  amount,
		From:    "system_distributor",
		To:      "charity_pool",
		Hash:    generateHash("system_distributor" + "charity_pool" + strconv.FormatFloat(amount, 'f', 6, 64)),
		Status:  "confirmed",
	})

	return nil
}

// RecordCharityProposal records a proposal for a new charity into the ledger.
func (l *Ledger) RecordCharityProposal(charityID string, proposalDetails string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Record the charity proposal in the ledger
	l.State.TransactionHistory = append(l.State.TransactionHistory, TransactionRecord{
		From:    charityID,
		To:      "proposal",
		Hash:    generateHash(charityID + "proposal" + proposalDetails),
		Status:  "pending",
	})

	return nil
}

// RecordCharityRemoval clears the charity pool in the ledger.
func (l *Ledger) RecordCharityRemoval() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Clear the charity pool
	l.charityPool.totalBalance = 0

	// Record the removal in the ledger state
	l.State.TransactionHistory = append(l.State.TransactionHistory, TransactionRecord{
		From:    "charity_pool",
		To:      "system_removal",
		Hash:    generateHash("charity_pool" + "system_removal"),
		Status:  "removed",
	})

	return nil
}

// RecordTopCharities updates the record of the top charities based on total funds received.
func (l *Ledger) RecordTopCharities(topCharities []string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Record the top charities in the ledger
	l.State.TransactionHistory = append(l.State.TransactionHistory, TransactionRecord{
		From:    "system",
		To:      "top_charities",
		Hash:    generateHash("system" + "top_charities" + fmt.Sprintf("%v", topCharities)),
		Status:  "updated",
	})

	return nil
}

// RecordInternalCharityWalletAddition adds funds to the charity pool for future distribution.
func (l *Ledger) RecordInternalCharityWalletAddition(amount float64) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Add funds to the charity pool
	l.charityPool.internalPool += amount

	// Record the transaction in the ledger state
	l.State.TransactionHistory = append(l.State.TransactionHistory, TransactionRecord{
		Amount:  amount,
		From:    "internal_wallet",
		To:      "charity_pool",
		Hash:    generateHash("internal_wallet" + "charity_pool" + strconv.FormatFloat(amount, 'f', 6, 64)),
		Status:  "confirmed",
	})

	return nil
}

// RecordInternalCharityPoolUpdate updates the total funds in the charity pool.
func (l *Ledger) RecordInternalCharityPoolUpdate(newAmount float64) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Update the charity pool
	l.charityPool.internalPool = newAmount

	// Record the transaction in the ledger state
	l.State.TransactionHistory = append(l.State.TransactionHistory, TransactionRecord{
		From:    "internal_wallet",
		To:      "charity_pool",
		Hash:    generateHash("internal_wallet" + "charity_pool" + strconv.FormatFloat(newAmount, 'f', 6, 64)),
		Status:  "updated",
	})

	return nil
}

// RecordInternalCharityFundDistribution records the distribution of funds from the internal charity pool.
func (l *Ledger) RecordInternalCharityFundDistribution(amount float64) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Ensure the charity pool has enough funds for the distribution
	if l.charityPool.internalPool < amount {
		return fmt.Errorf("insufficient funds in charity pool")
	}

	// Update the pool by subtracting the distributed amount
	l.charityPool.internalPool -= amount

	// Record the transaction in the ledger state
	l.State.TransactionHistory = append(l.State.TransactionHistory, TransactionRecord{
		Amount:  amount,
		From:    "charity_pool",
		To:      "fund_distribution",
		Hash:    generateHash("charity_pool" + "fund_distribution" + strconv.FormatFloat(amount, 'f', 6, 64)),
		Status:  "confirmed",
	})

	return nil
}
