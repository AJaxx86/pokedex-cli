package main

import "strings"

func cleanInput(text string) []string {
	result := strings.Fields(text)
	for i, w := range result {
		result[i] = strings.ToLower(w)
	}
	return result
}
