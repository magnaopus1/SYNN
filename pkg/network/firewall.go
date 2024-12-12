package network

import (
	"fmt"
	"net"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewFirewall initializes a new Firewall with a whitelist and block duration
func NewFirewall(allowedIPs []string, blockDuration time.Duration, ledgerInstance *ledger.Ledger) *Firewall {
	ipWhitelist := make(map[string]bool)
	for _, ip := range allowedIPs {
		ipWhitelist[ip] = true
	}

	return &Firewall{
		allowedIPs:     ipWhitelist,
		blockedIPs:     make(map[string]time.Time),
		blockDuration:  blockDuration,
		ledgerInstance: ledgerInstance,
	}
}

// AllowIP adds an IP address to the normal firewall's whitelist
func (fw *Firewall) AllowIP(ip string) {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    fw.allowedIPs[ip] = true
    fmt.Printf("IP %s added to normal firewall whitelist.\n", ip)

    // Create an encryption instance
    enc := &common.Encryption{}

    // Encrypt the event message
    encryptedEvent, err := enc.EncryptData("AES", []byte(fmt.Sprintf("IP Whitelisted: %s", ip)), common.EncryptionKey)
    if err != nil {
        fmt.Println("Failed to encrypt event:", err)
        return
    }

    // Log the event to the ledger with the encrypted message
    fw.ledgerInstance.LogFirewallEvent(string(encryptedEvent))  // Assuming LogFirewallEvent takes only one argument
}


// BlockIP blocks an IP address in the normal firewall
func (fw *Firewall) BlockIP(ip string) {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    fw.blockedIPs[ip] = time.Now()
    fmt.Printf("IP %s blocked in normal firewall.\n", ip)

    // Create an encryption instance
    enc := &common.Encryption{}

    // Log the event to the ledger with encryption
    encryptedEvent, err := enc.EncryptData("AES", []byte(fmt.Sprintf("IP Blocked: %s", ip)), common.EncryptionKey)
    if err != nil {
        fmt.Println("Failed to encrypt event:", err)
        return
    }

    // Log only the encrypted message in the ledger
    fw.ledgerInstance.LogFirewallEvent(string(encryptedEvent))  // Assuming LogFirewallEvent takes only one argument
}

// UnblockIP removes an IP from the blocked list after the block duration has passed
func (fw *Firewall) UnblockIP(ip string) {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    delete(fw.blockedIPs, ip)
    fmt.Printf("IP %s unblocked by normal firewall.\n", ip)

    // Create an encryption instance
    enc := &common.Encryption{}

    // Encrypt the event message
    encryptedEvent, err := enc.EncryptData("AES", []byte(fmt.Sprintf("IP Unblocked: %s", ip)), common.EncryptionKey)
    if err != nil {
        fmt.Println("Failed to encrypt event:", err)
        return
    }

    // Log the event to the ledger (assuming LogFirewallEvent takes one string argument)
    fw.ledgerInstance.LogFirewallEvent(string(encryptedEvent))  // Log only the encrypted message
}

// IsAllowed checks if an IP is allowed to interact with the network
func (fw *Firewall) IsAllowed(ip string) bool {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	// Check if the IP is in the blocked list and if the block duration has expired
	if blockedTime, blocked := fw.blockedIPs[ip]; blocked {
		if time.Since(blockedTime) < fw.blockDuration {
			fmt.Printf("IP %s is blocked and cannot access the network.\n", ip)
			return false
		}

		// Unblock the IP if the block duration has passed
		fw.UnblockIP(ip)
	}

	// Check if the IP is in the whitelist
	if allowed, exists := fw.allowedIPs[ip]; exists {
		return allowed
	}

	// Default behavior: block all IPs not on the whitelist
	fmt.Printf("IP %s is not whitelisted and cannot access the network.\n", ip)
	return false
}

// MonitorConnections monitors incoming connections for the normal firewall
func (fw *Firewall) MonitorConnections(connPool *ConnectionPool) {
    go func() {
        for {
            time.Sleep(5 * time.Second) // Monitor every 5 seconds

            fw.mutex.Lock()
            for _, conn := range connPool.ActiveConns { // Use ActiveConns as the field for connections
                ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
                if err != nil {
                    fmt.Println("Error retrieving IP:", err)
                    continue
                }

                if !fw.IsAllowed(ip) {
                    fmt.Printf("Connection from %s closed by normal firewall.\n", ip)
                    conn.Close()

                    // Log the connection termination to the ledger
                    enc := &common.Encryption{}
                    encryptedEvent, err := enc.EncryptData("AES", []byte(fmt.Sprintf("Connection Terminated: %s", ip)), common.EncryptionKey)
                    if err != nil {
                        fmt.Println("Failed to encrypt event:", err)
                        continue
                    }

                    // Log the encrypted event to the ledger (single argument)
                    fw.ledgerInstance.LogFirewallEvent(string(encryptedEvent)) // Only pass the encrypted string
                }
            }
            fw.mutex.Unlock()
        }
    }()
}

// LogSuspiciousActivity logs any suspicious activities, such as repeated connection attempts from a blocked IP
func (fw *Firewall) LogSuspiciousActivity(ip string, reason string) {
    fmt.Printf("Suspicious activity detected from IP %s: %s\n", ip, reason)

    // Prepare the log message
    logMessage := fmt.Sprintf("Suspicious Activity: IP %s, Reason: %s, Time: %s", ip, reason, time.Now().Format(time.RFC3339))

    // Log the suspicious activity to the ledger (single string argument)
    fw.ledgerInstance.LogFirewallEvent(logMessage)
}



// InitializeFirewallManager initializes the Firewall Manager with all four types of firewalls
func InitializeFirewallManager(ledgerInstance *ledger.Ledger, blockDuration time.Duration, allowedIPs []string) *FirewallManager {
	// Initialize the normal firewall with the whitelist of allowed IPs
	normalFirewall := &Firewall{
		allowedIPs:     make(map[string]bool),
		blockedIPs:     make(map[string]time.Time),
		blockDuration:  blockDuration,
		ledgerInstance: ledgerInstance,
	}

	// Populate the whitelist
	for _, ip := range allowedIPs {
		normalFirewall.allowedIPs[ip] = true
	}

	return &FirewallManager{
		NormalFirewall: normalFirewall,
		DynamicFirewall: &DynamicFirewall{
			allowedIPs:     make(map[string]bool),
			blockedIPs:     make(map[string]time.Time),
			blockDuration:  blockDuration,
			ledgerInstance: ledgerInstance,
		},
		StatelessFirewall: &StatelessFirewall{
			allowedPorts:   []int{80, 443}, // Allow HTTP and HTTPS by default
			blockedIPs:     make(map[string]time.Time),
			ledgerInstance: ledgerInstance,
		},
		StatefulFirewall: &StatefulFirewall{
			allowedConnections: make(map[string]string),
			blockedIPs:         make(map[string]time.Time),
			ledgerInstance:     ledgerInstance,
		},
		ledgerInstance: ledgerInstance,
	}
}

// AllowIP for the dynamic firewall (similar to the normal firewall but adaptive logic could be added)
func (fw *DynamicFirewall) AllowIP(ip string) {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    fw.allowedIPs[ip] = true
    fmt.Printf("Dynamic Firewall: IP %s allowed.\n", ip)

    // Encrypt the event message (make sure the encryption package is properly set up)
    encryptedEvent, err := fw.encryptionService.EncryptData("AES", []byte(fmt.Sprintf("IP Whitelisted: %s", ip)), common.EncryptionKey)
    if err != nil {
        fmt.Println("Error encrypting event:", err)
        return
    }

    // Log the event to the ledger (single string argument)
    fw.ledgerInstance.LogFirewallEvent(string(encryptedEvent))
}

// BlockIP for the dynamic firewall
func (fw *DynamicFirewall) BlockIP(ip string) {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    fw.blockedIPs[ip] = time.Now()
    fmt.Printf("Dynamic Firewall: IP %s blocked.\n", ip)

    // Encrypt the event message (make sure the encryption package is properly set up)
    encryptedEvent, err := fw.encryptionService.EncryptData("AES", []byte(fmt.Sprintf("IP Blocked: %s", ip)), common.EncryptionKey)
    if err != nil {
        fmt.Println("Error encrypting event:", err)
        return
    }

    // Log the event to the ledger (single string argument)
    fw.ledgerInstance.LogFirewallEvent(string(encryptedEvent))
}


// UnblockIP removes an IP from the blocked list in the dynamic firewall
func (fw *DynamicFirewall) UnblockIP(ip string) {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    // Remove the IP from the blocked list
    delete(fw.blockedIPs, ip)
    fmt.Printf("Dynamic Firewall: IP %s unblocked.\n", ip)

    // Log the event to the ledger
    encryptedEvent, err := fw.encryptionService.EncryptData("AES", []byte(fmt.Sprintf("IP Unblocked: %s", ip)), common.EncryptionKey)
    if err != nil {
        fmt.Println("Error encrypting event:", err)
        return
    }

    // Log the unblocking event to the ledger
    fw.ledgerInstance.LogFirewallEvent(string(encryptedEvent))
}

// AdjustFirewallRules dynamically adjusts the rules based on traffic analysis
func (fw *DynamicFirewall) AdjustFirewallRules() {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    // Simulated traffic analysis logic (can be expanded to real-world logic with machine learning)
    for ip := range fw.blockedIPs {
        if time.Since(fw.blockedIPs[ip]) > fw.blockDuration {
            fw.UnblockIP(ip) // Unblock IPs after the block duration has passed
        }
    }

    fmt.Println("Dynamic Firewall rules adjusted based on traffic patterns.")
}


// BlockPort blocks a specific port in the stateless firewall
func (fw *StatelessFirewall) BlockPort(port int) {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    fw.allowedPorts = removePort(fw.allowedPorts, port)
    fmt.Printf("Stateless Firewall: Port %d blocked.\n", port)

    // Encrypt the event message
    encryptedEvent, err := fw.encryptionService.EncryptData("AES", []byte(fmt.Sprintf("Port Blocked: %d", port)), common.EncryptionKey)
    if err != nil {
        fmt.Println("Error encrypting event:", err)
        return
    }

    // Log the event to the ledger (assuming only one argument for LogFirewallEvent)
    fw.ledgerInstance.LogFirewallEvent(string(encryptedEvent))  // Correct the number of arguments
}


// IsAllowedPort checks if a port is allowed in the stateless firewall
func (fw *StatelessFirewall) IsAllowedPort(port int) bool {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	for _, allowedPort := range fw.allowedPorts {
		if allowedPort == port {
			return true
		}
	}

	fmt.Printf("Stateless Firewall: Port %d is blocked.\n", port)
	return false
}

// Utility function to remove a port from the allowedPorts slice
func removePort(ports []int, port int) []int {
	for i, p := range ports {
		if p == port {
			return append(ports[:i], ports[i+1:]...)
		}
	}
	return ports
}

// AllowConnection allows a connection in the stateful firewall
func (fw *StatefulFirewall) AllowConnection(ip, state string) {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    fw.allowedConnections[ip] = state
    fmt.Printf("Stateful Firewall: Connection from IP %s set to state '%s'.\n", ip, state)

    // Encrypt the event message
    encryptedEvent, err := fw.encryptionService.EncryptData("AES", []byte(fmt.Sprintf("Connection Allowed: %s, State: %s", ip, state)), common.EncryptionKey)
    if err != nil {
        fmt.Println("Failed to encrypt event:", err)
        return
    }

    // Log the event to the ledger (only pass the encrypted event)
    fw.ledgerInstance.LogFirewallEvent(string(encryptedEvent))  // Correcting the number of arguments
}

// BlockConnection blocks a connection in the stateful firewall
func (fw *StatefulFirewall) BlockConnection(ip string) {
    fw.mutex.Lock()
    defer fw.mutex.Unlock()

    delete(fw.allowedConnections, ip)
    fmt.Printf("Stateful Firewall: Connection from IP %s blocked.\n", ip)

    // Encrypt the event message
    encryptedEvent, err := fw.encryptionService.EncryptData("AES", []byte(fmt.Sprintf("Connection Blocked: %s", ip)), common.EncryptionKey)
    if err != nil {
        fmt.Println("Failed to encrypt event:", err)
        return
    }

    // Log the event to the ledger (only pass the encrypted event)
    fw.ledgerInstance.LogFirewallEvent(string(encryptedEvent))  // Correcting the number of arguments
}


