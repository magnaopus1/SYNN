package syn10

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewAccessControl initializes a new AccessControl instance.
func newAccessControl(ledgerInstance *ledger.SYN10Ledger, encryptionService *common.Encryption) *SYN10AccessControl {
    return &SYN10AccessControl{
        Ledger:     ledgerInstance,
        Encryption: encryptionService,
        Users:      make(map[string]SYN10User),
    }
}

// RestrictRole ensures that only central bank or government nodes have minting/burning privileges.
func (ac *SYN10AccessControl) RestrictRole(userID string, requiredRole Role) error {
    user, err := ac.GetUserByID(userID)
    if err != nil {
        return err
    }

    if restrictedRoles[user.Role] == false {
        return fmt.Errorf("user %s does not have the required role for this action", userID)
    }

    return nil
}


// AddUser adds a new user to the system, ensuring proper roles are assigned.
func (ac *SYN10AccessControl) AddUser(username, email, password string, role Role) (string, error) {
    userID := common.GenerateUUID()
    passwordHash, err := ac.Encryption.HashPassword(password)
    if err != nil {
        return "", err
    }

    user := SYN10User{
        ID:           userID,
        Username:     username,
        Email:        email,
        PasswordHash: passwordHash,
        Role:         role,
        CreatedAt:    time.Now(),
    }

    ac.Users[userID] = user
    err = ac.Ledger.SaveUser(user)
    if err != nil {
        return "", err
    }

    return userID, nil
}


// GetUserByID retrieves a user by their ID.
func (ac *SYN10AccessControl) GetUserByID(userID string) (*SYN10User, error) {
    user, exists := ac.Users[userID]
    if !exists {
        return nil, errors.New("user not found")
    }
    return &user, nil
}

// VerifyPassword verifies a user's password against the stored hash.
func (ac *SYN10AccessControl) VerifyPassword(password string, hash []byte) (bool, error) {
    return ac.Encryption.ComparePassword(password, hash)
}

// AuthenticateUser authenticates a user by their username and password.
func (ac *SYN10AccessControl) AuthenticateUser(username, password string) (*SYN10User, error) {
    var foundUser *SYN10User
    for _, user := range ac.Users {
        if user.Username == username {
            foundUser = &user
            break
        }
    }

    if foundUser == nil {
        return nil, errors.New("user not found")
    }

    match, err := ac.VerifyPassword(password, foundUser.PasswordHash)
    if err != nil || !match {
        return nil, fmt.Errorf("authentication failed for user: %s", username)
    }

    return foundUser, nil
}


// Authorize checks if a user has the necessary permissions for an action.
func (ac *SYN10AccessControl) Authorize(userID string, requiredRole Role) (bool, error) {
    user, err := ac.GetUserByID(userID)
    if err != nil {
        return false, err
    }

    return user.Role == requiredRole, nil
}


// EncryptData securely encrypts sensitive data using AES encryption.
func (ac *SYN10AccessControl) EncryptData(data []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := common.GenerateRandomBytes(gcm.NonceSize())
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return ciphertext, nil
}


// DecryptData securely decrypts AES-encrypted data.
func (ac *SYN10AccessControl) decryptData(ciphertext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}


// SetUserRole sets the role of a user, updating the ledger.
func (ac *SYN10AccessControl) SetUserRole(userID string, role Role, changedBy string) error {
    user, err := ac.getUserByID(userID)
    if err != nil {
        return err
    }

    user.Role = role
    ac.Users[userID] = *user

    return ac.Ledger.UpdateUserRole(userID, role, changedBy)
}


// RestrictMintingAndBurning ensures only authorized roles can mint or burn tokens.
func (ac *SYN10AccessControl) RestrictMintingAndBurning(userID string) error {
    return ac.RestrictRole(userID, CentralBankRole)
}
