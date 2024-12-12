package sustainability

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewEcoFriendlyNodeCertificationSystem initializes a new certification system
func NewEcoFriendlyNodeCertificationSystem(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *EcoFriendlyNodeCertificationSystem {
	return &EcoFriendlyNodeCertificationSystem{
		Certificates:      make(map[string]*EcoFriendlyNodeCertificate),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// IssueCertificate issues an eco-friendly certificate to a node
func (ecs *EcoFriendlyNodeCertificationSystem) IssueCertificate(certificateID, nodeID, owner, issuer string, validityPeriod time.Duration) (*EcoFriendlyNodeCertificate, error) {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	// Encrypt certificate data
	certificateData := fmt.Sprintf("CertificateID: %s, NodeID: %s, Owner: %s, Issuer: %s", certificateID, nodeID, owner, issuer)
	encryptedData, err := ecs.EncryptionService.EncryptData([]byte(certificateData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt certificate data: %v", err)
	}

	// Create the certificate
	certificate := &EcoFriendlyNodeCertificate{
		CertificateID: certificateID,
		NodeID:        nodeID,
		Owner:         owner,
		Issuer:        issuer,
		IssueDate:     time.Now(),
		ExpiryDate:    time.Now().Add(validityPeriod),
		IsRevoked:     false,
	}

	// Add the certificate to the system
	ecs.Certificates[certificateID] = certificate

	// Log the certificate issuance in the ledger
	err = ecs.Ledger.RecordEcoFriendlyCertificateIssuance(certificateID, nodeID, owner, issuer, time.Now(), certificate.ExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("failed to log certificate issuance: %v", err)
	}

	fmt.Printf("Eco-friendly certificate %s issued to node %s (owner: %s)\n", certificateID, nodeID, owner)
	return certificate, nil
}

// RevokeCertificate revokes an eco-friendly certificate before its expiry
func (ecs *EcoFriendlyNodeCertificationSystem) RevokeCertificate(certificateID string, revocationReason string) error {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	// Retrieve the certificate
	certificate, exists := ecs.Certificates[certificateID]
	if !exists {
		return fmt.Errorf("certificate %s not found", certificateID)
	}

	// Ensure the certificate has not already been revoked
	if certificate.IsRevoked {
		return fmt.Errorf("certificate %s has already been revoked", certificateID)
	}

	// Revoke the certificate
	certificate.IsRevoked = true
	certificate.RevokedDate = time.Now()

	// Log the revocation in the ledger
	err := ecs.Ledger.RecordCertificateRevocation(certificateID, certificate.NodeID, revocationReason, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log certificate revocation: %v", err)
	}

	fmt.Printf("Eco-friendly certificate %s revoked for node %s\n", certificateID, certificate.NodeID)
	return nil
}

// RenewCertificate renews an eco-friendly certificate, extending its validity period
func (ecs *EcoFriendlyNodeCertificationSystem) RenewCertificate(certificateID string, additionalPeriod time.Duration) (*EcoFriendlyNodeCertificate, error) {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	// Retrieve the certificate
	certificate, exists := ecs.Certificates[certificateID]
	if !exists {
		return nil, fmt.Errorf("certificate %s not found", certificateID)
	}

	// Ensure the certificate has not been revoked
	if certificate.IsRevoked {
		return nil, fmt.Errorf("certificate %s has been revoked and cannot be renewed", certificateID)
	}

	// Extend the expiry date
	oldExpiry := certificate.ExpiryDate
	certificate.ExpiryDate = certificate.ExpiryDate.Add(additionalPeriod)

	// Log the certificate renewal in the ledger
	err := ecs.Ledger.RecordCertificateRenewal(certificateID, certificate.NodeID, oldExpiry, certificate.ExpiryDate, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log certificate renewal: %v", err)
	}

	fmt.Printf("Eco-friendly certificate %s renewed for node %s (new expiry: %s)\n", certificateID, certificate.NodeID, certificate.ExpiryDate)
	return certificate, nil
}

// ViewCertificate allows viewing of an eco-friendly certificate's details
func (ecs *EcoFriendlyNodeCertificationSystem) ViewCertificate(certificateID string) (*EcoFriendlyNodeCertificate, error) {
	ecs.mu.Lock()
	defer ecs.mu.Unlock()

	// Retrieve the certificate
	certificate, exists := ecs.Certificates[certificateID]
	if !exists {
		return nil, fmt.Errorf("certificate %s not found", certificateID)
	}

	return certificate, nil
}

// generateUniqueID creates a cryptographically secure unique ID
func generateUniqueID() string {
	id := make([]byte, 16)
	rand.Read(id)
	return hex.EncodeToString(id)
}
