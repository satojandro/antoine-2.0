// Package config provides secure credential management for Antoine CLI
// This file handles API keys, tokens, and other sensitive configuration data
package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Credential types for different services
const (
	CredentialTypeAPI       = "api"
	CredentialTypeMCP       = "mcp"
	CredentialTypeGitHub    = "github"
	CredentialTypeOpenAI    = "openai"
	CredentialTypeAnthropic = "anthropic"
	CredentialTypeAuth      = "auth"
)

// CredentialManager handles secure storage and retrieval of credentials
type CredentialManager struct {
	serviceName string
	storage     CredentialStorage
	encryption  bool
	keyring     bool
}

// CredentialStorage defines the interface for credential storage backends
type CredentialStorage interface {
	Store(key string, credential Credential) error
	Retrieve(key string) (Credential, error)
	Delete(key string) error
	List() ([]string, error)
	Clear() error
}

// Credential represents a stored credential with metadata
type Credential struct {
	Type        string            `json:"type"`
	Value       string            `json:"value"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
	Encrypted   bool              `json:"encrypted"`
	Description string            `json:"description,omitempty"`
}

// NewCredentialManager creates a new credential manager
func NewCredentialManager(serviceName string, useKeyring bool, useEncryption bool) *CredentialManager {
	var storage CredentialStorage

	if useKeyring {
		storage = &KeyringStorage{serviceName: serviceName}
	} else {
		storage = &FileStorage{
			serviceName: serviceName,
			encryption:  useEncryption,
		}
	}

	return &CredentialManager{
		serviceName: serviceName,
		storage:     storage,
		encryption:  useEncryption,
		keyring:     useKeyring,
	}
}

// Store saves a credential securely
func (cm *CredentialManager) Store(credType, key, value string, metadata map[string]string) error {
	if key == "" || value == "" {
		return errors.New("credential key and value cannot be empty")
	}

	credential := Credential{
		Type:        credType,
		Value:       value,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Encrypted:   cm.encryption,
		Description: fmt.Sprintf("%s credential for %s", credType, key),
	}

	// Set expiration for tokens
	if credType == CredentialTypeAuth || strings.Contains(credType, "token") {
		if expiryStr, exists := metadata["expires_in"]; exists {
			if duration, err := time.ParseDuration(expiryStr); err == nil {
				expiry := time.Now().Add(duration)
				credential.ExpiresAt = &expiry
			}
		}
	}

	return cm.storage.Store(key, credential)
}

// Retrieve gets a credential by key
func (cm *CredentialManager) Retrieve(key string) (Credential, error) {
	if key == "" {
		return Credential{}, errors.New("credential key cannot be empty")
	}

	credential, err := cm.storage.Retrieve(key)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to retrieve credential '%s': %w", key, err)
	}

	// Check if credential has expired
	if credential.ExpiresAt != nil && time.Now().After(*credential.ExpiresAt) {
		return Credential{}, fmt.Errorf("credential '%s' has expired", key)
	}

	return credential, nil
}

// RetrieveValue gets just the credential value by key
func (cm *CredentialManager) RetrieveValue(key string) (string, error) {
	credential, err := cm.Retrieve(key)
	if err != nil {
		return "", err
	}
	return credential.Value, nil
}

// Update modifies an existing credential
func (cm *CredentialManager) Update(key, value string, metadata map[string]string) error {
	// Retrieve existing credential to preserve metadata
	existing, err := cm.storage.Retrieve(key)
	if err != nil {
		return fmt.Errorf("credential '%s' not found for update", key)
	}

	// Update the credential
	existing.Value = value
	existing.UpdatedAt = time.Now()

	// Merge metadata
	if existing.Metadata == nil {
		existing.Metadata = make(map[string]string)
	}
	for k, v := range metadata {
		existing.Metadata[k] = v
	}

	return cm.storage.Store(key, existing)
}

// Delete removes a credential
func (cm *CredentialManager) Delete(key string) error {
	return cm.storage.Delete(key)
}

// List returns all stored credential keys
func (cm *CredentialManager) List() ([]string, error) {
	return cm.storage.List()
}

// Clear removes all credentials
func (cm *CredentialManager) Clear() error {
	return cm.storage.Clear()
}

// IsExpired checks if a credential has expired
func (cm *CredentialManager) IsExpired(key string) (bool, error) {
	credential, err := cm.storage.Retrieve(key)
	if err != nil {
		return false, err
	}

	if credential.ExpiresAt == nil {
		return false, nil
	}

	return time.Now().After(*credential.ExpiresAt), nil
}

// Refresh updates a credential's expiration time
func (cm *CredentialManager) Refresh(key string, duration time.Duration) error {
	credential, err := cm.storage.Retrieve(key)
	if err != nil {
		return err
	}

	newExpiry := time.Now().Add(duration)
	credential.ExpiresAt = &newExpiry
	credential.UpdatedAt = time.Now()

	return cm.storage.Store(key, credential)
}

// KeyringStorage implements credential storage using the system keyring
type KeyringStorage struct {
	serviceName string
}

// Store saves a credential to the system keyring
func (ks *KeyringStorage) Store(key string, credential Credential) error {
	data, err := json.Marshal(credential)
	if err != nil {
		return fmt.Errorf("failed to marshal credential: %w", err)
	}

	err = keyring.Set(ks.serviceName, key, string(data))
	if err != nil {
		return fmt.Errorf("failed to store credential in keyring: %w", err)
	}

	return nil
}

// Retrieve gets a credential from the system keyring
func (ks *KeyringStorage) Retrieve(key string) (Credential, error) {
	data, err := keyring.Get(ks.serviceName, key)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to get credential from keyring: %w", err)
	}

	var credential Credential
	err = json.Unmarshal([]byte(data), &credential)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to unmarshal credential: %w", err)
	}

	return credential, nil
}

// Delete removes a credential from the system keyring
func (ks *KeyringStorage) Delete(key string) error {
	err := keyring.Delete(ks.serviceName, key)
	if err != nil {
		return fmt.Errorf("failed to delete credential from keyring: %w", err)
	}
	return nil
}

// List returns all credential keys (keyring doesn't support listing)
func (ks *KeyringStorage) List() ([]string, error) {
	return nil, errors.New("keyring storage does not support listing credentials")
}

// Clear removes all credentials (keyring doesn't support bulk operations)
func (ks *KeyringStorage) Clear() error {
	return errors.New("keyring storage does not support clearing all credentials")
}

// FileStorage implements credential storage using encrypted files
type FileStorage struct {
	serviceName string
	encryption  bool
	basePath    string
}

// NewFileStorage creates a new file-based credential storage
func NewFileStorage(serviceName string, encryption bool) *FileStorage {
	homeDir, _ := os.UserHomeDir()
	basePath := filepath.Join(homeDir, ".antoine", "credentials")

	return &FileStorage{
		serviceName: serviceName,
		encryption:  encryption,
		basePath:    basePath,
	}
}

// Store saves a credential to an encrypted file
func (fs *FileStorage) Store(key string, credential Credential) error {
	// Ensure directory exists
	err := os.MkdirAll(fs.basePath, 0700)
	if err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}

	data, err := json.Marshal(credential)
	if err != nil {
		return fmt.Errorf("failed to marshal credential: %w", err)
	}

	// Encrypt if enabled
	if fs.encryption {
		data, err = fs.encrypt(data, key)
		if err != nil {
			return fmt.Errorf("failed to encrypt credential: %w", err)
		}
	}

	filePath := filepath.Join(fs.basePath, fs.getFileName(key))
	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write credential file: %w", err)
	}

	return nil
}

// Retrieve gets a credential from an encrypted file
func (fs *FileStorage) Retrieve(key string) (Credential, error) {
	filePath := filepath.Join(fs.basePath, fs.getFileName(key))

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Credential{}, fmt.Errorf("credential not found")
		}
		return Credential{}, fmt.Errorf("failed to read credential file: %w", err)
	}

	// Decrypt if enabled
	if fs.encryption {
		data, err = fs.decrypt(data, key)
		if err != nil {
			return Credential{}, fmt.Errorf("failed to decrypt credential: %w", err)
		}
	}

	var credential Credential
	err = json.Unmarshal(data, &credential)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to unmarshal credential: %w", err)
	}

	return credential, nil
}

// Delete removes a credential file
func (fs *FileStorage) Delete(key string) error {
	filePath := filepath.Join(fs.basePath, fs.getFileName(key))

	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete credential file: %w", err)
	}

	return nil
}

// List returns all credential keys from files
func (fs *FileStorage) List() ([]string, error) {
	entries, err := os.ReadDir(fs.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read credentials directory: %w", err)
	}

	var keys []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".cred") {
			key := strings.TrimSuffix(entry.Name(), ".cred")
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// Clear removes all credential files
func (fs *FileStorage) Clear() error {
	entries, err := os.ReadDir(fs.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read credentials directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".cred") {
			filePath := filepath.Join(fs.basePath, entry.Name())
			err := os.Remove(filePath)
			if err != nil {
				return fmt.Errorf("failed to remove credential file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// getFileName creates a safe filename for a credential key
func (fs *FileStorage) getFileName(key string) string {
	// Replace unsafe characters and add extension
	safe := strings.ReplaceAll(key, "/", "_")
	safe = strings.ReplaceAll(safe, "\\", "_")
	safe = strings.ReplaceAll(safe, ":", "_")
	return safe + ".cred"
}

// encrypt encrypts data using AES-GCM
func (fs *FileStorage) encrypt(data []byte, key string) ([]byte, error) {
	// Derive key from credential key and service name
	hash := sha256.Sum256([]byte(fs.serviceName + ":" + key))

	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// decrypt decrypts data using AES-GCM
func (fs *FileStorage) decrypt(data []byte, key string) ([]byte, error) {
	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	// Derive key from credential key and service name
	hash := sha256.Sum256([]byte(fs.serviceName + ":" + key))

	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// Helper functions for common credential operations

// StoreAPIKey stores an API key for a service
func StoreAPIKey(service, apiKey string) error {
	cm := NewCredentialManager("antoine-cli", true, true)
	return cm.Store(CredentialTypeAPI, service+".api_key", apiKey, map[string]string{
		"service": service,
		"type":    "api_key",
	})
}

// GetAPIKey retrieves an API key for a service
func GetAPIKey(service string) (string, error) {
	cm := NewCredentialManager("antoine-cli", true, true)
	return cm.RetrieveValue(service + ".api_key")
}

// StoreToken stores an authentication token
func StoreToken(service, token string, expiresIn time.Duration) error {
	cm := NewCredentialManager("antoine-cli", true, true)
	metadata := map[string]string{
		"service": service,
		"type":    "token",
	}

	if expiresIn > 0 {
		metadata["expires_in"] = expiresIn.String()
	}

	return cm.Store(CredentialTypeAuth, service+".token", token, metadata)
}

// GetToken retrieves an authentication token
func GetToken(service string) (string, error) {
	cm := NewCredentialManager("antoine-cli", true, true)
	return cm.RetrieveValue(service + ".token")
}

// StoreMCPCredentials stores credentials for an MCP server
func StoreMCPCredentials(serverName, endpoint, apiKey string) error {
	cm := NewCredentialManager("antoine-cli", true, true)
	return cm.Store(CredentialTypeMCP, serverName+".credentials", apiKey, map[string]string{
		"server":   serverName,
		"endpoint": endpoint,
		"type":     "mcp_credentials",
	})
}

// GetMCPCredentials retrieves credentials for an MCP server
func GetMCPCredentials(serverName string) (string, error) {
	cm := NewCredentialManager("antoine-cli", true, true)
	return cm.RetrieveValue(serverName + ".credentials")
}

// ValidateCredentials checks if credentials are valid and not expired
func ValidateCredentials() error {
	cm := NewCredentialManager("antoine-cli", true, true)

	keys, err := cm.List()
	if err != nil {
		return fmt.Errorf("failed to list credentials: %w", err)
	}

	expiredKeys := []string{}
	for _, key := range keys {
		expired, err := cm.IsExpired(key)
		if err != nil {
			continue // Skip credentials we can't check
		}

		if expired {
			expiredKeys = append(expiredKeys, key)
		}
	}

	if len(expiredKeys) > 0 {
		return fmt.Errorf("found %d expired credentials: %v", len(expiredKeys), expiredKeys)
	}

	return nil
}

// CleanupExpiredCredentials removes all expired credentials
func CleanupExpiredCredentials() error {
	cm := NewCredentialManager("antoine-cli", true, true)

	keys, err := cm.List()
	if err != nil {
		return fmt.Errorf("failed to list credentials: %w", err)
	}

	cleanedCount := 0
	for _, key := range keys {
		expired, err := cm.IsExpired(key)
		if err != nil {
			continue // Skip credentials we can't check
		}

		if expired {
			err := cm.Delete(key)
			if err == nil {
				cleanedCount++
			}
		}
	}

	if cleanedCount > 0 {
		fmt.Printf("Cleaned up %d expired credentials\n", cleanedCount)
	}

	return nil
}

// GetCredentialStatus returns status information about stored credentials
func GetCredentialStatus() (map[string]interface{}, error) {
	cm := NewCredentialManager("antoine-cli", true, true)

	keys, err := cm.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list credentials: %w", err)
	}

	status := map[string]interface{}{
		"total_credentials": len(keys),
		"by_type":           make(map[string]int),
		"expired":           0,
		"valid":             0,
	}

	typeCount := make(map[string]int)
	expiredCount := 0
	validCount := 0

	for _, key := range keys {
		credential, err := cm.Retrieve(key)
		if err != nil {
			continue
		}

		typeCount[credential.Type]++

		expired, err := cm.IsExpired(key)
		if err != nil {
			continue
		}

		if expired {
			expiredCount++
		} else {
			validCount++
		}
	}

	status["by_type"] = typeCount
	status["expired"] = expiredCount
	status["valid"] = validCount

	return status, nil
}
