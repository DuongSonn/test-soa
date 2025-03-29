package utils

import (
	"regexp"
	"strings"
)

// ConvertToUpperCase convert string to uppercase (e.g: "hello world" -> "HELLO_WORLD")
func ConvertToUpperCase(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "_")
	return strings.ToUpper(str)
}

// ConvertToCamelCase convert string to camel case (e.g: "hello world" -> "helloWorld")
func ConvertToCamelCase(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, " ")
	parts := strings.Split(str, " ")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}

	// Concatenate the words
	result := strings.Join(parts, "")
	return result
}

// ConvertToSlug convert string to slug (e.g: "hello world" -> "hello-world")
func ConvertToSlug(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "-")
	return strings.ToLower(str)
}
