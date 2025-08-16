package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"

	"antoine-cli/internal/config"
	"antoine-cli/internal/core"
	"antoine-cli/pkg/ascii"
)

var (
	cfgFile string
	client  *core.AntoineClient
	version = "1.0.0"
)

// rootCmd representa el comando base cuando se llama sin subcomandos
var rootCmd = &cobra.Command{
	Use:   "antoine",
	Short: "Your Ultimate Hackathon Mentor",
	Long: ascii.GetHeader(getTerminalWidth(), ascii.SubtitleMain) + `

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

// Execute añade todos los comandos hijos al comando root y establece las flags apropiadamente.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Flags globales
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.antoine.yaml)")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")
	rootCmd.PersistentFlags().Bool("no-animation", false, "disable animations")
	rootCmd.PersistentFlags().String("format", "interactive", "output format (interactive, json, yaml, table)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	rootCmd.PersistentFlags().Bool("debug", false, "debug mode")

	// Flags del comando root
	rootCmd.Flags().BoolP("version", "v", false, "show version")

	// Vincular flags con viper
	viper.BindPFlag("no_color", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("no_animation", rootCmd.PersistentFlags().Lookup("no-animation"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Añadir subcomandos
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(mentorCmd)
	rootCmd.AddCommand(trendsCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
}

// initConfig lee el archivo de configuración y variables de entorno si están establecidas.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".antoine")
	}

	viper.SetEnvPrefix("ANTOINE")
	viper.AutomaticEnv()

	// Cargar configuración por defecto
	config.SetDefaults()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	// Inicializar cliente
	var err error
	client, err = core.NewAntoineClient(config.Get())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Antoine client: %v\n", err)
		os.Exit(1)
	}
}

// showDashboard muestra el dashboard interactivo principal
func showDashboard() {
	termWidth := getTerminalWidth()

	// Mostrar header
	fmt.Println(ascii.GetHeader(termWidth, ascii.SubtitleMain))
	fmt.Println()

	// Cargar estadísticas (mock por ahora)
	stats := map[string]interface{}{
		"hackathons":   42,
		"projects":     1337,
		"searches":     156,
		"success_rate": 94,
	}

	// Mostrar dashboard
	fmt.Println(ascii.GetDashboardStats(stats))
	fmt.Println()

	// Mostrar comandos disponibles
	showQuickCommands()
}

// showQuickCommands muestra una guía rápida de comandos
func showQuickCommands() {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ascii.Cyan).
		Padding(1, 2).
		Margin(1, 0)

	content := `🚀 Quick Commands:

• antoine search hackathons --tech "AI"     - Find AI hackathons
• antoine analyze repo <url>                - Analyze a GitHub repository  
• antoine mentor start                      - Start interactive mentorship
• antoine trends --tech "blockchain"        - See blockchain trends
• antoine config show                       - View current configuration

Type 'antoine --help' for complete command reference.`

	fmt.Println(style.Render(content))
}

// getTerminalWidth obtiene el ancho del terminal
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // fallback
	}
	return width
}

// versionCmd maneja el comando de versión
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Antoine version",
	Run: func(cmd *cobra.Command, args []string) {
		termWidth := getTerminalWidth()

		versionStyle := lipgloss.NewStyle().
			Foreground(ascii.Gold).
			Bold(true)

		fmt.Println(ascii.GetLogo(termWidth))
		fmt.Println()
		fmt.Printf("%s v%s\n", versionStyle.Render("Antoine CLI"), version)
		fmt.Printf("Built with ❤️  for hackathon champions\n")
		fmt.Printf("Build time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	},
}
