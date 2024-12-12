package network

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "net"
    "synnergy_network/pkg/ledger"
)

// NewHandshake initializes the handshake mechanism by generating a new key pair
func NewHandshake(ledger *ledger.Ledger) *Handshake {
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        fmt.Println("Error generating private key:", err)
        return nil
    }

    return &Handshake{
        PrivateKey: privateKey,
        PublicKey:  &privateKey.PublicKey,
        ledger:     ledger,
    }
}

// SendPublicKey sends the node's public key to the remote peer during the handshake
func (hs *Handshake) SendPublicKey(conn net.Conn) error {
    hs.mutex.Lock()
    defer hs.mutex.Unlock()

    // Marshal the public key to PEM format
    pubKeyBytes := x509.MarshalPKCS1PublicKey(hs.PublicKey)
    pemBlock := &pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: pubKeyBytes,
    }
    pemData := pem.EncodeToMemory(pemBlock)

    // Send the public key to the peer
    _, err := conn.Write(pemData)
    if err != nil {
        return fmt.Errorf("failed to send public key: %v", err)
    }

    fmt.Println("Public key sent to peer.")

    // Log the event (passing only one argument as expected by the ledger's LogConnectionEvent method)
    hs.ledger.LogConnectionEvent("PublicKey Sent")
    return nil
}


// ReceivePublicKey receives and stores the remote peer's public key during the handshake
func (hs *Handshake) ReceivePublicKey(conn net.Conn) (*rsa.PublicKey, error) {
    hs.mutex.Lock()
    defer hs.mutex.Unlock()

    buffer := make([]byte, 4096)
    n, err := conn.Read(buffer)
    if err != nil {
        return nil, fmt.Errorf("failed to receive public key: %v", err)
    }

    pemBlock, _ := pem.Decode(buffer[:n])
    if pemBlock == nil {
        return nil, fmt.Errorf("failed to decode PEM block containing public key")
    }

    pubKey, err := x509.ParsePKCS1PublicKey(pemBlock.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse public key: %v", err)
    }

    fmt.Println("Public key received from peer.")

    // Log the event (passing only one argument as expected by the ledger's LogConnectionEvent method)
    hs.ledger.LogConnectionEvent("PublicKey Received")
    return pubKey, nil
}


// EncryptMessage encrypts a message using the recipient's public key
func (hs *Handshake) EncryptMessage(message string, pubKey *rsa.PublicKey) ([]byte, error) {
    label := []byte("")
    hash := sha256.New()

    ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pubKey, []byte(message), label)
    if err != nil {
        return nil, fmt.Errorf("error encrypting message: %v", err)
    }

    fmt.Println("Message encrypted successfully.")
    return ciphertext, nil
}

// DecryptMessage decrypts a received message using the node's private key
func (hs *Handshake) DecryptMessage(ciphertext []byte) (string, error) {
    label := []byte("")
    hash := sha256.New()

    plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, hs.PrivateKey, ciphertext, label)
    if err != nil {
        return "", fmt.Errorf("error decrypting message: %v", err)
    }

    fmt.Println("Message decrypted successfully.")
    return string(plaintext), nil
}

// PerformHandshake establishes a secure connection by exchanging public keys and encrypting the communication
func (hs *Handshake) PerformHandshake(conn net.Conn) error {
    // Send the node's public key
    if err := hs.SendPublicKey(conn); err != nil {
        return err
    }

    // Receive the remote node's public key
    remotePubKey, err := hs.ReceivePublicKey(conn)
    if err != nil {
        return err
    }

    // Verify communication with encrypted test message
    testMessage := "synnergy_test"
    encryptedMessage, err := hs.EncryptMessage(testMessage, remotePubKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt test message: %v", err)
    }

    if _, err := conn.Write(encryptedMessage); err != nil {
        return fmt.Errorf("failed to send encrypted test message: %v", err)
    }

    // Log the handshake success (passing only one argument)
    hs.ledger.LogConnectionEvent("Handshake Completed")
    fmt.Println("Handshake completed successfully with", conn.RemoteAddr().String())
    return nil
}

