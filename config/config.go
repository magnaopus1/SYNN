package config

// Config holds the mainnet configuration
type Config struct {
    API struct {
        Port int
    }
    CLI struct {
        Port int
    }
    Node struct {  // Add the Node field back
        IP string
    }
}