package mcp

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func generateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func mapToStruct(input interface{}, output interface{}) error {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, output); err != nil {
		return fmt.Errorf("failed to unmarshal to output: %w", err)
	}

	return nil
}
