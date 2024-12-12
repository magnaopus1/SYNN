package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// SYN721Deployment handles the deployment of SYN721 tokens (NFTs).
type SYN721Deployment struct {
	mutex      sync.Mutex                 // For thread-safe operations
	Ledger     *ledger.Ledger             // Reference to the ledger for recording deployments
	Consensus  *synnergy_consensus.Engine // Synnergy Consensus engine for validation
	Encryption *encryption.Encryption     // Encryption service for securing deployment data
	Storage    *SYN721Storage             // Storage for SYN721 token data
}

// NewSYN721Deployment initializes a new SYN721 deployment manager.
func NewSYN721Deployment(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, storage *SYN721Storage) *SYN721Deployment {
	return &SYN721Deployment{
		Ledger:     ledgerInstance,
		Consensus:  consensusEngine,
		Encryption: encryptionService,
		Storage:    storage,
	}
}

// DeployToken deploys a new SYN721 token to the blockchain.
func (sd *SYN721Deployment) DeployToken(tokenID, tokenURI, owner string) (*SYN721Token, error) {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()

	// Check if the token ID already exists in storage
	if _, exists := sd.Storage.GetTokenData(tokenID); exists {
		return nil, errors.New("token ID already exists")
	}

	// Create the token
	token := &SYN721Token{
		TokenID:    tokenID,
		TokenURI:   tokenURI,
		Owner:      owner,
		Ledger:     sd.Ledger,
		Consensus:  sd.Consensus,
		Encryption: sd.Encryption,
	}

	// Encrypt the token data
	encryptedData, err := sd.Encryption.EncryptData(fmt.Sprintf("%v", token), "")
	if err != nil {
		return nil, fmt.Errorf("error encrypting token data: %v", err)
	}

	// Validate the token deployment using Synnergy Consensus
	if valid, err := sd.Consensus.ValidateTokenDeployment(tokenID, owner, encryptedData); !valid || err != nil {
		return nil, fmt.Errorf("token deployment failed consensus validation: %v", err)
	}

	// Store the token in storage
	err = sd.Storage.StoreTokenData(tokenID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to store token data: %v", err)
	}

	// Record the deployment in the ledger
	err = sd.Ledger.RecordTokenDeployment(tokenID, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to log token deployment in ledger: %v", err)
	}

	fmt.Printf("Token %s successfully deployed for owner %s.\n", tokenID, owner)
	return token, nil
}

// RedeployToken redeploys an existing SYN721 token, updating its metadata if necessary.
func (sd *SYN721Deployment) RedeployToken(tokenID, tokenURI, owner string) (*SYN721Token, error) {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()

	// Retrieve the existing token from storage
	token, exists := sd.Storage.GetTokenData(tokenID)
	if !exists {
		return nil, errors.New("token not found")
	}

	// Check if the requester is the owner of the token
	if token.Owner != owner {
		return nil, errors.New("only the owner can redeploy the token")
	}

	// Update the token URI if provided
	if tokenURI != "" {
		token.TokenURI = tokenURI
	}

	// Encrypt the updated token data
	encryptedData, err := sd.Encryption.EncryptData(fmt.Sprintf("%v", token), "")
	if err != nil {
		return nil, fmt.Errorf("error encrypting updated token data: %v", err)
	}

	// Validate the token redeployment using Synnergy Consensus
	if valid, err := sd.Consensus.ValidateTokenDeployment(tokenID, owner, encryptedData); !valid || err != nil {
		return nil, fmt.Errorf("token redeployment failed consensus validation: %v", err)
	}

	// Update the token data in storage
	err = sd.Storage.UpdateTokenData(tokenID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to update token data: %v", err)
	}

	// Record the redeployment in the ledger
	err = sd.Ledger.RecordTokenRedeployment(tokenID, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to log token redeployment in ledger: %v", err)
	}

	fmt.Printf("Token %s successfully redeployed for owner %s.\n", tokenID, owner)
	return token, nil
}

// ValidateDeployment ensures that the deployment or redeployment meets all necessary criteria.
func (sd *SYN721Deployment) ValidateDeployment(tokenID, owner string) error {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()

	// Retrieve the token from storage
	token, exists := sd.Storage.GetTokenData(tokenID)
	if !exists {
		return errors.New("token not found")
	}

	// Validate ownership
	if token.Owner != owner {
		return errors.New("token ownership validation failed")
	}

	// Validate token deployment using consensus
	valid, err := sd.Consensus.ValidateTokenDeployment(tokenID, owner, "")
	if !valid || err != nil {
		return fmt.Errorf("token deployment validation failed: %v", err)
	}

	fmt.Printf("Token deployment for %s validated successfully.\n", tokenID)
	return nil
}
