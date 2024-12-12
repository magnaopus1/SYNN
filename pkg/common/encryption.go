package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
)

// Encryption struct handles all encryption-related tasks
type Encryption struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewEncryption initializes a new Encryption struct with a private key
func NewEncryption(bits int) (*Encryption, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return &Encryption{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// EncryptRSA encrypts data using RSA and the public key
func (e *Encryption) EncryptRSA(plainText []byte) ([]byte, error) {
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, e.publicKey, plainText, nil)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

// DecryptRSA decrypts data using RSA and the private key
func (e *Encryption) DecryptRSA(cipherText []byte) ([]byte, error) {
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, e.privateKey, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// EncryptAES encrypts data using AES algorithm
func (e *Encryption) EncryptAES(key []byte, plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

// DecryptAES decrypts AES-encrypted data
func (e *Encryption) DecryptAES(key []byte, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("cipherText too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

// EncryptionKey must be 16, 24, or 32 bytes long
var EncryptionKey = []byte("this-is-32-byte-key-123456789012") // 32 bytes


// EncryptData encrypts the data payload using the provided algorithm and key
func (e *Encryption) EncryptData(algorithm string, data []byte, key []byte) ([]byte, error) {
    switch algorithm {
    case "AES":
        // Validate key length
        keyLength := len(key)
        if keyLength != 16 && keyLength != 24 && keyLength != 32 {
            return nil, fmt.Errorf("invalid AES key size: %d bytes", keyLength)
        }

        block, err := aes.NewCipher(key)
        if err != nil {
            return nil, fmt.Errorf("failed to create cipher: %v", err)
        }

        cipherText := make([]byte, aes.BlockSize+len(data))
        iv := cipherText[:aes.BlockSize]

        if _, err := io.ReadFull(rand.Reader, iv); err != nil {
            return nil, fmt.Errorf("failed to generate IV: %v", err)
        }

        stream := cipher.NewCFBEncrypter(block, iv)
        stream.XORKeyStream(cipherText[aes.BlockSize:], data)

        return cipherText, nil
    default:
        return nil, errors.New("unsupported encryption algorithm")
    }
}



// DecryptData decrypts the data using AES decryption
func (e *Encryption) DecryptData(data []byte, key []byte) ([]byte, error) {
    // Validate key length
    keyLength := len(key)
    if keyLength != 16 && keyLength != 24 && keyLength != 32 {
        return nil, fmt.Errorf("invalid AES key size: %d bytes", keyLength)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %v", err)
    }

    if len(data) < aes.BlockSize {
        return nil, fmt.Errorf("ciphertext too short")
    }

    iv := data[:aes.BlockSize]
    ciphertext := data[aes.BlockSize:]

    // Decrypt using CFB (Cipher Feedback Mode)
    stream := cipher.NewCFBDecrypter(block, iv)
    decryptedData := make([]byte, len(ciphertext))
    stream.XORKeyStream(decryptedData, ciphertext)

    return decryptedData, nil
}



// ExportPrivateKey exports the private key in PEM format
func (e *Encryption) ExportPrivateKey() ([]byte, error) {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(e.privateKey)
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	})
	return privKeyPEM, nil
}

// ExportPublicKey exports the public key in PEM format
func (e *Encryption) ExportPublicKey() ([]byte, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(e.publicKey)
	if err != nil {
		return nil, err
	}
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})
	return pubKeyPEM, nil
}

// EncryptPost encrypts the given ForumPost object using AES encryption
func EncryptPost(post interface{}, key []byte) ([]byte, error) {
	plainText, err := json.Marshal(post) // Serialize the post to JSON
	if err != nil {
		return nil, fmt.Errorf("failed to serialize post: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

// EncryptFeedback encrypts the feedback object using AES encryption
func EncryptFeedback(feedback interface{}, key []byte) ([]byte, error) {
	plainText, err := json.Marshal(feedback) // Serialize the feedback to JSON
	if err != nil {
		return nil, fmt.Errorf("failed to serialize feedback: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

// PublicKey represents a wrapper for RSA public key encryption
type PublicKey struct {
    Key *rsa.PublicKey
}

// NodeKey represents the private and public key pair for the node
type NodeKey struct {
    PrivateKey *rsa.PrivateKey
    PublicKey  *rsa.PublicKey
}

// GenerateNodeKeyPair generates a new RSA key pair for a node
func GenerateNodeKeyPair() (*NodeKey, error) {
    // Generate an RSA private key
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, fmt.Errorf("failed to generate private key: %v", err)
    }

    // Extract the public key from the private key
    publicKey := &privateKey.PublicKey

    // Return the key pair
    return &NodeKey{
        PrivateKey: privateKey,
        PublicKey:  publicKey,
    }, nil
}

// EncryptContract encrypts the smart contract data using the provided encryption key.
func EncryptContract(contract *SmartContract, key []byte) ([]byte, error) {
    // Marshal the smart contract data to JSON
    contractData, err := json.Marshal(contract)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal contract: %v", err)
    }

    // Create an instance of the Encryption struct
    encryptionInstance := &Encryption{}

    // Use the AES encryption method to encrypt the contract data
    encryptedData, err := encryptionInstance.EncryptData("AES", contractData, key)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt contract data: %v", err)
    }

    return encryptedData, nil
}



// EncryptContractExecution encrypts the contract execution data using AES encryption
func EncryptContractExecution(execution interface{}, key []byte) ([]byte, error) {
    // Convert the execution data into a JSON string
    executionData, err := json.Marshal(execution)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal execution data: %v", err)
    }

    // Create a new AES cipher using the provided key
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %v", err)
    }

    // Generate a new IV (initialization vector)
    cipherText := make([]byte, aes.BlockSize+len(executionData))
    iv := cipherText[:aes.BlockSize]

    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, fmt.Errorf("failed to generate IV: %v", err)
    }

    // Encrypt the execution data using CFB (Cipher Feedback Mode)
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(cipherText[aes.BlockSize:], executionData)

    return cipherText, nil
}

// EncryptContractState encrypts the given contract state using the provided encryption key.
// The state is first serialized to JSON, then encrypted using AES.
func EncryptContractState(state map[string]interface{}, key []byte) ([]byte, error) {
	// Serialize the state to JSON
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	// Encrypt the serialized state using AES encryption
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(stateJSON))
	iv := ciphertext[:aes.BlockSize] // initialization vector (IV)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], stateJSON)

	return ciphertext, nil
}

// EncryptPrivateKey encrypts the private key using AES encryption
func EncryptPrivateKey(privateKey *ecdsa.PrivateKey, encryptionKey []byte) ([]byte, error) {
	keyBytes := privateKey.D.Bytes() // Convert private key to bytes
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(keyBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], keyBytes)
	return cipherText, nil
}

// EncodeMessage encrypts a message using AES and encodes it to base64
func EncodeMessage(message string) (string, error) {
    block, err := aes.NewCipher(EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %v", err)
    }

    plaintext := []byte(message)

    // Generate a new AES GCM cipher
    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %v", err)
    }

    // Create a nonce for AES GCM
    nonce := make([]byte, aesGCM.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", fmt.Errorf("failed to generate nonce: %v", err)
    }

    // Encrypt the plaintext
    ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

    // Return the encrypted data as a base64 string
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecodeMessage decrypts a base64 encoded and AES-encrypted message
func DecodeMessage(encodedMessage string) (string, error) {
    block, err := aes.NewCipher(EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %v", err)
    }

    // Decode the base64 encoded message
    ciphertext, err := base64.StdEncoding.DecodeString(encodedMessage)
    if err != nil {
        return "", fmt.Errorf("failed to decode base64 message: %v", err)
    }

    // Generate a new AES GCM cipher
    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %v", err)
    }

    nonceSize := aesGCM.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

    // Decrypt the ciphertext
    plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt message: %v", err)
    }

    return string(plaintext), nil
}


// EncodePublicKey encodes an RSA public key into a base64 string
func (e *Encryption) EncodePublicKey(pub *rsa.PublicKey) (string, error) {
	// Marshal the RSA public key into ASN.1 PKIX format
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %v", err)
	}

	// Encode the public key bytes to a base64 string
	encodedPubKey := base64.StdEncoding.EncodeToString(pubBytes)
	return encodedPubKey, nil
}

// DecodePublicKey decodes a base64-encoded string into an *rsa.PublicKey
func (e *Encryption) DecodePublicKey(encodedKey string) (*rsa.PublicKey, error) {
	// Decode the base64 string
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, errors.New("failed to decode base64 public key: " + err.Error())
	}

	// Decode the PEM block
	block, _ := pem.Decode(decodedKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing the public key")
	}

	// Parse the public key
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse public key: " + err.Error())
	}

	// Assert that the key is an RSA public key
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPublicKey, nil
}

// EncryptWithKey encrypts data using AES encryption with the provided key.
func EncryptWithKey(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Generate a new AES-GCM cipher for encryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a random nonce for encryption.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the data and append the nonce to the beginning of the result.
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// DecryptWithKey decrypts data using AES encryption with the provided key.
func DecryptWithKey(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Generate a new AES-GCM cipher for decryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Separate the nonce and ciphertext from the encrypted data.
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt the ciphertext using AES-GCM.
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}