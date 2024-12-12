package interoperability

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/tokens/syn1200"
	"time"
)

// NewAtomicSwapManager initializes the atomic swap manager
func NewAtomicSwapManager(ledgerInstance *ledger.Ledger) *AtomicSwapManager {
	return &AtomicSwapManager{
		ActiveSwaps:   make(map[string]*AtomicSwap),
		LedgerInstance: ledgerInstance,
	}
}

// InitiateSwap initializes a new atomic swap between Chain A and Chain B using Syn500 tokens
func (asm *AtomicSwapManager) InitiateSwap(tokenA syn1200.SYN1200Token, amountA float64, chainAAddress string,
	tokenB syn1200.SYN1200Token, amountB float64, chainBAddress string, initiator string, secret string, expirationDuration time.Duration) (string, error) {

	asm.mutex.Lock()
	defer asm.mutex.Unlock()

	// Generate a unique swap ID
	swapID := asm.generateSwapID(initiator, tokenA.Symbol, tokenB.Symbol)

	// Hash the secret
	secretHash := asm.generateSecretHash(secret)

	// Create the atomic swap
	swap := &AtomicSwap{
		TokenA:         tokenA,
		AmountA:        amountA,
		ChainAAddress:  chainAAddress,
		TokenB:         tokenB,
		AmountB:        amountB,
		ChainBAddress:  chainBAddress,
		SecretHash:     secretHash,
		Secret:         "",
		ExpirationTime: time.Now().Add(expirationDuration),
		SwapInitiator:  initiator,
		Status:         "pending",
		LedgerInstance: asm.LedgerInstance,
	}

	// Record the swap in the ledger
	err := asm.recordSwapToLedger(swap)
	if err != nil {
		return "", fmt.Errorf("failed to record swap in the ledger: %v", err)
	}

	asm.ActiveSwaps[swapID] = swap

	fmt.Printf("Atomic swap initiated by %s. Swap ID: %s\n", initiator, swapID)
	return swapID, nil
}

// CompleteSwap completes the atomic swap by providing the correct secret for validation
func (asm *AtomicSwapManager) CompleteSwap(swapID string, secret string, responder string) error {
	asm.mutex.Lock()
	defer asm.mutex.Unlock()

	swap, exists := asm.ActiveSwaps[swapID]
	if !exists {
		return errors.New("swap not found")
	}

	if swap.Status != "pending" {
		return errors.New("swap is not in a pending state")
	}

	if time.Now().After(swap.ExpirationTime) {
		swap.Status = "expired"
		return errors.New("swap has expired")
	}

	// Validate the secret
	if asm.generateSecretHash(secret) != swap.SecretHash {
		return errors.New("invalid secret")
	}

	// Complete the swap
	swap.Secret = secret
	swap.SwapResponder = responder
	swap.Status = "completed"

	// Log the swap completion to the ledger
	err := asm.logSwapCompletionToLedger(swap)
	if err != nil {
		return fmt.Errorf("failed to log swap completion to ledger: %v", err)
	}

	fmt.Printf("Atomic swap completed. Swap ID: %s\n", swapID)
	return nil
}

// ExpireSwap marks a swap as expired if the expiration time is reached
func (asm *AtomicSwapManager) ExpireSwap(swapID string) error {
	asm.mutex.Lock()
	defer asm.mutex.Unlock()

	swap, exists := asm.ActiveSwaps[swapID]
	if !exists {
		return errors.New("swap not found")
	}

	if swap.Status != "pending" {
		return errors.New("swap is not in a pending state")
	}

	if time.Now().Before(swap.ExpirationTime) {
		return errors.New("swap has not yet expired")
	}

	swap.Status = "expired"
	fmt.Printf("Atomic swap expired. Swap ID: %s\n", swapID)

	// Log the swap expiration to the ledger
	return asm.logSwapExpirationToLedger(swap)
}

// generateSwapID generates a unique swap ID based on the initiator and tokens involved
func (asm *AtomicSwapManager) generateSwapID(initiator string, tokenASymbol string, tokenBSymbol string) string {
	hashInput := fmt.Sprintf("%s%s%s%d", initiator, tokenASymbol, tokenBSymbol, time.Now().UnixNano())
	hash := sha256.New()
	hash.Write([]byte(hashInput))
	return hex.EncodeToString(hash.Sum(nil))
}

// generateSecretHash generates a SHA-256 hash for the provided secret
func (asm *AtomicSwapManager) generateSecretHash(secret string) string {
	hash := sha256.New()
	hash.Write([]byte(secret))
	return hex.EncodeToString(hash.Sum(nil))
}

// recordSwapToLedger logs the initiation of an atomic swap to the ledger.
func (asm *AtomicSwapManager) recordSwapToLedger(swap *AtomicSwap) error {
    // Generate a unique ID for the swap if not already set
    if swap.SwapID == "" {
        swap.SwapID = generateSwapID()
    }

    // Serialize swap data for encryption (excluding sensitive data)
    swapData := fmt.Sprintf("Atomic swap initiated with ID: %s, ChainAAddress: %s, ChainBAddress: %s, AmountA: %.2f, AmountB: %.2f, ExpirationTime: %s",
        swap.SwapID, swap.ChainAAddress, swap.ChainBAddress, swap.AmountA, swap.AmountB, swap.ExpirationTime)

    // Create an encryption instance and handle any errors
    encryptInstance, err := common.NewEncryption(256) // Assuming 256 bits; adjust as required
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Ensure the encryption key length is valid (16, 24, or 32 bytes)
    if len(common.EncryptionKey) != 16 && len(common.EncryptionKey) != 24 && len(common.EncryptionKey) != 32 {
        return fmt.Errorf("encryption key length must be 16, 24, or 32 bytes")
    }

    // Encrypt the swap data (assuming EncryptData now requires IV or additional parameters)
    iv := make([]byte, 16) // Example initialization vector; this should be securely generated if required
    rand.Read(iv)
    _, err = encryptInstance.EncryptData(swapData, common.EncryptionKey, iv)
    if err != nil {
        return fmt.Errorf("failed to encrypt swap data: %v", err)
    }

    // Record the atomic swap initiation in the ledger with swap details
    asm.LedgerInstance.RecordAtomicSwapInitiation(
        swap.SwapID,
        swap.SwapInitiator,
        swap.SwapResponder,
        "ChainA",      // Assuming "ChainA" represents the chain ID for Chain A
        swap.AmountA,  // Amount of Token A
        swap.ExpirationTime,
    )

    return nil
}


// generateSwapID generates a unique ID for atomic swap transactions.
func generateSwapID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}

// logSwapCompletionToLedger logs the completion of an atomic swap to the ledger.
func (asm *AtomicSwapManager) logSwapCompletionToLedger(swap *AtomicSwap) error {
    // Serialize swap data for logging/audit purposes
    swapData := fmt.Sprintf("Completed atomic swap: %+v", swap)

    // Create an encryption instance
    encryptInstance, err := common.NewEncryption(256)
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt the swap data (optional - only for secure audit)
    _, err = encryptInstance.EncryptData(swapData, common.EncryptionKey, nil)
    if err != nil {
        return fmt.Errorf("failed to encrypt swap completion data: %v", err)
    }

    // Record the atomic swap completion in the ledger using only the swapID
    asm.LedgerInstance.RecordAtomicSwapCompletion(swap.SwapID)

    return nil
}

// logSwapExpirationToLedger logs the expiration of an atomic swap to the ledger.
func (asm *AtomicSwapManager) logSwapExpirationToLedger(swap *AtomicSwap) error {
    // Serialize swap data for logging/audit purposes
    swapData := fmt.Sprintf("Expired atomic swap: %+v", swap)

    // Create an encryption instance
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt the swap data (only for secure logging/audit purposes if needed)
    _, err = encryptInstance.EncryptData(swapData, common.EncryptionKey, nil)
    if err != nil {
        return fmt.Errorf("failed to encrypt swap expiration data: %v", err)
    }

    // Record the atomic swap expiration in the ledger using only the swapID
    asm.LedgerInstance.RecordAtomicSwapExpiration(swap.SwapID)

    return nil
}

