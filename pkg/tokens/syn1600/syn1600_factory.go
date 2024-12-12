package syn1600

import (
	"errors"
	"time"
)

// TokenFactory is responsible for creating, managing, and interacting with SYN1600 tokens.
type TokenFactory struct {
	Ledger ledger.Ledger // Reference to the ledger for storing tokens
}

// CreateSYN1600Token creates a new SYN1600Token for a music royalty asset.
func (f *TokenFactory) CreateSYN1600Token(owner string, metadata common.MusicAssetMetadata, fractionalShare float64, forecast common.RevenueForecast, encryptedMetadata []byte) (*common.SYN1600Token, error) {
	if owner == "" || fractionalShare <= 0 {
		return nil, errors.New("invalid owner or fractional share")
	}

	token := &SYN1600Token{
		TokenID:            generateUniqueID(),
		Owner:              owner,
		MusicAssetMetadata: metadata,
		OwnershipRights: []common.OwnershipRecord{
			{
				OwnerID:        owner,
				FractionalShare: fractionalShare,
				PurchaseDate:   time.Now(),
			},
		},
		FutureRevenueStreams:  []common.RevenueStream{},
		RevenueDistribution:   []common.RoyaltyDistributionLog{},
		RevenueTracking:       []common.RevenueRecord{},
		CustomRoyaltySplits:   make(map[string]float64),
		ComplianceStatus:      "Compliant",
		AuditTrail:            []common.AuditLog{},
		OwnershipHistory:      []common.OwnershipHistory{},
		ImmutableRecords:      []common.ImmutableRecord{},
		RestrictedTransfers:   false,
		ApprovalRequired:      false,
		RevenueForecast:       forecast,
		EncryptedMetadata:     encryptedMetadata,
	}

	// Store the token in the ledger
	if err := f.Ledger.StoreToken(token.TokenID, token); err != nil {
		return nil, err
	}

	return token, nil
}

// AddRevenueStream adds a new revenue stream to the SYN1600 token.
func (f *TokenFactory) AddRevenueStream(tokenID string, stream common.RevenueStream) error {
	token, err := f.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	token.(*common.SYN1600Token).FutureRevenueStreams = append(token.(*common.SYN1600Token).FutureRevenueStreams, stream)

	// Update ledger
	return f.Ledger.UpdateToken(tokenID, token)
}

// DistributeRoyalties distributes the royalties among the rights holders based on their fractional ownership.
func (f *TokenFactory) DistributeRoyalties(tokenID string, totalRoyalty float64) error {
	token, err := f.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	for _, owner := range token.(*common.SYN1600Token).OwnershipRights {
		amount := totalRoyalty * owner.FractionalShare
		token.(*common.SYN1600Token).RevenueDistribution = append(token.(*common.SYN1600Token).RevenueDistribution, common.RoyaltyDistributionLog{
			DistributionID:   generateUniqueID(),
			RecipientID:      owner.OwnerID,
			Amount:           amount,
			DistributionDate: time.Now(),
			DistributionType: "Automatic",
		})
	}

	// Update ledger
	return f.Ledger.UpdateToken(tokenID, token)
}

// TransferOwnership transfers ownership of the SYN1600 token to another party.
func (f *TokenFactory) TransferOwnership(tokenID, newOwner string, fractionalShare float64) error {
	token, err := f.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	// Update ownership rights
	for i, record := range token.(*common.SYN1600Token).OwnershipRights {
		if record.OwnerID == newOwner {
			token.(*common.SYN1600Token).OwnershipRights[i].FractionalShare += fractionalShare
		} else {
			token.(*common.SYN1600Token).OwnershipRights = append(token.(*common.SYN1600Token).OwnershipRights, common.OwnershipRecord{
				OwnerID:        newOwner,
				FractionalShare: fractionalShare,
				PurchaseDate:   time.Now(),
			})
		}
	}

	// Log ownership history
	token.(*common.SYN1600Token).OwnershipHistory = append(token.(*common.SYN1600Token).OwnershipHistory, common.OwnershipHistory{
		HistoryID:       generateUniqueID(),
		PreviousOwnerID: token.(*common.SYN1600Token).Owner,
		NewOwnerID:      newOwner,
		TransferDate:    time.Now(),
		TransferTerms:   "Transfer of ownership",
	})

	// Update ledger
	return f.Ledger.UpdateToken(tokenID, token)
}

// EncryptMetadata encrypts the metadata using AES encryption.
func EncryptMetadata(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return []byte(base64.URLEncoding.EncodeToString(ciphertext)), nil
}

// DecryptMetadata decrypts the metadata using AES encryption.
func DecryptMetadata(encryptedData []byte, key []byte) ([]byte, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(string(encryptedData))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// Helper function to generate a unique ID (placeholder)
func generateUniqueID() string {
	return "UNIQUE_ID_" + time.Now().Format("20060102150405")
}
