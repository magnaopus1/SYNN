package syn1600

import (
	"errors"
	"time"
)

// TokenManagement handles operations related to the SYN1600 token standard.
type TokenManagement struct {
	Ledger ledger.Ledger
}

// AddOwnershipRecord adds a new ownership record to the SYN1600 token.
func (m *TokenManagement) AddOwnershipRecord(tokenID string, newOwner string, fractionalShare float64) error {
	token, err := m.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	// Add new ownership record
	token.(*common.SYN1600Token).OwnershipRights = append(token.(*common.SYN1600Token).OwnershipRights, common.OwnershipRecord{
		OwnerID:        newOwner,
		FractionalShare: fractionalShare,
		PurchaseDate:   time.Now(),
	})

	// Log the ownership change in the history
	token.(*common.SYN1600Token).OwnershipHistory = append(token.(*common.SYN1600Token).OwnershipHistory, common.OwnershipHistory{
		HistoryID:       generateUniqueID(),
		PreviousOwnerID: token.(*common.SYN1600Token).Owner,
		NewOwnerID:      newOwner,
		TransferDate:    time.Now(),
		TransferTerms:   "New ownership acquired",
	})

	// Update the ledger with the new ownership information
	return m.Ledger.UpdateToken(tokenID, token)
}

// RemoveOwnershipRecord removes a fractional share of ownership from a SYN1600 token.
func (m *TokenManagement) RemoveOwnershipRecord(tokenID string, ownerID string, fractionalShare float64) error {
	token, err := m.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	for i, record := range token.(*common.SYN1600Token).OwnershipRights {
		if record.OwnerID == ownerID {
			if record.FractionalShare < fractionalShare {
				return errors.New("insufficient fractional ownership to remove")
			}

			// Update the fractional ownership
			token.(*common.SYN1600Token).OwnershipRights[i].FractionalShare -= fractionalShare

			// If ownership share reaches zero, remove the ownership record
			if token.(*common.SYN1600Token).OwnershipRights[i].FractionalShare == 0 {
				token.(*common.SYN1600Token).OwnershipRights = append(token.(*common.SYN1600Token).OwnershipRights[:i], token.(*common.SYN1600Token).OwnershipRights[i+1:]...)
			}

			break
		}
	}

	// Update the ledger after modification
	return m.Ledger.UpdateToken(tokenID, token)
}

// UpdateRevenueStream updates the revenue stream associated with the SYN1600 token.
func (m *TokenManagement) UpdateRevenueStream(tokenID string, streamID string, updatedStream common.RevenueStream) error {
	token, err := m.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	for i, stream := range token.(*common.SYN1600Token).FutureRevenueStreams {
		if stream.RevenueID == streamID {
			token.(*common.SYN1600Token).FutureRevenueStreams[i] = updatedStream
			break
		}
	}

	// Update the ledger with the modified revenue stream
	return m.Ledger.UpdateToken(tokenID, token)
}

// DistributeRoyalties processes and distributes royalties to the owners of the SYN1600 token.
func (m *TokenManagement) DistributeRoyalties(tokenID string, totalRoyalty float64) error {
	token, err := m.Ledger.GetToken(tokenID)
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

	// Store distribution logs in the ledger
	return m.Ledger.UpdateToken(tokenID, token)
}

// EncryptSensitiveMetadata encrypts sensitive metadata related to the SYN1600 token.
func (m *TokenManagement) EncryptSensitiveMetadata(tokenID string, metadata []byte, key []byte) error {
	encryptedData, err := encryptData(metadata, key)
	if err != nil {
		return err
	}

	token, err := m.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	token.(*common.SYN1600Token).EncryptedMetadata = encryptedData

	// Update ledger with the encrypted metadata
	return m.Ledger.UpdateToken(tokenID, token)
}

// DecryptSensitiveMetadata decrypts sensitive metadata related to the SYN1600 token.
func (m *TokenManagement) DecryptSensitiveMetadata(tokenID string, key []byte) ([]byte, error) {
	token, err := m.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, err
	}

	return decryptData(token.(*common.SYN1600Token).EncryptedMetadata, key)
}

// Helper function to encrypt data using AES
func encryptData(data []byte, key []byte) ([]byte, error) {
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

// Helper function to decrypt data using AES
func decryptData(encryptedData []byte, key []byte) ([]byte, error) {
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

// ValidateTransactionForSynnergy validates the transactions as part of the Synnergy Consensus, splitting into sub-blocks.
func (m *TokenManagement) ValidateTransactionForSynnergy(transactionID string) error {
	// Example implementation of how a transaction could be validated with Synnergy Consensus
	valid, err := synnergy.ValidateSubBlockTransaction(transactionID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("transaction validation failed")
	}

	// Record validation result in the ledger
	err = m.Ledger.RecordTransactionValidation(transactionID)
	if err != nil {
		return err
	}

	return nil
}
