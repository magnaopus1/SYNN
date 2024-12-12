package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "synnergy_network/api"
    "synnergy_network/cli"
    "synnergy_network/config"
    "github.com/spf13/viper"
)

// LoadConfig loads the configuration from a YAML file
func LoadConfig(path string) (*config.Config, error) {
    viper.SetConfigFile(path)
    err := viper.ReadInConfig()
    if err != nil {
        return nil, err
    }

    var cfg config.Config
    err = viper.Unmarshal(&cfg)
    if err != nil {
        return nil, err
    }

    return &cfg, nil
}

func main() {
    // Load configuration
    cfg, err := LoadConfig("cmd/configs/mainnet_config.yaml")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Shared state
    var stateMutex sync.RWMutex

    // Initialize API Server
    apiServer := api.NewAPIServer(cfg, &stateMutex)

    // Start API Server in a separate goroutine
    go func() {
        apiServer.Start()
    }()

    // Wait for API server to start before running the script
    waitForServer(fmt.Sprintf("http://localhost:%d", cfg.API.Port))

    // Run the mainnet startup script
    err = runMainnetSetupScript()
    if err != nil {
        log.Fatalf("Failed to run mainnet setup script: %v", err)
    }

    // Start CLI Server in a separate goroutine
    go func() {
        cliServer := cli.NewCLIServer(cfg)
        cliServer.Start()
    }()

    // Handle graceful shutdown
    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
    <-shutdown
    fmt.Println("\nShutting down Mainnet...")

    fmt.Println("Mainnet successfully shut down.")
}

// runMainnetSetupScript executes the mainnet startup script
func runMainnetSetupScript() error {
    cmd := exec.Command("/bin/bash", "cmd/scripts/mainnet_start_up_scripts/mainnet_start_up.sh")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("failed to run mainnet startup script: %w", err)
    }
    return nil
}

// waitForServer waits until the API server is responding
func waitForServer(url string) {
    for {
        resp, err := http.Get(url)
        if err == nil && resp.StatusCode == http.StatusOK {
            fmt.Println("API server is up and running.")
            resp.Body.Close()
            return
        }
        if resp != nil {
            resp.Body.Close()
        }
        fmt.Println("Waiting for API server to start...")
        time.Sleep(1 * time.Second)
    }
}
