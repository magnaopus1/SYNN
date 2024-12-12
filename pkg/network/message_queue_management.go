package network

import (
	"container/list"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewMessageQueue initializes a new message queue with a maximum size
func NewMessageQueue(maxSize int, ledger *ledger.Ledger) *MessageQueue {
	return &MessageQueue{
		queue:          list.New(),
		maxQueueSize:   maxSize,
		ledgerInstance: ledger,
	}
}

// AddMessage adds a new message to the queue
func (mq *MessageQueue) AddMessage(from, to, content string, encrypt bool) (string, error) {
    mq.lock.Lock()
    defer mq.lock.Unlock()

    if mq.queue.Len() >= mq.maxQueueSize {
        return "", fmt.Errorf("message queue is full")
    }

    // Encrypt the content if requested
    var encryptedContent string
    var err error
    if encrypt {
        encryptedContent, err = mq.encryptContent(content)
        if err != nil {
            return "", err
        }
    } else {
        encryptedContent = content
    }

    // Create the message with a unique ID and SHA-256 hash of the content
    messageID := generateMessageID(from, to, content)
    hash := hashMessageContent(content)

    // Add the message to the queue
    message := Message{
        ID:        messageID,
        Timestamp: time.Now(),
        From:      from,
        To:        to,
        Content:   encryptedContent,
        Hash:      hash,
        Encrypted: encrypt,
    }
    mq.queue.PushBack(message)

    // Log the event with only one argument (e.g., message ID)
    mq.ledgerInstance.LogMessageEvent("MessageAdded")

    fmt.Printf("Message added from %s to %s: %s\n", from, to, message.Content)
    return messageID, nil
}


// ProcessMessage processes the first message in the queue and removes it
func (mq *MessageQueue) ProcessMessage() (*Message, error) {
    mq.lock.Lock()
    defer mq.lock.Unlock()

    if mq.queue.Len() == 0 {
        return nil, fmt.Errorf("no messages to process")
    }

    // Get the first message in the queue
    elem := mq.queue.Front()
    message := elem.Value.(Message)

    // Remove the message from the queue after processing
    mq.queue.Remove(elem)

    // Log the event in the ledger with only one argument (e.g., the event type)
    mq.ledgerInstance.LogMessageEvent("MessageProcessed")

    fmt.Printf("Processed message from %s to %s\n", message.From, message.To)
    return &message, nil
}


// VerifyMessageHash verifies the integrity of the message content
func (mq *MessageQueue) VerifyMessageHash(message Message) bool {
	expectedHash := hashMessageContent(message.Content)
	if message.Hash != expectedHash {
		fmt.Println("Message hash verification failed.")
		return false
	}
	fmt.Println("Message hash verified successfully.")
	return true
}

// EncryptMessageContent encrypts the message content using the receiver's public key
func (mq *MessageQueue) encryptContent(content string) (string, error) {
    pubKey := GetNetworkPublicKey() // Retrieve the public key
    encryptedContent, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(content))
    if err != nil {
        return "", fmt.Errorf("encryption failed: %v", err)
    }
    return hex.EncodeToString(encryptedContent), nil
}

// DecryptMessageContent decrypts an encrypted message using the node's private key
func (mq *MessageQueue) DecryptMessageContent(encryptedContent string) (string, error) {
    privKey := GetNodePrivateKey() // Retrieve the private key
    contentBytes, err := hex.DecodeString(encryptedContent)
    if err != nil {
        return "", fmt.Errorf("failed to decode encrypted content: %v", err)
    }

    decryptedContent, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, contentBytes)
    if err != nil {
        return "", fmt.Errorf("decryption failed: %v", err)
    }
    return string(decryptedContent), nil
}

// generateMessageID generates a unique message ID based on sender, receiver, and content
func generateMessageID(from, to, content string) string {
	return fmt.Sprintf("%s_%s_%s", from, to, hashMessageContent(content))
}

// hashMessageContent calculates the SHA-256 hash of the message content
func hashMessageContent(content string) string {
	hash := sha256.New()
	hash.Write([]byte(content))
	return hex.EncodeToString(hash.Sum(nil))
}

// GetNetworkPublicKey retrieves the RSA public key from a file or environment
func GetNetworkPublicKey() *rsa.PublicKey {
    pubKeyFile, err := ioutil.ReadFile("path/to/public_key.pem")
    if err != nil {
        log.Fatalf("Failed to read public key: %v", err)
    }

    block, _ := pem.Decode(pubKeyFile)
    if block == nil || block.Type != "PUBLIC KEY" {
        log.Fatalf("Failed to decode PEM block containing public key")
    }

    pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
    if err != nil {
        log.Fatalf("Failed to parse public key: %v", err)
    }

    return pubKey
}

// GetNodePrivateKey retrieves the RSA private key from a file or environment
func GetNodePrivateKey() *rsa.PrivateKey {
    privKeyFile, err := ioutil.ReadFile("path/to/private_key.pem")
    if err != nil {
        log.Fatalf("Failed to read private key: %v", err)
    }

    block, _ := pem.Decode(privKeyFile)
    if block == nil || block.Type != "PRIVATE KEY" {
        log.Fatalf("Failed to decode PEM block containing private key")
    }

    privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        log.Fatalf("Failed to parse private key: %v", err)
    }

    return privKey
}