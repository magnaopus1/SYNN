package defi

import (
	"fmt"
	"log"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewSyntheticAssetManager initializes the manager for synthetic assets.
// This function sets up the necessary maps for asset tracking and integrates the ledger and encryption services.
func NewSyntheticAssetManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *SyntheticAssetManager {
    if ledgerInstance == nil || encryptionService == nil {
        log.Fatalf("[ERROR] Ledger instance or encryption service cannot be nil")
    }

    log.Printf("[INFO] Initializing SyntheticAssetManager")
    return &SyntheticAssetManager{
        Assets:            make(map[string]*SyntheticAsset),
        Ledger:            ledgerInstance,
        EncryptionService: encryptionService,
    }
}


// CreateSyntheticAsset allows the creation of a new synthetic asset backed by collateral.
// It validates input parameters, encrypts asset data, and logs the creation in the ledger.
func (sam *SyntheticAssetManager) CreateSyntheticAsset(assetName, underlyingAsset string, price, collateralRatio, initialSupply float64) (*SyntheticAsset, error) {
    log.Printf("[INFO] Creating synthetic asset: %s (Underlying: %s)", assetName, underlyingAsset)

    // Step 1: Input Validation
    if assetName == "" || underlyingAsset == "" {
        err := fmt.Errorf("assetName and underlyingAsset cannot be empty")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }
    if price <= 0 || collateralRatio <= 0 || initialSupply < 0 {
        err := fmt.Errorf("price, collateralRatio, and initialSupply must be positive values")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Step 2: Generate Unique Asset ID
    assetID := generateUniqueID()
    log.Printf("[INFO] Generated unique AssetID: %s", assetID)

    // Step 3: Encrypt Asset Data
    assetData := fmt.Sprintf("AssetID: %s, Name: %s, Underlying: %s, Price: %.2f, CollateralRatio: %.2f, Supply: %.2f",
        assetID, assetName, underlyingAsset, price, collateralRatio, initialSupply)
    encryptedData, err := sam.EncryptionService.EncryptData("AES", []byte(assetData), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt asset data: %v", err)
        return nil, fmt.Errorf("failed to encrypt asset data: %w", err)
    }
    log.Printf("[INFO] Asset data encrypted successfully for AssetID: %s", assetID)

    // Step 4: Create the Synthetic Asset
    asset := &SyntheticAsset{
        AssetID:         assetID,
        AssetName:       assetName,
        UnderlyingAsset: underlyingAsset,
        Price:           price,
        CollateralRatio: collateralRatio,
        TotalSupply:     initialSupply,
        CreatedAt:       time.Now(),
        Status:          "Active",
        EncryptedData:   string(encryptedData),
    }

    // Step 5: Update Asset Manager State
    sam.mu.Lock()
    sam.Assets[assetID] = asset
    sam.mu.Unlock()

    // Step 6: Log Asset Creation in the Ledger
    log.Printf("[INFO] Recording asset creation in the ledger for AssetID: %s", assetID)
    err = sam.Ledger.DeFiLedger.RecordSyntheticAssetCreation(assetID, assetName, underlyingAsset, price, collateralRatio)
    if err != nil {
        log.Printf("[ERROR] Failed to log asset creation in the ledger: %v", err)
        return nil, fmt.Errorf("failed to log asset creation in the ledger: %w", err)
    }

    // Step 7: Success Logging
    log.Printf("[SUCCESS] Synthetic asset %s (AssetID: %s) created successfully with initial supply: %.2f", assetName, assetID, initialSupply)
    return asset, nil
}


// MintSyntheticAsset allows a user to mint a new supply of a synthetic asset, backed by collateral.
// It ensures sufficient collateral is provided, updates the asset supply, and logs the event in the ledger.
func (sam *SyntheticAssetManager) MintSyntheticAsset(assetID string, additionalSupply, collateralAmount float64) error {
    log.Printf("[INFO] Initiating minting process for AssetID: %s", assetID)

    sam.mu.Lock()
    defer sam.mu.Unlock()

    // Step 1: Validate Inputs
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if additionalSupply <= 0 || collateralAmount <= 0 {
        err := fmt.Errorf("additionalSupply and collateralAmount must be positive values")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Retrieve Synthetic Asset
    asset, exists := sam.Assets[assetID]
    if !exists {
        err := fmt.Errorf("synthetic asset %s not found", assetID)
        log.Printf("[ERROR] %v", err)
        return err
    }
    log.Printf("[INFO] Retrieved asset: %s", asset.AssetName)

    // Step 3: Validate Collateral Sufficiency
    requiredCollateral := additionalSupply * asset.Price * asset.CollateralRatio
    if collateralAmount < requiredCollateral {
        err := fmt.Errorf("insufficient collateral: %.2f required, %.2f provided", requiredCollateral, collateralAmount)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 4: Update Asset Supply
    asset.TotalSupply += additionalSupply
    log.Printf("[INFO] Updated supply for AssetID: %s. New TotalSupply: %.2f", assetID, asset.TotalSupply)

    // Step 5: Log Minting Event in Ledger
    err := sam.Ledger.DeFiLedger.RecordMintingEvent(assetID, additionalSupply)
    if err != nil {
        log.Printf("[ERROR] Failed to log minting event in the ledger: %v", err)
        return fmt.Errorf("failed to log minting event in the ledger: %w", err)
    }

    // Step 6: Finalize and Log
    log.Printf("[SUCCESS] Minted %.2f units of synthetic asset %s", additionalSupply, asset.AssetName)
    return nil
}


// BurnSyntheticAsset allows a user to burn synthetic assets, reducing supply and reclaiming collateral.
// It validates inputs, updates the total supply, and logs the event in the ledger.
func (sam *SyntheticAssetManager) BurnSyntheticAsset(assetID string, burnAmount, collateralToReclaim float64) error {
    log.Printf("[INFO] Initiating burn process for AssetID: %s", assetID)

    sam.mu.Lock()
    defer sam.mu.Unlock()

    // Step 1: Validate Inputs
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if burnAmount <= 0 || collateralToReclaim < 0 {
        err := fmt.Errorf("burnAmount must be positive, and collateralToReclaim must be non-negative")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Retrieve Synthetic Asset
    asset, exists := sam.Assets[assetID]
    if !exists {
        err := fmt.Errorf("synthetic asset %s not found", assetID)
        log.Printf("[ERROR] %v", err)
        return err
    }
    log.Printf("[INFO] Retrieved asset: %s", asset.AssetName)

    // Step 3: Validate Burn Amount
    if burnAmount > asset.TotalSupply {
        err := fmt.Errorf("cannot burn more than the total supply of %s", asset.AssetName)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 4: Update Asset Supply
    asset.TotalSupply -= burnAmount
    log.Printf("[INFO] Updated supply for AssetID: %s. New TotalSupply: %.2f", assetID, asset.TotalSupply)

    // Step 5: Log Burning Event in Ledger
    err := sam.Ledger.DeFiLedger.RecordBurningEvent(assetID, burnAmount)
    if err != nil {
        log.Printf("[ERROR] Failed to log burning event in the ledger: %v", err)
        return fmt.Errorf("failed to log burning event in the ledger: %w", err)
    }

    // Step 6: Finalize and Log
    log.Printf("[SUCCESS] Burned %.2f units of synthetic asset %s and reclaimed collateral: %.2f", burnAmount, asset.AssetName, collateralToReclaim)
    return nil
}


// GetAssetDetails retrieves details of a synthetic asset by its ID.
// It ensures the asset exists before returning its details.
func (sam *SyntheticAssetManager) GetAssetDetails(assetID string) (*SyntheticAsset, error) {
    log.Printf("[INFO] Fetching details for AssetID: %s", assetID)

    sam.mu.Lock()
    defer sam.mu.Unlock()

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Step 2: Retrieve Synthetic Asset
    asset, exists := sam.Assets[assetID]
    if !exists {
        err := fmt.Errorf("synthetic asset %s not found", assetID)
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Step 3: Return Asset Details
    log.Printf("[SUCCESS] Retrieved details for AssetID: %s, AssetName: %s", assetID, asset.AssetName)
    return asset, nil
}


// SyntheticAssetMint handles the minting of synthetic assets through the ledger.
// It validates the amount and logs the minting event in the ledger.
func SyntheticAssetMint(assetID string, amount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating minting process for AssetID: %s", assetID)

    // Step 1: Validate Inputs
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if amount <= 0 {
        err := fmt.Errorf("mint amount must be positive")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Mint Synthetic Asset
    err := ledgerInstance.DeFiLedger.MintSyntheticAsset(assetID, amount)
    if err != nil {
        log.Printf("[ERROR] Failed to mint synthetic asset: %v", err)
        return fmt.Errorf("failed to mint synthetic asset: %w", err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Minted %.2f units of synthetic asset %s", amount, assetID)
    return nil
}


// SyntheticAssetBurn handles the burning of synthetic assets through the ledger.
// It validates the amount and logs the burning event in the ledger.
func SyntheticAssetBurn(assetID string, amount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating burn process for AssetID: %s", assetID)

    // Step 1: Validate Inputs
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if amount <= 0 {
        err := fmt.Errorf("burn amount must be positive")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Execute Burn in Ledger
    err := ledgerInstance.DeFiLedger.BurnSyntheticAsset(assetID, amount)
    if err != nil {
        log.Printf("[ERROR] Failed to burn synthetic asset: %v", err)
        return fmt.Errorf("failed to burn synthetic asset: %w", err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Burned %.2f units of synthetic asset %s", amount, assetID)
    return nil
}


// SyntheticAssetPriceFeed updates the price feed of a synthetic asset in the ledger.
// It validates the new price and logs the update in the ledger.
func SyntheticAssetPriceFeed(assetID string, newPrice float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Updating price feed for AssetID: %s", assetID)

    // Step 1: Validate Inputs
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if newPrice <= 0 {
        err := fmt.Errorf("new price must be positive")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Update Price Feed in Ledger
    err := ledgerInstance.DeFiLedger.UpdateSyntheticAssetPrice(assetID, newPrice)
    if err != nil {
        log.Printf("[ERROR] Failed to update price feed for synthetic asset: %v", err)
        return fmt.Errorf("failed to update price feed: %w", err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Price feed for synthetic asset %s updated to %.2f", assetID, newPrice)
    return nil
}



// SyntheticAssetSetCollateral sets the collateral amount for a synthetic asset.
// It validates inputs and logs the operation in the ledger.
func SyntheticAssetSetCollateral(assetID string, collateralAmount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating collateral update for AssetID: %s", assetID)

    // Step 1: Validate Inputs
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if collateralAmount <= 0 {
        err := fmt.Errorf("collateralAmount must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Update Collateral in Ledger
    err := ledgerInstance.DeFiLedger.SetAssetCollateral(assetID, collateralAmount)
    if err != nil {
        log.Printf("[ERROR] Failed to set collateral for AssetID %s: %v", assetID, err)
        return fmt.Errorf("failed to set collateral for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Collateral of %.2f set for synthetic asset %s", collateralAmount, assetID)
    return nil
}

// SyntheticAssetVerifyCollateral verifies the collateral backing a synthetic asset.
// It validates the assetID and logs the verification process.
func SyntheticAssetVerifyCollateral(assetID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating collateral verification for AssetID: %s", assetID)

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Verify Collateral in Ledger
    err := ledgerInstance.DeFiLedger.VerifyAssetCollateral(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to verify collateral for AssetID %s: %v", assetID, err)
        return fmt.Errorf("failed to verify collateral for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Collateral for synthetic asset %s verified successfully", assetID)
    return nil
}


// SyntheticAssetLiquidate liquidates a synthetic asset due to insufficient collateral or other conditions.
// It validates the assetID and logs the liquidation event.
func SyntheticAssetLiquidate(assetID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating liquidation process for AssetID: %s", assetID)

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Perform Liquidation in Ledger
    err := ledgerInstance.DeFiLedger.LiquidateSyntheticAsset(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to liquidate synthetic asset %s: %v", assetID, err)
        return fmt.Errorf("failed to liquidate synthetic asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Synthetic asset %s successfully liquidated", assetID)
    return nil
}


// SyntheticAssetDistributeDividends distributes dividends for a synthetic asset.
// It validates the dividend amount and logs the distribution in the ledger.
func SyntheticAssetDistributeDividends(assetID string, dividendAmount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating dividend distribution for AssetID: %s", assetID)

    // Step 1: Validate Inputs
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if dividendAmount <= 0 {
        err := fmt.Errorf("dividendAmount must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Distribute Dividends in Ledger
    err := ledgerInstance.DeFiLedger.DistributeAssetDividends(assetID, dividendAmount)
    if err != nil {
        log.Printf("[ERROR] Failed to distribute dividends for AssetID %s: %v", assetID, err)
        return fmt.Errorf("failed to distribute dividends for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Dividends of %f distributed for synthetic asset %s", dividendAmount, assetID)
    return nil
}


// SyntheticAssetSetDividendRate sets the dividend rate for a synthetic asset.
// It validates the rate and logs the update in the ledger.
func SyntheticAssetSetDividendRate(assetID string, rate float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating dividend rate update for AssetID: %s", assetID)

    // Step 1: Validate Inputs
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if rate < 0 {
        err := fmt.Errorf("dividend rate cannot be negative")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Update Dividend Rate in Ledger
    err := ledgerInstance.DeFiLedger.SetDividendRate(assetID, rate)
    if err != nil {
        log.Printf("[ERROR] Failed to set dividend rate for AssetID %s: %v", assetID, err)
        return fmt.Errorf("failed to set dividend rate for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Dividend rate of %f set for synthetic asset %s", rate, assetID)
    return nil
}


// SyntheticAssetFetchDividendRate retrieves the current dividend rate for a synthetic asset.
// It validates the assetID and logs the retrieval process.
func SyntheticAssetFetchDividendRate(assetID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Fetching dividend rate for AssetID: %s", assetID)

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Fetch Dividend Rate from Ledger
    rate, err := ledgerInstance.DeFiLedger.GetDividendRate(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch dividend rate for AssetID %s: %v", assetID, err)
        return 0, fmt.Errorf("failed to fetch dividend rate for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Fetched dividend rate of %f for synthetic asset %s", rate, assetID)
    return rate, nil
}


// SyntheticAssetTrackMarketCap enables tracking of the market capitalization for a synthetic asset.
func SyntheticAssetTrackMarketCap(assetID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating market cap tracking for AssetID: %s", assetID)

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Enable Market Cap Tracking in Ledger
    err := ledgerInstance.DeFiLedger.TrackMarketCap(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to enable market cap tracking for AssetID %s: %v", assetID, err)
        return fmt.Errorf("failed to enable market cap tracking for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Market cap tracking enabled for synthetic asset %s", assetID)
    return nil
}


// SyntheticAssetFetchMarketCap retrieves the market capitalization of a synthetic asset.
func SyntheticAssetFetchMarketCap(assetID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Fetching market cap for AssetID: %s", assetID)

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Fetch Market Cap from Ledger
    marketCap, err := ledgerInstance.DeFiLedger.GetMarketCap(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch market cap for AssetID %s: %v", assetID, err)
        return 0, fmt.Errorf("failed to fetch market cap for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Fetched market cap of %f for synthetic asset %s", marketCap, assetID)
    return marketCap, nil
}


// SyntheticAssetFetchCollateralRatio retrieves the collateral ratio for a synthetic asset.
func SyntheticAssetFetchCollateralRatio(assetID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Fetching collateral ratio for AssetID: %s", assetID)

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Fetch Collateral Ratio from Ledger
    ratio, err := ledgerInstance.DeFiLedger.GetCollateralRatio(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch collateral ratio for AssetID %s: %v", assetID, err)
        return 0, fmt.Errorf("failed to fetch collateral ratio for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Fetched collateral ratio of %f for synthetic asset %s", ratio, assetID)
    return ratio, nil
}


// SyntheticAssetSetCollateralRatio sets a new collateral ratio for a synthetic asset.
func SyntheticAssetSetCollateralRatio(assetID string, newRatio float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Setting new collateral ratio for AssetID: %s to %f", assetID, newRatio)

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if newRatio <= 0 {
        err := fmt.Errorf("collateral ratio must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Update Collateral Ratio in Ledger
    err := ledgerInstance.DeFiLedger.SetCollateralRatio(assetID, newRatio)
    if err != nil {
        log.Printf("[ERROR] Failed to set collateral ratio for AssetID %s: %v", assetID, err)
        return fmt.Errorf("failed to set collateral ratio for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Collateral ratio set to %f for synthetic asset %s", newRatio, assetID)
    return nil
}


// SyntheticAssetTrackVolatility enables volatility tracking for a synthetic asset.
func SyntheticAssetTrackVolatility(assetID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating volatility tracking for AssetID: %s", assetID)

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Enable Volatility Tracking in Ledger
    err := ledgerInstance.DeFiLedger.TrackAssetVolatility(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to enable volatility tracking for AssetID %s: %v", assetID, err)
        return fmt.Errorf("failed to enable volatility tracking for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Volatility tracking enabled for synthetic asset %s", assetID)
    return nil
}


// SyntheticAssetFetchVolatility retrieves the current volatility of a synthetic asset.
func SyntheticAssetFetchVolatility(assetID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Fetching volatility for AssetID: %s", assetID)

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Fetch Volatility from Ledger
    volatility, err := ledgerInstance.DeFiLedger.GetAssetVolatility(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch volatility for AssetID %s: %v", assetID, err)
        return 0, fmt.Errorf("failed to fetch volatility for asset %s: %w", assetID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Fetched volatility of %f for synthetic asset %s", volatility, assetID)
    return volatility, nil
}


// SyntheticAssetAutoAdjustCollateral automatically adjusts collateral for a synthetic asset based on volatility.
func SyntheticAssetAutoAdjustCollateral(assetID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Starting collateral auto-adjustment for AssetID: %s", assetID)

    // Step 1: Validate Input
    if assetID == "" {
        err := fmt.Errorf("assetID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Fetch Current Volatility (Pre-Adjustment Validation)
    volatility, err := ledgerInstance.DeFiLedger.GetAssetVolatility(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch volatility for AssetID %s: %v", assetID, err)
        return fmt.Errorf("failed to fetch volatility for asset %s: %w", assetID, err)
    }
    log.Printf("[INFO] Fetched current volatility for AssetID %s: %f", assetID, volatility)

    // Step 3: Perform Auto Adjustment in Ledger
    err = ledgerInstance.DeFiLedger.AutoAdjustCollateral(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to auto-adjust collateral for AssetID %s: %v", assetID, err)
        return fmt.Errorf("failed to auto-adjust collateral for asset %s: %w", assetID, err)
    }

    // Step 4: Fetch New Collateral Details (Post-Adjustment Verification)
    newCollateralRatio, err := ledgerInstance.DeFiLedger.GetCollateralRatio(assetID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch new collateral ratio for AssetID %s after adjustment: %v", assetID, err)
        return fmt.Errorf("failed to fetch new collateral ratio for asset %s: %w", assetID, err)
    }
    log.Printf("[INFO] New collateral ratio for AssetID %s: %f", assetID, newCollateralRatio)

    // Step 5: Log Success
    log.Printf("[SUCCESS] Collateral auto-adjusted for synthetic asset %s based on volatility. New Collateral Ratio: %f", assetID, newCollateralRatio)

    return nil
}

