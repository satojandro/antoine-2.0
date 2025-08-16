package models

import "time"

type AnalysisRequest struct {
	Type      string                 `json:"type"` // repo, trends, market
	Target    string                 `json:"target"`
	Options   AnalysisOptions        `json:"options"`
	Timestamp time.Time              `json:"timestamp"`
	RequestID string                 `json:"request_id"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type AnalysisOptions struct {
	Depth               string   `json:"depth"`
	IncludeDependencies bool     `json:"include_dependencies"`
	Focus               []string `json:"focus"`
	Timeframe           string   `json:"timeframe"`
	Technologies        []string `json:"technologies"`
	Market              string   `json:"market"`
	CompareWith         []string `json:"compare_with"`
}

type AnalysisResult struct {
	ID              string                 `json:"id"`
	RequestID       string                 `json:"request_id"`
	Type            string                 `json:"type"`
	Status          string                 `json:"status"`
	Progress        int                    `json:"progress"`
	Results         interface{}            `json:"results"`
	Summary         string                 `json:"summary"`
	Insights        []Insight              `json:"insights"`
	Recommendations []Recommendation       `json:"recommendations"`
	StartTime       time.Time              `json:"start_time"`
	EndTime         *time.Time             `json:"end_time,omitempty"`
	Duration        time.Duration          `json:"duration"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type Insight struct {
	Type        string      `json:"type"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Impact      string      `json:"impact"` // high, medium, low
	Confidence  float64     `json:"confidence"`
	Data        interface{} `json:"data"`
}

type Recommendation struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"` // high, medium, low
	Effort      string   `json:"effort"`   // low, medium, high
	Impact      string   `json:"impact"`   // low, medium, high
	Category    string   `json:"category"`
	Actions     []Action `json:"actions"`
}

type Action struct {
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	Resources    []string `json:"resources"`
	TimeEstimate string   `json:"time_estimate"`
}
