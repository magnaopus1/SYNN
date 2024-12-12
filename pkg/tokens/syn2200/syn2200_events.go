package syn2200

import (
	"errors"
	"time"

)

// MintTokenEvent triggers the minting of a new SYN2200 token for real-time payments.
func MintTokenEvent(tokenID string, amount float64, currency string, senderID string, recipientID string) (*common.SYN2200Token, error) {
	// Check if all necessary details are present
	if tokenID == "" || senderID == "" || recipientID == "" || amount <= 0 {
		return nil, errors.New("invalid token creation parameters")
	}

	// Create new token
	token := &common.SYN2200Token{
		TokenID:     tokenID,
		Currency:    currency,
		Amount:      amount,
		Sender:      senderID,
		Recipient:   recipientID,
		CreatedTime: time.Now(),
		Executed:    false,
	}

	// Add the token to the ledger
	err := ledger.MintToken(token)
	if err != nil {
		return nil, errors.New("error minting token: " + err.Error())
	}

	// Log the event in the consensus mechanism
	err = consensus.RecordMintingEvent(token)
	if err != nil {
		return nil, errors.New("error recording minting event: " + err.Error())
	}

	// Notify users of the minted token
	notifications.SendRealTimeNotification(senderID, recipientID, "SYN2200 token minted for real-time payment")

	return token, nil
}

// TransferTokenEvent manages the transfer of a SYN2200 token between users.
func TransferTokenEvent(tokenID string, newOwner string) (*common.SYN2200Token, error) {
	// Retrieve token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("token not found: " + err.Error())
	}

	// Ensure the token is not already settled
	if token.Executed {
		return nil, errors.New("token has already been settled and cannot be transferred")
	}

	// Update the token's owner
	oldOwner := token.Recipient
	token.Recipient = newOwner

	// Record the transfer event in the ledger
	err = ledger.RecordTokenTransfer(tokenID, oldOwner, newOwner)
	if err != nil {
		return nil, errors.New("error recording token transfer: " + err.Error())
	}

	// Notify users of the transfer
	notifications.SendRealTimeNotification(oldOwner, newOwner, "SYN2200 token transferred to new recipient")

	// Log the transfer event in consensus
	err = consensus.RecordTransferEvent(token)
	if err != nil {
		return nil, errors.New("error recording transfer event in consensus: " + err.Error())
	}

	return token, nil
}

// SettleTokenEvent finalizes the transaction associated with a SYN2200 token, ensuring real-time payment settlement.
func SettleTokenEvent(tokenID string) (*common.SYN2200Token, error) {
	// Retrieve token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("token not found: " + err.Error())
	}

	// Ensure the token has not been settled yet
	if token.Executed {
		return nil, errors.New("token has already been settled")
	}

	// Mark the token as executed (settled)
	token.Executed = true
	token.ExecutedTime = time.Now()

	// Record the settlement in the ledger
	err = ledger.RecordTokenSettlement(token)
	if err != nil {
		return nil, errors.New("error recording token settlement: " + err.Error())
	}

	// Notify both sender and recipient of the settlement
	notifications.SendRealTimeNotification(token.Sender, token.Recipient, "SYN2200 token payment settled")

	// Record the settlement event in consensus
	err = consensus.RecordSettlementEvent(token)
	if err != nil {
		return nil, errors.New("error recording settlement event in consensus: " + err.Error())
	}

	return token, nil
}

// CancelTokenEvent allows the sender to cancel a pending SYN2200 token before it is settled.
func CancelTokenEvent(tokenID string, senderID string) (*common.SYN2200Token, error) {
	// Retrieve the token
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("token not found: " + err.Error())
	}

	// Ensure the token has not been settled yet
	if token.Executed {
		return nil, errors.New("token has already been settled and cannot be canceled")
	}

	// Only the sender can cancel the transaction
	if token.Sender != senderID {
		return nil, errors.New("only the sender can cancel the token")
	}

	// Mark the token as canceled in the ledger
	err = ledger.CancelToken(tokenID)
	if err != nil {
		return nil, errors.New("error canceling token: " + err.Error())
	}

	// Log the cancellation event in consensus
	err = consensus.RecordCancellationEvent(token)
	if err != nil {
		return nil, errors.New("error recording cancellation event in consensus: " + err.Error())
	}

	// Notify users of the cancellation
	notifications.SendRealTimeNotification(token.Sender, token.Recipient, "SYN2200 token payment has been canceled")

	return token, nil
}

// AddEventAuditLog adds a new audit log entry for lifecycle events such as minting, transfer, or settlement.
func AddEventAuditLog(token *common.SYN2200Token, event string, performedBy string, description string) error {
	// Add a new audit entry
	token.AddAuditLog(event, performedBy, description)

	// Record the audit log in the ledger
	err := ledger.RecordEventAuditLog(token, event, performedBy, description)
	if err != nil {
		return errors.New("failed to record event audit log in ledger: " + err.Error())
	}
	return nil
}

// EncryptEventData encrypts sensitive event data for SYN2200 tokens.
func EncryptEventData(token *common.SYN2200Token, data []byte) error {
	// Encrypt event data
	encryptedData, err := encryption.EncryptData(data)
	if err != nil {
		return err
	}

	// Store encrypted data
	token.EncryptedMetadata = encryptedData
	return nil
}

// DecryptEventData decrypts sensitive event data for SYN2200 tokens.
func DecryptEventData(token *common.SYN2200Token, encryptedData []byte) ([]byte, error) {
	// Decrypt the event data
	return encryption.DecryptData(encryptedData)
}
