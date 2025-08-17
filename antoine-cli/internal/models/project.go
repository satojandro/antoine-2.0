package models

import "time"

type Project struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	ShortDescription string                 `json:"short_description"`
	HackathonID      string                 `json:"hackathon_id"`
	HackathonName    string                 `json:"hackathon_name"`
	Team             Team                   `json:"team"`
	Repository       Repository             `json:"repository"`
	LiveURL          string                 `json:"live_url,omitempty"`
	DemoURL          string                 `json:"demo_url,omitempty"`
	PitchURL         string                 `json:"pitch_url,omitempty"`
	Technologies     []string               `json:"technologies"`
	Categories       []string               `json:"categories"`
	Tags             []string               `json:"tags"`
	Awards           []Award                `json:"awards"`
	Metrics          ProjectMetrics         `json:"metrics"`
	Feedback         []Feedback             `json:"feedback"`
	Status           string                 `json:"status"`
	SubmissionDate   time.Time              `json:"submission_date"`
	LastUpdated      time.Time              `json:"last_updated"`
	Media            MediaAssets            `json:"media"`
	Innovation       InnovationScore        `json:"innovation"`
	TechnicalDepth   TechnicalAnalysis      `json:"technical_depth"`
	MarketPotential  MarketAnalysis         `json:"market_potential"`
	Metadata         map[string]interface{} `json:"metadata"`
}

type Team struct {
	Members []TeamMember `json:"members"`
	Size    int          `json:"size"`
	Lead    string       `json:"lead"`
}

type TeamMember struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Role     string   `json:"role"`
	Skills   []string `json:"skills"`
	GitHub   string   `json:"github,omitempty"`
	LinkedIn string   `json:"linkedin,omitempty"`
	Twitter  string   `json:"twitter,omitempty"`
}

type Repository struct {
	URL         string         `json:"url"`
	Platform    string         `json:"platform"` // github, gitlab, etc.
	Stars       int            `json:"stars"`
	Forks       int            `json:"forks"`
	Language    string         `json:"primary_language"`
	Languages   map[string]int `json:"languages"`
	Commits     int            `json:"commits"`
	LastCommit  time.Time      `json:"last_commit"`
	Size        int            `json:"size_kb"`
	License     string         `json:"license"`
	Topics      []string       `json:"topics"`
	CodeQuality CodeQuality    `json:"code_quality"`
}

type CodeQuality struct {
	Score           float64           `json:"score"`
	Complexity      ComplexityMetrics `json:"complexity"`
	Documentation   float64           `json:"documentation_coverage"`
	TestCoverage    float64           `json:"test_coverage"`
	Security        SecurityAnalysis  `json:"security"`
	Maintainability float64           `json:"maintainability"`
	Architecture    ArchitectureScore `json:"architecture"`
}

type ComplexityMetrics struct {
	Cyclomatic float64 `json:"cyclomatic"`
	Cognitive  float64 `json:"cognitive"`
	Lines      int     `json:"lines_of_code"`
	Files      int     `json:"file_count"`
	Functions  int     `json:"function_count"`
}

type SecurityAnalysis struct {
	Score           float64            `json:"score"`
	Vulnerabilities []Vulnerability    `json:"vulnerabilities"`
	Dependencies    DependencyAnalysis `json:"dependencies"`
}

type Vulnerability struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	File        string `json:"file"`
	Line        int    `json:"line"`
}

type DependencyAnalysis struct {
	Total      int                `json:"total"`
	Outdated   int                `json:"outdated"`
	Vulnerable int                `json:"vulnerable"`
	Licenses   map[string]int     `json:"licenses"`
	Details    []DependencyDetail `json:"details"`
}

type DependencyDetail struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Latest     string `json:"latest"`
	License    string `json:"license"`
	Vulnerable bool   `json:"vulnerable"`
	Severity   string `json:"severity,omitempty"`
}

type ArchitectureScore struct {
	Score      float64         `json:"score"`
	Patterns   []string        `json:"patterns"`
	Violations []string        `json:"violations"`
	Modularity float64         `json:"modularity"`
	Coupling   CouplingMetrics `json:"coupling"`
}

type CouplingMetrics struct {
	Afferent  int     `json:"afferent"`
	Efferent  int     `json:"efferent"`
	Stability float64 `json:"stability"`
}

type Award struct {
	Position    string `json:"position"`
	Category    string `json:"category"`
	Prize       int    `json:"prize"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Sponsor     string `json:"sponsor,omitempty"`
}

type ProjectMetrics struct {
	Views      int     `json:"views"`
	Likes      int     `json:"likes"`
	Comments   int     `json:"comments"`
	Shares     int     `json:"shares"`
	Downloads  int     `json:"downloads"`
	Forks      int     `json:"forks"`
	Stars      int     `json:"stars"`
	Popularity float64 `json:"popularity_score"`
}

type Feedback struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Role      string    `json:"role"` // judge, mentor, participant
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	Category  string    `json:"category"`
	Timestamp time.Time `json:"timestamp"`
	Helpful   int       `json:"helpful_votes"`
}

type MediaAssets struct {
	Screenshots []string `json:"screenshots"`
	Videos      []string `json:"videos"`
	Documents   []string `json:"documents"`
	Logo        string   `json:"logo,omitempty"`
}

type InnovationScore struct {
	Score       float64  `json:"score"`
	Factors     []string `json:"factors"`
	Novelty     float64  `json:"novelty"`
	Creativity  float64  `json:"creativity"`
	Impact      float64  `json:"impact"`
	Feasibility float64  `json:"feasibility"`
}

type TechnicalAnalysis struct {
	Score          float64     `json:"score"`
	Complexity     float64     `json:"complexity"`
	Implementation float64     `json:"implementation_quality"`
	Innovation     float64     `json:"technical_innovation"`
	Scalability    float64     `json:"scalability"`
	Performance    Performance `json:"performance"`
}

type Performance struct {
	Score        float64 `json:"score"`
	LoadTime     float64 `json:"load_time_ms"`
	ResponseTime float64 `json:"response_time_ms"`
	Throughput   float64 `json:"throughput_rps"`
	Memory       float64 `json:"memory_usage_mb"`
	CPU          float64 `json:"cpu_usage_percent"`
}

type MarketAnalysis struct {
	Score       float64  `json:"score"`
	MarketSize  string   `json:"market_size"`
	Competition int      `json:"competition_level"`
	Opportunity float64  `json:"opportunity_score"`
	Viability   float64  `json:"commercial_viability"`
	UserNeed    float64  `json:"user_need_score"`
	Trends      []string `json:"relevant_trends"`
	Competitors []string `json:"competitors"`
	Advantages  []string `json:"competitive_advantages"`
}
