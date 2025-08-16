package cmd

import (
	"antoine-cli/internal/ui/views"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for hackathons, projects, and trends",
	Long: `Search through Antoine's vast database of hackathons, winning projects,
and developer trends. Powered by advanced MCP servers for real-time data.`,
}

var searchHackathonsCmd = &cobra.Command{
	Use:   "hackathons",
	Short: "Search for hackathons",
	Long: `Find the perfect hackathon based on your interests, skills, and timeline.
Antoine analyzes thousands of hackathons to find your ideal match.`,

	Run: func(cmd *cobra.Command, args []string) {
		options := &views.SearchOptions{
			Tech:     cmd.Flag("tech").Value.String(),
			Location: cmd.Flag("location").Value.String(),
			PrizeMin: cmd.Flag("prize-min").Value.String(),
			DateFrom: cmd.Flag("date-from").Value.String(),
			DateTo:   cmd.Flag("date-to").Value.String(),
			Online:   cmd.Flag("online").Changed,
			Format:   viper.GetString("format"),
		}

		view := views.NewSearchView(client)
		view.SearchHackathons(options)
	},
}

var searchProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Search for hackathon projects",
	Long: `Discover winning projects from past hackathons. Learn from the best
and get inspiration for your next project.`,

	Run: func(cmd *cobra.Command, args []string) {
		options := &views.SearchOptions{
			Hackathon: cmd.Flag("hackathon").Value.String(),
			Category:  cmd.Flag("category").Value.String(),
			Tech:      cmd.Flag("tech").Value.String(),
			Sort:      cmd.Flag("sort").Value.String(),
			Format:    viper.GetString("format"),
		}

		view := views.NewSearchView(client)
		view.SearchProjects(options)
	},
}

func init() {
	// Flags para hackathons
	searchHackathonsCmd.Flags().StringSlice("tech", []string{}, "technologies (e.g., --tech AI,Blockchain)")
	searchHackathonsCmd.Flags().String("location", "", "location (city, country, or 'online')")
	searchHackathonsCmd.Flags().Int("prize-min", 0, "minimum prize amount")
	searchHackathonsCmd.Flags().String("date-from", "", "start date (YYYY-MM-DD)")
	searchHackathonsCmd.Flags().String("date-to", "", "end date (YYYY-MM-DD)")
	searchHackathonsCmd.Flags().Bool("online", false, "online hackathons only")
	searchHackathonsCmd.Flags().String("difficulty", "", "difficulty level (beginner, intermediate, advanced)")

	// Flags para proyectos
	searchProjectsCmd.Flags().String("hackathon", "", "specific hackathon name")
	searchProjectsCmd.Flags().StringSlice("category", []string{}, "project categories")
	searchProjectsCmd.Flags().StringSlice("tech", []string{}, "technologies used")
	searchProjectsCmd.Flags().String("sort", "relevance", "sort by (relevance, popularity, recent, prize)")
	searchProjectsCmd.Flags().Int("limit", 10, "number of results")

	searchCmd.AddCommand(searchHackathonsCmd)
	searchCmd.AddCommand(searchProjectsCmd)
}
