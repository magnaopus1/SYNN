package syn3400

import (
	"errors"
	"sync"
	"time"

)

// Syn3400Token represents the SYN3400 Forex Pair Token standard.
type Syn3400Token struct {
	TokenID        string
	ForexPair      ForexPair
	Owner          string
	PositionSize   float64
	OpenRate       float64
	LongShort      string
	OpenedDate     time.Time
	LastUpdated    time.Time
	TransactionIDs []string
}

// ForexPair represents the structure of a forex trading pair.
type ForexPair struct {
	PairID        string    `json:"pair_id"`
	BaseCurrency  string    `json:"base_currency"`
	QuoteCurrency string    `json:"quote_currency"`
	CurrentRate   float64   `json:"current_rate"`
	LastUpdated   time.Time `json:"last_updated"`
}

// ForexMetadata represents the metadata associated with forex pairs.
type ForexMetadata struct {
	Pairs         map[string]ForexPair
	mutex         sync.Mutex
	encryptionKey []byte
}

// PairLinking manages the linking between tokens and forex pairs.
type PairLinking struct {
	PairLinks map[string]PairLink
	mutex     sync.Mutex
}

// PairLink represents a link between a token and a forex pair.
type PairLink struct {
	TokenID        string    `json:"token_id"`
	ForexPairID    string    `json:"forex_pair_id"`
	Linked         bool      `json:"linked"`
	LastLinkedTime time.Time `json:"last_linked_time"`
}

// TokenFactory handles the creation and management of SYN3400 tokens.
type TokenFactory struct {
	mu          sync.Mutex
	tokens      map[string]*Syn3400Token
	metadata    *ForexMetadata
	ledger      *ledger.Ledger
	encryptor   *encryption.Encryptor
	consensus   *consensus.SynnergyConsensus
}

// NewTokenFactory creates a new instance of TokenFactory.
func NewTokenFactory(metadata *ForexMetadata, ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *TokenFactory {
	return &TokenFactory{
		tokens:    make(map[string]*Syn3400Token),
		metadata:  metadata,
		ledger:    ledger,
		encryptor: encryptor,
		consensus: consensus,
	}
}

// CreateToken creates a new SYN3400Token and stores it.
func (tf *TokenFactory) CreateToken(owner string, pair ForexPair, positionSize, openRate float64, longShort string) (*Syn3400Token, error) {
	tf.mu.Lock()
	defer tf.mu.Unlock()

	// Validate input
	if owner == "" || pair.PairID == "" || positionSize <= 0 || openRate <= 0 || (longShort != "long" && longShort != "short") {
		return nil, errors.New("invalid token parameters")
	}

	tokenID := generateUniqueTokenID()
	token := &Syn3400Token{
		TokenID:      tokenID,
		ForexPair:    pair,
		Owner:        owner,
		PositionSize: positionSize,
		OpenRate:     openRate,
		LongShort:    longShort,
		OpenedDate:   time.Now(),
		LastUpdated:  time.Now(),
	}

	// Encrypt the token data before storing it.
	encryptedToken, err := tf.encryptor.EncryptData(token)
	if err != nil {
		return nil, err
	}

	// Store the encrypted token in memory.
	tf.tokens[tokenID] = encryptedToken.(*Syn3400Token)

	// Log the creation event in the ledger.
	tf.ledger.LogTokenEvent("TokenCreated", time.Now(), tokenID)

	// Validate the token creation using consensus.
	err = tf.consensus.ValidateSubBlock(tokenID)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetToken retrieves a SYN3400Token by its TokenID.
func (tf *TokenFactory) GetToken(tokenID string) (*Syn3400Token, error) {
	tf.mu.Lock()
	defer tf.mu.Unlock()

	encryptedToken, exists := tf.tokens[tokenID]
	if !exists {
		return nil, errors.New("token not found")
	}

	// Decrypt the token data before returning it.
	token, err := tf.encryptor.DecryptData(encryptedToken)
	if err != nil {
		return nil, err
	}

	return token.(*Syn3400Token), nil
}

// UpdateToken updates an existing SYN3400Token.
func (tf *TokenFactory) UpdateToken(token *Syn3400Token) error {
	tf.mu.Lock()
	defer tf.mu.Unlock()

	if _, exists := tf.tokens[token.TokenID]; !exists {
		return errors.New("token not found")
	}

	// Encrypt the token data before updating it.
	encryptedToken, err := tf.encryptor.EncryptData(token)
	if err != nil {
		return err
	}

	tf.tokens[token.TokenID] = encryptedToken.(*Syn3400Token)

	// Log the update event in the ledger.
	tf.ledger.LogTokenEvent("TokenUpdated", time.Now(), token.TokenID)

	// Sync the token update with consensus.
	return tf.consensus.ValidateSubBlock(token.TokenID)
}

// LinkTokenToPair links a SYN3400Token to a specific ForexPair.
func (tf *TokenFactory) LinkTokenToPair(tokenID, forexPairID string) error {
	tf.mu.Lock()
	defer tf.mu.Unlock()

	token, exists := tf.tokens[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	pair, exists := tf.metadata.Pairs[forexPairID]
	if !exists {
		return errors.New("forex pair not found")
	}

	token.ForexPair = pair
	token.LastUpdated = time.Now()

	// Encrypt and update the token.
	encryptedToken, err := tf.encryptor.EncryptData(token)
	if err != nil {
		return err
	}
	tf.tokens[tokenID] = encryptedToken.(*Syn3400Token)

	// Log the linking in the ledger.
	tf.ledger.LogTokenEvent("TokenLinkedToPair", time.Now(), tokenID)

	return tf.consensus.ValidateSubBlock(tokenID)
}

// generateUniqueTokenID generates a unique identifier for SYN3400 tokens.
func generateUniqueTokenID() string {
	return "TKN" + time.Now().Format("20060102150405")
}
