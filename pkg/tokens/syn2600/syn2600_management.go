package syn2600

import (

	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

)

// SYN2600TokenManagement provides management functions for SYN2600 tokens.
type SYN2600TokenManagement struct {
	TokenID      string
	AssetDetails string
	Owner        string
	Shares       float64
	IssuedDate   time.Time
	ExpiryDate   time.Time
	Active       bool
	Signature    string // Encryption signature for verification
}

// CreateNewInvestorToken handles the creation of a new SYN2600 investor token, integrates with the ledger, and performs encryption.
func CreateNewInvestorToken(assetDetails string, owner string, shares float64, expiryDate time.Time) (string, error) {
	tokenID := generateTokenID(assetDetails, owner)
	issuedDate := time.Now()

	// Create new token structure
	newToken := SYN2600TokenManagement{
		TokenID:      tokenID,
		AssetDetails: assetDetails,
		Owner:        owner,
		Shares:       shares,
		IssuedDate:   issuedDate,
		ExpiryDate:   expiryDate,
		Active:       true,
		Signature:    generateTokenSignature(tokenID, assetDetails, owner),
	}

	// Encrypt token data before storing
	encryptedToken, err := encryption.EncryptTokenData(&newToken)
	if err != nil {
		return "", errors.New("failed to encrypt token data before storing")
	}

	// Store encrypted token in the ledger
	err = ledger.StoreInvestorToken(encryptedToken)
	if err != nil {
		return "", errors.New("failed to store the investor token in the ledger")
	}

	// Validate the new token via Synnergy Consensus and validate sub-blocks
	err = synconsensus.ValidateSubBlock(tokenID)
	if err != nil {
		return "", errors.New("token validation failed in Synnergy Consensus")
	}

	return tokenID, nil
}

// TransferOwnership handles the transfer of ownership of SYN2600 tokens.
func TransferOwnership(tokenID string, newOwner string) (string, error) {
	// Fetch the token from the ledger
	token, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return "", errors.New("failed to fetch token for ownership transfer")
	}

	// Decrypt token before transferring ownership
	decryptedToken, err := encryption.DecryptTokenData(token)
	if err != nil {
		return "", errors.New("failed to decrypt token data for ownership transfer")
	}

	// Update owner details
	decryptedToken.Owner = newOwner

	// Re-encrypt token and update in the ledger
	encryptedToken, err := encryption.EncryptTokenData(decryptedToken)
	if err != nil {
		return "", errors.New("failed to encrypt token data after ownership transfer")
	}
	err = ledger.UpdateInvestorToken(encryptedToken)
	if err != nil {
		return "", errors.New("failed to update token ownership in the ledger")
	}

	// Record the transfer event
	eventID, err := RecordEvent(tokenID, "TRANSFER", "Ownership transferred to "+newOwner, newOwner)
	if err != nil {
		return "", errors.New("failed to record transfer event")
	}

	return eventID, nil
}

// RedeemInvestorToken allows for the redemption of an investor token when it reaches its expiry date.
func RedeemInvestorToken(tokenID string) (string, error) {
	// Fetch token from the ledger
	token, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return "", errors.New("failed to fetch token for redemption")
	}

	// Decrypt token data
	decryptedToken, err := encryption.DecryptTokenData(token)
	if err != nil {
		return "", errors.New("failed to decrypt token data for redemption")
	}

	// Check if the token has expired
	if time.Now().After(decryptedToken.ExpiryDate) {
		decryptedToken.Active = false

		// Re-encrypt token and update in the ledger
		encryptedToken, err := encryption.EncryptTokenData(decryptedToken)
		if err != nil {
			return "", errors.New("failed to encrypt token data after redemption")
		}
		err = ledger.UpdateInvestorToken(encryptedToken)
		if err != nil {
			return "", errors.New("failed to update token status in the ledger")
		}

		// Record the redemption event
		eventID, err := RecordEvent(tokenID, "REDEMPTION", "Token redeemed and deactivated", decryptedToken.Owner)
		if err != nil {
			return "", errors.New("failed to record redemption event")
		}

		return eventID, nil
	}

	return "", errors.New("token has not yet expired")
}

// FetchTokenDetails retrieves the full details of a SYN2600 token.
func FetchTokenDetails(tokenID string) (*SYN2600TokenManagement, error) {
	// Fetch encrypted token from the ledger
	encryptedToken, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to fetch token details from the ledger")
	}

	// Decrypt token data
	decryptedToken, err := encryption.DecryptTokenData(encryptedToken)
	if err != nil {
		return nil, errors.New("failed to decrypt token data")
	}

	return decryptedToken, nil
}

// RecordReturns handles the process of recording returns or dividends for a specific token.
func RecordReturns(tokenID string, returnAmount float64) (string, error) {
	// Fetch token from the ledger
	token, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return "", errors.New("failed to fetch token for return recording")
	}

	// Decrypt token data
	decryptedToken, err := encryption.DecryptTokenData(token)
	if err != nil {
		return "", errors.New("failed to decrypt token data for return recording")
	}

	// Record return event
	returnDetails := "Return of " + common.FloatToString(returnAmount) + " for owner " + decryptedToken.Owner
	eventID, err := RecordEvent(tokenID, "RETURN_UPDATE", returnDetails, decryptedToken.Owner)
	if err != nil {
		return "", errors.New("failed to record return update")
	}

	return eventID, nil
}

// generateTokenID generates a unique token ID based on the asset details and owner information.
func generateTokenID(assetDetails string, owner string) string {
	hashInput := assetDetails + owner + time.Now().String()
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// generateTokenSignature generates a signature for each token to ensure authenticity.
func generateTokenSignature(tokenID string, assetDetails string, owner string) string {
	signatureInput := tokenID + assetDetails + owner
	hash := sha256.Sum256([]byte(signatureInput))
	return hex.EncodeToString(hash[:])
}
