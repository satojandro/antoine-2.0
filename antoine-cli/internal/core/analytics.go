package core

import (
	"fmt"
	"sync"
	"time"
)

type AnalyticsManager struct {
	metrics map[string]interface{}
	mu      sync.RWMutex
}

func NewAnalyticsManager() *AnalyticsManager {
	return &AnalyticsManager{
		metrics: make(map[string]interface{}),
	}
}

func (am *AnalyticsManager) RecordSearch(searchType string, resultCount int) {
	am.mu.Lock()
	defer am.mu.Unlock()

	key := fmt.Sprintf("searches_%s", searchType)
	if count, exists := am.metrics[key]; exists {
		am.metrics[key] = count.(int) + 1
	} else {
		am.metrics[key] = 1
	}

	am.metrics[fmt.Sprintf("last_search_%s", searchType)] = time.Now()
	am.metrics[fmt.Sprintf("results_%s", searchType)] = resultCount
}

func (am *AnalyticsManager) RecordAnalysis(analysisType, target string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	key := fmt.Sprintf("analysis_%s", analysisType)
	if count, exists := am.metrics[key]; exists {
		am.metrics[key] = count.(int) + 1
	} else {
		am.metrics[key] = 1
	}

	am.metrics[fmt.Sprintf("last_analysis_%s", analysisType)] = time.Now()
}

func (am *AnalyticsManager) RecordTrends(technologies []string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	for _, tech := range technologies {
		key := fmt.Sprintf("trend_requests_%s", tech)
		if count, exists := am.metrics[key]; exists {
			am.metrics[key] = count.(int) + 1
		} else {
			am.metrics[key] = 1
		}
	}
}

func (am *AnalyticsManager) GetMetrics() map[string]interface{} {
	am.mu.RLock()
	defer am.mu.RUnlock()

	result := make(map[string]interface{})
	for k, v := range am.metrics {
		result[k] = v
	}

	return result
}
