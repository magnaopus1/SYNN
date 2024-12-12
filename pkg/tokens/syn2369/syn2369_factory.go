package syn2369

import (
	"time"
	"errors"
)

// SYN2369Token represents a virtual world item or property under the SYN2369 Token Standard.
type SYN2369Token struct {
	TokenID            string            // Unique identifier for the virtual item
	ItemName           string            // Name of the virtual item (e.g., "Sword of Valor")
	ItemType           string            // Type of the virtual item (e.g., "Weapon", "Real Estate", "3D Model")
	Description        string            // Detailed description of the virtual item
	Attributes         map[string]string // Custom attributes such as color, size, power, etc.
	Owner              string            // Current owner of the item
	Creator            string            // Original creator of the virtual item
	CreatedAt          time.Time         // Creation timestamp of the virtual item
	UpdatedAt          time.Time         // Timestamp of the last update
	Customizable       bool              // Indicates if the item can be customized
	Locked             bool              // Indicates if the item is locked from transfer or modification
	MultiSigRequired   bool              // Indicates if multi-signature approval is required for transactions
	EncryptedMetadata  []byte            // Encrypted metadata for secure information storage
	OwnershipHistory   []OwnershipRecord // Historical ownership records
	TransactionHistory []TransactionLog  // Logs of all transactions for this virtual item
	EventLogs          []EventLog        // Logs of all significant events related to this item
	OffChainMetadata   OffChainStorage   // References to 3D models and other media stored off-chain
}

// OffChainStorage represents off-chain storage metadata such as links to 3D models, textures, and multimedia.
type OffChainStorage struct {
	ModelURL        string    // URL or IPFS hash to the 3D model (e.g., .glb, .obj)
	TextureURL      string    // URL or IPFS hash to textures related to the model (e.g., .png, .jpg)
AdditionalFiles []string    // List of additional assets (e.g., audio, shaders) stored off-chain
	Hash           string    // Hash of the off-chain content for integrity verification
	LastVerified   time.Time // Last time the off-chain data was verified
}

// OwnershipRecord represents the ownership history of the virtual item.
type OwnershipRecord struct {
	PreviousOwner string    // Previous owner of the item
	NewOwner      string    // New owner after transfer
	TransferDate  time.Time // Date of ownership transfer
}

// TransactionLog represents a transaction made with the virtual item.
type TransactionLog struct {
	TransactionID   string    // Unique identifier for the transaction
	TransactionType string    // Type of transaction (e.g., "Transfer", "Auction", "Sale")
	Sender          string    // Sender of the virtual item in the transaction
	Recipient       string    // Recipient of the virtual item in the transaction
	Amount          float64   // Amount associated with the transaction (if applicable)
	Timestamp       time.Time // Timestamp of the transaction
}

// EventLog represents a significant event related to the virtual item.
type EventLog struct {
	EventID     string    // Unique identifier for the event
	EventType   string    // Type of event (e.g., "Customization", "Attribute Update")
	Description string    // Description of the event
	Timestamp   time.Time // Timestamp of the event
}

// AddCustomAttribute adds a custom attribute to the virtual item.
func (token *SYN2369Token) AddCustomAttribute(key, value string) error {
	if token.Customizable {
		token.Attributes[key] = value
		token.UpdatedAt = time.Now()
		return nil
	}
	return errors.New("item customization is not allowed")
}

// UpdateAttributes updates the attributes of the virtual item.
func (token *SYN2369Token) UpdateAttributes(attributes map[string]string) error {
	if token.Customizable {
		for key, value := range attributes {
			token.Attributes[key] = value
		}
		token.UpdatedAt = time.Now()
		return nil
	}
	return errors.New("item customization is not allowed")
}

// TransferOwnership transfers the ownership of the virtual item to a new owner.
func (token *SYN2369Token) TransferOwnership(newOwner string) error {
	if token.Locked {
		return errors.New("item is locked and cannot be transferred")
	}
	token.OwnershipHistory = append(token.OwnershipHistory, OwnershipRecord{
		PreviousOwner: token.Owner,
		NewOwner:      newOwner,
		TransferDate:  time.Now(),
	})
	token.Owner = newOwner
	token.UpdatedAt = time.Now()
	return nil
}

// LockItem locks the virtual item from being transferred or modified.
func (token *SYN2369Token) LockItem() {
	token.Locked = true
	token.UpdatedAt = time.Now()
}

// UnlockItem unlocks the virtual item, allowing transfers and modifications.
func (token *SYN2369Token) UnlockItem() {
	token.Locked = false
	token.UpdatedAt = time.Now()
}

// LogTransaction records a transaction in the token's transaction history.
func (token *SYN2369Token) LogTransaction(transactionType, sender, recipient string, amount float64) {
	token.TransactionHistory = append(token.TransactionHistory, TransactionLog{
		TransactionID:   generateUniqueID(),
		TransactionType: transactionType,
		Sender:          sender,
		Recipient:       recipient,
		Amount:          amount,
		Timestamp:       time.Now(),
	})
}

// LogEvent records a significant event in the token's event logs.
func (token *SYN2369Token) LogEvent(eventType, description string) {
	token.EventLogs = append(token.EventLogs, EventLog{
		EventID:     generateUniqueID(),
		EventType:   eventType,
		Description: description,
		Timestamp:   time.Now(),
	})
}

// EncryptMetadata encrypts the token's sensitive metadata.
func (token *SYN2369Token) EncryptMetadata() error {
	encryptedData, err := encryption.EncryptData([]byte(token.ItemID + token.Owner + token.Creator))
	if err != nil {
		return err
	}
	token.EncryptedMetadata = encryptedData
	return nil
}

// DecryptMetadata decrypts the token's encrypted metadata.
func (token *SYN2369Token) DecryptMetadata() (string, error) {
	decryptedData, err := encryption.DecryptData(token.EncryptedMetadata)
	if err != nil {
		return "", err
	}
	return string(decryptedData), nil
}

// AddOffChainModel adds off-chain storage references (e.g., 3D model) for the token.
func (token *SYN2369Token) AddOffChainModel(modelURL, textureURL, hash string, additionalFiles []string) {
	token.OffChainMetadata = OffChainStorage{
		ModelURL:        modelURL,
		TextureURL:      textureURL,
		Hash:            hash,
		AdditionalFiles: additionalFiles,
		LastVerified:    time.Now(),
	}
	token.UpdatedAt = time.Now()
}

// VerifyOffChainData verifies the off-chain data using its hash.
func (token *SYN2369Token) VerifyOffChainData(currentHash string) bool {
	return token.OffChainMetadata.Hash == currentHash
}

// CreateTokenFactory creates a new SYN2369Token with the provided attributes.
func CreateTokenFactory(
	itemName, itemType, description string,
	attributes map[string]string,
	creator, owner string,
	customizable, multiSigRequired bool,
) (common.SYN2369Token, error) {

	// Generate a unique TokenID for the new token
	tokenID, err := generateUniqueTokenID()
	if err != nil {
		return common.SYN2369Token{}, err
	}

	// Initialize the token's base structure
	token := common.SYN2369Token{
		TokenID:          tokenID,
		ItemName:         itemName,
		ItemType:         itemType,
		Description:      description,
		Attributes:       attributes,
		Owner:            owner,
		Creator:          creator,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Customizable:     customizable,
		MultiSigRequired: multiSigRequired,
		Locked:           false,
	}

	// Log the creation event in the token's event logs
	token.LogEvent("Token Created", "Initial creation of the SYN2369 token")

	// Add the new token to the ledger
	err = ledger.AddTokenToLedger(token)
	if err != nil {
		return token, err
	}

	// Encrypt the token's metadata for security
	err = token.EncryptMetadata()
	if err != nil {
		return token, err
	}

	return token, nil
}

// UpdateTokenAttributes updates the attributes of an existing SYN2369Token.
func UpdateTokenAttributes(token *common.SYN2369Token, attributes map[string]string) error {
	// Ensure the item is customizable
	if !token.Customizable {
		return errors.New("this item is not customizable")
	}

	// Update the attributes
	err := token.UpdateAttributes(attributes)
	if err != nil {
		return err
	}

	// Log the attribute update event
	token.LogEvent("Attributes Updated", "Item attributes were updated")

	// Sync the update with the ledger
	err = ledger.UpdateTokenInLedger(*token)
	if err != nil {
		return err
	}

	return nil
}

// TransferTokenOwnership securely transfers the ownership of a SYN2369Token.
func TransferTokenOwnership(token *common.SYN2369Token, newOwner string) error {
	// Ensure the item is not locked for transfer
	if token.Locked {
		return errors.New("item is locked and cannot be transferred")
	}

	// Transfer ownership
	err := token.TransferOwnership(newOwner)
	if err != nil {
		return err
	}

	// Log the ownership transfer event
	token.LogEvent("Ownership Transferred", "Ownership transferred to "+newOwner)

	// Sync the transfer with the ledger
	err = ledger.UpdateTokenInLedger(*token)
	if err != nil {
		return err
	}

	return nil
}

// AddOffChainAssets securely adds off-chain storage references (e.g., 3D models) to the token.
func AddOffChainAssets(token *common.SYN2369Token, modelURL, textureURL, hash string, additionalFiles []string) error {
	// Add the off-chain assets
	token.AddOffChainModel(modelURL, textureURL, hash, additionalFiles)

	// Log the off-chain asset addition
	token.LogEvent("Off-Chain Assets Added", "3D models and textures added to the token")

	// Sync the update with the ledger
	err := ledger.UpdateTokenInLedger(*token)
	if err != nil {
		return err
	}

	return nil
}

// GenerateUniqueTokenID generates a unique token ID.
func generateUniqueTokenID() (string, error) {
	id := make([]byte, 16)
	_, err := rand.Read(id)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(id), nil
}

// MintNewTokens mints new tokens for virtual world assets.
func MintNewTokens(items []common.SYN2369Token) error {
	for _, item := range items {
		// Encrypt metadata for security
		err := item.EncryptMetadata()
		if err != nil {
			return err
		}

		// Log the minting event
		item.LogEvent("Token Minted", "New SYN2369 token minted for " + item.ItemName)

		// Add each item to the ledger
		err = ledger.AddTokenToLedger(item)
		if err != nil {
			return err
		}
	}
	return nil
}

// BurnToken burns a token by removing it from the ledger.
func BurnToken(tokenID string) error {
	// Retrieve the token from the ledger
	token, err := ledger.GetTokenFromLedger(tokenID)
	if err != nil {
		return err
	}

	// Log the burning event
	token.LogEvent("Token Burned", "Token " + tokenID + " has been burned")

	// Remove the token from the ledger
	err = ledger.RemoveTokenFromLedger(tokenID)
	if err != nil {
		return err
	}

	return nil
}
