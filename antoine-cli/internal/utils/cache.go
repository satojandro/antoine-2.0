// Package utils provides caching utilities for Antoine CLI
// This file implements a flexible caching system with memory, disk, and hybrid storage
package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

// CacheType defines the type of cache storage
type CacheType string

const (
	CacheTypeMemory CacheType = "memory"
	CacheTypeDisk   CacheType = "disk"
	CacheTypeHybrid CacheType = "hybrid"
)

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled         bool                     `yaml:"enabled"`
	Type            CacheType                `yaml:"type"`
	MaxSizeMB       int64                    `yaml:"max_size_mb"`
	MaxEntries      int64                    `yaml:"max_entries"`
	CleanupInterval time.Duration            `yaml:"cleanup_interval"`
	TTL             map[string]time.Duration `yaml:"ttl"`
	Disk            DiskCacheConfig          `yaml:"disk"`
}

// DiskCacheConfig holds disk cache configuration
type DiskCacheConfig struct {
	Path        string `yaml:"path"`
	Compression bool   `yaml:"compression"`
}

// CacheItem represents a cached item with metadata
type CacheItem struct {
	Key         string                 `json:"key"`
	Value       interface{}            `json:"value"`
	Type        string                 `json:"type"`
	CreatedAt   time.Time              `json:"created_at"`
	ExpiresAt   time.Time              `json:"expires_at"`
	AccessCount int64                  `json:"access_count"`
	LastAccess  time.Time              `json:"last_access"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Cache interface defines the cache operations
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration) error
	SetWithType(key string, value interface{}, itemType string, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	Keys() []string
	Stats() CacheStatistics
	Close() error
}

// CacheStatistics provides cache statistics (renamed to avoid conflict)
type CacheStatistics struct {
	TotalEntries   int64   `json:"total_entries"`
	MemoryUsageMB  float64 `json:"memory_usage_mb"`
	DiskUsageMB    float64 `json:"disk_usage_mb"`
	HitRatio       float64 `json:"hit_ratio"`
	MissRatio      float64 `json:"miss_ratio"`
	TotalHits      int64   `json:"total_hits"`
	TotalMisses    int64   `json:"total_misses"`
	ExpiredEntries int64   `json:"expired_entries"`
	EvictedEntries int64   `json:"evicted_entries"`
}

// CacheManager manages different cache implementations
type CacheManager struct {
	config      CacheConfig
	cache       Cache
	stopCleanup chan bool
	mu          sync.RWMutex
}

// NewCacheManager creates a new cache manager
func NewCacheManager(config CacheConfig) (*CacheManager, error) {
	if !config.Enabled {
		return &CacheManager{
			config: config,
			cache:  &NoOpCache{},
		}, nil
	}

	var cache Cache
	var err error

	switch config.Type {
	case CacheTypeMemory:
		cache, err = NewMemoryCache(config)
	case CacheTypeDisk:
		cache, err = NewDiskCache(config)
	case CacheTypeHybrid:
		cache, err = NewHybridCache(config)
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", config.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	manager := &CacheManager{
		config:      config,
		cache:       cache,
		stopCleanup: make(chan bool),
	}

	// Start cleanup routine
	if config.CleanupInterval > 0 {
		go manager.cleanupRoutine()
	}

	return manager, nil
}

// Get retrieves a value from cache
func (cm *CacheManager) Get(key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cache.Get(key)
}

// Set stores a value in cache with default TTL
func (cm *CacheManager) Set(key string, value interface{}) error {
	return cm.SetWithTTL(key, value, cm.getDefaultTTL(key))
}

// SetWithTTL stores a value in cache with specific TTL
func (cm *CacheManager) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.cache.Set(key, value, ttl)
}

// SetWithType stores a value in cache with type and TTL
func (cm *CacheManager) SetWithType(key string, value interface{}, itemType string, ttl time.Duration) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.cache.SetWithType(key, value, itemType, ttl)
}

// Delete removes a value from cache
func (cm *CacheManager) Delete(key string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.cache.Delete(key)
}

// Clear removes all values from cache
func (cm *CacheManager) Clear() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.cache.Clear()
}

// Keys returns all cache keys
func (cm *CacheManager) Keys() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cache.Keys()
}

// Stats returns cache statistics
func (cm *CacheManager) Stats() CacheStatistics {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cache.Stats()
}

// Close closes the cache manager
func (cm *CacheManager) Close() error {
	close(cm.stopCleanup)
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.cache.Close()
}

// getDefaultTTL returns the default TTL for a key based on its type
func (cm *CacheManager) getDefaultTTL(key string) time.Duration {
	for prefix, ttl := range cm.config.TTL {
		if strings.HasPrefix(key, prefix) {
			return ttl
		}
	}
	return 1 * time.Hour // Default TTL
}

// cleanupRoutine periodically cleans up expired entries
func (cm *CacheManager) cleanupRoutine() {
	ticker := time.NewTicker(cm.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.cleanupExpired()
		case <-cm.stopCleanup:
			return
		}
	}
}

// cleanupExpired removes expired entries from cache
func (cm *CacheManager) cleanupExpired() {
	// This method would be implemented by specific cache types
	// For now, we'll log the cleanup attempt
	WithComponent("cache").Debug("Running cache cleanup")
}

// MemoryCache implements in-memory caching using ristretto
type MemoryCache struct {
	cache   *ristretto.Cache
	config  CacheConfig
	metrics CacheMetrics
	mu      sync.RWMutex
}

// CacheMetrics tracks cache performance metrics
type CacheMetrics struct {
	Hits         int64
	Misses       int64
	Expired      int64
	Evicted      int64
	TotalEntries int64
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache(config CacheConfig) (*MemoryCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: config.MaxEntries * 10,
		MaxCost:     config.MaxSizeMB << 20, // Convert MB to bytes
		BufferItems: 64,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ristretto cache: %w", err)
	}

	return &MemoryCache{
		cache:  cache,
		config: config,
	}, nil
}

// Get retrieves a value from memory cache
func (mc *MemoryCache) Get(key string) (interface{}, bool) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	value, found := mc.cache.Get(key)
	if found {
		mc.metrics.Hits++

		// Check if item is expired
		if item, ok := value.(*CacheItem); ok {
			if time.Now().After(item.ExpiresAt) {
				mc.cache.Del(key)
				mc.metrics.Expired++
				mc.metrics.Misses++
				return nil, false
			}
			item.AccessCount++
			item.LastAccess = time.Now()
			return item.Value, true
		}
		return value, true
	}

	mc.metrics.Misses++
	return nil, false
}

// Set stores a value in memory cache
func (mc *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	return mc.SetWithType(key, value, "default", ttl)
}

// SetWithType stores a value with type in memory cache
func (mc *MemoryCache) SetWithType(key string, value interface{}, itemType string, ttl time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	item := &CacheItem{
		Key:         key,
		Value:       value,
		Type:        itemType,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(ttl),
		AccessCount: 0,
		LastAccess:  time.Now(),
	}

	// Estimate cost based on JSON size
	cost := estimateSize(item)

	success := mc.cache.Set(key, item, cost)
	if success {
		mc.metrics.TotalEntries++
	}

	return nil
}

// Delete removes a value from memory cache
func (mc *MemoryCache) Delete(key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.cache.Del(key)
	mc.metrics.TotalEntries--
	return nil
}

// Clear removes all values from memory cache
func (mc *MemoryCache) Clear() error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.cache.Clear()
	mc.metrics = CacheMetrics{} // Reset metrics
	return nil
}

// Keys returns all keys in memory cache
func (mc *MemoryCache) Keys() []string {
	// Ristretto doesn't support key listing, so we'll return empty slice
	// In a real implementation, you might maintain a separate key set
	return []string{}
}

// Stats returns memory cache statistics
func (mc *MemoryCache) Stats() CacheStatistics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metrics := mc.cache.Metrics
	total := mc.metrics.Hits + mc.metrics.Misses
	hitRatio := 0.0
	if total > 0 {
		hitRatio = float64(mc.metrics.Hits) / float64(total)
	}

	return CacheStatistics{
		TotalEntries:   mc.metrics.TotalEntries,
		MemoryUsageMB:  float64(metrics.CostAdded()) / (1024 * 1024), // Fixed: call CostAdded as function
		HitRatio:       hitRatio,
		MissRatio:      1.0 - hitRatio,
		TotalHits:      mc.metrics.Hits,
		TotalMisses:    mc.metrics.Misses,
		ExpiredEntries: mc.metrics.Expired,
		EvictedEntries: mc.metrics.Evicted,
	}
}

// Close closes the memory cache
func (mc *MemoryCache) Close() error {
	mc.cache.Close()
	return nil
}

// DiskCache implements disk-based caching
type DiskCache struct {
	config   CacheConfig
	basePath string
	metrics  CacheMetrics
	mu       sync.RWMutex
}

// NewDiskCache creates a new disk cache
func NewDiskCache(config CacheConfig) (*DiskCache, error) {
	basePath := config.Disk.Path
	if strings.HasPrefix(basePath, "~/") {
		home, _ := os.UserHomeDir()
		basePath = filepath.Join(home, basePath[2:])
	}

	// Ensure directory exists
	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &DiskCache{
		config:   config,
		basePath: basePath,
	}, nil
}

// Get retrieves a value from disk cache
func (dc *DiskCache) Get(key string) (interface{}, bool) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	filePath := dc.getFilePath(key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		dc.metrics.Misses++
		return nil, false
	}

	var item CacheItem
	err = json.Unmarshal(data, &item)
	if err != nil {
		dc.metrics.Misses++
		return nil, false
	}

	// Check expiration
	if time.Now().After(item.ExpiresAt) {
		os.RemoveAll(filePath)
		dc.metrics.Expired++
		dc.metrics.Misses++
		return nil, false
	}

	// Update access info
	item.AccessCount++
	item.LastAccess = time.Now()

	// Write back updated item
	updatedData, _ := json.Marshal(item)
	os.WriteFile(filePath, updatedData, 0600)

	dc.metrics.Hits++
	return item.Value, true
}

// Set stores a value in disk cache
func (dc *DiskCache) Set(key string, value interface{}, ttl time.Duration) error {
	return dc.SetWithType(key, value, "default", ttl)
}

// SetWithType stores a value with type in disk cache
func (dc *DiskCache) SetWithType(key string, value interface{}, itemType string, ttl time.Duration) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	item := &CacheItem{
		Key:         key,
		Value:       value,
		Type:        itemType,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(ttl),
		AccessCount: 0,
		LastAccess:  time.Now(),
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal cache item: %w", err)
	}

	filePath := dc.getFilePath(key)
	err = os.WriteFile(filePath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	dc.metrics.TotalEntries++
	return nil
}

// Delete removes a value from disk cache
func (dc *DiskCache) Delete(key string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	filePath := dc.getFilePath(key)
	err := os.RemoveAll(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete cache file: %w", err)
	}

	dc.metrics.TotalEntries--
	return nil
}

// Clear removes all values from disk cache
func (dc *DiskCache) Clear() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	err := os.RemoveAll(dc.basePath)
	if err != nil {
		return fmt.Errorf("failed to clear cache directory: %w", err)
	}

	// Recreate directory
	err = os.MkdirAll(dc.basePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to recreate cache directory: %w", err)
	}

	dc.metrics = CacheMetrics{} // Reset metrics
	return nil
}

// Keys returns all keys in disk cache
func (dc *DiskCache) Keys() []string {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	var keys []string
	filepath.WalkDir(dc.basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".cache") {
			key := strings.TrimSuffix(d.Name(), ".cache")
			keys = append(keys, key)
		}
		return nil
	})

	return keys
}

// Stats returns disk cache statistics
func (dc *DiskCache) Stats() CacheStatistics {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	diskUsage := dc.calculateDiskUsage()
	total := dc.metrics.Hits + dc.metrics.Misses
	hitRatio := 0.0
	if total > 0 {
		hitRatio = float64(dc.metrics.Hits) / float64(total)
	}

	return CacheStatistics{
		TotalEntries:   dc.metrics.TotalEntries,
		DiskUsageMB:    diskUsage,
		HitRatio:       hitRatio,
		MissRatio:      1.0 - hitRatio,
		TotalHits:      dc.metrics.Hits,
		TotalMisses:    dc.metrics.Misses,
		ExpiredEntries: dc.metrics.Expired,
	}
}

// Close closes the disk cache
func (dc *DiskCache) Close() error {
	return nil
}

// getFilePath generates a file path for a cache key
func (dc *DiskCache) getFilePath(key string) string {
	// Hash the key to create a safe filename
	hash := sha256.Sum256([]byte(key))
	filename := fmt.Sprintf("%x.cache", hash)
	return filepath.Join(dc.basePath, filename)
}

// calculateDiskUsage calculates total disk usage in MB
func (dc *DiskCache) calculateDiskUsage() float64 {
	var totalSize int64
	filepath.WalkDir(dc.basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				totalSize += info.Size()
			}
		}
		return nil
	})
	return float64(totalSize) / (1024 * 1024)
}

// HybridCache combines memory and disk caching
type HybridCache struct {
	memory *MemoryCache
	disk   *DiskCache
	config CacheConfig
}

// NewHybridCache creates a new hybrid cache
func NewHybridCache(config CacheConfig) (*HybridCache, error) {
	memory, err := NewMemoryCache(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory cache: %w", err)
	}

	disk, err := NewDiskCache(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create disk cache: %w", err)
	}

	return &HybridCache{
		memory: memory,
		disk:   disk,
		config: config,
	}, nil
}

// Get retrieves a value from hybrid cache (memory first, then disk)
func (hc *HybridCache) Get(key string) (interface{}, bool) {
	// Try memory first
	if value, found := hc.memory.Get(key); found {
		return value, true
	}

	// Try disk cache
	if value, found := hc.disk.Get(key); found {
		// Store in memory for faster future access
		hc.memory.Set(key, value, time.Hour) // Use default TTL
		return value, true
	}

	return nil, false
}

// Set stores a value in both memory and disk cache
func (hc *HybridCache) Set(key string, value interface{}, ttl time.Duration) error {
	return hc.SetWithType(key, value, "default", ttl)
}

// SetWithType stores a value with type in hybrid cache
func (hc *HybridCache) SetWithType(key string, value interface{}, itemType string, ttl time.Duration) error {
	// Store in both caches
	err := hc.memory.SetWithType(key, value, itemType, ttl)
	if err != nil {
		return err
	}

	return hc.disk.SetWithType(key, value, itemType, ttl)
}

// Delete removes a value from both caches
func (hc *HybridCache) Delete(key string) error {
	hc.memory.Delete(key)
	return hc.disk.Delete(key)
}

// Clear removes all values from both caches
func (hc *HybridCache) Clear() error {
	hc.memory.Clear()
	return hc.disk.Clear()
}

// Keys returns all keys from both caches
func (hc *HybridCache) Keys() []string {
	diskKeys := hc.disk.Keys()
	memoryKeys := hc.memory.Keys()

	// Combine and deduplicate keys
	keySet := make(map[string]bool)
	for _, key := range diskKeys {
		keySet[key] = true
	}
	for _, key := range memoryKeys {
		keySet[key] = true
	}

	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}

	return keys
}

// Stats returns combined statistics from both caches
func (hc *HybridCache) Stats() CacheStatistics {
	memStats := hc.memory.Stats()
	diskStats := hc.disk.Stats()

	return CacheStatistics{
		TotalEntries:   memStats.TotalEntries + diskStats.TotalEntries,
		MemoryUsageMB:  memStats.MemoryUsageMB,
		DiskUsageMB:    diskStats.DiskUsageMB,
		HitRatio:       (memStats.HitRatio + diskStats.HitRatio) / 2,
		MissRatio:      (memStats.MissRatio + diskStats.MissRatio) / 2,
		TotalHits:      memStats.TotalHits + diskStats.TotalHits,
		TotalMisses:    memStats.TotalMisses + diskStats.TotalMisses,
		ExpiredEntries: memStats.ExpiredEntries + diskStats.ExpiredEntries,
		EvictedEntries: memStats.EvictedEntries + diskStats.EvictedEntries,
	}
}

// Close closes both caches
func (hc *HybridCache) Close() error {
	hc.memory.Close()
	return hc.disk.Close()
}

// NoOpCache implements a no-operation cache for when caching is disabled
type NoOpCache struct{}

func (nc *NoOpCache) Get(key string) (interface{}, bool)                         { return nil, false }
func (nc *NoOpCache) Set(key string, value interface{}, ttl time.Duration) error { return nil }
func (nc *NoOpCache) SetWithType(key string, value interface{}, itemType string, ttl time.Duration) error {
	return nil
}
func (nc *NoOpCache) Delete(key string) error { return nil }
func (nc *NoOpCache) Clear() error            { return nil }
func (nc *NoOpCache) Keys() []string          { return []string{} }
func (nc *NoOpCache) Stats() CacheStatistics  { return CacheStatistics{} }
func (nc *NoOpCache) Close() error            { return nil }

// Helper functions

// estimateSize estimates the size of an object for cache cost calculation
func estimateSize(obj interface{}) int64 {
	data, err := json.Marshal(obj)
	if err != nil {
		return 1024 // Default size estimate
	}
	return int64(len(data))
}

// CacheKey generates a consistent cache key from components
func CacheKey(components ...string) string {
	return strings.Join(components, ":")
}

// Global cache manager
var globalCache *CacheManager

// InitGlobalCache initializes the global cache manager
func InitGlobalCache(config CacheConfig) error {
	cache, err := NewCacheManager(config)
	if err != nil {
		return err
	}
	globalCache = cache
	return nil
}

// GetGlobalCache returns the global cache manager
func GetGlobalCache() *CacheManager {
	if globalCache == nil {
		// Fallback to disabled cache
		config := CacheConfig{Enabled: false}
		globalCache, _ = NewCacheManager(config)
	}
	return globalCache
}

// Convenience functions using global cache
func CacheGet(key string) (interface{}, bool) {
	return GetGlobalCache().Get(key)
}

func CacheSet(key string, value interface{}) error {
	return GetGlobalCache().Set(key, value)
}

func CacheSetWithTTL(key string, value interface{}, ttl time.Duration) error {
	return GetGlobalCache().SetWithTTL(key, value, ttl)
}

func CacheDelete(key string) error {
	return GetGlobalCache().Delete(key)
}

func CacheStats() CacheStatistics {
	return GetGlobalCache().Stats()
}
