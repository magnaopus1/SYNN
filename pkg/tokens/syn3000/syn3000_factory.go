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

// SYN3000Token struct for House Rental Token Standard
type SYN3000Token struct {
    TokenID            string    `json:"token_id"`
    PropertyID         string    `json:"property_id"`
    OwnerID            string    `json:"owner_id"`
    TenantID           string    `json:"tenant_id"`
    LeaseStartDate     time.Time `json:"lease_start_date"`
    LeaseEndDate       time.Time `json:"lease_end_date"`
    MonthlyRent        float64   `json:"monthly_rent"`
    DepositAmount      float64   `json:"deposit_amount"`
    IssuedDate         time.Time `json:"issued_date"`
    ActiveStatus       bool      `json:"active_status"`
    LastUpdateDate     time.Time `json:"last_update_date"`
    PropertyDetails    Property  `json:"property_details"`
    RentPaymentHistory []Payment `json:"rent_payment_history"`
}

// Property struct to store detailed property information.
type Property struct {
    PropertyID      string `json:"property_id"`
    Address         string `json:"address"`
    Owner           string `json:"owner"`
    Description     string `json:"description"`
    Bedrooms        int    `json:"bedrooms"`
    Bathrooms       int    `json:"bathrooms"`
    SquareFootage   int    `json:"square_footage"`
    Availability    bool   `json:"availability"`
}

// Payment struct to manage rent payments.
type Payment struct {
    PaymentID    string    `json:"payment_id"`
    Amount       float64   `json:"amount"`
    PaymentDate  time.Time `json:"payment_date"`
    Status       string    `json:"status"`
}

// SYN3000Manager handles operations related to SYN3000 tokens.
type SYN3000Manager struct {
    mu sync.Mutex
}

// NewSYN3000Manager initializes the token manager.
func NewSYN3000Manager() *SYN3000Manager {
    return &SYN3000Manager{}
}

// CreateToken creates a new SYN3000 token for a rental agreement.
func (tm *SYN3000Manager) CreateToken(ownerID string, tenantID string, property Property, leaseStartDate, leaseEndDate time.Time, monthlyRent, depositAmount float64) (*SYN3000Token, error) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    token := &SYN3000Token{
        TokenID:        common.GenerateTokenID(),
        PropertyID:     property.PropertyID,
        OwnerID:        ownerID,
        TenantID:       tenantID,
        LeaseStartDate: leaseStartDate,
        LeaseEndDate:   leaseEndDate,
        MonthlyRent:    monthlyRent,
        DepositAmount:  depositAmount,
        IssuedDate:     time.Now(),
        ActiveStatus:   true,
        LastUpdateDate: time.Now(),
        PropertyDetails: property,
    }

    // Store token in the ledger
    if err := ledger.StoreToken(token); err != nil {
        return nil, err
    }

    return token, nil
}

// TransferLease transfers the lease rights from one tenant to another.
func (tm *SYN3000Manager) TransferLease(tokenID, newTenantID string) error {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    token, err := ledger.GetToken(tokenID)
    if err != nil {
        return errors.New("token not found")
    }

    if !token.ActiveStatus {
        return errors.New("lease is inactive")
    }

    // Update the tenant information
    token.TenantID = newTenantID
    token.LastUpdateDate = time.Now()

    // Store updated token in ledger
    if err := ledger.UpdateToken(token); err != nil {
        return err
    }

    return nil
}

// ProcessRentPayment processes rent payments for a specific token.
func (tm *SYN3000Manager) ProcessRentPayment(tokenID string, amount float64) (*Payment, error) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    token, err := ledger.GetToken(tokenID)
    if err != nil {
        return nil, errors.New("token not found")
    }

    if !token.ActiveStatus {
        return nil, errors.New("lease is inactive")
    }

    // Create a payment record
    payment := Payment{
        PaymentID:   common.GeneratePaymentID(),
        Amount:      amount,
        PaymentDate: time.Now(),
        Status:      "Paid",
    }

    // Append payment to token's payment history
    token.RentPaymentHistory = append(token.RentPaymentHistory, payment)
    token.LastUpdateDate = time.Now()

    // Store updated token in ledger
    if err := ledger.UpdateToken(token); err != nil {
        return nil, err
    }

    return &payment, nil
}

// TerminateLease terminates a rental agreement token, setting it inactive.
func (tm *SYN3000Manager) TerminateLease(tokenID string) error {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    token, err := ledger.GetToken(tokenID)
    if err != nil {
        return errors.New("token not found")
    }

    // Mark lease as inactive
    token.ActiveStatus = false
    token.LastUpdateDate = time.Now()

    // Store the updated token
    if err := ledger.UpdateToken(token); err != nil {
        return err
    }

    return nil
}

// GenerateReport generates a rental report for a property.
func (tm *SYN3000Manager) GenerateReport(propertyID string) (string, error) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    // Retrieve all tokens associated with the property
    tokens, err := ledger.GetTokensByPropertyID(propertyID)
    if err != nil {
        return "", errors.New("property not found")
    }

    // Compile rental report
    report := "Rental Report for Property ID: " + propertyID + "\n"
    for _, token := range tokens {
        report += "Lease Start: " + token.LeaseStartDate.String() + ", Lease End: " + token.LeaseEndDate.String() + "\n"
        report += "Tenant: " + token.TenantID + ", Monthly Rent: " + common.FormatCurrency(token.MonthlyRent) + "\n"
        report += "Deposit: " + common.FormatCurrency(token.DepositAmount) + "\n"
        report += "Payment History: \n"
        for _, payment := range token.RentPaymentHistory {
            report += "  Payment ID: " + payment.PaymentID + ", Amount: " + common.FormatCurrency(payment.Amount) + ", Date: " + payment.PaymentDate.String() + "\n"
        }
        report += "---------------------------------------------\n"
    }

    return report, nil
}

// SYN3000Factory is responsible for creating and managing SYN3000 tokens.
type SYN3000Factory struct{}

// NewSYN3000Factory initializes the token factory.
func NewSYN3000Factory() *SYN3000Factory {
    return &SYN3000Factory{}
}

// CreateRentalToken creates a new SYN3000 token for a rental agreement.
func (f *SYN3000Factory) CreateRentalToken(
    property common.Property,
    ownerID string,
    tenantID string,
    leaseStartDate time.Time,
    leaseEndDate time.Time,
    monthlyRent float64,
    depositAmount float64,
) (*SYN3000Token, error) {
    // Validate inputs
    if ownerID == "" || tenantID == "" || property.PropertyID == "" || monthlyRent <= 0 || depositAmount < 0 {
        return nil, errors.New("invalid input parameters")
    }

    // Create a new SYN3000Token
    token := &SYN3000Token{
        TokenID:            common.GenerateTokenID(),
        PropertyID:         property.PropertyID,
        OwnerID:            ownerID,
        TenantID:           tenantID,
        LeaseStartDate:     leaseStartDate,
        LeaseEndDate:       leaseEndDate,
        MonthlyRent:        monthlyRent,
        DepositAmount:      depositAmount,
        IssuedDate:         time.Now(),
        ActiveStatus:       true,
        LastUpdateDate:     time.Now(),
        PropertyDetails:    property,
        RentPaymentHistory: []common.Payment{},
    }

    // Encrypt sensitive data before storing it in the ledger.
    encryptedToken, err := f.encryptTokenData(token)
    if err != nil {
        return nil, err
    }

    // Store encrypted token in the ledger
    if err := ledger.StoreToken(encryptedToken); err != nil {
        return nil, err
    }

    return token, nil
}

// GetRentalToken retrieves a rental token by its ID.
func (f *SYN3000Factory) GetRentalToken(tokenID string) (*SYN3000Token, error) {
    // Retrieve encrypted token from ledger
    encryptedToken, err := ledger.GetToken(tokenID)
    if err != nil {
        return nil, errors.New("token not found")
    }

    // Decrypt token data
    decryptedToken, err := f.decryptTokenData(encryptedToken)
    if err != nil {
        return nil, err
    }

    return decryptedToken, nil
}

// TerminateRentalToken terminates a rental agreement and sets the token as inactive.
func (f *SYN3000Factory) TerminateRentalToken(tokenID string) error {
    // Retrieve the token
    token, err := f.GetRentalToken(tokenID)
    if err != nil {
        return err
    }

    if !token.ActiveStatus {
        return errors.New("rental agreement is already inactive")
    }

    // Set the token as inactive
    token.ActiveStatus = false
    token.LastUpdateDate = time.Now()

    // Encrypt token data before updating the ledger
    encryptedToken, err := f.encryptTokenData(token)
    if err != nil {
        return err
    }

    // Update the ledger
    if err := ledger.UpdateToken(encryptedToken); err != nil {
        return err
    }

    return nil
}

// Encrypt token data for storage
func (f *SYN3000Factory) encryptTokenData(token *SYN3000Token) (string, error) {
    key := []byte(common.GetEncryptionKey()) // Fetch encryption key from common package
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

// Decrypt token data from storage
func (f *SYN3000Factory) decryptTokenData(encryptedData string) (*SYN3000Token, error) {
    key := []byte(common.GetEncryptionKey()) // Fetch encryption key from common package
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
