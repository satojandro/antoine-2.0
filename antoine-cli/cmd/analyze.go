package cmd

import (
	"antoine-cli/internal/ui/views"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze repositories, trends, and technologies",
	Long: `Deep analysis powered by Antoine's AI. Get insights on code quality,
innovation potential, market trends, and competitive landscape.`,
}

var analyzeRepoCmd = &cobra.Command{
	Use:   "repo [repository-url]",
	Short: "Analyze a GitHub repository",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]

		options := &views.AnalysisOptions{
			RepoURL:             repoURL,
			Depth:               cmd.Flag("depth").Value.String(),
			IncludeDependencies: cmd.Flag("include-dependencies").Changed,
			GenerateReport:      cmd.Flag("generate-report").Changed,
			Focus:               cmd.Flag("focus").Value.String(),
			Format:              viper.GetString("format"),
		}

		view := views.NewAnalysisView(client)
		view.AnalyzeRepository(options)
	},
}

var analyzeTrendsCmd = &cobra.Command{
	Use:   "trends",
	Short: "Analyze technology and market trends",

	Run: func(cmd *cobra.Command, args []string) {
		options := &views.AnalysisOptions{
			Tech:      cmd.Flag("tech").Value.String(),
			Timeframe: cmd.Flag("timeframe").Value.String(),
			Metrics:   cmd.Flag("include-metrics").Changed,
			Format:    viper.GetString("format"),
		}

		view := views.NewAnalysisView(client)
		view.AnalyzeTrends(options)
	},
}

func init() {
	// Flags para análisis de repositorio
	analyzeRepoCmd.Flags().String("depth", "standard", "analysis depth (quick, standard, deep)")
	analyzeRepoCmd.Flags().Bool("include-dependencies", false, "analyze dependencies")
	analyzeRepoCmd.Flags().Bool("generate-report", false, "generate detailed report")
	analyzeRepoCmd.Flags().StringSlice("focus", []string{}, "focus areas (architecture, security, performance)")

	// Flags para análisis de tendencias
	analyzeTrendsCmd.Flags().StringSlice("tech", []string{}, "technologies to analyze")
	analyzeTrendsCmd.Flags().String("timeframe", "6months", "analysis timeframe (1month, 3months, 6months, 1year)")
	analyzeTrendsCmd.Flags().Bool("include-metrics", false, "include detailed metrics")
	analyzeTrendsCmd.Flags().String("market", "", "specific market (defi, gaming, ai, etc.)")

	analyzeCmd.AddCommand(analyzeRepoCmd)
	analyzeCmd.AddCommand(analyzeTrendsCmd)
}
