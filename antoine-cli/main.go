// Antoine CLI - Your Ultimate Hackathon Mentor
// This is the main entry point for the Antoine CLI application
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"antoine-cli/cmd"
	"antoine-cli/internal/config"
	"antoine-cli/internal/utils"
	"antoine-cli/pkg/ascii"
	"antoine-cli/pkg/terminal"
)

var (
	// Build information (set by ldflags during build)
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"

	// Configuration
	configFile string
	verbose    bool
	quiet      bool
	noColor    bool
	theme      string
)

func main() {
	// Initialize the application
	if err := initializeApp(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize Antoine CLI: %v\n", err)
		os.Exit(1)
	}

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		// Error is already handled by Cobra, just exit
		os.Exit(1)
	}
}

// initializeApp initializes the Antoine CLI application
func initializeApp() error {
	// Detect terminal capabilities first
	termInfo := terminal.DetectTerminal()

	// Initialize configuration
	if err := initConfig(); err != nil {
		return fmt.Errorf("config initialization failed: %w", err)
	}

	// Initialize logging
	if err := initLogging(); err != nil {
		return fmt.Errorf("logging initialization failed: %w", err)
	}

	// Initialize cache
	if err := initCache(); err != nil {
		return fmt.Errorf("cache initialization failed: %w", err)
	}

	// Set up global flags and environment
	setupGlobalEnvironment(termInfo)

	// Display welcome message if appropriate
	if shouldShowWelcome() {
		displayWelcome(termInfo)
	}

	return nil
}

// initConfig initializes the configuration system
func initConfig() error {
	// Set config file search paths
	viper.SetConfigName("antoine")
	viper.SetConfigType("yaml")

	// Add config search paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(home, ".antoine"))
		viper.AddConfigPath(filepath.Join(home, ".config", "antoine"))
	}

	// Set environment variable prefix
	viper.SetEnvPrefix("ANTOINE")
	viper.AutomaticEnv()

	// Set defaults from our config system
	setConfigDefaults()

	// If a config file is specified, use it
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK, we'll use defaults
	}

	return nil
}

// setConfigDefaults sets default configuration values
func setConfigDefaults() {
	defaults := config.GetDefaultConfig()

	// API settings
	viper.SetDefault("api.base_url", defaults.API.BaseURL)
	viper.SetDefault("api.timeout", defaults.API.Timeout)
	viper.SetDefault("api.retry_count", defaults.API.RetryCount)

	// UI settings
	viper.SetDefault("ui.theme", defaults.UI.Theme)
	viper.SetDefault("ui.animations", defaults.UI.Animations)
	viper.SetDefault("ui.colors", defaults.UI.Colors)
	viper.SetDefault("ui.unicode_support", defaults.UI.UnicodeSupport)

	// Logging settings
	viper.SetDefault("logging.level", defaults.Logging.Level)
	viper.SetDefault("logging.format", defaults.Logging.Format)
	viper.SetDefault("logging.output", defaults.Logging.Output)

	// Cache settings
	viper.SetDefault("cache.enabled", defaults.Cache.Enabled)
	viper.SetDefault("cache.type", defaults.Cache.Type)
	viper.SetDefault("cache.max_size_mb", defaults.Cache.MaxSizeMB)

	// Analytics settings
	viper.SetDefault("analytics.enabled", defaults.Analytics.Enabled)
	viper.SetDefault("analytics.anonymous", defaults.Analytics.Anonymous)
}

// initLogging initializes the logging system
func initLogging() error {
	logConfig := utils.LoggerConfig{
		Level:       utils.LogLevel(viper.GetString("logging.level")),
		Format:      utils.LogFormat(viper.GetString("logging.format")),
		Output:      utils.LogOutput(viper.GetString("logging.output")),
		Development: viper.GetBool("debug.enabled"),
		Caller:      viper.GetBool("logging.caller"),
		StackTrace:  viper.GetBool("logging.stack_trace"),
	}

	// Set file logging configuration
	logConfig.File.Path = viper.GetString("logging.file.path")
	logConfig.File.MaxSizeMB = viper.GetInt("logging.file.max_size_mb")
	logConfig.File.MaxBackups = viper.GetInt("logging.file.max_backups")
	logConfig.File.MaxAgeDays = viper.GetInt("logging.file.max_age_days")
	logConfig.File.Compress = viper.GetBool("logging.file.compress")

	// Override with command line flags
	if verbose {
		logConfig.Level = utils.LogLevelDebug
	}
	if quiet {
		logConfig.Level = utils.LogLevelError
	}

	return utils.InitGlobalLogger(logConfig)
}

// initCache initializes the cache system
func initCache() error {
	if !viper.GetBool("cache.enabled") {
		// Initialize with disabled cache
		return utils.InitGlobalCache(utils.CacheConfig{Enabled: false})
	}

	cacheConfig := utils.CacheConfig{
		Enabled:         true,
		Type:            utils.CacheType(viper.GetString("cache.type")),
		MaxSizeMB:       viper.GetInt64("cache.max_size_mb"),
		MaxEntries:      viper.GetInt64("cache.max_entries"),
		CleanupInterval: viper.GetDuration("cache.cleanup_interval"),
		TTL:             make(map[string]time.Duration),
	}

	// Set TTL values
	cacheConfig.TTL["hackathons"] = viper.GetDuration("cache.ttl.hackathons")
	cacheConfig.TTL["projects"] = viper.GetDuration("cache.ttl.projects")
	cacheConfig.TTL["repositories"] = viper.GetDuration("cache.ttl.repositories")
	cacheConfig.TTL["trends"] = viper.GetDuration("cache.ttl.trends")
	cacheConfig.TTL["analysis"] = viper.GetDuration("cache.ttl.analysis")

	// Disk cache settings
	cacheConfig.Disk.Path = viper.GetString("cache.disk.path")
	cacheConfig.Disk.Compression = viper.GetBool("cache.disk.compression")

	return utils.InitGlobalCache(cacheConfig)
}

// setupGlobalEnvironment sets up global environment variables and settings
func setupGlobalEnvironment(termInfo *terminal.TerminalInfo) {
	// Handle color settings
	if noColor || !termInfo.SupportsColor {
		os.Setenv("NO_COLOR", "1")
		viper.Set("ui.colors", false)
	}

	// Handle theme settings
	if theme != "" {
		viper.Set("ui.theme", theme)
	}

	// Set terminal size for responsive design
	terminal.RefreshTerminalSize()

	// Start terminal size monitoring
	terminal.AddSizeChangeCallback(func(width, height int) {
		utils.WithFields(map[string]interface{}{
			"width":  width,
			"height": height,
		}).Debug("Terminal size changed")
	})
}

// shouldShowWelcome determines if we should show the welcome message
func shouldShowWelcome() bool {
	// Don't show welcome if:
	// - Running in non-interactive mode
	// - Quiet flag is set
	// - NO_COLOR is set (might be in a script)
	// - Not a TTY

	if quiet || os.Getenv("NO_COLOR") != "" {
		return false
	}

	if !terminal.IsRunningInTerminal() {
		return false
	}

	// Check if we're running a specific command (not just help)
	args := os.Args[1:]
	if len(args) > 0 {
		// Don't show welcome for help commands
		if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
			return false
		}
		// Don't show welcome for version commands
		if args[0] == "version" || args[0] == "--version" || args[0] == "-v" {
			return false
		}
		// Don't show welcome for config commands
		if args[0] == "config" {
			return false
		}
	}

	return true
}

// displayWelcome shows the welcome message
func displayWelcome(termInfo *terminal.TerminalInfo) {
	// Only show welcome for interactive sessions and first-time usage
	width := termInfo.Width
	if width > 120 {
		width = 120
	}

	// Create a simple welcome banner
	welcome := ascii.WelcomeBanner(width)
	fmt.Println(welcome)

	// Show quick usage hint
	if termInfo.SupportsColor {
		hint := terminal.Cyan("💡 Tip: Use 'antoine --help' to see all available commands")
		fmt.Println(hint)
	} else {
		fmt.Println("Tip: Use 'antoine --help' to see all available commands")
	}

	fmt.Println() // Add some spacing
}

// setupRootCommand configures the root command with global flags
func init() {
	// Get the root command from cmd package
	rootCmd := cmd.GetRootCommand()

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "",
		"config file (default searches for antoine.yaml in current dir and $HOME/.antoine/)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false,
		"suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false,
		"disable colored output")
	rootCmd.PersistentFlags().StringVar(&theme, "theme", "",
		"UI theme (dark, light, minimal)")

	// Bind flags to viper
	viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("ui.colors", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("ui.theme", rootCmd.PersistentFlags().Lookup("theme"))

	// Set version information
	rootCmd.Version = fmt.Sprintf("%s (built %s, commit %s)", Version, BuildTime, GitCommit)

	// Add pre-run hook for validation
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Validate theme if specified
		if theme != "" {
			validThemes := []string{"dark", "light", "minimal"}
			valid := false
			for _, validTheme := range validThemes {
				if theme == validTheme {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid theme '%s'. Valid themes: %v", theme, validThemes)
			}
		}

		return nil
	}

	// Add completion command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for Antoine CLI.

The command for each shell will print completion script to stdout.
You can source it or write it to a file and source it from your shell profile.`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	})
}

// Emergency error handler for panics
func init() {
	// Set up panic recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Antoine CLI encountered a fatal error: %v\n", r)

			// Try to log the panic if logging is initialized
			if logger := utils.GetGlobalLogger(); logger != nil {
				logger.WithFields(map[string]interface{}{
					"panic":   r,
					"runtime": runtime.Version(),
					"goos":    runtime.GOOS,
					"goarch":  runtime.GOARCH,
					"version": Version,
				}).Error("Application panic")
			}

			os.Exit(1)
		}
	}()
}

// Helper functions for environment detection

// isRunningInCI detects if we're running in a CI environment
func isRunningInCI() bool {
	ciEnvs := []string{"CI", "CONTINUOUS_INTEGRATION", "GITHUB_ACTIONS", "GITLAB_CI", "CIRCLECI"}
	for _, env := range ciEnvs {
		if os.Getenv(env) != "" {
			return true
		}
	}
	return false
}

// isRunningInContainer detects if we're running in a container
func isRunningInContainer() bool {
	// Check for common container indicators
	if os.Getenv("DOCKER_CONTAINER") != "" {
		return true
	}

	// Check for /.dockerenv file
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

// Configuration validation helpers

// validateConfig validates the loaded configuration
func validateConfig() error {
	validator := utils.NewConfigValidation()

	// Validate logging configuration
	validator.ValidateLogLevel(viper.GetString("logging.level"))
	validator.ValidateLogFormat(viper.GetString("logging.format"))
	validator.ValidateLogOutput(viper.GetString("logging.output"))

	// Validate cache configuration
	validator.ValidateCacheType(viper.GetString("cache.type"))

	// Validate UI configuration
	validator.ValidateUITheme(viper.GetString("ui.theme"))

	if validator.HasErrors() {
		return fmt.Errorf("configuration validation failed: %v", validator.Errors())
	}

	return nil
}

// Runtime information for debugging

// getRuntimeInfo returns runtime information for debugging
func getRuntimeInfo() map[string]interface{} {
	termInfo := terminal.DetectTerminal()

	return map[string]interface{}{
		"version":    Version,
		"build_time": BuildTime,
		"git_commit": GitCommit,
		"go_version": runtime.Version(),
		"goos":       runtime.GOOS,
		"goarch":     runtime.GOARCH,
		"terminal": map[string]interface{}{
			"program":          termInfo.Program,
			"type":             termInfo.Type,
			"supports_color":   termInfo.SupportsColor,
			"supports_unicode": termInfo.SupportsUnicode,
			"width":            termInfo.Width,
			"height":           termInfo.Height,
		},
		"environment": map[string]interface{}{
			"ci":        isRunningInCI(),
			"container": isRunningInContainer(),
			"no_color":  os.Getenv("NO_COLOR") != "",
		},
	}
}

// Debug command for troubleshooting
func init() {
	debugCmd := &cobra.Command{
		Use:    "debug",
		Short:  "Show debug information",
		Hidden: true, // Hide from main help
		Run: func(cmd *cobra.Command, args []string) {
			info := getRuntimeInfo()

			if data, err := utils.PrettyPrintJSON(info); err == nil {
				fmt.Println(data)
			} else {
				fmt.Printf("Runtime Info: %+v\n", info)
			}

			// Show configuration
			fmt.Println("\nConfiguration:")
			fmt.Printf("Config file: %s\n", viper.ConfigFileUsed())

			// Show cache stats if available
			if cache := utils.GetGlobalCache(); cache != nil {
				stats := cache.Stats()
				fmt.Printf("Cache stats: %+v\n", stats)
			}
		},
	}

	cmd.GetRootCommand().AddCommand(debugCmd)
}
