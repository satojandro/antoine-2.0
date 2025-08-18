package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"antoine-cli/internal/config"
	"antoine-cli/internal/core"
	"antoine-cli/pkg/ascii"
	"antoine-cli/pkg/terminal"
)

var (
	cfgFile string
	client  *core.AntoineClient
	version = "1.0.0"
	rootCmd *cobra.Command
)

// init inicializa el comando root
func init() {
	// Crear el comando root
	rootCmd = &cobra.Command{
		Use:   "antoine",
		Short: "Your Ultimate Hackathon Mentor",
		Long: getHeaderWithTerminalWidth() + `

Antoine is an AI-powered CLI tool that helps developers excel in hackathons.
From discovering the perfect hackathon to analyzing winning projects and 
providing personalized mentorship - Antoine is your ultimate competitive edge.

Powered by cutting-edge MCP servers and advanced AI analysis, Antoine learns
from every hackathon to help you build what's next.`,

		Run: func(cmd *cobra.Command, args []string) {
			// Si no hay argumentos, mostrar el dashboard interactivo
			if len(args) == 0 {
				showDashboard()
				return
			}

			// Mostrar ayuda por defecto
			cmd.Help()
		},
	}

	// Configurar inicialización
	cobra.OnInitialize(initConfig)

	// Flags globales
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.antoine.yaml)")
	rootCmd.PersistentFlags().Bool("no-color", false,
		"disable colored output")
	rootCmd.PersistentFlags().Bool("no-animation", false,
		"disable animations")
	rootCmd.PersistentFlags().String("format", "interactive",
		"output format (interactive, json, yaml, table)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false,
		"verbose output")
	rootCmd.PersistentFlags().Bool("debug", false,
		"debug mode")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false,
		"suppress non-essential output")
	rootCmd.PersistentFlags().String("theme", "",
		"UI theme (dark, light, minimal)")

	// Flags del comando root
	rootCmd.Flags().Bool("version", false, "show version")

	// Vincular flags con viper
	viper.BindPFlag("ui.colors", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("ui.animations", rootCmd.PersistentFlags().Lookup("no-animation"))
	viper.BindPFlag("output.format", rootCmd.PersistentFlags().Lookup("format"))
	viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("debug.enabled", rootCmd.PersistentFlags().Lookup("debug"))

	// Añadir subcomandos
	initSubcommands()
}

// GetRootCommand returns the root command for use in main.go
func GetRootCommand() *cobra.Command {
	return rootCmd
}

// Execute añade todos los comandos hijos al comando root y establece las flags apropiadamente.
func Execute() error {
	return rootCmd.Execute()
}

// initConfig lee el archivo de configuración usando el nuevo sistema
func initConfig() {
	// Si se especifica un archivo de configuración via flag, usarlo
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// Usar el nuevo sistema de configuración
	if err := config.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Aplicar overrides de flags de línea de comandos
	applyFlagOverrides()

	// Mostrar archivo de configuración en uso si está en modo verbose
	if viper.GetBool("debug.enabled") || viper.GetString("logging.level") == "debug" {
		if configFile := viper.ConfigFileUsed(); configFile != "" {
			fmt.Fprintf(os.Stderr, "Using config file: %s\n", configFile)
		}
	}

	// Inicializar cliente Antoine
	initializeClient()
}

// applyFlagOverrides aplica los overrides de flags de línea de comandos
func applyFlagOverrides() {
	// Manejar no-color flag (invertir la lógica)
	if noColor, _ := rootCmd.PersistentFlags().GetBool("no-color"); noColor {
		viper.Set("ui.colors", false)
	}

	// Manejar no-animation flag (invertir la lógica)
	if noAnimation, _ := rootCmd.PersistentFlags().GetBool("no-animation"); noAnimation {
		viper.Set("ui.animations", false)
	}

	// Manejar verbose flag
	if verbose, _ := rootCmd.PersistentFlags().GetBool("verbose"); verbose {
		viper.Set("logging.level", "debug")
	}

	// Manejar debug flag
	if debug, _ := rootCmd.PersistentFlags().GetBool("debug"); debug {
		viper.Set("debug.enabled", true)
		viper.Set("logging.level", "debug")
		viper.Set("debug.verbose_logging", true)
	}
}

// initializeClient inicializa el cliente Antoine
func initializeClient() {
	cfg := config.Get()

	var err error
	client, err = core.NewAntoineClient(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Antoine client: %v\n", err)
		os.Exit(1)
	}
}

// initSubcommands inicializa todos los subcomandos
func initSubcommands() {
	// Añadir subcomandos principales
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(mentorCmd)
	rootCmd.AddCommand(trendsCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(getVersionCommand())

	// Comando de completion
	rootCmd.AddCommand(getCompletionCommand())
}

// getHeaderWithTerminalWidth obtiene el header ajustado al terminal
func getHeaderWithTerminalWidth() string {
	width := terminal.GetTerminalWidth()
	if width <= 0 {
		width = 80 // fallback seguro
	}
	return ascii.GetHeader(width, ascii.SubtitleMain)
}

// showDashboard muestra el dashboard interactivo principal
func showDashboard() {
	cfg := config.Get()
	termWidth := terminal.GetTerminalWidth()

	// Aplicar configuración de colores
	if !cfg.UI.Colors {
		lipgloss.SetColorProfile(termenv.Ascii)
	} else {
		// Detectar el perfil de color apropiado
		lipgloss.SetColorProfile(termenv.ColorProfile())
	}

	// Mostrar header
	fmt.Println(ascii.GetHeader(termWidth, ascii.SubtitleMain))
	fmt.Println()

	// Cargar estadísticas (mock por ahora)
	stats := getDashboardStats()

	// Mostrar dashboard
	fmt.Println(ascii.GetDashboardStats(stats))
	fmt.Println()

	// Mostrar comandos disponibles
	showQuickCommands(cfg)
}

// getDashboardStats obtiene las estadísticas del dashboard
func getDashboardStats() map[string]interface{} {
	// TODO: Implementar estadísticas reales desde el cliente
	return map[string]interface{}{
		"hackathons":   42,
		"projects":     1337,
		"searches":     156,
		"success_rate": 94,
		"last_update":  time.Now().Format("2006-01-02 15:04"),
	}
}

// showQuickCommands muestra una guía rápida de comandos
func showQuickCommands(cfg *config.Config) {
	var style lipgloss.Style

	if cfg.UI.Colors {
		style = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ascii.Cyan).
			Padding(1, 2).
			Margin(1, 0)
	} else {
		style = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Margin(1, 0)
	}

	content := `🚀 Quick Commands:

• antoine search hackathons --tech "AI"     - Find AI hackathons
• antoine analyze repo <url>                - Analyze a GitHub repository  
• antoine mentor start                      - Start interactive mentorship
• antoine trends --tech "blockchain"        - See blockchain trends
• antoine config show                       - View current configuration

Type 'antoine --help' for complete command reference.`

	fmt.Println(style.Render(content))
}

// getVersionCommand crea el comando de versión
func getVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show Antoine version",
		Long:  "Display version information for Antoine CLI",
		Run: func(cmd *cobra.Command, args []string) {
			showVersionInfo()
		},
	}
}

// showVersionInfo muestra la información de versión
func showVersionInfo() {
	cfg := config.Get()
	termWidth := terminal.GetTerminalWidth()

	var versionStyle lipgloss.Style
	if cfg.UI.Colors {
		versionStyle = lipgloss.NewStyle().
			Foreground(ascii.Gold).
			Bold(true)
	} else {
		versionStyle = lipgloss.NewStyle().
			Bold(true)
	}

	// Mostrar logo si ASCII art está habilitado
	if cfg.UI.ASCIIArt {
		fmt.Println(ascii.GetLogo(termWidth))
		fmt.Println()
	}

	// Información de versión
	fmt.Printf("%s v%s\n", versionStyle.Render("Antoine CLI"), version)
	fmt.Printf("Built with ❤️  for hackathon champions\n")

	// Información técnica
	fmt.Printf("\nBuild Information:\n")
	fmt.Printf("  Go version: %s\n", runtime.Version())
	fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  Build time: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// Información de configuración
	if cfg.Debug.Enabled {
		fmt.Printf("\nConfiguration:\n")
		fmt.Printf("  Config file: %s\n", viper.ConfigFileUsed())
		fmt.Printf("  Theme: %s\n", cfg.UI.Theme)
		fmt.Printf("  Cache enabled: %v\n", cfg.Cache.Enabled)
	}
}

// getCompletionCommand crea el comando de completion
func getCompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for Antoine CLI.

The command for each shell will print completion script to stdout.
You can source it or write it to a file and source it from your shell profile.

Examples:
  # Bash
  source <(antoine completion bash)
  
  # Zsh
  source <(antoine completion zsh)
  
  # Fish
  antoine completion fish | source
  
  # PowerShell
  antoine completion powershell | Out-String | Invoke-Expression`,

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
	}
}

// Helper function para validar comandos
func validateCommand(cmd *cobra.Command, args []string) error {
	cfg := config.Get()

	// Validar que las features necesarias estén habilitadas
	switch cmd.Name() {
	case "search":
		if !cfg.Features.SearchEnabled {
			return fmt.Errorf("search feature is disabled in configuration")
		}
	case "analyze":
		if !cfg.Features.AnalysisEnabled {
			return fmt.Errorf("analysis feature is disabled in configuration")
		}
	case "mentor":
		if !cfg.Features.MentorEnabled {
			return fmt.Errorf("mentor feature is disabled in configuration")
		}
	}

	return nil
}

// GetClient returns the initialized Antoine client
func GetClient() *core.AntoineClient {
	return client
}

// GetConfig returns the current configuration
func GetConfig() *config.Config {
	return config.Get()
}
