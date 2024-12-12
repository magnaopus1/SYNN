package common

import "crypto/ecdsa"

// Validator represents a participant in the Proof of Stake system
type Validator struct {
	Address string  // Validator's wallet address
	Stake   float64 // Amount of SYNN staked by the validator
	ReputationScore float64 // Validator's reputation score
    PublicKey     *ecdsa.PublicKey // Validator's public key for signature verification

}

// ValidatorSelectionRecord holds information about a validator's selection in a specific epoch
type ValidatorSelectionRecord struct {
	ValidatorAddress string // The address of the selected validator
	Epoch            int    // The epoch in which the validator was selected
}