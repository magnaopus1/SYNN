package syn10

import (
    "errors"
    "math/big"
    "time"

    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)



// NewAuditComplianceManager initializes a new AuditComplianceManager.
func NewAuditComplianceManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *SYN10AuditComplianceManager {
    return &SYN10AuditComplianceManager{
        Ledger:            ledgerInstance,
        Encryption:        encryptionService,
        AuditLogs:         make(map[string]SYN10AuditLog),
        RegulatoryReports: make(map[string]SYN10RegulatoryReport),
    }
}

// LogActivity logs an audit activity and stores it in the ledger.
func (acm *SYN10AuditComplianceManager) LogActivity(activity, userID, details string) error {
    acm.mutex.Lock()
    defer acm.mutex.Unlock()

    log := SYN10AuditLog{
        LogID:     generateLogID(activity, userID),
        Timestamp: time.Now(),
        Activity:  activity,
        UserID:    userID,
        Details:   details,
    }

    // Encrypt the audit log details before storing them in the ledger
    encryptedDetails, err := acm.Encryption.EncryptData(details, common.EncryptionKey)
    if err != nil {
        return err
    }
    log.Details = encryptedDetails

    acm.AuditLogs[log.LogID] = log

    return acm.Ledger.AddAuditLog(log)
}

// VerifyAuditLog verifies the signature of an audit log.
func (acm *SYN10AuditComplianceManager) VerifyAuditLog(logID string) (bool, error) {
    acm.mutex.Lock()
    defer acm.mutex.Unlock()

    log, exists := acm.AuditLogs[logID]
    if !exists {
        return false, errors.New("audit log not found")
    }

    return acm.Encryption.VerifySignature(log.Details, log.Signature), nil
}

// NewMonetaryPolicyManager initializes a new MonetaryPolicyManager.
func NewMonetaryPolicyManager(easingMechanism *QuantitativeEasingMechanism, tighteningMechanism *MonetaryTighteningMechanism) *MonetaryPolicyManager {
    return &MonetaryPolicyManager{
        TokenSupply:        new(big.Float),
        InterestRates:      make(map[string]*big.Float),
        MintedTokens:       make(map[string]*big.Float),
        BurnedTokens:       make(map[string]*big.Float),
        TransactionLog:     []*MonetaryTransaction{},
        EasingMechanism:    easingMechanism,
        TighteningMechanism: tighteningMechanism,
    }
}

// MintTokens mints new tokens, logs the transaction, and updates the supply.
func (mpm *MonetaryPolicyManager) MintTokens(tokenType string, amount *big.Float, details string) error {
    mpm.mutex.Lock()
    defer mpm.mutex.Unlock()

    mpm.TokenSupply.Add(mpm.TokenSupply, amount)
    if mpm.MintedTokens[tokenType] == nil {
        mpm.MintedTokens[tokenType] = new(big.Float)
    }
    mpm.MintedTokens[tokenType].Add(mpm.MintedTokens[tokenType], amount)

    transaction := &MonetaryTransaction{
        TransactionType: "Mint",
        TokenType:       tokenType,
        Amount:          amount,
        Timestamp:       time.Now(),
        Details:         details,
    }
    mpm.TransactionLog = append(mpm.TransactionLog, transaction)

    return nil
}

// BurnTokens burns tokens, logs the transaction, and updates the supply.
func (mpm *MonetaryPolicyManager) BurnTokens(tokenType string, amount *big.Float, details string) error {
    mpm.mutex.Lock()
    defer mpm.mutex.Unlock()

    mpm.TokenSupply.Sub(mpm.TokenSupply, amount)
    if mpm.BurnedTokens[tokenType] == nil {
        mpm.BurnedTokens[tokenType] = new(big.Float)
    }
    mpm.BurnedTokens[tokenType].Add(mpm.BurnedTokens[tokenType], amount)

    transaction := &MonetaryTransaction{
        TransactionType: "Burn",
        TokenType:       tokenType,
        Amount:          amount,
        Timestamp:       time.Now(),
        Details:         details,
    }
    mpm.TransactionLog = append(mpm.TransactionLog, transaction)

    return nil
}

// NewInterestRateManager initializes an InterestRateManager.
func NewInterestRateManager(savingsRate, commercialRate, userRate *big.Float, updateInterval time.Duration) *InterestRateManager {
    return &InterestRateManager{
        SavingsBaseRate:         savingsRate,
        CommercialBorrowingRate: commercialRate,
        UserBorrowingRate:       userRate,
        RateUpdateInterval:      updateInterval,
        LastUpdated:             time.Now(),
    }
}

// UpdateRates updates savings, borrowing, and commercial rates periodically.
func (irm *InterestRateManager) UpdateRates(savingsRate, commercialRate, userRate *big.Float) error {
    irm.mutex.Lock()
    defer irm.mutex.Unlock()

    if time.Since(irm.LastUpdated) < irm.RateUpdateInterval {
        return errors.New("rate update interval not elapsed")
    }

    irm.SavingsBaseRate = savingsRate
    irm.CommercialBorrowingRate = commercialRate
    irm.UserBorrowingRate = userRate
    irm.LastUpdated = time.Now()

    return nil
}

// ConductQuantitativeEasing implements quantitative easing via asset purchases.
func (mpm *MonetaryPolicyManager) ConductQuantitativeEasing(amount *big.Float, details string) error {
    mpm.mutex.Lock()
    defer mpm.mutex.Unlock()

    if err := mpm.EasingMechanism.BuyAssets(amount); err != nil {
        return err
    }

    mpm.TokenSupply.Add(mpm.TokenSupply, amount)

    transaction := &MonetaryTransaction{
        TransactionType: "QuantitativeEasing",
        TokenType:       "SYN10",
        Amount:          amount,
        Timestamp:       time.Now(),
        Details:         details,
    }
    mpm.TransactionLog = append(mpm.TransactionLog, transaction)

    return nil
}

// NewPeggingMechanism initializes a new PeggingMechanism for SYN10 tokens.
func NewPeggingMechanism(fiatCurrency string, pegValue, initialValue *big.Float) *SYN10PeggingMechanism {
    return &SYN10PeggingMechanism{
        FiatCurrency:      fiatCurrency,
        PegValue:          pegValue,
        CurrentValue:      initialValue,
        CollateralReserves: make(map[string]*big.Float),
        StabilizationActive: true,
        RemovalDate:       time.Time{},
    }
}

// UpdatePegValue updates the pegged value and adjusts collateral.
func (pm *SYN10PeggingMechanism) UpdatePegValue(newPegValue *big.Float) error {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    if !pm.StabilizationActive {
        return errors.New("peg stabilization inactive")
    }

    pm.PegValue = newPegValue
    pm.adjustCollateral()
    pm.CurrentValue = newPegValue

    return nil
}

// AddCollateral adds collateral to the reserves for stabilization.
func (pm *SYN10PeggingMechanism) AddCollateral(asset string, amount *big.Float) {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    if pm.CollateralReserves[asset] == nil {
        pm.CollateralReserves[asset] = new(big.Float)
    }
    pm.CollateralReserves[asset].Add(pm.CollateralReserves[asset], amount)
}

