package defi

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewLendingManager initializes a new LendingManager for DeFi lending protocols.
func NewLendingManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *LendingManager {
    if ledgerInstance == nil || encryptionService == nil {
        log.Fatalf("[ERROR] Ledger instance and encryption service must not be nil")
    }

    log.Printf("[INFO] Initializing LendingManager")
    return &LendingManager{
        LendingPools:      make(map[string]*LendingPool),
        Loans:             make(map[string]*Loan),
        Ledger:            ledgerInstance,
        EncryptionService: encryptionService,
        mu:                &sync.Mutex{},
    }
}


// CreateLendingPool creates a new lending pool with specified liquidity and interest rate.
func (lm *LendingManager) CreateLendingPool(poolID string, liquidity, interestRate float64) (*LendingPool, error) {
    log.Printf("[INFO] Creating lending pool with PoolID: %s, Liquidity: %.2f, Interest Rate: %.2f", poolID, liquidity, interestRate)

    // Step 1: Input Validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }
    if liquidity <= 0 {
        err := fmt.Errorf("liquidity must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }
    if interestRate <= 0 {
        err := fmt.Errorf("interest rate must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Step 2: Thread-Safe Access
    lm.mu.Lock()
    defer lm.mu.Unlock()

    // Step 3: Check for Existing Pool
    if _, exists := lm.LendingPools[poolID]; exists {
        err := fmt.Errorf("lending pool with PoolID %s already exists", poolID)
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Step 4: Encrypt Pool Data
    poolData := fmt.Sprintf("PoolID: %s, Liquidity: %.2f, InterestRate: %.2f", poolID, liquidity, interestRate)
    encryptedData, err := lm.EncryptionService.EncryptData("AES", []byte(poolData), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt pool data for PoolID %s: %v", poolID, err)
        return nil, fmt.Errorf("failed to encrypt pool data: %w", err)
    }

    // Step 5: Create Lending Pool
    pool := &LendingPool{
        PoolID:         poolID,
        TotalLiquidity: liquidity,
        InterestRate:   interestRate,
        AvailableFunds: liquidity,
        ActiveLoans:    []*Loan{},
        EncryptedData:  string(encryptedData),
    }

    // Step 6: Add Pool to Manager
    lm.LendingPools[poolID] = pool

    // Step 7: Record Creation in Ledger
    if err := lm.Ledger.DeFiLedger.RecordLendingPoolCreation(poolID, liquidity, interestRate); err != nil {
        log.Printf("[ERROR] Failed to log lending pool creation for PoolID %s in ledger: %v", poolID, err)
        return nil, fmt.Errorf("failed to log pool creation in the ledger: %w", err)
    }

    // Step 8: Log Success
    log.Printf("[SUCCESS] Lending pool created successfully. PoolID: %s, Liquidity: %.2f, Interest Rate: %.2f", poolID, liquidity, interestRate)
    return pool, nil
}


// RequestLoan allows a borrower to request a loan from a specific lending pool.
func (lm *LendingManager) RequestLoan(poolID, borrower string, amount, collateral float64, duration time.Duration) (*Loan, error) {
    log.Printf("[INFO] Requesting loan. PoolID: %s, Borrower: %s, Amount: %.2f, Collateral: %.2f, Duration: %v", poolID, borrower, amount, collateral, duration)

    // Step 1: Input Validation
    if poolID == "" || borrower == "" {
        err := fmt.Errorf("poolID and borrower cannot be empty")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }
    if amount <= 0 || collateral <= 0 {
        err := fmt.Errorf("amount and collateral must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }
    if duration <= 0 {
        err := fmt.Errorf("duration must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Step 2: Lock for Thread Safety
    lm.mu.Lock()
    defer lm.mu.Unlock()

    // Step 3: Retrieve and Validate Lending Pool
    pool, exists := lm.LendingPools[poolID]
    if !exists {
        err := fmt.Errorf("lending pool %s not found", poolID)
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    if pool.AvailableFunds < amount {
        err := fmt.Errorf("insufficient liquidity in pool %s to cover the loan", poolID)
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Step 4: Generate Unique Loan ID
    loanID := generateUniqueID()

    // Step 5: Encrypt Loan Data
    loanData := fmt.Sprintf("LoanID: %s, Borrower: %s, Amount: %.2f, Collateral: %.2f, Duration: %v", loanID, borrower, amount, collateral, duration)
    encryptedData, err := lm.EncryptionService.EncryptData("AES", []byte(loanData), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt loan data: %v", err)
        return nil, fmt.Errorf("failed to encrypt loan data: %w", err)
    }

    // Step 6: Create and Register Loan
    loan := &Loan{
        LoanID:        loanID,
        Lender:        poolID,
        Borrower:      borrower,
        Amount:        amount,
        Collateral:    collateral,
        InterestRate:  pool.InterestRate,
        Duration:      duration,
        StartDate:     time.Now(),
        ExpiryDate:    time.Now().Add(duration),
        Status:        "Active",
        EncryptedData: string(encryptedData),
    }

    pool.AvailableFunds -= amount
    pool.ActiveLoans = append(pool.ActiveLoans, loan)
    lm.Loans[loanID] = loan

    // Step 7: Record Loan in Ledger
    if err := lm.Ledger.DeFiLedger.RecordLoanRequest(loanID, poolID, borrower, amount, collateral, pool.InterestRate, duration); err != nil {
        log.Printf("[ERROR] Failed to log loan request in ledger: %v", err)
        return nil, fmt.Errorf("failed to log loan request: %w", err)
    }

    // Step 8: Log Success
    log.Printf("[SUCCESS] Loan %s created for borrower %s. Amount: %.2f, Collateral: %.2f, Duration: %v", loanID, borrower, amount, collateral, duration)
    return loan, nil
}


// RepayLoan allows the borrower to repay their loan and close it.
func (lm *LendingManager) RepayLoan(loanID string) error {
    log.Printf("[INFO] Processing loan repayment. LoanID: %s", loanID)

    // Step 1: Input Validation
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Lock for Thread Safety
    lm.mu.Lock()
    defer lm.mu.Unlock()

    // Step 3: Retrieve and Validate Loan
    loan, exists := lm.Loans[loanID]
    if !exists {
        err := fmt.Errorf("loan %s not found", loanID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    if loan.Status != "Active" {
        err := fmt.Errorf("loan %s is not active", loanID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 4: Retrieve Lending Pool
    pool, exists := lm.LendingPools[loan.Lender]
    if !exists {
        err := fmt.Errorf("lending pool %s not found", loan.Lender)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 5: Calculate Repayment Amount
    repaymentAmount := loan.Amount * (1 + loan.InterestRate)

    // Step 6: Update Pool and Loan Status
    pool.AvailableFunds += repaymentAmount
    loan.Status = "Repaid"

    // Step 7: Log Repayment in Ledger
    if err := lm.Ledger.DeFiLedger.RecordLoanRepayment(loanID, loan.Borrower); err != nil {
        log.Printf("[ERROR] Failed to log repayment in ledger: %v", err)
        return fmt.Errorf("failed to log repayment in ledger: %w", err)
    }

    // Step 8: Log Success
    log.Printf("[SUCCESS] Loan %s repaid by borrower %s. Repayment amount: %.2f", loanID, loan.Borrower, repaymentAmount)
    return nil
}


// LendingCreateLoan creates a new loan and stores it in the ledger.
func LendingCreateLoan(loanID, borrowerID string, principal, interestRate float64, duration time.Duration, collateral string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating loan creation. LoanID: %s, BorrowerID: %s, Principal: %.2f, InterestRate: %.2f, Duration: %v", loanID, borrowerID, principal, interestRate, duration)

    // Step 1: Input Validation
    if loanID == "" || borrowerID == "" || collateral == "" {
        err := fmt.Errorf("loanID, borrowerID, and collateral cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if principal <= 0 || interestRate <= 0 || duration <= 0 {
        err := fmt.Errorf("principal, interestRate, and duration must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt Sensitive Data
    encryptedBorrowerID, err := encryption.EncryptString(borrowerID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt borrower ID: %v", err)
        return fmt.Errorf("failed to encrypt borrower ID: %w", err)
    }
    encryptedCollateral, err := encryption.EncryptString(collateral)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt collateral: %v", err)
        return fmt.Errorf("failed to encrypt collateral: %w", err)
    }

    // Step 3: Create Loan Struct
    loan := Loan{
        LoanID:       loanID,
        BorrowerID:   encryptedBorrowerID,
        Principal:    principal,
        InterestRate: interestRate,
        Duration:     duration,
        Collateral:   encryptedCollateral,
        Status:       "Pending",
        CreatedAt:    time.Now(),
    }

    // Step 4: Record Loan in Ledger
    if err := ledgerInstance.DeFiLedger.CreateLoan(loan); err != nil {
        log.Printf("[ERROR] Failed to record loan in ledger: %v", err)
        return fmt.Errorf("failed to create loan in ledger: %w", err)
    }

    // Step 5: Log Success
    log.Printf("[SUCCESS] Loan %s successfully created for borrower %s. Principal: %.2f, Duration: %v", loanID, borrowerID, principal, duration)
    return nil
}


// LendingApplyInterest applies accrued interest to the specified loan.
func LendingApplyInterest(loanID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Applying interest to loan. LoanID: %s", loanID)

    // Step 1: Input Validation
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Retrieve Loan from Ledger
    loan, err := ledgerInstance.DeFiLedger.GetLoanByID(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to retrieve loan %s: %v", loanID, err)
        return fmt.Errorf("failed to retrieve loan %s: %w", loanID, err)
    }

    // Step 3: Validate Loan Status
    if loan.Status != "Active" {
        err := fmt.Errorf("loan %s is not active; cannot apply interest", loanID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 4: Calculate Interest
    accruedInterest := loan.Principal * loan.InterestRate
    newBalance := loan.Principal + accruedInterest

    // Step 5: Update Loan in Ledger
    if err := ledgerInstance.DeFiLedger.UpdateLoanBalance(loanID, newBalance); err != nil {
        log.Printf("[ERROR] Failed to update loan balance for loan %s: %v", loanID, err)
        return fmt.Errorf("failed to update loan balance: %w", err)
    }

    // Step 6: Log Success
    log.Printf("[SUCCESS] Interest applied to loan %s. New Balance: %.2f (Accrued Interest: %.2f)", loanID, newBalance, accruedInterest)
    return nil
}

// LendingRepayLoan processes the repayment for a specified loan.
func LendingRepayLoan(loanID, borrowerID string, repaymentAmount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Starting repayment process. LoanID: %s, BorrowerID: %s, RepaymentAmount: %.2f", loanID, borrowerID, repaymentAmount)

    // Step 1: Input Validation
    if loanID == "" || borrowerID == "" {
        err := fmt.Errorf("loanID and borrowerID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if repaymentAmount <= 0 {
        err := fmt.Errorf("repaymentAmount must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt Borrower ID
    encryptedBorrowerID, err := encryption.EncryptString(borrowerID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt borrower ID: %v", err)
        return fmt.Errorf("failed to encrypt borrower ID: %w", err)
    }

    // Step 3: Retrieve Loan Details
    loan, err := ledgerInstance.DeFiLedger.GetLoanByID(loanID)
    if err != nil {
        log.Printf("[ERROR] Loan retrieval failed for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to retrieve loan %s: %w", loanID, err)
    }

    // Step 4: Validate Borrower and Status
    if loan.BorrowerID != encryptedBorrowerID {
        err := fmt.Errorf("unauthorized repayment attempt by BorrowerID %s for LoanID %s", borrowerID, loanID)
        log.Printf("[ERROR] %v", err)
        return err
    }
    if loan.Status != "Active" {
        err := fmt.Errorf("loan %s is not active and cannot be repaid", loanID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 5: Process Repayment
    remainingBalance := loan.Principal + loan.AccruedInterest - repaymentAmount
    if remainingBalance < 0 {
        remainingBalance = 0 // Ensure no negative balance
    }

    // Step 6: Update Ledger
    if err := ledgerInstance.DeFiLedger.UpdateLoanBalance(loanID, remainingBalance); err != nil {
        log.Printf("[ERROR] Failed to update loan balance for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to update loan balance: %w", err)
    }

    // Step 7: Finalize Repayment
    if remainingBalance == 0 {
        loan.Status = "Repaid"
        if err := ledgerInstance.DeFiLedger.MarkLoanAsRepaid(loanID); err != nil {
            log.Printf("[ERROR] Failed to mark loan %s as repaid: %v", loanID, err)
            return fmt.Errorf("failed to mark loan as repaid: %w", err)
        }
    }

    log.Printf("[SUCCESS] Loan %s repaid by Borrower %s. Remaining Balance: %.2f", loanID, borrowerID, remainingBalance)
    return nil
}


// LendingLiquidateLoan liquidates a loan in the event of default.
func LendingLiquidateLoan(loanID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating loan liquidation. LoanID: %s", loanID)

    // Step 1: Input Validation
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Retrieve Loan Details
    loan, err := ledgerInstance.DeFiLedger.GetLoanByID(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to retrieve loan %s: %v", loanID, err)
        return fmt.Errorf("failed to retrieve loan %s: %w", loanID, err)
    }

    // Step 3: Validate Loan Status
    if loan.Status != "Defaulted" {
        err := fmt.Errorf("loan %s is not in default and cannot be liquidated", loanID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 4: Liquidate Collateral
    if err := ledgerInstance.DeFiLedger.LiquidateCollateral(loanID, loan.Collateral); err != nil {
        log.Printf("[ERROR] Failed to liquidate collateral for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to liquidate collateral for loan %s: %w", loanID, err)
    }

    // Step 5: Mark Loan as Liquidated
    if err := ledgerInstance.DeFiLedger.MarkLoanAsLiquidated(loanID); err != nil {
        log.Printf("[ERROR] Failed to mark loan %s as liquidated: %v", loanID, err)
        return fmt.Errorf("failed to mark loan as liquidated: %w", err)
    }

    log.Printf("[SUCCESS] Loan %s successfully liquidated. Collateral processed.", loanID)
    return nil
}


// LendingTrackRepayment tracks the repayment history of a loan.
func LendingTrackRepayment(loanID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating repayment tracking for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Fetch Loan from Ledger
    repaymentHistory, err := ledgerInstance.DeFiLedger.TrackRepayment(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to retrieve repayment history for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to track repayment for loan %s: %w", loanID, err)
    }

    // Step 3: Log Repayment History Details
    for _, record := range repaymentHistory {
        log.Printf("[INFO] LoanID: %s | Amount: %.2f | Date: %s", loanID, record.Amount, record.Date.Format(time.RFC3339))
    }

    log.Printf("[SUCCESS] Repayment tracking completed for LoanID: %s", loanID)
    return nil
}


// LendingFetchLoanStatus retrieves the current status of a loan.
func LendingFetchLoanStatus(loanID string, ledgerInstance *ledger.Ledger) (string, error) {
    log.Printf("[INFO] Fetching status for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return "", err
    }

    // Step 2: Retrieve Loan Status from Ledger
    status, err := ledgerInstance.DeFiLedger.FetchLoanStatus(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch loan status for LoanID %s: %v", loanID, err)
        return "", fmt.Errorf("failed to fetch status for loan %s: %w", loanID, err)
    }

    // Step 3: Log Retrieved Status
    log.Printf("[INFO] LoanID: %s | Status: %s", loanID, status)
    return status, nil
}


// LendingAdjustInterestRate adjusts the interest rate of a specified loan.
func LendingAdjustInterestRate(loanID string, newRate float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Adjusting interest rate for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if newRate <= 0 {
        err := fmt.Errorf("newRate must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Adjust the Interest Rate in the Ledger
    if err := ledgerInstance.DeFiLedger.AdjustInterestRate(loanID, newRate); err != nil {
        log.Printf("[ERROR] Failed to adjust interest rate for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to adjust interest rate for loan %s: %w", loanID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Interest rate for LoanID %s adjusted to %.2f%%", loanID, newRate)
    return nil
}


// LendingEscrowCollateral escrows the collateral for a specified loan.
func LendingEscrowCollateral(loanID, collateral string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Escrowing collateral for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" || collateral == "" {
        err := fmt.Errorf("loanID and collateral cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt Collateral Data
    encryptedCollateral, err := encryption.EncryptString(collateral)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt collateral for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to encrypt collateral for loan %s: %w", loanID, err)
    }

    // Step 3: Store the Collateral in the Ledger
    if err := ledgerInstance.DeFiLedger.EscrowCollateral(loanID, encryptedCollateral); err != nil {
        log.Printf("[ERROR] Failed to escrow collateral for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to escrow collateral for loan %s: %w", loanID, err)
    }

    // Step 4: Log Success
    log.Printf("[SUCCESS] Collateral for LoanID %s successfully escrowed", loanID)
    return nil
}

// LendingReleaseCollateral releases the collateral for a specified loan.
func LendingReleaseCollateral(loanID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating collateral release for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Verify Loan Status
    loanStatus, err := ledgerInstance.DeFiLedger.FetchLoanStatus(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch loan status for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to fetch loan status for loan %s: %w", loanID, err)
    }
    if loanStatus != "Repaid" {
        err := fmt.Errorf("loan %s is not in a state to release collateral, current status: %s", loanID, loanStatus)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 3: Release the Collateral
    if err := ledgerInstance.DeFiLedger.ReleaseCollateral(loanID); err != nil {
        log.Printf("[ERROR] Failed to release collateral for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to release collateral for loan %s: %w", loanID, err)
    }

    // Step 4: Log Success
    log.Printf("[SUCCESS] Collateral for LoanID %s successfully released", loanID)
    return nil
}


// LendingVerifyCollateral verifies the collateral for a specified loan.
func LendingVerifyCollateral(loanID, collateral string, ledgerInstance *ledger.Ledger) (bool, error) {
    log.Printf("[INFO] Initiating collateral verification for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" || collateral == "" {
        err := fmt.Errorf("loanID and collateral cannot be empty")
        log.Printf("[ERROR] %v", err)
        return false, err
    }

    // Step 2: Encrypt Collateral Data
    encryptedCollateral, err := encryption.EncryptString(collateral)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt collateral for LoanID %s: %v", loanID, err)
        return false, fmt.Errorf("failed to encrypt collateral for loan %s: %w", loanID, err)
    }

    // Step 3: Verify the Collateral in the Ledger
    verified, err := ledgerInstance.DeFiLedger.VerifyCollateral(loanID, encryptedCollateral)
    if err != nil {
        log.Printf("[ERROR] Failed to verify collateral for LoanID %s: %v", loanID, err)
        return false, fmt.Errorf("failed to verify collateral for loan %s: %w", loanID, err)
    }

    // Step 4: Log Result
    log.Printf("[SUCCESS] Collateral verification for LoanID %s: %v", loanID, verified)
    return verified, nil
}


// LendingAudit audits the details of a specified loan.
func LendingAudit(loanID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating audit for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Perform Loan Audit
    err := ledgerInstance.DeFiLedger.AuditLoan(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to audit LoanID %s: %v", loanID, err)
        return fmt.Errorf("audit failed for loan %s: %w", loanID, err)
    }

    // Step 3: Log Audit Success
    log.Printf("[SUCCESS] Audit completed for LoanID %s", loanID)
    return nil
}


// LendingCalculateLoanHealth calculates the health of a specified loan.
func LendingCalculateLoanHealth(loanID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Calculating loan health for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Calculate Loan Health
    health, err := ledgerInstance.DeFiLedger.CalculateLoanHealth(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to calculate loan health for LoanID %s: %v", loanID, err)
        return 0, fmt.Errorf("loan health calculation failed for loan %s: %w", loanID, err)
    }

    // Step 3: Log Calculation Result
    log.Printf("[SUCCESS] Loan health for LoanID %s: %.2f", loanID, health)
    return health, nil
}


// LendingTrackBorrower tracks the activity of a specified borrower.
func LendingTrackBorrower(borrowerID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating borrower activity tracking for BorrowerID: %s", borrowerID)

    // Step 1: Validate Input Parameters
    if borrowerID == "" {
        err := fmt.Errorf("borrowerID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt Borrower ID
    encryptedBorrowerID, err := encryption.EncryptString(borrowerID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt BorrowerID %s: %v", borrowerID, err)
        return fmt.Errorf("encryption error for borrower %s: %w", borrowerID, err)
    }

    // Step 3: Track Borrower Activity
    err = ledgerInstance.DeFiLedger.TrackBorrower(encryptedBorrowerID)
    if err != nil {
        log.Printf("[ERROR] Failed to track activity for BorrowerID %s: %v", borrowerID, err)
        return fmt.Errorf("failed to track activity for borrower %s: %w", borrowerID, err)
    }

    // Step 4: Log Success
    log.Printf("[SUCCESS] Activity tracking completed for BorrowerID %s", borrowerID)
    return nil
}



// LendingSetCollateralRequirement sets the collateral requirement for a specific loan.
func LendingSetCollateralRequirement(loanID string, requirement float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Setting collateral requirement for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if requirement <= 0 {
        err := fmt.Errorf("requirement must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Update Collateral Requirement
    err := ledgerInstance.DeFiLedger.SetCollateralRequirement(loanID, requirement)
    if err != nil {
        log.Printf("[ERROR] Failed to set collateral requirement for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to set collateral requirement for loan %s: %w", loanID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Collateral requirement for LoanID %s set to %.2f", loanID, requirement)
    return nil
}


// LendingFetchCollateralRequirement fetches the collateral requirement for a specific loan.
func LendingFetchCollateralRequirement(loanID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Fetching collateral requirement for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Fetch Collateral Requirement
    requirement, err := ledgerInstance.DeFiLedger.FetchCollateralRequirement(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch collateral requirement for LoanID %s: %v", loanID, err)
        return 0, fmt.Errorf("failed to fetch collateral requirement for loan %s: %w", loanID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Collateral requirement for LoanID %s: %.2f", loanID, requirement)
    return requirement, nil
}


// LendingUpdateRepaymentSchedule updates the repayment schedule for a specific loan.
func LendingUpdateRepaymentSchedule(loanID string, newSchedule []time.Time, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Updating repayment schedule for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if len(newSchedule) == 0 {
        err := fmt.Errorf("newSchedule must contain at least one repayment date")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Sort and Validate the Repayment Schedule
    sort.Slice(newSchedule, func(i, j int) bool {
        return newSchedule[i].Before(newSchedule[j])
    })
    now := time.Now()
    for _, date := range newSchedule {
        if date.Before(now) {
            err := fmt.Errorf("repayment dates cannot be in the past")
            log.Printf("[ERROR] %v", err)
            return err
        }
    }

    // Step 3: Update Repayment Schedule in the Ledger
    err := ledgerInstance.DeFiLedger.UpdateRepaymentSchedule(loanID, newSchedule)
    if err != nil {
        log.Printf("[ERROR] Failed to update repayment schedule for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to update repayment schedule for loan %s: %w", loanID, err)
    }

    // Step 4: Log Success
    log.Printf("[SUCCESS] Repayment schedule for LoanID %s updated successfully", loanID)
    return nil
}


// LendingPauseRepayments pauses repayments for a specific loan.
func LendingPauseRepayments(loanID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Attempting to pause repayments for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Pause Repayments in the Ledger
    err := ledgerInstance.DeFiLedger.PauseRepayments(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to pause repayments for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to pause repayments for loan %s: %w", loanID, err)
    }

    // Step 3: Post-Action Validation and Logging
    log.Printf("[SUCCESS] Repayments for LoanID %s paused successfully", loanID)
    return nil
}


// LendingResumeRepayments resumes repayments for a specific loan.
func LendingResumeRepayments(loanID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Attempting to resume repayments for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Resume Repayments in the Ledger
    err := ledgerInstance.DeFiLedger.ResumeRepayments(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to resume repayments for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to resume repayments for loan %s: %w", loanID, err)
    }

    // Step 3: Post-Action Validation and Logging
    log.Printf("[SUCCESS] Repayments for LoanID %s resumed successfully", loanID)
    return nil
}


// LendingTrackLatePayments tracks late payments for a specific loan and logs relevant details.
func LendingTrackLatePayments(loanID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating late payment tracking for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Track Late Payments in the Ledger
    err := ledgerInstance.DeFiLedger.TrackLatePayments(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to track late payments for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to track late payments for loan %s: %w", loanID, err)
    }

    // Step 3: Post-Action Validation and Logging
    log.Printf("[SUCCESS] Late payments for LoanID %s tracked successfully", loanID)
    return nil
}


// LendingFetchLatePaymentHistory retrieves the late payment history for a specific loan.
func LendingFetchLatePaymentHistory(loanID string, ledgerInstance *ledger.Ledger) ([]ledger.LatePaymentRecord, error) {
    log.Printf("[INFO] Fetching late payment history for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Step 2: Fetch Late Payment History from the Ledger
    history, err := ledgerInstance.DeFiLedger.FetchLatePaymentHistory(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch late payment history for LoanID %s: %v", loanID, err)
        return nil, fmt.Errorf("failed to fetch late payment history for loan %s: %w", loanID, err)
    }

    // Step 3: Post-Action Logging and Return
    log.Printf("[SUCCESS] Late payment history for LoanID %s fetched successfully. Records: %d", loanID, len(history))
    return history, nil
}


// LendingFetchInterestRate retrieves the current interest rate for a specific loan.
func LendingFetchInterestRate(loanID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Fetching interest rate for LoanID: %s", loanID)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Fetch Interest Rate from the Ledger
    rate, err := ledgerInstance.DeFiLedger.FetchInterestRate(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch interest rate for LoanID %s: %v", loanID, err)
        return 0, fmt.Errorf("failed to fetch interest rate for loan %s: %w", loanID, err)
    }

    // Step 3: Post-Action Logging and Return
    log.Printf("[SUCCESS] Interest rate for LoanID %s fetched successfully: %.2f%%", loanID, rate)
    return rate, nil
}


// LendingAdjustInterestRatePeriod adjusts the period for interest rate application for a specific loan.
func LendingAdjustInterestRatePeriod(loanID string, period time.Duration, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Adjusting interest rate period for LoanID: %s to %s", loanID, period)

    // Step 1: Validate Input Parameters
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if period <= 0 {
        err := fmt.Errorf("interest rate period must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Adjust Interest Rate Period in the Ledger
    err := ledgerInstance.DeFiLedger.AdjustInterestRatePeriod(loanID, period)
    if err != nil {
        log.Printf("[ERROR] Failed to adjust interest rate period for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to adjust interest rate period for loan %s: %w", loanID, err)
    }

    // Step 3: Post-Action Logging
    log.Printf("[SUCCESS] Interest rate period for LoanID %s adjusted successfully to %s", loanID, period)
    return nil
}


// LendingAutoAdjustRepayments enables automatic repayment adjustments for a specific loan.
func LendingAutoAdjustRepayments(loanID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating auto-adjustment of repayments for LoanID: %s", loanID)

    // Step 1: Input Validation
    if loanID == "" {
        err := fmt.Errorf("loanID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Perform Pre-Action Verification
    log.Printf("[INFO] Verifying loan state for auto-repayment adjustments (LoanID: %s)", loanID)
    isEligible, err := ledgerInstance.DeFiLedger.CheckLoanEligibilityForAutoAdjustments(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to verify loan eligibility for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to verify loan eligibility for LoanID %s: %w", loanID, err)
    }
    if !isEligible {
        err := fmt.Errorf("loan %s is not eligible for auto-repayment adjustments", loanID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 3: Enable Auto-Repayment Adjustments in the Ledger
    log.Printf("[INFO] Enabling auto-repayment adjustments for LoanID: %s", loanID)
    err = ledgerInstance.DeFiLedger.EnableAutoRepaymentAdjustments(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to enable auto-repayment adjustments for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to enable auto repayment adjustments for LoanID %s: %w", loanID, err)
    }

    // Step 4: Post-Action Verification
    log.Printf("[INFO] Verifying auto-repayment adjustments were successfully enabled (LoanID: %s)", loanID)
    isEnabled, err := ledgerInstance.DeFiLedger.IsAutoAdjustmentEnabled(loanID)
    if err != nil {
        log.Printf("[ERROR] Failed to verify auto-repayment adjustment status for LoanID %s: %v", loanID, err)
        return fmt.Errorf("failed to verify auto-repayment adjustment status for LoanID %s: %w", loanID, err)
    }
    if !isEnabled {
        err := fmt.Errorf("auto repayment adjustments were not successfully enabled for LoanID %s", loanID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 5: Log Success
    log.Printf("[SUCCESS] Auto repayment adjustments successfully enabled for LoanID: %s", loanID)
    return nil
}

