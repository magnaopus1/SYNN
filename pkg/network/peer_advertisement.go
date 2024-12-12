package network

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"sync"
	"time"
)

// PeerAdvertisement contains information about a peer that is being advertised in the network
type PeerAdvertisement struct {
	PeerAddress  string    // IP address or domain of the peer
	PublicKey    string    // Public key of the peer for encryption
	Timestamp    time.Time // Timestamp when the advertisement was created
	Signature    string    // Encrypted signature to ensure authenticity
}

// PeerAdvertiser handles advertising and discovering peers in the network
type PeerAdvertiser struct {
	Peers         map[string]PeerAdvertisement // Known peers in the network
	Mutex         sync.Mutex                   // Mutex for thread-safe access
	Network       *P2PNetwork                  // Pointer to the P2P network
	AdvertiseFreq time.Duration                // Frequency to broadcast peer advertisements
}

// NewPeerAdvertiser initializes a new PeerAdvertiser
func NewPeerAdvertiser(network *P2PNetwork, advertiseFreq time.Duration) *PeerAdvertiser {
	return &PeerAdvertiser{
		Peers:         make(map[string]PeerAdvertisement),
		Network:       network,
		AdvertiseFreq: advertiseFreq,
	}
}

// AdvertisePeer broadcasts the current node's peer information to all connected peers
func (pa *PeerAdvertiser) AdvertisePeer() {
    for {
        pa.Mutex.Lock()

        // Step 1: Encode the public key (converting *rsa.PublicKey to string)
        encodedPublicKey, err := EncodePublicKey(pa.Network.NodeKey.PublicKey)
        if err != nil {
            fmt.Printf("Failed to encode public key: %v\n", err)
            pa.Mutex.Unlock()
            return
        }

        // Step 2: Create the peer advertisement message
        advertisement := PeerAdvertisement{
            PeerAddress: pa.Network.Address,  // Assuming Network.Address stores the peer's IP/domain
            PublicKey:   encodedPublicKey,    // Use the encoded public key
            Timestamp:   time.Now(),
        }

        // Step 3: Generate a signature to authenticate the advertisement
        advertisement.Signature = pa.generateSignature(advertisement)

        // Step 4: Broadcast the advertisement to all peers
        for peerAddress := range pa.Network.Peers {
            err := pa.Network.SendMessage(peerAddress, pa.serializeAdvertisement(advertisement))
            if err != nil {
                fmt.Printf("Failed to advertise peer to %s: %v\n", peerAddress, err)
            }
        }

        pa.Mutex.Unlock()

        // Wait before advertising again
        time.Sleep(pa.AdvertiseFreq)
    }
}


// EncodePublicKey converts an *rsa.PublicKey to a PEM-encoded string
func EncodePublicKey(pubKey *rsa.PublicKey) (string, error) {
    // Convert the public key to ASN.1 DER-encoded format
    pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
    if err != nil {
        return "", fmt.Errorf("error encoding public key: %v", err)
    }

    // Create a PEM block with the DER-encoded public key
    pemBlock := &pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: pubKeyBytes,
    }

    // Encode the PEM block into a string
    return string(pem.EncodeToMemory(pemBlock)), nil
}


// ReceivePeerAdvertisement processes incoming peer advertisements and adds them to the peer list
func (pa *PeerAdvertiser) ReceivePeerAdvertisement(advertisementStr string) {
	pa.Mutex.Lock()
	defer pa.Mutex.Unlock()

	// Deserialize the incoming advertisement
	advertisement := pa.deserializeAdvertisement(advertisementStr)

	// Validate the advertisement's signature
	if !pa.validateAdvertisement(advertisement) {
		fmt.Println("Invalid peer advertisement received, ignoring.")
		return
	}

	// Add the peer to the list if it doesn't already exist
	if _, exists := pa.Peers[advertisement.PeerAddress]; !exists {
		pa.Peers[advertisement.PeerAddress] = advertisement
		fmt.Printf("Added new peer: %s\n", advertisement.PeerAddress)
	}
}

// generateSignature generates an RSA signature for a peer advertisement
func (pa *PeerAdvertiser) generateSignature(advertisement PeerAdvertisement) string {
    // Hash the peer address and timestamp
    hashInput := fmt.Sprintf("%s%d", advertisement.PeerAddress, advertisement.Timestamp.UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    signatureHash := hash.Sum(nil)

    // Sign the hash with the node's private key using RSA
    signature, err := rsa.SignPKCS1v15(rand.Reader, pa.Network.NodeKey.PrivateKey, crypto.SHA256, signatureHash)
    if err != nil {
        fmt.Printf("Error signing advertisement: %v\n", err)
        return ""
    }

    // Return the signature as a hex-encoded string
    return hex.EncodeToString(signature)
}



// parseRSAPublicKey parses a base64-encoded RSA public key string into an *rsa.PublicKey
func parseRSAPublicKey(pubKeyStr string) (*rsa.PublicKey, error) {
    // Decode the public key from base64
    pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyStr)
    if err != nil {
        return nil, fmt.Errorf("failed to decode public key: %v", err)
    }

    // Parse the public key bytes
    pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse public key: %v", err)
    }

    // Ensure the parsed key is an RSA public key
    rsaPubKey, ok := pubKey.(*rsa.PublicKey)
    if !ok {
        return nil, fmt.Errorf("not an RSA public key")
    }

    return rsaPubKey, nil
}

// validateAdvertisement validates the RSA signature of a peer advertisement
func (pa *PeerAdvertiser) validateAdvertisement(advertisement PeerAdvertisement) bool {
    // Parse the public key from the advertisement
    pubKey, err := parseRSAPublicKey(advertisement.PublicKey)
    if err != nil {
        fmt.Printf("Error parsing public key: %v\n", err)
        return false
    }

    // Hash the peer address and timestamp
    hashInput := fmt.Sprintf("%s%d", advertisement.PeerAddress, advertisement.Timestamp.UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    expectedHash := hash.Sum(nil)

    // Decode the encrypted signature from hex
    signature, err := hex.DecodeString(advertisement.Signature)
    if err != nil {
        fmt.Printf("Error decoding signature: %v\n", err)
        return false
    }

    // Use the public key to verify the RSA signature
    err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, expectedHash, signature)
    if err != nil {
        fmt.Printf("Signature verification failed: %v\n", err)
        return false
    }

    return true
}



// serializeAdvertisement converts a peer advertisement to a string format for transmission
func (pa *PeerAdvertiser) serializeAdvertisement(advertisement PeerAdvertisement) string {
	return fmt.Sprintf("%s|%s|%d|%s", advertisement.PeerAddress, advertisement.PublicKey, advertisement.Timestamp.UnixNano(), advertisement.Signature)
}

// deserializeAdvertisement converts a string back into a PeerAdvertisement struct
func (pa *PeerAdvertiser) deserializeAdvertisement(advertisementStr string) PeerAdvertisement {
	var advertisement PeerAdvertisement
	fmt.Sscanf(advertisementStr, "%s|%s|%d|%s", &advertisement.PeerAddress, &advertisement.PublicKey, &advertisement.Timestamp, &advertisement.Signature)
	advertisement.Timestamp = time.Unix(0, advertisement.Timestamp.UnixNano())
	return advertisement
}

// GetPeers returns the list of known peers in the network
func (pa *PeerAdvertiser) GetPeers() map[string]PeerAdvertisement {
	pa.Mutex.Lock()
	defer pa.Mutex.Unlock()
	return pa.Peers
}

// serializeAdvertisement converts a peer advertisement to a string format for transmission
func (pa *PeerAdvertiser) SerializeAdvertisement(advertisement PeerAdvertisement) string {
	return fmt.Sprintf("%s|%s|%d|%s", advertisement.PeerAddress, advertisement.PublicKey, advertisement.Timestamp.UnixNano(), advertisement.Signature)
}

// deserializeAdvertisement converts a string back into a PeerAdvertisement struct
func (pa *PeerAdvertiser) DeserializeAdvertisement(advertisementStr string) PeerAdvertisement {
	var advertisement PeerAdvertisement
	fmt.Sscanf(advertisementStr, "%s|%s|%d|%s", &advertisement.PeerAddress, &advertisement.PublicKey, &advertisement.Timestamp, &advertisement.Signature)
	advertisement.Timestamp = time.Unix(0, advertisement.Timestamp.UnixNano())
	return advertisement
}
