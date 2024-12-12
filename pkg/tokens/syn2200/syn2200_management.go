package syn2200

import (
	"errors"
	"time"

)

// IssueToken manages the issuance of a new SYN2200 token.
func IssueToken(tokenID string, currency string, amount float64, senderID string, recipientID string) (*common.SYN2200Token, error) {
	// Validate input parameters
	if tokenID == "" || amount <= 0 || senderID == "" || recipientID == "" {
		return nil, errors.New("invalid parameters for issuing a token")
	}

	// Create the new SYN2200 token
	newToken := &common.SYN2200Token{
		TokenID:     tokenID,
		Currency:    currency,
		Amount:      amount,
		Sender:      senderID,
		Recipient:   recipientID,
		CreatedTime: time.Now(),
		Executed:    false,
	}

	// Add the token to the ledger
	err := ledger.MintToken(newToken)
	if err != nil {
		return nil, errors.New("error minting token: " + err.Error())
	}

	// Register the event in the Synnergy Consensus
	err = consensus.RegisterTokenCreation(newToken)
	if err != nil {
		return nil, errors.New("error registering token creation in consensus: " + err.Error())
	}

	// Store encrypted sensitive data related to the token
	encryptedData, err := encryption.EncryptData([]byte(newToken.TokenID + newToken.Currency))
	if err != nil {
		return nil, errors.New("error encrypting token data: " + err.Error())
	}
	newToken.EncryptedMetadata = encryptedData

	return newToken, nil
}

// ManageTokenLifecycle ensures compliance and lifecycle events are managed for SYN2200 tokens.
func ManageTokenLifecycle(tokenID string) (*common.SYN2200Token, error) {
	// Retrieve the token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("token not found: " + err.Error())
	}

	// Validate the token's compliance status
	compliant, err := compliance.VerifyCompliance(token)
	if err != nil || !compliant {
		return nil, errors.New("token is not compliant with regulatory requirements: " + err.Error())
	}

	// Proceed with token lifecycle management based on state
	if !token.Executed {
		// If the token is not yet executed (settled), proceed with settlement
		token.Executed = true
		token.ExecutedTime = time.Now()

		// Log the settlement in the ledger and consensus
		err = ledger.RecordTokenSettlement(token)
		if err != nil {
			return nil, errors.New("error recording token settlement in ledger: " + err.Error())
		}

		err = consensus.RecordSettlementEvent(token)
		if err != nil {
			return nil, errors.New("error recording settlement event in consensus: " + err.Error())
		}
	} else {
		// Token already settled, notify of final state
		return token, errors.New("token is already settled")
	}

	return token, nil
}

// CancelToken handles cancellation of an unexecuted SYN2200 token.
func CancelToken(tokenID string, cancelBy string) error {
	// Retrieve the token
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("token not found: " + err.Error())
	}

	// Ensure only the sender or the system can cancel the token
	if token.Sender != cancelBy && !compliance.IsSystemCancelAllowed() {
		return errors.New("only the sender or the system can cancel the token")
	}

	// Ensure token hasn't been executed (settled)
	if token.Executed {
		return errors.New("token is already settled, cannot be canceled")
	}

	// Cancel the token in the ledger
	err = ledger.CancelToken(tokenID)
	if err != nil {
		return errors.New("error canceling token: " + err.Error())
	}

	// Log the cancellation in the consensus
	err = consensus.RecordCancellationEvent(token)
	if err != nil {
		return errors.New("error recording cancellation in consensus: " + err.Error())
	}

	return nil
}

// GetTokenDetails retrieves full token details, including decrypted sensitive data.
func GetTokenDetails(tokenID string) (*common.SYN2200Token, error) {
	// Retrieve the token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("token not found: " + err.Error())
	}

	// Decrypt the sensitive data
	decryptedData, err := encryption.DecryptData(token.EncryptedMetadata)
	if err != nil {
		return nil, errors.New("error decrypting token data: " + err.Error())
	}
	token.DecryptedMetadata = string(decryptedData)

	return token, nil
}

// ManageCompliance ensures compliance checks are performed regularly for SYN2200 tokens.
func ManageCompliance(tokenID string) (bool, error) {
	// Retrieve the token
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return false, errors.New("token not found: " + err.Error())
	}

	// Perform compliance checks
	compliant, err := compliance.PerformComplianceCheck(token)
	if err != nil {
		return false, errors.New("error performing compliance check: " + err.Error())
	}

	// Log the compliance status
	err = ledger.RecordComplianceStatus(tokenID, compliant)
	if err != nil {
		return false, errors.New("error recording compliance status in ledger: " + err.Error())
	}

	return compliant, nil
}

// TransferOwnership handles the ownership transfer of SYN2200 tokens.
func TransferOwnership(tokenID string, newOwnerID string) error {
	// Retrieve the token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("token not found: " + err.Error())
	}

	// Ensure the token is not already settled
	if token.Executed {
		return errors.New("token has already been settled and cannot be transferred")
	}

	// Record the ownership transfer
	oldOwner := token.Recipient
	token.Recipient = newOwnerID

	err = ledger.RecordTokenTransfer(tokenID, oldOwner, newOwnerID)
	if err != nil {
		return errors.New("error recording token transfer: " + err.Error())
	}

	// Log the transfer in the consensus
	err = consensus.RecordTransferEvent(token)
	if err != nil {
		return errors.New("error recording transfer in consensus: " + err.Error())
	}

	return nil
}
