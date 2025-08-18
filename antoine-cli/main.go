// Antoine CLI - Your Ultimate Hackathon Mentor
// This is the main entry point for the Antoine CLI application
package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

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

	// Command line flags
	//configFile string
	//verbose    bool
	//quiet      bool
	//noColor    bool
	//theme      string
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
	//termInfo := terminal.DetectTerminal()
	termInfo := terminal.SafeDetectTerminal()

	// Initialize configuration using the new config system
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

// initConfig initializes the configuration system using the new consolidated config
func initConfig() error {
	// Obtener el archivo de configuración desde la flag ya definida en cmd/root.go
	rootCmd := cmd.GetRootCommand()
	if configFlag := rootCmd.PersistentFlags().Lookup("config"); configFlag != nil {
		if configFile := configFlag.Value.String(); configFile != "" {
			viper.SetConfigFile(configFile)
		}
	}

	// Use the new config system's LoadConfig function
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Apply command line flag overrides
	applyCommandLineOverrides()

	// Validate the final configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

// applyCommandLineOverrides applies command line flags to override config values
func applyCommandLineOverrides() {
	// Acceder a las flags ya definidas en cmd/root.go
	rootCmd := cmd.GetRootCommand()

	// Obtener valores de flags usando los métodos de Cobra
	if verboseFlag := rootCmd.PersistentFlags().Lookup("verbose"); verboseFlag != nil && verboseFlag.Changed {
		viper.Set("logging.level", "debug")
		viper.Set("debug.enabled", true)
		viper.Set("debug.verbose_logging", true)
	}

	if quietFlag := rootCmd.PersistentFlags().Lookup("quiet"); quietFlag != nil && quietFlag.Changed {
		viper.Set("logging.level", "error")
	}

	if noColorFlag := rootCmd.PersistentFlags().Lookup("no-color"); noColorFlag != nil && noColorFlag.Changed {
		viper.Set("ui.colors", false)
	}

	if themeFlag := rootCmd.PersistentFlags().Lookup("theme"); themeFlag != nil && themeFlag.Changed {
		viper.Set("ui.theme", themeFlag.Value.String())
	}
}

// initLogging initializes the logging system
func initLogging() error {
	// Get configuration values
	cfg := config.Get()

	logConfig := utils.LoggerConfig{
		Level:       utils.LogLevel(cfg.Logging.Level),
		Format:      utils.LogFormat(cfg.Logging.Format),
		Output:      utils.LogOutput(cfg.Logging.Output),
		Development: cfg.Debug.Enabled,
		Caller:      cfg.Logging.Caller,
		StackTrace:  cfg.Logging.StackTrace,
	}

	// Set file logging configuration
	logConfig.File.Path = cfg.Logging.File.Path
	logConfig.File.MaxSizeMB = cfg.Logging.File.MaxSizeMB
	logConfig.File.MaxBackups = cfg.Logging.File.MaxBackups
	logConfig.File.MaxAgeDays = cfg.Logging.File.MaxAgeDays
	logConfig.File.Compress = cfg.Logging.File.Compress

	return utils.InitGlobalLogger(logConfig)
}

// initCache initializes the cache system
func initCache() error {
	cfg := config.Get()

	if !cfg.Cache.Enabled {
		// Initialize with disabled cache
		return utils.InitGlobalCache(utils.CacheConfig{Enabled: false})
	}

	cacheConfig := utils.CacheConfig{
		Enabled:         true,
		Type:            utils.CacheType(cfg.Cache.Type),
		MaxSizeMB:       int64(cfg.Cache.MaxSizeMB),
		MaxEntries:      int64(cfg.Cache.MaxSize),
		CleanupInterval: parseDurationSafely(cfg.Cache.CleanupInterval),
		TTL:             make(map[string]time.Duration),
	}

	// Set TTL values - convert from string to time.Duration
	if cfg.Cache.TTLByType != nil {
		for key, ttlStr := range cfg.Cache.TTLByType {
			if duration := parseDurationSafely(ttlStr); duration > 0 {
				cacheConfig.TTL[key] = duration
			}
		}
	}

	// Disk cache settings
	cacheConfig.Disk.Path = cfg.Cache.Disk.Path
	cacheConfig.Disk.Compression = cfg.Cache.Disk.Compression

	return utils.InitGlobalCache(cacheConfig)
}

// setupGlobalEnvironment sets up global environment variables and settings
func setupGlobalEnvironment(termInfo *terminal.TerminalInfo) {
	cfg := config.Get()

	// Handle color settings
	if !cfg.UI.Colors || !termInfo.SupportsColor {
		os.Setenv("NO_COLOR", "1")
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
	// Obtener flags del comando root
	rootCmd := cmd.GetRootCommand()

	// Verificar flag quiet
	quietFlag := rootCmd.PersistentFlags().Lookup("quiet")
	isQuiet := quietFlag != nil && quietFlag.Changed

	// Don't show welcome if:
	// - Running in non-interactive mode
	// - Quiet flag is set
	// - NO_COLOR is set (might be in a script)
	// - Not a TTY

	if isQuiet || os.Getenv("NO_COLOR") != "" {
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
	cfg := config.Get()

	// Only show welcome for interactive sessions and first-time usage
	width := termInfo.Width
	if width > cfg.UI.MaxWidth {
		width = cfg.UI.MaxWidth
	}

	// Create a simple welcome banner
	welcome := ascii.WelcomeBanner(width)
	fmt.Println(welcome)

	// Show quick usage hint
	if cfg.UI.Colors && termInfo.SupportsColor {
		colorize := terminal.GetGlobalColorize()
		hint := colorize.Cyan("💡 Tip: Use 'antoine --help' to see all available commands")
		fmt.Println(hint)
	} else {
		fmt.Println("💡 Tip: Use 'antoine --help' to see all available commands")
	}

	fmt.Println() // Add some spacing
}

//// setupRootCommand configures the root command with global flags
//func init() {
//	// Get the root command from cmd package
//	rootCmd := cmd.GetRootCommand()
//
//	// Add global flags
//	rootCmd.PersistentFlags().StringVar(&configFile, "config", "",
//		"config file (default searches for antoine.yaml in current dir and $HOME/.antoine/)")
//	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
//		"enable verbose output")
//	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false,
//		"suppress non-essential output")
//	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false,
//		"disable colored output")
//	rootCmd.PersistentFlags().StringVar(&theme, "theme", "",
//		"UI theme (dark, light, minimal)")
//
//	// Set version information
//	rootCmd.Version = fmt.Sprintf("%s (built %s, commit %s)", Version, BuildTime, GitCommit)
//
//	// Add pre-run hook for validation
//	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
//		// Validate theme if specified
//		if theme != "" {
//			validThemes := []string{"dark", "light", "minimal"}
//			valid := false
//			for _, validTheme := range validThemes {
//				if theme == validTheme {
//					valid = true
//					break
//				}
//			}
//			if !valid {
//				return fmt.Errorf("invalid theme '%s'. Valid themes: %v", theme, validThemes)
//			}
//		}
//
//		return nil
//	}
//
//	// Add completion command
//	rootCmd.AddCommand(&cobra.Command{
//		Use:   "completion [bash|zsh|fish|powershell]",
//		Short: "Generate shell completion scripts",
//		Long: `Generate shell completion scripts for Antoine CLI.
//
//The command for each shell will print completion script to stdout.
//You can source it or write it to a file and source it from your shell profile.`,
//		DisableFlagsInUseLine: true,
//		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
//		Args:                  cobra.ExactValidArgs(1),
//		Run: func(cmd *cobra.Command, args []string) {
//			switch args[0] {
//			case "bash":
//				cmd.Root().GenBashCompletion(os.Stdout)
//			case "zsh":
//				cmd.Root().GenZshCompletion(os.Stdout)
//			case "fish":
//				cmd.Root().GenFishCompletion(os.Stdout, true)
//			case "powershell":
//				cmd.Root().GenPowerShellCompletion(os.Stdout)
//			}
//		},
//	})
//}

//// Emergency error handler for panics
//func init() {
//	// Set up panic recovery
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Fprintf(os.Stderr, "Antoine CLI encountered a fatal error: %v\n", r)
//
//			// Try to log the panic if logging is initialized
//			if logger := utils.GetGlobalLogger(); logger != nil {
//				logger.WithFields(map[string]interface{}{
//					"panic":   r,
//					"runtime": runtime.Version(),
//					"goos":    runtime.GOOS,
//					"goarch":  runtime.GOARCH,
//					"version": Version,
//				}).Error("Application panic")
//			}
//
//			os.Exit(1)
//		}
//	}()
//}

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

// getRuntimeInfo returns runtime information for debugging
func getRuntimeInfo() map[string]interface{} {
	termInfo := terminal.DetectTerminal()
	cfg := config.Get()

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
		"config": map[string]interface{}{
			"file":  viper.ConfigFileUsed(),
			"theme": cfg.UI.Theme,
			"cache": cfg.Cache.Enabled,
		},
	}
}

//// Debug command for troubleshooting
//func init() {
//	debugCmd := &cobra.Command{
//		Use:    "debug",
//		Short:  "Show debug information",
//		Hidden: true, // Hide from main help
//		Run: func(cmd *cobra.Command, args []string) {
//			info := getRuntimeInfo()
//
//			if data, err := utils.PrettyPrintJSON(info); err == nil {
//				fmt.Println(data)
//			} else {
//				fmt.Printf("Runtime Info: %+v\n", info)
//			}
//
//			// Show configuration using the new config system
//			fmt.Println("\nConfiguration:")
//			config.Show() // Use the enhanced Show() method
//
//			// Show cache stats if available
//			if cache := utils.GetGlobalCache(); cache != nil {
//				stats := cache.Stats()
//				fmt.Printf("\nCache stats: %+v\n", stats)
//			}
//		},
//	}
//
//	cmd.GetRootCommand().AddCommand(debugCmd)
//}

// Helper function to safely parse duration strings
func parseDurationSafely(durationStr string) time.Duration {
	if duration, err := time.ParseDuration(durationStr); err == nil {
		return duration
	}
	return 0
}
