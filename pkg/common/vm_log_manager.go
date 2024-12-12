package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// LogEntry represents a structured log entry.
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// LogManager manages logging within the VM.
type LogManager struct {
	level          LogLevel
	output         io.Writer
	mutex          sync.Mutex
	logFile        *os.File
	logFilePath    string
	maxFileSize    int64  // Max size in bytes before rotation
	currentSize    int64  // Current size of the log file
	rotate         bool
	consoleLogging bool
}

// NewLogManager initializes a new LogManager.
func NewLogManager(level LogLevel, logToFile bool, logFilePath string, maxFileSize int64, consoleLogging bool) (*LogManager, error) {
	lm := &LogManager{
		level:          level,
		rotate:         logToFile,
		logFilePath:    logFilePath,
		maxFileSize:    maxFileSize,
		consoleLogging: consoleLogging,
	}

	// Set up log output
	if logToFile {
		err := lm.initializeLogFile()
		if err != nil {
			return nil, err
		}
	} else {
		lm.output = os.Stdout
	}

	return lm, nil
}

// initializeLogFile opens or creates the log file and sets up the output writer.
func (lm *LogManager) initializeLogFile() error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	// Ensure the directory exists
	dir := filepath.Dir(lm.logFilePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open or create the log file
	file, err := os.OpenFile(lm.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	lm.logFile = file
	lm.output = file

	// Get current file size
	fi, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get log file info: %v", err)
	}
	lm.currentSize = fi.Size()

	return nil
}

// rotateLogFile rotates the log file when it reaches the maximum size.
func (lm *LogManager) rotateLogFile() error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if lm.logFile != nil {
		lm.logFile.Close()
	}

	// Rename the current log file
	timestamp := time.Now().Format("20060102_150405")
	rotatedFileName := fmt.Sprintf("%s.%s", lm.logFilePath, timestamp)
	err := os.Rename(lm.logFilePath, rotatedFileName)
	if err != nil {
		return fmt.Errorf("failed to rotate log file: %v", err)
	}

	// Reinitialize the log file
	err = lm.initializeLogFile()
	if err != nil {
		return err
	}

	lm.currentSize = 0
	return nil
}

// log writes a log entry if the level is appropriate.
func (lm *LogManager) log(level LogLevel, message string, fields map[string]interface{}) {
	if level < lm.level {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     levelToString(level),
		Message:   message,
		Fields:    fields,
	}

	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	// Encode entry as JSON
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}

	data = append(data, '\n')

	// Write to output
	n, err := lm.output.Write(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log entry: %v\n", err)
		return
	}

	lm.currentSize += int64(n)

	// Also log to console if enabled
	if lm.consoleLogging && lm.output != os.Stdout {
		os.Stdout.Write(data)
	}

	// Rotate log file if needed
	if lm.rotate && lm.currentSize >= lm.maxFileSize {
		err = lm.rotateLogFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to rotate log file: %v\n", err)
		}
	}
}

func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Debug logs a debug-level message.
func (lm *LogManager) Debug(message string, fields map[string]interface{}) {
	lm.log(DEBUG, message, fields)
}

// Info logs an info-level message.
func (lm *LogManager) Info(message string, fields map[string]interface{}) {
	lm.log(INFO, message, fields)
}

// Warn logs a warn-level message.
func (lm *LogManager) Warn(message string, fields map[string]interface{}) {
	lm.log(WARN, message, fields)
}

// Error logs an error-level message.
func (lm *LogManager) Error(message string, fields map[string]interface{}) {
	lm.log(ERROR, message, fields)
}

// Close closes the LogManager and releases resources.
func (lm *LogManager) Close() {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	if lm.logFile != nil {
		lm.logFile.Close()
	}
}
