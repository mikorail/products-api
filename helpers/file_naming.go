package helpers

import (
	"fmt"
	"strings"
	"time"
)

func GenerateFileName(filter map[string]interface{}) string {
	// Create a string with all filters
	fmt.Println(filter)
	var filterParts []string
	for key, value := range filter {
		fmt.Println("key:", key, "value:", value)
		filterParts = append(filterParts, fmt.Sprintf("%s-%v", key, value))
	}

	// Join the filters with underscores and append the timestamp
	filterStr := strings.Join(filterParts, "_")
	timestamp := time.Now().Format("20060102_150405") // Format as YYYYMMDD_HHMMSS

	// Clean up the filename (remove spaces and special characters)
	filename := strings.ReplaceAll(filterStr, " ", "_")
	filename = strings.ReplaceAll(filename, ":", "-")

	// Combine filter and timestamp for the filename
	return fmt.Sprintf("%s_%s", filename, timestamp)
}
