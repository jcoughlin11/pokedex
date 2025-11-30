package main

import "strings"

func cleanInput(text string) []string {
	cleaned := make([]string, 0)

	trimmed := strings.TrimSpace(text)
	parts := strings.Split(trimmed, " ")

	for _, part := range parts {
		if part == "" {
			continue
		}
		cleaned = append(cleaned, strings.ToLower(part))
	}

	return cleaned
}
