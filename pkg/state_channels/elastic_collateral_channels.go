package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewElasticCollateralChannel initializes a new Elastic Collateral Channel (ECC)
func NewElasticCollateralChannel(channelID string, participants []string, initialCollateral map[string]float64, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, collateralCap float64) *common.ElasticCollateralChannel {
	participantMap := make(map[string]*common.CollateralRecord)
	for _, participant := range participants {
		participantMap[participant] = &common.CollateralRecord{
			Participant: participant,
			Collateral:  initialCollateral[participant],
			LastUpdated: time.Now(),
		}
	}

	return &common.ElasticCollateralChannel{
		ChannelID:      channelID,
		Participants:   participantMap,
		State:          make(map[string]interface{}),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		CollateralCap:  collateralCap,
	}
}

// OpenChannel opens the elastic collateral state channel
func (ecc *common.ElasticCollateralChannel) OpenChannel() error {
	ecc.mu.Lock()
	defer ecc.mu.Unlock()

	if ecc.IsOpen {
		return errors.New("state channel is already open")
	}

	ecc.IsOpen = true

	// Log channel opening in the ledger
	err := ecc.Ledger.RecordChannelOpening(ecc.ChannelID, ecc.getParticipantsList(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Elastic Collateral Channel %s opened with participants: %v\n", ecc.ChannelID, ecc.getParticipantsList())
	return nil
}

// CloseChannel closes the elastic collateral state channel and settles collateral
func (ecc *common.ElasticCollateralChannel) CloseChannel() error {
	ecc.mu.Lock()
	defer ecc.mu.Unlock()

	if !ecc.IsOpen {
		return errors.New("state channel is already closed")
	}

	// Final validation of collateral before closure
	err := ecc.ValidateCollateral()
	if err != nil {
		return fmt.Errorf("failed to validate collateral before closing: %v", err)
	}

	ecc.IsOpen = false

	// Log channel closure in the ledger
	err = ecc.Ledger.RecordChannelClosure(ecc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Elastic Collateral Channel %s closed\n", ecc.ChannelID)
	return nil
}

// AdjustCollateral allows a participant to adjust their collateral in the channel
func (ecc *common.ElasticCollateralChannel) AdjustCollateral(participant string, newCollateral float64) error {
	ecc.mu.Lock()
	defer ecc.mu.Unlock()

	record, exists := ecc.Participants[participant]
	if !exists {
		return fmt.Errorf("participant %s not found in channel", participant)
	}

	if newCollateral > ecc.CollateralCap {
		return fmt.Errorf("new collateral exceeds channel's collateral cap")
	}

	// Adjust the participant's collateral
	record.Collateral = newCollateral
	record.LastUpdated = time.Now()

	// Encrypt the collateral data before storing it
	encryptedCollateral, err := ecc.Encryption.EncryptData([]byte(fmt.Sprintf("%f", newCollateral)), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt collateral: %v", err)
	}

	// Log the collateral adjustment in the ledger
	err = ecc.Ledger.RecordCollateralAdjustment(ecc.ChannelID, participant, encryptedCollateral, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log collateral adjustment: %v", err)
	}

	fmt.Printf("Participant %s adjusted collateral to %f in channel %s\n", participant, newCollateral, ecc.ChannelID)
	return nil
}

// ValidateCollateral ensures all participants have valid collateral amounts
func (ecc *common.ElasticCollateralChannel) ValidateCollateral() error {
	ecc.mu.Lock()
	defer ecc.mu.Unlock()

	for participant, record := range ecc.Participants {
		if record.Collateral > ecc.CollateralCap {
			return fmt.Errorf("participant %s has exceeded the collateral cap", participant)
		}

		fmt.Printf("Participant %s has valid collateral: %f\n", participant, record.Collateral)
	}

	// Log the collateral validation in the ledger
	err := ecc.Ledger.RecordCollateralValidation(ecc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log collateral validation: %v", err)
	}

	fmt.Printf("All participants in channel %s have valid collateral\n", ecc.ChannelID)
	return nil
}

// getParticipantsList returns a list of all participants in the channel
func (ecc *common.ElasticCollateralChannel) getParticipantsList() []string {
	var participants []string
	for participant := range ecc.Participants {
		participants = append(participants, participant)
	}
	return participants
}

// RetrieveCollateral retrieves the collateral of a specific participant
func (ecc *common.ElasticCollateralChannel) RetrieveCollateral(participant string) (float64, error) {
	ecc.mu.Lock()
	defer ecc.mu.Unlock()

	record, exists := ecc.Participants[participant]
	if !exists {
		return 0, fmt.Errorf("participant %s not found in channel", participant)
	}

	fmt.Printf("Participant %s has %f collateral in channel %s\n", participant, record.Collateral, ecc.ChannelID)
	return record.Collateral, nil
}

// MonitorCollateral continuously monitors the collateral status of participants
func (ecc *common.ElasticCollateralChannel) MonitorCollateral(interval time.Duration) {
	for {
		time.Sleep(interval)

		// Check the collateral status for all participants
		for participant, record := range ecc.Participants {
			if record.Collateral > ecc.CollateralCap {
				fmt.Printf("Warning: Participant %s has exceeded the collateral cap\n", participant)
			}
		}
	}
}
