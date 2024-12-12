package ledger

import (
	"errors"
	"fmt"
	"time"
)

// RecordCarbonCreditIssuance logs the issuance of a carbon credit.
func (l *SustainabilityLedger) RecordCarbonCreditIssuance(creditID, issuer string, amount float64) error {
	l.Lock()
	defer l.Unlock()

	// Create and store the carbon credit
	credit := CarbonCredit{
		CreditID: creditID,
		Issuer:   issuer,
		Amount:   amount,
		Status:   "active",
		IssuedAt: time.Now(),
	}

	l.CarbonCredits[creditID] = credit
	return nil
}

// RecordCarbonCreditTransfer logs the transfer of a carbon credit between accounts.
func (l *SustainabilityLedger) RecordCarbonCreditTransfer(creditID, newOwner string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the credit exists
	credit, exists := l.CarbonCredits[creditID]
	if !exists {
		return errors.New("carbon credit not found")
	}

	// Transfer ownership
	credit.Owner = newOwner
	l.CarbonCredits[creditID] = credit
	return nil
}

// RecordCarbonCreditRetirement logs the retirement of a carbon credit.
func (l *SustainabilityLedger) RecordCarbonCreditRetirement(creditID string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the credit exists and retire it
	credit, exists := l.CarbonCredits[creditID]
	if !exists {
		return errors.New("carbon credit not found")
	}

	credit.Status = "retired"
	credit.RetiredAt = time.Now()
	l.CarbonCredits[creditID] = credit
	return nil
}

// RecordOffsetRequest logs a request for carbon offset.
func (l *SustainabilityLedger) RecordOffsetRequest(requestID, requester string, amount float64) error {
	if requestID == "" || requester == "" || amount <= 0 {
		return fmt.Errorf("invalid offset request data")
	}

	l.Lock()
	defer l.Unlock()

	l.OffsetRequests[requestID] = OffsetRequest{
		RequestID:  requestID,
		Requester:  requester,
		Amount:     amount,
		Timestamp:  time.Now(),
		Status:     "Pending",
	}
	return nil
}

// RecordAvailableCredit logs available carbon credits for offset matching.
func (l *SustainabilityLedger) RecordAvailableCredit(creditID string, availableAmount float64) error {
	if creditID == "" || availableAmount <= 0 {
		return fmt.Errorf("invalid available credit data")
	}

	l.Lock()
	defer l.Unlock()

	l.CarbonCredits[creditID] = CarbonCredit{
		CreditID:    creditID,
		Amount:      availableAmount,
		Allocated:   false,
		AllocatedTo: "",
		Timestamp:   time.Now(),
	}
	return nil
}

// RecordOffsetMatch logs the matching of an offset request with available credits.
func (l *SustainabilityLedger) RecordOffsetMatch(requestID, creditID string, matchedAmount float64) error {
	if requestID == "" || creditID == "" || matchedAmount <= 0 {
		return fmt.Errorf("invalid match data")
	}

	l.Lock()
	defer l.Unlock()

	request, exists := l.OffsetRequests[requestID]
	if !exists {
		return fmt.Errorf("offset request not found")
	}

	credit, exists := l.CarbonCredits[creditID]
	if !exists {
		return fmt.Errorf("carbon credit not found")
	}

	if credit.Amount < matchedAmount || credit.Allocated {
		return fmt.Errorf("insufficient or already allocated credits")
	}

	credit.Amount -= matchedAmount
	credit.Allocated = true
	credit.AllocatedTo = requestID
	request.Status = "Matched"

	l.OffsetRequests[requestID] = request
	l.CarbonCredits[creditID] = credit

	return nil
}


// RecordEcoFriendlyCertificateIssuance logs the issuance of an eco-friendly certificate.
func (l *SustainabilityLedger) RecordEcoFriendlyCertificateIssuance(certID, recipient, certType string) error {
	l.Lock()
	defer l.Unlock()

	// Create and store the certificate
	certificate := EcoCertificate{
		CertificateID: certID,
		Recipient:     recipient,
		Type:          certType,
		IssuedAt:      time.Now(),
	}

	l.EcoCertificates[certID] = certificate
	return nil
}

// RecordCertificateRevocation logs the revocation of an eco-friendly certificate.
func (l *SustainabilityLedger) RecordCertificateRevocation(certID string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the certificate exists and revoke it
	certificate, exists := l.EcoCertificates[certID]
	if !exists {
		return errors.New("certificate not found")
	}

	certificate.RevokedAt = time.Now()
	l.EcoCertificates[certID] = certificate
	return nil
}

// RecordCertificateRenewal logs the renewal of an eco-friendly certificate.
func (l *SustainabilityLedger) RecordCertificateRenewal(certID string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the certificate exists and renew it
	certificate, exists := l.EcoCertificates[certID]
	if !exists {
		return errors.New("certificate not found")
	}

	certificate.RenewedAt = time.Now()
	l.EcoCertificates[certID] = certificate
	return nil
}

// RecordEnergyEfficiencyRating logs an energy efficiency rating for hardware or software.
func (l *SustainabilityLedger) RecordEnergyEfficiencyRating(entityID string, rating float64) error {
	if entityID == "" || rating <= 0 {
		return fmt.Errorf("invalid energy efficiency rating data")
	}

	l.Lock()
	defer l.Unlock()

	l.EnergyEfficiencyRatings[entityID] = EnergyEfficiencyRating{
		EntityID:   entityID,
		Rating:     rating,
		Timestamp:  time.Now(),
		IsRevoked:  false,
	}
	return nil
}

// RecordRatingRevocation logs the revocation of an energy efficiency rating.
func (l *SustainabilityLedger) RecordRatingRevocation(entityID string) error {
	if entityID == "" {
		return fmt.Errorf("invalid entity ID")
	}

	l.Lock()
	defer l.Unlock()

	rating, exists := l.EnergyEfficiencyRatings[entityID]
	if !exists {
		return fmt.Errorf("rating not found")
	}

	rating.IsRevoked = true
	rating.Timestamp = time.Now()

	l.EnergyEfficiencyRatings[entityID] = rating
	return nil
}

// RecordRatingRenewal logs the renewal of an energy efficiency rating.
func (l *SustainabilityLedger) RecordRatingRenewal(entityID string) error {
	if entityID == "" {
		return fmt.Errorf("invalid entity ID")
	}

	l.Lock()
	defer l.Unlock()

	rating, exists := l.EnergyEfficiencyRatings[entityID]
	if !exists || rating.IsRevoked {
		return fmt.Errorf("rating not found or revoked")
	}

	rating.Timestamp = time.Now()
	rating.IsRevoked = false

	l.EnergyEfficiencyRatings[entityID] = rating
	return nil
}

// RecordEnergyUsage logs the energy usage for an entity.
func (l *SustainabilityLedger) RecordEnergyUsage(usageID, entityID string, energyUsed float64) error {
	l.Lock()
	defer l.Unlock()

	// Log the energy usage
	usage := EnergyUsage{
		UsageID:    usageID,
		EntityID:   entityID,
		EnergyUsed: energyUsed,
		RecordedAt: time.Now(),
	}

	l.EnergyUsageRecords[usageID] = usage
	return nil
}

// RecordEnergyUsageUpdate logs updates to existing energy usage records.
func (l *SustainabilityLedger) RecordEnergyUsageUpdate(usageID string, updatedUsage float64) error {
	l.Lock()
	defer l.Unlock()

	// Check if the energy usage record exists
	usage, exists := l.EnergyUsageRecords[usageID]
	if !exists {
		return errors.New("energy usage record not found")
	}

	// Update the usage
	usage.EnergyUsed = updatedUsage
	l.EnergyUsageRecords[usageID] = usage
	return nil
}

// RecordGreenHardwareRegistration logs the registration of green hardware.
func (l *SustainabilityLedger) RecordGreenHardwareRegistration(hardwareID, entityID string) error {
	if hardwareID == "" || entityID == "" {
		return fmt.Errorf("invalid hardware registration data")
	}

	l.Lock()
	defer l.Unlock()

	l.GreenHardwareRegistry[hardwareID] = GreenHardware{
		HardwareID: hardwareID,
		EntityID:   entityID,
		Registered: true,
		Timestamp:  time.Now(),
	}
	return nil
}

// RecordSoftwareRegistration logs the registration of eco-friendly software.
func (l *SustainabilityLedger) RecordSoftwareRegistration(softwareID, entityID string) error {
	if softwareID == "" || entityID == "" {
		return fmt.Errorf("invalid software registration data")
	}

	l.Lock()
	defer l.Unlock()

	l.EcoFriendlySoftwareRegistry[softwareID] = EcoFriendlySoftware{
		SoftwareID: softwareID,
		EntityID:   entityID,
		Registered: true,
		Timestamp:  time.Now(),
	}
	return nil
}

// RecordCircularEconomyProgram logs the registration of a circular economy program.
func (l *SustainabilityLedger) RecordCircularEconomyProgram(programID, description string) error {
	if programID == "" || description == "" {
		return fmt.Errorf("invalid program data")
	}

	l.Lock()
	defer l.Unlock()

	l.CircularEconomyPrograms[programID] = CircularEconomyProgram{
		ProgramID:   programID,
		Description: description,
		Timestamp:   time.Now(),
	}
	return nil
}

// RecordEcoFriendlyCertificateAward logs the awarding of eco-friendly certificates.
func (l *SustainabilityLedger) RecordEcoFriendlyCertificateAward(certID, recipient string) error {
	if certID == "" || recipient == "" {
		return fmt.Errorf("invalid certificate award data")
	}

	l.Lock()
	defer l.Unlock()

	l.EcoFriendlyCertificates[certID] = EcoFriendlyCertificate{
		CertificateID: certID,
		Recipient:     recipient,
		AwardedAt:     time.Now(),
	}
	return nil
}

// RecordConservationInitiative logs a conservation initiative.
func (l *SustainabilityLedger) RecordConservationInitiative(initID, name, description string) error {
	l.Lock()
	defer l.Unlock()

	// Log the conservation initiative
	initiative := ConservationInitiative{
		InitiativeID: initID,
		Name:         name,
		Description:  description,
		LaunchedAt:   time.Now(),
	}

	l.ConservationPrograms = append(l.ConservationPrograms, initiative)
	return nil
}

// RecordCoolingSolutionOptimization logs optimization of cooling solutions.
func (l *SustainabilityLedger) RecordCoolingSolutionOptimization(optID, details string) error {
	l.Lock()
	defer l.Unlock()

	// Log the cooling solution optimization
	optimization := OptimizationRecord{
		OptimizationID: optID,
		Details:        details,
		OptimizedAt:    time.Now(),
	}

	l.OptimizationRecords = append(l.OptimizationRecords, optimization)
	return nil
}

// RecordEnergyConsumptionOptimization logs the optimization of energy consumption.
func (l *SustainabilityLedger) RecordEnergyConsumptionOptimization(optID, details string) error {
	l.Lock()
	defer l.Unlock()

	// Log the energy consumption optimization
	optimization := OptimizationRecord{
		OptimizationID: optID,
		Details:        details,
		OptimizedAt:    time.Now(),
	}

	l.OptimizationRecords = append(l.OptimizationRecords, optimization)
	return nil
}

// RecordRenewableEnergySourceRegistration logs the registration of a renewable energy source.
func (l *SustainabilityLedger) RecordRenewableEnergySourceRegistration(sourceID, sourceType string, energyAmount float64) error {
	l.Lock()
	defer l.Unlock()

	// Create and store the renewable energy source
	source := RenewableEnergySource{
		SourceID:       sourceID,
		SourceType:     sourceType,     // Use SourceType instead of Contributor
		EnergyProduced: energyAmount,   // Use EnergyProduced instead of EnergyAmount
		IntegrationDate: time.Now(),
	}

	l.RenewableEnergy[sourceID] = source
	return nil
}


// RecordRenewableEnergyContribution logs contributions from a renewable energy source.
func (l *SustainabilityLedger) RecordRenewableEnergyContribution(sourceID string, contributionAmount float64) error {
	l.Lock()
	defer l.Unlock()

	// Check if the source exists
	source, exists := l.RenewableEnergy[sourceID]
	if !exists {
		return errors.New("renewable energy source not found")
	}

	// Update the energy amount contributed
	source.EnergyProduced += contributionAmount  // Use EnergyProduced instead of EnergyAmount
	l.RenewableEnergy[sourceID] = source
	return nil
}
