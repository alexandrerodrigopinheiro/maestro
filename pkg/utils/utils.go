package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Atoi converts a string to an integer.
// If the conversion fails, it logs a fatal error.
func Atoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid port number: %s", s)
	}
	return value
}

// AtoiWithDefault converts a string to an integer, and if the conversion fails,
// it returns a provided default value instead of logging a fatal error.
func AtoiWithDefault(s string, defaultValue int) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Invalid port number '%s', using default: %d\n", s, defaultValue)
		return defaultValue
	}
	return value
}

// toSnakeCase converts a string to snake_case and makes it lowercase.
func ToSnakeCase(str string) string {
	// Regular expression to identify uppercase letters
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(str, "${1}_${2}")
	// Convert to lowercase and return
	return strings.ToLower(snake)
}

// ToPascalCase converts a snake_case string to PascalCase.
func ToPascalCase(str string) string {
	// Split the string by underscores
	words := strings.Split(str, "_")
	titleCaser := cases.Title(language.Und)
	for i, word := range words {
		// Capitalize the first letter of each word and append to result
		words[i] = titleCaser.String(word)
	}
	return strings.Join(words, "")
}

// LoadEnv loads environment variables from a .env file.
func LoadEnv(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// Split key and value by the first '=' character
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Set the environment variable
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}
