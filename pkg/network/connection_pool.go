package network

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ConnectionPool manages reusable connections between nodes
type ConnectionPool struct {
    connections map[string]net.Conn  // Active network connections keyed by node ID
    mutex       sync.Mutex           // Mutex for thread-safe access
    maxIdleTime time.Duration        // Maximum idle time before closing a connection
	ActiveConns []net.Conn // Use a list of net.Conn to store active connections
}

// NewConnectionPool initializes a new connection pool with a specified maximum idle time
func NewConnectionPool(maxIdleTime time.Duration) *ConnectionPool {
    return &ConnectionPool{
        connections: make(map[string]net.Conn),
        maxIdleTime: maxIdleTime,
    }
}

// GetConnection retrieves a connection from the pool or establishes a new one if none exists
func (cp *ConnectionPool) GetConnection(nodeID, address string) (net.Conn, error) {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    // Check if a connection already exists
    if conn, exists := cp.connections[nodeID]; exists {
        fmt.Printf("Reusing existing connection to node %s at %s\n", nodeID, address)
        return conn, nil
    }

    // Establish a new connection
    fmt.Printf("Establishing new connection to node %s at %s\n", nodeID, address)
    conn, err := net.DialTimeout("tcp", address, 5*time.Second)
    if err != nil {
        return nil, fmt.Errorf("failed to establish connection to %s: %v", address, err)
    }

    // Store the new connection in the pool
    cp.connections[nodeID] = conn
    return conn, nil
}

// ReleaseConnection releases a connection back to the pool or closes it if the idle time exceeds the limit
func (cp *ConnectionPool) ReleaseConnection(nodeID string, conn net.Conn) {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    fmt.Printf("Releasing connection to node %s\n", nodeID)

    // Close the connection if it has been idle for too long
    lastActivity := time.Now() // This would typically track the last activity timestamp
    if time.Since(lastActivity) > cp.maxIdleTime {
        fmt.Printf("Closing idle connection to node %s\n", nodeID)
        conn.Close()
        delete(cp.connections, nodeID)
    } else {
        // Otherwise, keep it in the pool for reuse
        cp.connections[nodeID] = conn
    }
}

// RemoveConnection forcibly closes and removes a connection from the pool
func (cp *ConnectionPool) RemoveConnection(nodeID string) {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    if conn, exists := cp.connections[nodeID]; exists {
        fmt.Printf("Forcibly closing connection to node %s\n", nodeID)
        conn.Close()
        delete(cp.connections, nodeID)
    } else {
        fmt.Printf("Connection to node %s not found in pool.\n", nodeID)
    }
}

// CloseAllConnections closes all connections in the pool and clears the pool
func (cp *ConnectionPool) CloseAllConnections() {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    fmt.Println("Closing all active connections in the pool...")
    for nodeID, conn := range cp.connections {
        fmt.Printf("Closing connection to node %s\n", nodeID)
        conn.Close()
        delete(cp.connections, nodeID)
    }
}

// GetActiveConnections returns the number of currently active connections in the pool
func (cp *ConnectionPool) GetActiveConnections() int {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    return len(cp.connections)
}

// MaintainConnectionPool periodically checks and closes idle connections
func (cp *ConnectionPool) MaintainConnectionPool() {
    go func() {
        for {
            time.Sleep(cp.maxIdleTime / 2) // Check connections periodically
            cp.mutex.Lock()
            for nodeID, conn := range cp.connections {
                lastActivity := time.Now() // This would track the actual last activity
                if time.Since(lastActivity) > cp.maxIdleTime {
                    fmt.Printf("Closing idle connection to node %s\n", nodeID)
                    conn.Close()
                    delete(cp.connections, nodeID)
                }
            }
            cp.mutex.Unlock()
        }
    }()
}

// AddConnection adds a new connection to the pool.
func (pool *ConnectionPool) AddConnection(conn net.Conn) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	// Add the connection to the ActiveConns slice
	pool.ActiveConns = append(pool.ActiveConns, conn)
	fmt.Printf("Connection added to pool from %s\n", conn.RemoteAddr().String())
}

// Close closes all active connections in the ConnectionPool
func (cp *ConnectionPool) Close() error {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	for _, conn := range cp.ActiveConns {
		err := conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close connection: %v", err)
		}
	}
	cp.ActiveConns = nil
	return nil
}

// SendPacket sends a packet to the specified route through the active connection
func (cp *ConnectionPool) SendPacket(route string, packet *Packet) error {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    // Check if there are active connections available
    if len(cp.ActiveConns) == 0 {
        return fmt.Errorf("no active connections available")
    }

    // Use the first active connection for simplicity, or implement load balancing if needed
    conn := cp.ActiveConns[0]

    // Convert the packet to bytes (you can use encoding like JSON, binary, etc.)
    packetData := packet.Data

    // Send the packet over the active connection
    _, err := conn.Write(packetData)
    if err != nil {
        return fmt.Errorf("failed to send packet: %v", err)
    }

    fmt.Printf("Packet sent to route %s\n", route)
    return nil
}


// Send sends data over all active connections in the pool
func (cp *ConnectionPool) Send(data []byte) error {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    for _, conn := range cp.ActiveConns {
        _, err := conn.Write(data)
        if err != nil {
            return fmt.Errorf("failed to send data: %v", err)
        }
    }
    return nil
}

// Receive receives data from any active connection in the pool
func (cp *ConnectionPool) Receive() ([]byte, error) {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    buf := make([]byte, 1024)
    for _, conn := range cp.ActiveConns {
        n, err := conn.Read(buf)
        if err != nil {
            return nil, fmt.Errorf("failed to receive data: %v", err)
        }
        return buf[:n], nil
    }
    return nil, fmt.Errorf("no active connections to receive data from")
}

// Add IsAlive method to ConnectionPool
func (cp *ConnectionPool) IsAlive() bool {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    // Check if there's at least one active connection
    for _, conn := range cp.ActiveConns {
        if conn != nil {
            return true
        }
    }
    return false
}
