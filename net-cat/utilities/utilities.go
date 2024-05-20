package utilities

import (
	"fmt"
	"log"
	"os"
	"time"
)

// GetTime returns the current time in a formatted string.
func GetTime() string {
	return time.Now().Format("[2006-01-02 15:04:05]")
}

// ValidatePort checks if a given port is valid.
func ValidatePort(port string) bool {
	// You can implement port validation logic here.
	// For simplicity, let's assume any non-empty port is valid.
	return port != ""
}

// LogError logs an error message along with the error to the standard logger.
func LogError(message string, err error) {
	log.Printf("[%s] ERROR: %s - %v\n", GetTime(), message, err)
}

// SaveToFile saves a slice of strings to a file.
func SaveToFile(logs []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		LogError("Failed to create file", err)
		return err
	}
	defer file.Close()

	for _, logEntry := range logs {
		_, err := file.WriteString(logEntry + "\n")
		if err != nil {
			LogError("Failed to write to file", err)
			return err
		}
	}

	fmt.Println("Chat log saved to", filename)
	return nil
}
