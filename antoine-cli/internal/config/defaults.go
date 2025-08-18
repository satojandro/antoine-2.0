// // Package config provides default configuration values for Antoine CLI
// // This file contains all default settings and fallback values
package config

//
//import (
//	"time"
//)
//
//// Default configuration values for Antoine CLI
//// These values are used when no configuration file is present
//// or when specific values are not set in the configuration
//
//// Application defaults
//const (
//	DefaultAppName        = "antoine-cli"
//	DefaultAppVersion     = "1.0.0"
//	DefaultAppDescription = "Your Ultimate Hackathon Mentor"
//	DefaultAppAuthor      = "Antoine Team"
//)
//
//// API defaults
//const (
//	DefaultAPIBaseURL    = "https://api.antoine.ai"
//	DefaultAPITimeout    = 30 * time.Second
//	DefaultAPIRetryCount = 3
//	DefaultAPIRetryDelay = 1 * time.Second
//	DefaultAPIRateLimit  = 100
//	DefaultAPIUserAgent  = "Antoine-CLI/1.0.0"
//)
//
//// MCP (Model Context Protocol) defaults
//const (
//	DefaultMCPTimeout        = 30 * time.Second
//	DefaultMCPRetryCount     = 3
//	DefaultMCPRetryDelay     = 2 * time.Second
//	DefaultMCPMaxConnections = 10
//	DefaultMCPKeepAlive      = true
//)
//
//// MCP Server endpoints and configurations
//var DefaultMCPServers = map[string]MCPServerConfig{
//	"exa": {
//		Endpoint:    "mcp://localhost:8001",
//		Description: "Semantic web search and content discovery",
//		Enabled:     true,
//		Timeout:     30 * time.Second,
//		Features:    []string{"search_hackathons", "search_projects", "trend_analysis"},
//	},
//	"github": {
//		Endpoint:    "mcp://localhost:8002",
//		Description: "GitHub repository analysis and tools",
//		Enabled:     true,
//		Timeout:     45 * time.Second,
//		Features:    []string{"repo_analysis", "file_reading", "commit_history", "issue_tracking"},
//	},
//	"deepwiki": {
//		Endpoint:    "mcp://localhost:8003",
//		Description: "Intelligent repository summarization",
//		Enabled:     true,
//		Timeout:     60 * time.Second,
//		Features:    []string{"repo_overview", "documentation_generation", "code_explanation"},
//	},
//	"e2b": {
//		Endpoint:    "mcp://localhost:8004",
//		Description: "Code execution in sandboxed environments",
//		Enabled:     true,
//		Timeout:     120 * time.Second,
//		Features:    []string{"code_execution", "data_analysis", "custom_scripts"},
//	},
//	"browserbase": {
//		Endpoint:    "mcp://localhost:8005",
//		Description: "Automated web browsing and data extraction",
//		Enabled:     false, // Disabled by default
//		Timeout:     90 * time.Second,
//		Features:    []string{"web_automation", "screenshot_capture", "form_interaction"},
//	},
//	"firecrawl": {
//		Endpoint:    "mcp://localhost:8006",
//		Description: "Structured content extraction from web pages",
//		Enabled:     false, // Disabled by default
//		Timeout:     60 * time.Second,
//		Features:    []string{"content_extraction", "structured_scraping", "markdown_conversion"},
//	},
//}
//
//// MCPServerConfig represents configuration for an MCP server
//type MCPServerConfig struct {
//	Endpoint    string        `yaml:"endpoint"`
//	Description string        `yaml:"description"`
//	Enabled     bool          `yaml:"enabled"`
//	Timeout     time.Duration `yaml:"timeout"`
//	Features    []string      `yaml:"features"`
//}
//
//// UI defaults
//const (
//	DefaultUITheme                     = "dark"
//	DefaultUIAnimations                = true
//	DefaultUIASCIIArt                  = true
//	DefaultUIColors                    = true
//	DefaultUIUnicodeSupport            = true
//	DefaultUIMaxWidth                  = 120
//	DefaultUIMaxHeight                 = 30
//	DefaultUISidebarWidth              = 25
//	DefaultUIShowHelpHints             = true
//	DefaultUIShowKeyBindings           = true
//	DefaultUIConfirmDestructiveActions = true
//)
//
//// UI Table defaults
//const (
//	DefaultTableMaxRows       = 20
//	DefaultTableShowBorders   = true
//	DefaultTableZebraStriping = true
//	DefaultTableSortable      = true
//)
//
//// UI Progress defaults
//const (
//	DefaultProgressShowPercentage = true
//	DefaultProgressShowETA        = true
//	DefaultProgressAnimationSpeed = 100 * time.Millisecond
//)
//
//// Cache defaults
//const (
//	DefaultCacheEnabled         = true
//	DefaultCacheType            = "memory" // Options: memory, disk, hybrid
//	DefaultCacheMaxSizeMB       = 100
//	DefaultCacheMaxEntries      = 1000
//	DefaultCacheCleanupInterval = 1 * time.Hour
//)
//
//// Cache TTL (Time To Live) defaults
//var DefaultCacheTTL = map[string]time.Duration{
//	"hackathons":     30 * time.Minute,
//	"projects":       1 * time.Hour,
//	"repositories":   2 * time.Hour,
//	"trends":         6 * time.Hour,
//	"analysis":       24 * time.Hour,
//	"search_results": 15 * time.Minute,
//}
//
//// Cache disk defaults
//const (
//	DefaultCacheDiskPath        = "~/.antoine/cache"
//	DefaultCacheDiskCompression = true
//)
//
//// Logging defaults
//const (
//	DefaultLogLevel       = "info"
//	DefaultLogFormat      = "text"   // Options: text, json
//	DefaultLogOutput      = "stderr" // Options: stdout, stderr, file
//	DefaultLogDevelopment = false
//	DefaultLogCaller      = false
//	DefaultLogStackTrace  = false
//)
//
//// Log file defaults
//const (
//	DefaultLogFilePath       = "~/.antoine/logs/antoine.log"
//	DefaultLogFileMaxSizeMB  = 10
//	DefaultLogFileMaxBackups = 5
//	DefaultLogFileMaxAgeDays = 30
//	DefaultLogFileCompress   = true
//)
//
//// Analytics defaults
//const (
//	DefaultAnalyticsEnabled       = true
//	DefaultAnalyticsAnonymous     = true
//	DefaultAnalyticsEndpoint      = "https://analytics.antoine.ai"
//	DefaultAnalyticsBatchSize     = 100
//	DefaultAnalyticsFlushInterval = 5 * time.Minute
//)
//
//// Analytics collection defaults
//const (
//	DefaultAnalyticsCollectUsageStats         = true
//	DefaultAnalyticsCollectPerformanceMetrics = true
//	DefaultAnalyticsCollectErrorReports       = true
//	DefaultAnalyticsCollectFeatureUsage       = true
//)
//
//// Analytics privacy defaults
//const (
//	DefaultAnalyticsAnonymizePaths = true
//	DefaultAnalyticsAnonymizeRepos = true
//	DefaultAnalyticsHashUserData   = true
//)
//
//// Search defaults
//const (
//	DefaultSearchLimit    = 20
//	DefaultSearchMaxLimit = 100
//)
//
//// Hackathon search defaults
//const (
//	DefaultHackathonSortBy         = "start_date" // Options: start_date, prize_pool, popularity
//	DefaultHackathonSortOrder      = "desc"       // Options: asc, desc
//	DefaultHackathonIncludePast    = false
//	DefaultHackathonIncludeVirtual = true
//	DefaultHackathonMinPrize       = 0
//)
//
//// Project search defaults
//const (
//	DefaultProjectSortBy       = "popularity" // Options: popularity, recent, stars
//	DefaultProjectSortOrder    = "desc"       // Options: asc, desc
//	DefaultProjectIncludeForks = false
//	DefaultProjectMinStars     = 0
//)
//
//// Analysis defaults
//const (
//	DefaultAnalysisParallelAnalysis  = true
//	DefaultAnalysisMaxConcurrentJobs = 4
//	DefaultAnalysisTimeout           = 10 * time.Minute
//)
//
//// Repository analysis defaults
//const (
//	DefaultRepoAnalysisMaxFileSizeMB       = 10
//	DefaultRepoAnalysisMaxFiles            = 1000
//	DefaultRepoAnalysisIncludeDependencies = true
//	DefaultRepoAnalysisIncludeSecurity     = true
//	DefaultRepoAnalysisIncludePerformance  = true
//	DefaultRepoAnalysisAnalyzeCode         = true
//	DefaultRepoAnalysisAnalyzeDocs         = true
//	DefaultRepoAnalysisAnalyzeConfig       = true
//	DefaultRepoAnalysisShallowClone        = true
//	DefaultRepoAnalysisMaxCommitHistory    = 100
//)
//
//// Trend analysis defaults
//const (
//	DefaultTrendsTimeframe           = "6months" // Options: 1month, 3months, 6months, 1year
//	DefaultTrendsMinDataPoints       = 10
//	DefaultTrendsConfidenceThreshold = 0.8
//)
//
//// Mentor defaults
//const (
//	DefaultMentorModel             = "gpt-4"
//	DefaultMentorMaxTokens         = 2048
//	DefaultMentorTemperature       = 0.7
//	DefaultMentorMaxHistory        = 50
//	DefaultMentorContextWindow     = 10
//	DefaultMentorSaveConversations = true
//	DefaultMentorPersonality       = "helpful"      // Options: helpful, casual, professional, technical
//	DefaultMentorExpertiseLevel    = "intermediate" // Options: beginner, intermediate, advanced
//	DefaultMentorCodeExamples      = true
//	DefaultMentorStepByStep        = true
//	DefaultMentorProvideResources  = true
//	DefaultMentorSessionTimeout    = 30 * time.Minute
//	DefaultMentorAutoSave          = true
//)
//
//// Security defaults
//const (
//	DefaultSecurityCredentialsStorage     = "keyring" // Options: keyring, file, env
//	DefaultSecurityCredentialsEncryption  = true
//	DefaultSecurityCredentialsAutoRefresh = true
//	DefaultSecurityAPIVerifySSL           = true
//	DefaultSecurityAPITimeout             = 30 * time.Second
//	DefaultSecurityAPIMaxRedirects        = 3
//	DefaultSecurityDataAnonymizeLogs      = true
//	DefaultSecurityDataEncryptCache       = false
//	DefaultSecurityDataSecureDelete       = true
//)
//
//// Performance defaults
//const (
//	DefaultPerformanceMaxGoroutines     = 10
//	DefaultPerformanceWorkerPoolSize    = 5
//	DefaultPerformanceMaxMemoryMB       = 512
//	DefaultPerformanceGCPercentage      = 100
//	DefaultPerformanceMaxConnections    = 20
//	DefaultPerformanceConnectionTimeout = 10 * time.Second
//	DefaultPerformanceIdleTimeout       = 5 * time.Minute
//)
//
//// Debug defaults
//const (
//	DefaultDebugEnabled          = false
//	DefaultDebugProfiling        = false
//	DefaultDebugTraceRequests    = false
//	DefaultDebugVerboseLogging   = false
//	DefaultDebugMockMCPServers   = false
//	DefaultDebugMockAPIResponses = false
//	DefaultDebugHotReload        = false
//	DefaultDebugDevMode          = false
//)
//
//// Feature flags defaults
//const (
//	DefaultFeatureExperimentalUI    = false
//	DefaultFeatureBetaAnalysis      = false
//	DefaultFeatureAdvancedMentor    = false
//	DefaultFeatureSearchEnabled     = true
//	DefaultFeatureAnalysisEnabled   = true
//	DefaultFeatureMentorEnabled     = true
//	DefaultFeatureTrendsEnabled     = true
//	DefaultFeatureGitHubIntegration = true
//	DefaultFeatureWebScraping       = false
//	DefaultFeatureCodeExecution     = true
//)
//
//// Notification defaults
//const (
//	DefaultNotificationsEnabled     = true
//	DefaultNotificationsDesktop     = false
//	DefaultNotificationsSound       = false
//	DefaultNotificationsDuration    = 3 * time.Second
//	DefaultNotificationsAutoDismiss = true
//)
//
//// Notification type defaults
//const (
//	DefaultNotificationsSuccess = true
//	DefaultNotificationsWarning = true
//	DefaultNotificationsError   = true
//	DefaultNotificationsInfo    = true
//)
//
//// Update defaults
//const (
//	DefaultUpdatesCheckForUpdates = true
//	DefaultUpdatesAutoUpdate      = false
//	DefaultUpdatesUpdateChannel   = "stable" // Options: stable, beta, alpha
//	DefaultUpdatesCheckInterval   = 24 * time.Hour
//	DefaultUpdatesGitHubReleases  = true
//)
//
//// Backup defaults
//const (
//	DefaultBackupEnabled       = false
//	DefaultBackupConfig        = true
//	DefaultBackupCache         = false
//	DefaultBackupLogs          = false
//	DefaultBackupConversations = true
//	DefaultBackupInterval      = "daily"
//	DefaultBackupRetention     = "30d"
//	DefaultBackupCompression   = true
//	DefaultBackupSyncEnabled   = false
//)
//
//// Directory and file defaults
//var DefaultPaths = map[string]string{
//	"config_dir": "~/.antoine",
//	"cache_dir":  "~/.antoine/cache",
//	"log_dir":    "~/.antoine/logs",
//	"data_dir":   "~/.antoine/data",
//	"temp_dir":   "~/.antoine/tmp",
//}
//
//// Environment variable defaults
//var DefaultEnvVars = map[string]string{
//	"ANTOINE_ENV":               "production",
//	"ANTOINE_LOG_LEVEL":         "info",
//	"ANTOINE_DEBUG_ENABLED":     "false",
//	"ANTOINE_ANALYTICS_ENABLED": "true",
//	"ANTOINE_CACHE_ENABLED":     "true",
//	"ANTOINE_UI_THEME":          "dark",
//	"ANTOINE_NO_COLOR":          "false",
//	"ANTOINE_NO_ANIMATIONS":     "false",
//}
//
//// GetDefaultConfig returns a complete default configuration
//func GetDefaultConfig() *Config {
//	return &Config{
//		App: AppConfig{
//			Name:        DefaultAppName,
//			Version:     DefaultAppVersion,
//			Description: DefaultAppDescription,
//			Author:      DefaultAppAuthor,
//		},
//		API: APIConfig{
//			BaseURL:    DefaultAPIBaseURL,
//			Timeout:    DefaultAPITimeout,
//			RetryCount: DefaultAPIRetryCount,
//			RetryDelay: DefaultAPIRetryDelay,
//			RateLimit:  DefaultAPIRateLimit,
//			UserAgent:  DefaultAPIUserAgent,
//		},
//		MCP: MCPConfig{
//			Timeout:        DefaultMCPTimeout,
//			RetryCount:     DefaultMCPRetryCount,
//			RetryDelay:     DefaultMCPRetryDelay,
//			MaxConnections: DefaultMCPMaxConnections,
//			KeepAlive:      DefaultMCPKeepAlive,
//			Servers:        DefaultMCPServers,
//		},
//		UI: UIConfig{
//			Theme:                     DefaultUITheme,
//			Animations:                DefaultUIAnimations,
//			ASCIIArt:                  DefaultUIASCIIArt,
//			Colors:                    DefaultUIColors,
//			UnicodeSupport:            DefaultUIUnicodeSupport,
//			MaxWidth:                  DefaultUIMaxWidth,
//			MaxHeight:                 DefaultUIMaxHeight,
//			SidebarWidth:              DefaultUISidebarWidth,
//			ShowHelpHints:             DefaultUIShowHelpHints,
//			ShowKeyBindings:           DefaultUIShowKeyBindings,
//			ConfirmDestructiveActions: DefaultUIConfirmDestructiveActions,
//		},
//		Cache: CacheConfig{
//			Enabled:         DefaultCacheEnabled,
//			Type:            DefaultCacheType,
//			MaxSizeMB:       DefaultCacheMaxSizeMB,
//			MaxEntries:      DefaultCacheMaxEntries,
//			CleanupInterval: DefaultCacheCleanupInterval,
//			TTL:             DefaultCacheTTL,
//			Disk: CacheDiskConfig{
//				Path:        DefaultCacheDiskPath,
//				Compression: DefaultCacheDiskCompression,
//			},
//		},
//		Logging: LoggingConfig{
//			Level:       DefaultLogLevel,
//			Format:      DefaultLogFormat,
//			Output:      DefaultLogOutput,
//			Development: DefaultLogDevelopment,
//			Caller:      DefaultLogCaller,
//			StackTrace:  DefaultLogStackTrace,
//			File: LogFileConfig{
//				Path:       DefaultLogFilePath,
//				MaxSizeMB:  DefaultLogFileMaxSizeMB,
//				MaxBackups: DefaultLogFileMaxBackups,
//				MaxAgeDays: DefaultLogFileMaxAgeDays,
//				Compress:   DefaultLogFileCompress,
//			},
//		},
//		Analytics: AnalyticsConfig{
//			Enabled:       DefaultAnalyticsEnabled,
//			Anonymous:     DefaultAnalyticsAnonymous,
//			Endpoint:      DefaultAnalyticsEndpoint,
//			BatchSize:     DefaultAnalyticsBatchSize,
//			FlushInterval: DefaultAnalyticsFlushInterval,
//			Collect: AnalyticsCollectConfig{
//				UsageStats:         DefaultAnalyticsCollectUsageStats,
//				PerformanceMetrics: DefaultAnalyticsCollectPerformanceMetrics,
//				ErrorReports:       DefaultAnalyticsCollectErrorReports,
//				FeatureUsage:       DefaultAnalyticsCollectFeatureUsage,
//			},
//			AnonymizePaths: DefaultAnalyticsAnonymizePaths,
//			AnonymizeRepos: DefaultAnalyticsAnonymizeRepos,
//			HashUserData:   DefaultAnalyticsHashUserData,
//		},
//		Search: SearchConfig{
//			DefaultLimit: DefaultSearchLimit,
//			MaxLimit:     DefaultSearchMaxLimit,
//			Hackathons: HackathonSearchConfig{
//				SortBy:         DefaultHackathonSortBy,
//				SortOrder:      DefaultHackathonSortOrder,
//				IncludePast:    DefaultHackathonIncludePast,
//				IncludeVirtual: DefaultHackathonIncludeVirtual,
//				MinPrize:       DefaultHackathonMinPrize,
//			},
//			Projects: ProjectSearchConfig{
//				SortBy:       DefaultProjectSortBy,
//				SortOrder:    DefaultProjectSortOrder,
//				IncludeForks: DefaultProjectIncludeForks,
//				MinStars:     DefaultProjectMinStars,
//			},
//		},
//		Analysis: AnalysisConfig{
//			ParallelAnalysis:  DefaultAnalysisParallelAnalysis,
//			MaxConcurrentJobs: DefaultAnalysisMaxConcurrentJobs,
//			Timeout:           DefaultAnalysisTimeout,
//			Repository: RepositoryAnalysisConfig{
//				MaxFileSizeMB:       DefaultRepoAnalysisMaxFileSizeMB,
//				MaxFiles:            DefaultRepoAnalysisMaxFiles,
//				IncludeDependencies: DefaultRepoAnalysisIncludeDependencies,
//				IncludeSecurity:     DefaultRepoAnalysisIncludeSecurity,
//				IncludePerformance:  DefaultRepoAnalysisIncludePerformance,
//				AnalyzeCode:         DefaultRepoAnalysisAnalyzeCode,
//				AnalyzeDocs:         DefaultRepoAnalysisAnalyzeDocs,
//				AnalyzeConfig:       DefaultRepoAnalysisAnalyzeConfig,
//				ShallowClone:        DefaultRepoAnalysisShallowClone,
//				MaxCommitHistory:    DefaultRepoAnalysisMaxCommitHistory,
//			},
//			Trends: TrendsAnalysisConfig{
//				DefaultTimeframe:    DefaultTrendsTimeframe,
//				MinDataPoints:       DefaultTrendsMinDataPoints,
//				ConfidenceThreshold: DefaultTrendsConfidenceThreshold,
//			},
//		},
//		Mentor: MentorConfig{
//			Model:             DefaultMentorModel,
//			MaxTokens:         DefaultMentorMaxTokens,
//			Temperature:       DefaultMentorTemperature,
//			MaxHistory:        DefaultMentorMaxHistory,
//			ContextWindow:     DefaultMentorContextWindow,
//			SaveConversations: DefaultMentorSaveConversations,
//			Personality:       DefaultMentorPersonality,
//			ExpertiseLevel:    DefaultMentorExpertiseLevel,
//			CodeExamples:      DefaultMentorCodeExamples,
//			StepByStep:        DefaultMentorStepByStep,
//			ProvideResources:  DefaultMentorProvideResources,
//			SessionTimeout:    DefaultMentorSessionTimeout,
//			AutoSave:          DefaultMentorAutoSave,
//		},
//		Security: SecurityConfig{
//			Credentials: CredentialsConfig{
//				Storage:     DefaultSecurityCredentialsStorage,
//				Encryption:  DefaultSecurityCredentialsEncryption,
//				AutoRefresh: DefaultSecurityCredentialsAutoRefresh,
//			},
//			API: SecurityAPIConfig{
//				VerifySSL:    DefaultSecurityAPIVerifySSL,
//				Timeout:      DefaultSecurityAPITimeout,
//				MaxRedirects: DefaultSecurityAPIMaxRedirects,
//			},
//			Data: SecurityDataConfig{
//				AnonymizeLogs: DefaultSecurityDataAnonymizeLogs,
//				EncryptCache:  DefaultSecurityDataEncryptCache,
//				SecureDelete:  DefaultSecurityDataSecureDelete,
//			},
//		},
//		Performance: PerformanceConfig{
//			MaxGoroutines:     DefaultPerformanceMaxGoroutines,
//			WorkerPoolSize:    DefaultPerformanceWorkerPoolSize,
//			MaxMemoryMB:       DefaultPerformanceMaxMemoryMB,
//			GCPercentage:      DefaultPerformanceGCPercentage,
//			MaxConnections:    DefaultPerformanceMaxConnections,
//			ConnectionTimeout: DefaultPerformanceConnectionTimeout,
//			IdleTimeout:       DefaultPerformanceIdleTimeout,
//		},
//		Debug: DebugConfig{
//			Enabled:          DefaultDebugEnabled,
//			Profiling:        DefaultDebugProfiling,
//			TraceRequests:    DefaultDebugTraceRequests,
//			VerboseLogging:   DefaultDebugVerboseLogging,
//			MockMCPServers:   DefaultDebugMockMCPServers,
//			MockAPIResponses: DefaultDebugMockAPIResponses,
//			HotReload:        DefaultDebugHotReload,
//			DevMode:          DefaultDebugDevMode,
//		},
//		Features: FeaturesConfig{
//			ExperimentalUI:    DefaultFeatureExperimentalUI,
//			BetaAnalysis:      DefaultFeatureBetaAnalysis,
//			AdvancedMentor:    DefaultFeatureAdvancedMentor,
//			SearchEnabled:     DefaultFeatureSearchEnabled,
//			AnalysisEnabled:   DefaultFeatureAnalysisEnabled,
//			MentorEnabled:     DefaultFeatureMentorEnabled,
//			TrendsEnabled:     DefaultFeatureTrendsEnabled,
//			GitHubIntegration: DefaultFeatureGitHubIntegration,
//			WebScraping:       DefaultFeatureWebScraping,
//			CodeExecution:     DefaultFeatureCodeExecution,
//		},
//		Notifications: NotificationsConfig{
//			Enabled:     DefaultNotificationsEnabled,
//			Desktop:     DefaultNotificationsDesktop,
//			Sound:       DefaultNotificationsSound,
//			Duration:    DefaultNotificationsDuration,
//			AutoDismiss: DefaultNotificationsAutoDismiss,
//			Types: NotificationTypesConfig{
//				Success: DefaultNotificationsSuccess,
//				Warning: DefaultNotificationsWarning,
//				Error:   DefaultNotificationsError,
//				Info:    DefaultNotificationsInfo,
//			},
//		},
//		Updates: UpdatesConfig{
//			CheckForUpdates: DefaultUpdatesCheckForUpdates,
//			AutoUpdate:      DefaultUpdatesAutoUpdate,
//			UpdateChannel:   DefaultUpdatesUpdateChannel,
//			CheckInterval:   DefaultUpdatesCheckInterval,
//			GitHubReleases:  DefaultUpdatesGitHubReleases,
//		},
//		Backup: BackupConfig{
//			Enabled:       DefaultBackupEnabled,
//			Config:        DefaultBackupConfig,
//			Cache:         DefaultBackupCache,
//			Logs:          DefaultBackupLogs,
//			Conversations: DefaultBackupConversations,
//			Interval:      DefaultBackupInterval,
//			Retention:     DefaultBackupRetention,
//			Compression:   DefaultBackupCompression,
//			SyncEnabled:   DefaultBackupSyncEnabled,
//		},
//	}
//}
//
//// Helper functions for environment-specific defaults
//
//// GetDevelopmentDefaults returns defaults optimized for development
//func GetDevelopmentDefaults() map[string]interface{} {
//	return map[string]interface{}{
//		"logging.level":                  "debug",
//		"logging.development":            true,
//		"logging.caller":                 true,
//		"debug.enabled":                  true,
//		"debug.verbose_logging":          true,
//		"analytics.enabled":              false,
//		"cache.ttl.hackathons":           "5m",
//		"cache.ttl.projects":             "10m",
//		"ui.confirm_destructive_actions": false,
//		"api.timeout":                    "60s",
//		"performance.max_goroutines":     5,
//	}
//}
//
//// GetProductionDefaults returns defaults optimized for production
//func GetProductionDefaults() map[string]interface{} {
//	return map[string]interface{}{
//		"logging.level":                   "info",
//		"logging.format":                  "json",
//		"logging.output":                  "file",
//		"debug.enabled":                   false,
//		"analytics.enabled":               true,
//		"security.credentials.encryption": true,
//		"security.api.verify_ssl":         true,
//		"cache.type":                      "hybrid",
//		"performance.max_memory_mb":       256,
//		"features.experimental_ui":        false,
//		"features.beta_analysis":          false,
//	}
//}
