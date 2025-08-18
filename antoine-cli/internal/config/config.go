package config

import (
	"antoine-cli/pkg/ascii"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Config represents the complete application configuration
type Config struct {
	App           AppConfig           `mapstructure:"app"`
	API           APIConfig           `mapstructure:"api"`
	MCP           MCPConfig           `mapstructure:"mcp"`
	UI            UIConfig            `mapstructure:"ui"`
	Cache         CacheConfig         `mapstructure:"cache"`
	Logging       LoggingConfig       `mapstructure:"logging"`
	Analytics     AnalyticsConfig     `mapstructure:"analytics"`
	Search        SearchConfig        `mapstructure:"search"`
	Analysis      AnalysisConfig      `mapstructure:"analysis"`
	Mentor        MentorConfig        `mapstructure:"mentor"`
	Security      SecurityConfig      `mapstructure:"security"`
	Performance   PerformanceConfig   `mapstructure:"performance"`
	Debug         DebugConfig         `mapstructure:"debug"`
	Features      FeaturesConfig      `mapstructure:"features"`
	Notifications NotificationsConfig `mapstructure:"notifications"`
}

// AppConfig represents application-level configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Description string `mapstructure:"description"`
	Author      string `mapstructure:"author"`
}

// APIConfig represents API client configuration
type APIConfig struct {
	BaseURL    string `mapstructure:"base_url"`
	Timeout    string `mapstructure:"timeout"`
	RetryCount int    `mapstructure:"retry_count"`
	RateLimit  int    `mapstructure:"rate_limit"`
	UserAgent  string `mapstructure:"user_agent"`
}

// MCPConfig represents MCP (Model Context Protocol) configuration
type MCPConfig struct {
	Servers        map[string]MCPServerConfig `mapstructure:"servers"`
	Timeout        string                     `mapstructure:"timeout"`
	RetryCount     int                        `mapstructure:"retry_count"`
	MaxConnections int                        `mapstructure:"max_connections"`
	KeepAlive      bool                       `mapstructure:"keep_alive"`
}

// MCPServerConfig represents configuration for an MCP server
type MCPServerConfig struct {
	Endpoint    string   `mapstructure:"endpoint"`
	Description string   `mapstructure:"description"`
	Enabled     bool     `mapstructure:"enabled"`
	Timeout     string   `mapstructure:"timeout"`
	Features    []string `mapstructure:"features"`
}

// UIConfig represents user interface configuration
type UIConfig struct {
	Theme                     string `mapstructure:"theme"`
	Animations                bool   `mapstructure:"animations"`
	ASCIIArt                  bool   `mapstructure:"ascii_art"`
	Colors                    bool   `mapstructure:"colors"`
	UnicodeSupport            bool   `mapstructure:"unicode_support"`
	MaxWidth                  int    `mapstructure:"max_width"`
	MaxHeight                 int    `mapstructure:"max_height"`
	ShowHelpHints             bool   `mapstructure:"show_help_hints"`
	ConfirmDestructiveActions bool   `mapstructure:"confirm_destructive_actions"`
}

// CacheConfig represents caching configuration
type CacheConfig struct {
	Enabled         bool              `mapstructure:"enabled"`
	Type            string            `mapstructure:"type"`
	TTL             string            `mapstructure:"ttl"`
	MaxSize         int               `mapstructure:"max_size"`
	MaxSizeMB       int               `mapstructure:"max_size_mb"`
	CleanupInterval string            `mapstructure:"cleanup_interval"`
	TTLByType       map[string]string `mapstructure:"ttl_by_type"`
	Disk            CacheDiskConfig   `mapstructure:"disk"`
}

// CacheDiskConfig represents disk cache configuration
type CacheDiskConfig struct {
	Path        string `mapstructure:"path"`
	Compression bool   `mapstructure:"compression"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level      string        `mapstructure:"level"`
	Format     string        `mapstructure:"format"`
	Output     string        `mapstructure:"output"`
	File       LogFileConfig `mapstructure:"file"`
	Caller     bool          `mapstructure:"caller"`
	StackTrace bool          `mapstructure:"stack_trace"`
}

// LogFileConfig represents log file configuration
type LogFileConfig struct {
	Path       string `mapstructure:"path"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
	Compress   bool   `mapstructure:"compress"`
}

// AnalyticsConfig represents analytics configuration
type AnalyticsConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	Anonymous     bool   `mapstructure:"anonymous"`
	Endpoint      string `mapstructure:"endpoint"`
	BatchSize     int    `mapstructure:"batch_size"`
	FlushInterval string `mapstructure:"flush_interval"`
}

// SearchConfig represents search configuration
type SearchConfig struct {
	DefaultLimit int                   `mapstructure:"default_limit"`
	MaxLimit     int                   `mapstructure:"max_limit"`
	Hackathons   HackathonSearchConfig `mapstructure:"hackathons"`
	Projects     ProjectSearchConfig   `mapstructure:"projects"`
}

// HackathonSearchConfig represents hackathon search configuration
type HackathonSearchConfig struct {
	SortBy         string `mapstructure:"sort_by"`
	SortOrder      string `mapstructure:"sort_order"`
	IncludePast    bool   `mapstructure:"include_past"`
	IncludeVirtual bool   `mapstructure:"include_virtual"`
	MinPrize       int    `mapstructure:"min_prize"`
}

// ProjectSearchConfig represents project search configuration
type ProjectSearchConfig struct {
	SortBy       string `mapstructure:"sort_by"`
	SortOrder    string `mapstructure:"sort_order"`
	IncludeForks bool   `mapstructure:"include_forks"`
	MinStars     int    `mapstructure:"min_stars"`
}

// AnalysisConfig represents analysis configuration
type AnalysisConfig struct {
	ParallelAnalysis  bool                     `mapstructure:"parallel_analysis"`
	MaxConcurrentJobs int                      `mapstructure:"max_concurrent_jobs"`
	Timeout           string                   `mapstructure:"timeout"`
	Repository        RepositoryAnalysisConfig `mapstructure:"repository"`
}

// RepositoryAnalysisConfig represents repository analysis configuration
type RepositoryAnalysisConfig struct {
	MaxFileSizeMB       int  `mapstructure:"max_file_size_mb"`
	MaxFiles            int  `mapstructure:"max_files"`
	IncludeDependencies bool `mapstructure:"include_dependencies"`
	IncludeSecurity     bool `mapstructure:"include_security"`
	IncludePerformance  bool `mapstructure:"include_performance"`
	ShallowClone        bool `mapstructure:"shallow_clone"`
}

// MentorConfig represents AI mentor configuration
type MentorConfig struct {
	Model             string  `mapstructure:"model"`
	MaxTokens         int     `mapstructure:"max_tokens"`
	Temperature       float64 `mapstructure:"temperature"`
	MaxHistory        int     `mapstructure:"max_history"`
	SaveConversations bool    `mapstructure:"save_conversations"`
	Personality       string  `mapstructure:"personality"`
	SessionTimeout    string  `mapstructure:"session_timeout"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	CredentialsStorage string `mapstructure:"credentials_storage"`
	EncryptCredentials bool   `mapstructure:"encrypt_credentials"`
	VerifySSL          bool   `mapstructure:"verify_ssl"`
	AnonymizeLogs      bool   `mapstructure:"anonymize_logs"`
}

// PerformanceConfig represents performance configuration
type PerformanceConfig struct {
	MaxGoroutines     int    `mapstructure:"max_goroutines"`
	MaxMemoryMB       int    `mapstructure:"max_memory_mb"`
	MaxConnections    int    `mapstructure:"max_connections"`
	ConnectionTimeout string `mapstructure:"connection_timeout"`
}

// DebugConfig represents debug configuration
type DebugConfig struct {
	Enabled          bool `mapstructure:"enabled"`
	VerboseLogging   bool `mapstructure:"verbose_logging"`
	MockMCPServers   bool `mapstructure:"mock_mcp_servers"`
	MockAPIResponses bool `mapstructure:"mock_api_responses"`
}

// FeaturesConfig represents feature flags configuration
type FeaturesConfig struct {
	SearchEnabled     bool `mapstructure:"search_enabled"`
	AnalysisEnabled   bool `mapstructure:"analysis_enabled"`
	MentorEnabled     bool `mapstructure:"mentor_enabled"`
	GitHubIntegration bool `mapstructure:"github_integration"`
	CodeExecution     bool `mapstructure:"code_execution"`
}

// NotificationsConfig represents notifications configuration
type NotificationsConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Desktop  bool   `mapstructure:"desktop"`
	Sound    bool   `mapstructure:"sound"`
	Duration string `mapstructure:"duration"`
}

// Get returns the current configuration
func Get() *Config {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Sprintf("Unable to decode config: %v", err))
	}
	return &cfg
}

// SetDefaults configures all default values
func SetDefaults() {
	// App defaults
	viper.SetDefault("app.name", "antoine-cli")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.description", "Your Ultimate Hackathon Mentor")
	viper.SetDefault("app.author", "Antoine Team")

	// API defaults
	viper.SetDefault("api.base_url", "https://api.antoine.ai")
	viper.SetDefault("api.timeout", "30s")
	viper.SetDefault("api.retry_count", 3)
	viper.SetDefault("api.rate_limit", 100)
	viper.SetDefault("api.user_agent", "Antoine-CLI/1.0.0")

	// MCP defaults
	setMCPDefaults()

	// UI defaults
	viper.SetDefault("ui.theme", "dark")
	viper.SetDefault("ui.animations", true)
	viper.SetDefault("ui.ascii_art", true)
	viper.SetDefault("ui.colors", true)
	viper.SetDefault("ui.unicode_support", true)
	viper.SetDefault("ui.max_width", 120)
	viper.SetDefault("ui.max_height", 30)
	viper.SetDefault("ui.show_help_hints", true)
	viper.SetDefault("ui.confirm_destructive_actions", true)

	// Cache defaults
	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.type", "memory")
	viper.SetDefault("cache.ttl", "30m")
	viper.SetDefault("cache.max_size", 1000)
	viper.SetDefault("cache.max_size_mb", 100)
	viper.SetDefault("cache.cleanup_interval", "1h")
	viper.SetDefault("cache.disk.path", "~/.antoine/cache")
	viper.SetDefault("cache.disk.compression", true)

	// Cache TTL by type
	viper.SetDefault("cache.ttl_by_type.hackathons", "30m")
	viper.SetDefault("cache.ttl_by_type.projects", "1h")
	viper.SetDefault("cache.ttl_by_type.repositories", "2h")
	viper.SetDefault("cache.ttl_by_type.trends", "6h")
	viper.SetDefault("cache.ttl_by_type.analysis", "24h")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "text")
	viper.SetDefault("logging.output", "stderr")
	viper.SetDefault("logging.caller", false)
	viper.SetDefault("logging.stack_trace", false)
	viper.SetDefault("logging.file.path", "~/.antoine/logs/antoine.log")
	viper.SetDefault("logging.file.max_size_mb", 10)
	viper.SetDefault("logging.file.max_backups", 5)
	viper.SetDefault("logging.file.max_age_days", 30)
	viper.SetDefault("logging.file.compress", true)

	// Analytics defaults
	viper.SetDefault("analytics.enabled", true)
	viper.SetDefault("analytics.anonymous", true)
	viper.SetDefault("analytics.endpoint", "https://analytics.antoine.ai")
	viper.SetDefault("analytics.batch_size", 100)
	viper.SetDefault("analytics.flush_interval", "5m")

	// Search defaults
	viper.SetDefault("search.default_limit", 20)
	viper.SetDefault("search.max_limit", 100)
	viper.SetDefault("search.hackathons.sort_by", "start_date")
	viper.SetDefault("search.hackathons.sort_order", "desc")
	viper.SetDefault("search.hackathons.include_past", false)
	viper.SetDefault("search.hackathons.include_virtual", true)
	viper.SetDefault("search.hackathons.min_prize", 0)
	viper.SetDefault("search.projects.sort_by", "popularity")
	viper.SetDefault("search.projects.sort_order", "desc")
	viper.SetDefault("search.projects.include_forks", false)
	viper.SetDefault("search.projects.min_stars", 0)

	// Analysis defaults
	viper.SetDefault("analysis.parallel_analysis", true)
	viper.SetDefault("analysis.max_concurrent_jobs", 4)
	viper.SetDefault("analysis.timeout", "10m")
	viper.SetDefault("analysis.repository.max_file_size_mb", 10)
	viper.SetDefault("analysis.repository.max_files", 1000)
	viper.SetDefault("analysis.repository.include_dependencies", true)
	viper.SetDefault("analysis.repository.include_security", true)
	viper.SetDefault("analysis.repository.include_performance", true)
	viper.SetDefault("analysis.repository.shallow_clone", true)

	// Mentor defaults
	viper.SetDefault("mentor.model", "gpt-4")
	viper.SetDefault("mentor.max_tokens", 2048)
	viper.SetDefault("mentor.temperature", 0.7)
	viper.SetDefault("mentor.max_history", 50)
	viper.SetDefault("mentor.save_conversations", true)
	viper.SetDefault("mentor.personality", "helpful")
	viper.SetDefault("mentor.session_timeout", "30m")

	// Security defaults
	viper.SetDefault("security.credentials_storage", "keyring")
	viper.SetDefault("security.encrypt_credentials", true)
	viper.SetDefault("security.verify_ssl", true)
	viper.SetDefault("security.anonymize_logs", true)

	// Performance defaults
	viper.SetDefault("performance.max_goroutines", 10)
	viper.SetDefault("performance.max_memory_mb", 512)
	viper.SetDefault("performance.max_connections", 20)
	viper.SetDefault("performance.connection_timeout", "10s")

	// Debug defaults
	viper.SetDefault("debug.enabled", false)
	viper.SetDefault("debug.verbose_logging", false)
	viper.SetDefault("debug.mock_mcp_servers", false)
	viper.SetDefault("debug.mock_api_responses", false)

	// Feature flags defaults
	viper.SetDefault("features.search_enabled", true)
	viper.SetDefault("features.analysis_enabled", true)
	viper.SetDefault("features.mentor_enabled", true)
	viper.SetDefault("features.github_integration", true)
	viper.SetDefault("features.code_execution", true)

	// Notifications defaults
	viper.SetDefault("notifications.enabled", true)
	viper.SetDefault("notifications.desktop", false)
	viper.SetDefault("notifications.sound", false)
	viper.SetDefault("notifications.duration", "3s")
}

// setMCPDefaults configures MCP server defaults
func setMCPDefaults() {
	viper.SetDefault("mcp.timeout", "30s")
	viper.SetDefault("mcp.retry_count", 3)
	viper.SetDefault("mcp.max_connections", 10)
	viper.SetDefault("mcp.keep_alive", true)

	// MCP Servers with detailed configuration
	servers := map[string]map[string]interface{}{
		"exa": {
			"endpoint":    "mcp://localhost:8001",
			"description": "Semantic web search and content discovery",
			"enabled":     true,
			"timeout":     "30s",
			"features":    []string{"search_hackathons", "search_projects", "trend_analysis"},
		},
		"github": {
			"endpoint":    "mcp://localhost:8002",
			"description": "GitHub repository analysis and tools",
			"enabled":     true,
			"timeout":     "45s",
			"features":    []string{"repo_analysis", "file_reading", "commit_history"},
		},
		"deepwiki": {
			"endpoint":    "mcp://localhost:8003",
			"description": "Intelligent repository summarization",
			"enabled":     true,
			"timeout":     "60s",
			"features":    []string{"repo_overview", "documentation_generation"},
		},
		"e2b": {
			"endpoint":    "mcp://localhost:8004",
			"description": "Code execution in sandboxed environments",
			"enabled":     true,
			"timeout":     "120s",
			"features":    []string{"code_execution", "data_analysis"},
		},
		"browserbase": {
			"endpoint":    "mcp://localhost:8005",
			"description": "Automated web browsing and data extraction",
			"enabled":     false,
			"timeout":     "90s",
			"features":    []string{"web_automation", "screenshot_capture"},
		},
		"firecrawl": {
			"endpoint":    "mcp://localhost:8006",
			"description": "Structured content extraction from web pages",
			"enabled":     false,
			"timeout":     "60s",
			"features":    []string{"content_extraction", "structured_scraping"},
		},
	}

	for name, config := range servers {
		for key, value := range config {
			viper.SetDefault(fmt.Sprintf("mcp.servers.%s.%s", name, key), value)
		}
	}
}

// Show displays the current configuration in a beautiful format
func Show() {
	cfg := Get()

	fmt.Println(ascii.GetBanner("Antoine Configuration", ascii.EmojiCode, 80))
	fmt.Println()

	// API Configuration
	fmt.Println("🌐 API Configuration:")
	fmt.Printf("  Base URL: %s\n", cfg.API.BaseURL)
	fmt.Printf("  Timeout: %s\n", cfg.API.Timeout)
	fmt.Printf("  Retry Count: %d\n", cfg.API.RetryCount)
	fmt.Printf("  Rate Limit: %d req/min\n", cfg.API.RateLimit)
	fmt.Printf("  User Agent: %s\n", cfg.API.UserAgent)
	fmt.Println()

	// MCP Configuration
	fmt.Println("🔗 MCP Servers:")
	for name, server := range cfg.MCP.Servers {
		status := "🔴 Disabled"
		if server.Enabled {
			status = "🟢 Enabled"
			// TODO: En implementación real, verificar conexión actual
		}
		fmt.Printf("  %-12s: %s %s\n", name, server.Endpoint, status)
		if server.Description != "" {
			fmt.Printf("    %s\n", server.Description)
		}
	}
	fmt.Printf("  Timeout: %s, Max Connections: %d\n", cfg.MCP.Timeout, cfg.MCP.MaxConnections)
	fmt.Println()

	// UI Configuration
	fmt.Println("🎨 UI Configuration:")
	fmt.Printf("  Theme: %s, Colors: %v, Animations: %v\n", cfg.UI.Theme, cfg.UI.Colors, cfg.UI.Animations)
	fmt.Printf("  ASCII Art: %v, Unicode: %v\n", cfg.UI.ASCIIArt, cfg.UI.UnicodeSupport)
	fmt.Printf("  Max Size: %dx%d\n", cfg.UI.MaxWidth, cfg.UI.MaxHeight)
	fmt.Println()

	// Cache Configuration
	fmt.Println("💾 Cache Configuration:")
	fmt.Printf("  Enabled: %v, Type: %s\n", cfg.Cache.Enabled, cfg.Cache.Type)
	fmt.Printf("  TTL: %s, Max Size: %d items (%d MB)\n", cfg.Cache.TTL, cfg.Cache.MaxSize, cfg.Cache.MaxSizeMB)
	if len(cfg.Cache.TTLByType) > 0 {
		fmt.Println("  Custom TTL:")
		for dataType, ttl := range cfg.Cache.TTLByType {
			fmt.Printf("    %s: %s\n", dataType, ttl)
		}
	}
	fmt.Println()

	// Analytics Configuration
	fmt.Println("📊 Analytics Configuration:")
	fmt.Printf("  Enabled: %v, Anonymous: %v\n", cfg.Analytics.Enabled, cfg.Analytics.Anonymous)
	fmt.Printf("  Endpoint: %s\n", cfg.Analytics.Endpoint)
	fmt.Println()

	// Feature Flags
	fmt.Println("🚀 Feature Flags:")
	fmt.Printf("  Search: %v, Analysis: %v, Mentor: %v\n",
		cfg.Features.SearchEnabled, cfg.Features.AnalysisEnabled, cfg.Features.MentorEnabled)
	fmt.Printf("  GitHub Integration: %v, Code Execution: %v\n",
		cfg.Features.GitHubIntegration, cfg.Features.CodeExecution)
	fmt.Println()

	fmt.Printf("📁 Config file: %s\n", viper.ConfigFileUsed())
}

// Set updates a configuration value and saves it
func Set(key, value string) error {
	viper.Set(key, value)
	return saveConfig()
}

// Reset restores all configuration to defaults
func Reset() error {
	viper.Reset()
	SetDefaults()
	return saveConfig()
}

// Validate checks if the current configuration is valid
func Validate() error {
	cfg := Get()

	// Validate API configuration
	if cfg.API.BaseURL == "" {
		return fmt.Errorf("api.base_url cannot be empty")
	}

	// Validate MCP servers
	if len(cfg.MCP.Servers) == 0 {
		return fmt.Errorf("at least one MCP server must be configured")
	}

	// Validate enabled servers have endpoints
	for name, server := range cfg.MCP.Servers {
		if server.Enabled && server.Endpoint == "" {
			return fmt.Errorf("MCP server '%s' is enabled but has no endpoint", name)
		}
	}

	// Validate cache configuration
	if cfg.Cache.MaxSize <= 0 {
		return fmt.Errorf("cache.max_size must be greater than 0")
	}

	// Validate UI dimensions
	if cfg.UI.MaxWidth <= 0 || cfg.UI.MaxHeight <= 0 {
		return fmt.Errorf("UI dimensions must be positive")
	}

	return nil
}

// saveConfig saves the current configuration to file
func saveConfig() error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// Create config file if it doesn't exist
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configFile = filepath.Join(home, ".antoine.yaml")
	}
	return viper.WriteConfigAs(configFile)
}

// LoadConfig initializes and loads configuration
func LoadConfig() error {
	// Set configuration file name and paths
	viper.SetConfigName("antoine")
	viper.SetConfigType("yaml")

	// Add configuration paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/.antoine")
	viper.AddConfigPath("/etc/antoine")

	// Set environment variable prefix
	viper.SetEnvPrefix("ANTOINE")
	viper.AutomaticEnv()

	// Set defaults first
	SetDefaults()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, use defaults
			return nil
		}
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Validate configuration
	return Validate()
}
