package syn3000

import (
    "errors"
    "time"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/hex"
    "io"
)

// SYN3000Compliance handles compliance checks, validations, and auditing for SYN3000 tokens.
type SYN3000Compliance struct{}

// NewSYN3000Compliance initializes the compliance system for SYN3000 tokens.
func NewSYN3000Compliance() *SYN3000Compliance {
    return &SYN3000Compliance{}
}

// ValidateRentalToken ensures the rental token complies with legal and regulatory requirements.
func (c *SYN3000Compliance) ValidateRentalToken(token *SYN3000Token) error {
    // Check if token is active
    if !token.ActiveStatus {
        return errors.New("rental token is not active")
    }

    // Validate lease dates
    if time.Now().After(token.LeaseEndDate) {
        return errors.New("lease end date has passed")
    }

    if time.Now().Before(token.LeaseStartDate) {
        return errors.New("lease start date is in the future")
    }

    // Ensure that rental terms are within reasonable legal limits (dummy check here)
    if token.MonthlyRent <= 0 {
        return errors.New("monthly rent must be greater than zero")
    }

    if token.DepositAmount < 0 {
        return errors.New("deposit amount cannot be negative")
    }

    return nil
}

// CheckCompliance performs a full compliance audit for a given token, ensuring regulatory alignment.
func (c *SYN3000Compliance) CheckCompliance(tokenID string) (bool, error) {
    // Retrieve the encrypted token from the ledger
    encryptedToken, err := ledger.GetToken(tokenID)
    if err != nil {
        return false, errors.New("token not found")
    }

    // Decrypt token
    decryptedToken, err := c.decryptTokenData(encryptedToken)
    if err != nil {
        return false, err
    }

    // Validate the token
    if err := c.ValidateRentalToken(decryptedToken); err != nil {
        return false, err
    }

    // Additional compliance checks can be added here

    return true, nil
}

// Encrypt sensitive token data before storing in the ledger
func (c *SYN3000Compliance) encryptTokenData(token *SYN3000Token) (string, error) {
    key := []byte(common.GetEncryptionKey()) // Fetch encryption key from the common package
    plaintext, err := common.Serialize(token)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return hex.EncodeToString(ciphertext), nil
}

// Decrypt token data from the ledger
func (c *SYN3000Compliance) decryptTokenData(encryptedData string) (*SYN3000Token, error) {
    key := []byte(common.GetEncryptionKey()) // Fetch encryption key from the common package
    ciphertext, err := hex.DecodeString(encryptedData)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    if len(ciphertext) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    var token SYN3000Token
    if err := common.Deserialize(ciphertext, &token); err != nil {
        return nil, err
    }

    return &token, nil
}

// GenerateComplianceReport generates a detailed report of all tokens' compliance status
func (c *SYN3000Compliance) GenerateComplianceReport() ([]common.ComplianceReport, error) {
    var reports []common.ComplianceReport

    // Retrieve all tokens from the ledger
    allTokens, err := ledger.GetAllTokens("SYN3000")
    if err != nil {
        return nil, errors.New("unable to retrieve tokens from the ledger")
    }

    // Perform compliance checks on each token
    for _, encryptedToken := range allTokens {
        decryptedToken, err := c.decryptTokenData(encryptedToken)
        if err != nil {
            continue // Skip token if there's a decryption error
        }

        complianceStatus, err := c.CheckCompliance(decryptedToken.TokenID)
        if err != nil {
            complianceStatus = false
        }

        // Create a compliance report entry
        report := common.ComplianceReport{
            TokenID:    decryptedToken.TokenID,
            PropertyID: decryptedToken.PropertyID,
            TenantID:   decryptedToken.TenantID,
            Status:     complianceStatus,
            Date:       time.Now(),
        }

        reports = append(reports, report)
    }

    return reports, nil
}
