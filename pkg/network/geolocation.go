package network

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewGeoLocationManager initializes the geolocation manager with node locations
func NewGeoLocationManager(ledger *ledger.Ledger) *GeoLocationManager {
    return &GeoLocationManager{
        NodeLocations: make(map[string]GeoLocation),
        ledger:        ledger,
    }
}

// RegisterNodeLocation registers the geolocation of a node
func (gl *GeoLocationManager) RegisterNodeLocation(nodeID string, location GeoLocation) {
    gl.NodeLocations[nodeID] = location
    fmt.Printf("Node %s registered at location: (%f, %f)\n", nodeID, location.Latitude, location.Longitude)

    // Log the geo event as a single string
    logMessage := fmt.Sprintf("Node %s registered at location: (%f, %f) at %v", nodeID, location.Latitude, location.Longitude, time.Now())
    gl.ledger.LogGeoEvent(logMessage)
}

// CalculateDistance calculates the great-circle distance between two geolocations
func (gl *GeoLocationManager) CalculateDistance(loc1, loc2 GeoLocation) float64 {
    const earthRadiusKm = 6371.0 // Earth's radius in kilometers

    latDiff := degreesToRadians(loc2.Latitude - loc1.Latitude)
    lonDiff := degreesToRadians(loc2.Longitude - loc1.Longitude)

    a := math.Sin(latDiff/2)*math.Sin(latDiff/2) +
        math.Cos(degreesToRadians(loc1.Latitude))*math.Cos(degreesToRadians(loc2.Latitude))*
            math.Sin(lonDiff/2)*math.Sin(lonDiff/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    return earthRadiusKm * c
}

// EncryptGeoData encrypts geolocation data before broadcasting or storing
func (gl *GeoLocationManager) EncryptGeoData(location GeoLocation, pubKey *common.PublicKey) (string, error) {
    geoString := fmt.Sprintf("%f,%f", location.Latitude, location.Longitude)
    
    // Convert geolocation data to bytes
    geoBytes := []byte(geoString)

    // Encrypt the geolocation data using the provided public key
    encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey.Key, geoBytes, nil)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt geo data: %v", err)
    }

    // Encode encrypted data to a hex string
    encryptedGeoData := hex.EncodeToString(encryptedBytes)

    fmt.Println("Geolocation data encrypted for secure transmission.")
    return encryptedGeoData, nil
}

// DecryptGeoData decrypts the encrypted geolocation data for use in the system
func (gl *GeoLocationManager) DecryptGeoData(encryptedData string) (GeoLocation, error) {
    geoBytes, err := hex.DecodeString(encryptedData)
    if err != nil {
        return GeoLocation{}, fmt.Errorf("failed to decrypt geolocation data: %v", err)
    }

    geoString := string(geoBytes)
    var location GeoLocation
    _, err = fmt.Sscanf(geoString, "%f,%f", &location.Latitude, &location.Longitude)
    if err != nil {
        return GeoLocation{}, fmt.Errorf("failed to parse decrypted geolocation data: %v", err)
    }

    return location, nil
}

// FindClosestNode identifies the closest node to a given location
func (gl *GeoLocationManager) FindClosestNode(targetLocation GeoLocation) (string, float64) {
    closestNode := ""
    closestDistance := math.MaxFloat64

    for nodeID, location := range gl.NodeLocations {
        distance := gl.CalculateDistance(location, targetLocation)
        if distance < closestDistance {
            closestNode = nodeID
            closestDistance = distance
        }
    }

    fmt.Printf("Closest node to location (%f, %f) is node %s at a distance of %.2f km.\n",
        targetLocation.Latitude, targetLocation.Longitude, closestNode, closestDistance)

    return closestNode, closestDistance
}

// BroadcastNodeLocationChange broadcasts any significant node location changes
func (gl *GeoLocationManager) BroadcastNodeLocationChange(nodeID string, oldLocation, newLocation GeoLocation) {
    distanceMoved := gl.CalculateDistance(oldLocation, newLocation)
    if distanceMoved > 10.0 { // Arbitrary threshold to broadcast changes
        fmt.Printf("Node %s has moved %.2f km. Broadcasting location update.\n", nodeID, distanceMoved)
        
        // Log the geo event as a single string to match the LogGeoEvent function's expected argument
        logMessage := fmt.Sprintf("Node %s changed location to %v (moved %.2f km)", nodeID, newLocation, distanceMoved)
        gl.ledger.LogGeoEvent(logMessage)
    }
}


// degreesToRadians converts degrees to radians
func degreesToRadians(degrees float64) float64 {
    return degrees * math.Pi / 180
}
